# DR-039: Catalog Index File

**Date:** 2025-01-13
**Status:** Accepted
**Category:** Asset Management

## Decision

Use a pre-built CSV index file (`assets/index.csv`) to enable fast metadata-rich searching of the GitHub catalog without downloading individual `.meta.toml` files. Generated via `start assets index` command, sorted alphabetically for predictable updates.

## Problem

Searching the GitHub catalog by description and tags requires accessing metadata that's distributed across hundreds of `.meta.toml` files. Options considered:

1. **Download all .meta.toml files** - Works but slow (200+ HTTP requests for large catalog)
2. **GitHub Search API** - Rate limited (10-30 requests/min), indexing delays, requires auth for reasonable limits
3. **Pre-built index** - Single fast download, proven pattern (npm, cargo, homebrew all do this)

Without an index, searching by description/tags is impractical at scale.

## Solution

### Index File Structure

**Location:** `assets/index.csv` (in catalog repository)

**Format:** CSV with header row

**Columns:**
```csv
type,category,name,description,tags,sha,size,created,updated
```

**Field definitions:**

- **type** - Asset type (`tasks`, `roles`, `agents`, `contexts`)
- **category** - Category subdirectory (e.g., `git-workflow`, `general`)
- **name** - Asset name (filename without extension)
- **description** - Human-readable summary (from .meta.toml)
- **tags** - Semicolon-separated keywords (e.g., `git;review;quality`)
- **sha** - Git blob SHA of content file (for update detection)
- **size** - File size in bytes
- **created** - ISO 8601 timestamp when asset was created
- **updated** - ISO 8601 timestamp when asset was last modified

**Sorting:** Alphabetical by `type` → `category` → `name`

**Example:**
```csv
type,category,name,description,tags,sha,size,created,updated
agents,anthropic,claude,Anthropic Claude AI via Claude Code CLI,claude;anthropic;ai,a1b2c3d4,1024,2025-01-10T00:00:00Z,2025-01-10T00:00:00Z
roles,general,code-reviewer,Expert code reviewer focusing on security,review;security;quality,b2c3d4e5,2048,2025-01-10T00:00:00Z,2025-01-12T00:00:00Z
tasks,git-workflow,commit-message,Generate conventional commit message,git;commit;conventional,c3d4e5f6,1536,2025-01-10T00:00:00Z,2025-01-10T00:00:00Z
tasks,git-workflow,pre-commit-review,Review staged changes before committing,git;review;quality;pre-commit,d4e5f6a1,2048,2025-01-10T00:00:00Z,2025-01-11T00:00:00Z
```

### Generation Process

**Command:** `start assets index`

**Validation:**
1. Check for `.git/` directory - Must be a git repository
2. Check for `assets/` directory - Must be catalog repository structure
3. Error if either missing: "Must be run in catalog repository"

**Algorithm:**
```
1. Scan assets/ directory recursively
2. Find all *.meta.toml files
3. For each .meta.toml:
   - Parse TOML
   - Extract: type, category (from path), name, description, tags, sha, created, updated
   - Validate all required fields present
4. Sort entries alphabetically by type → category → name
5. Write CSV to assets/index.csv with header row
6. Report: "Generated index with N assets"
```

**Error handling:**
- Missing .meta.toml for an asset → Warning, skip asset
- Invalid TOML syntax → Error, show filename and line number
- Missing required fields → Error, show which fields and which file
- File system errors → Error and exit

**Example usage:**
```bash
# Clone catalog repository
git clone https://github.com/grantcarthew/start.git
cd start

# Add new asset with metadata
mkdir -p assets/tasks/my-category
cat > assets/tasks/my-category/my-task.meta.toml <<EOF
[metadata]
name = "my-task"
description = "My new task"
tags = ["workflow", "custom"]
sha = "..."
created = "2025-01-13T00:00:00Z"
updated = "2025-01-13T00:00:00Z"
EOF

# Generate index
start assets index

# Output:
# Scanning assets/...
# Found 46 assets
# ✓ Generated index with 46 assets
# Updated: assets/index.csv

# Commit
git add assets/
git commit -m "Add my-task and regenerate index"
```

