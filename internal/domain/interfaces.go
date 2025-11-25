package domain

import (
	"context"
	"os"
	"time"
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

// Runner abstracts command execution
type Runner interface {
	Run(ctx context.Context, shell, command string, timeout time.Duration) (stdout, stderr string, err error)
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
