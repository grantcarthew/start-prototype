# start show

## Name

start show - Preview command execution or display resolved configuration content

## Synopsis

```bash
# Execution preview (shows what would execute, but doesn't execute)
start show [flags]
start show task <name> [instructions] [flags]
start show prompt <text> [flags]

# Content viewer (shows resolved content after UTD processing)
start show role [name] [flags]
start show context [name] [flags]
start show agent [name] [flags]
start show task [name] [flags]
```

## Description

The `start show` command has two distinct modes:

**1. Execution Preview Mode** - Preview what commands would execute without running them:

- `start show` - Preview `start` interactive session
- `start show task <name>` - Preview `start task <name>` execution
- `start show prompt <text>` - Preview `start prompt <text>` execution

Shows normal terminal output plus extra metadata (contexts loaded, file sizes, command to execute), with content truncated to 10 lines unless `--verbose` is used.

**2. Content Viewer Mode** - Display resolved content after UTD processing and config merging:

- `start show role [name]` - Show role content (after UTD processing)
- `start show context [name]` - Show context content (after UTD processing)
- `start show agent [name]` - Show agent effective configuration
- `start show task [name]` - Show task resolved prompt

Shows the final processed content that would be used by commands.

**For asset management**, see `start assets` commands:

- `start assets browse` - Open catalog in browser
- `start assets search` - Search catalog
- `start assets info` - Show asset details
- `start assets update` - Update cached assets
- `start assets clean` - Remove unused cache

**Key difference from `start config <type> show`:**

- `start config <type> show` - Shows raw TOML configuration structure
- `start show <type>` - Shows resolved/processed content (after UTD, placeholders, merging)

## Execution Preview Mode

### start show

Preview what the `start` command would execute without running the agent.

**Synopsis:**

```bash
start show [flags]
```

**Behavior:**

Displays:

- Which agent would be used
- Which role would be used
- Which contexts would be loaded (all: required + optional)
- File paths, sizes, existence checks
- Resolved role content (truncated to 10 lines)
- Final composed prompt (truncated to 10 lines)
- Exact agent command that would execute

Does NOT execute the agent.

**Output:**

```
Starting AI Agent (PREVIEW - NOT EXECUTING)
===============================================================================================
Agent: claude (model: claude-3-7-sonnet-20250219)
Role: code-reviewer (from global config)

Context documents (all):
  ✓ environment     ~/reference/ENVIRONMENT.md (2.3 KB, required)
  ✓ index           ~/reference/INDEX.csv (456 bytes, required)
  ✓ agents          ./AGENTS.md (1.2 KB, required)
  ✓ project         ./PROJECT.md (3.4 KB, optional)

Resolved role content (first 10 lines):
─────────────────────────────────────────────────
You are an expert code reviewer...
Focus on security, performance...
Check for edge cases...
Verify error handling...
Review test coverage...
Ensure documentation...
Look for code smells...
Consider maintainability...
Validate input handling...
Check resource cleanup...
... (347 more lines) - Use --verbose to see full content
─────────────────────────────────────────────────

Composed prompt (first 10 lines):
─────────────────────────────────────────────────
Read /Users/grant/reference/ENVIRONMENT.md for environment context.
Read /Users/grant/reference/INDEX.csv for documentation index.
Read /Users/grant/Projects/myapp/AGENTS.md for repository overview.
Read /Users/grant/Projects/myapp/PROJECT.md. Respond with summary.
... (125 more lines) - Use --verbose to see full content
─────────────────────────────────────────────────

Command that would execute:
❯ claude --model claude-3-7-sonnet-20250219 --append-system-prompt '...' 'Read...'

PREVIEW ONLY: Agent not executed
Use 'start' to execute, or 'start show --verbose' for full content
```

**With --verbose:**

```bash
start show --verbose
```

Shows full content (no truncation):

