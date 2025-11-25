package mocks

import (
	"context"
	"fmt"
	"time"
)

// MockRunner is a mock implementation of the Runner interface
type MockRunner struct {
	Outputs      map[string]string // command -> output
	CalledWith   []CallRecord
	ShouldError  bool
	ErrorMessage string
}

// CallRecord tracks a single call to Run
type CallRecord struct {
	Shell   string
	Command string
	Timeout time.Duration
}

func NewMockRunner() *MockRunner {
	return &MockRunner{
		Outputs: make(map[string]string),
	}
}

// Run executes a mock command
func (m *MockRunner) Run(ctx context.Context, shell, command string, timeout time.Duration) (string, string, error) {
	// Record the call
	if m.CalledWith == nil {
		m.CalledWith = []CallRecord{}
	}
	m.CalledWith = append(m.CalledWith, CallRecord{
		Shell:   shell,
		Command: command,
		Timeout: timeout,
	})

	// Return error if configured
	if m.ShouldError {
		return "", "", fmt.Errorf("%s", m.ErrorMessage)
	}

	// Return configured output if available
	if m.Outputs != nil {
		if output, ok := m.Outputs[command]; ok {
			return output, "", nil
		}
	}

	// Default: return empty output
	return "", "", nil
}
