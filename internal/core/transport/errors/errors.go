package core_errors

import "errors"

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrInvalidID       = errors.New("invalid id")
	ErrInvalidName     = errors.New("invalid name")
	ErrInvalidPrice    = errors.New("invalid price")
	ErrNotFound        = errors.New("not found")
)
