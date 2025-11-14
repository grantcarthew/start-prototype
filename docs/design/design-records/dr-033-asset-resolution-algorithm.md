# DR-033: Asset Resolution Algorithm

**Date:** 2025-01-10
**Status:** Accepted
**Category:** Asset Management

## Decision

Assets resolve in priority order: local config → global config → cache → GitHub catalog, with automatic download controlled by `asset_download` setting and `--asset-download` flag.

## What This Means

### Resolution Priority Order

When user runs `start task <name>` (or role, agent, etc.):

```
1. Local config      .start/tasks.toml
2. Global config     ~/.config/start/tasks.toml
3. Asset cache       ~/.config/start/assets/tasks/*/<name>.toml
4. GitHub catalog    Query catalog, download if allowed
5. Error             Not found anywhere
```

**First match wins** - No merging, no fallback to next level once found.

### Resolution Algorithm

```go
func ResolveAsset(assetType, name string, opts ResolveOptions) (*Asset, error) {
    // Step 1: Local config (exact match, then alias)
    if asset := localConfig.Get(assetType, name); asset != nil {
        log.Debug("Found in local config: %s", name)
        return asset, nil
    }

    // Step 2: Global config (exact match, then alias)
    if asset := globalConfig.Get(assetType, name); asset != nil {
        log.Debug("Found in global config: %s", name)
        return asset, nil
    }

    // Step 3: Asset cache
    cachePath := findInCache(assetType, name)
    if cachePath != "" {
        log.Info("Using cached asset: %s", name)
        asset := loadFromCache(cachePath)
        return asset, nil
    }

    // Step 4: Check if downloads allowed
    if !opts.AssetDownload {
        return nil, &ErrNotFound{
            AssetType: assetType,
            Name:      name,
            Reason:    "download disabled",
        }
    }

    // Step 5: GitHub catalog
    if !isOnline() {
        return nil, &ErrNetworkRequired{
            AssetType: assetType,
            Name:      name,
        }
    }

    // Get catalog (in-memory cached)
    catalog := getCatalog()
    githubPath := catalog.Find(assetType, name)
    if githubPath == "" {
        return nil, &ErrNotFound{
            AssetType: assetType,
            Name:      name,
            Reason:    "not in catalog",
        }
    }

    // Download and cache
    log.Info("Found in GitHub catalog: %s", githubPath)
    asset := downloadAsset(githubPath)
    cacheAsset(asset)

    // Add to config
    scope := "global"
    if opts.Local {
        scope = "local"
    }
    addToConfig(asset, scope)
    log.Info("Added to %s config: %s", scope, name)

    return asset, nil
}
```

### Configuration Settings

**In config.toml:**
```toml
[settings]
asset_download = true  # Default: auto-download from GitHub
```

**Behavior:**
- `asset_download = true` - Auto-download if asset not found in config/cache
- `asset_download = false` - Fail if asset not found in config/cache

### Command-Line Flags

**Override setting per-command:**
```bash
start task <name> [flags]

Flags:
  --local                  Add downloaded asset to local config (default: global)
  --asset-download[=bool]  Download from GitHub if not found (default: from settings)
```

**Flag precedence:**
- Flag explicitly set → use flag value
- Flag not set → use `asset_download` setting
- Setting not set → default to `true`

### Behavior Matrix

| Setting | Flag | Not in Config/Cache | Action |
|---------|------|---------------------|--------|
| `true`  | (none)                   | Found in GitHub     | Download, cache, add to global |
| `false` | (none)                   | Found in GitHub     | Error (download disabled)      |
| (any)   | `--asset-download`       | Found in GitHub     | Download, cache, add to global |
| (any)   | `--asset-download=true`  | Found in GitHub     | Download, cache, add to global |
| (any)   | `--asset-download=false` | Found in GitHub     | Error (download disabled)      |
| (any)   | `--asset-download --local` | Found in GitHub     | Download, cache, add to local  |

## Examples

### Example 1: Default Behavior (asset_download = true)

