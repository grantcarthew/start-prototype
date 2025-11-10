# DR-003: Named Documents for Context

**Date:** 2025-01-03
**Status:** Accepted
**Category:** Configuration

## Decision

Use named document sections instead of arrays

## Structure

```toml
[context.environment]
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."

[context.project]
file = "./PROJECT.md"
prompt = "Read {file}. Respond with summary."
```

## Rationale

- Names allow local config to override specific documents
- Can't target array items for override
- Enables both override (same name) and add (new name) patterns
- More explicit and readable

## Example Use Case

- Global defines "project" document as `./PROJECT.md`
- Local overrides to `~/multi-repo/BIG-PROJECT.md`
- Local adds new "vision" document as `./docs/vision.md`

## Related Decisions

- [DR-002](./dr-002-config-merge.md) - Merge strategy this enables
- [DR-012](./dr-012-context-required.md) - Context document required field
