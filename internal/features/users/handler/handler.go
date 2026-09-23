package handler

import (
	"go-marketplace/internal/core/domain"
	"go-marketplace/internal/core/transport/middleware"
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

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux, auth *middleware.AuthMiddleware) {
	mux.Handle("GET /users", auth.RequireRole(domain.UserAdmin, http.HandlerFunc(h.GetUsers)))
	mux.Handle("GET /users/{id}", auth.RequireRole(domain.UserAdmin, http.HandlerFunc(h.GetUser)))
	mux.Handle("POST /users", auth.RequireRole(domain.UserAdmin, http.HandlerFunc(h.PostUser)))
	mux.Handle("PATCH /users/{id}", auth.RequireRole(domain.UserAdmin, http.HandlerFunc(h.PatchUser)))
	mux.Handle("DELETE /users/{id}", auth.RequireRole(domain.UserAdmin, http.HandlerFunc(h.DeleteUser)))
}
