# Configuration Reference

Complete reference for `start` configuration files.

## Overview

`start` uses TOML configuration files with a two-tier hierarchy:

- **Global:** `~/.config/start/config.toml` - User-wide settings, agents, shared context
- **Local:** `./.start/config.toml` - Project-specific settings and context

**Merge behavior:**
- Settings: Local values override global values
- Agents: Combined (global + local), local overrides global for same name
- Contexts: Combined (global + local, names must be unique)
- Roles: Global only
- Tasks: Combined (global + local), local overrides global for same name

## Complete Example

### Global Config (~/.config/start/config.toml)

```toml
# Global settings
[settings]
default_agent = "claude"
log_level = "normal"
shell = "bash"
command_timeout = 30

# Agent configurations
[agents.claude]
description = "Anthropic's Claude AI assistant via Claude Code CLI"
url = "https://docs.claude.com/claude-code"
models_url = "https://docs.anthropic.com/en/docs/about-claude/models"
command = "claude --model {model} --append-system-prompt '{system_prompt}' '{prompt}'"
default_model = "sonnet"

  [agents.claude.models]
  haiku = "claude-3-5-haiku-20241022"
  sonnet = "claude-3-7-sonnet-20250219"
  opus = "claude-opus-4-20250514"

[agents.gemini]
description = "Google's Gemini AI via CLI"
url = "https://github.com/example/gemini-cli"
models_url = "https://ai.google.dev/models/gemini"
command = "gemini --model {model} '{prompt}'"
default_model = "flash"

  [agents.gemini.models]
  flash = "gemini-2.0-flash-exp"
  pro-exp = "gemini-2.0-pro-exp"

  [agents.gemini.env]
  GEMINI_SYSTEM_MD = "{system_prompt}"

[agents.aichat]
description = "All-in-one multi-provider AI chat tool"
url = "https://github.com/sigoden/aichat"
command = "aichat --model {model} '{prompt}'"
default_model = "gpt4-mini"

  [agents.aichat.models]
  gpt4-mini = "gpt-4o-mini"
  gpt4 = "gpt-4o"
  claude = "claude-3-5-sonnet-20241022"

# System prompt
[system_prompt]
file = "~/.config/start/roles/default.md"

# Global context documents
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

# Roles
[roles.coder]
file = "~/.config/start/roles/coder.md"

[roles.reviewer]
file = "~/.config/start/roles/reviewer.md"
```

### Local Config (./.start/config.toml)

```toml
# Project-specific settings (overrides global)
[settings]
log_level = "debug"

# Override system prompt for this project
[system_prompt]
file = "./ROLE.md"

# Project-specific context documents (combined with global)
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

**System prompt:**
- `./ROLE.md` (from local, overrides global `~/.config/start/roles/default.md`)

**Context documents (in definition order):**
1. environment - `~/reference/ENVIRONMENT.md` (global, required)
2. index - `~/reference/INDEX.csv` (global, required)
3. readme - `README.md` (global, optional)
4. agents - `./AGENTS.md` (local, required)
5. project - `./PROJECT.md` (local, optional)
6. design - `./docs/design-record.md` (local, optional)

**Roles:**
- coder, reviewer (from global only - local cannot define roles)

## File Locations

### Global Config

```
~/.config/start/config.toml
```

**Purpose:**
- Agent configurations (command templates, models)
- Global settings (default agent, log_level)
- Shared context documents (ENVIRONMENT.md, INDEX.csv, etc.)
- Role templates
- User-wide defaults

**Created by:** `start init`

### Local Config

```
./.start/config.toml
```

**Purpose:**
- Project-specific context documents (PROJECT.md, AGENTS.md, etc.)
- Project-specific settings overrides
- Local customizations

**Created by:** Manual creation or `start init` in project directory

## Configuration Sections

### [settings]

Global behavior settings. Local overrides global.

**Fields:**

**default_agent** (string, optional)
: Which agent to use when `--agent` flag not provided. Must match an agent name defined in `[agents]` section.

```toml
[settings]
default_agent = "claude"
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

**Validation:**

All fields use soft validation with fallback defaults:

- **default_agent** misconfigured or agent not found → **Warning**, fall back to first agent in config (TOML order)
- **log_level** invalid value → **Warning**, fall back to `"normal"`
- **shell** not found in PATH → **Warning**, fall back to auto-detected shell (`bash` or `sh`)
- **command_timeout** invalid → **Warning**, fall back to 30 seconds
- Missing fields → Silent, use defaults

**Example:**

```toml
[settings]
default_agent = "claude"
log_level = "normal"
shell = "bash"
command_timeout = 30
```

**Merge behavior:**

Local settings override global settings. If a setting is omitted in local config, the global value is used.

---

