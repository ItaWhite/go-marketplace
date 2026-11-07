package handler

import (
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/logger"
	"go-marketplace/internal/core/transport/response"
	"go-marketplace/internal/core/transport/utils"
	"net/http"
)

type GetProductsResponse []ProductResponse

func toDTOs(domains []domain.Product) []ProductResponse {
	dtos := make([]ProductResponse, len(domains))

	for i, d := range domains {
		dtos[i] = toDTO(d)
	}

	return dtos
}

// GetProducts godoc
// @Summary Получить список товаров
// @Tags products
// @Produce json
// @Param limit query int false "Размер страницы с товарами"
// @Param offset query int false "Смещение страницы с товарами"
// @Success 200 {object} GetProductsResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /products [get]
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	logger := core_logger.FromContext(r.Context())
	rh := response.NewResponseHandler(logger, w)

	sellerID, err := utils.GetQueryParam(r, "seller_id")
	if err != nil {
		rh.HandleError(fmt.Errorf("seller_id: %w", err))
		return
	}

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

	productDomains, err := h.service.GetProducts(r.Context(), sellerID, limit, offset)
	if err != nil {
		rh.HandleError(err)
		return
	}

	productsResponse := GetProductsResponse(toDTOs(productDomains))

	rh.SendResponse(http.StatusOK, productsResponse)
}
