# start task

## Name

start task - Run predefined AI workflow tasks

## Synopsis

```bash
start task                     # List all tasks
start task <name>              # Run task
start task <name> [instructions] [flags]
```

## Description

Executes predefined AI workflow tasks configured in `tasks.toml`. Tasks are reusable workflows with optional system prompt overrides, automatic required context inclusion, and dynamic content from shell commands.

Like all assets (including roles and agents), tasks can be lazy-loaded from the GitHub catalog on first use.

**Task resolution**:

1. Exact match: local → global → cache → GitHub (lazy fetch)
2. Prefix match: local → global → cache → GitHub (short-circuit at first source with matches)
   - Single match → use it
   - Multiple matches → interactive selection (TTY) or error (non-TTY)

```bash
start task pre-commit-review  # Exact match
start task pre                # Prefix match (if unambiguous)
start task code               # Ambiguous: interactive picker or error
```

Tasks can be downloaded from the GitHub catalog on first use, cached locally for offline use, and added to your configuration automatically.

**Task components:**

- **Role selection** - Optional `role` field to specify which role (system prompt) to use
- **Required contexts** - Automatically includes all contexts where `required = true`
- **Task prompt** - UTD fields for prompt template (`file`, `command`, `prompt`)
- **Alias** - Optional short name for quick access
- **Shell config** - Optional shell and timeout overrides

**Common use cases:**

- Code reviews with specific focus
- Git diff analysis
- Documentation reviews
- Comment cleanup
- Any repeatable AI-assisted workflow

Tasks are defined in `tasks.toml` and can be customized per user or per project. See [start-config-task(1)](start-config-task.md) for configuration details and [start-assets(1)](start-assets.md) for catalog information.

## Arguments

**name**
: Task name or alias to execute. Supports exact match or prefix matching. Use `start task` to see available tasks.

```bash
start task code-review        # Exact match
start task code               # Prefix match (if unambiguous)
start task c                  # Ambiguous: interactive picker or error
```

**instructions** (optional)
: Additional instructions passed to the task's prompt template via `{instructions}` placeholder. Multi-word instructions must be quoted. If omitted, `{instructions}` is replaced with "None".

```bash
start task gdr "focus on security issues"
start task code-review "check error handling"
```

## Flags

All global flags from `start` command are supported. See `start --help` for a full list of global flags like `--agent`, `--role`, `--model`, `--directory`, `--verbose`, and `--debug`.

**--quiet**, **-q**
: Suppress task summary and context list. Useful for scripting or when you only want the agent's output.

**-l, --local**
: When lazy-loading a new task from the catalog, this flag adds the task to the local config (`./.start/tasks.toml`) instead of the global config.

**--help**, **-h**
: Shows help for the `task` command. Use `start task <name> --help` to see details for a specific task.

```bash
start task pre-commit-review --asset-download     # Force download
start task pre-commit-review --asset-download=false  # Fail if not found
```

## Behavior

### Task Resolution and Lazy Loading

When you run `start task <name>`, the CLI follows this resolution order:

**1. Local config** (`./.start/tasks.toml`)

- Project-specific tasks
- Highest priority

**2. Global config** (`~/.config/start/tasks.toml`)

- Your personal tasks
- Available across all projects

**3. Asset cache** (`~/.config/start/assets/tasks/`)

- Previously downloaded catalog tasks
- Used immediately without prompting

**4. GitHub catalog** (if `asset_download = true`)

- Query GitHub for task
- Prompt to download (or auto-download if configured)
- Cache locally and add to config
- Global by default, local with `--local` flag

**5. Error if not found**

**Example - Lazy loading:**

```bash
$ start task pre-commit-review

Task 'pre-commit-review' not found locally.
Found in GitHub catalog: tasks/git-workflow/pre-commit-review
Downloading...

✓ Cached to ~/.config/start/assets/tasks/git-workflow/
✓ Added to global config as 'pre-commit-review'

Running task 'pre-commit-review'...
[task executes]
```

**Next time:**

```bash
$ start task pre-commit-review

Running task 'pre-commit-review'...
[task executes immediately from config]
```

### No Arguments - List Tasks

```bash
start task
```

Displays all configured tasks with aliases and descriptions. To browse the GitHub catalog, use `start assets browse` or `start assets add`.

### With Task Name - Execute Task

```bash
start task <name> [instructions]
```

**Execution flow:**

1. Resolve task using resolution algorithm (see above)
2. Load and merge configuration (global + local)
3. Find task by name or alias
4. Determine agent using precedence rules:
   - CLI `--agent` flag (highest priority)
   - Task `agent` field (if configured)
   - `default_agent` setting
   - First agent in config (TOML order)
