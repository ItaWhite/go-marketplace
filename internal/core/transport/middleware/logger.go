package middleware

import (
	"context"
	"go-marketplace/internal/core/logger"
	"log/slog"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestID := r.Header.Get("X-Request-ID")
		logger := slog.With("request_id", requestID, "method", r.Method, "url", r.URL.String())
		ctx := context.WithValue(r.Context(), core_logger.LoggerKey, logger)

		logger.Info("request started")

		next.ServeHTTP(w, r.WithContext(ctx))

		logger.Info("request completed", "duration", time.Since(start))
	})
}