### Usage Pattern

**When to fetch:**
- Download `assets/index.csv` on every CLI execution that needs catalog search
- No local caching (file is small, ~10-50KB)
- Fresh data every time

**How to fetch:**
```
GET https://raw.githubusercontent.com/{org}/{repo}/{branch}/assets/index.csv
```

Uses raw.githubusercontent.com (no API rate limits).

**Parsing:**
```go
func loadCatalogIndex() (*CatalogIndex, error) {
    // Download index.csv
    resp, err := http.Get(indexURL)
    if err != nil {
        return nil, fmt.Errorf("download index: %w", err)
    }
    defer resp.Body.Close()

    // Parse CSV
    reader := csv.NewReader(resp.Body)
    records, err := reader.ReadAll()
    if err != nil {
        return nil, fmt.Errorf("parse CSV: %w", err)
    }

    // Skip header row, parse entries
    var assets []Asset
    for i, record := range records[1:] {
        asset := Asset{
            Type:        record[0],
            Category:    record[1],
            Name:        record[2],
            Description: record[3],
            Tags:        strings.Split(record[4], ";"),
            SHA:         record[5],
            Size:        parseInt(record[6]),
            Created:     parseTime(record[7]),
            Updated:     parseTime(record[8]),
        }
        assets = append(assets, asset)
    }

    return &CatalogIndex{Assets: assets}, nil
}
```

**In-memory search:**
```go
func (idx *CatalogIndex) Search(query string) []Asset {
    query = strings.ToLower(query)
    var results []Asset

    for _, asset := range idx.Assets {
        // Substring match: name, description, or tags
        if strings.Contains(strings.ToLower(asset.Name), query) ||
           strings.Contains(strings.ToLower(asset.Description), query) ||
           containsTag(asset.Tags, query) {
            results = append(results, asset)
        }
    }

    return results
}
```

### Fallback Behavior

**When index.csv is missing/corrupted:**

Progressive degradation strategy:

```
1. Attempt to download index.csv
   - Success → use index for rich search
   - Fail → log warning, fall back to Tree API

2. Fall back to Tree API directory listing
   - GET /repos/{org}/{repo}/git/trees/{branch}?recursive=1
   - Filter paths: assets/{type}/*/*.toml (exclude .meta.toml)
   - Extract names from paths only
   - Limited search: name/path matching only (no descriptions/tags)

3. Return results with degraded metadata
   - Name and path available
   - Description: "(metadata unavailable)"
   - Tags: empty
```

**User experience:**
```bash
$ start assets search "commit"

Searching catalog...
⚠ Index unavailable, using directory listing (limited metadata)

Found 2 matches:
  tasks/git-workflow/commit-message
  tasks/git-workflow/pre-commit-review
```

**When asset not in index:**

Even if index.csv exists, newly added assets may not be indexed yet:

```
1. Search index.csv → not found
2. Search Tree API directory listing → found!
3. Return result with note:
   Found: tasks/workflows/new-feature
   (Not in index - metadata unavailable)
```

This handles PRs that add assets without regenerating the index.

### Source of Truth

**.meta.toml files are canonical:**
- Each asset has its own `.meta.toml` with complete metadata
- index.csv is a **derived file** for performance
- If index.csv and .meta.toml conflict → .meta.toml wins

**Update workflow:**
1. Contributor edits .meta.toml files (source of truth)
2. Run `start assets index` to regenerate index.csv (derived)
3. Commit both .meta.toml changes and updated index.csv

### Multi-File Assets

**Index entry per asset, not per file:**

Assets may have multiple files:
```
assets/tasks/git-workflow/pre-commit-review.toml
assets/tasks/git-workflow/pre-commit-review.md
assets/tasks/git-workflow/pre-commit-review.meta.toml
```

Index contains **one entry**:
```csv
tasks,git-workflow,pre-commit-review,Review staged changes,git;review,abc123,2048,2025-01-10T00:00:00Z,2025-01-10T00:00:00Z
```

When downloading the asset, fetch all files with matching name prefix.

## Benefits

