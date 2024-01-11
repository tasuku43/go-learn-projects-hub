package round_robin

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"
)

type LoadBalancer struct {
	Backends    []*url.URL
	NextBackend int32
}

func NewLoadBalancer(backends []*url.URL) *LoadBalancer {
	return &LoadBalancer{
		Backends: backends,
	}
}

func (r *LoadBalancer) Next() *url.URL {
	n := atomic.AddInt32(&r.NextBackend, 1)
	return r.Backends[int(n)%len(r.Backends)]
}

func Handler(logger *slog.Logger, backends []*url.URL) http.Handler {
	l := NewLoadBalancer(backends)

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
