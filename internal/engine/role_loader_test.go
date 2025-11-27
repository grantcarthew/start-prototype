package engine

import (
	"fmt"
	"os"
	"testing"

	"github.com/grantcarthew/start/internal/domain"
)

// mockCommandRunner for testing
type mockCommandRunner struct {
	output string
	err    error
}

func (m *mockCommandRunner) Run(shell string, command string, timeout int) (string, error) {
	return m.output, m.err
}

// mockFileSystem for testing
type mockFileSystem struct {
	files        map[string]string
	tempCounter  int
	writeError   error
	tempError    error
	removeError  error
	removeCount  int
	removedPaths []string
}

func newMockFileSystem() *mockFileSystem {
	return &mockFileSystem{
		files:        make(map[string]string),
		removedPaths: []string{},
	}
}

func (m *mockFileSystem) ReadFile(path string) ([]byte, error) {
	content, ok := m.files[path]
	if !ok {
		return nil, os.ErrNotExist
	}
	return []byte(content), nil
}

func (m *mockFileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	if m.writeError != nil {
		return m.writeError
	}
	m.files[path] = string(data)
	return nil
}

func (m *mockFileSystem) Exists(path string) bool {
	_, ok := m.files[path]
	return ok
}

func (m *mockFileSystem) Glob(pattern string) ([]string, error) {
	return nil, nil
}

func (m *mockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return nil
}

func (m *mockFileSystem) TempFile(pattern string) (string, error) {
	if m.tempError != nil {
		return "", m.tempError
	}
	m.tempCounter++
	return fmt.Sprintf("/tmp/test-%d.md", m.tempCounter), nil
}

func (m *mockFileSystem) Remove(path string) error {
	m.removeCount++
	m.removedPaths = append(m.removedPaths, path)
	if m.removeError != nil {
		return m.removeError
	}
	delete(m.files, path)
	return nil
}

func TestRoleLoader_LoadRole_SimpleFileRole(t *testing.T) {
	// Simple role with just a file path
	fs := newMockFileSystem()
	fs.files["/role.md"] = "Role content"

	cmdRunner := &mockCommandRunner{}
	utdProcessor := NewUTDProcessor(fs, cmdRunner, "/workdir")

	loader := NewRoleLoader(utdProcessor, fs)

	role := domain.Role{
		Name: "test-role",
		File: "/role.md",
	}

	result, err := loader.LoadRole(role, "bash", 30)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.Name != "test-role" {
		t.Errorf("Expected name %q, got %q", "test-role", result.Name)
	}

	if result.Content != "Role content" {
		t.Errorf("Expected content %q, got %q", "Role content", result.Content)
	}

	if result.FilePath != "/role.md" {
		t.Errorf("Expected file path %q, got %q", "/role.md", result.FilePath)
	}

	if result.IsTemp {
		t.Error("Expected IsTemp to be false for simple file role")
	}
}

func TestRoleLoader_LoadRole_ComplexRole(t *testing.T) {
	// Complex role with prompt or command - should create temp file
	fs := newMockFileSystem()

	cmdRunner := &mockCommandRunner{}
	utdProcessor := NewUTDProcessor(fs, cmdRunner, "/workdir")

	loader := NewRoleLoader(utdProcessor, fs)

	role := domain.Role{
		Name:   "complex-role",
		Prompt: "Some prompt text",
	}

	result, err := loader.LoadRole(role, "bash", 30)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.Name != "complex-role" {
		t.Errorf("Expected name %q, got %q", "complex-role", result.Name)
	}

	if result.Content != "Some prompt text" {
		t.Errorf("Expected content %q, got %q", "Some prompt text", result.Content)
	}

	if !result.IsTemp {
		t.Error("Expected IsTemp to be true for complex role")
	}

	if result.FilePath == "" {
		t.Error("Expected temp file path to be set")
	}

	// Verify temp file was written
	content := fs.files[result.FilePath]
	if content != "Some prompt text" {
		t.Errorf("Expected temp file content %q, got %q", "Some prompt text", content)
	}
}

func TestRoleLoader_LoadRole_FileWithCommand(t *testing.T) {
	// Role with both file and command - should create temp file
	fs := newMockFileSystem()
	fs.files["/original.md"] = "File content"

	cmdRunner := &mockCommandRunner{output: "\nextra"}
	utdProcessor := NewUTDProcessor(fs, cmdRunner, "/workdir")

	loader := NewRoleLoader(utdProcessor, fs)

	role := domain.Role{
		Name:    "combined-role",
		File:    "/original.md",
		Command: "echo 'extra'",
	}

	result, err := loader.LoadRole(role, "bash", 30)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !result.IsTemp {
		t.Error("Expected IsTemp to be true for role with file and command")
	}

	if result.FilePath == "/original.md" {
		t.Error("Expected temp file path, not original file path")
	}
}

