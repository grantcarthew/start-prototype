package config

import (
	"testing"

	"github.com/grantcarthew/start/internal/domain"
	"github.com/grantcarthew/start/test/mocks"
)

func TestTOMLHelper_ReadWriteAgents(t *testing.T) {
	fs := mocks.NewMockFileSystem()
	helper := NewTOMLHelper(fs)

	// Test writing agents
	agents := map[string]domain.Agent{
		"claude": {
			Name:        "claude",
			Bin:         "claude",
			Description: "Claude AI",
			Command:     "{bin} {prompt}",
			Models: map[string]string{
				"sonnet": "claude-sonnet-4-20250929",
			},
			DefaultModel: "sonnet",
		},
	}

	dir := "/test"
	err := helper.WriteAgentsFile(dir, agents)
	if err != nil {
		t.Fatalf("WriteAgentsFile failed: %v", err)
	}

	// Verify file was created
	path := dir + "/agents.toml"
	if !fs.Exists(path) {
		t.Fatal("agents.toml was not created")
	}

	// Test reading agents back
	readAgents, err := helper.ReadAgentsFile(dir)
	if err != nil {
		t.Fatalf("ReadAgentsFile failed: %v", err)
	}

	if len(readAgents) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(readAgents))
	}

	claude, ok := readAgents["claude"]
	if !ok {
		t.Fatal("claude agent not found")
	}

	if claude.Bin != "claude" {
		t.Errorf("expected bin=claude, got %s", claude.Bin)
	}

	if claude.Description != "Claude AI" {
		t.Errorf("expected description=Claude AI, got %s", claude.Description)
	}

	if len(claude.Models) != 1 {
		t.Errorf("expected 1 model, got %d", len(claude.Models))
	}
}

func TestTOMLHelper_ReadAgentsFile_Empty(t *testing.T) {
	fs := mocks.NewMockFileSystem()
	helper := NewTOMLHelper(fs)

	// Read from non-existent file should return empty map
	agents, err := helper.ReadAgentsFile("/nonexistent")
	if err != nil {
		t.Fatalf("expected no error for non-existent file, got: %v", err)
	}

	if len(agents) != 0 {
		t.Errorf("expected empty map, got %d agents", len(agents))
	}
}

func TestTOMLHelper_ReadWriteSettings(t *testing.T) {
	fs := mocks.NewMockFileSystem()
	helper := NewTOMLHelper(fs)

	settings := domain.Settings{
		DefaultAgent: "claude",
		DefaultRole:  "default",
		Shell:        "bash",
	}

	dir := "/test"
	err := helper.WriteSettingsFile(dir, settings)
	if err != nil {
		t.Fatalf("WriteSettingsFile failed: %v", err)
	}

	// Read back
	readSettings, err := helper.ReadSettingsFile(dir)
	if err != nil {
		t.Fatalf("ReadSettingsFile failed: %v", err)
	}

	if readSettings.DefaultAgent != "claude" {
		t.Errorf("expected DefaultAgent=claude, got %s", readSettings.DefaultAgent)
	}

	if readSettings.DefaultRole != "default" {
		t.Errorf("expected DefaultRole=default, got %s", readSettings.DefaultRole)
	}
}

func TestTOMLHelper_ReadWriteRoles(t *testing.T) {
	fs := mocks.NewMockFileSystem()
	helper := NewTOMLHelper(fs)

	// Test writing roles
	roles := map[string]domain.Role{
		"code-reviewer": {
			Name:        "code-reviewer",
			Description: "Code review expert",
			File:        "~/roles/code-reviewer.md",
			Prompt:      "{file_contents}\n\nFocus on security.",
		},
		"inline": {
			Name:   "inline",
			Prompt: "You are a helpful assistant.",
		},
	}

	dir := "/test"
	err := helper.WriteRolesFile(dir, roles)
	if err != nil {
		t.Fatalf("WriteRolesFile failed: %v", err)
	}

	// Verify file was created
	path := dir + "/roles.toml"
	if !fs.Exists(path) {
		t.Fatal("roles.toml was not created")
	}

	// Test reading roles back
	readRoles, err := helper.ReadRolesFile(dir)
	if err != nil {
		t.Fatalf("ReadRolesFile failed: %v", err)
	}

	if len(readRoles) != 2 {
		t.Fatalf("expected 2 roles, got %d", len(readRoles))
	}

	reviewer, ok := readRoles["code-reviewer"]
	if !ok {
		t.Fatal("code-reviewer role not found")
	}

	if reviewer.Description != "Code review expert" {
		t.Errorf("expected description=Code review expert, got %s", reviewer.Description)
	}

	if reviewer.File != "~/roles/code-reviewer.md" {
		t.Errorf("expected file=~/roles/code-reviewer.md, got %s", reviewer.File)
	}

	inline, ok := readRoles["inline"]
	if !ok {
		t.Fatal("inline role not found")
	}

	if inline.Prompt != "You are a helpful assistant." {
		t.Errorf("expected prompt=You are a helpful assistant., got %s", inline.Prompt)
	}
}

func TestTOMLHelper_ReadRolesFile_Empty(t *testing.T) {
	fs := mocks.NewMockFileSystem()
	helper := NewTOMLHelper(fs)

	// Read from non-existent file should return empty map
	roles, err := helper.ReadRolesFile("/nonexistent")
	if err != nil {
		t.Fatalf("expected no error for non-existent file, got: %v", err)
	}

	if len(roles) != 0 {
		t.Errorf("expected empty map, got %d roles", len(roles))
	}
}
