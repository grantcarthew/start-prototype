# Unified Template Design (UTD)

A consistent pattern for defining content across `start` configuration sections.

## Overview

The Unified Template Design (UTD) provides a flexible way to combine static files, dynamic command output, and template text. It's used throughout the configuration for:

- `[roles.<name>]` - Role (system prompt) definitions
- `[contexts.<name>]` - Context documents for sessions
- `[tasks.<name>]` - Task prompt fields

## Core Concept

UTD uses three optional fields that work together via placeholders:

**Fields:**

- `file` - Path to a file
- `command` - Shell command to execute
- `prompt` - Template text

**Placeholders:**

- `{file}` - Replaced with file path (absolute, with ~ expanded)
- `{file_contents}` - Replaced with file contents
- `{command}` - Replaced with command string as written
- `{command_output}` - Replaced with command execution output (stdout/stderr)

**At least one field must be present.**

## Fields

### file (string, optional)

Path to a file. Supports tilde (`~`) expansion and relative paths.

```toml
file = "./ROLE.md"
file = "~/reference/ENVIRONMENT.md"
```

**Behavior:**

- If file exists → Path available via `{file}`, contents available via `{file_contents}`
- If file missing and prompt uses `{file}` or `{file_contents}` → **Warning**, skip entire section

### command (string, optional)

Shell command to execute. Supports single-line and multi-line (triple-quote) commands.

```toml
command = "git status --short"
```

```toml
command = """
git status --short
echo "---"
git log -5 --oneline
"""
```

**Behavior:**

- Command executed in working directory
- stdout and stderr both captured
- Command string available via `{command}`, output available via `{command_output}`
- Subject to timeout (see `command_timeout`)

### prompt (string, optional)

Template text that can contain placeholders: `{file}`, `{file_contents}`, `{command}`, `{command_output}`.

```toml
prompt = "Let's process this step by step."
```

```toml
prompt = "Read {file} for context."
```

```toml
prompt = """
Project Status:
{file_contents}

Recent Activity:
{command_output}
"""
```

## Shell Configuration

### Global Shell Setting

Define default shell in `[settings]`:

```toml
[settings]
shell = "bash"
command_timeout = 30  # seconds
```

**Defaults if not specified:**

- `shell` → Auto-detect: `bash` if available, otherwise `sh`
- `command_timeout` → 30 seconds

### Per-Section Shell Override

Override shell for specific contexts/tasks:

```toml
[contexts.node-version]
command = "console.log(process.version)"
shell = "node"
command_timeout = 5
```

**Priority:**

1. Section-specific `shell` field (if present)
2. Global `[settings] shell` (if configured)
3. Auto-detected shell (`bash` or `sh`)

### Supported Shells

`start` automatically handles argument flags for common shells:

| Shell/Interpreter | Flag   | Example                         |
| ----------------- | ------ | ------------------------------- |
| **Shells**        |        |                                 |
| `bash`            | `-c`   | `bash -c "git status"`          |
| `sh`              | `-c`   | `sh -c "git status"`            |
| `zsh`             | `-c`   | `zsh -c "git status"`           |
| `fish`            | `-c`   | `fish -c "git status"`          |
| **JavaScript**    |        |                                 |
| `node`            | `-e`   | `node -e "console.log('hi')"`   |
| `nodejs`          | `-e`   | `nodejs -e "console.log('hi')"` |
| `bun`             | `-e`   | `bun -e "console.log('hi')"`    |
| `deno`            | `eval` | `deno eval "console.log('hi')"` |
| **Python**        |        |                                 |
| `python`          | `-c`   | `python -c "print('hi')"`       |
| `python2`         | `-c`   | `python2 -c "print('hi')"`      |
| `python3`         | `-c`   | `python3 -c "print('hi')"`      |
| **Other**         |        |                                 |
| `ruby`            | `-e`   | `ruby -e "puts 'hi'"`           |
| `perl`            | `-E`   | `perl -E "say 'hi'"`            |
| **Unknown**       | `-c`   | Falls back to `-c` flag         |

### Command Timeout

Commands are subject to timeout limits:

