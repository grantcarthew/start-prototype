package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grantcarthew/start/test/assert"
)

// TestPhase2_BasicExecution tests that start can execute smith with basic placeholders
func TestPhase2_BasicExecution(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Ensure smith binary exists
	ensureSmithBinary(t)

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
`

	err = os.WriteFile(filepath.Join(configDir, "agents.toml"), []byte(agentsConfig), 0644)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(configDir, "config.toml"), []byte(settingsConfig), 0644)
	assert.NoError(t, err)

	// Set environment
	env := []string{
		"HOME=" + tempDir,
		"SMITH_OUTPUT_DIR=" + outputDir,
		"PATH=" + os.Getenv("PATH"),
	}

	// Run: start "hello world"
	startPath := filepath.Join("..", "..", "bin", "start")
	cmd := exec.Command(startPath, "hello world")
	cmd.Env = env
	_, err = cmd.CombinedOutput()
	assert.NoError(t, err)

	// Verify smith captured the correct args
	argsPath := filepath.Join(outputDir, "args.txt")
	argsData, err := os.ReadFile(argsPath)
	assert.NoError(t, err)

	argsLines := strings.Split(string(argsData), "\n")
	assert.Contains(t, argsLines[0], "smith")

	foundModel := false
	foundModelID := false
	for _, line := range argsLines {
		if line == "--model" {
			foundModel = true
		}
		if line == "test-model-123" {
			foundModelID = true
		}
	}
	assert.True(t, foundModel, "Args should contain --model flag")
	assert.True(t, foundModelID, "Args should contain resolved model ID")

	// Verify prompt contains expected content
	promptPath := filepath.Join(outputDir, "prompt.md")
	promptData, err := os.ReadFile(promptPath)
	assert.NoError(t, err)

	assert.Equal(t, "hello world", string(promptData))
}

// TestPhase2_ModelResolution tests that model names are resolved correctly
func TestPhase2_ModelResolution(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ensureSmithBinary(t)

	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", "start")
	err := os.MkdirAll(configDir, 0755)
	assert.NoError(t, err)

	outputDir := filepath.Join(tempDir, "smith-output")
	err = os.MkdirAll(outputDir, 0755)
	assert.NoError(t, err)

	// Get absolute path to smith binary
	smithBinPath, err := filepath.Abs(filepath.Join("..", "..", "bin", "smith"))
	assert.NoError(t, err)

	// Write test config with multiple models
	agentsConfig := `[agents.smith]
bin = "` + smithBinPath + `"
command = "{bin} --model {model} '{prompt}'"
default_model = "sonnet"

  [agents.smith.models]
  haiku = "haiku-123"
  sonnet = "sonnet-456"
  opus = "opus-789"
`

	settingsConfig := `[settings]
default_agent = "smith"
`

	err = os.WriteFile(filepath.Join(configDir, "agents.toml"), []byte(agentsConfig), 0644)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(configDir, "config.toml"), []byte(settingsConfig), 0644)
	assert.NoError(t, err)

	env := []string{
		"HOME=" + tempDir,
		"SMITH_OUTPUT_DIR=" + outputDir,
		"PATH=" + os.Getenv("PATH"),
	}

	// Run: start --model haiku "test prompt"
	startPath := filepath.Join("..", "..", "bin", "start")
	cmd := exec.Command(startPath, "--model", "haiku", "test prompt")
	cmd.Env = env
	_, err = cmd.CombinedOutput()
	assert.NoError(t, err)

	// Verify the resolved model
	argsPath := filepath.Join(outputDir, "args.txt")
	argsData, err := os.ReadFile(argsPath)
	assert.NoError(t, err)

	argsStr := string(argsData)
	assert.Contains(t, argsStr, "haiku-123")
	assert.NotContains(t, argsStr, "sonnet-456")
}

// Helper function to ensure smith binary exists
func ensureSmithBinary(t *testing.T) {
	t.Helper()

	// Get project root (two levels up from test/integration)
	root := filepath.Join("..", "..")

	// Check if smith binary exists
	smithPath := filepath.Join(root, "bin", "smith")
	if _, err := os.Stat(smithPath); os.IsNotExist(err) {
		t.Log("Building smith binary...")
		cmd := exec.Command("go", "build", "-o", smithPath, "./cmd/smith/")
		cmd.Dir = root
		_, err := cmd.CombinedOutput()
		assert.NoError(t, err)
	}

	// Also ensure start binary exists
	startPath := filepath.Join(root, "bin", "start")
	if _, err := os.Stat(startPath); os.IsNotExist(err) {
		t.Log("Building start binary...")
		cmd := exec.Command("go", "build", "-o", startPath, "./cmd/start/")
		cmd.Dir = root
		_, err := cmd.CombinedOutput()
		assert.NoError(t, err)
	}
}
