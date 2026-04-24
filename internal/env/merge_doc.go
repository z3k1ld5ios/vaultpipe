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
//
// # MergeConflictError
//
// When StrategyError is active, a MergeConflictError is returned containing
// the name of the first conflicting key encountered. Callers can inspect the
// error with errors.As to retrieve the key name for diagnostic output.
//
// # Usage
//
//	merger := env.NewMerger(env.StrategySecretWins)
//	result, err := merger.Merge(baseEnv, secretEnv)
//	if err != nil {
//		var conflict *env.MergeConflictError
//		if errors.As(err, &conflict) {
//			log.Printf("conflict on key: %s", conflict.Key)
//		}
//		return err
//	}
package env
