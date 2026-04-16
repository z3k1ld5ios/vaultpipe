// Package middleware provides request-level controls for vault operations.
package middleware

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// RateLimiter enforces a maximum number of Vault secret reads per interval.
type RateLimiter struct {
	mu       sync.Mutex
	max      int
	interval time.Duration
	count    int
	reset    time.Time
}

// NewRateLimiter creates a RateLimiter allowing max calls per interval.
func NewRateLimiter(max int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		max:      max,
		interval: interval,
		reset:    time.Now().Add(interval),
	}
}

// Allow checks whether the next operation is within the rate limit.
// Returns an error if the limit has been exceeded.
func (r *RateLimiter) Allow(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if ctx.Err() != nil {
		return ctx.Err()
	}

	now := time.Now()
	if now.After(r.reset) {
		r.count = 0
		r.reset = now.Add(r.interval)
	}

	r.count++
	if r.count > r.max {
		return fmt.Errorf("rate limit exceeded: %d requests per %s", r.max, r.interval)
	}
	return nil
}

// Remaining returns the number of allowed calls left in the current window.
func (r *RateLimiter) Remaining() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	rem := r.max - r.count
	if rem < 0 {
		return 0
	}
	return rem
}
