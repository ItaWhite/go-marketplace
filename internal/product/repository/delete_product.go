package repository

import (
	"context"
	product2 "go-marketplace/internal/product"
)

func (r *productRepository) DeleteProduct(ctx context.Context, id int) error {
	cmd, err := r.db.Exec(ctx, "delete from products where id = $1", id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return product2.ErrNotFound
	}
	return nil
}
