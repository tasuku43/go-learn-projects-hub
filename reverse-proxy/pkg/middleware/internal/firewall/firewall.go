package firewall

import (
	"github.com/tasuku43/go-learn-projects-hub/waf/pkg/middleware/internal/firewall/rules"
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, rule := range newRules() {
			if !rule.Apply(r) {
				http.Error(w, "Access denied", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

type Rule interface {
	Apply(r *http.Request) bool
}

type Rules []Rule

func newRules() []Rule {
	return []Rule{
		&rules.NoopRule{},
	}
}
