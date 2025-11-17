# DR-019: Task Loading and Merging Algorithm

- Date: 2025-01-06
- Status: Updated by DR-031 and DR-033
- Category: Tasks

Note: This DR remains valid but has been extended by DR-031 (catalog architecture) and DR-033 (asset resolution algorithm including cache and GitHub catalog).

## Problem

Task loading needs clear rules for multiple sources and conflict resolution. The system must:

- Support global tasks (shared across all projects)
- Support local tasks (project-specific)
- Support asset catalog tasks (curated, downloadable)
- Define clear precedence when same task name appears in multiple sources
- Handle alias conflicts across scopes
- Provide transparency about task sources (security consideration)
- Allow local projects to override global tasks
- Enable task discovery without cluttering runtime execution
- Follow consistent replacement patterns with other config types (roles, agents)
- Support on-demand downloading from catalog

## Decision

Tasks load from multiple sources with clear precedence: local config → global config → asset cache → GitHub catalog (per DR-033). Local completely replaces global for same task name (no field merging). Asset tasks available via catalog resolution.

Task sources:

Runtime execution (`start task <name>`):
- Local config: `./.start/tasks.toml`
- Global config: `~/.config/start/tasks.toml`
- Asset cache: `~/.config/start/assets/tasks/`
- GitHub catalog: Query and download if `asset_download = true`

Management commands (`start config task list/add/edit`):
- Shows all configured tasks (local + global merged)
- Asset tasks visible via `start assets search` or `start assets add`

## Why

Clear precedence hierarchy:

