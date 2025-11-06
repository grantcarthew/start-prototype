# DR-015: Atomic Update Mechanism

**Date:** 2025-01-06
**Status:** Accepted
**Category:** Asset Management

## Decision

Use SHA-filtered incremental downloads with batch atomic install and rollback capability

## Update Flow

```
1. Fetch remote tree (GitHub API - 1 call)
2. Load local asset-version.toml
3. Compare SHAs → Identify changed files only
4. Download changed files to temp directory (N API calls)
5. Batch install with rollback safety
6. Update asset-version.toml
7. Cleanup temp and backups
```

## Atomic Install Mechanism

```go
Phase 1: Backup
  - For each changed file being replaced
  - Rename: file → file.backup

Phase 2: Install
  - Move from temp to assets/
  - If any fail → rollback all from .backup

Phase 3: Commit
  - Update asset-version.toml with new SHAs
  - If fails → rollback all from .backup

Phase 4: Cleanup
  - Remove all .backup files
  - Remove temp directory
```

## Failure Handling

| Failure Point | Result | Recovery |
|---|---|---|
| Download fails | Temp cleaned up, assets/ untouched | None needed |
| Backup fails | Abort, no changes made | None needed |
| Install fails | Restore from .backup files | Automatic rollback |
| Version file write fails | Restore from .backup files | Automatic rollback |
| Process killed | Orphaned .backup files | Auto-cleanup next run |

## Disk Space

- Temp directory: Only changed files (typically 3-5 files, ~20 KB)
- Backup files: Only files being replaced (same size as changed)
- Total overhead: 2x size of changed files (not entire asset library)

## Benefits

- ✅ **Failure-safe:** All-or-nothing installation
- ✅ **Efficient:** Only downloads changed files (SHA filtering)
- ✅ **Recoverable:** Automatic rollback on any error
- ✅ **Minimal overhead:** Backups only for changed files
- ✅ **Simple:** No transaction log parsing needed

## Example (Incremental Update)

```
Remote has 28 files
Local has 25 matching SHAs, 3 different

Download: 3 files to temp (3 API calls)
Backup: 2 existing files (.backup)
Install: 3 files from temp
Update: asset-version.toml
Cleanup: 2 .backup files, temp directory

Result: 4 total API calls, minimal disk usage
```

## Rationale

Combining SHA filtering (DR-014) with batch atomic install provides:
- Best API efficiency (only changed files)
- Best safety (rollback on failure)
- Simplest implementation (no complex state tracking)
- Clear failure recovery (restore from .backup)

## Related Decisions

- [DR-014](./dr-014-github-tree-api.md) - SHA-based incremental downloads
- [DR-018](./dr-018-init-update-integration.md) - Shared with init command
