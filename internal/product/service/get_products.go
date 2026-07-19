package service

import (
	"context"
	"go-marketplace/internal/core/domain"
)

func (s *ProductService) GetProducts(ctx context.Context) ([]domain.Product, error) {
	return s.repo.GetAll(ctx)
}
