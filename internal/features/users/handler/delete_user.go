package handler

import (
	"fmt"
	"go-marketplace/internal/core/logger"
	"go-marketplace/internal/core/transport/response"
	"go-marketplace/internal/core/transport/utils"
	"net/http"
)

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