```toml
[settings]
command_timeout = 30  # Global default: 30 seconds

[contexts.quick-check]
command = "git status"
command_timeout = 5   # Override: 5 seconds

[contexts.slow-analysis]
command = "npm run analyze"
command_timeout = 120  # Override: 2 minutes
```

**Behavior:**

- Command exceeds timeout → Killed, **Warning**: `"Command timeout after {N} seconds"`
- Output captured up to timeout point is used
- Empty output if command produces nothing before timeout

## Validation Rules

UTD validates field combinations and placeholder usage. At least one of `file`, `command`, or `prompt` must be present.

### 1. Only `file`

```toml
[roles.code-reviewer]
file = "./ROLE.md"
```

**Behavior:**

- File exists → **Use file contents directly** ✓
- File missing → **Warning**: `"File not found: {file}"`, ignore section

### 2. Only `command`

```toml
[contexts.git-status]
command = "git status --short"
```

**Behavior:**

- Execute command → **Use command output directly** ✓
- Command fails → **Warning**: `"Command failed: {error}"`, ignore section
- Command timeout → **Warning**: `"Command timeout"`, use partial output

### 3. Only `prompt`

```toml
[contexts.note]
prompt = "Important: This project uses Go 1.21"
```

**Behavior:**

- Prompt contains `{file}` or `{file_contents}` → **Warning**: `"No file defined but prompt uses {file}"`, ignore section
- Prompt contains `{command}` or `{command_output}` → **Warning**: `"No command defined but prompt uses {command}"`, ignore section
- Otherwise → **Use prompt as-is** (inline text) ✓

### 4. `file` + `command` (no `prompt`)

```toml
[contexts.project]
file = "./PROJECT.md"
command = "git log -5 --oneline"
```

**Behavior:**

- Read file contents
- If file contains `{command}` or `{command_output}` → Execute command, inject into placeholder, **use result** ✓
- If file doesn't contain `{command}` or `{command_output}` → **Warning**: `"Command defined but not used"`, ignore command, use file as-is ✓

**Example file with `{command_output}`:**

```markdown
# Project Status

Last 5 commits:
{command_output}
```

Result:

```markdown
# Project Status

Last 5 commits:
abc1234 Fix bug
def5678 Add feature
...
```

### 5. `file` + `prompt`

```toml
[contexts.environment]
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
```

**Behavior:**

- Prompt contains `{file}` or `{file_contents}` → Read file, inject into placeholder, **use prompt** ✓
- Prompt doesn't contain `{file}` or `{file_contents}` → **Warning**: `"File defined but not used"`, ignore file, use prompt as-is

### 6. `command` + `prompt`

```toml
[contexts.status]
command = "git status --short"
prompt = "Current status:\n{command_output}"
```

**Behavior:**

- Prompt contains `{command}` or `{command_output}` → Execute command, inject into placeholder, **use prompt** ✓
- Prompt doesn't contain `{command}` or `{command_output}` → **Warning**: `"Command defined but not used"`, ignore command, use prompt as-is

### 7. All three (`file` + `command` + `prompt`)

```toml
[contexts.full-state]
file = "./PROJECT.md"
command = "git status --short"
prompt = """
Project Documentation:
{file_contents}

Current Status:
{command_output}
"""
```

**Behavior:**

- Placeholders used in prompt → **Inject values, use prompt** ✓
- `{file}` or `{file_contents}` missing from prompt → **Warning**: `"File defined but not used"`, ignore file
- `{command}` or `{command_output}` missing from prompt → **Warning**: `"Command defined but not used"`, ignore command

### 8. Empty section (no fields)

```toml
[contexts.broken]
# No fields defined
```

**Behavior:**

- **Warning**: `"Empty section: at least one of file, command, or prompt required"`, ignore section

## Examples

### Simple File

```toml
[roles.code-reviewer]
file = "./ROLE.md"
```

Uses file contents directly.

### Simple Command

```toml
[contexts.git-status]
command = "git status --short"
description = "Working tree status"
```

Uses command output directly.

### Inline Prompt

```toml
[contexts.note]
prompt = "Important: This project uses Go 1.21"
description = "Go version requirement"
```

