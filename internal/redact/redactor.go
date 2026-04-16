// Package redact provides utilities for redacting sensitive secret values
// from log output, error messages, and other string representations.
package redact

import "strings"

// Redactor holds a set of known secret values and replaces them in strings.
type Redactor struct {
	secrets []string
}

// New creates a Redactor seeded with the given secret values.
func New(secrets map[string]string) *Redactor {
	vals := make([]string, 0, len(secrets))
	for _, v := range secrets {
		if v != "" {
			vals = append(vals, v)
		}
	}
	return &Redactor{secrets: vals}
}

// Redact replaces any known secret values in s with "***REDACTED***".
func (r *Redactor) Redact(s string) string {
	for _, secret := range r.secrets {
		s = strings.ReplaceAll(s, secret, "***REDACTED***")
	}
	return s
}

// Add registers additional secret values to be redacted.
func (r *Redactor) Add(values ...string) {
	for _, v := range values {
		if v != "" {
			r.secrets = append(r.secrets, v)
		}
	}
}

// RedactMap returns a copy of m with all values replaced by "***REDACTED***".
func RedactMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k := range m {
		out[k] = "***REDACTED***"
	}
	return out
}
