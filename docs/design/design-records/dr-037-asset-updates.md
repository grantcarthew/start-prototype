# DR-037: Asset Update Mechanism

**Date:** 2025-01-10
**Status:** Accepted
**Category:** Asset Management

## Decision

Use SHA-based update detection via GitHub Tree API to check for asset updates. Update only the cache, never the user's config. Provide manual-only update command with clear messaging.

**Note:** Per [DR-041](./dr-041-asset-command-reorganization.md), command moved from `start assets update` to `start assets update` with optional query parameter for selective updates.

## What This Means

### Update Command

**Manual update only:**
```bash
start assets update
```

**Consistent with DR-025 (no automatic checks):**
- User explicitly runs `start assets update`
- No background checks
- No automatic downloads
- No version checking on CLI startup

### Update Process

**When user runs `start assets update`:**

1. Fetch GitHub catalog tree (SHA for every file)
2. For each cached `.meta.toml` file:
   - Read local SHA
   - Find corresponding file in GitHub tree by path
   - Compare local SHA with remote SHA
   - If different → download new version
   - Update cached files (content + metadata)
3. Report what was updated
4. **Never touch user's config files**

### SHA-Based Versioning

**No semantic versioning** - Git blob SHA IS the version:

```toml
# Local cached metadata
# ~/.config/start/assets/tasks/git-workflow/pre-commit-review.meta.toml
sha = "a1b2c3d4e5f6789012345678901234567890abcd"
updated = "2025-01-10T12:00:00Z"
```

**GitHub Tree API provides current SHA:**
```json
{
  "path": "assets/tasks/git-workflow/pre-commit-review.toml",
  "sha": "b2c3d4e5f6789012345678901234567890abcdef"  // Different!
}
```

**SHA mismatch = update available**

This DR defines how `start assets update` works.

The scope of this command is limited to refreshing assets that the user has already acquired (i.e., are present in the local cache). It does not discover or add new assets from the catalog.

## User Config Never Automatically Overwritten

**Cache updates are separate from config:**

When you run `start assets update`:
1. The asset cache gets new versions.
2. Your configuration files (`tasks.toml`, etc.) remain unchanged.
3. If a task in your config references a cached file (e.g., via `prompt_file`), it will automatically use the new content on the next run.
4. If a task in your config has inlined content copied from an asset, you must manually update it to reflect the changes.


## Implementation

### Update Algorithm

```go
func RunUpdate() error {
    fmt.Println("Checking for asset updates...")

    // Fetch current catalog tree
    catalog, err := fetchGitHubTree()
    if err != nil {
        return fmt.Errorf("failed to fetch catalog: %w", err)
    }

    // Find all cached metadata files
    metaFiles := findCachedMetadata()

    var updated, unchanged int

    for _, metaPath := range metaFiles {
        // Load local metadata
        localMeta, err := LoadMetadata(metaPath)
        if err != nil {
            log.Warn("Failed to load %s: %v", metaPath, err)
            continue
        }

        // Find asset in GitHub tree
        assetPath := reconstructAssetPath(localMeta)
        remoteSHA := catalog.GetSHA(assetPath)

        if remoteSHA == "" {
            log.Warn("Asset not found in catalog: %s", assetPath)
            continue
        }

        // Compare SHAs
        if localMeta.SHA == remoteSHA {
            unchanged++
            log.Debug("%s (up to date)", localMeta.Name)
            continue
        }

        // Download new version
        fmt.Printf("  ⬇ Updating %s/%s...\n", localMeta.AssetType, localMeta.Name)
        asset := downloadAsset(assetPath)
        remoteMeta := downloadMetadata(assetPath + ".meta.toml")

        // Update cache
        cacheAsset(asset, remoteMeta)
        updated++
    }

    // Report results
    fmt.Printf("\n✓ Update complete\n")
    fmt.Printf("  Updated: %d assets\n", updated)
    fmt.Printf("  Unchanged: %d assets\n", unchanged)

    if updated > 0 {
        fmt.Println("\nNote: Your configuration files are unchanged.")
        fmt.Println("Review updated assets and manually update config if desired.")
    }

    return nil
}
```

