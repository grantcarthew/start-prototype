# start assets add

## Name

start assets add - Search and install asset from GitHub catalog

## Synopsis

```bash
start assets add [flags]                       # Interactive browsing
start assets add <query> [flags]               # Search and install
start assets add <path> [flags]                # Direct install by path
```

## Description

Search for assets in the GitHub catalog and install them to your configuration. Supports two modes: interactive browsing (no query) or search-based selection (with query). Downloaded assets are cached locally and added to global or local configuration.

> **Tip: Lazy Loading**
> You don't always need to run `start assets add` explicitly. If you try to run a task that isn't installed (e.g., `start task pre-commit-review`), the CLI will automatically search the catalog and prompt you to download it. Use `start assets add` when you want to browse options or install to a specific scope (local/global) ahead of time.

**Two operational modes:**

**Interactive mode** (no query or < 3 characters):

- Step-by-step numbered selection
- Browse by Asset Type → Category → Asset
- Preview metadata before downloading

**Search mode** (3+ character query):

- Substring search across name, path, description, and tags
- Auto-select if single match found
- Interactive selection if multiple matches
- Download and install selected asset

**Installation process:**

1. Search GitHub catalog for matching assets
2. Download all asset files (content + metadata) to cache (`~/.config/start/assets/`)
3. Add configuration entry to the appropriate file (e.g., `tasks.toml`) in `~/.config/start/` or `./.start/`
4. Report installation location and usage instructions

## Arguments

**\<query\>** (optional)
: Search query for finding assets.

**Query behavior:**

- **Omitted or < 3 chars** - Interactive mode (browse all assets)
- **3+ characters** - Search mode (substring matching)

**Query matching:** Case-insensitive substring match against:

- Asset name (highest priority)
- Full path (type/category/name)
- Description
- Tags

## Flags

**-l, --local**
: Install to local project configuration (`./.start/`) instead of global (`~/.config/start/`).

**Examples:**

```bash
start assets add "pre-commit" --local    # Install to ./.start/
start assets add "pre-commit" -l         # Short flag
start assets add "go-expert"             # Install to ~/.config/start/
```

## Behavior

### Search Mode (Query Provided)

**Single match (exact):**

```
1. Search catalog
2. Find single exact match
3. Auto-select
4. Download to cache
5. Add to config
6. Show usage instructions
```

**Multiple matches:**

```
1. Search catalog
2. Display matches grouped by type/category
3. User selects from numbered list
4. Download selected asset
5. Add to config
6. Show usage instructions
```

**No matches:**

```
1. Search catalog
2. No results found
3. Display error with suggestions
4. Exit code 2
```

**Query too short (<3 chars):**

```
1. Detect short query
2. Fall back to interactive mode
3. Display warning
4. Proceed with interactive browsing
```

### Interactive Mode (No Query)

When no query is provided (or query is < 3 chars), `start assets add` launches a numbered selection interface for browsing the catalog.

**Features:**

- **Categorized View**: Browse assets grouped by Type and Category.
- **Details**: View asset descriptions and tags in the list.
- **Confirmation**: Review metadata before downloading.

**Interaction:**

- **Numbered Selection**: Type the number of your choice and press Enter.
- **Navigation**: Select Type → Category to drill down, or `[view all]` to see everything.
- **Cancellation**: Type `q` to quit at any time.

### Installation Locations

**Global installation (default):**

- **Cache:** `~/.config/start/assets/{type}/{category}/{name}.*`
- **Config:** `~/.config/start/{type}.toml`

**Local installation (--local flag):**

- **Cache:** `~/.config/start/assets/{type}/{category}/{name}.*` (always global/shared)
- **Config:** `./.start/{type}.toml`

**Agent Assets Note:**
Installing an agent configuration (e.g., `claude-3-5-sonnet`) downloads the configuration file but **does not install the AI tool binary** (e.g., `claude` CLI). You must ensure the required binary is installed and discoverable on your system.

