package internal

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func CreateProxy(logger *slog.Logger, url *url.URL) *httputil.ReverseProxy {
	proxy := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			logger.Info("Proxying request to: " + url.String())
			r.URL.Scheme = url.Scheme
			r.URL.Host = url.Host
		},
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     90 * time.Second,
			DisableCompression:  true,
			MaxIdleConnsPerHost: 10,
		},
	}
	return proxy
}
