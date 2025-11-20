# start assets update

## Name

start assets update - Update cached assets from GitHub catalog

## Synopsis

```bash
start assets update            # Update all cached assets
start assets update <query>    # Update matching assets
```

## Description

Check for updates to cached assets and download new versions from the GitHub catalog. Compares local asset SHAs with catalog SHAs to detect changes. Updates only the asset cache, never modifies user configuration files.

**Update process:**

1. Fetch catalog tree from GitHub (contains SHAs for all files)
2. Find all cached assets (`.meta.toml` files in cache)
3. Compare local SHA with catalog SHA for each asset
4. Download and update assets with SHA mismatches
5. Report update summary

**Note:** While configuration files (`tasks.toml`) are not modified, assets are typically installed as file references (e.g., `prompt_file = "~/.config/start/assets/..."`). Updating the cached files **automatically updates the behavior** of these tasks without requiring config changes.

## Arguments

**\<query\>** (optional)
: Search query to filter which assets to update. Uses substring matching.

**No query** - Update all cached assets

**With query** - Update only assets matching the query (name, path, category)

**Query matching:**

- Case-insensitive substring match
- Searches: asset name, path, category
- Same matching algorithm as `start assets search`

## Behavior

### Update All Assets (No Query)

```
1. Fetch GitHub catalog tree
2. Find all .meta.toml files in cache
3. For each cached asset:
   - Compare local SHA with catalog SHA
   - If different → download new version
   - Update cache (files + metadata)
4. Report results
```

### Selective Update (With Query)

```
1. Fetch GitHub catalog tree
2. Find cached assets matching query
3. For each matching asset:
   - Compare local SHA with catalog SHA
   - If different → download new version
   - Update cache
4. Report results
```

### What Gets Updated

**Updated:**

- Asset files in `~/.config/start/assets/`
- `.meta.toml` files (SHA, updated timestamp)

**Not Updated:**

- `~/.config/start/*.toml` (configuration files)
- `./.start/*.toml` (local configuration)
- User-created custom assets

### SHA-Based Version Detection

**Git blob SHA as version:**

- No semantic versioning (no v1.0.0, etc.)
- Git blob SHA uniquely identifies file content
- SHA mismatch = content changed = update available

**Metadata tracking:**

```toml
# .meta.toml file
sha = "a1b2c3d4e5f6..."  # Git blob SHA
updated = "2025-01-12T14:30:00Z"
```

## Output

### Updates Available

```bash
$ start assets update

Checking for asset updates...

✓ Fetched catalog (46 assets)

Comparing cached assets with catalog...

  ⬇ Updating tasks/git-workflow/pre-commit-review...
     SHA: a1b2c3d4... → b2c3d4e5...
     Updated: 2025-01-10 → 2025-01-12

  ⬇ Updating roles/general/code-reviewer...
     SHA: c3d4e5f6... → d4e5f6a7...
     Updated: 2025-01-05 → 2025-01-12

✓ Update complete
  Updated: 2 assets
  Unchanged: 10 assets

Note: Your configuration files are unchanged.
The updated assets are now available in the cache.
```

### No Updates

```bash
$ start assets update

Checking for asset updates...

✓ Fetched catalog (46 assets)

Comparing cached assets with catalog...

✓ Update complete
  Updated: 0 assets
  Unchanged: 12 assets

All cached assets are up to date.
```

### Selective Update (With Query)

```bash
$ start assets update "commit"

Checking for updates to assets matching 'commit'...

✓ Fetched catalog (46 assets)

Found 3 cached assets matching 'commit':
  - tasks/git-workflow/commit-message
  - tasks/git-workflow/pre-commit-review
  - tasks/git-workflow/post-commit-hook

Comparing with catalog...

  ⬇ Updating tasks/git-workflow/pre-commit-review...
     SHA: a1b2c3d4... → b2c3d4e5...

✓ Update complete
  Updated: 1 asset
  Unchanged: 2 assets
```

### Category Update

```bash
$ start assets update "git-workflow"

Checking for updates to assets matching 'git-workflow'...

✓ Fetched catalog (46 assets)

Found 6 cached assets in git-workflow category

Comparing with catalog...

  ⬇ Updating tasks/git-workflow/pre-commit-review...
  ⬇ Updating tasks/git-workflow/commit-message...

✓ Update complete
  Updated: 2 assets
  Unchanged: 4 assets
```

### No Cached Assets

```bash
$ start assets update

Checking for asset updates...

No cached assets found.

Use 'start assets add <query>' to install assets.
```

Exit code: 0

### Network Error

```bash
$ start assets update

Checking for asset updates...
✗ Network error

Cannot connect to GitHub:
  dial tcp: no route to host

Check your internet connection and try again.
```

