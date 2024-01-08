package main

import (
	"github.com/joho/godotenv"
	"github.com/tasuku43/go-learn-projects-hub/waf/pkg/middleware"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
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

	customClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,               // 各ホストへの最大アイドルコネクション数
			IdleConnTimeout:     90 * time.Second, // アイドルコネクションのタイムアウト期間
			DisableCompression:  true,             // 圧縮を無効化
			MaxIdleConnsPerHost: 10,               // ホストごとの最大アイドルコネクション数
		},
	}

	proxy := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = backendUrl.Scheme
			r.URL.Host = backendUrl.Host
		},
		Transport: customClient.Transport,
	}

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
