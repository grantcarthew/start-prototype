# DR-021: GitHub Version Checking

- Date: 2025-01-06
- Status: Accepted
- Category: Version Management

## Problem

Users need to know when newer CLI versions are available. The system must:

- Check GitHub for latest CLI release version
- Compare current version with latest release
- Provide clear update instructions
- Work with different installation methods (brew, go install, manual)
- Handle rate limiting gracefully
- Not block or fail commands when version check fails
- Only check when user explicitly runs commands (no automatic background checks)
- Support both release and development builds
- Detect network errors and provide friendly messages
- Respect GitHub API rate limits

## Decision

Check GitHub Releases API for latest CLI version on every execution of `start doctor` and `start assets update`, with no caching.

API endpoint: `GET /repos/grantcarthew/start/releases/latest`

Version checks occur:
- `start doctor` - as part of health check
- `start assets update` - after updating cached assets

No caching (fresh API call each time), with rate limit protection (skip check if remaining < 10 requests).

## Why

User-initiated checks only:

- No automatic background checks (respectful, non-intrusive)
- Users explicitly run doctor or update (expect network calls)
- Clear intent to check system status
- No surprise network traffic

No caching for simplicity:

- Always shows current information (no stale data)
- No cache file management needed
- No cache invalidation logic
- Simpler implementation
- Rate limit check prevents abuse

Rate limiting protection:

- Check remaining API calls before version check
- Skip if < 10 remaining (preserve quota)
- Show friendly message when rate limited
- Respects `GH_TOKEN` env var for authenticated requests
- Anonymous: 60 requests/hour (usually sufficient)
- Authenticated: 5000 requests/hour

Non-blocking error handling:

- Version check failures don't stop commands
- Network errors show warning, continue
- Rate limiting shows friendly message, continue
- Parse errors log and show "Unknown"
- Only informational, not critical

Installation method detection:

- Detects how user installed (brew, go install, manual)
- Provides appropriate update command
- Helpful guidance for users
- Reduces friction for updates

Development build handling:

- Parses version strings like `v1.2.3-5-gabc1234`
- Extracts base version for comparison
- Shows status appropriate to build type
- Handles dirty builds, ahead-of-release builds

## Trade-offs

Accept:

- API call on every doctor/update execution
- Requires network connection for version check
- Slower than cached check (additional network round-trip)
- Two API calls (rate limit check + release check)
- May be rate limited if user runs frequently

Gain:

