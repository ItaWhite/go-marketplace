package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"go-marketplace/internal/product/handler"
	"go-marketplace/internal/product/repository"
	"go-marketplace/internal/product/service"
	"go-marketplace/internal/storage"
	"go-marketplace/internal/transport"
	"log"
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

	db, err := storage.NewPostgres(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_DB"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	addr := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))

	mux := transport.Router(productHandler)
	chain := transport.Chain(
		transport.Logging,
		transport.SecurityHeaders,
	)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

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
