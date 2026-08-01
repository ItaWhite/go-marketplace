package main

import (
	"context"
	"crypto/tls"
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
	"log"
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
		log.Fatal("Error loading .env file")
	}

	newLogger := core_logger.NewLogger(os.Getenv("LOGGER_LEVEL"), os.Getenv("LOGGER_FORMAT"))
	slog.SetDefault(newLogger)

	db, err := storage.NewPostgres(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_DB"))
	if err != nil {
		log.Fatal(err)
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

	s := http.Server{
		Addr:         addr,
		Handler:      chain(mux),
		TLSConfig:    tlsConfig,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	cert := "cmd/api/cert.pem"
	key := "cmd/api/key.pem"

	go func() {
		fmt.Printf("Server started at port %s...\n", os.Getenv("SERVER_PORT"))
		log.Fatal(s.ListenAndServeTLS(cert, key))
	}()

	<-ctx.Done()

	ctx, stop := context.WithTimeout(context.Background(), time.Second*5)
	defer stop()

	log.Fatal(s.Shutdown(ctx))
}
