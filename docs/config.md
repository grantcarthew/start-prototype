# Configuration Reference

Complete reference for `start` configuration files.

## Overview

`start` uses [TOML](https://toml.io/) configuration files with a two-tier hierarchy:

- **Global:** `~/.config/start/` - User-wide configurations (multi-file)
  - `config.toml` - Settings only
  - `roles.toml` - Role definitions
  - `tasks.toml` - Task definitions
  - `agents.toml` - Agent configurations
  - `contexts.toml` - Context configurations
- **Local:** `./.start/` - Project-specific configurations (multi-file)
  - `config.toml` - Settings overrides
  - `roles.toml` - Project roles
  - `tasks.toml` - Project tasks
  - `agents.toml` - Project agents
  - `contexts.toml` - Project contexts

**Merge behavior:**
- Settings: Merged per-field, local overrides global for same field
- Agents: Combined (global + local), local overrides global for same name
- Contexts: Combined (global + local), local overrides global for same name
- Roles: Combined (global + local), local overrides global for same name
## Asset Resolution & Lazy-Loading

`start` treats configurations for agents, roles, tasks, and contexts as "assets". When you ask for an asset (e.g., by running `start task <name>`, `start --role <name>`, or `start config context add <name>`), the CLI resolves it using the following order of precedence:

1.  **Local Config:** (`./.start/`) - Highest priority. Allows project-specific overrides.
2.  **Global Config:** (`~/.config/start/`) - Your personal, user-wide configurations.
3.  **Asset Cache:** (`~/.config/start/assets/`) - Assets you have previously downloaded.
4.  **GitHub Catalog:** If an asset is not found in any of the above locations, the CLI will search the official GitHub asset catalog.

If the asset is found in the GitHub catalog, it will be "lazy-loaded": downloaded, saved to the asset cache, and added to your global configuration for immediate and future offline use. This makes the entire library of official assets available on-demand without requiring you to install them all upfront.

## Complete Example

### Global Config

**config.toml** (`~/.config/start/config.toml`)

```toml
# Global settings
[settings]
default_agent = "claude"
default_role = "code-reviewer"
log_level = "normal"
shell = "bash"
command_timeout = 30
asset_download = true
asset_repo = "grantcarthew/start"
asset_path = "~/.config/start/assets"
```

**roles.toml** (`~/.config/start/roles.toml`)

```toml
[roles.code-reviewer]
description = "Expert code reviewer"
file = "~/.config/start/assets/roles/general/code-reviewer.md"

[roles.coder]
file = "~/.config/start/roles/coder.md"

[roles.reviewer]
file = "~/.config/start/roles/reviewer.md"
```

**agents.toml** (`~/.config/start/agents.toml`)

```toml
[agents.claude]
description = "Anthropic's Claude AI assistant via Claude Code CLI"
url = "https://docs.claude.com/claude-code"
models_url = "https://docs.anthropic.com/en/docs/about-claude/models"
command = "claude --model {model} --append-system-prompt '{role}' '{prompt}'"
default_model = "sonnet"

  [agents.claude.models]
  haiku = "claude-3-5-haiku-20241022"
  sonnet = "claude-3-7-sonnet-20250219"
  opus = "claude-opus-4-20250514"

[agents.gemini]
description = "Google's Gemini AI via CLI"
url = "https://github.com/example/gemini-cli"
models_url = "https://ai.google.dev/models/gemini"
command = "GEMINI_SYSTEM_MD='{role_file}' gemini --model {model} '{prompt}'"
default_model = "flash"

  [agents.gemini.models]
  flash = "gemini-2.0-flash-exp"
  pro-exp = "gemini-2.0-pro-exp"

[agents.aichat]
description = "All-in-one multi-provider AI chat tool"
url = "https://github.com/sigoden/aichat"
command = "aichat --model {model} '{prompt}'"
default_model = "gpt4-mini"

  [agents.aichat.models]
  gpt4-mini = "gpt-4o-mini"
  gpt4 = "gpt-4o"
  claude = "claude-3-5-sonnet-20241022"
```

**contexts.toml** (`~/.config/start/contexts.toml`)

```toml
[context.environment]
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
required = true

[context.index]
file = "~/reference/INDEX.csv"
prompt = "Read {file} for documentation index."
required = true

[context.readme]
file = "README.md"
prompt = "Project overview from {file}"
required = false
```

### Local Config

**config.toml** (`./.start/config.toml`)

```toml
# Project-specific settings (overrides global)
[settings]
log_level = "debug"
```

**roles.toml** (`./.start/roles.toml`)

```toml
# Project-specific role (overrides global code-reviewer)
[roles.code-reviewer]
file = "./ROLE.md"
description = "Project-specific code reviewer"
```

**contexts.toml** (`./.start/contexts.toml`)

```toml
[context.agents]
file = "./AGENTS.md"
prompt = "Read {file} for repository instructions and agent guidance."
required = true

[context.project]
file = "./PROJECT.md"
prompt = "Read {file}. Respond with the project title and shortest possible summary of work required."
required = false

[context.design]
file = "./docs/design-record.md"
prompt = "Read {file} for design decisions."
required = false
```

### Merged Result

When both configs exist, the effective configuration combines them:

**Settings:**
- `default_agent = "claude"` (from global)
- `log_level = "debug"` (from local, overrides global "normal")

**Agents:**
- claude, gemini, aichat (from global; local can override or add agents)

**Roles:**
- code-reviewer: `./ROLE.md` (from local, overrides global)
- go-expert: `~/.config/start/roles/go-base.md` (from global)

**Context documents (in definition order):**
1. environment - `~/reference/ENVIRONMENT.md` (global, required)
2. index - `~/reference/INDEX.csv` (global, required)
3. readme - `README.md` (global, optional)
4. agents - `./AGENTS.md` (local, required)
5. project - `./PROJECT.md` (local, optional)
6. design - `./docs/design-record.md` (local, optional)

## File Locations

### Global Config

```
~/.config/start/
├── config.toml      # Settings only
├── roles.toml       # Role definitions
├── tasks.toml       # Task definitions
├── agents.toml      # Agent configurations
├── contexts.toml    # Context configurations
└── assets/          # Cached catalog assets
```

**Purpose:**
- User-wide defaults and configurations
- Shared across all projects
- Personal preferences and common workflows

**Created by:** `start init`

### Local Config

```
./.start/
├── config.toml      # Settings overrides
├── roles.toml       # Project-specific roles
├── tasks.toml       # Project-specific tasks
├── agents.toml      # Project-specific agents (optional)
└── contexts.toml    # Project-specific contexts
```

**Purpose:**
- Project-specific overrides and customizations
- Team-shareable configurations (can be committed)
- Project-specific workflows

**Created by:** Manual creation or `start init` in project directory

## Configuration Sections

### [settings]

Global behavior settings. Local overrides global.

**Fields:**

**default_agent** (string, optional)
: Which agent to use when `--agent` flag not provided. Resolved using asset resolution algorithm (local config → global config → cache → GitHub catalog).

```toml
[settings]
default_agent = "claude"
```

**default_role** (string, optional)
: Which role to use when `--role` flag not provided and task doesn't specify a role. Resolved using asset resolution algorithm (local config → global config → cache → GitHub catalog). If omitted, first role in config (TOML order) is used.

```toml
[settings]
default_role = "code-reviewer"
```

**log_level** (string, optional)
: Default output level. Overridden by command-line flags (`--quiet`, `--verbose`, `--debug`).

Valid values: `"quiet"`, `"normal"`, `"verbose"`, `"debug"`

```toml
[settings]
log_level = "normal"
```

**shell** (string, optional)
: Default shell for executing commands in `command` fields. Overridden by section-specific `shell` field.

See [Unified Template Design](./unified-template-design.md#shell-configuration) for supported shells.

```toml
[settings]
shell = "bash"
```

**command_timeout** (integer, optional)
: Default timeout in seconds for command execution. Overridden by section-specific `command_timeout` field.

Default: 30 seconds

```toml
[settings]
command_timeout = 30
```

**asset_download** (boolean, optional)
: Enable automatic download of assets from GitHub catalog when not found locally. Can be overridden by `--asset-download` flag. Default: `true`

```toml
[settings]
asset_download = true
```

**asset_repo** (string, optional)
: GitHub repository for asset catalog. Default: `"grantcarthew/start"`

```toml
[settings]
asset_repo = "grantcarthew/start"
```

**asset_path** (string, optional)
: Local directory for caching downloaded assets. Default: `"~/.config/start/assets"`

```toml
[settings]
asset_path = "~/.config/start/assets"
```

**Validation:**

All fields use soft validation with fallback defaults:

- **default_agent** misconfigured or agent not found → **Warning**, fall back to first agent in config (TOML order)
- **default_role** misconfigured or role not found → **Warning**, fall back to first role in config (TOML order)
- **log_level** invalid value → **Warning**, fall back to `"normal"`
- **shell** not found in PATH → **Warning**, fall back to auto-detected shell (`bash` or `sh`)
- **command_timeout** invalid → **Warning**, fall back to 30 seconds
- **asset_download** invalid → **Warning**, fall back to `true`
- **asset_repo** invalid format → **Warning**, fall back to `"grantcarthew/start"`
- **asset_path** invalid or inaccessible → **Warning**, fall back to `"~/.config/start/assets"`
- Missing fields → Silent, use defaults

**Example:**

```toml
[settings]
default_agent = "claude"
default_role = "code-reviewer"
log_level = "normal"
shell = "bash"
command_timeout = 30
asset_download = true
asset_repo = "grantcarthew/start"
asset_path = "~/.config/start/assets"
```

**Merge behavior:**

Local settings override global settings. If a setting is omitted in local config, the global value is used.

---

### [agents.\<name\>]

AI agent tool configurations. Can be defined in both global and local configs.

Each agent section defines how to invoke an AI tool. Agent names should match the actual tool binary name (e.g., `claude`, `gemini`, `aichat`).

**Fields:**

**command** (string, required)
: Command template to execute the agent. Should contain `{prompt}` placeholder to receive the composed prompt. Supports additional placeholders: `{model}`, `{role}`, `{role_file}`, `{date}`.

```toml
[agents.claude]
command = "claude --model {model} --append-system-prompt '{role}' '{prompt}'"
```

**description** (string, optional)
: Human-readable description of the agent. Displayed in `start config agent list`.

```toml
[agents.claude]
description = "Anthropic's Claude AI assistant via Claude Code CLI"
```

**url** (string, optional)
: Documentation or homepage URL for the agent tool.

```toml
[agents.claude]
url = "https://docs.claude.com/claude-code"
```

**models_url** (string, optional)
: URL to model documentation. Helps users understand available models and capabilities.

```toml
[agents.claude]
models_url = "https://docs.anthropic.com/en/docs/about-claude/models"
```

**default_model** (string, optional)
: Model name to use when `--model` flag not provided. Must be a key in the `[agents.<name>.models]` table. If omitted, first model in `models` table is used.

```toml
[agents.claude]
default_model = "sonnet"
```

**Validation:**

- **command** missing `{prompt}` placeholder → **Warning**: "Command doesn't contain {prompt} - composed prompt won't be passed to agent"
- If command uses `{model}` placeholder:
  - **[agents.\<name\>.models]** section MUST exist with ≥1 model → **Error** if missing
- If **default_model** defined but not in models table → **Warning**, fall back to first model (TOML order)
- Unknown placeholders in command → **Warning**: `"Unknown placeholder {mdoel} (did you mean {model}?)"`
- Same agent name in global and local → **Info**: Local overrides global

**Scope:**

Agents can be defined in both **global and local** configs (per DR-004 update):

- **Global agents:** `~/.config/start/agents.toml` - Personal configurations
- **Local agents:** `./.start/agents.toml` - Team/project configurations (can be committed)
- **Merge behavior:** Local overrides global for same agent name
- **Use case:** Teams can commit `.start/` with standard configs; individuals maintain personal preferences

**Example agent (full):**

```toml
[agents.claude]
description = "Anthropic's Claude AI assistant via Claude Code CLI"
url = "https://docs.claude.com/claude-code"
models_url = "https://docs.anthropic.com/en/docs/about-claude/models"
command = "claude --model {model} --append-system-prompt '{role}' '{prompt}'"
default_model = "sonnet"

  [agents.claude.models]
  haiku = "claude-3-5-haiku-20241022"
  sonnet = "claude-3-7-sonnet-20250219"
  opus = "claude-opus-4-20250514"
```

**Example agent (minimal):**

```toml
[agents.simple-agent]
command = "simple-agent '{prompt}'"
```

#### [agents.\<name\>.models]

Model name mappings for the agent. User-defined names map to full model identifiers.

**Structure:**

```toml
[agents.<name>.models]
<name> = "<full-model-identifier>"
```

**Names:**
- User-defined (not hardcoded)
- lowercase, alphanumeric, hyphens
- Each agent defines its own names
- Common patterns: `haiku`, `sonnet`, `opus`, `flash`, `fast`, `best`, etc.

**Examples:**

```toml
[agents.claude.models]
haiku = "claude-3-5-haiku-20241022"
sonnet = "claude-3-7-sonnet-20250219"
opus = "claude-opus-4-20250514"

[agents.gemini.models]
flash = "gemini-2.0-flash-exp"
pro-exp = "gemini-2.0-pro-exp"

[agents.aichat.models]
gpt4-mini = "gpt-4o-mini"
gpt4 = "gpt-4o"
claude = "claude-3-5-sonnet-20241022"
```

**Environment Variables:**

Use standard shell syntax to set environment variables in the `command` field:

```toml
# Single variable
[agents.gemini]
command = "GEMINI_SYSTEM_MD='{role_file}' gemini --model {model} '{prompt}'"

# Multiple variables
[agents.custom]
command = "VAR1='{role}' VAR2='{date}' custom-ai --model {model} '{prompt}'"
```

**Benefits:**
- Standard and familiar syntax
- No separate env section needed
- Supports multiple variables easily
- Works with any shell command syntax

See [DR-005](./design/decisions/dr-005-role-configuration.md) for details on `{role}` and `{role_file}` placeholders.

---

### [roles.\<name\>]

Named role (system prompt) configurations. Global and local configs are combined; local overrides global for same role name.

Uses **[Unified Template Design (UTD)](./unified-template-design.md)** pattern.

**UTD Fields:**

- `file` (string, optional) - Path to role content file
- `command` (string, optional) - Shell command for dynamic content
- `prompt` (string, optional) - Template text with placeholders: `{file}`, `{file_contents}`, `{command}`, `{command_output}`

At least one UTD field must be present. See [UTD documentation](./design/unified-template-design.md) for complete validation rules.

**Role-Specific Fields:**

**description** (string, optional)
: Human-readable description of this role. Displayed in `start config role list`.

```toml
[roles.code-reviewer]
description = "Expert code reviewer focusing on security"
```

**Additional Fields:**

- `shell` (string, optional) - Override global shell for command execution
- `command_timeout` (integer, optional) - Override global timeout for command execution

**Role Selection:**

Roles are selected using precedence rules:
1. CLI `--role` flag (highest priority)
2. Task `role` field (if executing a task)
3. `default_role` setting
4. First role in config (TOML order)

**Default Role:**

```toml
[settings]
default_role = "code-reviewer"
```

**Placeholders:**

Roles are passed to agents via two placeholders:
- `{role}` - Inline content (fully resolved, for Claude, aichat)
- `{role_file}` - File path (for Gemini and file-based agents)

See [DR-005](./design/decisions/dr-005-role-configuration.md) and [DR-007](./design/decisions/dr-007-placeholders.md) for details.

**Merge behavior:**

- Global + local roles are combined
- Local role completely replaces global role with same name (no field merging)
- All roles available for selection

**Validation:**

- Role name must match: `/^[a-z0-9]+(-[a-z0-9]+)*$/`
- At least one UTD field required (`file`, `command`, or `prompt`)
- `default_role` must reference existing role (if specified)

**Examples:**

```toml
# Simple role (file only)
[roles.general-assistant]
description = "General purpose AI assistant"
file = "~/.config/start/roles/general.md"
```

```toml
# Role with template framing
[roles.code-reviewer]
description = "Code reviewer with context"
file = "~/.config/start/roles/reviewer.md"
prompt = """
{file_contents}

Additional Instructions:
- Focus on security and performance
- Check edge cases
- Verify error handling
"""
```

```toml
# Role with dynamic content
[roles.go-expert]
description = "Go language expert"
file = "~/.config/start/roles/go-base.md"
command = "go version 2>/dev/null || echo 'Go not installed'"
prompt = """
{file_contents}

Environment: {command_output}

Apply Go-specific best practices.
"""
```

```toml
# Inline role (no file)
[roles.documentation-writer]
description = "Technical documentation specialist"
prompt = """
You are a technical documentation specialist.

Guidelines:
- Clear, concise language
- Include code examples
- Focus on user needs
- Use active voice
"""
```

```toml
# Project-specific role (local config)
[roles.project-reviewer]
description = "Project-specific reviewer"
file = "./ROLE.md"
```

See [UTD Examples](./unified-template-design.md#examples) for more patterns.

---

### [context.\<name\>]

Context documents to include in prompts. Named sections allow targeted overrides.

Local and global contexts are **combined** (not replaced).

Uses **[Unified Template Design (UTD)](./unified-template-design.md)** pattern.

**UTD Fields:**

- `file` (string, optional) - Path to context document file
- `command` (string, optional) - Shell command for dynamic content
- `prompt` (string, optional) - Template text with placeholders: `{file}`, `{file_contents}`, `{command}`, `{command_output}`

At least one UTD field must be present. See [UTD documentation](./design/unified-template-design.md) for complete validation rules.

**Context-Specific Fields:**

**description** (string, optional)
: Human-readable description of this context. Displayed in `start config show` and validation output.

```toml
[context.environment]
description = "User environment and tool configuration"
```

**required** (boolean, optional, default: false)
: Whether this document is required context.

**Behavior by command:**
- `required = true`: Included in `start`, `start prompt`, and `start task`
- `required = false`: Included in `start` only, excluded from `start prompt` and `start task`

**Use cases:**
- `required = true`: Critical context needed for all interactions (ENVIRONMENT.md, AGENTS.md, INDEX.csv)
- `required = false`: Optional context for full sessions only (PROJECT.md, README.md)

```toml
[context.environment]
required = true  # Always included

[context.project]
required = false  # Only in 'start' command
```

**shell** (string, optional)
: Override global shell for command execution in this context.

**command_timeout** (integer, optional)
: Override global timeout for command execution in this context.

**Context names:**

- Lowercase, alphanumeric, hyphens only
- Pattern: `/^[a-z0-9]+(-[a-z0-9]+)*$/`
- Examples: `environment`, `project`, `design-docs`

**Document order:**

Documents appear in the prompt in **definition order**:
1. Global contexts (in TOML order)
2. Local contexts (in TOML order)

Rearrange config definitions to change prompt order.

**Merge behavior:**

- Global + local contexts are combined
- Order: Global contexts first, then local contexts
- If name conflict: Local overrides global (intentional override, not an error)

**Examples:**

```toml
# Simple file
[context.environment]
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
description = "User environment and tools"
required = true

# Inline text
[context.note]
prompt = "Important: This project uses Go 1.21"
description = "Project-specific note"
required = true

# Dynamic content from command
[context.git-status]
command = "git status --short"
prompt = "Working tree status:\n{command_output}"
description = "Current git status"
required = false

# Combined: file + command
[context.project-state]
file = "./PROJECT.md"
command = "git log -5 --oneline"
prompt = """
{file_contents}

Recent commits:
{command_output}
"""
required = true
```

See [UTD Examples](./unified-template-design.md#examples) for more patterns.

**Merged result:**

When both global and local configs exist, contexts are combined in order:

1. environment (global, required)
2. index (global, required)
3. readme (global, optional)
4. agents (local, required)
5. project (local, optional)
6. note (local, required, inline)

---

### [tasks.\<name\>]

Predefined workflow tasks. Can be defined in both global and local configs.

Tasks define reusable workflows with specific system prompts, prompt templates, and optional dynamic content. Tasks use the **[Unified Template Design (UTD)](./unified-template-design.md)** pattern for both system prompts and task prompts.

**Metadata Fields:**

**alias** (string, optional)
: Short name for quick access. Must be unique across all tasks (global + local merged).

```toml
[tasks.git-diff-review]
alias = "gdr"
```

**description** (string, optional)
: Help text displayed in task list and help output.

```toml
[tasks.git-diff-review]
description = "Review staged git changes"
```

**agent** (string, optional)
: Preferred agent for this task. Must reference an agent defined in `[agents.<name>]` configuration. Agent selection precedence: CLI `--agent` flag > task `agent` field > `default_agent` setting.

```toml
[tasks.go-review]
agent = "go-expert"
description = "Review Go code with specialized agent"
```

Validated at task execution time and by `start doctor` / `start config validate`.

**role** (string, optional)
: Preferred role for this task. Must reference a role defined in `[roles.<name>]` configuration. Role selection precedence: CLI `--role` flag > task `role` field > `default_role` setting > first role in config.

```toml
[tasks.security-audit]
role = "security-auditor"
agent = "claude"
description = "Security-focused code audit"
```

Validated at task execution time and by `start doctor` / `start config validate`. See [DR-005](./design/decisions/dr-005-role-configuration.md) for role configuration details.

**Task Prompt (UTD Pattern):**

At least one of `file`, `command`, or `prompt` must be present.

**file** (string, optional)
: Path to prompt template file.

**command** (string, optional)
: Shell command to generate dynamic content (e.g., `git diff --staged`). Output available via `{command_output}` placeholder.

**prompt** (string, optional)
: Template text with `{file}`, `{file_contents}`, `{command}`, `{command_output}`, and `{instructions}` placeholders.

`````toml
[tasks.git-diff-review]
command = "git diff --staged"
prompt = """
Review the following changes:

## Instructions
{instructions}

## Staged Changes
```diff
{command_output}
```
"""
`````

**Additional Fields:**

**shell** (string, optional)
: Override global shell for command execution. See [UTD shell configuration](./unified-template-design.md#shell-configuration).

**command_timeout** (integer, optional)
: Override global timeout (in seconds) for command execution.

**Context Inclusion:**

Tasks automatically include **all contexts where `required = true`**.

This ensures critical context (like AGENTS.md, ENVIRONMENT.md) is always present while excluding optional contexts. See the `required` field documentation in the **[context.\<name\>]** section above for complete behavior across all commands.

**Example task (full):**

`````toml
[tasks.git-diff-review]
alias = "gdr"
description = "Review staged git changes"
role = "code-reviewer"
agent = "claude"

# Task prompt (UTD)
command = "git diff --staged"
shell = "bash"
command_timeout = 10
prompt = """
Review the following changes:

## Instructions
{instructions}

## Staged Changes
```diff
{command_output}
```
"""
`````

**Example task (minimal):**

```toml
[tasks.simple]
alias = "s"
description = "Simple help task"
prompt = "Help me with: {instructions}"
# No role field = uses default_role or first role in config
# Auto-includes all contexts with required = true
```

**Placeholder behavior:**

In task prompt templates:
- `{file}` - File path from task `file`
- `{file_contents}` - Content from task `file`
- `{command}` - Command string from task `command`
- `{command_output}` - Output from task `command`
- `{instructions}` - Command-line args (or `"None"`)
- `{model}`, `{date}` - Global placeholders

See [DR-007](./design/design-records/dr-007-placeholders.md) for complete placeholder documentation.

**Usage:**

```bash
start task git-diff-review "focus on security"
start task gdr "ignore formatting changes"
start task simple "explain Go interfaces"
```

---

## Placeholders

Placeholders are variables expanded during command execution.

### Global Placeholders

Available in agent command templates, prompts, and environment variables:

**{model}**
: Resolved model name (full identifier)

Example: `claude-3-7-sonnet-20250219`

**{role}**
: Fully resolved role content (inline text). Use for agents that accept system prompts inline (Claude, aichat).

**{role_file}**
: File path to role content. Use for agents that require file-based system prompts (Gemini).
  - **Simple roles (file only):** Points to original file path
  - **UTD roles (with command/prompt):** Role is evaluated, saved to temp file (`/tmp/start-role-*.md`), path points to temp file. Temp file cleaned up after agent execution.

**{prompt}**
: Assembled prompt text from context documents and custom prompts.

**{date}**
: Current timestamp in ISO 8601 format with timezone.

Example: `2025-01-04T14:30:00+10:00`

### UTD Pattern Placeholders

Used in `prompt` field of `[context.<name>]`, `[roles.<name>]`, and `[tasks.<name>]`:

**{file}**
: File path from the `file` field (absolute, with ~ expanded).

Example: `file = "~/reference/ENVIRONMENT.md"` → `{file}` = `"/Users/username/reference/ENVIRONMENT.md"`

**{file_contents}**
: Content from the `file` field. The file is read and its contents replace the placeholder.

Example for injecting contents:
```toml
[context.environment]
file = "~/reference/ENVIRONMENT.md"
prompt = """
Environment Context:
{file_contents}
"""
```

Example for instructing AI to read (for agents with file access):
```toml
[context.environment]
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
```

**{command}**
: Command string from the `command` field.

Example: `command = "git status"` → `{command}` = `"git status"`

**{command_output}**
: Output from executing the `command` field (stdout/stderr combined).

Example: `command = "git diff --staged"` → `{command_output}` = output of git diff

### Task-Specific Placeholders

Available in task prompt templates:

**{instructions}**
: Command-line arguments after task name. Value is `"None"` if no arguments provided.

Example: `start task gdr "focus on security"` → `{instructions}` = `"focus on security"`

Tasks also have access to all UTD Pattern Placeholders: `{file}`, `{file_contents}`, `{command}`, `{command_output}`.

---

## Path Resolution

### Tilde Expansion

Paths with `~` expand to user's home directory:

```toml
file = "~/reference/ENVIRONMENT.md"
# Expands to: /Users/username/reference/ENVIRONMENT.md
```

### Relative Paths

Relative paths resolve based on context:

**In global config:**
- Relative to home directory or config directory (context-dependent)

**In local config:**
- Relative to working directory (current directory or `--directory` flag)

**Examples:**

```toml
file = "./AGENTS.md"           # Relative (same as "AGENTS.md")
file = "AGENTS.md"             # Relative
file = "/absolute/path.md"     # Absolute
file = "~/reference/file.md"   # Home-relative (tilde expansion)
```

---

## Validation Rules

### Required Fields

**[agents.\<name\>]:**
- `command` must be present
- `command` should contain `{prompt}` placeholder (warns if missing)

**[context.\<name\>]:**
- At least one of `file`, `command`, or `prompt` must be present (UTD pattern)
- UTD validation rules apply (see [Unified Template Design](./unified-template-design.md#validation-rules))

**[tasks.\<name\>]:**
- At least one of `file`, `command`, or `prompt` must be present (task prompt)
- UTD validation rules apply (see [Unified Template Design](./unified-template-design.md#validation-rules))
- `agent` field (if present) must reference an existing `[agents.<name>]` section

### Field Constraints

**Agent names:**
- Lowercase alphanumeric with hyphens
- Pattern: `/^[a-z0-9]+(-[a-z0-9]+)*$/`
- Examples: `claude`, `gemini`, `my-agent`

**Model names:**
- Same constraints as agent names
- Examples: `haiku`, `sonnet`, `gpt4-mini`

**Task names:**
- Same constraints as agent names
- Examples: `code-review`, `gdr`, `ct`

**Task aliases:**
- Same constraints as agent names
- Must be unique across all tasks

**Context document names:**
- Same constraints as agent names
- Must be unique across global + local configs

### Scope Constraints

**Allowed in both global and local:**
- `[settings]` - Local overrides global
- `[roles.<name>]` - Combined (global + local), local overrides global for same name
- `[agents.<name>]` - Combined (global + local), local overrides global for same name
- `[context.<name>]` - Combined (global + local)
- `[tasks.<name>]` - Combined (global + local), local overrides global for same name

---

## Complete Example

(See earlier example at top of document - same structure applies here)

### Merged Result

When both configs exist, the effective configuration is:

**Settings:**
- `default_agent = "claude"` (from global)
- `log_level = "verbose"` (from local, overrides global)

**Agents:**
- claude, gemini, aichat (from global; local can override or add agents)

**Roles:**
- code-reviewer: `./ROLE.md` (from local, overrides global)
- coder: `~/.config/start/roles/coder.md` (from global)
- reviewer: `~/.config/start/roles/reviewer.md` (from global)

**Context documents (in order):**
1. environment (global, required)
2. index (global, required)
3. readme (global, optional)
4. agents (local, required)
5. project (local, optional)
6. design (local, optional)

---

## Best Practices

### Agent Configuration

**Use descriptive model names:**

```toml
# Good - clear purpose
[agents.claude.models]
fast = "claude-3-5-haiku-20241022"
balanced = "claude-3-7-sonnet-20250219"
powerful = "claude-opus-4-20250514"

# Also good - model family names
[agents.claude.models]
haiku = "claude-3-5-haiku-20241022"
sonnet = "claude-3-7-sonnet-20250219"
opus = "claude-opus-4-20250514"
```

**Include metadata for discoverability:**

```toml
[agents.claude]
description = "Anthropic's Claude AI assistant via Claude Code CLI"
url = "https://docs.claude.com/claude-code"
models_url = "https://docs.anthropic.com/en/docs/about-claude/models"
# ... rest of config
```

### Context Documents

**Mark essential documents as required:**

```toml
# Always needed for context
[context.environment]
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
required = true

# Nice to have, but not essential
[context.changelog]
file = "CHANGELOG.md"
prompt = "Recent changes in {file}"
required = false
```

**Order documents by importance:**

Define most important documents first - they appear first in the prompt:

```toml
# First - critical context
[context.environment]
# ...

# Second - project overview
[context.agents]
# ...

# Third - current work
[context.project]
# ...

# Last - supplementary
[context.design]
# ...
```

**Use clear, actionable prompts:**

```toml
# Good - specific instruction
prompt = "Read {file} for environment context including user info and tools."

# Less good - vague
prompt = "Read {file}"
```

### Settings

**Keep global config minimal:**

Only define what's truly global. Let local configs customize:

```toml
# Global: Just the defaults
[settings]
default_agent = "claude"
log_level = "normal"
```

```toml
# Local: Override when needed
[settings]
log_level = "verbose"  # This project needs detailed output
```

### File Organization

**Global config structure:**

```toml
# 1. Settings
[settings]
# ...

# 2. Agents (main section)
[agents.claude]
# ...

# 3. Roles (system prompts)
[roles.code-reviewer]
# ...

# 4. Shared contexts
[context.environment]
# ...
```

**Local config structure:**

```toml
# 1. Settings overrides
[settings]
# ...

# 2. Project-specific roles (override global)
[roles.code-reviewer]
# ...

# 3. Project contexts
[context.agents]
# ...
```

---

## Selection Precedence

### Agent Selection

Agent selection follows this priority order:

1. `--agent` CLI flag (highest priority)
2. Task `agent` field (if executing a task)
3. `default_agent` setting
4. First agent in config (TOML order)

**Example:**

```bash
start --agent gemini              # Uses gemini (CLI flag)
start task code-review            # Uses task's agent or default_agent
start                             # Uses default_agent or first agent
```

### Role Selection

Role selection follows this priority order:

1. `--role` CLI flag (highest priority)
2. Task `role` field (if executing a task)
3. `default_role` setting
4. First role in config (TOML order)

**Example:**

```bash
start --role security-auditor     # Uses security-auditor (CLI flag)
start task code-review            # Uses task's role or default_role
start                             # Uses default_role or first role
```

### Model Selection

Model selection follows this priority order:

1. `--model` CLI flag (highest priority)
2. Agent's `default_model` field
3. First model in agent's models table (TOML order)

**Example:**

```bash
start --model opus                # Uses opus (CLI flag)
start --agent claude              # Uses claude's default_model
start                             # Uses default agent's default_model
```

---

## See Also

- [Design Record](./design/design-record.md) - Design decisions (DR-001 through DR-029)
- [CLI Reference](./cli/) - Command-line usage documentation
- [Vision](./vision.md) - Product vision and goals
- [Task Configuration](./tasks.md) - Task-specific documentation
