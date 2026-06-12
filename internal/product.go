package internal

import "github.com/jackc/pgx/v5/pgtype"

type Product struct {
	Id        int              `json:"id"`
	Name      string           `json:"name"`
	Price     int              `json:"price"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}
