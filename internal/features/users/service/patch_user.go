package service

import (
	"context"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/errors"
	"regexp"
	"strings"
)

func (s *UserService) PatchUser(ctx context.Context, id int, userPatch domain.UserPatch) (domain.User, error) {
	if id <= 0 {
		return domain.User{}, fmt.Errorf("invalid id=%d: %w", id, core_errors.ErrInvalidID)
	}

	if userPatch.Name.Set {
		if userPatch.Name.Value == nil {
			return domain.User{}, fmt.Errorf("patch name to null: %w", core_errors.ErrNullNotAllowed)
		}
		if strings.TrimSpace(*userPatch.Name.Value) == "" {
			return domain.User{}, fmt.Errorf("empty name: %w", core_errors.ErrInvalidName)
		}
	}
	if userPatch.Phone.Set && userPatch.Phone.Value != nil {
		if len([]rune(*userPatch.Phone.Value)) != 11 {
			return domain.User{}, fmt.Errorf("patch user with phone=%s: %w", *userPatch.Phone.Value, core_errors.ErrInvalidPhone)
		}
		if !regexp.MustCompile(`^[0-9]+$`).MatchString(*userPatch.Phone.Value) {
			return domain.User{}, fmt.Errorf("patch user with phone=%s: %w", *userPatch.Phone.Value, core_errors.ErrInvalidPhone)
		}
	}
	if userPatch.Role.Set {
		if userPatch.Role.Value == nil {
			return domain.User{}, fmt.Errorf("patch role to null: %w", core_errors.ErrNullNotAllowed)
		}
		if !domain.UserRole(*userPatch.Role.Value).IsValid() {
			return domain.User{}, fmt.Errorf("patch user with role=%s: %w", *userPatch.Role.Value, core_errors.ErrInvalidRole)
		}
	}

	user, err := s.repo.PatchUser(ctx, id, userPatch)
	if err != nil {
		return domain.User{}, fmt.Errorf("patch user: %w", err)
	}

	return user, err
}
