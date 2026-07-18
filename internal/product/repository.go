package product

import (
	"context"
	"errors"

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

func (r *productRepository) GetAll(ctx context.Context) ([]Product, error) {
	rows, err := r.db.Query(ctx, "select * from products;")
	if err != nil {
		return []Product{}, err
	}
	defer rows.Close()
	var productList []Product
	for rows.Next() {
		var product Product
		err = rows.Scan(&product.Id, &product.Name, &product.Price, &product.CreatedAt)
		if err != nil {
			return productList, err
		}
		productList = append(productList, product)
	}
	return productList, nil
}

func (r *productRepository) GetByID(ctx context.Context, id int) (Product, error) {
	var product Product
	err := r.db.QueryRow(ctx, "select * from products where id = $1", id).
		Scan(&product.Id, &product.Name, &product.Price, &product.CreatedAt)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return Product{}, ErrNotFound
		default:
			return Product{}, err
		}
	}
	return product, nil
}

func (r *productRepository) Create(ctx context.Context, product Product) (Product, error) {
	err := r.db.QueryRow(ctx, "insert into products (name, price) values ($1, $2) returning id, created_at",
		product.Name, product.Price).Scan(&product.Id, &product.CreatedAt)
	if err != nil {
		return Product{}, err
	}
	return product, nil
}

func (r *productRepository) Update(ctx context.Context, id int, product Product) error {
	cmd, err := r.db.Exec(ctx, "update products set name=$1, price=$2 where id=$3",
		product.Name, product.Price, id)
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
