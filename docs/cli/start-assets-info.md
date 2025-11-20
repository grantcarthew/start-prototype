# start assets info

## Name

start assets info - Show detailed asset information

## Synopsis

```bash
start assets info <query>
```

## Description

Display detailed information about a specific asset from the GitHub catalog without installing it. Shows complete metadata including description, tags, file sizes, timestamps, and installation status.

Use this command to preview asset details before installation, or to inspect already-installed assets.

**Information displayed:**

- Asset name, type, and category
- Full description
- Tags (keywords)
- File list with sizes
- Creation and update timestamps
- Installation status (cached, in global/local config)
- SHA for version tracking
- Usage instructions

**Search behavior:**

- Uses substring matching (same as `start assets search`)
- Auto-displays if single match
- Interactive selection if multiple matches

## Arguments

**\<query\>** (required)
: Search query to find asset. Minimum 3 characters.

**Query matching:**

- Case-insensitive substring match
- Searches: name, path, description, tags
- Single match → auto-display
- Multiple matches → interactive selection

## Output

### Single Match (Auto-Display)

```bash
$ start assets info "pre-commit-review"

Asset: pre-commit-review
═══════════════════════════════════════════════════════════
Type: tasks
Category: git-workflow
Path: tasks/git-workflow/pre-commit-review

Description:
  Review staged changes before committing. Analyzes git diff
  output and provides feedback on code quality, security, and
  best practices.

Tags:
  git, review, quality, pre-commit

Files:
  pre-commit-review.toml (2.1 KB)
  pre-commit-review.md (1.3 KB)

Created: 2025-01-10T00:00:00Z
Updated: 2025-01-12T14:30:00Z
SHA: a1b2c3d4e5f6...

Installation Status:
  ✓ Cached in ~/.config/start/assets/tasks/git-workflow/
  ✓ Installed in global config
  ✗ Not in local config

Use 'start task pre-commit-review' to run.
Use 'start assets add pre-commit-review --local' to add to local config.
```

### Multiple Matches (Interactive Selection)

```bash
$ start assets info "commit"

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

Select asset [1-5] (or 'q' to quit): 2

Asset: pre-commit-review
═══════════════════════════════════════════════════════════
[... full details as shown above ...]
```

### Not Installed

```bash
$ start assets info "security-audit"

Asset: security-audit
═══════════════════════════════════════════════════════════
Type: tasks
Category: quality
Path: tasks/quality/security-audit

Description:
  Comprehensive security vulnerability scan using multiple
  tools and best practices checklists.

Tags:
  security, audit, vulnerabilities, owasp

Files:
  security-audit.toml (3.2 KB)
  security-audit.md (2.8 KB)

Created: 2025-01-08T00:00:00Z
Updated: 2025-01-08T00:00:00Z
SHA: b2c3d4e5f6a7...

Installation Status:
  ✗ Not cached
  ✗ Not installed in any config

Use 'start assets add security-audit' to install.
Browse: https://github.com/grantcarthew/start/tree/main/assets/tasks/quality/security-audit.toml
```

### Installed in Multiple Scopes

```bash
$ start assets info "go-expert"

Asset: go-expert
═══════════════════════════════════════════════════════════
Type: roles
Category: languages
Path: roles/languages/go-expert

Description:
  Go programming language expert with deep knowledge of
  concurrency patterns, standard library, and performance
  optimization.

Tags:
  go, golang, language, expert, concurrency

Files:
  go-expert.md (4.5 KB)

Created: 2025-01-05T00:00:00Z
Updated: 2025-01-10T00:00:00Z
SHA: c3d4e5f6a7b8...

Installation Status:
  ✓ Cached in ~/.config/start/assets/roles/languages/
  ✓ Installed in global config
  ✓ Installed in local config

Use 'start --role go-expert' to use this role.
```

### No Matches

```bash
$ start assets info "nonexistent"

No matches found for 'nonexistent'

Suggestions:
- Check spelling
- Try a shorter or different query
- Use 'start assets search <query>' to explore
- Use 'start assets browse' to view catalog
```

Exit code: 2

### Query Too Short

```bash
$ start assets info "ab"

Error: Query too short (minimum 3 characters)

Please provide at least 3 characters for meaningful search.
Use 'start assets browse' to explore the catalog visually.
```

Exit code: 1

### Network Error

```bash
$ start assets info "pre-commit"

Searching catalog...
✗ Network error

Cannot connect to GitHub:
  dial tcp: no route to host

Check your internet connection and try again.
```

Exit code: 1

### With Update Available

```bash
$ start assets info "code-reviewer"

Asset: code-reviewer
═══════════════════════════════════════════════════════════
Type: roles
Category: general
Path: roles/general/code-reviewer

Description:
  Expert code reviewer focusing on security, performance,
  and best practices.

Tags:
  review, security, quality, best-practices

Files:
  code-reviewer.md (3.2 KB)

Created: 2025-01-05T00:00:00Z
Updated: 2025-01-12T00:00:00Z
SHA: d4e5f6a7b8c9...

Installation Status:
  ✓ Cached (SHA: d4e5f6a7... from 2025-01-05)
  ⚠ Update available (SHA: e5f6a7b8... from 2025-01-12)
  ✓ Installed in global config

Use 'start assets update code-reviewer' to update.
Use 'start --role code-reviewer' to use this role.
```

