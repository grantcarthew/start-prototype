# start assets search

## Name

start assets search - Search GitHub catalog by keyword

## Synopsis

```bash
start assets search <query>
```

## Description

Search the GitHub asset catalog by keyword and display matching assets without installing them. Terminal-based, non-interactive command that lists matching assets and exits.

Uses substring matching across multiple fields:

- Asset name
- Full path (type/category/name)
- Description
- Tags

Results are displayed grouped by type and category, sorted by relevance. This is a read-only discovery tool - use `start assets add` to install assets.

**Key differences from `start assets add`:**

- **Read-only** - Lists matches without installation
- **Non-interactive** - Prints results and exits
- **Discovery-focused** - For exploration before installation

## Arguments

**\<query\>** (required)
: Search query string. Minimum 3 characters.

**Query matching:**

- Case-insensitive
- Substring match (not prefix)
- Searches: name, path, description, tags

**Examples:**

```bash
start assets search commit       # Finds "pre-commit-review", "commit-message"
start assets search workflow     # Finds assets in "*-workflow" categories
start assets search security     # Finds assets with "security" in description
start assets search quality      # Finds assets tagged with "quality"
```

## Output

### Matches Found

```bash
$ start assets search "commit"

Found 5 matches:

tasks/
  git-workflow/
    commit-message         Generate conventional commit message
    pre-commit-review      Review staged changes before committing
    post-commit-hook       Post-commit validation workflow

  quality/
    commit-lint            Lint commit messages for conventions

roles/
  git/
    commit-specialist      Expert in git commit best practices
```

**Format:**

- Grouped by type and category (hierarchical)
- Asset name + description (if available)
- Sorted by relevance (exact name > name substring > path > description > tags)
- Within each group, alphabetically by name

### Single Match

```bash
$ start assets search "pre-commit-review"

Found 1 match:

tasks/
  git-workflow/
    pre-commit-review      Review staged changes before committing
```

No auto-selection or installation (read-only command).

### No Matches

```bash
$ start assets search "nonexistent"

No matches found for 'nonexistent'

Suggestions:
- Check spelling
- Try a shorter or different query
- Use 'start assets browse' to explore the catalog
- Visit: https://github.com/grantcarthew/start/tree/main/assets
```

Exit code: 2

### Query Too Short

```bash
$ start assets search "ab"

Error: Query too short (minimum 3 characters)

Please provide at least 3 characters for meaningful search results.
Alternatively, use 'start assets browse' for interactive browsing.
```

Exit code: 1

### With Metadata

```bash
$ start assets search "security" --verbose

Found 3 matches:

tasks/
  quality/
    security-audit         Comprehensive security vulnerability scan
      Tags: security, audit, vulnerabilities, owasp
      Size: 3.2 KB
      Updated: 2025-01-10

    dependency-check       Check dependencies for security issues
      Tags: security, dependencies, vulnerabilities
      Size: 2.1 KB
      Updated: 2025-01-08

roles/
  specialized/
    security-expert        Security specialist with penetration testing expertise
      Tags: security, pentesting, owasp, vulnerabilities
      Size: 4.5 KB
      Updated: 2025-01-12
```

### Network Error

```bash
$ start assets search "commit"

Searching catalog...
✗ Network error

Cannot connect to GitHub:
  dial tcp: no route to host

Check your internet connection and try again.
```

Exit code: 1

### Index Unavailable (Fallback)

```bash
$ start assets search "commit"

Searching catalog...
⚠ Index unavailable, using directory listing (limited metadata)

Found 3 matches (name/path only):

tasks/
  git-workflow/
    commit-message         (Metadata unavailable)
    pre-commit-review      (Metadata unavailable)
    post-commit-hook       (Metadata unavailable)

Note: Description and tag search unavailable without index.
Use 'start assets browse' for full catalog exploration.
```

**Degraded search:** Only name and path matching, no description or tag search.

## Exit Codes

**0** - Success (matches found and displayed)

**1** - Network error or invalid query

**2** - No matches found

## Flags

**--verbose**
: Show detailed metadata for each match (tags, size, dates).

## Examples

### Search by Name

```bash
$ start assets search "review"

Found 4 matches:

tasks/
  git-workflow/
    code-review            Review code for quality and best practices
    pre-commit-review      Review staged changes before committing
    doc-review             Review and improve documentation

roles/
  general/
    code-reviewer          Expert code reviewer focusing on security
```

### Search by Category

```bash
$ start assets search "git-workflow"

Found 6 matches:

tasks/
  git-workflow/
    commit-message         Generate conventional commit message
    pre-commit-review      Review staged changes before committing
    post-commit-hook       Post-commit validation workflow
    pr-ready               Complete PR preparation
    branch-cleanup         Clean up old branches
    rebase-helper          Interactive rebase assistant
```

### Search by Description

```bash
$ start assets search "security"

Found 3 matches:

tasks/
  quality/
    security-audit         Comprehensive security vulnerability scan
    dependency-check       Check dependencies for security issues

roles/
  specialized/
    security-expert        Security specialist with penetration testing expertise
```

### Search by Tag

```bash
$ start assets search "quality"

Found 5 matches:

tasks/
  git-workflow/
    pre-commit-review      Review staged changes before committing
    code-review            Review code for quality and best practices

  testing/
    test-coverage          Analyze test coverage and quality

  quality/
    lint-check             Run code quality linters
    complexity-check       Check code complexity and quality metrics
```

