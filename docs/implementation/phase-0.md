# Phase 0: Foundation & Smith

**Status:** Not Started
**Dependencies:** None
**Estimated Effort:** 2-3 hours

---

## Required Reading

Before starting this phase, review these documents:

**Architecture & Design:**
- [docs/architecture.md](../architecture.md) - Hexagonal architecture pattern, dependency injection
- [docs/testing.md](../testing.md) - Testing strategy, smith agent concept

**Design Records:**
- [DR-001: TOML Format](../design/design-records/dr-001-toml-format.md) - Understanding TOML tags for structs
- [DR-007: Placeholders](../design/design-records/dr-007-placeholders.md) - Placeholder system context

**Reference:**
- [Design Records Index](../design/design-records/README.md) - Overview of all design decisions
- [examples/minimal/](../../examples/minimal/) - Target configuration structure for Phase 0

---

## Goal

Build project scaffolding and testing infrastructure. This phase establishes the foundation for all future development.

---

## Deliverables

- [ ] Project structure created (directories, go.mod)
- [ ] Smith agent implemented and tested
- [ ] Domain models defined (structs only, no logic)
- [ ] Test harness setup (test.sh script)
- [ ] Build succeeds, tests run (even if empty)

---

## Tasks

### 1. Project Initialization

**Create Go module:**

```bash
go mod init github.com/grantcarthew/start
```

**Create directory structure:**

```bash
mkdir -p cmd/start cmd/smith
mkdir -p internal/{domain,config,engine,assets,cli,adapters}
mkdir -p test/{fixtures,mocks,integration}
mkdir -p bin
mkdir -p docs/implementation
```

**Create .gitignore:**

```
# Binaries
bin/

# Test outputs
test/output/
*.out
coverage.out

# Go
*.exe
*.test
*.prof

# IDE
.vscode/
.idea/
*.swp
```

### 2. Domain Models

**Create `internal/domain/models.go`:**

Define all structs with TOML tags. See [examples/minimal/global/](../../examples/minimal/global/) for the target configuration structure these models will load.

```go
package domain

import "time"

// Config represents the merged configuration
type Config struct {
    Settings Settings
    Agents   map[string]Agent
    Roles    map[string]Role
    Contexts map[string]Context
    Tasks    map[string]Task
}

// Settings from config.toml [settings]
type Settings struct {
    DefaultAgent    string `toml:"default_agent"`
    DefaultRole     string `toml:"default_role"`
    LogLevel        string `toml:"log_level"`
    Shell           string `toml:"shell"`
    CommandTimeout  int    `toml:"command_timeout"`
    AssetDownload   bool   `toml:"asset_download"`
    AssetRepo       string `toml:"asset_repo"`
    AssetPath       string `toml:"asset_path"`
}

// Agent from agents.toml [agents.<name>]
type Agent struct {
    Name         string
    Bin          string            `toml:"bin"`
    Command      string            `toml:"command"`
    Description  string            `toml:"description"`
    URL          string            `toml:"url"`
    ModelsURL    string            `toml:"models_url"`
    DefaultModel string            `toml:"default_model"`
    Models       map[string]string `toml:"models"`
}

// Role from roles.toml [roles.<name>] (UTD pattern)
type Role struct {
    Name           string
    Description    string `toml:"description"`
    File           string `toml:"file"`
    Command        string `toml:"command"`
    Prompt         string `toml:"prompt"`
    Shell          string `toml:"shell"`
    CommandTimeout int    `toml:"command_timeout"`
}

// Context from contexts.toml [contexts.<name>] (UTD pattern)
type Context struct {
    Name           string
    Description    string `toml:"description"`
    File           string `toml:"file"`
    Command        string `toml:"command"`
    Prompt         string `toml:"prompt"`
    Required       bool   `toml:"required"`
    Shell          string `toml:"shell"`
    CommandTimeout int    `toml:"command_timeout"`
}

// Task from tasks.toml [tasks.<name>] (UTD pattern)
type Task struct {
    Name           string
    Alias          string `toml:"alias"`
    Description    string `toml:"description"`
    Role           string `toml:"role"`
    Agent          string `toml:"agent"`
    File           string `toml:"file"`
    Command        string `toml:"command"`
    Prompt         string `toml:"prompt"`
    Shell          string `toml:"shell"`
    CommandTimeout int    `toml:"command_timeout"`
}

// AssetMeta from .meta.toml files
type AssetMeta struct {
    Type        string    `toml:"type"`
    Category    string    `toml:"category"`
    Name        string    `toml:"name"`
    Description string    `toml:"description"`
    Tags        string    `toml:"tags"`
    Bin         string    `toml:"bin"`
    SHA         string    `toml:"sha"`
    Size        int64     `toml:"size"`
    Created     time.Time `toml:"created"`
    Updated     time.Time `toml:"updated"`
}

// CachedAsset represents an asset in the cache
type CachedAsset struct {
    Type     string
    Category string
    Name     string
    Meta     AssetMeta
}
```

**Notes:**

- All fields are exported (capitalized)
- TOML tags for unmarshaling
- No methods yet, just data structures
- Comments describe where each type comes from

### 3. Domain Interfaces

**Create `internal/domain/interfaces.go`:**

```go
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
```

**Notes:**

- Interface-first design
- All methods return error for failure cases
- Context passed for cancellation
- Small, focused interfaces

### 4. Smith Agent

**Create `cmd/smith/main.go`:**

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