### [agents.\<name\>]

AI agent tool configurations. Can be defined in both global and local configs.

Each agent section defines how to invoke an AI tool. Agent names should match the actual tool binary name (e.g., `claude`, `gemini`, `aichat`).

**Fields:**

**command** (string, required)
: Command template to execute the agent. Must contain `{prompt}` placeholder. Supports additional placeholders: `{model}`, `{system_prompt}`, `{date}`.

```toml
[agents.claude]
command = "claude --model {model} --append-system-prompt '{system_prompt}' '{prompt}'"
```

**description** (string, optional)
: Human-readable description of the agent. Displayed in `start agent list`.

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
: Model alias to use when `--model` flag not provided. Must be a key in the `[agents.<name>.models]` table. If omitted, first model in `models` table is used.

```toml
[agents.claude]
default_model = "sonnet"
```

**Validation:**

- **command** must contain `{prompt}` placeholder → **Error** if missing
- If command uses `{model}` placeholder:
  - **[agents.\<name\>.models]** section MUST exist with ≥1 model → **Error** if missing
- If **default_model** defined but not in models table → **Warning**, fall back to first model (TOML order)
- Unknown placeholders in command → **Warning**: `"Unknown placeholder {mdoel} (did you mean {model}?)"`
- Same agent name in global and local → **Info**: Local overrides global

**Scope:**

Agents can be defined in both **global and local** configs (per DR-004 update):

- **Global agents:** `~/.config/start/config.toml` - Personal configurations
- **Local agents:** `./.start/config.toml` - Team/project configurations (can be committed)
- **Merge behavior:** Local overrides global for same agent name
- **Use case:** Teams can commit `.start/` with standard configs; individuals maintain personal preferences

**Example agent (full):**

```toml
[agents.claude]
description = "Anthropic's Claude AI assistant via Claude Code CLI"
url = "https://docs.claude.com/claude-code"
models_url = "https://docs.anthropic.com/en/docs/about-claude/models"
command = "claude --model {model} --append-system-prompt '{system_prompt}' '{prompt}'"
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

Model alias mappings for the agent. User-defined aliases map to full model identifiers.

**Structure:**

```toml
[agents.<name>.models]
<alias> = "<full-model-identifier>"
```

**Alias names:**
- User-defined (not hardcoded)
- lowercase, alphanumeric, hyphens
- Each agent defines its own aliases
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

#### [agents.\<name\>.env]

Environment variables to set when executing the agent command. Optional section.

**Structure:**

```toml
[agents.<name>.env]
<KEY> = "<value>"
```

**Placeholder support:**
- Values can contain placeholders: `{model}`, `{system_prompt}`, `{prompt}`, `{date}`
- Expanded before setting environment variable

**Example:**

```toml
[agents.gemini]
command = "gemini --model {model} '{prompt}'"

  [agents.gemini.env]
  GEMINI_SYSTEM_MD = "{system_prompt}"
  GEMINI_API_KEY = "your-api-key-here"
```

---

### [system_prompt]

System prompt configuration. Optional section. Local overrides global.

Uses **[Unified Template Design (UTD)](./unified-template-design.md)** pattern.

**UTD Fields:**

- `file` (string, optional) - Path to system prompt file
- `command` (string, optional) - Shell command for dynamic content
- `prompt` (string, optional) - Template text with `{file}` and `{command}` placeholders

At least one field must be present. See [UTD documentation](./unified-template-design.md) for complete validation rules and examples.

**Additional Fields:**

- `shell` (string, optional) - Override global shell for command execution
- `command_timeout` (integer, optional) - Override global timeout for command execution

**Behavior:**

- System prompt passed to agent via `{system_prompt}` placeholder in agent command
- Not all agents support system prompts
- Section can be omitted entirely (no warning)

**Merge behavior:**

Local section completely replaces global section. If local section missing, use global.

**Examples:**

```toml
# Simple file
[system_prompt]
file = "./ROLE.md"
```

```toml
# Inline text
[system_prompt]
prompt = """
You are an expert code reviewer.
Focus on security and performance.
"""
```

```toml
# File with framing
[system_prompt]
file = "./ROLE.md"
prompt = """
Role Definition:
{file}

Follow these instructions carefully.
"""
```

```toml
# Dynamic content from command
[system_prompt]
command = "git log -1 --format='%s'"
prompt = """
You are a code reviewer.
Current commit: {command}
"""
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
- `prompt` (string, optional) - Template text with `{file}` and `{command}` placeholders

At least one UTD field must be present. See [UTD documentation](./unified-template-design.md) for complete validation rules.

**Context-Specific Fields:**

**description** (string, optional)
: Human-readable description of this context. Displayed in `start config show` and validation output.

