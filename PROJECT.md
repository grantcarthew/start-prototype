# Project: CLI Documentation Review and Correction

**Status:** Active - Documentation Review
**Date Started:** 2025-01-11
**Previous Phase:** Catalog Design (completed, see PROJECT-backlog.md)

## Overview

Comprehensive review and correction of all CLI documentation in `docs/cli/` to ensure accuracy, consistency, and alignment with design decisions before implementation begins.

## Important

I am reviewing the documents. Your task is to fix issues defined by me during this review process. It is an interactive session.

We will be working through this step by step. If you see something that should be addressed, bring it up in the discussion.

As we fix issues, I will be needing you to ensure the docs/design/design-records/\* document that relates to the changes is in-sync with our updates.

This repository is not large yet. Before you do anything, run a `lsd --tree` to get a list of the documents.

## Objectives

1. Review all documentation files (3 root docs + 11 CLI commands) for accuracy
2. Fix inconsistencies, outdated information, and errors
3. Ensure alignment with design decisions (DRs)
4. Verify examples, flags, and usage patterns are correct
5. Prepare clean, accurate documentation for implementation phase
6. Design records will be updated as needed based on changes to other documents

## Review Process

1. I will identify and issue
2. We will discuss it and decide on the fix
3. You will fix the issue in the document
4. You will review related documents in docs/cli/ and docs/design/design-records/ to make sure there are no inconsistencies
5. You will update related documents
6. Add the fix to the bottom of this document in the ## Fixed section
7. Ask me to commit the changes
8. Next issue

## Documents to Review

Active document: `docs/cli/start.md`

### Root Documentation

- [ ] `docs/config.md` - Configuration reference
- [ ] `docs/tasks.md` - Task-specific documentation
- [ ] `docs/vision.md` - Product vision and goals

### Main Commands

- [ ] `docs/cli/start.md` - Main entry point, interactive sessions
- [ ] `docs/cli/start-prompt.md` - Prompt composition and execution
- [ ] `docs/cli/start-task.md` - Task execution (needs lazy loading updates)
- [ ] `docs/cli/start-init.md` - Configuration initialization
- [ ] `docs/cli/start-update.md` - Asset updates (needs per-asset update model)
- [ ] `docs/cli/start-doctor.md` - System diagnostics

### Configuration Commands

- [ ] `docs/cli/start-config.md` - Configuration management overview
- [ ] `docs/cli/start-config-agent.md` - Agent configuration
- [ ] `docs/cli/start-config-context.md` - Context configuration
- [ ] `docs/cli/start-config-role.md` - Role configuration
- [ ] `docs/cli/start-config-task.md` - Task configuration

## Design Alignment

Ensure all documentation aligns with these key design decisions:

- **DR-031**: Catalog-based asset architecture
- **DR-032**: Asset metadata schema (.meta.toml files)
- **DR-033**: Asset resolution (local → global → cache → GitHub)
- **DR-034**: GitHub API strategy (Tree API + raw.githubusercontent.com)
- **DR-035**: Interactive browsing (numbered selection)
- **DR-036**: Cache management (invisible, manual delete only)
- **DR-037**: Update mechanism (manual, SHA-based)

## Success Criteria

Documentation review is complete when:

- [ ] All 14 documents (3 root + 11 CLI) reviewed and corrected
- [ ] Examples are accurate and tested conceptually
- [ ] Flags and options are consistent across commands
- [ ] Multi-file config structure correctly documented
- [ ] Catalog behavior (lazy loading, browsing) accurately described
- [ ] Design decisions referenced where relevant
- [ ] No contradictions between documents
- [ ] Design records updated to reflect changes
- [ ] Ready for implementation phase

## Notes

- This review may identify gaps requiring new design decisions
- Some issues may require updates to design documents (docs/design/)
- Focus on correctness over completeness - better to have accurate docs than comprehensive but wrong docs
- Track design questions separately for resolution before implementation

## Fixed

