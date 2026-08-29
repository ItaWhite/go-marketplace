package repository

import (
	"context"
	"errors"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/errors"

	"github.com/jackc/pgx/v5"
)

func (r *userRepository) GetUser(ctx context.Context, id int) (domain.User, error) {
	row := r.db.QueryRow(ctx, "select id, version, name, phone, role, created_at from users where id = $1", id)

	var userModel UserModel

	err := row.Scan(&userModel.ID, &userModel.Version, &userModel.Name, &userModel.Phone, &userModel.Role, &userModel.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with id=%d not found: %w", id, core_errors.ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("get user query: %w", err)
	}

	userDomain := toDomain(userModel)

	return userDomain, nil
}
