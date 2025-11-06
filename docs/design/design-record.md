# Design Record

This document tracks design decisions made during the development of `start`.

## Overview

See [vision.md](./vision.md) for the product vision and goals.

---

## Decisions

### DR-001: Configuration File Format (2025-01-03)

**Decision:** Use TOML for all configuration files

**Rationale:**

- Human-readable and editable
- No whitespace sensitivity (unlike YAML)
- Excellent Go support via BurntSushi/toml
- Supports comments and complex nested structures
- Used by similar tools (mise)

**Alternatives considered:**

- YAML: Too error-prone with whitespace
- JSON: No comments, less human-friendly
- Custom key-value: Too limited for nested structures

---

### DR-002: Configuration File Structure (2025-01-03)

**Decision:** Single configuration file with global + local merge strategy

**Files:**

- Global: `~/.config/start/config.toml`
- Local: `./.start/config.toml` (project-specific)

**Merge behavior:**

- Local config merges with global
- Same keys in local override global values
- New keys in local are added
- Omitted keys use global defaults

**Rationale:**

- Single file simpler than multiple files
- Merge allows both defaults and project-specific overrides
- CLI commands will manage config, so complexity is hidden from users

---

### DR-003: Named Documents for Context (2025-01-03)

**Decision:** Use named document sections instead of arrays

**Structure:**

```toml
[context.documents.environment]
path = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."

[context.documents.project]
path = "./PROJECT.md"
prompt = "Read {file}. Respond with summary."
```

**Rationale:**

- Names allow local config to override specific documents
- Can't target array items for override
- Enables both override (same name) and add (new name) patterns
- More explicit and readable

**Example use case:**

- Global defines "project" document as `./PROJECT.md`
- Local overrides to `~/multi-repo/BIG-PROJECT.md`
- Local adds new "vision" document as `./docs/vision.md`

---

### DR-004: Agent Configuration Scope (2025-01-03, Updated 2025-01-05)

**Decision:** Agents can be defined in both global and local configs with merge behavior

**Structure:**

```toml
[settings]
default_agent = "claude"

[agents.claude]
command = "claude --model {model} --append-system-prompt '{system_prompt}' '{prompt}'"
default_model = "sonnet"  # Default alias to use

  [agents.claude.models]
  haiku = "claude-3-5-haiku-20241022"
  sonnet = "claude-3-7-sonnet-20250219"
  opus = "claude-opus-4-20250514"

[agents.gemini]
command = "gemini --model {model} '{prompt}'"
default_model = "flash"

  [agents.gemini.models]
  flash = "gemini-2.0-flash-exp"
  pro-exp = "gemini-2.0-pro-exp"
```

**Model alias behavior:**

- Alias names are user-defined (not hardcoded tiers)
- Each agent has its own set of model aliases
- `default_model` specifies which alias to use when `--model` flag not provided
- Users can use `--model <alias>` or `--model <full-model-name>`

**Rationale:**

- Agent names are the actual tool names (claude, gemini, opencode) not arbitrary aliases
- Self-documenting - clear which agents are available
- Flexible aliases allow users to name models meaningfully for their workflow

**Scope:**

**Global agents:** `~/.config/start/config.toml`
- Personal agent configurations
- Individual preferences (model aliases, default models)
- Private configurations

**Local agents:** `./.start/config.toml`
- Team-standardized configurations (can be committed to git)
- Project-specific agent wrappers or custom tools
- Consistent team experience (clone and go)

**Merge behavior:**

- Global + local agents are combined
- Same agent name: **local overrides global**
- Enables team standardization while allowing personal overrides
- Local config in version control ensures consistent team setup

**Example scenario:**

```toml
# Global: ~/.config/start/config.toml (personal preference)
[agents.claude]
default_model = "haiku"  # Fast model for personal use

# Local: ./.start/config.toml (team standard, committed)
[agents.claude]
command = "claude --model {model} --append-system-prompt '{system_prompt}' '{prompt}'"
default_model = "sonnet"  # Team uses better model
  [agents.claude.models]
  sonnet = "claude-3-7-sonnet-20250219"
  opus = "claude-opus-4-20250514"

# Result: Local overrides global when working in this project
```

**Benefits:**

- Team can commit `.start/` directory for consistent configuration
- New team members: clone repo and `start` works immediately (if agents installed)
- Personal global configs for individual workflows
- Project-specific agents for custom tooling

**Security note:**

Don't commit secrets in local agent configs. Use environment variable references:

