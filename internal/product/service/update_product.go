package service

import (
	"context"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/transport/errors"
	"strings"
)

func (s *ProductService) UpdateProduct(ctx context.Context, id int, product domain.Product) error {
	if id <= 0 {
		return core_errors.ErrInvalidID
	}
	if strings.TrimSpace(product.Name) == "" {
		return core_errors.ErrInvalidName
	}
	if product.Price < 0 {
		return core_errors.ErrInvalidPrice
	}
	return s.repo.UpdateProduct(ctx, id, product)
}
