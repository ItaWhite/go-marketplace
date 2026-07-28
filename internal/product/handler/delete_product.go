package handler

import (
	"fmt"
	"go-marketplace/internal/core/logger"
	"go-marketplace/internal/core/transport/response"
	"go-marketplace/internal/core/transport/utils"
	"net/http"
)

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	logger := core_logger.FromContext(r.Context())
	rh := response.NewResponseHandler(logger, w)

	productID, err := utils.GetPathValue(r, "id")
	if err != nil {
		rh.HandleError(fmt.Errorf("get path value: %w", err))
		return
	}

	err = h.service.DeleteProduct(r.Context(), productID)
	if err != nil {
		rh.HandleError(err)
		return
	}

	rh.SendResponse(http.StatusNoContent, nil)
}
