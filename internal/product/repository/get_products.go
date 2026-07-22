package repository

import (
	"context"
	"fmt"
	"go-marketplace/internal/core/domain"
)

func (r *productRepository) GetProducts(ctx context.Context, limit, offset int) ([]domain.Product, error) {
	query := "select id, version, name, description, price, created_at  from products order by id"

	var args []any

	if limit != 0 {
		args = append(args, limit)
		query += fmt.Sprintf(" limit $%d", len(args))
	}
	if offset != 0 {
		args = append(args, offset)
		query += fmt.Sprintf(" offset $%d", len(args))
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return []domain.Product{}, fmt.Errorf("error select: %w", err)
	}
	defer rows.Close()

	var productList []ProductModel

	for rows.Next() {
		var product ProductModel

		err = rows.Scan(&product.ID, &product.Version, &product.Name, &product.Description, &product.Price, &product.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scan: %w", err)
		}

		productList = append(productList, product)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("rows next error: %w", err)
	}

	productDomainsList := make([]domain.Product, len(productList))

	for i, m := range productList {
		productDomainsList[i] = domain.Product{
			ID:          m.ID,
			Version:     m.Version,
			Name:        m.Name,
			Description: m.Description,
			Price:       m.Price,
			CreatedAt:   m.CreatedAt,
		}
	}

	return productDomainsList, nil
}
