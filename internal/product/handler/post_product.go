package handler

import (
	"encoding/json"
	"errors"
	"go-marketplace/internal/core/domain"
	productfeat "go-marketplace/internal/product"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type PostProductRequest struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type PostProductResponse struct {
	ID        int       `json:"id"`
	Version   int64     `json:"version"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

func toDomain(dto PostProductRequest) domain.Product {
	return domain.Product{
		Name:  dto.Name,
		Price: dto.Price,
	}
}

func toDTO(domain domain.Product) PostProductResponse {
	return PostProductResponse{
		ID:        domain.ID,
		Version:   domain.Version,
		Name:      domain.Name,
		Price:     domain.Price,
		CreatedAt: domain.CreatedAt,
	}
}

func (h *ProductHandler) PostProduct(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		slog.Error("PostProductsHandler", "error", "wrong content type")
		http.Error(w, "content type must be application/json", http.StatusBadRequest)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var productRequest PostProductRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&productRequest)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		slog.Error("PostProductsHandler", "error", err)
		return
	}

	productDomain := toDomain(productRequest)

	productDomain, err = h.service.CreateProduct(r.Context(), productDomain)
	if err != nil {
		slog.Error("create product failed", "error", err)
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

	productResponse := toDTO(productDomain)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(productResponse)
	if err != nil {
		slog.Error("encode product response failed", "error", err)
	}
}
