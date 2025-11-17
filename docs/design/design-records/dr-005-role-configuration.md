# DR-005: Role Configuration & Selection

- Date: 2025-01-03 (original), 2025-01-07 (rewritten)
- Status: Accepted
- Category: Configuration

## Problem

System prompts (roles) need to be:

- Configurable and reusable across multiple tasks
- Selectable at runtime (task-specific, CLI override, defaults)
- Compatible with different agent APIs (inline content vs file-based)
- Support both simple (static file) and complex (dynamic, templated) definitions
- Work in both global (personal) and local (project/team) scopes

The design must balance simplicity for basic use cases with flexibility for advanced scenarios.

## Decision

Roles (system prompts) are configured as named entities using the Unified Template Design (UTD) pattern, selected via precedence rules, and made available to agents through two placeholders: `{role}` (content) and `{role_file}` (file path).

Role selection precedence:

1. CLI `--role` flag (highest priority)
2. Task `role` field (if executing a task)
3. `default_role` setting
4. First role in config (TOML order)

Agent integration:

- `{role}` placeholder - inline content (for agents with inline API)
- `{role_file}` placeholder - file path (for agents with file-based API)

## Why

Named roles for reusability:

- Define once, reference in multiple tasks
- Consistent with `[agents.<name>]` and `[tasks.<name>]` pattern
- Clear mental model: "Role" describes what the AI is
- Easy to swap roles via `--role` flag

Two-placeholder design for agent compatibility:

- Some agents accept inline content (Claude: `--append-system-prompt`)
- Other agents require files (Gemini: `GEMINI_SYSTEM_MD` env var pointing to file)
- Agents with file-based APIs can read original file directly (efficient)
- Agents with inline APIs get content directly (no extra I/O)

UTD pattern for flexibility:

- Simple case: just a file reference
- Advanced: dynamic content via commands, templating
- One pattern supports both static and dynamic roles

Precedence rules for flexibility:

- Tasks can specify default role
- Users can override via CLI flag
- Global default prevents requiring explicit role selection
- Clear, predictable order

Temporary files for UTD roles:

- `{role_file}` always returns a file path (consistency)
- Simple roles return original file path (no overhead)
- UTD roles create temp file with resolved content
- File-based agents work with both simple and complex roles

## Trade-offs

Accept:

- Two placeholders to learn (`{role}` vs `{role_file}`)
- Temporary file management for UTD roles with `{role_file}`
- Users must understand precedence rules for role selection
- Slightly more complex than single system_prompt field

Gain:

- Works with all agent API styles (inline and file-based)
- Reusable roles across tasks
- Runtime selection flexibility
- Supports both static and dynamic roles
- Natural merge behavior with global/local configs
- Efficient (no temp files unless needed)

## Configuration Structure

Roles use `[roles.<name>]` section with full UTD support:

```toml
[roles.code-reviewer]
description = "Expert code reviewer"
file = "~/.config/start/roles/code-reviewer.md"
command = "date"
prompt = """
{file}

Review date: {command}
"""
```

Fields:

- `description` (optional) - Human-readable description
- `file` (optional) - Path to role content file
- `command` (optional) - Shell command for dynamic content
- `prompt` (optional) - Template with `{file}` and `{command}` placeholders

UTD Requirement: At least one of `file`, `command`, or `prompt` must be present.

Default role:

```toml
[settings]
default_role = "code-reviewer"
```

If `default_role` is not specified, the first role defined in the config (TOML order) is used.

File: `roles.toml`

No `[system_prompt]` section - all system prompts are defined as named roles.

## Role Selection Precedence

When executing `start`, `start prompt`, or `start task`:

1. CLI `--role` flag (highest priority)
2. Task `role` field (if executing a task)
3. `default_role` setting
4. First role in config (TOML order)

Examples:

```bash
# Uses default_role
start

# Uses specified role
start --role security-auditor

# Task with no role field → uses default_role
start task code-review

# Task with role field → uses task's role
[tasks.security-scan]
role = "security-auditor"

start task security-scan

# Override task's role with CLI flag
start task security-scan --role code-reviewer
```

## Agent Integration

Roles are made available to agents through two placeholders:

`{role}` placeholder:

- Contains fully resolved role content (inline text)
- Use for agents that accept system prompts inline via command flags or env vars
- Resolution: file content (simple role) or full UTD processing result (complex role)

Example:

```toml
[agents.claude]
command = "claude --model {model} --append-system-prompt '{role}' '{prompt}'"

[agents.aichat]
command = "AICHAT_SYSTEM_PROMPT='{role}' aichat --model {model} '{prompt}'"
```

`{role_file}` placeholder:

- Contains a file path to the role content
- Use for agents that require system prompts from a file
- Simple role (file only): returns the original file path (expanded, absolute)
- UTD role (complex): creates temporary file with resolved content, returns temp file path

Example:

```toml
[agents.gemini]
command = "GEMINI_SYSTEM_MD='{role_file}' gemini --model {model} --prompt '{prompt}'"
```

Temporary file handling for UTD roles:

1. Create temp file at `/tmp/start-role-{random}.md`
2. Write fully resolved content (after UTD processing)
3. Set permissions to 0600 (owner read/write only)
4. Replace `{role_file}` with temp file path
5. Execute agent (reads from temp file)
6. Delete temp file after agent completes

Simple role example (no temp file):

```toml
[roles.reviewer]
file = "~/.config/start/roles/reviewer.md"
```

`{role_file}` → `/Users/username/.config/start/roles/reviewer.md`

