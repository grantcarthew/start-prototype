# DR-031: Catalog-Based Asset Architecture

- Date: 2025-01-10
- Status: Accepted
- Category: Asset Management

## Problem

Asset distribution needs a scalable, user-friendly system. The design must address:

- Distribution model (how users get assets)
- Discovery mechanism (how users find available assets)
- Installation workflow (download all upfront vs on-demand)
- Update management (version tracking, update detection, user control)
- Offline usability (work without network after initial setup)
- Cache management (where assets stored, how updates tracked)
- User customization (how catalog assets relate to user config)
- Storage efficiency (avoid duplicating content)
- Scalability (support growing catalog without performance issues)
- GitHub dependency (API limits, availability, fallback strategies)
- Search performance (fast metadata-rich searching across catalog)

## Decision

Use catalog-driven asset distribution where assets are discovered via GitHub catalog index, downloaded on-demand, cached locally, and lazily loaded at runtime.

Core architecture:

- GitHub repository IS the catalog database
- Catalog index (assets/index.csv) enables fast metadata-rich search
- Download individual assets on first use (lazy loading)
- Cache locally at ~/.config/start/assets/ for offline use
- User config separate from cached assets
- No local tracking files (filesystem IS the state)
- SHA-based update detection (no version manifest)

Key principles:

1. GitHub as source of truth: Assets live in GitHub repo, not bundled with binary
2. Index-driven discovery: Single CSV index file (assets/index.csv) for fast search
3. Lazy loading: Download on first use, not upfront
4. Filesystem as state: No tracking files; if cached file exists, asset is cached
5. On-demand installation: Browse and install only what you need
6. Hash-based versioning: SHA comparison for update detection
7. Offline-friendly: Cached assets work offline, manual config always possible

Catalog index system:

Index file structure:
- Location: assets/index.csv in catalog repository
- Format: CSV with columns (type, category, name, description, tags, bin, sha, size, created, updated)
- Downloaded fresh on every search/browse operation (no local caching)
- Small file (~10-50KB for hundreds of assets)
- Enables fast metadata-rich searching without downloading individual .meta.toml files

Fallback mechanism:
- If index.csv unavailable: Fall back to GitHub Tree API for directory listing
- Tree API provides name/path only (no descriptions/tags)
- Degraded but functional search capability

Asset resolution order (when user runs start task <name>):

1. Local config (.start/tasks.toml)
2. Global config (~/.config/start/tasks.toml)
3. Asset cache (~/.config/start/assets/tasks/)
4. GitHub catalog (query index.csv, prompt, download if asset_download enabled)
5. Error (not found anywhere)

Cache behavior:

Cache is invisible to users:
- Automatically populated when downloading from GitHub
- Updated via start assets update (SHA-based detection)
- Never auto-cleaned (files are tiny text files)
- No user-facing cache management commands
- User can manually delete cache directory if desired

Cache contains:
- Asset content files (.toml, .md)
- Sidecar metadata files (.meta.toml)
- SHA tracking for update detection

Cache does NOT contain:
- User modifications (those go in config files)
- State tracking files (filesystem IS the state)
- Catalog index (downloaded fresh each time)

Update workflow:

- start assets update checks cached assets for updates
- Compares local SHA with GitHub SHA (from index.csv)
- Downloads updated assets to cache
- Never automatically overwrites user config
- User reviews changes and manually updates config if desired

## Why

Lazy loading provides better user experience:

- Users download only what they use (no large upfront download)
- Immediate value from browsing catalog
- Faster initial setup (no bulk download during init)
- Progressive asset discovery (learn as you go)

Index-driven search scales efficiently:

- Single ~10-50KB download vs 200+ individual .meta.toml requests
- Fast in-memory search with rich metadata (descriptions, tags)
- No API rate limits (uses raw.githubusercontent.com)
- Always fresh data (downloaded on each use)
- Graceful degradation (fallback to Tree API if index unavailable)

On-demand installation scales:

- Catalog can grow to hundreds of assets without impacting performance
- Users not forced to download entire catalog
- Add new assets without affecting existing users
- Bandwidth-efficient for users who need few assets

GitHub as catalog enables rapid iteration:

- Update assets without CLI releases
- Community can contribute assets via PRs
- Asset improvements immediately available
- Simple content management (files in repo, index generated via start assets index)

Cache separation maintains flexibility:

- User config independent of catalog
- Mix catalog and custom assets seamlessly
- Catalog updates don't overwrite user customizations
- Clear separation of concerns (config vs cache)

Offline support after first use:

- Cached assets work offline
- Manual config always possible
- No network required for execution
- Graceful degradation when network unavailable

Simple implementation and maintenance:

- No version tracking files needed
- Filesystem IS the state (if file exists, cached)
- SHA comparison for updates (simple and reliable)
- CSV index for fast search (standard format)
- Text files (easy to inspect and debug)

## Trade-offs

Accept:

