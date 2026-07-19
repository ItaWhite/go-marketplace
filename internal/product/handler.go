package product

import (
	"encoding/json"
	"errors"
	"go-marketplace/internal/core/domain"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

type ProductHandler struct {
	service *productService
}

func NewProductHandler(s *productService) *ProductHandler {
	return &ProductHandler{
		service: s,
	}
}

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
		case errors.Is(err, ErrInvalidID):
			http.Error(w, "invalid product id", http.StatusBadRequest)
		case errors.Is(err, ErrNotFound):
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
		case errors.Is(err, ErrInvalidName):
			http.Error(w, "invalid name", http.StatusBadRequest)
		case errors.Is(err, ErrInvalidPrice):
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

func (h *ProductHandler) PutProductsHandler(w http.ResponseWriter, r *http.Request) {
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
		case errors.Is(err, ErrInvalidID):
			http.Error(w, "invalid product id", http.StatusBadRequest)
		case errors.Is(err, ErrInvalidName):
			http.Error(w, "invalid name", http.StatusBadRequest)
		case errors.Is(err, ErrInvalidPrice):
			http.Error(w, "invalid price", http.StatusBadRequest)
		case errors.Is(err, ErrNotFound):
			http.Error(w, "product not found", http.StatusNotFound)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Error("DeleteProductHandler", "error", err)
		switch {
		case errors.Is(err, ErrInvalidID):
			http.Error(w, "invalid product id", http.StatusBadRequest)
		case errors.Is(err, ErrNotFound):
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
