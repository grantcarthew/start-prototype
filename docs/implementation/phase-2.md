# Phase 2: Simple Agent Execution

**Status:** Not Started
**Dependencies:** Phase 1
**Estimated Effort:** 4-6 hours

---

## Required Reading

Before starting this phase, review these documents:

**Design Records:**

- [DR-006: Cobra CLI](../design/design-records/dr-006-cobra-cli.md) - CLI command structure
- [DR-007: Placeholders](../design/design-records/dr-007-placeholders.md) - Placeholder resolution system
- [DR-013: Agent Templates](../design/design-records/dr-013-agent-templates.md) - Agent configuration structure

**CLI Documentation:**

- [start.md](../cli/start.md) - Main command and interactive mode

---

## Goal

Execute an agent with minimal features (no UTD yet, just basic placeholders).

---

## Deliverables

- [ ] Basic prompt assembly
- [ ] Placeholder resolution ({model}, {prompt}, {bin}, {date})
- [ ] Agent command construction
- [ ] Process execution via Runner interface
- [ ] Root command executes agent

---

## Testing Criteria

- [ ] `start "hello world"` executes smith
- [ ] Smith receives correct model and prompt
- [ ] Placeholders resolve correctly
- [ ] Errors handled gracefully

---

_Next: [Phase 3](phase-3.md)_
