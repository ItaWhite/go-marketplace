package internal

import "net/http"

func Router(h *productHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /products", h.GetProductsHandler)
	mux.HandleFunc("GET /products/{id}", h.GetProductByIdHandler)
	mux.HandleFunc("POST /products", h.PostProductsHandler)
	mux.HandleFunc("DELETE /products/{id}", h.DeleteProductHandler)

	return mux
}
