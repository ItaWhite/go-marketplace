package service

import (
	"context"
	"go-marketplace/internal/core/domain"
)

type ProductRepository interface {
	GetProducts(ctx context.Context) ([]domain.Product, error)
	GetProduct(ctx context.Context, id int) (domain.Product, error)
	CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error)
	UpdateProduct(ctx context.Context, id int, product domain.Product) error
	DeleteProduct(ctx context.Context, id int) error
}

type ProductService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}
