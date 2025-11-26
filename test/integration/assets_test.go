package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grantcarthew/start/test/assert"
)

// TestPhase7_AssetsSearch tests `start assets search <query>`
// This test requires network access to GitHub. It will skip if network is unavailable.
func TestPhase7_AssetsSearch(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ensureStartBinary(t)

	// Run: start assets search claude
	startPath := filepath.Join("..", "..", "bin", "start")
	cmd := exec.Command(startPath, "assets", "search", "claude")
	output, err := cmd.CombinedOutput()
	t.Logf("Command output: %s", string(output))

	outputStr := string(output)

	// Handle expected failure scenarios
	if err != nil {
		// Network errors - skip test
		if strings.Contains(outputStr, "failed to parse catalog index") ||
			strings.Contains(outputStr, "Failed to fetch") {
			t.Skip("Skipping test - network unavailable or GitHub catalog not accessible")
		}

		t.Fatalf("assets search failed unexpectedly: %v\nOutput: %s", err, outputStr)
	}

	// Verify search results shown
	assert.Contains(t, outputStr, "Search results")
}

// TestPhase7_AssetsSearchNoResults tests search with no matches
// This test requires network access to GitHub. It will skip if network is unavailable.
func TestPhase7_AssetsSearchNoResults(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ensureStartBinary(t)

	// Run: start assets search nonexistentassetxyz123
	startPath := filepath.Join("..", "..", "bin", "start")
	cmd := exec.Command(startPath, "assets", "search", "nonexistentassetxyz123")
	output, err := cmd.CombinedOutput()
	t.Logf("Command output: %s", string(output))

	outputStr := string(output)

	// Handle expected failure scenarios
	if err != nil {
		// Network errors - skip test
		if strings.Contains(outputStr, "failed to parse catalog index") ||
			strings.Contains(outputStr, "Failed to fetch") {
			t.Skip("Skipping test - network unavailable or GitHub catalog not accessible")
		}

		// "No assets found" is expected (exits with error)
		if strings.Contains(outputStr, "No assets found") {
			return
		}

		t.Fatalf("assets search failed unexpectedly: %v\nOutput: %s", err, outputStr)
	}
}

// TestPhase7_AssetsSearchMinLength tests query minimum length validation
func TestPhase7_AssetsSearchMinLength(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ensureStartBinary(t)

	// Run: start assets search ab (too short)
	startPath := filepath.Join("..", "..", "bin", "start")
	cmd := exec.Command(startPath, "assets", "search", "ab")
	output, err := cmd.CombinedOutput()
	t.Logf("Command output: %s", string(output))

	outputStr := string(output)

	// Should fail with validation error
	if err == nil {
		t.Fatal("Expected error for query too short, got nil")
	}

	assert.Contains(t, outputStr, "at least 3 characters")
}

// TestPhase7_AssetsBrowse tests `start assets browse`
func TestPhase7_AssetsBrowse(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ensureStartBinary(t)

	// Note: This test can't verify the browser opened, only that the command succeeds
	// We'll just verify it doesn't error

	// Run: start assets browse
	startPath := filepath.Join("..", "..", "bin", "start")
	cmd := exec.Command(startPath, "assets", "browse")
	output, err := cmd.CombinedOutput()
	t.Logf("Command output: %s", string(output))

	if err != nil {
		outputStr := string(output)
		// It's okay if browser can't be opened (headless environment)
		if !strings.Contains(outputStr, "Failed to open browser") {
			t.Fatalf("assets browse failed unexpectedly: %v\nOutput: %s", err, outputStr)
		}
	}
}

// TestPhase7_AssetsInfo tests `start assets info <query>`
// This test requires network access to GitHub. It will skip if network is unavailable.
func TestPhase7_AssetsInfo(t *testing.T) {
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

	// Run: start assets info claude
	startPath := filepath.Join("..", "..", "bin", "start")
	cmd := exec.Command(startPath, "assets", "info", "claude")
	cmd.Env = env
	cmd.Stdin = strings.NewReader("1\n") // Select first result if multiple
	output, err := cmd.CombinedOutput()
	t.Logf("Command output: %s", string(output))

	outputStr := string(output)

	// Handle expected failure scenarios
	if err != nil {
		// Network errors - skip test
		if strings.Contains(outputStr, "failed to parse catalog index") ||
			strings.Contains(outputStr, "Failed to fetch") {
			t.Skip("Skipping test - network unavailable or GitHub catalog not accessible")
		}

		// "No assets found" is valid
		if strings.Contains(outputStr, "No assets found") {
			return
		}

		t.Fatalf("assets info failed unexpectedly: %v\nOutput: %s", err, outputStr)
	}

	// If successful, should show asset details
	if !strings.Contains(outputStr, "No assets found") {
		assert.Contains(t, outputStr, "Name:")
	}
}

