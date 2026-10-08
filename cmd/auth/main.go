package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"go-marketplace/internal/core/logger"
	"go-marketplace/internal/core/security"
	"go-marketplace/internal/core/storage"
	"go-marketplace/internal/core/transport/middleware"
	auth_grpc_client "go-marketplace/internal/features/auth/grpc"
	auth_handler "go-marketplace/internal/features/auth/handler"
	auth_repository "go-marketplace/internal/features/auth/repository"
	auth_service "go-marketplace/internal/features/auth/service"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		slog.Error("load .env", "error", err)
		os.Exit(1)
	}

	newLogger := core_logger.NewLogger(os.Getenv("LOGGER_LEVEL"), os.Getenv("LOGGER_FORMAT"))
	slog.SetDefault(newLogger)

	db, err := storage.NewPostgres(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_DB"))
	if err != nil {
		slog.Error("new postgres", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	authRepository := auth_repository.NewAuthRepository(db)

	publicKey, err := security.LoadPublicKey(os.Getenv("AUTH_PUBLIC_KEY_PATH"))
	if err != nil {
		slog.Error("load public key", "error", err)
		os.Exit(1)
	}

	privateKey, err := security.LoadPrivateKey(os.Getenv("AUTH_PRIVATE_KEY_PATH"))
	if err != nil {
		slog.Error("load private key", "error", err)
		os.Exit(1)
	}

	accessTTL, err := time.ParseDuration(os.Getenv("AUTH_ACCESS_TTL"))
	if err != nil {
		slog.Error("parse access TTL", "error", err)
		os.Exit(1)
	}

	refreshTTL, err := time.ParseDuration(os.Getenv("AUTH_REFRESH_TTL"))
	if err != nil {
		slog.Error("parse refresh TTL", "error", err)
		os.Exit(1)
	}

	usersClient, err := auth_grpc_client.NewUserServiceClient(os.Getenv("GRPC_ADDR"))
	if err != nil {
		slog.Error("new user service client", "error", err)
		os.Exit(1)
	}

	tokenService := auth_service.NewTokenService(os.Getenv("AUTH_ISSUER"), privateKey, accessTTL, refreshTTL)

	authService := auth_service.NewAuthService(
		authRepository,
		authRepository,
		usersClient,
		tokenService,
	)

	authHandler := auth_handler.NewAuthHandler(authService)

	mux := http.NewServeMux()

	validator := security.NewJWTValidator(os.Getenv("AUTH_ISSUER"), publicKey)
	authMiddleware := middleware.NewAuthMiddleware(validator)

	authHandler.RegisterRoutes(mux, authMiddleware)

	chain := middleware.Chain(
		middleware.RequestID,
		middleware.Logger,
		middleware.Panic,
		middleware.SecurityHeaders,
	)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	readTimeout, err := time.ParseDuration(os.Getenv("SERVER_READ_TIMEOUT"))
	if err != nil {
		slog.Error("parse read timeout", "error", err)
		os.Exit(1)
	}

	writeTimeout, err := time.ParseDuration(os.Getenv("SERVER_WRITE_TIMEOUT"))
	if err != nil {
		slog.Error("parse write timeout", "error", err)
		os.Exit(1)
	}

	idleTimeout, err := time.ParseDuration(os.Getenv("SERVER_IDLE_TIMEOUT"))
	if err != nil {
		slog.Error("parse idle timeout", "error", err)
		os.Exit(1)
	}

	s := http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("AUTH_SERVER_PORT")),
		Handler:      chain(mux),
		TLSConfig:    tlsConfig,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	cert := os.Getenv("TLS_CERT_PATH")
	key := os.Getenv("TLS_KEY_PATH")
	if cert == "" || key == "" {
		slog.Error("empty certificate or key path")
		os.Exit(1)
	}

	slog.Info("auth http server started", "port", os.Getenv("AUTH_SERVER_PORT"))

	go func() {
		err := s.ListenAndServeTLS(cert, key)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("auth http server failed", "error", err)
		}

		cancel()
	}()

	<-ctx.Done()

	ctx, stop := context.WithTimeout(context.Background(), 5*time.Second)
	defer stop()

	slog.Info("shutting down auth http server")

	start := time.Now()

	err = s.Shutdown(ctx)
	if err != nil {
		slog.Error("auth http server shutdown error", "error", err)
	} else {
		slog.Info("auth http server stopped", "duration", time.Since(start))
	}
}