**Note:** Asset content files (Markdown, templates) are **always downloaded to the global cache**, even for local installations. This prevents duplication and saves disk space. Only the configuration definition is added to the local project scope.

**Example:**

```bash
# Global Role
start assets add "code-reviewer"
# → Cache: ~/.config/start/assets/roles/general/code-reviewer.md
# → Config: ~/.config/start/roles.toml

# Global Agent
start assets add "claude-3-5-sonnet"
# → Cache: ~/.config/start/assets/agents/anthropic/claude-3-5-sonnet.toml
# → Config: ~/.config/start/agents.toml

# Global Context
start assets add "docker-context"
# → Cache: ~/.config/start/assets/contexts/devops/docker-context.toml
# → Config: ~/.config/start/contexts.toml

# Local Task
start assets add "pre-commit" --local
# → Cache: ~/.config/start/assets/tasks/git-workflow/pre-commit-review.toml
# → Config: ./.start/tasks.toml
```

### Multi-File Assets

Assets may consist of multiple files:

```
tasks/git-workflow/pre-commit-review.toml       # Task definition
tasks/git-workflow/pre-commit-review.md         # Prompt or instructions
tasks/git-workflow/pre-commit-review.meta.toml  # Metadata
```

**All files downloaded** (including metadata) when installing an asset. This ensures full offline capability and metadata availability.

## Output

### Single Match (Auto-Select)

```bash
$ start assets add "pre-commit-review"

Searching catalog...
✓ Found 1 match (exact)

tasks/git-workflow/pre-commit-review
  Description: Review staged changes before committing
  Tags: git, review, quality, pre-commit

Downloading...
✓ Cached to ~/.config/start/assets/tasks/git-workflow/
✓ Added to global config (~/.config/start/tasks.toml)

Use 'start task pre-commit-review' to run.
```

### Multiple Matches (Interactive Selection)

```bash
$ start assets add "commit"

Searching catalog...
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

Selected: tasks/git-workflow/pre-commit-review

Downloading...
✓ Cached to ~/.config/start/assets/tasks/git-workflow/
✓ Added to global config (~/.config/start/tasks.toml)

Use 'start task pre-commit-review' to run.
```

### No Matches

```bash
$ start assets add "nonexistent"

Searching catalog...
No matches found for 'nonexistent'

Suggestions:
- Check spelling
- Try a shorter or different query
- Use 'start assets browse' to explore the catalog
- Visit: https://github.com/grantcarthew/start/tree/main/assets
```

Exit code: 2

### Query Too Short (Fallback)

```bash
$ start assets add "ab"

⚠ Query too short (minimum 3 characters)
  Falling back to interactive mode...

[ Launches TUI Browser ]
```

### With --local Flag

```bash
$ start assets add "code-reviewer" --local

Searching catalog...
✓ Found 1 match

roles/general/code-reviewer
  Description: Expert code reviewer focusing on security

Downloading...
✓ Cached to ~/.config/start/assets/roles/general/
✓ Added to local config (./.start/roles.toml)

Use 'start --role code-reviewer' to use this role.
```

### Already Installed

```bash
$ start assets add "pre-commit-review"

Searching catalog...
✓ Found 1 match

tasks/git-workflow/pre-commit-review
  ✓ Already installed in global config

Options:
  1) Reinstall (update cache and config)
  2) Add to local config (keep global)
  3) Cancel

Select [1-3]: 1

Downloading...
✓ Updated cache
✓ Config unchanged (already present)

Use 'start task pre-commit-review' to run.
```

### Network Error

```bash
$ start assets add "pre-commit"

Searching catalog...
✗ Network error

Cannot connect to GitHub:
  dial tcp: no route to host

Check your internet connection and try again.
```

Exit code: 1

## Exit Codes

**0** - Success (asset installed)

**1** - Network error, catalog unavailable, or user cancelled

**2** - Asset not found

**3** - File system error (cache write failed, config write failed)

## Examples

### Install Task to Global Config

