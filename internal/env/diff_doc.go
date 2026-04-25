// Package env provides utilities for constructing, transforming, and
// inspecting environment variable maps sourced from HashiCorp Vault secrets.
//
// The diff sub-feature compares two snapshots of an env map and surfaces
// granular key-level changes (added, removed, updated). It is intentionally
// decoupled from the snapshot package so that callers can use plain
// map[string]string values without importing additional types.
//
// Typical usage:
//
//	changes := env.Diff(previousEnv, currentEnv)
//	for _, c := range changes {
//		fmt.Printf("%s %s\n", c.Type, c.Key)
//	}
package env
