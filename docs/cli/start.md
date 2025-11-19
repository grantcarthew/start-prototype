# start

## Name

start - Launch AI agent with important context

## Synopsis

```bash
start [flags]
```

## Description

Launches an AI agent with automatically detected project context. Reads ALL configured context documents (both required and optional), builds an intelligent initial prompt, and delegates to the configured AI agent tool (claude, gemini, etc.).

**Context document behavior:**

- Includes ALL context documents (required + optional)
- Missing files generate warnings and are skipped
- Use `start prompt` for required documents only

This is the primary command for launching an AI session with full context. For custom prompts with minimal context, use `start prompt` subcommand. For predefined workflows, use `start task` subcommand. To discover and install new agents, roles, tasks, and contexts from the catalog, use `start assets` subcommand.

Note: _You can always launch your agent directly for full custom controls if the role and context documents are not needed for one-off use._

## Global Flags

These flags work on all `start` commands.

**--agent** _name_, **-a** _name_
: Which agent to use. Overrides default agent from config. Resolution order:

1. Exact match: local → global → cache → GitHub (lazy fetch)
2. Prefix match: local → global → cache → GitHub (short-circuit at first source with matches)
   - Single match → use it
   - Multiple matches → interactive selection (TTY) or error (non-TTY)

```bash
start --agent anthropic  # Exact match
start -a anth            # Prefix match (if unambiguous)
start --agent a          # Ambiguous: interactive picker or error
```

See [DR-038](../design/design-records/dr-038-flag-value-resolution.md) for full resolution algorithm.

**--role** _name_, **-r** _name_
: Which role to use for the system prompt. Overrides default role from config. Resolution order:

1. Exact match: local → global → cache → GitHub (lazy fetch)
2. Prefix match: local → global → cache → GitHub (short-circuit at first source with matches)
   - Single match → use it
   - Multiple matches → interactive selection (TTY) or error (non-TTY)

```bash
start --role go-expert       # Exact match
start -r go                  # Prefix match (if unambiguous)
start --role code            # Ambiguous: interactive picker or error
```

See [DR-038](../design/design-records/dr-038-flag-value-resolution.md) for full resolution algorithm.

**--model** _name_, **-m** _name_
: Model to use (from agent configuration). Resolution order:

1. Exact match on any configured model name → use it
2. Prefix match (short-circuit at first match) → use it
   - Single match → use it
   - Multiple matches → interactive selection (TTY) or error (non-TTY)
3. No match → pass string to agent as-is (passthrough, agent errors if invalid)

```bash
start --model claude-sonnet-4           # Exact match
start -m claude                         # Prefix match (if unambiguous)
start --model gpt-5-experimental        # No match, passthrough to agent
```

See [DR-038](../design/design-records/dr-038-flag-value-resolution.md) for full resolution algorithm.

**--directory** _path_, **-d** _path_
: Working directory for context detection. Relative paths in config resolve to this directory. Default: current directory (pwd).

```bash
start --directory ~/my-project
start -d ~/my-project
```

**--quiet**, **-q**
: Quiet mode. No output, launches agent directly. Use when you don't want to see context summary.

**--verbose**
: Verbose output. Shows config resolution, file detection details, full paths, and context being sent.

**--debug**
: Debug mode. Shows everything: config merging, placeholder resolution, command construction, environment variables.

**--asset-download[=bool]**
: Enable or disable downloading assets from the GitHub catalog on-demand. Defaults to `true`. Use `--asset-download=false` to prevent network requests for missing assets.

**--local**
: When downloading assets from the catalog (roles, agents, contexts), add them to local config (`./.start/`) instead of global config (`~/.config/start/`). Only applies when an asset is downloaded; has no effect if asset already exists in config or cache.

**--help**, **-h**
: Show help text.

**--version**, **-v**
: Show version information.

## Behavior

### Execution Flow

1. Load global config (`~/.config/start/`)
2. Load local config (`./.start/`) if exists
3. Merge configs (local overrides global)
4. Detect ALL context documents (check which files exist)
   - Includes both `required = true` and `required = false` documents
   - Missing files generate warnings and are skipped
   - Order determined by config definition order (see below)
5. Build prompt from document prompts
6. Resolve placeholders in agent command template
7. Display context summary (unless `--quiet`)
8. Execute agent command

### Document Order

Context documents appear in the prompt in the **order they are defined in the config file**.

**Example config:**