- Network dependency for discovery (browsing/searching catalog requires network, but cached assets work offline)
- Index downloaded on every search (not cached locally, ~10-50KB per search, acceptable for always-fresh data)
- Potential content duplication (same asset in cache and config, acceptable for tiny text files)
- No automatic config updates (user must manually apply asset updates to their config)
- GitHub dependency (relies on GitHub availability, but graceful degradation via Tree API fallback)
- Cache grows over time (never auto-cleaned, but files are tiny)
- More complex resolution (check multiple locations: local, global, cache, GitHub)
- Manual index regeneration (contributors must run start assets index after changes)

Gain:

- Immediate user value (browse and install curated assets on-demand)
- Highly discoverable (fast metadata-rich searching via index.csv)
- Download only what you use (no forced bulk download)
- Fresh content anytime (check for updates via start assets update)
- Seamlessly customizable (mix catalog and custom assets)
- Offline-friendly (cached assets work without network)
- No tracking files (user config is self-contained, simple)
- Extensible architecture (add asset types easily)
- Scalable to hundreds of assets (single index download, fast search)
- Maintainable (update assets without CLI releases, regenerate index)
- Community-ready (others can contribute assets, index generated automatically)
- Focused separation (binary is code, content is assets)
- Graceful degradation (Tree API fallback if index unavailable)

## Alternatives

Bundled assets in binary:

Example: Include all assets in CLI binary at build time
- Assets compiled into executable
- No network dependency ever
- Simple deployment

Pros:
- Zero network dependency (100% offline)
- Guaranteed availability (assets always present)
- Fast access (no download or cache lookup)
- Simple distribution (single binary)

