package handler

import (
	"encoding/json"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/errors"
	"go-marketplace/internal/core/logger"
	"go-marketplace/internal/core/transport/response"
	"go-marketplace/internal/core/transport/utils"
	"net/http"
	"strings"
)

type PatchProductRequest struct {
	Name        domain.Nullable[string] `json:"name"`
	Description domain.Nullable[string] `json:"description"`
	Price       domain.Nullable[int]    `json:"price"`
}

type PatchProductResponse ProductResponse

func toProductPatchDomain(dto PatchProductRequest) domain.ProductPatch {
	return domain.ProductPatch{
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
	}
}

func (h *ProductHandler) PatchProduct(w http.ResponseWriter, r *http.Request) {
	logger := core_logger.FromContext(r.Context())
	rh := response.NewResponseHandler(logger, w)

	productID, err := utils.GetPathValue(r, "id")
	if err != nil {
		rh.HandleError(fmt.Errorf("get path value: %w", err))
		return
	}

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		err = fmt.Errorf("content type: %s: %w", r.Header.Get("Content-Type"), core_errors.ErrInvalidContentType)
		rh.HandleError(fmt.Errorf(""))
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var productRequest PatchProductRequest

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err = dec.Decode(&productRequest)
	if err != nil {
		rh.HandleError(fmt.Errorf("invalid request body: %w: %v", core_errors.ErrInvalidRequestBody, err))
		return
	}

	productPatch := toProductPatchDomain(productRequest)

	productDomain, err := h.service.PatchProduct(r.Context(), productID, productPatch)
	if err != nil {
		rh.HandleError(err)
		return
	}

	productResponse := PatchProductResponse(ToDTO(productDomain))

	rh.SendResponse(http.StatusOK, productResponse)
}
