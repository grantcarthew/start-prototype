# DR-022: Asset Branch Strategy

- Date: 2025-01-06
- Status: Updated by DR-031 and DR-039
- Category: Asset Management

Note: Core decision (main branch vs releases) remains valid. Implementation updated by DR-031 (catalog architecture) and DR-039 (catalog index from main branch).

## Problem

Asset catalog needs to know which GitHub branch to fetch assets from. The system must:

- Define which branch contains the authoritative asset catalog
- Support rapid iteration on asset content (tasks, roles, agents, contexts)
- Decouple asset updates from CLI binary releases
- Allow fixing typos and adding new assets without cutting releases
- Provide users with current assets without requiring CLI updates
- Balance stability with freshness for content vs code
- Work with catalog index file (index.csv) and individual assets
- Support different update cycles for content (frequent) vs binaries (stable)

## Decision

Asset catalog always fetched from **main branch**, not GitHub Releases.

Catalog operations using main branch:

- `start assets search` - downloads `assets/index.csv` from main
- `start assets add` - downloads catalog index and assets from main
- `start assets update` - checks for updates on main branch
- Lazy loading - downloads individual assets from main on first use

CLI binary versioning (separate concern):

- CLI releases use GitHub Releases (semantic versioning)
- Users update CLI via `brew upgrade` or `go install`
- Version checks compare against latest release (DR-021)

## Why

Assets are content not code:

- Roles: Markdown prompt documents
- Agents: TOML configuration templates
- Tasks: TOML task definitions
- Contexts: TOML context configurations
- Content updates have lower risk than binary code updates
- Benefit from rapid iteration without release overhead

Decoupled update cycles:

- CLI binary: Tied to GitHub Releases (stable, tested, versioned)
- Asset catalog: Tied to main branch commits (frequent, current)
- Users get content improvements immediately via `start assets update`
- No CLI release needed to fix typo in role prompt
- No CLI release needed to add new task template

Catalog index from main:

- `index.csv` fetched from main branch via raw.githubusercontent.com
- Always shows current catalog state
- No rate limits (raw content, not API)
- Individual assets also fetched from main
- Consistent branch across all catalog operations

Faster feedback loop:

- Iterate on task templates without cutting releases
- Add new roles as they're created
- Update agent configs for new model versions
- Fix documentation typos immediately
- Community contributions available faster

User control:

- Users control when they get updates (via `start assets update`)
- Not forced to take changes immediately
- Can review asset changes before updating
- Lazy loading means download on first use only

Simple mental model:

- Code (CLI binary) = GitHub Releases (stable)
- Content (assets) = main branch (current)
- Clear separation of concerns

## Trade-offs

Accept:

- Broken assets could reach users if merged to main
- No "stable" vs "bleeding edge" channel choice for assets
- Assets always from latest main (no version pinning)
- Must keep main branch stable for assets
- Content changes visible immediately after merge

Gain:

- Faster iteration on asset content without release overhead
- Users get improvements immediately
- Simple architecture (one branch, no channel complexity)
- Lower risk (content files vs executable binaries)
- Community contributions flow faster
- No release coordination needed for content updates

## Alternatives

Assets from GitHub Releases (matching CLI version):

```bash
# Assets tied to CLI release tags
GET /repos/grantcarthew/start/releases/latest
# Extract tag: v1.3.0
GET https://raw.githubusercontent.com/grantcarthew/start/v1.3.0/assets/index.csv
```

- Pro: Assets perfectly matched to CLI version
- Pro: More stable (only released versions)
- Pro: Version pinning built-in
- Con: Cannot update assets without cutting CLI release
- Con: Typo fixes require full release process
- Con: Slower feedback for content improvements
- Con: Couples asset updates to code releases
- Rejected: Too slow for content iteration

Assets from develop branch:

```bash
GET https://raw.githubusercontent.com/grantcarthew/start/develop/assets/index.csv
```

- Pro: Even more bleeding edge than main
- Pro: Could have main for stable, develop for experimental
- Con: Adds complexity (which branch to use?)
- Con: Users must choose stability level
- Con: Two catalogs to maintain
- Con: Main would become stale if develop is where changes go
- Rejected: Unnecessary complexity, keep main stable instead

Per-user branch selection:

```toml
[settings]
asset_branch = "main"  # or "develop", "stable", etc.
```

- Pro: User choice of stability level
- Pro: Power users can use bleeding edge
- Con: Configuration complexity
- Con: Multiple branches to maintain
- Con: Support burden (which branch did user use?)
- Con: Testing across multiple branches
- Rejected: Over-engineering for minimal benefit

Asset channels (stable/latest/develop):

```toml
[settings]
asset_channel = "stable"   # From releases
asset_channel = "latest"   # From main (default)
asset_channel = "develop"  # From develop branch
```

- Pro: Flexibility for different user preferences
- Pro: Could satisfy both conservative and experimental users
- Con: Multiple catalogs to maintain
- Con: Which channel gets tested?
- Con: Support complexity
- Con: Configuration burden on users
- Rejected: Assets are templates, one stable channel sufficient

