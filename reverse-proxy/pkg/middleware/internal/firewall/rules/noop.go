package rules

import (
	"fmt"
	"net/http"
)

type NoopRule struct{}

func (r *NoopRule) Apply(req *http.Request) bool {
	fmt.Println("NoopRule")
	return true
}
