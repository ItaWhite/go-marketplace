package core_errors

import "errors"

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrNullNotAllowed  = errors.New("null not allowed")
	ErrInvalidID       = errors.New("invalid id")
	ErrInvalidName     = errors.New("invalid name")
	ErrInvalidPrice    = errors.New("invalid price")
	ErrNotFound        = errors.New("not found")
)
