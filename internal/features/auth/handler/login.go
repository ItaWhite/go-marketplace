package handler

import (
	"encoding/json"
	"fmt"
	"go-marketplace/internal/core/errors"
	"go-marketplace/internal/core/logger"
	"go-marketplace/internal/core/transport/response"
	"go-marketplace/internal/features/auth/service"
	"net/http"
	"strings"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func toLoginRequest(dto LoginRequest) service.LoginInput {
	return service.LoginInput{
		Login:    dto.Login,
		Password: dto.Password,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	logger := core_logger.FromContext(r.Context())
	rh := response.NewResponseHandler(logger, w)

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		rh.HandleError(fmt.Errorf("content type: %s: %w", r.Header.Get("Content-Type"), core_errors.ErrInvalidContentType))
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var loginRequest LoginRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&loginRequest)
	if err != nil {
		rh.HandleError(fmt.Errorf("invalid request body: %w: %v", core_errors.ErrInvalidRequestBody, err))
		return
	}

	pair, err := h.service.Login(r.Context(), toLoginRequest(loginRequest))
	if err != nil {
		rh.HandleError(err)
		return
	}

	resp := toDTO(pair)
	rh.SendResponse(http.StatusOK, resp)
}
