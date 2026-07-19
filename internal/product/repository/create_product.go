package repository

import (
	"context"
	"go-marketplace/internal/core/domain"
)

func (r *productRepository) CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	err := r.db.QueryRow(ctx, "insert into products (name, price) values ($1, $2) returning id, version, created_at",
		product.Name, product.Price).Scan(&product.ID, &product.Version, &product.CreatedAt)
	if err != nil {
		return domain.Product{}, err
	}
	return product, nil
}
