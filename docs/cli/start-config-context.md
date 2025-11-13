# start config context

## Name

start config context - Manage context document configurations

## Synopsis

```bash
start config context list [scope]
start config context new [scope]
start config context new [scope]
start config context show [name] [scope]
start config context test <name>
start config context edit [name] [scope]
start config context remove [name] [scope]
```

## Description

Manages context document configurations in config files. Context documents provide background information to AI agents. Documents can be configured in global config (`~/.config/start/contexts.toml`) or local project config (`./.start/contexts.toml`).

**Context management operations:**

- **list** - Display all configured contexts with details
- **add** - Add new context interactively
- **show** - Display context configuration structure
- **test** - Test context configuration and file availability
- **edit** - Modify existing context configuration
- **remove** - Delete context from configuration

**Note:** Per DR-017, context documents can be defined in both global and local configs. These commands can manage either scope using the `[scope]` argument. If scope is omitted, the command prompts interactively.

## Context Configuration Structure

Contexts are defined using the **[Unified Template Design (UTD)](../design/unified-template-design.md)** pattern:

```toml
[context.environment]
description = "User environment and tool configuration"
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
required = true
```

**UTD Fields (at least one required):**

**file** (optional)
: Path to context document file. Supports `~` expansion and relative paths.

**command** (optional)
: Shell command to execute for dynamic content. Command string available via `{command}`, output available via `{command_output}`.

**prompt** (optional)
: Template text with placeholders: `{file}`, `{file_contents}`, `{command}`, `{command_output}`.

**Context-Specific Fields:**

**description** (optional)
: Human-readable description of this context.

**required** (optional, default: false)
: Whether this document is required context.
  - `true` - Included by both `start` and `start prompt`
  - `false` - Included by `start`, excluded by `start prompt`

**shell** (optional)
: Override global shell for command execution.

**command_timeout** (optional)
: Override global timeout (seconds) for command execution.

## Subcommands

### start config context list

Display all configured contexts with their details.

**Synopsis:**

```bash
start config context list          # Select scope interactively
start config context list global   # List global contexts only
start config context list local    # List local contexts only
start config context list merged   # Show merged view (global + local)
```

**Behavior:**

Lists all contexts defined in the selected scope(s) with:

- Context name
- Description
- Required status
- Source (file, command, or prompt)
- Scope (global, local, or override)

**Output (merged view):**

```
Configured contexts (merged):
═══════════════════════════════════════════════════════════

Required contexts (3):
  environment [global]
    User environment and tool configuration
    File: ~/reference/ENVIRONMENT.md
    Prompt: Read {file} for environment context.

  index [global]
    Documentation index
    File: ~/reference/INDEX.csv
    Prompt: Read {file} for documentation index.

  agents [local]
    Repository instructions and agent guidance
    File: ./AGENTS.md
    Prompt: Read {file} for repository instructions.

Optional contexts (2):
  project [local]
    Project status and progress
    File: ./PROJECT.md
    Prompt: Read {file}. Respond with summary.

  design [local]
    Design decisions and rationale
    File: ./docs/design-record.md
    Prompt: Read {file} for design decisions.
```

**Output (global only):**

```bash
start config context list global
```

```
Configured contexts (global):
═══════════════════════════════════════════════════════════

Required contexts (2):
  environment
    User environment and tool configuration
    File: ~/reference/ENVIRONMENT.md
    Prompt: Read {file} for environment context.

  index
    Documentation index
    File: ~/reference/INDEX.csv
    Prompt: Read {file} for documentation index.

Optional contexts (1):
  readme
    Project overview
    File: README.md
    Prompt: Project overview from {file}
```

**Output (local only):**

```bash
start config context list local
```

```
Configured contexts (local):
═══════════════════════════════════════════════════════════

Required contexts (1):
  agents
    Repository instructions and agent guidance
    File: ./AGENTS.md
    Prompt: Read {file} for repository instructions.

Optional contexts (2):
  project
    Project status and progress
    File: ./PROJECT.md
    Prompt: Read {file}. Respond with summary.

  design
    Design decisions and rationale
    File: ./docs/design-record.md
    Prompt: Read {file} for design decisions.
```

