package least_connections

import (
	"net/url"
	"testing"
)

func TestNextServer(t *testing.T) {
	mockURL1, _ := url.Parse("http://server1")
	mockURL2, _ := url.Parse("http://server2")
	mockURL3, _ := url.Parse("http://server3")

	lb := newLoadBalancer([]*url.URL{mockURL1, mockURL2, mockURL3})

	lb.Servers[0].Connections = 5
	lb.Servers[1].Connections = 3
	lb.Servers[2].Connections = 10

	got := lb.nextServer()
	want := lb.Servers[1]

	if got != want {
		t.Errorf("nextServer() = %v, want %v", got.url, want.url)
	}
}
