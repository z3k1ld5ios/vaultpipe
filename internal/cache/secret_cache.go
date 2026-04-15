// Package cache provides a simple in-memory TTL cache for Vault secrets
// to reduce redundant API calls during process execution.
package cache

import (
	"sync"
	"time"
)

// Entry holds a cached secret map and its expiry time.
type Entry struct {
	Secrets   map[string]string
	ExpiresAt time.Time
}

// SecretCache is a thread-safe in-memory cache for secret maps keyed by path.
type SecretCache struct {
	mu      sync.RWMutex
	entries map[string]Entry
	ttl     time.Duration
}

// New creates a new SecretCache with the given TTL duration.
// A TTL of zero disables caching (every Get returns a miss).
func New(ttl time.Duration) *SecretCache {
	return &SecretCache{
		entries: make(map[string]Entry),
		ttl:     ttl,
	}
}

// Set stores the secret map for the given path.
// If TTL is zero the entry is not stored.
func (c *SecretCache) Set(path string, secrets map[string]string) {
	if c.ttl == 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[path] = Entry{
		Secrets:   secrets,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// Get returns the cached secret map for path and whether it was a valid hit.
func (c *SecretCache) Get(path string) (map[string]string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[path]
	if !ok {
		return nil, false
	}
	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}
	return entry.Secrets, true
}

// Invalidate removes the cache entry for the given path.
func (c *SecretCache) Invalidate(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, path)
}

// Flush removes all entries from the cache.
func (c *SecretCache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]Entry)
}
