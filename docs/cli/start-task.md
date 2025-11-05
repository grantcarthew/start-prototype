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

Executes predefined AI workflow tasks configured in `config.toml`. Tasks are reusable workflows with specific roles, prompt templates, context documents, and optional dynamic content from shell commands.

**Task components:**

- **Role** - System prompt specific to the task
- **Documents** - Subset of context documents to include
- **Content command** - Optional shell command (e.g., `git diff --staged`)
- **Prompt template** - Template with placeholders for instructions and content
- **Alias** - Optional short name for quick access

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

**--model** _alias|name_
: Model to use. Accepts either model alias or full model name.

```bash
start task gdr --model opus
start task gdr --model claude-opus-4-20250514
```

**--directory** _path_, **-d** _path_
: Working directory for context detection and content_command execution.

```bash
start task gdr --directory ~/my-project
```

**--verbose**, **-v**
: Show detailed output including config resolution, document detection, and content_command execution.

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
3. Determine agent (task override, `--agent` flag, or default)
4. Load role/system prompt:
   - If `role` is file path → read file contents
   - If `role` is inline text → use directly
   - Error if file path doesn't exist
5. Load documents specified in task's `documents` array:
   - Resolves document names to `[context.documents.*]` config
   - Skips undefined document names (no error)
   - Skips missing files (same as `start` command behavior)
   - Order: As listed in task's `documents` array
6. Run `content_command` if configured:
   - Execute in working directory
   - Capture stdout and stderr
   - Error and exit if command fails (non-zero exit code)
7. Build prompt from task's prompt template:
   - Insert context documents first (with their prompts)
   - Append task prompt template
   - Replace `{instructions}` with user's args (or "None")
   - Replace `{content}` with content_command output (or empty string)
   - Replace global placeholders ({model}, {system_prompt}, {date})
8. Display task summary (unless verbose/debug)
9. Execute agent command

### Task-Specific Help

```bash
start task <name> --help
```

Displays task configuration including role, documents, content command, and usage examples.

## Task Configuration

Tasks are defined in `config.toml`:

````toml
[task.git-diff-review]
alias = "gdr"
description = "Review git diff changes"
role = "./tasks/code-reviewer.md"
documents = ["environment", "agents"]
content_command = "git diff --staged"
prompt = """
Review the following changes:

## Instructions
{instructions}

## Changes
```diff
{content}
````

"""

````

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

Task documents:
  ✓ environment     ~/reference/ENVIRONMENT.md
  ✓ agents          ./AGENTS.md

Role: ./tasks/code-reviewer.md
Content: (none)
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

### Task with Content Command

```bash
start task git-diff-review
```

Output:

```
Starting Task: git-diff-review
===============================================================================================
Agent: claude (model: claude-3-7-sonnet-20250219)

Task documents:
  ✓ environment     ~/reference/ENVIRONMENT.md
  ✓ agents          ./AGENTS.md

Role: ./tasks/code-reviewer.md
Content: git diff --staged (127 lines)
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

Context documents and content_command execute relative to specified directory.

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
  Role: ./tasks/code-reviewer.md
  Documents: environment, agents
  Content command: (none)

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
  Role: ./tasks/code-reviewer.md
  Documents: environment, agents
  Content command: git diff --staged

Loading role...
  Path: ./tasks/code-reviewer.md → /Users/gcarthew/Projects/my-app/tasks/code-reviewer.md
  Size: 847 bytes

Detecting task documents:
  environment: ~/reference/ENVIRONMENT.md → /Users/gcarthew/reference/ENVIRONMENT.md (exists)
  agents: ./AGENTS.md → /Users/gcarthew/Projects/my-app/AGENTS.md (exists)

Executing content command...
  Command: git diff --staged
  Working directory: /Users/gcarthew/Projects/my-app
  Output: 127 lines

Building prompt...
  Document prompts: 2 documents
  Task prompt template: 487 characters
  Instructions: "None"
  Content: 3.2 KB (git diff output)
  Final prompt size: 4.1 KB

Starting Task: git-diff-review
===============================================================================================
Agent: claude (model: claude-3-7-sonnet-20250219)

Task documents:
  ✓ environment     ~/reference/ENVIRONMENT.md
  ✓ agents          ./AGENTS.md

Role: ./tasks/code-reviewer.md
Content: git diff --staged (127 lines)
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

Task documents:
  ✓ <name>     <path>
  ✗ <name>     <path> (not found)

Role: <path>
Content: <command> (<lines> lines) | (none)
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

- Role file not found
- Working directory doesn't exist
- Config file permissions error

**4** - Runtime error

- Content command failed
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

### Content Command Failed

```
Error: Content command failed: git diff --staged

Command output:
fatal: not a git repository

Exit code: 128
```

Exit code: 4

Task execution stops. Agent is not launched.

### Role File Not Found

```
Error: Role file not found: ./tasks/code-reviewer.md

Task 'code-review' references missing file.
Update task configuration or create the file.
```

Exit code: 3

### Invalid Task Configuration

```
Error: Task 'code-review' has invalid configuration.

Missing required field: 'role'

Update configuration: start config edit
```

Exit code: 1

## Notes

### Task vs Root Command Differences

**`start task <name>`:**

- Uses task's specific role (not `[context.system_prompt]`)
- Includes only documents listed in task's `documents` array
- Can run content_command for dynamic content
- Prompt template from task configuration
- `{instructions}` placeholder available

**`start` (root):**

- Uses `[context.system_prompt]` (if configured)
- Includes ALL context documents (required + optional)
- No content_command
- Simple document prompt concatenation
- No instructions placeholder

### Task Placeholders

Task prompt templates support these placeholders:

**Task-specific:**

- `{instructions}` - User's command-line arguments (or "None")
- `{content}` - Output from content_command (or empty string)

**Global:**

- `{model}` - Model name
- `{system_prompt}` - Role file contents
- `{date}` - Current timestamp (ISO 8601)
- `{file}` - Document file path (in document prompts only)

### Document Handling

**Document resolution:**

1. Task specifies: `documents = ["environment", "agents", "project"]`
2. Each name resolves to `[context.documents.<name>]` section in config
3. If document name not found in config → skip silently (no error)
4. If document file doesn't exist → skip with status display (same as `start` command)

**Document order in prompt:**

- Documents appear in order specified in task's `documents` array
- Document prompts appear BEFORE task's prompt template
- Example: environment prompts, agents prompts, then task prompt

### Role File Location

By convention, task role files are stored in:

- Global: `~/.config/start/tasks/*.md`
- Local (per-project): `./.start/tasks/*.md` or `./tasks/*.md`

Example task configuration:

```toml
[task.code-review]
role = "./tasks/code-reviewer.md"  # Project-specific role
# or
role = "~/.config/start/tasks/code-reviewer.md"  # Shared role
```

Role can also be inline:

```toml
[task.quick-review]
role = """
You are a code reviewer.
Focus on critical issues only.
"""
```

### Content Command Execution

**Working directory:**

- Defaults to current directory (`pwd`)
- Override with `--directory` flag
- Paths in config resolve relative to working directory

**Command execution:**

- Runs in shell (same as `bash -c "command"`)
- Environment variables inherited
- Timeout: None (command runs to completion)
- Exit code: Non-zero = error and task stops

**Common content commands:**

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

- Tasks involve dynamic content (content_command output)
- Users need visibility into what content is being analyzed
- Summary is already concise
- Use `--verbose` or `--debug` for more detail if needed

## See Also

- start(1) - Launch with context
- start-prompt(1) - Launch with custom prompt
- start-config(1) - Manage configuration
- task.md - Task configuration reference
