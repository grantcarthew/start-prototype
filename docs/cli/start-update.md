# start update

## Name

start update - Update cached assets from GitHub

## Synopsis

```bash
start update
start update [flags]
```

## Description

Checks cached assets for updates and downloads newer versions from GitHub. Updates only the asset cache (`~/.config/start/assets/`), never modifies user configuration files.

**Behavior:**

- **This command only checks for updates to assets already present in your local cache. It does not discover or download new assets from the catalog.**
- Compares local asset SHAs with remote SHAs from GitHub catalog
- Downloads new versions of assets if SHAs differ

**What gets updated:**

- **Cached roles** - System prompt markdown files
- **Cached tasks** - Task definition files
- **Cached agents** - Agent configuration templates
- **Asset metadata** - `.meta.toml` sidecar files

**What doesn't change:**

- User's config files (config.toml, tasks.toml, agents.toml, contexts.toml)
- Custom files outside asset cache
- Installed agent binaries
- User customizations

**Use cases:**

- Get latest improvements to catalog assets
- Update role prompts with better content
- Get bug fixes in task definitions
- Refresh cache after repository updates

**Safety:**

- Non-destructive: Only updates asset cache
- User config files untouched
- User modifications preserved
- Can re-run safely anytime

## Flags

This command supports the standard global flags for controlling output verbosity and showing help: `--verbose`, `--quiet`, and `--help`. See `start --help` for more details.

## Behavior

### Normal Update

```bash
start update
```

**Process:**

1. Check network connectivity
2. Fetch GitHub catalog tree (all files and SHAs)
3. Find all cached `.meta.toml` files
4. Compare local SHA with remote SHA for each asset
5. Download assets where SHA differs
6. Update cache files and metadata
7. Report changes

**Output:**

```
Checking for asset updates...

  ⬇ Updating tasks/pre-commit-review...
  ⬇ Updating roles/code-reviewer...

✓ Update complete
  Updated: 2 assets
  Unchanged: 10 assets

Note: Your configuration files are unchanged.
Review updated assets and manually update config if desired.
```

### No Updates Available

```bash
start update
```

When all cached assets are up to date:

```
Checking for asset updates...

✓ Update complete
  Updated: 0 assets
  Unchanged: 12 assets

All cached assets are up to date.
```

### First Run (No Cache)

```bash
start update
```

When `~/.config/start/assets/` is empty:

```
Checking for asset updates...

✓ Update complete
  Updated: 0 assets
  Unchanged: 0 assets

No cached assets found.

To download assets:
  - Browse catalog: start config task add
  - Use a task: start task <name>
```

### Verbose Mode

```bash
start update --verbose
```

Shows detailed SHA comparison:

```
Checking for asset updates...

Fetching catalog from GitHub...
  Repository: grantcarthew/start
  Branch: main
  ✓ Fetched tree (156 files)

Scanning cached assets...
  Found 12 cached assets

Comparing SHAs:
  tasks/git-workflow/pre-commit-review
    Local:  a1b2c3d4...
    Remote: b2c3d4e5...
    → UPDATE NEEDED

  tasks/git-workflow/pr-ready
    Local:  e5f6g7h8...
    Remote: e5f6g7h8...
    → UP TO DATE

  roles/general/code-reviewer
    Local:  i9j0k1l2...
    Remote: j0k1l2m3...
    → UPDATE NEEDED

  [... 9 more assets ...]

Downloading updates:
  ⬇ tasks/git-workflow/pre-commit-review.toml (1.2 KB)
  ⬇ tasks/git-workflow/pre-commit-review.md (3.4 KB)
  ⬇ tasks/git-workflow/pre-commit-review.meta.toml (245 bytes)
  ⬇ roles/general/code-reviewer.md (4.8 KB)
  ⬇ roles/general/code-reviewer.meta.toml (198 bytes)

✓ Update complete
  Updated: 2 assets
  Unchanged: 10 assets
```

### Quiet Mode

```bash
start update --quiet
```

Minimal output:

```
✓ Updated 2 assets, 10 unchanged
```

## SHA-Based Versioning

Assets use Git blob SHA for versioning (per DR-032):

**Metadata file example:**

```toml
# pre-commit-review.meta.toml
name = "pre-commit-review"
description = "Review staged changes before committing"
tags = ["git", "review", "quality"]
sha = "a1b2c3d4e5f6789012345678901234567890abcd"
created = "2025-01-10T00:00:00Z"
updated = "2025-01-10T12:30:00Z"
```

**Update detection:**

1. Read local SHA from cached `.meta.toml`
2. Get remote SHA from GitHub Tree API
3. If SHAs differ → download new version
4. Update cached files with new content and SHA

**Benefits:**

- Reliable: Content hash guarantees integrity
- Simple: No version number conflicts
- Efficient: Single API call for all SHAs

## Exit Codes

**0** - Success (assets checked, updates applied if available)

**1** - Network error
- Cannot reach GitHub
- Repository not found
- API rate limit exceeded

**2** - File system error
- Cannot write to asset cache
- Permission denied
- Disk full

**3** - Partial failure
- Some assets updated, some failed
- Cache may be incomplete

## Error Handling

### Network Errors

**Cannot reach GitHub:**

```
Error: Cannot connect to GitHub

  Network error: dial tcp: no route to host

Check your internet connection and try again.
```

Exit code: 1

**Rate limit exceeded:**

```
Error: GitHub API rate limit exceeded

  Limit: 60 requests/hour (anonymous)
  Reset: 2025-01-10 12:00:00 (in 45 minutes)

Solutions:
1. Set GITHUB_TOKEN for 5,000 requests/hour:
   export GITHUB_TOKEN=ghp_xxxxxxxxxxxx

2. Wait until rate limit resets

3. Use cached assets (if available)
```

