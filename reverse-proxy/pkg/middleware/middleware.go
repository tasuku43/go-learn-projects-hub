package middleware

import (
	"github.com/tasuku43/go-learn-projects-hub/waf/pkg/middleware/internal/firewall"
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
}
