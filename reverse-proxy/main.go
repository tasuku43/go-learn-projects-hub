package main

import (
	"github.com/joho/godotenv"
	"github.com/tasuku43/go-learn-projects-hub/waf/pkg/load_balancer/least_connections"
	"github.com/tasuku43/go-learn-projects-hub/waf/pkg/load_balancer/round_robin"
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

	lbType := strings.ToLower(os.Getenv("LOAD_BALANCER_TYPE"))
	http.Handle("/", handler(logger, lbType))

	port := os.Getenv("PROXY_PORT")
	logger.Info("Starting server on port: " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		return
	}
}

func handler(logger *slog.Logger, lbType string) http.Handler {
	var handler http.Handler
	switch lbType {
	case "round_robin":
		logger.Info("Using Round Robin Load Balancer")
		handler = round_robin.NewLoadBalancerHandler(logger, getBackendUrls())
	case "least_connections":
		logger.Info("Using Least Connections Load Balancer")
		handler = least_connections.NewLoadBalancerHandler(logger, getBackendUrls())
	}

	chainedHandler := middleware.Chain(
		handler,
		middleware.NewMiddlewares(logger)...,
	)
	return chainedHandler
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
