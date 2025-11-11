# start config role

## Name

start config role - Manage role configuration

## Synopsis

```bash
start config role list
start config role add [path]
start config role new [scope]
start config role show [scope]
start config role edit [scope]
start config role remove [scope]
start config role default [name]
start config role test
```

## Description

Manages role configuration in config files. Roles define the AI agent's persona and behavior (system prompts). Roles are passed to agents via the `{role}` and `{role_file}` placeholders in agent commands.

**Role management operations:**

- **list** - Display all configured roles
- **add** - Add role from catalog (interactive) or by path
- **show** - Display current role configurations
- **edit** - Modify a role (create if doesn't exist)
- **remove** - Remove a role configuration
- **default** - Set or show default role
- **test** - Verify role configuration and file availability

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
: Shell command to execute for dynamic content. Command string available via `{command}`, output available via `{command_output}`.

**prompt** (optional)
: Template text with placeholders: `{file}`, `{file_contents}`, `{command}`, `{command_output}`.

**Additional Fields:**

**shell** (optional)
: Override global shell for command execution.

**command_timeout** (optional)
: Override global timeout (seconds) for command execution.

**Merge behavior:**

Global and local `[roles.<name>]` sections are **combined**. If a role with the same name exists in both configs, local overrides global for that role name.

## Subcommands

### start config role list

Display all roles (configured and catalog).

**Synopsis:**

```bash
start config role list
```

**Behavior:**

Lists roles from three sources:

1. **Global config** (`~/.config/start/roles.toml`) - Personal roles
2. **Local config** (`./.start/roles.toml`) - Project-specific roles
3. **Asset catalog** (`~/.config/start/assets/roles/`) - Available catalog roles

**Output:**

```
Configured roles:
═══════════════════════════════════════════════════════════

Global (2):
  code-reviewer
    Expert code reviewer focused on security
    File: ~/.config/start/assets/roles/general/code-reviewer.md

  default
    Balanced helpful assistant
    File: ~/.config/start/assets/roles/general/default.md

Local (1):
  project-specific
    Project-specific role definition
    File: ./ROLE.md

Available catalog roles (4):
  general/code-reviewer
    Expert code reviewer focused on security

  general/default
    Balanced helpful assistant

  languages/go-expert
    Go programming language expert

  specialized/rubber-duck
    Socratic method questioning assistant
```

**Exit codes:**

- 0 - Success (roles listed)
- 1 - No config file exists

### start config role add

Add a role from the official GitHub asset catalog.

**Synopsis:**

```bash
start config role add                   # Interactive catalog browser
start config role add general/code-reviewer  # Direct install from catalog
start config role add --local           # Add to local config
```

**Arguments:**

**[path]** (optional)
: The full catalog path of the role to install directly (e.g., `general/code-reviewer`). If omitted, an interactive browser is shown.

**Flags:**

**--local**
: Add to local config (`./.start/config.toml`) instead of global.

**Behavior:**

**Interactive mode** (no path specified):

1. Fetch catalog from GitHub
2. Show categories and roles
3. User selects role
4. Download and cache role file
5. Add role configuration to config file

**Direct mode** (path specified):

1. Check if role exists in catalog
2. Download and cache role file
3. Add role configuration to config file

**Interactive flow:**

```
Fetching catalog from GitHub...
✓ Found 8 roles across 3 categories

Select category:
  1. general (4 roles)
  2. languages (2 roles)
  3. specialized (2 roles)

Choice [1-3]: 1

general roles:
  1. code-reviewer - Expert code reviewer focused on security
  2. default - Balanced helpful assistant
  3. pair-programmer - Collaborative programming assistant
  4. explainer - Teaching mode that simplifies concepts

Choice [1-4]: 1

Selected: code-reviewer
Description: Expert code reviewer focused on security
Tags: review, quality, security, strict

Download and add to config? [Y/n] y

Downloading...
✓ Cached to ~/.config/start/assets/roles/general/
✓ Added to global config as 'code-reviewer'

Try it: start --role code-reviewer
    or: start task <taskname>  # Task-specific role overrides work too
```

**Direct install:**

```bash
start config role add general/code-reviewer
```

Output:

```
Downloading general/code-reviewer...
✓ Cached to ~/.config/start/assets/roles/general/
✓ Added to global config as 'code-reviewer'
```

**Add to local config:**

```bash
start config role add general/code-reviewer --local
```

Output:

```
Downloading general/code-reviewer...
✓ Cached to ~/.config/start/assets/roles/general/
✓ Added to local config as 'code-reviewer'
```

**Resulting configuration:**

```toml
# In config.toml (or tasks.toml, agents.toml, contexts.toml as appropriate)
[roles.code-reviewer]
description = "Expert code reviewer focused on security"
file = "~/.config/start/assets/roles/general/code-reviewer.md"
```

**Error handling:**

**Network unavailable:**

```
Error: Cannot fetch catalog from GitHub

  Network error: dial tcp: no route to host

To resolve:
- Check internet connection
- Add custom role: start config role edit
```

**Role not found in catalog:**

```
Error: Role 'nonexistent/role' not found in catalog

Available roles:
  - general/code-reviewer
  - general/default
  - languages/go-expert
  - specialized/rubber-duck

Try: start config role add  # Browse interactively
```

**Exit codes:**

- 0 - Success (role added)
- 1 - Network error (catalog unavailable)
- 2 - Role not found
- 3 - User cancelled

### start config role new

Interactively create a new custom role configuration.

**Synopsis:**

```bash
start config role new [scope]
```

**Behavior:**

This command launches an interactive wizard to help you create a new role from scratch. It will prompt you for the content source (file, command, or inline prompt) and other configuration details. This is for creating your own custom roles, as opposed to adding existing ones from the asset catalog.

(See `start config agent new` for a detailed example of the interactive wizard flow).


### start config role show

Display current role configuration.

**Synopsis:**

```bash
start config role show          # Select scope interactively
start config role show global   # Show global role only
start config role show local    # Show local role only
start config role show merged   # Show effective role (with override info)
```

**Behavior:**

Displays role configuration from the selected scope(s) with:

- Scope (global, local, or merged)
- Source type (file, command, inline, or combination)
- File path (if configured)
- Command (if configured)
- Prompt template (if configured)
- Shell and timeout overrides (if configured)

**Output (merged view):**

```
Role configuration (merged):
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
Role configuration (global):
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
Role configuration (local):
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
Role configuration (global):
═══════════════════════════════════════════════════════════

Type: Inline prompt

Prompt:
  You are an expert code reviewer.
  Focus on security and performance.
```

**Output (command-based):**

```
Role configuration (local):
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

**No role configured:**

```
No role configured in global config.

Configure: start config role edit global
```

**Exit codes:**

- 0 - Success (role shown)
- 1 - No role configured
- 2 - Invalid scope argument

### start config role edit

Edit or create role configuration interactively.

**Synopsis:**

```bash
start config role edit          # Select scope interactively
start config role edit global   # Edit global role
start config role edit local    # Edit local role
```

**Behavior:**

Prompts for role configuration and updates the selected config file:

1. **Select scope** (if not provided)
   - global - Edit `~/.config/start/roles.toml`
   - local - Edit `./.start/roles.toml`

2. **Content source** (choose one or more)
   - File path (static role document)
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
   - Writes role to config
   - Shows success message

**Interactive flow (create from file):**

```
Edit role
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

Saving role to ~/.config/start/roles.toml...
✓ Role configured successfully

Use 'start config role show global' to verify.
Use 'start config role test' to validate.
```

**Interactive flow (file with template):**

```
Edit role
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

Saving role to ./.start/roles.toml...
✓ Role updated successfully

Use 'start config role show local' to verify.
```

**Interactive flow (inline prompt):**

```
Edit role
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

Saving role to ~/.config/start/roles.toml...
✓ Role configured successfully
```

**Interactive flow (combination - file + command):**

```
Edit role
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

Saving role to ./.start/roles.toml...
✓ Role configured successfully
```

**Resulting config (simple file):**

```toml
[roles.code-reviewer]
file = "~/.config/start/roles/code-reviewer.md"
```

**Resulting config (file with template):**

```toml
[roles.project-default]
file = "./ROLE.md"
prompt = """
{file}

Additional context: Focus on code quality.
"""
```

**Resulting config (inline prompt):**

```toml
[roles.inline-reviewer]
prompt = """
You are an expert code reviewer.
Focus on security and performance.
"""
```

**Resulting config (combination):**

```toml
[roles.git-aware-reviewer]
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

- 0 - Success (role configured)
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
  Role will fail at runtime if file is not found.

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

Existing config preserved at: ~/.config/start/roles.toml
Role not configured.
```

Exit code: 3

### start config role remove

Remove role configuration.

**Synopsis:**

```bash
start config role remove          # Select scope interactively
start config role remove global   # Remove from global config
start config role remove local    # Remove from local config
```

**Behavior:**

Removes a `[roles.<name>]` section from the selected config file. You'll be prompted to select which role to remove.

**Interactive flow:**

```bash
start config role remove
```

Output:

```
Remove role
─────────────────────────────────────────────────

Select scope:
  1) global (all projects)
  2) local (this project only)

