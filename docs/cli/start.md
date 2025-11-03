# start

## Name

start - Launch AI agent with project context

## Synopsis

```bash
start [prompt] [flags]
```

## Description

Launches an AI agent with automatically detected project context. Reads configured context documents, builds an intelligent initial prompt, and delegates to the configured AI agent tool (claude, gemini, etc.).

When run without arguments, uses default context documents and launches interactive session. When run with a custom prompt argument, includes that prompt along with context.

## Arguments

**prompt** (optional)
: Custom prompt to send to the agent. If provided, this becomes the initial prompt along with context document instructions.

```bash
start "analyze this codebase for security vulnerabilities"
```

## Global Flags

These flags work on all `start` commands.

**--agent** *name*
: Which agent to use. Overrides default agent from config.

```bash
start --agent gemini
```

**--model** *tier|name*
: Model to use. Accepts either:
- Tier names: `fast`, `mid`, `pro` (uses agent's configured model for that tier)
- Full model names: `claude-3-5-haiku-20241022`, `gemini-2.0-flash-exp`

```bash
start --model fast                      # Use fast tier
start --model claude-3-5-haiku-20241022 # Use specific model
```

**--directory** *path*, **-d** *path*
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

### Default Execution (No Arguments)

1. Load global config (`~/.config/start/config.toml`)
2. Load local config (`./.start/config.toml`) if exists
3. Merge configs (local overrides global)
4. Detect context documents (check which files exist)
5. Build prompt from document suffixes
6. Resolve placeholders in agent command template
7. Display context summary (unless `--quiet`)
8. Execute agent command

### Custom Prompt Execution (With Argument)

1-6. Same as above
7. Prepend custom prompt to document instructions
8. Display context summary (unless `--quiet`)
9. Execute agent command

### Model Flag Resolution

When `--model fast` (tier name):
1. Look up agent's `models.fast` value from config
2. Use that model name
3. Error if tier not configured for agent

When `--model claude-3-5-haiku-20241022` (full name):
1. Use exact model name provided
2. Bypass tier configuration
3. Agent must support this model

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
  ✓ agents          ./AGENTS.md
  ✗ project         ./PROJECT.md (not found)

System prompt: ./ROLE.md

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
  Command template: claude --model {model} --append-system-prompt '{system_prompt}' '{prompt}'
  Model tier: mid → claude-3-7-sonnet-20250219

Detecting context documents (working directory: /Users/gcarthew/Projects/my-app):
  environment: ~/reference/ENVIRONMENT.md → /Users/gcarthew/reference/ENVIRONMENT.md (exists)
  agents: ./AGENTS.md → /Users/gcarthew/Projects/my-app/AGENTS.md (exists)
  project: ./PROJECT.md → /Users/gcarthew/Projects/my-app/PROJECT.md (not found, skipped)

Loading system prompt:
  Path: ./ROLE.md → /Users/gcarthew/Projects/my-app/ROLE.md
  Size: 1.2 KB

Building prompt:
  Document suffixes: 2 documents
  Custom prompt: (none)
  Final prompt size: 450 characters

Starting AI Agent
===============================================================================================
Agent: claude (model: claude-3-7-sonnet-20250219)

Context documents:
  ✓ environment     ~/reference/ENVIRONMENT.md
  ✓ agents          ./AGENTS.md
  ✗ project         ./PROJECT.md (not found)

System prompt: ./ROLE.md

Executing command...
❯ claude --model claude-3-7-sonnet-20250219 --append-system-prompt '...' '2025-11-03...'
```

**--debug** (everything):
```
[DEBUG] Config loader initialized
[DEBUG] Reading global config: ~/.config/start/config.toml
[DEBUG] Global config loaded: 245 lines, 3 sections (settings, agents, context)
[DEBUG] Reading local config: ./.start/config.toml
[DEBUG] Local config loaded: 12 lines, 1 section (context)
[DEBUG] Merging configs...
[DEBUG]   settings.default_agent: "claude" (from global)
[DEBUG]   agents.claude: (from global)
[DEBUG]   context.documents.environment: (from global)
[DEBUG]   context.documents.agents: "AGENTS.md" → "./AGENTS.md" (overridden by local)
[DEBUG]   context.documents.project: (from global)

[DEBUG] Agent resolution:
[DEBUG]   Requested: (none, using default)
[DEBUG]   Default agent: claude
[DEBUG]   Command template: claude --model {model} --append-system-prompt '{system_prompt}' '{prompt}'
[DEBUG]   Environment variables: (none)

[DEBUG] Model resolution:
[DEBUG]   Requested: (none, using default)
[DEBUG]   Default tier: mid
[DEBUG]   Agent tiers: fast=claude-3-5-haiku-20241022, mid=claude-3-7-sonnet-20250219, pro=claude-opus-4-20250514
[DEBUG]   Resolved model: claude-3-7-sonnet-20250219

[DEBUG] Working directory: /Users/gcarthew/Projects/my-app

[DEBUG] Context document detection:
[DEBUG]   environment:
[DEBUG]     Configured path: ~/reference/ENVIRONMENT.md
[DEBUG]     Expanded path: /Users/gcarthew/reference/ENVIRONMENT.md
[DEBUG]     Exists: true
[DEBUG]     Suffix: "Read {file} for environment context."
[DEBUG]   agents:
[DEBUG]     Configured path: ./AGENTS.md
[DEBUG]     Resolved path: /Users/gcarthew/Projects/my-app/AGENTS.md
[DEBUG]     Exists: true
[DEBUG]     Suffix: "Read {file} for repository overview."
[DEBUG]   project:
[DEBUG]     Configured path: ./PROJECT.md
[DEBUG]     Resolved path: /Users/gcarthew/Projects/my-app/PROJECT.md
[DEBUG]     Exists: false (skipped)

[DEBUG] System prompt resolution:
[DEBUG]   Configured path: ./ROLE.md
[DEBUG]   Resolved path: /Users/gcarthew/Projects/my-app/ROLE.md
[DEBUG]   Exists: true
[DEBUG]   Size: 1247 bytes

[DEBUG] Prompt construction:
[DEBUG]   Custom prompt: (none)
[DEBUG]   Document 1: "Read /Users/gcarthew/reference/ENVIRONMENT.md for environment context."
[DEBUG]   Document 2: "Read /Users/gcarthew/Projects/my-app/AGENTS.md for repository overview."
[DEBUG]   Final prompt: "Read /Users/gcarthew/reference/ENVIRONMENT.md for...[448 chars total]"

[DEBUG] Placeholder resolution:
[DEBUG]   {model} → "claude-3-7-sonnet-20250219"
[DEBUG]   {system_prompt} → "[1247 chars from ./ROLE.md]"
[DEBUG]   {prompt} → "[448 chars]"
[DEBUG]   {date} → "2025-11-03T16:30:45+10:00"

[DEBUG] Final command:
[DEBUG]   claude --model claude-3-7-sonnet-20250219 --append-system-prompt 'You are a senior...[truncated]' 'Read /Users/gcarthew/reference...[truncated]'

Starting AI Agent
===============================================================================================
Agent: claude (model: claude-3-7-sonnet-20250219)

Context documents:
  ✓ environment     ~/reference/ENVIRONMENT.md
  ✓ agents          ./AGENTS.md
  ✗ project         ./PROJECT.md (not found)

System prompt: ./ROLE.md

Executing command...
❯ claude --model claude-3-7-sonnet-20250219 --append-system-prompt '...' '2025-11-03...'

[DEBUG] Executing command in shell
[DEBUG] Working directory: /Users/gcarthew/Projects/my-app
[DEBUG] Environment inherited from parent
[DEBUG] Executing...
```

## Examples

### Basic Usage

Launch with default configuration:
```bash
start
```

Launch with custom prompt:
```bash
start "analyze this codebase for security vulnerabilities"
```

### Agent Selection

Use specific agent:
```bash
start --agent gemini
start --agent opencode
```

Custom prompt with specific agent:
```bash
start --agent gemini "review the API design"
```

### Model Selection

Use tier name:
```bash
start --model fast
start --model mid
start --model pro
```

Use full model name:
```bash
start --model claude-3-5-haiku-20241022
start --model gemini-2.0-flash-exp
```

Combined:
```bash
start --agent claude --model pro "comprehensive code review"
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
start --directory ~/my-project "what is the project status?"
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
start --agent gemini --model pro --directory ~/my-project --verbose "review architecture"
```

Quick and quiet:
```bash
start --agent claude --model fast --quiet
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

### Invalid Model Tier

If tier not configured for agent:
```
Error: Model tier 'pro' not configured for agent 'claude'.

Available tiers for claude:
  - fast: claude-3-5-haiku-20241022
  - mid: claude-3-7-sonnet-20250219

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

### Custom Prompt + Context Documents

When providing a custom prompt, context documents are still included. The final prompt structure:

```
{custom_prompt}

{document_suffix_1}
{document_suffix_2}
...
```

To launch with ONLY custom prompt (no context), use a task with empty `documents` array.

### Model Override Behavior

The `--model` flag overrides the default model tier but respects the agent. If you want to use a specific model with a specific agent:

```bash
start --agent claude --model pro
```

If the model name doesn't match agent's expected format, the agent may error. Use model names appropriate for the selected agent.

### Tilde Expansion

Paths with `~` are expanded to user's home directory. This happens for:
- Config paths
- Context document paths
- System prompt path
- Working directory

### Relative Path Resolution

Relative paths in config (e.g., `./AGENTS.md`) resolve relative to:
- Working directory (default: `pwd`)
- `--directory` path if specified

Absolute paths and `~` paths resolve independently of working directory.

## See Also

- start-init(1) - Initialize configuration
- start-task(1) - Run predefined tasks
- start-agent(1) - Manage agents
- start-config(1) - Manage configuration
