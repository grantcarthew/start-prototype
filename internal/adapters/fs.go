package adapters

import (
	"os"
	"path/filepath"
	"strings"
)

// RealFileSystem implements the FileSystem interface using the real OS filesystem
type RealFileSystem struct{}

// ReadFile reads a file from the filesystem
func (fs *RealFileSystem) ReadFile(path string) ([]byte, error) {
	expanded := expandPath(path)
	return os.ReadFile(expanded)
}

// WriteFile writes data to a file
func (fs *RealFileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	expanded := expandPath(path)
	return os.WriteFile(expanded, data, perm)
}

// Exists checks if a path exists
func (fs *RealFileSystem) Exists(path string) bool {
	expanded := expandPath(path)
	_, err := os.Stat(expanded)
	return err == nil
}

// Glob returns paths matching a pattern
func (fs *RealFileSystem) Glob(pattern string) ([]string, error) {
	expanded := expandPath(pattern)
	return filepath.Glob(expanded)
}

// MkdirAll creates a directory and all parents
func (fs *RealFileSystem) MkdirAll(path string, perm os.FileMode) error {
	expanded := expandPath(path)
	return os.MkdirAll(expanded, perm)
}

// TempFile creates a temporary file
func (fs *RealFileSystem) TempFile(pattern string) (string, error) {
	tmpFile, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", err
	}
	name := tmpFile.Name()
	tmpFile.Close()
	return name, nil
}

// Remove removes a file or directory
func (fs *RealFileSystem) Remove(path string) error {
	expanded := expandPath(path)
	return os.Remove(expanded)
}

// expandPath expands ~ to the user's home directory
func expandPath(path string) string {
	if !strings.HasPrefix(path, "~/") {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		// If we can't get home dir, return path as-is
		return path
	}

	return filepath.Join(home, path[2:])
}
