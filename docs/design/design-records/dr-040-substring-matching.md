# DR-040: Substring Matching Algorithm for Asset Search

- Date: 2025-01-13
- Status: Accepted
- Category: Asset Management

## Problem

Users need to discover assets in the catalog by partial names, categories, descriptions, and tags. The search algorithm must address:

- Match type (exact vs prefix vs substring)
- Search fields (name only vs multiple fields)
- Minimum query length (prevent too-broad matches)
- Case sensitivity (case-sensitive vs case-insensitive)
- Result ranking (how to order matches)
- Performance (fast enough for large catalogs)
- Relationship to prefix matching (different purpose)
- Interactive vs non-interactive behavior
- Error handling (no matches, query too short)

## Decision

Implement substring-based search for discovering assets with 3-character minimum query length and multi-field matching (name, path, description, tags). Distinct from prefix matching used for flag value resolution.

Substring matching algorithm:

Definition: Query matches if found anywhere in search fields (not just beginning)
Case sensitivity: Case-insensitive matching (normalize to lowercase)
Minimum length: 3 characters minimum (prevents overly broad matches)

Search fields (in order):
1. Name: Asset filename without extension
2. Full path: Complete path including type and category
3. Description: Human-readable summary
4. Tags: Metadata keywords (semicolon-separated in index.csv)

Match priority: All matches returned, sorted by relevance
- Exact name matches (highest)
- Name substring matches
- Path matches
- Description matches
- Tag matches (lowest)
- Within each category: alphabetical by asset name

Multiple match handling:

Interactive (TTY):
- Display tree-structured selection menu with numbers
- User selects by number or quits
- Shows type/category grouping for clarity

Non-interactive (piped or --non-interactive flag):
- List all matches and exit with code 0
- No interactive prompt
- Machine-readable output

Single match:
- Auto-select in interactive mode
- Proceed to action (add/info/etc)

No matches:
- Error with suggestions (check spelling, try shorter query, use browse)

Query too short (< 3 chars):
- Error: "Query too short (minimum 3 characters)"
- Suggest using start assets browse for exploration

Relationship to prefix matching (DR-038):

Purpose: Substring for discovery, prefix for execution
Used by: Substring for search/add/info/update commands, prefix for --agent/--role/--task flags
Match type: Substring anywhere in string, prefix at beginning only
Sources: Substring searches GitHub catalog only, prefix searches local → global → cache → GitHub
Fields: Substring searches name/path/description/tags, prefix searches name only

## Why

Substring matching improves discovery:

- Find assets by partial names (commit finds pre-commit-review and commit-message)
- Search by category (workflow finds git-workflow assets)
- Search by description (security finds assets with security in description)
- Search by tags (review finds all review-related assets)
- Flexible querying without exact names

Multi-field search increases findability:

- Name matching (most common)
- Path matching (category-based search)
- Description matching (feature-based search)
- Tag matching (keyword-based search)
- Comprehensive coverage of metadata

3-character minimum prevents noise:

- Avoids overly broad matches (a matches 50% of catalog)
- Reduces false positives
- Still allows useful shortcuts (git, rev, sec)
- Clear error message guides users to browse mode

Case-insensitive matching simplifies UX:

- Users don't think about case
- commit and Commit and COMMIT all work
- Reduces friction (no case concerns)
- Standard pattern for search

Relevance sorting improves results:

- Exact name matches first (most relevant)
- Name substring matches next (still name-based)
- Path matches (category-level relevance)
- Description matches (content-level relevance)
- Tag matches (metadata-level relevance)
- Alphabetical within each level (predictable order)

Distinct from prefix matching serves different purposes:

- Substring for discovery (exploring what exists)
- Prefix for execution (running known assets)
- Different match algorithms for different goals
- Complementary features

Interactive selection improves UX:

- Tree-structured display (type → category grouping)
- Numbered selection (simple input)
- Auto-select for single match (efficiency)
- Cancel option (user control)

## Trade-offs

Accept:

