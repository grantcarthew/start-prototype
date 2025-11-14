# DR-021: GitHub Version Checking

**Date:** 2025-01-06
**Status:** Accepted
**Category:** Version Management

## Decision

Check GitHub Releases API for latest CLI version on every execution of `start doctor` and `start assets update`, with no caching.

## API Endpoint

**Latest Release:**
```
GET /repos/grantcarthew/start/releases/latest
```

**Response:**
```json
{
  "tag_name": "v1.3.0",
  "name": "Release v1.3.0",
  "published_at": "2025-01-06T10:30:00Z",
  "html_url": "https://github.com/grantcarthew/start/releases/tag/v1.3.0"
}
```

## When Checks Occur

### `start doctor`
Checks and reports CLI version status as part of health check:
```bash
$ start doctor
Version Information:
  CLI Version:     1.2.3
  Commit:          abc1234
  Build Date:      2025-01-06T10:30:00Z

Asset Information:
  Asset Version:   1.1.0 (commit: def5678)
  Last Updated:    2 days ago
  Status:          ✓ Up to date

CLI Version Check:
  Current:         v1.2.3
  Latest Release:  v1.3.0
  Status:          ⚠ Update available
  Update via:      brew upgrade grantcarthew/tap/start
```

### `start assets update`
Checks CLI version after updating assets:
```bash
$ start assets update
Checking for asset updates...
✓ Downloaded 3 updated files
✓ Asset library updated to commit def5678

CLI Version Check:
  Current:         v1.2.3
  Latest Release:  v1.3.0
  Status:          ⚠ Update available
  Update via:      brew upgrade grantcarthew/tap/start
```

If up to date:
```bash
$ start assets update
Asset library is up to date (commit: def5678)

CLI Version Check:
  Current:         v1.3.0
  Latest Release:  v1.3.0
  Status:          ✓ Up to date
```

## Rate Limiting Strategy

Before checking latest release:

```
1. GET /rate_limit
   → Check remaining API calls

2. If remaining < 10:
   → Skip version check
   → Show: "Latest Release: (rate limited - set GH_TOKEN to check)"
   → Continue with command

3. If remaining >= 10:
   → GET /repos/grantcarthew/start/releases/latest
   → Compare versions
```

**Authentication:**
- Respect `GH_TOKEN` environment variable (like DR-014)
- Anonymous: 60 requests/hour (usually sufficient)
- Authenticated: 5000 requests/hour

**No Caching:**
- Every execution makes fresh API call
- Simple implementation (no cache file management)
- Always shows current data
- Rate limit check prevents abuse

## Version Comparison Logic

### Semantic Version Parsing

```go
// Strip 'v' prefix from tags
currentVer := strings.TrimPrefix(version.Version, "v")
latestVer := strings.TrimPrefix(release.TagName, "v")

// Handle development builds: v1.2.3-5-gabc1234
// Extract base version: 1.2.3
if strings.Contains(currentVer, "-") {
    parts := strings.Split(currentVer, "-")
    currentVer = parts[0]
    isDev = true
}

// Compare using semver library
current := semver.MustParse(currentVer)
latest := semver.MustParse(latestVer)

if current.LessThan(latest) {
    status = "Update available"
} else if current.Equal(latest) {
    status = "Up to date"
} else {
    status = "Ahead of latest release" // Local build newer than release
}
```

### Build Type Handling

**Release build:**
```
Current: v1.2.3
Status: ✓ Up to date
```

**Development build:**
```
Current: v1.2.3-5-gabc1234 (development build)
Status: ✓ Base version up to date
```

**Dirty build:**
```
Current: v1.2.3-dirty (uncommitted changes)
Status: ✓ Up to date
```

**Ahead of release:**
```
Current: v1.4.0-dirty
Latest:  v1.3.0
Status: ℹ Ahead of latest release
```

## Output Formats

### When Update Available

```
CLI Version Check:
  Current:         v1.2.3
  Latest Release:  v1.3.0
  Status:          ⚠ Update available
  Update via:      brew upgrade grantcarthew/tap/start
```

### When Up to Date

```
CLI Version Check:
  Current:         v1.3.0
  Latest Release:  v1.3.0
  Status:          ✓ Up to date
```

### When Rate Limited

```
CLI Version Check:
  Current:         v1.2.3
  Latest Release:  (rate limited - set GH_TOKEN to check)
  Status:          Unknown
```

### When Network Error

```
CLI Version Check:
  Current:         v1.2.3
  Latest Release:  (network error - check connection)
  Status:          Unknown
```

## Implementation Notes

### Error Handling

Version check failures should **not** cause command to fail:

- Network errors: Show warning, continue
- API errors: Show warning, continue
- Rate limiting: Show friendly message, continue
- Parse errors: Log error, show "Unknown"

Only `start assets update` should fail if **asset** update fails. CLI version check is informational only.

### Installation Method Detection

Update suggestion should detect installation method:

```go
// Check if installed via Homebrew
if _, err := exec.LookPath("brew"); err == nil {
    if output, err := exec.Command("brew", "list", "--formula").Output(); err == nil {
        if strings.Contains(string(output), "grantcarthew/tap/start") {
            return "brew upgrade grantcarthew/tap/start"
        }
    }
}

// Check if installed via go install
if gopath := os.Getenv("GOPATH"); gopath != "" {
    return "go install github.com/grantcarthew/start/cmd/start@latest"
}

// Default fallback
return "See https://github.com/grantcarthew/start#installation"
```

### Package Structure

Create `internal/version/checker.go`:

```go
package version

type ReleaseInfo struct {
    TagName     string
    PublishedAt time.Time
    HTMLURL     string
}

// CheckLatestRelease queries GitHub for the latest release
func CheckLatestRelease(ctx context.Context) (*ReleaseInfo, error) {
    // Implementation
}

// CompareVersions returns status and message
func CompareVersions(current, latest string) (status, message string) {
    // Implementation
}

// DetectInstallMethod returns update command
func DetectInstallMethod() string {
    // Implementation
}
```

## Benefits

- ✅ **Always fresh** - No stale cached data
- ✅ **Simple** - No cache file management
- ✅ **User-initiated** - Only checks when user runs doctor/update
- ✅ **Non-blocking** - Errors don't stop commands
- ✅ **Informative** - Clear update instructions
- ✅ **Respectful** - No automatic nagging or background checks

## Trade-offs Accepted

- ❌ API call on every doctor/update (mitigated: rate limit check)
- ❌ Requires network connection (acceptable: shows friendly error)
- ❌ Slower than cached check (acceptable: users expect network call)

## Rationale

No caching simplifies implementation and ensures users always see current information. Since checks are user-initiated (not automatic), the API call overhead is acceptable. Rate limiting protection prevents abuse of GitHub API.

## Related Decisions

- [DR-020](./dr-020-version-injection.md) - Binary version injection (source of current version)
- [DR-014](./dr-014-github-tree-api.md) - GitHub API usage patterns and rate limiting

## Implementation Checklist

- [ ] Create `internal/version/checker.go`
- [ ] Implement `CheckLatestRelease()` with rate limit check
- [ ] Implement `CompareVersions()` with semver parsing
- [ ] Implement `DetectInstallMethod()` for update suggestions
- [ ] Add version check to `start doctor` command
- [ ] Add version check to `start assets update` command
- [ ] Handle all error cases gracefully (network, API, parsing)
- [ ] Add unit tests for version comparison logic
- [ ] Document GH_TOKEN usage for rate limit increases
