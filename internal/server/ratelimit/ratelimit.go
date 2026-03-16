// Package ratelimit provides rate limiting middleware for the HTTP server.
// Uses a sliding window algorithm with in-memory storage.
package ratelimit

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// maxTrackedKeys bounds the number of IPs/keys tracked to prevent memory exhaustion.
const maxTrackedKeys = 100_000

// Limiter tracks request counts per key using a sliding window.
type Limiter struct {
	mu       sync.RWMutex
	windows  map[string]*window
	limit    int           // max requests per window
	window   time.Duration // window duration
	cleanupT *time.Ticker
	stopCh   chan struct{}
}

type window struct {
	count     int
	startTime time.Time
}

// Config holds rate limiter configuration.
type Config struct {
	Enabled bool

	// Login rate limiting (strict)
	LoginLimit  int           // requests per window
	LoginWindow time.Duration // window duration

	// Public API rate limiting (moderate)
	APILimit  int
	APIWindow time.Duration

	// Submission rate limiting (comments, feedback, newsletter)
	SubmitLimit  int
	SubmitWindow time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Enabled:      true,
		LoginLimit:   5,
		LoginWindow:  15 * time.Minute,
		APILimit:     100,
		APIWindow:    time.Minute,
		SubmitLimit:  10,
		SubmitWindow: time.Minute,
	}
}

// New creates a new rate limiter.
func New(limit int, windowDuration time.Duration) *Limiter {
	l := &Limiter{
		windows: make(map[string]*window),
		limit:   limit,
		window:  windowDuration,
		stopCh:  make(chan struct{}),
	}

	// Start cleanup goroutine
	l.cleanupT = time.NewTicker(windowDuration)
	go l.cleanup()

	return l
}

// cleanup removes expired windows periodically.
func (l *Limiter) cleanup() {
	for {
		select {
		case <-l.cleanupT.C:
			l.mu.Lock()
			now := time.Now()
			for key, w := range l.windows {
				if now.Sub(w.startTime) > l.window*2 {
					delete(l.windows, key)
				}
			}
			l.mu.Unlock()
		case <-l.stopCh:
			l.cleanupT.Stop()
			return
		}
	}
}

// Stop stops the cleanup goroutine.
func (l *Limiter) Stop() {
	close(l.stopCh)
}

// evictOldest removes the oldest entries when the map exceeds maxTrackedKeys.
// Caller must hold l.mu write lock.
func (l *Limiter) evictOldest() {
	if len(l.windows) <= maxTrackedKeys {
		return
	}
	now := time.Now()
	// First pass: remove all expired entries
	for key, w := range l.windows {
		if now.Sub(w.startTime) > l.window {
			delete(l.windows, key)
		}
	}
	// If still over capacity, remove oldest entries until under limit
	for len(l.windows) > maxTrackedKeys {
		var oldestKey string
		var oldestTime time.Time
		first := true
		for key, w := range l.windows {
			if first || w.startTime.Before(oldestTime) {
				oldestKey = key
				oldestTime = w.startTime
				first = false
			}
		}
		delete(l.windows, oldestKey)
	}
}

// Allow checks if a request should be allowed for the given key.
// Returns (allowed, remaining, resetTime).
func (l *Limiter) Allow(key string) (bool, int, time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	w, exists := l.windows[key]

	if !exists || now.Sub(w.startTime) > l.window {
		// Evict before adding a new entry
		l.evictOldest()

		// New window
		l.windows[key] = &window{
			count:     1,
			startTime: now,
		}
		return true, l.limit - 1, now.Add(l.window)
	}

	// Existing window
	if w.count >= l.limit {
		remaining := 0
		resetTime := w.startTime.Add(l.window)
		return false, remaining, resetTime
	}

	w.count++
	remaining := l.limit - w.count
	resetTime := w.startTime.Add(l.window)
	return true, remaining, resetTime
}

// Reset clears the rate limit for a key (e.g., after successful login).
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.windows, key)
}

// KeyFunc extracts a rate limit key from a request.
type KeyFunc func(c *gin.Context) string

// IPKeyFunc uses the client IP as the rate limit key.
func IPKeyFunc(c *gin.Context) string {
	return c.ClientIP()
}

// IPAndPathKeyFunc uses IP + path as the key (for per-endpoint limiting).
func IPAndPathKeyFunc(c *gin.Context) string {
	return c.ClientIP() + ":" + c.Request.URL.Path
}

// Middleware creates a Gin middleware that applies rate limiting.
func (l *Limiter) Middleware(keyFunc KeyFunc) gin.HandlerFunc {
	if keyFunc == nil {
		keyFunc = IPKeyFunc
	}

	return func(c *gin.Context) {
		key := keyFunc(c)
		allowed, remaining, resetTime := l.Allow(key)

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", itoa(l.limit))
		c.Header("X-RateLimit-Remaining", itoa(remaining))
		c.Header("X-RateLimit-Reset", itoa(int(resetTime.Unix())))

		if !allowed {
			retryAfter := int(time.Until(resetTime).Seconds())
			if retryAfter < 1 {
				retryAfter = 1
			}
			c.Header("Retry-After", itoa(retryAfter))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "rate limit exceeded",
				"retry_after": retryAfter,
			})
			return
		}

		c.Next()
	}
}

// itoa converts int to string without importing strconv.
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	var b [20]byte
	pos := len(b)
	for i > 0 {
		pos--
		b[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		b[pos] = '-'
	}
	return string(b[pos:])
}
