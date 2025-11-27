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

// TestPhase8c_ConfigRoleList tests role list command
func TestPhase8c_ConfigRoleList(t *testing.T) {
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
			name:         "role list help",
			args:         []string{"config", "role", "list", "--help"},
			expectOutput: "List all roles",
		},
		{
			name:         "role list",
			args:         []string{"config", "role", "list"},
			expectOutput: "roles",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()

			if err != nil && !strings.Contains(string(output), "failed to load") {
				t.Errorf("Command failed: %v\nOutput: %s", err, output)
			}

			if !strings.Contains(string(output), tt.expectOutput) {
				t.Errorf("Output does not contain expected string %q\nGot: %s", tt.expectOutput, output)
			}
		})
	}
}

// TestPhase8c_ConfigRoleShow tests role show command
func TestPhase8c_ConfigRoleShow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "role show help",
			args:        []string{"config", "role", "show", "--help"},
			expectError: false,
		},
		{
			name:        "role show nonexistent",
			args:        []string{"config", "role", "show", "nonexistent"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()

			if tt.expectError && err == nil {
				t.Errorf("Expected error but command succeeded\nOutput: %s", output)
			}

			if !tt.expectError && err != nil {
				t.Errorf("Command failed: %v\nOutput: %s", err, output)
			}
		})
	}
}

// TestPhase8c_ConfigRoleTest tests role test command
func TestPhase8c_ConfigRoleTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "role test help",
			args:        []string{"config", "role", "test", "--help"},
			expectError: false,
		},
		{
			name:        "role test nonexistent",
			args:        []string{"config", "role", "test", "nonexistent"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()

			if tt.expectError && err == nil {
				t.Errorf("Expected error but command succeeded\nOutput: %s", output)
			}

			if !tt.expectError && err != nil {
				t.Errorf("Command failed: %v\nOutput: %s", err, output)
			}
		})
	}
}

// TestPhase8c_ConfigRoleDefault tests role default command
func TestPhase8c_ConfigRoleDefault(t *testing.T) {
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
			name:         "role default help",
			args:         []string{"config", "role", "default", "--help"},
			expectOutput: "default role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()

			if err != nil && !strings.Contains(string(output), "failed to") {
				t.Errorf("Command failed: %v\nOutput: %s", err, output)
			}

			if !strings.Contains(string(output), tt.expectOutput) {
				t.Errorf("Output does not contain expected string %q\nGot: %s", tt.expectOutput, output)
			}
		})
	}
}

// TestPhase8c_ConfigRoleEdit tests role edit command
func TestPhase8c_ConfigRoleEdit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "config", "role", "edit", "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Command failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "Edit") {
		t.Errorf("Output does not contain expected string 'Edit'\nGot: %s", output)
	}
}

// TestPhase8c_ConfigRoleRemove tests role remove command
func TestPhase8c_ConfigRoleRemove(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "config", "role", "remove", "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Command failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "Remove") {
		t.Errorf("Output does not contain expected string 'Remove'\nGot: %s", output)
	}
}

// TestPhase8c_ConfigRoleNew tests role new command
func TestPhase8c_ConfigRoleNew(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "config", "role", "new", "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Command failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "Interactive wizard") {
		t.Errorf("Output does not contain expected string 'Interactive wizard'\nGot: %s", output)
	}
}

// TestPhase8b_ConfigAgentList tests agent list command
func TestPhase8b_ConfigAgentList(t *testing.T) {
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
			name:         "agent list help",
			args:         []string{"config", "agent", "list", "--help"},
			expectOutput: "List all agents",
		},
		{
			name:         "agent list",
			args:         []string{"config", "agent", "list"},
			expectOutput: "agents",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()

			if err != nil && !strings.Contains(string(output), "failed to load") {
				t.Errorf("Command failed: %v\nOutput: %s", err, output)
			}

			if !strings.Contains(string(output), tt.expectOutput) {
				t.Errorf("Output does not contain expected string %q\nGot: %s", tt.expectOutput, output)
			}
		})
	}
}

