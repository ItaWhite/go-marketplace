package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const LoggerKey = "logger"

func LoggerFromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(LoggerKey).(*slog.Logger)
	if !ok {
		slog.Error("error getting logger from context")
		os.Exit(1)
	}

	return logger
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestID := r.Header.Get("X-Request-ID")
		logger := slog.With("request_id", requestID, "method", r.Method, "url", r.URL.String())
		ctx := context.WithValue(r.Context(), LoggerKey, logger)

		logger.Info("request started")

		next.ServeHTTP(w, r.WithContext(ctx))

		logger.Info("request completed", "duration", time.Since(start))
	})
}
