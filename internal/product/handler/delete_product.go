package handler

import (
	"errors"
	productfeat "go-marketplace/internal/product"
	"log/slog"
	"net/http"
	"strconv"
)

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Error("DeleteProductHandler", "error", err)
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
	err = h.service.DeleteProduct(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Error("DeleteProductHandler", "error", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
