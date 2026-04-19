// Package env provides utilities for constructing and manipulating process
// environment maps used when spawning child processes.
//
// The Override type allows callers to layer explicit key-value pairs on top of
// a base environment (e.g. secrets fetched from Vault), ensuring that
// operator-supplied values always take precedence. This is useful for
// injecting one-off tunables or test values without modifying the underlying
// secret map.
package env
