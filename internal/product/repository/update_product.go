package repository

import (
	"context"
	"go-marketplace/internal/core/domain"
	product2 "go-marketplace/internal/product"
)

func (r *productRepository) UpdateProduct(ctx context.Context, id int, product domain.Product) error {
	cmd, err := r.db.Exec(ctx, "update products set name=$1, price=$2 where id=$3 and version=$4",
		product.Name, product.Price, id, product.Version)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return product2.ErrNotFound
	}
	return nil
}
