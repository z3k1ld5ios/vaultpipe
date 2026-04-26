// Package env provides utilities for managing environment variable construction,
// injection, and transformation in the vaultpipe pipeline.
package env

import (
	"strings"
	"unicode"
)

// SanitizeKey normalises an arbitrary string into a valid POSIX environment
// variable name: uppercase, only [A-Z0-9_], leading digit prefixed with '_'.
//
// Rules applied in order:
//  1. Convert to uppercase.
//  2. Replace hyphens, dots, and spaces with underscores.
//  3. Replace any remaining non-alphanumeric/non-underscore character with '_'.
//  4. Prefix with '_' if the first character is a digit.
func SanitizeKey(key string) string {
	if key == "" {
		return ""
	}

	key = strings.ToUpper(key)

	var b strings.Builder
	b.Grow(len(key))

	for _, r := range key {
		switch {
		case r == '-' || r == '.' || r == ' ':
			b.WriteRune('_')
		case unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_':
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}

	result := b.String()
	if len(result) > 0 && unicode.IsDigit(rune(result[0])) {
		result = "_" + result
	}

	return result
}

// SanitizeMap applies SanitizeKey to every key in src, returning a new map.
// If two source keys normalise to the same sanitized key, the last value
// (in iteration order) wins.
func SanitizeMap(src map[string]string) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[SanitizeKey(k)] = v
	}
	return out
}
