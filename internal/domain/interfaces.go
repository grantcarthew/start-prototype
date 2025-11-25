package domain

import (
	"context"
	"os"
)

// FileSystem abstracts all file operations
type FileSystem interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte, perm os.FileMode) error
	Exists(path string) bool
	Glob(pattern string) ([]string, error)
	MkdirAll(path string, perm os.FileMode) error
	TempFile(pattern string) (name string, err error)
	Remove(path string) error
}

// Runner abstracts command execution via process replacement
type Runner interface {
	// Exec replaces the current process with the command
	// After successful exec, this function never returns
	// Only returns on error (before exec)
	Exec(shell, command string) error
}

// CommandRunner abstracts command execution with output capture
type CommandRunner interface {
	// Run executes a command and returns stdout+stderr combined output
	// Returns error if command fails or times out
	Run(shell, command string, timeoutSeconds int) (string, error)
}

// GitHubClient abstracts GitHub HTTP operations
type GitHubClient interface {
	FetchIndex(ctx context.Context, repo, branch string) ([]byte, error)
	FetchAsset(ctx context.Context, repo, branch, path string) ([]byte, error)
}

// Cache abstracts asset cache operations
type Cache interface {
	Get(assetType, name string) ([]byte, error)
	Set(assetType, name string, content []byte, meta AssetMeta) error
	List(assetType string) ([]CachedAsset, error)
	Delete(assetType, name string) error
}
