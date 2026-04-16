// Package snapshot provides point-in-time captures of resolved secret maps.
// Snapshots can be compared for equality or diffed to identify which keys
// changed between two fetch cycles, enabling conditional process restarts
// or audit logging of secret rotation events.
package snapshot
