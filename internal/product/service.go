package product

import (
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

func (s *productService) GetAllProducts() ([]Product, error) {
	return s.repo.GetAll()
}

func (s *productService) GetProductByID(id int) (Product, error) {
	if id <= 0 {
		return Product{}, fmt.Errorf("invalid id")
	}
	return s.repo.GetByID(id)
}

func (s *productService) CreateProduct(product Product) (Product, error) {
	if strings.TrimSpace(product.Name) == "" {
		return Product{}, fmt.Errorf("invalid name")
	}
	if product.Price < 0 {
		return Product{}, fmt.Errorf("invalid price")
	}
	return s.repo.Create(product)
}

func (s *productService) DeleteProduct(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid id")
	}
	return s.repo.Delete(id)
}
