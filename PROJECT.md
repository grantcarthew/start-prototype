# Project: Design Record Schema Alignment

**Status:** In Progress - HIGHLY SENSITIVE
**Date Started:** 2025-01-17

## CRITICAL IMPORTANCE

This project is HIGHLY SENSITIVE and requires UTMOST CARE.

We are aligning all Design Records (DRs) with a new standardized schema. These documents capture the fundamental architectural decisions for the entire project. Any mistakes, omissions, or misalignment could:

- Lose critical design rationale
- Create inconsistencies that break the design
- Invalidate hours of careful decision-making
- Require expensive rework during implementation

PROCEED WITH EXTREME CAUTION. When in doubt, ask questions. It is better to pause and clarify than to proceed incorrectly.

## Objectives

Align all 37 Design Records in `docs/design/design-records/` with the new DR schema defined in `docs/design/dr-writing-guide.md`.

Critical requirements:

1. PRESERVE ALL DESIGN DECISIONS - Never lose important design information
2. Align with new schema structure (Problem, Decision, Why, Trade-offs, Alternatives)
3. Remove inappropriate content (implementation code, cross-links, duplication)
4. Fix any inconsistencies or outdated information discovered during alignment
5. Maintain document quality and completeness

## Reference Documents

MUST READ before starting work:

1. `docs/design/dr-writing-guide.md` - Complete DR writing guidelines and schema
   - Defines required sections: Problem, Decision, Why, Trade-offs, Alternatives
   - Defines optional sections: Structure, Scope, Validation, Usage Examples, etc.
   - Lists what belongs (config examples, field descriptions) and what doesn't (implementation code, cross-links)
   - Provides when to create a DR and writing principles

## Alignment Guidelines

### Required Sections (every DR must have)

1. Problem - What constraint or issue drove this decision?
2. Decision - Clear, specific statement of what was decided
3. Why - Core reasoning behind this choice
4. Trade-offs - Accept (costs) and Gain (benefits)
5. Alternatives - Other options considered with rejection reasoning

### Optional Sections (add as needed)

- Structure - Schema definitions, field descriptions
- Scope - Where/how this applies (global vs local)
- Validation - Rules for correctness
- Usage Examples - How to use the decision in practice
- Execution Flow - Step-by-step behavior
- Breaking Changes - Updates from previous versions
- Updates - Historical changes with dates

### What to KEEP

✅ Configuration examples (TOML, JSON schemas)
✅ Usage examples (bash commands, CLI usage)
✅ Field descriptions and constraints
✅ Validation rules
✅ Execution flows and algorithms
✅ Tables and matrices
✅ All design decisions and rationale
✅ Trade-off analysis
✅ Alternative approaches with rejection reasoning

### What to REMOVE

❌ Implementation code (Go, Python, etc.) - Exception: pseudocode in Why section for complex logic
❌ Cross-links to other DRs - Exception: status changes (Superseded by DR-XXX)
❌ User documentation duplication
❌ "Related Decisions" sections
❌ Step-by-step implementation instructions
❌ Low-level code guidance

### Formatting Rules

1. NO BOLD FORMATTING - Do not use `**text**` anywhere in DRs
2. Header format: `# DR-NNN: Title` with metadata bullets (Date, Status, Category)
3. Use plain text with structure via headings and lists
4. Code blocks for examples (toml, bash, etc.)
5. Clear, concise language

### Process for Each DR

1. Read the existing DR completely
2. Identify what needs to change:
   - Missing required sections (Problem, Trade-offs, etc.)
   - Sections to rename (Rationale → Why)
   - Content to remove (Related Decisions, implementation code)
   - Inconsistencies or outdated information
3. If unclear or inconsistencies found: STOP and ask the user
4. Restructure into new schema while preserving all design decisions
5. Write the aligned DR
6. Mark as completed in checklist below

## Inconsistencies Found and Fixed

During alignment, we discovered and fixed these issues:

1. DR-002: Said "Single configuration file" but system uses multi-file structure (settings.toml, agents.toml, tasks.toml, roles.toml, contexts.toml) - FIXED
2. DR-004: Referenced config.toml instead of agents.toml - FIXED
3. DR-004: Included env variable substitution feature that doesn't exist in design - REMOVED
4. DR-007: Had {model} placeholder available universally (agents, roles, contexts, tasks) but no legitimate use case outside agent commands - FIXED to agent-commands-only

