# DR-040: Substring Matching Algorithm for Asset Search

**Date:** 2025-01-13
**Status:** Accepted
**Category:** Asset Management

## Decision

Implement substring-based search for discovering assets in the GitHub catalog, with a 3-character minimum query length and multi-field matching (name, path, description, tags). Distinct from prefix matching used for flag value resolution.

## Problem

Users need to discover assets in the catalog by:
- **Partial names** - "commit" should find "pre-commit-review" and "commit-message"
- **Categories** - "workflow" should find assets in git-workflow category
- **Descriptions** - "security" should find assets with security in description
- **Tags** - "review" should find all review-related assets

This is different from command execution (which uses prefix matching) - this is for **discovery and exploration**.

## Solution

### Substring Matching Algorithm

**Definition:** Query matches if found **anywhere** in the search fields (not just at the beginning).

**Case sensitivity:** Case-insensitive matching (normalize both query and fields to lowercase).

**Minimum length:** 3 characters minimum (prevents overly broad matches).

**Search fields (in order):**
1. **Name** - Asset filename without extension
2. **Full path** - Complete path including type and category
3. **Description** - Human-readable summary
4. **Tags** - Metadata keywords (semicolon-separated in index.csv)

**Match priority:** All matches returned, sorted by relevance (exact name > name substring > path > description > tags).

### Algorithm Implementation

```go
type SearchResult struct {
    Asset Asset
    MatchField string  // "name", "path", "description", "tag"
    MatchType  string  // "exact", "substring"
}

func searchAssets(query string, index *CatalogIndex) ([]SearchResult, error) {
    // Validate minimum length
    if len(query) < 3 {
        return nil, fmt.Errorf("query too short (minimum 3 characters)")
    }

    // Normalize query
    query = strings.ToLower(strings.TrimSpace(query))

    var results []SearchResult

    for _, asset := range index.Assets {
        // Build full path: type/category/name
        fullPath := fmt.Sprintf("%s/%s/%s", asset.Type, asset.Category, asset.Name)

        // Check name (exact match first)
        if strings.ToLower(asset.Name) == query {
            results = append(results, SearchResult{
                Asset: asset,
                MatchField: "name",
                MatchType: "exact",
            })
            continue
        }

        // Check name (substring)
        if strings.Contains(strings.ToLower(asset.Name), query) {
            results = append(results, SearchResult{
                Asset: asset,
                MatchField: "name",
                MatchType: "substring",
            })
            continue
        }

        // Check full path (substring)
        if strings.Contains(strings.ToLower(fullPath), query) {
            results = append(results, SearchResult{
                Asset: asset,
                MatchField: "path",
                MatchType: "substring",
            })
            continue
        }

        // Check description (substring)
        if strings.Contains(strings.ToLower(asset.Description), query) {
            results = append(results, SearchResult{
                Asset: asset,
                MatchField: "description",
                MatchType: "substring",
            })
            continue
        }

        // Check tags (substring in any tag)
        for _, tag := range asset.Tags {
            if strings.Contains(strings.ToLower(tag), query) {
                results = append(results, SearchResult{
                    Asset: asset,
                    MatchField: "tag",
                    MatchType: "substring",
                })
                break  // Only count once per asset
            }
        }
    }

    // Sort results by relevance
    sortByRelevance(results)

    return results, nil
}
```

### Sorting by Relevance

**Priority order:**
1. Exact name matches (highest relevance)
2. Name substring matches
3. Path matches
4. Description matches
5. Tag matches (lowest relevance)

Within each category, sort alphabetically by asset name.

```go
func sortByRelevance(results []SearchResult) {
    sort.Slice(results, func(i, j int) bool {
        // Priority scores
        scoreI := matchScore(results[i])
        scoreJ := matchScore(results[j])

        if scoreI != scoreJ {
            return scoreI > scoreJ  // Higher score first
        }

        // Same priority - sort alphabetically
        return results[i].Asset.Name < results[j].Asset.Name
    })
}

func matchScore(result SearchResult) int {
    if result.MatchType == "exact" {
        return 100
    }

    switch result.MatchField {
    case "name":
        return 80
    case "path":
        return 60
    case "description":
        return 40
    case "tag":
        return 20
    default:
        return 0
    }
}
```

