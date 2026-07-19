package product

import (
	"context"
	"errors"
	"go-marketplace/internal/core/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type productRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *productRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) GetAll(ctx context.Context) ([]domain.Product, error) {
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

func (r *productRepository) GetByID(ctx context.Context, id int) (domain.Product, error) {
	var product domain.Product
	err := r.db.QueryRow(ctx, "select id, version, name, price, created_at from products where id = $1", id).
		Scan(&product.ID, &product.Version, &product.Name, &product.Price, &product.CreatedAt)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return domain.Product{}, ErrNotFound
		default:
			return domain.Product{}, err
		}
	}
	return product, nil
}

func (r *productRepository) Create(ctx context.Context, product domain.Product) (domain.Product, error) {
	err := r.db.QueryRow(ctx, "insert into products (name, price) values ($1, $2) returning id, version, created_at",
		product.Name, product.Price).Scan(&product.ID, &product.Version, &product.CreatedAt)
	if err != nil {
		return domain.Product{}, err
	}
	return product, nil
}

func (r *productRepository) Update(ctx context.Context, id int, product domain.Product) error {
	cmd, err := r.db.Exec(ctx, "update products set name=$1, price=$2 where id=$3 and version=$4",
		product.Name, product.Price, id, product.Version)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *productRepository) Delete(ctx context.Context, id int) error {
	cmd, err := r.db.Exec(ctx, "delete from products where id = $1", id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
