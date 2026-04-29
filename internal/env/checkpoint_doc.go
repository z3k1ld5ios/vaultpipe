// Package env provides utilities for managing environment variable maps
// throughout the vaultpipe secret injection pipeline.
//
// Checkpointer allows callers to save named snapshots of an env map at
// arbitrary points during processing and later diff any two checkpoints
// to understand what changed between pipeline stages.
//
// Example:
//
//	cp := env.NewCheckpointer()
//	cp.Save("before", rawEnv)
//	cp.Save("after", processedEnv)
//	changes, err := cp.Between("before", "after")
package env
