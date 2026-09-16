package repository

import (
	"context"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/errors"
)

func (r *authRepository) RotateRefreshToken(ctx context.Context, oldTokenID int, newToken domain.RefreshToken) error {
	tx, err := r.db.Begin(ctx)

	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, "UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = $1 AND revoked_at IS NULL", oldTokenID)
	if err != nil {
		return fmt.Errorf("revoke old refresh token: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return core_errors.ErrNotFound
	}

	_, err = tx.Exec(ctx, "INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)",
		newToken.UserID, newToken.TokenHash, newToken.ExpiresAt)
	if err != nil {
		return fmt.Errorf("create new refresh token: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit refresh token rotation: %w", err)
	}

	return nil
}
