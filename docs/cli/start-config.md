# start config

## Name

start config - Manage configuration files

## Synopsis

```bash
start config show
start config edit [scope]
start config path
start config validate
```

## Description

Manages `start` configuration files (global and local). Provides tools for viewing merged configuration, editing config files, locating config files, and validating configuration syntax.

**Configuration management operations:**

- **show** - Display merged configuration with sources
- **edit** - Open config.toml (settings) file in editor
- **path** - Show config directory paths and files
- **validate** - Check configuration syntax and semantics across all config files

**Configuration hierarchy** (per DR-031):

- **Global:** `~/.config/start/` - User-wide configuration files
  - `config.toml` - Settings
  - `agents.toml` - Agent configurations
  - `roles.toml` - Role definitions
  - `contexts.toml` - Context document references
  - `tasks.toml` - Task definitions
- **Local:** `./.start/` - Project-specific configuration files (same structure)
- **Merge behavior:** Local overrides global for matching items, collections are combined

## Subcommands

### start config show

Display merged configuration showing both global and local values with source indicators.

**Synopsis:**

```bash
start config show
```

**Behavior:**

Shows the effective configuration after merging global and local configs from all configuration files (config.toml, agents.toml, roles.toml, contexts.toml, tasks.toml). Displays:

- Which values come from global vs local config
- All sections: settings, agents, contexts, roles, tasks
- Missing optional files are noted but not shown as errors

**Output (both global and local exist):**

```
Configuration
─────────────────────────────────────────────────

Global: ~/.config/start/
Local:  ./.start/

Settings
────────
  default_agent = "claude"              (global)
  verbosity = "normal"                  (global)
  required_only = false                 (local, overrides global)

Agents (3)
──────────
  claude                                (global)
    Command: claude --model {model} --append-system-prompt '{role}' '{prompt}'
    Default model: sonnet
    Models: haiku, sonnet, opus

  gemini                                (global)
    Command: gemini --model {model} '{prompt}'
    Default model: flash
    Models: flash, pro-exp

  aichat                                (global)
    Command: aichat --model {model} '{prompt}'
    Default model: gpt4-mini
    Models: gpt4-mini, gpt4, claude

Contexts (5)
────────────
  environment                           (global)
    Path: ~/reference/ENVIRONMENT.md
    Required: yes
    Prompt: Read for environment context

  project                               (local)
    Path: PROJECT.md
    Required: yes
    Prompt: Read the PROJECT.md document

  agents                                (local)
    Path: AGENTS.md
    Required: no
    Prompt: Read for agent instructions

  readme                                (global)
    Path: README.md
    Required: no
    Prompt: Project overview and setup

  changelog                             (global)
    Path: CHANGELOG.md
    Required: no
    Prompt: Recent changes

Roles (2)
─────────
  coder                                 (global)
    Path: ~/.config/start/roles/coder.md

  reviewer                              (global)
    Path: ~/.config/start/roles/reviewer.md
```

**Output (global only):**

```
Configuration
─────────────────────────────────────────────────

Global: ~/.config/start/
Local:  (none)

Settings
────────
  default_agent = "claude"
  verbosity = "normal"

Agents (2)
──────────
  [... agents ...]

Contexts (3)
────────────
  [... contexts ...]

Roles (1)
─────────
  [... roles ...]
```

**Output (local only, no global):**

```
Configuration
─────────────────────────────────────────────────

Global: (none)
Local:  ./.start/

Settings
────────
  required_only = false

Contexts (2)
────────────
  project                               (local)
    Path: PROJECT.md
    Required: yes
    Prompt: Read the PROJECT.md document

  agents                                (local)
    Path: AGENTS.md
    Required: no
    Prompt: Read for agent instructions

Agents: (none - agents must be defined in global config)
Roles: (none)
```

**Verbose output:**

```bash
start config show --verbose
```

Adds:

- Full file paths for roles
- Agent command templates
- Description and URL fields for agents
- Internal merge order and precedence details

**Exit codes:**

- 0 - Success (config displayed)
- 1 - No config files found (neither global nor local)

**Error handling:**

**No config files:**

```
Error: No configuration found.

Checked locations:
  Global: ~/.config/start/ (not found)
  Local:  ./.start/ (not found)

Run 'start init' to create initial configuration.
```

