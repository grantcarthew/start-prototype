package mocks

import (
	"context"
	"fmt"
	"time"
)

type MockRunner struct {
	Outputs map[string]string // command -> output
}

func NewMockRunner() *MockRunner {
	return &MockRunner{
		Outputs: make(map[string]string),
	}
}

func (m *MockRunner) Run(ctx context.Context, shell, command string, timeout time.Duration) (string, string, error) {
	output, ok := m.Outputs[command]
	if !ok {
		return "", "", fmt.Errorf("command not found in mock: %s", command)
	}
	return output, "", nil
}
