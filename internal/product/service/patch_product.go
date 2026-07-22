package service

import (
	"context"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/transport/errors"
	"strings"
)

func (s *ProductService) PatchProduct(ctx context.Context, id int, productPatch domain.ProductPatch) (domain.Product, error) {
	if id <= 0 {
		return domain.Product{}, core_errors.ErrInvalidID
	}
	if productPatch.Name.Set {
		if productPatch.Name.Value == nil {
			return domain.Product{}, fmt.Errorf("patch name to null: %w", core_errors.ErrNullNotAllowed)
		}
		if strings.TrimSpace(*productPatch.Name.Value) == "" {
			return domain.Product{}, core_errors.ErrInvalidName
		}
	}
	if productPatch.Price.Set {
		if productPatch.Price.Value == nil {
			return domain.Product{}, fmt.Errorf("patch price to null: %w", core_errors.ErrNullNotAllowed)
		}
		if *productPatch.Price.Value < 0 {
			return domain.Product{}, core_errors.ErrInvalidPrice
		}
	}

	productDomain, err := s.repo.PatchProduct(ctx, id, productPatch)
	if err != nil {
		return domain.Product{}, fmt.Errorf("patch product: %w", err)
	}

	return productDomain, nil
}
