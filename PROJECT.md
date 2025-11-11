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
4. You will review related documents to make sure there are no inconsistencies
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