Uses prompt text directly.

### File with Template

```toml
[contexts.environment]
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
```

Injects file contents into prompt template.

### Command with Template

```toml
[contexts.recent-changes]
command = "git log -5 --oneline"
prompt = """
Recent commits:
{command_output}

Focus on these changes during the session.
"""
```

Injects command output into prompt template.

### File with Command Injection

Create a PROJECT.md file:

```markdown
# My Project

## Recent Activity

{command_output}

## Status

Work in progress.
```

Config:

```toml
[contexts.project]
file = "./PROJECT.md"
command = "git log -3 --oneline"
```

Result includes git log output where `{command_output}` appears in the file.

### Combined: File + Command + Prompt

```toml
[contexts.complete-status]
file = "./PROJECT.md"
command = "git status --short"
prompt = """
# Full Project Context

## Documentation
{file_contents}

## Working Tree
{command_output}

Use this context to understand current project state.
"""
```

Both file contents and command output injected into prompt template.

### Multi-line Script with Node.js

```toml
[contexts.package-info]
shell = "node"
command = """
const pkg = require('./package.json');
console.log(`${pkg.name}@${pkg.version}`);
console.log(`Dependencies: ${Object.keys(pkg.dependencies).length}`);
"""
prompt = "Package details:\n{command_output}"
```

### Python Analysis

```toml
[contexts.python-files]
shell = "python3"
command_timeout = 10
command = """
import os
py_files = [f for f in os.listdir('.') if f.endswith('.py')]
print(f"Python files: {len(py_files)}")
print('\\n'.join(py_files))
"""
prompt = "Python project files:\n{command_output}"
```

### Bun Runtime

```toml
[contexts.bun-version]
shell = "bun"
command = "console.log(Bun.version)"
prompt = "Using Bun {command_output}"
```

### Deno Example

```toml
[contexts.deno-check]
shell = "deno"
command = "console.log(Deno.version.deno)"
prompt = "Deno runtime: {command_output}"
```

## Where UTD is Used

### [roles.code-reviewer]

```toml
[roles.code-reviewer]
file = "./ROLE.md"
command = "git log -1 --format='%s'"
prompt = """
Role: {file_contents}

Current task: {command_output}
"""
```

### [contexts.\<name\>]

```toml
[contexts.environment]
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for context."
required = true

[contexts.git-status]
command = "git status --short"
prompt = "Repository state:\n{command_output}"
```

### [tasks.\<name\>] - System Prompt Override

```toml
[tasks.code-review]
role = "code-reviewer"  # References [roles.code-reviewer]
command = "git diff --staged"
prompt = """
Review these changes:

{command_output}

Instructions: {instructions}
"""
```

Note: Tasks reference roles by name using the `role` field. The role itself can use UTD for dynamic content.

### [tasks.\<name\>] - Task Prompt

```toml
[tasks.git-review]
# Standard task prompt fields
file = "./prompts/review-template.md"
command = "git diff --staged"
prompt = "Review this diff:\n{command_output}"
# ... other task fields
```

Note: Task prompt fields use standard UTD field names (`file`, `command`, `prompt`).

## Implementation Notes

### Working Directory

Commands execute in:

- Current working directory (default)
- Directory specified by `--directory` flag
- Project root (if detected)

### Error Handling

- File not found and prompt uses `{file}` or `{file_contents}` → Warning, skip entire section
- Command fails → Warning, `{command_output}` = empty string (or partial output)
- Command timeout → Warning, `{command_output}` = output captured before timeout
- Shell not found → Error, section ignored

### Security Considerations

**Command execution runs shell scripts with full system access:**

1. **Validate command sources** - Only execute commands from trusted configs
2. **Review local configs** - Local `./.start/` config files can execute arbitrary commands
3. **Be cautious with shared configs** - Review before using configs from others
4. **Timeout protection** - Commands are killed after timeout
5. **No automatic sudo** - Commands run with current user permissions

## See Also

- [Configuration Reference](./config.md) - Main config documentation
- [Design Record](./design-record.md) - Design decisions
- [Tasks](./task.md) - Task configuration details