Exit code: 1

**Invalid TOML syntax:**

```
Error: Configuration file has invalid syntax.

File: ~/.config/start/config.toml
Line 23: expected '=' after key, found ']'

Fix the configuration file or restore from backup.
Use 'start config validate' for detailed validation.
```

Exit code: 1

### start config edit

Open settings configuration file (config.toml) in editor. For editing other config files, use specialized commands: `start config task edit`, `start config agent edit`, `start config role edit`, `start config context edit`.

**Synopsis:**

```bash
start config edit [scope]
```

**Arguments:**

**[scope]** (optional)
: Which config.toml to edit. If omitted, edit will detect and ask.

- `global` - Edit `~/.config/start/config.toml` (settings)
- `local` - Edit `./.start/config.toml` (settings)

**Behavior:**

Opens the configuration file in your preferred editor:

1. **Editor detection:**

   - Checks `$VISUAL` environment variable
   - Falls back to `$EDITOR` environment variable
   - Falls back to `vi`
   - Shows info message if neither `$VISUAL` nor `$EDITOR` set

2. **Config selection:**

   - With `global`: Edit `~/.config/start/config.toml`
   - With `local`: Edit `./.start/config.toml`
   - No scope, both exist: Ask which to edit
   - No scope, only one exists: Edit that one
   - No scope, neither exists: Ask which to create

3. **After editing:**
   - Validates config syntax and semantics
   - Shows warnings for issues (soft warnings, file already saved)
   - Does not prevent saving or re-open editor

**Interactive selection (both configs exist):**

```bash
start config edit
```

Output:

```
Edit configuration
─────────────────────────────────────────────────

Both global and local configs exist:
  1) Global: ~/.config/start/config.toml
  2) Local:  ./.start/config.toml

Select [1-2] (or 'q' to quit): 1

Opening ~/.config/start/config.toml in vi...
Set $EDITOR to use your preferred editor.
```

**Edit specific config:**

```bash
start config edit global
```

Output:

```
Opening ~/.config/start/config.toml in nvim...
```

