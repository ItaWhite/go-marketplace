package service

import (
	"context"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/errors"
	"strings"
)

func (s *ProductService) CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	if strings.TrimSpace(product.Name) == "" {
		return domain.Product{}, fmt.Errorf("create product with name=%s: %w", product.Name, core_errors.ErrInvalidName)
	}
	if product.Description != nil {
		desc := strings.TrimSpace(*product.Description)

		if desc == "" {
			return domain.Product{}, fmt.Errorf("create product with description=%s: %w", *product.Description, core_errors.ErrInvalidDescription)
		}
	}
	if product.Price <= 0 {
		return domain.Product{}, fmt.Errorf("create product with price=%d: %w", product.Price, core_errors.ErrInvalidPrice)
	}
	if product.SellerID <= 0 {
		return domain.Product{}, fmt.Errorf("create product with seller id=%d: %w", product.SellerID, core_errors.ErrInvalidID)
	}
	
	productDomain, err := s.repo.CreateProduct(ctx, product)
	if err != nil {
		return domain.Product{}, fmt.Errorf("create product: %w", err)
	}

	return productDomain, nil
}
