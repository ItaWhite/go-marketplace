package product

import "github.com/jackc/pgx/v5/pgtype"

type Product struct {
	ID        int                `json:"id"`
	Version   int64              `json:"version"`
	Name      string             `json:"name"`
	Price     int                `json:"price"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}