### 3-Character Minimum

**Rationale:**
- Prevents overly broad matches (e.g., "a" matching 50% of catalog)
- Reduces false positives
- Still allows useful shortcuts ("git", "rev", "sec")

**Behavior for queries < 3 chars:**

```bash
$ start assets search "ab"

Error: Query too short (minimum 3 characters)

Please provide at least 3 characters for meaningful search results.
Alternatively, use 'start assets browse' for interactive browsing.
```

**Alternative:** Fall back to interactive browse mode instead of error.

### Multiple Match Handling

**Interactive (TTY):**
Display tree-structured selection menu:

```bash
$ start assets search "commit"

Found 5 matches:

tasks/
  git-workflow/
    [1] commit-message         Generate conventional commit message
    [2] pre-commit-review      Review staged changes before committing
    [3] post-commit-hook       Post-commit validation workflow

  quality/
    [4] commit-lint            Lint commit messages for conventions

roles/
  git/
    [5] commit-specialist      Expert in git commit best practices

Select asset [1-5] (or 'q' to quit): _
```

**Non-interactive (piped, --non-interactive):**
List all matches and exit with code 0:

```bash
$ start assets search "commit" --non-interactive

Found 5 matches:

tasks/git-workflow/commit-message
tasks/git-workflow/pre-commit-review
tasks/git-workflow/post-commit-hook
tasks/quality/commit-lint
roles/git/commit-specialist
```

**Single match:**
In interactive mode, auto-select if only one match:

```bash
$ start assets search "pre-commit-review"

Found 1 match:
tasks/git-workflow/pre-commit-review

✓ Selected: tasks/git-workflow/pre-commit-review
[proceeds to action - add/info/etc.]
```

**No matches:**
Error with suggestions:

```bash
$ start assets search "nonexistent"

No matches found for 'nonexistent'

Suggestions:
- Check spelling
- Try a shorter query (minimum 3 characters)
- Use 'start assets browse' to explore interactively
- Use 'start assets search --help' for search tips
```

## Relationship to DR-038 (Prefix Matching)

These are **complementary but distinct** matching algorithms:

| Feature | DR-038: Prefix Matching | DR-040: Substring Matching |
|---------|------------------------|---------------------------|
| **Purpose** | Command execution | Asset discovery |
| **Used by** | `--agent`, `--role`, `--task` flags | `start assets search/add/info/update` |
| **Match type** | Beginning of string only | Anywhere in string |
| **Sources** | Local → global → cache → GitHub | GitHub catalog only (via index.csv) |
| **Behavior** | Exact → prefix (short-circuit) | Substring across all assets |
| **Fields** | Asset name only | Name, path, description, tags |
| **Examples** | `--task pre` matches "pre-commit" | `search commit` matches "pre-**commit**-review" |

**Why different?**

- **Prefix matching** is for **speed** when you already know what you want
- **Substring matching** is for **discovery** when you're exploring

**Example comparison:**

```bash
# Prefix matching (DR-038) - for command execution
start task pre              # Matches "pre-commit-review" (starts with "pre")
start task commit           # No match (doesn't start with "commit")

# Substring matching (DR-040) - for discovery
start assets search pre     # Matches "pre-commit-review", "prepare-release"
start assets search commit  # Matches "pre-commit-review", "commit-message", "post-commit-hook"
```

## Data Source: Catalog Index (DR-039)

Substring matching uses the catalog index file for rich metadata:

**Primary:** `assets/index.csv` (downloaded from GitHub)
- Contains: name, path, description, tags
- Enables searching by all fields
- Fast in-memory search

**Fallback:** GitHub Tree API (if index.csv unavailable)
- Contains: name, path only
- Limited search (no description/tags)
- Still functional, just degraded

See [DR-039](./dr-039-catalog-index.md) for index structure and usage.

## Commands Using Substring Matching

**`start assets search <query>`**
- Primary use case: search and display results
- Returns all matches with metadata

**`start assets add <query>`**
- Search for asset by query
- Interactive selection if multiple matches
- Downloads and installs selected asset

**`start assets info <query>`**
- Search for asset by query
- Display detailed information
- Shows .meta.toml contents

