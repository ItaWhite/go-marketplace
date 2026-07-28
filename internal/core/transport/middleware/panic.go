package middleware

import (
	"go-marketplace/internal/core/logger"
	"go-marketplace/internal/core/transport/response"
	"net/http"
)

func Panic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := core_logger.FromContext(r.Context())
		responseHandler := response.NewResponseHandler(logger, w)

		defer func() {
			err := recover()
			if err != nil {
				responseHandler.HandlePanic(err)
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
