package handler

import (
	"encoding/json"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/logger"
	"go-marketplace/internal/core/transport/errors"
	"go-marketplace/internal/core/transport/response"
	"net/http"
	"strings"
)

type PostUserRequest struct {
	Name  string  `json:"name"`
	Phone *string `json:"phone"`
	Role  string  `json:"role"`
}

type PostUserResponse UserResponse

func toDomain(dto PostUserRequest) domain.User {
	return domain.User{
		Name:  dto.Name,
		Phone: dto.Phone,
		Role:  domain.UserRole(dto.Role),
	}
}

func (h *UserHandler) PostUser(w http.ResponseWriter, r *http.Request) {
	logger := core_logger.FromContext(r.Context())
	rh := response.NewResponseHandler(logger, w)

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		rh.HandleError(fmt.Errorf("content type: %s: %w", r.Header.Get("Content-Type"), core_errors.ErrInvalidContentType))
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var userRequest PostUserRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&userRequest)
	if err != nil {
		rh.HandleError(fmt.Errorf("invalid request body: %w: %v", core_errors.ErrInvalidRequestBody, err))
		return
	}

	userDomain := toDomain(userRequest)

	userDomain, err = h.service.CreateUser(r.Context(), userDomain)
	if err != nil {
		rh.HandleError(err)
		return
	}

	userResponse := PostUserResponse(ToDTO(userDomain))

	rh.SendResponse(http.StatusCreated, userResponse)
}