Scope [1-2]: 2

Current configuration (local):
  File: ./ROLE.md
  Prompt template: {file}\n\nFocus on code quality.

Remove role from local config? [y/N]: y

Backing up config to config.2025-01-06-112012.toml...
✓ Backup created

Removing [roles.project-default] from ./.start/roles.toml...
✓ Role removed successfully

Global role will now be used (if configured).

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

Remove role from global config? [y/N]: y

Backing up config to config.2025-01-06-112045.toml...
✓ Backup created

Removing [roles.code-reviewer] from ~/.config/start/roles.toml...
✓ Role removed successfully
⚠ No role configured

Agents will run without roles.

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

⚠ Note: Removing local role will revert to global configuration.

Remove role from local config? [y/N]: y

Backing up config to config.2025-01-06-112123.toml...
✓ Backup created

Removing [roles.project-default] from ./.start/roles.toml...
✓ Role removed successfully
✓ Now using global role

Global configuration:
  File: ~/.config/start/roles/default.md

Use 'start config role show merged' to verify.
```

**Declining confirmation:**

```
Remove role from global config? [y/N]: n

Role not removed.
```

Exit code: 0

**Exit codes:**

- 0 - Success (role removed, or user declined)
- 1 - No role configured
- 2 - Invalid scope
- 3 - File system error (cannot write config, backup failed)

**Error handling:**

**No role configured:**

```
Error: No role configured in global config.

