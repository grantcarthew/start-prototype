# DR-009: Task Structure and Placeholders

- Date: 2025-01-03 (original), 2025-01-07 (updated)
- Status: Accepted
- Category: Tasks

## Problem

The tool needs a task system that allows users to define reusable workflows. The system must:

- Support predefined workflows with custom prompts and dynamic content
- Reference roles for system prompts without duplicating role definitions
- Allow both simple static tasks and complex dynamic tasks with command execution
- Work with task-specific user instructions provided at runtime
- Enable task-specific agent and role preferences while allowing CLI overrides
- Support both global (personal) and local (project-specific) task definitions
- Handle context inclusion consistently without requiring per-task configuration
- Follow the same configuration patterns as agents and roles

## Decision

Tasks are predefined workflows that reference roles by name, use the Unified Template Design (UTD) pattern for prompts, and support task-specific placeholders including `{instructions}`, `{file}`, `{file_contents}`, `{command}`, and `{command_output}`.

Tasks automatically include all contexts where `required = true`.

## Why

Named task pattern for reusability:

- Define once, execute repeatedly with different instructions
- Consistent with `[agents.<name>]` and `[roles.<name>]` pattern
- Easy to reference and execute by name or alias
- Supports both personal and team-shared workflows

Role reference instead of inline definition:

- Tasks reference roles by name, not file paths
- Enables role reusability across multiple tasks
- Separates role definition from task definition
- Allows role override via `--role` flag at execution time
- Consistent with agent field pattern

UTD pattern for flexibility:

- Simple case: just a prompt template
- Advanced: dynamic content via commands, file templates, or combinations
- One pattern supports both static and dynamic tasks
- Consistent with roles and contexts design

Task-specific placeholders for user interaction:

- `{instructions}` captures runtime user input
- Defaults to "None" when empty (explicit indication)
- Allows parameterized tasks without hardcoding values
- Mirrors existing bash script behavior

Automatic context inclusion:

- Required contexts (`required = true`) always included in tasks
- Ensures critical context (ENVIRONMENT.md, AGENTS.md) is always present
- Simplifies task configuration (no per-task context management)
- Clear separation: required contexts for tasks, optional contexts for full sessions
- Tasks cannot exclude required contexts or cherry-pick optional contexts

Agent field for task preferences:

- Tasks can specify preferred agent
- CLI override available via `--agent` flag
- Follows same precedence pattern as roles
- Enables task-specific agent optimization (specific task works better with specific agent)

Shell and timeout overrides:

- Different tasks may need different shells (bash vs node vs python)
- Long-running commands need higher timeouts
- Per-task control without changing global settings

## Trade-offs

Accept:

- More configuration structure than simple script execution
- Users must understand role references, UTD pattern, and placeholder syntax
- Task definitions can become verbose for complex workflows
- Automatic context inclusion means tasks cannot exclude required contexts
- Role and agent fields add configuration complexity

Gain:

- Reusable workflows across projects
- Dynamic content generation with command execution
- Consistent pattern with roles and contexts
- Runtime parameterization via `{instructions}`
- Role and agent override flexibility at execution time
- Clear separation between task definition and role definition
- Natural merge behavior with global/local scopes
- Guaranteed context inclusion for tasks

## Alternatives

Inline role definitions in tasks:

```toml
[tasks.review]
system_prompt_file = "ROLE.md"
command = "git diff"
prompt = "Review: {command_output}"
```

- Pro: Everything in one place, simpler to understand initially
- Pro: No separate role to reference
- Con: Role definitions duplicated across tasks using the same role
- Con: Cannot easily swap roles for experimentation
- Con: Difficult to maintain (change role = update all tasks)
- Con: No reusability across tasks
- Rejected: Poor reusability and maintenance burden

Array-based task context selection:

```toml
[tasks.review]
contexts = ["environment", "project", "agents"]
prompt = "Review code"
```

- Pro: Fine-grained control over which contexts included
- Pro: Can cherry-pick contexts per task
- Con: Must manage context list for every task
- Con: Easy to forget critical contexts
- Con: Inconsistent context inclusion across tasks
- Con: More configuration burden
- Rejected: Automatic required context inclusion is simpler and more consistent

Simple script execution (no UTD):

```toml
[tasks.review]
command = "git diff"
```

- Pro: Extremely simple configuration
- Pro: Minimal learning curve
- Con: No prompt customization
- Con: No dynamic content composition
- Con: Cannot combine file templates with command output
- Con: Limited flexibility
- Rejected: Too limited for diverse task needs

Command-line positional arguments instead of `{instructions}`:

```bash
start task review $1 $2 $3
```

- Pro: Standard shell pattern
- Pro: Multiple separate arguments
- Con: Must handle variable argument counts
- Con: Argument ordering becomes significant
- Con: Spaces in arguments require quoting complexity
- Con: Less flexible than single instruction string
- Rejected: Single instruction string is simpler and more flexible

## Structure

Complete task configuration:

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

Task fields:

Metadata fields:

- `alias` (string, optional) - Short name for quick access, must be unique across all tasks (global + local merged), pattern: `/^[a-z0-9]+(-[a-z0-9]+)*$/`, example: `"gdr"`, `"cr"`, `"sec"`
- `description` (string, optional) - Human-readable help text shown in task list and help output

Role and agent selection:

- `role` (string, optional) - Name reference to a role defined in `[roles.<name>]`, does NOT reference a file path, example: `"code-reviewer"` (references `[roles.code-reviewer]`), if omitted uses `default_role` from settings (or first role in config)
- `agent` (string, optional) - Name reference to an agent defined in `[agents.<name>]`, example: `"claude"` (references `[agents.claude]`), if omitted uses `default_agent` from settings (or first agent in config)

