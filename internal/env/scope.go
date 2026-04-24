package env

import "strings"

// Scope restricts which environment keys are visible or injectable
// based on an explicit allowlist of key names or prefixes.
type Scope struct {
	keys     map[string]struct{}
	prefixes []string
}

// NewScope creates a Scope from a list of allowed keys or prefix patterns.
// Entries ending with "*" are treated as prefix matches (e.g. "APP_*").
func NewScope(entries []string) *Scope {
	s := &Scope{
		keys: make(map[string]struct{}),
	}
	for _, e := range entries {
		e = strings.TrimSpace(e)
		if strings.HasSuffix(e, "*") {
			s.prefixes = append(s.prefixes, strings.TrimSuffix(e, "*"))
		} else if e != "" {
			s.keys[e] = struct{}{}
		}
	}
	return s
}

// Allows reports whether the given key is permitted by the scope.
// If the scope is empty (no keys, no prefixes), all keys are allowed.
func (s *Scope) Allows(key string) bool {
	if len(s.keys) == 0 && len(s.prefixes) == 0 {
		return true
	}
	if _, ok := s.keys[key]; ok {
		return true
	}
	for _, p := range s.prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}

// Filter returns a new map containing only entries permitted by the scope.
func (s *Scope) Filter(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if s.Allows(k) {
			out[k] = v
		}
	}
	return out
}

// Keys returns the sorted list of explicitly allowed keys (non-prefix entries).
func (s *Scope) Keys() []string {
	result := make([]string, 0, len(s.keys))
	for k := range s.keys {
		result = append(result, k)
	}
	return result
}
