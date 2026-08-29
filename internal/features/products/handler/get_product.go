package handler

import (
	"fmt"
	"go-marketplace/internal/core/logger"
	"go-marketplace/internal/core/transport/response"
	"go-marketplace/internal/core/transport/utils"
	"net/http"
)

type GetProductResponse ProductResponse

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	logger := core_logger.FromContext(r.Context())
	rh := response.NewResponseHandler(logger, w)

	productID, err := utils.GetPathValue(r, "id")
	if err != nil {
		rh.HandleError(fmt.Errorf("id param: %w", err))
		return
	}

	productDomain, err := h.service.GetProduct(r.Context(), productID)
	if err != nil {
		rh.HandleError(err)
		return
	}

	productResponse := GetProductResponse(toDTO(productDomain))

	rh.SendResponse(http.StatusOK, productResponse)
}
