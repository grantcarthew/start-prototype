package engine

import (
	"strings"

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

// ExecuteParams holds parameters for execution
type ExecuteParams struct {
	Agent         domain.Agent
	Model         string
	UserPrompt    string
	RoleContent   string
	RoleFilePath  string
	Contexts      []LoadedContext
	Shell         string
}

// Execute runs an agent command with the given parameters
// This function replaces the current process and never returns on success
func (e *Executor) Execute(params ExecuteParams) error {
	// Build final prompt by combining contexts, role, and user prompt
	finalPrompt := e.buildFinalPrompt(params.Contexts, params.UserPrompt)

	// Prepare placeholder values
	values := map[string]string{
		"bin":       params.Agent.Bin,
		"model":     params.Model,
		"prompt":    finalPrompt,
		"role":      params.RoleContent,
		"role_file": params.RoleFilePath,
	}

	// Resolve placeholders in command template
	command := e.resolver.Resolve(params.Agent.Command, values)

	// Replace process with agent (never returns on success)
	return e.runner.Exec(params.Shell, command)
}

// buildFinalPrompt combines contexts and user prompt
func (e *Executor) buildFinalPrompt(contexts []LoadedContext, userPrompt string) string {
	var parts []string

	// Add context documents
	for _, ctx := range contexts {
		if ctx.Content != "" {
			parts = append(parts, ctx.Content)
		}
	}

	// Add user prompt
	if userPrompt != "" {
		parts = append(parts, userPrompt)
	}

	return strings.Join(parts, "\n\n")
}
