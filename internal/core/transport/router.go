package transport

import (
	"go-marketplace/internal/product/handler"
	"net/http"
)

func Router(h *handler.ProductHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /products", h.GetProducts)
	mux.HandleFunc("GET /products/{id}", h.GetProduct)
	mux.HandleFunc("POST /products", h.PostProduct)
	mux.HandleFunc("PUT /products/{id}", h.PutProduct)
	mux.HandleFunc("DELETE /products/{id}", h.DeleteProduct)

	return mux
}
