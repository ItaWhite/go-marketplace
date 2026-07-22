package repository

import (
	"context"
	"fmt"
	"go-marketplace/internal/core/transport/errors"
)

func (r *productRepository) DeleteProduct(ctx context.Context, id int) error {
	cmd, err := r.db.Exec(ctx, "delete from products where id = $1", id)
	if err != nil {
		return fmt.Errorf("exec sql query: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("product with id=%d not found: %w", id, core_errors.ErrNotFound)
	}

	return nil
}
