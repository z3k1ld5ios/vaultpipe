package render

import (
	"fmt"
	"strings"
)

// ErrInvalidEnvKey is returned when a key contains characters that are
// unsafe for use as POSIX environment variable names.
type ErrInvalidEnvKey struct {
	Key string
}

func (e *ErrInvalidEnvKey) Error() string {
	return fmt.Sprintf("invalid environment variable key: %q", e.Key)
}

// SanitizeKey checks that an environment variable key is non-empty,
// contains only alphanumeric characters and underscores, and does not
// start with a digit.
func SanitizeKey(key string) error {
	if key == "" {
		return &ErrInvalidEnvKey{Key: key}
	}
	for i, ch := range key {
		switch {
		case ch >= 'A' && ch <= 'Z':
		case ch >= 'a' && ch <= 'z':
		case ch == '_':
		case ch >= '0' && ch <= '9' && i > 0:
		default:
			return &ErrInvalidEnvKey{Key: key}
		}
	}
	return nil
}

// SanitizeMap validates all keys in the provided map and returns an error
// for the first invalid key found.
func SanitizeMap(env map[string]string) error {
	for k := range env {
		if err := SanitizeKey(k); err != nil {
			return err
		}
	}
	return nil
}

// MaskValue replaces all but the first two characters of a secret value
// with asterisks, suitable for debug logging.
func MaskValue(v string) string {
	if len(v) <= 2 {
		return strings.Repeat("*", len(v))
	}
	return v[:2] + strings.Repeat("*", len(v)-2)
}