- `docs/cli/start.md`: Removed "alias" terminology from --model flag, clarified resolution as exact match → prefix match → passthrough
- `docs/cli/start-prompt.md`: Updated --model flag to match start.md (removed "alias|name", added resolution order)
- **Short flags added**: Added `-a` (--agent), `-r` (--role), `-m` (--model) short flags across all CLI docs
- **Version flag corrected**: Changed `-v` from --verbose to --version across all CLI docs (--verbose has no short form)
- `docs/cli/start.md`: Added short flags -a, -r, -m; moved -v to --version; removed -v from --verbose
- `docs/cli/start-prompt.md`: Added short flags -a, -r, -m; added --version with -v
- `docs/cli/start-config.md`: Moved -v to --version; removed from --verbose
- `docs/cli/start-config-agent.md`: Moved -v to --version; removed from --verbose
- `docs/cli/start-config-task.md`: Moved -v to --version; removed from --verbose
- `docs/cli/start-config-role.md`: Moved -v to --version; removed from --verbose
- `docs/cli/start-config-context.md`: Moved -v to --version; removed from --verbose
- `docs/design/design-records/dr-006-cobra-cli.md`: Updated Global Flags section with complete flag list including short forms
- `docs/design/design-records/dr-028-shell-completion.md`: Updated flag completion examples to include all short flags
- **--context flag**: Decided not to add --context/-c flag; current config-driven design is sufficient
- **Missing file behavior**: Standardized wording across all docs - missing context files always generate warnings and are skipped
- `docs/cli/start.md`: Updated context behavior description and execution flow; changed output examples from ✗ to ⚠ for missing files
- `docs/cli/start-config.md`: Changed from "silently skipped" to "generate warnings and are skipped"
- `docs/cli/start-config-context.md`: Updated test output example for missing files
- `docs/cli/start-task.md`: Updated execution flow for missing file handling
- **Runtime behavior**: Missing files show `⚠ context-name file-path (not found, skipped)` - warnings, not errors
- `docs/design/design-records/dr-008-file-handling.md`: Updated to reflect warning behavior instead of "silently skipped"; updated output examples to use ⚠ symbol; added rationale about catching config errors
- **Flag value prefix matching**: Implemented intelligent prefix matching for --agent, --role, --task, and --model flags with ambiguity detection and TTY-aware interactive selection
- `docs/design/design-records/dr-038-flag-value-resolution.md`: Created new DR defining two-phase resolution (exact → prefix), short-circuit evaluation, ambiguity handling (interactive/error), and passthrough for --model
- `docs/cli/start.md`: Updated --agent, --role, and --model flag descriptions with full resolution algorithm and examples
- `docs/cli/start-prompt.md`: Updated --agent, --role, and --model flag descriptions to reference DR-038
- `docs/cli/start-task.md`: Updated task resolution section and name argument to describe prefix matching behavior
- `docs/design/design-records/dr-033-asset-resolution-algorithm.md`: Added DR-038 to related decisions (prefix matching extends exact match)
- **UTD Placeholder Design**: Enhanced placeholder system with four new placeholders for better flexibility
- New placeholders: `{file}` (path), `{file_contents}` (contents), `{command}` (string), `{command_output}` (output)
- Pattern: Short form = reference/source, long form = result/contents
- `docs/design/unified-template-design.md`: Updated all placeholder definitions and examples with new four-placeholder system
- `docs/design/design-records/dr-007-placeholders.md`: Updated UTD Pattern Placeholders section with detailed descriptions and use cases for all four placeholders
- `docs/config.md`: Updated all UTD field descriptions and examples throughout roles, contexts, and tasks sections
- `docs/tasks.md`: Updated task prompt placeholders and all code examples
- `docs/cli/start-prompt.md`: Examples already correct (using {file} for paths)
- `docs/cli/start-task.md`: Updated execution flow and placeholder documentation
- `docs/cli/start-config-context.md`: Updated field descriptions and placeholder section
- `docs/cli/start-config-role.md`: Updated field descriptions and placeholder section
- `docs/cli/start-config-task.md`: Updated field descriptions and placeholder section
- `docs/design/design-records/dr-009-task-structure.md`: Updated decision summary, field descriptions, placeholder definitions, execution flow, and all examples
- **Contexts as lazy-loadable assets**: Confirmed contexts can be lazy-loaded from GitHub catalog (config templates, not content)
- `docs/design/design-records/dr-031-catalog-based-assets.md`: Added `contexts/` to cache structure; removed contexts from Future Considerations
- `docs/config.md`: Updated Asset Resolution & Lazy-Loading section to include contexts as downloadable asset type
- `docs/design/design-records/dr-032-asset-metadata-schema.md`: Added contexts to asset types list
- `docs/design/design-records/dr-022-asset-branch-strategy.md`: Added contexts to asset types in three locations (content list, benefits, rationale)
- **Command name corrections**: Fixed all instances of `start agent` to `start config agent` (there is no `start agent` command)
- `docs/config.md`: Changed `start agent list` to `start config agent list`
- `docs/cli/start.md`: Changed `start agent list` to `start config agent list`
- `docs/cli/start-doctor.md`: Changed `start agent remove` to `start config agent remove` and `start agent test` to `start config agent test`
- `docs/cli/start-config-agent.md`: Changed `start agent list` to `start config agent list` and two instances of `start agent remove` to `start config agent remove`
- **Local-only config support**: Fixed documentation to reflect that local config can be used without global config
- `docs/cli/start-config.md`: Changed "Agents: (none - agents must be defined in global config)" to "Agents: (none configured)" (line 178)
- `docs/cli/start-config.md`: Fixed Configuration Merge Behavior section - agents, roles, and tasks can all be in both global and local configs with merge behavior (lines 889-893)
- `docs/cli/start.md`: Fixed "Local config only" section - removed error message, clarified that local-only config is valid (no global config required)
