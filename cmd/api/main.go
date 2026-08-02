package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"go-marketplace/internal/core/logger"
	"go-marketplace/internal/core/storage"
	"go-marketplace/internal/core/transport/middleware"
	products_handler "go-marketplace/internal/features/products/handler"
	products_repository "go-marketplace/internal/features/products/repository"
	products_service "go-marketplace/internal/features/products/service"
	users_handler "go-marketplace/internal/features/users/handler"
	users_repository "go-marketplace/internal/features/users/repository"
	users_service "go-marketplace/internal/features/users/service"
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
		newLogger.Error("new postgres", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	productRepository := products_repository.NewProductRepository(db)
	productService := products_service.NewProductService(productRepository)
	productHandler := products_handler.NewProductHandler(productService)

	userRepository := users_repository.NewUserRepository(db)
	userService := users_service.NewUserService(userRepository)
	userHandler := users_handler.NewUserHandler(userService)

	mux := http.NewServeMux()

	productHandler.RegisterRoutes(mux)
	userHandler.RegisterRoutes(mux)

	chain := middleware.Chain(
		middleware.RequestID,
		middleware.Logger,
		middleware.Panic,
		middleware.SecurityHeaders,
	)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	addr := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))

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
		Addr:         addr,
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

	slog.Info("server started", "port", os.Getenv("SERVER_PORT"))

	go func() {
		err := s.ListenAndServeTLS(cert, key)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "error", err)
		}
		cancel()
	}()

	<-ctx.Done()

	ctx, stop := context.WithTimeout(context.Background(), 5*time.Second)
	defer stop()

	slog.Info("shutting down server")

	start := time.Now()

	err = s.Shutdown(ctx)
	if err != nil {
		slog.Error("server shutdown error", "error", err)
	} else {
		slog.Info("server stopped", "duration", time.Since(start))
	}
}
