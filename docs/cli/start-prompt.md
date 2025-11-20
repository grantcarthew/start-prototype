# start prompt

## Name

start prompt - Launch AI agent with custom prompt and context

## Synopsis

```bash
start prompt [text] [flags]
```

## Description

Launches an AI agent with an optional custom prompt combined with required project context documents. This is useful for one-off queries or exploratory sessions.

**Context document behavior:**

- Only **required** context documents are included (documents with `required = true`)
- Optional documents (default behavior) are excluded
- This keeps the prompt focused for specific queries

**The final prompt sent to the agent:**

1. Required context document instructions (from config)
2. Your custom prompt text (if provided, appended last)

**If you want ONLY a custom prompt with no context:**

- Don't use `start` - use your agent directly: `claude "your prompt"`

## Arguments

**text** (optional)
: Custom prompt text to send to the agent. Multi-word prompts must be quoted.

```bash
start prompt "analyze this codebase for security vulnerabilities"
start prompt  # Launch with required context only, no custom prompt
```

## Flags

All global flags from `start` command are supported:

**--agent** _name_, **-a** _name_
: Which agent to use. Supports exact match or prefix matching with ambiguity handling.

**--role** _name_, **-r** _name_
: Which role (system prompt) to use. Supports exact match or prefix matching with ambiguity handling.

**--model** _name_, **-m** _name_
: Model to use (from agent configuration). Supports exact match, prefix matching, and passthrough.

**--directory** _path_, **-d** _path_
: Working directory for context detection

**--quiet**, **-q**
: Quiet mode (no output)

**--verbose**
: Verbose output

**--debug**
: Debug mode

**--asset-download[=bool]**
: Enable or disable downloading assets from the GitHub catalog on-demand. Defaults to `true`. Use `--asset-download=false` to prevent network requests for missing assets.

**-l, --local**
: When downloading assets from the catalog (roles, agents, contexts), download them to the global cache (`~/.config/start/assets/`) but add the configuration entry to the local config (`./.start/`) instead of global. Only applies when an asset is downloaded.

**--help**, **-h**
: Show help

**--version**, **-v**
: Show version information

## Behavior

### Execution Flow

1. Load and merge configuration (global + local)
2. Filter context documents (only include `required = true`)
3. Detect which required documents exist
4. Build prompt combining:
   - Required context document prompts (first, in config definition order)
   - Your custom prompt text (appended last, if provided)
5. Resolve placeholders in agent command template
6. Display context summary (unless `--quiet`)
7. Execute agent command

**Document order:** Required documents appear in the prompt in the order they are defined in the config file. See `start` command documentation for details.

### Prompt Structure

The final prompt sent to the agent:

```
Read /Users/gcarthew/reference/ENVIRONMENT.md for environment context.
Read /Users/gcarthew/reference/INDEX.csv for documentation index.
Read /Users/gcarthew/Projects/my-app/AGENTS.md for repository overview.

{your_custom_text}
```

**Note:**

- Only documents with `required = true` are included
- Documents appear in config definition order
- Document prompts come from your config's `prompt` field with `{file}` replaced by actual paths
- Custom prompt is appended LAST

Example config:

```toml
[contexts.environment]
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
required = true    # Included

[contexts.index]
file = "~/reference/INDEX.csv"
prompt = "Read {file} for documentation index."
required = true    # Included

[contexts.agents]
file = "./AGENTS.md"
prompt = "Read {file} for repository overview."
required = true    # Included

[contexts.project]
file = "./PROJECT.md"
prompt = "Read {file}. Respond with summary."
# required = false (default) - NOT included with start prompt
```

### Shell Quoting

Multi-word prompts **must be quoted** to prevent shell interpretation:

```bash
# ✗ Wrong - shell splits into multiple arguments
start prompt analyze this codebase

# ✓ Correct - quoted as single argument
start prompt "analyze this codebase"
```

## Examples

### Basic Usage

No custom prompt (required context only):

```bash
start prompt
```

Simple prompt with required context:

```bash
start prompt "analyze security vulnerabilities"
```

Multi-line prompt (using quotes):

```bash
start prompt "Review this codebase and:
1. Identify security issues
2. Suggest performance improvements
3. Check for Go best practices"
```

### With Agent Selection

Use specific agent:

```bash
start prompt "review the API design" --agent gemini
start prompt "optimize this algorithm" --agent claude
```

### With Model Selection

Use specific model tier:

```bash
start prompt "quick code review" --model fast
start prompt "comprehensive analysis" --model pro
```

Use full model identifier:

```bash
start prompt "review auth flow" --model claude-opus-4-20250514
```

### With Directory Override