func TestRoleLoader_LoadRole_UTDProcessingFailed(t *testing.T) {
	fs := newMockFileSystem()
	// Don't add file to fs - will cause UTD to fail

	cmdRunner := &mockCommandRunner{}
	utdProcessor := NewUTDProcessor(fs, cmdRunner, "/workdir")

	loader := NewRoleLoader(utdProcessor, fs)

	// Role with file referenced in prompt - will cause skipped
	role := domain.Role{
		Name:   "failed-role",
		File:   "/nonexistent.md",
		Prompt: "Use file: {file_contents}",
	}

	_, err := loader.LoadRole(role, "bash", 30)

	if err == nil {
		t.Error("Expected error when UTD processing fails")
		return
	}

	if !contains(err.Error(), "role processing failed") {
		t.Errorf("Expected error about role processing, got: %v", err)
	}
}

func TestRoleLoader_LoadRole_TempFileCreationFails(t *testing.T) {
	fs := newMockFileSystem()
	fs.tempError = fmt.Errorf("cannot create temp file")

	cmdRunner := &mockCommandRunner{}
	utdProcessor := NewUTDProcessor(fs, cmdRunner, "/workdir")

	loader := NewRoleLoader(utdProcessor, fs)

	role := domain.Role{
		Name:   "test-role",
		Prompt: "Some prompt",
	}

	_, err := loader.LoadRole(role, "bash", 30)

	if err == nil {
		t.Error("Expected error when temp file creation fails")
	}

	if !contains(err.Error(), "failed to create temp file") {
		t.Errorf("Expected error about temp file creation, got: %v", err)
	}
}

func TestRoleLoader_LoadRole_TempFileWriteFails(t *testing.T) {
	fs := newMockFileSystem()
	fs.writeError = fmt.Errorf("write permission denied")

	cmdRunner := &mockCommandRunner{}
	utdProcessor := NewUTDProcessor(fs, cmdRunner, "/workdir")

	loader := NewRoleLoader(utdProcessor, fs)

	role := domain.Role{
		Name:   "test-role",
		Prompt: "Some prompt",
	}

	_, err := loader.LoadRole(role, "bash", 30)

	if err == nil {
		t.Error("Expected error when temp file write fails")
	}

	if !contains(err.Error(), "failed to write temp file") {
		t.Errorf("Expected error about temp file write, got: %v", err)
	}

	// Verify temp file was removed after write error
	if fs.removeCount != 1 {
		t.Errorf("Expected 1 remove call after write error, got %d", fs.removeCount)
	}
}

func TestRoleLoader_LoadRole_UTDPassesDefaultsCorrectly(t *testing.T) {
	fs := newMockFileSystem()
	fs.files["/role.md"] = "Role content"

	cmdRunner := &mockCommandRunner{}
	utdProcessor := NewUTDProcessor(fs, cmdRunner, "/workdir")

	loader := NewRoleLoader(utdProcessor, fs)

	role := domain.Role{
		Name: "test-role",
		File: "/role.md",
	}

	result, err := loader.LoadRole(role, "bash", 30)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify basic processing worked
	if result.Content != "Role content" {
		t.Errorf("Expected content %q, got %q", "Role content", result.Content)
	}
}

func TestRoleLoader_CleanupRole_TempFile(t *testing.T) {
	fs := newMockFileSystem()
	loader := NewRoleLoader(nil, fs)

	tempPath := "/tmp/role-123.md"
	fs.files[tempPath] = "content"

	role := LoadedRole{
		Name:     "test",
		FilePath: tempPath,
		IsTemp:   true,
	}

	err := loader.CleanupRole(role)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if fs.removeCount != 1 {
		t.Errorf("Expected 1 remove call, got %d", fs.removeCount)
	}

	if len(fs.removedPaths) != 1 || fs.removedPaths[0] != tempPath {
		t.Errorf("Expected temp file %q to be removed", tempPath)
	}
}

func TestRoleLoader_CleanupRole_NotTemp(t *testing.T) {
	fs := newMockFileSystem()
	loader := NewRoleLoader(nil, fs)

	role := LoadedRole{
		Name:     "test",
		FilePath: "/original.md",
		IsTemp:   false,
	}

	err := loader.CleanupRole(role)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if fs.removeCount != 0 {
		t.Errorf("Expected no remove calls for non-temp file, got %d", fs.removeCount)
	}
}

func TestRoleLoader_CleanupRole_EmptyPath(t *testing.T) {
	fs := newMockFileSystem()
	loader := NewRoleLoader(nil, fs)

	role := LoadedRole{
		Name:     "test",
		FilePath: "",
		IsTemp:   true,
	}

	err := loader.CleanupRole(role)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if fs.removeCount != 0 {
		t.Errorf("Expected no remove calls for empty path, got %d", fs.removeCount)
	}
}
