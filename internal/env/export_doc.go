// Package env provides utilities for managing environment variables within
// vaultpipe's secret injection pipeline.
//
// # Exporter
//
// Exporter serialises a resolved env map into one of three text formats:
//
//   - FormatShell  — POSIX shell export statements (export KEY="VALUE")
//   - FormatDotenv — .env file syntax (KEY=VALUE)
//   - FormatJSON   — compact JSON object ({"KEY":"VALUE"})
//
// Keys are always emitted in lexicographic order so output is deterministic.
//
// Example:
//
//	exporter := env.NewExporter(env.FormatDotenv, true)
//	out, err := exporter.Export(secrets)
package env
