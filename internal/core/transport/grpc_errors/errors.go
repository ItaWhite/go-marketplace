package grpc_errors

import (
	"errors"
	"go-marketplace/internal/core/errors"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Reason string

const (
	ReasonInvalidID    Reason = "INVALID_ID"
	ReasonInvalidName  Reason = "INVALID_NAME"
	ReasonInvalidRole  Reason = "INVALID_ROLE"
	ReasonInvalidPhone Reason = "INVALID_PHONE"
)

func BusinessError(code codes.Code, message string, field string, reason Reason) error {
	st := status.New(code, message)

	detailed, err := st.WithDetails(&errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: string(reason),
	})
	if err != nil {
		return st.Err()
	}

	return detailed.Err()
}

func ToGRPCError(err error) error {
	switch {
	case errors.Is(err, core_errors.ErrInvalidName):
		return BusinessError(codes.InvalidArgument, err.Error(), "name", ReasonInvalidName)
	case errors.Is(err, core_errors.ErrInvalidRole):
		return BusinessError(codes.InvalidArgument, err.Error(), "role", ReasonInvalidRole)
	case errors.Is(err, core_errors.ErrInvalidPhone):
		return BusinessError(codes.InvalidArgument, err.Error(), "phone", ReasonInvalidPhone)
	case errors.Is(err, core_errors.ErrInvalidID):
		return BusinessError(codes.InvalidArgument, err.Error(), "id", ReasonInvalidID)
	case errors.Is(err, core_errors.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
