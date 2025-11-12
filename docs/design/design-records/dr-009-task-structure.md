# DR-009: Task Structure and Placeholders

**Date:** 2025-01-03 (original), 2025-01-07 (updated)
**Status:** Accepted
**Category:** Tasks

## Decision

Tasks are predefined workflows that reference roles by name, use the Unified Template Design (UTD) pattern for prompts, and support task-specific placeholders including `{instructions}`, `{file}`, `{file_contents}`, `{command}`, and `{command_output}`.

## Complete Task Configuration

`````toml
[tasks.git-diff-review]
alias = "gdr"
description = "Review git diff changes"
role = "code-reviewer"           # References [roles.code-reviewer]
agent = "claude"                  # Optional: preferred agent
command = "git diff --staged"
shell = "bash"                    # Optional: override global shell
command_timeout = 30              # Optional: override global timeout
prompt = """
Analyze the following git diff and act as a code reviewer.

## Special Instructions

{instructions}

## Git Diff

```diff
{command_output}
```

"""
`````

## Task Fields

### Metadata Fields

**alias** (string, optional)
- Short name for quick access
- Must be unique across all tasks (global + local merged)
- Pattern: `/^[a-z0-9]+(-[a-z0-9]+)*$/`
- Example: `"gdr"`, `"cr"`, `"sec"`

**description** (string, optional)
- Human-readable help text
- Shown in task list and help output
- Example: `"Review git diff changes"`

### Role and Agent Selection

**role** (string, optional)
- Name reference to a role defined in `[roles.<name>]`
- Does NOT reference a file path (changed from original design)
- Example: `"code-reviewer"` (references `[roles.code-reviewer]`)
- If omitted: Uses `default_role` from settings (or first role in config)
- See [DR-005](./dr-005-role-configuration.md) for role configuration

**agent** (string, optional)
- Name reference to an agent defined in `[agents.<name>]`
- Example: `"claude"` (references `[agents.claude]`)
- If omitted: Uses `default_agent` from settings (or first agent in config)
- See [DR-029](./dr-029-task-agent-field.md) for agent selection

### Task Prompt (UTD Pattern)

Tasks use the **Unified Template Design (UTD)** pattern for prompt templates.

**file** (string, optional)
- Path to prompt template file
- File path available via `{file}`, contents available via `{file_contents}`

**command** (string, optional)
- Shell command to generate dynamic content
- Example: `"git diff --staged"`
- Command string available via `{command}`, output available via `{command_output}`
- Stderr and stdout both captured

**prompt** (string, optional)
- Prompt template text
- Can use `{file}`, `{file_contents}`, `{command}`, `{command_output}`, and `{instructions}` placeholders
- Can be multi-line string

**UTD Requirement:** At least one of `file`, `command`, or `prompt` must be present.

**shell** (string, optional)
- Override global shell for command execution
- Example: `"bash"`, `"node"`, `"python"`
- Defaults to `[settings] shell` or auto-detected shell

**command_timeout** (integer, optional)
- Override global timeout for command execution (seconds)
- Example: `30`
- Defaults to `[settings] command_timeout` or 30 seconds

See [Unified Template Design](../unified-template-design.md) for UTD details.

## Context Inclusion

Tasks **automatically include all contexts where `required = true`**.

There is **no `documents` array** in task configuration (removed from original design).

**Rationale:**
- Ensures critical context (like AGENTS.md, ENVIRONMENT.md) is always present
- Simplifies task configuration
- Tasks cannot exclude required contexts or cherry-pick optional contexts
- Separation: Required contexts = task context, Optional contexts = full session context

**Example:**

```toml
# Global config
[context.environment]
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
required = true  # Automatically included in all tasks

[context.project]
file = "./PROJECT.md"
prompt = "Read {file} for project context."
required = false  # NOT included in tasks

# Task configuration
[tasks.code-review]
role = "code-reviewer"
prompt = "Review code: {instructions}"
# No documents array needed - gets 'environment' automatically
```

## Task-Specific Placeholders

Available in task prompt templates only (not in role prompts or context prompts):

**{instructions}** - User's command-line arguments
- Value: Arguments provided after task name
- Empty case: `"None"` (not empty string)
- Example: `start task gdr "focus on security"` → `{instructions}` = `"focus on security"`
- Example: `start task gdr` → `{instructions}` = `"None"`

