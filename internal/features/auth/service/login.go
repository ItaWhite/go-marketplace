package service

import (
	"context"
	"errors"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/errors"
	"strings"
)

type LoginInput struct {
	Login    string
	Password string
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (domain.TokenPair, error) {
	if strings.TrimSpace(input.Login) == "" {
		return domain.TokenPair{}, fmt.Errorf("login with login=%s: %w", input.Login, core_errors.ErrInvalidLogin)
	}
	if len(strings.TrimSpace(input.Password)) < 8 || len(strings.TrimSpace(input.Password)) > 72 {
		return domain.TokenPair{}, core_errors.ErrInvalidPassword
	}

	login := strings.TrimSpace(input.Login)

	credentials, err := s.credRepo.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return domain.TokenPair{}, core_errors.ErrInvalidCredentials
		}

		return domain.TokenPair{}, fmt.Errorf("get credentials by login: %w", err)
	}

	if !CheckPasswordHash(input.Password, credentials.PasswordHash) {
		return domain.TokenPair{}, fmt.Errorf("invalid password: %w", core_errors.ErrInvalidCredentials)
	}

	role, err := s.userService.GetRole(ctx, credentials.UserID)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("get user: %w", err)
	}

	accessTokenStr, err := s.tokenService.GenerateAccessToken(credentials.UserID, role)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("generate access token: %w", err)
	}

	refreshTokenStr, expiresAt, err := s.tokenService.GenerateRefreshToken()
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("generate refresh token: %w", err)
	}

	refreshToken := domain.RefreshToken{
		UserID:    credentials.UserID,
		TokenHash: HashRefreshToken(refreshTokenStr),
		ExpiresAt: expiresAt,
		RevokedAt: nil,
	}

	err = s.tokenRepo.CreateRefreshToken(ctx, refreshToken)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("create refresh token: %w", err)
	}

	pair := domain.TokenPair{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	}

	return pair, nil
}