## Exit Codes

**0** - Success (asset information displayed)

**1** - Network error, query too short, or user cancelled

**2** - Asset not found

## Examples

### Show Task Info

```bash
$ start assets info "pre-commit-review"

Asset: pre-commit-review
═══════════════════════════════════════════════════════════
Type: tasks
Category: git-workflow
Path: tasks/git-workflow/pre-commit-review

Description:
  Review staged changes before committing

Tags:
  git, review, quality, pre-commit

Files:
  pre-commit-review.toml (2.1 KB)
  pre-commit-review.md (1.3 KB)

Created: 2025-01-10
Updated: 2025-01-12
SHA: a1b2c3d4...

Installation Status:
  ✓ Cached
  ✓ Installed in global config

Use 'start task pre-commit-review' to run.
```

### Show Role Info

```bash
$ start assets info "code-reviewer"

Asset: code-reviewer
═══════════════════════════════════════════════════════════
Type: roles
Category: general
Path: roles/general/code-reviewer

Description:
  Expert code reviewer focusing on security

Tags:
  review, security, quality

Files:
  code-reviewer.md (3.2 KB)

Installation Status:
  ✓ Installed in global config

Use 'start --role code-reviewer' to use this role.
```

### Multiple Matches Selection

```bash
$ start assets info "workflow"

Found 8 matches:

tasks/
  git-workflow/
    [1] commit-message
    [2] pre-commit-review
    [3] post-commit-hook
    [4] pr-ready

  ci-workflow/
    [5] test-pipeline
    [6] deploy-check
    [7] release-prep
    [8] build-verify

Select asset [1-8] (or 'q' to quit): 4

Asset: pr-ready
═══════════════════════════════════════════════════════════
[... full details ...]
```

### User Cancellation

```bash
$ start assets info "commit"

Found 5 matches:
[... list ...]

Select asset [1-5] (or 'q' to quit): q

Cancelled.
```

Exit code: 1

## Use Cases

### Preview Before Installation

**Problem:** Want to see details before installing.

```bash
# Preview asset
start assets info "security-audit"

# Review description, tags, file sizes

# Install if suitable
start assets add "security-audit"
```

### Check Installation Status

**Problem:** Forgot if asset is installed.

```bash
start assets info "pre-commit-review"
```

Shows installation status (cached, global, local).

### Verify Asset Updates

**Problem:** Want to know if installed asset has updates.

```bash
start assets info "code-reviewer"
```

Shows update availability with version comparison.

### Explore Asset Details

**Problem:** Want complete information about an asset.

```bash
start assets info "go-expert"
```

Displays comprehensive metadata and usage instructions.

## Comparison with Other Commands

### vs `start assets search`

**`start assets search`** - List multiple matches briefly

```bash
start assets search "commit"
# Shows list with short descriptions
```

**`start assets info`** - Detailed single asset view

```bash
start assets info "commit"
# Shows complete metadata for selected asset
```

Search for discovery, info for inspection.

### vs `start assets add`

**`start assets info`** - View details only (read-only)

```bash
start assets info "pre-commit"
# Displays information, no installation
```

**`start assets add`** - View and install

```bash
start assets add "pre-commit"
# Displays info, prompts for installation
```

Info is non-invasive inspection.

### vs `start show task`

**`start assets info`** - Catalog asset metadata

```bash
start assets info "pre-commit-review"
# Shows catalog metadata (GitHub)
```

**`start show task`** - Resolved configuration

```bash
start show task pre-commit-review
# Shows effective configuration after UTD processing
```

Info shows catalog source, show displays runtime config.

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

### Installation Status Detection

**Status checks:**

1. **Cached** - File exists in `~/.config/start/assets/{type}/{category}/`
2. **Global config** - Entry exists in `~/.config/start/{type}.toml`
3. **Local config** - Entry exists in `./.start/{type}.toml`

**Independent checks:**

- Asset can be cached but not configured
- Asset can be in global but not local (or vice versa)
- Asset can be in both global and local

### Update Detection

**SHA comparison:**

- Cached asset SHA (from `.meta.toml`)
- Catalog asset SHA (from `index.csv`)
- Different → update available
- Same → up to date

### GitHub-Only Source

Searches **only the GitHub catalog** for asset metadata.

**For installed assets:** Shows installation status by checking local filesystem.

### Network Required

Requires network to download catalog index and metadata.

**Offline:** Cannot fetch catalog data (installation status still shown from local checks).

### Substring Matching

Uses same algorithm as `start assets search` and `start assets add`.

**Minimum length:** 3 characters

## See Also

- start-assets(1) - Asset management overview
- start-assets-search(1) - Search catalog
- start-assets-add(1) - Install asset
- start-assets-update(1) - Update cached assets
- start-show(1) - Display resolved configuration
