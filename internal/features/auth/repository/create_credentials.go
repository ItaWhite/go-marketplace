package repository

import (
	"context"
	"fmt"
	"go-marketplace/internal/core/domain"
)

func (r *authRepository) CreateCredentials(ctx context.Context, credentials domain.Credentials) error {
	_, err := r.db.Exec(ctx, "insert into credentials (user_id, login, password_hash) values ($1, $2, $3)",
		credentials.UserID, credentials.Login, credentials.PasswordHash)
	if err != nil {
		return fmt.Errorf("create credentials error: %w", err)
	}

	return nil
}
