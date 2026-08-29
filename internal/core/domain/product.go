package domain

import (
	"time"
)

type Product struct {
	ID          int
	Version     int64
	Name        string
	Description *string
	Price       int
	CreatedAt   time.Time
	SellerID    int
}

type ProductPatch struct {
	Name        Nullable[string]
	Description Nullable[string]
	Price       Nullable[int]
}
