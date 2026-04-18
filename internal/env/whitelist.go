package env

import (
	"fmt"
	"strings"
)

// Whitelist controls which environment variable names are permitted
// to be injected into a process environment.
type Whitelist struct {
	allowed map[string]struct{}
	prefixes []string
}

// NewWhitelist creates a Whitelist from explicit keys and prefixes.
// Keys are matched case-insensitively.
func NewWhitelist(keys []string, prefixes []string) *Whitelist {
	allowed := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		allowed[strings.ToUpper(k)] = struct{}{}
	}
	norm := make([]string, len(prefixes))
	for i, p := range prefixes {
		norm[i] = strings.ToUpper(p)
	}
	return &Whitelist{allowed: allowed, prefixes: norm}
}

// Allow returns true when key is permitted by the whitelist.
// If the whitelist has no keys and no prefixes every key is allowed.
func (w *Whitelist) Allow(key string) bool {
	if len(w.allowed) == 0 && len(w.prefixes) == 0 {
		return true
	}
	upper := strings.ToUpper(key)
	if _, ok := w.allowed[upper]; ok {
		return true
	}
	for _, p := range w.prefixes {
		if strings.HasPrefix(upper, p) {
			return true
		}
	}
	return false
}

// Filter returns only the entries from env whose key is allowed.
func (w *Whitelist) Filter(env []string) []string {
	out := make([]string, 0, len(env))
	for _, entry := range env {
		key, _, _ := strings.Cut(entry, "=")
		if w.Allow(key) {
			out = append(out, entry)
		}
	}
	return out
}

// Validate checks that every key in secrets is allowed, returning an
// error listing any rejected keys.
func (w *Whitelist) Validate(secrets map[string]string) error {
	var rejected []string
	for k := range secrets {
		if !w.Allow(k) {
			rejected = append(rejected, k)
		}
	}
	if len(rejected) > 0 {
		return fmt.Errorf("whitelist rejected keys: %s", strings.Join(rejected, ", "))
	}
	return nil
}
