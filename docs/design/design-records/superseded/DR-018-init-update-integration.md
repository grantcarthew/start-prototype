# DR-018: Init and Update Command Integration

- Date: 2025-01-06, Superseded 2025-01-17
- Status: Superseded by DR-031 (Catalog-Based Assets)
- Category: Asset Management

## Superseded

This design record is superseded by DR-031 (Catalog-Based Assets). The bulk download integration model is replaced by catalog-driven, on-demand asset loading.

The catalog system eliminates the need for init/update integration because:

- NO bulk downloads during `start init` - init creates config files only, no asset downloads
- NO `asset-version.toml` tracking file - filesystem IS the state
- NO shared update logic for bulk downloads - assets are lazy loaded on first use
- `start init` focuses on config creation and optional agent detection (prompts to download detected agents)
- `start assets update` updates individual cached assets via per-asset SHA comparison (not bulk update)
- Assets downloaded on-demand: via `start task <name>`, `start assets add`, or auto-download when `asset_download = true`

The bulk download model (init always downloads all assets) is incompatible with the catalog model (query and download on-demand).

## Integration Strategy

```
start init workflow:
1. Interactive setup wizard (agent selection, context detection)
2. Write config.toml
3. Call shared asset update logic
4. Display success message

start assets update workflow:
1. Call shared asset update logic
2. Display what changed
```

## Shared Update Implementation

Both commands use identical asset update logic:

```go
// Package: internal/assets
func UpdateAssets() error {
    // 1. Fetch remote tree from GitHub (DR-014)
    // 2. Load local asset-version.toml
    // 3. Compare SHAs → identify changed files
    // 4. Download changed files to temp
    // 5. Atomic install with rollback (DR-015)
    // 6. Update asset-version.toml
    // 7. Cleanup
}
```

## No Conditional Logic

- `start init` **always** updates assets (no staleness checks)
- No flags: `--skip-assets`, `--force`, etc. (KISS principle)
- SHA comparison naturally skips unchanged files (efficient by default)

## Network Failure Handling

```
start init (network fails):
  ✓ Config created successfully
  ⚠ Warning: Asset download failed (network unavailable)
    Assets can be downloaded later with 'start assets update'
  Exit code: 0 (success)

start assets update (network fails):
  ✗ Error: Cannot reach GitHub (network unavailable)
    Check network connection and try again
  Exit code: 1 (failure)
```

## First Run Scenario

```bash
$ start init
# Creates config, downloads assets (all 28 files)
# API calls: 1 tree + 28 contents = 29 total

$ start init  # Run again accidentally
# Config exists, backs up, rewrites
# Assets: SHA comparison shows 0 changes
# API calls: 1 tree + 0 contents = 1 total
# No unnecessary downloads
```

## Benefits

- ✅ **No complexity:** Removed all conditional flags and staleness logic
- ✅ **Self-optimizing:** SHA comparison prevents redundant downloads
- ✅ **Predictable:** `init` always tries to update, `update` requires network
- ✅ **Offline-friendly:** Init can proceed without assets (warns user)
- ✅ **Efficient:** First-time downloads all, subsequent checks are minimal

## Rationale

- **Simplicity:** No special cases, staleness checks, or flags
- **Consistency:** One update algorithm, two entry points
- **Efficiency:** SHA comparison ensures minimal downloads
- **User experience:** init succeeds even if network fails, update requires network
- **Maintainability:** Single source of truth for update logic

## Related Decisions

- [DR-014](./dr-014-github-tree-api.md) - SHA-based download strategy
- [DR-015](./dr-015-atomic-updates.md) - Atomic update mechanism
