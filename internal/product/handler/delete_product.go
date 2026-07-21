package handler

import (
	"errors"
	"go-marketplace/internal/core/transport/errors"
	"go-marketplace/internal/core/transport/utils"
	"log/slog"
	"net/http"
)

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productID, err := utils.GetPathValue(r, "id")
	if err != nil {
		slog.Warn("get path value", "error", err)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteProduct(r.Context(), productID)
	if err != nil {
		switch {
		case errors.Is(err, core_errors.ErrInvalidID):
			slog.Warn("invalid id", "error", err)
			http.Error(w, "invalid product id", http.StatusBadRequest)
		case errors.Is(err, core_errors.ErrNotFound):
			slog.Warn("product not found", "error", err)
			http.Error(w, "product not found", http.StatusNotFound)
		default:
			slog.Error("get product", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
