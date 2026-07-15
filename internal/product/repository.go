package product

import (
	"context"

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

func (r *productRepository) InitSchema() error {
	_, err := r.db.Exec(context.Background(), `
create table if not exists products(
    id int generated always as identity primary key,
    name varchar(255) not null,
    price int not null,
    created_at timestamp default now()
);`)

	_, err = r.db.Exec(context.Background(), `
insert into products (name, price) values ('Test 1', 100), ('Test 2', 200);`)
	return err
}

func (r *productRepository) DropSchema() error {
	_, err := r.db.Exec(context.Background(), `
drop table if exists products;
`)
	return err
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
		return Product{}, err
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
		return err
	}
	return err
}

func (r *productRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, "delete from products where id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