Configure: start config role edit global
```

Exit code: 1

**Backup failed:**

```
Remove role from global config? [y/N]: y

Backing up config to config.2025-01-06-112156.toml...
✗ Failed to backup config: permission denied

Existing config preserved at: ~/.config/start/roles.toml
Role not removed.
```

Exit code: 3

### start config role default

Set or show the default role.

**Synopsis:**

```bash
start config role default          # Show current default
start config role default <name>   # Set default role
```

**Behavior:**

Without a role name, shows the current default role. With a role name, sets it as the default in the `[settings]` section of global config.

**Output (show current default):**

```bash
start config role default
```

Output:

```
Default role: code-reviewer
  Expert code reviewer focusing on security
  Source: ~/.config/start/roles.toml

Use 'start config role default <name>' to change.
```

**No default set:**

If `default_role` is not configured in `[settings]`:

```
No default role configured.

First role in config will be used: code-reviewer

Use 'start config role default <name>' to set explicitly.
```

Exit code: 0

**Setting default (no previous default):**

```bash
start config role default security-auditor
```

Output:

```
Backing up config to config.2025-01-06-113045.toml...
✓ Backup created

Setting default role to 'security-auditor'...
✓ Default role set to 'security-auditor'

Use 'start' to launch with default role.
Use 'start config role default' to confirm.
```

Exit code: 0

**Updating existing default:**

```bash
start config role default go-expert
```

Output:

```
Current default: code-reviewer

