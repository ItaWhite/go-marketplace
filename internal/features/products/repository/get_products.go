package repository

import (
	"context"
	"fmt"
	"go-marketplace/internal/core/domain"
)

func (r *productRepository) GetProducts(ctx context.Context, sellerID, limit, offset *int) ([]domain.Product, error) {
	query := "select id, version, name, description, price, created_at, seller_id from products"

	var args []any

	if sellerID != nil {
		args = append(args, *sellerID)
		query += fmt.Sprintf(" where seller_id = $%d", len(args))
	}

	query += " order by id"

	if limit != nil {
		args = append(args, *limit)
		query += fmt.Sprintf(" limit $%d", len(args))
	}
	if offset != nil {
		args = append(args, *offset)
		query += fmt.Sprintf(" offset $%d", len(args))
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error select: %w", err)
	}
	defer rows.Close()

	var productModels []ProductModel

	for rows.Next() {
		var product ProductModel

		err = rows.Scan(&product.ID, &product.Version, &product.Name, &product.Description, &product.Price, &product.CreatedAt, &product.SellerID)
		if err != nil {
			return nil, fmt.Errorf("error scan: %w", err)
		}

		productModels = append(productModels, product)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("rows next error: %w", rows.Err())
	}

	productDomains := make([]domain.Product, len(productModels))

	for i, m := range productModels {
		productDomains[i] = toDomain(m)
	}

	return productDomains, nil
}
