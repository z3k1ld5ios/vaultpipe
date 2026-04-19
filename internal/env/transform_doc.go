// Package env provides utilities for injecting, resolving, transforming,
// and filtering environment variables sourced from Vault secrets.
//
// The Transformer type applies named string transformations (upper, lower, trim)
// to selected keys within a secret map, enabling lightweight normalization
// before secrets are injected into a child process environment.
package env
