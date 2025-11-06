# DR-025: No Automatic Checks or Caching

**Date:** 2025-01-06
**Status:** Accepted
**Category:** CLI Design

## Decision

The CLI does **not** perform automatic background health checks or cache health check results. All checks are user-initiated only.

## What This Means

### No Automatic Checks

**No background checking:**
- No automatic version checks on every command execution
- No periodic background health checks
- No "you haven't run doctor in 30 days" nag messages
- No silent network calls during normal operation

**User-initiated only:**
- `start doctor` - explicitly checks health
- `start update` - explicitly updates assets and checks CLI version
- All other commands (`start`, `start task`, etc.) - no health checks

### No Result Caching

**No cached health check results:**
- `start doctor` always performs fresh checks (network calls)
- No "last checked 2 hours ago" cached status
- No "check is still valid for X hours" logic
- Every execution makes fresh API calls (with rate limit protection)

**Asset download caching only:**
- DR-014 SHA-based caching still applies (for asset downloads)
- This prevents re-downloading unchanged files
- This is file content caching, not health check caching

## Rationale

**Respectful of user's workflow:**
- No unexpected network calls
- No performance impact on normal commands
- No nagging or interruptions
- User controls when checks happen

**Simpler implementation:**
- No cache file management for health checks
- No cache invalidation logic
- No staleness calculations for cached results
- No background job scheduling

**Always fresh data:**
- When user runs doctor, they get current information
- No confusion about "last checked" vs "current status"
- No stale cached results

**Consistent with design:**
- DR-021: CLI version checks have no caching
- DR-023: Asset staleness checks have no caching
- This decision makes that pattern explicit

## Comparison with Other Tools

Many CLI tools do automatic checks. We explicitly choose not to:

**Homebrew:**
```bash
# Auto-updates on brew install/upgrade
# Can disable with HOMEBREW_NO_AUTO_UPDATE
```

**npm:**
```bash
# Checks for npm updates periodically
# Shows "npm update available" message
```

**rustup:**
```bash
# Checks for toolchain updates
# Can configure update behavior
```

**Our approach (start):**
```bash
# No automatic checks
# User runs 'start doctor' when they want to check
# User runs 'start update' when they want to update
```

## User Experience

### Normal Command Execution

```bash
$ start
# Launches agent immediately
# No version checks
# No network calls
# No delays

$ start task review
# Executes task immediately
# No health checks
# Fast execution
```

### Explicit Health Check

```bash
$ start doctor
# User explicitly requests check
# Makes network calls to GitHub
# Shows comprehensive status
# Fresh data every time
```

### Explicit Update

```bash
$ start update
# User explicitly requests update
# Downloads new assets
# Checks CLI version
# Shows what changed
```

## Benefits

- ✅ **Fast** - Normal commands have no check overhead
- ✅ **Predictable** - Checks only when user requests
- ✅ **Respectful** - No unexpected network activity
- ✅ **Simple** - No cache management complexity
- ✅ **Offline-friendly** - Commands work without network (except doctor/update)
- ✅ **No nagging** - No "you should update" messages during work

## Trade-offs Accepted

- ❌ User might miss available updates (acceptable: they can run doctor)
- ❌ User might use outdated assets (acceptable: they control update timing)
- ❌ No "last checked" information (acceptable: always fresh when checking)

## Implementation Notes

### What NOT to Implement

**Don't add:**
- Background check threads/goroutines
- Cache files for health check results (e.g., `~/.config/start/doctor-cache.json`)
- "Last checked" timestamps
- Automatic update notifications
- Periodic check timers
- "Check on first use each day" logic

### What to Keep

**Do keep:**
- SHA-based asset download caching (DR-014)
- Rate limit checking before API calls
- Network error handling

### Command Behavior

**Commands that check health:**
- `start doctor` - always checks, always fresh
- `start update` - always checks, always fresh

**Commands that don't check:**
- `start` - no checks
- `start prompt` - no checks
- `start task` - no checks
- `start config *` - no checks
- `start init` - no checks (just downloads assets)

## Documentation

### Help Text

Make this clear in `--help`:

```
start doctor - Check system health and version status

This command always performs fresh checks (no caching).
It requires network access to check for updates.

Run 'start doctor' periodically to check for updates.
```

### README

Document the philosophy:

```markdown
## Updates

`start` does not automatically check for updates. This keeps
normal commands fast and respects your workflow.

To check for updates:
- Run `start doctor` to check health and version status
- Run `start update` to update the asset library

We recommend running `start doctor` occasionally to check
for CLI and asset updates.
```

## Related Decisions

- [DR-014](./dr-014-github-tree-api.md) - SHA-based asset caching (file content caching)
- [DR-021](./dr-021-github-version-check.md) - CLI version checking (no result caching)
- [DR-023](./dr-023-asset-staleness-check.md) - Asset staleness checking (no result caching)
- [DR-024](./dr-024-doctor-exit-codes.md) - Doctor exit codes

## Future Considerations

If users request automatic checks, we could add:

**Optional background checks** (not implementing now):
```toml
[settings]
auto_check = false  # Default: disabled
auto_check_frequency = "daily"  # If enabled
```

**But this conflicts with our philosophy:**
- Adds complexity
- Requires cache management
- Makes commands less predictable
- Goes against "user-initiated only" principle

**Current stance:** Don't implement unless users strongly request it.
