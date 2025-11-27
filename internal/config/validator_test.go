package config_test

import (
	"strings"
	"testing"

	"github.com/grantcarthew/start/internal/config"
	"github.com/grantcarthew/start/internal/domain"
)

func TestValidateAgent(t *testing.T) {
	validator := config.NewValidator()

	// Valid agent
	validConfig := domain.Config{
		Agents: map[string]domain.Agent{
			"claude": {
				Name:    "claude",
				Bin:     "claude",
				Command: "{bin} --model {model} '{prompt}'",
				Models: map[string]string{
					"sonnet": "claude-3-7-sonnet-20250219",
				},
				DefaultModel: "sonnet",
			},
		},
	}

	err := validator.Validate(validConfig)
	if err != nil {
		t.Errorf("Expected no error for valid config, got: %v", err)
	}

	// Missing bin
	invalidConfig := domain.Config{
		Agents: map[string]domain.Agent{
			"test": {
				Name:    "test",
				Command: "{bin} --model {model} '{prompt}'",
				Models: map[string]string{
					"default": "test-model",
				},
			},
		},
	}

	err = validator.Validate(invalidConfig)
	if err == nil {
		t.Error("Expected error for missing bin")
	}

	// Missing {bin} placeholder
	invalidConfig = domain.Config{
		Agents: map[string]domain.Agent{
			"test": {
				Name:    "test",
				Bin:     "test",
				Command: "test --model {model} '{prompt}'",
				Models: map[string]string{
					"default": "test-model",
				},
			},
		},
	}

	err = validator.Validate(invalidConfig)
	if err == nil {
		t.Error("Expected error for missing {bin} placeholder")
	}
	if err != nil && !strings.Contains(err.Error(), "{bin}") {
		t.Errorf("Expected error about {bin}, got: %v", err)
	}

	// Missing {model} placeholder
	invalidConfig = domain.Config{
		Agents: map[string]domain.Agent{
			"test": {
				Name:    "test",
				Bin:     "test",
				Command: "{bin} test '{prompt}'",
				Models: map[string]string{
					"default": "test-model",
				},
			},
		},
	}

	err = validator.Validate(invalidConfig)
	if err == nil {
		t.Error("Expected error for missing {model} placeholder")
	}
	if err != nil && !strings.Contains(err.Error(), "{model}") {
		t.Errorf("Expected error about {model}, got: %v", err)
	}

	// No models
	invalidConfig = domain.Config{
		Agents: map[string]domain.Agent{
			"test": {
				Name:    "test",
				Bin:     "test",
				Command: "{bin} --model {model} '{prompt}'",
				Models:  map[string]string{},
			},
		},
	}

	err = validator.Validate(invalidConfig)
	if err == nil {
		t.Error("Expected error for no models")
	}
}

func TestValidateRole(t *testing.T) {
	validator := config.NewValidator()

	// Valid role with file
	validConfig := domain.Config{
		Roles: map[string]domain.Role{
			"code-reviewer": {
				Name: "code-reviewer",
				File: "~/.config/start/roles/reviewer.md",
			},
		},
	}

	err := validator.Validate(validConfig)
	if err != nil {
		t.Errorf("Expected no error for valid role, got: %v", err)
	}

	// Valid role with prompt
	validConfig = domain.Config{
		Roles: map[string]domain.Role{
			"simple": {
				Name:   "simple",
				Prompt: "You are a code reviewer",
			},
		},
	}

	err = validator.Validate(validConfig)
	if err != nil {
		t.Errorf("Expected no error for valid role with prompt, got: %v", err)
	}

	// Invalid role - no UTD fields
	invalidConfig := domain.Config{
		Roles: map[string]domain.Role{
			"invalid": {
				Name: "invalid",
			},
		},
	}

	err = validator.Validate(invalidConfig)
	if err == nil {
		t.Error("Expected error for role without UTD fields")
	}
	if err != nil && !strings.Contains(err.Error(), "UTD") {
		t.Errorf("Expected error about UTD pattern, got: %v", err)
	}

	// Invalid role name
	invalidConfig = domain.Config{
		Roles: map[string]domain.Role{
			"Invalid_Name": {
				Name:   "Invalid_Name",
				Prompt: "test",
			},
		},
	}

	err = validator.Validate(invalidConfig)
	if err == nil {
		t.Error("Expected error for invalid role name")
	}
}

