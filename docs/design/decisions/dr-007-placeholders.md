# DR-007: Command Interpolation and Placeholders

**Date:** 2025-01-03
**Status:** Accepted
**Category:** Configuration

## Decision

Use single-brace placeholders with specific supported variables

## Supported Placeholders

- `{model}` - Model name (e.g., "claude-3-7-sonnet-20250219")
- `{system_prompt}` - System prompt file contents
- `{prompt}` - Built prompt text from context documents
- `{date}` - Current timestamp (ISO 8601 format with timezone)

## Path Expansion

- `~` - Expands to user's home directory

## Usage Examples

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

## Substitution Behavior

- `{model}` - Replaced with selected model tier value
- `{system_prompt}` - Replaced with file contents (empty string if not configured)
- `{prompt}` - Replaced with assembled prompt text
- `{date}` - Replaced with current timestamp (e.g., "2025-01-03T10:30:00+10:00")
- `~` - Expanded before command execution

## Rationale

- Single braces simpler than double (`{}` vs `{{}}`)
- "system_prompt" is clear and matches standard terminology
- All placeholders optional - agents can use what they need
- Tilde expansion more concise than `{home}` placeholder
- No environment variable substitution (`{env:...}`) - agents inherit environment naturally

## Not Supported

- `{env:VAR}` - Environment variables (use agent `env` section instead)
- `{home}` - Use `~` instead
- `{cwd}` - Use `--directory` flag if needed

## Related Decisions

- [DR-009](./dr-009-task-structure.md) - Task-specific placeholders