// TestPhase7_AssetsInfoMinLength tests query minimum length validation
func TestPhase7_AssetsInfoMinLength(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ensureStartBinary(t)

	// Run: start assets info xy (too short)
	startPath := filepath.Join("..", "..", "bin", "start")
	cmd := exec.Command(startPath, "assets", "info", "xy")
	output, err := cmd.CombinedOutput()
	t.Logf("Command output: %s", string(output))

	outputStr := string(output)

	// Should fail with validation error
	if err == nil {
		t.Fatal("Expected error for query too short, got nil")
	}

	assert.Contains(t, outputStr, "at least 3 characters")
}

// TestPhase7_AssetsUpdate tests `start assets update`
// This test requires network access to GitHub. It will skip if network is unavailable.
func TestPhase7_AssetsUpdate(t *testing.T) {
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

	// Run: start assets update (no cache, should skip)
	startPath := filepath.Join("..", "..", "bin", "start")
	cmd := exec.Command(startPath, "assets", "update")
	cmd.Env = env
	output, err := cmd.CombinedOutput()
	t.Logf("Command output: %s", string(output))

	outputStr := string(output)

	// Handle expected failure scenarios
	if err != nil {
		// Network errors - skip test
		if strings.Contains(outputStr, "failed to parse catalog index") ||
			strings.Contains(outputStr, "Failed to fetch") {
			t.Skip("Skipping test - network unavailable or GitHub catalog not accessible")
		}

		// "No cached assets" is expected for empty cache
		if strings.Contains(outputStr, "No cached assets") {
			return
		}

		t.Fatalf("assets update failed unexpectedly: %v\nOutput: %s", err, outputStr)
	}
}

// TestPhase7_AssetsUpdateQuery tests `start assets update <query>`
// This test requires network access to GitHub. It will skip if network is unavailable.
func TestPhase7_AssetsUpdateQuery(t *testing.T) {
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

	// Run: start assets update claude (no cache, should skip)
	startPath := filepath.Join("..", "..", "bin", "start")
	cmd := exec.Command(startPath, "assets", "update", "claude")
	cmd.Env = env
	output, err := cmd.CombinedOutput()
	t.Logf("Command output: %s", string(output))

	outputStr := string(output)

	// Handle expected failure scenarios
	if err != nil {
		// Network errors - skip test
		if strings.Contains(outputStr, "failed to parse catalog index") ||
			strings.Contains(outputStr, "Failed to fetch") {
			t.Skip("Skipping test - network unavailable or GitHub catalog not accessible")
		}

		// "No assets found matching" is expected for empty cache
		if strings.Contains(outputStr, "No assets found matching") {
			return
		}

		t.Fatalf("assets update <query> failed unexpectedly: %v\nOutput: %s", err, outputStr)
	}
}

