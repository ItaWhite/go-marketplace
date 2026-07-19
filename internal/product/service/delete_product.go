package service

import (
	"context"
	productfeat "go-marketplace/internal/product"
)

func (s *ProductService) DeleteProduct(ctx context.Context, id int) error {
	if id <= 0 {
		return productfeat.ErrInvalidID
	}
	return s.repo.DeleteProduct(ctx, id)
}