5. Determine role using precedence rules:
   - CLI `--role` flag (highest priority)
   - Task `role` field (if configured)
   - `default_role` setting
   - First role in config (TOML order)
6. Validate agent and role exist in configuration
   - UTD supports: file, command, prompt with placeholders
7. Load required contexts:
   - Auto-includes all contexts where `required = true`
   - Missing files generate warnings and are skipped (same as `start` command behavior)
   - Order: Config definition order
8. Run task `command` if configured (UTD):
   - Execute in working directory
   - Capture stdout and stderr
   - Error and exit if command fails (non-zero exit code)
9. Build prompt from task's prompt template using UTD:
   - Load from `file` if specified
   - Execute `command` if specified
   - Process `prompt` template with placeholders
   - Replace `{instructions}` with user's args (or "None")
   - Replace `{file}` with file path, `{file_contents}` with file contents
   - Replace `{command}` with command string, `{command_output}` with command output
   - Replace global placeholders ({date})
   - Insert required context document prompts first
10. Display task summary (unless verbose/debug)
11. Execute agent command

### Task-Specific Help

```bash
start task <name> --help
```

Displays task configuration including system prompt override, required contexts, command, and usage examples.

## Task Configuration

Tasks are defined in `tasks.toml` using the **Unified Template Design (UTD)** pattern:

`````toml
[tasks.git-diff-review]
alias = "gdr"
agent = "claude"                        # Optional: Preferred agent
role = "code-reviewer"                  # Optional: Preferred role
description = "Review git diff changes"

# Task prompt (UTD - at least one required)
command = "git diff --staged"
prompt = """
Review the following changes:

## Instructions
{instructions}

## Changes
```diff
{command_output}
```

"""

# Shell config (optional)
shell = "bash"
command_timeout = 30

`````

**Agent Selection:**

- Tasks can specify preferred agent with `agent` field
- Precedence: `--agent` flag > task `agent` field > `default_agent` setting > first agent in config
- Agent must exist in `[agents.<name>]` configuration

**Role Selection:**

- Tasks can specify preferred role with `role` field
- Precedence: `--role` flag > task `role` field > `default_role` setting > first role in config
- Role must exist in `[roles.<name>]` configuration

**Note:** Tasks automatically include all contexts where `required = true`.

See [tasks.md](../tasks.md) for complete configuration documentation.

## Examples

### List Available Tasks

```bash
start task
```

Output:

```
Available tasks:
  code-review (cr)      - Review code for quality and best practices
  git-diff-review (gdr) - Review git diff changes
  comment-tidy (ct)     - Review and tidy code comments
  doc-review (dr)       - Review and improve documentation

Use 'start task <name>' to run a task.
Use 'start task --help' for more information.
```

### Run Task by Name

```bash
start task code-review
```

Output:

```
Starting Task: code-review
===============================================================================================
Agent: claude (model: claude-3-7-sonnet-20250219)

Required contexts:
  ✓ environment     ~/reference/ENVIRONMENT.md
  ✓ index           ~/reference/INDEX.csv

Role: (from task - code-reviewer.md)
Command output: (none)
Instructions: None

Executing command...
❯ claude --model claude-3-7-sonnet-20250219 --append-system-prompt '...' '...'
```

### Run Task by Alias

```bash
start task gdr
```

Same as `start task git-diff-review`.

### Run Task with Instructions

```bash
start task git-diff-review "focus on security vulnerabilities"
start task gdr "ignore formatting changes"
```

Instructions replace `{instructions}` placeholder in task's prompt template.

### Task with Command

```bash
start task git-diff-review
```

Output:

```
Starting Task: git-diff-review
===============================================================================================
Agent: claude (model: claude-3-7-sonnet-20250219)

Required contexts:
  ✓ environment     ~/reference/ENVIRONMENT.md
  ✓ index           ~/reference/INDEX.csv

Role: (from task - code-reviewer.md)
Command output: git diff --staged (127 lines)
Instructions: None

Executing command...
❯ claude --model claude-3-7-sonnet-20250219 --append-system-prompt '...' '...'
```

### Override Agent

```bash
start task code-review --agent gemini
start task gdr --agent opencode "check performance"
```

### Override Model

```bash
start task code-review --model haiku
start task gdr --model claude-3-5-haiku-20241022
```

### Different Directory

```bash
start task gdr --directory ~/other-project
```

Required contexts and task command execute relative to specified directory.

### Task-Specific Help

```bash
start task code-review --help
```

Output:

```
Task: code-review (alias: cr)

Description:
  Review code for quality and best practices

Configuration:
  Role: (from task configuration)
  Required contexts: Auto-included (environment, index, etc.)
  Command: (none)

Usage:
  start task code-review
  start task code-review "special instructions"
  start task cr --agent gemini

Use 'start task --help' for general help.
```

### Verbose Mode

