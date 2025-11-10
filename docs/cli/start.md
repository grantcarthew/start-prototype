# start

## Name

start - Launch AI agent with project context

## Synopsis

```bash
start [flags]
```

## Description

Launches an AI agent with automatically detected project context. Reads ALL configured context documents (both required and optional), builds an intelligent initial prompt, and delegates to the configured AI agent tool (claude, gemini, etc.).

**Context document behavior:**

- Includes ALL context documents (required + optional)
- Missing files are skipped with notification
- Use `start prompt` for required documents only

This is the primary command for launching an AI session with full context. For custom prompts with minimal context, use `start prompt` subcommand. For predefined workflows, use `start task` subcommand.

## Global Flags

These flags work on all `start` commands.

**--agent** _name_
: Which agent to use. Overrides default agent from config. If the agent is not found in the local or global configuration, it will be searched for in the GitHub asset catalog and can be lazy-loaded on first use.

```bash
start --agent gemini
```

**--role** _name_
: Which role to use for the system prompt. Overrides default role from config. If the role is not found in the local or global configuration, it will be searched for in the GitHub asset catalog and can be lazy-loaded on first use.

```bash
start --role security-auditor
start --role go-expert
```

**--model** _alias|name_
: Model to use. Accepts either:

- Model alias: User-defined aliases from agent's config (e.g., `sonnet`, `haiku`, `opus`)
- Full model name: Complete model identifier (e.g., `claude-3-5-haiku-20241022`, `gemini-2.0-flash-exp`)

```bash
start --model sonnet                    # Use alias from config
start --model claude-3-5-haiku-20241022 # Use specific model
```

**--directory** _path_, **-d** _path_
: Working directory for context detection. Relative paths in config resolve to this directory. Default: current directory (pwd).

```bash
start --directory ~/my-project
start -d ~/my-project
```

**--quiet**, **-q**
: Quiet mode. No output, launches agent directly. Use when you don't want to see context summary.

**--verbose**, **-v**
: Verbose output. Shows config resolution, file detection details, full paths, and context being sent.

**--debug**
: Debug mode. Shows everything: config merging, placeholder resolution, command construction, environment variables.

**--help**, **-h**
: Show help text.

**--version**
: Show version information.

## Behavior

### Execution Flow

1. Load global config (`~/.config/start/config.toml`)
2. Load local config (`./.start/config.toml`) if exists
3. Merge configs (local overrides global)
4. Detect ALL context documents (check which files exist)
   - Includes both `required = true` and `required = false` documents
   - Missing files are skipped (not errors)
   - Order determined by config definition order (see below)
5. Build prompt from document prompts
6. Resolve placeholders in agent command template
7. Display context summary (unless `--quiet`)
8. Execute agent command

### Document Order

Context documents appear in the prompt in the **order they are defined in the config file**.

**Example config:**

```toml
[context.environment]  # First
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."

[context.index]        # Second
file = "~/reference/INDEX.csv"
prompt = "Read {file} for documentation index."

[context.agents]       # Third
file = "./AGENTS.md"
prompt = "Read {file} for repository overview."

[context.project]      # Fourth
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

**When using alias** (`--model sonnet`):

1. Look up agent's `models.sonnet` value from config
2. Use the resolved model name
3. Error if alias not defined for agent

**When using full model name** (`--model claude-3-5-haiku-20241022`):

1. Use exact model name provided
2. Bypass alias resolution
3. Agent must support this model

**When no `--model` flag**:

1. Use agent's `default_model` alias from config
2. Error if no default_model configured

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
  ✗ project         ./PROJECT.md (not found)

Role: code-reviewer (from ~/.config/start/roles/code-reviewer.md)

Executing command...
❯ claude --model claude-3-7-sonnet-20250219 --append-system-prompt '...' '2025-11-03...'
```

**--verbose** (detailed):

```
Loading configuration...
  Global: ~/.config/start/config.toml
  Local:  ./.start/config.toml (found)
  Merged: 3 sections

Resolving agent: claude (default)
  Command template: claude --model {model} --append-system-prompt '{role}' '{prompt}'
  Model alias: sonnet → claude-3-7-sonnet-20250219

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
  ✗ project         ./PROJECT.md (not found)

Role: code-reviewer (from ~/.config/start/roles/code-reviewer.md)

Executing command...
❯ claude --model claude-3-7-sonnet-20250219 --append-system-prompt '...' '2025-11-03...'
```

