# DR-007: Command Interpolation and Placeholders

**Date:** 2025-01-03 (original), 2025-01-07 (updated)
**Status:** Accepted
**Category:** Configuration

## Decision

Use single-brace placeholders (`{name}`) with specific supported variables for command interpolation, role handling, and template processing.

## Global Placeholders

Available in agent commands and all template contexts:

**{model}** - Currently selected model identifier
- Value: Full model identifier after name resolution
- Example: `"claude-3-7-sonnet-20250219"`
- Example: `"gemini-2.0-flash-exp"`

**{prompt}** - Assembled prompt text
- Value: Final prompt after context document inclusion and template processing
- Used by: Agent commands
- Contains: Required context prompts + optional context prompts + custom prompt (if using `start prompt`)

**{date}** - Current timestamp
- Value: ISO 8601 format with timezone
- Example: `"2025-01-07T14:30:00+10:00"`
- Updated each execution

**{role}** - Resolved role content (inline text)
- Value: Complete role content after UTD processing
- Use for: Agents that accept system prompts inline (Claude, aichat)
- Simple role (file only): File contents
- UTD role: Fully resolved content (file + command + template)

**{role_file}** - Role file path
- Value: Absolute file path to role content
- Use for: Agents that require system prompt files (Gemini)
- Simple role: Original file path (expanded to absolute path)
- UTD role: Temporary file path with resolved content
- Temp files auto-created and cleaned up

## Role Placeholder Details

**When to use {role}:**
```toml
[agents.claude]
command = "claude --model {model} --append-system-prompt '{role}' '{prompt}'"

[agents.aichat]
command = "AICHAT_SYSTEM_PROMPT='{role}' aichat --model {model} '{prompt}'"
```

**When to use {role_file}:**
```toml
[agents.gemini]
command = "GEMINI_SYSTEM_MD='{role_file}' gemini --model {model} '{prompt}'"
```

**Temporary file behavior for {role_file}:**

Simple role (file only):
```toml
[roles.reviewer]
file = "~/.config/start/roles/reviewer.md"
```
→ `{role_file}` = `/Users/username/.config/start/roles/reviewer.md` (no temp file)

UTD role (complex):
```toml
[roles.dynamic]
file = "base.md"
command = "date"
prompt = "{file}\n\nDate: {command}"
```
→ `{role_file}` = `/tmp/start-role-a8f3b2c1.md` (temp file with resolved content)

Temp files:
- Created before agent execution
- Permissions: 0600 (owner read/write only)
- Deleted after agent execution (success or failure)
- Random suffix prevents collisions

## UTD Pattern Placeholders

Available in UTD template contexts (roles, contexts, tasks):

**{file}** - File path from UTD `file` field
- Replaced with absolute file path (~ expanded)
- Example: `file = "~/ref/ENV.md"` → `{file}` = `"/Users/username/ref/ENV.md"`
- Use case: Instructing AI agents with file access to read files
- Example: `prompt = "Read {file} for context."`

**{file_contents}** - Content from UTD `file` field
- Replaced with complete file contents
- Example in role: `{file_contents}` replaced with role file contents
- Example in context: `{file_contents}` replaced with context file contents
- Use case: Injecting file content directly into prompts for AI agents without file access
- Example: `prompt = "Environment info:\n{file_contents}"`

**{command}** - Command string from UTD `command` field
- Replaced with the command string as written
- Example: `command = "git status"` → `{command}` = `"git status"`
- Use case: Documenting what command was run
- Example: `prompt = "Output of '{command}':\n{command_output}"`

**{command_output}** - Output from UTD `command` field
- Stdout and stderr captured and combined
- Empty string if command fails or produces no output
- Example: `command = "date"` → `{command_output}` = `"2025-01-07 14:30:00"`
- Use case: Injecting dynamic command results into prompts

## Task-Specific Placeholders

Available only in task prompt templates (not in roles or contexts):

**{instructions}** - User's command-line arguments
- Value: Arguments after task name
- Empty case: `"None"` (not empty string)
- Example: `start task review "focus on security"` → `{instructions}` = `"focus on security"`
- Example: `start task review` → `{instructions}` = `"None"`

Tasks also have access to all UTD Pattern Placeholders: `{file}`, `{file_contents}`, `{command}`, `{command_output}`.

See [DR-009](./dr-009-task-structure.md) for task structure details.

## Path Expansion

**Tilde (~) expansion:**
- `~` expands to user's home directory
- Applied before all file operations
- Example: `~/reference/file.md` → `/Users/username/reference/file.md`

**Relative paths:**
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

**No separate env section:** The `[agents.<name>.env]` pattern from earlier designs is **removed**. Use shell syntax in the `command` field.

## Usage Examples

### Agent with Inline Role

```toml
[agents.claude]
description = "Claude with inline system prompt"
command = "claude --model {model} --append-system-prompt '{role}' '{prompt}'"
default_model = "sonnet"
```

