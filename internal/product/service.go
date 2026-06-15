package product

import (
	"context"
	"fmt"
	"strings"
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
		return Product{}, fmt.Errorf("invalid id")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *productService) CreateProduct(ctx context.Context, product Product) (Product, error) {
	if strings.TrimSpace(product.Name) == "" {
		return Product{}, fmt.Errorf("invalid name")
	}
	if product.Price < 0 {
		return Product{}, fmt.Errorf("invalid price")
	}
	return s.repo.Create(ctx, product)
}

func (s *productService) DeleteProduct(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid id")
	}
	return s.repo.Delete(ctx, id)
}