**`start assets update <query>`**
- Search for installed assets by query
- Check for updates
- Download new versions

All use the same substring matching algorithm.

## Examples

### Example 1: Name Matching

```bash
$ start assets search "review"

Found 4 matches:

tasks/
  git-workflow/
    [1] code-review            Review code for quality and best practices
    [2] pre-commit-review      Review staged changes before committing
    [3] doc-review             Review and improve documentation

roles/
  general/
    [4] code-reviewer          Expert code reviewer focusing on security

Select asset [1-4]: _
```

**Matches:**
- "code-**review**" - name substring
- "pre-commit-**review**" - name substring
- "doc-**review**" - name substring
- "code-**reviewer**" - name substring

### Example 2: Path Matching (Category)

```bash
$ start assets search "workflow"

Found 6 matches:

tasks/
  git-workflow/
    [1] commit-message         Generate conventional commit message
    [2] pre-commit-review      Review staged changes before committing
    [3] post-commit-hook       Post-commit validation workflow

  ci-workflow/
    [4] test-pipeline          Continuous integration test workflow
    [5] deploy-check           Deployment workflow validation
    [6] release-prep           Prepare release workflow

Select asset [1-6]: _
```

**Matches:**
- All assets in "git-**workflow**" and "ci-**workflow**" categories (path matching)

### Example 3: Description Matching

```bash
$ start assets search "security"

Found 3 matches:

tasks/
  quality/
    [1] security-audit         Comprehensive security vulnerability scan
    [2] dependency-check       Check dependencies for security issues

roles/
  specialized/
    [3] security-expert        Security specialist with penetration testing expertise

Select asset [1-3]: _
```

**Matches:**
- Descriptions containing "**security**"

### Example 4: Tag Matching

```bash
$ start assets search "quality"

Found 5 matches:

tasks/
  git-workflow/
    [1] pre-commit-review      Review staged changes before committing
                               Tags: git, review, quality, pre-commit
    [2] code-review            Review code for quality and best practices
                               Tags: review, quality, best-practices

  testing/
    [3] test-coverage          Analyze test coverage and quality
                               Tags: testing, quality, coverage

  quality/
    [4] lint-check             Run code quality linters
                               Tags: quality, linting, standards
    [5] complexity-check       Check code complexity and quality metrics
                               Tags: quality, metrics, complexity

Select asset [1-5]: _
```

**Matches:**
- Assets with "**quality**" in tags

### Example 5: Exact Match (Auto-Select)

```bash
$ start assets add "pre-commit-review"

Searching catalog...
Found 1 match (exact):
  tasks/git-workflow/pre-commit-review

✓ Auto-selected

Downloading pre-commit-review...
✓ Cached to ~/.config/start/assets/tasks/git-workflow/
✓ Added to global config

Use 'start task pre-commit-review' to run.
```

### Example 6: Query Too Short

```bash
$ start assets search "ab"

Error: Query too short (minimum 3 characters)

Please provide at least 3 characters for meaningful search results.
Alternatively, use 'start assets browse' for interactive browsing.
```

### Example 7: No Matches

```bash
$ start assets search "nonexistent"

No matches found for 'nonexistent'

Suggestions:
- Check spelling
- Try a shorter or different query
- Use 'start assets browse' to explore interactively
- Visit: https://github.com/grantcarthew/start/tree/main/assets
```

### Example 8: Non-Interactive Mode

```bash
$ start assets search "commit" --non-interactive

Found 5 matches:

tasks/git-workflow/commit-message
tasks/git-workflow/pre-commit-review
tasks/git-workflow/post-commit-hook
tasks/quality/commit-lint
roles/git/commit-specialist
```

### Example 9: Fallback (Index Unavailable)

```bash
$ start assets search "commit"

Searching catalog...
⚠ Index unavailable, using directory listing (limited metadata)

Found 3 matches (name/path only):

tasks/git-workflow/commit-message
  (Metadata unavailable - not in index)

tasks/git-workflow/pre-commit-review
  (Metadata unavailable - not in index)

tasks/git-workflow/post-commit-hook
  (Metadata unavailable - not in index)

Note: Description and tag search unavailable without index.
```

**Limited matching:** Only searches name and path, no description or tags.

