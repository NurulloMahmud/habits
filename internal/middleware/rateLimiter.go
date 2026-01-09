package middleware

import (
	"net/http"

	"github.com/NurulloMahmud/habits/pkg/response"
	"golang.org/x/time/rate"
)

func (m *Middleware) RateLimit(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(2, 4)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			response.RateLimitExceeded(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