// TestPhase8b_ConfigAgentShow tests agent show command
func TestPhase8b_ConfigAgentShow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "agent show help",
			args:        []string{"config", "agent", "show", "--help"},
			expectError: false,
		},
		{
			name:        "agent show nonexistent",
			args:        []string{"config", "agent", "show", "nonexistent"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()

			if tt.expectError && err == nil {
				t.Errorf("Expected error but command succeeded\nOutput: %s", output)
			}

			if !tt.expectError && err != nil {
				t.Errorf("Command failed: %v\nOutput: %s", err, output)
			}
		})
	}
}

// TestPhase8b_ConfigAgentTest tests agent test command
func TestPhase8b_ConfigAgentTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "agent test help",
			args:        []string{"config", "agent", "test", "--help"},
			expectError: false,
		},
		{
			name:        "agent test nonexistent",
			args:        []string{"config", "agent", "test", "nonexistent"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()

			if tt.expectError && err == nil {
				t.Errorf("Expected error but command succeeded\nOutput: %s", output)
			}

			if !tt.expectError && err != nil {
				t.Errorf("Command failed: %v\nOutput: %s", err, output)
			}
		})
	}
}

// TestPhase8b_ConfigAgentDefault tests agent default command
func TestPhase8b_ConfigAgentDefault(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "agent default help",
			args: []string{"config", "agent", "default", "--help"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()

			if err != nil && !strings.Contains(string(output), "failed to load") {
				t.Errorf("Command failed: %v\nOutput: %s", err, output)
			}
		})
	}
}

// TestPhase8b_ConfigAgentEdit tests agent edit command help
func TestPhase8b_ConfigAgentEdit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "config", "agent", "edit", "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Command failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "help for edit") {
		t.Errorf("Output does not contain expected string 'help for edit'\nGot: %s", output)
	}
}

// TestPhase8b_ConfigAgentRemove tests agent remove command help
func TestPhase8b_ConfigAgentRemove(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "config", "agent", "remove", "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Command failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "Remove") {
		t.Errorf("Output does not contain expected string 'Remove'\nGot: %s", output)
	}
}

// TestPhase8b_ConfigAgentNew tests agent new command help
func TestPhase8b_ConfigAgentNew(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "config", "agent", "new", "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Command failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "Interactive wizard") {
		t.Errorf("Output does not contain expected string 'Interactive wizard'\nGot: %s", output)
	}
}

// TestPhase8d_ConfigContextList tests context list command
func TestPhase8d_ConfigContextList(t *testing.T) {
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
			name:         "context list help",
			args:         []string{"config", "context", "list", "--help"},
			expectOutput: "List all context",
		},
		{
			name:         "context list",
			args:         []string{"config", "context", "list"},
			expectOutput: "No contexts configured.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()

			if err != nil && !strings.Contains(string(output), "failed to load") {
				t.Errorf("Command failed: %v\nOutput: %s", err, output)
			}

			if !strings.Contains(string(output), tt.expectOutput) {
				t.Errorf("Output does not contain expected string %q\nGot: %s", tt.expectOutput, output)
			}
		})
	}
}

// TestPhase8d_ConfigContextShow tests context show command
func TestPhase8d_ConfigContextShow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "context show help",
			args:        []string{"config", "context", "show", "--help"},
			expectError: false,
		},
		{
			name:        "context show nonexistent",
			args:        []string{"config", "context", "show", "nonexistent"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()

			if tt.expectError && err == nil {
				t.Errorf("Expected error but command succeeded\nOutput: %s", output)
			}

			if !tt.expectError && err != nil {
				t.Errorf("Command failed: %v\nOutput: %s", err, output)
			}
		})
	}
}