```toml
[contexts.environment]  # First
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."

[contexts.index]        # Second
file = "~/reference/INDEX.csv"
prompt = "Read {file} for documentation index."

[contexts.agents]       # Third
file = "./AGENTS.md"
prompt = "Read {file} for repository overview."

[contexts.project]      # Fourth
file = "./PROJECT.md"
prompt = "Read {file} for current project status."
```

**Resulting prompt order:**

```
Read ~/reference/ENVIRONMENT.md for environment context.
Read ~/reference/INDEX.csv for documentation index.
Read ./AGENTS.md for repository overview.
Read ./PROJECT.md for current project status.
```

**Controlling order:**

- Rearrange document definitions in config file
- TOML preserves declaration order within sections
- No alphabetical or other automatic sorting

### Model Flag Resolution

**When `--model` provided:**

1. Check for exact match in selected agent's models → use it
2. Check for prefix match in selected agent's models → use first match (by config order)
3. No match → pass string directly to agent (agent handles validation)

**Examples:**

Agent config has models: `sonnet`, `sonnet-new`, `haiku`

- `--model sonnet` → exact match, uses `sonnet`
- `--model s` → prefix match, uses `sonnet` (first in config)
- `--model haiku` → exact match, uses `haiku`
- `--model xyz-model` → no match, passes "xyz-model" to agent

**When no `--model` flag:**

1. Use agent's `default_model` from config
2. If no default_model, use first model in agent's models table (config order)

### Verbosity Levels

**--quiet** (minimal):

```
(no output - launches agent directly)
```

**Normal** (default):

```
Starting AI Agent
===============================================================================================
Agent: claude (model: claude-3-7-sonnet-20250219)

Context documents:
  ✓ environment     ~/reference/ENVIRONMENT.md
  ✓ index           ~/reference/INDEX.csv
  ✓ agents          ./AGENTS.md
  ⚠ project         ./PROJECT.md (not found, skipped)

Role: code-reviewer (from ~/.config/start/roles/code-reviewer.md)

Executing command...
❯ claude --model claude-3-7-sonnet-20250219 --append-system-prompt '...' '2025-11-03...'
```

**--verbose** (detailed):

```
Loading configuration...
  Global: ~/.config/start/ (5 files)
  Local:  ./.start/ (found, 2 files)
  Merged: 3 sections

Resolving agent: claude (default)
  Command template: claude --model {model} --append-system-prompt '{role}' '{prompt}'
  Model name: sonnet → claude-3-7-sonnet-20250219

Detecting context documents (working directory: /Users/gcarthew/Projects/my-app):
  environment: ~/reference/ENVIRONMENT.md → /Users/gcarthew/reference/ENVIRONMENT.md (exists)
  index: ~/reference/INDEX.csv → /Users/gcarthew/reference/INDEX.csv (exists)
  agents: ./AGENTS.md → /Users/gcarthew/Projects/my-app/AGENTS.md (exists)
  project: ./PROJECT.md → /Users/gcarthew/Projects/my-app/PROJECT.md (not found, skipped)

Loading role:
  Role: code-reviewer
  Source: ~/.config/start/roles/code-reviewer.md → /Users/gcarthew/.config/start/roles/code-reviewer.md
  Size: 1.2 KB

Building prompt:
  Document prompts: 3 documents
  Custom prompt: (none)
  Final prompt size: 520 characters

Starting AI Agent
===============================================================================================
Agent: claude (model: claude-3-7-sonnet-20250219)

Context documents:
  ✓ environment     ~/reference/ENVIRONMENT.md
  ✓ index           ~/reference/INDEX.csv
  ✓ agents          ./AGENTS.md
  ⚠ project         ./PROJECT.md (not found, skipped)

Role: code-reviewer (from ~/.config/start/roles/code-reviewer.md)

Executing command...
❯ claude --model claude-3-7-sonnet-20250219 --append-system-prompt '...' '2025-11-03...'
```

**--debug** (everything):

