package handler

import (
	productfeat "go-marketplace/internal/product/service"
)

type ProductHandler struct {
	service *productfeat.ProductService
}

func NewProductHandler(s *productfeat.ProductService) *ProductHandler {
	return &ProductHandler{
		service: s,
	}
}
