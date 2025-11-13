# DR-041: Asset Command Reorganization

**Date:** 2025-01-13
**Status:** Accepted
**Category:** CLI Design

## Decision

Create a unified `start assets` command suite for all asset discovery, installation, and management operations. Deprecate scattered asset-related commands (`start config [type] add`, `start update`, `start show assets`).

## Problem

Asset management functionality is scattered across multiple commands with inconsistent interfaces:

**Old structure:**
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

**Issues:**

1. **Discovery is hard** - No search or browse functionality
2. **Inconsistent UX** - Each asset type uses different command pattern
3. **No asset information** - Can't preview before installing
4. **Update is unclear** - `start update` name conflicts with "update CLI binary"
5. **No cache management** - Can't clean old/unused assets
6. **Type coupling** - Must know asset type before discovering

## Solution

### Unified Command Structure

Create `start assets` as the central hub for all asset operations:

```bash
# Discovery and search
start assets browse              # Interactive catalog browser (tree UI)
start assets search <query>      # Search by name/description/tags
start assets info <query>        # Show detailed asset information

# Installation
start assets add <query>         # Search and install asset

# Management
start assets update [query]      # Update cached assets (all or matching query)
start assets clean               # Remove unused cached assets

# Maintenance (for catalog contributors)
start assets index               # Generate catalog index.csv
```

**Key principle:** Asset-first design - discover assets, then install to appropriate config.

### Command Details

#### start assets browse

Interactive tree-based catalog browser.

**Behavior:**
1. Download catalog index
2. Display hierarchical tree by type → category
3. Navigate with arrow keys, select with enter
4. Show asset details in preview pane
5. Select action: add, info, cancel

**Example:**
```
Catalog Browser (46 assets)
═══════════════════════════════════════════════════════════

▼ tasks (28 assets)
  ▼ git-workflow (12 assets)
    ▸ commit-message
    ▸ pre-commit-review      ← Preview
    ▸ post-commit-hook
  ▸ quality (8 assets)
  ▸ documentation (8 assets)

▼ roles (12 assets)
▸ agents (4 assets)
▸ contexts (2 assets)

─────────────────────────────────────────────────────────────
Preview: pre-commit-review

Description: Review staged changes before committing
Tags: git, review, quality, pre-commit
Size: 2.1 KB
Updated: 2025-01-10

[a] Add  [i] Info  [q] Quit
```

See [DR-035](./dr-035-interactive-browsing.md) for detailed UI specification.

#### start assets search \<query\>

Search catalog by substring matching (name, path, description, tags).

Uses [DR-040](./dr-040-substring-matching.md) substring matching algorithm.

**Example:**
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

#### start assets add \<query\>

Search for asset and install to config.

**Behavior:**
1. Search catalog using substring matching
2. If single match → auto-select
3. If multiple matches → interactive selection
4. Show asset details and confirm installation
5. Download to cache
6. Add to global or local config (via --local flag)

**Example:**
```bash
$ start assets add "pre-commit-review"

Found 1 match (exact):
  tasks/git-workflow/pre-commit-review

Downloading...
✓ Cached to ~/.config/start/assets/tasks/git-workflow/
✓ Added to global config as 'pre-commit-review'

Use 'start task pre-commit-review' to run.
```

**With --local flag:**
```bash
$ start assets add "code-reviewer" --local

Found 1 match:
  roles/general/code-reviewer

Downloading...
✓ Cached to ~/.config/start/assets/roles/general/
✓ Added to local config (./.start/roles.toml)

Use 'start --role code-reviewer' to use this role.
```

#### start assets info \<query\>

Show detailed information about an asset without installing.

**Displays:**
- Full metadata from .meta.toml
- Description and tags
- File size and dates
- Installation status (cached, installed in global/local)
- Full file contents preview (optional)

**Example:**
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

#### start assets update [query]

Update cached assets from GitHub catalog.

**Without query (update all):**
```bash
$ start assets update

Checking for updates...

✓ tasks/git-workflow/pre-commit-review  (v1.0 → v1.1)
✓ roles/general/code-reviewer           (up to date)
✓ tasks/quality/commit-lint             (v2.3 → v2.4)

Updated 2 assets, 1 up to date.
```

