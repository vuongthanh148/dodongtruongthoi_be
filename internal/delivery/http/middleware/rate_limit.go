package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/vuongthanh148/dodongtruongthoi_be/pkg/response"
)

// RateLimiter holds in-memory rate limit state keyed by client IP.
type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new in-memory rate limiter.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// RateLimit returns a middleware factory that applies per-IP rate limiting.
func (rl *RateLimiter) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getClientIP(r)
			if !rl.allowRequest(ip) {
				response.Error(w, http.StatusTooManyRequests, "too many requests, try again later")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// allowRequest checks if the IP has remaining quota. It removes stale entries and
// returns true if the request is allowed (count < limit), false otherwise.
func (rl *RateLimiter) allowRequest(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	times, exists := rl.requests[ip]

	// Remove stale entries outside the window.
	var validTimes []time.Time
	for _, t := range times {
		if now.Sub(t) < rl.window {
			validTimes = append(validTimes, t)
		}
	}

	// Check limit and add current request if allowed.
	if len(validTimes) < rl.limit {
		validTimes = append(validTimes, now)
		rl.requests[ip] = validTimes
		return true
	}

	// Limit exceeded; keep the valid times for future checks.
	if len(validTimes) > 0 {
		rl.requests[ip] = validTimes
	} else if exists {
		// Clean up empty entries to prevent memory leak.
		delete(rl.requests, ip)
	}
	return false
}

// getClientIP extracts the client IP from the request, checking X-Forwarded-For first
// (for proxy scenarios like Render) and falling back to RemoteAddr.
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (first entry is the original client).
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs; take the first.
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Fall back to RemoteAddr and strip the port.
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}
