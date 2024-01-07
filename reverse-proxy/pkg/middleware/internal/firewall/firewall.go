package firewall

import (
	"fmt"
	"github.com/tasuku43/go-learn-projects-hub/waf/pkg/middleware/internal/firewall/rules"
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, rule := range newRules() {
			_, err := rule.Apply(r)
			if err != nil {
				message := err.Error()
				fmt.Println("Access denied: " + message)
				http.Error(w, "Access denied.", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

type Rule interface {
	Apply(r *http.Request) (bool, error)
}

type Rules []Rule

func newRules() []Rule {
	return []Rule{
		&rules.PathAllowRule{
			AllowedPaths: []string{
				"/tasks",
			},
		},
		&rules.MethodBlockRule{
			AllowedMethods: []string{
				http.MethodPost,
				http.MethodPut,
				http.MethodDelete,
				http.MethodGet,
			},
		},
		&rules.HeaderRule{
			RequiredHeaders: map[string]string{
				"X-Required-Header": "True",
			},
		},
	}
}
