package server

import (
	"sync"
	"time"
)

// RateLimiterConfig holds configuration for the rate limiter.
type RateLimiterConfig struct {
	RequestsPerSecond int           // Maximum requests per second per client
	BurstSize         int           // Maximum burst size (tokens in bucket)
	CleanupInterval   time.Duration // How often to clean up expired entries
}

// DefaultRateLimiterConfig returns sensible defaults for rate limiting.
func DefaultRateLimiterConfig() *RateLimiterConfig {
	return &RateLimiterConfig{
		RequestsPerSecond: 10,               // 10 requests per second per client
		BurstSize:         20,               // Allow bursts up to 20 requests
		CleanupInterval:   5 * time.Minute,  // Clean up every 5 minutes
	}
}

// TokenBucket represents a token bucket for a single client.
type TokenBucket struct {
	tokens    float64   // Current number of tokens
	capacity  float64   // Maximum capacity
	refillRate float64  // Tokens per second
	lastRefill time.Time // Last time bucket was refilled
	mutex     sync.Mutex // Protects bucket state
}

// NewTokenBucket creates a new token bucket with the specified parameters.
func NewTokenBucket(capacity float64, refillRate float64) *TokenBucket {
	return &TokenBucket{
		tokens:     capacity,
		capacity:   capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow checks if a request should be allowed and consumes a token if so.
func (tb *TokenBucket) Allow() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	
	// Refill tokens based on elapsed time
	tb.tokens += elapsed * tb.refillRate
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}
	tb.lastRefill = now

	// Check if we have tokens available
	if tb.tokens >= 1.0 {
		tb.tokens -= 1.0
		return true
	}

	return false
}

// RateLimiter manages rate limiting for multiple clients using token buckets.
type RateLimiter struct {
	config  *RateLimiterConfig
	buckets map[string]*TokenBucket
	mutex   sync.RWMutex
	stopCh  chan struct{}
}

// NewRateLimiter creates a new rate limiter with the specified configuration.
func NewRateLimiter(config *RateLimiterConfig) *RateLimiter {
	if config == nil {
		config = DefaultRateLimiterConfig()
	}

	rl := &RateLimiter{
		config:  config,
		buckets: make(map[string]*TokenBucket),
		stopCh:  make(chan struct{}),
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// Allow checks if a request from the given client should be allowed.
func (rl *RateLimiter) Allow(clientIP string) bool {
	if clientIP == "" {
		clientIP = "unknown"
	}

	rl.mutex.RLock()
	bucket, exists := rl.buckets[clientIP]
	rl.mutex.RUnlock()

	if !exists {
		// Create new bucket for this client
		rl.mutex.Lock()
		// Double-check pattern to avoid race condition
		bucket, exists = rl.buckets[clientIP]
		if !exists {
			bucket = NewTokenBucket(
				float64(rl.config.BurstSize),
				float64(rl.config.RequestsPerSecond),
			)
			rl.buckets[clientIP] = bucket
		}
		rl.mutex.Unlock()
	}

	return bucket.Allow()
}

// cleanup periodically removes inactive buckets to prevent memory leaks.
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.performCleanup()
		case <-rl.stopCh:
			return
		}
	}
}

// performCleanup removes buckets that haven't been used recently.
func (rl *RateLimiter) performCleanup() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	cutoff := 10 * time.Minute // Remove buckets inactive for 10+ minutes

	for clientIP, bucket := range rl.buckets {
		bucket.mutex.Lock()
		inactive := now.Sub(bucket.lastRefill) > cutoff
		bucket.mutex.Unlock()

		if inactive {
			delete(rl.buckets, clientIP)
		}
	}
}

// Stop shuts down the rate limiter and stops background cleanup.
func (rl *RateLimiter) Stop() {
	close(rl.stopCh)
}

// GetStats returns current rate limiter statistics.
func (rl *RateLimiter) GetStats() RateLimiterStats {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()

	return RateLimiterStats{
		ActiveClients:     len(rl.buckets),
		RequestsPerSecond: rl.config.RequestsPerSecond,
		BurstSize:         rl.config.BurstSize,
	}
}

// RateLimiterStats contains statistics about the rate limiter.
type RateLimiterStats struct {
	ActiveClients     int // Number of clients with active buckets
	RequestsPerSecond int // Configured requests per second limit
	BurstSize         int // Configured burst size limit
}