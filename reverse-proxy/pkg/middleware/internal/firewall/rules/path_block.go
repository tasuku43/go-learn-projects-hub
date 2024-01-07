package rules

import (
	"fmt"
	"net/http"
)

type PathAllowRule struct {
	AllowedPaths []string
}

func (rule *PathAllowRule) Apply(r *http.Request) (bool, error) {
	for _, allowedPath := range rule.AllowedPaths {
		if r.URL.Path == allowedPath {
			return true, nil
		}
	}
	return false, fmt.Errorf("access to path %s is not allowed", r.URL.Path)
}
