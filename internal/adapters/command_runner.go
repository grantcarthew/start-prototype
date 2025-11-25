package adapters

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

// RealCommandRunner executes commands with output capture
type RealCommandRunner struct{}

// NewRealCommandRunner creates a new command runner
func NewRealCommandRunner() *RealCommandRunner {
	return &RealCommandRunner{}
}

// Run executes a command and returns combined stdout+stderr output
func (r *RealCommandRunner) Run(shell, command string, timeoutSeconds int) (string, error) {
	// Determine shell flag
	flag := getShellFlag(shell)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	// Create command
	cmd := exec.CommandContext(ctx, shell, flag, command)

	// Capture combined output
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	// Execute
	err := cmd.Run()

	// Check for timeout
	if ctx.Err() == context.DeadlineExceeded {
		return output.String(), fmt.Errorf("command timeout after %d seconds", timeoutSeconds)
	}

	// Return output even on error (partial output may be useful)
	return output.String(), err
}

// getShellFlag returns the appropriate flag for the shell
func getShellFlag(shell string) string {
	switch shell {
	case "node", "nodejs", "bun":
		return "-e"
	case "deno":
		return "eval"
	case "ruby":
		return "-e"
	case "perl":
		return "-E"
	default:
		// bash, sh, zsh, fish, python, python2, python3, and unknown shells
		return "-c"
	}
}
