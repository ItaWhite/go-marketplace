package handler

import (
	"fmt"
	"go-marketplace/internal/core/logger"
	"go-marketplace/internal/core/transport/response"
	"go-marketplace/internal/core/transport/utils"
	"net/http"
)

type GetProductResponse ProductResponse

// GetProduct godoc
// @Summary Получить товар по ID
// @Tags products
// @Produce json
// @Param id path int true "ID товара"
// @Success 200 {object} GetProductResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /products/{id} [get]
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
