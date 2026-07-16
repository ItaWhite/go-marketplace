package product

import (
	"context"
	"strings"
)

type ProductRepository interface {
	GetAll(ctx context.Context) ([]Product, error)
	GetByID(ctx context.Context, id int) (Product, error)
	Create(ctx context.Context, product Product) (Product, error)
	Update(ctx context.Context, id int, product Product) error
	Delete(ctx context.Context, id int) error
}

type productService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) *productService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) GetAllProducts(ctx context.Context) ([]Product, error) {
	return s.repo.GetAll(ctx)
}

func (s *productService) GetProductByID(ctx context.Context, id int) (Product, error) {
	if id <= 0 {
		return Product{}, ErrInvalidID
	}
	return s.repo.GetByID(ctx, id)
}

func (s *productService) CreateProduct(ctx context.Context, product Product) (Product, error) {
	if strings.TrimSpace(product.Name) == "" {
		return Product{}, ErrInvalidName
	}
	if product.Price < 0 {
		return Product{}, ErrInvalidPrice
	}
	return s.repo.Create(ctx, product)
}

func (s *productService) UpdateProduct(ctx context.Context, id int, product Product) error {
	if id <= 0 {
		return ErrInvalidID
	}
	if strings.TrimSpace(product.Name) == "" {
		return ErrInvalidName
	}
	if product.Price < 0 {
		return ErrInvalidPrice
	}
	return s.repo.Update(ctx, id, product)
}

func (s *productService) DeleteProduct(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrInvalidID
	}
	return s.repo.Delete(ctx, id)
}
