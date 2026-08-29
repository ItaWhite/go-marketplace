package service

import (
	"context"
	"errors"
	"testing"

	"go-marketplace/internal/core/domain"
	core_errors "go-marketplace/internal/core/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockProductRepo struct {
	products map[int]domain.Product
	nextID   int

	createErr error
	getErr    error
	getAllErr error
	patchErr  error
	deleteErr error
}

func newMockProductRepo() *mockProductRepo {
	return &mockProductRepo{
		products: make(map[int]domain.Product),
		nextID:   1,
	}
}

func (r *mockProductRepo) GetProducts(
	ctx context.Context,
	sellerID, limit, offset *int,
) ([]domain.Product, error) {
	if r.getAllErr != nil {
		return nil, r.getAllErr
	}

	products := make([]domain.Product, 0, len(r.products))

	for _, product := range r.products {
		if sellerID != nil && product.SellerID != *sellerID {
			continue
		}

		products = append(products, product)
	}

	return products, nil
}

func (r *mockProductRepo) GetProduct(
	ctx context.Context,
	id int,
) (domain.Product, error) {
	if r.getErr != nil {
		return domain.Product{}, r.getErr
	}

	product, ok := r.products[id]
	if !ok {
		return domain.Product{}, core_errors.ErrNotFound
	}

	return product, nil
}

func (r *mockProductRepo) CreateProduct(
	ctx context.Context,
	product domain.Product,
) (domain.Product, error) {
	if r.createErr != nil {
		return domain.Product{}, r.createErr
	}

	product.ID = r.nextID
	product.Version = 1

	r.products[product.ID] = product
	r.nextID++

	return product, nil
}

func (r *mockProductRepo) PatchProduct(
	ctx context.Context,
	id int,
	patch domain.ProductPatch,
) (domain.Product, error) {
	if r.patchErr != nil {
		return domain.Product{}, r.patchErr
	}

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

func (r *mockProductRepo) DeleteProduct(
	ctx context.Context,
	id int,
) error {
	if r.deleteErr != nil {
		return r.deleteErr
	}

	if _, ok := r.products[id]; !ok {
		return core_errors.ErrNotFound
	}

	delete(r.products, id)

	return nil
}

func TestProductService_CreateProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockProductRepo()
		service := NewProductService(repo)

		description := "Test description"

		product, err := service.CreateProduct(
			context.Background(),
			domain.Product{
				Name:        "Test product",
				Description: &description,
				Price:       100,
				SellerID:    1,
			},
		)

		require.NoError(t, err)
		assert.Equal(t, 1, product.ID)
		assert.Equal(t, int64(1), product.Version)
		assert.Equal(t, "Test product", product.Name)
		assert.Equal(t, 100, product.Price)
		assert.Equal(t, 1, product.SellerID)
		assert.Equal(t, &description, product.Description)
	})

	t.Run("blank name", func(t *testing.T) {
		repo := newMockProductRepo()
		service := NewProductService(repo)

		_, err := service.CreateProduct(
			context.Background(),
			domain.Product{
				Name:     "   ",
				Price:    100,
				SellerID: 1,
			},
		)

		require.ErrorIs(t, err, core_errors.ErrInvalidName)
	})

	t.Run("invalid description", func(t *testing.T) {
		repo := newMockProductRepo()
		service := NewProductService(repo)

		description := "   "

		_, err := service.CreateProduct(
			context.Background(),
			domain.Product{
				Name:        "Test",
				Description: &description,
				Price:       100,
				SellerID:    1,
			},
		)

		require.ErrorIs(t, err, core_errors.ErrInvalidDescription)
	})

	t.Run("zero price", func(t *testing.T) {
		repo := newMockProductRepo()
		service := NewProductService(repo)

		_, err := service.CreateProduct(
			context.Background(),
			domain.Product{
				Name:     "Test",
				Price:    0,
				SellerID: 1,
			},
		)

		require.ErrorIs(t, err, core_errors.ErrInvalidPrice)
	})

	t.Run("negative price", func(t *testing.T) {
		repo := newMockProductRepo()
		service := NewProductService(repo)

		_, err := service.CreateProduct(
			context.Background(),
			domain.Product{
				Name:     "Test",
				Price:    -100,
				SellerID: 1,
			},
		)

		require.ErrorIs(t, err, core_errors.ErrInvalidPrice)
	})

	t.Run("invalid seller id", func(t *testing.T) {
		repo := newMockProductRepo()
		service := NewProductService(repo)

		_, err := service.CreateProduct(
			context.Background(),
			domain.Product{
				Name:     "Test",
				Price:    100,
				SellerID: 0,
			},
		)

		require.ErrorIs(t, err, core_errors.ErrInvalidID)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := newMockProductRepo()
		repo.createErr = errors.New("database error")

		service := NewProductService(repo)

		_, err := service.CreateProduct(
			context.Background(),
			domain.Product{
				Name:     "Test",
				Price:    100,
				SellerID: 1,
			},
		)

		require.Error(t, err)
		assert.ErrorContains(t, err, "database error")
	})
}

