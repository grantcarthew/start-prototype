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
[context.environment]
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."

[context.project]
file = "./PROJECT.md"
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
[system_prompt]
file = "./ROLE.md"

[context.environment]
file = "~/reference/ENVIRONMENT.md"
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
[system_prompt]
file = "~/shared-roles/senior-go-dev.md"
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
--model <alias>       # Model alias or full model name
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

[context.environment]
file = "~/reference/ENVIRONMENT.md"
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
file = "./AGENTS.md"   # Same
file = "AGENTS.md"     # Same
file = "/absolute/path/file.md"  # Absolute
file = "~/reference/file.md"     # Home-relative
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
command = "git diff --staged"
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
- **command** (optional) - Shell command to run, output becomes {content}
- **prompt** (required) - Prompt template (file path or inline text)

**Task-specific placeholders:**
- `{instructions}` - Command-line arguments after task name
  - Value: User's arguments or "None" if not provided
  - Usage: `start task gdr "focus on security"` → `{instructions}` = "focus on security"
  - Usage: `start task gdr` → `{instructions}` = "None"
- `{content}` - Output from `command`
  - Value: Command output or empty string if no command
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
- Flexible: simple tasks don't need command or instructions
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
├── asset-version.toml       # Track asset library version
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

`asset-version.toml` file format:
```toml
# Asset version tracking - managed by 'start update'

commit = "abc123def456"
timestamp = "2025-01-06T10:30:00Z"
repository = "github.com/grantcarthew/start"
branch = "main"

[files]
"agents/claude.toml" = "a1b2c3d4e5f6..."
"agents/gemini.toml" = "e5f6g7h8i9j0..."
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
[context.environment]  # First in prompt
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
required = true    # Always included

[context.index]        # Second in prompt
file = "~/reference/INDEX.csv"
prompt = "Read {file} for documentation index."
required = true    # Always included

[context.agents]       # Third in prompt
file = "./AGENTS.md"
prompt = "Read {file} for repository context."
required = true    # Always included

[context.project]      # Fourth in prompt
file = "./PROJECT.md"
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

### DR-014: GitHub Asset Download Strategy (2025-01-06)

**Decision:** Use GitHub Tree API with SHA-based caching for incremental asset updates

**Download mechanism:**

```
1. GET /repos/{owner}/{repo}/git/trees/{branch}?recursive=1
   → Returns complete file tree with SHA hashes for all files

2. Load local asset-version.toml
   → Contains last downloaded commit + file SHAs

3. Compare local vs remote SHAs:
   - Skip files with matching SHA (already up to date)
   - Download only changed/new files via Contents API
   - Track removed files (exist locally but not in remote tree)

4. Update asset-version.toml with new commit + all file SHAs
```

**Asset version tracking file:**

Location: `~/.config/start/asset-version.toml`

Format:
```toml
# Asset version tracking - managed by 'start update'
# Last updated: 2025-01-06T10:30:00Z

commit = "abc123def456"
timestamp = "2025-01-06T10:30:00Z"
repository = "github.com/grantcarthew/start"
branch = "main"

