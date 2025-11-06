# DR-019: Task Loading and Merging Algorithm

**Date:** 2025-01-06
**Status:** Accepted
**Category:** Tasks

## Decision

Tasks load from global and local configs only; asset tasks serve as templates; local completely replaces global for same task name

## Task Sources

**Runtime task execution** (`start task <name>`):
- Global config: `~/.config/start/config.toml`
- Local config: `./.start/config.toml`
- Assets are NOT automatically loaded at runtime

**Management commands** (`start config task list/add/edit`):
- Assets at `~/.config/start/assets/tasks/*.toml` shown as available templates
- Users explicitly add asset tasks to their config to activate them

## Loading Algorithm (Runtime)

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

## Replacement Behavior

When local defines same task name as global, local **completely replaces** global (no field merging):

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

This matches system_prompt replacement behavior (DR-005).

## Task Resolution Algorithm

When user runs `start task <input>`, resolution priority:

```go
func ResolveTask(input string) *Task {
    // 1. Local task name (exact match)
    if task, exists := localTasks[input]; exists {
        return task
    }

    // 2. Local alias
    if taskName, exists := localAliases[input]; exists {
        return localTasks[taskName]
    }

    // 3. Global task name (exact match)
    if task, exists := globalTasks[input]; exists {
        return task
    }

    // 4. Global alias
    if taskName, exists := globalAliases[input]; exists {
        return globalTasks[taskName]
    }

    return nil // Task not found
}
```

## Alias Priority

Aliases are treated as short names for tasks. Local always wins.

Resolution:
- `start task cr` → Local alias wins over global alias
- Warning displayed at runtime about shadowed global alias

## Source Metadata Tracking

Every loaded task includes source metadata for transparency/security:

```go
type Task struct {
    Name        string
    Alias       string
    Description string
    // ... UTD fields ...

    // Metadata (not in config file)
    Source      string  // "global" or "local"
    SourcePath  string  // Full path to config file
}
```

## Runtime Output

Normal output shows task source for transparency:

```
Starting task: code-review
─────────────────────────────────────────────────
Task source: local (./.start/config.toml)
Agent: claude (model: sonnet)
System prompt: custom override

Context documents:
  ✓ environment     ~/reference/ENVIRONMENT.md
  ✓ agents          ./AGENTS.md

Executing...
```

**Use case:** If repository contains malicious `.start/config.toml` with suspicious task, user sees it's from local config and can stop execution.

## Task Listing (start config task list)

Show all three sources:

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

Available asset tasks (4):
  code-review (cr)
    Review code for quality and best practices

  git-diff-review (gdr)
    Review staged git changes
```

Note: Asset tasks with same name as user tasks are shown separately (no "override" concept since assets aren't active).

## Task Adding Workflow (start config task add)

Interactive flow with template option:

```
Add new task
─────────────────────────────────────────────────

Select scope:
  1) global (all projects)
  2) local (this project only)

Scope [1-2]: 1

Start from template or create new?
  1) Use asset template
  2) Create from scratch

Select [1-2]: 1

Available templates:
  1) code-review (cr) - Review code for quality and best practices
  2) git-diff-review (gdr) - Review staged git changes
  ...

(Interactive prompts for each field with template values as defaults)
```

## Rationale

- **Simplicity:** Two sources only (global + local), assets as templates
- **Predictability:** Local always wins (name or alias)
- **Transparency:** Source shown at runtime for security
- **Flexibility:** Asset templates make common tasks easy to adopt
- **Consistency:** Replacement behavior matches system_prompt (DR-005)
- **Discoverability:** Asset tasks visible in list command
- **Safety:** Conflict warnings prevent accidental shadowing

## Implementation Notes

```go
// Package: internal/config

type TaskRegistry struct {
    tasks   map[string]*Task
    aliases map[string]string  // alias -> task name
}

func (r *TaskRegistry) Load(globalCfg, localCfg *Config) {
    // Load global first
    // Load local second (overwrites)
    // Build alias map with local precedence
}

func (r *TaskRegistry) Resolve(input string) (*Task, error) {
    // Check local name -> local alias -> global name -> global alias
}

func (r *TaskRegistry) CheckConflicts(scope, name, alias string) []Conflict {
    // Return warnings about shadowing
}
```

## Related Decisions

- [DR-005](./dr-005-system-prompt.md) - Same replacement behavior
- [DR-009](./dr-009-task-structure.md) - Task structure
- [DR-011](./dr-011-asset-distribution.md) - Asset templates
- [DR-016](./dr-016-asset-discovery.md) - Task directory structure
