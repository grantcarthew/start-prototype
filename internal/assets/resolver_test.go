package assets

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/grantcarthew/start/internal/config"
	"github.com/grantcarthew/start/internal/domain"
)

// mockCache for testing
type mockCache struct {
	data  map[string]map[string][]byte // assetType -> name -> data
	setError error
	getError error
}

func newMockCache() *mockCache {
	return &mockCache{
		data: make(map[string]map[string][]byte),
	}
}

func (m *mockCache) Get(assetType, name string) ([]byte, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	if typeData, ok := m.data[assetType]; ok {
		if data, ok := typeData[name]; ok {
			return data, nil
		}
	}
	return nil, os.ErrNotExist
}

func (m *mockCache) Set(assetType, name string, data []byte, meta domain.AssetMeta) error {
	if m.setError != nil {
		return m.setError
	}
	if m.data[assetType] == nil {
		m.data[assetType] = make(map[string][]byte)
	}
	m.data[assetType][name] = data
	return nil
}

func (m *mockCache) List(assetType string) ([]domain.CachedAsset, error) {
	return nil, nil
}

func (m *mockCache) Delete(assetType, name string) error {
	return nil
}

// mockGitHubClient for testing
type mockGitHubClient struct {
	indexData []byte
	assetData []byte
	indexErr  error
	assetErr  error
}

func (m *mockGitHubClient) FetchIndex(ctx context.Context, repo, branch string) ([]byte, error) {
	if m.indexErr != nil {
		return nil, m.indexErr
	}
	return m.indexData, nil
}

func (m *mockGitHubClient) FetchAsset(ctx context.Context, repo, branch, path string) ([]byte, error) {
	if m.assetErr != nil {
		return nil, m.assetErr
	}
	return m.assetData, nil
}

// mockFileSystem for testing
type mockFS struct {
	files map[string]string
}

func newMockFS() *mockFS {
	return &mockFS{
		files: make(map[string]string),
	}
}

func (m *mockFS) ReadFile(path string) ([]byte, error) {
	if content, ok := m.files[path]; ok {
		return []byte(content), nil
	}
	return nil, os.ErrNotExist
}

func (m *mockFS) WriteFile(path string, data []byte, perm os.FileMode) error {
	m.files[path] = string(data)
	return nil
}

func (m *mockFS) Exists(path string) bool {
	_, ok := m.files[path]
	return ok
}

func (m *mockFS) Glob(pattern string) ([]string, error) {
	return nil, nil
}

func (m *mockFS) MkdirAll(path string, perm os.FileMode) error {
	return nil
}

func (m *mockFS) TempFile(pattern string) (string, error) {
	return "/tmp/test.tmp", nil
}

func (m *mockFS) Remove(path string) error {
	delete(m.files, path)
	return nil
}

func TestResolver_ResolveTask_LocalConfig(t *testing.T) {
	fs := newMockFS()
	cache := newMockCache()
	github := &mockGitHubClient{}
	configLoader := config.NewLoader(fs)

	resolver := NewResolver(fs, cache, github, configLoader)

	// Task exists in local config
	cfg := domain.Config{
		Tasks: map[string]domain.Task{
			"local-task": {
				Prompt:      "Local task prompt",
				Description: "Local task",
			},
		},
		Settings: domain.Settings{},
	}

	task, found, err := resolver.ResolveTask(context.Background(), "local-task", cfg, true)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !found {
		t.Fatal("Expected task to be found in local config")
	}

	if task.Prompt != "Local task prompt" {
		t.Errorf("Expected prompt %q, got %q", "Local task prompt", task.Prompt)
	}
}

func TestResolver_ResolveTask_FromCache(t *testing.T) {
	fs := newMockFS()
	cache := newMockCache()
	github := &mockGitHubClient{}
	configLoader := config.NewLoader(fs)

	resolver := NewResolver(fs, cache, github, configLoader)

	// Add task to cache
	taskTOML := `[task]
prompt = "Cached task prompt"
description = "Cached task"
`
	cache.data["tasks"] = map[string][]byte{
		"cached-task": []byte(taskTOML),
	}

	cfg := domain.Config{
		Tasks:    map[string]domain.Task{},
		Settings: domain.Settings{},
	}

	task, found, err := resolver.ResolveTask(context.Background(), "cached-task", cfg, true)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !found {
		t.Fatal("Expected task to be found in cache")
	}

	if task.Prompt != "Cached task prompt" {
		t.Errorf("Expected prompt %q, got %q", "Cached task prompt", task.Prompt)
	}
}

func TestResolver_ResolveTask_DownloadFromGitHub(t *testing.T) {
	fs := newMockFS()
	cache := newMockCache()
	configLoader := config.NewLoader(fs)

	// Setup mock GitHub client
	indexCSV := `type,category,name,description,tags,bin,sha,size,created,updated
tasks,development,github-task,GitHub task,dev;util,,abc123,100,2024-01-01T00:00:00Z,2024-01-01T00:00:00Z
`
	taskTOML := `[task]
prompt = "GitHub task prompt"
description = "GitHub task"
`
	github := &mockGitHubClient{
		indexData: []byte(indexCSV),
		assetData: []byte(taskTOML),
	}

	resolver := NewResolver(fs, cache, github, configLoader)

	cfg := domain.Config{
		Tasks:    map[string]domain.Task{},
		Settings: domain.Settings{},
	}

	task, found, err := resolver.ResolveTask(context.Background(), "github-task", cfg, true)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !found {
		t.Fatal("Expected task to be found from GitHub")
	}

	if task.Prompt != "GitHub task prompt" {
		t.Errorf("Expected prompt %q, got %q", "GitHub task prompt", task.Prompt)
	}

	// Verify task was cached
	cachedData, err := cache.Get("tasks", "github-task")
	if err != nil {
		t.Errorf("Expected task to be cached, got error: %v", err)
	}
	if string(cachedData) != taskTOML {
		t.Error("Cached data doesn't match downloaded data")
	}
}

