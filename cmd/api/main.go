package main

import (
	"context"
	"fmt"
	"go-marketplace/internal"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("cmd/api/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := internal.ConnectDb(os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal("Error connecting to database")
	}
	defer db.Close()

	err = db.Ping(context.Background())
	if err != nil {
		log.Fatal("Error pinging the database")
	}
	repository := internal.NewProductRepository(db)
	service := internal.NewProductService(repository)
	handler := internal.NewProductHandler(service)

	err = repository.DropSchema()
	if err != nil {
		log.Fatal(err)
	}
	err = repository.InitSchema()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /products", handler.GetProductsHandler)
	mux.HandleFunc("GET /products/{id}", handler.GetProductByIdHandler)
	mux.HandleFunc("POST /products", handler.PostProductsHandler)
	mux.HandleFunc("DELETE /products/{id}", handler.DeleteProductHandler)

	middlewareMux := internal.SecurityHeaders(mux)
	addr := fmt.Sprintf(":%s", os.Getenv("API_PORT"))
	s := http.Server{
		Addr:    addr,
		Handler: middlewareMux,
	}

	fmt.Printf("Server started at port %s...\n", os.Getenv("API_PORT"))
	log.Fatal(s.ListenAndServe())
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)
	})
}
