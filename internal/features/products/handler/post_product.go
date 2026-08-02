package handler

import (
	"encoding/json"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/errors"
	"go-marketplace/internal/core/logger"
	"go-marketplace/internal/core/transport/response"
	"net/http"
	"strings"
)

type PostProductRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Price       int     `json:"price"`
	SellerID    int     `json:"seller_id"`
}

type PostProductResponse ProductResponse

func toDomain(dto PostProductRequest) domain.Product {
	return domain.Product{
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
		SellerID:    dto.SellerID,
	}
}

func (h *ProductHandler) PostProduct(w http.ResponseWriter, r *http.Request) {
	logger := core_logger.FromContext(r.Context())
	rh := response.NewResponseHandler(logger, w)

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		rh.HandleError(fmt.Errorf("content type: %s: %w", r.Header.Get("Content-Type"), core_errors.ErrInvalidContentType))
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var productRequest PostProductRequest

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&productRequest)
	if err != nil {
		rh.HandleError(fmt.Errorf("invalid request body: %w: %v", core_errors.ErrInvalidRequestBody, err))
		return
	}

	productDomain := toDomain(productRequest)

	productDomain, err = h.service.CreateProduct(r.Context(), productDomain)
	if err != nil {
		rh.HandleError(err)
		return
	}

	productResponse := PostProductResponse(toDTO(productDomain))

	rh.SendResponse(http.StatusCreated, productResponse)
}
