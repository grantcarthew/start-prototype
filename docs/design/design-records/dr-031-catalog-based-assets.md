# DR-031: Catalog-Based Asset Architecture

**Date:** 2025-01-10
**Status:** Accepted
**Category:** Asset Management

## Decision

Transform asset distribution from a bulk download model to a catalog-driven system where assets are discovered via GitHub, downloaded on-demand, cached locally, and lazily loaded at runtime.

## What This Means

### Core Architecture Shift

**Previous model (DR-014, DR-015):**
- Bulk download all assets during `start init` or `start assets update`
- Track versions in `asset-version.toml`
- All-or-nothing update mechanism

**New catalog model:**
- GitHub repository IS the catalog database
- Download individual assets on first use (lazy loading)
- Cache locally for offline use (`~/.config/start/assets/`)
- User config is separate from cached assets
- In-memory catalog cache (no disk tracking)

### Key Principles

1. **GitHub as Source of Truth** - Assets live in GitHub repo, not bundled with binary
2. **Lazy Loading** - Download on first use, not upfront
3. **Filesystem = State** - No tracking files; if cached, it exists
4. **On-Demand Installation** - Browse and install only what you need
5. **Hash-Based Versioning** - SHA comparison for update detection
6. **Offline-Friendly** - Cached assets work offline, manual config always possible

### User Workflows

**Browse and install (DR-041):**
```bash
start assets browse                      # Interactive catalog browser
start assets add "pre-commit"            # Search and install by query
start assets search "commit"             # Search by description/tags
start assets info "pre-commit-review"    # Preview before installing
```

**Lazy loading:**
```bash
start task pre-commit-review  # Not in config? Download from GitHub and add
```

**Update cached assets (DR-041):**
```bash
start assets update                      # Check cached assets for updates via SHA comparison
start assets update "git"                # Update only matching assets
```

**Note:** Original commands (`start config task add`, `start assets update`) deprecated in favor of unified `start assets` command suite. See [DR-041](./dr-041-asset-command-reorganization.md).

### Configuration Structure

**Multi-file configuration:**
```
~/.config/start/
├── config.toml      # Settings only
├── roles.toml       # All role definitions
├── tasks.toml       # All task definitions
├── agents.toml      # All agent configurations
├── contexts.toml    # All context configurations
└── assets/          # Cached catalog assets
    ├── tasks/
    │   └── git-workflow/
    │       ├── pre-commit-review.toml
    │       ├── pre-commit-review.md
    │       └── pre-commit-review.meta.toml
    ├── roles/
    ├── agents/
    └── contexts/
```

**Settings (config.toml):**
```toml
[settings]
default_agent = "claude"
default_role = "default"
log_level = "normal"
shell = "bash"
command_timeout = 30
asset_download = true  # Download from GitHub if asset not found
asset_path = "~/.config/start/assets"
asset_repo = "grantcarthew/start"
```

**Task configuration (tasks.toml):**
```toml
# Downloaded from catalog - inlined completely
[tasks.pre-commit-review]
alias = "pcr"
description = "Review staged changes before committing"
command = "git diff --staged"
prompt_file = "~/.config/start/assets/tasks/git-workflow/pre-commit-review.md"

# Custom user task - same structure
[tasks.my-review]
alias = "mr"
description = "My custom review"
command = "git diff"
prompt = "Check for issues"
```

### Asset Resolution Order

When user runs `start task <name>`:

1. **Local config** (`.start/tasks.toml`)
2. **Global config** (`~/.config/start/tasks.toml`)
3. **Asset cache** (`~/.config/start/assets/tasks/`)
4. **GitHub catalog** (query, prompt, download if `asset_download = true`)
5. **Error** - Not found anywhere

### Cache Behavior

**Cache is invisible:**
- Automatically populated when downloading from GitHub
- Updated via `start assets update` (SHA-based detection)
- Never auto-cleaned (files are tiny)
- No user-facing cache management commands
- User can manually delete `~/.config/start/assets/` if desired

**Cache contains:**
- Asset content files (`.toml`, `.md`)
- Sidecar metadata files (`.meta.toml`)
- SHA tracking for update detection

**Cache does NOT contain:**
- User modifications (those go in config files)
- State tracking files (filesystem IS the state)

### Update Workflow

```bash
$ start assets update

Checking for asset updates...
  ✓ tasks/git-workflow/pre-commit-review (updated v1.0 → v1.1)
  ✓ roles/general/code-reviewer (up to date)

Cache updated with 1 new version.

Note: Your task configurations are unchanged.
Review changes and manually update tasks.toml if desired.
```

Your existing task configurations are **never automatically overwritten**. Updates only affect the cache.

## Benefits

**For Users:**
- ✅ **Immediate value** - Browse and install curated assets
- ✅ **Discoverable** - Interactive catalog browsing and search (DR-041)
- ✅ **On-demand** - Only download what you use
- ✅ **Fresh** - Check for updates anytime via `start assets update`
- ✅ **Customizable** - Mix catalog + custom assets seamlessly
- ✅ **Offline-friendly** - Cached assets work offline
- ✅ **No tracking** - User config is self-contained

