package config_test

import (
	"os"
	"testing"

	"github.com/grantcarthew/start/internal/config"
	"github.com/grantcarthew/start/test/mocks"
)

func TestLoadGlobal(t *testing.T) {
	mockFS := mocks.NewMockFileSystem()

	// Get the actual home directory for the test
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot get home directory")
	}

	// Setup mock files
	mockFS.Files[home+"/.config/start/config.toml"] = `
[settings]
default_agent = "claude"
log_level = "normal"
`

	mockFS.Files[home+"/.config/start/agents.toml"] = `
[agents.claude]
bin = "claude"
command = "{bin} --model {model} '{prompt}'"
default_model = "sonnet"

  [agents.claude.models]
  sonnet = "claude-3-7-sonnet-20250219"
`

	loader := config.NewLoader(mockFS)

	cfg, err := loader.LoadGlobal()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Check settings
	if cfg.Settings.DefaultAgent != "claude" {
		t.Errorf("Expected default_agent 'claude', got '%s'", cfg.Settings.DefaultAgent)
	}

	if cfg.Settings.LogLevel != "normal" {
		t.Errorf("Expected log_level 'normal', got '%s'", cfg.Settings.LogLevel)
	}

	// Check agents
	if len(cfg.Agents) != 1 {
		t.Errorf("Expected 1 agent, got %d", len(cfg.Agents))
	}

	claude, ok := cfg.Agents["claude"]
	if !ok {
		t.Fatal("Expected 'claude' agent to exist")
	}

	if claude.Bin != "claude" {
		t.Errorf("Expected bin 'claude', got '%s'", claude.Bin)
	}

	if claude.DefaultModel != "sonnet" {
		t.Errorf("Expected default_model 'sonnet', got '%s'", claude.DefaultModel)
	}

	if len(claude.Models) != 1 {
		t.Errorf("Expected 1 model, got %d", len(claude.Models))
	}
}

func TestLoadLocal(t *testing.T) {
	mockFS := mocks.NewMockFileSystem()

	// Setup mock files
	mockFS.Files["/project/.start/config.toml"] = `
[settings]
log_level = "debug"
`

	mockFS.Files["/project/.start/roles.toml"] = `
[roles.code-reviewer]
file = "./ROLE.md"
description = "Project code reviewer"
`

	loader := config.NewLoader(mockFS)

	cfg, err := loader.LoadLocal("/project")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Check settings
	if cfg.Settings.LogLevel != "debug" {
		t.Errorf("Expected log_level 'debug', got '%s'", cfg.Settings.LogLevel)
	}

	// Check roles
	if len(cfg.Roles) != 1 {
		t.Errorf("Expected 1 role, got %d", len(cfg.Roles))
	}

	reviewer, ok := cfg.Roles["code-reviewer"]
	if !ok {
		t.Fatal("Expected 'code-reviewer' role to exist")
	}

	if reviewer.File != "./ROLE.md" {
		t.Errorf("Expected file './ROLE.md', got '%s'", reviewer.File)
	}
}

func TestLoadMissingFiles(t *testing.T) {
	mockFS := mocks.NewMockFileSystem()

	// Empty filesystem
	loader := config.NewLoader(mockFS)

	cfg, err := loader.LoadGlobal()

	// Should not error, just return empty config
	if err != nil {
		t.Fatalf("Expected no error for missing files, got: %v", err)
	}

	if len(cfg.Agents) != 0 {
		t.Errorf("Expected empty agents, got %d", len(cfg.Agents))
	}

	if len(cfg.Roles) != 0 {
		t.Errorf("Expected empty roles, got %d", len(cfg.Roles))
	}
}

func TestLoadInvalidTOML(t *testing.T) {
	mockFS := mocks.NewMockFileSystem()

	// Get the actual home directory for the test
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot get home directory")
	}

	// Invalid TOML
	mockFS.Files[home+"/.config/start/config.toml"] = `
[settings
this is not valid TOML
`

	loader := config.NewLoader(mockFS)

	_, err = loader.LoadGlobal()

	if err == nil {
		t.Fatal("Expected error for invalid TOML, got nil")
	}
}
