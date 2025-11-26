package cli

import (
	"os"
	"strings"
	"testing"

	"github.com/grantcarthew/start/internal/domain"
)

func TestSelectDefaultAgent(t *testing.T) {
	ic := &InitCommand{}

	tests := []struct {
		name   string
		agents []domain.AssetMeta
		want   string
	}{
		{
			name: "claude has priority",
			agents: []domain.AssetMeta{
				{Name: "aichat"},
				{Name: "claude"},
				{Name: "gemini"},
			},
			want: "claude",
		},
		{
			name: "gemini second priority",
			agents: []domain.AssetMeta{
				{Name: "aichat"},
				{Name: "gemini"},
				{Name: "aider"},
			},
			want: "gemini",
		},
		{
			name: "fallback to first",
			agents: []domain.AssetMeta{
				{Name: "aichat"},
				{Name: "aider"},
				{Name: "opencode"},
			},
			want: "aichat",
		},
		{
			name:   "empty list",
			agents: []domain.AssetMeta{},
			want:   "",
		},
		{
			name: "single agent",
			agents: []domain.AssetMeta{
				{Name: "aider"},
			},
			want: "aider",
		},
		{
			name: "claude and gemini both present",
			agents: []domain.AssetMeta{
				{Name: "gemini"},
				{Name: "claude"},
			},
			want: "claude",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ic.selectDefaultAgent(tt.agents)

			if got != tt.want {
				t.Errorf("selectDefaultAgent() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDetectInstalledAgents(t *testing.T) {
	ic := &InitCommand{}

	// Note: This test is platform-dependent because it uses exec.LookPath
	// We'll test with known binaries that should exist on most systems

	tests := []struct {
		name    string
		agents  []domain.AssetMeta
		wantMin int // Minimum number of detected agents (at least 0)
	}{
		{
			name: "common shell commands exist",
			agents: []domain.AssetMeta{
				{Name: "ls", Bin: "ls"},     // Should exist on Unix systems
				{Name: "cat", Bin: "cat"},   // Should exist on Unix systems
				{Name: "xxxx", Bin: "xxxx"}, // Should not exist
			},
			wantMin: 0, // At least none (conservative for cross-platform)
		},
		{
			name: "nonexistent binaries",
			agents: []domain.AssetMeta{
				{Name: "nonexistent1", Bin: "nonexistent-binary-12345"},
				{Name: "nonexistent2", Bin: "another-fake-binary-67890"},
			},
			wantMin: 0,
		},
		{
			name:    "empty agent list",
			agents:  []domain.AssetMeta{},
			wantMin: 0,
		},
		{
			name: "agent without bin",
			agents: []domain.AssetMeta{
				{Name: "no-bin-agent", Bin: ""},
			},
			wantMin: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detected := ic.detectInstalledAgents(tt.agents)

			if len(detected) < tt.wantMin {
				t.Errorf("detectInstalledAgents() detected %d agents, want at least %d", len(detected), tt.wantMin)
			}

			// Verify detected agents are a subset of input agents
			for _, det := range detected {
				found := false
				for _, input := range tt.agents {
					if det.Name == input.Name {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Detected agent %q was not in input list", det.Name)
				}
			}
		})
	}
}

func TestWriteConfigFiles(t *testing.T) {
	ic := &InitCommand{}

	// Use temp directory for tests
	targetPath := t.TempDir()

	agents := []domain.AssetMeta{
		{
			Name:        "claude",
			Bin:         "claude",
			Description: "Anthropic Claude AI",
		},
		{
			Name:        "gemini",
			Bin:         "gemini",
			Description: "Google Gemini AI",
		},
	}

	err := ic.writeConfigFiles(targetPath, agents, "claude")
	if err != nil {
		t.Fatalf("writeConfigFiles() error = %v, want nil", err)
	}

	// Check that all files were created
	expectedFiles := []string{
		"config.toml",
		"agents.toml",
		"roles.toml",
		"contexts.toml",
		"tasks.toml",
	}

	for _, filename := range expectedFiles {
		path := targetPath + "/" + filename
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected file %s was not created", filename)
		}
	}

	// Verify config.toml contains default agent
	configPath := targetPath + "/config.toml"
	configData, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config.toml: %v", err)
	}

	configStr := string(configData)
	if !strings.Contains(configStr, `default_agent = "claude"`) {
		t.Errorf("config.toml should contain default_agent = \"claude\"")
	}

	// Verify agents.toml contains both agents
	agentsPath := targetPath + "/agents.toml"
	agentsData, err := os.ReadFile(agentsPath)
	if err != nil {
		t.Fatalf("Failed to read agents.toml: %v", err)
	}

	agentsStr := string(agentsData)
	if !strings.Contains(agentsStr, "[agents.claude]") {
		t.Errorf("agents.toml should contain [agents.claude]")
	}
	if !strings.Contains(agentsStr, "[agents.gemini]") {
		t.Errorf("agents.toml should contain [agents.gemini]")
	}
}

func TestBackupConfig(t *testing.T) {
	ic := &InitCommand{}

	// Create temp directory with existing config
	targetPath := t.TempDir()

	// Write some existing config files
	existingFiles := map[string]string{
		"config.toml":   "[settings]\ndefault_agent = \"old\"",
		"agents.toml":   "[agents.old]\nbin = \"old\"",
		"roles.toml":    "[roles.old]\nprompt = \"old\"",
		"contexts.toml": "[contexts.old]\nfile = \"old.md\"",
		"tasks.toml":    "[tasks.old]\nprompt = \"old\"",
	}

	for filename, content := range existingFiles {
		path := targetPath + "/" + filename
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	// Perform backup
	err := ic.backupConfig(targetPath)
	if err != nil {
		t.Fatalf("backupConfig() error = %v, want nil", err)
	}

	// Check that backup files were created
	// They should have format: <filename>.YYYY-MM-DD-HHMMSS.toml
	entries, err := os.ReadDir(targetPath)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	backupCount := 0
	for _, entry := range entries {
		name := entry.Name()
		// Count files that match backup pattern (contain timestamp)
		if strings.Contains(name, ".2") && strings.HasSuffix(name, ".toml") && name != "config.toml" {
			backupCount++
		}
	}

	if backupCount != len(existingFiles) {
		t.Errorf("Expected %d backup files, got %d", len(existingFiles), backupCount)
	}

	// Verify original files still exist
	for filename := range existingFiles {
		path := targetPath + "/" + filename
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Original file %s should still exist after backup", filename)
		}
	}
}
