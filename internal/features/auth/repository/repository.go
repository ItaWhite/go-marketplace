package repository

import "github.com/jackc/pgx/v5/pgxpool"

type authRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *authRepository {
	return &authRepository{
		db: db,
	}
}
