package logging

import (
	"log/slog"
	"net/http"
)

func NewMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("started handling request",
				"method", r.Method,
				"request_url", r.RequestURI,
			)
			next.ServeHTTP(w, r)
		})
	}
}
