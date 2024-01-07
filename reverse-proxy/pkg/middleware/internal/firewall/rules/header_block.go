package rules

import (
	"fmt"
	"net/http"
)

type HeaderRule struct {
	RequiredHeaders map[string]string
}

func (rule *HeaderRule) Apply(r *http.Request) (bool, error) {
	for key, value := range rule.RequiredHeaders {
		if r.Header.Get(key) != value {
			return false, fmt.Errorf("required header %s with value %s is missing", key, value)
		}
	}
	return true, nil
}
