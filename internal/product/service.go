package product

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrInvalidID    = errors.New("invalid id")
	ErrInvalidName  = errors.New("invalid name")
	ErrInvalidPrice = errors.New("invalid price")
)

type productService struct {
	repo *productRepository
}

func NewProductService(repo *productRepository) *productService {
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

func (s *productService) DeleteProduct(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrInvalidID
	}
	return s.repo.Delete(ctx, id)
}