func TestValidateTask(t *testing.T) {
	validator := config.NewValidator()

	// Valid task
	validConfig := domain.Config{
		Agents: map[string]domain.Agent{
			"claude": {
				Name:    "claude",
				Bin:     "claude",
				Command: "{bin} --model {model} '{prompt}'",
				Models: map[string]string{
					"sonnet": "test",
				},
			},
		},
		Roles: map[string]domain.Role{
			"reviewer": {
				Name:   "reviewer",
				Prompt: "test",
			},
		},
		Tasks: map[string]domain.Task{
			"review": {
				Name:   "review",
				Alias:  "cr",
				Agent:  "claude",
				Role:   "reviewer",
				Prompt: "Review this code",
			},
		},
	}

	err := validator.Validate(validConfig)
	if err != nil {
		t.Errorf("Expected no error for valid task, got: %v", err)
	}

	// Task references non-existent agent
	invalidConfig := domain.Config{
		Tasks: map[string]domain.Task{
			"test": {
				Name:   "test",
				Agent:  "nonexistent",
				Prompt: "test",
			},
		},
	}

	err = validator.Validate(invalidConfig)
	if err == nil {
		t.Error("Expected error for non-existent agent")
	}
	if err != nil && !strings.Contains(err.Error(), "nonexistent") {
		t.Errorf("Expected error about nonexistent agent, got: %v", err)
	}

	// Task references non-existent role
	invalidConfig = domain.Config{
		Tasks: map[string]domain.Task{
			"test": {
				Name:   "test",
				Role:   "nonexistent",
				Prompt: "test",
			},
		},
	}

	err = validator.Validate(invalidConfig)
	if err == nil {
		t.Error("Expected error for non-existent role")
	}
}

func TestValidateSettings(t *testing.T) {
	validator := config.NewValidator()

	// Valid settings
	validConfig := domain.Config{
		Settings: domain.Settings{
			DefaultAgent: "claude",
			DefaultRole:  "reviewer",
			LogLevel:     "debug",
		},
		Agents: map[string]domain.Agent{
			"claude": {
				Name:    "claude",
				Bin:     "claude",
				Command: "{bin} --model {model} '{prompt}'",
				Models: map[string]string{
					"sonnet": "test",
				},
			},
		},
		Roles: map[string]domain.Role{
			"reviewer": {
				Name:   "reviewer",
				Prompt: "test",
			},
		},
	}

	err := validator.Validate(validConfig)
	if err != nil {
		t.Errorf("Expected no error for valid settings, got: %v", err)
	}

	// Invalid default_agent
	invalidConfig := domain.Config{
		Settings: domain.Settings{
			DefaultAgent: "nonexistent",
		},
	}

	err = validator.Validate(invalidConfig)
	if err == nil {
		t.Error("Expected error for non-existent default_agent")
	}

	// Invalid log_level
	invalidConfig = domain.Config{
		Settings: domain.Settings{
			LogLevel: "invalid",
		},
	}

	err = validator.Validate(invalidConfig)
	if err == nil {
		t.Error("Expected error for invalid log_level")
	}
}

