# Phase 3: Roles & Contexts

**Status:** Not Started
**Dependencies:** Phase 2
**Estimated Effort:** 5-7 hours

---

## Required Reading

Before starting this phase, review these documents:

**Design Records:**

- [DR-005: Role Configuration](../design/design-records/dr-005-role-configuration.md) - Role config and selection precedence
- [DR-003: Named Documents](../design/design-records/dr-003-named-documents.md) - Named document structure
- [DR-012: Context Required](../design/design-records/dr-012-context-required.md) - Required field and document order
- [DR-008: File Handling](../design/design-records/dr-008-file-handling.md) - Path resolution and missing files

**Design Documents:**

- [unified-template-design.md](../design/unified-template-design.md) - UTD pattern (file/command/prompt)

---

## Goal

Add role system prompts and context document loading.

---

## Deliverables

- [ ] Role loading and selection
- [ ] Context document detection and ordering
- [ ] Role placeholders ({role}, {role_file})
- [ ] File reading and path resolution
- [ ] Required vs optional contexts

---

## Testing Criteria

- [ ] Roles selected correctly (precedence rules)
- [ ] Context documents appear in prompt
- [ ] Order correct (contexts first)
- [ ] Missing files warn but continue
- [ ] {role} placeholder works

---

_Next: [Phase 4](phase-4.md)_
