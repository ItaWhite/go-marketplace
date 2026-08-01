package service

import (
	"context"
	"fmt"
	"go-marketplace/internal/core/errors"
)

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid id=%d: %w", id, core_errors.ErrInvalidID)
	}

	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}