Analyze different project:

```bash
start prompt "what is this project about?" --directory ~/other-project
```

### Combined Flags

Full example:

```bash
start prompt "security audit" --agent claude --model pro --directory ~/api-server --verbose
```

Quiet mode:

```bash
start prompt "review error handling" --quiet
```

## Output

### Normal Output (with custom prompt)

```
Starting AI Agent
===============================================================================================
Agent: claude (model: claude-3-7-sonnet-20250219)

Custom prompt: analyze security vulnerabilities

Context documents (required only):
  ✓ environment     ~/reference/ENVIRONMENT.md
  ✓ index           ~/reference/INDEX.csv
  ✓ agents          ./AGENTS.md

System prompt: ./ROLE.md

Executing command...
❯ claude --model claude-3-7-sonnet-20250219 --append-system-prompt '...' 'Read...'
```

### Normal Output (no custom prompt)

```
Starting AI Agent
===============================================================================================
Agent: claude (model: claude-3-7-sonnet-20250219)

Context documents (required only):
  ✓ environment     ~/reference/ENVIRONMENT.md
  ✓ index           ~/reference/INDEX.csv
  ✓ agents          ./AGENTS.md

System prompt: ./ROLE.md

Executing command...
❯ claude --model claude-3-7-sonnet-20250219 --append-system-prompt '...' 'Read...'
```

### Quiet Output

```
(no output - launches agent directly)
```

## Exit Codes

**0** - Success (agent launched successfully)

**1** - Configuration error (invalid config, missing fields)

**2** - Agent error (agent not found, model not configured)

**3** - File error (working directory doesn't exist)

**4** - Runtime error (agent tool not installed, command failed)

## Common Patterns

### Code Analysis

```bash
start prompt "analyze this codebase for common bug patterns"
start prompt "identify code duplication and refactoring opportunities"
start prompt "review for Go idioms and best practices"
```

### Architecture Review

```bash
start prompt "explain the overall architecture of this project"
start prompt "identify architectural issues and suggest improvements"
start prompt "diagram the component relationships"
```

### Documentation Questions

```bash
start prompt "what is missing from the documentation?"
start prompt "generate API documentation from the code"
start prompt "create a getting started guide"
```

### Performance Analysis

```bash
start prompt "identify performance bottlenecks"
start prompt "suggest optimization opportunities"
start prompt "review memory usage patterns"
```

## Comparison with Other Commands

### vs `start` (root command)

**`start`** - ALL context documents (required + optional), no custom prompt

```bash
start  # Launches with all context documents
```

**`start prompt`** - ONLY required context documents, optional custom prompt

```bash
start prompt                      # Required context only, no custom prompt
start prompt "analyze security"   # Required context + custom prompt
```

**Key difference:** `start` includes ALL documents, `start prompt` includes ONLY required documents.

### vs `start task`

**`start prompt`** - One-off custom prompt with required context

```bash
start prompt "review this specific function"
```

**`start task`** - Reusable workflow with predefined prompt template and configurable documents

```bash
start task code-review  # Uses predefined template and specified documents
```

**When to use `start prompt`:**

- Exploratory questions with minimal context
- One-off analysis
- Custom requests not covered by tasks
- When you want required context only (not all documents)

**When to use `start task`:**

- Repeatable workflows
- Standardized reviews
- When you want dynamic content (e.g., git diff)
- When you need specific document combinations
- Team-shared workflows

## Notes

### Required Context Documents

`start prompt` ONLY includes documents marked with `required = true`. This provides focused context for one-off queries.

**To include all documents:** Use `start` (root command) instead.

**To include no documents:** Use your agent directly:

```bash
claude "your prompt here"
gemini "your prompt here"
```

**To include specific documents:** Use `start task` (auto-includes all contexts where `required = true`).

### Prompt Length Limits

Be aware of:

- Shell command line length limits (~100KB on most systems)
- Agent-specific prompt length limits
- Model context windows

Very long prompts may need to be:

- Split into multiple messages
- Put into a file and referenced via task's `command`
- Sent through the agent's native interface instead

### Escaping Special Characters

When your prompt contains shell special characters, use appropriate quoting:

```bash
# Single quotes - most literal (no variable expansion)
start prompt 'analyze $variable usage'

# Double quotes - allows some expansion
start prompt "analyze the ${PROJECT} codebase"

# Escape special characters
start prompt "analyze \$dollar and \"quotes\""
```

## See Also

- start(1) - Launch with context only
- start-task(1) - Run predefined tasks
- start-init(1) - Initialize configuration
- start-config-agent(1) - Manage agents
- start-config(1) - Manage configuration