```toml
# Bad
[agents.custom.env]
API_KEY = "sk-1234567890"  # DON'T COMMIT

# Good
[agents.custom.env]
API_KEY = "${CUSTOM_API_KEY}"  # Reference user's env var
```

**Update (2025-01-04):**

Changed from hardcoded tier names (fast/mid/pro) to flexible user-defined aliases.

**Update (2025-01-05):**

Changed from global-only to allowing both global and local agents. Enables team standardization via committed `.start/` directory while maintaining personal preferences in global config.

---

### DR-005: System Prompt Handling (2025-01-03)

**Decision:** System prompt configured separately from context documents, and is optional

**Structure:**

```toml
[context.system_prompt]
path = "./ROLE.md"

[context.documents.environment]
path = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
```

**Rationale:**

- System prompt is conceptually different from context documents
- Passed to agents differently (via `--system-prompt` flag or `{role}` placeholder in command)
- Separate section makes it clear and allows easy override
- Can be overridden in local config like other context settings

**Optional behavior:**

- System prompt section can be missing entirely
- Path can be empty
- Not all AI agents support system prompts
- `start` will skip system prompt handling if not configured or file doesn't exist

**Example local override:**

```toml
# Local ./.start/config.toml
[context.system_prompt]
path = "~/shared-roles/senior-go-dev.md"
```

---

### DR-006: CLI Command Structure (2025-01-03)

**Decision:** Use Cobra with subcommand pattern and global flags

**Pattern:**

```bash
start <subcommand> [args] [flags]
```

**Core commands:**

```bash
# Root command (no subcommand)
start                              # Launch default session with context
start --agent gemini               # Launch with specific agent

# Task subcommand
start task <name>                  # Run predefined task
start task code-review             # By name
start task cr                      # By alias
start task cr --agent gemini       # With specific agent

# Agent management
start agent add                    # Add new agent (interactive)
start agent list                   # List configured agents
start agent test <name>            # Test agent configuration

# Config management
start config show                  # Display current config
start config init                  # Create default config
start config edit                  # Open config in editor
```

**Global flags (work on all commands):**

```bash
--agent <name>        # Which agent to use (overrides default)
--model <tier>        # Model tier: fast, mid, pro
--directory <path>    # Working directory (default: pwd)
```

**Rationale:**

- Cobra provides robust subcommand support
- Persistent flags work across all subcommands
- Follows kubectl/git patterns (familiar to developers)
- Tasks discovered dynamically from config
- Extensible for future subcommands

**Task implementation:**

- Tasks defined in config are loaded at startup
- Cobra subcommands generated dynamically
- Each task becomes `start task <name>` with alias support
- See [task.md](./task.md) for task configuration details

---

### DR-007: Command Interpolation and Placeholders (2025-01-03)

**Decision:** Use single-brace placeholders with specific supported variables

**Supported placeholders:**

- `{model}` - Model name (e.g., "claude-3-7-sonnet-20250219")
- `{system_prompt}` - System prompt file contents
- `{prompt}` - Built prompt text from context documents
- `{date}` - Current timestamp (ISO 8601 format with timezone)

**Path expansion:**

- `~` - Expands to user's home directory

**Usage examples:**

```toml
[agents.claude]
command = "claude --model {model} --append-system-prompt '{system_prompt}' '{prompt}'"

[agents.gemini]
command = "gemini --model {model} --include-directories ~/reference '{prompt}'"

  [agents.gemini.env]
  GEMINI_SYSTEM_MD = "{system_prompt}"

[context.documents.environment]
path = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
```

**Substitution behavior:**

- `{model}` - Replaced with selected model tier value
- `{system_prompt}` - Replaced with file contents (empty string if not configured)
- `{prompt}` - Replaced with assembled prompt text
- `{date}` - Replaced with current timestamp (e.g., "2025-01-03T10:30:00+10:00")
- `~` - Expanded before command execution

**Rationale:**

- Single braces simpler than double (`{}` vs `{{}}`)
- "system_prompt" is clear and matches standard terminology
- All placeholders optional - agents can use what they need
- Tilde expansion more concise than `{home}` placeholder
- No environment variable substitution (`{env:...}`) - agents inherit environment naturally

**Not supported:**

- `{env:VAR}` - Environment variables (use agent `env` section instead)
- `{home}` - Use `~` instead
- `{cwd}` - Use `--directory` flag if needed

---

### DR-008: Context File Detection and Handling (2025-01-03)

**Decision:** Relative paths resolve to working directory; missing files skipped with status display

