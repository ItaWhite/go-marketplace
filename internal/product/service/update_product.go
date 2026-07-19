package service

import (
	"context"
	"go-marketplace/internal/core/domain"
	productfeat "go-marketplace/internal/product"
	"strings"
)

func (s *ProductService) UpdateProduct(ctx context.Context, id int, product domain.Product) error {
	if id <= 0 {
		return productfeat.ErrInvalidID
	}
	if strings.TrimSpace(product.Name) == "" {
		return productfeat.ErrInvalidName
	}
	if product.Price < 0 {
		return productfeat.ErrInvalidPrice
	}
	return s.repo.Update(ctx, id, product)
}
