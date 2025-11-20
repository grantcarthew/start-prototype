# DR-039: Catalog Index File

- Date: 2025-01-13
- Status: Accepted
- Category: Asset Management

## Problem

Searching the GitHub catalog by description and tags requires accessing metadata distributed across hundreds of .meta.toml files. The search strategy must address:

- Performance (downloading 200+ individual .meta.toml files is slow)
- Rate limits (GitHub API limits constrain search operations)
- User experience (fast metadata-rich searching expected)
- Maintainability (solution must be simple to maintain)
- Reliability (must work even if solution partially fails)
- Freshness (stale data vs performance tradeoff)
- Contributor workflow (how to keep index updated)
- File format (CSV, JSON, or custom)
- Sorting (predictable diffs for git)

## Decision

Use pre-built CSV index file (assets/index.csv) for fast metadata-rich searching, generated via start assets index command, sorted alphabetically for predictable diffs.

Index file structure:

Location: assets/index.csv (in catalog repository)
Format: CSV with header row
Columns: type, category, name, description, tags, bin, sha, size, created, updated

Field definitions:

- type: Asset type (tasks, roles, agents, contexts)
- category: Category subdirectory (git-workflow, general, etc.)
- name: Asset name (filename without extension)
- description: Human-readable summary (from .meta.toml)
- tags: Semicolon-separated keywords (git;review;quality)
- bin: Binary name for agent auto-detection (agents only, empty for others)
- sha: Git blob SHA of content file (for update detection)
- size: File size in bytes
- created: ISO 8601 timestamp when created
- updated: ISO 8601 timestamp when last modified

Sorting: Alphabetical by type → category → name (predictable diffs)

Generation process:

Command: start assets index
Validation: Must be git repository with assets/ directory
Algorithm:

1. Scan assets/ directory recursively
2. Find all .meta.toml files
3. Parse each .meta.toml, extract metadata
4. For agents: extract bin from corresponding .toml file
5. Sort alphabetically (type → category → name)
6. Write CSV to assets/index.csv with header row
7. Report: "Generated index with N assets"

Usage pattern:

When to fetch: Every CLI execution needing catalog search
How to fetch: GET <https://raw.githubusercontent.com/{org}/{repo}/{branch}/assets/index.csv>
No local caching: File is small (10-50KB), always fresh data
Zero rate limits: Uses raw.githubusercontent.com (not API)

Fallback behavior:

If index.csv missing/corrupted:

1. Attempt to download index.csv
2. If fail: Log warning, fall back to Tree API
3. Tree API provides name/path only (no descriptions/tags)
4. Return degraded results with note

Source of truth:

- .meta.toml files are canonical (each asset has complete metadata)
- index.csv is derived file (for performance)
- If conflict: .meta.toml wins
- Update workflow: Edit .meta.toml → run start assets index → commit both

## Why

Index pattern proven by package managers:

- npm, cargo, homebrew all use index files
- Single download vs hundreds of requests
- Fast in-memory search
- Predictable performance
- Industry standard approach

Performance dramatically improved:

- Single ~10-50KB download vs 200+ .meta.toml requests
- Fast in-memory search with rich metadata
- Zero API rate limits (raw.githubusercontent.com)
- Typical download <100ms
- No network calls after initial download

Graceful degradation ensures reliability:

- Index missing: Fall back to Tree API (name/path only)
- New assets not indexed: Still discoverable via Tree API
- PRs without index regeneration: Assets still work
- Clear error messages guide contributors

Alphabetical sorting enables predictable diffs:

- Git diffs show meaningful changes (not random reordering)
- Code review easier (clear what added/removed/changed)
- Merge conflicts rare (alphabetical order stable)
- Consistent ordering across all tools

Simple CSV format reduces complexity:

- Standard format, widely supported
- Easy to inspect and debug
- No custom parsing needed (standard library)
- Human-readable for manual inspection
- RFC 4180 standard (automatic escaping)

Built-in generation tool improves contributor workflow:

- Clear process: start assets index command
- Validation catches errors early (missing fields, invalid TOML)
- Consistent output (sorted, formatted correctly)
- Simple to run (one command)

Always-fresh data preferred over caching:

- No stale results (always current catalog)
- No cache invalidation complexity
- File small enough (~10-50KB) to download each time
- Download fast (<100ms typical)
- Simplicity over micro-optimization

## Trade-offs

Accept:

- Manual generation (requires running start assets index after changes, but clear documentation and fallback to Tree API)
- CSV limitations (tags as semicolon-separated strings not native array, but standard CSV escaping handles it)
- No local cache (downloads index.csv on every use ~10-50KB, but file small and download fast <100ms, always fresh)
- Potential staleness (if contributor forgets to regenerate index, but fallback to Tree API finds new assets without rich metadata)
- Contributor discipline (PRs must include both .meta.toml and index.csv updates, but validation and fallback help)

