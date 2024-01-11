package least_connections

import (
	"fmt"
	"github.com/tasuku43/go-learn-projects-hub/waf/pkg/load_balancer/internal"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
)

type LoadBalancer struct {
	Servers []*Server
}

func newLoadBalancer(backends []*url.URL) *LoadBalancer {
	var servers []*Server
	for _, backend := range backends {
		servers = append(servers, &Server{
			url: backend,
		})
	}
	return &LoadBalancer{
		Servers: servers,
	}
}

type Server struct {
	url         *url.URL
	Connections int
	mu          sync.Mutex
}

func (s *Server) increment() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Connections++
}

func (s *Server) decrement() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Connections--
}

func (l *LoadBalancer) nextServer() *Server {
	var minServer *Server
	for _, server := range l.Servers {
		if minServer == nil {
			minServer = server
			continue
		}
		if server.Connections < minServer.Connections {
			minServer = server
		}
	}
	fmt.Printf("Next server: %s, with connections: %d\n", minServer.url.String(), minServer.Connections)

	return minServer
}

type LoadBalancerHandler struct {
	logger *slog.Logger
	lb     *LoadBalancer
}

func NewLoadBalancerHandler(logger *slog.Logger, backends []*url.URL) *LoadBalancerHandler {
	logger.Info("Initializing Least Connections Load Balancer Handler")
	lb := newLoadBalancer(backends)

	return &LoadBalancerHandler{
		logger,
		lb,
	}
}

func (h *LoadBalancerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server := h.lb.nextServer()

	server.increment()
	defer server.decrement()

	internal.CreateProxy(h.logger, server.url).ServeHTTP(w, r)
}
