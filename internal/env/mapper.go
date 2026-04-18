// Package env provides utilities for environment variable management.
package env

import "fmt"

// Mapping defines how a secret key maps to an environment variable name.
type Mapping struct {
	SecretKey string
	EnvKey    string
}

// Mapper translates secret maps into env var maps using explicit key mappings.
type Mapper struct {
	mappings []Mapping
	prefix   string
}

// NewMapper creates a Mapper with an optional prefix applied to auto-generated env keys.
func NewMapper(prefix string, mappings []Mapping) *Mapper {
	return &Mapper{prefix: prefix, mappings: mappings}
}

// Apply converts a secrets map into an env var map.
// If explicit mappings are provided, only mapped keys are included.
// If no mappings are provided, all keys are included with the prefix applied.
func (m *Mapper) Apply(secrets map[string]string) (map[string]string, error) {
	result := make(map[string]string)

	if len(m.mappings) == 0 {
		for k, v := range secrets {
			envKey := SanitizeKey(m.prefix + k)
			result[envKey] = v
		}
		return result, nil
	}

	for _, mapping := range m.mappings {
		v, ok := secrets[mapping.SecretKey]
		if !ok {
			return nil, fmt.Errorf("mapper: secret key %q not found", mapping.SecretKey)
		}
		envKey := SanitizeKey(mapping.EnvKey)
		result[envKey] = v
	}

	return result, nil
}

// SanitizeKey is re-exported from render for convenience within env package.
// It delegates to the render package's SanitizeKey via an internal wrapper.
func sanitizeForMapper(key string) string {
	return SanitizeKey(key)
}
