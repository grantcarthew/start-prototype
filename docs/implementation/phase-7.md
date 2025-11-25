# Phase 7: Init & Asset Management

**Status:** Not Started
**Dependencies:** Phase 6
**Estimated Effort:** 8-10 hours

---

## Required Reading

Before starting this phase, review these documents:

**Design Records:**
- [DR-035: Interactive Browsing](../design/design-records/dr-035-interactive-browsing.md) - Interactive asset browsing
- [DR-037: Asset Updates](../design/design-records/dr-037-asset-updates.md) - Update mechanism
- [DR-040: Substring Matching](../design/design-records/dr-040-substring-matching.md) - Search algorithm
- [DR-041: Asset Command Reorganization](../design/design-records/dr-041-asset-command-reorganization.md) - Command structure

**CLI Documentation:**
- [start-init.md](../cli/start-init.md) - Init wizard specification
- [start-assets-browse.md](../cli/start-assets-browse.md) - Browse command
- [start-assets-info.md](../cli/start-assets-info.md) - Info command
- [start-assets-index.md](../cli/start-assets-index.md) - Index command
- [start-assets-update.md](../cli/start-assets-update.md) - Update command

---

## Goal

Setup wizard and full asset tooling.

---

## Deliverables

- [ ] `start init` interactive wizard
- [ ] Agent auto-detection
- [ ] Default config generation
- [ ] `start assets browse/update/info`
- [ ] Asset index generation

---

## Testing Criteria

- [ ] `start init` works (interactive and --force)
- [ ] All asset commands work
- [ ] Config generation correct
- [ ] Agent detection accurate

---

_Next: [Phase 8](phase-8.md)_