```toml
[context.environment]
description = "User environment and tool configuration"
```

**required** (boolean, optional, default: false)
: Whether this document is required context.

- `true` - Included by both `start` and `start prompt`
- `false` - Included by `start`, excluded by `start prompt`

```toml
[context.environment]
required = true
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
prompt = "Working tree status:\n{command}"
description = "Current git status"
required = false

# Combined: file + command
[context.project-state]
file = "./PROJECT.md"
command = "git log -5 --oneline"
prompt = """
{file}

Recent commits:
{command}
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

**System Prompt Override (UTD Pattern):**

Tasks can override the global/local system prompt using UTD fields. All fields are optional - if omitted, uses global/local `[system_prompt]`.

**system_prompt_file** (string, optional)
: Path to role definition file.

**system_prompt_command** (string, optional)
: Shell command to generate dynamic role content.

**system_prompt** (string, optional)
: Template text with `{file}` and `{command}` placeholders.

```toml
[tasks.code-review]
system_prompt_file = "~/.config/start/roles/code-reviewer.md"
system_prompt = """
{file}

CRITICAL: Focus on security vulnerabilities.
"""
```

See [UTD documentation](./unified-template-design.md#validation-rules) for field combination behavior.

**Task Prompt (UTD Pattern):**

At least one of `file`, `command`, or `prompt` must be present.

**file** (string, optional)
: Path to prompt template file.

**command** (string, optional)
: Shell command to generate dynamic content (e.g., `git diff --staged`). Output available via `{command}` placeholder.

**prompt** (string, optional)
: Template text with `{file}`, `{command}`, and `{instructions}` placeholders.

```toml
[tasks.git-diff-review]
command = "git diff --staged"
prompt = """
Review the following changes:

## Instructions
{instructions}

## Staged Changes
```diff
{command}
```
"""
```

**Additional Fields:**

**shell** (string, optional)
: Override global shell for command execution. See [UTD shell configuration](./unified-template-design.md#shell-configuration).

**command_timeout** (integer, optional)
: Override global timeout (in seconds) for command execution.

**Context Inclusion:**

Tasks automatically include **all contexts where `required = true`**. There is no `documents` array. This ensures critical context (like AGENTS.md) is always present.

**Example task (full):**

````toml
[tasks.git-diff-review]
alias = "gdr"
description = "Review staged git changes"

# System prompt override (UTD)
system_prompt_file = "~/.config/start/roles/code-reviewer.md"
system_prompt = """
{file}

Additional context: Focus on security and performance.
"""

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
{command}
```
"""
````

**Example task (minimal):**

```toml
[tasks.simple]
alias = "s"
description = "Simple help task"
prompt = "Help me with: {instructions}"
# No system_prompt_* fields = uses global/local [system_prompt]
# Auto-includes all contexts with required = true
```

**Placeholder behavior:**

In system_prompt templates:
- `{file}` - Content from `system_prompt_file`
- `{command}` - Output from `system_prompt_command`
- `{model}`, `{date}` - Global placeholders

In task prompt templates:
- `{file}` - Content from task `file`
- `{command}` - Output from task `command`
- `{instructions}` - Command-line args (or `"None"`)
- `{model}`, `{date}` - Global placeholders

**Usage:**

```bash
start task git-diff-review "focus on security"
start task gdr "ignore formatting changes"
start task simple "explain Go interfaces"
```

---

### [roles.\<name\>]

Role template definitions. **Global-only**.

Roles are reusable system prompts stored as named references.

**Fields:**

**path** (string, required)
: File path to role definition.

```toml
[roles.coder]
path = "~/.config/start/roles/coder.md"

[roles.reviewer]
path = "~/.config/start/roles/reviewer.md"
```

**Usage:**

Roles are referenced by name in task definitions:

```toml
[tasks.code-review]
role = "reviewer"  # References [roles.reviewer]
prompt = "Review this code..."
```

Or used directly:

```toml
[tasks.custom]
role = "./CUSTOM-ROLE.md"  # Direct file path
prompt = "..."
```

---

## Placeholders

Placeholders are variables expanded during command execution.

### Global Placeholders

Available in agent command templates, prompts, and environment variables:

**{model}**
: Resolved model name (full identifier, not alias)

Example: `claude-3-7-sonnet-20250219`

**{system_prompt}**
: Contents of system prompt file. Empty string if not configured or file doesn't exist.

**{prompt}**
: Assembled prompt text from context documents and custom prompts.

**{date}**
: Current timestamp in ISO 8601 format with timezone.

Example: `2025-01-04T14:30:00+10:00`

### Context-Specific Placeholders

**{file}** (context documents only)
: File path of the document. Used in `prompt` field of `[context.<name>]`.

Example: `"Read {file} for context"` → `"Read ~/reference/ENVIRONMENT.md for context"`

### Task-Specific Placeholders

Available in task prompt templates:

**{instructions}**
: Command-line arguments after task name. Value is `"None"` if no arguments provided.

Example: `start task gdr "focus on security"` → `{instructions}` = `"focus on security"`

**{command}**
: Output from task's `command`. Empty string if no command defined.

Example: `command = "git diff --staged"` → `{command}` = output of git diff

Note: In task prompts, use `{command}` not `{content}`. The `{content}` placeholder was from an earlier design.

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
- `command` must contain `{prompt}` placeholder

**[context.\<name\>]:**
- At least one of `file`, `command`, or `prompt` must be present (UTD pattern)
- UTD validation rules apply (see [Unified Template Design](./unified-template-design.md#validation-rules))

**[tasks.\<name\>]:**
- At least one of `file`, `command`, or `prompt` must be present (task prompt)
- UTD validation rules apply (see [Unified Template Design](./unified-template-design.md#validation-rules))

### Field Constraints

**Agent names:**
- Lowercase alphanumeric with hyphens
- Pattern: `/^[a-z0-9]+(-[a-z0-9]+)*$/`
- Examples: `claude`, `gemini`, `my-agent`

**Model aliases:**
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

**Global-only sections:**
- `[roles]` - Cannot appear in local config

**Allowed in both global and local:**
- `[settings]` - Local overrides global
- `[system_prompt]` - Local overrides global
- `[agents]` - Combined (global + local), local overrides global for same name
- `[context]` - Combined (global + local)
- `[tasks]` - Combined (global + local), local overrides global for same name

---

## Complete Example

### Global Config (~/.config/start/config.toml)

```toml
# Global settings
[settings]
default_agent = "claude"
log_level = "normal"
shell = "bash"
command_timeout = 30

# Agent configurations
[agents.claude]
description = "Anthropic's Claude AI assistant via Claude Code CLI"
url = "https://docs.claude.com/claude-code"
models_url = "https://docs.anthropic.com/en/docs/about-claude/models"
command = "claude --model {model} --append-system-prompt '{system_prompt}' '{prompt}'"
default_model = "sonnet"

  [agents.claude.models]
  haiku = "claude-3-5-haiku-20241022"
  sonnet = "claude-3-7-sonnet-20250219"
  opus = "claude-opus-4-20250514"

[agents.gemini]
description = "Google's Gemini AI via CLI"
url = "https://github.com/example/gemini-cli"
models_url = "https://ai.google.dev/models/gemini"
command = "gemini --model {model} '{prompt}'"
default_model = "flash"

  [agents.gemini.models]
  flash = "gemini-2.0-flash-exp"
  pro-exp = "gemini-2.0-pro-exp"

  [agents.gemini.env]
  GEMINI_SYSTEM_MD = "{system_prompt}"

[agents.aichat]
description = "All-in-one multi-provider AI chat tool"
url = "https://github.com/sigoden/aichat"
command = "aichat --model {model} '{prompt}'"
default_model = "gpt4-mini"

  [agents.aichat.models]
  gpt4-mini = "gpt-4o-mini"
  gpt4 = "gpt-4o"
  claude = "claude-3-5-sonnet-20241022"

# System prompt
[system_prompt]
file = "~/.config/start/roles/default.md"

# Global context documents
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

# Roles
[roles.coder]
file = "~/.config/start/roles/coder.md"

[roles.reviewer]
file = "~/.config/start/roles/reviewer.md"
```

### Local Config (./.start/config.toml)

```toml
# Project-specific settings
[settings]
log_level = "verbose"

# Override system prompt for this project
[system_prompt]
file = "./ROLE.md"

# Project-specific context documents
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

When both configs exist, the effective configuration is:

**Settings:**
- `default_agent = "claude"` (from global)
- `log_level = "verbose"` (from local, overrides global)

**Agents:**
- claude, gemini, aichat (from global; local can override or add agents)

**System prompt:**
- `./ROLE.md` (from local, overrides global)

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

**Use descriptive model aliases:**

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

# 3. System prompt
[system_prompt]
# ...

# 4. Shared contexts
[context.environment]
# ...

# 5. Roles
[roles.coder]
# ...
```

**Local config structure:**

```toml
# 1. Settings overrides
[settings]
# ...

# 2. System prompt override
[system_prompt]
# ...

# 3. Project contexts
[context.agents]
# ...
```

---

## See Also

- [Design Record](./design-record.md) - Design decisions (DR-001 through DR-013)
- [CLI Reference](./cli/) - Command-line usage documentation
- [Vision](./vision.md) - Product vision and goals
- [Task Configuration](./task.md) - Task-specific documentation
