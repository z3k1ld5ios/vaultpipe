package process

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// Runner executes a subprocess with an injected environment.
type Runner struct {
	env []string
}

// NewRunner creates a Runner with the provided environment slice.
func NewRunner(env []string) *Runner {
	return &Runner{env: env}
}

// Run executes the given command with the injected environment.
// It forwards OS signals to the child process and returns the exit code.
func (r *Runner) Run(name string, args []string) (int, error) {
	if name == "" {
		return 1, errors.New("process: command name must not be empty")
	}

	cmd := exec.Command(name, args...)
	cmd.Env = r.env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return 1, fmt.Errorf("process: failed to start command: %w", err)
	}

	// Forward signals to the child process.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range sigCh {
			_ = cmd.Process.Signal(sig)
		}
	}()

	err := cmd.Wait()
	signal.Stop(sigCh)
	close(sigCh)

	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode(), nil
		}
		return 1, fmt.Errorf("process: command error: %w", err)
	}

	return 0, nil
}