**Path resolution:**

- Relative paths (e.g., `./AGENTS.md` or `AGENTS.md`) resolve to working directory
- Working directory defaults to current directory (`pwd`)
- Override with `--directory` flag
- Absolute paths and `~` paths resolve independently of working directory

**Path equivalence:**

```toml
path = "./AGENTS.md"   # Same
path = "AGENTS.md"     # Same
path = "/absolute/path/file.md"  # Absolute
path = "~/reference/file.md"     # Home-relative
```

**Missing file behavior:**

- Files that don't exist are skipped silently
- No error, no exit
- Missing files excluded from prompt
- Status displayed to user before execution

**Output format:**

```
Starting AI Agent
===============================================================================================
Agent: claude (model: claude-sonnet-4-5@20250929)

Context documents:
  ✓ environment     ~/reference/ENVIRONMENT.md
  ✓ index          ~/reference/INDEX.csv
  ✗ agents         ./AGENTS.md (not found)
  ✗ project        ./PROJECT.md (not found)

System prompt: ./ROLE.md

Executing command...
❯ claude --model claude-sonnet-4-5@20250929 --append-system-prompt '...' '2025-11-03...'
```

**Rationale:**

- Users see exactly what context is being used
- Can diagnose path issues easily
- Optional files work naturally (not all projects have all documents)
- No false errors for legitimately missing optional files
- Command display truncates system prompt ('...') to avoid noise
- Full prompt visible in agent chat once started

**Working directory examples:**

```bash
# Default - uses pwd
cd ~/my-project
start  # Looks for ~/my-project/AGENTS.md

# Override working directory
start --directory ~/my-project

# From anywhere
cd ~
start --directory ~/my-project  # Still finds ~/my-project/AGENTS.md
```

---

### DR-009: Task Structure and Placeholders (2025-01-03)

**Decision:** Tasks include prompt templates with {instructions} and {content} placeholders

**Complete task configuration:**

````toml
[task.git-diff-review]
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
````