Backing up config to config.2025-01-06-113112.toml...
✓ Backup created

Setting default role to 'go-expert'...
✓ Default role changed: code-reviewer → go-expert

Use 'start' to launch with new default.
```

Exit code: 0

**Exit codes:**

- 0 - Success (default shown or set)
- 1 - No roles configured
- 2 - Role not found
- 3 - File system error (cannot write config, backup failed)

**Error handling:**

**Role not found:**

```
Error: Role 'nonexistent' not found in configuration.

Available roles:
  - code-reviewer
  - security-auditor
  - go-expert

Use 'start config role list' for details.
```

Exit code: 2

**No roles configured:**

```
Error: No roles configured.

Use 'start config role add' to add a role.
Use 'start init' to set up roles automatically.
```

Exit code: 1

### start config role test

Test role configuration and file availability.

**Synopsis:**

```bash
start config role test
```

**Behavior:**

Validates effective role configuration (merged global + local). Performs checks:

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

**Does NOT pass role to any agent** - only validates configuration.

**Output (file-based, success):**

```
Testing role configuration...
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

✓ Role is configured correctly
```

**Output (file with template, success):**

```
Testing role configuration...
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

✓ Role is configured correctly
```

**Output (command-based, success):**

```
Testing role configuration...
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

✓ Role is configured correctly
```

**Output (inline prompt, success):**

```
Testing role configuration...
─────────────────────────────────────────────────

Effective configuration:
  Scope: global
  Type: Inline prompt

Prompt:
  You are an expert code reviewer.
  Focus on security and performance.
  ✓ Valid inline prompt (87 characters)

✓ Role is configured correctly
```

**Output (file not found):**

```
Testing role configuration...
─────────────────────────────────────────────────

Effective configuration:
  Scope: local
  Type: File-based

File:
  Path: ./MISSING.md
  Resolved: /Users/grant/Projects/myapp/MISSING.md
  ✗ File not found

Prompt: (file content only)

✗ Role has errors
  File not found - will fail at runtime
  Fix: Create file or update configuration
```

**Output (command failed):**

```
Testing role configuration...
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

✗ Role has errors
  Command execution will fail at runtime
```

**Output (no role configured):**

```
Testing role configuration...
─────────────────────────────────────────────────

No role configured (global or local).

Agents will run without roles.

Configure: start config role edit global
```

Exit code: 1

**Output (configuration error):**

```
Testing role configuration...
─────────────────────────────────────────────────

Effective configuration:
  Scope: global
  Type: (invalid - no UTD fields)

✗ No content source defined
  At least one field required: file, command, or prompt

✗ Role has configuration errors
  Fix configuration: start config role edit global
```

**Verbose output:**

```bash
start config role test --verbose
```

```
Testing role configuration...
─────────────────────────────────────────────────

Loading configuration...
  Global config: ~/.config/start/roles.toml
  Local config: ./.start/roles.toml

Configuration merge:
  Global [roles]: 3 roles configured
  Local [roles]: 1 role configured
  Combined: 4 total roles (1 override)

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

✓ Role is configured correctly
```

**Exit codes:**

- 0 - Success (role valid, file exists, command succeeds)
- 1 - No role configured
- 2 - Configuration error (invalid configuration)
- 3 - File not found (config valid but file missing)
- 4 - Command failed (config valid but command execution failed)

**Error handling:**

**Multiple errors:**

```
Testing role configuration...
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

✗ Role has multiple errors:
  - File not found
  - Command execution failed
  - Invalid placeholder usage
