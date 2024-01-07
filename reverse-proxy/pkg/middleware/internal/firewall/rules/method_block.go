package rules

import (
	"fmt"
	"net/http"
)

type MethodBlockRule struct {
	AllowedMethods []string
}

func (rule *MethodBlockRule) Apply(r *http.Request) (bool, error) {
	for _, allowedMethod := range rule.AllowedMethods {
		if r.Method == allowedMethod {
			return true, nil
		}
	}
	return false, fmt.Errorf("method %s is not allowed", r.Method)
}
