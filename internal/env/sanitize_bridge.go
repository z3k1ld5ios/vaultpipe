package env

import (
	"strings"
	"unicode"
)

// SanitizeKey converts a string into a valid environment variable key.
// It uppercases the string and replaces any non-alphanumeric character with '_'.
func SanitizeKey(key string) string {
	key = strings.ToUpper(key)
	var b strings.Builder
	for _, r := range key {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else {
			b.WriteRune('_')
		}
	}
	return b.String()
}
