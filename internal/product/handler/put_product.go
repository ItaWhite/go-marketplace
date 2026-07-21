package handler

import (
	"encoding/json"
	"errors"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/transport/errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

func (h *ProductHandler) PutProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		slog.Error("PutProductsHandler", "error", err)
		return
	}
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		http.Error(w, "content type must be application/json", http.StatusBadRequest)
		slog.Error("PutProductsHandler", "error", "wrong content type")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var product domain.Product
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&product)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		slog.Error("PutProductsHandler", "error", err)
		return
	}
	if product.Name == "" || product.Price == 0 {
		http.Error(w, "all fields are required", http.StatusBadRequest)
		slog.Error("PutProductsHandler")
		return
	}

	err = h.service.UpdateProduct(r.Context(), id, product)
	if err != nil {
		slog.Error("PutProductsHandler", "error", err)
		switch {
		case errors.Is(err, core_errors.ErrInvalidID):
			http.Error(w, "invalid product id", http.StatusBadRequest)
		case errors.Is(err, core_errors.ErrInvalidName):
			http.Error(w, "invalid name", http.StatusBadRequest)
		case errors.Is(err, core_errors.ErrInvalidPrice):
			http.Error(w, "invalid price", http.StatusBadRequest)
		case errors.Is(err, core_errors.ErrNotFound):
			http.Error(w, "product not found", http.StatusNotFound)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
