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
	mux.HandleFunc("GET /users", h.GetUsers)
	mux.HandleFunc("GET /users/{id}", h.GetUser)
	mux.HandleFunc("POST /users", h.PostUser)
	mux.HandleFunc("PATCH /users/{id}", h.PatchUser)
}
