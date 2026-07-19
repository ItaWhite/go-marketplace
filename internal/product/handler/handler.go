package handler

import (
	productfeat "go-marketplace/internal/product"
)

type ProductHandler struct {
	service *productfeat.ProductService
}

func NewProductHandler(s *productfeat.ProductService) *ProductHandler {
	return &ProductHandler{
		service: s,
	}
}
