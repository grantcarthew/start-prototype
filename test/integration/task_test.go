package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/grantcarthew/start/test/assert"
)

// TestPhase5_TaskExecution tests task execution with instructions
func TestPhase5_TaskExecution(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Ensure binaries exist
	ensureSmithBinary(t)
	ensureStartBinary(t)

	// Create temp directory for test config
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", "start")
	err := os.MkdirAll(configDir, 0755)
	assert.NoError(t, err)

	// Create output directory for smith
	outputDir := filepath.Join(tempDir, "smith-output")
	err = os.MkdirAll(outputDir, 0755)
	assert.NoError(t, err)

	// Get absolute path to smith binary
	smithBinPath, err := filepath.Abs(filepath.Join("..", "..", "bin", "smith"))
	assert.NoError(t, err)

	// Write test config
	agentsConfig := `[agents.smith]
bin = "` + smithBinPath + `"
command = "{bin} --model {model} '{prompt}'"
default_model = "test"

  [agents.smith.models]
  test = "test-model-123"
`

	settingsConfig := `[settings]
default_agent = "smith"
default_role = "test-role"
`

	rolesConfig := `[roles.test-role]
description = "Test role"
prompt = "You are a test assistant."
`

	tasksConfig := `[tasks.help]
alias = "h"
description = "Get help"
prompt = "Help me with: {instructions}"

[tasks.code-review]
alias = "cr"
description = "Review code"
prompt = "Review this code. Focus: {instructions}"
`

	err = os.WriteFile(filepath.Join(configDir, "agents.toml"), []byte(agentsConfig), 0644)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(configDir, "config.toml"), []byte(settingsConfig), 0644)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(configDir, "roles.toml"), []byte(rolesConfig), 0644)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(configDir, "tasks.toml"), []byte(tasksConfig), 0644)
	assert.NoError(t, err)

	// Set environment
	env := []string{
		"HOME=" + tempDir,
		"SMITH_OUTPUT_DIR=" + outputDir,
		"PATH=" + os.Getenv("PATH"),
	}

	// Run: start task help "debugging code"
	startPath := filepath.Join("..", "..", "bin", "start")
	cmd := exec.Command(startPath, "task", "help", "debugging code")
	cmd.Env = env
	output, err := cmd.CombinedOutput()
	t.Logf("Command output: %s", string(output))
	if err != nil {
		t.Logf("Command error: %v", err)
	}
	assert.NoError(t, err)

	// Verify smith captured the prompt
	promptPath := filepath.Join(outputDir, "prompt.md")
	promptData, err := os.ReadFile(promptPath)
	assert.NoError(t, err)

	prompt := string(promptData)

	// Check that instructions were replaced
	assert.Contains(t, prompt, "debugging code")
	assert.Contains(t, prompt, "Help me with: debugging code")

	// Clean up for next test
	os.Remove(promptPath)
}

// TestPhase5_TaskWithNoInstructions tests task with empty instructions (should use "None")
func TestPhase5_TaskWithNoInstructions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Ensure binaries exist
	ensureSmithBinary(t)
	ensureStartBinary(t)

	// Create temp directory for test config
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", "start")
	err := os.MkdirAll(configDir, 0755)
	assert.NoError(t, err)

	// Create output directory for smith
	outputDir := filepath.Join(tempDir, "smith-output")
	err = os.MkdirAll(outputDir, 0755)
	assert.NoError(t, err)

	// Get absolute path to smith binary
	smithBinPath, err := filepath.Abs(filepath.Join("..", "..", "bin", "smith"))
	assert.NoError(t, err)

	// Write test config
	agentsConfig := `[agents.smith]
bin = "` + smithBinPath + `"
command = "{bin} --model {model} '{prompt}'"
default_model = "test"

  [agents.smith.models]
  test = "test-model-123"
`

	settingsConfig := `[settings]
default_agent = "smith"
default_role = "test-role"
`

	rolesConfig := `[roles.test-role]
description = "Test role"
prompt = "You are a test assistant."
`

	tasksConfig := `[tasks.help]
prompt = "Help requested. Instructions: {instructions}"
`

	err = os.WriteFile(filepath.Join(configDir, "agents.toml"), []byte(agentsConfig), 0644)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(configDir, "config.toml"), []byte(settingsConfig), 0644)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(configDir, "roles.toml"), []byte(rolesConfig), 0644)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(configDir, "tasks.toml"), []byte(tasksConfig), 0644)
	assert.NoError(t, err)

	// Set environment
	env := []string{
		"HOME=" + tempDir,
		"SMITH_OUTPUT_DIR=" + outputDir,
		"PATH=" + os.Getenv("PATH"),
	}

	// Run: start task help (no instructions)
	startPath := filepath.Join("..", "..", "bin", "start")
	cmd := exec.Command(startPath, "task", "help")
	cmd.Env = env
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Command output: %s", string(output))
	}
	assert.NoError(t, err)

	// Verify smith captured the prompt
	promptPath := filepath.Join(outputDir, "prompt.md")
	promptData, err := os.ReadFile(promptPath)
	assert.NoError(t, err)

	prompt := string(promptData)

	// Check that instructions default to "None"
	assert.Contains(t, prompt, "Instructions: None")
}

