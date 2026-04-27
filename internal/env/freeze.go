package env

import (
	"errors"
	"fmt"
	"sync"
)

// ErrFrozen is returned when a mutation is attempted on a frozen environment map.
var ErrFrozen = errors.New("env: map is frozen and cannot be modified")

// FrozenMap wraps an immutable snapshot of an environment map.
// Once frozen, any attempt to set or delete a key returns ErrFrozen.
// It is safe for concurrent reads.
type FrozenMap struct {
	mu   sync.RWMutex
	data map[string]string
}

// Freeze creates a FrozenMap from the provided map.
// A deep copy is taken so the caller's original map is not aliased.
func Freeze(src map[string]string) *FrozenMap {
	copy := make(map[string]string, len(src))
	for k, v := range src {
		copy[k] = v
	}
	return &FrozenMap{data: copy}
}

// Get returns the value for key and whether it was present.
func (f *FrozenMap) Get(key string) (string, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	v, ok := f.data[key]
	return v, ok
}

// Set always returns ErrFrozen; mutations are not permitted.
func (f *FrozenMap) Set(key, value string) error {
	return fmt.Errorf("%w: attempted to set key %q", ErrFrozen, key)
}

// Delete always returns ErrFrozen; mutations are not permitted.
func (f *FrozenMap) Delete(key string) error {
	return fmt.Errorf("%w: attempted to delete key %q", ErrFrozen, key)
}

// Keys returns a sorted slice of all keys in the frozen map.
func (f *FrozenMap) Keys() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	keys := make([]string, 0, len(f.data))
	for k := range f.data {
		keys = append(keys, k)
	}
	sortStrings(keys)
	return keys
}

// ToMap returns a mutable deep copy of the frozen map.
func (f *FrozenMap) ToMap() map[string]string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	out := make(map[string]string, len(f.data))
	for k, v := range f.data {
		out[k] = v
	}
	return out
}

// Len returns the number of entries in the frozen map.
func (f *FrozenMap) Len() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.data)
}
