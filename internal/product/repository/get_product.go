package repository

import (
	"context"
	"errors"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/transport/errors"

	"github.com/jackc/pgx/v5"
)

func (r *productRepository) GetProduct(ctx context.Context, id int) (domain.Product, error) {
	row := r.db.QueryRow(ctx, "select id, version, name, description, price, created_at from products where id = $1", id)

	var product ProductModel

	err := row.Scan(&product.ID, &product.Version, &product.Name, &product.Description, &product.Price, &product.CreatedAt)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return domain.Product{}, fmt.Errorf("product with id=%d not found: %w", id, core_errors.ErrNotFound)
		}
		return domain.Product{}, err
	}

	productDomain := domain.Product{
		ID:          product.ID,
		Version:     product.Version,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   product.CreatedAt,
	}

	return productDomain, nil
}
