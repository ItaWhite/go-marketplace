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
		dtos[i] = ToDTO(d)
	}

	return dtos
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
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

	productDomainsList, err := h.service.GetProducts(r.Context(), limit, offset)
	if err != nil {
		rh.HandleError(err)
		return
	}

	productsResponse := GetProductsResponse(toDTOs(productDomainsList))

	rh.SendResponse(http.StatusOK, productsResponse)
}
