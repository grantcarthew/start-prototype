# DR-025: No Automatic Checks or Caching

- Date: 2025-01-06
- Status: Accepted
- Category: CLI Design

## Problem

The CLI needs to decide when and how to perform health checks (version updates, asset staleness, configuration validation). The strategy must balance:

- Keeping users informed about available updates and issues
- Respecting user workflow (no unexpected delays or network calls)
- Command execution performance (fast startup, no overhead)
- Offline usability (commands work without network access)
- Implementation complexity (cache management, staleness logic)
- User control over when checks occur

## Decision

The CLI does not perform automatic background health checks or cache health check results. All checks are user-initiated only via explicit commands.

No automatic checks:

- No automatic version checks on command execution
- No periodic background health checks
- No "you haven't run doctor in 30 days" nag messages
- No silent network calls during normal operation

No result caching:

- start doctor always performs fresh checks (no cached results)
- No "last checked 2 hours ago" cached status
- No "check is still valid for X hours" logic
- Every execution makes fresh API calls

User-initiated only:

- start doctor - explicitly checks health (always fresh, requires network)
- start assets update - explicitly updates assets (always fresh, requires network)
- All other commands - no health checks (fast, offline-friendly)

## Why

Respectful of user workflow:

- No unexpected network calls during work
- No performance impact on normal commands
- No interruptions or nagging messages
- User controls exactly when checks happen
- Commands work offline (except doctor and update)

Simpler implementation:

- No cache file management for health check results
- No cache invalidation logic or staleness calculations
- No background job scheduling or goroutines
- Fewer edge cases to handle
- Less state to manage

Always fresh data:

- When user runs doctor, they get current information
- No confusion about "last checked" vs "current status"
- No stale cached results misleading users
- Clear expectation: doctor checks now, not from cache

Fast command execution:

- start launches agent immediately (no delays)
- start task executes immediately (no overhead)
- No network round-trips during normal operation
- Predictable performance

## Trade-offs

Accept:

- User might miss available updates (they must remember to run doctor)
- User might use outdated assets (no notification when updates available)
- No "last checked" information (always fresh, but no historical context)
- Repeated doctor calls make repeated API requests (no caching efficiency)
- Users must develop their own update check cadence

Gain:

- Extremely fast command execution (no check overhead, no network calls)
- Predictable behavior (checks only when user explicitly requests)
- Respectful of user workflow (no unexpected network activity or interruptions)
- Simple implementation (no cache management, invalidation, or staleness logic)
- Offline-friendly (all commands work without network except doctor/update)
- No nagging messages during work sessions
- Clear mental model (doctor = check now, other commands = execute now)

## Alternatives

Automatic checks on every command (like npm, homebrew):

Example approaches:
- npm checks for updates periodically, shows "update available" message
- homebrew auto-updates on brew install/upgrade (can disable with HOMEBREW_NO_AUTO_UPDATE)
- rustup checks for toolchain updates with configurable behavior

Pros:
- Users stay informed about updates automatically
- No need to remember to check manually
- Updates are discovered quickly

Cons:
- Adds network latency to every command execution
- Unexpected network calls (privacy and performance concerns)
- Commands slower and less predictable
- Requires network for normal operation
- Interrupts user workflow with messages
- Requires cache management to avoid excessive API calls

Rejected: Performance and predictability more important than automatic update awareness. Users can develop their own cadence for running doctor.

Cached results with time-to-live:

Example: Cache doctor results for 24 hours, show cached status if fresh
- start doctor within 24h: "Last checked 2 hours ago (cached), status: healthy"
- start doctor after 24h: Perform fresh check

Pros:
- Reduces API calls for repeated doctor checks
- Could enable "check on first use each day" pattern
- Provides "last checked" information

Cons:
- Adds cache file management complexity
- Requires cache invalidation logic
- Stale results can mislead users
- TTL value is arbitrary (too short = ineffective, too long = stale)
- Users unsure if status is current or cached
- More state to manage and debug

Rejected: Complexity outweighs benefits. Fresh checks provide clear, unambiguous status.

Optional automatic check setting:

Example config:
```toml
[settings]
auto_check = false  # Default: disabled
auto_check_frequency = "daily"
```

Pros:
- Users can opt-in if they want automatic checks
- Flexibility for different workflows
- Advanced users can configure behavior

Cons:
- Adds implementation complexity (scheduling, cache, staleness)
- Two code paths to maintain and test
- Makes command behavior less predictable (depends on config)
- Conflicts with "user-initiated only" principle
- Most users would leave it disabled anyway

Rejected: Adds significant complexity for uncertain benefit. Current design is simple and clear.

## Structure

Commands that perform health checks:

- start doctor - Always checks, always fresh, requires network
- start assets update - Always checks, always fresh, requires network

Commands that do not perform checks:

- start - No checks, fast execution, offline-friendly
- start prompt - No checks
- start task - No checks
- start config (all subcommands) - No checks
- start init - No checks (creates config only)

Asset file caching (separate concern):

- Asset downloads use catalog-based lazy loading
- Downloaded assets cached locally to avoid re-downloading
- This is file content caching (prevents redundant downloads)
- Not health check result caching (different concern)

## Usage Examples

Normal command execution (no checks):

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

Explicit health check:

```bash
$ start doctor
# User explicitly requests check
# Makes fresh network calls to GitHub
# Shows comprehensive current status
# No cached results
```

Explicit update:

```bash
$ start assets update
# User explicitly requests update
# Downloads new assets from catalog
# Shows what changed
# Always fresh check
```

Comparison with other tools:

Homebrew (automatic):
```bash
# Auto-updates on brew install/upgrade
# Can disable with HOMEBREW_NO_AUTO_UPDATE
```

npm (automatic):
```bash
# Checks for npm updates periodically
# Shows "npm update available" message
```

Our approach (manual):
```bash
# No automatic checks
# User runs 'start doctor' when they want to check
# User runs 'start assets update' when they want to update
```

## Scope

This decision applies to:

- Health check execution (when checks occur)
- Result caching (whether check results are cached)
- All commands in the CLI (which commands check, which don't)

This decision does not apply to:

- Asset file caching (downloaded assets are cached to avoid re-downloading)
- Configuration file loading (configs are read on every execution)
- Version injection at build time (embedded in binary)

## Updates

- 2025-01-17: Removed references to superseded DR-014 and DR-023, updated asset caching description for catalog system
