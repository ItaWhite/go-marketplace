package handler

import (
	"fmt"
	"go-marketplace/internal/core/logger"
	"go-marketplace/internal/core/transport/response"
	"go-marketplace/internal/core/transport/utils"
	"net/http"
)

type GetUserResponse UserResponse

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	logger := core_logger.FromContext(r.Context())
	rh := response.NewResponseHandler(logger, w)

	userID, err := utils.GetPathValue(r, "id")
	if err != nil {
		rh.HandleError(fmt.Errorf("id param: %w", err))
		return
	}

	userDomain, err := h.service.GetUser(r.Context(), userID)
	if err != nil {
		rh.HandleError(err)
		return
	}

	userResponse := GetUserResponse(ToDTO(userDomain))

	rh.SendResponse(http.StatusOK, userResponse)
}
