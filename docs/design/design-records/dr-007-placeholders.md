# DR-007: Command Interpolation and Placeholders

- Date: 2025-01-03 (original), 2025-01-07 (updated)
- Status: Accepted
- Category: Configuration

## Problem

The tool needs a template system for dynamic content composition. The system must:

- Allow agents to receive runtime values (model, prompt, role, etc.)
- Support both static files and dynamic content (command output, timestamps)
- Work across different contexts (agent commands, roles, tasks, contexts)
- Be simple enough for users to understand and use
- Prevent conflicts with shell syntax or other template systems
- Support both inline content and file-based patterns for roles
- Handle path expansion consistently (home directory, relative paths)
- Keep placeholders available only where they make logical sense

## Decision

Use single-brace placeholders (`{name}`) with specific supported variables for command interpolation, role handling, and template processing.

Placeholder categories:

1. Universal placeholder - available everywhere: `{date}`
2. Agent command placeholders - available in agent commands only: `{bin}`, `{model}`, `{prompt}`, `{role}`, `{role_file}`
3. UTD pattern placeholders - available where UTD is supported (roles, contexts, tasks): `{file}`, `{file_contents}`, `{command}`, `{command_output}`
4. Task-specific placeholders - available in tasks only: `{instructions}`

Path expansion:

- Tilde (~) expands to home directory
- Relative paths resolve from working directory (local) or config directory (global)

Environment variables:

- Use standard shell syntax in command field: `VAR=value command`
- No separate env section or env placeholders

## Why

Single-brace syntax:

- Simpler than double braces `{{name}}`
- Clear and readable
- Standard in many template systems
- Less likely to conflict with shell syntax than $VAR

Scoped placeholders (available only where needed):

- Prevents confusing or fragile configurations
- `{model}` in agent commands makes sense (where model is used)
- `{model}` in roles/contexts doesn't make sense (creates model-dependent, brittle configs)
- Roles and contexts should be reusable across different models/agents
- Clearer mental model: placeholders available where logically appropriate

Universal {date}:

- Timestamping is universally useful
- Roles: "Review conducted on {date}"
- Tasks: "Analysis performed on {date}"
- Agent commands: timestamp execution
- Harmless everywhere, useful in many places

Separate role placeholders ({role} and {role_file}):

- Some agents accept inline content (Claude: --append-system-prompt)
- Other agents require files (Gemini: GEMINI_SYSTEM_MD env var)
- Two placeholders support both patterns
- Allows agents to choose their preferred method
- Clear which pattern each agent uses

UTD pattern separation:

- Different placeholders for reference vs content: `{file}` (path) vs `{file_contents}` (contents)
- Different placeholders for command vs output: `{command}` (string) vs `{command_output}` (result)
- Short form = source/reference, long form = result/contents
- Gives users flexibility in how they compose prompts

Shell syntax for environment variables:

- Standard and familiar to all users
- No need to implement separate env section
- Supports multiple variables easily: `VAR1=a VAR2=b command`
- Works with any shell command syntax
- One less config section to document

Tilde expansion:

- Standard Unix convention
- More concise than `{home}` placeholder
- Less typing
- Users already understand this pattern

No {cwd} placeholder:

- Use --directory flag to control working directory
- Keeps placeholders focused on content, not execution context
- Working directory is a runtime setting, not a template value

## Trade-offs

Accept:

- Users must learn which placeholders work in which contexts
- Two role placeholders ({role} and {role_file}) instead of one
- Placeholder scope varies by context (prevents some technically-possible but illogical uses)
- Path expansion rules differ for global vs local configs

Gain:

- Works with all agent API styles (inline and file-based)
- Flexible content composition (file references vs inline content)
- Simple, readable syntax
- No conflicts with shell syntax
- Standard path expansion patterns
- Minimal concepts to learn (no env sections, no home placeholders)
- Prevents brittle, model-dependent roles/contexts
- Clear logical boundaries for placeholder usage

## Universal Placeholder

Available everywhere (agent commands, roles, contexts, tasks):

