package round_robin

import (
	"github.com/tasuku43/go-learn-projects-hub/waf/pkg/load_balancer/internal"
	"log/slog"
	"net/http"
	"net/url"
	"sync/atomic"
)

type LoadBalancer struct {
	Backends    []*url.URL
	NextBackend int32
}

func newLoadBalancer(backends []*url.URL) *LoadBalancer {
	return &LoadBalancer{
		Backends: backends,
	}
}

func (r *LoadBalancer) next() *url.URL {
	n := atomic.AddInt32(&r.NextBackend, 1)
	return r.Backends[int(n)%len(r.Backends)]
}

func Handler(logger *slog.Logger, backends []*url.URL) http.Handler {
	l := newLoadBalancer(backends)

	proxy := internal.CreateProxy(logger, l.next())

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})
}
