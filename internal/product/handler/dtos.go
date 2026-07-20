package handler

import (
	"go-marketplace/internal/core/domain"
	"time"
)

type ProductResponse struct {
	ID        int       `json:"id"`
	Version   int64     `json:"version"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

func ToDTO(domain domain.Product) ProductResponse {
	return ProductResponse{
		ID:        domain.ID,
		Version:   domain.Version,
		Name:      domain.Name,
		Price:     domain.Price,
		CreatedAt: domain.CreatedAt,
	}
}
