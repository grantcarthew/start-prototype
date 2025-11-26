package integration

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestPhase8_PrefixMatching tests prefix matching for commands
func TestPhase8_PrefixMatching(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	tests := []struct {
		name         string
		args         []string
		expectError  bool
		expectOutput string
	}{
		{
			name:         "prefix match doctor",
			args:         []string{"d", "--help"},
			expectError:  false,
			expectOutput: "health check",
		},
		{
			name:         "prefix match completion",
			args:         []string{"com", "--help"},
			expectError:  false,
			expectOutput: "shell completion",
		},
		{
			name:         "prefix match config",
			args:         []string{"con", "--help"},
			expectError:  false,
			expectOutput: "configuration",
		},
		{
			name:         "prefix match task",
			args:         []string{"t", "--help"},
			expectError:  false,
			expectOutput: "workflow tasks",
		},
		{
			name:         "prefix match assets",
			args:         []string{"ass", "--help"},
			expectError:  false,
			expectOutput: "catalog",
		},
		{
			name:         "prefix match init",
			args:         []string{"i", "--help"},
			expectError:  false,
			expectOutput: "configuration files",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()

			if tt.expectError && err == nil {
				t.Errorf("Expected error but command succeeded")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Command failed: %v\nOutput: %s", err, output)
			}

			if tt.expectOutput != "" && !strings.Contains(string(output), tt.expectOutput) {
				t.Errorf("Output does not contain expected string %q\nGot: %s", tt.expectOutput, output)
			}
		})
	}
}

// TestPhase8_CompletionCommand tests the completion command
func TestPhase8_CompletionCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	tests := []struct {
		name         string
		args         []string
		expectError  bool
		expectOutput string
	}{
		{
			name:         "completion help",
			args:         []string{"completion", "--help"},
			expectError:  false,
			expectOutput: "Generate shell completion scripts",
		},
		{
			name:         "completion bash",
			args:         []string{"completion", "bash"},
			expectError:  false,
			expectOutput: "bash completion for start",
		},
		{
			name:         "completion zsh",
			args:         []string{"completion", "zsh"},
			expectError:  false,
			expectOutput: "#compdef start",
		},
		{
			name:         "completion fish",
			args:         []string{"completion", "fish"},
			expectError:  false,
			expectOutput: "complete -c start",
		},
		{
			name:         "completion invalid shell",
			args:         []string{"completion", "invalid"},
			expectError:  true,
			expectOutput: "",
		},
		{
			name:         "completion install help",
			args:         []string{"completion", "install", "--help"},
			expectError:  false,
			expectOutput: "Auto-install completion",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()

			if tt.expectError && err == nil {
				t.Errorf("Expected error but command succeeded")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Command failed: %v\nOutput: %s", err, output)
			}

			if tt.expectOutput != "" && !strings.Contains(string(output), tt.expectOutput) {
				t.Errorf("Output does not contain expected string %q\nGot: %s", tt.expectOutput, output)
			}
		})
	}
}

// TestPhase8_DoctorCommand tests the doctor command
func TestPhase8_DoctorCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	tests := []struct {
		name         string
		args         []string
		expectOutput string
	}{
		{
			name:         "doctor help",
			args:         []string{"doctor", "--help"},
			expectOutput: "health check",
		},
		{
			name:         "doctor run",
			args:         []string{"doctor"},
			expectOutput: "Diagnosing start installation",
		},
		{
			name:         "doctor verbose",
			args:         []string{"doctor", "--verbose"},
			expectOutput: "Diagnosing start installation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()

			// Doctor may exit with 1 if issues found, that's ok
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					if exitErr.ExitCode() != 1 && exitErr.ExitCode() != 0 {
						t.Errorf("Unexpected exit code: %d\nOutput: %s", exitErr.ExitCode(), output)
					}
				}
			}

			if !strings.Contains(string(output), tt.expectOutput) {
				t.Errorf("Output does not contain expected string %q\nGot: %s", tt.expectOutput, output)
			}
		})
	}
}

// TestPhase8_DoctorExitCodes tests doctor exit codes
func TestPhase8_DoctorExitCodes(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	// Doctor should exit with 0 or 1 (per DR-024)
	cmd := exec.Command(binary, "doctor", "--quiet")
	err := cmd.Run()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode := exitErr.ExitCode()
			if exitCode != 0 && exitCode != 1 {
				t.Errorf("Expected exit code 0 or 1, got %d", exitCode)
			}
		} else {
			t.Errorf("Command failed with non-exit error: %v", err)
		}
	}
	// Exit code 0 is also valid (no issues found)
}

// TestPhase8_DoctorQuietMode tests doctor quiet mode
func TestPhase8_DoctorQuietMode(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "doctor", "--quiet")
	output, _ := cmd.CombinedOutput()

	// Quiet mode should have less output (no section headers)
	if strings.Contains(string(output), "═══════════════") {
		t.Error("Quiet mode should not show section separators")
	}
}

// TestPhase8_CompletionBashOutput tests bash completion output structure
func TestPhase8_CompletionBashOutput(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "completion", "bash")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to generate bash completion: %v\nOutput: %s", err, output)
	}

	// Check for essential bash completion components
	requiredStrings := []string{
		"bash completion for start",
		"__start_",
		"complete -o default",
	}

	for _, required := range requiredStrings {
		if !strings.Contains(string(output), required) {
			t.Errorf("Bash completion missing required string: %q", required)
		}
	}
}

// TestPhase8_CompletionZshOutput tests zsh completion output structure
func TestPhase8_CompletionZshOutput(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "completion", "zsh")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to generate zsh completion: %v\nOutput: %s", err, output)
	}

	// Check for essential zsh completion components
	requiredStrings := []string{
		"#compdef start",
		"_start",
	}

	for _, required := range requiredStrings {
		if !strings.Contains(string(output), required) {
			t.Errorf("Zsh completion missing required string: %q", required)
		}
	}
}

// TestPhase8_CompletionFishOutput tests fish completion output structure
func TestPhase8_CompletionFishOutput(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "completion", "fish")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to generate fish completion: %v\nOutput: %s", err, output)
	}

	// Check for essential fish completion components
	requiredStrings := []string{
		"complete -c start",
	}

	for _, required := range requiredStrings {
		if !strings.Contains(string(output), required) {
			t.Errorf("Fish completion missing required string: %q", required)
		}
	}
}

// getBinaryPath returns the path to the test binary
func getBinaryPath(t *testing.T) string {
	t.Helper()

	// Check if binary exists
	binary := "../../bin/start"
	if _, err := os.Stat(binary); err != nil {
		t.Fatalf("Binary not found at %s. Run 'go build -o bin/start cmd/start/main.go' first", binary)
	}

	return binary
}
