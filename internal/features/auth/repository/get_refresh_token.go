package repository

import (
	"context"
	"errors"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/errors"

	"github.com/jackc/pgx/v5"
)

func (r *authRepository) GetRefreshToken(ctx context.Context, tokenHash string) (domain.RefreshToken, error) {
	var token domain.RefreshToken

	err := r.db.QueryRow(ctx, "SELECT id, user_id, token_hash, created_at, expires_at, revoked_at FROM refresh_tokens WHERE token_hash = $1", tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.CreatedAt,
		&token.ExpiresAt,
		&token.RevokedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.RefreshToken{}, core_errors.ErrNotFound
		}
		return domain.RefreshToken{}, fmt.Errorf("get refresh token: %w", err)
	}

	return token, nil
}