## Performance Characteristics

**Best case (cached index):**
- Download index.csv: ~50-100ms
- Parse CSV: ~5-10ms
- In-memory search: <1ms
- **Total: ~60-110ms**

**Worst case (index unavailable, fallback to Tree API):**
- Tree API call: ~100-200ms
- Parse response: ~10-20ms
- Path-only search: <1ms
- **Total: ~110-220ms**

**Complexity:**
- O(n) linear search through all assets
- Acceptable for catalogs up to 1000+ assets
- In-memory = fast even on large catalogs

## Benefits

**Discovery-friendly:**
- ✅ Find assets by partial names, categories, descriptions, tags
- ✅ Flexible querying without knowing exact names
- ✅ Relevant results sorted by match quality

**User experience:**
- ✅ 3-character minimum prevents garbage matches
- ✅ Interactive selection for multiple matches
- ✅ Auto-select for single matches
- ✅ Clear error messages with suggestions

**Implementation:**
- ✅ Simple algorithm, easy to understand and maintain
- ✅ Fast in-memory search
- ✅ Graceful fallback without index

## Trade-offs Accepted

**No fuzzy matching:**
- ❌ Typos don't match (e.g., "comit" won't find "commit")
- **Mitigation:** Clear error messages, suggestions to check spelling
- **Future:** Could add Levenshtein distance for "did you mean?"

**3-character minimum:**
- ❌ Can't search for 2-letter queries (e.g., "go")
- **Mitigation:** Use `start assets browse` for exploration
- **Trade-off:** Prevents overly broad matches

**No ranking by relevance sophistication:**
- ❌ Simple field-based priority, not TF-IDF or other advanced scoring
- **Mitigation:** Alphabetical sort within each priority level
- **Trade-off:** Good enough for small to medium catalogs

**Linear search (O(n)):**
- ❌ Slower for very large catalogs (1000+ assets)
- **Mitigation:** In-memory search is fast enough for expected catalog sizes
- **Future:** Could add indexing data structure if needed

## Implementation Details

### Normalization

```go
func normalize(s string) string {
    return strings.ToLower(strings.TrimSpace(s))
}
```

**Removes:**
- Leading/trailing whitespace
- Case differences

**Preserves:**
- Internal whitespace
- Special characters (hyphens, underscores)

### Tag Matching

Tags stored as semicolon-separated string in CSV:
```csv
...,git;review;quality;pre-commit,...
```

Parse and match:
```go
func matchesTags(tags string, query string) bool {
    tagList := strings.Split(tags, ";")
    for _, tag := range tagList {
        if strings.Contains(normalize(tag), query) {
            return true
        }
    }
    return false
}
```

### Path Construction

Full path for matching:
```go
fullPath := fmt.Sprintf("%s/%s/%s", asset.Type, asset.Category, asset.Name)
// Example: "tasks/git-workflow/pre-commit-review"
```

Matches queries like:
- "git-workflow" (category)
- "tasks/git" (type + category prefix)
- "workflow/pre" (category + name prefix)

## Related Decisions

- [DR-038](./dr-038-flag-value-resolution.md) - Prefix matching for flag values (different purpose)
- [DR-039](./dr-039-catalog-index.md) - Catalog index file (data source for matching)
- [DR-034](./dr-034-github-catalog-api.md) - GitHub API strategy (Tree API fallback)
- [DR-035](./dr-035-interactive-browsing.md) - Interactive browsing (alternative to search)

## Future Considerations

**Fuzzy matching:**
- Could add Levenshtein distance for typo tolerance
- Example: "comit" → "Did you mean 'commit'?"
- Trade-off: More complex implementation vs better UX

**Advanced ranking:**
- Could use TF-IDF or other relevance scoring
- Boost frequently used assets
- Trade-off: Complexity vs marginal UX improvement

**Query language:**
- Could support filters: `type:tasks tag:git`
- Boolean operators: `commit AND review`
- Trade-off: More powerful vs more complex

**Caching search results:**
- Could cache popular queries
- Trade-off: Memory vs speed improvement

**Current stance:** Ship with described substring matching. Monitor usage patterns and iterate based on user feedback. Simple is better than complex for v1.