- Always shows current information (no stale cache)
- Simple implementation (no cache management)
- User-initiated only (respectful, no nagging)
- Clear error messages when network unavailable
- Rate limit protection (doesn't exhaust quota)
- Non-blocking (errors don't stop commands)

## Alternatives

Cached version check with TTL:

```go
// Cache latest release for 24 hours
type VersionCache struct {
    LatestVersion string
    CachedAt      time.Time
    TTL           time.Duration
}
```

- Pro: Fewer API calls (respects rate limits better)
- Pro: Faster response (no network call if cached)
- Con: May show stale data (up to TTL old)
- Con: Cache file management needed
- Con: Cache invalidation complexity
- Con: Must handle cache corruption/missing
- Rejected: Simplicity more important than speed for infrequent checks

Automatic background checks:

```go
// Check on every command execution
if time.Since(lastCheck) > 24*time.Hour {
    go checkVersionInBackground()
}
```

- Pro: Users always informed of updates
- Pro: More discoverable
- Con: Intrusive (surprise network calls)
- Con: Violates user expectations
- Con: Privacy concern (phones home automatically)
- Con: Can't control when checks happen
- Rejected: Too intrusive, violates user control

Version file bundled with binary:

```
start/
├── start (binary)
└── LATEST_VERSION
```

- Pro: No network call needed
- Pro: Always available offline
- Con: Always stale (bundle time, not current)
- Con: Deployment complexity (two files)
- Con: Defeats purpose (can't show actual latest)
- Rejected: Misses the point of checking latest

Environment variable for version URL:

```bash
START_VERSION_URL=https://custom.com/version start doctor
```

- Pro: Customizable for private forks
- Pro: Can point to different sources
- Con: Most users won't set this
- Con: Adds configuration complexity
- Con: Not discoverable
- Rejected: Over-engineering for uncommon use case

Only check on explicit command:

```bash
start version check
```

- Pro: Even more explicit user control
- Pro: Separate concern from doctor/update
- Con: Users must know about separate command
- Con: Less discoverable
- Con: Extra command to remember
- Rejected: Doctor/update are natural places to check

## Structure

API endpoint:

```
GET /repos/grantcarthew/start/releases/latest
```

Response:
```json
{
  "tag_name": "v1.3.0",
  "name": "Release v1.3.0",
  "published_at": "2025-01-06T10:30:00Z",
  "html_url": "https://github.com/grantcarthew/start/releases/tag/v1.3.0"
}
```

Rate limiting strategy:

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

Authentication:
- Respects `GH_TOKEN` environment variable
- Anonymous: 60 requests/hour
- Authenticated: 5000 requests/hour

Version comparison logic:

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

## Usage Examples

start doctor with update available:

```bash
$ start doctor
Version Information:
  CLI Version:     1.2.3
  Commit:          abc1234
  Build Date:      2025-01-06T10:30:00Z
  Go Version:      go1.22.0

Asset Information:
  Asset Cache:     ~/.config/start/assets/
  Cached Assets:   12 tasks, 8 roles, 4 agents

CLI Version Check:
  Current:         v1.2.3
  Latest Release:  v1.3.0
  Status:          ⚠ Update available
  Update via:      brew upgrade grantcarthew/tap/start
```

start assets update with version check:

```bash
$ start assets update
Checking for asset updates...
✓ Updated 3 cached assets
✓ Asset cache refreshed

CLI Version Check:
  Current:         v1.2.3
  Latest Release:  v1.3.0
  Status:          ⚠ Update available
  Update via:      brew upgrade grantcarthew/tap/start
```

When up to date:

```bash
$ start assets update
Asset cache is up to date

CLI Version Check:
  Current:         v1.3.0
  Latest Release:  v1.3.0
  Status:          ✓ Up to date
```

Build type handling:

Release build:
```
Current: v1.2.3
Status: ✓ Up to date
```

Development build:
```
Current: v1.2.3-5-gabc1234 (development build)
Status: ✓ Base version up to date
```

Dirty build:
```
Current: v1.2.3-dirty (uncommitted changes)
Status: ✓ Up to date
```

Ahead of release:
```
Current: v1.4.0-dirty
Latest:  v1.3.0
Status: ℹ Ahead of latest release
```

When rate limited:

```
CLI Version Check:
  Current:         v1.2.3
  Latest Release:  (rate limited - set GH_TOKEN to check)
  Status:          Unknown
```

When network error:

```
CLI Version Check:
  Current:         v1.2.3
  Latest Release:  (network error - check connection)
  Status:          Unknown
```

## Error Handling

Version check failures should not cause command to fail:

- Network errors: Show warning, continue
- API errors: Show warning, continue
- Rate limiting: Show friendly message, continue
- Parse errors: Log error, show "Unknown"

Only `start assets update` should fail if asset update fails. CLI version check is informational only.

## Installation Method Detection

Update suggestion detects installation method:

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

## Implementation

Package structure - create `internal/version/checker.go`:

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

Implementation checklist:

- Create `internal/version/checker.go`
- Implement `CheckLatestRelease()` with rate limit check
- Implement `CompareVersions()` with semver parsing
- Implement `DetectInstallMethod()` for update suggestions
- Add version check to `start doctor` command
- Add version check to `start assets update` command
- Handle all error cases gracefully (network, API, parsing)
- Add unit tests for version comparison logic
- Document GH_TOKEN usage for rate limit increases

## Updates

- 2025-01-17: Updated doctor output to show asset cache info instead of asset version tracking (aligns with catalog system per DR-031)
