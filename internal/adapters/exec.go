package adapters

import (
	"bytes"
	"context"
	"os/exec"
	"time"
)

// RealRunner implements the Runner interface using os/exec
type RealRunner struct{}

// Run executes a command with the given shell and returns stdout, stderr, and any error
func (r *RealRunner) Run(ctx context.Context, shell, command string, timeout time.Duration) (string, string, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Create command
	cmd := exec.CommandContext(ctx, shell, "-c", command)

	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute
	err := cmd.Run()

	return stdout.String(), stderr.String(), err
}
