# DR-002: Configuration File Structure

**Date:** 2025-01-03
**Status:** Accepted
**Category:** Configuration

## Decision

Single configuration file with global + local merge strategy

## Files

- **Global:** `~/.config/start/config.toml`
- **Local:** `./.start/config.toml` (project-specific)

## Merge Behavior

- Local config merges with global
- Same keys in local override global values
- New keys in local are added
- Omitted keys use global defaults

## Rationale

- Single file simpler than multiple files
- Merge allows both defaults and project-specific overrides
- CLI commands will manage config, so complexity is hidden from users

## Related Decisions

- [DR-001](./dr-001-toml-format.md) - TOML format choice
- [DR-003](./dr-003-named-documents.md) - Named documents enable override pattern
- [DR-004](./dr-004-agent-scope.md) - Agents can be in both scopes
