package handler

import (
	"go-marketplace/internal/core/transport/middleware"
	"go-marketplace/internal/features/auth/service"
	"net/http"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: s,
	}
}

func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux, auth *middleware.AuthMiddleware) {
	mux.HandleFunc("POST /auth/register", h.Register)
}
