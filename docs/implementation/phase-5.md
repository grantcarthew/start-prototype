# Phase 5: Tasks

**Status:** Not Started
**Dependencies:** Phase 4
**Estimated Effort:** 4-6 hours

---

## Required Reading

Before starting this phase, review these documents:

**Design Records:**

- [DR-009: Task Structure](../design/design-records/dr-009-task-structure.md) - Task configuration and placeholders
- [DR-019: Task Loading](../design/design-records/dr-019-task-loading.md) - Task loading and merging algorithm
- [DR-029: Task Agent Field](../design/design-records/dr-029-task-agent-field.md) - Optional agent field

**CLI Documentation:**

- [start-task.md](../cli/start-task.md) - Task command specification

---

## Goal

Predefined workflow tasks with {instructions} placeholder.

---

## Deliverables

- [ ] Task configuration loading
- [ ] {instructions} placeholder
- [ ] `start task <name>` command
- [ ] Task resolution (local → global)
- [ ] Alias support
- [ ] Auto-include required contexts

---

## Testing Criteria

- [ ] Tasks execute correctly
- [ ] Instructions passed through
- [ ] Aliases work
- [ ] Role/agent selection works
- [ ] Required contexts included

---

_Next: [Phase 6](phase-6.md)_