**Performance:**
- ✅ Single ~10-50KB download vs 200+ individual .meta.toml requests
- ✅ Fast in-memory search with rich metadata
- ✅ No rate limits (raw.githubusercontent.com)

**Reliability:**
- ✅ Graceful degradation if index missing
- ✅ Still works for brand-new assets not yet indexed
- ✅ Tree API fallback ensures catalog always accessible

**Maintainability:**
- ✅ Standard pattern used by package managers
- ✅ Simple CSV format, easy to inspect/debug
- ✅ Alphabetical sorting = predictable diffs
- ✅ Built-in generation tool (`start assets index`)

**Contributor friendly:**
- ✅ Clear workflow: add .meta.toml → run `start assets index` → commit
- ✅ Validation during generation catches errors early
- ✅ Even if index not regenerated, assets still discoverable

## Trade-offs Accepted

**Manual generation:**
- ❌ Requires running `start assets index` after changes
- ❌ PRs must include both .meta.toml and index.csv updates
- **Mitigation:** Clear documentation, validation in generation tool, fallback to Tree API

**CSV limitations:**
- ❌ Tags as semicolon-separated strings (not native array)
- ❌ Escaping needed for descriptions with commas/quotes
- **Mitigation:** Standard CSV escaping rules, proven format

**No local cache:**
- ❌ Downloads index.csv on every use
- ❌ ~10-50KB download each time
- **Mitigation:** File is small, download is fast (<100ms typical), always fresh data

**Potential staleness:**
- ❌ If contributor forgets to run `start assets index`, index is stale
- **Mitigation:** Fallback to Tree API still finds new assets, just without rich metadata

## Implementation

### CSV Escaping

Follow RFC 4180 CSV standard:

**Fields containing commas:**
```csv
tasks,workflow,my-task,"Review code, check tests, verify docs",testing,abc123,1024,2025-01-13T00:00:00Z,2025-01-13T00:00:00Z
```

**Fields containing quotes:**
```csv
tasks,workflow,my-task,"Review ""special"" cases",testing,abc123,1024,2025-01-13T00:00:00Z,2025-01-13T00:00:00Z
```

Standard Go `encoding/csv` package handles this automatically.

### Generation Tool

**Context validation:**
```go
func validateCatalogRepo() error {
    // Check .git directory
    if _, err := os.Stat(".git"); os.IsNotExist(err) {
        return fmt.Errorf("not a git repository")
    }

    // Check assets directory
    if _, err := os.Stat("assets"); os.IsNotExist(err) {
        return fmt.Errorf("assets/ directory not found")
    }

    return nil
}
```

**Directory scanning:**
```go
func scanAssets() ([]AssetMeta, error) {
    var assets []AssetMeta

    err := filepath.Walk("assets", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // Only process .meta.toml files
        if !strings.HasSuffix(path, ".meta.toml") {
            return nil
        }

        // Parse metadata
        meta, err := parseMetadata(path)
        if err != nil {
            return fmt.Errorf("%s: %w", path, err)
        }

        // Extract type and category from path
        // e.g., assets/tasks/git-workflow/pre-commit.meta.toml
        parts := strings.Split(filepath.Dir(path), string(filepath.Separator))
        if len(parts) >= 3 {
            meta.Type = parts[1]      // "tasks"
            meta.Category = parts[2]  // "git-workflow"
        }

        assets = append(assets, meta)
        return nil
    })

    return assets, err
}
```

**Sorting:**
```go
func sortAssets(assets []AssetMeta) {
    sort.Slice(assets, func(i, j int) bool {
        // Primary: type
        if assets[i].Type != assets[j].Type {
            return assets[i].Type < assets[j].Type
        }
        // Secondary: category
        if assets[i].Category != assets[j].Category {
            return assets[i].Category < assets[j].Category
        }
        // Tertiary: name
        return assets[i].Name < assets[j].Name
    })
}
```

