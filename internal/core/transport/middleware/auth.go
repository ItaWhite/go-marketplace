package middleware

import (
	"context"
	"fmt"
	"go-marketplace/internal/core/errors"
	"go-marketplace/internal/core/logger"
	"go-marketplace/internal/core/security"
	"go-marketplace/internal/core/transport/response"
	"net/http"
	"strconv"
	"strings"
)

const prefix = "Bearer "

type AuthMiddleware struct {
	validator *security.JWTValidator
}

func NewAuthMiddleware(validator *security.JWTValidator) *AuthMiddleware {
	return &AuthMiddleware{
		validator: validator,
	}
}

func (m *AuthMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := core_logger.FromContext(r.Context())
		rh := response.NewResponseHandler(logger, w)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, prefix) {
			rh.HandleError(fmt.Errorf("invalid authorization header: %w", core_errors.ErrUnauthorized))
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))

		if tokenString == "" {
			rh.HandleError(fmt.Errorf("invalid token string: %w", core_errors.ErrUnauthorized))
			return
		}

		claims, err := m.validator.Validate(tokenString)
		if err != nil {
			rh.HandleError(fmt.Errorf("validate access token: %w", err))
			return
		}

		userID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			rh.HandleError(fmt.Errorf("invalid user ID=%s: %w", claims.Subject, core_errors.ErrUnauthorized))
			return
		}

		ctx := context.WithValue(r.Context(), security.UserIDKey, userID)
		ctx = context.WithValue(ctx, security.RoleKey, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