// TestPhase7_AssetsIndex tests `start assets index`
// This test requires network access to GitHub. It will skip if network is unavailable.
func TestPhase7_AssetsIndex(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ensureStartBinary(t)

	// Create temp directory for test
	tempDir := t.TempDir()

	// Create mock assets directory structure
	assetsDir := filepath.Join(tempDir, "assets")
	agentsDir := filepath.Join(assetsDir, "agents")
	err := os.MkdirAll(agentsDir, 0755)
	assert.NoError(t, err)

	// Create a test agent file
	testAgent := filepath.Join(agentsDir, "testagent.toml")
	testContent := `bin = "testagent"
description = "Test agent for integration testing"
`
	err = os.WriteFile(testAgent, []byte(testContent), 0644)
	assert.NoError(t, err)

	// Set environment
	env := []string{
		"HOME=" + tempDir,
		"PATH=" + os.Getenv("PATH"),
	}

	// Run: start assets index
	// Get absolute path to start binary
	startPath, err := filepath.Abs(filepath.Join("..", "..", "bin", "start"))
	assert.NoError(t, err)

	cmd := exec.Command(startPath, "assets", "index")
	cmd.Dir = tempDir // Run in temp dir where assets/ exists
	cmd.Env = env
	output, err := cmd.CombinedOutput()
	t.Logf("Command output: %s", string(output))

	outputStr := string(output)

	// Handle expected failure scenarios
	if err != nil {
		// Network errors - skip test
		if strings.Contains(outputStr, "failed to parse catalog index") ||
			strings.Contains(outputStr, "Failed to fetch") {
			t.Skip("Skipping test - network unavailable or GitHub catalog not accessible")
		}

		// "No assets directory found" is valid if assets/ doesn't exist in CWD
		if strings.Contains(outputStr, "No assets directory found") {
			return
		}

		// "not a git repository" is expected in temp directory
		if strings.Contains(outputStr, "not a git repository") {
			t.Skip("Skipping test - requires git repository")
		}

		t.Fatalf("assets index failed unexpectedly: %v\nOutput: %s", err, outputStr)
	}

	// Verify index file was created
	indexPath := filepath.Join(assetsDir, "INDEX.csv")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Error("Expected INDEX.csv to be created")
	} else {
		// Read and verify content
		indexData, err := os.ReadFile(indexPath)
		assert.NoError(t, err)

		indexStr := string(indexData)
		// Should contain header and test asset
		assert.Contains(t, indexStr, "type,name,description")
		assert.Contains(t, indexStr, "testagent")
	}
}

// TestPhase7_AssetsAdd tests `start assets add [query]`
// This test requires network access to GitHub. It will skip if network is unavailable.
func TestPhase7_AssetsAdd(t *testing.T) {
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

	// Run: start assets add claude
	startPath := filepath.Join("..", "..", "bin", "start")
	cmd := exec.Command(startPath, "assets", "add", "claude")
	cmd.Env = env
	cmd.Stdin = strings.NewReader("1\n") // Select first result if multiple
	output, err := cmd.CombinedOutput()
	t.Logf("Command output: %s", string(output))

	outputStr := string(output)

	// Handle expected failure scenarios
	if err != nil {
		// Network errors - skip test
		if strings.Contains(outputStr, "failed to parse catalog index") ||
			strings.Contains(outputStr, "Failed to fetch") {
			t.Skip("Skipping test - network unavailable or GitHub catalog not accessible")
		}

		// "No assets found" is valid
		if strings.Contains(outputStr, "No assets found") {
			return
		}

		t.Fatalf("assets add failed unexpectedly: %v\nOutput: %s", err, outputStr)
	}

	// If successful, should download to cache
	// Check if cache was created
	cacheDir := filepath.Join(tempDir, ".cache", "start", "assets")
	if _, err := os.Stat(cacheDir); err == nil {
		t.Logf("Cache directory created at %s", cacheDir)
	}
}

// TestPhase7_AssetsAddMinLength tests query minimum length validation
func TestPhase7_AssetsAddMinLength(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ensureStartBinary(t)

	// Run: start assets add ab (too short)
	startPath := filepath.Join("..", "..", "bin", "start")
	cmd := exec.Command(startPath, "assets", "add", "ab")
	output, err := cmd.CombinedOutput()
	t.Logf("Command output: %s", string(output))

	outputStr := string(output)

	// Should fail with validation error
	if err == nil {
		t.Fatal("Expected error for query too short, got nil")
	}

	assert.Contains(t, outputStr, "at least 3 characters")
}

// TestPhase7_AssetsHelp tests `start assets --help`
func TestPhase7_AssetsHelp(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ensureStartBinary(t)

	// Run: start assets --help
	startPath := filepath.Join("..", "..", "bin", "start")
	cmd := exec.Command(startPath, "assets", "--help")
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err)

	outputStr := string(output)

	// Verify help text contains key information
	assert.Contains(t, outputStr, "assets")
	assert.Contains(t, outputStr, "search")
	assert.Contains(t, outputStr, "add")
	assert.Contains(t, outputStr, "browse")
	assert.Contains(t, outputStr, "info")
	assert.Contains(t, outputStr, "update")
	assert.Contains(t, outputStr, "index")
}