### Finding Cached Metadata

```go
func findCachedMetadata() []string {
    pattern := filepath.Join(assetPath, "**/*.meta.toml")
    matches, _ := filepath.Glob(pattern)
    return matches
}
```

### Reconstructing Asset Path

```go
func reconstructAssetPath(meta *AssetMetadata) string {
    // From: ~/.config/start/assets/tasks/git-workflow/pre-commit-review.meta.toml
    // To: assets/tasks/git-workflow/pre-commit-review.toml

    return fmt.Sprintf("assets/%s/%s/%s.toml",
        meta.AssetType,
        meta.Category,
        meta.Name,
    )
}
```

### Comparing SHAs

```go
func (t *GitHubTree) GetSHA(path string) string {
    for _, item := range t.Tree {
        if item.Path == path && item.Type == "blob" {
            return item.SHA
        }
    }
    return ""
}
```

## User Experience

### Example 1: Updates Available

```bash
$ start assets update

Checking for asset updates...

  ⬇ Updating tasks/pre-commit-review...
  ⬇ Updating roles/code-reviewer...

✓ Update complete
  Updated: 2 assets
  Unchanged: 10 assets

Note: Your configuration files are unchanged.
Review updated assets and manually update config if desired.

To see changes:
  diff ~/.config/start/assets/tasks/git-workflow/pre-commit-review.toml \
       <(git show b2c3d4e:assets/tasks/git-workflow/pre-commit-review.toml)
```

### Example 2: No Updates

```bash
$ start assets update

Checking for asset updates...

✓ Update complete
  Updated: 0 assets
  Unchanged: 12 assets

All cached assets are up to date.
```

### Example 3: Network Error

```bash
$ start assets update

Checking for asset updates...

Error: Cannot connect to GitHub

  Network error: dial tcp: no route to host

Check your internet connection and try again.
```

### Example 4: First Run (No Cache)

```bash
$ start assets update

Checking for asset updates...

✓ Update complete
  Updated: 0 assets
  Unchanged: 0 assets

No cached assets found.

To download assets:
  - Browse catalog: start assets add
  - Use a task: start task <name>
```

## Comparing Updates

**User wants to see what changed:**

```bash
# View updated asset content
cat ~/.config/start/assets/tasks/git-workflow/pre-commit-review.toml

# View updated metadata
cat ~/.config/start/assets/tasks/git-workflow/pre-commit-review.meta.toml

# Diff with GitHub (if git available)
cd /tmp
git clone --depth=1 https://github.com/grantcarthew/start.git
diff ~/. config/start/assets/tasks/git-workflow/pre-commit-review.toml \
     /tmp/start/assets/tasks/git-workflow/pre-commit-review.toml
```

**Future enhancement:** `start assets update --diff` to show changes

## Update Notifications

**Per DR-025 (no automatic checks):**
- ❌ No "updates available" message on CLI startup
- ❌ No background update checking
- ❌ No version checking during `start task`

**User must explicitly run `start assets update` to check**

**Future consideration:** Opt-in notifications
```toml
[settings]
check_updates_on_start = false  # Default: false (compliant with DR-025)
```

## Selective Updates (DR-040, DR-041)

**No arguments - update all:**
```bash
start assets update              # Update all cached assets
```

**With query - update matching (DR-040):**
```bash
start assets update "commit"     # Update assets matching 'commit' (substring)
start assets update git-workflow # Update all in git-workflow category
start assets update pre-commit-review  # Update specific asset
```

Uses substring matching algorithm (DR-040) to find matching assets in cache, then updates only those.

## Rollback

**User wants previous version:**

