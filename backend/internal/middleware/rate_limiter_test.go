package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRateLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	// Allow 3 requests per 100ms
	r.POST("/test-login", RateLimit(3, 100*time.Millisecond, "Too many requests"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Request 1: should pass (200)
	req1 := httptest.NewRequest(http.MethodPost, "/test-login", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", w1.Code)
	}
	if rem := w1.Header().Get("X-RateLimit-Remaining"); rem != "2" {
		t.Errorf("Expected remaining 2, got %s", rem)
	}

	// Request 2: should pass (200)
	req2 := httptest.NewRequest(http.MethodPost, "/test-login", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", w2.Code)
	}
	if rem := w2.Header().Get("X-RateLimit-Remaining"); rem != "1" {
		t.Errorf("Expected remaining 1, got %s", rem)
	}

	// Request 3: should pass (200, last allowed in window)
	req3 := httptest.NewRequest(http.MethodPost, "/test-login", nil)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	if w3.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", w3.Code)
	}
	if rem := w3.Header().Get("X-RateLimit-Remaining"); rem != "0" {
		t.Errorf("Expected remaining 0, got %s", rem)
	}

	// Request 4: should be blocked with 429 Too Many Requests
	req4 := httptest.NewRequest(http.MethodPost, "/test-login", nil)
	w4 := httptest.NewRecorder()
	r.ServeHTTP(w4, req4)
	if w4.Code != http.StatusTooManyRequests {
		t.Fatalf("Expected 429 Too Many Requests, got %d", w4.Code)
	}
	if retryAfter := w4.Header().Get("Retry-After"); retryAfter == "" {
		t.Errorf("Expected Retry-After header to be set")
	}

	// Wait for window to expire
	time.Sleep(120 * time.Millisecond)

	// Request 5: should pass again after window reset (200)
	req5 := httptest.NewRequest(http.MethodPost, "/test-login", nil)
	w5 := httptest.NewRecorder()
	r.ServeHTTP(w5, req5)
	if w5.Code != http.StatusOK {
		t.Fatalf("Expected 200 after window reset, got %d", w5.Code)
	}
}
