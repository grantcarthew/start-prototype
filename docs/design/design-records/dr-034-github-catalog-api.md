# DR-034: GitHub Catalog API Strategy

- Date: 2025-01-10
- Status: Accepted
- Category: Asset Management

## Problem

The CLI needs to access the GitHub catalog for browsing, downloading, and updating assets. The API strategy must address:

- API selection (which GitHub APIs to use for different operations)
- Rate limiting (avoid hitting GitHub API limits)
- Authentication (support both anonymous and authenticated access)
- Download efficiency (fast asset downloads)
- Caching strategy (reduce redundant API calls)
- Index download (catalog index.csv for fast searching)
- Offline behavior (work without network when possible)
- Error handling (rate limits, network failures, authentication issues)
- Update detection (efficiently check for asset updates)
- Reliability (fallback mechanisms when APIs fail)

## Decision

Use GitHub Tree API for catalog browsing with in-memory caching, raw.githubusercontent.com URLs for index and asset downloads to avoid rate limits, and Contents API as fallback.

API strategy by operation:

Catalog index download (fast metadata searching):

- Use raw.githubusercontent.com to download assets/index.csv
- Endpoint: GET <https://raw.githubusercontent.com/{owner}/{repo}/main/assets/index.csv>
- Downloaded fresh on every search/browse operation (no local caching)
- Not subject to API rate limits
- Enables fast substring matching across asset metadata

Browsing catalog structure (directory tree, file SHAs):

- Use GitHub Tree API (single call gets entire repository tree)
- Endpoint: GET /repos/{owner}/{repo}/git/trees/{branch}?recursive=1
- Returns complete file structure with SHAs and sizes
- Cache tree in-memory for current session
- Rate limit: 60/hour (anonymous) or 5,000/hour (authenticated)
- Used for resolution and as fallback when index unavailable

Downloading assets (file content):

- Use raw.githubusercontent.com URLs (not subject to API rate limits)
- Endpoint: GET <https://raw.githubusercontent.com/{owner}/{repo}/{branch}/{path}>
- Direct HTTP GET, no authentication needed
- Fallback to Contents API if raw URL fails
- No rate limit on raw URLs

Update checking (SHA comparison):

- Download index.csv via raw.githubusercontent.com (contains SHAs)
- Compare index SHAs with cached .meta.toml SHA values
- No Tree API call needed (index has all SHAs)
- Zero API calls against rate limit

In-memory tree cache:

Cache structure:

- Stores GitHub Tree API response
- Persists for duration of CLI invocation only
- Cleared when CLI exits
- Never written to disk

Cache lifecycle:

- Created: First time Tree API accessed (if needed)
- Used: Fallback when index.csv unavailable, resolution operations
- Cleared: When CLI exits

Note: Most operations use index.csv instead of Tree API cache.

Authentication:

Environment variable: GITHUB_TOKEN

- Optional but recommended
- Increases rate limit from 60 to 5,000 requests/hour for Tree/Contents APIs
- Uses standard GitHub authentication header
- Works with personal access tokens or fine-grained tokens
- Not needed for raw.githubusercontent.com URLs

Rate limiting:

Anonymous access:

- Limit: 60 requests/hour for API calls
- Applies to: Tree API, Contents API
- Does NOT apply to: raw.githubusercontent.com (index.csv, asset downloads)

Authenticated access (with GITHUB_TOKEN):

- Limit: 5,000 requests/hour for API calls
- Applies to: Tree API, Contents API
- Does NOT apply to: raw.githubusercontent.com (index.csv, asset downloads)

Rate limit headers:

- X-RateLimit-Limit: Total limit
- X-RateLimit-Remaining: Requests remaining
- X-RateLimit-Reset: Unix timestamp when resets

## Why

Index-first approach eliminates most API calls:

- Index.csv downloaded via raw URLs (no rate limit)
- Contains all metadata (name, description, tags, SHA)
- Enables fast searching without API calls
- Downloaded fresh each time (always current)
- Small file (~10-50KB, fast download)

