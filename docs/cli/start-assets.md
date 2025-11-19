# start assets

## Name

start assets - Discover and manage catalog assets

## Synopsis

```bash
start assets browse            # Open GitHub catalog in browser
start assets search <query>    # Search catalog by name/description/tags
start assets add <query>       # Search and install asset
start assets info <query>      # Show detailed asset information
start assets update [query]    # Update cached assets
start assets index             # Generate catalog index (contributors)
```

## Description

Unified command suite for discovering, installing, and managing assets from the GitHub catalog. Provides semantic separation from `start config` commands: `start assets` is for shopping the catalog, while `start config` is for managing your configuration.

**Asset types:**
- **tasks** - Predefined AI workflow tasks
- **roles** - System prompts and personas
- **agents** - AI agent configurations
- **contexts** - Context document templates

**Core principles:**
- `start assets` commands **only search GitHub catalog** (not local/global/cache)
- Discovery-first design (browse, search, preview, then install)
- Type-agnostic operations (one command for all asset types)
- Cache transparency (automatic caching, invisible to user)

See individual subcommand documentation for detailed usage.

## Subcommands

### start assets browse

Open GitHub catalog in default web browser for visual exploration.

```bash
start assets browse
```

Opens: `https://github.com/{asset_repo}/tree/main/assets`

Uses `[settings] asset_repo` value (default: `grantcarthew/start`).

See [start-assets-browse(1)](./start-assets-browse.md) for details.

### start assets search

Search catalog by substring matching (name, path, description, tags). Terminal output only, non-interactive.

```bash
start assets search <query>    # Minimum 3 characters
```

**Examples:**
```bash
start assets search "commit"
start assets search "security"
start assets search git
```

See [start-assets-search(1)](./start-assets-search.md) for details.

### start assets add

Search and install asset from catalog. Interactive when no query provided, or with query using substring matching.

```bash
start assets add               # Interactive TUI browser
start assets add <query>       # Search and install
```

**Examples:**
```bash
start assets add                              # Browse all assets
start assets add "commit"                     # Search for 'commit'
start assets add git-workflow/pre-commit-review  # Direct install
```

See [start-assets-add(1)](./start-assets-add.md) for details.

### start assets info

Show detailed metadata for specific asset without installing. Searches catalog by substring matching.

```bash
start assets info <query>
```

**Examples:**
```bash
start assets info "pre-commit-review"
start assets info "code-reviewer"
```

See [start-assets-info(1)](./start-assets-info.md) for details.

### start assets update

Update cached assets by comparing SHAs with GitHub catalog.

```bash
start assets update            # Update all cached assets
start assets update <query>    # Update matching assets
```

**Examples:**
```bash
start assets update                   # Update all
start assets update "commit"          # Update matching 'commit'
start assets update git-workflow      # Update category
```

See [start-assets-update(1)](./start-assets-update.md) for details.

### start assets index

Generate `assets/index.csv` for catalog contributors. Must be run in catalog repository with `.git/` and `assets/` directories.

```bash
start assets index
```

**For catalog maintainers only.** See [start-assets-index(1)](./start-assets-index.md) for details.

## Common Workflows

### Discovery Workflow

**Browse visually in GitHub:**
```bash
start assets browse
```

**Search by keyword:**
```bash
start assets search "commit"
start assets search "security review"
```

**Preview before installing:**
```bash
start assets info "pre-commit-review"
```

### Installation Workflow

**Interactive browsing:**
```bash
start assets add
# Navigate categories, select asset
```

**Search and install:**
```bash
start assets add "commit"
# Select from matching results
```

**Direct install:**
```bash
start assets add git-workflow/pre-commit-review
# Installs immediately
```

### Maintenance Workflow

**Update all cached assets:**
```bash
start assets update
```

**Update specific assets:**
```bash
start assets update git-workflow
```

## Design References

- [DR-039](../design/design-records/dr-039-catalog-index.md) - Catalog index file (CSV schema)
- [DR-040](../design/design-records/dr-040-substring-matching.md) - Substring matching algorithm
- [DR-041](../design/design-records/dr-041-asset-command-reorganization.md) - Asset command reorganization

## Examples

### Browse GitHub Catalog