- No fuzzy matching (typos don't match like comit won't find commit, but clear error messages suggest checking spelling, could add Levenshtein distance later)
- 3-character minimum (can't search for 2-letter queries like go, use start assets browse instead, prevents overly broad matches)
- No advanced ranking (simple field-based priority not TF-IDF, alphabetical within levels, good enough for expected catalog sizes)
- Linear search O(n) (slower for very large catalogs 1000+ assets, but in-memory search fast enough, could add indexing if needed)
- No query language (can't use type:tasks tag:git filters, keep simple for v1, add if users request)

Gain:

- Discovery-friendly (find by partial names, categories, descriptions, tags, flexible querying, relevant results sorted)
- Good UX (3-char minimum prevents garbage matches, interactive selection, auto-select for single match, clear errors)
- Simple implementation (easy to understand and maintain, fast in-memory search, graceful fallback without index)
- Fast performance (60-110ms typical with index, 110-220ms fallback to Tree API, acceptable for user-facing search)
- Complements prefix matching (different tools for different jobs, substring for discovery, prefix for execution)

## Alternatives

Prefix matching for search:

Example: Only match beginning of strings
```bash
start assets search "commit"   # Matches commit-message, NOT pre-commit-review
start assets search "pre"       # Matches pre-commit-review
```

Pros:
- Faster (can short-circuit on first character mismatch)
- More predictable (only matches beginning)
- Simpler algorithm

Cons:
- Less discoverable (must know how names start)
- Poor for category search (workflow doesn't find git-workflow)
- Poor for description search (can't search content)
- Less flexible (users must know naming patterns)

Rejected: Substring matching much better for discovery. Prefix matching available via flag resolution for execution.

Fuzzy matching with typo tolerance:

Example: Use Levenshtein distance for "did you mean?"
```bash
start assets search "comit"     # Did you mean "commit"?
start assets search "reviw"     # Did you mean "review"?
```

Pros:
- Typo-tolerant (helps with spelling mistakes)
- Better UX (users don't need perfect spelling)
- Discoverable (suggests corrections)

Cons:
- More complex (Levenshtein distance calculation)
- False positives possible (wrong suggestions)
- Performance impact (more computation)
- May suggest too many alternatives

Rejected: Keep v1 simple. Substring matching with clear error messages sufficient. Add fuzzy matching in v2 if users request.

No minimum query length:

Example: Allow 1 or 2 character queries
```bash
start assets search "a"   # Returns 50% of catalog
start assets search "ab"  # Still very broad
```

Pros:
- More flexible (no artificial restriction)
- Users can try any query

Cons:
- Overly broad matches (too many results)
- Poor UX (overwhelming result lists)
- Performance impact (searching for short strings)
- False positives (a matches almost everything)

Rejected: 3-character minimum better UX. Prevents garbage matches. Users can use browse mode for exploration.

Advanced ranking (TF-IDF, popularity):

Example: Score results by frequency, usage stats
- Boost frequently used assets
- Weight by how unique terms are
- Consider user's past usage

Pros:
- Better relevance (most useful results first)
- Learns from usage patterns
- More sophisticated

Cons:
- Complex implementation (TF-IDF calculation, usage tracking)
- Requires usage analytics (privacy concerns)
- More state to manage (statistics)
- Marginal benefit for small catalogs

Rejected: Simple field-based priority sufficient for v1. Alphabetical within levels predictable. Add advanced ranking if catalog grows large.

Query language with filters:

Example: Support structured queries
```bash
start assets search "type:tasks tag:git"
start assets search "commit AND review"
start assets search "category:git-workflow"
```

Pros:
- More powerful (precise queries)
- Advanced users can be specific
- Boolean operators (AND, OR, NOT)

Cons:
- Complex to implement (query parser)
- Harder to use (learning curve)
- Over-engineering for v1 (most users just want simple search)
- More documentation needed

Rejected: Simple substring search better for v1. Keep it simple. Add query language later if users request advanced features.

## Structure

Substring matching algorithm:

Match criteria:
- Query found anywhere in search field (not just beginning)
- Case-insensitive comparison (normalize to lowercase)
- Minimum 3 characters (error if < 3)

Search fields in priority order:
1. Name (exact match): Highest priority, exact string equality
2. Name (substring): High priority, query anywhere in name
3. Full path: Medium priority, includes type/category/name
4. Description: Lower priority, human-readable summary
5. Tags: Lowest priority, metadata keywords

Result sorting:
- Group by match field (name > path > description > tag)
- Within group: alphabetical by asset name
- Exact name matches always first

Data sources:

Primary: assets/index.csv (via raw.githubusercontent.com)
- Contains name, path, description, tags
- Enables rich multi-field search
- Fast in-memory search after download

Fallback: GitHub Tree API
- Contains name, path only (no description/tags)
- Limited search capability
- Degraded but functional

Interactive display:

Tree structure (TTY):
```
Found N matches:

type1/
  category1/
    [1] asset-name1    Description text
    [2] asset-name2    Description text
  category2/
    [3] asset-name3    Description text

type2/
  category3/
    [4] asset-name4    Description text

Select asset [1-N] (or 'q' to quit): _
```

List format (non-interactive):
```
Found N matches:

type1/category1/asset-name1
type1/category1/asset-name2
type1/category2/asset-name3
type2/category3/asset-name4
```

Error handling:

Query too short (<3 chars):
- Message: "Query too short (minimum 3 characters)"
- Suggestion: "Use 'start assets browse' for interactive browsing"
- Exit code: 1

No matches found:
- Message: "No matches found for '{query}'"
- Suggestions:
  - Check spelling
  - Try shorter or different query
  - Use start assets browse
  - Visit GitHub catalog URL
- Exit code: 1

Single match (auto-select in TTY):
- Message: "Found 1 match (exact): {path}"
- Message: "✓ Auto-selected"
- Proceed to action

Multiple matches (interactive):
- Display tree structure
- Prompt for selection
- Validate input
- Proceed with selected asset

Commands using substring matching:

start assets search <query>:
- Search and display results
- Returns all matches with metadata
- Exit after display

start assets add <query>:
- Search for asset
- Interactive selection if multiple matches
- Download and install selected

start assets info <query>:
- Search for asset
- Display detailed information
- Show .meta.toml contents

start assets update <query>:
- Search installed assets
- Check for updates
- Download new versions

Path construction:

Full path format: {type}/{category}/{name}
Example: tasks/git-workflow/pre-commit-review

Matches:
- git-workflow (category name)
- tasks/git (type + category prefix)
- workflow/pre (category + name prefix)
- Any substring of full path

Tag matching:

CSV format: git;review;quality;pre-commit
Split by semicolon, match each tag
Return true if any tag contains query substring

Performance:

Best case (with index):
- Download index.csv: ~50-100ms
- Parse CSV: ~5-10ms
- In-memory search: <1ms
- Total: ~60-110ms

Worst case (fallback to Tree API):
- Tree API call: ~100-200ms
- Parse response: ~10-20ms
- Path-only search: <1ms
- Total: ~110-220ms

Complexity: O(n) linear search through all assets
Acceptable for catalogs up to 1000+ assets

## Usage Examples

Name matching:

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

Path matching (category):

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

Description matching:

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

Tag matching:

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

Exact match (auto-select):

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

Query too short:

```bash
$ start assets search "ab"

Error: Query too short (minimum 3 characters)

Please provide at least 3 characters for meaningful search results.
Alternatively, use 'start assets browse' for interactive browsing.
```

No matches:

```bash
$ start assets search "nonexistent"

No matches found for 'nonexistent'

Suggestions:
- Check spelling
- Try a shorter or different query
- Use 'start assets browse' to explore interactively
- Visit: https://github.com/grantcarthew/start/tree/main/assets
```

Non-interactive mode:

```bash
$ start assets search "commit" --non-interactive

Found 5 matches:

tasks/git-workflow/commit-message
tasks/git-workflow/pre-commit-review
tasks/git-workflow/post-commit-hook
tasks/quality/commit-lint
roles/git/commit-specialist
```

Fallback (index unavailable):

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

## Updates

- 2025-01-17: Initial version aligned with schema; removed implementation code, Related Decisions, and Future Considerations sections
