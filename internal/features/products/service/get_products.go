package service

import (
	"context"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/errors"
)

func (s *ProductService) GetProducts(ctx context.Context, sellerID, limit, offset *int) ([]domain.Product, error) {
	if sellerID != nil && *sellerID < 0 {
		return nil, fmt.Errorf("sellerID is negative: %w", core_errors.ErrInvalidQueryParam)
	}
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf("limit is negative: %w", core_errors.ErrInvalidQueryParam)
	}
	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("offset is negative: %w", core_errors.ErrInvalidQueryParam)
	}

	products, err := s.repo.GetProducts(ctx, sellerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get products: %w", err)
	}

	return products, nil
}