```bash
start task gdr --verbose
```

Output:

```
Loading configuration...
  Global: ~/.config/start/ (5 files)
  Local:  ./.start/ (found, 3 files)

Resolving task: git-diff-review (alias: gdr)
  Description: Review git diff changes
  Role override: code-reviewer.md with template
  Task command: git diff --staged
  Auto-includes required contexts

Loading role (UTD)...
  File: ~/.config/start/roles/code-reviewer.md → /Users/gcarthew/.config/start/roles/code-reviewer.md
  Size: 847 bytes
  Template: Applied with framing text

Detecting required contexts:
  environment: ~/reference/ENVIRONMENT.md → /Users/gcarthew/reference/ENVIRONMENT.md (exists)
  index: ~/reference/INDEX.csv → /Users/gcarthew/reference/INDEX.csv (exists)

Executing task command...
  Command: git diff --staged
  Working directory: /Users/gcarthew/Projects/my-app
  Output: 127 lines

Building prompt (UTD)...
  Required context prompts: 2 documents
  Task prompt template: 487 characters
  Instructions: "None"
  Command output: 3.2 KB (git diff)
  Final prompt size: 4.1 KB

Starting Task: git-diff-review
===============================================================================================
Agent: claude (model: claude-3-7-sonnet-20250219)

Required contexts:
  ✓ environment     ~/reference/ENVIRONMENT.md
  ✓ index           ~/reference/INDEX.csv

Role: (from task - code-reviewer.md)
Command output: git diff --staged (127 lines)
Instructions: None

Executing command...
❯ claude --model claude-3-7-sonnet-20250219 --append-system-prompt '...' '...'
```

## Output

### Task List (No Arguments)

```
Available tasks:
  code-review (cr)      - Review code for quality and best practices
  git-diff-review (gdr) - Review git diff changes
  comment-tidy (ct)     - Review and tidy code comments
  doc-review (dr)       - Review and improve documentation

Use 'start task <name>' to run a task.
Use 'start task --help' for more information.
```

### Task Execution (Normal)

```
Starting Task: <name>
===============================================================================================
Agent: <agent> (model: <model>)

Required contexts:
  ✓ <name>     <path>
  ✗ <name>     <path> (not found)

Role: <source>
Command output: <command> (<lines> lines) | (none)
Instructions: <text> | None

Executing command...
❯ <agent command>
```

### No Tasks Configured

```
No tasks configured.

Add tasks to your configuration:
  start config edit

See: https://github.com/grantcarthew/start#tasks
```

## Exit Codes

**0** - Success (task executed successfully)

**1** - Configuration error

- No tasks configured
- Invalid task configuration
- Config file syntax error

**2** - Task/agent error

- Task not found
- Agent not found in config
- Model not configured

**3** - File error

- Role file not found (if role references missing file)
- Working directory doesn't exist
- Config file permissions error

**4** - Runtime error

- Task command failed
- Agent tool not installed
- Agent command execution failed

## Error Handling

### Task Not Found

```
Error: Task 'nonexistent' not found.

Available tasks:
  code-review (cr)
  git-diff-review (gdr)
  comment-tidy (ct)
  doc-review (dr)

Use 'start task' to see all tasks.
```

Exit code: 2

### No Tasks Configured

```
No tasks configured.

Add tasks to your configuration:
  start config edit

See: https://github.com/grantcarthew/start#tasks
```

Exit code: 1

### Task Command Failed

```
Error: Task command failed: git diff --staged

Command output:
fatal: not a git repository

Exit code: 128
```

Exit code: 4

Task execution stops. Agent is not launched.

### Role File Not Found

```
Error: Role file not found: ~/.config/start/roles/code-reviewer.md

Task 'code-review' uses role 'code-reviewer' which references missing file.
Update role configuration or create the file.
```

Exit code: 3

### Invalid Task Configuration

```
Error: Task 'code-review' has invalid configuration.

Task prompt must have at least one of: file, command, or prompt

Update configuration: start config edit
```

Exit code: 1

### Agent Not Found

```
Error: Agent 'go-expert' not found (required by task 'go-review').

Configured agents:
  claude
  opencode

Add agent: start assets add go-expert
Or override: start task go-review --agent claude
```

Exit code: 2

## Notes

### Task vs Root Command Differences

**`start task <name>`:**

- Can override role with task-specific `role` field
- Auto-includes ONLY contexts where `required = true`
- Can run task `command` for dynamic content (UTD)
- Task prompt from UTD fields (`file`, `command`, `prompt`)
- Task-specific placeholders: `{instructions}`, `{file}`, `{file_contents}`, `{command}`, `{command_output}`

**`start` (root):**

