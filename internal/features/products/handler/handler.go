package handler

import (
	productfeat "go-marketplace/internal/features/products/service"
)

type ProductHandler struct {
	service *productfeat.ProductService
}

func NewProductHandler(s *productfeat.ProductService) *ProductHandler {
	return &ProductHandler{
		service: s,
	}
}
