package mocks

import (
	"fmt"
)

// MockRunner is a mock implementation of the Runner interface
type MockRunner struct {
	CalledWith   []CallRecord
	ShouldError  bool
	ErrorMessage string
}

// CallRecord tracks a single call to Exec
type CallRecord struct {
	Shell   string
	Command string
}

func NewMockRunner() *MockRunner {
	return &MockRunner{}
}

// Exec simulates process replacement (doesn't actually replace for testing)
func (m *MockRunner) Exec(shell, command string) error {
	// Record the call
	if m.CalledWith == nil {
		m.CalledWith = []CallRecord{}
	}
	m.CalledWith = append(m.CalledWith, CallRecord{
		Shell:   shell,
		Command: command,
	})

	// Return error if configured
	if m.ShouldError {
		return fmt.Errorf("%s", m.ErrorMessage)
	}

	// In real implementation, this never returns on success
	// For testing, we just return nil to simulate success
	return nil
}