### Agent with File-Based Role

```toml
[agents.gemini]
description = "Gemini with file-based system prompt"
command = "GEMINI_SYSTEM_MD='{role_file}' gemini --model {model} --include-directories ~/reference '{prompt}'"
default_model = "flash"
```

### Agent with Env Var Content

```toml
[agents.aichat]
description = "aichat with system prompt in env var"
command = "AICHAT_SYSTEM_PROMPT='{role}' aichat --model {model} '{prompt}'"
default_model = "gpt4-mini"
```

### Role with UTD Pattern

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

### Context with Placeholders

```toml
[context.environment]
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
```

### Task with Placeholders

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

## Substitution Behavior

**Order of operations:**

1. **Role resolution:**
   - Select role (precedence rules)
   - Resolve role UTD (file, command, prompt)
   - Generate `{role}` content
   - Generate `{role_file}` path (create temp file if needed)

2. **Context resolution:**
   - Load required and optional contexts
   - Resolve each context's UTD
   - Build context prompts

3. **Task resolution** (if executing task):
   - Resolve task UTD (file, command, prompt)
   - Replace `{instructions}` with user args
   - Replace `{file}` with task file path
   - Replace `{file_contents}` with task file content
   - Replace `{command}` with task command string
   - Replace `{command_output}` with task command output

4. **Final assembly:**
   - Assemble `{prompt}` from contexts + task/custom prompt
   - Replace `{date}` with current timestamp
   - Replace `{model}` with selected model name

5. **Agent command:**
   - Replace all placeholders in agent command
   - Execute command
   - Cleanup temp files

**Placeholder scope:**

| Placeholder       | Agent Commands | Roles | Contexts | Tasks |
|-------------------|----------------|-------|----------|-------|
| {model}           | ✓              | ✓     | ✓        | ✓     |
| {date}            | ✓              | ✓     | ✓        | ✓     |
| {prompt}          | ✓              | -     | -        | -     |
| {role}            | ✓              | -     | -        | -     |
| {role_file}       | ✓              | -     | -        | -     |
| {file}            | -              | ✓     | ✓        | ✓     |
| {file_contents}   | -              | ✓     | ✓        | ✓     |
| {command}         | -              | ✓     | ✓        | ✓     |
| {command_output}  | -              | ✓     | ✓        | ✓     |
| {instructions}    | -              | -     | -        | ✓     |

## Rationale

**Single braces ({}):**
- Simpler syntax than double braces
- Clear and readable
- Standard in many template systems

**{role} and {role_file} instead of {system_prompt}:**
- "Role" better describes what it is (AI's persona/instructions)
- Two placeholders support both inline and file-based agents
- Allows agents to choose their preferred pattern
- See [DR-005](./dr-005-role-configuration.md) for role design

**Shell syntax for environment variables:**
- Standard and familiar to all users
- No need to implement separate env section
- Supports multiple variables easily
- Works with any shell command syntax

**No environment variable placeholders ({env:VAR}):**
- Agents inherit environment naturally
- Use shell syntax instead: `VAR=value command`
- Simpler implementation
- Fewer concepts to document

**Tilde (~) instead of {home}:**
- More concise
- Standard Unix convention
- Less typing

**No {cwd} placeholder:**
- Use `--directory` flag to control working directory
- Keeps placeholders focused on content, not paths
- Working directory is an execution context, not a template value

## Not Supported

**Environment variable placeholders:**
```toml
# NOT SUPPORTED
command = "tool --path {env:HOME}/file"

# USE INSTEAD (shell syntax)
command = "tool --path ~/file"
# or
command = "HOME_PATH=~ tool --path \"$HOME_PATH/file\""
```

**Home directory placeholder:**
```toml
# NOT SUPPORTED
file = "{home}/reference/file.md"

# USE INSTEAD
file = "~/reference/file.md"
```

**Current working directory:**
```toml
# NOT SUPPORTED
file = "{cwd}/ROLE.md"

# USE INSTEAD
file = "./ROLE.md"  # Relative path
# or use --directory flag to change working directory
```

## Breaking Changes from Original

This updates the original DR-007 with:

1. **Removed:** `{system_prompt}` placeholder
2. **Added:** `{role}` placeholder (inline content)
3. **Added:** `{role_file}` placeholder (file path)
4. **Removed:** `[agents.<name>.env]` section reference
5. **Added:** Shell environment variable syntax documentation
6. **Added:** Temporary file behavior for UTD roles
7. **Updated:** All examples to use role-based design

## Related Decisions

- [DR-005](./dr-005-role-configuration.md) - Role configuration and selection
- [DR-009](./dr-009-task-structure.md) - Task structure with {instructions} placeholder
- [DR-002](./dr-002-config-merge.md) - Config merge affecting placeholder resolution
