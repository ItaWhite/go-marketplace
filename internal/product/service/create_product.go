package service

import (
	"context"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/transport/errors"
	"strings"
)

func (s *ProductService) CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	if strings.TrimSpace(product.Name) == "" {
		return domain.Product{}, core_errors.ErrInvalidName
	}
	if product.Price < 0 {
		return domain.Product{}, core_errors.ErrInvalidPrice
	}
	productDomain, err := s.repo.CreateProduct(ctx, product)
	if err != nil {
		return domain.Product{}, fmt.Errorf("create product: %w", err)
	}

	return productDomain, nil
}
