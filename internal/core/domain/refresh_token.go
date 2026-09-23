package domain

import "time"

type RefreshToken struct {
	ID        int
	UserID    int
	TokenHash string
	CreatedAt time.Time
	ExpiresAt time.Time
	RevokedAt *time.Time
}
