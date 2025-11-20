# DR-037: Asset Update Mechanism

- Date: 2025-01-10
- Status: Accepted
- Category: Asset Management

## Problem

Cached assets need an update mechanism to get new versions from GitHub. The update strategy must address:

- Update detection (how to know if new version available)
- Update scope (what gets updated - cache, config, or both)
- User control (automatic vs manual updates)
- Versioning approach (semantic versioning vs content hashing)
- Selective updates (all assets vs specific assets)
- Config preservation (never overwrite user customizations)
- Error handling (partial failures, removed assets, network issues)
- Performance (efficient update checking without excessive API calls)
- Notification strategy (when/how to inform user of updates)

## Decision

Use SHA-based update detection via catalog index to check for asset updates. Update only the cache, never user config. Provide manual-only update command with optional selective updates.

Key aspects:

Command: start assets update [query]

- Manual-only (user explicitly runs command)
- No automatic checks (consistent with no automatic operations principle)
- Optional query parameter for selective updates
- Updates cache only (never modifies user config)

Update process:

1. Download index.csv from GitHub (contains current SHAs)
2. For each cached .meta.toml file:
   - Read local SHA
   - Find asset in index by path
   - Compare SHAs
   - Download if different
3. Update cached files (content + metadata)
4. Report what was updated
5. Never touch user config files

SHA-based versioning:

- No semantic versioning (v1.2.3)
- Git blob SHA IS the version
- SHA comparison for update detection
- Content hash guarantees integrity
- Simple and reliable

Cache-only updates:

- Cache gets new versions
- User config remains unchanged
- Tasks referencing cached files (via prompt_file) automatically use new content
- Tasks with inlined content require manual update
- Clear separation (cache vs config)

Selective updates:

- No query: Update all cached assets
- With query: Update matching assets only (substring matching)
- Examples: start assets update "commit", start assets update git-workflow

## Why

SHA-based versioning is reliable:

