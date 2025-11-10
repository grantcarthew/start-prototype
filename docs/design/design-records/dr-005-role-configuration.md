# DR-005: Role Configuration & Selection

**Date:** 2025-01-03 (original), 2025-01-07 (rewritten)
**Status:** Accepted
**Category:** Configuration

## Decision

Roles (system prompts) are configured as named entities using the Unified Template Design (UTD) pattern, selected via precedence rules, and made available to agents through two placeholders: `{role}` (content) and `{role_file}` (file path).

## Configuration Structure

### Role Definition

Roles use the `[roles.<name>]` section with full UTD support:

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

**Fields:**
- `description` (optional) - Human-readable description
- `file` (optional) - Path to role content file
- `command` (optional) - Shell command for dynamic content
- `prompt` (optional) - Template with `{file}` and `{command}` placeholders

**UTD Requirement:** At least one of `file`, `command`, or `prompt` must be present.

See [DR-007](./dr-007-placeholders.md) for UTD pattern details.

### Default Role

```toml
[settings]
default_role = "code-reviewer"
```

If `default_role` is not specified, the first role defined in the config (TOML order) is used.

### No [system_prompt] Section

The `[system_prompt]` section from earlier designs is **removed**. All system prompts are defined as named roles in `[roles.<name>]` sections.

## Role Selection Precedence

When executing `start`, `start prompt`, or `start task`, the role is selected using this priority order:

1. **CLI `--role` flag** (highest priority)
2. **Task `role` field** (if executing a task)
3. **`default_role` setting**
4. **First role in config** (TOML order)

**Examples:**

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

Roles are made available to agents through two placeholders in agent commands:

### {role} Placeholder

Contains the **fully resolved role content** (inline text).

**Use for:** Agents that accept system prompts inline via command flags or environment variables.

**Resolution:**
- Simple role (file only): Content of the file
- UTD role: Content after full UTD processing (file + command execution + template rendering)

**Example:**

```toml
[agents.claude]
command = "claude --model {model} --append-system-prompt '{role}' '{prompt}'"

[agents.aichat]
command = "AICHAT_SYSTEM_PROMPT='{role}' aichat --model {model} '{prompt}'"
```

### {role_file} Placeholder

Contains a **file path** to the role content.

**Use for:** Agents that require system prompts to be read from a file.

**Resolution:**
- **Simple role (file only):** Returns the original file path (expanded, absolute)
- **UTD role (complex):** Creates a temporary file with resolved content, returns temp file path

**Example:**

```toml
[agents.gemini]
command = "GEMINI_SYSTEM_MD='{role_file}' gemini --model {model} --prompt '{prompt}'"
```

### Temporary File Handling

When `{role_file}` is used with a UTD role (one that uses `command` or complex `prompt` processing):

1. **Before execution:** Create temp file at `/tmp/start-role-{random}.md`
2. **Write content:** Fully resolved role content (after UTD processing)
3. **Set permissions:** 0600 (owner read/write only)
4. **Replace placeholder:** `{role_file}` → temp file path
5. **Execute agent:** Agent reads from temp file
6. **Cleanup:** Delete temp file after agent completes (success or failure)

**Simple role example (no temp file):**

```toml
[roles.reviewer]
file = "~/.config/start/roles/reviewer.md"
```

`{role_file}` → `/Users/username/.config/start/roles/reviewer.md`

**UTD role example (creates temp file):**

```toml
[roles.dynamic]
file = "base.md"
command = "date"
prompt = "{file}\n\nDate: {command}"
```

`{role_file}` → `/tmp/start-role-a8f3b2c1.md` (contains: base.md content + date output)

## Scope and Merge Behavior

**Allowed scopes:**
- Global: `~/.config/start/roles.toml`
- Local: `./.start/roles.toml`

**Merge behavior:**
- Global + local roles are **combined**
- Local role **replaces** global role with same name (no field merging)
- All roles from both configs available for selection

**Example:**

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

**Role selection for tasks:**
1. CLI `--role` flag
2. Task `role` field
3. `default_role` setting
4. First role in config

See [DR-009](./dr-009-task-structure.md) for task structure details.

## Validation

**At configuration load time:**
- Verify role has at least one of: `file`, `command`, or `prompt` (UTD requirement)
- Role names must match pattern: `/^[a-z0-9]+(-[a-z0-9]+)*$/`
- `default_role` must reference an existing role (if specified)

**At execution time:**
- Selected role must exist in merged config
- If role has `file` field, file must exist (or warning + empty content)
- Task `role` field must reference existing role (if specified)

**Validation commands:**
- `start doctor` - Checks role configuration validity
- `start config validate` - Validates all config including roles

## Examples

### Simple Role (File Only)

```toml
[roles.general-assistant]
description = "General purpose AI assistant"
file = "~/.config/start/roles/general.md"
```

### Role with Template Framing

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

### Role with Dynamic Content

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

### Inline Role (No File)

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

### Project-Specific Role (Local Config)

