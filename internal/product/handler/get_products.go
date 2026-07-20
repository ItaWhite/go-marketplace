package handler

import (
	"encoding/json"
	"fmt"
	"go-marketplace/internal/core/domain"
	productfeat "go-marketplace/internal/product"
	"log/slog"
	"net/http"
	"strconv"
)

type GetProductsResponse []ProductResponse

func toDTOs(domains []domain.Product) []ProductResponse {
	dtos := make([]ProductResponse, len(domains))

	for i, d := range domains {
		dtos[i] = ToDTO(d)
	}

	return dtos
}

func getQueryParam(r *http.Request, key string) (int, error) {
	valueStr := r.URL.Query().Get(key)
	if valueStr == "" {
		return -1, nil
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return -2, fmt.Errorf("param %s is not integer: %v: %w", valueStr, err, productfeat.ErrInvalidArgument)
	}

	return value, nil
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	limit, err := getQueryParam(r, "limit")
	if err != nil {
		slog.Warn("invalid limit", "error", err)
		http.Error(w, "invalid limit", http.StatusBadRequest)
		return
	}

	offset, err := getQueryParam(r, "offset")
	if err != nil {
		slog.Warn("invalid offset", "error", err)
		http.Error(w, "invalid offset", http.StatusBadRequest)
		return
	}

	productDomainsList, err := h.service.GetProducts(r.Context(), limit, offset)
	if err != nil {
		slog.Error("GetProductsHandler", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	productsResponse := GetProductsResponse(toDTOs(productDomainsList))

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(productsResponse)
	if err != nil {
		slog.Error("GetProductsHandler", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
