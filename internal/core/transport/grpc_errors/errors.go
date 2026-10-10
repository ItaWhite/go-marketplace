package grpc_errors

type Reason string

const (
	ReasonInvalidID    Reason = "INVALID_ID"
	ReasonInvalidName  Reason = "INVALID_NAME"
	ReasonInvalidRole  Reason = "INVALID_ROLE"
	ReasonInvalidPhone Reason = "INVALID_PHONE"
)
