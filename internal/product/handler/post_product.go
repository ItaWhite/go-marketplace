package handler

import (
	"encoding/json"
	"errors"
	"go-marketplace/internal/core/domain"
	productfeat "go-marketplace/internal/product"
	"log/slog"
	"net/http"
	"strings"
)

func (h *ProductHandler) PostProductsHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		slog.Error("PostProductsHandler", "error", "wrong content type")
		http.Error(w, "content type must be application/json", http.StatusBadRequest)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var product domain.Product
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&product)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		slog.Error("PostProductsHandler", "error", err)
		return
	}
	if product.Name == "" || product.Price == 0 {
		http.Error(w, "all fields are required", http.StatusBadRequest)
		return
	}

	product, err = h.service.CreateProduct(r.Context(), product)
	if err != nil {
		slog.Error("PostProductsHandler", "error", err)
		switch {
		case errors.Is(err, productfeat.ErrInvalidName):
			http.Error(w, "invalid name", http.StatusBadRequest)
		case errors.Is(err, productfeat.ErrInvalidPrice):
			http.Error(w, "invalid price", http.StatusBadRequest)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)

		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		slog.Error("PostProductsHandler", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
