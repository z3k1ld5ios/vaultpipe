// Package signal handles OS-level termination signals (SIGINT, SIGTERM) for
// vaultpipe. It provides a Handler that accepts shutdown callbacks, enabling
// clean teardown of Vault leases, cache flushes, and process cleanup before
// the process exits.
package signal
