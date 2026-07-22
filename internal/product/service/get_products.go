package service

import (
	"context"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/transport/errors"
)

func (s *ProductService) GetProducts(ctx context.Context, limit, offset int) ([]domain.Product, error) {
	if limit < 0 {
		return nil, fmt.Errorf("limit is negative: %w", core_errors.ErrInvalidArgument)
	}
	if offset < 0 {
		return nil, fmt.Errorf("offset is negative: %w", core_errors.ErrInvalidArgument)
	}

	products, err := s.repo.GetProducts(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get products: %w", err)
	}

	return products, nil
}
