package service

import (
	"context"
	"go-marketplace/internal/core/domain"

	users_service "go-marketplace/internal/features/users/service"
)

type AuthUserService struct {
	users *users_service.UserService
}

func NewAuthUserService(users *users_service.UserService) *AuthUserService {
	return &AuthUserService{
		users: users,
	}
}

type CreateUserInput struct {
	Name  string
	Phone *string
	Role  string
}

func (c *AuthUserService) Create(ctx context.Context, input CreateUserInput) (int, string, error) {
	user, err := c.users.CreateUser(ctx, domain.User{
		Name:  input.Name,
		Phone: input.Phone,
		Role:  domain.UserRole(input.Role),
	})
	if err != nil {
		return 0, "", err
	}

	return user.ID, string(user.Role), nil
}

func (c *AuthUserService) GetRole(ctx context.Context, id int) (string, error) {
	user, err := c.users.GetUser(ctx, id)
	if err != nil {
		return "", err
	}

	return string(user.Role), nil
}