**{file}** - File path from task's `file` field
- Value: Absolute file path (~ expanded)
- Empty case: Empty string if no file defined
- Example: `file = "./template.md"` → `{file}` = `"/Users/username/project/template.md"`

**{file_contents}** - Content from task's `file` field
- Value: File contents
- Empty case: Empty string if no file defined or file missing
- Example: `file = "./template.md"` → `{file_contents}` = template contents

**{command}** - Command string from task's `command` field
- Value: Command string as written
- Empty case: Empty string if no command defined
- Example: `command = "git diff --staged"` → `{command}` = `"git diff --staged"`

**{command_output}** - Output from task's `command` execution
- Value: Stdout and stderr from command execution
- Empty case: Empty string if no command defined or command fails
- Example: `command = "git diff --staged"` → `{command_output}` = git diff output

## All Placeholders Available in Task Prompts

**Task-specific:**
- `{instructions}` - User's CLI arguments
- `{file}` - Task file path
- `{file_contents}` - Task file contents
- `{command}` - Task command string
- `{command_output}` - Task command output

**Global (available everywhere):**
- `{model}` - Current model name
- `{date}` - Current timestamp (ISO 8601)

**Not available in task prompts:**
- `{role}` - Only in agent commands
- `{role_file}` - Only in agent commands
- `{prompt}` - Only in agent commands

See [DR-007](./dr-007-placeholders.md) for complete placeholder reference.

## Execution Flow

When `start task <name>` is executed:

1. **Select role:**
   - CLI `--role` flag → use it
   - Else task `role` field → use it
   - Else `default_role` setting → use it
   - Else first role in config → use it

2. **Select agent:**
   - CLI `--agent` flag → use it
   - Else task `agent` field → use it
   - Else `default_agent` setting → use it
   - Else first agent in config → use it

3. **Load required contexts:**
   - All contexts where `required = true`
   - Resolve each context's UTD (file, command, prompt)
   - Build context prompts

4. **Execute task command** (if defined):
   - Run in working directory
   - Capture stdout and stderr
   - Error and exit if non-zero exit code

5. **Build task prompt:**
   - Load from `file` if specified
   - Execute `command` if specified
   - Process `prompt` template
   - Replace `{instructions}` with user args (or "None")
   - Replace `{file}` with file path, `{file_contents}` with file contents
   - Replace `{command}` with command string, `{command_output}` with command output
   - Replace global placeholders (`{model}`, `{date}`)

6. **Assemble final prompt:**
   - Required context prompts (in order)
   - Task prompt

7. **Resolve role:**
   - Load selected role from config
   - Resolve role's UTD (file, command, prompt)
   - Generate `{role}` content
   - Generate `{role_file}` path

8. **Execute agent:**
   - Replace placeholders in agent command
   - Execute command
   - Cleanup temp files (if created for `{role_file}`)

## Usage Examples

### Simple Task (Default Role and Agent)

```toml
[tasks.help]
alias = "h"
description = "General help"
prompt = "Help me with: {instructions}"
```

Usage:
```bash
start task help "how to structure this code"
start task h "explain error handling"
```

### Task with Specific Role

```toml
[tasks.code-review]
alias = "cr"
role = "code-reviewer"
description = "Code quality review"
prompt = "Review the code. Focus: {instructions}"
```

Usage:
```bash
start task code-review "security issues"
start task cr "performance"
```

### Task with Role and Agent

```toml
[tasks.go-review]
alias = "gor"
role = "go-expert"
agent = "claude"
description = "Go code review"
command = "git diff --staged -- '*.go'"
prompt = """
Review Go code changes:

{command_output}

Focus: {instructions}
"""
```

Usage:
```bash
start task go-review "check for goroutine leaks"
start task gor "verify error handling"
```

### Task with Dynamic Content

`````toml
[tasks.git-diff-review]
alias = "gdr"
role = "code-reviewer"
description = "Review git diff"
command = "git diff --staged"
prompt = """
Review these changes:

## Instructions
{instructions}

## Changes
```diff
{command_output}
```
"""
`````

Usage:
```bash
start task gdr "focus on security"
start task git-diff-review "check for bugs"
```

### Task with File Template

```toml
[tasks.doc-review]
alias = "dr"
role = "documentation-writer"
description = "Review documentation"
file = "./README.md"
prompt = """
Review this documentation:

{file_contents}

Improvements needed: {instructions}
"""
```