func TestProductService_GetProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockProductRepo()
		repo.products[1] = domain.Product{
			ID:       1,
			Name:     "Test",
			Price:    100,
			SellerID: 2,
		}

		service := NewProductService(repo)

		product, err := service.GetProduct(context.Background(), 1)

		require.NoError(t, err)
		assert.Equal(t, 1, product.ID)
		assert.Equal(t, "Test", product.Name)
	})

	t.Run("zero id", func(t *testing.T) {
		repo := newMockProductRepo()
		service := NewProductService(repo)

		_, err := service.GetProduct(context.Background(), 0)

		require.ErrorIs(t, err, core_errors.ErrInvalidID)
	})

	t.Run("negative id", func(t *testing.T) {
		repo := newMockProductRepo()
		service := NewProductService(repo)

		_, err := service.GetProduct(context.Background(), -1)

		require.ErrorIs(t, err, core_errors.ErrInvalidID)
	})

	t.Run("not found", func(t *testing.T) {
		repo := newMockProductRepo()
		service := NewProductService(repo)

		_, err := service.GetProduct(context.Background(), 1)

		require.ErrorIs(t, err, core_errors.ErrNotFound)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := newMockProductRepo()
		repo.getErr = errors.New("database error")

		service := NewProductService(repo)

		_, err := service.GetProduct(context.Background(), 1)

		require.Error(t, err)
		assert.ErrorContains(t, err, "database error")
	})
}

func TestProductService_GetProducts(t *testing.T) {
	t.Run("success without filters", func(t *testing.T) {
		repo := newMockProductRepo()

		repo.products[1] = domain.Product{
			ID:       1,
			SellerID: 10,
			Name:     "Product 1",
		}
		repo.products[2] = domain.Product{
			ID:       2,
			SellerID: 20,
			Name:     "Product 2",
		}

		service := NewProductService(repo)

		products, err := service.GetProducts(
			context.Background(),
			nil,
			nil,
			nil,
		)

		require.NoError(t, err)
		assert.Len(t, products, 2)
	})

	t.Run("negative seller id", func(t *testing.T) {
		repo := newMockProductRepo()
		service := NewProductService(repo)

		sellerID := -1

		_, err := service.GetProducts(
			context.Background(),
			&sellerID,
			nil,
			nil,
		)

		require.ErrorIs(t, err, core_errors.ErrInvalidQueryParam)
	})

	t.Run("negative limit", func(t *testing.T) {
		repo := newMockProductRepo()
		service := NewProductService(repo)

		limit := -1

		_, err := service.GetProducts(
			context.Background(),
			nil,
			&limit,
			nil,
		)

		require.ErrorIs(t, err, core_errors.ErrInvalidQueryParam)
	})

	t.Run("negative offset", func(t *testing.T) {
		repo := newMockProductRepo()
		service := NewProductService(repo)

		offset := -1

		_, err := service.GetProducts(
			context.Background(),
			nil,
			nil,
			&offset,
		)

		require.ErrorIs(t, err, core_errors.ErrInvalidQueryParam)
	})

	t.Run("seller filter", func(t *testing.T) {
		repo := newMockProductRepo()

		repo.products[1] = domain.Product{
			ID:       1,
			SellerID: 10,
		}
		repo.products[2] = domain.Product{
			ID:       2,
			SellerID: 20,
		}

		service := NewProductService(repo)

		sellerID := 10

		products, err := service.GetProducts(
			context.Background(),
			&sellerID,
			nil,
			nil,
		)

		require.NoError(t, err)
		require.Len(t, products, 1)
		assert.Equal(t, 10, products[0].SellerID)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := newMockProductRepo()
		repo.getAllErr = errors.New("database error")

		service := NewProductService(repo)

		_, err := service.GetProducts(
			context.Background(),
			nil,
			nil,
			nil,
		)

		require.Error(t, err)
		assert.ErrorContains(t, err, "database error")
	})
}

