package service

import (
	"context"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/transport/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing"
)

type mockProductRepo struct {
	products map[int]domain.Product
	nextID   int
}

func newMockProductRepo() *mockProductRepo {
	return &mockProductRepo{
		products: make(map[int]domain.Product),
		nextID:   1,
	}
}

func (r *mockProductRepo) GetProducts(ctx context.Context, limit, offset int) ([]domain.Product, error) {
	products := make([]domain.Product, 0, len(r.products))

	for _, p := range r.products {
		products = append(products, p)
	}

	return products, nil
}

func (r *mockProductRepo) GetProduct(ctx context.Context, id int) (domain.Product, error) {
	product, ok := r.products[id]
	if !ok {
		return domain.Product{}, core_errors.ErrNotFound
	}

	return product, nil
}

func (r *mockProductRepo) CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	product.ID = r.nextID
	product.Version = 1

	r.products[product.ID] = product
	r.nextID++

	return product, nil
}

func (r *mockProductRepo) PatchProduct(ctx context.Context, id int, patch domain.ProductPatch) (domain.Product, error) {
	product, ok := r.products[id]
	if !ok {
		return domain.Product{}, core_errors.ErrNotFound
	}

	if patch.Name.Set {
		product.Name = *patch.Name.Value
	}

	if patch.Description.Set {
		product.Description = patch.Description.Value
	}

	if patch.Price.Set {
		product.Price = *patch.Price.Value
	}

	product.Version++

	r.products[id] = product

	return product, nil
}

func (r *mockProductRepo) DeleteProduct(ctx context.Context, id int) error {
	delete(r.products, id)
	return nil
}

func TestGetAll(t *testing.T) {
	mock := newMockProductRepo()

	s := NewProductService(mock)
	_, err := s.CreateProduct(context.Background(), domain.Product{
		Name:  "Test",
		Price: 100,
	})
	require.NoError(t, err)
	products, err := s.GetProducts(context.Background(), 0, 0)
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
	require.ErrorIs(t, err, core_errors.ErrInvalidName)
}

func TestCreate_NegativePrice(t *testing.T) {
	mock := newMockProductRepo()

	s := NewProductService(mock)
	_, err := s.CreateProduct(context.Background(), domain.Product{Name: "Test", Price: -100})
	require.ErrorIs(t, err, core_errors.ErrInvalidPrice)
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
	require.ErrorIs(t, err, core_errors.ErrInvalidID)
}

func TestPatch(t *testing.T) {
	mock := newMockProductRepo()

	s := NewProductService(mock)

	product, _ := s.CreateProduct(context.Background(), domain.Product{
		Name:  "Old",
		Price: 100,
	})
	newName := "New"
	updated, err := s.PatchProduct(
		context.Background(),
		product.ID,
		domain.ProductPatch{
			Name: domain.Nullable[string]{
				Set:   true,
				Value: &newName,
			},
		},
	)
	require.NoError(t, err)
	assert.Equal(t, "New", updated.Name)
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
	require.ErrorIs(t, err, core_errors.ErrInvalidID)
}