```
Starting AI Agent (PREVIEW - NOT EXECUTING)
===============================================================================================
Agent: claude (model: claude-3-7-sonnet-20250219)
Role: code-reviewer (from global config)

Context documents (all):
  ✓ environment     ~/reference/ENVIRONMENT.md (2.3 KB, required)
  ✓ index           ~/reference/INDEX.csv (456 bytes, required)
  ✓ agents          ./AGENTS.md (1.2 KB, required)
  ✓ project         ./PROJECT.md (3.4 KB, optional)

Resolved role content (full):
─────────────────────────────────────────────────
You are an expert code reviewer...
[... full 357 lines of role content ...]
─────────────────────────────────────────────────

Composed prompt (full):
─────────────────────────────────────────────────
Read /Users/grant/reference/ENVIRONMENT.md for environment context.
[... full 135 lines of prompt ...]
─────────────────────────────────────────────────

Command that would execute:
❯ claude --model claude-3-7-sonnet-20250219 --append-system-prompt '...' 'Read...'

PREVIEW ONLY: Agent not executed
```

### start show task

Preview what a task would execute without running it.

**Synopsis:**

```bash
start show task <name> [instructions] [flags]
```

**Arguments:**

**name** (required)
: Task name or alias to preview.

**instructions** (optional)
: Instructions to pass to the task (used in `{instructions}` placeholder).

**Behavior:**

Displays:

- Which task would be used
- Task's configured role (if specified)
- Which agent would be used
- Which contexts would be loaded (required only)
- Task's command execution (if configured)
- Task's resolved prompt template
- Final composed prompt
- Exact agent command that would execute

Does NOT execute the task or agent.

**Output:**

```
Starting Task: git-diff-review (PREVIEW - NOT EXECUTING)
===============================================================================================
Agent: claude (model: claude-3-7-sonnet-20250219)
Role: code-reviewer (from task configuration)
Task alias: gdr

Context documents (required only):
  ✓ environment     ~/reference/ENVIRONMENT.md (2.3 KB)
  ✓ index           ~/reference/INDEX.csv (456 bytes)
  ✓ agents          ./AGENTS.md (1.2 KB)

Task configuration:
  Command: git diff --staged
  Shell: bash
  Timeout: 10 seconds

Executing task command...
  ❯ git diff --staged
  ✓ Executed successfully (1.8 KB output)

Task prompt (first 10 lines):
─────────────────────────────────────────────────
Review changes:

## Instructions
focus on security

  ## Changes
  ```diff
  {command_output}
  ```

  ... (45 more lines) - Use --verbose to see full content─────────────────────────────────────────────────

Final composed prompt (first 10 lines):
─────────────────────────────────────────────────
Read /Users/grant/reference/ENVIRONMENT.md for environment context.
Read /Users/grant/reference/INDEX.csv for documentation index.
Read /Users/grant/Projects/myapp/AGENTS.md for repository overview.

Review changes:
... (52 more lines) - Use --verbose to see full content
─────────────────────────────────────────────────

Command that would execute:
❯ claude --model claude-3-7-sonnet-20250219 --append-system-prompt '...' 'Read...'

PREVIEW ONLY: Task not executed

```

**Error cases shown in preview:**

```

Starting Task: broken-task (PREVIEW - NOT EXECUTING)
===============================================================================================

Task configuration:
  Command: nonexistent-command
  Shell: bash
  Timeout: 10 seconds

Executing task command...
  ❯ nonexistent-command
  ✗ Command failed (exit code 127)
  Error: nonexistent-command: command not found

⚠ This task will fail when executed due to command errors

```

### start show prompt

Preview what a custom prompt would execute without running it.

**Synopsis:**

```bash
start show prompt <text> [flags]
```

**Arguments:**

**text** (required)
: Custom prompt text to preview.

**Behavior:**

Displays:

- Which agent would be used
- Which role would be used
- Which contexts would be loaded (required only)
- Custom prompt text
- Final composed prompt
- Exact agent command that would execute

Does NOT execute the agent.

**Output:**

```
Starting AI Agent with Custom Prompt (PREVIEW - NOT EXECUTING)
===============================================================================================
Agent: claude (model: claude-3-7-sonnet-20250219)
Role: code-reviewer (from global config)

Custom prompt: analyze security vulnerabilities

Context documents (required only):
  ✓ environment     ~/reference/ENVIRONMENT.md (2.3 KB)
  ✓ index           ~/reference/INDEX.csv (456 bytes)
  ✓ agents          ./AGENTS.md (1.2 KB)

Final composed prompt (first 10 lines):
─────────────────────────────────────────────────
Read /Users/grant/reference/ENVIRONMENT.md for environment context.
Read /Users/grant/reference/INDEX.csv for documentation index.
Read /Users/grant/Projects/myapp/AGENTS.md for repository overview.

analyze security vulnerabilities
... (8 more lines) - Use --verbose to see full content
─────────────────────────────────────────────────

Command that would execute:
❯ claude --model claude-3-7-sonnet-20250219 --append-system-prompt '...' 'Read...'

PREVIEW ONLY: Agent not executed
```

