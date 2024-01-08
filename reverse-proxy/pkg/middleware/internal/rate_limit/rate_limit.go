package rate_limit

import (
	"golang.org/x/time/rate"
	"log/slog"
	"net/http"
)

type RateLimiter struct {
	limiter *rate.Limiter
	logger  *slog.Logger
}

func (r *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, rec *http.Request) {
		if !r.limiter.Allow() {
			r.logger.Info("Too many requests")
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, rec)
	})
}

func NewRateLimiter(r *rate.Limiter, logger *slog.Logger) *RateLimiter {
	return &RateLimiter{
		limiter: r,
		logger:  logger,
	}
}
