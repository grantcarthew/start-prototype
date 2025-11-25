package adapters

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// RealRunner implements the Runner interface using process replacement
type RealRunner struct{}

// Exec replaces the current process with the command
// This never returns on success - the process is replaced
func (r *RealRunner) Exec(shell, command string) error {
	// Find shell binary
	shellPath, err := exec.LookPath(shell)
	if err != nil {
		return fmt.Errorf("shell not found: %w", err)
	}

	// Replace current process with shell running command
	// Args: [0] = shell name, [1] = "-c", [2] = command
	// Env: inherit current environment
	err = syscall.Exec(shellPath, []string{shell, "-c", command}, os.Environ())

	// Only reached if exec fails
	return fmt.Errorf("exec failed: %w", err)
}
