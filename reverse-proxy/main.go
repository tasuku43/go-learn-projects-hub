package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/tasuku43/go-learn-projects-hub/waf/pkg/middleware"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	backendUrl, err := url.Parse(os.Getenv("BACKEND_URL"))
	if err != nil {
		fmt.Printf("Error parsing BACKEND_URL: %v\n", err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(backendUrl)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	chainedHandler := middleware.Chain(handler, middleware.Middlewares...)

	http.Handle("/", chainedHandler)

	port := os.Getenv("PROXY_PORT")
	fmt.Println("Starting server on port", port)
	if err = http.ListenAndServe(":"+port, nil); err != nil {
		return
	}
}