Gain:

- Performance (single ~10-50KB download vs 200+ requests, fast in-memory search, zero rate limits)
- Reliability (graceful degradation if index missing, Tree API fallback, new assets still discoverable)
- Maintainability (standard CSV format proven by package managers, easy to inspect/debug, built-in generation tool)
- Contributor friendly (clear workflow with validation, even forgotten index regeneration still works via fallback)
- Predictable diffs (alphabetical sorting, meaningful git diffs, easier code review)
- Always fresh (no cache invalidation complexity, no stale data, simple design)

## Alternatives

Download all .meta.toml files individually:

Example: Fetch each .meta.toml on search

- Query Tree API for file list
- Download each .meta.toml (200+ HTTP requests)
- Parse and search in memory

Pros:

- No index file to maintain (one less thing)
- Always absolutely fresh (no derived file)
- Source of truth directly accessed

Cons:

- Extremely slow (200+ HTTP requests for large catalog)
- Poor user experience (long wait for searches)
- Doesn't scale (worse as catalog grows)
- Network-intensive (bandwidth waste)

Rejected: Performance unacceptable. Index.csv pattern proven by npm, cargo, homebrew.

Use GitHub Search API:

Example: Use GitHub's built-in code search

- Search API for file contents
- Filter results by repository/path
- Parse search results

Pros:

- No index file needed (GitHub maintains index)
- Powerful search (regex, code search)
- Always fresh (GitHub indexes automatically)

Cons:

- Rate limited (10-30 requests/min authenticated, 60/hour)
- Indexing delays (new files not immediately searchable)
- Requires authentication for reasonable limits
- Complex query syntax (not user-friendly)
- Unreliable for programmatic use (rate limits)

Rejected: Rate limits too restrictive. Index.csv provides better control and performance.

JSON index instead of CSV:

Example: Use assets/index.json with structured data

```json
[
  {
    "type": "tasks",
    "category": "git-workflow",
    "name": "pre-commit-review",
    "description": "Review staged changes",
    "tags": ["git", "review", "quality"],
    "sha": "abc123",
    ...
  }
]
```

Pros:

- Native array support (tags not semicolon-separated)
- More structured (nested objects possible)
- Familiar format (widely used)

Cons:

- Larger file size (JSON more verbose than CSV)
- Harder to inspect visually (no tabular view)
- Git diffs less readable (closing braces, commas)
- No real benefit (CSV handles our needs fine)

Rejected: CSV simpler, smaller, easier to diff. Tags as semicolon-separated acceptable.

Cache index.csv locally with TTL:

Example: Save index.csv with 5-minute expiry

- Download once, reuse for 5 minutes
- Refresh after TTL expires
- Faster repeated operations

Pros:

- Fewer downloads (reuse cached index)
- Faster for repeated searches
- Could work partially offline (within TTL)

Cons:

- Stale data possible (user sees old asset list)
- Cache invalidation complexity (when to refresh?)
- Confusing (new assets added, not visible until cache expires)
- More state to manage (TTL, expiry, refresh logic)

Rejected: Always-fresh data preferred. Index small (~10-50KB), download fast. Staleness worse than latency.

Auto-generate index via GitHub Actions:

Example: GitHub workflow regenerates index on every PR

- PR submitted with .meta.toml changes
- Action runs start assets index
- Commits index.csv automatically

Pros:

- Contributors can't forget (automated)
- Always up to date (regenerated on every change)
- No manual step (workflow handles it)

Cons:

- Bot commits add noise (extra commit per PR)
- Workflow complexity (setup and maintenance)
- Can fail silently (workflow breaks, PRs still merge)
- Contributor confusion (auto-commits not obvious)

Rejected: Manual generation acceptable. Clear documentation guides contributors. Fallback handles forgotten regeneration.

## Structure

Index file format:

CSV structure:

```csv
type,category,name,description,tags,bin,sha,size,created,updated
```

Example entries:

```csv
agents,anthropic,claude,Anthropic Claude AI via Claude Code CLI,claude;anthropic;ai,claude,a1b2c3d4,1024,2025-01-10T00:00:00Z,2025-01-10T00:00:00Z
roles,general,code-reviewer,Expert code reviewer focusing on security,review;security;quality,,b2c3d4e5,2048,2025-01-10T00:00:00Z,2025-01-12T00:00:00Z
tasks,git-workflow,commit-message,Generate conventional commit message,git;commit;conventional,,c3d4e5f6,1536,2025-01-10T00:00:00Z,2025-01-10T00:00:00Z
tasks,git-workflow,pre-commit-review,Review staged changes before committing,git;review;quality;pre-commit,,d4e5f6a1,2048,2025-01-10T00:00:00Z,2025-01-11T00:00:00Z
```

