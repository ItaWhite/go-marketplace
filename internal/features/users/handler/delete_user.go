package handler

import (
	"fmt"
	"go-marketplace/internal/core/logger"
	"go-marketplace/internal/core/transport/response"
	"go-marketplace/internal/core/transport/utils"
	"net/http"
)

// DeleteUser godoc
// @Title Удалить пользователя
// @Tags users
// @Param id path int true "ID пользователя"
// @Success 204
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	logger := core_logger.FromContext(r.Context())
	rh := response.NewResponseHandler(logger, w)

	userID, err := utils.GetPathValue(r, "id")
	if err != nil {
		rh.HandleError(fmt.Errorf("id param: %w", err))
		return
	}

	err = h.service.DeleteUser(r.Context(), userID)
	if err != nil {
		rh.HandleError(err)
		return
	}

	rh.SendResponse(http.StatusNoContent, nil)
}
