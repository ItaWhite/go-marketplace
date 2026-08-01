package handler

import (
	"encoding/json"
	"fmt"
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/logger"
	"go-marketplace/internal/core/transport/errors"
	"go-marketplace/internal/core/transport/response"
	"go-marketplace/internal/core/transport/utils"
	"net/http"
	"strings"
)

type PatchUserRequest struct {
	Name  domain.Nullable[string] `json:"name"`
	Phone domain.Nullable[string] `json:"phone"`
	Role  domain.Nullable[string] `json:"role"`
}

type PatchUserResponse UserResponse

func toUserPatchDomain(dto PatchUserRequest) domain.UserPatch {
	return domain.UserPatch{
		Name:  dto.Name,
		Phone: dto.Phone,
		Role:  dto.Role,
	}
}

func (h *UserHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	logger := core_logger.FromContext(r.Context())
	rh := response.NewResponseHandler(logger, w)

	userID, err := utils.GetPathValue(r, "id")
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

	var userRequest PatchUserRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&userRequest)
	if err != nil {
		rh.HandleError(fmt.Errorf("invalid request body: %w: %v", core_errors.ErrInvalidRequestBody, err))
		return
	}

	userPatch := toUserPatchDomain(userRequest)

	userDomain, err := h.service.PatchUser(r.Context(), userID, userPatch)
	if err != nil {
		rh.HandleError(err)
		return
	}

	userResponse := PatchUserResponse(ToDTO(userDomain))

	rh.SendResponse(http.StatusOK, userResponse)
}
