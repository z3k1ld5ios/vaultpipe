package env

import "strings"

// PrefixFilter applies or strips a namespace prefix from environment variable keys.
type PrefixFilter struct {
	prefix string
}

// NewPrefixFilter creates a PrefixFilter with the given prefix.
func NewPrefixFilter(prefix string) *PrefixFilter {
	return &PrefixFilter{prefix: strings.ToUpper(prefix)}
}

// Apply adds the prefix to each key in the map, returning a new map.
func (p *PrefixFilter) Apply(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		key := k
		if p.prefix != "" {
			key = p.prefix + k
		}
		out[key] = v
	}
	return out
}

// Strip removes the prefix from keys that have it, dropping keys without the prefix.
// If prefix is empty all keys are returned unchanged.
func (p *PrefixFilter) Strip(env map[string]string) map[string]string {
	out := make(map[string]string)
	for k, v := range env {
		if p.prefix == "" {
			out[k] = v
			continue
		}
		if strings.HasPrefix(k, p.prefix) {
			out[strings.TrimPrefix(k, p.prefix)] = v
		}
	}
	return out
}

// HasPrefix reports whether the key begins with the configured prefix.
func (p *PrefixFilter) HasPrefix(key string) bool {
	if p.prefix == "" {
		return true
	}
	return strings.HasPrefix(strings.ToUpper(key), p.prefix)
}