[files]
"agents/claude.toml" = "a1b2c3d4e5f6..."
"agents/gemini.toml" = "e5f6g7h8i9j0..."
"roles/code-reviewer.md" = "i9j0k1l2m3n4..."
```

**API endpoints used:**

1. **Tree API** (discovery + SHAs):
   - `GET /repos/{owner}/{repo}/git/trees/{branch}?recursive=1`
   - Returns: Complete file tree with SHA-256 hashes
   - Single API call gets entire repository structure

2. **Contents API** (download):
   - `GET /repos/{owner}/{repo}/contents/{path}?ref={branch}`
   - Downloads individual files
   - One call per changed/new file

**Rate limiting strategy:**

- **Anonymous:** 60 requests/hour
- **Authenticated:** 5000 requests/hour via `GH_TOKEN` env var
- **Check before download:** Query `/rate_limit` endpoint
- **Abort if insufficient:** Error with reset time if < 50 requests remaining
- **Smart caching:** SHA comparison reduces API calls dramatically

**Incremental update example:**

First update (cold cache):
```
- Tree API: 1 call
- 28 asset files: 28 calls
- Total: 29 API calls
```

Subsequent update (3 files changed):
```
- Tree API: 1 call
- 3 changed files: 3 calls
- 25 unchanged: 0 calls (skipped via SHA match)
- Total: 4 API calls
```

**Benefits:**

- ✅ **Automatic discovery** - No manifest file to maintain
- ✅ **Incremental updates** - Only downloads changed files
- ✅ **Integrity verification** - SHA comparison validates files
- ✅ **Extensible** - New asset types discovered automatically
- ✅ **Efficient** - Caching minimizes API usage
- ✅ **No external dependencies** - Pure GitHub API, no git/tar needed

**Trade-offs accepted:**

- ❌ Multiple API calls required (mitigated by caching)
- ❌ Rate limiting for anonymous users (GH_TOKEN recommended)
- ❌ First update downloads all files (one-time cost)

**Alternatives considered:**

- **Release archives:** Single download but requires tarball extraction
- **Manifest file:** Requires manual maintenance, can drift out of sync
- **Git clone:** Heavyweight, downloads entire repo history

**Rationale:**

SHA-based caching provides best balance of:
- Developer ergonomics (no manifest to maintain)
- User experience (fast incremental updates)
- Resource efficiency (minimal API calls after first download)
- Implementation simplicity (no external dependencies)

---

### DR-015: Atomic Update Mechanism (2025-01-06)

**Decision:** Use SHA-filtered incremental downloads with batch atomic install and rollback capability

**Update flow:**

```
1. Fetch remote tree (GitHub API - 1 call)
2. Load local asset-version.toml
3. Compare SHAs → Identify changed files only
4. Download changed files to temp directory (N API calls)
5. Batch install with rollback safety
6. Update asset-version.toml
7. Cleanup temp and backups
```

**Atomic install mechanism:**

```go
Phase 1: Backup
  - For each changed file being replaced
  - Rename: file → file.backup

Phase 2: Install
  - Move from temp to assets/
  - If any fail → rollback all from .backup

Phase 3: Commit
  - Update asset-version.toml with new SHAs
  - If fails → rollback all from .backup

Phase 4: Cleanup
  - Remove all .backup files
  - Remove temp directory
```

**Failure handling:**

| Failure Point | Result | Recovery |
|---|---|---|
| Download fails | Temp cleaned up, assets/ untouched | None needed |
| Backup fails | Abort, no changes made | None needed |
| Install fails | Restore from .backup files | Automatic rollback |
| Version file write fails | Restore from .backup files | Automatic rollback |
| Process killed | Orphaned .backup files | Auto-cleanup next run |

**Disk space:**

- Temp directory: Only changed files (typically 3-5 files, ~20 KB)
- Backup files: Only files being replaced (same size as changed)
- Total overhead: 2x size of changed files (not entire asset library)

**Benefits:**

- ✅ **Failure-safe:** All-or-nothing installation
- ✅ **Efficient:** Only downloads changed files (SHA filtering)
- ✅ **Recoverable:** Automatic rollback on any error
- ✅ **Minimal overhead:** Backups only for changed files
- ✅ **Simple:** No transaction log parsing needed

**Example (incremental update):**

```
Remote has 28 files
Local has 25 matching SHAs, 3 different

Download: 3 files to temp (3 API calls)
Backup: 2 existing files (.backup)
Install: 3 files from temp
Update: asset-version.toml
Cleanup: 2 .backup files, temp directory

Result: 4 total API calls, minimal disk usage
```

**Alternatives considered:**

- **Full temp directory copy:** Wasteful (downloads all 28 files every time)
- **Transaction log:** Complex, error-prone, hard to recover
- **Checkpointed download:** Leaves asset directory in inconsistent state

**Rationale:**

Combining SHA filtering (DR-014) with batch atomic install provides:
- Best API efficiency (only changed files)
- Best safety (rollback on failure)
- Simplest implementation (no complex state tracking)
- Clear failure recovery (restore from .backup)

---

### DR-016: Asset Discovery and Directory Structure (2025-01-06)

**Decision:** No central asset discovery system; each feature checks its own directory with graceful fallbacks

**Asset directory structure:**

```
~/.config/start/
├── config.toml
├── asset-version.toml
└── assets/
    ├── agents/              # Agent configuration templates
    │   ├── claude.toml
    │   ├── gemini.toml
    │   └── aichat.toml
    ├── roles/               # System prompt markdown files
    │   ├── code-reviewer.md
    │   ├── doc-reviewer.md
    │   └── security-reviewer.md
    ├── tasks/               # Default task definitions
    │   ├── code-review.toml
    │   ├── git-diff-review.toml
    │   └── doc-review.toml
    └── examples/            # Example configuration files
        ├── global-config.toml
        └── local-config.toml
