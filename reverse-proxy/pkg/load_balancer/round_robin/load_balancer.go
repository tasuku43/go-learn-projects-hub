package round_robin

import (
	"net/url"
	"sync/atomic"
)

type LoadBalancer struct {
	Backends    []*url.URL
	NextBackend int32
}

func (r *LoadBalancer) Next() *url.URL {
	n := atomic.AddInt32(&r.NextBackend, 1)
	return r.Backends[int(n)%len(r.Backends)]
}
