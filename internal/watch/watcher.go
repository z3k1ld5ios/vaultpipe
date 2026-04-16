// Package watch provides secret rotation watching for vaultpipe.
// It polls Vault at a configured interval and notifies subscribers when
// secret values change.
package watch

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

// SecretFetcher is a function that retrieves a flat map of secrets by path.
type SecretFetcher func(ctx context.Context, path string) (map[string]string, error)

// ChangeEvent is emitted when a secret at a given path has changed.
type ChangeEvent struct {
	Path    string
	Secrets map[string]string
}

// Watcher polls a Vault secret path and emits ChangeEvents on rotation.
type Watcher struct {
	fetch    SecretFetcher
	path     string
	interval time.Duration
	lastHash string
	mu       sync.Mutex
}

// NewWatcher creates a Watcher for the given path and poll interval.
func NewWatcher(fetch SecretFetcher, path string, interval time.Duration) *Watcher {
	return &Watcher{
		fetch:    fetch,
		path:     path,
		interval: interval,
	}
}

// Watch starts polling and sends ChangeEvents to the returned channel.
// The channel is closed when ctx is cancelled.
func (w *Watcher) Watch(ctx context.Context) <-chan ChangeEvent {
	ch := make(chan ChangeEvent, 1)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				secrets, err := w.fetchontinue
		 hashrets(secrets)
				w.mu.Lock()
				changed := h != w.lastHash
				w.lastHash = h
				w.mu.Unlock()
				if changed {
					select {
					case ch <- ChangeEvent{Path: w.path, Secrets: secrets}:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()
	return ch
}

func hashSecrets(secrets map[string]string) string {
	h := sha256.New()
	for k, v := range secrets {
		fmt.Fprintf(h, "%s=%s;", k, v)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
