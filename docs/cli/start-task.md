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

Executes predefined AI workflow tasks configured in `config.toml`. Tasks are reusable workflows with optional system prompt overrides, automatic required context inclusion, and dynamic content from shell commands.

**Task components:**

- **System prompt override** - Optional UTD fields to override global/local system prompt (`system_prompt_file`, `system_prompt_command`, `system_prompt`)
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

Tasks are defined in `config.toml` and can be customized per user or per project. See [task.md](../task.md) for configuration details.

## Arguments

**name**
: Task name or alias to execute. Use `start task` to see available tasks.

```bash
start task code-review        # By name
start task cr                 # By alias
```

**instructions** (optional)
: Additional instructions passed to the task's prompt template via `{instructions}` placeholder. Multi-word instructions must be quoted. If omitted, `{instructions}` is replaced with "None".

```bash
start task gdr "focus on security issues"
start task code-review "check error handling"
```

## Flags

All global flags from `start` command are supported except `--quiet`.

**--agent** _name_
: Override agent for this task. Overrides task's agent setting and default agent.

```bash
start task code-review --agent gemini
```

**--role** _name_
: Override role for this task. Overrides task's role setting and default role.

```bash
start task code-review --role security-auditor
```

**--model** _alias|name_
: Model to use. Accepts either model alias or full model name.

```bash
start task gdr --model opus
start task gdr --model claude-opus-4-20250514
```

**--directory** _path_, **-d** _path_
: Working directory for context detection and command execution.

```bash
start task gdr --directory ~/my-project
```

**--verbose**, **-v**
: Show detailed output including config resolution, context detection, and command execution.

**--debug**
: Debug mode. Shows all internal operations, placeholder resolution, and command construction.

**--help**, **-h**
: Show help. Behavior depends on usage:

- `start task --help` - Show general task command help
- `start task <name> --help` - Show specific task configuration and usage

**--quiet**, **-q**
: Silently ignored for tasks (no effect).

## Behavior

### No Arguments - List Tasks

```bash
start task
```

Displays all available tasks (default and custom) with aliases and descriptions.

### With Task Name - Execute Task

```bash
start task <name> [instructions]
```

**Execution flow:**

1. Load and merge configuration (global + local)
2. Find task by name or alias
3. Determine agent using precedence rules:
   - CLI `--agent` flag (highest priority)
   - Task `agent` field (if configured)
   - `default_agent` setting (fallback)
4. Determine role using precedence rules:
   - CLI `--role` flag (highest priority)
   - Task `role` field (if configured)
   - `default_role` setting
   - First role in config (TOML order)
5. Validate agent and role exist in configuration
   - UTD supports: file, command, prompt with placeholders
6. Load required contexts:
   - Auto-includes all contexts where `required = true`
   - No `documents` array needed
   - Skips missing files (same as `start` command behavior)
   - Order: Config definition order
7. Run task `command` if configured (UTD):
   - Execute in working directory
   - Capture stdout and stderr
   - Error and exit if command fails (non-zero exit code)
8. Build prompt from task's prompt template using UTD:
   - Load from `file` if specified
   - Execute `command` if specified
   - Process `prompt` template with placeholders
   - Replace `{instructions}` with user's args (or "None")
   - Replace `{command}` with command output (or empty string)
   - Replace global placeholders ({model}, {date})
   - Insert required context document prompts first
9. Display task summary (unless verbose/debug)
10. Execute agent command

### Task-Specific Help

```bash
start task <name> --help
```

Displays task configuration including system prompt override, required contexts, command, and usage examples.

## Task Configuration

Tasks are defined in `config.toml` using the **Unified Template Design (UTD)** pattern:

````toml
[tasks.git-diff-review]
alias = "gdr"
agent = "claude"                        # Optional: Preferred agent
description = "Review git diff changes"

# System prompt override (optional, UTD)
system_prompt_file = "~/.config/start/roles/code-reviewer.md"
system_prompt = """
{file}

Focus on security and maintainability.
"""

# Task prompt (UTD - at least one required)
command = "git diff --staged"
prompt = """
Review the following changes:

## Instructions
{instructions}

## Changes
```diff
{command}
````

"""

# Shell config (optional)
shell = "bash"
command_timeout = 30

````