CSV escaping (RFC 4180):

Fields with commas (quoted):

```csv
tasks,workflow,my-task,"Review code, check tests, verify docs",testing,,abc123,1024,2025-01-13T00:00:00Z,2025-01-13T00:00:00Z
```

Fields with quotes (doubled quotes):

```csv
tasks,workflow,my-task,"Review ""special"" cases",testing,,abc123,1024,2025-01-13T00:00:00Z,2025-01-13T00:00:00Z
```

Multi-file assets:

One index entry per asset (not per file):

```
assets/tasks/git-workflow/pre-commit-review.toml
assets/tasks/git-workflow/pre-commit-review.md
assets/tasks/git-workflow/pre-commit-review.meta.toml
```

Index contains single entry:

```csv
tasks,git-workflow,pre-commit-review,Review staged changes,git;review,,abc123,2048,2025-01-10T00:00:00Z,2025-01-10T00:00:00Z
```

Fetching and parsing:

Download:

```
GET https://raw.githubusercontent.com/{org}/{repo}/{branch}/assets/index.csv
```

Parse CSV:

- Read CSV with header row
- Skip header, parse entries
- Split tags by semicolon
- Parse timestamps as ISO 8601
- Build in-memory asset list

Search:

- Substring match on name, description, tags (case-insensitive)
- Return all matching assets with full metadata

Fallback behavior:

Progressive degradation:

1. Attempt to download index.csv
   - Success: Use index for rich search
   - Fail: Log warning, fall back to Tree API

2. Fall back to Tree API directory listing
   - GET /repos/{org}/{repo}/git/trees/{branch}?recursive=1
   - Filter paths: assets/{type}/_/_.toml (exclude .meta.toml)
   - Extract names from paths only
   - No descriptions/tags (limited search)

3. Return results with degraded metadata
   - Name and path available
   - Description: "(metadata unavailable)"
   - Tags: empty

New assets not in index:

1. Search index.csv: not found
2. Search Tree API: found
3. Return with note: "(Not in index - metadata unavailable)"
4. Suggest: "Ask maintainer to run 'start assets index'"

Generation command:

Validation:

- Check for .git/ directory (must be git repo)
- Check for assets/ directory (must be catalog structure)
- Error if missing: "Must be run in catalog repository"

Algorithm:

1. Scan assets/ recursively for .meta.toml files
2. Parse each .meta.toml (extract metadata)
3. Derive type and category from path
4. For agents: Extract bin from .toml file
5. Validate required fields present
6. Sort alphabetically (type → category → name)
7. Write CSV to assets/index.csv with header
8. Report: "Generated index with N assets"

Error handling:

- Missing .meta.toml: Warning, skip asset
- Invalid TOML: Error with filename and line number
- Missing required fields: Error with which fields and file
- File system errors: Error and exit

Size estimates:

Minimal catalog (28 assets):

- 28 rows × ~120 bytes/row = ~3.4 KB
- With header and CSV overhead = ~4 KB

Large catalog (200 assets):

- 200 rows × ~120 bytes/row = ~24 KB
- With header and CSV overhead = ~30 KB

Very large catalog (500 assets):

- 500 rows × ~120 bytes/row = ~60 KB
- With header and CSV overhead = ~70 KB

All well within acceptable download sizes (<100ms on typical connections).

## Usage Examples

Searching with index (normal operation):

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

Searching without index (fallback to Tree API):

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

New asset not in index:

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

Generating index:

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

Generation with validation error:

```bash
$ start assets index

Validating repository structure...
✗ Error: Not a git repository

This command must be run in the catalog repository.
Run: git clone https://github.com/grantcarthew/start.git
```

Generation with TOML parsing error:

```bash
$ start assets index

Scanning assets/...
✗ Error: tasks/broken/bad-task.meta.toml:5: invalid TOML syntax

Fix the metadata file and try again.
```

Contributor workflow:

```bash
# Clone catalog repository
$ git clone https://github.com/grantcarthew/start.git
$ cd start

# Add new asset with metadata
$ mkdir -p assets/tasks/my-category
$ cat > assets/tasks/my-category/my-task.meta.toml <<EOF
[metadata]
name = "my-task"
description = "My new task"
tags = ["workflow", "custom"]
sha = "..."
created = "2025-01-13T00:00:00Z"
updated = "2025-01-13T00:00:00Z"
EOF

# Generate index
$ start assets index

# Output:
# Scanning assets/...
# Found 46 assets
# ✓ Generated index with 46 assets
# Updated: assets/index.csv

# Commit both .meta.toml and index.csv
$ git add assets/
$ git commit -m "Add my-task and regenerate index"
```

## Updates

- 2025-01-17: Initial version aligned with schema; removed implementation code, Related Decisions, and Future Considerations sections
