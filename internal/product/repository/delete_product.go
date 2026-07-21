package repository

import (
	"context"
	"go-marketplace/internal/core/transport/errors"
)

func (r *productRepository) DeleteProduct(ctx context.Context, id int) error {
	cmd, err := r.db.Exec(ctx, "delete from products where id = $1", id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return core_errors.ErrNotFound
	}
	return nil
}