{date} - Current timestamp:

- Value: ISO 8601 format with timezone
- Example: `"2025-01-07T14:30:00+10:00"`
- Updated each execution
- Use case: Timestamping any operation

## Agent Command Placeholders

Available in agent commands only:

{bin} - Agent binary name:

- Value: Binary name from agent's `bin` field
- Example: `"claude"`, `"gemini"`
- Required in agent commands (enforces DRY, prevents bin/command mismatch)

{model} - Currently selected model identifier:

- Value: Full model identifier after name resolution
- Example: `"claude-3-7-sonnet-20250219"`, `"gemini-2.0-flash-exp"`
- Available in agent commands only (execution detail, not content concern)

{prompt} - Assembled prompt text:

- Value: Final prompt after context document inclusion and template processing
- Contains: Required context prompts + optional context prompts + custom prompt (if using `start prompt`)

{role} - Resolved role content (inline text):

- Value: Complete role content after UTD processing
- Use for: Agents that accept system prompts inline (Claude, aichat)
- Simple role (file only): File contents
- UTD role: Fully resolved content (file + command + template)

{role_file} - Role file path:

- Value: Absolute file path to role content
- Use for: Agents that require system prompt files (Gemini)
- Simple role: Original file path (expanded to absolute path)
- UTD role: Temporary file path with resolved content
- Temp files auto-created and cleaned up

## Role Placeholder Details

When to use {role}:

```toml
[agents.claude]
bin = "claude"
command = "{bin} --model {model} --append-system-prompt '{role}' '{prompt}'"

[agents.aichat]
bin = "aichat"
command = "AICHAT_SYSTEM_PROMPT='{role}' {bin} --model {model} '{prompt}'"
```

When to use {role_file}:

```toml
[agents.gemini]
bin = "gemini"
command = "GEMINI_SYSTEM_MD='{role_file}' {bin} --model {model} '{prompt}'"
```

Temporary file behavior for {role_file}:

Simple role (file only):

```toml
[roles.reviewer]
file = "~/.config/start/roles/reviewer.md"
```

`{role_file}` = `/Users/username/.config/start/roles/reviewer.md` (no temp file)

UTD role (complex):

```toml
[roles.dynamic]
file = "base.md"
command = "date"
prompt = "{file}\n\nDate: {command}"
```

`{role_file}` = `/tmp/start-role-a8f3b2c1.md` (temp file with resolved content)

Temp files:

- Created before agent execution
- Permissions: 0600 (owner read/write only)
- Deleted after agent execution (success or failure)
- Random suffix prevents collisions

## UTD Pattern Placeholders

Available where UTD is supported (roles, contexts, tasks):

{file} - File path from UTD `file` field:

- Replaced with absolute file path (~ expanded)
- Example: `file = "~/reference/ENVIRONMENT.md"` → `{file}` = `"/Users/username/reference/ENVIRONMENT.md"`
- Use case: Instructing AI agents with file access to read files
- Example: `prompt = "Read {file} for context."`

{file_contents} - Content from UTD `file` field:

- Replaced with complete file contents
- Example in role: `{file_contents}` replaced with role file contents
- Example in context: `{file_contents}` replaced with context file contents
- Use case: Injecting file content directly into prompts for AI agents without file access
- Example: `prompt = "Environment info:\n{file_contents}"`

{command} - Command string from UTD `command` field:

- Replaced with the command string as written
- Example: `command = "git status"` → `{command}` = `"git status"`
- Use case: Documenting what command was run
- Example: `prompt = "Output of '{command}':\n{command_output}"`

{command_output} - Output from UTD `command` field:

- Stdout and stderr captured and combined
- Empty string if command fails or produces no output
- Example: `command = "date"` → `{command_output}` = `"2025-01-07 14:30:00"`
- Use case: Injecting dynamic command results into prompts

## Task-Specific Placeholders

Available in tasks only:

{instructions} - User's command-line arguments:

- Value: Arguments after task name
- Empty case: `"None"` (not empty string)
- Example: `start task review "focus on security"` → `{instructions}` = `"focus on security"`
- Example: `start task review` → `{instructions}` = `"None"`

