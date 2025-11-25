package mocks

import (
	"fmt"
	"os"
	"path/filepath"
)

type MockFileSystem struct {
	Files map[string]string // path -> content
}

func NewMockFileSystem() *MockFileSystem {
	return &MockFileSystem{
		Files: make(map[string]string),
	}
}

func (m *MockFileSystem) ReadFile(path string) ([]byte, error) {
	content, ok := m.Files[path]
	if !ok {
		return nil, os.ErrNotExist
	}
	return []byte(content), nil
}

func (m *MockFileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	if m.Files == nil {
		m.Files = make(map[string]string)
	}
	m.Files[path] = string(data)
	return nil
}

func (m *MockFileSystem) Exists(path string) bool {
	_, ok := m.Files[path]
	return ok
}

func (m *MockFileSystem) Glob(pattern string) ([]string, error) {
	var matches []string
	for path := range m.Files {
		matched, _ := filepath.Match(pattern, path)
		if matched {
			matches = append(matches, path)
		}
	}
	return matches, nil
}

func (m *MockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return nil
}

func (m *MockFileSystem) TempFile(pattern string) (string, error) {
	path := fmt.Sprintf("/tmp/mock-%s-%d", pattern, len(m.Files))
	m.Files[path] = ""
	return path, nil
}

func (m *MockFileSystem) Remove(path string) error {
	delete(m.Files, path)
	return nil
}
