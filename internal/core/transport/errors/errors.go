package core_errors

import "errors"

var (
	ErrNullNotAllowed     = errors.New("null not allowed")
	ErrInvalidID          = errors.New("invalid id")
	ErrInvalidName        = errors.New("invalid name")
	ErrInvalidPrice       = errors.New("invalid price")
	ErrNotFound           = errors.New("not found")
	ErrInvalidContentType = errors.New("invalid content type")
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrInvalidPathValue   = errors.New("invalid path value")
	ErrInvalidQueryParam  = errors.New("invalid query param")
)