**Writing CSV:**
```go
func writeIndex(assets []AssetMeta) error {
    file, err := os.Create("assets/index.csv")
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    // Header row
    writer.Write([]string{
        "type", "category", "name", "description",
        "tags", "sha", "size", "created", "updated",
    })

    // Data rows
    for _, asset := range assets {
        writer.Write([]string{
            asset.Type,
            asset.Category,
            asset.Name,
            asset.Description,
            strings.Join(asset.Tags, ";"),
            asset.SHA,
            strconv.Itoa(asset.Size),
            asset.Created.Format(time.RFC3339),
            asset.Updated.Format(time.RFC3339),
        })
    }

    return writer.Error()
}
```

## Examples

### Searching with Index

```bash
$ start assets search "commit"

Searching catalog...
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

### Searching without Index (Fallback)

```bash
$ start assets search "commit"

Searching catalog...
⚠ Index unavailable, using directory listing (limited metadata)

Found 3 matches:

tasks/git-workflow/commit-message
  (Metadata unavailable - not in index)

tasks/git-workflow/pre-commit-review
  (Metadata unavailable - not in index)

tasks/workflows/post-commit-hook
  (Metadata unavailable - not in index)
```

### New Asset Not in Index

```bash
$ start assets search "brand-new"

Searching catalog...
✓ Loaded index (46 assets)
✗ Not found in index

Checking directory tree...
✓ Found 1 match:

tasks/experimental/brand-new-feature
  (Not in index - metadata unavailable)

Tip: Ask maintainer to run 'start assets index'
```

### Generating Index

```bash
$ cd ~/projects/start-catalog-fork

$ start assets index

Validating repository structure...
✓ Git repository detected
✓ Assets directory found

Scanning assets/...
  Found: agents/anthropic/claude.meta.toml
  Found: roles/general/code-reviewer.meta.toml
  Found: tasks/git-workflow/commit-message.meta.toml
  Found: tasks/git-workflow/pre-commit-review.meta.toml
  ... (42 more)

Sorting assets (type → category → name)...
Writing index to assets/index.csv...

✓ Generated index with 46 assets
Updated: assets/index.csv

Ready to commit:
  git add assets/
  git commit -m "Regenerate catalog index"
```

### Generation with Errors

```bash
$ start assets index

Validating repository structure...
✗ Error: Not a git repository

This command must be run in the catalog repository.
Run: git clone https://github.com/grantcarthew/start.git
```

```bash
$ start assets index

Scanning assets/...
✗ Error: tasks/broken/bad-task.meta.toml:5: invalid TOML syntax

Fix the metadata file and try again.
```

## Size Estimates

**Minimal catalog (28 assets):**
- 28 rows × ~120 bytes/row = ~3.4 KB
- With header and CSV overhead = ~4 KB

**Large catalog (200 assets):**
- 200 rows × ~120 bytes/row = ~24 KB
- With header and CSV overhead = ~30 KB

**Very large catalog (500 assets):**
- 500 rows × ~120 bytes/row = ~60 KB
- With header and CSV overhead = ~70 KB

All well within acceptable download sizes (< 100ms on typical connections).

## Related Decisions

- [DR-032](./dr-032-asset-metadata-schema.md) - Asset metadata schema (defines .meta.toml structure)
- [DR-034](./dr-034-github-catalog-api.md) - GitHub API strategy (Tree API fallback)
- [DR-035](./dr-035-interactive-browsing.md) - Interactive browsing (uses index for search)
- [DR-036](./dr-036-cache-management.md) - Cache management (index.csv not cached locally)

## Future Considerations

**GitHub Actions automation:**
- Could add workflow to auto-generate index on PR
- Validates contributors didn't forget to run `start assets index`
- Current: Manual generation is acceptable

**Compression:**
- Could gzip index.csv for even smaller downloads
- Trade-off: Parsing complexity vs ~50% size reduction
- Current: Uncompressed is simple and fast enough

**Local caching:**
- Could cache index.csv for 5-10 minutes
- Trade-off: Stale data vs fewer downloads
- Current: Always fresh is preferred, file is small

**Metadata in Tree API:**
- GitHub doesn't expose file contents in Tree API
- Would require custom GitHub App or different approach
- Current: index.csv is simpler and proven

**Current stance:** Ship with described behavior. Monitor performance and iterate based on user feedback.