```
[DEBUG] Config loader initialized
[DEBUG] Reading global config: ~/.config/start/
[DEBUG]   config.toml: 45 lines (settings)
[DEBUG]   agents.toml: 67 lines (2 agents)
[DEBUG]   roles.toml: 34 lines (1 role)
[DEBUG]   contexts.toml: 28 lines (2 contexts)
[DEBUG]   tasks.toml: 71 lines (4 tasks)
[DEBUG] Reading local config: ./.start/
[DEBUG]   contexts.toml: 12 lines (1 context)
[DEBUG] Merging configs...
[DEBUG]   context.agents: "./AGENTS.md" (overridden by local)
[DEBUG]   ... (other merge details)

[DEBUG] Agent resolution:
[DEBUG]   Default agent: claude
[DEBUG]   Command template: claude --model {model} --append-system-prompt '{role}' '{prompt}'

[DEBUG] Model resolution:
[DEBUG]   Default model name: sonnet
[DEBUG]   Resolved model: claude-3-7-sonnet-20250219

[DEBUG] Context document detection:
[DEBUG]   environment: ~/reference/ENVIRONMENT.md → /Users/gcarthew/reference/ENVIRONMENT.md (exists)
[DEBUG]   index: ~/reference/INDEX.csv → /Users/gcarthew/reference/INDEX.csv (exists)
[DEBUG]   agents: ./AGENTS.md → /Users/gcarthew/Projects/my-app/AGENTS.md (exists)
[DEBUG]   project: ./PROJECT.md → /Users/gcarthew/Projects/my-app/PROJECT.md (not found)

[DEBUG] Placeholder resolution:
[DEBUG]   {model} → "claude-3-7-sonnet-20250219"
[DEBUG]   {role} → "[1247 chars from code-reviewer.md]"
[DEBUG]   {prompt} → "[448 chars]"
[DEBUG]   {date} → "2025-11-03T16:30:45+10:00"

Starting AI Agent
===============================================================================================
Agent: claude (model: claude-3-7-sonnet-20250219)

Context documents:
  ✓ environment     ~/reference/ENVIRONMENT.md
  ✓ index           ~/reference/INDEX.csv
  ✓ agents          ./AGENTS.md
  ⚠ project         ./PROJECT.md (not found, skipped)

Role: code-reviewer (from ~/.config/start/roles/code-reviewer.md)

Executing command...
❯ claude --model claude-3-7-sonnet-20250219 --append-system-prompt '...' '2025-11-03...'

[DEBUG] Executing command in shell
[DEBUG] Working directory: /Users/gcarthew/Projects/my-app
```

**Note:** Debug output shows detailed information for each step including config merging, path resolution, and placeholder substitution. Useful for troubleshooting configuration issues. Documents appear in config definition order.

### Lazy-Loading of Assets

If an asset (like an agent, role, or task) is requested but not found in the local or global configuration, the CLI will attempt to find it in the GitHub asset catalog. If found, it will be downloaded, cached, and added to your global configuration for future use. This allows for seamless, on-demand use of the entire asset library.

This behavior can be controlled with the `--asset-download` flag.

## Examples

### Basic Usage

Launch with default configuration:

```bash
start
```

### Agent Selection

Use specific agent:

```bash
start --agent gemini
start --agent opencode
```

### Model Selection

Use configured model names:

```bash
start --model haiku
start --model sonnet
start --model opus
```

Use prefix matching:

```bash
start --model s      # Matches 'sonnet' if unambiguous
start --model h      # Matches 'haiku' if unambiguous
```

Use any model string (passthrough):

```bash
start --model claude-3-5-haiku-20241022
start --model gemini-2.0-flash-exp
```

Combined with agent:

```bash
start --agent claude --model sonnet
```

### Directory Override

Work from different directory:

```bash
start --directory ~/my-project
start -d ~/projects/work/api-server
```

Useful when running from outside project:

```bash
cd ~
start --directory ~/my-project
```

### Verbosity Control

Quiet mode (no output):

```bash
start --quiet
start -q
```

Verbose mode:

```bash
start --verbose
```

Debug mode:

```bash
start --debug
```

### Combined Examples

Full power:

```bash
start --agent claude --model opus --directory ~/my-project --verbose
```

Quick and quiet:

```bash
start --agent gemini --model flash --quiet
```

## Output

See **Verbosity Levels** section above for detailed output examples.

## Exit Codes

**0** - Success (agent launched successfully)

**1** - Configuration error

- Config file syntax error
- Invalid TOML
- Missing required fields

**2** - Agent error

- Agent not found in config
- Agent command template invalid
- Model tier not configured

**3** - File error

- Working directory doesn't exist
- Config file permissions error

**4** - Runtime error

- Agent tool not installed
- Agent command failed to execute

## Environment

**EDITOR**
: Used by `start config edit` to open config files. Not used by root command.

**HOME**
: Used to resolve `~` in file paths.

**PWD**
: Default working directory if `--directory` not specified.

## Files

