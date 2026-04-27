// Package env provides utilities for building, transforming, and inspecting
// environment variable maps sourced from HashiCorp Vault secrets.
//
// # Tracer
//
// Tracer records fine-grained step-by-step mutations applied to an env map
// during pipeline execution. Each TraceEntry captures the step name, the
// affected key, old and new values, a severity level, and a human-readable
// message.
//
// Usage:
//
//	tr := env.NewTracer()
//	tr.Record("coerce", "PORT", "8080", "8080", env.TraceLevelInfo, "no change")
//	fmt.Println(tr.Summary())
//
// Tracer is safe for concurrent use.
package env