```bash
$ start assets browse

Opening GitHub catalog in browser...
✓ https://github.com/grantcarthew/start/tree/main/assets
```

### Search for Assets

```bash
$ start assets search "commit"

Found 3 matches:

tasks/git-workflow/commit-message
  Description: Generate conventional commit message
  Tags: git, commit, conventional

tasks/git-workflow/pre-commit-review
  Description: Review staged changes before committing
  Tags: git, review, quality, pre-commit

tasks/git-workflow/post-commit-hook
  Description: Post-commit validation workflow
  Tags: git, commit, hooks, validation
```

### Add Asset Interactively

```bash
$ start assets add

Fetching catalog from GitHub...
✓ Found 46 assets across 4 types

Select category:
  1. git-workflow (12 tasks)
  2. quality (8 tasks)
  3. security (6 tasks)
  4. debugging (2 tasks)

> 1

git-workflow tasks:
  1. commit-message - Generate conventional commit message
  2. pre-commit-review - Review staged changes
  3. pr-ready - Complete PR preparation
  4. post-commit-hook - Post-commit validation

> 2

Selected: pre-commit-review
Description: Review staged changes before committing
Tags: git, review, quality, pre-commit

Download and add to config? [Y/n] y

Downloading...
✓ Cached to ~/.config/start/assets/tasks/git-workflow/
✓ Added to global config as 'pre-commit-review'

Try it: start task pre-commit-review
```

### Add Asset by Query

```bash
$ start assets add "pre-commit"

Found 1 match (exact):
  tasks/git-workflow/pre-commit-review

Downloading...
✓ Cached to ~/.config/start/assets/tasks/git-workflow/
✓ Added to global config

Use 'start task pre-commit-review' to run.
```

### Show Asset Info

```bash
$ start assets info "code-reviewer"

Asset: code-reviewer
═══════════════════════════════════════
Type: roles
Category: general
Path: roles/general/code-reviewer

Description:
  Expert code reviewer focusing on security

Tags:
  review, security, quality

Files:
  code-reviewer.md (847 bytes)

Created: 2025-01-10T00:00:00Z
Updated: 2025-01-12T00:00:00Z

Status:
  ✓ Installed in global config
  ✓ Cached locally
  ✗ Not in local config

Use 'start assets add code-reviewer --local' to add to local config.
```

### Update Assets

```bash
$ start assets update

Checking for updates...

✓ tasks/git-workflow/pre-commit-review  (v1.0 → v1.1)
✓ roles/general/code-reviewer           (up to date)
✓ tasks/quality/commit-lint             (v2.3 → v2.4)

Updated 2 assets, 1 up to date.
```

### Update Specific Assets

```bash
$ start assets update "commit"

Checking for updates to assets matching 'commit'...

✓ tasks/git-workflow/commit-message     (v1.2 → v1.3)
✓ tasks/git-workflow/pre-commit-review  (up to date)
✓ tasks/git-workflow/post-commit-hook   (v1.0 → v1.1)

Updated 2 assets, 1 up to date.
```

## Exit Codes

**0** - Success

**1** - Network error (catalog unavailable, download failed)

**2** - Asset not found

**3** - User cancelled or file system error

## Notes

### GitHub-Only Search

`start assets` commands only search the GitHub catalog, not your local/global configs or cache.

**Rationale:**
- Discovery is about exploring what's **available** in the catalog
- Local/global configs are already known to you (`start config <type> list`)
- Cache is a subset of GitHub catalog
- Searching GitHub provides complete, fresh view

**For local assets:**
```bash
start config task list         # See your configured tasks
start config role list         # See your configured roles
```

### Cache Transparency

Cache is automatically managed:
- Populated when downloading from GitHub
- Updated via `start assets update`
- No manual inspection needed

Cache location: `~/.config/start/assets/`

### Installation Scope

Assets install to **global config** by default:

```bash
start assets add "task-name"           # → ~/.config/start/tasks.toml
start assets add "task-name" --local   # → ./.start/tasks.toml
```

## See Also

- start-assets-browse(1) - Open catalog in browser
- start-assets-search(1) - Search catalog
- start-assets-add(1) - Add asset from catalog
- start-assets-info(1) - Show asset information
- start-assets-update(1) - Update cached assets
- start-config(1) - Manage configuration
- start-task(1) - Run tasks