**--debug** (everything):

```
[DEBUG] Config loader initialized
[DEBUG] Reading global config: ~/.config/start/config.toml
[DEBUG] Global config loaded: 245 lines, 3 sections
[DEBUG] Reading local config: ./.start/config.toml
[DEBUG] Local config loaded: 12 lines, 1 section
[DEBUG] Merging configs...
[DEBUG]   context.documents.agents: "./AGENTS.md" (overridden by local)
[DEBUG]   ... (other merge details)

[DEBUG] Agent resolution:
[DEBUG]   Default agent: claude
[DEBUG]   Command template: claude --model {model} --append-system-prompt '{role}' '{prompt}'

[DEBUG] Model resolution:
[DEBUG]   Default alias: sonnet
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
  ✗ project         ./PROJECT.md (not found)

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

Use model alias (from config):

```bash
start --model haiku
start --model sonnet
start --model opus
```

Use full model name:

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
start -v
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

**~/.config/start/config.toml**
: Global configuration file

**./.start/config.toml**
: Local (project-specific) configuration file

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

Use 'start agent list' to see details.
```

Exit code: 2

### Invalid Model Alias

If alias not configured for agent:

```
Error: Model alias 'pro' not configured for agent 'gemini'.

Available aliases for gemini:
  - flash: gemini-2.0-flash-exp
  - pro-exp: gemini-2.0-pro-exp

Update config or use full model name with --model.
```

Exit code: 2

### Agent Tool Not Found

If agent binary not in PATH:

```
Error: Agent tool 'gemini' not found.

Command attempted: gemini --model gemini-2.0-flash-exp '...'

Make sure 'gemini' is installed and in your PATH.
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

If the `[context.documents]` section is empty or missing:

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

If neither global (`~/.config/start/config.toml`) nor local (`./.start/config.toml`) exists:

```
Error: No configuration found.

Run 'start init' to create initial configuration.
```

**Local config only:**

If local `./.start/config.toml` exists but no global config:

```
Error: No global configuration found at ~/.config/start/config.toml

Run 'start init' to create global configuration, or ensure
agents are defined in local config.
```

Note: Per DR-004, agents can be defined in both global and local configs. If using only local config, at least one agent must be defined. For shared team configurations, defining agents in local `./.start/config.toml` allows the config to be committed to version control.

**All context documents missing:**

If all configured documents don't exist:

```
Context documents:
  ✗ environment    ~/reference/ENVIRONMENT.md (not found)
  ✗ index          ~/reference/INDEX.csv (not found)
  ✗ agents         ./AGENTS.md (not found)
  ✗ project        ./PROJECT.md (not found)
```

The agent launches with no context. This is valid - useful for general AI sessions.

### Context Document Configuration

Documents can be marked as `required` to control which commands include them:

```toml
[context.environment]
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
required = true    # Always included

[context.index]
file = "~/reference/INDEX.csv"
prompt = "Read {file} for documentation index."
required = true    # Always included

[context.agents]
file = "./AGENTS.md"
prompt = "Read {file} for repository context."
required = true    # Always included

[context.project]
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

The `--model` flag overrides the default model but respects the agent. Each agent has its own model aliases defined in config.

**Agent-specific aliases:**

- Claude: `haiku`, `sonnet`, `opus`
- Gemini: `flash`, `pro-exp`
- Other agents: user-defined

When using an alias, it must be defined for the selected agent:

```bash
start --agent claude --model opus     # ✓ Works if opus defined for claude
start --agent gemini --model opus     # ✗ Error if opus not defined for gemini
start --agent gemini --model flash    # ✓ Works if flash defined for gemini
```

When using full model names, the agent must support that model:

```bash
start --agent claude --model claude-opus-4-20250514  # Agent must support this
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

- start-prompt(1) - Launch with custom prompt
- start-task(1) - Run predefined tasks
- start-init(1) - Initialize configuration
- start-agent(1) - Manage agents
- start-config(1) - Manage configuration