Exit code: 1

### Asset Not in Catalog

```bash
$ start assets update

Checking for asset updates...

✓ Fetched catalog (46 assets)

Comparing cached assets with catalog...

  ⚠ tasks/custom/my-task not found in catalog (skipped)

✓ Update complete
  Updated: 0 assets
  Unchanged: 11 assets
  Skipped: 1 asset (not in catalog)
```

Custom/local assets not in catalog are skipped.

## Exit Codes

**0** - Success (updates checked and applied)

**1** - Network error (catalog unavailable)

**2** - File system error (cache write failed)

**3** - Partial failure (some updates failed)

## Flags

**--dry-run**, **-n**
: Show what would be updated without downloading.

**--verbose**, **-v**
: Show detailed SHA comparison and file changes.

**--force**, **-f**
: Force re-download even if SHAs match.

## Examples

### Update All Cached Assets

```bash
$ start assets update

Checking for asset updates...

✓ Fetched catalog (46 assets)

Comparing 12 cached assets...

  ⬇ Updating tasks/git-workflow/pre-commit-review...
  ⬇ Updating roles/general/code-reviewer...

✓ Update complete
  Updated: 2 assets
  Unchanged: 10 assets
```

### Update Specific Asset

```bash
$ start assets update "pre-commit-review"

Checking for updates to assets matching 'pre-commit-review'...

Found 1 cached asset:
  - tasks/git-workflow/pre-commit-review

Comparing with catalog...

  ⬇ Updating tasks/git-workflow/pre-commit-review...
     SHA: a1b2c3d4... → b2c3d4e5...
     Files updated: pre-commit-review.toml (2.1 KB → 2.3 KB)

✓ Update complete
  Updated: 1 asset
```

### Update by Category

```bash
$ start assets update "git-workflow"

Checking for updates to assets matching 'git-workflow'...

Found 6 cached assets in git-workflow category

Comparing with catalog...

  ⬇ Updating tasks/git-workflow/pre-commit-review...
  ⬇ Updating tasks/git-workflow/commit-message...

✓ Update complete
  Updated: 2 assets
  Unchanged: 4 assets
```

### Dry Run Mode

```bash
$ start assets update --dry-run

Checking for asset updates (DRY RUN)...

✓ Fetched catalog (46 assets)

Would update:
  ⬇ tasks/git-workflow/pre-commit-review
     Current: a1b2c3d4... (2025-01-10)
     Latest:  b2c3d4e5... (2025-01-12)

  ⬇ roles/general/code-reviewer
     Current: c3d4e5f6... (2025-01-05)
     Latest:  d4e5f6a7... (2025-01-12)

Would update: 2 assets
Would skip: 10 assets (unchanged)

No changes made (dry run).
Run without --dry-run to apply updates.
```

### Verbose Output

```bash
$ start assets update "pre-commit" --verbose

Checking for updates to assets matching 'pre-commit'...

Fetching catalog from GitHub...
  URL: https://api.github.com/repos/grantcarthew/start/git/trees/main?recursive=1
  ✓ Fetched 846 files

Loading cached metadata...
  Found: ~/.config/start/assets/tasks/git-workflow/pre-commit-review.meta.toml
  Local SHA: a1b2c3d4e5f6789012345678901234567890abcd
  Updated: 2025-01-10T12:00:00Z

Finding in catalog...
  Path: assets/tasks/git-workflow/pre-commit-review.toml
  Remote SHA: b2c3d4e5f6789012345678901234567890abcdef
  Updated: 2025-01-12T14:30:00Z

SHA mismatch detected:
  Local:  a1b2c3d4...
  Remote: b2c3d4e5...
  Status: UPDATE AVAILABLE

Downloading new version...
  ⬇ pre-commit-review.toml (2.1 KB → 2.3 KB)
  ⬇ pre-commit-review.md (1.3 KB → 1.4 KB)
  ⬇ pre-commit-review.meta.toml

Writing to cache...
  ✓ ~/.config/start/assets/tasks/git-workflow/pre-commit-review.toml
  ✓ ~/.config/start/assets/tasks/git-workflow/pre-commit-review.md
  ✓ ~/.config/start/assets/tasks/git-workflow/pre-commit-review.meta.toml

✓ Update complete
  Updated: 1 asset
```

### Force Update

```bash
$ start assets update "code-reviewer" --force

Forcing re-download of assets matching 'code-reviewer'...

  ⬇ Downloading roles/general/code-reviewer... (even though SHA matches)

✓ Update complete
  Updated: 1 asset (forced)
```

### No Matching Query