func TestResolver_ResolveTask_DownloadNotAllowed(t *testing.T) {
	fs := newMockFS()
	cache := newMockCache()
	github := &mockGitHubClient{}
	configLoader := config.NewLoader(fs)

	resolver := NewResolver(fs, cache, github, configLoader)

	cfg := domain.Config{
		Tasks:    map[string]domain.Task{},
		Settings: domain.Settings{},
	}

	// downloadAllowed = false
	task, found, err := resolver.ResolveTask(context.Background(), "unknown-task", cfg, false)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if found {
		t.Error("Expected task not to be found when downloads not allowed")
	}

	if task.Name != "" {
		t.Error("Expected empty task when not found")
	}
}

func TestResolver_ResolveTask_PrecedenceLocalOverCache(t *testing.T) {
	fs := newMockFS()
	cache := newMockCache()
	github := &mockGitHubClient{}
	configLoader := config.NewLoader(fs)

	resolver := NewResolver(fs, cache, github, configLoader)

	// Same task in both local and cache
	cfg := domain.Config{
		Tasks: map[string]domain.Task{
			"task-name": {Prompt: "Local version"},
		},
		Settings: domain.Settings{},
	}

	taskTOML := `[task]
prompt = "Cached version"
`
	cache.data["tasks"] = map[string][]byte{
		"task-name": []byte(taskTOML),
	}

	task, found, err := resolver.ResolveTask(context.Background(), "task-name", cfg, true)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !found {
		t.Fatal("Expected task to be found")
	}

	// Should use local version
	if task.Prompt != "Local version" {
		t.Errorf("Expected local version, got %q", task.Prompt)
	}
}

func TestResolver_ResolveTask_CacheError(t *testing.T) {
	fs := newMockFS()
	cache := newMockCache()
	cache.getError = fmt.Errorf("cache read error")
	github := &mockGitHubClient{}
	configLoader := config.NewLoader(fs)

	resolver := NewResolver(fs, cache, github, configLoader)

	cfg := domain.Config{
		Tasks:    map[string]domain.Task{},
		Settings: domain.Settings{},
	}

	_, _, err := resolver.ResolveTask(context.Background(), "unknown-task", cfg, true)

	if err == nil {
		t.Error("Expected error when cache fails")
	}

	if !contains(err.Error(), "failed to check cache") {
		t.Errorf("Expected cache error message, got: %v", err)
	}
}

func TestResolver_ResolveTask_GitHubFetchError(t *testing.T) {
	fs := newMockFS()
	cache := newMockCache()
	github := &mockGitHubClient{
		indexErr: fmt.Errorf("network error"),
	}
	configLoader := config.NewLoader(fs)

	resolver := NewResolver(fs, cache, github, configLoader)

	cfg := domain.Config{
		Tasks:    map[string]domain.Task{},
		Settings: domain.Settings{},
	}

	_, _, err := resolver.ResolveTask(context.Background(), "unknown-task", cfg, true)

	if err == nil {
		t.Error("Expected error when GitHub fetch fails")
	}

	if !contains(err.Error(), "failed to fetch catalog index") {
		t.Errorf("Expected GitHub fetch error message, got: %v", err)
	}
}

func TestResolver_ResolveTask_TaskNotInCatalog(t *testing.T) {
	fs := newMockFS()
	cache := newMockCache()
	configLoader := config.NewLoader(fs)

	// Empty catalog
	indexCSV := `type,category,name,description,tags,bin,sha,size,created,updated
`
	github := &mockGitHubClient{
		indexData: []byte(indexCSV),
	}

	resolver := NewResolver(fs, cache, github, configLoader)

	cfg := domain.Config{
		Tasks:    map[string]domain.Task{},
		Settings: domain.Settings{},
	}

	_, found, err := resolver.ResolveTask(context.Background(), "nonexistent-task", cfg, true)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if found {
		t.Error("Expected task not to be found in empty catalog")
	}
}

func TestResolver_ResolveTask_InvalidCachedTOML(t *testing.T) {
	fs := newMockFS()
	cache := newMockCache()
	github := &mockGitHubClient{}
	configLoader := config.NewLoader(fs)

	resolver := NewResolver(fs, cache, github, configLoader)

	// Invalid TOML in cache
	cache.data["tasks"] = map[string][]byte{
		"bad-task": []byte("invalid toml [[["),
	}

	cfg := domain.Config{
		Tasks:    map[string]domain.Task{},
		Settings: domain.Settings{},
	}

	_, _, err := resolver.ResolveTask(context.Background(), "bad-task", cfg, true)

	if err == nil {
		t.Error("Expected error when parsing invalid cached TOML")
	}

	if !contains(err.Error(), "failed to parse cached task") {
		t.Errorf("Expected parse error message, got: %v", err)
	}
}

// Helper function (reused from other tests)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && indexOfSubstring(s, substr) >= 0))
}

func indexOfSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