// TestPhase8d_ConfigContextTest tests context test command
func TestPhase8d_ConfigContextTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "context test help",
			args:        []string{"config", "context", "test", "--help"},
			expectError: false,
		},
		{
			name:        "context test nonexistent",
			args:        []string{"config", "context", "test", "nonexistent"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()

			if tt.expectError && err == nil {
				t.Errorf("Expected error but command succeeded\nOutput: %s", output)
			}

			if !tt.expectError && err != nil {
				t.Errorf("Command failed: %v\nOutput: %s", err, output)
			}
		})
	}
}

// TestPhase8d_ConfigContextEdit tests context edit command help
func TestPhase8d_ConfigContextEdit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "config", "context", "edit", "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Command failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "help for edit") {
		t.Errorf("Output does not contain expected string 'help for edit'\nGot: %s", output)
	}
}

// TestPhase8d_ConfigContextRemove tests context remove command help
func TestPhase8d_ConfigContextRemove(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "config", "context", "remove", "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Command failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "Remove") {
		t.Errorf("Output does not contain expected string 'Remove'\nGot: %s", output)
	}
}

// TestPhase8d_ConfigContextNew tests context new command help
func TestPhase8d_ConfigContextNew(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "config", "context", "new", "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Command failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "Interactive wizard") {
		t.Errorf("Output does not contain expected string 'Interactive wizard'\nGot: %s", output)
	}
}

// TestPhase8e_ConfigTaskList tests task list command
func TestPhase8e_ConfigTaskList(t *testing.T) {
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
			name:         "task list help",
			args:         []string{"config", "task", "list", "--help"},
			expectOutput: "List all tasks",
		},
		{
			name:         "task list",
			args:         []string{"config", "task", "list"},
			expectOutput: "No tasks configured.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()

			if err != nil && !strings.Contains(string(output), "failed to load") {
				t.Errorf("Command failed: %v\nOutput: %s", err, output)
			}

			if !strings.Contains(string(output), tt.expectOutput) {
				t.Errorf("Output does not contain expected string %q\nGot: %s", tt.expectOutput, output)
			}
		})
	}
}

// TestPhase8e_ConfigTaskShow tests task show command
func TestPhase8e_ConfigTaskShow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "task show help",
			args:        []string{"config", "task", "show", "--help"},
			expectError: false,
		},
		{
			name:        "task show nonexistent",
			args:        []string{"config", "task", "show", "nonexistent"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()

			if tt.expectError && err == nil {
				t.Errorf("Expected error but command succeeded\nOutput: %s", output)
			}

			if !tt.expectError && err != nil {
				t.Errorf("Command failed: %v\nOutput: %s", err, output)
			}
		})
	}
}

// TestPhase8e_ConfigTaskTest tests task test command
func TestPhase8e_ConfigTaskTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "task test help",
			args:        []string{"config", "task", "test", "--help"},
			expectError: false,
		},
		{
			name:        "task test nonexistent",
			args:        []string{"config", "task", "test", "nonexistent"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()

			if tt.expectError && err == nil {
				t.Errorf("Expected error but command succeeded\nOutput: %s", output)
			}

			if !tt.expectError && err != nil {
				t.Errorf("Command failed: %v\nOutput: %s", err, output)
			}
		})
	}
}

// TestPhase8e_ConfigTaskEdit tests task edit command help
func TestPhase8e_ConfigTaskEdit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "config", "task", "edit", "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Command failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "help for edit") {
		t.Errorf("Output does not contain expected string 'help for edit'\nGot: %s", output)
	}
}

// TestPhase8e_ConfigTaskRemove tests task remove command help
func TestPhase8e_ConfigTaskRemove(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "config", "task", "remove", "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Command failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "Remove") {
		t.Errorf("Output does not contain expected string 'Remove'\nGot: %s", output)
	}
}

// TestPhase8e_ConfigTaskNew tests task new command help
func TestPhase8e_ConfigTaskNew(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "config", "task", "new", "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Command failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "Interactive wizard") {
		t.Errorf("Output does not contain expected string 'Interactive wizard'\nGot: %s", output)
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
