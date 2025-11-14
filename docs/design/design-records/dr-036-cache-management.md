# DR-036: Cache Management

**Date:** 2025-01-10
**Status:** Accepted
**Category:** Asset Management

## Decision

Asset cache is an invisible implementation detail with no user-facing management commands. Users can manually delete the cache directory if needed.

## What This Means

### Cache is Invisible

**No cache commands:**
- ❌ No `start cache list`
- ❌ No `start cache clean`
- ❌ No `start cache info`
- ❌ No `start cache prune`

**Why no commands?**
- Asset files are tiny text files (typically < 5KB)
- Even hundreds of assets = < 1MB total
- No need for size limits or cleanup
- Adds complexity for minimal benefit
- KISS principle

**If user wants to clear cache:**
```bash
# Simple manual deletion
rm -rf ~/.config/start/assets

# Assets re-download automatically on next use
```

### Cache Structure

```
~/.config/start/assets/
├── roles/
│   ├── general/
│   │   ├── code-reviewer.md
│   │   ├── code-reviewer.meta.toml
│   │   ├── default.md
│   │   └── default.meta.toml
│   └── languages/
│       ├── go-expert.md
│       └── go-expert.meta.toml
├── tasks/
│   ├── git-workflow/
│   │   ├── pre-commit-review.toml
│   │   ├── pre-commit-review.md
│   │   ├── pre-commit-review.meta.toml
│   │   ├── pr-ready.toml
│   │   ├── pr-ready.md
│   │   └── pr-ready.meta.toml
│   └── code-quality/
│       └── ...
└── agents/
    └── claude/
        ├── sonnet.toml
        └── sonnet.meta.toml
```

**Properties:**
- **Location:** `~/.config/start/assets/` (configurable via `asset_path`)
- **Structure:** `{type}/{category}/{name}.{ext}`
- **Contents:** Asset files + sidecar metadata
- **State:** Filesystem IS the state (no tracking files)

### Cache Behavior

**Automatic population:**
```bash
# User downloads asset
start task pre-commit-review

# Cache populated automatically:
#   ~/.config/start/assets/tasks/git-workflow/pre-commit-review.toml
#   ~/.config/start/assets/tasks/git-workflow/pre-commit-review.md
#   ~/.config/start/assets/tasks/git-workflow/pre-commit-review.meta.toml
```

**Automatic updates:**
```bash
# User updates cache
start assets update

# Checks GitHub for newer SHAs
# Downloads updated assets to cache
# User config remains unchanged
```

**No automatic cleanup:**
- Cache persists until manually deleted
- No age-based cleanup
- No size-based cleanup
- Consistent with DR-025 (no automatic operations)

### Cache Operations

**Read (during resolution):**
```go
func findInCache(assetType, name string) string {
    pattern := filepath.Join(
        assetPath,
        assetType,
        "*",  // Any category
        fmt.Sprintf("%s.toml", name),
    )
    matches, _ := filepath.Glob(pattern)
    if len(matches) > 0 {
        return matches[0]
    }
    return ""
}
```

**Write (after download):**
```go
func cacheAsset(asset *Asset, metadata *AssetMetadata) error {
    // Determine cache location
    cachePath := filepath.Join(
        assetPath,
        asset.Type,
        metadata.Category,
    )

    // Create directory structure
    os.MkdirAll(cachePath, 0755)

    // Write asset content
    contentPath := filepath.Join(cachePath, fmt.Sprintf("%s.toml", asset.Name))
    os.WriteFile(contentPath, asset.Content, 0644)

    // Write prompt file (if exists)
    if asset.PromptContent != nil {
        promptPath := filepath.Join(cachePath, fmt.Sprintf("%s.md", asset.Name))
        os.WriteFile(promptPath, asset.PromptContent, 0644)
    }

    // Write metadata
    metaPath := filepath.Join(cachePath, fmt.Sprintf("%s.meta.toml", asset.Name))
    writeMetadata(metaPath, metadata)

    return nil
}
```

**Update (during `start assets update`):**
```go
func updateCachedAsset(asset *Asset, newMetadata *AssetMetadata) error {
    // Overwrite existing files
    return cacheAsset(asset, newMetadata)
}
```

### Cache Persistence

**Cache is local only:**
- Not synchronized across machines
- Not backed up automatically
- Each machine has its own cache