"""

````

**Task fields:**
- **alias** (optional) - Short name for quick access
- **description** (optional) - Help text
- **role** (required) - System prompt (file path or inline text)
- **documents** (optional) - Array of named context documents to include
- **content_command** (optional) - Shell command to run, output becomes {content}
- **prompt** (required) - Prompt template (file path or inline text)

**Task-specific placeholders:**
- `{instructions}` - Command-line arguments after task name
  - Value: User's arguments or "None" if not provided
  - Usage: `start task gdr "focus on security"` → `{instructions}` = "focus on security"
  - Usage: `start task gdr` → `{instructions}` = "None"
- `{content}` - Output from `content_command`
  - Value: Command output or empty string if no content_command
  - Example: `git diff --staged` output

**Prompt field format:**
```toml
# Option 1: File path
prompt = "./prompts/gdr-prompt.md"

# Option 2: Inline text
prompt = """
Review this code...
{content}
"""
````

**All placeholders available in task prompts:**

- Task-specific: `{instructions}`, `{content}`
- Global: `{model}`, `{system_prompt}`, `{prompt}`, `{date}`

**Usage examples:**

```bash
# Basic task
start task git-diff-review

# Task with instructions
start task git-diff-review "only focus on security issues"
start task gdr "ignore comment changes"

# Task with agent override
start task gdr --agent gemini "check for performance issues"
```

**Rationale:**

- Mirrors existing bash script pattern (gdr, ucl, etc.)
- Flexible: simple tasks don't need content_command or instructions
- Clear placeholder names ({instructions}, {content})
- "None" default matches existing bash script behavior
- Supports both dynamic content (git diff) and static prompts

---

### DR-010: Default Task Definitions (2025-01-03)

**Decision:** Ship four interactive review tasks as defaults

**Default tasks:**

1. **code-review** (alias: `cr`) - General code review
2. **git-diff-review** (alias: `gdr`) - Review git diff output
3. **comment-tidy** (alias: `ct`) - Review and tidy code comments
4. **doc-review** (alias: `dr`) - Review and improve documentation

**Rationale:**

- All tasks are **interactive reviews** - user works with agent in chat
- No tasks that write files or require orchestration (commit messages, gitignore generation)
- Stays true to vision: launcher only, not a workflow orchestrator
- Users can add non-interactive tasks in their own config if desired

**Tasks NOT included:**

- commit-message - Requires committing after generation (orchestration)
- gitignore - Requires saving to file (file I/O)
- update-changelog - Requires writing to CHANGELOG.md (file I/O)

**User customization:**

- Users can override any default task by defining same name in config
- Users can add additional tasks (including non-interactive ones)
- Users can remove defaults by not including them

**Implementation:**

- Default tasks embedded in binary
- Loaded first, then merged with user config
- User config takes precedence

---

### DR-011: Asset Distribution and Update System (2025-01-03, updated 2025-01-06)

**Decision:** Assets fetched from GitHub repository; `start init` downloads on first run; `start update` refreshes asset library

**Asset structure in repo:**

```
start/
└── assets/
    ├── agents/              # Agent configuration templates
    │   ├── claude.toml
    │   ├── gemini.toml
    │   ├── aichat.toml
    │   ├── openai.toml
    │   └── deepseek.toml
    ├── roles/               # System prompt markdown files
    │   ├── code-reviewer.md
    │   ├── doc-reviewer.md
    │   ├── security-reviewer.md
    │   └── ...
    ├── tasks/               # Default task configurations
    │   ├── code-review.toml
    │   ├── git-diff-review.toml
    │   ├── comment-tidy.toml
    │   ├── doc-review.toml
    │   └── security-review.toml
    └── examples/            # Example configuration files
        ├── global-config.toml
        └── local-config.toml
```

**Asset installation location:**

```
~/.config/start/
├── config.toml              # User's global config
├── .asset-version           # Track asset library version
└── assets/                  # Downloaded asset library
    ├── agents/
    ├── roles/
    ├── tasks/
    └── examples/
```

**Distribution:**

- Assets stored in GitHub repository (`/assets` directory)
- Downloaded on-demand (not embedded in binary)
- Updateable without new release
- `start init` performs initial download
- `start update` refreshes asset library
- Network required for download (can work offline after initial setup)

**Asset version tracking:**

`.asset-version` file format:
```
commit=abc123def456
timestamp=2025-01-06T10:30:00Z
repository=github.com/grantcarthew/start
branch=main
```

**`start init` behavior:**

```bash
$ start init

Welcome to start! Let's configure your AI agents.

Existing config found. Backing up to config.toml.backup

Which agents do you have installed? (select multiple)
[ ] claude
[ ] gemini
[x] opencode
[ ] aichat
[ ] codex
[ ] aider
[ ] Other (custom)

[...continues with setup wizard...]

Config created at ~/.config/start/config.toml
Run 'start' to launch!
```

**File operations:**

1. Check if `~/.config/start/` exists, create if needed
2. If `config.toml` exists, backup to `config.toml.backup` (overwrites old backup)
3. Extract embedded assets to `~/.config/start/`
4. Run interactive wizard to customize config
5. Display success message

**Backup behavior:**

- Automatic backup of existing config.toml
- Terminal message: "Existing config found. Backing up to config.toml.backup"
- Old backup is overwritten (only keep most recent)

**Agent list order (from search results):**

1. claude (Claude Code - Anthropic)
2. gemini (Gemini CLI - Google)
3. aichat (All-in-one multi-provider)
4. opencode (Open-source coding agent)
5. codex (OpenAI Codex CLI)
6. aider (Popular coding assistant)

**Default document detection:**
Check these paths, add to config if they exist:

- `~/reference/ENVIRONMENT.md`
- `~/reference/INDEX.csv`
- `./AGENTS.md`
- `./PROJECT.md`

If none exist, only add `./AGENTS.md` (most common local file)

**Asset usage patterns:**

**Agent templates:**
- Located in `~/.config/start/assets/agents/`
- Used during `start agent add` to pre-fill configurations
- User selects template, values are copied to `config.toml`

**Role files:**
- Located in `~/.config/start/assets/roles/`
- Referenced in config: `file = "~/.config/start/assets/roles/code-reviewer.md"`
- Updates flow automatically when `start update` is run

**Task definitions:**
- Located in `~/.config/start/assets/tasks/`
- Merged with user's task definitions (user tasks take precedence)
- New tasks available immediately after `start update`

**Example configs:**
- Located in `~/.config/start/assets/examples/`
- Reference only, not automatically loaded
- Users manually copy sections to their config

**Update workflow:**

```bash
# Check for updates
start doctor
# Shows: "⚠ Assets are 45 days old. Run 'start update'"

# Update assets
start update
# Downloads latest from GitHub, reports changes

# Assets auto-update next run
start
# References to roles/* automatically use new content
```

**Rationale:**

- Assets updateable without binary release
- New agent configs, roles, tasks available immediately
- Network dependency acceptable (one-time per update)
- Offline work after initial download
- Separation: binary vs content
- Users control update timing (not forced)
- Automatic backup prevents accidental config loss
- Interactive wizard better UX than manual config editing
- Agent order reflects popularity and completeness

---

### DR-012: Context Document Required Field and Order (2025-01-04)

**Decision:** Add optional `required` field to context documents to control inclusion behavior; documents appear in config definition order

**Structure:**

```toml
[context.documents.environment]  # First in prompt
path = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
required = true    # Always included

[context.documents.index]        # Second in prompt
path = "~/reference/INDEX.csv"
prompt = "Read {file} for documentation index."
required = true    # Always included

[context.documents.agents]       # Third in prompt
path = "./AGENTS.md"
prompt = "Read {file} for repository context."
required = true    # Always included

[context.documents.project]      # Fourth in prompt
path = "./PROJECT.md"
prompt = "Read {file}. Respond with summary."
required = false   # Optional (default) - excluded from start prompt
```

**Behavior by command:**

- `start` (root) → Includes ALL documents (required + optional)
- `start prompt` → Includes ONLY required documents
- `start task` → Includes documents specified in task's `documents` array (ignores `required` field)

**Default value:**

- If `required` field is omitted, defaults to `false` (optional document)

**Document order:**

- Documents appear in prompt in the order they are defined in config file
- TOML preserves declaration order within sections
- Users control order by arranging config file
- Predictable and explicit - no alphabetical or other automatic sorting
- Consistent across all commands (start, start prompt, tasks)

**Rationale:**

- `start` provides full context for comprehensive sessions (all documents)
- `start prompt` provides minimal context for focused queries (required only)
- Allows users to designate "essential" vs "nice-to-have" context
- Reduces noise for one-off questions while maintaining critical context
- Tasks maintain full control via explicit `documents` array
- Definition order gives users control over context priority

**Use cases:**

- `~/reference/ENVIRONMENT.md` marked required: Always provides user/environment context (first)
- `~/reference/INDEX.csv` marked required: Always provides documentation index (second)
- `AGENTS.md` marked required: Always provides repository overview (third)
- `PROJECT.md` marked optional: Included for full sessions, excluded for quick queries
- `./DESIGN.md` marked optional: Only for comprehensive reviews

**Zero context scenario:**

Users wanting ONLY custom prompt (no context at all) should use agent directly:

```bash
claude "your prompt"
gemini "your prompt"
```

---

### DR-013: Agent Configuration Distribution via GitHub (2025-01-04)

**Decision:** Fetch agent configurations from GitHub during `start init` rather than embedding in binary

**Structure:**

Agent configs stored in repository:

```
start/
├── assets/
│   ├── agents/
│   │   ├── claude.toml
│   │   ├── gemini.toml
│   │   ├── aichat.toml
│   │   └── ...
│   ├── tasks/
│   └── roles/
```

**Init behavior:**

1. Fetch agent list from GitHub API: `GET /repos/grantcarthew/start/contents/assets/agents`
2. Auto-detect installed agents using `command -v`
3. Download config files for selected agents
4. Merge into user's `~/.config/start/config.toml`

**Technical details:**

- API endpoint: `https://api.github.com/repos/grantcarthew/start/contents/assets/agents`
- Timeout: 10 seconds
- No caching between runs
- Network required (error if offline)
- Rate limit: 60 requests/hour (unauthenticated)

**Rationale:**

**Why fetch instead of embed:**

- Model names change frequently (claude-3-5 → claude-3-7 → claude-4)
- Agent command flags evolve over time
- New agents emerge regularly
- Embedding means stale configs until next release
- Users get current configs without waiting for release

**Trade-offs accepted:**

- Requires network during init (acceptable for one-time setup)
- Dependency on GitHub availability (agents need network anyway)
- No offline init (manual config documented as alternative)

**Update workflow:**

- Model names stale? Update TOML file in repo
- New agent released? Add new config file
- Flag changes? Update command template
- No code changes or releases needed

**Community benefits:**

- Easy to contribute new agent configs (PR a TOML file)
- Clear separation: code vs configuration data
- Living documentation (configs show current best practices)

**Supersedes:** DR-011's embedded assets approach for agent configs. Tasks and roles may still be embedded (TBD).

---

## Pending Decisions

No major decisions remaining. Ready to begin implementation.

---

## References

- [Vision](./vision.md)
- [Design Thoughts](./thoughts.md)
