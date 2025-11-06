# DR-010: Default Task Definitions

**Date:** 2025-01-03
**Status:** Accepted
**Category:** Tasks

## Decision

Ship four interactive review tasks as defaults

## Default Tasks

1. **code-review** (alias: `cr`) - General code review
2. **git-diff-review** (alias: `gdr`) - Review git diff output
3. **comment-tidy** (alias: `ct`) - Review and tidy code comments
4. **doc-review** (alias: `dr`) - Review and improve documentation

## Rationale

- All tasks are **interactive reviews** - user works with agent in chat
- No tasks that write files or require orchestration (commit messages, gitignore generation)
- Stays true to vision: launcher only, not a workflow orchestrator
- Users can add non-interactive tasks in their own config if desired

## Tasks NOT Included

- **commit-message** - Requires committing after generation (orchestration)
- **gitignore** - Requires saving to file (file I/O)
- **update-changelog** - Requires writing to CHANGELOG.md (file I/O)

## User Customization

- Users can override any default task by defining same name in config
- Users can add additional tasks (including non-interactive ones)
- Users can remove defaults by not including them

## Implementation

- Default tasks embedded in binary
- Loaded first, then merged with user config
- User config takes precedence

## Related Decisions

- [DR-009](./dr-009-task-structure.md) - Task structure
- [DR-011](./dr-011-asset-distribution.md) - Asset distribution
- [DR-019](./dr-019-task-loading.md) - Task loading (updated: assets are templates)