```bash
$ start assets add "pre-commit-review"

Searching catalog...
✓ Found 1 match (exact)

tasks/git-workflow/pre-commit-review
  Description: Review staged changes before committing

Downloading...
✓ Cached to ~/.config/start/assets/tasks/git-workflow/
✓ Added to global config (tasks.toml)

Use 'start task pre-commit-review' to run.
```

### Install Role to Local Config

```bash
$ start assets add "go-expert" --local

Searching catalog...
✓ Found 1 match

roles/languages/go-expert
  Description: Go programming language expert

Downloading...
✓ Cached to ~/.config/start/assets/roles/languages/
✓ Added to local config (./.start/roles.toml)

Use 'start --role go-expert' for Go development.
```

### Search by Partial Name

```bash
$ start assets add "commit"

Searching catalog...
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

### Search by Category

```bash
$ start assets add "git-workflow"

Searching catalog...
Found 6 matches in git-workflow:

tasks/git-workflow/
  [1] commit-message         Generate conventional commit message
  [2] pre-commit-review      Review staged changes before committing
  [3] post-commit-hook       Post-commit validation workflow
  [4] pr-ready               Complete PR preparation
  [5] branch-cleanup         Clean up old branches
  [6] rebase-helper          Interactive rebase assistant

Select asset [1-6] (or 'q' to quit): _
```

### Search by Description

```bash
$ start assets add "security"

Searching catalog...
Found 3 matches:

tasks/
  quality/
    [1] security-audit         Comprehensive security vulnerability scan
    [2] dependency-check       Check dependencies for security issues

roles/
  specialized/
    [3] security-expert        Security specialist with penetration testing expertise

Select asset [1-3] (or 'q' to quit): _
```

### Search by Tag

```bash
$ start assets add "quality"

Searching catalog...
Found 5 matches:

tasks/
  git-workflow/
    [1] pre-commit-review      Tags: git, review, quality, pre-commit
    [2] code-review            Tags: review, quality, best-practices

  testing/
    [3] test-coverage          Tags: testing, quality, coverage

  quality/
    [4] lint-check             Tags: quality, linting, standards
    [5] complexity-check       Tags: quality, metrics, complexity

Select asset [1-5] (or 'q' to quit): _
```

### Direct Path Install

```bash
$ start assets add "tasks/git-workflow/pre-commit-review"

Searching catalog...
✓ Found 1 match (exact)

tasks/git-workflow/pre-commit-review

Downloading...
✓ Cached to ~/.config/start/assets/tasks/git-workflow/
✓ Added to global config (tasks.toml)

Use 'start task pre-commit-review' to run.
```

### User Cancellation

```bash
$ start assets add "commit"

Searching catalog...
Found 5 matches:

tasks/
  git-workflow/
    [1] commit-message
    [2] pre-commit-review
    [3] post-commit-hook
  quality/
    [4] commit-lint
roles/
  git/
    [5] commit-specialist

Select asset [1-5] (or 'q' to quit): q

Cancelled.
```

Exit code: 1

## Use Cases

### Quick Installation

**Problem:** Know exact asset name, want to install quickly.

```bash
start assets add "pre-commit-review"
```

Auto-selects and installs immediately.

**Tip:** You can also just run `start task pre-commit-review`. If it's not installed, `start` will offer to download it automatically.

### Exploration and Discovery

**Problem:** Not sure what's available, want to search by topic.

```bash
# Search by topic
start assets add "security"

# Browse results, select one
```

Discovers assets by description/tags.

### Team Standardization

**Problem:** Want project-specific assets for team consistency.

```bash
# Install to local config (can be committed)
start assets add "go-expert" --local
start assets add "pre-commit-review" --local

# Team members get same assets via git
```

Local config can be committed to version control.

**Note:** When team members check out the repo and run a task, `start` will automatically detect any missing asset files (tasks, roles) referenced by the config and restore them from the catalog (Asset Restoration).

### Multi-Project Workflow

**Problem:** Use same assets across many projects.

```bash
# Install commonly used assets to global config
start assets add "code-reviewer"
start assets add "commit-message"

