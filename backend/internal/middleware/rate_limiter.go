package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// ipBucket tracks request attempts and reset window for a specific client IP.
type ipBucket struct {
	count    int
	lastSeen time.Time
	resetAt  time.Time
}

// RateLimiter manages in-memory rate limiting counters per IP address.
type RateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*ipBucket
	limit   int
	window  time.Duration
	message string
}

// NewRateLimiter creates a new RateLimiter instance with automatic background cleanup.
func NewRateLimiter(limit int, window time.Duration, customMessage ...string) *RateLimiter {
	msg := "Too many requests. Please wait a moment before trying again."
	if len(customMessage) > 0 && customMessage[0] != "" {
		msg = customMessage[0]
	}

	rl := &RateLimiter{
		buckets: make(map[string]*ipBucket),
		limit:   limit,
		window:  window,
		message: msg,
	}

	// Periodically evict stale client IP records to prevent memory leak
	go rl.cleanupLoop()

	return rl
}

func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.window * 2)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, bucket := range rl.buckets {
			if now.Sub(bucket.lastSeen) > rl.window*3 {
				delete(rl.buckets, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Middleware returns a gin.HandlerFunc enforcing the rate limit.
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if ip == "" {
			ip = "unknown"
		}

		rl.mu.Lock()
		now := time.Now()
		bucket, exists := rl.buckets[ip]

		// Initialize or reset window if expired
		if !exists || now.After(bucket.resetAt) {
			rl.buckets[ip] = &ipBucket{
				count:    1,
				lastSeen: now,
				resetAt:  now.Add(rl.window),
			}
			remaining := rl.limit - 1
			resetUnix := now.Add(rl.window).Unix()
			rl.mu.Unlock()

			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
			c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetUnix))
			c.Next()
			return
		}

		bucket.lastSeen = now

		// Check if quota exceeded
		if bucket.count >= rl.limit {
			retryAfter := int(time.Until(bucket.resetAt).Seconds())
			if retryAfter <= 0 {
				retryAfter = 1
			}
			resetUnix := bucket.resetAt.Unix()
			rl.mu.Unlock()

			c.Header("Retry-After", fmt.Sprintf("%d", retryAfter))
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetUnix))

			c.JSON(http.StatusTooManyRequests, gin.H{
				"success":    false,
				"message":    rl.message,
				"retryAfter": retryAfter,
			})
			c.Abort()
			return
		}

		bucket.count++
		remaining := rl.limit - bucket.count
		resetUnix := bucket.resetAt.Unix()
		rl.mu.Unlock()

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetUnix))

		c.Next()
	}
}

// RateLimit is a convenient constructor helper for inline router use.
func RateLimit(limit int, window time.Duration, customMessage ...string) gin.HandlerFunc {
	return NewRateLimiter(limit, window, customMessage...).Middleware()
}
