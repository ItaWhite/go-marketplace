package service

import (
	"context"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/product"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing"
)

type mockProductRepo struct {
	products map[int]domain.Product
	id       int
}

func newMockProductRepo() *mockProductRepo {
	return &mockProductRepo{
		products: make(map[int]domain.Product),
		id:       1,
	}
}

func (r *mockProductRepo) GetAll(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product
	for _, p := range r.products {
		products = append(products, p)
	}
	return products, nil
}

func (r *mockProductRepo) GetByID(ctx context.Context, id int) (domain.Product, error) {
	return r.products[id], nil
}

func (r *mockProductRepo) Create(ctx context.Context, product domain.Product) (domain.Product, error) {
	product.ID = r.id
	r.products[r.id] = product
	r.id++
	return product, nil
}

func (r *mockProductRepo) Update(ctx context.Context, id int, product domain.Product) error {
	return nil
}

func (r *mockProductRepo) Delete(ctx context.Context, id int) error {
	delete(r.products, id)
	return nil
}

func TestGetAll(t *testing.T) {
	mock := newMockProductRepo()
	s := NewProductService(mock)
	mock.Create(context.Background(), domain.Product{Name: "Test", Price: 100})
	products, err := s.GetProducts(context.Background())
	require.NoError(t, err)
	assert.Len(t, products, 1)

}

func TestCreate(t *testing.T) {
	mock := newMockProductRepo()
	s := NewProductService(mock)
	product, err := s.CreateProduct(context.Background(), domain.Product{Name: "Test", Price: 100})
	require.NoError(t, err)
	assert.Equal(t, 1, product.ID)
	assert.Equal(t, "Test", product.Name)
	assert.Equal(t, 100, product.Price)
}

func TestCreate_BlankName(t *testing.T) {
	mock := newMockProductRepo()
	s := NewProductService(mock)
	_, err := s.CreateProduct(context.Background(), domain.Product{Name: "", Price: 100})
	require.ErrorIs(t, err, product.ErrInvalidName)
}

func TestCreate_NegativePrice(t *testing.T) {
	mock := newMockProductRepo()
	s := NewProductService(mock)
	_, err := s.CreateProduct(context.Background(), domain.Product{Name: "Test", Price: -100})
	require.ErrorIs(t, err, product.ErrInvalidPrice)
}

func TestGetByID(t *testing.T) {
	mock := newMockProductRepo()
	s := NewProductService(mock)
	s.CreateProduct(context.Background(), domain.Product{Name: "Test", Price: 100})
	product, err := s.GetProduct(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, 1, product.ID)
}

func TestGetById_NegativeId(t *testing.T) {
	mock := newMockProductRepo()
	s := NewProductService(mock)
	_, err := s.GetProduct(context.Background(), -1)
	require.ErrorIs(t, err, product.ErrInvalidID)
}

func TestDelete(t *testing.T) {
	mock := newMockProductRepo()
	s := NewProductService(mock)
	s.CreateProduct(context.Background(), domain.Product{Name: "Test", Price: 100})
	err := s.DeleteProduct(context.Background(), 1)
	require.NoError(t, err)
}

func TestDelete_NegativeId(t *testing.T) {
	mock := newMockProductRepo()
	s := NewProductService(mock)
	err := s.DeleteProduct(context.Background(), -1)
	require.ErrorIs(t, err, product.ErrInvalidID)
}