Usage:
```bash
start task doc-review "check clarity"
start task dr "add examples"
```

### Task with Shell Override

```toml
[tasks.api-check]
alias = "api"
role = "code-reviewer"
description = "Check API endpoints"
shell = "bash"
command_timeout = 60
command = """
echo "=== Routes ==="
grep -r "router\." src/ | head -20
echo "=== Recent Changes ==="
git log --oneline --grep="api" -5
"""
prompt = """
Review API endpoints:

{command_output}

Focus: {instructions}
"""
```

Usage:
```bash
start task api-check "verify REST compliance"
```

### Task with CLI Overrides

```bash
# Override task's role
start task code-review --role security-auditor

# Override task's agent
start task go-review --agent gemini

# Override both
start task gdr --role go-expert --agent gemini "check concurrency"

# Override model
start task code-review --model haiku "quick check"
```

## Scope and Merge Behavior

Tasks can be defined in both global and local configs:

**Global:** `~/.config/start/config.toml`
- Shared across all projects

**Local:** `./.start/config.toml`
- Project-specific tasks

**Merge behavior:**
- Global + local tasks **combined**
- Same task name: Local **completely replaces** global (no field merging)
- Task list alphabetically sorted
- Alias conflicts: First in TOML order wins (after merge)

**Example:**

```toml
# Global
[tasks.code-review]
role = "code-reviewer"
prompt = "Review: {instructions}"

# Local (completely replaces global)
[tasks.code-review]
role = "go-expert"
command = "git diff --staged"
prompt = "Review Go code: {command_output}"
```

Result: Local task used, global ignored.

## Validation

**At configuration load:**
- Task name matches pattern: `/^[a-z0-9]+(-[a-z0-9]+)*$/`
- At least one of `file`, `command`, or `prompt` present (UTD requirement)
- `role` field (if present) references existing role name
- `agent` field (if present) references existing agent name
- Alias (if present) matches same pattern as task name

**At execution time:**
- Selected task exists in merged config
- Task's role (if specified) exists
- Task's agent (if specified) exists
- Task file (if specified) exists (or warning)
- Task command (if specified) executes successfully

**Validation commands:**
- `start doctor` - Checks task configuration validity
- `start config validate` - Validates all config including tasks

## Rationale

**Why role field references role name (not file path):**
- Enables role reusability across multiple tasks
- Separates role definition from task definition
- Allows role override via `--role` flag
- Consistent with agent field pattern
- See [DR-005](./dr-005-role-configuration.md)

**Why no documents array:**
- Required contexts always included (ensures critical context present)
- Optional contexts for full sessions only (not tasks)
- Simpler task configuration
- Clear separation of concerns

**Why {instructions} defaults to "None":**
- Mirrors existing bash script behavior
- Clearer than empty string in prompts
- Explicit indication of no user instructions

**Why {command} instead of {content}:**
- More descriptive name
- Matches field name (`command = "..."`)
- Consistent with UTD pattern terminology
- See [DR-007](./dr-007-placeholders.md)

**Why UTD pattern for task prompts:**
- Consistent with contexts and roles
- Full flexibility (file, command, prompt, or combinations)
- Dynamic content generation support
- Template composition support

**Why shell and command_timeout overrides:**
- Different tasks may need different shells (bash vs node vs python)
- Long-running commands need higher timeouts
- Per-task control without global changes

## Breaking Changes from Original

This updates the original DR-009 with:

1. **Changed:** `role` field now references role name (not file path)
2. **Removed:** `documents` array (auto-includes required contexts)
3. **Added:** `agent` field for agent selection
4. **Changed:** `{content}` placeholder → `{command}`
5. **Added:** Full UTD pattern (file, command, prompt)
6. **Added:** `shell` and `command_timeout` overrides
7. **Updated:** All examples to use role-based design
8. **Added:** Role and agent selection precedence rules

## Related Decisions

- [DR-005](./dr-005-role-configuration.md) - Role configuration and selection
- [DR-007](./dr-007-placeholders.md) - Global and task-specific placeholders
- [DR-010](./dr-010-default-tasks.md) - Default task definitions
- [DR-012](./dr-012-context-required.md) - Required context field
- [DR-019](./dr-019-task-loading.md) - Task loading algorithm
- [DR-029](./dr-029-task-agent-field.md) - Task agent field
