package main

import (
	"github.com/joho/godotenv"
	"github.com/tasuku43/go-learn-projects-hub/waf/pkg/middleware"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := godotenv.Load(); err != nil {
		logger.Warn("Error loading .env file:", err)
		return
	}

	backendUrl, err := url.Parse(os.Getenv("BACKEND_URL"))
	if err != nil {
		logger.Warn("Error parsing BACKEND_URL: %v\n", err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(backendUrl)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	chainedHandler := middleware.Chain(handler, middleware.NewMiddlewares(logger)...)

	http.Handle("/", chainedHandler)

	port := os.Getenv("PROXY_PORT")
	logger.Info("Starting server on port: " + port)
	if err = http.ListenAndServe(":"+port, nil); err != nil {
		return
	}
}
