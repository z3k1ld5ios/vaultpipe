// Package redact provides string redaction utilities for vaultpipe.
//
// It is used to ensure that secret values do not appear in log lines,
// error messages, or diagnostic output. A Redactor is seeded with known
// secret values and can scrub them from arbitrary strings at runtime.
package redact