Tasks also have access to all UTD Pattern Placeholders and {date}.

## Path Expansion

Tilde (~) expansion:

- `~` expands to user's home directory
- Applied before all file operations
- Example: `~/reference/file.md` → `/Users/username/reference/file.md`

Relative paths:

- Global config: Relative to home or config directory (context-dependent)
- Local config: Relative to working directory
- Example in local: `./ROLE.md` → `/Users/username/Projects/myapp/ROLE.md`

## Environment Variables in Commands

Use standard shell syntax to set environment variables:

```toml
# Single variable
[agents.gemini]
command = "GEMINI_SYSTEM_MD='{role_file}' gemini --model {model} '{prompt}'"

# Multiple variables
[agents.custom]
command = "VAR1='{role}' VAR2='{date}' custom-ai --model {model} '{prompt}'"
```

No separate env section - the `[agents.<name>.env]` pattern from earlier designs is removed. Use shell syntax in the `command` field.

## Placeholder Scope

| Placeholder       | Agent Commands | Roles | Contexts | Tasks |
|-------------------|----------------|-------|----------|-------|
| {date}            | ✓              | ✓     | ✓        | ✓     |
| {bin}             | ✓              | -     | -        | -     |
| {model}           | ✓              | -     | -        | -     |
| {prompt}          | ✓              | -     | -        | -     |
| {role}            | ✓              | -     | -        | -     |
| {role_file}       | ✓              | -     | -        | -     |
| {file}            | -              | ✓     | ✓        | ✓     |
| {file_contents}   | -              | ✓     | ✓        | ✓     |
| {command}         | -              | ✓     | ✓        | ✓     |
| {command_output}  | -              | ✓     | ✓        | ✓     |
| {instructions}    | -              | -     | -        | ✓     |

## Usage Examples

Agent with inline role:

```toml
[agents.claude]
bin = "claude"
description = "Claude with inline system prompt"
command = "{bin} --model {model} --append-system-prompt '{role}' '{prompt}'"
default_model = "sonnet"
```

Agent with file-based role:

```toml
[agents.gemini]
bin = "gemini"
description = "Gemini with file-based system prompt"
command = "GEMINI_SYSTEM_MD='{role_file}' {bin} --model {model} --include-directories ~/reference '{prompt}'"
default_model = "flash"
```

Agent with env var content:

```toml
[agents.aichat]
bin = "aichat"
description = "aichat with system prompt in env var"
command = "AICHAT_SYSTEM_PROMPT='{role}' {bin} --model {model} '{prompt}'"
default_model = "gpt4-mini"
```

Role with UTD pattern:

```toml
[roles.code-reviewer]
file = "~/.config/start/roles/reviewer.md"
command = "date '+%Y-%m-%d'"
prompt = """
{file_contents}

Review Date: {command_output}

Focus on security and performance.
"""
```

When this role is used:
- `{role}` = fully resolved content (file contents + date + framing)
- `{role_file}` = temp file path (because it uses `command`)

Context with placeholders:

```toml
[contexts.environment]
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
```

Task with placeholders:

```toml
[tasks.git-review]
role = "code-reviewer"
command = "git diff --staged"
prompt = """
Review these changes:

{command_output}

Instructions: {instructions}
"""
```

## Substitution Order

Order of operations:

1. Role resolution:
   - Select role (precedence rules)
   - Resolve role UTD (file, command, prompt)
   - Replace {date} in role prompt
   - Generate `{role}` content
   - Generate `{role_file}` path (create temp file if needed)

2. Context resolution:
   - Load required and optional contexts
   - Resolve each context's UTD
   - Replace {date} in context prompts
   - Build context prompts

3. Task resolution (if executing task):
   - Resolve task UTD (file, command, prompt)
   - Replace `{instructions}` with user args
   - Replace `{file}` with task file path
   - Replace `{file_contents}` with task file content
   - Replace `{command}` with task command string
   - Replace `{command_output}` with task command output
   - Replace `{date}` in task prompt

