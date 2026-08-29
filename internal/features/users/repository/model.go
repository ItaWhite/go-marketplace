package repository

import (
	"go-marketplace/internal/core/domain"
	"time"
)

type UserModel struct {
	ID        int
	Version   int64
	Name      string
	Phone     *string
	Role      string
	CreatedAt time.Time
}

func toDomain(model UserModel) domain.User {
	return domain.User{
		ID:        model.ID,
		Version:   model.Version,
		Name:      model.Name,
		Phone:     model.Phone,
		Role:      domain.UserRole(model.Role),
		CreatedAt: model.CreatedAt,
	}
}
