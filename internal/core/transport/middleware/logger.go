package middleware

import (
	"context"
	"log/slog"
	"net/http"
)

const LoggerKey = "logger"

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		logger := slog.With("request_id", requestID, "method", r.Method, "url", r.URL.String())
		ctx := context.WithValue(r.Context(), LoggerKey, logger)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