4. Final assembly:
   - Assemble `{prompt}` from contexts + task/custom prompt
   - Replace `{date}` with current timestamp in final prompt
   - Replace `{model}` with selected model identifier

5. Agent command:
   - Replace all placeholders in agent command ({bin}, {model}, {prompt}, {role}, {role_file}, {date})
   - Execute command
   - Cleanup temp files

## Not Supported

Environment variable placeholders:

```toml
# NOT SUPPORTED
command = "tool --path {env:HOME}/file"

# USE INSTEAD (shell syntax)
command = "tool --path ~/file"
# or
command = "HOME_PATH=~ tool --path \"$HOME_PATH/file\""
```

Home directory placeholder:

```toml
# NOT SUPPORTED
file = "{home}/reference/file.md"

# USE INSTEAD
file = "~/reference/file.md"
```

Current working directory:

```toml
# NOT SUPPORTED
file = "{cwd}/ROLE.md"

# USE INSTEAD
file = "./ROLE.md"  # Relative path
# or use --directory flag to change working directory
```

## Alternatives

Double-brace placeholders `{{name}}`:

- Pro: Less likely to conflict with shell syntax
- Pro: Common in some template systems (Jinja2, Handlebars)
- Con: More verbose, harder to type
- Con: Less readable in compact command strings
- Con: Still possible to conflict with shell patterns
- Rejected: Single braces are simpler and more readable

Dollar-sign variables `$name` or `${name}`:

- Pro: Familiar to shell users
- Pro: Very concise
- Con: Direct conflict with shell variable syntax
- Con: Would require escaping in shell commands
- Con: Ambiguous whether it's a shell var or our placeholder
- Rejected: Too much conflict with shell syntax

Environment variable placeholders `{env:VAR}`:

- Pro: Could access any environment variable
- Pro: Flexible for user customization
- Con: Agents inherit environment naturally anyway
- Con: Shell syntax already supports this: `VAR=value command`
- Con: Adds complexity for minimal benefit
- Rejected: Shell syntax is sufficient

Home directory placeholder `{home}`:

- Pro: Explicit and clear
- Pro: Cross-platform (works on Windows)
- Con: More verbose than tilde
- Con: Tilde is Unix standard, well-understood
- Con: One more placeholder to learn
- Rejected: Tilde is simpler and standard

Working directory placeholder `{cwd}`:

- Pro: Could make relative paths explicit
- Pro: Useful for documentation
- Con: Adds execution context to templates
- Con: Relative paths already work naturally
- Con: --directory flag controls this already
- Rejected: Relative paths and --directory are sufficient

Single role placeholder (no {role_file}):

- Pro: Simpler - only one placeholder to learn
- Pro: Fewer concepts
- Con: Must create temp files for ALL agents (wasteful)
- Con: Extra I/O for agents that support inline (Claude, aichat)
- Con: Hides whether agent uses inline or file-based API
- Rejected: Two placeholders give agents flexibility and efficiency

Make {model} universally available:

- Pro: Technically possible, simpler scoping rules
- Pro: Users could reference model anywhere
- Con: Encourages model-dependent roles/contexts (brittle)
- Con: Roles should be reusable across different models
- Con: Contexts describe environment, not execution details
- Con: No legitimate use case in roles/contexts/tasks
- Rejected: Model is an execution detail, belongs in agent commands only

## Breaking Changes from Original

This updates the original DR-007 with:

1. Removed: `{system_prompt}` placeholder
2. Added: `{role}` placeholder (inline content)
3. Added: `{role_file}` placeholder (file path)
4. Removed: `[agents.<name>.env]` section reference
5. Added: Shell environment variable syntax documentation
6. Added: Temporary file behavior for UTD roles
7. Updated: All examples to use role-based design
8. Changed: `{model}` scope from universal to agent-commands-only

## Updates

- 2025-01-07: Role-based placeholders, removed system_prompt, added temp file handling
- 2025-01-17: Changed {model} scope from universal to agent-commands-only (design consistency)