**No contexts configured:**

```
No contexts configured in global config.

Create contexts: start config context new global
```

**Exit codes:**

- 0 - Success (contexts listed)
- 1 - No config file exists
- 2 - Invalid scope argument

### start config context new

Interactively add a new context to the configuration.

**Synopsis:**

```bash
start config context new          # Select scope interactively
start config context new global   # Add to global config
start config context new local    # Add to local config
```

**Behavior:**

Prompts for context details and adds to the selected config file:

1. **Select scope** (if not provided)
   - global - Add to `~/.config/start/contexts.toml`
   - local - Add to `./.start/contexts.toml`

2. **Context name** (required)
   - Validation: lowercase alphanumeric with hyphens
   - Pattern: `/^[a-z0-9]+(-[a-z0-9]+)*$/`
   - Must be unique within scope
   - Examples: `environment`, `project`, `design-docs`

3. **Description** (optional)
   - Human-readable description
   - Press enter to skip

4. **Content source** (choose one or more)
   - File path (static document)
   - Command (dynamic content)
   - Inline prompt text
   - At least one is required

5. **Prompt template** (required if file or command specified)
   - Template with `{file}` and/or `{command}` placeholders
   - Or inline prompt text

6. **Required context?** (yes/no, default: no)
   - `true` - Included by both `start` and `start prompt`
   - `false` - Included by `start`, excluded by `start prompt`

7. **Advanced options?** (yes/no, default: no)
   - Shell override
   - Command timeout

8. **Backup and save**
   - Backs up existing config to `config.YYYY-MM-DD-HHMMSS.toml`
   - Writes new context to config
   - Shows success message

**Interactive flow (simple file-based context):**

```
Add new context
─────────────────────────────────────────────────

Select scope:
  1) global (all projects)
  2) local (this project only)

Scope [1-2]: 1

Context name: environment
Description (optional): User environment and tool configuration

Content source:
  1) File path
  2) Command
  3) Inline prompt
  4) Combination

Select [1-4]: 1

File path: ~/reference/ENVIRONMENT.md
✓ File exists

Prompt template: Read {file} for environment context.
✓ Valid template (uses {file} placeholder)

Required context? [y/N]: y

Advanced options? [y/N]: n

Backing up config to config.2025-01-06-091523.toml...
✓ Backup created

Saving context 'environment' to ~/.config/start/contexts.toml...
✓ Context added successfully

Use 'start config context list global' to see all contexts.
Use 'start config context test environment' to verify.
```

**Interactive flow (dynamic command-based context):**

```
Add new context
─────────────────────────────────────────────────

Select scope:
  1) global (all projects)
  2) local (this project only)

Scope [1-2]: 2

Context name: git-status
Description (optional): Current git working tree status

Content source:
  1) File path
  2) Command
  3) Inline prompt
  4) Combination

Select [1-4]: 2

Command: git status --short
✓ Valid command

Prompt template: Working tree status:\n{command}
✓ Valid template (uses {command} placeholder)

Required context? [y/N]: n

Advanced options? [y/N]: y

Shell override (or enter for default): bash
Command timeout in seconds (or enter for default): 5

Backing up config to config.2025-01-06-091645.toml...
✓ Backup created

Saving context 'git-status' to ./.start/contexts.toml...
✓ Context added successfully

Use 'start config context list local' to see all contexts.
Use 'start config context test git-status' to verify.
```

**Interactive flow (inline prompt):**

```
Add new context
─────────────────────────────────────────────────

Select scope:
  1) global (all projects)
  2) local (this project only)

Scope [1-2]: 2

Context name: project-note
Description (optional): Project-specific guidance

Content source:
  1) File path
  2) Command
  3) Inline prompt
  4) Combination

Select [1-4]: 3

Prompt text: Important: This project uses Go 1.21 and follows standard Go conventions.
✓ Valid prompt

Required context? [y/N]: y

Advanced options? [y/N]: n

Backing up config to config.2025-01-06-091712.toml...
✓ Backup created

Saving context 'project-note' to ./.start/contexts.toml...
✓ Context added successfully

Use 'start config context list local' to see all contexts.
```

**Resulting config (file-based):**

