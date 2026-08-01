package handler

import (
	"go-marketplace/internal/core/domain"
	"time"
)

type UserResponse struct {
	ID        int       `json:"id"`
	Version   int64     `json:"version"`
	Name      string    `json:"name"`
	Phone     *string   `json:"phone"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func ToDTO(domain domain.User) UserResponse {
	return UserResponse{
		ID:        domain.ID,
		Version:   domain.Version,
		Name:      domain.Name,
		Phone:     domain.Phone,
		Role:      string(domain.Role),
		CreatedAt: domain.CreatedAt,
	}
}
