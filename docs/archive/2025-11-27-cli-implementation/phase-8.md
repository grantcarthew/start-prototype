# Phase 8: Config Management & Doctor

**Status:** Phase 8a Complete, Phase 8b In Progress
**Dependencies:** Phase 7
**Estimated Effort:** 30-35 hours total (split into sub-phases)

---

## Phase Breakdown

This phase is split into focused sub-phases:

### Phase 8a: Diagnostics & CLI UX ✅ **COMPLETE**
- Prefix matching (DR-030)
- Shell completion (DR-028)
- Version checker (DR-021)
- Doctor command (DR-024)
- **Effort:** 3-4 hours
- **Status:** ✅ Complete (2025-11-26)

### Phase 8b: Agent Config Commands ✅ **COMPLETE**
- `start config agent list/new/show/test/edit/remove/default`
- TOML backup/manipulation helpers
- Interactive wizards with validation
- Complete with unit + integration tests
- **Effort:** 6-8 hours (actual: ~7 hours)
- **Status:** ✅ Complete (2025-11-26)

### Phase 8c: Role Config Commands 🔄 **NEXT**
- `start config role list/new/show/test/edit/remove/default`
- Reuse helper functions from 8b
- **Effort:** 4-5 hours
- **Status:** Not Started

### Phase 8d: Context Config Commands
- `start config context list/new/show/test/edit/remove`
- Similar to roles
- **Effort:** 4-5 hours
- **Status:** Not Started

### Phase 8e: Task Config Commands
- `start config task list/new/show/test/edit/remove`
- Most complex (includes role/agent selection)
- **Effort:** 5-6 hours
- **Status:** Not Started

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

## Phase 8b Deliverables

- [ ] TOML backup helper (timestamped backups)
- [ ] TOML manipulation helpers (read/write with preservation)
- [ ] `start config agent list` - Display all configured agents
- [ ] `start config agent new` - Interactive agent creation wizard
- [ ] `start config agent show <name>` - Display agent configuration
- [ ] `start config agent test <name>` - Validate agent configuration
- [ ] `start config agent edit <name>` - Interactive agent editing
- [ ] `start config agent remove <name>` - Remove agent with confirmation
- [ ] `start config agent default <name>` - Set/show default agent
- [ ] Unit tests for all agent commands
- [ ] Integration tests for agent commands

---

## Phase 8b Testing Criteria

- [ ] Agent list shows all agents from global and local configs
- [ ] Agent new wizard creates valid agent configs
- [ ] Agent show displays correct configuration details
- [ ] Agent test validates binary availability and command templates
- [ ] Agent edit wizard modifies existing configs correctly
- [ ] Agent remove creates backups and deletes correctly
- [ ] Agent default sets default_agent in settings
- [ ] All commands handle missing configs gracefully
- [ ] Backup files are created before modifications
- [ ] Interactive prompts work correctly
- [ ] --local flag targets local config
- [ ] Validation catches invalid configurations

---

_Next: [Phase 9](phase-9.md)_
