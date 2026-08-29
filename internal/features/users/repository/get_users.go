package repository

import (
	"context"
	"fmt"
	"go-marketplace/internal/core/domain"
)

func (r *userRepository) GetUsers(ctx context.Context, limit, offset *int) ([]domain.User, error) {
	query := "select id, version, name, phone, role, created_at from users order by id"

	var args []any

	if limit != nil {
		args = append(args, *limit)
		query += fmt.Sprintf(" limit $%d", len(args))
	}
	if offset != nil {
		args = append(args, *offset)
		query += fmt.Sprintf(" offset $%d", len(args))
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error select: %w", err)
	}
	defer rows.Close()

	var usersList []UserModel

	for rows.Next() {
		var user UserModel

		err = rows.Scan(&user.ID, &user.Version, &user.Name, &user.Phone, &user.Role, &user.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scan: %w", err)
		}

		usersList = append(usersList, user)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("rows next error: %w", rows.Err())
	}

	userDomainsList := make([]domain.User, len(usersList))

	for i, m := range usersList {
		userDomainsList[i] = domain.User{
			ID:        m.ID,
			Version:   m.Version,
			Name:      m.Name,
			Phone:     m.Phone,
			Role:      domain.UserRole(m.Role),
			CreatedAt: m.CreatedAt,
		}
	}

	return userDomainsList, nil
}
