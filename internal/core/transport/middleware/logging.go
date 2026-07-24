package middleware

import (
	"log/slog"
	"net/http"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("HTTP Request", "Method", r.Method, "Path", r.URL, "IP", r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
