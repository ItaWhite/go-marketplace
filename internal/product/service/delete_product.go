package service

import (
	"context"
	"fmt"
	"go-marketplace/internal/core/transport/errors"
)

func (s *ProductService) DeleteProduct(ctx context.Context, id int) error {
	if id <= 0 {
		return core_errors.ErrInvalidID
	}

	err := s.repo.DeleteProduct(ctx, id)
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}

	return nil
}
