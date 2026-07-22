package repository

import (
	"context"
	"errors"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/transport/errors"

	"github.com/jackc/pgx/v5"
)

func (r *productRepository) PatchProduct(ctx context.Context, id int, productPatch domain.ProductPatch) (domain.Product, error) {
	row := r.db.QueryRow(ctx, "select id, version, name, description, price, created_at from products where id=$1", id)

	var product ProductModel

	err := row.Scan(&product.ID, &product.Version, &product.Name, &product.Description, &product.Price, &product.CreatedAt)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return domain.Product{}, fmt.Errorf("product with id=%d not found: %w", id, core_errors.ErrNotFound)
		}
		return domain.Product{}, err
	}

	if productPatch.Name.Set {
		product.Name = *productPatch.Name.Value
	}
	if productPatch.Description.Set {
		product.Description = productPatch.Description.Value
	}
	if productPatch.Price.Set {
		product.Price = *productPatch.Price.Value
	}

	row = r.db.QueryRow(ctx, "update products set version=version+1, name=$1, description=$2, price=$3 where id=$4 and version=$5 returning version",
		product.Name, product.Description, product.Price, product.ID, product.Version)
	err = row.Scan(&product.Version)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return domain.Product{}, fmt.Errorf("product with id=%d accessed from several requests: %w", id, core_errors.ErrNotFound)
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