**With query (update matching):**
```bash
$ start assets update "commit"

Checking for updates to assets matching 'commit'...

✓ tasks/git-workflow/commit-message     (v1.2 → v1.3)
✓ tasks/git-workflow/pre-commit-review  (up to date)
✓ tasks/git-workflow/post-commit-hook   (v1.0 → v1.1)

Updated 2 assets, 1 up to date.
```

See [DR-037](./dr-037-asset-updates.md) for update algorithm.

#### start assets clean

Remove unused cached assets.

**Definition of "unused":**
- Cached but NOT referenced in global or local config
- Orphaned cache files from removed assets

**Behavior:**
```bash
$ start assets clean

Scanning cache...

Unused assets (not in any config):
  tasks/experimental/old-task (1.2 KB)
  roles/deprecated/old-role (850 bytes)

Remove 2 unused assets (2.0 KB)? [y/N]: y

✓ Removed 2 assets
✓ Freed 2.0 KB

Cache: 46 assets (125 KB)
```

**Dry run:**
```bash
$ start assets clean --dry-run

Would remove:
  tasks/experimental/old-task (1.2 KB)
  roles/deprecated/old-role (850 bytes)

Total: 2 assets (2.0 KB)

Use 'start assets clean' to remove.
```

#### start assets index

Generate `assets/index.csv` for catalog contributors.

**Context:** Must be run in catalog repository clone.

**Behavior:**
1. Validate git repository structure
2. Scan assets/ directory for .meta.toml files
3. Parse metadata
4. Generate sorted CSV index
5. Write to assets/index.csv

See [DR-039](./dr-039-catalog-index.md) for complete specification.

**Example:**
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

### Deprecated Commands

**Removed commands:**

| Old Command | Replacement | Notes |
|-------------|-------------|-------|
| `start config task add` | `start assets add` | Universal asset installer |
| `start config role add` | `start assets add` | Universal asset installer |
| `start config agent add` | `start assets add` | Universal asset installer |
| `start update` | `start assets update` | Clearer naming |
| `start show assets` | `start assets browse` | Better discovery UX |

**Backward compatibility:**

Old commands remain as **deprecated aliases** for one major version:

```bash
$ start config task add "code-review"

⚠ Deprecated: 'start config task add' is deprecated.
  Use: start assets add "code-review"

Redirecting to 'start assets add'...

[proceeds with new command]
```

After deprecation period, commands removed entirely with error:

```bash
$ start config task add "code-review"

Error: Command removed. Use 'start assets add "code-review"'
See: https://github.com/grantcarthew/start#migration-guide
```

### Migration Path

**Phase 1: Introduce (current release)**
- Add `start assets` command suite
- Mark old commands as deprecated (with warnings)
- Update all documentation to use new commands

**Phase 2: Deprecation (next major version)**
- Old commands removed from help output
- Deprecated warnings become more prominent
- Migration guide published

**Phase 3: Removal (major version after that)**
- Old commands removed entirely
- Error messages with migration instructions

**Migration guide provided:**
```bash
# Old → New command mapping

# Adding assets
start config task add "foo"  →  start assets add "foo"
start config role add "bar"  →  start assets add "bar"
start config agent add "baz" →  start assets add "baz"

# Updating assets
start update                 →  start assets update

# Browsing assets
start show assets            →  start assets browse

# New capabilities (no old equivalent)
start assets search "query"  # Search by description/tags
start assets info "asset"    # Preview before installing
start assets clean           # Remove unused cache
```

## Benefits

**User experience:**
- ✅ **Unified discovery** - One place for all asset operations
- ✅ **Searchable** - Find assets by description, tags, not just name
- ✅ **Preview before install** - See details with `start assets info`
- ✅ **Type-agnostic** - Don't need to know if something is a task/role/agent
- ✅ **Clear organization** - Logical grouping under `assets` command

**Developer experience:**
- ✅ **Consistent implementation** - Same code paths for all asset types
- ✅ **Easier to extend** - Add new asset types without new commands
- ✅ **Better testing** - Unified test suite for asset operations