func main() {
    outputDir := os.Getenv("SMITH_OUTPUT_DIR")

    if outputDir == "" {
        // Manual mode: print prompt to stdout
        if len(os.Args) > 1 {
            fmt.Println(os.Args[len(os.Args)-1])
        }
        os.Exit(0)
    }

    // Testing mode: write args and prompt to files
    argsFile := filepath.Join(outputDir, "args.txt")
    promptFile := filepath.Join(outputDir, "prompt.md")

    // Write args (one per line)
    argsContent := strings.Join(os.Args, "\n")
    if err := os.WriteFile(argsFile, []byte(argsContent), 0644); err != nil {
        fmt.Fprintf(os.Stderr, "Error writing args: %v\n", err)
        os.Exit(1)
    }

    // Write prompt (last arg)
    if len(os.Args) > 1 {
        prompt := os.Args[len(os.Args)-1]
        if err := os.WriteFile(promptFile, []byte(prompt), 0644); err != nil {
            fmt.Fprintf(os.Stderr, "Error writing prompt: %v\n", err)
            os.Exit(1)
        }
    }

    os.Exit(0)
}
```

**Build and test:**

```bash
# Build
go build -o bin/smith cmd/smith/main.go

# Test manual mode
./bin/smith arg1 arg2 "test prompt"
# Should print: test prompt

# Test testing mode
mkdir -p /tmp/smith-test
SMITH_OUTPUT_DIR=/tmp/smith-test ./bin/smith arg1 arg2 "test prompt"
cat /tmp/smith-test/args.txt
cat /tmp/smith-test/prompt.md
```

### 5. Test Harness

**Create `test.sh`:**

```bash
#!/usr/bin/env bash
set -e

echo "=== Building binaries ==="

# Build smith if needed
if [ ! -f bin/smith ]; then
    echo "Building smith..."
    go build -o bin/smith cmd/smith/main.go
fi

# Build start (when it exists)
if [ -f cmd/start/main.go ]; then
    echo "Building start..."
    go build -o bin/start cmd/start/main.go
fi

echo ""
echo "=== Running unit tests ==="
go test -v -short ./... || true

echo ""
echo "=== Running integration tests ==="
go test -v ./test/integration/... || true

echo ""
echo "✓ Tests complete!"
```

**Make executable:**

```bash
chmod +x test.sh
```

### 6. Mock Implementations

**Create `test/mocks/fs.go`:**

```go
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
```

**Create `test/mocks/runner.go`:**

```go
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
```

**Create `test/mocks/github.go`:**

```go
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
```

### 7. Verify Everything Compiles

**Run:**

```bash
go mod tidy
go build ./...
./test.sh
```

**Expected output:**

- All packages compile
- Smith builds successfully
- Test script runs (even with no tests yet)

---

## Testing Criteria

- [ ] `go mod tidy` succeeds
- [ ] `go build ./...` compiles all packages
- [ ] `go build -o bin/smith cmd/smith/main.go` succeeds
- [ ] Smith writes correct output files when `SMITH_OUTPUT_DIR` set
- [ ] Smith prints to stdout when `SMITH_OUTPUT_DIR` not set
- [ ] `./test.sh` runs without errors
- [ ] All domain models compile
- [ ] All interfaces compile
- [ ] Mock implementations compile

---

## Manual Testing

### Test Smith Agent

```bash
# Build smith
go build -o bin/smith cmd/smith/main.go

# Test manual mode (no SMITH_OUTPUT_DIR)
./bin/smith "Hello world"
# Expected: prints "Hello world" to stdout

# Test with arguments
./bin/smith --model test-model --flag value "Final prompt"
# Expected: prints "Final prompt" to stdout

# Test testing mode
mkdir -p /tmp/smith-output
SMITH_OUTPUT_DIR=/tmp/smith-output ./bin/smith arg1 arg2 "test prompt"

# Verify args.txt
cat /tmp/smith-output/args.txt
# Expected:
# smith
# arg1
# arg2
# test prompt

# Verify prompt.md
cat /tmp/smith-output/prompt.md
# Expected: test prompt

# Cleanup
rm -rf /tmp/smith-output
```

### Verify Project Structure

```bash
tree -L 3 -I 'bin'
```

**Expected output:**

```
.
├── cmd
│   ├── smith
│   │   └── main.go
│   └── start
├── docs
│   └── implementation
│       ├── phase-0.md
│       └── ...
├── internal
│   ├── adapters
│   ├── assets
│   ├── cli
│   ├── config
│   ├── domain
│   │   ├── interfaces.go
│   │   └── models.go
│   └── engine
├── test
│   ├── fixtures
│   ├── integration
│   └── mocks
│       ├── fs.go
│       ├── github.go
│       └── runner.go
├── go.mod
├── go.sum
└── test.sh
```

---

## Completion Checklist

- [ ] All directories created
- [ ] go.mod initialized
- [ ] Domain models defined in `internal/domain/models.go`
- [ ] Domain interfaces defined in `internal/domain/interfaces.go`
- [ ] Smith agent implemented in `cmd/smith/main.go`
- [ ] Smith agent tested manually
- [ ] Mock implementations created
- [ ] `test.sh` script created and executable
- [ ] `.gitignore` created
- [ ] All code compiles without errors
- [ ] Test script runs successfully

---

## Next Phase

Once this phase is complete, proceed to [Phase 1: Config Loading & Validation](phase-1.md).

---

_Phase Status: Not Started_
_Last Updated: 2025-11-24_
