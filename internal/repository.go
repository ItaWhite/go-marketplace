package internal

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

func (r *productRepository) GetAll() ([]Product, error) {
	rows, err := r.db.Query(context.Background(), "select * from products;")
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

func (r *productRepository) GetByID(id int) (Product, error) {
	var product Product
	err := r.db.QueryRow(context.Background(), "select * from products where id = $1", id).
		Scan(&product.Id, &product.Name, &product.Price, &product.CreatedAt)
	if err != nil {
		return Product{}, err
	}
	return product, nil
}

func (r *productRepository) Create(product Product) (Product, error) {
	err := r.db.QueryRow(context.Background(), "insert into products (name, price) values ($1, $2) returning id, created_at",
		product.Name, product.Price).Scan(&product.Id, &product.CreatedAt)
	if err != nil {
		return Product{}, err
	}
	return product, nil
}

func (r *productRepository) Delete(id int) error {
	_, err := r.db.Exec(context.Background(), "delete from products where id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
