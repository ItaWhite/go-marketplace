package core_swagger

import (
	"net/http"

	"github.com/swaggo/http-swagger"
)

func RegisterSwagger(mux *http.ServeMux, readDoc func() string) {
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	mux.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(readDoc()))
	})
}
