# DR-012: Context Document Required Field and Order

**Date:** 2025-01-04
**Status:** Accepted
**Category:** Configuration

## Decision

Add optional `required` field to context documents to control inclusion behavior; documents appear in config definition order

## Structure

```toml
[context.environment]  # First in prompt
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
required = true    # Always included

[context.project]      # Second in prompt
file = "./PROJECT.md"
prompt = "Read {file}. Respond with summary."
required = false   # Optional (default)
```

## Behavior by Command

- `start` (root) → Includes ALL documents (required + optional)
- `start prompt` → Includes ONLY required documents
- `start task` → Includes documents specified in task's `documents` array

## Default Value

If `required` field is omitted, defaults to `false` (optional document)

## Document Order

- Documents appear in prompt in the order they are defined in config file
- TOML preserves declaration order within sections
- Users control order by arranging config file
- Predictable and explicit - no alphabetical or other automatic sorting
- Consistent across all commands (start, start prompt, tasks)

## Rationale

- `start` provides full context for comprehensive sessions (all documents)
- `start prompt` provides minimal context for focused queries (required only)
- Allows users to designate "essential" vs "nice-to-have" context
- Reduces noise for one-off questions while maintaining critical context
- Tasks maintain full control via explicit `documents` array
- Definition order gives users control over context priority

## Use Cases

- `~/reference/ENVIRONMENT.md` marked required: Always provides user/environment context (first)
- `~/reference/INDEX.csv` marked required: Always provides documentation index (second)
- `AGENTS.md` marked required: Always provides repository overview (third)
- `PROJECT.md` marked optional: Included for full sessions, excluded for quick queries

## Related Decisions

- [DR-003](./dr-003-named-documents.md) - Named documents pattern
- [DR-008](./dr-008-file-handling.md) - File detection and handling