Raw URLs bypass rate limits entirely:

- Unlimited downloads (not counted against API limits)
- No base64 decoding required (direct file content)
- Fast and simple (plain HTTP GET)
- Works unauthenticated (no token needed)
- Proven reliable infrastructure

Tree API as fallback provides reliability:

- Complete directory structure in one call if index unavailable
- SHA for every file (version tracking built-in)
- Size information (useful for validation)
- In-memory cache for session
- Graceful degradation when index missing

Minimal API usage overall:

- Search/browse: 0 API calls (uses index.csv via raw URL)
- Download assets: 0 API calls (raw URLs)
- Update check: 0 API calls (index.csv has SHAs)
- Resolution: 0 API calls (uses index.csv or cached tree)
- Most operations never hit rate limits

Authenticated access scales better when needed:

- 5,000 requests/hour (83x more than anonymous)
- Simple configuration (environment variable)
- Standard across GitHub tooling (CLI, Actions, etc.)
- Recommended for power users
- No cost (uses existing GitHub account)

Fallback mechanism provides reliability:

- Raw URL fails: Fall back to Contents API
- Index unavailable: Fall back to Tree API
- Contents API provides base64-encoded content
- Graceful degradation
- Higher reliability overall

## Trade-offs

Accept:

- Index downloaded every search (not cached locally, but ~10-50KB is fast and ensures freshness)
- Session-only tree cache (tree fetched if needed per invocation, but most operations use index instead)
- No disk cache for catalog (can't browse offline, but catalog browsing requires network by design)
- Anonymous rate limits low (60 requests/hour without token, but most operations use raw URLs that don't count)
- GitHub dependency (catalog unavailable if GitHub down, but cached assets work offline and manual config always possible)

Gain:

- Zero API calls for common operations (search/browse/update/download all use raw URLs, no rate limits)
- Efficient with index.csv (single small file download, complete metadata, fast substring search)
- Reliable downloads (fallback from raw to Contents API, fallback from index to Tree API)
- Simple implementation (standard HTTP GET requests, no external dependencies, no git binary)
- Scalable design (handles hundreds of assets, raw downloads unlimited, index approach proven)
- Fast operations (raw URL downloads, minimal network round-trips, no rate limit concerns)
- Tree API rarely needed (index provides everything for most operations)

## Alternatives

Use Tree API for all operations:

Example: Always use Tree API, never use index.csv

- Download tree for every operation
- Parse tree for metadata
- No index.csv needed

Pros:

- Single API approach
- Always has latest data
- No index maintenance

Cons:

- Counts against rate limit (60/hour anonymous)
- No description/tags in tree (just file paths and SHAs)
- Must download .meta.toml files individually for search
- Multiple API calls for rich searching
- Slower user experience

Rejected: Index.csv via raw URL is much more efficient. Zero API calls for searching.

Use Contents API for everything:

Example: Use /repos/{owner}/{repo}/contents/{path} for all operations

- Browse catalog: Recursive calls to get directory structure
- Download assets: Get base64-encoded content from API
- Update checking: Get file metadata for SHAs

Pros:

- Single API to learn and use
- Consistent approach
- All operations authenticated same way

Cons:

- Multiple API calls to browse (one per directory level)
- Downloads count against rate limit (60/hour anonymous)
- Base64 decoding overhead (extra processing)
- Slower catalog browsing (many round-trips)
- Hits rate limits faster
- No index.csv support

Rejected: Raw URLs + index.csv is much more efficient. Zero API calls for most operations.

Clone repository with git:

Example: Use git clone or git archive to get catalog

- Clone full repository
- Read files from disk
- Use git for updates

Pros:

- Complete offline catalog (all files local)
- No API calls after initial clone
- Can use git for version control

Cons:

- Requires git binary (external dependency)
- Large initial download (entire repository)
- Disk space overhead (full git history)
- Complex update logic (git pull, merge conflicts)
- Overkill for text file catalog

Rejected: Too heavy for simple asset catalog. Raw URLs + index.csv is lighter and simpler.

Cache index.csv locally:

Example: Save index.csv to ~/.cache/start/index.csv with TTL

- Download once, reuse for period
- Refresh after TTL expires

Pros:

- Fewer downloads (reuse cached index)
- Faster for repeated operations
- Could work partially offline

Cons:

- Stale data possible (cache may be outdated)
- Cache invalidation complexity (when to refresh?)
- Index shows old assets (confusing)
- More state to manage

Rejected: Always-fresh data preferred. Index is small (~10-50KB) and downloads fast.

## Structure

GitHub API endpoints:

Index download (primary for search/browse):

- Endpoint: GET <https://raw.githubusercontent.com/{owner}/{repo}/main/assets/index.csv>
- Returns: CSV with asset metadata (type, category, name, description, tags, sha, etc.)
- Rate limit: None (not subject to API limits)
- Use: Catalog searching, browsing, update checking

Raw content (asset downloads):

- Endpoint: GET <https://raw.githubusercontent.com/{owner}/{repo}/{branch}/{path}>
- Returns: Plain text file content (no encoding)
- Rate limit: None (not subject to API limits)
- Use: Primary method for downloading assets

Tree API (fallback for structure):

- Endpoint: GET <https://api.github.com/repos/{owner}/{repo}/git/trees/{branch}?recursive=1>
- Returns: Complete file tree with paths, SHAs, sizes, types
- Rate limit: 60/hour anonymous, 5,000/hour authenticated
- Use: Fallback when index.csv unavailable, in-memory cache for session

Contents API (fallback for downloads):

- Endpoint: GET <https://api.github.com/repos/{owner}/{repo}/contents/{path}?ref={branch}>
- Returns: File metadata plus base64-encoded content
- Rate limit: 60/hour anonymous, 5,000/hour authenticated
- Use: Fallback when raw URL fails

In-memory tree cache:

Lifecycle:

- Created if Tree API accessed (rare - most operations use index)
- Persists for duration of process only
- Cleared when CLI exits
- Never saved to disk

Usage:

- Fallback when index.csv unavailable
- Resolution operations if tree already fetched
- All operations prefer index.csv first

Benefits:

- Rarely needed (index.csv handles most cases)
- Zero disk I/O
- Simple lifecycle management

Authentication:

Environment variable: GITHUB_TOKEN

- Hardcoded (not configurable)
- Optional but recommended for power users
- Standard across GitHub tooling
- Only affects Tree/Contents API calls
- Not needed for raw URLs (index.csv, asset downloads)

Rate limits:

- Anonymous: 60 requests/hour (Tree/Contents API only)
- Authenticated: 5,000 requests/hour (Tree/Contents API only)
- Raw URLs: Unlimited (index.csv, asset downloads)

Configuration (config.toml):

```toml
[settings]
asset_repo = "grantcarthew/start"  # GitHub repository
```

Download strategy:

For index.csv:

1. Always use raw.githubusercontent.com
2. No fallback needed (if unavailable, use Tree API for structure)

For assets:

1. Primary: raw.githubusercontent.com (no rate limit, direct content)
2. Fallback: Contents API (if raw fails, base64 decode)

For catalog structure:

1. Primary: index.csv via raw URL (contains all metadata)
2. Fallback: Tree API (if index unavailable)

Rate limit handling:

Check headers in API responses:

- X-RateLimit-Remaining: Requests left
- X-RateLimit-Reset: When limit resets

Error on limit exceeded:

- Show requests used (e.g., 60/60)
- Show reset time (e.g., "in 45 minutes")
- Suggest GITHUB_TOKEN for higher limits
- Note: Most operations use raw URLs (no rate limit)

Error handling:

Rate limit exceeded (rare - most operations use raw URLs):

- Message: "GitHub rate limit exceeded"
- Show: Requests used, reset time
- Solutions: Set GITHUB_TOKEN, wait, use cached assets
- Note: Search/browse/download don't count against limits

Network error:

- Message: "Cannot connect to GitHub"
- Show: Network error details
- Solution: Check internet connection

Authentication error:

- Warning: "GITHUB_TOKEN authentication failed"
- Fallback: Use anonymous access (60 requests/hour)
- Solution: Check token format (should start with ghp_or github_pat_)
- Note: Raw URLs work without authentication

Index unavailable:

- Fallback to Tree API automatically
- Warning: "Index unavailable, using directory listing (limited metadata)"
- Continue with degraded search (name/path only, no descriptions/tags)

## Usage Examples

Search catalog (zero API calls):

```bash
$ start assets search "commit"

# Behind the scenes:
# 1. Download index.csv via raw.githubusercontent.com (0 API calls, no rate limit)
# 2. Parse CSV into memory
# 3. Substring search across name, description, tags
# 4. Display results with rich metadata
#
# Total API calls: 0
# Rate limit impact: None
```

Browse catalog (zero API calls):

```bash
$ start assets browse

# Behind the scenes:
# 1. Download index.csv via raw.githubusercontent.com (0 API calls, no rate limit)
# 2. Parse and organize by type/category
# 3. Display interactive browser
# 4. User selects task
# 5. downloadAsset() via raw.githubusercontent.com (0 API calls, no rate limit)
#
# Total API calls: 0
# Rate limit impact: None
```

Update checking (zero API calls):

```bash
$ start assets update

# Behind the scenes:
# 1. Download index.csv via raw.githubusercontent.com (0 API calls, no rate limit)
# 2. For each cached .meta.toml:
#    - Read local SHA
#    - Find in index by path
#    - Compare SHAs
#    - If different → downloadAsset() via raw.githubusercontent.com (0 API calls)
#
# Total API calls: 0
# Rate limit impact: None
```

Lazy loading (zero API calls):

```bash
$ start task pre-commit-review

# Behind the scenes:
# 1. Check local config → Not found
# 2. Check global config → Not found
# 3. Check cache → Not found
# 4. Download index.csv via raw.githubusercontent.com (0 API calls)
# 5. Find asset in index
# 6. downloadAsset() via raw.githubusercontent.com (0 API calls)
# 7. Add to config and run
#
# Total API calls: 0
# Rate limit impact: None
```

Fallback when index unavailable (1 API call):

```bash
$ start assets search "commit"

# Behind the scenes:
# 1. Attempt to download index.csv via raw URL → 404 Not Found
# 2. Fallback to Tree API (1 API call, counts against rate limit)
# 3. Search tree paths (name/path only, no descriptions/tags)
# 4. Display degraded results
#
# Total API calls: 1 (Tree API)
# Rate limit impact: 1 request used

Searching catalog...
⚠ Index unavailable, using directory listing (limited metadata)

Found 3 matches:
  tasks/git-workflow/commit-message
  tasks/git-workflow/pre-commit-review
  (Metadata unavailable - not in index)
```

API call comparison:

Without index.csv:

```bash
# Search catalog: 1 Tree API call + N .meta.toml downloads
# Browse catalog: 1 Tree API call + N .meta.toml downloads
# Update check: 1 Tree API call
# Download assets: 0 API calls (raw URLs)
# Total for typical session: ~10-20 API calls
```

With index.csv (current design):

```bash
# Search catalog: 0 API calls (index via raw URL)
# Browse catalog: 0 API calls (index via raw URL)
# Update check: 0 API calls (index via raw URL)
# Download assets: 0 API calls (raw URLs)
# Total for typical session: 0 API calls
```

Configuration:

```toml
# config.toml
[settings]
asset_repo = "grantcarthew/start"  # GitHub repository
```

```bash
# Environment (optional, only needed if hitting Tree/Contents API rate limits)
export GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

Network error:

```bash
$ start assets browse

Error: Cannot connect to GitHub

Network error: dial tcp: no route to host

Check your internet connection and try again.
```

## Updates

- 2025-01-17: Initial version aligned with schema; incorporated index.csv download strategy from DR-039
