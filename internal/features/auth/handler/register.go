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

type RegisterRequest struct {
	Name     string  `json:"name"`
	Phone    *string `json:"phone"`
	Login    string  `json:"login"`
	Password string  `json:"password"`
	Role     string  `json:"role"`
}

func toRegisterInput(dto RegisterRequest) service.RegisterInput {
	return service.RegisterInput{
		Name:     dto.Name,
		Phone:    dto.Phone,
		Login:    dto.Login,
		Password: dto.Password,
		Role:     dto.Role,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	logger := core_logger.FromContext(r.Context())
	rh := response.NewResponseHandler(logger, w)

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		rh.HandleError(fmt.Errorf("content type: %s: %w", r.Header.Get("Content-Type"), core_errors.ErrInvalidContentType))
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var registerRequest RegisterRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&registerRequest)
	if err != nil {
		rh.HandleError(fmt.Errorf("invalid request body: %w: %v", core_errors.ErrInvalidRequestBody, err))
		return
	}

	pair, err := h.service.Register(r.Context(), toRegisterInput(registerRequest))
	if err != nil {
		rh.HandleError(err)
		return
	}

	response := toDTO(pair)
	rh.SendResponse(http.StatusCreated, response)
}
