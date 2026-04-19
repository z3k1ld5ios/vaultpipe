// Package env provides utilities for environment variable management.
package env

import (
	"fmt"
	"os"
	"strings"
)

// Expander resolves ${VAR} and $VAR references within string values,
// substituting from a merged set of secrets and OS environment variables.
type Expander struct {
	secrets map[string]string
	fallbackToOS bool
}

// NewExpander creates an Expander backed by the provided secrets map.
// If fallbackToOS is true, unresolved references fall back to os.Getenv.
func NewExpander(secrets map[string]string, fallbackToOS bool) *Expander {
	return &Expander{secrets: secrets, fallbackToOS: fallbackToOS}
}

// Expand replaces variable references in s with their resolved values.
// Returns an error if a reference cannot be resolved.
func (e *Expander) Expand(s string) (string, error) {
	var expandErr error
	result := os.Expand(s, func(key string) string {
		if v, ok := e.secrets[key]; ok {
			return v
		}
		if e.fallbackToOS {
			if v := os.Getenv(key); v != "" {
				return v
			}
		}
		expandErr = fmt.Errorf("expand: unresolved variable %q", key)
		return ""
	})
	if expandErr != nil {
		return "", expandErr
	}
	return result, nil
}

// ExpandMap applies Expand to every value in the provided map.
// Returns a new map and the first error encountered, if any.
func (e *Expander) ExpandMap(m map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(m))
	for k, v := range m {
		expanded, err := e.Expand(v)
		if err != nil {
			return nil, fmt.Errorf("expand map key %q: %w", k, err)
		}
		out[k] = expanded
	}
	return out, nil
}

// ContainsReference reports whether s contains at least one $VAR or ${VAR} reference.
func ContainsReference(s string) bool {
	return strings.Contains(s, "$")
}