```toml
[context.environment]
description = "User environment and tool configuration"
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
required = true
```

**Resulting config (command-based):**

```toml
[context.git-status]
description = "Current git working tree status"
command = "git status --short"
prompt = "Working tree status:\n{command}"
required = false
shell = "bash"
command_timeout = 5
```

**Resulting config (inline prompt):**

```toml
[context.project-note]
description = "Project-specific guidance"
prompt = "Important: This project uses Go 1.21 and follows standard Go conventions."
required = true
```

**Exit codes:**

- 0 - Success (context added)
- 1 - Validation error (invalid name, duplicate context, invalid configuration)
- 2 - Scope error (invalid scope, local config directory doesn't exist)
- 3 - File system error (cannot write config, backup failed)

**Error handling:**

**Invalid context name:**

```
Context name: My-Context
✗ Invalid context name. Use lowercase alphanumeric with hyphens.
  Examples: environment, project, design-docs

Context name: my-context
✓ Valid name
```

**Duplicate context:**

```
Context name: environment
✗ Context 'environment' already exists in global config.

Use 'start config context edit environment global' to modify existing context.
```

Exit code: 1

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
  Context will be skipped at runtime if file is not found.

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

**Backup failed:**

```
Backing up config to config.2025-01-06-091523.toml...
✗ Failed to backup config: permission denied

Existing config preserved at: ~/.config/start/contexts.toml
Context not added.
```

Exit code: 3

### start config context show

Display current context configuration.

**Synopsis:**

```bash
start config context show                 # Select context and scope interactively
start config context show <name>          # Select scope for named context
start config context show <name> global   # Show global context only
start config context show <name> local    # Show local context only
```

**Behavior:**

Displays context configuration from the selected scope with:

- Scope (global or local)
- Context name
- Description (if configured)
- Source type (file, command, inline, or combination)
- File path (if configured)
- Command (if configured)
- Prompt template (if configured)
- Required flag (true/false)
- Shell and timeout overrides (if configured)

**Output (global context - file):**

```
Context configuration: environment (global)
═══════════════════════════════════════════════════════════

Description: User environment and tool configuration
Required: true

Source: File-based

File:
  Path: ~/reference/ENVIRONMENT.md
  Resolved: /Users/grant/reference/ENVIRONMENT.md

Prompt template:
  Read {file} for environment context.
```

**Output (local context - command-based):**

```bash
start config context show git-status local
```

```
Context configuration: git-status (local)
═══════════════════════════════════════════════════════════

Description: Current git working tree status
Required: false

Source: Command-based

Command:
  Shell: bash
  Timeout: 5 seconds
  Command: git status --short

Prompt template:
  Working tree status:
  {command_output}
```

**Output (inline context):**

```
Context configuration: note (local)
═══════════════════════════════════════════════════════════

Required: true

Source: Inline prompt

Prompt:
  Important: This project uses Go 1.21
```

**No context configured:**

```
No context 'nonexistent' found in global config.

Configure: start config context new global
```

**Interactive selection:**

```bash
start config context show
```

```
Show context configuration
─────────────────────────────────────────────────

Select context:
  1) environment
  2) index
  3) agents
  4) project

Select [1-4]: 1

Select scope:
  1) global
  2) local

Scope [1-2]: 1

(displays context configuration)
```

**Exit codes:**

- 0 - Success (context shown)
- 1 - No context configured
- 2 - Invalid scope argument
- 3 - Context not found

**Error handling:**

**Context not found:**

```
Error: Context 'nonexistent' not found in configuration.

Use 'start config context list' to see available contexts.
```

Exit code: 3

### start config context test

Test context configuration and file availability.

**Synopsis:**

```bash
start config context test <name>
```

**Behavior:**

Validates context configuration without using it. Performs three checks:

1. **File availability** (if `file` field present)
   - Checks if file exists at specified path
   - Resolves `~` and relative paths
   - Reports: found (with resolved path) or not found

2. **Command execution** (if `command` field present)
   - Executes command in configured shell
   - Reports: success (with output size) or failure (with error)
   - Does NOT display command output (security)

3. **Configuration validation**
   - At least one UTD field present (file, command, or prompt)
   - Prompt template uses valid placeholders (`{file}`, `{command}`)
   - Unknown placeholders detected (likely typos)
   - Shell and timeout settings valid

**Does NOT include the context in any prompt** - only validates and shows what would be available.

**Output (file-based, success):**

```
Testing context: environment
─────────────────────────────────────────────────

Configuration:
  Scope: global
  Description: User environment and tool configuration
  Required: yes
  Type: File-based

File:
  Path: ~/reference/ENVIRONMENT.md
  Resolved: /Users/grant/reference/ENVIRONMENT.md
  ✓ File exists (2.3 KB)

Prompt template:
  Read {file} for environment context.
  ✓ Valid template
  ✓ Uses {file} placeholder

✓ Context 'environment' is configured correctly
```

**Output (command-based, success):**

```
Testing context: git-status
─────────────────────────────────────────────────

Configuration:
  Scope: local
  Description: Current git working tree status
  Required: no
  Type: Command-based

Command:
  Shell: bash
  Timeout: 5 seconds
  Command: git status --short
  ✓ Executed successfully (42 bytes output)

Prompt template:
  Working tree status:
  {command}
  ✓ Valid template
  ✓ Uses {command} placeholder

✓ Context 'git-status' is configured correctly
```

**Output (inline prompt, success):**

```
Testing context: project-note
─────────────────────────────────────────────────

Configuration:
  Scope: local
  Description: Project-specific guidance
  Required: yes
  Type: Inline prompt

Prompt:
  Important: This project uses Go 1.21...
  ✓ Valid inline prompt (87 characters)

✓ Context 'project-note' is configured correctly
```

**Output (file not found):**

```
Testing context: missing-doc
─────────────────────────────────────────────────

Configuration:
  Scope: local
  Description: Missing documentation
  Required: no
  Type: File-based

File:
  Path: ./MISSING.md
  Resolved: /Users/grant/Projects/myapp/MISSING.md
  ✗ File not found

Prompt template:
  Read {file} for context.
  ✓ Valid template
  ✓ Uses {file} placeholder

⚠ Context 'missing-doc' has warnings
  File will generate warning and be skipped at runtime
```

**Output (command failed):**

```
Testing context: broken-cmd
─────────────────────────────────────────────────

Configuration:
  Scope: local
  Description: Broken command example
  Required: yes
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
  ✓ Uses {command} placeholder

✗ Context 'broken-cmd' has errors
  Command execution will fail at runtime
```

**Output (configuration error):**

```
Testing context: invalid
─────────────────────────────────────────────────

Configuration:
  Scope: global
  Description: Invalid context configuration
  Required: yes
  Type: (invalid - no UTD fields)

✗ No content source defined
  At least one field required: file, command, or prompt

✗ Context 'invalid' has configuration errors
  Fix configuration: start config context edit invalid global
```

**Verbose output:**

```bash
start config context test environment --verbose
```

```
Testing context: environment
─────────────────────────────────────────────────

Loading configuration...
  Config file: ~/.config/start/contexts.toml
  Context section: [context.environment]

Configuration details:
  Name: environment
  Scope: global
  Description: User environment and tool configuration
  Required: true

UTD fields:
  file: ~/reference/ENVIRONMENT.md
  prompt: Read {file} for environment context.

File resolution:
  Original path: ~/reference/ENVIRONMENT.md
  Home expansion: /Users/grant/reference/ENVIRONMENT.md
  ✓ File exists
  Size: 2,345 bytes
  Modified: 2025-01-05 14:23:10

Prompt template analysis:
  Template: Read {file} for environment context.
  Placeholders found: {file}
  ✓ Valid placeholder usage
  ✓ {file} placeholder matches UTD file field

✓ Context 'environment' is configured correctly
```

**Exit codes:**

- 0 - Success (context valid, file exists or command succeeds)
- 1 - Configuration error (invalid configuration)
- 2 - Context not found in config
- 3 - File not found (config valid but file missing)
- 4 - Command failed (config valid but command execution failed)

**Error handling:**

**Context not in config:**

```
Error: Context 'nonexistent' not found in configuration.

Use 'start config context list' to see available contexts.
Use 'start assets add' to install from catalog or 'start config context new' to create custom.
```

Exit code: 2

**Multiple errors:**

```
Testing context: broken
─────────────────────────────────────────────────

Configuration:
  ✗ No UTD fields present (no file, command, or prompt)
  ⚠ Unknown placeholder {unknown} in prompt template

File:
  ✗ File not found: ./missing.md

✗ Context 'broken' has multiple errors:
  - Invalid configuration (no content source)
  - File not found
  - Invalid placeholder usage
```

Exit code: 1 (configuration errors take precedence)

### start config context edit

Edit context configuration interactively.

**Synopsis:**

```bash
start config context edit                  # Select context and scope
start config context edit <name>           # Select scope for named context
start config context edit <name> global    # Edit in global config
start config context edit <name> local     # Edit in local config
```

**Behavior:**

**Without context name (interactive selection):**

Shows list of configured contexts for selection:

```bash
start config context edit
```

Output:

```
Edit context
─────────────────────────────────────────────────

Select context to edit:

Global contexts:
  1) environment (required)
  2) index (required)
  3) readme (optional)

Local contexts:
  4) agents (required)
  5) project (optional)
  6) design (optional)

Select [1-6] (or 'q' to quit): 1

(continues to interactive edit flow for 'environment' in global config)
```

**With context name only:**

If context exists in only one scope, edits that config. If exists in both, prompts for scope:

```bash
start config context edit environment
```

Context exists in global only:

```
Editing context 'environment' in global config...
(continues to interactive edit flow)
```

Context exists in both:

```
Context 'environment' exists in multiple scopes.

Select scope to edit:
  1) global - ~/reference/ENVIRONMENT.md
  2) local - ./ENVIRONMENT.md

Select [1-2]: 1
(continues to interactive edit flow for global)
```

**With context name and scope:**

Interactive prompts to edit specific context. Shows current values as defaults - press enter to keep current value.

1. **Description** - Current value shown in brackets
2. **Content source changes** - Modify file, command, or prompt
   - Keep existing content source
   - Add new content source
   - Remove content source (if multiple exist)
3. **Prompt template** - Current value shown in brackets
4. **Required flag** - Current value shown
5. **Advanced options** - Shell, timeout
6. **Backup and save** - Backs up to `config.YYYY-MM-DD-HHMMSS.toml`

**Interactive flow:**

```
Edit context: environment (global)
─────────────────────────────────────────────────

Current configuration:
  Description: User environment and tool configuration
  File: ~/reference/ENVIRONMENT.md
  Prompt: Read {file} for environment context.
  Required: yes
  Shell: (default)
  Timeout: (default)

Press enter to keep current value, or type new value:

Description [User environment and tool configuration]:
File path [~/reference/ENVIRONMENT.md]:
Prompt template [Read {file} for environment context.]:
Required [yes]:

Advanced options? [y/N]: n

Backing up config to config.2025-01-06-092312.toml...
✓ Backup created

Saving changes to ~/.config/start/contexts.toml...
✓ Context 'environment' updated successfully

Use 'start config context list global' to see changes.
Use 'start config context test environment' to validate.
```

**Interactive flow (adding command to file-based context):**

```
Edit context: environment (global)
─────────────────────────────────────────────────

Current configuration:
  Description: User environment and tool configuration
  File: ~/reference/ENVIRONMENT.md
  Prompt: Read {file} for environment context.
  Required: yes

Press enter to keep current value, or type new value:

Description [User environment and tool configuration]:
File path [~/reference/ENVIRONMENT.md]:

Add command for dynamic content? [y/N]: y

Command: echo "Current time: $(date)"
✓ Valid command

Prompt template [Read {file} for environment context.]: {file}\n\nCurrent context: {command}
✓ Valid template (uses {file} and {command})

Required [yes]:

Advanced options? [y/N]: n

Backing up config to config.2025-01-06-092415.toml...
✓ Backup created

Saving changes to ~/.config/start/contexts.toml...
✓ Context 'environment' updated successfully
```

**Exit codes:**

- 0 - Success (context edited)
- 1 - Validation error (invalid configuration)
- 2 - Context not found or scope error
- 3 - File system error (cannot write config, backup failed)

**Error handling:**

**Context not found:**

```
Error: Context 'nonexistent' not found in configuration.

Use 'start config context list' to see available contexts.
Use 'start assets add' to install from catalog or 'start config context new' to create custom.
```

Exit code: 2

**Invalid prompt template:**

```
Prompt template [Read {file} for context.]: Invalid {unknown} placeholder
⚠ Warning: Unknown placeholder {unknown}
  Valid placeholders: {file}, {command}

Continue anyway? [y/N]: n

Prompt template [Read {file} for context.]: Read {file} for environment.
✓ Valid template
```

**No changes made:**

```
No changes detected.

Context 'environment' not modified.
```

Exit code: 0 (no backup created, no write)

### start config context remove

Remove context from configuration.

**Synopsis:**

```bash
start config context remove                  # Select context and scope
start config context remove <name>           # Select scope for named context
start config context remove <name> global    # Remove from global config
start config context remove <name> local     # Remove from local config
```

**Behavior:**

**Without context name:**
Shows list of configured contexts for selection:

```
Remove context
─────────────────────────────────────────────────

Select context to remove:

Global contexts:
  1) environment (required)
  2) index (required)
  3) readme (optional)

Local contexts:
  4) agents (required)
  5) project (optional)
  6) design (optional)

Select [1-6] (or 'q' to quit): 3

Remove context 'readme' from global config? [y/N]: y

Backing up config to config.2025-01-06-093012.toml...
✓ Backup created

Removing context 'readme' from ~/.config/start/contexts.toml...
✓ Context 'readme' removed successfully

Use 'start config context list global' to see remaining contexts.
```

**With context name only:**
If context exists in only one scope, removes from that config. If exists in both, prompts for scope:

```bash
start config context remove readme
```

Context exists in global only:

```
Remove context 'readme' from global config? [y/N]: y

Backing up config to config.2025-01-06-093045.toml...
✓ Backup created

Removing context 'readme' from ~/.config/start/contexts.toml...
✓ Context 'readme' removed successfully

Use 'start config context list global' to see remaining contexts.
```

Context exists in both:

```
Context 'environment' exists in multiple scopes.

Select scope to remove from:
  1) global - ~/reference/ENVIRONMENT.md
  2) local - ./ENVIRONMENT.md
  3) both

Select [1-3]: 1

Remove context 'environment' from global config? [y/N]: y
(continues with removal from global)
```

**Removing required context (warning):**

```bash
start config context remove environment global
```

Output:

```
⚠ Warning: 'environment' is marked as required context.
  Removing it may affect agent behavior.

Remove context 'environment' from global config? [y/N]: y

Backing up config to config.2025-01-06-093123.toml...
✓ Backup created

Removing context 'environment' from ~/.config/start/contexts.toml...
✓ Context 'environment' removed successfully

Use 'start config context list global' to see remaining contexts.
```

**Declining confirmation:**

```
Remove context 'readme' from global config? [y/N]: n

Context 'readme' not removed.
```

Exit code: 0

**Exit codes:**

- 0 - Success (context removed, or user declined)
- 1 - No contexts configured
- 2 - Context not found or scope error
- 3 - File system error (cannot write config, backup failed)

**Error handling:**

**Context not found:**

```
Error: Context 'nonexistent' not found in configuration.

Use 'start config context list' to see available contexts.
```

Exit code: 2

**No contexts configured:**

```
No contexts configured in global config.

Use 'start config context new global' to create a context.
```

Exit code: 1

**Backup failed:**

```
Remove context 'readme' from global config? [y/N]: y

Backing up config to config.2025-01-06-093156.toml...
✗ Failed to backup config: permission denied

Existing config preserved at: ~/.config/start/contexts.toml
Context not removed.
```

Exit code: 3

## Global Flags

These flags work on all `start config context` subcommands where applicable.

**--help**, **-h**
: Show help for the subcommand.

**--verbose**
: Verbose output. Shows config file paths and additional details.

**--debug**
: Debug mode. Shows all internal operations.

**--version**, **-v**
: Show version information.

## Examples

### List All Contexts (Merged View)

```bash
start config context list merged
```

Show all contexts from both global and local configs.

### List Global Contexts Only

```bash
start config context list global
```

### Create Context in Global Config

```bash
start config context new global
```

### Create Context in Local Config

```bash
start config context new local
```

### Test Context

```bash
start config context test environment
```

Verify context configuration and file availability.

### Edit Context

```bash
start config context edit environment global
```

### Remove Context

```bash
start config context remove readme global
```

### Interactive Context Selection

```bash
start config context edit
```

Shows list of all contexts to choose from.

## Files

**~/.config/start/contexts.toml**
: Global configuration file containing context definitions.

**./.start/contexts.toml**
: Local project configuration file containing project-specific contexts.

Both files can contain `[context.<name>]` sections. Contexts are combined when both configs exist (global + local).

## Error Handling

### No Configuration File

```
Error: No configuration found at ~/.config/start/

Run 'start init' to create initial configuration.
```

Exit code: 1

### Invalid TOML Syntax

```
Error: Configuration file has invalid syntax.

File: ~/.config/start/contexts.toml
Line 42: invalid TOML syntax

Fix the configuration file or restore from backup.
```

Exit code: 1

## Notes

### Context Merge Behavior

Per DR-002, contexts from global and local configs are combined:

**Global contexts:** `~/.config/start/contexts.toml`
- User-wide context documents
- Managed by `start config context` commands
- Shared across all projects

**Local contexts:** `./.start/contexts.toml`
- Project-specific context documents
- Managed by `start config context` commands with `local` scope
- Can override global contexts (same name)

**Merge behavior:**
- Global + local contexts are combined
- Order: Global contexts first, then local contexts
- If name conflict: Local overrides global (intentional, not an error)

**Document order:**
Documents appear in the prompt in definition order:
1. Global contexts (in TOML order)
2. Local contexts (in TOML order)

Rearrange config definitions to change prompt order.

### Required vs Optional Contexts

**Required contexts** (`required = true`):
- Included by both `start` and `start prompt`
- Auto-included in all tasks
- Essential background information

**Optional contexts** (`required = false`):
- Included by `start` only
- Excluded from `start prompt` (for focused queries)
- Nice-to-have information

See DR-012 for full rationale.

### Unified Template Design (UTD)

Contexts use the UTD pattern for flexible content sourcing:

**File-based:**
```toml
[context.environment]
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
```

**Command-based:**
```toml
[context.git-status]
command = "git status --short"
prompt = "Working tree status:\n{command}"
```

**Inline prompt:**
```toml
[context.note]
prompt = "Important: This project uses Go 1.21"
```

**Combined:**
```toml
[context.project-state]
file = "./PROJECT.md"
command = "git log -5 --oneline"
prompt = "{file}\n\nRecent commits:\n{command}"
```

See [UTD documentation](../design/unified-template-design.md) for complete details.

### Placeholders

Context prompt templates support these placeholders:

- `{file}` - File path from `file` field (absolute, ~ expanded)
- `{file_contents}` - Content from `file` field (empty if file missing)
- `{command}` - Command string from `command` field
- `{command_output}` - Output from `command` execution (empty if command fails)

**Example:**
```toml
[context.environment]
file = "~/reference/ENVIRONMENT.md"
command = "date"
prompt = """
{file_contents}

Current timestamp: {command_output}
"""
```

### Shell Configuration

Contexts can override the global shell setting:

```toml
[context.git-status]
command = "git status --short"
prompt = "Working tree status:\n{command_output}"
shell = "bash"
command_timeout = 5
```

See [UTD shell configuration](../design/unified-template-design.md#shell-configuration) for supported shells.

## See Also

- start(1) - Launch with context
- start-prompt(1) - Launch with custom prompt
- start-task(1) - Run predefined tasks
- start-config(1) - Manage configuration files
- start-config-agent(1) - Manage AI agents
- start-config-task(1) - Manage task configurations
- start-config-role(1) - Manage system prompts
- DR-002 - Configuration file structure and merge behavior
- DR-003 - Named documents for context
- DR-008 - Context file detection and handling
- DR-012 - Context document required field and order
- DR-017 - CLI command reorganization
