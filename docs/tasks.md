# Tasks

Tasks are predefined AI workflows with custom system prompts and dynamic content.

## Purpose

Tasks allow you to define reusable AI-assisted workflows. Each task can:

- Override the system prompt with a task-specific role
- Generate dynamic content from shell commands
- Use template prompts with placeholders
- Optionally provide a short alias for quick access

**Key difference from contexts:**
- **Context** = Passive background information included in sessions
- **Task** = Active workflow with specific role and structured prompt

## Configuration

Tasks are defined in both global and local config files:

- **Global:** `~/.config/start/config.toml` - Shared tasks across all projects
- **Local:** `./.start/config.toml` - Project-specific tasks

Tasks use the section name `[tasks.<name>]` where `<name>` is a unique identifier.

### Basic Structure

```toml
[tasks.<name>]
alias = "..."                           # Optional: Short name
description = "..."                     # Optional: Help text
agent = "..."                           # Optional: Preferred agent
role = "..."                            # Optional: Preferred role

# Task prompt (UTD pattern - at least one required)
file = "..."                            # Optional: Path to prompt file
command = "..."                         # Optional: Dynamic content command
prompt = "..."                          # Optional: Template with {file}/{command}/{instructions}

# Shell overrides (optional)
shell = "..."                           # Optional: Override global shell
command_timeout = 30                    # Optional: Override global timeout
```

## Fields

### Metadata

**alias** (string, optional)
: Short name for quick access. Must be unique across all tasks (global + local). Uses lowercase-kebab-case.

```toml
[tasks.git-diff-review]
alias = "gdr"
```

**description** (string, optional)
: Human-readable description shown in task list and help output.

```toml
[tasks.git-diff-review]
description = "Review staged git changes"
```

**agent** (string, optional)
: Preferred agent for this task. Must reference an agent defined in `[agents.<name>]` configuration. Agent selection precedence: CLI `--agent` flag > task `agent` field > `default_agent` setting > first agent in config.

```toml
[tasks.go-review]
agent = "go-expert"
description = "Review Go code with specialized agent"
```

If omitted, uses the `default_agent` from settings, or the first agent defined in config (TOML order) if no default is set.

**Use cases:**
- Specialized agents for specific languages or domains
- Different model perspectives for alternative reviews
- Performance optimization (fast agents for quick checks)
- Tool-specific features (vision, artifacts, etc.)

**Validation:**
- Agent name must match an existing `[agents.<name>]` section
- Validated at task execution time
- Also checked by `start doctor` and `start config validate`

### Role Field

**role** (string, optional)
: Preferred role for this task. Must reference a role defined in `[roles.<name>]` configuration.

**Role Selection Precedence:**
1. CLI `--role` flag (highest priority)
2. Task `role` field
3. `default_role` setting
4. First role in config (TOML order)

```toml
[tasks.security-audit]
role = "security-auditor"
agent = "claude"
description = "Security-focused code audit"
```

**Validation:**
- Role name must match an existing `[roles.<name>]` section
- Validated at task execution time
- Also checked by `start doctor` and `start config validate`

**Use cases:**
- Task-specific AI personas (security auditor, code reviewer, documentation writer)
- Different perspectives for same codebase
- Specialized domain knowledge (Go expert, API validator, etc.)

See [DR-005](./design/decisions/dr-005-role-configuration.md) for role configuration details.

### Task Prompt

Tasks use the **[Unified Template Design (UTD)](./design/unified-template-design.md)** pattern for prompts.

**file** (string, optional)
: Path to a prompt template file. File contents available via `{file}` placeholder.

**command** (string, optional)
: Shell command to generate dynamic content (e.g., `git diff --staged`). Output available via `{command}` placeholder.

**prompt** (string, optional)
: Template text that can use `{file}`, `{command}`, and `{instructions}` placeholders.

**Requirement:** At least one of `file`, `command`, or `prompt` must be present.