```

**Usage pattern per asset type:**

**Agents (`assets/agents/`):**
- Accessed by: `start config agent add`
- Usage: Load as templates, show selection menu
- Fallback: Manual agent entry if directory empty/missing

**Roles (`assets/roles/`):**
- Accessed by: User config references: `file = "~/.config/start/assets/roles/code-reviewer.md"`
- Usage: Read file contents when loading system prompt
- Fallback: Standard file-not-found error

**Tasks (`assets/tasks/`):**
- Accessed by: `start task` command
- Usage: Parse and merge with user-defined tasks from config
- Fallback: Empty list if directory missing (use only user tasks)

**Examples (`assets/examples/`):**
- Accessed by: Never automatically (reference only)
- Usage: Users manually view/copy sections
- Fallback: N/A (optional resource)

**No discovery system needed:**

Each command handles its own asset directory:

```go
// Constants for asset directories (no magic strings)
const (
    AssetDirAgents   = "agents"
    AssetDirRoles    = "roles"
    AssetDirTasks    = "tasks"
    AssetDirExamples = "examples"
)

// Each feature loads independently
func LoadAgentTemplates() []AgentTemplate {
    dir := filepath.Join(assetDir, AssetDirAgents)
    files, err := filepath.Glob(filepath.Join(dir, "*.toml"))
    if err != nil || len(files) == 0 {
        return []AgentTemplate{} // Graceful fallback
    }
    return parseTemplates(files)
}
```

**Benefits:**

- ✅ **Simple:** No central registry or manifest
- ✅ **Decoupled:** Each feature self-contained
- ✅ **Graceful:** Missing directories don't break functionality
- ✅ **Maintainable:** Constants prevent magic strings
- ✅ **Extensible:** Add new directory + code that uses it

**Alternatives considered:**

- **Manifest file:** Over-engineering, requires maintenance
- **Dynamic discovery:** No way to know what to do with unknown types
- **Hardcoded registry:** Unnecessary coupling

**Rationale:**

Asset types are few and stable. Each type requires specific behavior (templates vs files vs merged config). Simple directory checks with fallbacks provide clean implementation without unnecessary abstraction.

---

### DR-017: CLI Command Reorganization (2025-01-06)

**Decision:** Reorganize commands by purpose - configuration management under `start config`, execution at top level

**Problem:**

Original structure was inconsistent:
- `start agent add` - Configuration management
- `start task code-review` - Execution

Different purposes, similar command structure = confusing.

**New structure:**

**Execution commands (top-level):**
```bash
start                        # Launch agent with all context
start prompt [text]          # Launch with required context + custom prompt
start task <name> [inst]     # Run predefined task
```

**Configuration management:**
```bash
start config show            # View merged configuration
start config edit [scope]    # Edit config file
start config path            # Show config file paths
start config validate        # Validate configuration

start config agent <sub>     # Manage agents (moved from start agent)
start config context <sub>   # Manage contexts (new)
start config task <sub>      # Manage tasks (new)
start config role <sub>      # Manage system prompts (new)
```

**Utility commands:**
```bash
start init [scope]           # Initialize configuration
start doctor                 # Diagnose installation
start update                 # Update asset library
```

**Configuration subcommands:**

All follow consistent pattern:

```bash
start config agent list
start config agent add
start config agent edit [name]
start config agent remove [name]
start config agent test <name>
start config agent default [name]

start config context list
start config context add
start config context edit [name]
start config context remove [name]
start config context test <name>

start config task list
start config task add
start config task edit [name]
start config task remove [name]

start config role list
start config role add
start config role edit [name]
start config role remove [name]
```

**Benefits:**

- ✅ **Clear mental model:** `start config X` = managing config, `start X` = executing
- ✅ **Consistent:** All configuration under one command
- ✅ **Discoverable:** Easy to find: "How do I manage X? → start config X"
- ✅ **Extensible:** New config sections fit the pattern naturally
- ✅ **No ambiguity:** Command purpose clear from structure

**Breaking change:**
- Design phase only - no existing users
- `start agent` → `start config agent`

**Rationale:**

Clear separation between configuration (managing settings) and execution (running the tool) provides better mental model and allows consistent expansion of configuration management without top-level command pollution.

---

## Pending Decisions

Design phase ongoing. Implementation decisions being finalized.

---

## References

- [Vision](./vision.md)
- [Design Thoughts](./thoughts.md)
