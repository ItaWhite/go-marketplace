package repository

import (
	"context"
	"errors"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/errors"

	"github.com/jackc/pgx/v5"
)

func (r *productRepository) PatchProduct(ctx context.Context, id int, productPatch domain.ProductPatch) (domain.Product, error) {
	row := r.db.QueryRow(ctx, "select id, version, name, description, price, created_at, seller_id from products where id=$1", id)

	var productModel ProductModel

	err := row.Scan(&productModel.ID, &productModel.Version, &productModel.Name, &productModel.Description, &productModel.Price, &productModel.CreatedAt, &productModel.SellerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Product{}, fmt.Errorf("product with id=%d not found: %w", id, core_errors.ErrNotFound)
		}

		return domain.Product{}, fmt.Errorf("select product query: %w", err)
	}

	if productPatch.Name.Set {
		productModel.Name = *productPatch.Name.Value
	}
	if productPatch.Description.Set {
		productModel.Description = productPatch.Description.Value
	}
	if productPatch.Price.Set {
		productModel.Price = *productPatch.Price.Value
	}

	row = r.db.QueryRow(ctx, "update products set version=version+1, name=$1, description=$2, price=$3 where id=$4 and version=$5 returning version",
		productModel.Name, productModel.Description, productModel.Price, productModel.ID, productModel.Version)
	err = row.Scan(&productModel.Version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Product{}, fmt.Errorf("product with id=%d accessed from several requests: %w", id, core_errors.ErrConflict)
		}

		return domain.Product{}, fmt.Errorf("scan error: %w", err)
	}

	productDomain := toDomain(productModel)

	return productDomain, nil
}
