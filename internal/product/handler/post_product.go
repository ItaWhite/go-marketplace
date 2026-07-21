package handler

import (
	"encoding/json"
	"errors"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/transport/errors"
	"log/slog"
	"net/http"
	"strings"
)

type PostProductRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Price       int     `json:"price"`
}

func toDomain(dto PostProductRequest) domain.Product {
	return domain.Product{
		Name:  dto.Name,
		Price: dto.Price,
	}
}

func (h *ProductHandler) PostProduct(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		slog.Warn("invalid content type", "content_type", r.Header.Get("Content-Type"))
		http.Error(w, "content type must be application/json", http.StatusBadRequest)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var productRequest PostProductRequest
	
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&productRequest)
	if err != nil {
		slog.Warn("invalid request body", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	productDomain := toDomain(productRequest)

	productDomain, err = h.service.CreateProduct(r.Context(), productDomain)
	if err != nil {
		switch {
		case errors.Is(err, core_errors.ErrInvalidName):
			slog.Warn("invalid name", "error", err)
			http.Error(w, "invalid name", http.StatusBadRequest)
		case errors.Is(err, core_errors.ErrInvalidPrice):
			slog.Warn("invalid price", "error", err)
			http.Error(w, "invalid price", http.StatusBadRequest)
		default:
			slog.Error("create product", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	productResponse := ToDTO(productDomain)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(productResponse)
	if err != nil {
		slog.Error("encode product response", "error", err)
	}
}
