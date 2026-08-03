package repository

import (
	"context"
	"errors"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/errors"

	"github.com/jackc/pgx/v5"
)

func (r *productRepository) GetProduct(ctx context.Context, id int) (domain.Product, error) {
	row := r.db.QueryRow(ctx, "select id, version, name, description, price, created_at, seller_id from products where id = $1", id)

	var productModel ProductModel

	err := row.Scan(&productModel.ID, &productModel.Version, &productModel.Name, &productModel.Description, &productModel.Price, &productModel.CreatedAt, &productModel.SellerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Product{}, fmt.Errorf("product with id=%d not found: %w", id, core_errors.ErrNotFound)
		}

		return domain.Product{}, fmt.Errorf("scan error: %w", err)
	}

	productDomain := toDomain(productModel)

	return productDomain, nil
}
