// Package env provides utilities for constructing, transforming, and
// injecting environment variable sets derived from Vault secrets.
//
// The Flattener type converts arbitrarily nested secret maps into flat
// key=value pairs suitable for process environments. Nested keys are joined
// with a configurable separator (default "__"), making it straightforward to
// map structured Vault KV data to conventional env var naming conventions.
//
// Example:
//
//	f := env.NewFlattener("__").WithPrefix("APP")
//	flat, err := f.Flatten(map[string]any{
//	    "db": map[string]any{"host": "localhost", "port": 5432},
//	})
//	// flat => {"APP__db__host": "localhost", "APP__db__port": "5432"}
package env
