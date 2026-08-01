package service

import (
	"context"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/transport/errors"
)

func (s *ProductService) GetProduct(ctx context.Context, id int) (domain.Product, error) {
	if id <= 0 {
		return domain.Product{}, core_errors.ErrInvalidID
	}

	product, err := s.repo.GetProduct(ctx, id)
	if err != nil {
		return domain.Product{}, fmt.Errorf("get product: %w", err)
	}

	return product, nil
}
