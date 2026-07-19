package domain

import (
	"time"
)

type Product struct {
	ID        int       `json:"id"`
	Version   int64     `json:"version"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}
