// Package validate provides pre-injection validation of resolved secrets.
package validate

import (
	"errors"
	"fmt"
	"strings"
)

// Rule defines a validation rule applied to a secret map.
type Rule struct {
	// RequiredKeys lists keys that must be present and non-empty.
	RequiredKeys []string
	// ForbidEmptyValues rejects any key whose value is an empty string.
	ForbidEmptyValues bool
	// MaxValueLength, if > 0, rejects values exceeding this length.
	MaxValueLength int
}

// ValidationError accumulates all rule violations found during Validate.
type ValidationError struct {
	Violations []string
}

func (e *ValidationError) Error() string {
	return "secret validation failed:\n  " + strings.Join(e.Violations, "\n  ")
}

// Validate applies r to secrets and returns a ValidationError if any
// violations are found, or nil when all checks pass.
func Validate(secrets map[string]string, r Rule) error {
	var violations []string

	for _, key := range r.RequiredKeys {
		v, ok := secrets[key]
		if !ok {
			violations = append(violations, fmt.Sprintf("required key %q is missing", key))
			continue
		}
		if v == "" {
			violations = append(violations, fmt.Sprintf("required key %q is empty", key))
		}
	}

	if r.ForbidEmptyValues {
		for k, v := range secrets {
			if v == "" {
				violations = append(violations, fmt.Sprintf("key %q has an empty value", k))
			}
		}
	}

	if r.MaxValueLength > 0 {
		for k, v := range secrets {
			if len(v) > r.MaxValueLength {
				violations = append(violations, fmt.Sprintf("key %q exceeds max length %d", k, r.MaxValueLength))
			}
		}
	}

	if len(violations) > 0 {
		return &ValidationError{Violations: violations}
	}
	return nil
}

// MustValidate is like Validate but panics on error — useful in init paths.
func MustValidate(secrets map[string]string, r Rule) {
	if err := Validate(secrets, r); err != nil {
		panic(errors.New("MustValidate: " + err.Error()))
	}
}
