package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"go-marketplace/internal/product"
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

	err := godotenv.Load("cmd/api/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := storage.NewPostgres(os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repository := product.NewProductRepository(db)
	service := product.NewProductService(repository)
	handler := product.NewProductHandler(service)

	err = repository.DropSchema()
	if err != nil {
		log.Fatal(err)
	}
	err = repository.InitSchema()
	if err != nil {
		log.Fatal(err)
	}

	mux := transport.Router(handler)

	cert := "cmd/api/cert.pem"
	key := "cmd/api/key.pem"

	addr := fmt.Sprintf(":%s", os.Getenv("API_PORT"))
	middlewareMux := transport.SecurityHeaders(mux)
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	s := http.Server{
		Addr:         addr,
		Handler:      middlewareMux,
		TLSConfig:    tlsConfig,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		fmt.Printf("Server started at port %s...\n", os.Getenv("API_PORT"))
		log.Fatal(s.ListenAndServeTLS(cert, key))
	}()

	<-ctx.Done()

	ctx, stop := context.WithTimeout(context.Background(), time.Second*5)
	defer stop()

	log.Fatal(s.Shutdown(ctx))
}
