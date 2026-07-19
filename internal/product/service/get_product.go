package service

import (
	"context"
	"go-marketplace/internal/core/domain"
	productfeat "go-marketplace/internal/product"
)

func (s *ProductService) GetProduct(ctx context.Context, id int) (domain.Product, error) {
	if id <= 0 {
		return domain.Product{}, productfeat.ErrInvalidID
	}
	return s.repo.GetByID(ctx, id)
}
