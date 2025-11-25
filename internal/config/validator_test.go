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
