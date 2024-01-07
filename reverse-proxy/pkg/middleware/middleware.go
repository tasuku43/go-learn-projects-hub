package middleware

import "net/http"

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
	// Noop
	// TODO: ここにミドルウェアを追加していく
	func(h http.Handler) http.Handler {
		return h
	},
}
