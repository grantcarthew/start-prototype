package config_test

import (
	"testing"

	"github.com/grantcarthew/start/internal/config"
	"github.com/grantcarthew/start/internal/domain"
)

func TestMergeSettings(t *testing.T) {
	global := domain.Config{
		Settings: domain.Settings{
			DefaultAgent:   "claude",
			DefaultRole:    "reviewer",
			LogLevel:       "normal",
			Shell:          "bash",
			CommandTimeout: 30,
		},
	}

	local := domain.Config{
		Settings: domain.Settings{
			LogLevel: "debug",
			Shell:    "zsh",
		},
	}

	result := config.Merge(global, local)

	// Local overrides
	if result.Settings.LogLevel != "debug" {
		t.Errorf("Expected log_level 'debug', got '%s'", result.Settings.LogLevel)
	}

	if result.Settings.Shell != "zsh" {
		t.Errorf("Expected shell 'zsh', got '%s'", result.Settings.Shell)
	}

	// Global preserved
	if result.Settings.DefaultAgent != "claude" {
		t.Errorf("Expected default_agent 'claude', got '%s'", result.Settings.DefaultAgent)
	}

	if result.Settings.CommandTimeout != 30 {
		t.Errorf("Expected command_timeout 30, got %d", result.Settings.CommandTimeout)
	}
}

func TestMergeAgents(t *testing.T) {
	global := domain.Config{
		Agents: map[string]domain.Agent{
			"claude": {
				Name: "claude",
				Bin:  "claude",
			},
			"gemini": {
				Name: "gemini",
				Bin:  "gemini",
			},
		},
	}

	local := domain.Config{
		Agents: map[string]domain.Agent{
			"claude": {
				Name: "claude",
				Bin:  "claude-local",
			},
			"custom": {
				Name: "custom",
				Bin:  "custom-agent",
			},
		},
	}

	result := config.Merge(global, local)

	// Should have 3 agents (claude overridden, gemini from global, custom from local)
	if len(result.Agents) != 3 {
		t.Errorf("Expected 3 agents, got %d", len(result.Agents))
	}

	// Claude should be overridden
	if result.Agents["claude"].Bin != "claude-local" {
		t.Errorf("Expected claude bin 'claude-local', got '%s'", result.Agents["claude"].Bin)
	}

	// Gemini should be from global
	if result.Agents["gemini"].Bin != "gemini" {
		t.Errorf("Expected gemini bin 'gemini', got '%s'", result.Agents["gemini"].Bin)
	}

	// Custom should be from local
	if result.Agents["custom"].Bin != "custom-agent" {
		t.Errorf("Expected custom bin 'custom-agent', got '%s'", result.Agents["custom"].Bin)
	}
}

func TestMergeRoles(t *testing.T) {
	global := domain.Config{
		Roles: map[string]domain.Role{
			"reviewer": {
				Name: "reviewer",
				File: "~/.config/start/roles/reviewer.md",
			},
		},
	}

	local := domain.Config{
		Roles: map[string]domain.Role{
			"reviewer": {
				Name: "reviewer",
				File: "./ROLE.md",
			},
		},
	}

	result := config.Merge(global, local)

	// Local should override global
	if result.Roles["reviewer"].File != "./ROLE.md" {
		t.Errorf("Expected file './ROLE.md', got '%s'", result.Roles["reviewer"].File)
	}
}

func TestMergeContexts(t *testing.T) {
	global := domain.Config{
		Contexts: map[string]domain.Context{
			"environment": {
				Name:     "environment",
				File:     "~/reference/ENVIRONMENT.md",
				Required: true,
			},
		},
		ContextOrder: []string{"environment"},
	}

	local := domain.Config{
		Contexts: map[string]domain.Context{
			"agents": {
				Name:     "agents",
				File:     "./AGENTS.md",
				Required: true,
			},
		},
		ContextOrder: []string{"agents"},
	}

	result := config.Merge(global, local)

	// Should have both contexts
	if len(result.Contexts) != 2 {
		t.Errorf("Expected 2 contexts, got %d", len(result.Contexts))
	}

	// Both should exist
	if _, ok := result.Contexts["environment"]; !ok {
		t.Error("Expected 'environment' context to exist")
	}

	if _, ok := result.Contexts["agents"]; !ok {
		t.Error("Expected 'agents' context to exist")
	}
}

func TestMergeTasks(t *testing.T) {
	global := domain.Config{
		Tasks: map[string]domain.Task{
			"review": {
				Name:  "review",
				Alias: "cr",
			},
		},
	}

	local := domain.Config{
		Tasks: map[string]domain.Task{
			"test": {
				Name: "test",
			},
		},
	}

	result := config.Merge(global, local)

	// Should have both tasks
	if len(result.Tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(result.Tasks))
	}

	if _, ok := result.Tasks["review"]; !ok {
		t.Error("Expected 'review' task to exist")
	}

	if _, ok := result.Tasks["test"]; !ok {
		t.Error("Expected 'test' task to exist")
	}
}
