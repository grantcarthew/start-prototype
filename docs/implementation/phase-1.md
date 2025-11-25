# Phase 1: Config Loading & Validation

**Status:** Not Started
**Dependencies:** Phase 0
**Estimated Effort:** 4-6 hours

See Phase 0 completion before starting. See [PROJECT.md](../../PROJECT.md) for context.

---

## Required Reading

Before starting this phase, review these documents:

**Design Records (Critical):**
- [DR-001: TOML Format](../design/design-records/dr-001-toml-format.md) - TOML configuration format and strict mode
- [DR-002: Config Merge](../design/design-records/dr-002-config-merge.md) - Global + local merge strategy
- [DR-003: Named Documents](../design/design-records/dr-003-named-documents.md) - Named document sections (not arrays)
- [DR-004: Agent Scope](../design/design-records/dr-004-agent-scope.md) - Agents in global and local configs
- [DR-005: Role Configuration](../design/design-records/dr-005-role-configuration.md) - Role config structure
- [DR-008: File Handling](../design/design-records/dr-008-file-handling.md) - Path resolution and missing files

**Design Records (Context):**
- [DR-007: Placeholders](../design/design-records/dr-007-placeholders.md) - Understanding placeholder fields
- [DR-012: Context Required](../design/design-records/dr-012-context-required.md) - Required field behavior

**CLI Documentation:**
- [start-config.md](../cli/start-config.md) - Config command overview
- [start-show.md](../cli/start-show.md) - Show command specification

**Reference:**
- [docs/config.md](../config.md) - Complete configuration reference
- [examples/](../../examples/) - Use all examples as test fixtures (minimal, complete, real-world)

---

## Goal

Load and merge TOML configurations correctly with validation.

---

## Deliverables

- [ ] Config loader (global + local)
- [ ] TOML parsing with go-toml/v2
- [ ] Config merge logic (DR-002)
- [ ] Validation (required fields, patterns)
- [ ] `start config show` command

---

## Testing Criteria

- [ ] Can load configs from test fixtures (use examples/ directory)
- [ ] Merge logic correct (local overrides global) - test with complete/ and real-world/
- [ ] Validation catches all error cases
- [ ] `start config show` displays config - test with minimal/ example
- [ ] Missing files handled gracefully

---

_See original PROJECT.md for detailed implementation tasks_
_Next: [Phase 2](phase-2.md)_
