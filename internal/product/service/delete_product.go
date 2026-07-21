package service

import (
	"context"
	"go-marketplace/internal/core/transport/errors"
)

func (s *ProductService) DeleteProduct(ctx context.Context, id int) error {
	if id <= 0 {
		return core_errors.ErrInvalidID
	}
	return s.repo.DeleteProduct(ctx, id)
}
