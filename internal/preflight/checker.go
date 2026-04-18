// Package preflight validates runtime requirements before secrets are fetched.
package preflight

import (
	"errors"
	"fmt"
	"strings"
)

// Check represents a single preflight requirement.
type Check struct {
	Name string
	Run  func() error
}

// Result holds the outcome of a single check.
type Result struct {
	Name    string
	Passed  bool
	Message string
}

// Runner executes a set of preflight checks.
type Runner struct {
	checks []Check
}

// NewRunner returns a Runner with the provided checks.
func NewRunner(checks ...Check) *Runner {
	return &Runner{checks: checks}
}

// Run executes all checks and returns results.
// It returns an error if any check fails.
func (r *Runner) Run() ([]Result, error) {
	results := make([]Result, 0, len(r.checks))
	var failed []string

	for _, c := range r.checks {
		err := c.Run()
		if err != nil {
			results = append(results, Result{Name: c.Name, Passed: false, Message: err.Error()})
			failed = append(failed, c.Name)
		} else {
			results = append(results, Result{Name: c.Name, Passed: true, Message: "ok"})
		}
	}

	if len(failed) > 0 {
		return results, fmt.Errorf("preflight failed: %s", strings.Join(failed, ", "))
	}
	return results, nil
}

// RequireEnv returns a Check that fails if the named env var is empty.
func RequireEnv(key string, env func(string) string) Check {
	return Check{
		Name: fmt.Sprintf("env:%s", key),
		Run: func() error {
			if env(key) == "" {
				return fmt.Errorf("required environment variable %q is not set", key)
			}
			return nil
		},
	}
}

// RequireNonEmptySecrets returns a Check that fails if the secrets map is empty.
func RequireNonEmptySecrets(secrets map[string]string) Check {
	return Check{
		Name: "secrets:non-empty",
		Run: func() error {
			if len(secrets) == 0 {
				return errors.New("no secrets were loaded")
			}
			return nil
		},
	}
}
