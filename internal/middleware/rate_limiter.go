package middleware

import (
	"net"
	"net/http"
	"strings"
	"time"
	"shortify/internal/cache"
)

func RateLimitMiddleware(redis *cache.RedisClient, limit int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ip := getIP(r)

			key := "rate:" + ip

			// increment count
			count, err := redis.Increment(key)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// set expiry if first request
			if count == 1 {
				_ = redis.Expire(key, time.Minute)
			}

			// check limit
			if count > int64(limit) {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// helper to extract IP
func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		parts := strings.Split(ip, ",")
		return strings.TrimSpace(parts[0])
	}

	ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	return ip
}