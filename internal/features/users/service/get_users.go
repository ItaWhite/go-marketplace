package service

import (
	"context"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/transport/errors"
)

func (s *UserService) GetUsers(ctx context.Context, limit, offset int) ([]domain.User, error) {
	if limit < 0 {
		return nil, fmt.Errorf("limit is negative: %w", core_errors.ErrInvalidQueryParam)
	}
	if offset < 0 {
		return nil, fmt.Errorf("offset is negative: %w", core_errors.ErrInvalidQueryParam)
	}

	users, err := s.repo.GetUsers(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}

	return users, nil
}
