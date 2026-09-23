package repository

import (
	"context"
	"fmt"
	"go-marketplace/internal/core/errors"
)

func (r *authRepository) RevokeRefreshToken(ctx context.Context, id int) error {
	tag, err := r.db.Exec(ctx, "UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = $1 AND revoked_at IS NULL", id)
	if err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return core_errors.ErrNotFound
	}

	return nil
}
