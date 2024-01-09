package load_balancer

import (
	"github.com/tasuku43/go-learn-projects-hub/waf/pkg/load_balancer/round_robin"
	"net/url"
)

type LoadBalancer interface {
	Next() *url.URL
}

func NewLoadBalancer(backends []*url.URL) LoadBalancer {
	return &round_robin.LoadBalancer{
		Backends: backends,
	}
}
