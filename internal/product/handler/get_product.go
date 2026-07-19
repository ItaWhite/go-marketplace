package handler

import (
	"encoding/json"
	"errors"
	productfeat "go-marketplace/internal/product"
	"log/slog"
	"net/http"
	"strconv"
)

func (h *ProductHandler) GetProductByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		slog.Error("GetProductByIdHandler", "error", err)
		return
	}
	product, err := h.service.GetProductByID(r.Context(), id)
	if err != nil {
		slog.Error("GetProductByIdHandler", "error", err)
		switch {
		case errors.Is(err, productfeat.ErrInvalidID):
			http.Error(w, "invalid product id", http.StatusBadRequest)
		case errors.Is(err, productfeat.ErrNotFound):
			http.Error(w, "product not found", http.StatusNotFound)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		slog.Error("GetProductByIdHandler", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
