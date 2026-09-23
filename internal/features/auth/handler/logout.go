package handler

import (
	"encoding/json"
	"fmt"
	"go-marketplace/internal/core/errors"
	"go-marketplace/internal/core/logger"
	"go-marketplace/internal/core/transport/response"
	"net/http"
	"strings"
)

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	logger := core_logger.FromContext(r.Context())
	rh := response.NewResponseHandler(logger, w)

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		rh.HandleError(fmt.Errorf("content type: %s: %w", r.Header.Get("Content-Type"), core_errors.ErrInvalidContentType))
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var logoutRequest LogoutRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&logoutRequest)
	if err != nil {
		rh.HandleError(fmt.Errorf("invalid request body: %w: %v", core_errors.ErrInvalidRequestBody, err))
		return
	}

	if strings.TrimSpace(logoutRequest.RefreshToken) == "" {
		rh.HandleError(core_errors.ErrInvalidRefreshToken)
		return
	}

	err = h.service.Logout(r.Context(), logoutRequest.RefreshToken)
	if err != nil {
		rh.HandleError(err)
		return
	}

	rh.SendResponse(http.StatusNoContent, nil)
}