UTD role example (creates temp file):

```toml
[roles.dynamic]
file = "base.md"
command = "date"
prompt = "{file}\n\nDate: {command}"
```

`{role_file}` → `/tmp/start-role-a8f3b2c1.md` (contains: base.md content + date output)

## Scope and Merge Behavior

Allowed scopes:

- Global: `~/.config/start/roles.toml`
- Local: `./.start/roles.toml`

Merge behavior:

- Global + local roles are combined
- Local role replaces global role with same name (no field merging)
- All roles from both configs available for selection

Example:

```toml
# Global: ~/.config/start/roles.toml
[roles.code-reviewer]
file = "~/.config/start/roles/general-reviewer.md"

# Local: ./.start/roles.toml
[roles.code-reviewer]
file = "./PROJECT_REVIEWER.md"
```

Result: Local completely replaces global. Final role uses `./PROJECT_REVIEWER.md`.

## Task Integration

Tasks can specify a role using the `role` field:

```toml
[tasks.security-audit]
role = "security-auditor"
description = "Security-focused audit"
command = "git diff --staged"
prompt = "Audit: {command}"
```

Role selection for tasks follows same precedence rules.

## Validation

At configuration load time:

- Verify role has at least one of: `file`, `command`, or `prompt` (UTD requirement)
- Role names must match pattern: `/^[a-z0-9]+(-[a-z0-9]+)*$/`
- `default_role` must reference an existing role (if specified)

At execution time:

- Selected role must exist in merged config
- If role has `file` field, file must exist (or warning + empty content)
- Task `role` field must reference existing role (if specified)

## Usage Examples

Simple role (file only):

```toml
[roles.general-assistant]
description = "General purpose AI assistant"
file = "~/.config/start/roles/general.md"
```

Role with template framing:

```toml
[roles.code-reviewer]
description = "Code reviewer with additional context"
file = "~/.config/start/roles/reviewer.md"
prompt = """
{file}

Additional Instructions:
- Focus on security and performance
- Check for edge cases
- Verify error handling
"""
```

Role with dynamic content:

```toml
[roles.go-expert]
description = "Go language expert"
file = "~/.config/start/roles/go-base.md"
command = "go version 2>/dev/null || echo 'Go not installed'"
prompt = """
{file}

Environment: {command}

Apply Go-specific best practices and idioms.
"""
```

Inline role (no file):

```toml
[roles.documentation-writer]
description = "Technical documentation specialist"
prompt = """
You are a technical documentation specialist.

Guidelines:
- Use clear, concise language
- Include code examples
- Focus on user needs
- Use active voice and present tense
"""
```

Project-specific role (local config):

```toml
# ./.start/roles.toml
[roles.project-reviewer]
description = "Project-specific reviewer"
file = "./ROLE.md"
```

Command line usage:

```bash
# Use default role
start

# Use specific role
start --role security-auditor

# Use role with custom prompt
start --role code-reviewer prompt "Review authentication code"

# Task inherits role
start task code-review

# Override task's role
start task code-review --role go-expert
```

## Alternatives

Single `[system_prompt]` section:

```toml
[system_prompt]
file = "ROLE.md"
```

- Pro: Simpler - only one system prompt
- Pro: Fewer concepts to learn
- Con: No reusability across tasks
- Con: No selection flexibility
- Con: Cannot have different roles for different tasks
- Con: Inconsistent with `[agents.<name>]` and `[tasks.<name>]` pattern
- Rejected: Too limiting, inconsistent naming pattern

Single placeholder `{system_prompt}` (content only):

- Pro: Only one placeholder to learn
- Pro: Simpler agent config
- Con: Must create temp files for all agents (wasteful for inline agents)
- Con: Extra I/O overhead for common case (Claude, aichat support inline)
- Con: Less explicit about agent API requirements
- Rejected: Inefficient, hides agent API requirements

Task `system_prompt_*` fields (original design):

```toml
[tasks.review]
system_prompt_file = "..."
system_prompt_command = "..."
system_prompt = "..."
```

- Pro: System prompt defined directly in task
- Pro: No separate role to reference
- Con: Verbose and repetitive
- Con: No reusability (duplicate role definitions across tasks)
- Con: Cannot easily swap roles for experimentation
- Con: Difficult to maintain (change role = update all tasks using it)
- Rejected: Poor reusability, maintenance burden

Agent `[agents.<name>.env]` section for env vars:

```toml
[agents.gemini.env]
GEMINI_SYSTEM_MD = "{role_file}"
```

- Pro: Explicit env var configuration
- Pro: Structured config for environment
- Con: Adds complexity (new config section)
- Con: Standard shell syntax works fine: `VAR=value command`
- Con: Less flexible (can't easily set multiple env vars)
- Con: One more concept to implement and document
- Rejected: Standard shell syntax is simpler and well-understood

## Breaking Changes from Original Design

This is a complete redesign of DR-005. Changes from original:

1. Removed: `[system_prompt]` section
2. Added: `[roles.<name>]` sections with UTD
3. Added: `default_role` setting
4. Added: `{role}` and `{role_file}` placeholders
5. Removed: `{system_prompt}` placeholder
6. Added: Role selection precedence rules
7. Added: `--role` CLI flag
8. Added: `role` field for tasks
9. Removed: Task `system_prompt_*` fields (moved to role reference)
10. Removed: `[agents.<name>.env]` sections

## Updates

- 2025-01-07: Complete redesign - named roles with UTD, two placeholders, precedence rules
