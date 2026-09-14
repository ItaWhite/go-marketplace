package core_errors

import "errors"

var (
	ErrNullNotAllowed      = errors.New("null not allowed")
	ErrInvalidID           = errors.New("invalid id")
	ErrInvalidName         = errors.New("invalid name")
	ErrInvalidDescription  = errors.New("invalid description")
	ErrInvalidPrice        = errors.New("invalid price")
	ErrInvalidPhone        = errors.New("invalid phone")
	ErrInvalidRole         = errors.New("invalid role")
	ErrNotFound            = errors.New("not found")
	ErrConflict            = errors.New("conflict")
	ErrForeignKeyViolation = errors.New("foreign key violation")
	ErrInvalidContentType  = errors.New("invalid content type")
	ErrInvalidRequestBody  = errors.New("invalid request body")
	ErrInvalidPathValue    = errors.New("invalid path value")
	ErrInvalidQueryParam   = errors.New("invalid query param")
	ErrInvalidLogin        = errors.New("invalid login")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrLoginAlreadyExists  = errors.New("login already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
)