(User's editor opens, they make changes and save)

```
Validating configuration...
✓ Configuration is valid

Changes saved to ~/.config/start/config.toml
```

**Validation warnings after edit:**

```
Opening ~/.config/start/config.toml in nvim...
```

(User's editor opens, they introduce issues and save)

```
Validating configuration...

⚠ Warnings found:

Context 'missing':
  File not found: ~/reference/MISSING.md
  This file will be skipped at runtime.

Agent 'broken':
  Command template missing required {prompt} placeholder

✓ Changes saved, but configuration has warnings
  Run 'start config validate' for full validation details.
```

**Config file doesn't exist:**

```bash
start config edit local
```

Output:

```
Local config does not exist: ./.start/config.toml

Create new local config? [y/N]: y

Creating ./.start directory...
✓ Directory created

Creating ./.start/config.toml...
✓ File created

Opening ./.start/config.toml in nvim...
```

**Editor not configured:**

```bash
start config edit global
```

(When `$VISUAL` and `$EDITOR` not set)

Output:

```
Opening ~/.config/start/config.toml in vi...

ℹ Using 'vi' (default editor)
  Set $EDITOR environment variable to use your preferred editor.

  Examples:
    export EDITOR=nvim
    export EDITOR="code --wait"
    export EDITOR="subl --wait"

  Note: GUI editors require '--wait' flag to work properly.
```

**Exit codes:**

- 0 - Success (file edited)
- 1 - User cancelled or quit
- 2 - Config directory not writable
- 3 - Editor command failed

**Error handling:**

**Directory not writable:**

```bash
start config edit local
```

Output:

```
Error: Cannot create ./.start directory
  Permission denied

Check directory permissions and try again.
```

Exit code: 2

**Editor command failed:**

```bash
start config edit global
```

(When `$EDITOR="badeditor"`)

Output:

```
Error: Editor command failed: badeditor
  Command not found in PATH

Set $EDITOR to a valid editor command.
Available editors: vi, vim, nvim, nano, emacs
```

Exit code: 3

### start config path

Show paths to configuration directories and files.

**Synopsis:**

```bash
start config path
```

**Behavior:**

Displays the paths to global and local config directories, listing all configuration files and indicating whether each exists.

**Output (both exist):**

```
Configuration locations:

Global: ~/.config/start/ ✓
  config.toml ✓
  agents.toml ✓
  roles.toml ✓
  contexts.toml ✓
  tasks.toml ✓

Local:  ./.start/ ✓
  config.toml ✓
  agents.toml ✓
  roles.toml ✓
  contexts.toml ✓
  tasks.toml ✓

Use 'start config edit [scope]' to edit settings.
Use specialized commands for other files (e.g., 'start config task edit').
```

**Output (only global exists):**

```
Configuration locations:

Global: ~/.config/start/ ✓
  config.toml ✓
  agents.toml ✓
  roles.toml ✓
  contexts.toml ✓
  tasks.toml ✓

Local:  ./.start/ (not found)

Use 'start config edit global' to edit global settings.
Use 'start init local' to create local config.
```

**Output (neither exists):**

```
Configuration locations:

Global: ~/.config/start/ (not found)
Local:  ./.start/ (not found)

Run 'start init' to create initial configuration.
```

**Verbose output:**

```bash
start config path --verbose
```

Output:

```
Configuration locations:

Global directory: /Users/grant/.config/start/
  config.toml: 1.2 KB (2025-01-04 14:30:15)
  agents.toml: 3.1 KB (2025-01-04 14:30:15)
  roles.toml: 847 bytes (2025-01-04 14:30:15)
  contexts.toml: 456 bytes (2025-01-04 14:30:15)
  tasks.toml: 2.8 KB (2025-01-04 14:30:15)
  Backups: 12 files

Local directory: /Users/grant/Projects/myproject/.start/
  config.toml: 234 bytes (2025-01-04 15:12:42)
  agents.toml: (not found)
  roles.toml: (not found)
  contexts.toml: 198 bytes (2025-01-04 15:12:42)
  tasks.toml: (not found)
  Backups: 2 files
```

**Exit codes:**

- 0 - Success (paths shown)

### start config validate

Validate configuration syntax and semantics across all config files without launching agent.

**Synopsis:**

```bash
start config validate
```

**Behavior:**

Performs comprehensive validation of all configuration files (config.toml, agents.toml, roles.toml, contexts.toml, tasks.toml) in both global and local locations:

1. **TOML syntax validation**

   - Valid TOML structure in all files
   - Proper section headers
   - Correct data types

2. **Semantic validation**

   - Required fields present
   - Valid field values
   - Agent command templates have `{prompt}` placeholder
   - Unknown placeholders detected
   - Context file paths exist
   - Role file paths exist
   - Task definitions valid

3. **Merge validation**
   - No conflicting settings between global and local
   - Combined context names unique
   - Local overrides properly structured

**Output (all valid):**

```
Validating configuration...
─────────────────────────────────────────────────

Global: ~/.config/start/
  ✓ config.toml - Settings valid
  ✓ agents.toml - 3 agents configured
    ✓ claude
    ✓ gemini
    ✓ aichat
  ✓ roles.toml - 2 roles configured
  ✓ contexts.toml - 3 contexts configured
  ✓ tasks.toml - 4 tasks configured

Local: ./.start/
  ✓ config.toml - Settings valid
  ✓ agents.toml - (not present)
  ✓ roles.toml - (not present)
  ✓ contexts.toml - 2 contexts configured
  ✓ tasks.toml - (not present)

Merged configuration:
  ✓ No conflicts
  ✓ 5 total contexts (3 global + 2 local)
  ✓ 3 agents available
  ✓ 2 roles available
  ✓ 4 tasks available
  ✓ Default agent set: claude

✓ Configuration is valid
```

**Output (warnings):**

```
Validating configuration...
─────────────────────────────────────────────────

Global: ~/.config/start/config.toml
  ✓ TOML syntax valid
  ✓ Settings section valid
  ✓ Agents section valid (2 agents)
    ✓ claude
    ⚠ gemini
      Binary not found in PATH: gemini
      Command will fail unless 'gemini' is installed
  ✓ Contexts section valid (3 contexts)
    ⚠ readme
      File not found: README.md
      This context will be skipped at runtime
  ✓ Roles section valid (1 role)

Local: ./.start/config.toml
  ✓ TOML syntax valid
  ✓ Settings section valid
  ✓ Contexts section valid (1 context)

Merged configuration:
  ✓ 4 total contexts (3 global + 1 local)
  ✓ 2 agents available
  ✓ Default agent set: claude

⚠ Configuration has 2 warnings (see above)
  Configuration will work, but some features may be unavailable.
```

**Output (errors):**

```
Validating configuration...
─────────────────────────────────────────────────

Global: ~/.config/start/config.toml
  ✓ TOML syntax valid
  ✓ Settings section valid
  ✗ Agents section has errors
    ✗ broken-agent
      Command template missing required {prompt} placeholder
      Command: broken-agent --model {model}
    ⚠ test-agent
      Unknown placeholder {mdoel} in command template
      Valid placeholders: {model}, {role}, {role_file}, {prompt}, {date}
  ✓ Contexts section valid (2 contexts)

Local: ./.start/config.toml
  ✗ TOML syntax error at line 15
    expected '=' after key, found ']'

✗ Configuration has errors
  Fix errors before using 'start'.

  Use 'start config edit global' to fix global config.
  Use 'start config edit local' to fix local config.
```

**Output (duplicate agents with warnings):**

```
Validating configuration...
─────────────────────────────────────────────────

Global: ~/.config/start/config.toml
  ✓ TOML syntax valid
  ✓ Agents section valid (2 agents)
    ✓ claude
    ✓ gemini

Local: ./.start/config.toml
  ✓ TOML syntax valid
  ✓ Agents section valid (1 agent)
    ⚠ claude
      Also defined in global config
      Local will override global for this agent

Merged configuration:
  ✓ 2 agents available (claude from local, gemini from global)
  ✓ Default agent set: claude

⚠ Configuration has 1 warning (see above)
  Local agent 'claude' overrides global definition.
```

**Verbose output:**

```bash
start config validate --verbose
```

Adds detailed checks:

- Full file paths checked
- Each agent command template parsed
- Each placeholder validated
- Each context file path checked
- Model alias validation
- Default model resolution

**Exit codes:**

- 0 - Success (config valid, or warnings only)
- 1 - Validation errors (config will not work)
- 2 - TOML syntax errors
- 3 - No config files found

**Error handling:**

**No config files:**

```
Error: No configuration found.

Checked locations:
  Global: ~/.config/start/ (not found)
  Local:  ./.start/ (not found)

Run 'start init' to create initial configuration.
```

Exit code: 3

**TOML syntax error:**

```
Error: TOML syntax error in global config

File: ~/.config/start/config.toml
Line 23: expected '=' after key, found ']'

  21 | [agents.claude]
  22 | description = "Claude AI"
  23 | command = ["broken syntax"]
     |          ^
  24 |
  25 | [agents.gemini]

Fix the syntax error and run validation again.
Use 'start config edit global' to edit the file.
```

Exit code: 2

## Global Flags

These flags work on all `start config` subcommands where applicable.

**--help**, **-h**
: Show help for the subcommand.

**--verbose**
: Verbose output. Shows additional details and full paths.

**--debug**
: Debug mode. Shows all internal operations and validation steps.

**--version**, **-v**
: Show version information.

## Examples

### View Merged Configuration

```bash
start config show
```

See effective configuration after global + local merge.

### Edit Global Config

```bash
start config edit global
```

Open global config in your editor (uses `$EDITOR` or `vi`).

### Edit Local Config

```bash
start config edit local
```

Open local project config in your editor.

### Find Config Files

```bash
start config path
```

Show paths to both global and local config files.

### Validate Before Launch

```bash
start config validate
```

Check configuration for errors before running `start`.

### Validate with Details

```bash
start config validate --verbose
```

Show detailed validation results including all checks performed.

## Files

**~/.config/start/**
: Global configuration directory containing:
  - `config.toml` - Settings only
  - `agents.toml` - Agent configurations
  - `roles.toml` - Role definitions
  - `contexts.toml` - Context document references
  - `tasks.toml` - Task definitions

**./.start/**
: Local project configuration directory (same file structure as global)
  - Local values override global for matching items
  - Collections are combined (e.g., contexts from both are merged)

**~/.config/start/\*.YYYY-MM-DD-HHMMSS.toml**
: Backup files created before modifications. Examples:
  - `config.2025-01-04-103045.toml`
  - `agents.2025-01-04-103045.toml`
  - `roles.2025-01-04-103045.toml`

## Editor Configuration

The `start config edit` command uses the following editor detection sequence:

1. `$VISUAL` - Full-screen editor (highest priority)
2. `$EDITOR` - Standard editor environment variable
3. `vi` - Universal fallback (always available on Unix)

**Setting your editor:**

```bash
# For terminal editors
export EDITOR=nvim
export EDITOR=vim
export EDITOR=nano
export EDITOR=emacs

# For GUI editors (requires --wait flag)
export EDITOR="code --wait"       # VS Code
export EDITOR="subl --wait"       # Sublime Text
export EDITOR="mate --wait"       # TextMate
```

Add to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.) to persist.

**GUI editor requirements:**

GUI editors must use the `--wait` flag to block until the file is closed. Without `--wait`:

- The editor command returns immediately
- `start` cannot validate the file after editing
- Validation is skipped

This is the user's responsibility. `start config edit` will launch the editor as configured.

## Validation Behavior

**After `start config edit`:**

- **Soft warnings** - Issues are reported but file is already saved
- **Non-blocking** - User is not forced to re-edit
- **Informational** - Helps identify problems for future fixes

**With `start config validate`:**

- **Hard errors** - Exit code 1 for issues that prevent `start` from working
- **Warnings** - Exit code 0 for issues that don't prevent operation
- **Comprehensive** - Checks syntax, semantics, and merge behavior

## Notes

### Configuration Merge Behavior

When both global and local configs exist:

1. **Settings:** Local values override global values
2. **Contexts:** Combined (global + local), names must be unique
3. **Agents:** Only from global (local cannot define agents per DR-004)
4. **Roles:** Only from global (local cannot define roles)

**Example:**

Global config:

```toml
[settings]
verbosity = "normal"
default_agent = "claude"

[contexts.environment]
path = "~/reference/ENVIRONMENT.md"
required = true
```

Local config:

```toml
[settings]
verbosity = "verbose"  # Overrides global

[contexts.project]
path = "PROJECT.md"
required = true
```

Merged result:

```toml
[settings]
verbosity = "verbose"      # From local (override)
default_agent = "claude"   # From global (not overridden)

# Contexts: environment (global) + project (local)
```

### Agent Configuration Scope

Per DR-004 (updated 2025-01-05), agents can be defined in **both global and local** configs:

- Global: `~/.config/start/config.toml` - Personal agent configurations
- Local: `./.start/config.toml` - Team/project agent configurations (can be committed)
- Merge behavior: Local overrides global for same agent name
- `start config validate` warns if local agent overrides global
- Rationale: Enables team standardization via committed configs while maintaining personal preferences

### Backup Files

Commands that modify configs create timestamped backups for each file:

- Format: `<filename>.YYYY-MM-DD-HHMMSS.toml`
- Location: Same directory as original config files
- Created before any modification
- Not automatically cleaned up (manual deletion safe)
- Examples:
  - `config.2025-01-04-103045.toml`
  - `agents.2025-01-04-103045.toml`
  - `tasks.2025-01-04-103045.toml`

**Backup creation triggers:**

- `start config edit` (after validation warnings)
- `start config task add|edit|remove` (before write to tasks.toml)
- `start config agent add|edit|remove` (before write to agents.toml)
- `start config role edit|remove` (before write to roles.toml)
- `start config context add|edit|remove` (before write to contexts.toml)
- `start init` (when modifying existing config files)

### Context File Validation

Validation checks if context files exist but:

- Missing files are **warnings**, not errors
- Files are checked relative to current directory (for local contexts)
- Tilde (`~`) expansion is supported
- Missing optional contexts (`required = false`) are silently skipped at runtime

## See Also

- start(1) - Launch with context
- start-agent(1) - Manage agents
- start-init(1) - Initialize configuration
- DR-004 - Agent configuration scope
- DR-001 - Configuration format (TOML)