# Available in all projects
```

Global config provides baseline across all projects.

### Category-Based Discovery

**Problem:** Want to see all assets in a category.

```bash
start assets add "git-workflow"
```

Matches all assets in git-workflow category.

## Comparison with Other Commands

### vs `start assets browse`

**`start assets browse`** - Visual catalog exploration (browser)

```bash
start assets browse
# Opens GitHub in browser, no installation
```

**`start assets add`** - Search and install (terminal)

```bash
start assets add "commit"
# Finds, selects, downloads, and installs
```

Workflow: Browse visually, then return to terminal to install.

### vs `start assets search`

**`start assets search`** - Find assets, display only

```bash
start assets search "commit"
# Lists matches, exits
```

**`start assets add`** - Find and install

```bash
start assets add "commit"
# Lists matches, prompts for selection, installs
```

Search is read-only, add performs installation.

## Configuration

**Asset repository:**

In `~/.config/start/config.toml`:

```toml
[settings]
asset_repo = "grantcarthew/start"    # Default
# asset_repo = "myorg/custom-assets"  # Custom
```

**Note:** This setting configures the _source_ repository. Downloaded assets are installed into `tasks.toml`, `roles.toml`, `agents.toml`, or `contexts.toml`, not into `config.toml`.

## Notes

### GitHub-Only Source

`start assets add` searches **only the GitHub catalog**:

- Does NOT search local configuration
- Does NOT search global configuration
- Does NOT search cache

**Rationale:** This is for discovering and adding new assets from the catalog.

**For local assets:**

```bash
start config task list         # See installed tasks
start config role list         # See installed roles
```

### Cache Location

All downloaded assets cached to:

```
~/.config/start/assets/{type}/{category}/{name}.*
```

**Cache is shared** between global and local configs:

- Download once, reference from multiple configs
- `--local` flag affects config location, not cache

### Network Required

This command requires network access to:

- Download catalog index (`assets/index.csv`)
- Download asset files

**Offline:** Cannot add new assets (catalog unavailable).

**Cached assets:** Use `start config task new` to create custom local assets.

### Substring Matching Details

**Minimum length:** 3 characters

**Match priority:**

1. Exact name match (auto-select if single)
2. Name substring
3. Path substring (category matching)
4. Description substring
5. Tag substring

**Case-insensitive:** "commit" matches "Commit", "COMMIT", etc.

### Installation Scope Best Practices

**Global (`~/.config/start/`):**

- Personal preferences
- Commonly used across projects
- Not committed to version control

**Local (`./.start/`):**

- Project-specific assets
- Team standardization
- Committed to version control
- Shared with collaborators

**Example workflow:**

```bash
# Personal preferences (global)
start assets add "code-reviewer"
start assets add "commit-message"

# Project requirements (local)
cd ~/projects/myapp
start assets add "go-expert" --local
start assets add "security-audit" --local
git add .start/
git commit -m "Add project-specific assets"
```

### Multi-File Asset Handling

Assets may consist of multiple files:

```
pre-commit-review.toml       # Required (task definition)
pre-commit-review.md         # Optional (documentation)
pre-commit-review.meta.toml  # Metadata
```

**All files downloaded** (including metadata) when installing an asset. This ensures full offline capability and metadata availability for updates.

**Only .toml added to config** (references other files as needed).

### Already Installed Behavior

If asset already exists in target config:

**Options presented:**

1. Reinstall (update cache, keep config)
2. Add to other scope (global → local or vice versa)
3. Cancel

**Reinstall use case:** Update cached files without changing config.

## See Also

- start-assets(1) - Asset management overview
- start-assets-search(1) - Search without installing
- start-assets-browse(1) - Visual catalog browsing
- start-assets-info(1) - Preview asset details
- start-config-task(1) - Manage local tasks
- start-config-role(1) - Manage local roles