### Verbose Output

```bash
$ start assets search "commit" --verbose

Found 5 matches:

tasks/git-workflow/commit-message
  Description: Generate conventional commit message
  Tags: git, commit, conventional
  Category: git-workflow
  Type: tasks
  Size: 2.1 KB
  Created: 2025-01-10
  Updated: 2025-01-10

tasks/git-workflow/pre-commit-review
  Description: Review staged changes before committing
  Tags: git, review, quality, pre-commit
  Category: git-workflow
  Type: tasks
  Size: 3.4 KB
  Created: 2025-01-10
  Updated: 2025-01-12

[... additional matches ...]
```

### Exact Match

```bash
$ start assets search "pre-commit-review"

Found 1 match:

tasks/
  git-workflow/
    pre-commit-review      Review staged changes before committing
```

No installation (use `start assets add` to install).

### Pipeline Usage

```bash
$ start assets search "commit" | grep "pre-commit"

    pre-commit-review      Review staged changes before committing
```

Output is grep-friendly for scripting.

## Use Cases

### Discovery Before Installation

**Problem:** Want to see what's available before installing.

```bash
# Search for assets
start assets search "commit"

# Review results, note interesting ones
# Then install specific assets
start assets add "pre-commit-review"
```

**Workflow:** Search → review → add.

### Explore by Category

**Problem:** Want to see all assets in a category.

```bash
start assets search "git-workflow"
```

Lists all assets in the git-workflow category.

### Find by Topic

**Problem:** Looking for assets related to a specific topic.

```bash
start assets search "security"
start assets search "testing"
start assets search "documentation"
```

Searches descriptions and tags.

### Scripting and Automation

**Problem:** Need to programmatically check if asset exists.

```bash
if start assets search "my-custom-task" > /dev/null 2>&1; then
    echo "Asset found in catalog"
else
    echo "Asset not found"
fi
```

Exit codes: 0 (Found), 1 (Error), 2 (Not found).

### Quick Reference

**Problem:** Forgot exact asset name.

```bash
start assets search "pre"
```

Shows all assets matching "pre" substring.

## Comparison with Other Commands

### vs `start assets add`

**`start assets search`** - Read-only discovery

```bash
start assets search "commit"
# Lists matches, exits (no installation)
```

**`start assets add`** - Search and install

```bash
start assets add "commit"
# Lists matches, prompts for selection, installs
```

Use search for exploration, add for installation.

### vs `start assets browse`

**`start assets search`** - Keyword-based terminal search

```bash
start assets search "commit"
# Terminal output, filtered by keyword
```

**`start assets browse`** - Visual catalog browsing

```bash
start assets browse
# Opens browser, graphical navigation
```

Search is faster for targeted queries, browse is better for exploration.

### vs `start assets info`

**`start assets search`** - List multiple matches

```bash
start assets search "commit"
# Shows brief list of all matches
```

**`start assets info`** - Detailed single asset view

```bash
start assets info "pre-commit-review"
# Shows complete metadata for one asset
```

Search is for finding, info is for detailed inspection.

## Configuration

**Asset repository:**

In `~/.config/start/config.toml`:

```toml
[settings]
asset_repo = "grantcarthew/start"    # Default
# asset_repo = "myorg/custom-assets"  # Custom
```

**No other configuration needed.**

## Notes

### GitHub-Only Search

Searches **only the GitHub catalog**:

- Does NOT search local configuration
- Does NOT search global configuration
- Does NOT search cache

**For local assets:**

```bash
start config task list         # List installed tasks
start config role list         # List installed roles
```

### Substring Matching Details

**Algorithm:**

- Case-insensitive
- Minimum 3 characters
- Matches anywhere in field (not just beginning)

**Search fields (priority order):**

1. Exact name match (highest)
2. Name substring
3. Path substring
4. Description substring
5. Tag substring (lowest)

### Performance

**Typical search time:**

- Download index.csv: ~50-100ms
- Parse CSV: ~5-10ms
- In-memory search: <1ms
- **Total: ~60-110ms**

**Index unavailable (fallback):**

- Tree API call: ~100-200ms
- Parse and search: ~10-20ms
- **Total: ~110-220ms**

### Output Format

**Human-readable by default:**

- Grouped by type/category
- Indented tree structure
- Descriptive text

**Machine-parsable with care:**

- Consistent indentation (2 spaces per level)
- Predictable format
- Consider using `--verbose` for structured output

**For scripting:** Parse exit codes rather than output text.

### Network Required

Requires network access to download catalog index.

**Offline:** Cannot search (catalog unavailable).

**Cached assets:** Use `start config <type> list` to see installed assets.

### Catalog Index Dependency

**With index.csv (preferred):**

- Rich search: name, path, description, tags
- Fast and complete results

**Without index.csv (fallback):**

- Limited search: name and path only
- No description or tag matching
- Degraded but functional

## See Also

- start-assets(1) - Asset management overview
- start-assets-add(1) - Search and install asset
- start-assets-browse(1) - Visual catalog browsing
- start-assets-info(1) - Detailed asset information
- start-config-task(1) - Manage local tasks