## Content Viewer Mode

### start show role

Display resolved role content after UTD processing.

**Synopsis:**

```bash
start show role              # Show default role
start show role <name>       # Show named role
start show role --scope global
start show role --scope local
```

**Arguments:**

**name** (optional)
: Role name to show. If omitted, shows the default role.

**Flags:**

**--scope** _scope_
: Show role from specific scope (`global` or `local`). If omitted, shows effective/merged role.

**Behavior:**

Shows the fully resolved role content after:

- UTD processing (file contents loaded, commands executed, placeholders replaced)
- Config merging (global + local)

This is the exact content that would be passed to the agent as the system prompt.

**Output (default role):**

```bash
start show role
```

```
Role: code-reviewer (effective)
═══════════════════════════════════════════════════════════
Source: global config
Type: File-based

You are an expert code reviewer with deep knowledge of software
engineering best practices, security vulnerabilities, and performance
optimization.

Focus on:
- Security vulnerabilities (OWASP Top 10, injection, XSS, etc.)
- Performance bottlenecks and optimization opportunities
- Code maintainability and readability
- Test coverage and edge cases
- Error handling and resource cleanup
- Documentation completeness

[... full resolved role content ...]
```

**Output (named role):**

```bash
start show role go-expert
```

```
Role: go-expert (effective)
═══════════════════════════════════════════════════════════
Source: global config
Type: File with command (UTD processed)

You are a Go programming language expert with expertise in:
- Go idioms and best practices
- Concurrency patterns (goroutines, channels, sync)
- Performance optimization
- Standard library usage
- Error handling patterns

Current Go environment:
go version go1.21.5 darwin/arm64

[... full resolved content ...]
```

**Output (with scope):**

```bash
start show role code-reviewer --scope global
```

```
Role: code-reviewer (global)
═══════════════════════════════════════════════════════════
Source: ~/.config/start/roles.toml
Type: File-based

[... global role content ...]
```

**No role configured:**

```
No default role configured.

Using first role in config: code-reviewer

Use 'start config role list' to see all roles.
Use 'start config role default <name>' to set default.
```

### start show context

Display resolved context content after UTD processing.

**Synopsis:**

```bash
start show context              # Show all contexts
start show context <name>       # Show named context
start show context --scope global
start show context --scope local
```

**Arguments:**

**name** (optional)
: Context name to show. If omitted, shows all contexts.

**Flags:**

**--scope** _scope_
: Show context from specific scope (`global` or `local`). If omitted, shows effective/merged contexts.

**Behavior:**

Shows the fully resolved context content after:

- UTD processing (file contents loaded, commands executed, placeholders replaced)
- Config merging (global + local)

This is the exact content that would be included in prompts.

**Output (all contexts):**

```bash
start show context
```

```
Contexts (effective - 4 total)
═══════════════════════════════════════════════════════════

environment (global, required)
─────────────────────────────────────────────────
Read /Users/grant/reference/ENVIRONMENT.md for environment context.
─────────────────────────────────────────────────

index (global, required)
─────────────────────────────────────────────────
Read /Users/grant/reference/INDEX.csv for documentation index.
─────────────────────────────────────────────────

agents (local, required)
─────────────────────────────────────────────────
Read /Users/grant/Projects/myapp/AGENTS.md for repository overview.
─────────────────────────────────────────────────

project (local, optional)
─────────────────────────────────────────────────
Read /Users/grant/Projects/myapp/PROJECT.md. Respond with summary.
─────────────────────────────────────────────────
```

**Output (named context):**

```bash
start show context environment
```

```
Context: environment (effective)
═══════════════════════════════════════════════════════════
Source: global config
Type: File-based
Required: true

Read /Users/grant/reference/ENVIRONMENT.md for environment context.
```

**Output (command-based context):**

```bash
start show context git-status
```

