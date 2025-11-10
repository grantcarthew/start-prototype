# DR-034: GitHub Catalog API Strategy

**Date:** 2025-01-10
**Status:** Accepted
**Category:** Asset Management

## Decision

Use GitHub Tree API for catalog browsing with in-memory caching, and raw.githubusercontent.com URLs for asset downloads to avoid rate limits.

## What This Means

### API Strategy Overview

**For browsing catalog:**
- Use GitHub **Tree API** to get complete file structure
- Single API call gets entire repository tree with SHAs
- Cache tree in-memory for current session
- Rate limit: 60/hour (anonymous) or 5,000/hour (authenticated)

**For downloading assets:**
- Use **raw.githubusercontent.com** URLs for file content
- Direct HTTP GET, no authentication needed
- **Not subject to API rate limits** (different infrastructure)
- Fallback to Contents API if raw URL fails

**For update checking:**
- Use Tree API to get current SHAs
- Compare with cached `.meta.toml` SHA values
- Single API call per update check

### GitHub API Endpoints

**1. Tree API (catalog browsing):**
```
GET https://api.github.com/repos/{owner}/{repo}/git/trees/{branch}?recursive=1
```

Response includes:
```json
{
  "sha": "eda8147f...",
  "tree": [
    {
      "path": "assets/tasks/git-workflow/pre-commit-review.toml",
      "mode": "100644",
      "type": "blob",
      "sha": "a1b2c3d4...",
      "size": 1024
    },
    ...
  ]
}
```

**Benefits:**
- Complete directory structure in one call
- SHA for every file (version tracking)
- Size information (useful for validation)
- Fast recursive traversal

**2. Raw Content (downloading):**
```
GET https://raw.githubusercontent.com/{owner}/{repo}/{branch}/{path}
```

Response is plain text file content (no encoding needed).

**Benefits:**
- No API rate limiting
- No base64 decoding required
- Fast and simple
- Works unauthenticated

**3. Contents API (fallback):**
```
GET https://api.github.com/repos/{owner}/{repo}/contents/{path}?ref={branch}
```

Response includes:
```json
{
  "name": "pre-commit-review.toml",
  "path": "assets/tasks/git-workflow/pre-commit-review.toml",
  "sha": "a1b2c3d4...",
  "size": 1024,
  "content": "W3Rhc2td...",  // base64 encoded
  "encoding": "base64"
}
```

**Use when:**
- Raw URL fails (rare)
- Need additional metadata during download
- Fallback mechanism for reliability

### In-Memory Catalog Cache

**Cache structure:**
```go
var catalogCache struct {
    tree      *GitHubTree
    timestamp time.Time
    mu        sync.RWMutex
}
```

**Lifecycle:**
- **Created:** First time catalog is accessed in a session
- **Persists:** For duration of CLI invocation only
- **Cleared:** When CLI exits
- **Not saved:** Never written to disk

**Cache hit:** Asset resolution, catalog browsing, update checks use cached tree.

**Cache miss:** Fetch from GitHub on first access.

### Rate Limiting Strategy

**Anonymous (default):**
```
Limit: 60 requests/hour
Applies to: Tree API, Contents API
Does NOT apply to: raw.githubusercontent.com
```

**Authenticated (recommended):**
```
Limit: 5,000 requests/hour
Requires: GITHUB_TOKEN environment variable
Applies to: Tree API, Contents API
Does NOT apply to: raw.githubusercontent.com
```

**Rate limit checking:**
```bash
# GitHub responds with headers
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 42
X-RateLimit-Reset: 1762740583  # Unix timestamp
```

**Graceful handling:**
```go
func checkRateLimit(resp *http.Response) error {
    remaining := resp.Header.Get("X-RateLimit-Remaining")
    if remaining == "0" {
        reset := resp.Header.Get("X-RateLimit-Reset")
        resetTime := parseUnixTimestamp(reset)
        return fmt.Errorf("GitHub rate limit exceeded. Resets at %s", resetTime)
    }
    return nil
}
```

### Authentication

