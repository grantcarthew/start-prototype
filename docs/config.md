# Configuration Reference

Complete reference for `start` configuration files.

## Overview

`start` uses TOML configuration files with a two-tier hierarchy:

- **Global:** `~/.config/start/config.toml` - User-wide settings, agents, shared context
- **Local:** `./.start/config.toml` - Project-specific settings and context

**Merge behavior:**
- Settings: Local values override global values
- Agents: Global only (local cannot define agents - see DR-004)
- Contexts: Combined (global + local, names must be unique)
- Roles: Global only
- Tasks: Global only (defined via embedded assets)

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
path = "~/.config/start/roles/coder.md"

[roles.reviewer]
path = "~/.config/start/roles/reviewer.md"
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
- claude, gemini, aichat (from global only - local cannot define agents)

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
- Global settings (default agent, verbosity)
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

AI agent tool configurations. **Global-only** (per DR-004).

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
- Agents defined in local config → **Warning**: `"Agents in local config ignored (global-only)"`

**Scope:**

Agents are **global-only** (per DR-004). Local configs cannot define agents.

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

Predefined workflow tasks. **Global-only** (embedded in binary, not user-configurable yet).

Tasks define reusable workflows with specific prompts, roles, and optional dynamic content.

**Fields:**

**alias** (string, optional)
: Short name for quick access. Must be unique across all tasks.

```toml
[tasks.git-diff-review]
alias = "gdr"
```

**description** (string, optional)
: Help text displayed in task list and help output.

```toml
[tasks.git-diff-review]
description = "Review git diff changes"
```

**role** (string, required)
: System prompt for this task. Can be:
- File path: `"./roles/code-reviewer.md"`
- Inline text: `"You are an expert code reviewer..."`

```toml
[tasks.git-diff-review]
role = "./roles/code-reviewer.md"
```

**documents** (array of strings, optional)
: Names of context documents to include. References document names from `[context.<name>]` sections.

```toml
[tasks.git-diff-review]
documents = ["environment", "agents"]
```

**content_command** (string, optional)
: Shell command to execute. Output becomes `{content}` placeholder in prompt template.

```toml
[tasks.git-diff-review]
content_command = "git diff --staged"
```

**prompt** (string, required)
: Prompt template for this task. Can be:
- File path: `"./prompts/review.md"`
- Inline text (multi-line): See example below

Supports placeholders:
- Task-specific: `{instructions}`, `{content}`
- Global: `{model}`, `{system_prompt}`, `{prompt}`, `{date}`

**Example task (full):**

````toml
[tasks.git-diff-review]
alias = "gdr"
description = "Review git diff changes"
role = "./roles/code-reviewer.md"
documents = ["environment", "agents"]
content_command = "git diff --staged"
prompt = """
Analyze the following git diff and act as a code reviewer.

## Special Instructions

{instructions}

## Git Diff

```diff
{content}
```
"""
````

**Example task (minimal):**

```toml
[tasks.simple]
role = "You are a helpful assistant."
prompt = "Help me with: {instructions}"
```

**Placeholder behavior:**

- `{instructions}` - Command-line args after task name, or `"None"` if not provided
- `{content}` - Output from `content_command`, or empty string if no command
- Standard placeholders also available

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

**{content}**
: Output from task's `content_command`. Empty string if no command defined.

Example: `content_command = "git diff --staged"` → `{content}` = output of git diff

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
path = "AGENTS.md"             # Relative
path = "/absolute/path.md"     # Absolute
path = "~/reference/file.md"   # Home-relative (tilde expansion)
```

---

## Validation Rules

### Required Fields

**[agents.\<name\>]:**
- `command` must be present
- `command` must contain `{prompt}` placeholder

**[context.\<name\>]:**
- `path` must be present
- `prompt` must be present

**[tasks.\<name\>]:**
- `role` must be present
- `prompt` must be present

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
- `[agents]` - Cannot appear in local config (DR-004)
- `[roles]` - Cannot appear in local config
- `[tasks]` - Cannot appear in local config (embedded assets)

**Allowed in both global and local:**
- `[settings]` - Local overrides global
- `[system_prompt]` - Local overrides global
- `[context.documents]` - Combined (global + local)

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
path = "~/.config/start/roles/coder.md"

[roles.reviewer]
path = "~/.config/start/roles/reviewer.md"
```

### Local Config (./.start/config.toml)

```toml
# Project-specific settings
[settings]
verbosity = "verbose"

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
- `verbosity = "verbose"` (from local, overrides global)

**Agents:**
- claude, gemini, aichat (from global only)

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
path = "CHANGELOG.md"
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
verbosity = "verbose"  # This project needs detailed output
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
