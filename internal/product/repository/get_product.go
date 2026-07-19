package repository

import (
	"context"
	"errors"
	"go-marketplace/internal/core/domain"
	product2 "go-marketplace/internal/product"

	"github.com/jackc/pgx/v5"
)

func (r *productRepository) GetProduct(ctx context.Context, id int) (domain.Product, error) {
	var product domain.Product
	err := r.db.QueryRow(ctx, "select id, version, name, price, created_at from products where id = $1", id).
		Scan(&product.ID, &product.Version, &product.Name, &product.Price, &product.CreatedAt)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return domain.Product{}, product2.ErrNotFound
		default:
			return domain.Product{}, err
		}
	}
	return product, nil
}
