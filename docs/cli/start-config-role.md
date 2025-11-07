# start config role

## Name

start config role - Manage system prompt configuration

## Synopsis

```bash
start config role show [scope]
start config role edit [scope]
start config role remove [scope]
start config role test
```

## Description

Manages role (system prompt) configuration in config files. Roles define the AI agent's persona and behavior. Roles are passed to agents via the `{role}` and `{role_file}` placeholders in agent commands.

**Role management operations:**

- **show** - Display current role configurations
- **edit** - Modify a role (create if doesn't exist)
- **remove** - Remove a role configuration
- **test** - Verify role configuration and file availability
- **list** - List all configured roles

**Note:** Roles use the **[Unified Template Design (UTD)](../design/unified-template-design.md)** pattern. Roles are defined as `[roles.<name>]` sections in both global and local configs. Global + local roles are combined; local overrides global for same role name.

## Role Configuration Structure

Roles are defined using the **[Unified Template Design (UTD)](../design/unified-template-design.md)** pattern:

```toml
[roles.code-reviewer]
description = "Expert code reviewer"
file = "~/.config/start/roles/code-reviewer.md"
prompt = """
{file}

Focus on security and performance.
"""
```

**UTD Fields (at least one required):**

**file** (optional)
: Path to system prompt file. Supports `~` expansion and relative paths.

**command** (optional)
: Shell command to execute for dynamic content. Output replaces `{command}` placeholder.

**prompt** (optional)
: Template text with `{file}` and `{command}` placeholders.

**Additional Fields:**

**shell** (optional)
: Override global shell for command execution.

**command_timeout** (optional)
: Override global timeout (seconds) for command execution.

**Merge behavior:**

Local `[system_prompt]` section **completely replaces** global section. If local section is missing, uses global section.

## Subcommands

### start config role show

Display current system prompt configuration.

**Synopsis:**

```bash
start config role show          # Select scope interactively
start config role show global   # Show global system prompt only
start config role show local    # Show local system prompt only
start config role show merged   # Show effective system prompt (with override info)
```

**Behavior:**

Displays system prompt configuration from the selected scope(s) with:

- Scope (global, local, or merged)
- Source type (file, command, inline, or combination)
- File path (if configured)
- Command (if configured)
- Prompt template (if configured)
- Shell and timeout overrides (if configured)

**Output (merged view):**

```
System prompt configuration (merged):
═══════════════════════════════════════════════════════════

Effective configuration:
  Source: local (overrides global)
  Type: File with template

File:
  Path: ./ROLE.md
  Resolved: /Users/grant/Projects/myapp/ROLE.md
  ✓ File exists (1,234 bytes)

Prompt template:
  {file}

  Additional context: Focus on code quality.

Global configuration (overridden):
  Source: global
  Type: File only
  File: ~/.config/start/roles/default.md
```

**Output (global only):**

```bash
start config role show global
```

```
System prompt configuration (global):
═══════════════════════════════════════════════════════════

File:
  Path: ~/.config/start/roles/default.md
  Resolved: /Users/grant/.config/start/roles/default.md
  ✓ File exists (847 bytes)

Prompt: (file content only, no template)
```

**Output (local only):**

```bash
start config role show local
```

```
System prompt configuration (local):
═══════════════════════════════════════════════════════════

File:
  Path: ./ROLE.md
  Resolved: /Users/grant/Projects/myapp/ROLE.md
  ✓ File exists (1,234 bytes)

Prompt template:
  {file}

  Additional context: Focus on code quality.
```

**Output (inline prompt):**

```
System prompt configuration (global):
═══════════════════════════════════════════════════════════

Type: Inline prompt

Prompt:
  You are an expert code reviewer.
  Focus on security and performance.
```

**Output (command-based):**

```
System prompt configuration (local):
═══════════════════════════════════════════════════════════

Type: Command-based

Command:
  Shell: bash
  Timeout: 5 seconds
  Command: git log -1 --format='%s'

Prompt template:
  You are a code reviewer.
  Current commit: {command}
```

**No system prompt configured:**

```
No system prompt configured in global config.

Configure: start config role edit global
```

**Exit codes:**

- 0 - Success (system prompt shown)
- 1 - No system prompt configured
- 2 - Invalid scope argument

### start config role edit

Edit or create system prompt configuration interactively.

**Synopsis:**

```bash
start config role edit          # Select scope interactively
start config role edit global   # Edit global system prompt
start config role edit local    # Edit local system prompt
```

**Behavior:**

Prompts for system prompt configuration and updates the selected config file:

1. **Select scope** (if not provided)
   - global - Edit `~/.config/start/config.toml`
   - local - Edit `./.start/config.toml`

2. **Content source** (choose one or more)
   - File path (static system prompt document)
   - Command (dynamic content)
   - Inline prompt text
   - At least one is required

3. **Prompt template** (required if file or command specified)
   - Template with `{file}` and/or `{command}` placeholders
   - Or inline prompt text

4. **Advanced options?** (yes/no, default: no)
   - Shell override
   - Command timeout

5. **Backup and save**
   - Backs up existing config to `config.YYYY-MM-DD-HHMMSS.toml`
   - Writes system prompt to config
   - Shows success message

**Interactive flow (create from file):**

```
Edit system prompt
─────────────────────────────────────────────────

Select scope:
  1) global (all projects)
  2) local (this project only)

Scope [1-2]: 1

Current configuration: (none)

Content source:
  1) File path
  2) Command
  3) Inline prompt
  4) Combination

Select [1-4]: 1

File path: ~/.config/start/roles/code-reviewer.md
✓ File exists

Use prompt template to frame file content? [y/N]: n
✓ Will use file content directly

Advanced options? [y/N]: n

Backing up config to config.2025-01-06-111234.toml...
✓ Backup created

Saving system prompt to ~/.config/start/config.toml...
✓ System prompt configured successfully

Use 'start config role show global' to verify.
Use 'start config role test' to validate.
```

**Interactive flow (file with template):**

```
Edit system prompt
─────────────────────────────────────────────────

Select scope:
  1) global (all projects)
  2) local (this project only)

Scope [1-2]: 2

Current configuration:
  File: ./ROLE.md
  Prompt: (file only, no template)

Content source:
  1) File path
  2) Command
  3) Inline prompt
  4) Combination

Select [1-4]: 1

File path [./ROLE.md]:

Use prompt template to frame file content? [y/N]: y

Prompt template: {file}\n\nAdditional context: Focus on code quality.
✓ Valid template (uses {file} placeholder)

Advanced options? [y/N]: n

Backing up config to config.2025-01-06-111345.toml...
✓ Backup created

Saving system prompt to ./.start/config.toml...
✓ System prompt updated successfully

Use 'start config role show local' to verify.
```

**Interactive flow (inline prompt):**

```
Edit system prompt
─────────────────────────────────────────────────

Select scope:
  1) global (all projects)
  2) local (this project only)

Scope [1-2]: 1

Current configuration: (none)

Content source:
  1) File path
  2) Command
  3) Inline prompt
  4) Combination

Select [1-4]: 3

Prompt text:
You are an expert code reviewer.
Focus on security and performance.

(Press Ctrl+D or enter empty line to finish)

✓ Valid prompt (87 characters)

Advanced options? [y/N]: n

Backing up config to config.2025-01-06-111456.toml...
✓ Backup created

Saving system prompt to ~/.config/start/config.toml...
✓ System prompt configured successfully
```

**Interactive flow (combination - file + command):**

```
Edit system prompt
─────────────────────────────────────────────────

Select scope:
  1) global (all projects)
  2) local (this project only)

Scope [1-2]: 2

Current configuration: (none)

Content source:
  1) File path
  2) Command
  3) Inline prompt
  4) Combination

Select [1-4]: 4

File path: ./ROLE.md
✓ File exists

Add command for dynamic content? [y/N]: y

Command: git log -1 --format='%s'
✓ Valid command

Prompt template: {file}\n\nCurrent commit: {command}
✓ Valid template (uses {file} and {command})

Advanced options? [y/N]: y

Shell override (or enter for default): bash
Command timeout in seconds (or enter for default): 5

Backing up config to config.2025-01-06-111567.toml...
✓ Backup created

Saving system prompt to ./.start/config.toml...
✓ System prompt configured successfully
```

**Resulting config (simple file):**

```toml
[system_prompt]
file = "~/.config/start/roles/code-reviewer.md"
```

**Resulting config (file with template):**

```toml
[system_prompt]
file = "./ROLE.md"
prompt = """
{file}

Additional context: Focus on code quality.
"""
```

**Resulting config (inline prompt):**

```toml
[system_prompt]
prompt = """
You are an expert code reviewer.
Focus on security and performance.
"""
```

**Resulting config (combination):**

```toml
[system_prompt]
file = "./ROLE.md"
command = "git log -1 --format='%s'"
prompt = """
{file}

Current commit: {command}
"""
shell = "bash"
command_timeout = 5
```

**Exit codes:**

- 0 - Success (system prompt configured)
- 1 - Validation error (invalid configuration)
- 2 - Scope error (invalid scope, local config directory doesn't exist)
- 3 - File system error (cannot write config, backup failed)

**Error handling:**

**No UTD fields:**

```
Content source:
  1) File path
  2) Command
  3) Inline prompt
  4) Combination

Select [1-4]: 1

File path:
✗ At least one content source is required (file, command, or prompt).
  Press enter to return to content source selection.
```

**File doesn't exist (warning only):**

```
File path: ./MISSING.md
⚠ Warning: File does not exist: ./MISSING.md
  System prompt will fail at runtime if file is not found.

Continue anyway? [y/N]: y
```

**Invalid placeholder:**

```
Prompt template: Invalid {unknown} text
⚠ Warning: Unknown placeholder {unknown}
  Valid placeholders: {file}, {command}

Continue anyway? [y/N]: n

Prompt template:
✓ Valid template
```

**Local config directory doesn't exist:**

```
Select scope:
  1) global (all projects)
  2) local (this project only)

Scope [1-2]: 2

✗ Local config directory doesn't exist: ./.start/
  Create it first: mkdir -p ./.start

Or configure global system prompt instead.
```

Exit code: 2

**Backup failed:**

```
Backing up config to config.2025-01-06-111234.toml...
✗ Failed to backup config: permission denied

Existing config preserved at: ~/.config/start/config.toml
System prompt not configured.
```

Exit code: 3

### start config role remove

Remove system prompt configuration.

**Synopsis:**

```bash
start config role remove          # Select scope interactively
start config role remove global   # Remove from global config
start config role remove local    # Remove from local config
```

**Behavior:**

Removes `[system_prompt]` section from the selected config file. After removal, agents will run without system prompts (or use global if removing from local).

**Interactive flow:**

```bash
start config role remove
```

Output:

```
Remove system prompt
─────────────────────────────────────────────────

Select scope:
  1) global (all projects)
  2) local (this project only)

Scope [1-2]: 2

Current configuration (local):
  File: ./ROLE.md
  Prompt template: {file}\n\nFocus on code quality.

Remove system prompt from local config? [y/N]: y

Backing up config to config.2025-01-06-112012.toml...
✓ Backup created

Removing [system_prompt] from ./.start/config.toml...
✓ System prompt removed successfully

Global system prompt will now be used (if configured).

Use 'start config role show merged' to see effective configuration.
```

**Direct scope removal:**

```bash
start config role remove global
```

Output:

```
Current configuration (global):
  File: ~/.config/start/roles/default.md

Remove system prompt from global config? [y/N]: y

Backing up config to config.2025-01-06-112045.toml...
✓ Backup created

Removing [system_prompt] from ~/.config/start/config.toml...
✓ System prompt removed successfully
⚠ No system prompt configured

Agents will run without system prompts.

Configure: start config role edit global
```

**Removing local (reverts to global):**

```bash
start config role remove local
```

Output:

```
Current configuration (local):
  File: ./ROLE.md

⚠ Note: Removing local system prompt will revert to global configuration.

Remove system prompt from local config? [y/N]: y

Backing up config to config.2025-01-06-112123.toml...
✓ Backup created

Removing [system_prompt] from ./.start/config.toml...
✓ System prompt removed successfully
✓ Now using global system prompt

Global configuration:
  File: ~/.config/start/roles/default.md

Use 'start config role show merged' to verify.
```

**Declining confirmation:**

```
Remove system prompt from global config? [y/N]: n

System prompt not removed.
```

Exit code: 0

**Exit codes:**

- 0 - Success (system prompt removed, or user declined)
- 1 - No system prompt configured
- 2 - Invalid scope
- 3 - File system error (cannot write config, backup failed)

**Error handling:**

**No system prompt configured:**

```
Error: No system prompt configured in global config.

Configure: start config role edit global
```

Exit code: 1

**Backup failed:**

```
Remove system prompt from global config? [y/N]: y

Backing up config to config.2025-01-06-112156.toml...
✗ Failed to backup config: permission denied

Existing config preserved at: ~/.config/start/config.toml
System prompt not removed.
```

Exit code: 3

### start config role test

Test system prompt configuration and file availability.

**Synopsis:**

```bash
start config role test
```

**Behavior:**

Validates effective system prompt configuration (merged global + local). Performs checks:

1. **File availability** (if `file` field present)
   - Checks if file exists at specified path
   - Resolves `~` and relative paths
   - Reports: found (with resolved path and size) or not found

2. **Command execution** (if `command` field present)
   - Executes command in configured shell
   - Reports: success (with output size) or failure (with error)
   - Does NOT display command output (may be large)

3. **Configuration validation**
   - At least one UTD field present (file, command, or prompt)
   - Prompt template uses valid placeholders (`{file}`, `{command}`)
   - Unknown placeholders detected (likely typos)
   - Shell and timeout settings valid

**Does NOT pass system prompt to any agent** - only validates configuration.

**Output (file-based, success):**

```
Testing system prompt configuration...
─────────────────────────────────────────────────

Effective configuration:
  Scope: global
  Type: File-based

File:
  Path: ~/.config/start/roles/code-reviewer.md
  Resolved: /Users/grant/.config/start/roles/code-reviewer.md
  ✓ File exists (847 bytes)
  Modified: 2025-01-05 10:23:15

Prompt: (file content only, no template)

✓ System prompt is configured correctly
```

**Output (file with template, success):**

```
Testing system prompt configuration...
─────────────────────────────────────────────────

Effective configuration:
  Scope: local (overrides global)
  Type: File with template

File:
  Path: ./ROLE.md
  Resolved: /Users/grant/Projects/myapp/ROLE.md
  ✓ File exists (1,234 bytes)
  Modified: 2025-01-06 09:15:20

Prompt template:
  {file}

  Additional context: Focus on code quality.
  ✓ Valid template
  ✓ Uses {file} placeholder

✓ System prompt is configured correctly
```

**Output (command-based, success):**

```
Testing system prompt configuration...
─────────────────────────────────────────────────

Effective configuration:
  Scope: local
  Type: Command-based

Command:
  Shell: bash
  Timeout: 5 seconds
  Command: git log -1 --format='%s'
  ✓ Executed successfully (42 bytes output)

Prompt template:
  You are a code reviewer.
  Current commit: {command}
  ✓ Valid template
  ✓ Uses {command} placeholder

✓ System prompt is configured correctly
```

**Output (inline prompt, success):**

```
Testing system prompt configuration...
─────────────────────────────────────────────────

Effective configuration:
  Scope: global
  Type: Inline prompt

Prompt:
  You are an expert code reviewer.
  Focus on security and performance.
  ✓ Valid inline prompt (87 characters)

✓ System prompt is configured correctly
```

**Output (file not found):**

```
Testing system prompt configuration...
─────────────────────────────────────────────────

Effective configuration:
  Scope: local
  Type: File-based

File:
  Path: ./MISSING.md
  Resolved: /Users/grant/Projects/myapp/MISSING.md
  ✗ File not found

Prompt: (file content only)

✗ System prompt has errors
  File not found - will fail at runtime
  Fix: Create file or update configuration
```

**Output (command failed):**

```
Testing system prompt configuration...
─────────────────────────────────────────────────

Effective configuration:
  Scope: global
  Type: Command-based

Command:
  Shell: bash
  Timeout: 30 seconds
  Command: nonexistent-command --flag
  ✗ Command failed (exit code 127)
  Error: nonexistent-command: command not found

Prompt template:
  Output: {command}
  ✓ Valid template

✗ System prompt has errors
  Command execution will fail at runtime
```

**Output (no system prompt configured):**

```
Testing system prompt configuration...
─────────────────────────────────────────────────

No system prompt configured (global or local).

Agents will run without system prompts.

Configure: start config role edit global
```

Exit code: 1

**Output (configuration error):**

```
Testing system prompt configuration...
─────────────────────────────────────────────────

Effective configuration:
  Scope: global
  Type: (invalid - no UTD fields)

✗ No content source defined
  At least one field required: file, command, or prompt

✗ System prompt has configuration errors
  Fix configuration: start config role edit global
```

**Verbose output:**

```bash
start config role test --verbose
```

```
Testing system prompt configuration...
─────────────────────────────────────────────────

Loading configuration...
  Global config: ~/.config/start/config.toml
  Local config: ./.start/config.toml

Configuration merge:
  Global [system_prompt]: configured
  Local [system_prompt]: configured (overrides global)
  Effective: local

Local configuration details:
  File field: ./ROLE.md
  File resolution:
    Working directory: /Users/grant/Projects/myapp
    Resolved path: /Users/grant/Projects/myapp/ROLE.md
    ✓ File exists
    Size: 1,234 bytes
    Modified: 2025-01-06 09:15:20

  Prompt field:
    {file}

    Additional context: Focus on code quality.

  Placeholders found: {file}
  ✓ Valid placeholder usage
  ✓ {file} placeholder matches UTD file field

✓ System prompt is configured correctly
```

**Exit codes:**

- 0 - Success (system prompt valid, file exists, command succeeds)
- 1 - No system prompt configured
- 2 - Configuration error (invalid configuration)
- 3 - File not found (config valid but file missing)
- 4 - Command failed (config valid but command execution failed)

**Error handling:**

**Multiple errors:**

```
Testing system prompt configuration...
─────────────────────────────────────────────────

Effective configuration:
  Scope: local
  Type: Combination (file + command)

File:
  ✗ File not found: ./missing.md

Command:
  ✗ Command failed (exit code 1)

Prompt template:
  ⚠ Unknown placeholder {unknown}

✗ System prompt has multiple errors:
  - File not found
  - Command execution failed
  - Invalid placeholder usage
```

Exit code: 2 (configuration errors take precedence)

## Global Flags

These flags work on all `start config role` subcommands where applicable.

**--help**, **-h**
: Show help for the subcommand.

**--verbose**, **-v**
: Verbose output. Shows config file paths and additional details.

**--debug**
: Debug mode. Shows all internal operations.

## Examples

### Show System Prompt (Merged View)

```bash
start config role show merged
```

Show effective system prompt with override information.

### Show Global System Prompt Only

```bash
start config role show global
```

### Show Local System Prompt Only

```bash
start config role show local
```

### Edit Global System Prompt

```bash
start config role edit global
```

### Edit Local System Prompt

```bash
start config role edit local
```

### Test System Prompt

```bash
start config role test
```

Verify system prompt configuration, file availability, and command execution.

### Remove Local System Prompt (Revert to Global)

```bash
start config role remove local
```

### Remove Global System Prompt

```bash
start config role remove global
```

### Interactive Scope Selection

```bash
start config role edit
```

Prompts for scope selection.

## Files

**~/.config/start/config.toml**
: Global configuration file containing `[system_prompt]` section.

**./.start/config.toml**
: Local project configuration file containing project-specific `[system_prompt]` section.

The local section completely replaces the global section (not merged). If local section is missing, uses global section.

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
Line 15: invalid TOML syntax

Fix the configuration file or restore from backup.
```

Exit code: 1

## Notes

### System Prompt Merge Behavior

Per DR-002 and DR-005, the `[system_prompt]` section has special merge behavior:

**Global system prompt:** `~/.config/start/config.toml`
- Personal default role definition
- Used across all projects (unless overridden)

**Local system prompt:** `./.start/config.toml`
- Project-specific role definition
- **Completely replaces** global section (not merged)

**Merge behavior:**
- If local `[system_prompt]` exists: use local only (ignore global)
- If local `[system_prompt]` missing: use global
- This is different from contexts (which are combined)

**Rationale:**
System prompts define the complete role. Merging doesn't make sense - you want either the global role OR the project-specific role, not a combination.

### Optional System Prompt

System prompts are completely optional:

- Section can be omitted entirely (no warning)
- Not all AI agents support system prompts
- If omitted, agents run without role definition
- Some agents have built-in default roles

### Unified Template Design (UTD)

System prompts use UTD pattern for flexible content sourcing:

**File-based:**
```toml
[system_prompt]
file = "~/.config/start/roles/code-reviewer.md"
```

**Command-based:**
```toml
[system_prompt]
command = "git log -1 --format='%s'"
prompt = "You are a code reviewer. Current commit: {command}"
```

**Inline prompt:**
```toml
[system_prompt]
prompt = """
You are an expert code reviewer.
Focus on security and performance.
"""
```

**File with template framing:**
```toml
[system_prompt]
file = "./ROLE.md"
prompt = """
Role Definition:
{file}

Follow these instructions carefully.
"""
```

**Combination (file + command):**
```toml
[system_prompt]
file = "./ROLE.md"
command = "date"
prompt = """
{file}

Current time: {command}
"""
```

See [UTD documentation](../design/unified-template-design.md) for complete details.

### Placeholders

System prompt templates support these placeholders:

- `{file}` - Content from `file` field (empty if not specified)
- `{command}` - Output from `command` field (empty if not specified)

**Example:**
```toml
[system_prompt]
file = "~/.config/start/roles/reviewer.md"
command = "git branch --show-current"
prompt = """
{file}

Current branch: {command}
"""
```

### Shell Configuration

System prompts can override the global shell setting:

```toml
[system_prompt]
command = "git log -1 --format='%s'"
prompt = "Current commit: {command}"
shell = "bash"
command_timeout = 5
```

See [UTD shell configuration](../design/unified-template-design.md#shell-configuration) for supported shells.

### Role Files Location

By convention, role definition files are stored in:

**Global:** `~/.config/start/roles/*.md`
- Personal role definitions
- Managed via `start update` (asset roles)
- Can also be user-created

**Local (per-project):** `./.start/roles/*.md` or `./roles/*.md`
- Project-specific role definitions
- Manually created

**Asset roles:** `~/.config/start/assets/roles/*.md`
- Provided by `start` as defaults
- Updated via `start update`

### Agent Support

Not all AI agents support system prompts:

- Check agent documentation for system prompt support
- Some agents use different terminology (role, instruction, etc.)
- Agent command templates use `{system_prompt}` placeholder
- If agent doesn't support it, placeholder is ignored

## See Also

- start(1) - Launch with context
- start-task(1) - Run predefined tasks
- start-config(1) - Manage configuration files
- start-config-agent(1) - Manage AI agents
- start-config-context(1) - Manage context documents
- start-config-task(1) - Manage task configurations
- start-update(1) - Update asset library
- DR-002 - Configuration file structure and merge behavior
- DR-005 - System prompt handling
- DR-017 - CLI command reorganization