```toml
# ./.start/config.toml
[roles.project-reviewer]
description = "Project-specific reviewer"
file = "./ROLE.md"
# Relative path resolves from working directory
```

### Role Usage Examples

```bash
# Use default role
start

# Use specific role
start --role security-auditor

# Use role with custom prompt
start --role code-reviewer prompt "Review authentication code"

# Task inherits role
start task code-review  # Uses task's role field

# Override task's role
start task code-review --role go-expert
```

## Rationale

**Why named roles instead of [system_prompt] section:**

1. **Consistency:** Matches pattern of `[agents.<name>]` and `[tasks.<name>]`
2. **Reusability:** Define once, use in multiple tasks
3. **Mix-and-match:** Easy to swap roles via `--role` flag
4. **Clear mental model:** "Role" is what the AI is, easier to understand than "system prompt"
5. **Selection flexibility:** Precedence rules allow override at multiple levels

**Why two placeholders ({role} and {role_file}):**

1. **Agent compatibility:** Some agents accept inline content, others require files
2. **Efficiency:** Agents with file-based APIs (Gemini) can read original file directly
3. **Transparency:** Clear which pattern each agent uses
4. **Flexibility:** Supports both content-based and file-based agent designs

**Why temp files for UTD roles with {role_file}:**

1. **Consistency:** `{role_file}` always returns a file path, regardless of role complexity
2. **Agent compatibility:** File-based agents work with both simple and complex roles
3. **Clean abstraction:** Agent config doesn't need to know if role is simple or complex

**Why remove [agents.<name>.env] section:**

1. **Simplicity:** Standard shell syntax `VAR=value command` is universally understood
2. **Fewer concepts:** One less config section to implement and document
3. **Flexibility:** Can set multiple env vars easily: `VAR1=a VAR2=b command`
4. **No special handling:** No need for separate env var processing logic

## Alternatives Considered

**Alternative 1: Keep [system_prompt] section**

```toml
[system_prompt]
file = "ROLE.md"
```

Rejected: Inconsistent with named pattern used for agents and tasks. No reusability, no selection flexibility.

**Alternative 2: Single placeholder {system_prompt}**

Use only `{system_prompt}` with content, create temp files implicitly for all agents.

Rejected: Wasteful for agents that accept inline content (Claude, aichat). Extra I/O for common case.

**Alternative 3: Keep [agents.<name>.env] section**

```toml
[agents.gemini.env]
GEMINI_SYSTEM_MD = "{role_file}"
```

Rejected: Adds complexity. Standard shell syntax works fine and is more familiar.

**Alternative 4: Task system_prompt_* fields (original design)**

```toml
[tasks.review]
system_prompt_file = "..."
system_prompt_command = "..."
system_prompt = "..."
```

Rejected: Verbose, duplicative. Defining reusable roles is cleaner. Changed in this DR.

## Breaking Changes from Original Design

This is a **complete redesign** of DR-005. Changes from original:

1. **Removed:** `[system_prompt]` section
2. **Added:** `[roles.<name>]` sections with UTD
3. **Added:** `default_role` setting
4. **Added:** `{role}` and `{role_file}` placeholders
5. **Removed:** `{system_prompt}` placeholder
6. **Added:** Role selection precedence rules
7. **Added:** `--role` CLI flag
8. **Added:** `role` field for tasks
9. **Removed:** Task `system_prompt_*` fields (moved to role reference)
10. **Removed:** `[agents.<name>.env]` sections

## Related Decisions

- [DR-002](./dr-002-config-merge.md) - Config merge behavior (roles follow same rules)
- [DR-003](./dr-003-named-documents.md) - Named sections pattern (roles use same pattern)
- [DR-004](./dr-004-agent-scope.md) - Agent scope rules (roles follow same scope rules)
- [DR-007](./dr-007-placeholders.md) - Placeholder system and UTD pattern
- [DR-009](./dr-009-task-structure.md) - Task structure (tasks reference roles)
- [DR-019](./dr-019-task-loading.md) - Task loading (role selection during task execution)
- [DR-029](./dr-029-task-agent-field.md) - Task agent field (parallel to role field)

## Implementation Notes

**Role Resolution Steps:**

1. Determine selected role (precedence rules)
2. Load role configuration from merged config
3. Resolve UTD fields (file, command, prompt)
4. Generate `{role}` content (fully resolved)
5. Generate `{role_file}` path:
   - Simple role: original file path (expanded)
   - UTD role: create temp file, return temp path
6. Replace placeholders in agent command
7. Execute agent
8. Cleanup temp files

**Error Handling:**

- Missing role file: Warning + use empty content
- Invalid role name: Error at config load
- Task references non-existent role: Error at execution
- `default_role` references non-existent role: Error at config load
- Temp file creation fails: Error at execution

**Performance Considerations:**

- Role content cached after first resolution (within single execution)
- Temp files only created when needed (`{role_file}` + UTD role)
- File I/O minimized (read once, cache content)
- Temp file cleanup in defer/finally blocks (guaranteed cleanup)