```bash
$ start assets update "nonexistent"

Checking for updates to assets matching 'nonexistent'...

No cached assets found matching 'nonexistent'

Try:
  - start assets search "nonexistent" (search catalog)
  - start assets update (update all)
```

Exit code: 0 (no error, just nothing to update)

## Use Cases

### Regular Maintenance

**Problem:** Want to keep cached assets up to date.

```bash
# Periodically check for updates
start assets update
```

**Frequency:** Weekly or monthly, as needed.

### Before Important Work

**Problem:** Want latest asset versions before starting a project.

```bash
# Update all assets
start assets update

# Start work with latest versions
start task pre-commit-review
```

### Selective Updates

**Problem:** Want to update only specific assets.

```bash
# Update only git-related tasks
start assets update "git"

# Update specific asset
start assets update "pre-commit-review"
```

### Preview Changes

**Problem:** Want to see what would be updated before applying.

```bash
# Dry run to preview
start assets update --dry-run

# Apply if acceptable
start assets update
```

### Troubleshooting

**Problem:** Asset behaving incorrectly, want to re-download.

```bash
# Force re-download
start assets update "problematic-asset" --force
```

## Comparison with Other Commands

### vs `start assets add`

**`start assets update`** - Updates existing cached assets

```bash
start assets update
# Only updates assets already in cache
```

**`start assets add`** - Adds new assets from catalog

```bash
start assets add "new-asset"
# Downloads and installs new asset
```

Update maintains existing, add acquires new.

## Configuration

**Asset repository:**

In `~/.config/start/config.toml`:

```toml
[settings]
asset_repo = "grantcarthew/start"    # Default
# asset_repo = "myorg/custom-assets"  # Custom
```

**No automatic update checks:**

Updates are **manual only**:

- No automatic checks on CLI startup
- No background update checking
- User must explicitly run `start assets update`

## Notes

### Configuration Files Never Modified

**Critical behavior:** Config files (`*.toml`) are **never automatically changed**.

**What updates:**

- Asset files in cache: `~/.config/start/assets/**/*`
- Metadata files: `*.meta.toml`

**What doesn't update:**

- `~/.config/start/tasks.toml`
- `~/.config/start/roles.toml`
- `./.start/*.toml`

**Rationale:** User owns configuration, cache is transient.

### Referenced vs Inlined Content

**File-referenced assets (automatic update):**

```toml
[tasks.my-task]
file = "~/.config/start/assets/tasks/my-category/my-task.toml"
# After update, next run uses new content automatically
```

**Inlined content (manual update required):**

```toml
[tasks.my-task]
prompt = """
Content copied from asset...
"""
# Update doesn't change this - you must manually update
```

**Best practice:** Use `prompt_file` for auto-updates.

### SHA-Based Versioning

**No semantic versions:**

- No v1.0.0, v2.0.0
- Git blob SHA is the version
- SHA change = content changed

**Why SHAs:**

- Uniquely identifies content
- No version number maintenance
- Git-native approach

### Custom Assets Skipped

**User-created assets** (not from catalog):

- Stored in cache but not in catalog
- Skipped during update (no catalog entry)
- Warning displayed

**Example:**

```
⚠ tasks/custom/my-task not found in catalog (skipped)
```

### Network Required

Requires network access to:

- Fetch catalog tree from GitHub
- Download updated asset files

**Offline:** Cannot check for or apply updates.

### Update Frequency

**No automatic checks** - User decides when to update.

**Recommended frequency:**

- Weekly for active users
- Monthly for occasional users
- Before important work
- When asset behaving unexpectedly

### Performance

**Typical update check:**

- Fetch Tree API: ~100-200ms
- Compare SHAs: <10ms per asset
- Download (if needed): ~50-100ms per file

**Example (12 cached assets, 2 updates):**

- Fetch tree: 150ms
- Compare 12 SHAs: 10ms
- Download 2 assets: 200ms
- **Total: ~360ms**

### Viewing Changes

**After update, to see what changed:**

```bash
# View updated file
cat ~/.config/start/assets/tasks/git-workflow/pre-commit-review.toml

# Check metadata
cat ~/.config/start/assets/tasks/git-workflow/pre-commit-review.meta.toml

# See update timestamp
grep updated ~/.config/start/assets/tasks/**/*.meta.toml
```

**Future enhancement:** `start assets diff <asset>` to show changes.

### Substring Matching for Selective Updates

Query uses substring matching:

- Minimum 3 characters
- Case-insensitive
- Matches name, path, category

**Examples:**

- `"commit"` - Matches all commit-related assets
- `"git-workflow"` - Matches all in category
- `"pre-commit-review"` - Matches specific asset

## See Also

- start-assets(1) - Asset management overview
- start-assets-add(1) - Install new assets
- start-assets-info(1) - View asset details
