package main

import (
	"github.com/joho/godotenv"
	"github.com/tasuku43/go-learn-projects-hub/waf/pkg/load_balancer/least_connections"
	"github.com/tasuku43/go-learn-projects-hub/waf/pkg/middleware"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := godotenv.Load(); err != nil {
		logger.Warn("Error loading .env file:", err)
		return
	}

	chainedHandler := middleware.Chain(
		//round_robin.NewLoadBalancerHandler(logger, getBackendUrls()),
		least_connections.NewLoadBalancerHandler(logger, getBackendUrls()),
		middleware.NewMiddlewares(logger)...,
	)

	http.Handle("/", chainedHandler)

	port := os.Getenv("PROXY_PORT")
	logger.Info("Starting server on port: " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		return
	}
}

func getBackendUrls() []*url.URL {
	backendUrlsStr := os.Getenv("BACKEND_URLS")
	if backendUrlsStr == "" {
		panic("BACKEND_URLS is not set")
	}
	urlStrings := strings.Split(backendUrlsStr, ",")

	var backendUrls []*url.URL
	for _, urlString := range urlStrings {
		backendUrl, err := url.Parse(urlString)
		if err != nil {
			panic(err)
		}
		backendUrls = append(backendUrls, backendUrl)
	}
	return backendUrls
}