func TestProductService_PatchProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockProductRepo()

		repo.products[1] = domain.Product{
			ID:      1,
			Version: 1,
			Name:    "Old",
			Price:   100,
		}

		service := NewProductService(repo)

		name := "New"

		product, err := service.PatchProduct(
			context.Background(),
			1,
			domain.ProductPatch{
				Name: domain.Nullable[string]{
					Set:   true,
					Value: &name,
				},
			},
		)

		require.NoError(t, err)
		assert.Equal(t, "New", product.Name)
		assert.Equal(t, int64(2), product.Version)
	})

	t.Run("invalid id", func(t *testing.T) {
		repo := newMockProductRepo()
		service := NewProductService(repo)

		_, err := service.PatchProduct(
			context.Background(),
			0,
			domain.ProductPatch{},
		)

		require.ErrorIs(t, err, core_errors.ErrInvalidID)
	})

	t.Run("name null", func(t *testing.T) {
		repo := newMockProductRepo()
		service := NewProductService(repo)

		_, err := service.PatchProduct(
			context.Background(),
			1,
			domain.ProductPatch{
				Name: domain.Nullable[string]{
					Set:   true,
					Value: nil,
				},
			},
		)

		require.ErrorIs(t, err, core_errors.ErrNullNotAllowed)
	})

	t.Run("blank name", func(t *testing.T) {
		repo := newMockProductRepo()
		service := NewProductService(repo)

		name := "   "

		_, err := service.PatchProduct(
			context.Background(),
			1,
			domain.ProductPatch{
				Name: domain.Nullable[string]{
					Set:   true,
					Value: &name,
				},
			},
		)

		require.ErrorIs(t, err, core_errors.ErrInvalidName)
	})

	t.Run("price null", func(t *testing.T) {
		repo := newMockProductRepo()
		service := NewProductService(repo)

		_, err := service.PatchProduct(
			context.Background(),
			1,
			domain.ProductPatch{
				Price: domain.Nullable[int]{
					Set:   true,
					Value: nil,
				},
			},
		)

		require.ErrorIs(t, err, core_errors.ErrNullNotAllowed)
	})

	t.Run("zero price", func(t *testing.T) {
		repo := newMockProductRepo()
		service := NewProductService(repo)

		price := 0

		_, err := service.PatchProduct(
			context.Background(),
			1,
			domain.ProductPatch{
				Price: domain.Nullable[int]{
					Set:   true,
					Value: &price,
				},
			},
		)

		require.ErrorIs(t, err, core_errors.ErrInvalidPrice)
	})

	t.Run("negative price", func(t *testing.T) {
		repo := newMockProductRepo()
		service := NewProductService(repo)

		price := -100

		_, err := service.PatchProduct(
			context.Background(),
			1,
			domain.ProductPatch{
				Price: domain.Nullable[int]{
					Set:   true,
					Value: &price,
				},
			},
		)

		require.ErrorIs(t, err, core_errors.ErrInvalidPrice)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := newMockProductRepo()
		repo.patchErr = errors.New("database error")

		service := NewProductService(repo)

		name := "New"

		_, err := service.PatchProduct(
			context.Background(),
			1,
			domain.ProductPatch{
				Name: domain.Nullable[string]{
					Set:   true,
					Value: &name,
				},
			},
		)

		require.Error(t, err)
		assert.ErrorContains(t, err, "database error")
	})
}

func TestProductService_DeleteProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockProductRepo()
		repo.products[1] = domain.Product{ID: 1}

		service := NewProductService(repo)

		err := service.DeleteProduct(context.Background(), 1)

		require.NoError(t, err)
		assert.NotContains(t, repo.products, 1)
	})

	t.Run("zero id", func(t *testing.T) {
		repo := newMockProductRepo()
		service := NewProductService(repo)

		err := service.DeleteProduct(context.Background(), 0)

		require.ErrorIs(t, err, core_errors.ErrInvalidID)
	})

	t.Run("negative id", func(t *testing.T) {
		repo := newMockProductRepo()
		service := NewProductService(repo)

		err := service.DeleteProduct(context.Background(), -1)

		require.ErrorIs(t, err, core_errors.ErrInvalidID)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := newMockProductRepo()
		repo.deleteErr = errors.New("database error")

		service := NewProductService(repo)

		err := service.DeleteProduct(context.Background(), 1)

		require.Error(t, err)
		assert.ErrorContains(t, err, "database error")
	})
}
