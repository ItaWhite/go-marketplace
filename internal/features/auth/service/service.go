package service

import (
	"context"
	"go-marketplace/internal/core/domain"
)

type TokenRepository interface {
	CreateRefreshToken(ctx context.Context, token domain.RefreshToken) error
	GetRefreshToken(ctx context.Context, tokenHash string) (domain.RefreshToken, error)
	RotateRefreshToken(ctx context.Context, oldTokenID int, newToken domain.RefreshToken) error
	RevokeRefreshToken(ctx context.Context, id int) error
}

type CredentialsRepository interface {
	CreateCredentials(ctx context.Context, credentials domain.Credentials) error
	GetByLogin(ctx context.Context, login string) (domain.Credentials, error)
}

type UserService interface {
	Create(ctx context.Context, input CreateUserInput) (int, string, error)
	GetRole(ctx context.Context, id int) (string, error)
}

type AuthService struct {
	tokenRepo    TokenRepository
	credRepo     CredentialsRepository
	userService  UserService
	tokenService *TokenService
}

func NewAuthService(tokenRepo TokenRepository, credRepo CredentialsRepository, userService UserService, tokenService *TokenService) *AuthService {
	return &AuthService{
		tokenRepo:    tokenRepo,
		credRepo:     credRepo,
		userService:  userService,
		tokenService: tokenService,
	}
}
