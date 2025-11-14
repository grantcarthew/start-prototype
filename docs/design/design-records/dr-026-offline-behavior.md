# DR-026: Offline Fallback and Network Unavailable Behavior

**Date:** 2025-01-07
**Status:** Accepted
**Category:** Asset Management

## Decision

The CLI handles network unavailability gracefully with clear error messages and no manual asset installation support. Network access is required only for asset download operations.

## What This Means

### No Manual Asset Installation

**Network-only approach:**
- No `--asset-path` flag for offline installation
- No support for manually dropping assets into `~/.config/start/assets/`
- No ZIP/tarball asset bundles for offline distribution
- Users must have network access for initial asset download

**Rationale:**
- KISS principle - one way to install assets
- Avoids security concerns with manual asset installation
- Simplifies implementation and testing
- Asset updates are optional, not critical path

### Network Failure Behavior

**start init:**
- Warns about network failure
- Creates config file anyway (succeeds)
- Informs user they can create their own config manually
- Exit code: 0 (success, config created despite asset download failure)

**start update:**
- Informs about network failure
- Exits immediately (fails)
- Does not modify existing assets
- Exit code: 1 (network error)

**start doctor (offline):**
- Displays what it can check locally
- Skips checks that require network (version, asset staleness)
- No errors for skipped checks
- Shows partial results
- Exit code: Based on local checks only

### Missing Assets Behavior

Defined in individual command specs:

**Pattern 1: Warn and continue**
- `start` / `start prompt` / `start task` - Work without assets
- Role references fail at runtime with clear error (expected)

**Pattern 2: Reduced functionality**
- `start assets add` - Can't browse catalog, manual creation still works
- `start config role` - No asset roles, user-defined roles still work

## Examples

### init with network failure

```bash
$ start init

Creating config at ~/.config/start/...

Warning: Unable to download assets from GitHub.
  Network error: dial tcp: no route to host

Config created successfully.

You can create your own configuration manually:
  start config edit global

Or try downloading assets later:
  start update
```

Exit code: 0 (success)

### update with network failure

```bash
$ start update

Error: Cannot connect to GitHub

  Network error: dial tcp: no route to host

Update requires network access.
Asset library not modified.

Check your internet connection and try again.
```

Exit code: 1 (network error)

### doctor offline

```bash
$ start doctor

Diagnosing start installation...
═══════════════════════════════════════════════════════════

Version
  start v1.2.3
  ⚠ Cannot check for updates (network unavailable)

Assets
  ✓ Asset library exists
  Last updated: 2024-12-15 (23 days ago)
  ⚠ Cannot check for newer assets (network unavailable)

Configuration
  ✓ Global config valid
  ✓ Local config valid

Agents (2 configured)
  ✓ claude - /usr/local/bin/claude
  ✓ gemini - /usr/local/bin/gemini

Contexts (2 required)
  ✓ environment - ~/reference/ENVIRONMENT.md
  ✓ index - ~/reference/INDEX.csv

Environment
  ✓ All checks passed

Summary
───────────────────────────────────────────────────────────
  ✓ No local issues found
  ⚠ Some checks skipped (network unavailable)

Run 'start doctor' with network access for complete check.
```

Exit code: 0 (local checks healthy)

### task add without assets

```bash
$ start assets add

Add new asset
─────────────────────────────────────────────────

⚠ Asset library not available (catalog unavailable)

Task name: my-task
Alias (optional): mt
Description (optional): My custom task

Task prompt: Do something useful with: {instructions}

[continues with normal flow, no templates offered]
```

Exit code: 0 (task added successfully)

## Rationale

### Network-Only Approach

**Why no manual installation?**
- Assets are optional content, not critical for core functionality
- Manual installation adds complexity (validation, security, UX)
- Self-optimizing SHA caching (DR-014) makes repeated downloads fast
- Offline scenarios are rare for development tools

**Why network failure in init is non-fatal?**
- Config file creation is the primary goal
- User can create config manually if needed
- Assets can be downloaded later via `start update`
- Per DR-018: init invokes update logic but continues on failure

**Why network failure in update is fatal?**
- Downloading assets is the explicit purpose of the command
- No point continuing if network is unavailable
- User expectations: update = download, or nothing

### Graceful Degradation

Commands that use assets follow the "warn and skip OR warn and continue" pattern established in existing command specs.

**Design philosophy:**
- Core functionality (launching agents) works without assets
- Assets enhance UX (templates, default roles/tasks)
- Missing assets = reduced convenience, not broken tool

## Implementation Notes

### What NOT to Implement

**Don't add:**
- Manual asset installation mechanisms
- Offline asset bundles or tarballs
- Asset mirroring or alternative download sources
- Complex offline fallback chains

### Command-Specific Behavior

**start init:**
```go
func runInit() error {
    // Create config structure
    if err := createConfigFile(); err != nil {
        return err
    }

    // Try to download assets
    if err := assets.Update(); err != nil {
        log.Warn("Unable to download assets: %v", err)
        log.Info("You can download assets later: start update")
        // Continue - don't return error
    }

    return nil  // Success
}
```

**start update:**
```go
func runUpdate() error {
    // Network is required - fail fast
    if err := assets.Update(); err != nil {
        return fmt.Errorf("update failed: %w", err)
    }
    return nil
}
```

**start doctor:**
```go
func runDoctor() error {
    checks := []Check{
        checkVersion(),        // Skip if network unavailable
        checkAssetAge(),       // Can check locally
        checkAssetStaleness(), // Skip if network unavailable
        checkConfig(),         // Always runs
        checkAgents(),         // Always runs
        checkContexts(),       // Always runs
        checkEnvironment(),    // Always runs
    }

    // Show results for all checks, mark skipped ones
    return displayResults(checks)
}
```

### User Guidance

**For missing assets:**
```
Warning: Asset templates not available.

To download asset templates:
  start update
```

**For network errors:**
```
Error: Cannot connect to GitHub

Check your internet connection and try again.

If you're behind a proxy, configure HTTP_PROXY and HTTPS_PROXY.
```

## Benefits

- ✅ **Simple** - One way to get assets (network download)
- ✅ **Secure** - No manual asset injection concerns
- ✅ **Clear** - Explicit about network requirements
- ✅ **Graceful** - Commands work without assets (reduced functionality)
- ✅ **User-controlled** - init succeeds even if assets fail
- ✅ **Explicit** - update requires network by design

## Trade-offs Accepted

- ❌ No air-gapped environment support (acceptable: rare for dev tools)
- ❌ No manual asset installation (acceptable: adds complexity)
- ❌ init might succeed with incomplete setup (acceptable: clear warnings)

## Related Decisions

- [DR-011](./dr-011-asset-distribution.md) - GitHub-fetched assets (network dependency established)
- [DR-014](./dr-014-github-tree-api.md) - SHA-based caching (optimizes repeated downloads)
- [DR-018](./dr-018-init-update-integration.md) - Init invokes update, warns on failure
- [DR-025](./dr-025-no-automatic-checks.md) - No automatic checks (predictable network usage)

## Future Considerations

If air-gapped environments become a requirement, we could add:

**Option 1: Manual asset directory support**
```bash
start update --from-directory /path/to/assets
```

**Option 2: Bundled asset snapshots**
```bash
start update --from-bundle assets-v1.2.3.tar.gz
```

**Current stance:** Don't implement unless users explicitly request it. The added complexity isn't justified by current use cases.
