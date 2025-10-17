package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	requests map[string]int
	mu        sync.Mutex
	resetTime time.Time
	limit     int
	window    time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string]int),
		resetTime: time.Now().Add(window),
		limit:     limit,
		window:    window,
	}
}

func (r *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		r.mu.Lock()
		defer r.mu.Unlock()

		// Normalize remote address to IP only (strip port) when possible
		ip := req.RemoteAddr
		if host, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
			ip = host
		}

		// Reset if window expired
		if time.Now().After(r.resetTime) {
			r.requests = make(map[string]int)
			r.resetTime = time.Now().Add(r.window)
		}

		r.requests[ip]++

		if r.requests[ip] > r.limit {
			http.Error(w, "Too many requests, please try again later.", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, req)
	})
}

// RateLimiter is a convenience wrapper so router code can call authMiddleware.RateLimiter(limit)
// It returns a middleware that enforces `limit` requests per 1 minute window.
func RateLimiter(limit int) func(http.Handler) http.Handler {
	rl := NewRateLimiter(limit, time.Minute)
	return rl.Middleware
}