```
Context: git-status (effective)
═══════════════════════════════════════════════════════════
Source: local config
Type: Command-based
Required: false

Command executed:
  ❯ git status --short
  ✓ Executed successfully (127 bytes output)

Resolved content:
─────────────────────────────────────────────────
Working tree status:
 M main.go
 M README.md
?? newfile.go
─────────────────────────────────────────────────
```

### start show agent

Display effective agent configuration after config merging.

**Synopsis:**

```bash
start show agent              # Show default agent
start show agent <name>       # Show named agent
start show agent --scope global
start show agent --scope local
```

**Arguments:**

**name** (optional)
: Agent name to show. If omitted, shows the default agent.

**Flags:**

**--scope** _scope_
: Show agent from specific scope (`global` or `local`). If omitted, shows effective/merged agent.

**Behavior:**

Shows the effective agent configuration after config merging (global + local).

This is the configuration that would be used when executing the agent.

**Output (default agent):**

```bash
start show agent
```

```
Agent: claude (effective)
═══════════════════════════════════════════════════════════
Source: global config
Description: Anthropic's Claude AI assistant via Claude Code CLI
URL: https://docs.claude.com/claude-code

Command template:
  {bin} --model {model} --append-system-prompt '{role}' '{prompt}'

Default model: claude-3-7-sonnet-20250219 (sonnet)

Models:
  haiku  → claude-3-5-haiku-20241022
  sonnet → claude-3-7-sonnet-20250219
  opus   → claude-opus-4-20250514

Model docs: https://docs.anthropic.com/en/docs/about-claude/models
```

**Output (named agent):**

```bash
start show agent gemini
```

```
Agent: gemini (effective)
═══════════════════════════════════════════════════════════
Source: global config
Description: Google's Gemini AI via CLI

Command template:
  GEMINI_SYSTEM_MD='{role_file}' {bin} --model {model} '{prompt}'

Default model: gemini-2.0-flash-exp (flash)

Models:
  flash   → gemini-2.0-flash-exp
  pro-exp → gemini-2.0-pro-exp
```

### start show task

Display resolved task prompt after UTD processing.

**Synopsis:**

```bash
start show task <name>       # Show named task
start show task --scope global
start show task --scope local
```

**Arguments:**

**name** (required)
: Task name or alias to show.

**Flags:**

**--scope** _scope_
: Show task from specific scope (`global` or `local`). If omitted, shows effective/merged task.

**Behavior:**

Shows the resolved task prompt template after UTD processing (without executing).

This shows what the task prompt looks like with placeholders visible (but not yet filled).

**Output:**

```bash
start show task git-diff-review
```

```
Task: git-diff-review (effective)
═══════════════════════════════════════════════════════════
Source: global config
Alias: gdr
Description: Review staged git changes
Role: code-reviewer
Agent: (uses default)

Task prompt template (command-based):
─────────────────────────────────────────────────
Command: git diff --staged
Shell: bash
Timeout: 10 seconds

Prompt template:
  Review changes:

  ## Instructions
  {instructions}

  ## Changes
  ```diff
  {command_output}
  ```

─────────────────────────────────────────────────

Context inclusion: All required contexts

```

## Flags

Execution preview commands support all main command flags:

**--agent** _name_, **-a** _name_
: Which agent to use. Supports exact match or prefix matching.

**--role** _name_, **-r** _name_
: Which role (system prompt) to use. Supports exact match or prefix matching.

**--model** _name_, **-m** _name_
: Model to use (from agent configuration).

**--directory** _path_, **-d** _path_
: Working directory for context detection.

**--verbose**
: Show full content (no truncation in execution preview mode).

**--quiet**, **-q**
: Quiet mode (minimal output).

**--debug**
: Debug mode (show all internal operations).

Content viewer commands support:

**--scope** _scope_
: Show content from specific scope (`global` or `local`). If omitted, shows effective/merged content.

**--help**, **-h**
: Show help.

**--version**, **-v**
: Show version information.

## Exit Codes

**0** - Success (preview shown)

**1** - Configuration error (invalid config, missing fields)

**2** - Invalid arguments (wrong command usage)