Exit code: 1

### File System Errors

**Permission denied:**

```
Error: Cannot write to asset cache

  Path: ~/.config/start/assets/
  Error: permission denied

Check directory permissions:
  chmod 755 ~/.config/start/assets
```

Exit code: 2

**Disk full:**

```
Error: Insufficient disk space

  Required: ~50 KB
  Available: 12 KB

Free up disk space and try again.
```

Exit code: 2

### Partial Failures

**Some assets failed:**

```
Warning: Update partially failed

  ✓ Updated: tasks/pre-commit-review
  ✗ Failed: roles/code-reviewer (network timeout)

Some assets updated successfully.
Re-run 'start update' to retry failed downloads.
```

Exit code: 3

Cache is partially updated. Successful downloads are kept.

### Asset Removed from Catalog

**Asset exists in cache but not in GitHub:**

```
Warning: Asset not found in catalog

  Asset: tasks/deprecated-task
  (exists in cache but removed from GitHub)

Keep cached version? (no action taken)
```

Exit code: 0 (warning only, not an error)

## Manual Config Updates

After updating cache, review changes and manually update config if desired:

**View updated asset:**

```bash
cat ~/.config/start/assets/tasks/git-workflow/pre-commit-review.toml
cat ~/.config/start/assets/tasks/git-workflow/pre-commit-review.meta.toml
```

**Compare with your config:**

```bash
# View your current task config
cat ~/.config/start/tasks.toml | grep -A 10 "\[tasks.pre-commit-review\]"

# View updated cached version
cat ~/.config/start/assets/tasks/git-workflow/pre-commit-review.toml
```

**Update your config if desired:**

```bash
start config task edit  # Manually update with new content
```

## Network Requirements

**Required:**
- HTTPS access to `github.com`
- HTTPS access to `api.github.com`
- HTTPS access to `raw.githubusercontent.com`

**Optional:**
- `GITHUB_TOKEN` environment variable (for higher rate limits)
  - Anonymous: 60 requests/hour
  - Authenticated: 5,000 requests/hour

**Offline behavior:**

Per DR-026, update requires network access:

```
Error: Cannot connect to GitHub

  Network error: no internet connection

Update requires network access.
Asset cache not modified.
```

## Performance

**Typical update:**
- Check 12 cached assets: < 1 second (1 API call)
- Download 2 updated assets: 1-2 seconds
- Network speed dependent

**Bandwidth:**
- Tree API call: ~5-10 KB
- Asset downloads: 2-5 KB per asset (via raw.githubusercontent.com)
- Total: ~10-30 KB for typical update

**Disk space:**
- Per asset: 2-5 KB (toml + md + meta)
- Cache growth: Minimal (old versions overwritten)

## Notes

### Update Frequency

**Manual only** (per DR-025):
- No automatic background checks
- No version checking on CLI startup
- User explicitly runs `start update`

**Recommended:**
- Update monthly for improvements
- Update when you need new features
- Update is optional, not required

### User Config Never Auto-Updated

**Cache updates are separate from config:**

When you run `start update`:
1. Cache gets new versions
2. Your config files remain unchanged
3. If task references cache file path, it uses new version automatically
4. If task has inlined content, you must manually update

**Example - Automatic update (file reference):**

```toml
# In your tasks.toml
[tasks.pre-commit-review]
prompt_file = "~/.config/start/assets/tasks/git-workflow/pre-commit-review.md"

# After 'start update':
# - Cache file updated with new content
# - Your config unchanged
# - Next task run uses new cached file
# ✓ Automatic
```

**Example - Manual update needed (inlined):**

```toml
# In your tasks.toml
[tasks.pre-commit-review]
prompt = """
Review the following changes...
(old inlined content)
"""

# After 'start update':
# - Cache file updated with new content
# - Your config unchanged
# - Still uses old inlined content
# ✗ Manual update needed
```

### Rollback

**No built-in rollback in v1:**

If you want previous version:
1. Delete cache: `rm -rf ~/.config/start/assets`
2. Re-download specific version (future feature)

Cache is disposable - can always re-download.

### Privacy

`start update` makes these network requests:
- GitHub Tree API: Repository tree with SHAs
- raw.githubusercontent.com: Asset file content

No telemetry, no tracking, no data sent to external services.

## Examples

### Regular Update

```bash
start update
```

Check for and apply asset updates.

### Quiet Update (CI/CD)

```bash
start update --quiet
if [ $? -ne 0 ]; then
  echo "Asset update failed"
  exit 1
fi
```

Minimal output, check exit code.

### Verbose Troubleshooting

```bash
start update --verbose
```

See SHA comparison and download details.

### Check Before Updating

```bash
# Use doctor to check asset status (future feature)
start doctor

# Update assets
start update
```

## Configuration

**Settings in config.toml:**

```toml
[settings]
asset_repo = "grantcarthew/start"   # GitHub repository
asset_download = true               # Download if not found (doesn't affect update)
```

**GitHub Authentication:** Uses `GITHUB_TOKEN` environment variable (optional, recommended to avoid rate limits).

## See Also

- [DR-031](../design/decisions/dr-031-catalog-based-assets.md) - Catalog architecture
- [DR-037](../design/decisions/dr-037-asset-updates.md) - Update mechanism
- [DR-032](../design/decisions/dr-032-asset-metadata-schema.md) - Metadata schema
- start-config-task(1) - Add tasks from catalog
- start-config-role(1) - Add roles from catalog
- start-init(1) - Initialize configuration
- start-doctor(1) - Diagnose installation
