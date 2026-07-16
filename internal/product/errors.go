package product

import "errors"

var (
	ErrInvalidID    = errors.New("invalid id")
	ErrInvalidName  = errors.New("invalid name")
	ErrInvalidPrice = errors.New("invalid price")
	ErrNotFound     = errors.New("not found")
)