**3** - File error (context file not found, working directory doesn't exist)

**4** - Command error (task command failed, though preview still shown)

## Examples

### Preview Interactive Session

```bash
start show
```

See what `start` would execute without running it.

### Preview with Verbose Output

```bash
start show --verbose
```

See full role and prompt content (no truncation).

### Preview Task Execution

```bash
start show task git-diff-review
start show task gdr "focus on security"
```

See what task would do without executing.

### Preview Custom Prompt

```bash
start show prompt "analyze security vulnerabilities"
```

See what would be sent to agent.

### View Resolved Role Content

```bash
start show role
start show role code-reviewer
start show role go-expert --scope global
```

See exactly what role content would be used.

### View All Contexts

```bash
start show context
```

See all resolved contexts and their content.

### View Specific Context

```bash
start show context environment
start show context git-status
```

See resolved content for specific context.

### View Agent Configuration

```bash
start show agent
start show agent claude
start show agent gemini --scope local
```

See effective agent configuration.

### View Task Prompt Template

```bash
start show task code-review
start show task git-diff-review --scope global
```

See task's prompt template structure.

## Use Cases

### Debugging Prompts

**Problem:** Not sure what's being sent to the agent.

```bash
start show --verbose
```

See exact prompt and role content before executing.

### Verifying Context Loading

**Problem:** Want to confirm which contexts are included.

```bash
start show
start show task code-review
```

See which contexts are loaded and their file paths.

### Understanding Task Behavior

**Problem:** Want to see what a task does before running it.

```bash
start show task git-diff-review "focus on security"
```

See task command execution and prompt composition.

### Checking Role Content

**Problem:** Want to see what role is actually being used.

```bash
start show role
start show role --verbose
```

See resolved role content after UTD processing.

### Inspecting Config Merging

**Problem:** Want to understand global vs local config behavior.

```bash
start show role code-reviewer --scope global
start show role code-reviewer --scope local
start show role code-reviewer  # merged
```

Compare configurations across scopes.

### Validating Dynamic Content

**Problem:** Want to see command output before sending to agent.

```bash
start show context git-status
start show task git-diff-review
```

See resolved output from dynamic commands.

## Comparison with Other Commands

### vs `start` (execution)

**`start`** - Executes the agent with all contexts

```bash
start  # Runs agent
```

**`start show`** - Previews what would execute, doesn't run

```bash
start show  # Shows preview only
```

### vs `start config <type> show` (config viewer)

**`start config role show`** - Shows TOML configuration structure

```bash
start config role show global
# Output: file = "~/.config/start/roles/code-reviewer.md"
#         prompt = "{file_contents}\n\nFocus on security."
```

**`start show role`** - Shows resolved content after processing

```bash
start show role
# Output: You are an expert code reviewer...
#         [actual role content after UTD processing]
```

**Key difference:** Config shows the configuration, show displays the content.

## Notes

### Truncation in Preview Mode

**Default behavior (without --verbose):**

- Role content: First 10 lines + line count
- Prompt content: First 10 lines + line count
- Command output: First 10 lines + line count

**With --verbose flag:**

- Shows full content (no truncation)
- Use when debugging or inspecting full prompts

### Content Viewer Mode vs Config Viewer

**Content viewer (`start show <type>`):**

- Shows resolved/processed content
- After UTD processing (files read, commands executed, placeholders replaced)
- After config merging (effective configuration)
- What the agent would actually receive

**Config viewer (`start config <type> show`):**

- Shows raw TOML configuration
- File paths, prompt templates, settings
- Configuration structure, not content

### Execution Preview Accuracy

Preview shows exactly what would execute, including:

- ✓ All context documents that would be loaded
- ✓ File existence checks and sizes
- ✓ Command execution (for tasks and UTD)
- ✓ Error messages if commands fail
- ✓ Exact agent command that would run

The only difference from actual execution:

- ✗ Agent is NOT invoked
- ✗ Agent binary is NOT executed

### Performance Considerations

**Execution preview mode:**

- Executes task commands (if task has `command` field)
- Reads all context files
- Processes UTD (executes commands in roles/contexts)
- Nearly identical cost to actual execution (minus agent invocation)

**Content viewer mode:**

- Faster - only resolves requested content
- Use when you just want to see specific role/context/task

## See Also

- start(1) - Launch with context
- start-task(1) - Run predefined tasks
- start-prompt(1) - Launch with custom prompt
- start-config-role(1) - Manage role configuration
- start-config-context(1) - Manage context configuration
- start-config-agent(1) - Manage agent configuration
- start-config-task(1) - Manage task configuration
