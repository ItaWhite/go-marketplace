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

func (h *ResponseHandler) SendResponse(code int, v any) {
	h.w.Header().Set("Content-Type", "application/json")

	h.w.WriteHeader(code)

	if code == http.StatusNoContent {
		return
	}

	err := json.NewEncoder(h.w).Encode(v)
	if err != nil {
		h.logger.Error("write HTTP response", "error", err)
	}
}

func (h *ResponseHandler) HandleError(err error) {
	var code int
	var msg string

	switch {
	case errors.Is(err, core_errors.ErrNullNotAllowed):
		msg = core_errors.ErrNullNotAllowed.Error()
		h.logger.Warn(msg, "error", err)
		code = http.StatusBadRequest
	case errors.Is(err, core_errors.ErrInvalidID):
		msg = core_errors.ErrInvalidID.Error()
		h.logger.Warn(msg, "error", err)
		code = http.StatusBadRequest
	case errors.Is(err, core_errors.ErrInvalidName):
		msg = core_errors.ErrInvalidName.Error()
		h.logger.Warn(msg, "error", err)
		code = http.StatusBadRequest
	case errors.Is(err, core_errors.ErrInvalidPrice):
		msg = core_errors.ErrInvalidPrice.Error()
		h.logger.Warn(msg, "error", err)
		code = http.StatusBadRequest
	case errors.Is(err, core_errors.ErrInvalidPhone):
		msg = core_errors.ErrInvalidPhone.Error()
		h.logger.Warn(msg, "error", err)
		code = http.StatusBadRequest
	case errors.Is(err, core_errors.ErrInvalidRole):
		msg = core_errors.ErrInvalidRole.Error()
		h.logger.Warn(msg, "error", err)
		code = http.StatusBadRequest
	case errors.Is(err, core_errors.ErrNotFound):
		msg = core_errors.ErrNotFound.Error()
		h.logger.Warn(msg, "error", err)
		code = http.StatusNotFound
	case errors.Is(err, core_errors.ErrInvalidContentType):
		msg = core_errors.ErrInvalidContentType.Error()
		h.logger.Warn(msg, "error", err)
		code = http.StatusBadRequest
	case errors.Is(err, core_errors.ErrInvalidRequestBody):
		msg = core_errors.ErrInvalidRequestBody.Error()
		h.logger.Warn(msg, "error", err)
		code = http.StatusBadRequest
	case errors.Is(err, core_errors.ErrInvalidPathValue):
		msg = core_errors.ErrInvalidPathValue.Error()
		h.logger.Warn(msg, "error", err)
		code = http.StatusBadRequest
	case errors.Is(err, core_errors.ErrInvalidQueryParam):
		msg = core_errors.ErrInvalidQueryParam.Error()
		h.logger.Warn(msg, "error", err)
		code = http.StatusBadRequest
	default:
		msg = "internal server error"
		h.logger.Error(msg, "error", err)
		code = http.StatusInternalServerError
	}

	errorResponse := ErrorResponse{Error: msg}

	h.SendResponse(code, errorResponse)
}

func (h *ResponseHandler) HandlePanic(p any) {
	err := fmt.Errorf("panic: %v", p)

	h.logger.Error("handle panic", "error", err)

	errorResponse := ErrorResponse{Error: "internal server error"}

	h.SendResponse(http.StatusInternalServerError, errorResponse)
}