func TestValidateContext(t *testing.T) {
	validator := config.NewValidator()

	// Valid context with file
	validConfig := domain.Config{
		Contexts: map[string]domain.Context{
			"environment": {
				Name: "environment",
				File: "ENVIRONMENT.md",
			},
		},
	}

	err := validator.Validate(validConfig)
	if err != nil {
		t.Errorf("Expected no error for valid context with file, got: %v", err)
	}

	// Valid context with command
	validConfig = domain.Config{
		Contexts: map[string]domain.Context{
			"git-status": {
				Name:    "git-status",
				Command: "git status",
			},
		},
	}

	err = validator.Validate(validConfig)
	if err != nil {
		t.Errorf("Expected no error for valid context with command, got: %v", err)
	}

	// Valid context with prompt
	validConfig = domain.Config{
		Contexts: map[string]domain.Context{
			"project-info": {
				Name:   "project-info",
				Prompt: "Project information",
			},
		},
	}

	err = validator.Validate(validConfig)
	if err != nil {
		t.Errorf("Expected no error for valid context with prompt, got: %v", err)
	}

	// Valid context name with hyphens
	validConfig = domain.Config{
		Contexts: map[string]domain.Context{
			"multi-word-name": {
				Name: "multi-word-name",
				File: "file.md",
			},
		},
	}

	err = validator.Validate(validConfig)
	if err != nil {
		t.Errorf("Expected no error for valid context name with hyphens, got: %v", err)
	}

	// Invalid context name - uppercase
	invalidConfig := domain.Config{
		Contexts: map[string]domain.Context{
			"InvalidName": {
				Name: "InvalidName",
				File: "file.md",
			},
		},
	}

	err = validator.Validate(invalidConfig)
	if err == nil {
		t.Error("Expected error for uppercase context name")
	} else if !containsValidationError(err, "must be lowercase") {
		t.Errorf("Expected error about lowercase, got: %v", err)
	}

	// Invalid context name - spaces
	invalidConfig = domain.Config{
		Contexts: map[string]domain.Context{
			"invalid name": {
				Name: "invalid name",
				File: "file.md",
			},
		},
	}

	err = validator.Validate(invalidConfig)
	if err == nil {
		t.Error("Expected error for context name with spaces")
	}

	// Invalid context name - underscores
	invalidConfig = domain.Config{
		Contexts: map[string]domain.Context{
			"invalid_name": {
				Name: "invalid_name",
				File: "file.md",
			},
		},
	}

	err = validator.Validate(invalidConfig)
	if err == nil {
		t.Error("Expected error for context name with underscores")
	}

	// Invalid context - empty UTD pattern
	invalidConfig = domain.Config{
		Contexts: map[string]domain.Context{
			"empty-context": {
				Name: "empty-context",
				// No file, command, or prompt
			},
		},
	}

	err = validator.Validate(invalidConfig)
	if err == nil {
		t.Error("Expected error for context with empty UTD pattern")
	} else if !containsValidationError(err, "UTD pattern") {
		t.Errorf("Expected error about UTD pattern, got: %v", err)
	}

	// Valid context - combination of file and prompt
	validConfig = domain.Config{
		Contexts: map[string]domain.Context{
			"combined": {
				Name:   "combined",
				File:   "file.md",
				Prompt: "Additional: {file_contents}",
			},
		},
	}

	err = validator.Validate(validConfig)
	if err != nil {
		t.Errorf("Expected no error for context with file and prompt, got: %v", err)
	}

	// Valid context - all three fields
	validConfig = domain.Config{
		Contexts: map[string]domain.Context{
			"all-fields": {
				Name:    "all-fields",
				File:    "file.md",
				Command: "echo test",
				Prompt:  "File: {file_contents}, Command: {command_output}",
			},
		},
	}

	err = validator.Validate(validConfig)
	if err != nil {
		t.Errorf("Expected no error for context with all UTD fields, got: %v", err)
	}
}

// Helper function to check if validation error contains a substring
func containsValidationError(err error, substr string) bool {
	if err == nil {
		return false
	}
	return len(err.Error()) >= len(substr) &&
		(err.Error() == substr || indexOfSubstr(err.Error(), substr) >= 0)
}

func indexOfSubstr(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