**Agent Selection:**
- Tasks can specify preferred agent with `agent` field
- Precedence: `--agent` flag > task `agent` field > `default_agent` setting
- Agent must exist in `[agents.<name>]` configuration

**Role Selection:**
- Tasks can specify preferred role with `role` field
- Precedence: `--role` flag > task `role` field > `default_role` setting > first role in config
- Role must exist in `[roles.<name>]` configuration

**Note:** Tasks automatically include all contexts where `required = true`. No `documents` array needed.

See [task.md](../task.md) for complete configuration documentation.

## Examples

### List Available Tasks

```bash
start task
````

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

System prompt: (from task - code-reviewer.md)
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

System prompt: (from task - code-reviewer.md)
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
  System prompt: (from task configuration)
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
  Global: ~/.config/start/config.toml
  Local:  ./.start/config.toml (found)

Resolving task: git-diff-review (alias: gdr)
  Description: Review git diff changes
  System prompt override: code-reviewer.md with template
  Task command: git diff --staged
  Auto-includes required contexts

Loading system prompt (UTD)...
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

System prompt: (from task - code-reviewer.md)
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

System prompt: <source>
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

- System prompt file not found (if using `system_prompt_file`)
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

### System Prompt File Not Found

```
Error: System prompt file not found: ./tasks/code-reviewer.md

Task 'code-review' references missing file in system_prompt_file.
Update task configuration or create the file.
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

Add agent: start config agent add go-expert
Or override: start task go-review --agent claude
```

Exit code: 2

## Notes

### Task vs Root Command Differences

**`start task <name>`:**

- Can override system prompt with task-specific `system_prompt_*` fields (UTD)
- Auto-includes ONLY contexts where `required = true`
- Can run task `command` for dynamic content (UTD)
- Task prompt from UTD fields (`file`, `command`, `prompt`)
- `{instructions}` and `{command}` placeholders available

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
- `{command}` - Output from task `command` field (or empty string)

**Global:**

- `{model}` - Model name
- `{date}` - Current timestamp (ISO 8601)
- `{file}` - File contents (in UTD file fields)

**System prompt placeholders:**

In task `system_prompt_*` templates:
- `{file}` - Content from `system_prompt_file`
- `{command}` - Output from `system_prompt_command`

**Task prompt placeholders:**

In task prompt (`file`, `command`, `prompt`):
- `{file}` - Content from task `file`
- `{command}` - Output from task `command`
- `{instructions}` - User's arguments

### Required Context Handling

**Automatic inclusion:**

Tasks automatically include all contexts where `required = true`. No `documents` array needed in task configuration.

**Context resolution:**

1. Task loads all `[context.<name>]` sections where `required = true`
2. Contexts resolved from merged global + local config
3. If context file doesn't exist → skip with status display (same as `start` command)

**Context order in prompt:**

- Contexts appear in config definition order (global first, then local)
- Context prompts appear BEFORE task's prompt template
- Example: environment, index, then task prompt

### System Prompt File Location

By convention, task system prompt files are stored in:

- Global: `~/.config/start/roles/*.md`
- Local (per-project): `./.start/roles/*.md` or `./roles/*.md`

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

### Default Tasks

`start` ships with four default interactive review tasks:

1. **code-review (cr)** - General code review
2. **git-diff-review (gdr)** - Review git diff output
3. **comment-tidy (ct)** - Review and tidy code comments
4. **doc-review (dr)** - Documentation review

**Customization:**

- Override defaults by defining tasks with same name in config
- Add custom tasks
- Tasks defined in local `./.start/config.toml` override global

### Task Discovery

Tasks are discovered from:

1. Default tasks (embedded in binary)
2. Global config: `~/.config/start/config.toml`
3. Local config: `./.start/config.toml`

**Merge behavior:**

- Local tasks override global tasks (same name)
- All tasks from all sources combined
- Alphabetically sorted in task list

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

The `--quiet` flag is accepted but silently ignored for tasks. Tasks always display summary output before launching the agent.

**Rationale:**

- Tasks involve dynamic content (command output)
- Users need visibility into what content is being analyzed
- Summary is already concise
- Use `--verbose` or `--debug` for more detail if needed

## See Also

- start(1) - Launch with context
- start-prompt(1) - Launch with custom prompt
- start-config(1) - Manage configuration
- task.md - Task configuration reference
