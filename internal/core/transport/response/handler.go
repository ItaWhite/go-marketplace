package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-marketplace/internal/core/transport/errors"
	"log/slog"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type ResponseHandler struct {
	logger *slog.Logger
	w      http.ResponseWriter
}

func NewResponseHandler(logger *slog.Logger, w http.ResponseWriter) *ResponseHandler {
	return &ResponseHandler{
		logger: logger,
		w:      w,
	}
}

func (h *ResponseHandler) SendJSON(code int, v any) {
	h.w.Header().Set("Content-Type", "application/json")

	h.w.WriteHeader(code)

	err := json.NewEncoder(h.w).Encode(v)
	if err != nil {
		h.logger.Error("write HTTP response", "error", err)
	}
}

func (h *ResponseHandler) HandleError(err error) {
	var code int
	var msg string

	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		msg = "invalid argument"
		h.logger.Warn(msg, "error", err)
		code = http.StatusBadRequest
	case errors.Is(err, core_errors.ErrNullNotAllowed):
		msg = "argument can not be null"
		h.logger.Warn(msg, "error", err)
		code = http.StatusBadRequest
	case errors.Is(err, core_errors.ErrInvalidID):
		msg = "invalid id"
		h.logger.Warn(msg, "error", err)
		code = http.StatusBadRequest
	case errors.Is(err, core_errors.ErrInvalidName):
		msg = "invalid name"
		h.logger.Warn(msg, "error", err)
		code = http.StatusBadRequest
	case errors.Is(err, core_errors.ErrInvalidPrice):
		msg = "invalid price"
		h.logger.Warn(msg, "error", err)
		code = http.StatusBadRequest
	case errors.Is(err, core_errors.ErrNotFound):
		msg = "not found"
		h.logger.Warn(msg, "error", err)
		code = http.StatusNotFound
	default:
		msg = "internal server error"
		h.logger.Error(msg, "error", err)
		code = http.StatusInternalServerError
	}

	errorResponse := ErrorResponse{Error: msg}

	h.SendJSON(code, errorResponse)
}

func (h *ResponseHandler) HandlePanic(p any) {
	err := fmt.Errorf("panic: %v", p)

	h.logger.Error("handle panic", "error", err)

	errorResponse := ErrorResponse{Error: "internal server error"}

	h.SendJSON(http.StatusInternalServerError, errorResponse)
}
