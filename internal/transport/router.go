package transport

import (
	"go-marketplace/internal/product/handler"
	"net/http"
)

func Router(h *handler.ProductHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /products", h.GetProductsHandler)
	mux.HandleFunc("GET /products/{id}", h.GetProductByIdHandler)
	mux.HandleFunc("POST /products", h.PostProductsHandler)
	mux.HandleFunc("PUT /products/{id}", h.PutProductsHandler)
	mux.HandleFunc("DELETE /products/{id}", h.DeleteProductHandler)

	return mux
}
