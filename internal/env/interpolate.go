package env

import (
	"fmt"
	"regexp"
	"strings"
)

// interpolatePattern matches ${KEY} and ${KEY:-default} style references.
var interpolatePattern = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)(?::-(.*?))?\}`)

// Interpolator replaces ${KEY} and ${KEY:-default} placeholders in string
// values using a provided secrets map as the source of truth, with an optional
// fallback to OS environment variables.
type Interpolator struct {
	secrets map[string]string
	useFallback bool
}

// NewInterpolator creates an Interpolator backed by the given secrets map.
// If useFallback is true, unresolved keys are looked up via os.Getenv.
func NewInterpolator(secrets map[string]string, useFallback bool) *Interpolator {
	return &Interpolator{secrets: secrets, useFallback: useFallback}
}

// Interpolate replaces all ${KEY} and ${KEY:-default} placeholders in s.
// Returns an error if a placeholder cannot be resolved and has no default.
func (i *Interpolator) Interpolate(s string) (string, error) {
	var resolveErr error
	result := interpolatePattern.ReplaceAllStringFunc(s, func(match string) string {
		parts := interpolatePattern.FindStringSubmatch(match)
		if len(parts) < 3 {
			return match
		}
		key, def := parts[1], parts[2]
		if v, ok := i.secrets[key]; ok {
			return v
		}
		if i.useFallback {
			if v, ok := lookupEnv(key); ok {
				return v
			}
		}
		if def != "" || strings.Contains(match, ":-") {
			return def
		}
		resolveErr = fmt.Errorf("interpolate: unresolved placeholder %q", key)
		return match
	})
	if resolveErr != nil {
		return "", resolveErr
	}
	return result, nil
}

// InterpolateMap applies Interpolate to every value in m, returning a new map.
func (i *Interpolator) InterpolateMap(m map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(m))
	for k, v := range m {
		resolved, err := i.Interpolate(v)
		if err != nil {
			return nil, fmt.Errorf("interpolate map key %q: %w", k, err)
		}
		out[k] = resolved
	}
	return out, nil
}
