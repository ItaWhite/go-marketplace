package service

import (
	"context"
	"go-marketplace/internal/core/domain"
	productfeat "go-marketplace/internal/product"
	"strings"
)

func (s *ProductService) CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	if strings.TrimSpace(product.Name) == "" {
		return domain.Product{}, productfeat.ErrInvalidName
	}
	if product.Price < 0 {
		return domain.Product{}, productfeat.ErrInvalidPrice
	}
	return s.repo.Create(ctx, product)
}
