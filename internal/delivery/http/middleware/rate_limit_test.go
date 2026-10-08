package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter_AllowRequest(t *testing.T) {
	tests := []struct {
		name      string
		limit     int
		window    time.Duration
		requests  int
		expectAllow bool
	}{
		{
			name:         "under limit",
			limit:        5,
			window:       time.Minute,
			requests:     3,
			expectAllow: true,
		},
		{
			name:         "at limit",
			limit:        5,
			window:       time.Minute,
			requests:     5,
			expectAllow: true,
		},
		{
			name:         "exceeds limit",
			limit:        5,
			window:       time.Minute,
			requests:     6,
			expectAllow: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rl := NewRateLimiter(tt.limit, tt.window)
			ip := "192.168.1.1"

			var lastAllowed bool
			for i := 0; i < tt.requests; i++ {
				lastAllowed = rl.allowRequest(ip)
			}

			if lastAllowed != tt.expectAllow {
				t.Errorf("expected %v, got %v after %d requests with limit %d",
					tt.expectAllow, lastAllowed, tt.requests, tt.limit)
			}
		})
	}
}

func TestRateLimiter_ResetAfterWindow(t *testing.T) {
	limit := 2
	window := 100 * time.Millisecond
	rl := NewRateLimiter(limit, window)
	ip := "192.168.1.1"

	// Make limit requests (should all pass).
	for i := 0; i < limit; i++ {
		if !rl.allowRequest(ip) {
			t.Errorf("request %d should be allowed", i+1)
		}
	}

	// Next request should be blocked.
	if rl.allowRequest(ip) {
		t.Error("request beyond limit should be blocked")
	}

	// Wait for window to expire.
	time.Sleep(window + 10*time.Millisecond)

	// Next request should be allowed after window reset.
	if !rl.allowRequest(ip) {
		t.Error("request should be allowed after window expires")
	}
}

func TestGetClientIP_RemoteAddr(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.1.1:8080"

	ip := getClientIP(req)
	if ip != "192.168.1.1" {
		t.Errorf("expected 192.168.1.1, got %s", ip)
	}
}

func TestGetClientIP_XForwardedFor(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "10.0.0.1:8080"
	req.Header.Set("X-Forwarded-For", "192.168.1.1, 10.0.0.1")

	ip := getClientIP(req)
	if ip != "192.168.1.1" {
		t.Errorf("expected 192.168.1.1 from X-Forwarded-For, got %s", ip)
	}
}

func TestMiddleware_Integration(t *testing.T) {
	limit := 2
	rl := NewRateLimiter(limit, time.Minute)
	middleware := rl.Middleware()

	// Test handler that always returns 200.
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware(nextHandler)

	// First two requests should succeed.
	for i := 0; i < limit; i++ {
		req := httptest.NewRequest("POST", "/api/v1/orders", nil)
		req.RemoteAddr = "192.168.1.1:8080"
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d: expected 200, got %d", i+1, w.Code)
		}
	}

	// Third request should be rate-limited (429).
	req := httptest.NewRequest("POST", "/api/v1/orders", nil)
	req.RemoteAddr = "192.168.1.1:8080"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429 for rate-limited request, got %d", w.Code)
	}
}

func TestRateLimiter_DifferentIPs(t *testing.T) {
	rl := NewRateLimiter(1, time.Minute)

	// IP 1 gets one request.
	if !rl.allowRequest("192.168.1.1") {
		t.Error("first request from IP1 should pass")
	}

	// IP 1 is now rate-limited.
	if rl.allowRequest("192.168.1.1") {
		t.Error("second request from IP1 should be blocked")
	}

	// IP 2 should still have quota (separate counter).
	if !rl.allowRequest("192.168.1.2") {
		t.Error("first request from IP2 should pass")
	}

	// IP 2 is now rate-limited.
	if rl.allowRequest("192.168.1.2") {
		t.Error("second request from IP2 should be blocked")
	}
}
