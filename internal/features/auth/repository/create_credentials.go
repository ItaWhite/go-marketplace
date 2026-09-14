package repository

import (
	"context"
	"errors"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func (r *authRepository) CreateCredentials(ctx context.Context, credentials domain.Credentials) error {
	_, err := r.db.Exec(ctx, "insert into credentials (user_id, login, password_hash) values ($1, $2, $3)",
		credentials.UserID, credentials.Login, credentials.PasswordHash)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return core_errors.ErrLoginAlreadyExists
			}
		}

		return fmt.Errorf("create credentials error: %w", err)
	}

	return nil
}