**For Project:**
- ✅ **Extensible** - Add asset types easily
- ✅ **Scalable** - Can grow to hundreds of assets
- ✅ **Maintainable** - Update assets without releases
- ✅ **Community-ready** - Others can contribute assets
- ✅ **Focused** - Binary is code, content is assets
- ✅ **Simple** - No manifest files, no version tracking

## Trade-offs Accepted

**Network dependency:**
- ❌ Initial download requires network access
- ❌ Browsing catalog requires network
- **Mitigation:** Cached assets work offline, manual config always possible

**Potential duplication:**
- ❌ Same asset content in cache and config
- **Mitigation:** Files are tiny text files, acceptable trade-off for simplicity

**No automatic overwrites of user config:**
- ❌ `start assets update` does not change your configured tasks. You must manually apply asset updates to your `tasks.toml` if desired.
- **Mitigation:** Preserves user customizations, prevents breaking changes.

**GitHub dependency:**
- ❌ Relies on GitHub API availability
- **Mitigation:** Graceful degradation per DR-026, raw.githubusercontent.com has no rate limit

## Implementation Notes

### Resolution Algorithm

```go
func ResolveAsset(assetType, name string, opts ResolveOptions) (*Asset, error) {
    // 1. Local config
    if asset := localConfig.Get(assetType, name); asset != nil {
        return asset, nil
    }

    // 2. Global config
    if asset := globalConfig.Get(assetType, name); asset != nil {
        return asset, nil
    }

    // 3. Asset cache
    if asset := loadFromCache(assetType, name); asset != nil {
        return asset, nil
    }

    // 4. Check if downloads allowed
    if !opts.AssetDownload {
        return nil, ErrNotFoundNoDownload
    }

    // 5. GitHub catalog (in-memory cached tree)
    catalog := getCatalog()  // See DR-034
    asset := downloadAsset(catalog.Find(assetType, name))
    cacheAsset(asset)

    // 6. Add to config (global by default, local if --local)
    scope := ternary(opts.Local, "local", "global")
    addToConfig(asset, scope)

    return asset, nil
}
```

### Command Flags

```bash
start task <name> [flags]

Flags:
  --local                  Add downloaded asset to local config (default: global)
  --asset-download[=bool]  Download from GitHub if not found (default: from settings)
```

### Multi-File Config Loading

```go
type Config struct {
    Settings map[string]interface{}
    Tasks    map[string]*Task
    Agents   map[string]*Agent
    Contexts map[string]*Context
}

func LoadConfig(dir string) (*Config, error) {
    cfg := &Config{}
    cfg.Settings = loadTOML(filepath.Join(dir, "config.toml"))
    cfg.Tasks = loadTOML(filepath.Join(dir, "tasks.toml"))
    cfg.Agents = loadTOML(filepath.Join(dir, "agents.toml"))
    cfg.Contexts = loadTOML(filepath.Join(dir, "contexts.toml"))
    return cfg
}
```

## Related Decisions

**Supersedes/updates:**
- [DR-014](./dr-014-github-tree-api.md) - GitHub Tree API (now for catalog browsing)
- [DR-015](./dr-015-atomic-updates.md) - Atomic updates (now per-asset)
- [DR-016](./dr-016-asset-discovery.md) - Asset discovery (now interactive browsing)
- [DR-019](./dr-019-task-loading.md) - Task loading (now includes cache resolution)
- [DR-023](./dr-023-asset-staleness-check.md) - Staleness checking (now per-asset SHA)

**New DRs (detailed decisions):**
- [DR-032](./dr-032-asset-metadata-schema.md) - Sidecar metadata format
- [DR-033](./dr-033-asset-resolution-algorithm.md) - Resolution priority and behavior
- [DR-034](./dr-034-github-catalog-api.md) - GitHub API strategy and caching
- [DR-035](./dr-035-interactive-browsing.md) - TUI catalog browsing
- [DR-036](./dr-036-cache-management.md) - Cache structure and behavior
- [DR-037](./dr-037-asset-updates.md) - Update mechanism and SHA comparison
- [DR-039](./dr-039-catalog-index.md) - Catalog index file (CSV schema for fast search)
- [DR-040](./dr-040-substring-matching.md) - Substring matching algorithm for asset search
- [DR-041](./dr-041-asset-command-reorganization.md) - Unified `start assets` command suite

**Consistent with:**
- [DR-026](./dr-026-offline-behavior.md) - Offline fallback (cache works offline)
- [DR-025](./dr-025-no-automatic-checks.md) - No automatic checks (explicit updates)

## Future Considerations

**Community contributions:**
- Asset submission process (PRs to catalog repo)
- Quality control and review process
- Asset namespacing for user repos

**Additional features:**
- Search/filter functionality across catalog
- Asset dependencies (tasks requiring specific roles)
- Workspace templates (preset configurations)
- Additional asset types (metaprompts, snippets)

**Current stance:** Ship v1 with minimal viable set (28 assets across 4 types), iterate based on user feedback.
