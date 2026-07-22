package handler

import (
	"encoding/json"
	"errors"
	"go-marketplace/internal/core/domain"
	core_errors "go-marketplace/internal/core/transport/errors"
	"go-marketplace/internal/core/transport/utils"
	"log/slog"
	"net/http"
	"strings"
)

type PatchProductRequest struct {
	Name        domain.Nullable[string] `json:"name"`
	Description domain.Nullable[string] `json:"description"`
	Price       domain.Nullable[int]    `json:"price"`
}

func toProductPatchDomain(dto PatchProductRequest) domain.ProductPatch {
	return domain.ProductPatch{
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
	}
}

func (h *ProductHandler) PatchProduct(w http.ResponseWriter, r *http.Request) {
	productID, err := utils.GetPathValue(r, "id")
	if err != nil {
		slog.Warn("get path value", "error", err)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		slog.Warn("invalid content type", "content_type", r.Header.Get("Content-Type"))
		http.Error(w, "content type must be application/json", http.StatusBadRequest)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var productRequest PatchProductRequest

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err = dec.Decode(&productRequest)
	if err != nil {
		slog.Warn("invalid request body", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	productPatch := toProductPatchDomain(productRequest)

	productDomain, err := h.service.PatchProduct(r.Context(), productID, productPatch)
	if err != nil {
		switch {
		case errors.Is(err, core_errors.ErrInvalidArgument):
			slog.Warn("argument can not be null", "error", err)
			http.Error(w, "argument can not be null", http.StatusBadRequest)
		case errors.Is(err, core_errors.ErrInvalidName):
			slog.Warn("invalid name", "error", err)
			http.Error(w, "invalid name", http.StatusBadRequest)
		case errors.Is(err, core_errors.ErrInvalidPrice):
			slog.Warn("invalid price", "error", err)
			http.Error(w, "invalid price", http.StatusBadRequest)
		default:
			slog.Error("patch product", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	productResponse := ToDTO(productDomain)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(productResponse)
	if err != nil {
		slog.Error("encode product response", "error", err)
	}
}
