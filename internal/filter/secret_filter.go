// Package filter provides key-based filtering of secret maps before injection.
package filter

import "strings"

// Config holds filter configuration.
type Config struct {
	// AllowKeys, if non-empty, only passes keys in this list.
	AllowKeys []string
	// DenyKeys excludes keys matching these prefixes or exact names.
	DenyKeys []string
}

// Filter applies allow/deny rules to a secret map and returns a filtered copy.
func Filter(secrets map[string]string, cfg Config) map[string]string {
	allowSet := toSet(cfg.AllowKeys)
	result := make(map[string]string, len(secrets))

	for k, v := range secrets {
		if len(allowSet) > 0 {
			if _, ok := allowSet[k]; !ok {
				continue
			}
		}
		if matchesAny(k, cfg.DenyKeys) {
			continue
		}
		result[k] = v
	}
	return result
}

func matchesAny(key string, patterns []string) bool {
	for _, p := range patterns {
		if strings.EqualFold(key, p) || strings.HasPrefix(strings.ToUpper(key), strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

func toSet(keys []string) map[string]struct{} {
	s := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		s[k] = struct{}{}
	}
	return s
}