Task prompt (UTD pattern):

- `file` (string, optional) - Path to prompt template file, file path available via `{file}`, contents available via `{file_contents}`
- `command` (string, optional) - Shell command to generate dynamic content, example: `"git diff --staged"`, command string available via `{command}`, output available via `{command_output}`, stderr and stdout both captured
- `prompt` (string, optional) - Prompt template text, can use `{file}`, `{file_contents}`, `{command}`, `{command_output}`, and `{instructions}` placeholders, can be multi-line string
- `shell` (string, optional) - Override global shell for command execution, example: `"bash"`, `"node"`, `"python"`, defaults to `[settings] shell` or auto-detected shell
- `command_timeout` (integer, optional) - Override global timeout for command execution (seconds), example: `30`, defaults to `[settings] command_timeout` or 30 seconds

UTD requirement: At least one of `file`, `command`, or `prompt` must be present.

Task-specific placeholders (available in task prompt templates only):

- `{instructions}` - User's command-line arguments, value is arguments provided after task name, empty case is `"None"` (not empty string), example: `start task gdr "focus on security"` → `{instructions}` = `"focus on security"`, example: `start task gdr` → `{instructions}` = `"None"`
- `{file}` - File path from task's `file` field, value is absolute file path (~ expanded), empty case is empty string if no file defined
- `{file_contents}` - Content from task's `file` field, value is file contents, empty case is empty string if no file defined or file missing
- `{command}` - Command string from task's `command` field, value is command string as written, empty case is empty string if no command defined
- `{command_output}` - Output from task's `command` execution, value is stdout and stderr from command execution, empty case is empty string if no command defined or command fails

All placeholders available in task prompts:

Task-specific:

- `{instructions}` - User's CLI arguments
- `{file}` - Task file path
- `{file_contents}` - Task file contents
- `{command}` - Task command string
- `{command_output}` - Task command output

Universal (available everywhere):

- `{date}` - Current timestamp (ISO 8601)

Not available in task prompts:

- `{bin}` - Only in agent commands
- `{model}` - Only in agent commands
- `{role}` - Only in agent commands
- `{role_file}` - Only in agent commands
- `{prompt}` - Only in agent commands

## Execution Flow

When `start task <name>` is executed:

1. Select role:
   - CLI `--role` flag → use it
   - Else task `role` field → use it
   - Else `default_role` setting → use it
   - Else first role in config → use it

2. Select agent:
   - CLI `--agent` flag → use it
   - Else task `agent` field → use it
   - Else `default_agent` setting → use it
   - Else first agent in config → use it

3. Load required contexts:
   - All contexts where `required = true`
   - Resolve each context's UTD (file, command, prompt)
   - Build context prompts

4. Execute task command (if defined):
   - Run in working directory
   - Capture stdout and stderr
   - Error and exit if non-zero exit code

5. Build task prompt:
   - Load from `file` if specified
   - Execute `command` if specified
   - Process `prompt` template
   - Replace `{instructions}` with user args (or "None")
   - Replace `{file}` with file path, `{file_contents}` with file contents
   - Replace `{command}` with command string, `{command_output}` with command output
   - Replace `{date}` placeholder

6. Assemble final prompt:
   - Required context prompts (in order)
   - Task prompt

7. Resolve role:
   - Load selected role from config
   - Resolve role's UTD (file, command, prompt)
   - Generate `{role}` content
   - Generate `{role_file}` path

8. Execute agent:
   - Replace placeholders in agent command
   - Execute command
   - Cleanup temp files (if created for `{role_file}`)

## Usage Examples

Simple task (default role and agent):

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

Task with specific role:

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

Task with role and agent:

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

Task with dynamic content:

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

Task with file template:

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

Task with shell override:

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

Task with CLI overrides:

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

## Scope

Tasks can be defined in both global and local configs:

Global: `~/.config/start/tasks.toml` - shared across all projects

Local: `./.start/tasks.toml` - project-specific tasks

Merge behavior:

- Global + local tasks combined
- Same task name: local completely replaces global (no field merging)
- Task list alphabetically sorted
- Alias conflicts: first in TOML order wins (after merge)

Example:

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

At configuration load:

- Task name matches pattern: `/^[a-z0-9]+(-[a-z0-9]+)*$/`
- At least one of `file`, `command`, or `prompt` present (UTD requirement)
- `role` field (if present) references existing role name
- `agent` field (if present) references existing agent name
- Alias (if present) matches same pattern as task name

At execution time:

- Selected task exists in merged config
- Task's role (if specified) exists
- Task's agent (if specified) exists
- Task file (if specified) exists (or warning)
- Task command (if specified) executes successfully

Validation commands:

- `start doctor` - checks task configuration validity
- `start config validate` - validates all config including tasks

## Breaking Changes from Original

This updates the original DR-009 with:

1. Changed: `role` field now references role name (not file path)
2. Removed: `documents` array (auto-includes required contexts)
3. Added: `agent` field for agent selection
4. Changed: `{content}` placeholder → `{command}`
5. Added: Full UTD pattern (file, command, prompt)
6. Added: `shell` and `command_timeout` overrides
7. Updated: All examples to use role-based design
8. Added: Role and agent selection precedence rules

## Updates

- 2025-01-17: Fixed {model} placeholder scope - removed from tasks (agent-commands-only per DR-007)
