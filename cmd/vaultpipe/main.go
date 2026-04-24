// Package main is the entry point for the vaultpipe CLI.
// It wires together all internal packages to fetch secrets from Vault
// and inject them into a subprocess environment without writing to disk.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/youorg/vaultpipe/internal/audit"
	"github.com/youorg/vaultpipe/internal/cache"
	"github.com/youorg/vaultpipe/internal/config"
	"github.com/youorg/vaultpipe/internal/diagnostics"
	"github.com/youorg/vaultpipe/internal/env"
	"github.com/youorg/vaultpipe/internal/output"
	"github.com/youorg/vaultpipe/internal/preflight"
	"github.com/youorg/vaultpipe/internal/process"
	"github.com/youorg/vaultpipe/internal/token"
	"github.com/youorg/vaultpipe/internal/vault"
)

var (
	cfgFile    string
	formatFlag string
	dryRun     bool
	verbose    bool
)

func main() {
	if err := rootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vaultpipe [flags] -- <command> [args...]",
		Short: "Pipe Vault secrets into a process environment",
		Long: `vaultpipe fetches secrets from HashiCorp Vault and injects them
as environment variables into the specified command without writing
any secret material to disk.`,
		SilenceUsage: true,
		RunE:         run,
	}

	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "path to config file (default: vaultpipe.yaml)")
	cmd.PersistentFlags().StringVarP(&formatFlag, "format", "f", "text", "output format: text or json")
	cmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "print resolved env vars without executing the command")
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging")

	cmd.AddCommand(versionCmd())

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Load configuration
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Resolve Vault token
	provider := token.NewProvider(cfg.Token.Env, cfg.Token.File)
	vaultToken, err := provider.Resolve()
	if err != nil {
		return fmt.Errorf("resolving vault token: %w", err)
	}

	// Build Vault client
	vaultClient, err := vault.NewClient(cfg.Vault.Address, vaultToken)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	// Run preflight checks
	checker := diagnostics.NewChecker(vaultClient)
	preflightRunner := preflight.NewRunner(
		preflight.ReachableVault(checker),
		preflight.TokenPresent(vaultToken),
		preflight.CommandNonEmpty(args),
	)
	if err := preflightRunner.Run(ctx); err != nil {
		return fmt.Errorf("preflight check failed: %w", err)
	}

	// Set up audit logger
	auditLog, err := audit.NewLogger(os.Stderr, verbose)
	if err != nil {
		return fmt.Errorf("creating audit logger: %w", err)
	}

	// Set up secret cache
	secretCache := cache.New(5 * time.Minute)

	// Fetch secrets for each configured path
	merged := make(map[string]string)
	for _, s := range cfg.Secrets {
		if cached, ok := secretCache.Get(s.Path); ok {
			for k, v := range cached {
				merged[k] = v
			}
			continue
		}

		secrets, err := vault.ReadSecretKV2(ctx, vaultClient, cfg.Vault.Mount, s.Path)
		if err != nil {
			return fmt.Errorf("reading secret %q: %w", s.Path, err)
		}

		secretCache.Set(s.Path, secrets)
		auditLog.SecretRead(s.Path, len(secrets))

		for k, v := range secrets {
			merged[k] = v
		}
	}

	// Build the injected environment
	injector := env.NewInjector(os.Environ())
	injectedEnv := injector.Merge(merged)

	// Output formatter for dry-run or verbose preview
	fmt_ := output.NewFormatter(formatFlag)

	if dryRun {
		fmt_.PrintEnvPreview(injectedEnv)
		return nil
	}

	if verbose {
		fmt_.PrintSecretKeys(merged)
	}

	// Launch the subprocess
	runner := process.NewRunner(args[0], args[1:], injectedEnv)
	auditLog.ProcessStart(args)

	if err := runner.Run(ctx); err != nil {
		return fmt.Errorf("process exited: %w", err)
	}

	return nil
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print vaultpipe version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("vaultpipe v0.1.0")
		},
	}
}