```bash
# Manual rollback (delete cache, re-download old version)
rm -rf ~/.config/start/assets/tasks/git-workflow/pre-commit-review.*

# Use config version (if user saved it)
# Config references cache, but user can inline content if needed
```

**Future enhancement:** Version history
```bash
start assets update --rollback tasks/pre-commit-review  # Restore previous SHA
```

**v1:** No rollback mechanism (delete and re-download if needed)

## Error Handling

### Partial Update Failure

```bash
$ start assets update

Checking for asset updates...

  ⬇ Updating tasks/pre-commit-review...
  ✓ Updated tasks/pre-commit-review

  ⬇ Updating roles/code-reviewer...
  ✗ Failed: download error

✓ Update partially complete
  Updated: 1 asset
  Failed: 1 asset
  Unchanged: 10 assets

Errors:
  - roles/code-reviewer: network timeout

Try running 'start assets update' again to retry failed downloads.
```

### Asset Removed from Catalog

```bash
$ start assets update

Checking for asset updates...

  ⚠ Warning: tasks/deprecated-task not found in catalog
    (asset exists in cache but removed from GitHub)

  Keep cached version? [Y/n] _
```

### Invalid SHA in Metadata

```bash
$ start assets update

Checking for asset updates...

  ✗ Error: Invalid SHA in tasks/pre-commit-review.meta.toml
    Expected 40-char hex, got: "invalid"

  Skipping this asset.
  To fix: Delete cache and re-download
    rm ~/.config/start/assets/tasks/git-workflow/pre-commit-review.*
```

## Benefits

**Reliable versioning:**
- ✅ SHA comparison is bulletproof
- ✅ No version number conflicts
- ✅ Content hash guarantees integrity

**User control:**
- ✅ Manual updates only (DR-025 compliant)
- ✅ User config never auto-modified
- ✅ Explicit opt-in to updates

**Simple implementation:**
- ✅ Single Tree API call
- ✅ SHA comparison is straightforward
- ✅ No complex version constraints

**Cache-only updates:**
- ✅ Cache refresh is safe
- ✅ User modifications preserved
- ✅ Clear separation of concerns

## Trade-offs Accepted

**No automatic notifications:**
- ❌ User must remember to run `start assets update`
- **Mitigation:** Consistent with DR-025, users can set calendar reminder

**No rollback in v1:**
- ❌ Can't easily revert to previous version
- **Mitigation:** Cache is disposable, can delete and re-download

**No selective updates in v1:**
- ❌ Must update all cached assets
- **Mitigation:** Updates are fast, can add selective updates in v2

**Manual config updates:**
- ❌ User must manually update config after cache update
- **Mitigation:** Preserves user customizations, explicit is better

## Configuration

**Settings in config.toml:**
```toml
[settings]
asset_repo = "grantcarthew/start"  # Repository to check for updates
```

**Custom repository:**
```toml
[settings]
asset_repo = "myorg/custom-assets"  # Use custom asset repository
```

## Related Decisions

- [DR-031](./dr-031-catalog-based-assets.md) - Catalog architecture (update context)
- [DR-032](./dr-032-asset-metadata-schema.md) - Metadata schema (SHA field)
- [DR-034](./dr-034-github-catalog-api.md) - GitHub API (tree fetching)
- [DR-036](./dr-036-cache-management.md) - Cache structure (where updates go)
- [DR-025](./dr-025-no-automatic-checks.md) - No automatic operations (manual updates only)
- [DR-040](./dr-040-substring-matching.md) - Substring matching algorithm (query parameter support)
- [DR-041](./dr-041-asset-command-reorganization.md) - Asset command reorganization (moved from `start update`)

## Future Considerations

**Enhanced update features:**
- Interactive update (show changes, confirm each)
- Diff view before/after
- ~~Selective updates by type or asset~~ **IMPLEMENTED** via query parameter (DR-040)
- Rollback to previous version
- Update history and changelog

**Current stance:** Keep v1 simple. Update all cached assets, report results. Enhance based on user feedback.
