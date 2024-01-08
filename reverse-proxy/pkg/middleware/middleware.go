package middleware

import (
	"github.com/tasuku43/go-learn-projects-hub/waf/pkg/middleware/internal/firewall"
	"github.com/tasuku43/go-learn-projects-hub/waf/pkg/middleware/internal/logging"
	"github.com/tasuku43/go-learn-projects-hub/waf/pkg/middleware/internal/rate_limit"
	"golang.org/x/time/rate"
	"log/slog"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	if len(middlewares) == 0 {
		return h
	}

	for _, middleware := range middlewares {
		h = middleware(h)
	}

	return h
}

func NewMiddlewares(logger *slog.Logger) []Middleware {
	return []Middleware{
		logging.NewMiddleware(logger),
		firewall.NewMiddleware(logger),
		rate_limit.NewRateLimiter(rate.NewLimiter(1, 1), logger).Middleware,
	}
}
