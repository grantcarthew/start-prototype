# DR-005: System Prompt Handling

**Date:** 2025-01-03
**Status:** Accepted
**Category:** Configuration

## Decision

System prompt configured separately from context documents, and is optional

## Structure

```toml
[system_prompt]
file = "./ROLE.md"

[context.environment]
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
```

## Rationale

- System prompt is conceptually different from context documents
- Passed to agents differently (via `--system-prompt` flag or `{role}` placeholder)
- Separate section makes it clear and allows easy override
- Can be overridden in local config like other context settings

## Optional Behavior

- System prompt section can be missing entirely
- Path can be empty
- Not all AI agents support system prompts
- `start` will skip system prompt handling if not configured or file doesn't exist

## Example Local Override

```toml
# Local ./.start/config.toml
[system_prompt]
file = "~/shared-roles/senior-go-dev.md"
```

## Related Decisions

- [DR-002](./dr-002-config-merge.md) - Local overrides global (replacement)
- [DR-019](./dr-019-task-loading.md) - Tasks use same replacement behavior
