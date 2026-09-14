package service

import (
	"context"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/errors"
	"strings"
)

type RegisterInput struct {
	Name     string
	Phone    *string
	Login    string
	Password string
	Role     string
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (domain.TokenPair, error) {
	if strings.TrimSpace(input.Login) == "" {
		return domain.TokenPair{}, fmt.Errorf("register with login=%s: %w", input.Login, core_errors.ErrInvalidLogin)
	}
	if len(strings.TrimSpace(input.Password)) < 8 || len(strings.TrimSpace(input.Password)) > 72 {
		return domain.TokenPair{}, fmt.Errorf("register with password=%s: %w", input.Login, core_errors.ErrInvalidPassword)
	}

	userID, userRole, err := s.userService.Create(ctx, CreateUserInput{
		Name:  input.Name,
		Phone: input.Phone,
		Role:  input.Role,
	})
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("create user: %w", err)
	}

	hash, err := HashPassword(input.Password)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("hash password: %w", err)
	}

	credentials := domain.Credentials{
		UserID:       userID,
		PasswordHash: hash,
		Login:        input.Login,
	}

	err = s.credRepo.CreateCredentials(ctx, credentials)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("register: %w", err)
	}

	accessTokenStr, err := s.tokenService.GenerateAccessToken(userID, userRole)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("generate access token: %w", err)
	}

	refreshTokenStr, expiresAt, err := s.tokenService.GenerateRefreshToken()
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("generate refresh token: %w", err)
	}

	refreshToken := domain.RefreshToken{
		UserID:    userID,
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