**Environment variable:**
```bash
export GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

**Usage in requests:**
```go
func fetchWithAuth(url string) (*http.Response, error) {
    req, _ := http.NewRequest("GET", url, nil)

    if token := os.Getenv("GITHUB_TOKEN"); token != "" {
        req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
    }

    return http.DefaultClient.Do(req)
}
```

**Benefits:**
- 5,000 requests/hour (vs 60)
- Recommended for all users
- Simple to configure
- Works with personal access tokens or fine-grained tokens

## Implementation

### Fetching Catalog Tree

```go
func fetchGitHubTree() (*GitHubTree, error) {
    repo := config.Settings["asset_repo"]  // "grantcarthew/start"
    url := fmt.Sprintf(
        "https://api.github.com/repos/%s/git/trees/main?recursive=1",
        repo,
    )

    req, _ := http.NewRequest("GET", url, nil)
    if token := os.Getenv("GITHUB_TOKEN"); token != "" {
        req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("fetch tree: %w", err)
    }
    defer resp.Body.Close()

    if err := checkRateLimit(resp); err != nil {
        return nil, err
    }

    var tree GitHubTree
    if err := json.NewDecoder(resp.Body).Decode(&tree); err != nil {
        return nil, fmt.Errorf("decode tree: %w", err)
    }

    return &tree, nil
}
```

### Downloading Asset Content

```go
func downloadAsset(githubPath string) ([]byte, error) {
    repo := config.Settings["asset_repo"]
    branch := "main"

    // Try raw.githubusercontent.com first (no rate limit)
    rawURL := fmt.Sprintf(
        "https://raw.githubusercontent.com/%s/%s/%s",
        repo, branch, githubPath,
    )

    resp, err := http.Get(rawURL)
    if err == nil && resp.StatusCode == 200 {
        defer resp.Body.Close()
        return io.ReadAll(resp.Body)
    }

    // Fallback to Contents API
    apiURL := fmt.Sprintf(
        "https://api.github.com/repos/%s/contents/%s?ref=%s",
        repo, githubPath, branch,
    )

    req, _ := http.NewRequest("GET", apiURL, nil)
    if token := os.Getenv("GITHUB_TOKEN"); token != "" {
        req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
    }

    resp, err = http.DefaultClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("download asset: %w", err)
    }
    defer resp.Body.Close()

    var content struct {
        Content  string `json:"content"`
        Encoding string `json:"encoding"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&content); err != nil {
        return nil, err
    }

    // Decode base64
    return base64.StdEncoding.DecodeString(content.Content)
}
```

### Finding Asset in Tree

```go
type GitHubTree struct {
    SHA  string `json:"sha"`
    Tree []struct {
        Path string `json:"path"`
        Mode string `json:"mode"`
        Type string `json:"type"`
        SHA  string `json:"sha"`
        Size int    `json:"size"`
    } `json:"tree"`
}

func (t *GitHubTree) Find(assetType, name string) string {
    // Look for: assets/{type}/*/{name}.toml
    prefix := fmt.Sprintf("assets/%s/", assetType)
    suffix := fmt.Sprintf("/%s.toml", name)

    for _, item := range t.Tree {
        if item.Type != "blob" {
            continue
        }
        if strings.HasPrefix(item.Path, prefix) && strings.HasSuffix(item.Path, suffix) {
            return item.Path
        }
    }
    return ""
}
```

## API Usage Patterns

### Pattern 1: Browse Catalog

```
User: start config task add

1. getCatalog() → Check in-memory cache
   - Cache miss → fetchGitHubTree() (1 API call)
   - Cache hit → Use cached tree (0 API calls)

2. Filter tree for tasks/

3. Display categories and tasks

4. User selects task

5. downloadAsset() via raw.githubusercontent.com (0 API calls)

Total API calls: 1 (tree) + 0 (download) = 1
```

### Pattern 2: Lazy Loading

```
User: start task pre-commit-review

1. Check local config → Not found
2. Check global config → Not found
3. Check cache → Not found

4. getCatalog() → In-memory cached from previous browse (0 API calls)
   - OR fetch if first access (1 API call)

5. Find in tree → assets/tasks/git-workflow/pre-commit-review.toml

6. downloadAsset() via raw.githubusercontent.com (0 API calls)

7. Add to config and run

Total API calls: 0-1 (tree) + 0 (download) = 0-1
```

### Pattern 3: Update Check

```
User: start update

1. getCatalog() → fetchGitHubTree() (1 API call)

2. For each cached .meta.toml:
   - Read local SHA
   - Find in tree by path
   - Compare SHAs
   - If different → downloadAsset() via raw.githubusercontent.com (0 API calls)

Total API calls: 1 (tree) + 0 (downloads) = 1
```

**Key insight:** Using raw.githubusercontent.com means unlimited asset downloads.

## Benefits

**Efficient:**
- ✅ Single Tree API call for complete catalog
- ✅ In-memory cache avoids repeated API calls
- ✅ Raw URLs bypass rate limits entirely

**Reliable:**
- ✅ Fallback from raw to Contents API
- ✅ Rate limit checking and clear errors
- ✅ Works authenticated or anonymous

**Simple:**
- ✅ Standard HTTP GET requests
- ✅ No external dependencies
- ✅ No git binary required

**Scalable:**
- ✅ Can handle hundreds of catalog assets
- ✅ Raw downloads have no rate limit
- ✅ Tree API scales with recursive flag

## Trade-offs Accepted

**Session-only cache:**
- ❌ Tree fetched on each CLI invocation
- **Mitigation:** Single fast API call, acceptable overhead

**No disk cache for catalog:**
- ❌ Can't browse catalog offline
- **Mitigation:** Per DR-026, catalog browsing requires network by design

**Anonymous rate limits:**
- ❌ 60 requests/hour without token
- **Mitigation:** Recommend GITHUB_TOKEN, raw URLs don't count

**GitHub dependency:**
- ❌ Catalog unavailable if GitHub down
- **Mitigation:** Cached assets work offline, manual config always possible

## Error Handling

### Rate Limit Exceeded

```
Error: GitHub rate limit exceeded

You've used 60/60 requests this hour.
Rate limit resets at: 2025-01-10 12:00:00 (in 45 minutes)

Solutions:
1. Set GITHUB_TOKEN for 5,000 requests/hour:
   export GITHUB_TOKEN=ghp_xxxxxxxxxxxx

2. Wait until rate limit resets

3. Use cached assets (if available)
```

### Network Error

```
Error: Cannot connect to GitHub

Network error: dial tcp: no route to host

Check your internet connection and try again.
```

### Authentication Error

```
Warning: GITHUB_TOKEN authentication failed

Using anonymous access (60 requests/hour).

Check your token:
  echo $GITHUB_TOKEN
  # Should start with ghp_ or github_pat_
```

## Configuration

**Settings in config.toml:**
```toml
[settings]
github_token_env = "GITHUB_TOKEN"    # Environment variable name
asset_repo = "grantcarthew/start"    # GitHub repository
```

**Future: Custom repositories**
```toml
# Not in v1, but structure allows it
asset_repo = "myorg/my-custom-assets"
```

## Testing

**API testing performed:**
- ✅ Tree API returns complete structure with SHAs
- ✅ raw.githubusercontent.com works without authentication
- ✅ Contents API provides fallback with base64 content
- ✅ Rate limit headers present in responses
- ✅ Anonymous access works (60/hour confirmed)

**Test repository:** `grantcarthew/start`

## Related Decisions

- [DR-031](./dr-031-catalog-based-assets.md) - Catalog architecture (API usage context)
- [DR-032](./dr-032-asset-metadata-schema.md) - Metadata schema (SHA comparison)
- [DR-033](./dr-033-asset-resolution-algorithm.md) - Resolution (catalog lookup)
- [DR-037](./dr-037-asset-updates.md) - Updates (SHA comparison from tree)
- [DR-026](./dr-026-offline-behavior.md) - Offline behavior (network requirements)