// TestPhase5_TaskByAlias tests task execution by alias
func TestPhase5_TaskByAlias(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Ensure binaries exist
	ensureSmithBinary(t)
	ensureStartBinary(t)

	// Create temp directory for test config
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", "start")
	err := os.MkdirAll(configDir, 0755)
	assert.NoError(t, err)

	// Create output directory for smith
	outputDir := filepath.Join(tempDir, "smith-output")
	err = os.MkdirAll(outputDir, 0755)
	assert.NoError(t, err)

	// Get absolute path to smith binary
	smithBinPath, err := filepath.Abs(filepath.Join("..", "..", "bin", "smith"))
	assert.NoError(t, err)

	// Write test config
	agentsConfig := `[agents.smith]
bin = "` + smithBinPath + `"
command = "{bin} --model {model} '{prompt}'"
default_model = "test"

  [agents.smith.models]
  test = "test-model-123"
`

	settingsConfig := `[settings]
default_agent = "smith"
default_role = "test-role"
`

	rolesConfig := `[roles.test-role]
description = "Test role"
prompt = "You are a test assistant."
`

	tasksConfig := `[tasks.code-review]
alias = "cr"
description = "Code review task"
prompt = "Review code: {instructions}"
`

	err = os.WriteFile(filepath.Join(configDir, "agents.toml"), []byte(agentsConfig), 0644)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(configDir, "config.toml"), []byte(settingsConfig), 0644)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(configDir, "roles.toml"), []byte(rolesConfig), 0644)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(configDir, "tasks.toml"), []byte(tasksConfig), 0644)
	assert.NoError(t, err)

	// Set environment
	env := []string{
		"HOME=" + tempDir,
		"SMITH_OUTPUT_DIR=" + outputDir,
		"PATH=" + os.Getenv("PATH"),
	}

	// Run: start task cr "security issues"
	startPath := filepath.Join("..", "..", "bin", "start")
	cmd := exec.Command(startPath, "task", "cr", "security issues")
	cmd.Env = env
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Command output: %s", string(output))
	}
	assert.NoError(t, err)

	// Verify smith captured the prompt
	promptPath := filepath.Join(outputDir, "prompt.md")
	promptData, err := os.ReadFile(promptPath)
	assert.NoError(t, err)

	prompt := string(promptData)

	// Check that task was executed with instructions
	assert.Contains(t, prompt, "Review code: security issues")
}

// TestPhase5_TaskNotFound tests error handling for missing tasks
func TestPhase5_TaskNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Ensure binaries exist
	ensureSmithBinary(t)
	ensureStartBinary(t)

	// Create temp directory for test config
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", "start")
	err := os.MkdirAll(configDir, 0755)
	assert.NoError(t, err)

	// Get absolute path to smith binary
	smithBinPath, err := filepath.Abs(filepath.Join("..", "..", "bin", "smith"))
	assert.NoError(t, err)

	// Write test config (no tasks)
	agentsConfig := `[agents.smith]
bin = "` + smithBinPath + `"
command = "{bin} --model {model} '{prompt}'"
default_model = "test"

  [agents.smith.models]
  test = "test-model-123"
`

	settingsConfig := `[settings]
default_agent = "smith"
default_role = "test-role"
`

	rolesConfig := `[roles.test-role]
description = "Test role"
prompt = "You are a test assistant."
`

	tasksConfig := `[tasks.existing-task]
prompt = "Existing task"
`

	err = os.WriteFile(filepath.Join(configDir, "agents.toml"), []byte(agentsConfig), 0644)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(configDir, "config.toml"), []byte(settingsConfig), 0644)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(configDir, "roles.toml"), []byte(rolesConfig), 0644)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(configDir, "tasks.toml"), []byte(tasksConfig), 0644)
	assert.NoError(t, err)

	// Set environment
	env := []string{
		"HOME=" + tempDir,
		"PATH=" + os.Getenv("PATH"),
	}

	// Run: start task nonexistent
	startPath := filepath.Join("..", "..", "bin", "start")
	cmd := exec.Command(startPath, "task", "nonexistent")
	cmd.Env = env
	output, err := cmd.CombinedOutput()

	// Should error
	if err == nil {
		t.Errorf("Expected error for nonexistent task, got nil")
	}

	// Check error message
	errMsg := string(output)
	assert.Contains(t, errMsg, "not found")
}

// ensureStartBinary ensures the start binary is built
func ensureStartBinary(t *testing.T) {
	startPath := filepath.Join("..", "..", "bin", "start")
	if _, err := os.Stat(startPath); os.IsNotExist(err) {
		t.Logf("Building start binary...")
		cmd := exec.Command("go", "build", "-o", "bin/start", "cmd/start/main.go")
		cmd.Dir = filepath.Join("..", "..")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to build start: %v\nOutput: %s", err, string(output))
		}
	}
}
