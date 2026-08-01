package handler

import (
	productfeat "go-marketplace/internal/features/products/service"
	"net/http"
)

type ProductHandler struct {
	service *productfeat.ProductService
}

func NewProductHandler(s *productfeat.ProductService) *ProductHandler {
	return &ProductHandler{
		service: s,
	}
}

func (h *ProductHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /products", h.GetProducts)
	mux.HandleFunc("GET /products/{id}", h.GetProduct)
	mux.HandleFunc("POST /products", h.PostProduct)
	mux.HandleFunc("PATCH /products/{id}", h.PatchProduct)
	mux.HandleFunc("DELETE /products/{id}", h.DeleteProduct)
}
