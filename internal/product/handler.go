package product

import (
	"encoding/json"
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
	productsList, err := h.service.GetAllProducts()
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
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "invalid product id", http.StatusBadRequest)
		slog.Error("GetProductByIdHandler", "error", err)
		return
	}
	product, err := h.service.GetProductByID(id)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Error("GetProductByIdHandler", "error", err)
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
		w.Header().Set("Content-Type", "text/plain")
		slog.Error("PostProductsHandler", "error", "wrong content type")
		http.Error(w, "content type must be application/json", http.StatusBadRequest)
		return
	}
	var product Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "invalid request body", http.StatusBadRequest)
		slog.Error("PostProductsHandler", "error", err)
		return
	}
	if product.Name == "" || product.Price == 0 {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "all fields are required", http.StatusBadRequest)
		return
	}

	product, err = h.service.CreateProduct(product)
	if err != nil {
		slog.Error("PostProductsHandler", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
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

func (h *ProductHandler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "invalid product id", http.StatusBadRequest)
		slog.Error("DeleteProductHandler", "error", err)
		return
	}
	err = h.service.DeleteProduct(id)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Error("DeleteProductHandler", "error", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