Continue to flag any inconsistencies discovered during alignment.

## Design Records Checklist

Location: `docs/design/design-records/`

Progress: 8 of 37 completed (22%)

- [x] DR-001: TOML Format
- [x] DR-002: Config Merge
- [x] DR-003: Named Documents
- [x] DR-004: Agent Scope
- [x] DR-005: Role Configuration
- [x] DR-006: Cobra CLI
- [x] DR-007: Placeholders
- [x] DR-008: File Handling
- [ ] DR-009: Task Structure
- [ ] DR-010: Default Tasks
- [ ] DR-011: Asset Distribution
- [ ] DR-012: Context Required
- [ ] DR-013: Agent Templates
- [ ] DR-017: CLI Reorganization
- [ ] DR-018: Init Update Integration
- [ ] DR-019: Task Loading
- [ ] DR-020: Version Injection
- [ ] DR-021: GitHub Version Check
- [ ] DR-022: Asset Branch Strategy
- [ ] DR-024: Doctor Exit Codes
- [ ] DR-025: No Automatic Checks
- [ ] DR-026: Offline Behavior
- [ ] DR-027: Security Trust Model
- [ ] DR-028: Shell Completion
- [ ] DR-029: Task Agent Field
- [ ] DR-030: Prefix Matching
- [ ] DR-031: Catalog-Based Assets
- [ ] DR-032: Asset Metadata Schema
- [ ] DR-033: Asset Resolution Algorithm
- [ ] DR-034: GitHub Catalog API
- [ ] DR-035: Interactive Browsing
- [ ] DR-036: Cache Management
- [ ] DR-037: Asset Updates
- [ ] DR-038: Flag Value Resolution
- [ ] DR-039: Catalog Index
- [ ] DR-040: Substring Matching
- [ ] DR-041: Asset Command Reorganization

## Notes for Continuation

### Working Context

The start project is an AI agent CLI orchestrator. Key concepts:

- Agents: AI CLI tools (claude, gemini, gpt, etc.)
- Roles: System prompts defining agent behavior
- Tasks: Reusable workflow definitions
- Contexts: Environment/project information loaded at runtime
- Assets: Downloadable catalog from GitHub (roles, tasks, agents)
- Multi-file config: 5 TOML files per scope (global ~/.config/start/, local ./.start/)

### File Structure

```
~/.config/start/          # Global config
  settings.toml
  agents.toml
  tasks.toml
  roles.toml
  contexts.toml

./.start/                 # Local config (project-specific)
  settings.toml
  agents.toml
  tasks.toml
  roles.toml
  contexts.toml
```

### Key Design Patterns

1. UTD (Unified Template Design): Pattern with file, command, prompt fields for dynamic content
2. Named sections: `[agents.<name>]`, `[roles.<name>]`, `[tasks.<name>]`, `[context.<name>]`
3. Merge behavior: Local completely replaces global for same name (no per-field merge)
4. Asset resolution: local → global → cache → GitHub
5. Precedence rules: CLI flags → task fields → settings → defaults

### Important Terminology

- "Model name" not "model alias" (user-defined friendly names)
- "Full model identifier" not "full model name" (provider's actual model ID)
- `contexts.toml` file contains `[context.<name>]` sections (file plural, section singular)
- Multi-file config, not single-file

## Success Criteria

Alignment is complete when:

- [ ] All 37 DRs restructured with required sections (Problem, Decision, Why, Trade-offs, Alternatives)
- [ ] All Related Decisions sections removed
- [ ] All cross-references to other DRs removed (except Superseded links)
- [ ] All implementation code removed (config examples kept)
- [ ] All bold formatting removed
- [ ] All design decisions preserved
- [ ] All inconsistencies fixed
- [ ] Each DR is complete and accurate

## Work Session Instructions

When continuing this work:

1. Read this PROJECT-dr-alignment.md file completely
2. Read docs/design/dr-writing-guide.md
3. Check the checklist above for next DR to align
4. Read the DR completely before making changes
5. Flag any inconsistencies or unclear content
6. Align the DR following the guidelines
7. Update the checklist in this file
8. Add summary of changes to "Completed DRs" section above
9. Continue to next DR

REMEMBER: This is HIGHLY SENSITIVE work. Preserve all design decisions. When in doubt, ask.
