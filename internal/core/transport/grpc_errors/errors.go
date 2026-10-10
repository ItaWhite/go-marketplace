package grpc_errors

import (
	"errors"
	"fmt"
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

func MapInvalidArgument(violation *errdetails.BadRequest_FieldViolation) error {
	switch Reason(violation.Reason) {
	case ReasonInvalidID:
		return core_errors.ErrInvalidID
	case ReasonInvalidName:
		return core_errors.ErrInvalidName
	case ReasonInvalidRole:
		return core_errors.ErrInvalidRole
	case ReasonInvalidPhone:
		return core_errors.ErrInvalidPhone
	default:
		return fmt.Errorf("internal error")
	}
}

func ToSentinelError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return fmt.Errorf("users grpc request: %w", err)
	}

	switch st.Code() {
	case codes.InvalidArgument:
		detail := st.Details()[0]

		violation, ok := detail.(*errdetails.BadRequest_FieldViolation)
		if !ok {
			return fmt.Errorf("users grpc request: %w", err)
		}

		return fmt.Errorf("%s: %w", st.Message(), MapInvalidArgument(violation))
	case codes.NotFound:
		return core_errors.ErrNotFound
	default:
		return fmt.Errorf("users grpc request: code=%s: %s", st.Code(), st.Message())
	}
}
