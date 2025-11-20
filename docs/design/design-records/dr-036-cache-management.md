# DR-036: Cache Management

- Date: 2025-01-10
- Status: Accepted
- Category: Asset Management

## Problem

Asset caching needs a management strategy. The approach must address:

- Cache location (where to store cached assets)
- User visibility (can users see what's cached)
- Management commands (list, clean, prune operations)
- Cleanup policies (age-based, size-based, manual only)
- Cache operations (read, write, update, delete)
- Persistence (per-machine vs synchronized)
- Configuration (customizable location)
- Troubleshooting (how to debug cache issues)
- Complexity vs benefit (KISS principle for tiny files)

## Decision

Asset cache is invisible implementation detail with no user-facing management commands. Cache stored at ~/.config/start/assets/ with manual deletion only.

Key aspects:

Location: ~/.config/start/assets/

- Separate subdirectory under config directory
- Configurable via asset_path setting
- Standard directory structure

No management commands:

- No start cache list
- No start cache clean
- No start cache info
- No start cache prune

Manual deletion:

- rm -rf ~/.config/start/assets/ to clear cache
- Assets re-download automatically on next use
- start assets update to bulk re-download

No cleanup policies:

- Cache persists until manually deleted
- No age-based cleanup
- No size-based cleanup
- Consistent with no automatic operations principle

Filesystem as state:

- No tracking files
- File existence IS the cache state
- Standard directory structure
- Human-readable text files

Configurable location:

- Default: ~/.config/start/assets/
- Customizable: asset_path setting in config.toml
- Useful for network drives, shared caches, testing

## Why

Files are tiny (no need for management):

- Asset files typically < 5KB each
- 28 assets = ~62KB total
- 200 assets = ~600KB total
- Negligible disk space on modern systems
- No need for size limits or cleanup policies
- Growth is minimal even with hundreds of assets

KISS principle:

- Don't build features speculatively
- Cache management adds complexity for minimal benefit
- Simple manual deletion works fine
- Users rarely need to manage cache
- Keeps binary small and code simple

Filesystem as state is simpler:

- No tracking files (no drift possible)
- File existence IS the cache state
- Easy to inspect with standard tools (ls, cat, tree)
- Easy to reset (delete directory)
- No database or index needed
- Reliable and transparent

Fast re-download makes cache disposable:

- Raw.githubusercontent.com has no rate limits
- Assets are small (< 5KB typically)
- Re-downloading is fast
- Network cost minimal
- Cache is disposable implementation detail
- No penalty for clearing cache

Transparency via filesystem:

- Human-readable text files
- Standard directory structure
- Can inspect manually with ls, cat, tree
- No special tools needed
- Debugging is straightforward
- Users understand how it works

## Trade-offs

Accept:

- No cache visibility commands (can't see what's cached via CLI, use ls or tree instead)
- No automatic cleanup (cache grows over time, but files tiny so negligible growth)
- No cache statistics (can't see cache size via CLI, can add to start doctor if requested)
- Manual deletion only (can't selectively remove via CLI, but cache is disposable so rm -rf works)
- Per-machine cache (not synchronized across machines, but re-download is fast and painless)

Gain:

- Simple implementation (no cache commands, no management logic, no size limits, no cleanup policies)
- Transparent operation (standard filesystem directory, human-readable text files, easy to inspect)
- Reliable state (filesystem IS the state, no tracking files to corrupt, easy to reset)
- Efficient lookups (filesystem cache fast, no database needed, minimal disk space < 1MB)
- Easy troubleshooting (use standard tools, delete and re-download, no complex debugging)
- KISS compliance (don't build features speculatively, cache management not needed for tiny files)

## Alternatives

Add cache management commands:

Example: Implement start cache list, start cache clean, start cache info

```bash
start cache list       # List cached assets
start cache clean      # Remove unused assets
start cache info       # Show cache statistics
start cache prune --older-than 30d  # Remove old assets
```

Pros:

- User visibility (can see what's cached via CLI)
- Selective cleanup (remove specific assets)
- Statistics available (cache size, asset count)
- Familiar pattern (like npm cache, yarn cache)

Cons:

- Added complexity (implement multiple commands, maintain cache metadata)
- Minimal benefit (files are tiny < 1MB, cleanup not needed)
- Over-engineering (cache management for negligible disk space)
- Violates KISS (building features speculatively)
- More code to maintain (commands, help text, error handling)

Rejected: Files are tiny (< 1MB even for large catalogs). Cache management adds complexity for no real benefit. Manual deletion sufficient.

Automatic cleanup policies:

Example: Age-based or size-based automatic cleanup

```toml
[settings]
cache_max_size = "10MB"
cache_max_age_days = 90
```

- Automatically delete assets older than 90 days
- Automatically delete oldest assets when size exceeds 10MB

Pros:

- Automatic maintenance (no user intervention)
- Bounded growth (cache won't grow indefinitely)
- Familiar pattern (like browser cache)

Cons:

- Violates no automatic operations principle
- Added complexity (implement cleanup logic, track timestamps)
- Unnecessary (files tiny, growth negligible)
- Surprising behavior (assets disappear unexpectedly)
- Can break workflows (delete assets user relies on)

Rejected: Violates no automatic operations principle. Unnecessary for tiny files. Manual deletion simpler and more predictable.

SQLite database for cache metadata:

Example: Track cache in database instead of filesystem

```
~/.config/start/cache.db (SQLite database)
- Table: cached_assets (name, type, category, sha, size, downloaded_at)
- Assets stored separately
```

Pros:

- Fast queries (can query cache by various fields)
- Statistics built-in (count, size, dates)
- Easy to add features (cache commands, cleanup policies)

Cons:

- Over-engineering (database for tiny text files)
- Added complexity (database schema, migrations)
- Another thing to corrupt (database can break)
- Drift possible (database vs filesystem state)
- Overkill for simple cache

Rejected: Filesystem IS the state. Database overkill for tiny files. KISS principle violated.

Synchronized cache across machines:

Example: Cloud-synced cache directory

```toml
[settings]
asset_path = "~/Dropbox/start-cache"  # Synced via Dropbox
```

- Cache synced across all machines
- One download, available everywhere

Pros:

- Convenience (download once, available everywhere)
- Faster setup (new machine has cache immediately)
- Bandwidth savings (no re-downloads)

Cons:

- Sync conflicts possible (multiple machines downloading simultaneously)
- Dependency on sync service (Dropbox, iCloud, etc.)
- More complex setup (user must configure)
- Can cause issues (sync service down, conflicts)

Rejected: Re-download is fast (no rate limits, tiny files). Sync adds complexity and potential issues. Local cache per machine is simpler.

## Structure

Cache location:

Default: ~/.config/start/assets/

- Subdirectory under config directory
- Configurable via asset_path setting
- Standard directory structure

Configuration:

```toml
# config.toml
[settings]
asset_path = "~/.config/start/assets"  # Default location
```

Custom location:

```toml
# config.toml
[settings]
asset_path = "/custom/path/to/assets"  # Custom location
```

Directory structure:

```
~/.config/start/assets/
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
├── roles/
│   ├── general/
│   │   ├── code-reviewer.md
│   │   ├── code-reviewer.meta.toml
│   │   ├── default.md
│   │   └── default.meta.toml
│   └── languages/
│       ├── go-expert.md
│       └── go-expert.meta.toml
├── agents/
│   └── anthropic/
│       ├── claude.toml
│       └── claude.meta.toml
└── contexts/
    └── reference/
        ├── environment.toml
        └── environment.meta.toml
```

Structure pattern: {cache}/{type}/{category}/{name}.{ext}

- {cache}: ~/.config/start/assets/ (or custom via asset_path)
- {type}: tasks, roles, agents, contexts
- {category}: git-workflow, general, anthropic, reference, etc.
- {name}: Asset name (matches metadata name field)
- {ext}: .toml (content), .md (prompt file), .meta.toml (metadata)

Cache operations:

Automatic population:

- User runs: start task pre-commit-review
- Asset not in config/cache
- Downloads from GitHub
- Caches to ~/.config/start/assets/tasks/git-workflow/
- Adds to config
- Uses cached version on subsequent runs

Automatic updates:

- User runs: start assets update
- Checks GitHub for newer SHAs (via index.csv)
- Downloads updated assets to cache
- User config remains unchanged
- User must manually apply updates to config

Manual deletion:

- User runs: rm -rf ~/.config/start/assets/
- Cache cleared completely
- Assets re-download automatically on next use
- OR user runs: start assets update to bulk re-download

Cache persistence:

Per-machine only:

- Not synchronized across machines
- Not backed up automatically
- Each machine has independent cache

Re-download is fast:

- Raw.githubusercontent.com (no rate limits)
- Assets small (< 5KB typically)
- Network cost minimal
- Cache is disposable

User config portable:

```toml
# tasks.toml references cached assets
[tasks.pre-commit-review]
prompt_file = "~/.config/start/assets/tasks/git-workflow/pre-commit-review.md"

# On new machine:
# 1. Copy config files
# 2. Run tasks - cache populated automatically
# 3. OR: start assets update (download all assets referenced in config)
```

## Usage Examples

Automatic cache population:

```bash
$ start task pre-commit-review

# If not in cache:
# 1. Downloads from GitHub
# 2. Caches to ~/.config/start/assets/tasks/git-workflow/
# 3. Adds to config
# 4. Runs task

# Subsequent runs use cache (no download)
```

Update cache:

```bash
$ start assets update

Downloading catalog index...
✓ Loaded index (46 assets)

Checking for asset updates...
  ✓ tasks/git-workflow/pre-commit-review (updated v1.0 → v1.1)
  ✓ roles/general/code-reviewer (up to date)

Cache updated with 1 new version.

Note: Your task configurations are unchanged.
Review changes and manually update tasks.toml if desired.
```

Manual cache inspection:

```bash
# List cached assets
$ find ~/.config/start/assets -name "*.toml" -not -name "*.meta.toml"
/Users/user/.config/start/assets/tasks/git-workflow/pre-commit-review.toml
/Users/user/.config/start/assets/tasks/git-workflow/pr-ready.toml
/Users/user/.config/start/assets/roles/general/code-reviewer.toml

# Or use tree
$ tree ~/.config/start/assets
/Users/user/.config/start/assets
├── tasks
│   └── git-workflow
│       ├── pre-commit-review.toml
│       ├── pre-commit-review.md
│       └── pre-commit-review.meta.toml
└── roles
    └── general
        ├── code-reviewer.md
        └── code-reviewer.meta.toml

# View metadata
$ cat ~/.config/start/assets/tasks/git-workflow/pre-commit-review.meta.toml
[metadata]
name = "pre-commit-review"
description = "Review staged changes before committing"
tags = ["git", "review", "quality", "pre-commit"]
sha = "a1b2c3d4e5f6789012345678901234567890abcd"
created = "2025-01-10T00:00:00Z"
updated = "2025-01-10T12:30:00Z"
```

Manual cache deletion:

```bash
# Clear entire cache
$ rm -rf ~/.config/start/assets/

# Assets re-download automatically on next use
$ start task pre-commit-review
Downloading...
✓ Cached to ~/.config/start/assets/tasks/git-workflow/
✓ Running task...
```

Troubleshooting - cache corrupted:

```bash
# Delete and re-download
$ rm -rf ~/.config/start/assets/
$ start assets update

Downloading catalog index...
Checking for asset updates...
✓ Downloaded 28 assets

Cache populated.
```

Troubleshooting - asset not found:

```bash
# Check if cached
$ ls ~/.config/start/assets/tasks/git-workflow/
pre-commit-review.toml
pre-commit-review.md
pre-commit-review.meta.toml

# If missing, re-download
$ start assets update git-workflow
```

Custom cache location:

```toml
# config.toml
[settings]
asset_path = "/mnt/shared/start-cache"  # Network drive
```

```bash
# Cache uses custom location
$ start task pre-commit-review
✓ Cached to /mnt/shared/start-cache/tasks/git-workflow/
```

Cache size estimates:

Minimal viable set (28 assets):

```
Tasks (12): 12 × 3KB = 36KB
Roles (8): 8 × 2KB = 16KB
Agents (6): 6 × 1KB = 6KB
Contexts (2): 2 × 2KB = 4KB
Total: ~62KB (tiny!)
```

Large catalog (200 assets):

```
200 × 3KB = 600KB
Still trivial for modern systems
No cleanup needed
```

## Updates

- 2025-01-17: Initial version aligned with schema; removed implementation code, Related Decisions, and Future Considerations sections
