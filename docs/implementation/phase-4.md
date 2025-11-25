# Phase 4: UTD Pattern Processing

**Status:** Not Started
**Dependencies:** Phase 3
**Estimated Effort:** 6-8 hours

---

## Required Reading

Before starting this phase, review these documents:

**Design Documents (Critical):**
- [unified-template-design.md](../design/unified-template-design.md) - Complete UTD pattern specification

**Design Records:**
- [DR-007: Placeholders](../design/design-records/dr-007-placeholders.md) - Placeholder system and resolution
- [DR-008: File Handling](../design/design-records/dr-008-file-handling.md) - File path handling and temp files

---

## Goal

Support Unified Template Design for dynamic content generation.

---

## Deliverables

- [ ] UTD file/command/prompt parsing
- [ ] Command execution (shell integration)
- [ ] UTD placeholders ({file}, {file_contents}, {command_output})
- [ ] Shell configuration and timeouts
- [ ] Temporary file handling

---

## Testing Criteria

- [ ] UTD works for roles and contexts
- [ ] Commands execute correctly
- [ ] Placeholders resolve in prompts
- [ ] Temp files cleaned up
- [ ] Command timeouts work

---

_Next: [Phase 5](phase-5.md)_