**Discoverability:**
- ✅ **Help is clearer** - `start assets --help` shows all operations
- ✅ **Intuitive naming** - "assets" clearly indicates catalog operations
- ✅ **Less cognitive load** - One command to remember, not four

## Trade-offs Accepted

**Learning curve for existing users:**
- ❌ Users must learn new command structure
- **Mitigation:** Deprecation warnings with command mapping, migration guide

**Command is longer:**
- ❌ `start assets add` is longer than `start config task add`
- **Mitigation:**
  - Still supports search shortcuts: `start assets add "pre"` (prefix matching)
  - Most users will use interactive browse or search anyway

**Breaking change:**
- ❌ Requires migration for existing scripts/workflows
- **Mitigation:** Phased deprecation with long transition period

## Implementation Strategy

### Phase 1: Add new commands (parallel)

```go
// New unified command
func assetsCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "assets",
        Short: "Discover and manage catalog assets",
    }

    cmd.AddCommand(assetsBrowseCmd())
    cmd.AddCommand(assetsSearchCmd())
    cmd.AddCommand(assetsAddCmd())
    cmd.AddCommand(assetsInfoCmd())
    cmd.AddCommand(assetsUpdateCmd())
    cmd.AddCommand(assetsCleanCmd())
    cmd.AddCommand(assetsIndexCmd())

    return cmd
}

// Shared implementation
func addAsset(query string, scope string) error {
    // Universal asset installer
    // Used by: start assets add
    // Also used by: deprecated start config [type] add
}
```

### Phase 2: Deprecate old commands

```go
// Deprecated command with warning
func configTaskAddCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:        "add",
        Short:      "Add task from catalog (DEPRECATED - use 'start assets add')",
        Deprecated: "Use 'start assets add' instead",
        Run: func(cmd *cobra.Command, args []string) {
            printDeprecationWarning("start config task add", "start assets add")
            // Delegate to new implementation
            addAsset(args[0], "tasks")
        },
    }
    return cmd
}
```

### Phase 3: Remove old commands

```go
// Remove command entirely, show error
// (In future version after deprecation period)
```

## Examples

### Discovery Workflow (New)

**Before (scattered, no search):**
```bash
# Can't search, must know exact name
start config task add "pre-commit-review"  # Hope it exists

# Can't browse, must check GitHub web
# No way to see what's available
```

**After (unified search and browse):**
```bash
# Browse interactively
start assets browse

# Or search by description
start assets search "commit"

# Preview before installing
start assets info "pre-commit-review"

# Install when ready
start assets add "pre-commit-review"
```

### Update Workflow

**Before (unclear command):**
```bash
start update              # Updates what? CLI? Assets? Both?
```

**After (explicit):**
```bash
start assets update       # Clear: updates cached assets
start assets update "git" # Update only git-related assets
```

### Cache Management

**Before (no cache management):**
```bash
# No way to clean unused assets
# Cache grows indefinitely
```

**After (explicit cleanup):**
```bash
start assets clean        # Remove unused assets
start assets clean --dry-run  # See what would be removed
```

## Related Decisions

- [DR-017](./dr-017-cli-reorganization.md) - CLI command reorganization (general structure)
- [DR-031](./dr-031-catalog-based-assets.md) - Catalog-based assets (asset architecture)
- [DR-035](./dr-035-interactive-browsing.md) - Interactive browsing (browse command)
- [DR-037](./dr-037-asset-updates.md) - Asset updates (update command)
- [DR-039](./dr-039-catalog-index.md) - Catalog index (powers search)
- [DR-040](./dr-040-substring-matching.md) - Substring matching (search algorithm)

## Future Considerations

**Asset collections/bundles:**
- Could add `start assets add-collection "go-dev"` for installing related asset groups
- Example: "go-dev" installs go-expert role, go-related tasks, etc.

**Asset ratings/popularity:**
- Could add metadata for download counts, ratings
- Show in search results: "★★★★☆ (234 downloads)"

**Local asset repositories:**
- Could support custom asset sources (not just GitHub)
- `start assets add "foo" --from private-repo`

**Asset dependencies:**
- Tasks could declare required roles/contexts
- Auto-install dependencies when installing asset

**Current stance:** Ship with core functionality. Monitor usage patterns and add advanced features based on user feedback.
