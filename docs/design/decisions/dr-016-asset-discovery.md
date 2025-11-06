# DR-016: Asset Discovery and Directory Structure

**Date:** 2025-01-06
**Status:** Accepted
**Category:** Asset Management

## Decision

No central asset discovery system; each feature checks its own directory with graceful fallbacks

## Asset Directory Structure

```
~/.config/start/
├── config.toml
├── asset-version.toml
└── assets/
    ├── agents/              # Agent configuration templates
    │   ├── claude.toml
    │   ├── gemini.toml
    │   └── aichat.toml
    ├── roles/               # System prompt markdown files
    │   ├── code-reviewer.md
    │   ├── doc-reviewer.md
    │   └── security-reviewer.md
    ├── tasks/               # Default task definitions
    │   ├── code-review.toml
    │   ├── git-diff-review.toml
    │   └── doc-review.toml
    └── examples/            # Example configuration files
        ├── global-config.toml
        └── local-config.toml
```

## Usage Pattern Per Asset Type

**Agents (`assets/agents/`):**
- Accessed by: `start config agent add`
- Usage: Load as templates, show selection menu
- Fallback: Manual agent entry if directory empty/missing

**Roles (`assets/roles/`):**
- Accessed by: User config references: `file = "~/.config/start/assets/roles/code-reviewer.md"`
- Usage: Read file contents when loading system prompt
- Fallback: Standard file-not-found error

**Tasks (`assets/tasks/`):**
- Accessed by: `start config task list`
- Usage: Show as templates (not auto-loaded)
- Fallback: Empty list if directory missing

**Examples (`assets/examples/`):**
- Accessed by: Never automatically (reference only)
- Usage: Users manually view/copy sections
- Fallback: N/A (optional resource)

## No Discovery System Needed

Each command handles its own asset directory:

```go
// Constants for asset directories (no magic strings)
const (
    AssetDirAgents   = "agents"
    AssetDirRoles    = "roles"
    AssetDirTasks    = "tasks"
    AssetDirExamples = "examples"
)

// Each feature loads independently
func LoadAgentTemplates() []AgentTemplate {
    dir := filepath.Join(assetDir, AssetDirAgents)
    files, err := filepath.Glob(filepath.Join(dir, "*.toml"))
    if err != nil || len(files) == 0 {
        return []AgentTemplate{} // Graceful fallback
    }
    return parseTemplates(files)
}
```

## Benefits

- ✅ **Simple:** No central registry or manifest
- ✅ **Decoupled:** Each feature self-contained
- ✅ **Graceful:** Missing directories don't break functionality
- ✅ **Maintainable:** Constants prevent magic strings
- ✅ **Extensible:** Add new directory + code that uses it

## Rationale

Asset types are few and stable. Each type requires specific behavior (templates vs files vs merged config). Simple directory checks with fallbacks provide clean implementation without unnecessary abstraction.

## Related Decisions

- [DR-011](./dr-011-asset-distribution.md) - Asset distribution strategy
- [DR-019](./dr-019-task-loading.md) - Task templates (not auto-loaded)
