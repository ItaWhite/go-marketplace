package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	productfeat "go-marketplace/internal/product"
	"log/slog"
	"net/http"
	"strconv"
)

type GetProductResponse ProductResponse

func getPathValue(r *http.Request, key string) (int, error) {
	valStr := r.PathValue(key)

	if valStr == "" {
		return 0, fmt.Errorf("no value by key %v: %w", key, productfeat.ErrInvalidArgument)
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, fmt.Errorf("invalid value by key %v: %w: %w", key, err, productfeat.ErrInvalidArgument)
	}

	return val, nil
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	productID, err := getPathValue(r, "id")
	if err != nil {
		slog.Warn("get path value", "error", err)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	productDomain, err := h.service.GetProduct(r.Context(), productID)
	if err != nil {
		switch {
		case errors.Is(err, productfeat.ErrInvalidID):
			slog.Warn("invalid id", "error", err)
			http.Error(w, "invalid product id", http.StatusBadRequest)
		case errors.Is(err, productfeat.ErrNotFound):
			slog.Warn("product not found", "error", err)
			http.Error(w, "product not found", http.StatusNotFound)
		default:
			slog.Warn("get product", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	productResponse := GetProductResponse(ToDTO(productDomain))

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(productResponse)
	if err != nil {
		slog.Error("encode product response", "error", err)
	}
}