```bash
$ start task pre-commit-review

Task 'pre-commit-review' not found locally.
Found in GitHub catalog: tasks/git-workflow/pre-commit-review
Downloading...

✓ Cached to ~/.config/start/assets/tasks/git-workflow/
✓ Added to global config as 'pre-commit-review'

Running task 'pre-commit-review'...
[task executes]
```

### Example 2: Add to Local Config

```bash
$ start task pre-commit-review --local

Task 'pre-commit-review' not found locally.
Found in GitHub catalog: tasks/git-workflow/pre-commit-review
Downloading...

✓ Cached to ~/.config/start/assets/tasks/git-workflow/
✓ Added to local config as 'pre-commit-review'

Running task 'pre-commit-review'...
[task executes]
```

### Example 3: Downloads Disabled (asset_download = false)

```bash
$ start task pre-commit-review

Error: Task 'pre-commit-review' not found

  ✗ Not in local config (.start/tasks.toml)
  ✗ Not in global config (~/.config/start/tasks.toml)
  ✗ Not in asset cache (~/.config/start/assets/)
  ⚠ GitHub download disabled (asset_download = false)

To resolve:
  - Enable: start task pre-commit-review --asset-download
  - Add manually: start assets add
  - Browse catalog: start assets add
```

### Example 4: Override with Flag

```bash
$ start task pre-commit-review --asset-download=false

Error: Task 'pre-commit-review' not found

  ✗ Not in local config
  ✗ Not in global config
  ✗ Not in asset cache
  ⚠ Download disabled by --asset-download=false flag

To resolve:
  - Remove flag to allow download
  - Add manually: start assets add
```

### Example 5: Found in Cache

```bash
$ start task pre-commit-review

Using cached asset: pre-commit-review
Running task 'pre-commit-review'...
[task executes immediately]
```

**Note:** Cached assets are used immediately without prompting.

### Example 6: Offline (No Network)

```bash
$ start task pre-commit-review

Error: Task 'pre-commit-review' not found

  ✗ Not in local config
  ✗ Not in global config
  ✗ Not in asset cache
  ⚠ Cannot check GitHub catalog (offline)

To resolve:
  - Check spelling: 'pre-commit-review'
  - Add manually when online: start assets add
  - Use a configured task: start config task list
```

### Example 7: Not in Catalog

```bash
$ start task nonexistent-task

Error: Task 'nonexistent-task' not found

  ✗ Not in local config
  ✗ Not in global config
  ✗ Not in asset cache
  ✗ Not in GitHub catalog

To resolve:
  - Check spelling: 'nonexistent-task'
  - Browse available: start assets add
  - Create custom: start assets add nonexistent-task
```

## Implementation Details

### Cache Lookup

```go
func findInCache(assetType, name string) string {
    // Cache structure: ~/.config/start/assets/{type}/{category}/{name}.toml
    // We don't know category, so glob all categories
    pattern := filepath.Join(
        assetPath,
        assetType,
        "*",  // Any category
        fmt.Sprintf("%s.toml", name),
    )

    matches, err := filepath.Glob(pattern)
    if err != nil || len(matches) == 0 {
        return ""
    }

    // Return first match (should only be one)
    return matches[0]
}
```

### Config Addition

```go
func addToConfig(asset *Asset, scope string) error {
    var configPath string
    if scope == "local" {
        configPath = ".start/tasks.toml"
        if !fileExists(configPath) {
            return fmt.Errorf("no local config found (.start/config.toml)")
        }
    } else {
        configPath = "~/.config/start/tasks.toml"
    }

    // Load existing config
    config := loadTOML(configPath)

    // Add task (inline all content)
    config.Tasks[asset.Name] = asset

    // Write back
    return saveTOML(configPath, config)
}
```

### In-Memory Catalog Cache

