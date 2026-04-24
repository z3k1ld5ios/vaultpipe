// Package env provides utilities for constructing, transforming, and
// validating process environment maps sourced from HashiCorp Vault secrets.
//
// The Merger type offers fine-grained control over how secret values are
// combined with an existing base environment, supporting three strategies:
//
//   - StrategySecretWins  – secret values overwrite any conflicting base key
//     (default, matches typical "inject secrets" semantics).
//   - StrategyBaseWins    – pre-existing base values are preserved; secrets
//     only fill in missing keys.
//   - StrategyError       – any key present in both maps is treated as a
//     configuration error and a MergeConflictError is returned.
package env
