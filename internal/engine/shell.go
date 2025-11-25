package engine

import "os/exec"

// DetectShell detects the default shell to use for command execution
// Returns "bash" if available, otherwise "sh"
func DetectShell() string {
	// Try bash first
	if _, err := exec.LookPath("bash"); err == nil {
		return "bash"
	}

	// Fall back to sh
	return "sh"
}
