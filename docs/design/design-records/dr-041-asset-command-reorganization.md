# DR-041: Asset Command Reorganization

- Date: 2025-01-13
- Status: Accepted
- Category: CLI Design

## Problem

Asset management functionality is scattered across multiple commands with inconsistent interfaces:

Old structure:
```bash
# Adding assets - different paths per type
start config task add [query]    # Add task from catalog
start config role add [query]    # Add role from catalog
start config agent add [query]   # Add agent from catalog

# Updating assets - separate command
start update                     # Update all cached assets

# Showing assets - buried in show command
start show assets                # List cached assets
```

Issues:

- Discovery is hard (no search or browse functionality)
- Inconsistent UX (each asset type uses different command pattern)
- No asset information (can't preview before installing)
- Update is unclear (start update name conflicts with "update CLI binary")
- Type coupling (must know asset type before discovering)

## Decision

Create unified start assets command suite for all asset discovery, installation, and management operations.

**Note:** `start assets clean` was considered but rejected. Assets cached for local projects cannot be safely identified from a global context, creating a risk of breaking local configurations by deleting their dependencies.

Command structure:

```bash
# Discovery and search
start assets browse              # Open GitHub catalog in web browser
start assets search <query>      # Search by name/description/tags
start assets info <query>        # Show detailed asset information

# Installation
start assets add <query>         # Search and install asset
start assets add                 # Interactive terminal browser (no args)

# Management
start assets update [query]      # Update cached assets (all or matching query)

# Maintenance (for catalog contributors)
start assets index               # Generate catalog index.csv
```

Key principle: Asset-first design - discover assets, then install to appropriate config.

Replaces scattered commands:

| Old Command | Replacement |
|-------------|-------------|
| start config task add | start assets add |
| start config role add | start assets add |
| start config agent add | start assets add |
| start update | start assets update |
| start show assets | start assets browse |

## Why

Unified discovery simplifies asset management:

- One place for all asset operations (not scattered across config/update/show)
- Searchable (find assets by description, tags, not just name)
- Preview before install (see details with start assets info)
- Type-agnostic (don't need to know if something is a task/role/agent)
- Clear organization (logical grouping under assets command)

Consistent implementation improves maintainability:

- Same code paths for all asset types (universal installer)
- Easier to extend (add new asset types without new commands)
- Better testing (unified test suite for asset operations)
- Shared functionality (search, download, cache all reused)

Clearer discoverability for users:

- Help is clearer (start assets --help shows all operations)
- Intuitive naming ("assets" clearly indicates catalog operations)
- Less cognitive load (one command to remember, not four)
- Natural grouping (related operations together)

Clearer update semantics:

- Explicit scope (start assets update clearly updates assets)
- No ambiguity (not "update CLI binary")
- Selective updates (update matching query only)
- Predictable behavior (update asset cache, not config)

## Trade-offs

Accept:

- Breaking change (requires migration for existing scripts/workflows, but clear replacement mapping provided)
- Command is longer (start assets add vs start config task add, but supports search shortcuts with prefix matching)
- Learning curve for existing users (must learn new command structure, but clearer long-term)
- Indefinite cache growth (no clean command, but assets are small text files)

Gain:

- Unified discovery (one place for all asset operations, searchable, preview before install, type-agnostic, clear organization)
- Consistent implementation (same code paths for all asset types, easier to extend, better testing, shared functionality)
- Clearer discoverability (help is clearer, intuitive naming, less cognitive load, natural grouping)
- Clearer update semantics (explicit scope, no ambiguity, selective updates, predictable behavior)

## Alternatives

Keep scattered commands:

Example: Maintain current structure with separate commands per type
```bash
start config task add "foo"
start config role add "bar"
start config agent add "baz"
start update
start show assets
```

Pros:
- No breaking changes (existing users unaffected)
- No migration needed (current scripts work)
- Shorter commands (config task add vs assets add)
- Familiar pattern (users already know it)

Cons:
- No unified discovery (must know asset type first)
- No search functionality (can't find by description/tags)
- No preview capability (can't see details before installing)
- Inconsistent UX (different patterns per type)
- Harder to extend (new asset type needs new command)
- Confusing update (start update ambiguous)
- Poor discoverability (commands scattered across CLI)

Rejected: Unified command suite provides much better UX. Breaking change is acceptable with clear migration mapping.

Use subcommands under config instead of new assets command:

Example: Keep under config but unify pattern
```bash
start config assets browse
start config assets search <query>
start config assets add <query>
start config assets update
```

Pros:
- Less top-level command pollution (nested under config)
- Still unifies asset operations
- Familiar location (near existing config commands)

Cons:
- Confusing semantics (assets command modifies cache, not config)
- Longer commands (start config assets add vs start assets add)
- Misleading grouping (config is about user settings, assets about catalog)
- Harder to discover (buried under config)
- Semantic mismatch (browsing catalog not configuring)

Rejected: start assets is clearer - assets are distinct from user config. Cache management and catalog browsing don't belong under config.

Add browse/search without deprecating old commands:

Example: Add new discovery commands but keep old add commands
```bash
start assets browse              # New
start assets search <query>      # New
start assets info <query>        # New
start config task add <query>    # Keep
start config role add <query>    # Keep
start update                     # Keep
```

Pros:
- No breaking changes (old commands still work)
- Adds discovery features (browse, search, info)
- Gradual migration (users can switch over time)

Cons:
- Duplicated functionality (two ways to add assets)
- Inconsistent UX (new way vs old way coexist)
- Confusing for new users (which command to use?)
- Maintenance burden (support both paths indefinitely)
- No incentive to migrate (old commands never go away)
- Cluttered help output (too many similar commands)

Rejected: Clean break better. Duplicated functionality creates confusion.

Use single command with type flag:

Example: Single add command with --type flag
```bash
start assets add "foo" --type task
start assets add "bar" --type role
start assets add "baz" --type agent
```

Pros:
- Single command (one add instead of multiple)
- Explicit type (clear what's being added)
- Consistent pattern (same command for all types)

Cons:
- Requires knowing type (defeats type-agnostic discovery)
- Extra typing (--type flag every time)
- Less ergonomic (more verbose)
- Type is implementation detail (users shouldn't need to specify)
- Search can infer type (from catalog path)

Rejected: Type-agnostic discovery is key feature. Asset type is metadata, not user concern. Let search/browse determine type automatically.

## Structure

Command structure:

start assets browse:
- Open GitHub catalog in web browser
- URL: https://github.com/{org}/{repo}/tree/main/assets
- Uses system default browser
- Fallback: If browser fails to open, print URL to terminal for manual click
- No terminal interaction (pure web browsing)

start assets search <query>:
- Search catalog using substring matching (DR-040)
- Match against name, path, description, tags
- Display results grouped by type/category
- Interactive selection if multiple matches
- Auto-select if single exact match
- Minimum 3 characters required

start assets add <query>:
- Search catalog using substring matching
- Single match: auto-select and proceed
- Multiple matches: interactive selection
- Show asset details and confirm installation
- Download to cache (~/.config/start/assets/)
- Add to global or local config (via --local flag)
- Display usage instructions after installation

start assets add (no arguments):
- Interactive terminal browser with numbered selection (DR-035)
- Download catalog index
- Category-first navigation
- Select asset from numbered list
- Confirm before installing
- Download and add to config

start assets info <query>:
- Search for asset (same matching as add)
- Display full metadata from .meta.toml
- Show description and tags
- Show file size and dates
- Show installation status (cached, in global/local config)
- Optional full file contents preview

start assets update [query]:
- Download catalog index.csv
- Find cached .meta.toml files
- Filter by query if provided (substring matching)
- Compare local SHA with index SHA
- Download updates for changed assets
- Update cache only (never modify user config)
- Report updated count and unchanged count
- See DR-037 for complete algorithm

start assets index:
- Validate git repository structure
- Scan assets/ directory for .meta.toml files
- Parse metadata from each file
- For agents: extract bin from .toml file
- Sort alphabetically (type → category → name)
- Write CSV to assets/index.csv with header row
- Report generated asset count
- See DR-039 for complete specification

Scope control:

Downloaded assets added to config based on --local flag:
- Default: global config (~/.config/start/tasks.toml)
- With --local flag: local config (.start/tasks.toml)
- Local creates or updates the .start/ directory

Cache behavior:

Cached assets used immediately without prompting:
- Cache is transparent implementation detail
- No user interaction needed for cached assets
- Reduces network calls automatically

## Usage Examples

Discovery workflow:

```bash
# Browse in web browser
start assets browse
# Opens https://github.com/grantcarthew/start/tree/main/assets

# Or search by description in terminal
start assets search "commit"

# Or use interactive terminal browser
start assets add
# Shows numbered category and asset selection

# Preview before installing
start assets info "pre-commit-review"

# Install when ready
start assets add "pre-commit-review"
```

Update workflow:

```bash
start assets update       # Clear: updates cached assets
start assets update "git" # Update only git-related assets
```

Open web browser to catalog (success):

```bash
$ start assets browse

Opening GitHub catalog in browser...
✓ Opened https://github.com/grantcarthew/start/tree/main/assets
```

Open web browser to catalog (fallback when browser fails):

```bash
$ start assets browse

Opening GitHub catalog in browser...
⚠ Could not open browser automatically

Visit the catalog at:
https://github.com/grantcarthew/start/tree/main/assets
```

Interactive terminal browser:

```bash
$ start assets add

Fetching catalog from GitHub...
✓ Found 46 assets across 4 types and 12 categories

Select category:
  1. git-workflow (4 tasks)
  2. code-quality (4 tasks)
  3. security (2 tasks)
  4. debugging (2 tasks)
  5. [view all tasks]

> 1

git-workflow tasks:
  1. pre-commit-review - Review staged changes before commit
  2. pr-ready - Complete PR preparation checklist
  3. commit-message - Generate conventional commit message
  4. explain-changes - Understand what changed in commits

> 1

Selected: pre-commit-review
Description: Review staged changes before commit
Tags: git, review, quality, pre-commit

Download and add to config? [Y/n] y

Downloading...
✓ Cached to ~/.config/start/assets/tasks/git-workflow/
✓ Added to global config as 'pre-commit-review'

Try it: start task pre-commit-review
```

Search and install:

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

Select asset [1-5] (or 'q' to quit): 2

Found 1 match:
  tasks/git-workflow/pre-commit-review

Downloading...
✓ Cached to ~/.config/start/assets/tasks/git-workflow/
✓ Added to global config as 'pre-commit-review'

Use 'start task pre-commit-review' to run.
```

Exact match (auto-select):

```bash
$ start assets add "pre-commit-review"

Found 1 match (exact):
  tasks/git-workflow/pre-commit-review

Downloading...
✓ Cached to ~/.config/start/assets/tasks/git-workflow/
✓ Added to global config as 'pre-commit-review'

Use 'start task pre-commit-review' to run.
```

Add to local config:

```bash
$ start assets add "code-reviewer" --local

Found 1 match:
  roles/general/code-reviewer

Downloading...
✓ Cached to ~/.config/start/assets/roles/general/
✓ Added to local config (./.start/roles.toml)

Use 'start --role code-reviewer' to use this role.
```

Asset info:

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
  pre-commit-review.toml (1.8 KB)
  pre-commit-review.md (512 bytes)

Created: 2025-01-10T00:00:00Z
Updated: 2025-01-12T00:00:00Z

Status:
  ✓ Cached in ~/.config/start/assets/tasks/git-workflow/
  ✓ Installed in global config
  ✗ Not in local config

Use 'start assets add pre-commit-review --local' to add to local config.
Use 'start task pre-commit-review' to run.
```

Update all cached assets:

```bash
$ start assets update

Checking for updates...

✓ tasks/git-workflow/pre-commit-review  (a1b2c3 → d4e5f6)
✓ roles/general/code-reviewer           (up to date)
✓ tasks/quality/commit-lint             (7890ab → cdef12)

Updated 2 assets, 1 up to date.
```

Update matching assets:

```bash
$ start assets update "commit"

Checking for updates to assets matching 'commit'...

✓ tasks/git-workflow/commit-message     (123456 → 7890ab)
✓ tasks/git-workflow/pre-commit-review  (up to date)
✓ tasks/git-workflow/post-commit-hook   (abcdef → 123456)

Updated 2 assets, 1 up to date.
```

Generate catalog index:

```bash
$ cd ~/projects/start-catalog-fork
$ start assets index

Validating repository structure...
✓ Git repository detected
✓ Assets directory found

Scanning assets/...
Found 46 assets

Sorting assets (type → category → name)...
Writing index to assets/index.csv...

✓ Generated index with 46 assets
Updated: assets/index.csv

Ready to commit:
  git add assets/
  git commit -m "Regenerate catalog index"
```

## Updates

- 2025-01-17: Initial version aligned with schema; removed implementation code; corrected start assets browse to open web browser (not terminal TUI); added URL fallback when browser fails to open; removed 'start assets clean' to prevent accidental deletion of assets used by local configurations
