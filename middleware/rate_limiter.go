package middleware

import (
	"context"
	"fmt"
	. "goThrottle/limiter"
	"log"
	"net"
	"net/http"
)

const rateLimitExceededMsg = "you have reached the maximum number of requests or actions allowed within a certain time frame"

func RateLimiter(
	ctx context.Context,
	limiter *Limiter,
	next http.Handler,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("API_KEY")
		if apiKey != "" {
			log.Default().Printf("Limiting based on API key")
			// Limit based on API key (token)
			allowed, err := limiter.CheckLimit(ctx, fmt.Sprintf("token:%s", apiKey), TokenLimit)
			if err != nil || !allowed {
				http.Error(w, rateLimitExceededMsg, http.StatusTooManyRequests)
				return
			}
		} else {
			log.Default().Printf("Limiting based on IP address")
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				http.Error(w, "invalid remote address", http.StatusBadRequest)
				return
			}
			allowed, err := limiter.CheckLimit(ctx, fmt.Sprintf("ip:%s", ip), IPLimit)
			if err != nil || !allowed {
				http.Error(w, rateLimitExceededMsg, http.StatusTooManyRequests)
				return
			}
		}

		// If not blocked, continue processing the request
		next.ServeHTTP(w, r)
	}
}
