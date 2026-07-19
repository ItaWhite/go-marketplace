package repository

import (
	"context"
	"fmt"
	"go-marketplace/internal/core/domain"
)

func (r *productRepository) CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	row := r.db.QueryRow(ctx, "insert into products (name, price) values ($1, $2) returning id, version, created_at",
		product.Name, product.Price)

	var productModel ProductModel

	err := row.Scan(&productModel.ID, &productModel.Version, &productModel.CreatedAt)
	if err != nil {
		return domain.Product{}, fmt.Errorf("scan error: %w", err)
	}

	productDomain := domain.Product{
		ID:        productModel.ID,
		Version:   productModel.Version,
		Name:      product.Name,
		Price:     product.Price,
		CreatedAt: productModel.CreatedAt,
	}

	return productDomain, nil
}