See [UTD validation rules](./design/unified-template-design.md#validation-rules) for field combination behavior.

### Shell Configuration

**shell** (string, optional)
: Override global shell for command execution in this task. Defaults to global `[settings] shell` or auto-detected shell.

See [UTD shell configuration](./design/unified-template-design.md#shell-configuration) for supported shells.

**command_timeout** (integer, optional)
: Override global timeout for command execution (in seconds). Defaults to global `[settings] command_timeout` or 30 seconds.

## Placeholders

### In Task Prompt Templates

Available in `file` content and `prompt` templates:

- `{file}` - Content from task `file`
- `{command}` - Output from task `command`
- **`{instructions}`** - User's command-line arguments (task-specific)
- `{model}` - Model name (global)
- `{date}` - Current timestamp (global)

**{instructions} behavior:**
- Value: User's arguments after task name
- If no arguments provided: `"None"`
- Example: `start task gdr "focus on security"` → `{instructions}` = `"focus on security"`
- Example: `start task gdr` → `{instructions}` = `"None"`

## Context Inclusion

Tasks automatically include **all contexts where `required = true`**.

There is no `documents` array in task configuration. Instead:

1. All contexts with `[context.<name>] required = true` are automatically included
2. Tasks cannot exclude required contexts or include optional contexts
3. This ensures critical context (like AGENTS.md) is always present

**Example:**

```toml
# Global config
[context.environment]
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
required = true  # Always included in tasks

[context.project]
file = "./PROJECT.md"
prompt = "Read {file} for project status."
required = false  # Never included in tasks

# Task automatically gets 'environment' context
[tasks.code-review]
prompt = "Review the code. {instructions}"
```

## Scope and Merge Behavior

Tasks can be defined in **both global and local** configs:

**Global tasks:** `~/.config/start/config.toml`
- Shared across all projects
- Default tasks (cr, gdr, ct, dr)

**Local tasks:** `./.start/config.toml`
- Project-specific workflows
- Override global tasks by using same name

**Merge behavior:**
- Global + local tasks are combined
- Same task name: **local overrides global**
- Task list alphabetically sorted
- Alias conflicts: First in TOML order wins (after merge)

## Validation

### Required Fields

At least one of `file`, `command`, or `prompt` must be present for the task prompt.

### Naming Constraints

**Task names:**
- Lowercase alphanumeric with hyphens
- Pattern: `/^[a-z0-9]+(-[a-z0-9]+)*$/`
- Examples: `code-review`, `git-diff-review`, `doc-review`

**Task aliases:**
- Same constraints as task names
- Must be unique across all tasks (global + local merged)
- Conflict resolution: First in TOML order wins

### Warnings

- **Role not found:** `"Role 'security-auditor' not found in configuration"` - Error at execution time
- **Alias conflict:** `"Alias 'cr' used by multiple tasks, using first: code-review"`
- **Field defined but not used:** `"file defined but not used in prompt template"`

See [UTD validation rules](./design/unified-template-design.md#validation-rules) for complete field validation behavior.

## Examples

### Simple Task (Inline System Prompt)

```toml
[tasks.code-review]
alias = "cr"
description = "General code quality review"
role = "code-reviewer"

prompt = "Review the code in this project. {instructions}"
```

Usage: `start task cr "check error handling"`

### Task with Role Reference

```toml
[tasks.git-diff-review]
alias = "gdr"
description = "Review staged git changes"
role = "code-reviewer"

command = "git diff --staged"
prompt = """
Review the following git diff.

## Instructions
{instructions}

## Staged Changes
```diff
{command}
```
"""
```

Usage: `start task gdr "focus on security"`

### System Prompt Template with File

```toml
[tasks.security-review]
alias = "sec"
description = "Security-focused code review"
role = "security-auditor"

command = "git diff --staged"
prompt = "Security review:\n{command}\n\n{instructions}"
```

### Task Without Specific Role

```toml
[tasks.quick-check]
alias = "qc"
description = "Quick review with default role"

# No role field = uses default_role or first role in config

prompt = "Quick code check: {instructions}"
```

### Multi-Line Shell Script

```toml
[tasks.api-check]
alias = "api"
description = "Validate API endpoints"
role = "api-validator"

shell = "bash"
command = """
echo "=== API Endpoints ==="
grep -r "router\." src/ | cut -d: -f2 | sort | uniq
echo ""
echo "=== Recent API Changes ==="
git log --oneline --grep="api" -5
"""
prompt = """
Validate these API endpoints and recent changes:

{command}

Focus: {instructions}
"""
```

### Node.js Script Task

```toml
[tasks.package-info]
alias = "pkg"
description = "Analyze package.json"
role = "nodejs-expert"

shell = "node"
command_timeout = 10
command = """
const pkg = require('./package.json');
console.log(`Name: ${pkg.name}`);
console.log(`Version: ${pkg.version}`);
console.log(`Dependencies: ${Object.keys(pkg.dependencies || {}).length}`);
console.log(`DevDependencies: ${Object.keys(pkg.devDependencies || {}).length}`);
"""
prompt = """
Analyze this package:

{command}

Recommendations: {instructions}
"""
```

### Project-Specific Task (Local Config)

```toml
# ./.start/config.toml
[tasks.validate-go]
alias = "vgo"
agent = "go-expert"
description = "Project-specific Go validation"

role = "project-reviewer"

command = """
go vet ./... 2>&1
echo "---"
golangci-lint run --fast 2>&1
"""
prompt = """
Go validation results:

{command}

Address: {instructions}
"""
```

### Task with Specialized Agent

```toml
[tasks.security-audit]
alias = "sec"
agent = "security-specialist"
description = "Security-focused code audit"

role = "security-auditor"

command = "git diff --staged"
prompt = """
Perform a security audit on these changes:

{command}

Focus areas: {instructions}
"""
```

### Task with Performance-Optimized Agent

```toml
[tasks.quick-lint]
alias = "ql"
agent = "haiku-agent"
role = "code-linter"
description = "Fast linting with lightweight agent"

command = "git diff --staged"
prompt = "Quick lint: {command}"
```

### Override Global Task (Local Config)

```toml
# Global: ~/.config/start/config.toml
[tasks.code-review]
alias = "cr"
role = "reviewer"
prompt = "Review code: {instructions}"

# Local: ./.start/config.toml (overrides global)
[tasks.code-review]
alias = "cr"
role = "go-expert"
command = "git diff --staged"
prompt = "Review Go code changes:\n{command}\n\n{instructions}"
```

The local task completely replaces the global task with the same name.

## Default Tasks

`start` ships with four default interactive review tasks:

1. **code-review** (alias: `cr`) - General code review for quality and best practices
2. **git-diff-review** (alias: `gdr`) - Review git diff output
3. **comment-tidy** (alias: `ct`) - Review and tidy code comments
4. **doc-review** (alias: `dr`) - Review and improve documentation

Users can:
- Override defaults by defining tasks with the same name in config
- Add custom tasks (global or local)
- The default tasks are embedded in the binary

## Usage

```bash
# List all tasks
start task

# Run task by name
start task code-review

# Run task by alias
start task cr

# With instructions
start task git-diff-review "focus on security"
start task gdr "ignore formatting"

# With agent override
start task code-review --agent gemini

# With model override
start task gdr --model opus

# Combined flags
start task cr --agent gemini --model flash "check error handling"

# Task-specific help
start task code-review --help
```

See [start-task.md](./cli/start-task.md) for complete CLI documentation.

## See Also

- [start-task CLI Reference](./cli/start-task.md) - Command usage and flags
- [Unified Template Design](./design/unified-template-design.md) - UTD pattern details
- [Configuration Reference](./config.md) - Complete config documentation
- [Design Record](./design/design-record.md) - Design decisions (DR-009, DR-010)
