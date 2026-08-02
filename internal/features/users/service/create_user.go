package service

import (
	"context"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/errors"
	"regexp"
	"strings"
)

func (s *UserService) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	if strings.TrimSpace(user.Name) == "" {
		return domain.User{}, fmt.Errorf("create user with name=%s: %w", user.Name, core_errors.ErrInvalidName)
	}
	if !user.Role.IsValid() {
		return domain.User{}, fmt.Errorf("create user with role=%s: %w", user.Role, core_errors.ErrInvalidRole)
	}
	if user.Phone != nil {
		if len([]rune(*user.Phone)) != 11 {
			return domain.User{}, fmt.Errorf("create user with phone=%s: %w", *user.Phone, core_errors.ErrInvalidPhone)
		}
		if !regexp.MustCompile(`^[0-9]+$`).MatchString(*user.Phone) {
			return domain.User{}, fmt.Errorf("create user with phone=%s: %w", *user.Phone, core_errors.ErrInvalidPhone)
		}
	}

	userDomain, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return userDomain, nil
}
