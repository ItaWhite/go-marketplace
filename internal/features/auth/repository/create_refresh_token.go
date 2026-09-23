package repository

import (
	"context"
	"fmt"
	"go-marketplace/internal/core/domain"
)

func (r *authRepository) CreateRefreshToken(ctx context.Context, token domain.RefreshToken) error {
	_, err := r.db.Exec(ctx, "insert into refresh_tokens (user_id, token_hash, expires_at) values ($1, $2, $3)",
		token.UserID, token.TokenHash, token.ExpiresAt)
	if err != nil {
		return fmt.Errorf("create refresh token error: %w", err)
	}

	return nil
}