```

Exit code: 2 (configuration errors take precedence)

## Global Flags

These flags work on all `start config role` subcommands where applicable.

**--help**, **-h**
: Show help for the subcommand.

**--verbose**
: Verbose output. Shows config file paths and additional details.

**--debug**
: Debug mode. Shows all internal operations.

**--version**, **-v**
: Show version information.

## Examples

### Show Role (Merged View)

```bash
start config role show merged
```

Show effective role with override information.

### Show Global Role Only

```bash
start config role show global
```

### Show Local Role Only

```bash
start config role show local
```

### Edit Global Role

```bash
start config role edit global
```

### Edit Local Role

```bash
start config role edit local
```

### Test Role

```bash
start config role test
```

Verify role configuration, file availability, and command execution.

### Remove Local Role (Revert to Global)

```bash
start config role remove local
```

### Remove Global Role

```bash
start config role remove global
```

### Interactive Scope Selection

```bash
start config role edit
```

Prompts for scope selection.

## Files

**~/.config/start/roles.toml**
: Global roles configuration file containing `[roles.<name>]` sections.

**./.start/roles.toml**
: Local project roles configuration file containing project-specific `[roles.<name>]` sections.

Global and local roles are combined. If a role with the same name exists in both configs, local overrides global for that role name.

## Error Handling

### No Configuration File

```
Error: No roles configuration found at ~/.config/start/roles.toml

Run 'start init' to create initial configuration.
```

Exit code: 1

### Invalid TOML Syntax

```
Error: Configuration file has invalid syntax.

File: ~/.config/start/roles.toml
Line 15: invalid TOML syntax

Fix the configuration file or restore from backup.
```

Exit code: 1

## Notes

### Role Merge Behavior

Per DR-002 and DR-005, `[roles.<name>]` sections have combine-and-override merge behavior:

**Global roles:** `~/.config/start/roles.toml`
- Personal default role definitions
- Used across all projects

**Local roles:** `./.start/roles.toml`
- Project-specific role definitions
- Added to global roles

**Merge behavior:**
- Global and local roles are **combined**
- If a role with the same name exists in both: local overrides global for that role
- All other roles from both configs remain available
- This allows projects to override specific roles while keeping others

**Rationale:**
Projects often need custom roles for specific workflows while still having access to global roles. The combine-and-override approach provides maximum flexibility.

### Optional Roles

Roles are completely optional:

- Section can be omitted entirely (no warning)
- Not all AI agents support roles (system prompts)
- If omitted, agents run without role definition
- Some agents have built-in default roles

### Unified Template Design (UTD)

Roles use UTD pattern for flexible content sourcing:

**File-based:**
```toml
[roles.code-reviewer]
file = "~/.config/start/roles/code-reviewer.md"
```

**Command-based:**
```toml
[roles.dynamic-reviewer]
command = "git log -1 --format='%s'"
prompt = "You are a code reviewer. Current commit: {command}"
```

**Inline prompt:**
```toml
[roles.security-reviewer]
prompt = """
You are an expert code reviewer.
Focus on security and performance.
"""
```

**File with template framing:**
```toml
[roles.project-role]
file = "./ROLE.md"
prompt = """
Role Definition:
{file}

Follow these instructions carefully.
"""
```

**Combination (file + command):**
```toml
[roles.time-aware-role]
file = "./ROLE.md"
command = "date"
prompt = """
{file}

Current time: {command}
"""
```

See [UTD documentation](../design/unified-template-design.md) for complete details.

### Placeholders

Role templates support these placeholders:

- `{file}` - File path from `file` field (absolute, ~ expanded)
- `{file_contents}` - Content from `file` field (empty if file missing)
- `{command}` - Command string from `command` field
- `{command_output}` - Output from `command` execution (empty if command fails)

**Example:**
```toml
[roles.branch-aware-reviewer]
file = "~/.config/start/roles/reviewer.md"
command = "git branch --show-current"
prompt = """
{file_contents}

Current branch: {command_output}
"""
```

### Shell Configuration

Roles can override the global shell setting:

```toml
[roles.git-reviewer]
command = "git log -1 --format='%s'"
prompt = "Current commit: {command_output}"
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

Not all AI agents support roles (system prompts):

- Check agent documentation for role/system prompt support
- Some agents use different terminology (role, instruction, etc.)
- Agent command templates use `{role}` and `{role_file}` placeholders
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
- DR-005 - Role configuration & selection
- DR-017 - CLI command reorganization
