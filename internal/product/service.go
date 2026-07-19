package product

import (
	"context"
	"go-marketplace/internal/core/domain"
	"strings"
)

type ProductRepository interface {
	GetAll(ctx context.Context) ([]domain.Product, error)
	GetByID(ctx context.Context, id int) (domain.Product, error)
	Create(ctx context.Context, product domain.Product) (domain.Product, error)
	Update(ctx context.Context, id int, product domain.Product) error
	Delete(ctx context.Context, id int) error
}

type ProductService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
	return s.repo.GetAll(ctx)
}

func (s *ProductService) GetProductByID(ctx context.Context, id int) (domain.Product, error) {
	if id <= 0 {
		return domain.Product{}, ErrInvalidID
	}
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	if strings.TrimSpace(product.Name) == "" {
		return domain.Product{}, ErrInvalidName
	}
	if product.Price < 0 {
		return domain.Product{}, ErrInvalidPrice
	}
	return s.repo.Create(ctx, product)
}

func (s *ProductService) UpdateProduct(ctx context.Context, id int, product domain.Product) error {
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

func (s *ProductService) DeleteProduct(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrInvalidID
	}
	return s.repo.Delete(ctx, id)
}
