package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grantcarthew/start/test/assert"
)

// TestPhase7_InitForce tests `start init --force` automatic configuration
// This test requires network access to GitHub. It will skip if network is unavailable.
func TestPhase7_InitForce(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ensureStartBinary(t)

	// Create temp directory for test
	tempDir := t.TempDir()

	// Set environment (no existing config)
	env := []string{
		"HOME=" + tempDir,
		"PATH=" + os.Getenv("PATH"),
	}

	// Run: start init --force
	startPath := filepath.Join("..", "..", "bin", "start")
	cmd := exec.Command(startPath, "init", "--force")
	cmd.Env = env
	output, err := cmd.CombinedOutput()
	t.Logf("Command output: %s", string(output))

	outputStr := string(output)

	// Handle expected failure scenarios
	if err != nil {
		// Network errors - skip test
		if strings.Contains(outputStr, "Failed to fetch agent configurations") ||
			strings.Contains(outputStr, "failed to parse catalog index") {
			t.Skip("Skipping test - network unavailable or GitHub catalog not accessible")
		}

		// "No agents detected" is valid (exits with error but message is informative)
		if strings.Contains(outputStr, "No agents detected") {
			// This is acceptable - test passes
			return
		}

		t.Fatalf("init --force failed unexpectedly: %v\nOutput: %s", err, outputStr)
	}

	// If successful, verify config files were created
	configDir := filepath.Join(tempDir, ".config", "start")

	expectedFiles := []string{
		"config.toml",
		"agents.toml",
		"roles.toml",
		"contexts.toml",
		"tasks.toml",
	}

	for _, filename := range expectedFiles {
		filePath := filepath.Join(configDir, filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Expected file %s was not created", filename)
		}
	}

	// Verify config.toml has basic structure
	configData, err := os.ReadFile(filepath.Join(configDir, "config.toml"))
	assert.NoError(t, err)

	configStr := string(configData)
	assert.Contains(t, configStr, "[settings]")
	assert.Contains(t, configStr, "default_agent")
}

// TestPhase7_InitLocal tests `start init --local --force`
// This test requires network access to GitHub. It will skip if network is unavailable.
func TestPhase7_InitLocal(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ensureStartBinary(t)

	// Create temp directory for test
	tempDir := t.TempDir()

	// Set environment
	env := []string{
		"HOME=" + tempDir,
		"PATH=" + os.Getenv("PATH"),
	}

	// Run: start init --local --force
	// Get absolute path to start binary
	startPath, err := filepath.Abs(filepath.Join("..", "..", "bin", "start"))
	assert.NoError(t, err)

	cmd := exec.Command(startPath, "init", "--local", "--force")
	cmd.Dir = tempDir // Run in temp dir so local config goes here
	cmd.Env = env
	output, err := cmd.CombinedOutput()
	t.Logf("Command output: %s", string(output))

	outputStr := string(output)

	// Handle expected failure scenarios
	if err != nil {
		// Network errors - skip test
		if strings.Contains(outputStr, "Failed to fetch agent configurations") ||
			strings.Contains(outputStr, "failed to parse catalog index") {
			t.Skip("Skipping test - network unavailable or GitHub catalog not accessible")
		}

		// "No agents detected" is valid
		if strings.Contains(outputStr, "No agents detected") {
			return
		}

		t.Fatalf("init --local --force failed unexpectedly: %v\nOutput: %s", err, outputStr)
	}

	// Verify local config directory was created
	localConfigDir := filepath.Join(tempDir, ".start")

	expectedFiles := []string{
		"config.toml",
		"agents.toml",
		"roles.toml",
		"contexts.toml",
		"tasks.toml",
	}

	for _, filename := range expectedFiles {
		filePath := filepath.Join(localConfigDir, filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Expected local file %s was not created", filename)
		}
	}
}

// TestPhase7_InitBackup tests that init backs up existing config
// This test requires network access to GitHub. It will skip if network is unavailable.
func TestPhase7_InitBackup(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ensureStartBinary(t)

	// Create temp directory
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", "start")
	err := os.MkdirAll(configDir, 0755)
	assert.NoError(t, err)

	// Create existing config file
	existingConfig := `[settings]
default_agent = "old-agent"
`
	err = os.WriteFile(filepath.Join(configDir, "config.toml"), []byte(existingConfig), 0644)
	assert.NoError(t, err)

	// Set environment
	env := []string{
		"HOME=" + tempDir,
		"PATH=" + os.Getenv("PATH"),
	}

	// Run: start init --force (should auto-backup)
	startPath := filepath.Join("..", "..", "bin", "start")
	cmd := exec.Command(startPath, "init", "--force")
	cmd.Env = env
	output, err := cmd.CombinedOutput()
	t.Logf("Command output: %s", string(output))

	outputStr := string(output)

	// Handle expected failure scenarios
	if err != nil {
		// Network errors - skip test (but backup should still have happened)
		if strings.Contains(outputStr, "Failed to fetch agent configurations") ||
			strings.Contains(outputStr, "failed to parse catalog index") {
			// Even with network error, backup should have been created
			// Continue to verify backup was created
		} else if strings.Contains(outputStr, "No agents detected") {
			// This is valid, backup should have been created
		} else {
			t.Fatalf("init --force with existing config failed unexpectedly: %v\nOutput: %s", err, outputStr)
		}
	}

	// Check that backup was created (timestamped file)
	entries, err := os.ReadDir(configDir)
	assert.NoError(t, err)

	hasBackup := false
	for _, entry := range entries {
		name := entry.Name()
		// Look for config.YYYY-MM-DD-HHMMSS.toml pattern
		if strings.HasPrefix(name, "config.20") && strings.HasSuffix(name, ".toml") && name != "config.toml" {
			hasBackup = true

			// Verify backup contains old content
			backupData, err := os.ReadFile(filepath.Join(configDir, name))
			assert.NoError(t, err)
			assert.Contains(t, string(backupData), "old-agent")
			break
		}
	}

	if !hasBackup {
		t.Error("Expected backup file to be created")
	}
}

// TestPhase7_InitHelp tests `start init --help`
func TestPhase7_InitHelp(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ensureStartBinary(t)

	// Run: start init --help
	startPath := filepath.Join("..", "..", "bin", "start")
	cmd := exec.Command(startPath, "init", "--help")
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err)

	outputStr := string(output)

	// Verify help text contains key information
	assert.Contains(t, outputStr, "init")
	assert.Contains(t, outputStr, "--local")
	assert.Contains(t, outputStr, "--force")
	assert.Contains(t, outputStr, "Interactive wizard")
}
