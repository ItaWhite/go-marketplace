package repository

import (
	"context"
	"fmt"
	"go-marketplace/internal/core/domain"
)

func (r *userRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	row := r.db.QueryRow(ctx, "insert into users (name, phone, role) values ($1, $2, $3) returning id, version, created_at",
		user.Name, user.Phone, user.Role)

	var userModel UserModel

	err := row.Scan(&userModel.ID, &userModel.Version, &userModel.CreatedAt)
	if err != nil {
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := domain.User{
		ID:        userModel.ID,
		Version:   userModel.Version,
		Name:      user.Name,
		Phone:     user.Phone,
		Role:      user.Role,
		CreatedAt: userModel.CreatedAt,
	}

	return userDomain, nil
}
