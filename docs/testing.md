# Testing Strategy

**Project:** start - AI Agent CLI
**Document:** Testing Strategy and Guidelines
**Last Updated:** 2025-11-24

---

## Table of Contents

1. [Overview](#overview)
2. [Testing Philosophy](#testing-philosophy)
3. [Smith Test Agent](#smith-test-agent)
4. [Unit Testing](#unit-testing)
5. [Integration Testing](#integration-testing)
6. [Test Coverage Goals](#test-coverage-goals)
7. [Mock Implementations](#mock-implementations)
8. [Running Tests](#running-tests)

---

## Overview

Three-tier testing approach:

1. **Unit Tests:** Fast, isolated, test individual components with mocks
2. **Integration Tests:** Use smith agent to verify end-to-end behavior
3. **Manual Testing:** After each phase, manual validation by developer

**Key Principles:**

- Test-first when possible
- Each component independently testable
- No external dependencies in unit tests (no network, no real filesystem)
- Integration tests use real binary + smith agent
- High coverage of critical paths

---

## Testing Philosophy

### What We Test

**Unit Test Coverage:**

- Config loading and merging
- TOML parsing and validation
- Placeholder resolution (all types)
- UTD pattern processing
- Prompt assembly
- Asset resolution algorithm
- Command construction

**Integration Test Coverage:**

- Full command execution (start, task, prompt, assets)
- Config file detection and merging
- Real binary behavior with smith agent
- Edge cases (missing files, invalid configs)

**Not Tested:**

- External services (GitHub API - mocked)
- Real AI agents (use smith instead)
- Terminal UI rendering (test output strings only)

### Why Smith?

**Problem:** Testing with real AI agents is:

- Expensive (API costs)
- Slow (network latency)
- Non-deterministic (different responses)
- Requires API keys
- Can't run offline

**Solution:** Smith test agent

- Deterministic (same input → same output)
- Fast (no network)
- Free (no API costs)
- Inspectable (writes outputs to files)
- Offline (works anywhere)

---

## Smith Test Agent

### Overview

Agent Smith is a minimal test double that captures what it receives from `start` without making external API calls.

**Named after:** Agent Smith from The Matrix (appropriate for a test agent that's everywhere when you need it)

### Location

`cmd/smith/main.go`

### Behavior

1. Check `SMITH_OUTPUT_DIR` environment variable
2. **If set (testing mode):**
   - Write all CLI arguments to `{DIR}/args.txt` (one per line)
   - Write last argument (the prompt) to `{DIR}/prompt.md`
   - Exit 0
3. **If not set (manual mode):**
   - Print last argument (prompt) to stdout
   - Exit 0

### Implementation

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

### Output Files

**args.txt** (one argument per line):

```
smith
--model
claude-3-7-sonnet-20250219
--append-system-prompt
You are a code reviewer...
Read ~/reference/ENVIRONMENT.md...
```

**prompt.md** (assembled prompt):

```markdown
Read ~/reference/ENVIRONMENT.md for environment context.
Read ~/reference/INDEX.csv for documentation index.
Read ./AGENTS.md for repository overview.

Review the following changes:

## Instructions
check security

## Staged Changes
```diff
...
```

```

### Building Smith

**Manual build:**
```bash
go build -o bin/smith cmd/smith/main.go
```

**Automatic build (in test.sh):**

```bash
if [ ! -f bin/smith ]; then
    echo "Building smith..."
    go build -o bin/smith cmd/smith/main.go
fi
```

### Using Smith Manually

**Test prompt assembly:**

```bash
# Build smith
go build -o bin/smith cmd/smith/main.go

# Run with output directory
mkdir -p /tmp/smith-test
SMITH_OUTPUT_DIR=/tmp/smith-test ./bin/smith arg1 arg2 "final prompt"

# Check outputs
cat /tmp/smith-test/args.txt
cat /tmp/smith-test/prompt.md
```

**Test with start:**

```bash
# Configure start to use smith
cat > ~/.config/start/agents.toml <<EOF
[agents.smith]
bin = "smith"
command = "{bin} --model {model} '{prompt}'"
default_model = "test"

  [agents.smith.models]
  test = "test-model"
EOF

# Run start with smith
SMITH_OUTPUT_DIR=/tmp/smith-test start "test prompt"

# Verify outputs
cat /tmp/smith-test/args.txt
cat /tmp/smith-test/prompt.md
```

---

## Unit Testing

### Approach

Test each component in isolation using mocked dependencies.

### Example: Config Loader

```go
// internal/config/loader_test.go
package config_test

import (
    "testing"

    "github.com/grantcarthew/start/internal/config"
    "github.com/grantcarthew/start/test/mocks"
    "github.com/stretchr/testify/assert"
)

func TestLoadGlobal(t *testing.T) {
    // Create mock filesystem
    mockFS := &mocks.MockFileSystem{
        Files: map[string]string{
            "/home/user/.config/start/config.toml": `
[settings]
default_agent = "claude"
`,
            "/home/user/.config/start/agents.toml": `
[agents.claude]
bin = "claude"
command = "claude --model {model} '{prompt}'"
default_model = "sonnet"

  [agents.claude.models]
  sonnet = "claude-3-7-sonnet-20250219"
`,
        },
    }

    loader := config.NewLoader(mockFS)

    cfg, err := loader.LoadGlobal()

    assert.NoError(t, err)
    assert.Equal(t, "claude", cfg.Settings.DefaultAgent)
    assert.Equal(t, "claude", cfg.Agents["claude"].Bin)
    assert.Equal(t, "sonnet", cfg.Agents["claude"].DefaultModel)
}

func TestLoadGlobal_MissingFiles(t *testing.T) {
    // Empty filesystem
    mockFS := &mocks.MockFileSystem{
        Files: map[string]string{},
    }

    loader := config.NewLoader(mockFS)

    cfg, err := loader.LoadGlobal()

    // Should not error, just return empty config
    assert.NoError(t, err)
    assert.Empty(t, cfg.Agents)
}
```

### Example: Prompt Assembly

```go
// internal/engine/prompt_test.go
package engine_test

import (
    "context"
    "testing"

    "github.com/grantcarthew/start/internal/domain"
    "github.com/grantcarthew/start/internal/engine"
    "github.com/grantcarthew/start/test/mocks"
    "github.com/stretchr/testify/assert"
)

func TestPromptAssembly_WithContexts(t *testing.T) {
    mockFS := &mocks.MockFileSystem{
        Files: map[string]string{
            "/home/user/.config/start/roles/reviewer.md": "You are a code reviewer",
            "/home/user/project/AGENTS.md":               "# Agent Instructions",
            "/home/user/reference/ENVIRONMENT.md":        "# Environment",
        },
    }

    mockRunner := &mocks.MockRunner{}

    promptEngine := engine.NewPromptEngine(mockFS, mockRunner)

    role := domain.Role{
        Name: "reviewer",
        File: "/home/user/.config/start/roles/reviewer.md",
    }

    contexts := []domain.Context{
        {
            Name:     "agents",
            File:     "/home/user/project/AGENTS.md",
            Prompt:   "Read {file} for agent instructions.",
            Required: true,
        },
        {
            Name:     "environment",
            File:     "/home/user/reference/ENVIRONMENT.md",
            Prompt:   "Read {file} for environment context.",
            Required: true,
        },
    }

    result, err := promptEngine.Assemble(context.Background(), role, contexts, "Review this code")

    assert.NoError(t, err)
    assert.Contains(t, result, "Read /home/user/project/AGENTS.md")
    assert.Contains(t, result, "Read /home/user/reference/ENVIRONMENT.md")
    assert.Contains(t, result, "Review this code")
}

func TestPromptAssembly_MissingFile(t *testing.T) {
    mockFS := &mocks.MockFileSystem{
        Files: map[string]string{},
    }

    mockRunner := &mocks.MockRunner{}

    promptEngine := engine.NewPromptEngine(mockFS, mockRunner)

    contexts := []domain.Context{
        {
            Name:   "missing",
            File:   "/nonexistent.md",
            Prompt: "Read {file}",
        },
    }

    result, err := promptEngine.Assemble(context.Background(), domain.Role{}, contexts, "test")

    // Should warn but continue
    assert.NoError(t, err)
    assert.Contains(t, result, "test")
}
```

### Example: Placeholder Resolution

```go
// internal/engine/placeholder_test.go
package engine_test

import (
    "testing"

    "github.com/grantcarthew/start/internal/engine"
    "github.com/stretchr/testify/assert"
)

func TestPlaceholderResolution_Basic(t *testing.T) {
    resolver := engine.NewPlaceholderResolver(nil, nil)

    template := "Model: {model}, Prompt: {prompt}"
    values := map[string]string{
        "model":  "claude-3-7-sonnet-20250219",
        "prompt": "Hello world",
    }

    result := resolver.Resolve(template, values)

    assert.Equal(t, "Model: claude-3-7-sonnet-20250219, Prompt: Hello world", result)
}

func TestPlaceholderResolution_MissingValue(t *testing.T) {
    resolver := engine.NewPlaceholderResolver(nil, nil)

    template := "Model: {model}, Missing: {missing}"
    values := map[string]string{
        "model": "claude-3-7-sonnet-20250219",
    }

    result := resolver.Resolve(template, values)

    // Missing placeholder left as-is (or replaced with empty string)
    assert.Contains(t, result, "Model: claude-3-7-sonnet-20250219")
}

func TestPlaceholderResolution_Date(t *testing.T) {
    resolver := engine.NewPlaceholderResolver(nil, nil)

    template := "Date: {date}"

    result := resolver.Resolve(template, map[string]string{})

    // Should contain ISO 8601 format date
    assert.Regexp(t, `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`, result)
}
```

### Test Organization

```
internal/
├── config/
│   ├── loader.go
│   ├── loader_test.go
│   ├── merge.go
│   ├── merge_test.go
│   ├── validator.go
│   └── validator_test.go
├── engine/
│   ├── prompt.go
│   ├── prompt_test.go
│   ├── placeholder.go
│   ├── placeholder_test.go
│   ├── utd.go
│   └── utd_test.go
```

**Naming Convention:**

- Test files: `*_test.go`
- Package: `<package>_test` (black box testing) or `<package>` (white box testing)
- Functions: `TestFunctionName_Scenario`

---

## Integration Testing

### Approach

Test the full system by running the real `start` binary with smith agent and verifying outputs.

### Test Structure

```
test/
├── fixtures/               # Test config files
│   ├── minimal/
│   │   ├── config.toml
│   │   └── agents.toml
│   ├── full/
│   │   ├── config.toml
│   │   ├── agents.toml
│   │   ├── roles.toml
│   │   ├── contexts.toml
│   │   └── tasks.toml
│   └── invalid/
│       └── config.toml
├── mocks/                  # Mock implementations
│   ├── fs.go
│   ├── runner.go
│   └── github.go
└── integration/            # Integration tests
    ├── task_test.go
    ├── prompt_test.go
    ├── assets_test.go
    └── helpers.go
```

### Example: Task Execution Test

```go
// test/integration/task_test.go
package integration

import (
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestTaskExecution(t *testing.T) {
    // Ensure binaries exist
    ensureBinaries(t)

    // Create temp directories
    configDir := t.TempDir()
    outputDir := t.TempDir()

    // Write test config
    writeTestConfig(t, configDir, testConfigFull)

    // Set environment
    env := []string{
        "HOME=" + configDir,
        "SMITH_OUTPUT_DIR=" + outputDir,
    }

    // Run: start task review "check security"
    cmd := exec.Command("./bin/start", "task", "review", "check security")
    cmd.Env = env
    output, err := cmd.CombinedOutput()

    require.NoError(t, err, "Command failed: %s", output)

    // Verify smith captured the correct args
    argsPath := filepath.Join(outputDir, "args.txt")
    args, err := os.ReadFile(argsPath)
    require.NoError(t, err, "Failed to read args.txt")

    argsLines := strings.Split(string(args), "\n")
    assert.Contains(t, argsLines, "smith")
    assert.Contains(t, argsLines, "--model")
    assert.Contains(t, argsLines, "test-model")

    // Verify prompt contains expected content
    promptPath := filepath.Join(outputDir, "prompt.md")
    prompt, err := os.ReadFile(promptPath)
    require.NoError(t, err, "Failed to read prompt.md")

    promptStr := string(prompt)
    assert.Contains(t, promptStr, "check security", "Instructions not in prompt")
    assert.Contains(t, promptStr, "ENVIRONMENT.md", "Context not in prompt")
}

func TestTaskExecution_WithAlias(t *testing.T) {
    ensureBinaries(t)

    configDir := t.TempDir()
    outputDir := t.TempDir()

    writeTestConfig(t, configDir, testConfigWithAlias)

    env := []string{
        "HOME=" + configDir,
        "SMITH_OUTPUT_DIR=" + outputDir,
    }

    // Use alias instead of full name
    cmd := exec.Command("./bin/start", "task", "cr", "check bugs")
    cmd.Env = env
    output, err := cmd.CombinedOutput()

    require.NoError(t, err, "Command failed: %s", output)

    // Verify it executed the same task
    promptPath := filepath.Join(outputDir, "prompt.md")
    prompt, err := os.ReadFile(promptPath)
    require.NoError(t, err)

    assert.Contains(t, string(prompt), "check bugs")
}

// Helpers

func ensureBinaries(t *testing.T) {
    t.Helper()

    if _, err := os.Stat("./bin/start"); os.IsNotExist(err) {
        t.Fatal("start binary not found. Run: go build -o bin/start cmd/start/main.go")
    }
    if _, err := os.Stat("./bin/smith"); os.IsNotExist(err) {
        t.Fatal("smith binary not found. Run: go build -o bin/smith cmd/smith/main.go")
    }
}

func writeTestConfig(t *testing.T, dir string, config string) {
    t.Helper()

    configPath := filepath.Join(dir, ".config", "start")
    err := os.MkdirAll(configPath, 0755)
    require.NoError(t, err)

    // Split config into separate files
    // (implementation details omitted for brevity)
}

const testConfigFull = `
[settings]
default_agent = "smith"

[agents.smith]
bin = "smith"
command = "{bin} --model {model} '{prompt}'"
default_model = "test"

  [agents.smith.models]
  test = "test-model"

[roles.reviewer]
prompt = "You are a code reviewer"

[contexts.environment]
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
required = true

[tasks.review]
alias = "cr"
description = "Code review"
role = "reviewer"
prompt = "Review: {instructions}"
`
```

### Test Script

```bash
#!/usr/bin/env bash
# test.sh
set -e

echo "=== Building binaries ==="

# Build smith if needed
if [ ! -f bin/smith ]; then
    echo "Building smith..."
    go build -o bin/smith cmd/smith/main.go
fi

# Build start
echo "Building start..."
go build -o bin/start cmd/start/main.go

echo ""
echo "=== Running unit tests ==="
go test -v -short ./...

echo ""
echo "=== Running integration tests ==="
go test -v ./test/integration/...

echo ""
echo "✓ All tests passed!"
```

**Usage:**

```bash
chmod +x test.sh
./test.sh
```

---

## Test Coverage Goals

### Overall Coverage

- **Target:** 80%+ for core logic
- **Minimum:** 70% overall

### Critical Paths (100% Coverage)

- Config loading and merging
- Placeholder resolution
- UTD pattern processing
- Asset resolution algorithm
- Agent command construction

### Lower Priority (60%+ Coverage)

- CLI command handlers (mostly integration tested)
- Error formatting
- Output rendering

### Excluded from Coverage

- `cmd/start/main.go` (DI wiring, tested via integration)
- `cmd/smith/main.go` (test tool)
- Adapters (thin wrappers, tested via integration)

### Measuring Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out

# Show coverage by package
go test -cover ./...
```

---

## Mock Implementations

### MockFileSystem

```go
// test/mocks/fs.go
package mocks

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

type MockFileSystem struct {
    Files map[string]string // path -> content
}

func (m *MockFileSystem) ReadFile(path string) ([]byte, error) {
    content, ok := m.Files[path]
    if !ok {
        return nil, os.ErrNotExist
    }
    return []byte(content), nil
}

func (m *MockFileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
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
    // Mock implementation (no-op)
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

### MockRunner

```go
// test/mocks/runner.go
package mocks

import (
    "context"
    "fmt"
    "time"
)

type MockRunner struct {
    Outputs map[string]string // command -> output
}

func (m *MockRunner) Run(ctx context.Context, shell, command string, timeout time.Duration) (string, string, error) {
    output, ok := m.Outputs[command]
    if !ok {
        return "", "", fmt.Errorf("command not found in mock: %s", command)
    }
    return output, "", nil
}
```

### MockGitHubClient

```go
// test/mocks/github.go
package mocks

import (
    "context"
    "fmt"
)

type MockGitHubClient struct {
    Index  []byte
    Assets map[string][]byte // path -> content
}

func (m *MockGitHubClient) FetchIndex(ctx context.Context, repo, branch string) ([]byte, error) {
    if m.Index == nil {
        return nil, fmt.Errorf("no index configured")
    }
    return m.Index, nil
}

func (m *MockGitHubClient) FetchAsset(ctx context.Context, repo, branch, path string) ([]byte, error) {
    content, ok := m.Assets[path]
    if !ok {
        return nil, fmt.Errorf("asset not found: %s", path)
    }
    return content, nil
}
```

---

## Running Tests

### All Tests

```bash
./test.sh
```

### Unit Tests Only

```bash
go test -short ./...
```

### Integration Tests Only

```bash
go test ./test/integration/...
```

### Specific Package

```bash
go test -v ./internal/config/
```

### Specific Test

```bash
go test -v -run TestPromptAssembly ./internal/engine/
```

### With Coverage

```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Verbose Output

```bash
go test -v ./...
```

### Parallel Execution

```bash
go test -parallel 4 ./...
```

### Watch Mode (with external tool)

```bash
# Install gotestsum
go install gotest.tools/gotestsum@latest

# Watch and run tests on file changes
gotestsum --watch
```

---

_Document Status: Complete_
_Last Updated: 2025-11-24_
