package mocks

import (
	"context"
	"fmt"
)

type MockGitHubClient struct {
	Index  []byte
	Assets map[string][]byte // path -> content
}

func NewMockGitHubClient() *MockGitHubClient {
	return &MockGitHubClient{
		Assets: make(map[string][]byte),
	}
}

func (m *MockGitHubClient) FetchIndex(ctx context.Context, repo, branch string) ([]byte, error) {
	if m.Index == nil {
		return nil, fmt.Errorf("no index configured in mock")
	}
	return m.Index, nil
}

func (m *MockGitHubClient) FetchAsset(ctx context.Context, repo, branch, path string) ([]byte, error) {
	content, ok := m.Assets[path]
	if !ok {
		return nil, fmt.Errorf("asset not found in mock: %s", path)
	}
	return content, nil
}
