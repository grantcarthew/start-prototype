package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/grantcarthew/start/internal/domain"
)

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors is a collection of validation errors
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return "no validation errors"
	}

	var sb strings.Builder
	sb.WriteString("configuration validation failed:\n")
	for _, err := range e {
		sb.WriteString(fmt.Sprintf("  - %s\n", err.Error()))
	}
	return sb.String()
}

// Validator validates configuration
type Validator struct{}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{}
}

// Validate validates a merged configuration
func (v *Validator) Validate(cfg domain.Config) error {
	var errors ValidationErrors

	// Validate agents
	for name, agent := range cfg.Agents {
		errors = append(errors, v.validateAgent(name, agent)...)
	}

	// Validate roles
	for name, role := range cfg.Roles {
		errors = append(errors, v.validateRole(name, role)...)
	}

	// Validate contexts
	for name, ctx := range cfg.Contexts {
		errors = append(errors, v.validateContext(name, ctx)...)
	}

	// Validate tasks
	for name, task := range cfg.Tasks {
		errors = append(errors, v.validateTask(name, task, cfg)...)
	}

	// Validate settings references
	errors = append(errors, v.validateSettings(cfg)...)

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// validateAgent validates an agent configuration
func (v *Validator) validateAgent(name string, agent domain.Agent) ValidationErrors {
	var errors ValidationErrors

	// Agent name pattern: lowercase alphanumeric with hyphens
	namePattern := regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
	if !namePattern.MatchString(name) {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("agents.%s", name),
			Message: "agent name must be lowercase alphanumeric with hyphens (e.g., 'my-agent')",
		})
	}

	// Bin is required
	if agent.Bin == "" {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("agents.%s.bin", name),
			Message: "bin field is required",
		})
	}

	// Command is required
	if agent.Command == "" {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("agents.%s.command", name),
			Message: "command field is required",
		})
	}

	// Command must contain {bin} placeholder
	if !strings.Contains(agent.Command, "{bin}") {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("agents.%s.command", name),
			Message: "command must contain {bin} placeholder",
		})
	}

	// Command must contain {model} placeholder
	if !strings.Contains(agent.Command, "{model}") {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("agents.%s.command", name),
			Message: "command must contain {model} placeholder",
		})
	}

	// Models table must exist and have at least one model
	if len(agent.Models) == 0 {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("agents.%s.models", name),
			Message: "agent requires at least one model definition",
		})
	}

	// Validate model names
	for modelName := range agent.Models {
		if !namePattern.MatchString(modelName) {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("agents.%s.models.%s", name, modelName),
				Message: "model name must be lowercase alphanumeric with hyphens",
			})
		}
	}

	// If default_model is set, it must exist in models
	if agent.DefaultModel != "" {
		if _, ok := agent.Models[agent.DefaultModel]; !ok {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("agents.%s.default_model", name),
				Message: fmt.Sprintf("default_model '%s' not found in models table", agent.DefaultModel),
			})
		}
	}

	return errors
}

// validateRole validates a role configuration
func (v *Validator) validateRole(name string, role domain.Role) ValidationErrors {
	var errors ValidationErrors

	// Role name pattern
	namePattern := regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
	if !namePattern.MatchString(name) {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("roles.%s", name),
			Message: "role name must be lowercase alphanumeric with hyphens",
		})
	}

	// UTD pattern: at least one of file, command, or prompt must be present
	if role.File == "" && role.Command == "" && role.Prompt == "" {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("roles.%s", name),
			Message: "at least one of 'file', 'command', or 'prompt' must be specified (UTD pattern)",
		})
	}

	return errors
}

// validateContext validates a context configuration
func (v *Validator) validateContext(name string, ctx domain.Context) ValidationErrors {
	var errors ValidationErrors

	// Context name pattern
	namePattern := regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
	if !namePattern.MatchString(name) {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("contexts.%s", name),
			Message: "context name must be lowercase alphanumeric with hyphens",
		})
	}

	// UTD pattern: at least one of file, command, or prompt must be present
	if ctx.File == "" && ctx.Command == "" && ctx.Prompt == "" {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("contexts.%s", name),
			Message: "at least one of 'file', 'command', or 'prompt' must be specified (UTD pattern)",
		})
	}

	return errors
}

// validateTask validates a task configuration
func (v *Validator) validateTask(name string, task domain.Task, cfg domain.Config) ValidationErrors {
	var errors ValidationErrors

	// Task name pattern
	namePattern := regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
	if !namePattern.MatchString(name) {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("tasks.%s", name),
			Message: "task name must be lowercase alphanumeric with hyphens",
		})
	}

	// Alias pattern (if present)
	if task.Alias != "" && !namePattern.MatchString(task.Alias) {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("tasks.%s.alias", name),
			Message: "alias must be lowercase alphanumeric with hyphens",
		})
	}

	// UTD pattern: at least one of file, command, or prompt must be present
	if task.File == "" && task.Command == "" && task.Prompt == "" {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("tasks.%s", name),
			Message: "at least one of 'file', 'command', or 'prompt' must be specified (UTD pattern)",
		})
	}

	// If agent is specified, it must exist
	if task.Agent != "" {
		if _, ok := cfg.Agents[task.Agent]; !ok {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("tasks.%s.agent", name),
				Message: fmt.Sprintf("agent '%s' not found in configuration", task.Agent),
			})
		}
	}

	// If role is specified, it must exist
	if task.Role != "" {
		if _, ok := cfg.Roles[task.Role]; !ok {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("tasks.%s.role", name),
				Message: fmt.Sprintf("role '%s' not found in configuration", task.Role),
			})
		}
	}

	return errors
}

// validateSettings validates settings and their references
func (v *Validator) validateSettings(cfg domain.Config) ValidationErrors {
	var errors ValidationErrors

	// If default_agent is set, it must exist
	if cfg.Settings.DefaultAgent != "" {
		if _, ok := cfg.Agents[cfg.Settings.DefaultAgent]; !ok {
			errors = append(errors, ValidationError{
				Field:   "settings.default_agent",
				Message: fmt.Sprintf("default_agent '%s' not found in agents", cfg.Settings.DefaultAgent),
			})
		}
	}

	// If default_role is set, it must exist
	if cfg.Settings.DefaultRole != "" {
		if _, ok := cfg.Roles[cfg.Settings.DefaultRole]; !ok {
			errors = append(errors, ValidationError{
				Field:   "settings.default_role",
				Message: fmt.Sprintf("default_role '%s' not found in roles", cfg.Settings.DefaultRole),
			})
		}
	}

	// Validate log_level values
	if cfg.Settings.LogLevel != "" {
		validLevels := map[string]bool{
			"quiet":   true,
			"normal":  true,
			"verbose": true,
			"debug":   true,
		}
		if !validLevels[cfg.Settings.LogLevel] {
			errors = append(errors, ValidationError{
				Field:   "settings.log_level",
				Message: "log_level must be one of: quiet, normal, verbose, debug",
			})
		}
	}

	return errors
}
