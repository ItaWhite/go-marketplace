package repository

import "time"

type ProductModel struct {
	ID          int
	Version     int64
	Name        string
	Description *string
	Price       int
	CreatedAt   time.Time
}
