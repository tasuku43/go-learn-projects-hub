package middleware

import (
	"github.com/tasuku43/go-learn-projects-hub/waf/pkg/middleware/internal/firewall"
	"github.com/tasuku43/go-learn-projects-hub/waf/pkg/middleware/internal/rate_limit"
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

var Middlewares = []Middleware{
	firewall.Middleware,
	rate_limit.NewRateLimiter(1, 1).Middleware,
}
