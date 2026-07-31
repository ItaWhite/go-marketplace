package repository

import (
	"context"
	"errors"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/transport/errors"

	"github.com/jackc/pgx/v5"
)

func (r *userRepository) PatchUser(ctx context.Context, id int, userPatch domain.UserPatch) (domain.User, error) {
	row := r.db.QueryRow(ctx, "select id, version, name, phone, role, created_at from users where id = $1", id)

	var userModel UserModel

	err := row.Scan(&userModel.ID, &userModel.Version, &userModel.Name, &userModel.Phone, &userModel.Role, &userModel.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with id=%d not found: %w", id, core_errors.ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("get user query: %w", err)
	}

	if userPatch.Name.Set {
		userModel.Name = *userPatch.Name.Value
	}
	if userPatch.Phone.Set {
		userModel.Phone = userPatch.Phone.Value
	}
	if userPatch.Role.Set {
		userModel.Role = *userPatch.Role.Value
	}

	row = r.db.QueryRow(ctx, "update users set version=version+1, name=$1, phone=$2, role=$3 where id=$4 and version=$5 returning version",
		userModel.Name, userModel.Phone, userModel.Role, userModel.ID, userModel.Version)
	err = row.Scan(&userModel.Version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with id=%d accessed from several requests: %w", id, core_errors.ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := domain.User{
		ID:        userModel.ID,
		Version:   userModel.Version,
		Name:      userModel.Name,
		Phone:     userModel.Phone,
		Role:      domain.UserRole(userModel.Role),
		CreatedAt: userModel.CreatedAt,
	}

	return userDomain, nil
}
