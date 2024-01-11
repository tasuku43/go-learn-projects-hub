package main

import (
	"github.com/joho/godotenv"
	"github.com/tasuku43/go-learn-projects-hub/waf/pkg/load_balancer"
	"github.com/tasuku43/go-learn-projects-hub/waf/pkg/middleware"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	loadBalancer = load_balancer.NewLoadBalancer(getBackendUrls())
)

func handler(logger *slog.Logger, l load_balancer.LoadBalancer) http.Handler {
	proxy := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			backendUrl := l.Next()
			logger.Info("Proxying request to: " + backendUrl.String())
			r.URL.Scheme = backendUrl.Scheme
			r.URL.Host = backendUrl.Host
		},
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     90 * time.Second,
			DisableCompression:  true,
			MaxIdleConnsPerHost: 10,
		},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})
}

func main() {
	var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := godotenv.Load(); err != nil {
		logger.Warn("Error loading .env file:", err)
		return
	}

	chainedHandler := middleware.Chain(handler(logger, loadBalancer), middleware.NewMiddlewares(logger)...)

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
