# start config task

## Name

start config task - Manage task configurations

## Synopsis

```bash
start config task list [scope]
start config task add [scope]
start config task test <name>
start config task edit [name] [scope]
start config task remove [name] [scope]
```

## Description

Manages predefined workflow task configurations in config files. Tasks define reusable AI-assisted workflows with optional system prompt overrides, automatic required context inclusion, and dynamic content from shell commands.

**Task management operations:**

- **list** - Display all configured tasks with details
- **add** - Add new task interactively
- **test** - Test task configuration and command execution
- **edit** - Modify existing task configuration
- **remove** - Delete task from configuration

**Note:** Per DR-017, tasks can be defined in both global and local configs, and are also loaded from asset library (`~/.config/start/assets/tasks/`). These commands manage user-defined tasks only (global or local config). Asset tasks are managed by `start update`.

## Task Configuration Structure

Tasks use the **[Unified Template Design (UTD)](../design/unified-template-design.md)** pattern for both system prompts and task prompts:

```toml
[tasks.git-diff-review]
alias = "gdr"
description = "Review staged git changes"

# System prompt override (UTD - all optional)
system_prompt_file = "~/.config/start/roles/code-reviewer.md"
system_prompt = """
{file}

Focus on security and performance.
"""

# Task prompt (UTD - at least one required)
command = "git diff --staged"
prompt = """
Review changes:

## Instructions
{instructions}

## Changes
```diff
{command}
```
"""

# Shell config (optional)
shell = "bash"
command_timeout = 10
```

**Metadata Fields:**

**alias** (optional)
: Short name for quick access. Must be unique across all tasks.

**description** (optional)
: Help text displayed in task list and help output.

**System Prompt Override (UTD):**

All fields optional. If omitted, uses global/local `[system_prompt]`.

- `system_prompt_file` - Path to role definition file
- `system_prompt_command` - Shell command for dynamic role content
- `system_prompt` - Template with `{file}` and `{command}` placeholders

**Task Prompt (UTD):**

At least one required:

- `file` - Path to prompt template file
- `command` - Shell command for dynamic content
- `prompt` - Template with `{file}`, `{command}`, `{instructions}` placeholders

**Additional Fields:**

- `shell` (optional) - Override global shell
- `command_timeout` (optional) - Override global timeout (seconds)

**Context Inclusion:**

Tasks automatically include **all contexts where `required = true`**. No `documents` array needed.

## Subcommands

### start config task list

Display all configured tasks with their details.

**Synopsis:**

```bash
start config task list          # Select scope interactively
start config task list global   # List global tasks only
start config task list local    # List local tasks only
start config task list merged   # Show merged view (assets + global + local)
```

**Behavior:**

Lists all tasks defined in the selected scope(s) with:

- Task name
- Alias
- Description
- System prompt override (yes/no)
- Task prompt type (file, command, inline, or combination)
- Source scope (asset, global, local, or override)

**Output (merged view):**

```
Configured tasks (merged):
═══════════════════════════════════════════════════════════

Asset tasks (4):
  code-review (cr)
    Review code for quality and best practices
    Source: ~/.config/start/assets/tasks/code-review.toml

  git-diff-review (gdr)
    Review staged git changes
    Source: ~/.config/start/assets/tasks/git-diff-review.toml

  comment-tidy (ct)
    Review and tidy code comments
    Source: ~/.config/start/assets/tasks/comment-tidy.toml

  doc-review (dr)
    Review and improve documentation
    Source: ~/.config/start/assets/tasks/doc-review.toml

Global tasks (1):
  security-review (sr)
    Security-focused code review
    System prompt: custom (code-reviewer.md + template)
    Task: command-based (git diff --staged)
    Source: ~/.config/start/config.toml

Local tasks (1):
  quick-help (qh)
    Quick help with instructions
    System prompt: default
    Task: inline prompt
    Source: ./.start/config.toml

Overridden tasks (1):
  code-review (cr) [local overrides asset]
    Project-specific code review
    System prompt: custom (project-reviewer.md)
    Task: combination (file + command)
    Source: ./.start/config.toml (overrides asset)
```

**Output (global only):**

```bash
start config task list global
```

```
Configured tasks (global):
═══════════════════════════════════════════════════════════

security-review (sr)
  Security-focused code review
  System prompt: custom (code-reviewer.md + template)
  Task: command-based (git diff --staged)
  Shell: bash
  Timeout: 10 seconds
```

