# Project: CLI Documentation Review and Correction

**Status:** Active - Documentation Review
**Date Started:** 2025-01-11
**Previous Phase:** Catalog Design (completed, see PROJECT-backlog.md)

## Overview

Comprehensive review and correction of all CLI documentation in `docs/cli/` to ensure accuracy, consistency, and alignment with design decisions before implementation begins.

## Important

I am reviewing the documents. Your task is to fix issues defined by me during this review process. It is an interactive session.

We will be working through this step by step. If you see something that should be addressed, bring it up in the discussion.

As we fix issues, I will be needing you to ensure the docs/design/design-records/* document that relates to the changes is in-sync with our updates.

## Objectives

1. Review all 11 CLI command documentation files for accuracy
2. Fix inconsistencies, outdated information, and errors
3. Ensure alignment with design decisions (DRs)
4. Verify examples, flags, and usage patterns are correct
5. Prepare clean, accurate documentation for implementation phase

## CLI Documents to Review

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

- [ ] All 11 CLI documents reviewed and corrected
- [ ] Examples are accurate and tested conceptually
- [ ] Flags and options are consistent across commands
- [ ] Multi-file config structure correctly documented
- [ ] Catalog behavior (lazy loading, browsing) accurately described
- [ ] Design decisions referenced where relevant
- [ ] No contradictions between documents
- [ ] Ready for implementation phase

## Notes

- This review may identify gaps requiring new design decisions
- Some issues may require updates to design documents (docs/design/)
- Focus on correctness over completeness - better to have accurate docs than comprehensive but wrong docs
- Track design questions separately for resolution before implementation
