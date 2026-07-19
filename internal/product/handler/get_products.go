package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func (h *ProductHandler) GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	productsList, err := h.service.GetAllProducts(r.Context())
	if err != nil {
		slog.Error("GetProductsHandler", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(productsList)
	if err != nil {
		slog.Error("GetProductsHandler", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
