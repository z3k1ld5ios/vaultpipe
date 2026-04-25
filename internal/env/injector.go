package env

import (
	"fmt"
	"strings"
)

// SecretMap holds key-value pairs resolved from Vault secrets.
type SecretMap map[string]string

// Injector merges Vault secrets into a base environment slice.
type Injector struct {
	baseEnv []string
}

// NewInjector creates an Injector seeded with the provided base environment.
// Typically called with os.Environ().
func NewInjector(baseEnv []string) *Injector {
	return &Injector{baseEnv: baseEnv}
}

// Merge combines the base environment with the provided secrets.
// Secrets take precedence over any existing keys in the base environment.
// Returns a new slice suitable for exec.Cmd.Env.
func (inj *Injector) Merge(secrets SecretMap) []string {
	// Build a map from the base env so we can override duplicates.
	result := make(map[string]string, len(inj.baseEnv)+len(secrets))

	for _, entry := range inj.baseEnv {
		key, value := splitEntry(entry)
		result[key] = value
	}

	for k, v := range secrets {
		result[k] = v
	}

	env := make([]string, 0, len(result))
	for k, v := range result {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	return env
}

// Lookup returns the value for the given key from the base environment,
// along with a boolean indicating whether the key was found.
func (inj *Injector) Lookup(key string) (string, bool) {
	for _, entry := range inj.baseEnv {
		k, v := splitEntry(entry)
		if k == key {
			return v, true
		}
	}
	return "", false
}

// splitEntry splits an environment string of the form KEY=VALUE.
// If no '=' is present the whole string is treated as the key with an empty value.
func splitEntry(entry string) (key, value string) {
	parts := strings.SplitN(entry, "=", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return parts[0], ""
}
