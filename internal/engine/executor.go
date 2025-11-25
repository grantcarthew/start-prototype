package engine

import (
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
// This function replaces the current process and never returns on success
func (e *Executor) Execute(agent domain.Agent, model, prompt, shell string) error {
	// Prepare placeholder values
	values := map[string]string{
		"bin":    agent.Bin,
		"model":  model,
		"prompt": prompt,
	}

	// Resolve placeholders in command template
	command := e.resolver.Resolve(agent.Command, values)

	// Replace process with agent (never returns on success)
	return e.runner.Exec(shell, command)
}
