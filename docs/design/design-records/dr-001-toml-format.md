# DR-001: Configuration File Format

**Date:** 2025-01-03
**Status:** Accepted
**Category:** Configuration

## Decision

Use TOML for all configuration files

## Rationale

- Human-readable and editable
- No whitespace sensitivity (unlike YAML)
- Excellent Go support via BurntSushi/toml
- Supports comments and complex nested structures
- Used by similar tools (mise)

## Alternatives Considered

- **YAML:** Too error-prone with whitespace
- **JSON:** No comments, less human-friendly
- **Custom key-value:** Too limited for nested structures

## Related Decisions

- [DR-002](./dr-002-config-merge.md) - Configuration file structure
- [DR-003](./dr-003-named-documents.md) - Named documents pattern
