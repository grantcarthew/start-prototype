# Architecture

**Project:** start - AI Agent CLI
**Document:** Architecture Deep Dive
**Last Updated:** 2025-11-24

---

## Table of Contents

1. [Overview](#overview)
2. [Directory Structure](#directory-structure)
3. [Hexagonal Architecture](#hexagonal-architecture)
4. [Domain Layer](#domain-layer)
5. [Engine Layer](#engine-layer)
6. [Adapters Layer](#adapters-layer)
7. [CLI Layer](#cli-layer)
8. [Dependency Injection](#dependency-injection)
9. [Design Patterns](#design-patterns)

---

## Overview

The `start` CLI follows **Hexagonal Architecture** (also called Ports and Adapters) to achieve:

- **Testability:** Core logic isolated from external dependencies
- **Flexibility:** Easy to swap implementations (real vs mock)
- **Maintainability:** Clear separation of concerns
- **Idiomatic Go:** Standard Go patterns throughout

**Key Principles:**

- Business logic depends on interfaces, not concrete implementations
- External dependencies (filesystem, HTTP, exec) are adapters
- Core domain is pure Go with no external dependencies
- Dependency injection wired in main.go

---

## Directory Structure

```
start/
├── cmd/
│   ├── start/              # Main entry point
│   │   └── main.go         # DI wiring, version injection
│   └── smith/              # Test agent
│       └── main.go         # Captures args/prompt for testing
│
├── internal/               # All implementation (no pkg/)
│   ├── domain/             # Core domain (pure Go)
│   │   ├── models.go       # Data structures
│   │   └── interfaces.go   # Contracts
│   │
│   ├── config/             # TOML configuration
│   │   ├── loader.go       # Load global + local
│   │   ├── merge.go        # Merge strategy
│   │   └── validator.go    # Validation rules
│   │
│   ├── engine/             # Business logic
│   │   ├── prompt.go       # Prompt assembly
│   │   ├── placeholder.go  # Placeholder resolution
│   │   ├── executor.go     # Agent execution
│   │   └── utd.go          # UTD processing
│   │
│   ├── assets/             # Asset management
│   │   ├── catalog.go      # GitHub catalog
│   │   ├── cache.go        # Local cache
│   │   └── resolver.go     # Resolution algorithm
│   │
│   ├── cli/                # Cobra commands
│   │   ├── root.go
│   │   ├── init.go
│   │   ├── task.go
│   │   ├── prompt.go
│   │   ├── assets.go
│   │   ├── config.go
│   │   └── doctor.go
│   │
│   └── adapters/           # Concrete implementations
│       ├── fs.go           # RealFileSystem
│       ├── exec.go         # RealRunner
│       └── github.go       # RealGitHubClient
│
├── test/                   # Integration tests
│   ├── fixtures/           # Test data
│   ├── mocks/              # Mock implementations
│   └── integration/        # End-to-end tests
│
├── bin/                    # Compiled binaries (gitignored)
├── go.mod
└── go.sum
```

**Package Responsibilities:**

| Package | Purpose | Dependencies |
|---------|---------|--------------|
| `domain/` | Core types and contracts | None (pure Go) |
| `config/` | TOML loading and merging | domain, go-toml |
| `engine/` | Business logic | domain interfaces |
| `assets/` | Asset management | domain interfaces |
| `cli/` | User interface | engine, config |
| `adapters/` | External integrations | os, net/http, os/exec |

---

## Hexagonal Architecture

```
┌─────────────────────────────────────────┐
│        CLI Layer (Cobra)                │
│  - Parse flags and arguments            │
│  - Route to command handlers            │
│  - Format output for user               │
│  - Minimal business logic               │
└─────────────────┬───────────────────────┘
                  │
                  │ Calls
                  ▼
┌─────────────────────────────────────────┐
│       Engine Layer (Business Logic)     │
│  - Prompt assembly                      │
│  - Placeholder resolution               │
│  - UTD processing                       │
│  - Asset resolution                     │
│  - Configuration merging                │
└─────────────────┬───────────────────────┘
                  │
                  │ Uses interfaces from
                  ▼
┌─────────────────────────────────────────┐
│      Domain Layer (Interfaces + Models) │
│  - FileSystem interface                 │
│  - Runner interface                     │
│  - GitHubClient interface               │
│  - Cache interface                      │
│  - Pure domain models (structs)         │
└─────────────────┬───────────────────────┘
                  │
                  │ Implemented by
                  ▼
┌─────────────────────────────────────────┐
│       Adapters Layer (Concrete Impls)   │
│  - RealFileSystem (uses os package)     │
│  - RealRunner (uses os/exec)            │
│  - RealGitHubClient (uses net/http)     │
│  - MockFS/MockRunner for tests          │
└─────────────────────────────────────────┘
```

**Flow Example: `start task review "check security"`**

1. **CLI Layer:** Parse command, extract task name and instructions
2. **Engine Layer:** Load config, resolve task, assemble prompt
3. **Domain Layer:** Interfaces used (FileSystem, Runner)
4. **Adapters Layer:** RealFileSystem reads files, RealRunner executes agent
5. **CLI Layer:** Format output, display to user

---

## Domain Layer

**Location:** `internal/domain/`

**Purpose:** Define core business entities and contracts with zero external dependencies

### Models (`models.go`)

```go
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
    DefaultAgent    string
    DefaultRole     string
    LogLevel        string
    Shell           string
    CommandTimeout  int
    AssetDownload   bool
    AssetRepo       string
    AssetPath       string
}

// Agent from agents.toml [agents.<name>]
type Agent struct {
    Name         string
    Bin          string
    Command      string
    Description  string
    URL          string
    ModelsURL    string
    DefaultModel string
    Models       map[string]string
}

// Role from roles.toml [roles.<name>] (UTD pattern)
type Role struct {
    Name           string
    Description    string
    File           string
    Command        string
    Prompt         string
    Shell          string
    CommandTimeout int
}

// Context from contexts.toml [contexts.<name>] (UTD pattern)
type Context struct {
    Name           string
    Description    string
    File           string
    Command        string
    Prompt         string
    Required       bool
    Shell          string
    CommandTimeout int
}

// Task from tasks.toml [tasks.<name>] (UTD pattern)
type Task struct {
    Name           string
    Alias          string
    Description    string
    Role           string
    Agent          string
    File           string
    Command        string
    Prompt         string
    Shell          string
    CommandTimeout int
}

// AssetMeta from .meta.toml files
type AssetMeta struct {
    Type        string
    Category    string
    Name        string
    Description string
    Tags        string
    Bin         string
    SHA         string
    Size        int64
    Created     time.Time
    Updated     time.Time
}
```

### Interfaces (`interfaces.go`)

```go
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

**Why interfaces?**

- Enables testing with mocks
- Allows different implementations (in-memory, real filesystem, etc.)
- Clear contracts for what operations are needed
- Decouples business logic from infrastructure

---

## Engine Layer

**Location:** `internal/engine/`

**Purpose:** Core business logic that uses domain interfaces

### Prompt Engine (`prompt.go`)

Assembles prompts from contexts, roles, and custom text.

**Responsibilities:**

- Load and process context documents
- Load and process roles
- Resolve placeholders
- Assemble final prompt in correct order
- Handle missing files gracefully

**Interface:**

```go
type PromptEngine struct {
    fs     FileSystem
    runner Runner
}

func NewPromptEngine(fs FileSystem, runner Runner) *PromptEngine

func (e *PromptEngine) Assemble(ctx context.Context, cfg Config, role Role, contexts []Context, customPrompt string) (string, error)
```

### Placeholder Resolver (`placeholder.go`)

Resolves all placeholder types with proper scoping.

**Placeholders:**

- Universal: `{date}`
- Agent commands: `{bin}`, `{model}`, `{prompt}`, `{role}`, `{role_file}`
- UTD pattern: `{file}`, `{file_contents}`, `{command}`, `{command_output}`
- Tasks: `{instructions}`

**Interface:**

```go
type PlaceholderResolver struct {
    fs     FileSystem
    runner Runner
}

func NewPlaceholderResolver(fs FileSystem, runner Runner) *PlaceholderResolver

func (r *PlaceholderResolver) Resolve(template string, values map[string]string) (string, error)
```

### UTD Processor (`utd.go`)

Processes Unified Template Design pattern (file, command, prompt).

**Responsibilities:**

- Execute commands if specified
- Read files if specified
- Resolve placeholders in prompt template
- Return final assembled content

**Interface:**

```go
type UTDProcessor struct {
    fs     FileSystem
    runner Runner
}

func NewUTDProcessor(fs FileSystem, runner Runner) *UTDProcessor

func (p *UTDProcessor) Process(ctx context.Context, file, command, prompt, shell string, timeout int) (string, error)
```

### Executor (`executor.go`)

Executes agent commands with resolved placeholders.

**Responsibilities:**

- Construct agent command from template
- Resolve all placeholders
- Create temp files if needed (`{role_file}`)
- Execute command via Runner
- Clean up temp files
- Stream output to stdout/stderr

**Interface:**

```go
type Executor struct {
    runner Runner
    fs     FileSystem
}

func NewExecutor(runner Runner, fs FileSystem) *Executor

func (e *Executor) Execute(ctx context.Context, agent Agent, model, role, roleFile, prompt string) error
```

---

## Adapters Layer

**Location:** `internal/adapters/`

**Purpose:** Concrete implementations of domain interfaces

### RealFileSystem (`fs.go`)

Implements FileSystem interface using Go's os package.

```go
type RealFileSystem struct{}

func (fs *RealFileSystem) ReadFile(path string) ([]byte, error) {
    expanded := expandPath(path) // Handle ~ expansion
    return os.ReadFile(expanded)
}

func (fs *RealFileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
    expanded := expandPath(path)
    return os.WriteFile(expanded, data, perm)
}

func (fs *RealFileSystem) Exists(path string) bool {
    expanded := expandPath(path)
    _, err := os.Stat(expanded)
    return err == nil
}

// ... other methods

func expandPath(path string) string {
    if strings.HasPrefix(path, "~/") {
        home, _ := os.UserHomeDir()
        return filepath.Join(home, path[2:])
    }
    return path
}
```

### RealRunner (`exec.go`)

Implements Runner interface using os/exec.

```go
type RealRunner struct{}

func (r *RealRunner) Run(ctx context.Context, shell, command string, timeout time.Duration) (string, string, error) {
    ctx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()

    cmd := exec.CommandContext(ctx, shell, "-c", command)

    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    err := cmd.Run()
    return stdout.String(), stderr.String(), err
}
```

### RealGitHubClient (`github.go`)

Implements GitHubClient interface using net/http.

```go
type RealGitHubClient struct {
    client *http.Client
}

func NewRealGitHubClient() *RealGitHubClient {
    return &RealGitHubClient{
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

func (c *RealGitHubClient) FetchIndex(ctx context.Context, repo, branch string) ([]byte, error) {
    url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/assets/index.csv", repo, branch)
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}

func (c *RealGitHubClient) FetchAsset(ctx context.Context, repo, branch, path string) ([]byte, error) {
    url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", repo, branch, path)
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}
```

### FileCache (`cache.go`)

Implements Cache interface using filesystem.

**Structure:**

```
~/.config/start/assets/
├── agents/
│   └── anthropic/
│       ├── claude.toml
│       └── claude.meta.toml
├── roles/
├── tasks/
└── contexts/
```

```go
type FileCache struct {
    FS   FileSystem
    Base string // ~/.config/start/assets
}

func (c *FileCache) Get(assetType, name string) ([]byte, error) {
    pattern := filepath.Join(c.Base, assetType, "*", name+".toml")
    matches, err := c.FS.Glob(pattern)
    if err != nil || len(matches) == 0 {
        return nil, os.ErrNotExist
    }
    return c.FS.ReadFile(matches[0])
}

func (c *FileCache) Set(assetType, name string, content []byte, meta AssetMeta) error {
    dir := filepath.Join(c.Base, assetType, meta.Category)
    c.FS.MkdirAll(dir, 0755)

    // Write asset content
    assetPath := filepath.Join(dir, name+".toml")
    c.FS.WriteFile(assetPath, content, 0644)

    // Write metadata
    metaPath := filepath.Join(dir, name+".meta.toml")
    metaBytes := marshalMeta(meta)
    c.FS.WriteFile(metaPath, metaBytes, 0644)

    return nil
}
```

---

## CLI Layer

**Location:** `internal/cli/`

**Purpose:** User interface using Cobra framework

### Root Command (`root.go`)

Main entry point, handles `start` with no subcommand.

**Responsibilities:**

- Parse global flags
- Load config (global + local)
- Select agent, role, model
- Assemble prompt
- Execute agent

**Structure:**

```go
type RootCommand struct {
    configLoader  *config.Loader
    promptEngine  *engine.PromptEngine
    executor      *engine.Executor
    assetResolver *assets.Resolver
    version       string
}

func NewRootCommand(...) *cobra.Command {
    rc := &RootCommand{...}

    cmd := &cobra.Command{
        Use:   "start",
        Short: "Launch AI agent with context",
        RunE:  rc.run,
    }

    // Add flags
    cmd.PersistentFlags().StringP("agent", "a", "", "Agent to use")
    cmd.PersistentFlags().StringP("role", "r", "", "Role to use")
    cmd.PersistentFlags().StringP("model", "m", "", "Model to use")
    // ... more flags

    // Add subcommands
    cmd.AddCommand(NewInitCommand(...))
    cmd.AddCommand(NewTaskCommand(...))
    cmd.AddCommand(NewPromptCommand(...))
    cmd.AddCommand(NewAssetsCommand(...))
    cmd.AddCommand(NewConfigCommand(...))
    cmd.AddCommand(NewDoctorCommand(...))

    return cmd
}

func (rc *RootCommand) run(cmd *cobra.Command, args []string) error {
    // Implementation
}
```

### Other Commands

Each subcommand follows similar pattern:

- Thin layer over engine logic
- Handle user input/output
- Format errors for display
- No business logic

**Commands:**

- `init.go` - `start init`
- `task.go` - `start task <name>`
- `prompt.go` - `start prompt <text>`
- `assets.go` - `start assets` (with subcommands)
- `config.go` - `start config` (with subcommands)
- `doctor.go` - `start doctor`

---

## Dependency Injection

**Location:** `cmd/start/main.go`

**Purpose:** Wire up all dependencies at application startup

```go
package main

import (
    "os"

    "github.com/grantcarthew/start/internal/adapters"
    "github.com/grantcarthew/start/internal/assets"
    "github.com/grantcarthew/start/internal/cli"
    "github.com/grantcarthew/start/internal/config"
    "github.com/grantcarthew/start/internal/engine"
)

var version = "dev" // Injected at build time via -ldflags

func main() {
    // Create adapters (real implementations)
    fs := &adapters.RealFileSystem{}
    runner := &adapters.RealRunner{}
    githubClient := adapters.NewRealGitHubClient()

    // Create cache
    cacheBase := expandPath("~/.config/start/assets")
    cache := &adapters.FileCache{
        FS:   fs,
        Base: cacheBase,
    }

    // Create config loader
    configLoader := config.NewLoader(fs)

    // Create engines
    utdProcessor := engine.NewUTDProcessor(fs, runner)
    placeholderResolver := engine.NewPlaceholderResolver(fs, runner)
    promptEngine := engine.NewPromptEngine(fs, runner, utdProcessor, placeholderResolver)
    executor := engine.NewExecutor(runner, fs)

    // Create asset resolver
    assetResolver := assets.NewResolver(fs, cache, githubClient, configLoader)

    // Create root command with all dependencies
    rootCmd := cli.NewRootCommand(
        configLoader,
        promptEngine,
        executor,
        assetResolver,
        version,
    )

    // Execute
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}

func expandPath(path string) string {
    if strings.HasPrefix(path, "~/") {
        home, _ := os.UserHomeDir()
        return filepath.Join(home, path[2:])
    }
    return path
}
```

**Benefits:**

- All dependencies created in one place
- Easy to swap for testing (create with mocks instead)
- Clear dependency graph
- Single source of truth for construction

---

## Design Patterns

### Constructor Pattern

All components use constructor functions:

```go
func NewPromptEngine(fs FileSystem, runner Runner) *PromptEngine {
    return &PromptEngine{
        fs:     fs,
        runner: runner,
    }
}
```

**Benefits:**

- Clear dependencies
- Easy to test (pass mocks)
- Explicit initialization

### Interface Segregation

Interfaces are small and focused:

```go
// Good: Small, focused interface
type FileReader interface {
    ReadFile(path string) ([]byte, error)
}

// Avoid: Large, monolithic interface
type EverythingDoer interface {
    ReadFile(path string) ([]byte, error)
    WriteFile(path string, data []byte) error
    ExecuteCommand(cmd string) (string, error)
    FetchFromGitHub(url string) ([]byte, error)
    // ...
}
```

### Error Wrapping

Errors are wrapped with context:

```go
if err != nil {
    return fmt.Errorf("failed to load config: %w", err)
}
```

**Benefits:**

- Error chain preserved
- Context added at each layer
- Can unwrap with errors.Is/As

### Context Propagation

All long-running operations accept context:

```go
func (e *Executor) Execute(ctx context.Context, ...) error {
    // Can be canceled
    // Has timeout
    // Carries request-scoped values
}
```

---

## Testing Architecture

See [testing.md](testing.md) for complete testing strategy.

**Key Points:**

- Unit tests use mock implementations (MockFileSystem, MockRunner)
- Integration tests use real `start` binary with smith agent
- All domain logic is testable without real filesystem/network
- Adapters have minimal logic (thin wrappers)

---

_Document Status: Complete_
_Last Updated: 2025-11-24_
