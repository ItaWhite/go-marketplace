package handler

import (
	"encoding/json"
	"errors"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/transport/utils"
	productfeat "go-marketplace/internal/product"
	"log/slog"
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
	limit, err := utils.GetQueryParam(r, "limit")
	if err != nil {
		slog.Warn("invalid query param", "param", "limit", "error", err)
		http.Error(w, "invalid limit", http.StatusBadRequest)
		return
	}

	offset, err := utils.GetQueryParam(r, "offset")
	if err != nil {
		slog.Warn("invalid query param", "param", "offset", "error", err)
		http.Error(w, "invalid offset", http.StatusBadRequest)
		return
	}

	productDomainsList, err := h.service.GetProducts(r.Context(), limit, offset)
	if err != nil {
		switch {
		case errors.Is(err, productfeat.ErrInvalidArgument):
			slog.Warn("get products", "error", err)
			http.Error(w, "invalid limit or offset", http.StatusBadRequest)
		default:
			slog.Error("get products", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	productsResponse := GetProductsResponse(toDTOs(productDomainsList))

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(productsResponse)
	if err != nil {
		slog.Error("encode products response", "error", err)
	}
}
