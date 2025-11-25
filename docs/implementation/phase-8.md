# Phase 8: Config Management & Doctor

**Status:** Not Started
**Dependencies:** Phase 7
**Estimated Effort:** 6-8 hours

---

## Required Reading

Before starting this phase, review these documents:

**Design Records:**

- [DR-024: Doctor Exit Codes](../design/design-records/dr-024-doctor-exit-codes.md) - Exit code system
- [DR-025: No Automatic Checks](../design/design-records/dr-025-no-automatic-checks.md) - No auto checks or caching
- [DR-028: Shell Completion](../design/design-records/dr-028-shell-completion.md) - Completion support
- [DR-030: Prefix Matching](../design/design-records/dr-030-prefix-matching.md) - Command prefix matching
- [DR-038: Flag Value Resolution](../design/design-records/dr-038-flag-value-resolution.md) - Flag resolution and prefix matching
- [DR-021: GitHub Version Check](../design/design-records/dr-021-github-version-check.md) - Version checking

**CLI Documentation:**

- [start-doctor.md](../cli/start-doctor.md) - Doctor command specification
- [start-config.md](../cli/start-config.md) - Config management overview
- [start-config-agent.md](../cli/start-config-agent.md) - Agent config commands
- [start-config-role.md](../cli/start-config-role.md) - Role config commands
- [start-config-context.md](../cli/start-config-context.md) - Context config commands
- [start-config-task.md](../cli/start-config-task.md) - Task config commands

---

## Goal

Configuration tooling and diagnostics.

---

## Deliverables

- [ ] `start config agent/role/task/context` commands
- [ ] `start doctor` with health checks
- [ ] Prefix matching
- [ ] Shell completion
- [ ] Version checking

---

## Testing Criteria

- [ ] All config commands work
- [ ] Prefix matching works
- [ ] Doctor identifies issues
- [ ] Shell completion works
- [ ] Version checking works

---

_Next: [Phase 9](phase-9.md)_
