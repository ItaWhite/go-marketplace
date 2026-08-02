package repository

import (
	"go-marketplace/internal/core/domain"
	"time"
)

type ProductModel struct {
	ID          int
	Version     int64
	Name        string
	Description *string
	Price       int
	CreatedAt   time.Time
	SellerID    int
}

func toDomain(model ProductModel) domain.Product {
	return domain.Product{
		ID:          model.ID,
		Version:     model.Version,
		Name:        model.Name,
		Description: model.Description,
		Price:       model.Price,
		CreatedAt:   model.CreatedAt,
		SellerID:    model.SellerID,
	}
}
