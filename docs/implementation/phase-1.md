# Phase 1: Config Loading & Validation

**Status:** Not Started
**Dependencies:** Phase 0
**Estimated Effort:** 4-6 hours

See Phase 0 completion before starting. See [PROJECT.md](../../PROJECT.md) for context.

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

- [ ] Can load configs from test fixtures
- [ ] Merge logic correct (local overrides global)
- [ ] Validation catches all error cases
- [ ] `start config show` displays config
- [ ] Missing files handled gracefully

---

_See original PROJECT.md for detailed implementation tasks_
_Next: [Phase 2](phase-2.md)_
