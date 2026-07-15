package product

import (
	"context"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing"
)

type mockProductRepo struct {
	products map[int]Product
	id       int
}

func newMockProductRepo() *mockProductRepo {
	return &mockProductRepo{
		products: make(map[int]Product),
		id:       1,
	}
}

func (r *mockProductRepo) GetAll(ctx context.Context) ([]Product, error) {
	var products []Product
	for _, p := range r.products {
		products = append(products, p)
	}
	return products, nil
}

func (r *mockProductRepo) GetByID(ctx context.Context, id int) (Product, error) {
	return r.products[id], nil
}

func (r *mockProductRepo) Create(ctx context.Context, product Product) (Product, error) {
	product.Id = r.id
	r.products[r.id] = product
	r.id++
	return product, nil
}

func (r *mockProductRepo) Update(ctx context.Context, id int, product Product) error {
	return nil
}

func (r *mockProductRepo) Delete(ctx context.Context, id int) error {
	delete(r.products, id)
	return nil
}

func TestGetAll(t *testing.T) {
	mock := newMockProductRepo()
	s := NewProductService(mock)
	mock.Create(context.Background(), Product{Name: "Test", Price: 100})
	products, err := s.GetAllProducts(context.Background())
	require.NoError(t, err)
	assert.Len(t, products, 1)

}

func TestCreate(t *testing.T) {
	mock := newMockProductRepo()
	s := NewProductService(mock)
	product, err := s.CreateProduct(context.Background(), Product{Name: "Test", Price: 100})
	require.NoError(t, err)
	assert.Equal(t, 1, product.Id)
	assert.Equal(t, "Test", product.Name)
	assert.Equal(t, 100, product.Price)
}

func TestCreate_BlankName(t *testing.T) {
	mock := newMockProductRepo()
	s := NewProductService(mock)
	_, err := s.CreateProduct(context.Background(), Product{Name: "", Price: 100})
	require.ErrorIs(t, err, ErrInvalidName)
}

func TestCreate_NegativePrice(t *testing.T) {
	mock := newMockProductRepo()
	s := NewProductService(mock)
	_, err := s.CreateProduct(context.Background(), Product{Name: "Test", Price: -100})
	require.ErrorIs(t, err, ErrInvalidPrice)
}

func TestGetByID(t *testing.T) {
	mock := newMockProductRepo()
	s := NewProductService(mock)
	s.CreateProduct(context.Background(), Product{Name: "Test", Price: 100})
	product, err := s.GetProductByID(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, 1, product.Id)
}

func TestGetById_NegativeId(t *testing.T) {
	mock := newMockProductRepo()
	s := NewProductService(mock)
	_, err := s.GetProductByID(context.Background(), -1)
	require.ErrorIs(t, err, ErrInvalidID)
}

func TestDelete(t *testing.T) {
	mock := newMockProductRepo()
	s := NewProductService(mock)
	s.CreateProduct(context.Background(), Product{Name: "Test", Price: 100})
	err := s.DeleteProduct(context.Background(), 1)
	require.NoError(t, err)
}

func TestDelete_NegativeId(t *testing.T) {
	mock := newMockProductRepo()
	s := NewProductService(mock)
	err := s.DeleteProduct(context.Background(), -1)
	require.ErrorIs(t, err, ErrInvalidID)
}
