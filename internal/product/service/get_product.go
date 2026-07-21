package service

import (
	"context"
	"fmt"
	"go-marketplace/internal/core/domain"
	productfeat "go-marketplace/internal/product"
)

func (s *ProductService) GetProduct(ctx context.Context, id int) (domain.Product, error) {
	if id <= 0 {
		return domain.Product{}, productfeat.ErrInvalidID
	}
	product, err := s.repo.GetProduct(ctx, id)
	if err != nil {
		return domain.Product{}, fmt.Errorf("get product: %w", err)
	}

	return product, nil
}