- Content hash IS the version (can't be wrong)
- No version number maintenance (automatic from git)
- SHA comparison bulletproof (same SHA = identical content)
- Guarantees integrity (can verify downloaded content)
- Simple implementation (just string comparison)

Index.csv enables efficient updates:

- Single download via raw.githubusercontent.com (no API rate limit)
- Contains SHAs for all assets
- Fast comparison (load index once, compare all cached assets)
- Zero API calls (raw URL download)
- Always fresh data (downloaded on each update check)

Cache-only updates preserve user customizations:

- User config never automatically modified
- Customizations preserved (user changes not lost)
- Explicit is better (user controls when config changes)
- Clear separation (cache is implementation detail, config is user's)
- Safe to update cache (can always delete and re-download)

Manual-only updates provide control:

- User explicitly opts in (no surprise updates)
- Consistent with no automatic operations principle
- User controls timing (update when convenient)
- No background checks (no network usage without consent)
- No startup delays (no version checking on CLI start)

Selective updates add flexibility:

- Update all: start assets update (default behavior)
- Update by query: start assets update "commit" (fewer downloads)
- Update specific asset: start assets update pre-commit-review
- Substring matching (consistent with search/browse)
- Fast for targeted updates (only download what's needed)

## Trade-offs

Accept:

- No automatic notifications (user must remember to run update, can set calendar reminder if desired)
- No rollback in v1 (can't revert to previous version, but cache is disposable so delete and re-download)
- Manual config updates (user must manually apply cache changes to config, but preserves customizations)
- No diff view in v1 (can't see what changed before updating, use file tools to compare if needed)
- No update history (can't see changelog, add if users request)

Gain:

- Reliable versioning (SHA comparison bulletproof, content hash guarantees integrity)
- User control (manual updates only, explicit opt-in, user controls timing)
- Simple implementation (single index download, SHA comparison straightforward, zero API calls)
- Cache-only updates (safe to update cache, user modifications preserved, clear separation)
- Efficient performance (index.csv via raw URL, zero rate limit impact, fast downloads)
- Selective updates (update all or specific assets, flexible workflow, saves bandwidth)

## Alternatives

Semantic versioning with version field:

Example: Use semver in .meta.toml instead of SHA

```toml
[metadata]
version = "1.2.3"
```

- Compare version numbers instead of SHAs
- Check if remote version > local version
- Download if newer version available

Pros:

- Human-readable versions (v1.2.3 easier to understand)
- Semantic meaning (major/minor/patch conveys change type)
- Familiar pattern (developers understand semver)
- Can detect breaking changes (major version bump)

Cons:

- Manual maintenance required (must remember to bump version)
- Can be wrong (version bumped but content unchanged, or vice versa)
- No automatic detection (can't tell if version should change)
- Requires discipline (easy to forget to update)
- Version conflicts possible (merge conflicts on version field)

Rejected: SHA is more reliable - automatic, always accurate, no manual maintenance. SHA comparison IS the version check.

Use Tree API for update checking:

Example: Fetch GitHub Tree API instead of index.csv

- Call Tree API to get all file SHAs
- Compare with cached SHAs
- Download updates

Pros:

- No index file needed (one less thing to maintain)
- Always absolutely fresh (direct from GitHub API)
- Source of truth directly accessed

Cons:

- Counts against rate limit (60/hour anonymous, 5,000/hour authenticated)
- Slower than raw URL (API overhead)
- Can hit rate limit with frequent checks
- Doesn't scale as well (API limits)

Rejected: Index.csv via raw URL is more efficient - zero API calls, no rate limits, fast downloads. Tree API only needed as fallback.

Automatic updates to cache:

Example: Auto-update cache on CLI startup or task execution

- Check for updates automatically
- Download silently in background
- Keep cache fresh without user action

Pros:

- Always up to date (cache automatically refreshed)
- No user action needed (convenient)
- Fresh content always available

Cons:

- Violates no automatic operations principle
- Surprise network usage (privacy/security concern)
- Startup delays (checking for updates adds latency)
- Unexpected changes (assets change without user knowledge)
- Can break workflows (new version incompatible)

Rejected: Violates no automatic operations principle. User control more important than convenience. Manual updates better.

Update both cache and config:

Example: Update cache and automatically update user config

- Download new asset version
- Update cache
- Update tasks.toml with new content
- User config always matches cache

Pros:

- Consistent (config always matches latest assets)
- No manual config updates (automatic sync)
- Simpler mental model (one source of truth)

Cons:

- Overwrites user customizations (modifications lost)
- Breaking changes forced (no user control)
- Surprising behavior (config changes unexpectedly)
- Loss of user work (custom modifications gone)
- Can break workflows (new version incompatible)

Rejected: Preserving user customizations critical. Cache-only updates keep config safe. Explicit manual updates better.

## Structure

Update command:

Syntax: start assets update [query]

- No query: Update all cached assets
- With query: Update assets matching query (substring matching)

Update algorithm:

1. Download catalog index
   - Fetch index.csv from raw.githubusercontent.com
   - Zero API calls (no rate limit)
   - Parse CSV into memory

2. Find cached assets
   - Glob ~/.config/start/assets/**/*.meta.toml
   - Read each .meta.toml file
   - Extract name, category, type, SHA

3. Filter by query (if provided)
   - Substring match against asset name, category, path
   - Only check matching assets
   - Skip non-matching assets

4. Compare SHAs
   - For each cached asset:
     - Find in index by path (assets/{type}/{category}/{name}.toml)
     - Compare local SHA with index SHA
     - If same: skip (up to date)
     - If different: mark for update

5. Download updates
   - For each asset with different SHA:
     - Download via raw.githubusercontent.com (zero API calls)
     - Download all asset files (.toml, .md if exists, .meta.toml)
     - Overwrite cached files
     - Report update

6. Report results
   - Show updated count
   - Show unchanged count
   - Remind user that config is unchanged
   - Suggest manual config review if updates occurred

Cache update behavior:

Update cache files:

- Overwrite existing .toml, .md, .meta.toml files
- Preserve directory structure
- Update SHA in .meta.toml
- Update timestamp in .meta.toml

Never modify config:

- ~/.config/start/tasks.toml unchanged
- ~/.config/start/roles.toml unchanged
- User customizations preserved
- User must manually review and apply changes

Effect on config:

- Tasks with prompt_file: Automatically use new content (references cache file)
- Tasks with inlined prompt: Require manual update (user must copy new content)

Selective updates:

Update all cached assets:

```bash
start assets update
```

Update assets matching query:

```bash
start assets update "commit"             # Match "commit" in name/category/path
start assets update git-workflow         # Match category
start assets update pre-commit-review    # Match specific asset name
```

Matching algorithm: Substring matching (case-insensitive)

- Search name, category, path
- Match anywhere in string
- Returns all matching assets from cache

Error handling:

Network unavailable:

- Message: "Cannot connect to GitHub"
- Show network error
- Exit with error code

Asset removed from catalog:

- Warning: "Asset not found in catalog: {name}"
- Keep cached version (don't delete)
- Continue with other assets

Invalid SHA in metadata:

- Error: "Invalid SHA in {name}.meta.toml"
- Skip this asset
- Suggest: Delete cache and re-download

Partial update failure:

- Update what succeeded
- Report failures separately
- Suggest retry

## Usage Examples

Update all cached assets:

```bash
$ start assets update

Downloading catalog index...
✓ Loaded index (46 assets)

Checking for asset updates...
  ✓ tasks/git-workflow/pre-commit-review (updated v1.0 → v1.1)
  ✓ roles/general/code-reviewer (up to date)
  ✓ tasks/git-workflow/pr-ready (up to date)

Cache updated with 1 new version.

Note: Your task configurations are unchanged.
Review changes and manually update tasks.toml if desired.
```

Update specific assets by query:

```bash
$ start assets update "commit"

Downloading catalog index...
Checking assets matching "commit"...
  ✓ tasks/git-workflow/pre-commit-review (updated)
  ✓ tasks/git-workflow/commit-message (up to date)

Cache updated with 1 new version.
```

No updates available:

```bash
$ start assets update

Downloading catalog index...
✓ Loaded index (46 assets)

Checking for asset updates...

✓ Update complete
  Updated: 0 assets
  Unchanged: 12 assets

All cached assets are up to date.
```

First run (no cache):

```bash
$ start assets update

Downloading catalog index...
✓ Loaded index (46 assets)

Checking for asset updates...

✓ Update complete
  Updated: 0 assets
  Unchanged: 0 assets

No cached assets found.

To download assets:
  - Browse catalog: start assets add
  - Use a task: start task <name>
```

Network error:

```bash
$ start assets update

Downloading catalog index...

Error: Cannot connect to GitHub

Network error: dial tcp: no route to host

Check your internet connection and try again.
```

Partial update failure:

```bash
$ start assets update

Downloading catalog index...
Checking for asset updates...

  ✓ Updated tasks/pre-commit-review
  ✗ Failed: roles/code-reviewer (download error)

✓ Update partially complete
  Updated: 1 asset
  Failed: 1 asset
  Unchanged: 10 assets

Errors:
  - roles/code-reviewer: network timeout

Try running 'start assets update' again to retry failed downloads.
```

Asset removed from catalog:

```bash
$ start assets update

Downloading catalog index...
Checking for asset updates...

  ⚠ Warning: tasks/deprecated-task not found in catalog
    (asset exists in cache but removed from GitHub)
    Keeping cached version.

✓ Update complete
  Updated: 0 assets
  Unchanged: 11 assets
  Not in catalog: 1 asset
```

Reviewing updated cache files:

```bash
# View updated asset content
$ cat ~/.config/start/assets/tasks/git-workflow/pre-commit-review.toml

# View updated metadata
$ cat ~/.config/start/assets/tasks/git-workflow/pre-commit-review.meta.toml
[metadata]
name = "pre-commit-review"
description = "Review staged changes before committing"
tags = ["git", "review", "quality", "pre-commit"]
sha = "b2c3d4e5f6789012345678901234567890abcdef"  # Updated SHA
created = "2025-01-10T00:00:00Z"
updated = "2025-01-15T14:30:00Z"  # Updated timestamp
```

Automatic use of updated cache:

```bash
# Task references cache via prompt_file
# tasks.toml:
[tasks.pre-commit-review]
prompt_file = "~/.config/start/assets/tasks/git-workflow/pre-commit-review.md"

# After update, task automatically uses new content
$ start task pre-commit-review
# Uses updated prompt from cache (no config change needed)
```

Manual config update required:

```bash
# Task has inlined prompt
# tasks.toml:
[tasks.pre-commit-review]
prompt = "Old prompt content here..."

# After cache update, user must manually update config
$ cat ~/.config/start/assets/tasks/git-workflow/pre-commit-review.md
# Review new content, copy to config if desired
```

Configuration:

```toml
# config.toml
[settings]
asset_repo = "grantcarthew/start"  # Repository to check for updates
```

Custom repository:

```toml
# config.toml
[settings]
asset_repo = "myorg/custom-assets"  # Use custom asset repository
```

## Updates

- 2025-01-17: Initial version aligned with schema; removed implementation code, Related Decisions, and Future Considerations sections; updated to use index.csv instead of Tree API per DR-034
