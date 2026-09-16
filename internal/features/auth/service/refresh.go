package service

import (
	"context"
	"errors"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/errors"
	"time"
)

func (s *AuthService) Refresh(ctx context.Context, refreshTokenStr string) (domain.TokenPair, error) {
	hashedToken := HashRefreshToken(refreshTokenStr)

	refreshToken, err := s.tokenRepo.GetRefreshToken(ctx, hashedToken)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return domain.TokenPair{}, core_errors.ErrInvalidRefreshToken
		}

		return domain.TokenPair{}, fmt.Errorf("get refresh token: %w", err)
	}

	now := time.Now()

	if refreshToken.RevokedAt != nil {
		return domain.TokenPair{}, fmt.Errorf("refresh token is revoked: %w", core_errors.ErrInvalidRefreshToken)
	}

	if !now.Before(refreshToken.ExpiresAt) {
		return domain.TokenPair{}, fmt.Errorf("refresh token is expired: %w", core_errors.ErrInvalidRefreshToken)
	}

	role, err := s.userService.GetRole(ctx, refreshToken.UserID)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("get user role: %w", err)
	}

	accessTokenStr, err := s.tokenService.GenerateAccessToken(refreshToken.UserID, role)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("generate access token: %w", err)
	}

	newRefreshTokenStr, expiresAt, err := s.tokenService.GenerateRefreshToken()
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("generate refresh token: %w", err)
	}

	newRefreshToken := domain.RefreshToken{
		UserID:    refreshToken.UserID,
		TokenHash: HashRefreshToken(newRefreshTokenStr),
		ExpiresAt: expiresAt,
		RevokedAt: nil,
	}

	err = s.tokenRepo.RotateRefreshToken(ctx, refreshToken.ID, newRefreshToken)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("rotate refresh token: %w", err)
	}

	pair := domain.TokenPair{
		AccessToken:  accessTokenStr,
		RefreshToken: newRefreshTokenStr,
	}

	return pair, nil
}
