// Package output provides display formatting for vaultpipe diagnostics.
//
// It supports two output formats:
//
//   - text: human-readable tabular output suitable for terminal use
//   - json: machine-readable JSON suitable for scripting or CI pipelines
//
// Secret values are never written to output; only key names and masked
// placeholders are rendered to avoid accidental exposure in logs or
// terminal history.
package output
