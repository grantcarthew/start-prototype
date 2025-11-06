# DR-014: GitHub Asset Download Strategy

**Date:** 2025-01-06
**Status:** Accepted
**Category:** Asset Management

## Decision

Use GitHub Tree API with SHA-based caching for incremental asset updates

## Download Mechanism

```
1. GET /repos/{owner}/{repo}/git/trees/{branch}?recursive=1
   → Returns complete file tree with SHA hashes for all files

2. Load local asset-version.toml
   → Contains last downloaded commit + file SHAs

3. Compare local vs remote SHAs:
   - Skip files with matching SHA (already up to date)
   - Download only changed/new files via Contents API
   - Track removed files (exist locally but not in remote tree)

4. Update asset-version.toml with new commit + all file SHAs
```

## Asset Version Tracking File

Location: `~/.config/start/asset-version.toml`

Format:
```toml
# Asset version tracking - managed by 'start update'
# Last updated: 2025-01-06T10:30:00Z

commit = "abc123def456"
timestamp = "2025-01-06T10:30:00Z"
repository = "github.com/grantcarthew/start"
branch = "main"

[files]
"agents/claude.toml" = "a1b2c3d4e5f6..."
"agents/gemini.toml" = "e5f6g7h8i9j0..."
"roles/code-reviewer.md" = "i9j0k1l2m3n4..."
```

## API Endpoints Used

1. **Tree API** (discovery + SHAs):
   - `GET /repos/{owner}/{repo}/git/trees/{branch}?recursive=1`
   - Returns: Complete file tree with SHA-256 hashes
   - Single API call gets entire repository structure

2. **Contents API** (download):
   - `GET /repos/{owner}/{repo}/contents/{path}?ref={branch}`
   - Downloads individual files
   - One call per changed/new file

## Rate Limiting Strategy

- **Anonymous:** 60 requests/hour
- **Authenticated:** 5000 requests/hour via `GH_TOKEN` env var
- **Check before download:** Query `/rate_limit` endpoint
- **Abort if insufficient:** Error with reset time if < 50 requests remaining
- **Smart caching:** SHA comparison reduces API calls dramatically

## Incremental Update Example

First update (cold cache):
```
- Tree API: 1 call
- 28 asset files: 28 calls
- Total: 29 API calls
```

Subsequent update (3 files changed):
```
- Tree API: 1 call
- 3 changed files: 3 calls
- 25 unchanged: 0 calls (skipped via SHA match)
- Total: 4 API calls
```

## Benefits

- ✅ **Automatic discovery** - No manifest file to maintain
- ✅ **Incremental updates** - Only downloads changed files
- ✅ **Integrity verification** - SHA comparison validates files
- ✅ **Extensible** - New asset types discovered automatically
- ✅ **Efficient** - Caching minimizes API usage
- ✅ **No external dependencies** - Pure GitHub API, no git/tar needed

## Trade-offs Accepted

- ❌ Multiple API calls required (mitigated by caching)
- ❌ Rate limiting for anonymous users (GH_TOKEN recommended)
- ❌ First update downloads all files (one-time cost)

## Rationale

SHA-based caching provides best balance of:
- Developer ergonomics (no manifest to maintain)
- User experience (fast incremental updates)
- Resource efficiency (minimal API calls after first download)
- Implementation simplicity (no external dependencies)

## Related Decisions

- [DR-011](./dr-011-asset-distribution.md) - Asset distribution system
- [DR-015](./dr-015-atomic-updates.md) - Atomic update mechanism
- [DR-018](./dr-018-init-update-integration.md) - Init/update integration
