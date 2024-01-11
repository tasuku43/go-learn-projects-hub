package round_robin

import (
	"net/url"
	"testing"
)

func TestNext(t *testing.T) {
	mockURL1, _ := url.Parse("http://server1")
	mockURL2, _ := url.Parse("http://server2")
	mockURL3, _ := url.Parse("http://server3")

	lb := newLoadBalancer([]*url.URL{mockURL1, mockURL2, mockURL3})

	got := lb.next()
	want := mockURL2

	if got.String() != want.String() {
		t.Errorf("next() = %v, want %v", got, want)
	}

	got = lb.next()
	want = mockURL3

	if got.String() != want.String() {
		t.Errorf("next() = %v, want %v", got, want)
	}

	got = lb.next()
	want = mockURL1

	if got.String() != want.String() {
		t.Errorf("next() = %v, want %v", got, want)
	}
}
