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

// MockCommandRunner is a mock implementation of the CommandRunner interface
type MockCommandRunner struct {
	output  string
	err     error
	Outputs map[string]string // command -> output mapping
}

func NewMockCommandRunner() *MockCommandRunner {
	return &MockCommandRunner{
		Outputs: make(map[string]string),
	}
}

// SetOutput sets the output and error to return from Run
func (m *MockCommandRunner) SetOutput(output string, err error) {
	m.output = output
	m.err = err
}

// Run simulates command execution with output capture
func (m *MockCommandRunner) Run(shell, command string, timeoutSeconds int) (string, error) {
	// Check if we have a specific output for this command
	if output, ok := m.Outputs[command]; ok {
		return output, nil
	}
	// Fall back to default output/error
	return m.output, m.err
}
