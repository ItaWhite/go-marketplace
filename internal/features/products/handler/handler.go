package handler

import (
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/transport/middleware"
	"go-marketplace/internal/features/products/service"
	"net/http"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{
		service: s,
	}
}

func (h *ProductHandler) RegisterRoutes(mux *http.ServeMux, auth *middleware.AuthMiddleware) {
	mux.HandleFunc("GET /products", h.GetProducts)
	mux.HandleFunc("GET /products/{id}", h.GetProduct)
	mux.Handle("POST /products", auth.RequireRole(domain.UserSeller, http.HandlerFunc(h.PostProduct)))
	mux.Handle("PATCH /products/{id}", auth.RequireRole(domain.UserSeller, http.HandlerFunc(h.PatchProduct)))
	mux.Handle("DELETE /products/{id}", auth.RequireRole(domain.UserSeller, http.HandlerFunc(h.DeleteProduct)))
}
