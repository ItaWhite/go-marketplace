package service

import (
	"context"
	"errors"
	"fmt"
	"go-marketplace/internal/core/errors"
	"time"
)

func (s *AuthService) Logout(ctx context.Context, refreshTokenStr string) error {
	hashedToken := HashRefreshToken(refreshTokenStr)

	refreshToken, err := s.tokenRepo.GetRefreshToken(ctx, hashedToken)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return core_errors.ErrInvalidRefreshToken
		}

		return fmt.Errorf("get refresh token: %w", err)
	}

	now := time.Now()

	if refreshToken.RevokedAt != nil {
		return fmt.Errorf("refresh token is revoked: %w", core_errors.ErrInvalidRefreshToken)
	}

	if !now.Before(refreshToken.ExpiresAt) {
		return fmt.Errorf("refresh token is expired: %w", core_errors.ErrInvalidRefreshToken)
	}

	err = s.tokenRepo.RevokeRefreshToken(ctx, refreshToken.ID)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return core_errors.ErrInvalidRefreshToken
		}

		return fmt.Errorf("revoke refresh token: %w", err)
	}

	return nil
}