**Re-downloading is fast:**
- Using raw.githubusercontent.com (no rate limit per DR-034)
- Assets are small (< 5KB typically)
- Network cost is minimal

**User config is portable:**
```toml
# User's tasks.toml can reference cache
[tasks.pre-commit-review]
prompt_file = "~/.config/start/assets/tasks/git-workflow/pre-commit-review.md"

# On new machine:
# 1. Copy config files
# 2. Run tasks - cache populated automatically
# 3. OR run: start assets update (download all assets)
```

## Size Estimates

### Minimal Viable Set (28 assets)

**Tasks (12):**
- 12 × 3KB (toml + md + meta) = 36KB

**Roles (8):**
- 8 × 2KB (md + meta) = 16KB

**Agents (6):**
- 6 × 1KB (toml + meta) = 6KB

**Templates (2):**
- 2 × 2KB (toml + meta) = 4KB

**Total: ~62KB** (tiny!)

### Large Catalog (200 assets)

**Assuming average 3KB per asset:**
- 200 × 3KB = 600KB

**Still trivial for modern systems.**

### Conclusion

No need for size limits or cleanup. Even a large catalog is negligible.

## Configuration

**Settings in config.toml:**
```toml
[settings]
asset_path = "~/.config/start/assets"  # Cache location
```

**Custom cache location:**
```toml
[settings]
asset_path = "/custom/path/to/assets"
```

**Useful for:**
- Network drives
- Shared team caches
- Testing/development

## Troubleshooting

**Cache corrupted or broken?**
```bash
# Delete and re-download
rm -rf ~/.config/start/assets
start assets update  # Re-download all configured assets
```

**Asset not found but should be cached?**
```bash
# Check cache contents
ls -lR ~/.config/start/assets/tasks/

# Check for metadata file
cat ~/.config/start/assets/tasks/git-workflow/pre-commit-review.meta.toml
```

**Want to see what's cached?**
```bash
# Simple directory listing
find ~/.config/start/assets -name "*.toml" -not -name "*.meta.toml"

# Or use tree
tree ~/.config/start/assets
```

**Future consideration:** If users request it, add `start doctor` check for cache health.

## Benefits

**Simple:**
- ✅ No user-facing cache commands
- ✅ No cache management logic
- ✅ No size limits to implement
- ✅ No cleanup policies

**Transparent:**
- ✅ Standard filesystem directory
- ✅ Human-readable text files
- ✅ Easy to inspect manually
- ✅ Can use standard tools (ls, cat, rm)

**Reliable:**
- ✅ Filesystem is the state (no drift)
- ✅ No tracking files to corrupt
- ✅ Easy to reset (delete directory)
- ✅ Automatic re-population

**Efficient:**
- ✅ Fast lookups (filesystem cache)
- ✅ No database or index needed
- ✅ Minimal disk space usage
- ✅ No background processes

## Trade-offs Accepted

**No visibility into cache:**
- ❌ Can't easily see what's cached without ls
- **Mitigation:** Users rarely need to know, can use filesystem tools

**No automatic cleanup:**
- ❌ Cache grows over time (slightly)
- **Mitigation:** Files are tiny, growth is negligible

**No cache statistics:**
- ❌ Can't see "cache size: 142KB, 28 assets"
- **Mitigation:** Not needed for tiny files, can add to `start doctor` if requested

**Manual deletion only:**
- ❌ Can't selectively remove cached assets via CLI
- **Mitigation:** User config is source of truth, cache is disposable

## Related Decisions

- [DR-031](./dr-031-catalog-based-assets.md) - Catalog architecture (cache role)
- [DR-033](./dr-033-asset-resolution-algorithm.md) - Resolution (cache lookup)
- [DR-034](./dr-034-github-catalog-api.md) - GitHub API (cache population)
- [DR-037](./dr-037-asset-updates.md) - Updates (cache refresh)
- [DR-025](./dr-025-no-automatic-checks.md) - No automatic operations (no auto-cleanup)

## Future Considerations

**If cache management becomes necessary:**

Potential additions (not in v1):
- `start doctor` cache health check
- `start doctor --fix` cache repair
- Cache statistics in `start doctor` output
- Shared cache for team environments

**Current stance:** Keep it simple. Cache is invisible. Users can always delete the directory manually if needed.
