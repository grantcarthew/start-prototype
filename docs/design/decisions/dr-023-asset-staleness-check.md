# DR-023: Asset Staleness Checking

**Date:** 2025-01-06
**Status:** Accepted
**Category:** Asset Management

## Decision

`start doctor` always checks GitHub for the latest asset commit SHA and compares it to the local version, with no caching.

## Check Strategy

**GitHub Comparison (Commit SHA):**
- Read local commit SHA from `~/.config/start/asset-version.toml`
- Fetch latest commit SHA from main branch via GitHub API
- Compare: if different, updates are available
- No caching - fresh check every time

**Not timestamp-based:**
- Don't use age alone (e.g., "assets > 30 days old")
- Timestamp shown for context, but not used for status determination
- Only warn if actual updates exist on GitHub

## API Endpoint

**Get latest commit on main branch:**

```
GET /repos/grantcarthew/start/commits/main
```

**Response (relevant fields):**
```json
{
  "sha": "def5678abc1234...",
  "commit": {
    "author": {
      "date": "2025-01-06T10:30:00Z"
    },
    "message": "Update code-reviewer role prompt"
  }
}
```

**Efficient alternative (lighter response):**
```
GET /repos/grantcarthew/start/git/refs/heads/main
```

**Response:**
```json
{
  "ref": "refs/heads/main",
  "object": {
    "sha": "def5678abc1234...",
    "type": "commit"
  }
}
```

Use the second endpoint (refs API) - lighter payload, faster response.

## Display Format

### When Updates Available

```bash
Asset Information:
  Current Commit:  abc1234 (45 days ago)
  Latest Commit:   def5678 (2 hours ago)
  Status:          ⚠ Updates available
  Action:          Run 'start update' to refresh
```

### When Up to Date

```bash
Asset Information:
  Current Commit:  abc1234 (2 days ago)
  Latest Commit:   abc1234
  Status:          ✓ Up to date
```

### When Rate Limited

```bash
Asset Information:
  Current Commit:  abc1234 (2 days ago)
  Latest Commit:   (rate limited - set GH_TOKEN to check)
  Status:          Unknown
```

### When Network Error

```bash
Asset Information:
  Current Commit:  abc1234 (2 days ago)
  Latest Commit:   (network error - check connection)
  Status:          Unknown
```

### When Assets Not Initialized

```bash
Asset Information:
  Status:          ⚠ Not initialized
  Action:          Run 'start init' to download assets
```

## Rate Limiting Strategy

Same approach as DR-021 (CLI version checking):

```
1. GET /rate_limit
   → Check remaining API calls

2. If remaining < 10:
   → Skip asset check
   → Show: "Latest Commit: (rate limited - set GH_TOKEN to check)"
   → Continue with command

3. If remaining >= 10:
   → GET /repos/grantcarthew/start/git/refs/heads/main
   → Compare commits
```

**Authentication:**
- Respect `GH_TOKEN` environment variable
- Anonymous: 60 requests/hour
- Authenticated: 5000 requests/hour

**No Caching:**
- Every `start doctor` execution makes fresh API call
- Consistent with DR-021 (CLI version checking)
- Always shows current data

## Implementation

### Commit Comparison Logic

```go
// Load local asset version
localCommit := loadAssetVersion() // From asset-version.toml
localTime := loadAssetTimestamp() // From asset-version.toml

// Fetch latest commit from main branch
latestCommit, latestTime, err := fetchLatestCommit(ctx)
if err != nil {
    // Handle network/rate limit errors
    return showUnknownStatus(localCommit, localTime, err)
}

// Compare commits (first 7 chars for display)
if localCommit == latestCommit {
    status = "✓ Up to date"
} else {
    status = "⚠ Updates available"
}

// Display information
showAssetInfo(localCommit, localTime, latestCommit, latestTime, status)
```

### Commit SHA Display

Show short SHA (7 characters) for readability:
- Full SHA: `def5678abc1234567890abcdef123456`
- Display: `def5678`

This matches git's default short SHA format.

### Timestamp Display

Show relative time for user-friendliness:
- `2 hours ago`
- `3 days ago`
- `45 days ago`
- `2 months ago`

Use the timestamp from `asset-version.toml` (when assets were last updated locally).

## Error Handling

Asset check failures should **not** cause `start doctor` to fail:

- **Network errors:** Show warning, continue with other checks
- **API errors:** Show warning, continue
- **Rate limiting:** Show friendly message, continue
- **Missing asset-version.toml:** Show "Not initialized", suggest `start init`

The doctor command should complete successfully even if asset check fails.

## Integration with `start doctor`

Full `start doctor` output:

```bash
$ start doctor

Version Information:
  CLI Version:     1.2.3
  Commit:          abc1234
  Build Date:      2025-01-06T10:30:00Z
  Go Version:      go1.22.0

Configuration:
  Global Config:   ✓ ~/.config/start/config.toml
  Local Config:    ✗ Not found
  Validation:      ✓ Valid

Asset Information:
  Current Commit:  abc1234 (45 days ago)
  Latest Commit:   def5678 (2 hours ago)
  Status:          ⚠ Updates available
  Action:          Run 'start update' to refresh

CLI Version Check:
  Current:         v1.2.3
  Latest Release:  v1.3.0
  Status:          ⚠ Update available
  Update via:      brew upgrade grantcarthew/tap/start

Environment:
  GH_TOKEN:        ✓ Set (authenticated API access)
  EDITOR:          ✓ Set (vim)
  Shell:           /bin/bash

Overall Status:   ⚠ Updates available
Exit Code:        1 (warnings present)
```

## Package Structure

Update `internal/assets/checker.go`:

```go
package assets

type CommitInfo struct {
    SHA       string
    ShortSHA  string    // First 7 chars
    Timestamp time.Time
}

// CheckLatestCommit queries GitHub for the latest commit on main branch
func CheckLatestCommit(ctx context.Context) (*CommitInfo, error) {
    // Implementation
}

// CompareCommits returns status based on local vs latest
func CompareCommits(local, latest *CommitInfo) (status, message string) {
    // Implementation
}
```

## Benefits

- ✅ **Accurate** - Only warns when updates actually exist
- ✅ **Consistent** - Same pattern as CLI version check (DR-021)
- ✅ **Simple** - No caching, no staleness threshold
- ✅ **Informative** - Shows current vs latest commits
- ✅ **Always fresh** - No stale cached data

## Trade-offs Accepted

- ❌ Requires network connection (acceptable: shows friendly error)
- ❌ API call every time (mitigated: rate limit check, efficient endpoint)
- ❌ Slower than timestamp-only check (acceptable: users expect network call)

## Rationale

Checking GitHub directly ensures `start doctor` only warns when updates are actually available. This is more useful than a simple age-based warning. Consistency with DR-021 (CLI version checking) keeps the implementation simple and predictable.

## Related Decisions

- [DR-014](./dr-014-github-tree-api.md) - GitHub API usage patterns
- [DR-021](./dr-021-github-version-check.md) - CLI version checking (same pattern)
- [DR-022](./dr-022-asset-branch-strategy.md) - Assets from main branch

## Implementation Checklist

- [ ] Create `internal/assets/checker.go`
- [ ] Implement `CheckLatestCommit()` with rate limit check
- [ ] Implement `CompareCommits()` for status determination
- [ ] Add asset check to `start doctor` command
- [ ] Handle all error cases gracefully
- [ ] Format timestamps as relative time ("2 days ago")
- [ ] Display short SHA (7 chars)
- [ ] Add unit tests for commit comparison logic
- [ ] Document GH_TOKEN usage for rate limit increases
