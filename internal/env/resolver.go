package env

import (
	"fmt"
	"os"
	"strings"
)

// Resolver resolves environment variable mappings from a secrets map.
// It supports direct key references and optional default values.
type Resolver struct {
	secrets map[string]string
}

// NewResolver creates a Resolver backed by the provided secrets map.
func NewResolver(secrets map[string]string) *Resolver {
	return &Resolver{secrets: secrets}
}

// Resolve takes a mapping of ENV_VAR -> secret key (or "secret_key:default")
// and returns a flat map of ENV_VAR -> resolved value.
func (r *Resolver) Resolve(mappings map[string]string) (map[string]string, error) {
	result := make(map[string]string, len(mappings))
	for envKey, ref := range mappings {
		value, err := r.resolveRef(ref)
		if err != nil {
			return nil, fmt.Errorf("resolving %q: %w", envKey, err)
		}
		result[envKey] = value
	}
	return result, nil
}

// ResolveWithOS merges resolved mappings on top of the current OS environment.
func (r *Resolver) ResolveWithOS(mappings map[string]string) ([]string, error) {
	resolved, err := r.Resolve(mappings)
	if err != nil {
		return nil, err
	}
	base := os.Environ()
	inj := NewInjector(base)
	return inj.Merge(resolved), nil
}

// resolveRef parses "key" or "key:default" and looks up the secret.
func (r *Resolver) resolveRef(ref string) (string, error) {
	key, def, hasDefault := strings.Cut(ref, ":")
	key = strings.TrimSpace(key)
	if val, ok := r.secrets[key]; ok {
		return val, nil
	}
	if hasDefault {
		return def, nil
	}
	return "", fmt.Errorf("secret key %q not found and no default provided", key)
}