Version pinning for assets:

```bash
start assets update --tag v1.2.0  # Pin to specific asset version
```

- Pro: Reproducible builds with specific asset versions
- Pro: Can rollback to known-good state
- Con: Requires tagging assets separately from CLI
- Con: More complex version management
- Con: Users must track asset versions
- Rejected: Assets change frequently, pinning adds overhead

## Structure

Catalog index fetch:

```bash
# Always from main branch
GET https://raw.githubusercontent.com/grantcarthew/start/main/assets/index.csv
```

Individual asset fetch:

```bash
# Download from main branch
GET https://raw.githubusercontent.com/grantcarthew/start/main/assets/tasks/git-workflow/pre-commit-review.toml
GET https://raw.githubusercontent.com/grantcarthew/start/main/assets/tasks/git-workflow/pre-commit-review.md
GET https://raw.githubusercontent.com/grantcarthew/start/main/assets/tasks/git-workflow/pre-commit-review.meta.toml
```

Constants in code:

```go
const (
    AssetRepository = "grantcarthew/start"
    AssetBranch     = "main"  // Always main
    AssetsBasePath  = "assets"
)

func FetchCatalogIndex(ctx context.Context) (*CatalogIndex, error) {
    url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/index.csv",
        AssetRepository, AssetBranch, AssetsBasePath)
    // Download and parse
}

func FetchAsset(path string) ([]byte, error) {
    url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s",
        AssetRepository, AssetBranch, path)
    // Download asset
}
```

## Asset Stability Strategy

To prevent broken assets reaching users:

Keep main branch stable:

- Test assets before merging to main
- Don't commit broken TOML or malformed Markdown
- Review asset changes like code
- Use PRs for asset contributions
- Validate metadata in PR reviews

Validation on catalog index generation:

- `start assets index` validates TOML syntax
- Checks required metadata fields
- Errors on missing/invalid data
- Catches issues before merge

Per-asset validation (future consideration):

- Could validate TOML syntax before caching
- Could check for required fields
- Could rollback individual asset on failure
- Not implemented initially (keep simple)

## Usage Examples

Searching catalog (downloads index.csv from main):

```bash
$ start assets search "commit"

Searching catalog...
✓ Loaded index from main branch (46 assets)

Found 3 matches:
  tasks/git-workflow/commit-message
  tasks/git-workflow/pre-commit-review
  tasks/workflows/post-commit-hook
```

Adding asset (downloads from main):

```bash
$ start assets add tasks/git-workflow/pre-commit-review

Found in catalog:
  Name: pre-commit-review
  Description: Review staged changes before committing
  Branch: main

Downloading from main branch...
✓ Downloaded to cache: ~/.config/start/assets/tasks/git-workflow/
✓ Added to global config

Ready to use: start task pre-commit-review
```

Updating cached assets (checks main branch):

```bash
$ start assets update

Checking for updates on main branch...
✓ Found 2 updated assets
  - tasks/git-workflow/commit-message (SHA changed)
  - roles/code-reviewer (SHA changed)

Downloading updates...
✓ Updated 2 cached assets

Cache refreshed from main branch
```

Lazy loading (downloads from main on first use):

```bash
$ start task pre-commit-review

Task 'pre-commit-review' not found locally.
Found in catalog (main branch): tasks/git-workflow/pre-commit-review
Downloading...
✓ Cached and added to config

Running task 'pre-commit-review'...
```

## Comparison: CLI vs Assets

| Aspect | CLI Binary | Asset Catalog |
|--------|------------|---------------|
| Source | GitHub Releases | Main branch |
| Version | Semantic (v1.2.3) | Branch commit (main) |
| Update Command | `brew upgrade` / `go install` | `start assets update` |
| Update Frequency | Manual releases | Every commit to main |
| Testing | Full release process | PR review before merge |
| Stability | High (release process) | Medium (main branch stability) |
| Risk | High (executable code) | Low (config templates) |
| Fetch Method | Binary download | raw.githubusercontent.com |

## Implementation

Code location: `internal/assets/catalog.go`

```go
const (
    AssetRepository = "grantcarthew/start"
    AssetBranch     = "main"  // Always main
    AssetsBasePath  = "assets"
)

// FetchCatalogIndex downloads index.csv from main branch
func FetchCatalogIndex(ctx context.Context) (*CatalogIndex, error) {
    url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/index.csv",
        AssetRepository, AssetBranch, AssetsBasePath)

    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("download index: %w", err)
    }
    defer resp.Body.Close()

    return parseCatalogIndex(resp.Body)
}
```

Implementation checklist:

- Use main branch constant in all catalog operations
- Fetch index.csv from main branch
- Download individual assets from main branch
- Update cached assets by comparing with main branch
- Document in README that assets come from main branch
- No version tracking file needed (catalog system handles SHA comparison per asset)

## Updates

- 2025-01-10: Updated by DR-031 (catalog architecture) - catalog-based implementation instead of bulk downloads
- 2025-01-13: Updated by DR-039 (catalog index) - index.csv fetched from main branch for catalog metadata
