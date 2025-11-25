package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/grantcarthew/start/internal/domain"
)

// Executor executes agent commands with resolved placeholders
type Executor struct {
	runner   domain.Runner
	resolver *PlaceholderResolver
}

// NewExecutor creates a new executor
func NewExecutor(runner domain.Runner, resolver *PlaceholderResolver) *Executor {
	return &Executor{
		runner:   runner,
		resolver: resolver,
	}
}

// Execute runs an agent command with the given parameters
func (e *Executor) Execute(ctx context.Context, agent domain.Agent, model, prompt string) error {
	// Prepare placeholder values
	values := map[string]string{
		"bin":    agent.Bin,
		"model":  model,
		"prompt": prompt,
	}

	// Resolve placeholders in command template
	command := e.resolver.Resolve(agent.Command, values)

	// Determine shell (default to bash)
	shell := "bash"

	// Execute command with 2 minute default timeout
	timeout := 2 * time.Minute
	stdout, stderr, err := e.runner.Run(ctx, shell, command, timeout)

	// Print output to console
	if stdout != "" {
		fmt.Print(stdout)
	}
	if stderr != "" {
		fmt.Print(stderr)
	}

	return err
}
