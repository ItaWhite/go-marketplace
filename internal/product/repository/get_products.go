package repository

import (
	"context"
	"go-marketplace/internal/core/domain"
)

func (r *productRepository) GetProducts(ctx context.Context, limit, offset int) ([]domain.Product, error) {
	rows, err := r.db.Query(ctx, "select id, version, name, price, created_at  from products;")
	if err != nil {
		return []domain.Product{}, err
	}
	defer rows.Close()
	var productList []domain.Product
	for rows.Next() {
		var product domain.Product
		err = rows.Scan(&product.ID, &product.Version, &product.Name, &product.Price, &product.CreatedAt)
		if err != nil {
			return productList, err
		}
		productList = append(productList, product)
	}
	return productList, nil
}
