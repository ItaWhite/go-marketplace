package handler

import (
	"go-marketplace/internal/features/users/service"
	"net/http"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{
		service: s,
	}
}

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /users", h.PostUser)
}
