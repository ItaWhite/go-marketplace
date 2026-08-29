package repository

import (
	"context"
	"errors"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func (r *productRepository) CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	row := r.db.QueryRow(ctx, "insert into products (name, description, price, seller_id) values ($1, $2, $3, $4) returning id, version, name, description, price, created_at, seller_id",
		product.Name, product.Description, product.Price, product.SellerID)

	var productModel ProductModel

	err := row.Scan(&productModel.ID, &productModel.Version, &productModel.Name, &productModel.Description, &productModel.Price, &productModel.CreatedAt, &productModel.SellerID)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return domain.Product{}, fmt.Errorf("create product: seller_id=%d: %w", product.SellerID, core_errors.ErrForeignKeyViolation)
			}
		}

		return domain.Product{}, fmt.Errorf("scan error: %w", err)
	}

	productDomain := toDomain(productModel)

	return productDomain, nil
}
