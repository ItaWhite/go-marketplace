package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type productRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *productRepository {
	return &productRepository{
		db: db,
	}
}
