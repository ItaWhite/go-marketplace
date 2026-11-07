package handler

import (
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/logger"
	"go-marketplace/internal/core/transport/response"
	"go-marketplace/internal/core/transport/utils"
	"net/http"
)

type GetUsersResponse []UserResponse

func toDTOs(domains []domain.User) []UserResponse {
	dtos := make([]UserResponse, len(domains))

	for i, d := range domains {
		dtos[i] = ToDTO(d)
	}

	return dtos
}

// GetUsers godoc
// @Summary Получить список пользователей
// @Tags users
// @Produce json
// @Param limit query int false "Размер страницы с пользователями"
// @Param offset query int false "Смещение страницы с пользователями"
// @Success 200 {object} GetUsersResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users [get]
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	logger := core_logger.FromContext(r.Context())
	rh := response.NewResponseHandler(logger, w)

	limit, err := utils.GetQueryParam(r, "limit")
	if err != nil {
		rh.HandleError(fmt.Errorf("limit: %w", err))
		return
	}

	offset, err := utils.GetQueryParam(r, "offset")
	if err != nil {
		rh.HandleError(fmt.Errorf("offset: %w", err))
		return
	}

	userDomainsList, err := h.service.GetUsers(r.Context(), limit, offset)
	if err != nil {
		rh.HandleError(err)
		return
	}

	usersResponse := toDTOs(userDomainsList)

	rh.SendResponse(http.StatusOK, usersResponse)
}