```go
var catalogCache struct {
    tree      *GitHubTree
    timestamp time.Time
    mu        sync.RWMutex
}

func getCatalog() *GitHubCatalog {
    catalogCache.mu.RLock()
    if catalogCache.tree != nil {
        catalogCache.mu.RUnlock()
        return catalogCache.tree  // Use cached
    }
    catalogCache.mu.RUnlock()

    // Fetch from GitHub
    catalogCache.mu.Lock()
    defer catalogCache.mu.Unlock()

    tree := fetchGitHubTree()  // See DR-034
    catalogCache.tree = tree
    catalogCache.timestamp = time.Now()

    return tree
}
```

## Benefits

**Simple and predictable:**
- ✅ Clear priority order (local > global > cache > GitHub)
- ✅ First match wins (no complex merging)
- ✅ User config always takes precedence

**Flexible download control:**
- ✅ Global setting for default behavior
- ✅ Per-command flag override
- ✅ Explicit control over network usage

**Offline-friendly:**
- ✅ Works without network if asset configured or cached
- ✅ Clear error messages when network needed
- ✅ Graceful degradation per DR-026

**Cache transparency:**
- ✅ Cached assets used immediately
- ✅ No prompts for assets user already downloaded
- ✅ Cache = implementation detail

## Resolution vs Discovery (DR-041)

This resolution algorithm applies to **command execution** (running tasks, roles, agents).

**Asset discovery commands** (`start assets search/browse/info`) work differently:

| Aspect | Resolution (this DR) | Discovery (DR-040, DR-041) |
|--------|---------------------|---------------------------|
| **Purpose** | Execute configured asset | Find assets in catalog |
| **Commands** | `start task <name>`, `start --role <name>` | `start assets search`, `start assets browse` |
| **Search sources** | Local → Global → Cache → GitHub | GitHub catalog only |
| **Match algorithm** | Exact match → prefix (DR-038) | Substring matching (DR-040) |
| **Search fields** | Name only | Name, path, description, tags |
| **Uses index.csv** | No | Yes (DR-039) |

**Rationale for GitHub-only search:**
- Discovery is about **exploring what's available** in the catalog
- Local/global configs are **already known** to the user
- Cache is a subset of GitHub catalog
- Searching GitHub provides complete, fresh catalog view

**Example distinction:**
```bash
# Resolution (local → global → cache → GitHub)
start task pre-commit   # Prefix match, checks all sources

# Discovery (GitHub catalog only)
start assets search "commit"  # Substring match, GitHub only
```

After discovering an asset via `start assets`, it can be added to config and will then be found via normal resolution.

## Trade-offs Accepted

**Network dependency for new assets:**
- ❌ First use of catalog asset requires network
- **Mitigation:** Clear error messages, manual config always possible

**No automatic config detection:**
- ❌ Can't auto-detect if user wants global vs local
- **Mitigation:** Sensible default (global), `--local` flag for override

**Cache glob on every lookup:**
- ❌ Glob all categories to find asset in cache
- **Mitigation:** Filesystem cache is small, glob is fast

**In-memory catalog doesn't persist:**
- ❌ Catalog re-fetched on each CLI invocation
- **Mitigation:** Single API call, in-memory for session only

## Related Decisions

- [DR-031](./dr-031-catalog-based-assets.md) - Catalog architecture (resolution context)
- [DR-032](./dr-032-asset-metadata-schema.md) - Metadata schema (SHA for cache)
- [DR-034](./dr-034-github-catalog-api.md) - GitHub API (catalog fetching)
- [DR-036](./dr-036-cache-management.md) - Cache structure (cache lookup)
- [DR-026](./dr-026-offline-behavior.md) - Offline fallback (network error handling)
- [DR-038](./dr-038-flag-value-resolution.md) - Prefix matching for flag values (extends exact match with prefix support)
- [DR-039](./dr-039-catalog-index.md) - Catalog index file (used by discovery commands, not resolution)
- [DR-040](./dr-040-substring-matching.md) - Substring matching algorithm (used by discovery, not resolution)
- [DR-041](./dr-041-asset-command-reorganization.md) - Asset command reorganization (distinction between resolution and discovery)
