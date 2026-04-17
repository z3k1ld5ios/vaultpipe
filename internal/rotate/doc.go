// Package rotate implements secret rotation detection for vaultpipe.
//
// It polls a Vault secret path at a configurable interval and invokes
// registered callbacks whenever the secret values change, enabling
// processes to react to credential rotation without restarting.
package rotate