**Output (local only):**

```bash
start config task list local
```

```
Configured tasks (local):
═══════════════════════════════════════════════════════════

code-review (cr)
  Project-specific code review (overrides asset task)
  System prompt: custom (project-reviewer.md)
  Task: combination (file + command)

quick-help (qh)
  Quick help with instructions
  System prompt: default
  Task: inline prompt
```

**No tasks configured:**

```
No tasks configured in global config.

Add tasks: start config task add global
View asset tasks: start config task list merged
```

**Exit codes:**

- 0 - Success (tasks listed)
- 1 - No config file exists
- 2 - Invalid scope argument

### start config task add

Interactively add a new task to the configuration.

**Synopsis:**

```bash
start config task add          # Select scope interactively
start config task add global   # Add to global config
start config task add local    # Add to local config
```

**Behavior:**

Prompts for task details and adds to the selected config file:

1. **Select scope** (if not provided)
   - global - Add to `~/.config/start/config.toml`
   - local - Add to `./.start/config.toml`

2. **Task name** (required)
   - Validation: lowercase alphanumeric with hyphens
   - Pattern: `/^[a-z0-9]+(-[a-z0-9]+)*$/`
   - Must be unique across all sources (assets + global + local)
   - Examples: `code-review`, `git-diff-review`, `my-task`

3. **Alias** (optional)
   - Short name for quick access
   - Same validation as task name
   - Must be globally unique

4. **Description** (optional)
   - Human-readable description
   - Press enter to skip

5. **System prompt override?** (yes/no, default: no)
   - If yes: Configure UTD fields (file, command, prompt)
   - If no: Uses global/local `[system_prompt]`

6. **Task prompt** (required)
   - Configure UTD fields (file, command, prompt)
   - At least one UTD field required

7. **Advanced options?** (yes/no, default: no)
   - Shell override
   - Command timeout

8. **Backup and save**
   - Backs up existing config to `config.YYYY-MM-DD-HHMMSS.toml`
   - Writes new task to config
   - Shows success message

**Interactive flow (simple task with inline prompt):**

```
Add new task
─────────────────────────────────────────────────

Select scope:
  1) global (all projects)
  2) local (this project only)

Scope [1-2]: 1

Task name: quick-help
Alias (optional): qh
Description (optional): Quick help with instructions

Override system prompt? [y/N]: n
✓ Will use default system prompt

Task prompt:

Content source:
  1) File path
  2) Command
  3) Inline prompt
  4) Combination

Select [1-4]: 3

Prompt template: Help me with: {instructions}
✓ Valid template (uses {instructions} placeholder)

Advanced options? [y/N]: n

Backing up config to config.2025-01-06-101234.toml...
✓ Backup created

Saving task 'quick-help' to ~/.config/start/config.toml...
✓ Task added successfully

Use 'start config task list global' to see all tasks.
Use 'start task quick-help "your question"' to run.
```

**Interactive flow (complex task with system prompt override and command):**

```
Add new task
─────────────────────────────────────────────────

Select scope:
  1) global (all projects)
  2) local (this project only)

Scope [1-2]: 1

Task name: git-diff-review
Alias (optional): gdr
Description (optional): Review staged git changes

Override system prompt? [y/N]: y

System prompt configuration:

Content source:
  1) File path
  2) Command
  3) Inline prompt
  4) Combination

Select [1-4]: 4

System prompt file: ~/.config/start/roles/code-reviewer.md
✓ File exists

System prompt template: {file}\n\nFocus on security and performance.
✓ Valid template (uses {file} placeholder)

Task prompt:

Content source:
  1) File path
  2) Command
  3) Inline prompt
  4) Combination

Select [1-4]: 2

Command: git diff --staged
✓ Valid command

Prompt template: Review changes:\n\n## Instructions\n{instructions}\n\n## Changes\n```diff\n{command}\n```
✓ Valid template (uses {instructions} and {command} placeholders)

Advanced options? [y/N]: y

Shell override (or enter for default): bash
Command timeout in seconds (or enter for default): 10

Backing up config to config.2025-01-06-101345.toml...
✓ Backup created

Saving task 'git-diff-review' to ~/.config/start/config.toml...
✓ Task added successfully

Use 'start config task list global' to see all tasks.
Use 'start task git-diff-review "focus on security"' to run.
```