- Local config wins (project-specific customization)
- Global config next (user's personal defaults)
- Asset cache next (previously downloaded from catalog)
- GitHub catalog last (on-demand discovery and download)
- First match wins, no merging (predictable behavior)

Complete replacement not field merging:

- Consistent with role replacement behavior (DR-005)
- Predictable: local task completely overrides global
- No hidden inheritance complexity
- Users know exactly what they configured

Source metadata for security:

- Shows where task came from (local vs global vs cache vs catalog)
- User can see if project contains suspicious local task
- Transparency before execution
- Security warning if downloading from catalog

Asset catalog integration:

- Tasks discoverable via `start assets search`
- Auto-download on first use when `asset_download = true`
- Cache populated transparently
- No bulk downloads required

Alias priority:

- Aliases treated as short names for tasks
- Same precedence as task names (local → global → cache → catalog)
- Warnings for shadowed aliases
- Prevents accidental conflicts

## Trade-offs

Accept:

- Local task completely replaces global (no partial override)
- Must configure task in local or global to make it "active"
- Alias conflicts require manual resolution
- Source metadata increases memory footprint slightly
- Four-level resolution adds complexity

Gain:

- Predictable replacement behavior (no hidden merging)
- Security transparency (source always visible)
- Consistent with role replacement pattern
- On-demand catalog access without bulk downloads
- Local projects can completely override global defaults
- Clear mental model (local wins, then global, then cache, then catalog)

## Alternatives

Field-level merging (shallow merge):

```toml
# Global
[tasks.code-review]
alias = "cr"
description = "Security review"
command = "git diff"

# Local
[tasks.code-review]
description = "Project review"
# Inherits: alias="cr", command="git diff"
```

- Pro: Less duplication when overriding single fields
- Pro: Can "extend" global tasks
- Con: Hidden inheritance complexity
- Con: Not obvious which fields come from where
- Con: Different behavior from roles (inconsistent)
- Con: Harder to reason about final configuration
- Rejected: Complete replacement is more predictable

Deep merging (recursive merge):

- Pro: Maximum inheritance capability
- Con: Very complex to reason about
- Con: Unclear which fields come from which source
- Con: Debugging becomes difficult
- Rejected: Too complex, unpredictable behavior

Global always wins (no local override):

- Pro: Simpler (one source only)
- Con: No project-specific customization
- Con: Can't adapt global tasks for specific projects
- Rejected: Too restrictive

Asset tasks auto-loaded at runtime:

```go
// Load order: global → local → assets (all active)
```

- Pro: Asset tasks immediately available without config
- Con: Performance impact (loading all asset files)
- Con: Namespace pollution (dozens of auto-loaded tasks)
- Con: Security risk (assets executed without user awareness)
- Rejected: Catalog resolution provides better balance

Three-way merge (global + local + assets):

- Pro: Could combine best of all sources
- Con: Extremely complex resolution rules
- Con: Unclear which field comes from which source
- Con: Performance impact
- Rejected: Too complex, unpredictable

## Loading Algorithm

Runtime execution:

```go
func LoadTasksForExecution() map[string]Task {
    tasks := make(map[string]Task)

    // 1. Load global config tasks
    for name, task := range globalConfig.Tasks {
        tasks[name] = task
        task.Source = "global"
        task.SourcePath = globalConfigPath
    }

    // 2. Load local config tasks (completely replaces global)
    for name, task := range localConfig.Tasks {
        tasks[name] = task  // Overwrites global if exists
        task.Source = "local"
        task.SourcePath = localConfigPath
    }

    return tasks
}
```

Task resolution (includes catalog per DR-033):

```go
func ResolveTask(input string) (*Task, error) {
    // 1. Local task name (exact match)
    if task, exists := localTasks[input]; exists {
        return task, nil
    }

    // 2. Local alias
    if taskName, exists := localAliases[input]; exists {
        return localTasks[taskName], nil
    }

    // 3. Global task name (exact match)
    if task, exists := globalTasks[input]; exists {
        return task, nil
    }

    // 4. Global alias
    if taskName, exists := globalAliases[input]; exists {
        return globalTasks[taskName], nil
    }

    // 5. Asset cache (DR-033)
    cachePath := findInCache("tasks", input)
    if cachePath != "" {
        task := loadFromCache(cachePath)
        return task, nil
    }

    // 6. GitHub catalog (DR-033)
    if !opts.AssetDownload {
        return nil, ErrNotFoundNoDownload
    }

    catalog := getCatalog()
    githubPath := catalog.Find("tasks", input)
    if githubPath == "" {
        return nil, ErrNotFound
    }

    // Download, cache, add to config
    task := downloadAsset(githubPath)
    cacheAsset(task)
    addToConfig(task, opts.Scope)

    return task, nil
}
```

## Replacement Behavior

When local defines same task name as global, local completely replaces global (no field merging):

```toml
# Global
[tasks.code-review]
alias = "cr"
description = "Security-focused review"
command = "git diff"
prompt = "Review for security issues"

# Local
[tasks.code-review]
description = "Project-specific review"
prompt = "Review for style"
# alias and command are now undefined (not inherited)
```

This matches role replacement behavior (DR-005).

## Source Metadata

Every loaded task includes source metadata for transparency and security:

```go
type Task struct {
    Name        string
    Alias       string
    Description string
    Role        string  // Optional role preference
    Agent       string  // Optional agent preference
    // ... UTD fields ...

    // Metadata (not in config file)
    Source      string  // "local", "global", "cache", "catalog"
    SourcePath  string  // Full path to config file or cache location
}
```

## Runtime Output

Normal output shows task source for transparency:

```
Starting task: code-review
─────────────────────────────────────────────────
Task source: local (./.start/tasks.toml)
Agent: claude (model: sonnet)
Role: code-reviewer

Context documents:
  ✓ environment     ~/reference/ENVIRONMENT.md
  ✓ agents          ./AGENTS.md

Executing...
```

Use case: If repository contains malicious `.start/tasks.toml` with suspicious task, user sees it's from local config and can stop execution.

## Task Listing

Show configured tasks with source indicators:

```
Configured tasks:
═══════════════════════════════════════════════════════════

Global tasks (2):
  security-review (sr)
    Security-focused code review

  quick-help (qh)
    Quick help with instructions

Local tasks (1):
  code-review (cr)
    Project-specific code review
```

Asset catalog tasks shown via:
```bash
start assets search "review"
start assets add  # Interactive TUI
```

## Alias Priority

Aliases are treated as short names for tasks. Local always wins.

Resolution:
- `start task cr` → Local alias wins over global alias
- Warning displayed at runtime about shadowed global alias

## Usage Examples

Complete replacement behavior:

```toml
# Global ~/.config/start/tasks.toml
[tasks.review]
alias = "r"
role = "code-reviewer"
command = "git diff --staged"
prompt = "Review: {command_output}"

# Local .start/tasks.toml
[tasks.review]
role = "security-auditor"
prompt = "Security check"
# alias and command NOT inherited - must be redefined if needed
```

Catalog resolution:

```bash
# Task not in config, auto-downloads from catalog
start task pre-commit-review

# Output:
# Task 'pre-commit-review' not found locally.
# Found in GitHub catalog: tasks/git-workflow/pre-commit-review
# Downloading...
# ✓ Cached to ~/.config/start/assets/tasks/git-workflow/
# ✓ Added to global config as 'pre-commit-review'
# Running task 'pre-commit-review'...
```

Source transparency:

```bash
start task code-review

# Output shows source:
# Starting task: code-review
# Task source: local (./.start/tasks.toml)
# [execution continues]
```

## Validation

At configuration load:

- Task names must be unique within scope (local or global)
- Aliases must be unique within scope
- Warn if local task/alias shadows global task/alias

At execution time:

- Selected task exists in merged config or catalog
- Task's role (if specified) exists
- Task's agent (if specified) exists

## Updates

- 2025-01-10: Extended by DR-031 (catalog architecture) and DR-033 (asset resolution algorithm including cache and GitHub catalog)