**~/.config/start/**
: Global configuration directory containing config.toml (settings), agents.toml, roles.toml, contexts.toml, and tasks.toml

**./.start/**
: Local (project-specific) configuration directory with same structure

## Error Handling

### Missing Config

If no config files exist:

```
Error: No configuration found.

Run 'start init' to create initial configuration.
```

Exit code: 1

### Invalid Agent

If specified agent not in config:

```
Error: Agent 'foo' not found in configuration.

Available agents:
  - claude
  - gemini
  - opencode

Use 'start config agent list' to see details.
```

Exit code: 2

### Invalid Model (from agent)

If model string is invalid (determined by agent, not CLI):

```
(Agent-specific error output)
```

Exit code: 4 (agent execution error)

### Agent Tool Not Found

If agent binary is not found:

```
Error: Agent tool 'gemini' not found.

Command attempted: gemini --model gemini-2.0-flash-exp '...'

Make sure 'gemini' is installed and available.
```

Exit code: 4

### Invalid Working Directory

If `--directory` path doesn't exist:

```
Error: Working directory not found: ~/nonexistent

Check the path and try again.
```

Exit code: 3

## Notes

### Edge Cases

**No context documents configured:**

If the `[contexts]` section is empty or missing:

```
Context documents: (none configured)
```

The agent launches with no context document instructions - only system prompt (if configured).

**No role configured:**

If no roles are configured or the role file doesn't exist:

```
Role: (none)
```

The agent launches without a system prompt. This is valid - not all agents require system prompts.

**No configuration files:**

If neither global (`~/.config/start/`) nor local (`./.start/`) exists:

```
Error: No configuration found.

Run 'start init' to create initial configuration.
```

**Local config only:**

If local `./.start/` exists but no global config, this is **valid** and the tool works normally. No error is shown. Per DR-004, agents can be defined in both global and local configs.

Use case: Team configurations where `./.start/` contains all necessary config (agents, roles, tasks, contexts) and is committed to version control. No global config required.

**All context documents missing:**

If all configured documents don't exist:

```
Context documents:
  ⚠ environment    ~/reference/ENVIRONMENT.md (not found, skipped)
  ⚠ index          ~/reference/INDEX.csv (not found, skipped)
  ⚠ agents         ./AGENTS.md (not found, skipped)
  ⚠ project        ./PROJECT.md (not found, skipped)
```

The agent launches with no context. This is valid - useful for general AI sessions.

### Context Document Configuration

Documents can be marked as `required` to control which commands include them:

```toml
[contexts.environment]
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
required = true    # Always included

[contexts.index]
file = "~/reference/INDEX.csv"
prompt = "Read {file} for documentation index."
required = true    # Always included

[contexts.agents]
file = "./AGENTS.md"
prompt = "Read {file} for repository context."
required = true    # Always included

[contexts.project]
file = "./PROJECT.md"
prompt = "Read {file}. Respond with summary."
required = false   # Optional - included by start, excluded by start prompt
```

**Behavior by command:**

- `start` (root) → Includes ALL documents (required + optional)
- `start prompt` → Includes ONLY required documents
- `start task` → Auto-includes all contexts where `required = true`

**Default value:**

- If `required` field is omitted, defaults to `false` (optional)

**Document order:**

- Documents appear in prompt in the order defined in config file
- See "Document Order" section above for details

### Model Override Behavior

The `--model` flag overrides the default model. Each agent has its own configured models.

**Example configured models:**

- Claude: `haiku`, `sonnet`, `opus`
- Gemini: `flash`, `pro-exp`
- Other agents: user-defined

Model resolution is agent-specific:

```bash
start --agent claude --model opus     # Matches 'opus' in claude config
start --agent gemini --model opus     # No match, passes "opus" to gemini (likely errors)
start --agent gemini --model flash    # Matches 'flash' in gemini config
```

Any unmatched string is passed to the agent:

```bash
start --agent claude --model claude-opus-4-20250514  # Passthrough to agent
```

### Tilde Expansion

Paths with `~` are expanded to user's home directory. This happens for:

- Config paths
- Context document paths
- Role file paths
- Working directory

### Relative Path Resolution

Relative paths in config (e.g., `./AGENTS.md`) resolve relative to:

- Working directory (default: `pwd`)
- `--directory` path if specified

Absolute paths and `~` paths resolve independently of working directory.

## See Also

- start-assets(1) - Manage catalog assets
- start-prompt(1) - Launch with custom prompt
- start-task(1) - Run predefined tasks
- start-init(1) - Initialize configuration
- start-config-agent(1) - Manage agents
- start-config(1) - Manage configuration
