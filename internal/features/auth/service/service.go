package service

import (
	"context"
	"go-marketplace/internal/core/domain"
)

type AuthRepository interface {
	CreateCredentials(ctx context.Context, credentials domain.Credentials) error
	CreateRefreshToken(ctx context.Context, token domain.RefreshToken) error
}

type UserService interface {
	Create(ctx context.Context, input CreateUserInput) (int, string, error)
}

type AuthService struct {
	repo   AuthRepository
	users  UserService
	tokens *TokenService
}

func NewAuthService(repo AuthRepository, userService UserService, tokenService *TokenService) *AuthService {
	return &AuthService{
		repo:   repo,
		users:  userService,
		tokens: tokenService,
	}
}
