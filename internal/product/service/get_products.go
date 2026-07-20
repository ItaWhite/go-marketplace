package service

import (
	"context"
	"go-marketplace/internal/core/domain"
)

func (s *ProductService) GetProducts(ctx context.Context, limit, offset int) ([]domain.Product, error) {
	return s.repo.GetProducts(ctx)
}
