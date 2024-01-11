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

type LoadBalancerHandler struct {
	logger *slog.Logger
	lb     *LoadBalancer
}

func NewLoadBalancerHandler(logger *slog.Logger, backends []*url.URL) *LoadBalancerHandler {
	logger.Info("Initializing Round Robin Load Balancer Handler")
	lb := newLoadBalancer(backends)

	return &LoadBalancerHandler{
		logger,
		lb,
	}
}

func (h *LoadBalancerHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	internal.CreateProxy(h.logger, h.lb.next()).ServeHTTP(rw, req)
}