**Resulting config (simple task):**

```toml
[tasks.quick-help]
alias = "qh"
description = "Quick help with instructions"
prompt = "Help me with: {instructions}"
```

**Resulting config (complex task):**

```toml
[tasks.git-diff-review]
alias = "gdr"
description = "Review staged git changes"

system_prompt_file = "~/.config/start/roles/code-reviewer.md"
system_prompt = """
{file}

Focus on security and performance.
"""

command = "git diff --staged"
prompt = """
Review changes:

## Instructions
{instructions}

## Changes
```diff
{command}
```
"""

shell = "bash"
command_timeout = 10
```

**Exit codes:**

- 0 - Success (task added)
- 1 - Validation error (invalid name, duplicate task, invalid configuration)
- 2 - Scope error (invalid scope, local config directory doesn't exist)
- 3 - File system error (cannot write config, backup failed)

**Error handling:**

**Invalid task name:**

```
Task name: My-Task
✗ Invalid task name. Use lowercase alphanumeric with hyphens.
  Examples: code-review, git-diff-review, my-task

Task name: my-task
✓ Valid name
```

**Duplicate task:**

```
Task name: code-review
⚠ Warning: Task 'code-review' exists in asset library.
  Your task will override the asset task.

Continue? [y/N]: y
```

**Duplicate alias:**

```
Alias (optional): gdr
✗ Alias 'gdr' already used by task 'git-diff-review'.

Alias (optional): gr
✓ Unique alias
```

**No UTD fields for task prompt:**

```
Content source:
  1) File path
  2) Command
  3) Inline prompt
  4) Combination

Select [1-4]: 1

File path:
✗ At least one task prompt field is required (file, command, or prompt).
  Press enter to return to content source selection.
```

**Command doesn't exist (warning only):**

```
Command: nonexistent-command --flag
⚠ Warning: Command may not be available.
  Test: command -v nonexistent-command

Continue anyway? [y/N]: y
```

**Local config directory doesn't exist:**

```
Select scope:
  1) global (all projects)
  2) local (this project only)

Scope [1-2]: 2

✗ Local config directory doesn't exist: ./.start/
  Create it first: mkdir -p ./.start

Or add to global config instead.
```

Exit code: 2

### start config task test

Test task configuration, file availability, and command execution.

**Synopsis:**

```bash
start config task test <name>
```

**Behavior:**

Validates task configuration without executing the task. Performs checks:

1. **System prompt validation** (if override configured)
   - File exists (if `system_prompt_file` present)
   - Command executes (if `system_prompt_command` present)
   - Template uses valid placeholders

2. **Task prompt validation**
   - File exists (if `file` present)
   - Command executes (if `command` present)
   - Template uses valid placeholders (`{file}`, `{command}`, `{instructions}`)

3. **Configuration validation**
   - At least one task prompt UTD field present
   - Shell and timeout settings valid
   - Alias unique (if configured)

**Does NOT run the full task** - only validates configuration and tests commands in isolation.

**Output (simple task, success):**

```
Testing task: quick-help
─────────────────────────────────────────────────

Configuration:
  Scope: global
  Alias: qh
  Description: Quick help with instructions
  System prompt: default (uses global/local)

Task prompt:
  Type: Inline prompt
  Prompt: Help me with: {instructions}
  ✓ Valid template
  ✓ Uses {instructions} placeholder

✓ Task 'quick-help' is configured correctly
```

**Output (complex task, success):**

```
Testing task: git-diff-review
─────────────────────────────────────────────────

Configuration:
  Scope: global
  Alias: gdr
  Description: Review staged git changes
  System prompt: custom override

System prompt override:
  File: ~/.config/start/roles/code-reviewer.md
  Resolved: /Users/grant/.config/start/roles/code-reviewer.md
  ✓ File exists (847 bytes)

  Template:
    {file}

    Focus on security and performance.
  ✓ Valid template
  ✓ Uses {file} placeholder

Task prompt:
  Type: Command-based
  Shell: bash
  Timeout: 10 seconds
  Command: git diff --staged
  ✓ Command executed successfully (1,234 bytes output)

  Template:
    Review changes:

    ## Instructions
    {instructions}

    ## Changes
    ```diff
    {command}
    ```
  ✓ Valid template
  ✓ Uses {instructions} and {command} placeholders

✓ Task 'git-diff-review' is configured correctly
```

**Output (file not found):**

```
Testing task: broken-task
─────────────────────────────────────────────────

Configuration:
  Scope: local
  Alias: bt
  Description: Broken task example
  System prompt: custom override

System prompt override:
  File: ./roles/missing.md
  Resolved: /Users/grant/Projects/myapp/roles/missing.md
  ✗ File not found

  Template: {file}
  ✓ Valid template

Task prompt:
  Type: Inline prompt
  Prompt: Do something
  ✓ Valid template

✗ Task 'broken-task' has errors
  System prompt file not found
```

**Output (command failed):**

```
Testing task: bad-command
─────────────────────────────────────────────────

Configuration:
  Scope: global
  Alias: bc
  Description: Task with broken command
  System prompt: default

Task prompt:
  Type: Command-based
  Shell: bash
  Timeout: 30 seconds
  Command: nonexistent-command --flag
  ✗ Command failed (exit code 127)
  Error: nonexistent-command: command not found

  Template: Output: {command}
  ✓ Valid template
  ✓ Uses {command} placeholder

✗ Task 'bad-command' has errors
  Task command execution will fail at runtime
```

**Output (configuration error):**

```
Testing task: invalid
─────────────────────────────────────────────────

Configuration:
  Scope: global
  Description: Invalid task
  System prompt: default

Task prompt:
  ✗ No UTD fields present (no file, command, or prompt)

✗ Task 'invalid' has configuration errors
  At least one task prompt field required
  Fix configuration: start config task edit invalid global
```

**Verbose output:**

```bash
start config task test git-diff-review --verbose
```

```
Testing task: git-diff-review
─────────────────────────────────────────────────

Loading configuration...
  Config file: ~/.config/start/config.toml
  Task section: [tasks.git-diff-review]

Configuration details:
  Name: git-diff-review
  Scope: global
  Alias: gdr
  Description: Review staged git changes

System prompt override:
  File field: ~/.config/start/roles/code-reviewer.md
  File resolution:
    Home expansion: /Users/grant/.config/start/roles/code-reviewer.md
    ✓ File exists
    Size: 847 bytes
    Modified: 2025-01-05 10:23:15

  Template field:
    {file}

    Focus on security and performance.
  Placeholders found: {file}
  ✓ Valid placeholder usage

Task prompt:
  Command field: git diff --staged
  Shell: bash (configured)
  Timeout: 10 seconds (configured)

  Command execution:
    Working directory: /Users/grant/Projects/myapp
    ✓ Executed successfully
    Output size: 1,234 bytes
    Exit code: 0

  Prompt field:
    Review changes:

    ## Instructions
    {instructions}

    ## Changes
    ```diff
    {command}
    ```

  Placeholders found: {instructions}, {command}
  ✓ Valid placeholder usage
  ✓ {instructions} - task-specific placeholder
  ✓ {command} - matches task command field

Required contexts:
  ✓ Auto-includes contexts where required = true
  Found: environment, index, agents (3 contexts)

✓ Task 'git-diff-review' is configured correctly
```

**Exit codes:**

- 0 - Success (task valid, files exist, commands succeed)
- 1 - Configuration error (invalid configuration)
- 2 - Task not found in config
- 3 - File not found (config valid but file missing)
- 4 - Command failed (config valid but command execution failed)

**Error handling:**

**Task not in config:**

```
Error: Task 'nonexistent' not found in configuration.

Use 'start config task list' to see available tasks.
Use 'start config task add' to add a new task.
```

Exit code: 2

**Multiple errors:**

```
Testing task: broken
─────────────────────────────────────────────────

Configuration:
  ✗ No task prompt UTD fields
  ⚠ Unknown placeholder {unknown} in template

System prompt override:
  ✗ File not found: ./missing.md

✗ Task 'broken' has multiple errors:
  - Invalid configuration (no task prompt source)
  - System prompt file not found
  - Invalid placeholder usage
```

Exit code: 1 (configuration errors take precedence)

### start config task edit

Edit task configuration interactively.

**Synopsis:**

```bash
start config task edit                  # Select task and scope
start config task edit <name>           # Select scope for named task
start config task edit <name> global    # Edit in global config
start config task edit <name> local     # Edit in local config
```

**Behavior:**

**Without task name (interactive selection):**

Shows list of configured tasks for selection:

```bash
start config task edit
```

Output:

```
Edit task
─────────────────────────────────────────────────

Select task to edit:

Global tasks:
  1) security-review (sr)
  2) quick-help (qh)

Local tasks:
  3) code-review (cr) - overrides asset
  4) project-review (pr)

Select [1-4] (or 'q' to quit): 1

(continues to interactive edit flow for 'security-review' in global config)
```

**With task name only:**

If task exists in only one user config, edits that config. If exists in multiple, prompts for scope:

```bash
start config task edit quick-help
```

Task exists in global only:

```
Editing task 'quick-help' in global config...
(continues to interactive edit flow)
```

Task exists in both:

```
Task 'code-review' exists in multiple scopes.

Select scope to edit:
  1) global - Custom security review
  2) local - Project-specific review (overrides asset)

Select [1-2]: 2
(continues to interactive edit flow for local)
```

**With task name and scope:**

Interactive prompts to edit specific task. Shows current values as defaults - press enter to keep current value.

1. **Alias** - Current value shown in brackets
2. **Description** - Current value shown in brackets
3. **System prompt override changes** - Modify or remove
4. **Task prompt changes** - Modify file, command, or prompt
5. **Advanced options** - Shell, timeout
6. **Backup and save** - Backs up to `config.YYYY-MM-DD-HHMMSS.toml`

**Interactive flow:**

```
Edit task: quick-help (global)
─────────────────────────────────────────────────

Current configuration:
  Alias: qh
  Description: Quick help with instructions
  System prompt: default
  Task prompt: inline (Help me with: {instructions})

Press enter to keep current value, or type new value:

Alias [qh]:
Description [Quick help with instructions]: Quick help for any question

System prompt override? [y/N]: n

Task prompt template [Help me with: {instructions}]:

Advanced options? [y/N]: n

No changes detected.

Task 'quick-help' not modified.
```

**Interactive flow (with changes):**

```
Edit task: git-diff-review (global)
─────────────────────────────────────────────────

Current configuration:
  Alias: gdr
  Description: Review staged git changes
  System prompt: custom (code-reviewer.md + template)
  Task prompt: command-based (git diff --staged)
  Shell: bash
  Timeout: 10 seconds

Press enter to keep current value, or type new value:

Alias [gdr]:
Description [Review staged git changes]:

Keep system prompt override? [Y/n]: y
System prompt file [~/.config/start/roles/code-reviewer.md]:
System prompt template [{file}\n\nFocus on security and performance.]:

Keep command for task prompt? [Y/n]: y
Task command [git diff --staged]:
Task prompt template [...]: Review changes:\n\n{instructions}\n\n{command}

Advanced options? [y/N]: y
Shell [bash]:
Timeout in seconds [10]: 15

Backing up config to config.2025-01-06-102345.toml...
✓ Backup created

Saving changes to ~/.config/start/config.toml...
✓ Task 'git-diff-review' updated successfully

Use 'start config task list global' to see changes.
Use 'start config task test git-diff-review' to validate.
```

**Exit codes:**

- 0 - Success (task edited or no changes)
- 1 - Validation error (invalid configuration)
- 2 - Task not found or scope error
- 3 - File system error (cannot write config, backup failed)

**Error handling:**

**Task not found:**

```
Error: Task 'nonexistent' not found in configuration.

Use 'start config task list' to see available tasks.
Use 'start config task add' to add a new task.
```

Exit code: 2

**Cannot edit asset task:**

```
Error: Task 'code-review' is from asset library.

Asset tasks cannot be edited directly.
To customize, create an override in global or local config:
  start config task add global

Or remove from assets: Remove the file from ~/.config/start/assets/tasks/
```

Exit code: 2

**Invalid template:**

```
Task prompt template [Help with: {instructions}]: Invalid {unknown} text
⚠ Warning: Unknown placeholder {unknown}
  Valid placeholders: {file}, {command}, {instructions}, {model}, {date}

Continue anyway? [y/N]: n

Task prompt template [Help with: {instructions}]:
✓ Valid template
```

**No changes made:**

```
No changes detected.

Task 'quick-help' not modified.
```

Exit code: 0 (no backup created, no write)

### start config task remove

Remove task from configuration.

**Synopsis:**

```bash
start config task remove                  # Select task and scope
start config task remove <name>           # Select scope for named task
start config task remove <name> global    # Remove from global config
start config task remove <name> local     # Remove from local config
```

**Behavior:**

**Without task name:**
Shows list of configured tasks for selection:

```
Remove task
─────────────────────────────────────────────────

Select task to remove:

Global tasks:
  1) security-review (sr)
  2) quick-help (qh)

Local tasks:
  3) code-review (cr) - overrides asset
  4) project-review (pr)

Select [1-4] (or 'q' to quit): 2

Remove task 'quick-help' from global config? [y/N]: y

Backing up config to config.2025-01-06-103012.toml...
✓ Backup created

Removing task 'quick-help' from ~/.config/start/config.toml...
✓ Task 'quick-help' removed successfully

Use 'start config task list global' to see remaining tasks.
```

**With task name only:**
If task exists in only one user config, removes from that config. If exists in multiple, prompts for scope:

```bash
start config task remove quick-help
```

Task exists in global only:

```
Remove task 'quick-help' from global config? [y/N]: y

Backing up config to config.2025-01-06-103045.toml...
✓ Backup created

Removing task 'quick-help' from ~/.config/start/config.toml...
✓ Task 'quick-help' removed successfully

Use 'start config task list global' to see remaining tasks.
```

Task exists in multiple:

```
Task 'code-review' exists in multiple scopes.

Select scope to remove from:
  1) global - Custom security review
  2) local - Project override (removes override, asset task remains)
  3) both

Select [1-3]: 2

Remove task 'code-review' from local config? [y/N]: y
(continues with removal from local)
```

**Removing local override (restores asset task):**

```bash
start config task remove code-review local
```

Output:

```
⚠ Note: 'code-review' overrides an asset task.
  Removing this will restore the asset task behavior.

Remove task 'code-review' from local config? [y/N]: y

Backing up config to config.2025-01-06-103123.toml...
✓ Backup created

Removing task 'code-review' from ./.start/config.toml...
✓ Task 'code-review' removed successfully
✓ Asset task 'code-review' is now active

Use 'start task code-review' to run asset version.
Use 'start config task list merged' to see all tasks.
```

**Declining confirmation:**

```
Remove task 'quick-help' from global config? [y/N]: n

Task 'quick-help' not removed.
```

Exit code: 0

**Exit codes:**

- 0 - Success (task removed, or user declined)
- 1 - No tasks configured
- 2 - Task not found, scope error, or asset task
- 3 - File system error (cannot write config, backup failed)

**Error handling:**

**Task not found:**

```
Error: Task 'nonexistent' not found in configuration.

Use 'start config task list' to see available tasks.
```

Exit code: 2

**Asset task (cannot remove):**

```
Error: Task 'code-review' is from asset library.

Asset tasks cannot be removed via this command.
To hide asset task, create empty override in local config,
or remove from asset directory: ~/.config/start/assets/tasks/

Use 'start config task list merged' to see all sources.
```

Exit code: 2

**No tasks configured:**

```
No tasks configured in global config.

Use 'start config task add global' to add a task.
View asset tasks: start config task list merged
```

Exit code: 1

**Backup failed:**

```
Remove task 'quick-help' from global config? [y/N]: y

Backing up config to config.2025-01-06-103156.toml...
✗ Failed to backup config: permission denied

Existing config preserved at: ~/.config/start/config.toml
Task not removed.
```

Exit code: 3

## Global Flags

These flags work on all `start config task` subcommands where applicable.

**--help**, **-h**
: Show help for the subcommand.

**--verbose**, **-v**
: Verbose output. Shows config file paths and additional details.

**--debug**
: Debug mode. Shows all internal operations.

## Examples

### List All Tasks (Merged View)

```bash
start config task list merged
```

Show all tasks from assets, global, and local configs.

### List User Tasks Only

```bash
start config task list global
start config task list local
```

### Add Task to Global Config

```bash
start config task add global
```

### Add Task to Local Config

```bash
start config task add local
```

### Test Task

```bash
start config task test git-diff-review
```

Verify task configuration, file availability, and command execution.

### Edit Task

```bash
start config task edit git-diff-review global
```

### Remove Task

```bash
start config task remove quick-help global
```

### Interactive Task Selection

```bash
start config task edit
```

Shows list of all tasks to choose from.

## Files

**~/.config/start/config.toml**
: Global configuration file containing user-defined task definitions.

**./.start/config.toml**
: Local project configuration file containing project-specific tasks.

**~/.config/start/assets/tasks/*.toml**
: Asset library tasks (updated via `start update`). Cannot be edited directly via `start config task` commands.

Tasks are merged from all three sources: assets + global + local. User tasks (global/local) take precedence over asset tasks with the same name.

## Error Handling

### No Configuration File

```
Error: No configuration found at ~/.config/start/config.toml

Run 'start init' to create initial configuration.
```

Exit code: 1

### Invalid TOML Syntax

```
Error: Configuration file has invalid syntax.

File: ~/.config/start/config.toml
Line 123: invalid TOML syntax

Fix the configuration file or restore from backup.
```

Exit code: 1

## Notes

### Task Merge Behavior

Tasks are discovered from three sources:

1. **Asset library:** `~/.config/start/assets/tasks/*.toml` (managed by `start update`)
2. **Global config:** `~/.config/start/config.toml` (managed by `start config task`)
3. **Local config:** `./.start/config.toml` (managed by `start config task`)

**Precedence:**
1. Local tasks override global tasks (same name)
2. Global tasks override asset tasks (same name)
3. Local tasks override asset tasks (same name)

**Result:** User-defined tasks take precedence, allowing customization of asset tasks.

### Task Source Labeling

In task lists, tasks are labeled by source:

- `[asset]` - From `~/.config/start/assets/tasks/`
- `[global]` - From `~/.config/start/config.toml`
- `[local]` - From `./.start/config.toml`
- `[local overrides asset]` - Local task overrides asset task

### Unified Template Design (UTD)

Tasks use UTD pattern for both system prompts and task prompts:

**System prompt override (all optional):**
```toml
[tasks.code-review]
system_prompt_file = "~/.config/start/roles/code-reviewer.md"
system_prompt_command = "date"
system_prompt = """
{file}

Current time: {command}
"""
```

**Task prompt (at least one required):**
```toml
[tasks.git-diff]
file = "./prompts/diff-template.md"
command = "git diff --staged"
prompt = """
{file}

## Changes
{command}

## Instructions
{instructions}
"""
```

See [UTD documentation](../design/unified-template-design.md) for complete details.

### Placeholders

**System prompt templates:**
- `{file}` - Content from `system_prompt_file`
- `{command}` - Output from `system_prompt_command`
- `{model}`, `{date}` - Global placeholders

**Task prompt templates:**
- `{file}` - Content from task `file`
- `{command}` - Output from task `command`
- `{instructions}` - User's command-line arguments (or "None")
- `{model}`, `{date}` - Global placeholders

### Required Context Auto-Inclusion

Tasks automatically include **all contexts where `required = true`**. There is no `documents` array in task configuration.

**Rationale:**
- Ensures critical context (like AGENTS.md) is always present
- Simplifies task configuration
- Users control what's "required" via context configuration

**Example:**
```toml
# In config
[context.agents]
file = "./AGENTS.md"
required = true

[tasks.code-review]
# No documents array needed
# agents context is automatically included
prompt = "Review this code"
```

### Shell Configuration

Tasks can override the global shell setting:

```toml
[tasks.git-diff]
command = "git diff --staged"
shell = "bash"
command_timeout = 10
```

Supported shells: bash, sh, zsh, fish, node, python, ruby, etc.

See [UTD shell configuration](../design/unified-template-design.md#shell-configuration) for complete list.

### Asset vs User Tasks

**Asset tasks:**
- Provided by `start` as defaults
- Located in `~/.config/start/assets/tasks/`
- Updated via `start update`
- Cannot be edited directly via `start config task`
- Can be overridden by creating user task with same name

**User tasks:**
- Created via `start config task add`
- Located in global or local config
- Fully customizable
- Take precedence over asset tasks

**To customize asset task:**
1. Run `start config task add global` (or local)
2. Use same name as asset task
3. Configure as desired
4. Your task overrides the asset task

## See Also

- start-task(1) - Run predefined tasks
- start-config(1) - Manage configuration files
- start-config-agent(1) - Manage AI agents
- start-config-context(1) - Manage context documents
- start-config-role(1) - Manage system prompts
- start-update(1) - Update asset library
- DR-009 - Task structure and placeholders
- DR-010 - Default task definitions
- DR-017 - CLI command reorganization
