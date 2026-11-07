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
	Name        domain.Nullable[string] `json:"name" swaggertype:"string"`
	Description domain.Nullable[string] `json:"description" swaggertype:"string"`
	Price       domain.Nullable[int]    `json:"price" swaggertype:"string"`
}

type PatchProductResponse ProductResponse

func toProductPatchDomain(dto PatchProductRequest) domain.ProductPatch {
	return domain.ProductPatch{
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
	}
}

// PatchProduct godoc
// @Summary Обновить товар
// @Description ### Логика обновления полей:
// @Description 1. **Поле отсутствует**: игнорируется, значение в БД не изменяется
// @Description 2. **Передано значение**: значение в БД изменяется на переданное
// @Description 3. **Передан `null`**: поле в БД очищается
// @Description Ограничения: `name` и `price` не могут быть `null`
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "ID товара"
// @Param request body PatchProductRequest true "Тело запроса"
// @Success 200 {object} PatchProductResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /products/{id} [patch]
func (h *ProductHandler) PatchProduct(w http.ResponseWriter, r *http.Request) {
	logger := core_logger.FromContext(r.Context())
	rh := response.NewResponseHandler(logger, w)

	productID, err := utils.GetPathValue(r, "id")
	if err != nil {
		rh.HandleError(fmt.Errorf("id param: %w", err))
		return
	}

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		err = fmt.Errorf("content type: %s: %w", r.Header.Get("Content-Type"), core_errors.ErrInvalidContentType)
		rh.HandleError(err)
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

	productResponse := PatchProductResponse(toDTO(productDomain))

	rh.SendResponse(http.StatusOK, productResponse)
}
