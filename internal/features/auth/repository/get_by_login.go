package repository

import (
	"context"
	"errors"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/errors"

	"github.com/jackc/pgx/v5"
)

func (r *authRepository) GetByLogin(ctx context.Context, login string) (domain.Credentials, error) {
	row := r.db.QueryRow(ctx, "select user_id, login, password_hash from credentials where login = $1", login)

	var credentials domain.Credentials

	err := row.Scan(&credentials.UserID, &credentials.Login, &credentials.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Credentials{}, fmt.Errorf("credentials with login=%s not found: %w", login, core_errors.ErrNotFound)
		}

		return domain.Credentials{}, fmt.Errorf("get credentials query: %w", err)
	}

	return credentials, nil
}