- Uses configured roles from config (via `default_role` or `--role` flag)
- Includes ALL context documents (required + optional)
- No dynamic command execution
- Simple context document prompt concatenation
- No task-specific placeholders

### Task Placeholders

Task prompt templates support these placeholders:

**Task-specific:**

- `{instructions}` - User's command-line arguments (or "None")

**UTD placeholders (from task's UTD fields):**

- `{file}` - File path from task `file` field (absolute, ~ expanded)
- `{file_contents}` - Content from task `file` field
- `{command}` - Command string from task `command` field
- `{command_output}` - Output from task `command` execution (or empty string)

**Global:**

- `{date}` - Current timestamp (ISO 8601)

### Required Context Handling

**Automatic inclusion:**

Tasks automatically include all contexts where `required = true`.

**Context resolution:**

1. Task loads all `[contexts.<name>]` sections where `required = true`
2. Contexts resolved from merged global + local config
3. If context file doesn't exist → skip with status display (same as `start` command)

**Context order in prompt:**

- Contexts appear in config definition order (global first, then local)
- Context prompts appear BEFORE task's prompt template
- Example: environment, index, then task prompt

### Role Reference

Tasks do not define system prompts directly. Instead, they reference a **Role** which contains the system prompt.

- Tasks reference roles by name (e.g., `role = "code-reviewer"`)
- Roles are resolved using the standard asset resolution (Local → Global → Cache → GitHub)
- Role files are stored in:
  - Global: `~/.config/start/roles/*.md`
  - Local: `./.start/roles/*.md`
  - Cache: `~/.config/start/assets/roles/*/*.md`

Example task configuration with system prompt override:

```toml
[tasks.code-review]
# Use specific role by name
role = "code-reviewer"

# Or use a different role for this task
role = "security-auditor"

# No role specified - uses default_role setting or first role in config
# (omit role field)
```

### Task Command Execution

**Working directory:**

- Defaults to current directory (`pwd`)
- Override with `--directory` flag
- Paths in config resolve relative to working directory

**Command execution:**

- Runs in configured shell (default: global `shell` setting or auto-detected)
- Override shell per-task with `shell` field
- Environment variables inherited
- Timeout: Configurable via `command_timeout` (default: 30 seconds or global setting)
- Exit code: Non-zero = error and task stops

**Shell configuration:**

```toml
[tasks.git-review]
command = "git diff --staged"
shell = "bash"
command_timeout = 10
```

See [UTD shell configuration](../design/unified-template-design.md#shell-configuration) for supported shells.

**Common task commands:**

- `git diff --staged` - Staged changes
- `git diff HEAD~1` - Last commit
- `git log --oneline -n 10` - Recent commits
- `cat file.md` - File contents
- `find . -name "*.go"` - File listing

### Available Tasks

Tasks are available from the GitHub asset catalog and automatically download on first use (if `asset_download = true`).

**Example tasks available:**

1. **code-review (cr)** - General code review
2. **git-diff-review (gdr)** - Review staged git changes
3. **comment-tidy (ct)** - Review and improve code comments
4. **doc-review (dr)** - Review and improve documentation

**How it works:**

- Run `start task <name>` for any catalog task
- Task downloads from GitHub and adds to your config automatically
- Subsequent uses run from your config (no network required)
- Browse available tasks: `start assets browse` or `start assets add`

**Customization:**

- Override catalog tasks by defining the same name in your config
- Create custom tasks (global or local)
- Disable auto-download: `asset_download = false` setting

### Task Discovery

Tasks are discovered from:

1. Local config: `./.start/tasks.toml` (highest priority)
2. Global config: `~/.config/start/tasks.toml`
3. Asset cache: `~/.config/start/assets/tasks/`
4. GitHub catalog: `grantcarthew/start` (if `asset_download = true`)

**Merge behavior:**

- Local tasks override global/cache/catalog tasks (same name)
- First source with the task wins (resolution order)
- All tasks from all sources combined for listing

### Instructions Handling

**With instructions:**

```bash
start task gdr "focus on security"
```

`{instructions}` → `"focus on security"`

**Without instructions:**

```bash
start task gdr
```

`{instructions}` → `"None"`

**Multi-word instructions:**

```bash
# ✓ Correct - quoted
start task gdr "check error handling and logging"

# ✗ Wrong - shell splits arguments
start task gdr check error handling
```

### Quiet Flag Behavior

The `--quiet` flag suppresses the task summary and context list.

- **With `--quiet`**: Outputs only the agent's execution output (and task command output if part of the prompt).
- **Without `--quiet`**: Displays task header, agent details, role source, and context summary before execution.

Use `--quiet` when piping output or running in scripts.

## See Also

- start(1) - Launch with context
- start-prompt(1) - Launch with custom prompt
- start-config(1) - Manage configuration
- tasks.md - Task configuration reference
