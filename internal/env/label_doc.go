// Package env provides utilities for constructing, transforming, and managing
// environment variable maps sourced from HashiCorp Vault secrets.
//
// # Labeler
//
// Labeler attaches arbitrary key-value metadata (labels) to individual
// environment variable keys. Labels are useful for tracking provenance,
// classification (e.g. "source=vault", "tier=secret"), or any other
// annotation that should travel alongside an env entry without being
// injected into the child process environment.
//
// Labels are stored in-memory and are never written to disk or injected
// into the process environment.
package env