Cons:
- Large binary size (grows with asset count)
- Can't update assets without CLI release
- No community contributions between releases
- Inflexible (users can't choose which assets to include)
- Couples content updates to code releases
- Binary grows indefinitely as catalog expands

Rejected: Couples content to code releases. Can't update assets independently. Binary size grows uncontrollably.

Plugin system with separate asset repositories:

Example: Users add asset repos like apt sources
```toml
[asset_sources]
official = "github.com/grantcarthew/start"
community = "github.com/start-community/assets"
personal = "github.com/myuser/my-assets"
```

Pros:
- Multiple asset sources (community ecosystems)
- Decentralized asset hosting
- Users can create private asset repos
- Flexible and extensible

Cons:
- Complex configuration (users must manage sources)
- Trust and security issues (which sources are safe?)
- Discovery fragmentation (assets spread across repos)
- Namespace collisions (same asset name in different repos)
- More implementation complexity
- Overwhelming for typical users

Rejected: Over-engineered for v1. Single trusted source is simpler and sufficient. Can add later if needed.

Download all .meta.toml files for search:

Example: Fetch all metadata files individually for search
- Query GitHub Tree API for asset list
- Download each .meta.toml file (200+ HTTP requests)
- Parse and search in memory

Pros:
- No index file needed (one less thing to maintain)
- Always absolutely fresh (no derived file)
- Source of truth directly accessed

Cons:
- Extremely slow (200+ HTTP requests for large catalog)
- API rate limits (GitHub API has limits)
- Poor user experience (long wait for searches)
- Doesn't scale (worse as catalog grows)

Rejected: Performance unacceptable. Index.csv pattern is proven (npm, cargo, homebrew all use it).

## Structure

Configuration structure:

Directory layout:
```
~/.config/start/
├── config.toml      # Settings only
├── roles.toml       # All role definitions
├── tasks.toml       # All task definitions
├── agents.toml      # All agent configurations
├── contexts.toml    # All context configurations

~/.config/start/assets/      # Cached catalog assets (separate from config)
├── tasks/
│   └── git-workflow/
│       ├── pre-commit-review.toml
│       ├── pre-commit-review.md
│       └── pre-commit-review.meta.toml
├── roles/
├── agents/
└── contexts/
```

Settings configuration (config.toml):
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

Task configuration (tasks.toml):
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

Catalog index system:

Index file (assets/index.csv):
- Location: assets/index.csv in catalog repository
- Format: CSV with header row
- Columns: type, category, name, description, tags, bin, sha, size, created, updated
- Sorted: Alphabetically by type → category → name
- Generated: Via start assets index command (run by maintainers)
- Downloaded: Fresh on every search/browse (no local caching)
- Size: ~10-50KB for hundreds of assets

Search workflow:
1. Download index.csv from raw.githubusercontent.com
2. Parse CSV into in-memory structure
3. Search by substring matching (name, description, tags)
4. Display results with rich metadata
5. If index unavailable: Fall back to Tree API (name/path only)

Asset resolution algorithm:

When user runs start task <name>:

1. Check local config (.start/tasks.toml)
   - If found: use it (highest priority)
   - If not found: continue to step 2

2. Check global config (~/.config/start/tasks.toml)
   - If found: use it
   - If not found: continue to step 3

3. Check asset cache (~/.config/start/assets/tasks/)
   - If found: use cached asset
   - If not found: continue to step 4

4. Check if downloads allowed (asset_download setting)
   - If disabled: return error (not found, no download)
   - If enabled: continue to step 5

5. Query GitHub catalog
   - Download index.csv
   - Search for asset by name
   - If found: download asset files, cache them, add to config
   - If not found: return error (not found anywhere)

Cache behavior:

Automatic population:
- Assets downloaded from GitHub cached automatically
- Cache location: ~/.config/start/assets/
- Organized by asset type (tasks/, roles/, agents/, contexts/)

Update detection:
- start assets update downloads index.csv
- Compares SHA from index with SHA of cached files
- Downloads updated assets to cache
- Shows what changed (version, description)

Preservation of user config:
- Cache updates never overwrite user config files
- User must manually apply updates to tasks.toml if desired
- Preserves customizations and prevents breaking changes

No cleanup:
- Cache never auto-cleaned
- Files are tiny text files (storage not a concern)
- User can manually delete cache directory if needed

Command flags and options:

start task <name> [flags]:
- --local: Add downloaded asset to local config (default: global)
- --asset-download[=bool]: Download from GitHub if not found (default: from settings)

start assets browse:
- Open GitHub catalog in web browser

start assets add <query>:
- Search catalog and install matching asset (uses index.csv)

start assets update [query]:
- Update cached assets (all or matching query, uses index.csv for SHA comparison)

start assets search <query>:
- Search catalog by description/tags (uses index.csv)

start assets info <name>:
- Preview asset before installing (uses index.csv)

start assets index:
- Generate index.csv from .meta.toml files (maintainer command)

## Usage Examples

Browse and install assets:

```bash
$ start assets browse
# Opens web browser to GitHub catalog

$ start assets add "pre-commit"
# Search for assets matching "pre-commit", prompt to install

$ start assets search "commit"
# Show all assets related to "commit" (from index.csv)

$ start assets info "pre-commit-review"
# Show details about specific asset before installing
```

Lazy loading on first use:

```bash
$ start task pre-commit-review
# Asset not in config? Download from GitHub and add automatically

Task 'pre-commit-review' not found locally.

Found in catalog:
  Name: pre-commit-review
  Description: Review staged changes before committing
  Tags: git, review, quality, pre-commit

Download and add to config? [Y/n] y

✓ Downloaded pre-commit-review
✓ Cached to ~/.config/start/assets/tasks/git-workflow/
✓ Added to ~/.config/start/tasks.toml

Running task...
```

Update cached assets:

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

Update specific assets by query:

```bash
$ start assets update "git"

Downloading catalog index...
Checking assets matching "git"...
  ✓ tasks/git-workflow/pre-commit-review (updated)
  ✓ tasks/git-workflow/commit-message (up to date)

Cache updated.
```

Search with index (normal operation):

```bash
$ start assets search "commit"

Downloading catalog index...
✓ Loaded index (46 assets)

Found 3 matches:

tasks/git-workflow/commit-message
  Description: Generate conventional commit message
  Tags: git, commit, conventional

tasks/git-workflow/pre-commit-review
  Description: Review staged changes before committing
  Tags: git, review, quality, pre-commit

tasks/workflows/post-commit-hook
  Description: Post-commit validation workflow
  Tags: git, commit, hooks, validation
```

Search without index (fallback to Tree API):

```bash
$ start assets search "commit"

Downloading catalog index...
⚠ Index unavailable, using directory listing (limited metadata)

Found 3 matches:

tasks/git-workflow/commit-message
  (Metadata unavailable - not in index)

tasks/git-workflow/pre-commit-review
  (Metadata unavailable - not in index)

tasks/workflows/post-commit-hook
  (Metadata unavailable - not in index)
```

Configuration showing catalog and custom assets:

```toml
# tasks.toml
# Mix of catalog and custom assets

[tasks.pre-commit-review]  # From catalog
alias = "pcr"
description = "Review staged changes before committing"
command = "git diff --staged"
prompt_file = "~/.config/start/assets/tasks/git-workflow/pre-commit-review.md"

[tasks.my-review]  # Custom user task
alias = "mr"
description = "My custom review workflow"
command = "git diff"
prompt = "Check for security issues and code quality"
```

Offline behavior with cached assets:

```bash
# Network available - browse catalog
$ start assets browse
[Opens web browser to GitHub catalog]

# Download and cache
$ start task pre-commit-review
✓ Downloaded and cached to ~/.config/start/assets/

# Network unavailable - cached assets still work
$ start task pre-commit-review
[Executes from cache, no network needed]

# Trying to browse without network
$ start assets browse
Error: Cannot connect to GitHub
Catalog browsing requires network access.

Cached assets are still available:
  start task pre-commit-review
```

Maintainer workflow - generating index:

```bash
$ cd ~/projects/start-catalog-fork

$ start assets index

Validating repository structure...
✓ Git repository detected
✓ Assets directory found

Scanning assets/...
Found 46 assets

Sorting assets (type → category → name)...
Writing index to assets/index.csv...

✓ Generated index with 46 assets
Updated: assets/index.csv

Ready to commit:
  git add assets/
  git commit -m "Regenerate catalog index"
```

## Updates

- 2025-01-17: Initial version aligned with schema; incorporated index.csv system from DR-039
