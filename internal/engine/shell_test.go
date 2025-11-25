package engine

import (
	"os/exec"
	"testing"
)

func TestDetectShell(t *testing.T) {
	shell := DetectShell()

	// Should return either "bash" or "sh"
	if shell != "bash" && shell != "sh" {
		t.Errorf("Expected 'bash' or 'sh', got %q", shell)
	}

	// If bash is available, should return "bash"
	if _, err := exec.LookPath("bash"); err == nil {
		if shell != "bash" {
			t.Errorf("bash is available but DetectShell returned %q", shell)
		}
	}

	// If we got "sh", bash should not be available
	if shell == "sh" {
		if _, err := exec.LookPath("bash"); err == nil {
			t.Errorf("bash is available but DetectShell returned 'sh'")
		}
	}
}
