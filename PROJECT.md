# Design Thoughts

This document is just my thoughts dumped for reference. Do not use this document as a concrete reference. It is most likely wrong.

## General Ideas

- ~~I want it to be able to use a prompt writer prompt to create role documents on the fly~~ → **RESOLVED**: Added to `docs/ideas/assets.md` as catalog assets (`roles/meta/role-writer.md` + `tasks/new-role.toml`)
- ~~Need an easy way to switch defaults~~ → **RESOLVED**: Already exists via `start config role default <name>` and `start config agent default <name>`
- ~~Need a config delete option if it does not exist (or remove), something like `start config agent rm xyz`~~ → **RESOLVED**: Already exists - `start config <type> remove` commands in design
- ~~**Dry-run flag**: Add `--dry-run` flag to preview aggregated context without calling the agent~~ → **RESOLVED**: Evolved into `start show` command - see `docs/cli/start-show.md`
  - Execution preview: `start show`, `start show task <name>`, `start show prompt <text>`
  - Content viewer: `start show role`, `start show context`, `start show agent`, `start show task`
- **Unified asset management**: Consider `start assets` subcommand to consolidate asset operations
  - `start assets browse` - Open GitHub catalog in browser (better discoverability than `start config <type> add`)
  - `start assets add <type> <name>` - Add from catalog to config (replaces `start config <type> add`)
  - `start assets update [name]` - Replace vague `start update` with clearer naming; optionally update specific asset
  - `start assets info <name>` - Show asset metadata (description, last updated, dependencies, source)
  - `start assets list` - Show all cached assets with status
  - `start assets clean` - Clear cache (vs manual rm -rf)
  - **Problem:** Current asset management fragmented across `start config <type> add` commands, not discoverable
  - **Benefit:** Single place for all asset operations, clearer commands, better UX
  - **Semantic separation:**
    - `start config` = Manage YOUR configuration (things you've defined/customized)
      - `new` - Create new custom asset
      - `edit` - Edit your asset
      - `remove` - Remove your asset
      - `test` - Test your asset
      - `list` - List your assets
    - `start assets` = Interact with the CATALOG (browse, add from GitHub, update cache)
      - `browse` - View catalog in browser
      - `add` - Add from catalog to config
      - `update` - Update cached assets
      - `info` - Show asset metadata
      - `list` - Show cached assets
  - **Migration:** `start config task add` → `start assets add task`
  - **STATUS:** APPROVED - See detailed specification below

---

## Task: Implement `start assets` Command Suite

**Status:** Ready for implementation
**Date Decided:** 2025-01-13
**Context:** Asset management UX design - resolving fragmented catalog operations into unified command namespace

### Decision Summary

Implement new `start assets` top-level command to consolidate all GitHub catalog operations. This provides semantic separation between managing user configurations (`start config`) and shopping the catalog (`start assets`).

### Final Command Set

**Approved commands:**
1. `start assets browse` - Open GitHub catalog in browser
2. `start assets add <query>` - Download and install asset with substring matching
3. `start assets search <term>` - Search catalog, terminal output only
4. `start assets info <query>` - Show detailed asset metadata
5. `start assets update [query]` - Update cached assets (replaces `start update`)
6. `start assets clean [--force]` - Selective cache cleanup with prompts

**Rejected commands:**
- ~~`start assets catalog`~~ - Folded into list, then list dropped
- ~~`start assets list`~~ - Use browse/search instead; config commands show installed assets
- ~~`start assets show`~~ - Dropped; info command is sufficient
- ~~`start assets diff`~~ - Dropped; unnecessary complexity
- ~~`start assets check`~~ - Dropped; not needed
- ~~`start assets prune`~~ - Dropped; clean does enough
- ~~`start assets verify/repair/test/validate`~~ - Dropped; covered by other commands
- ~~All versioning/pinning/export commands~~ - Dropped; premature optimization

### Command Specifications

#### 1. `start assets browse`

**Purpose:** Open GitHub catalog in default browser for visual exploration

**Behavior:**
- Opens browser to `https://github.com/{asset_repo}/tree/main/assets`
- Uses `[settings] asset_repo` value (default: `grantcarthew/start`)
- Expected regular use case for discovery

**Error Handling:**
- If browser fails to open: print URL with warning, exit 0
- No "continue" needed - command completes after attempt

**Exit Codes:**
- 0 - Success (browser opened or URL printed)

---

#### 2. `start assets add <query>`

**Purpose:** Download asset from catalog, cache it, and add to config

**Replaces:** `start config task add`, `start config role add`, etc.

**Matching Strategy:**
- Substring match against: name, full path, description, tags
- Minimum 3 characters required
- < 3 chars → fall back to interactive browse with warning message
- Uses substring matching DR (to be written)

**Modes:**

**A) Interactive mode** (no args or < 3 chars):
```bash
start assets add
```
- Browse all types and categories
- Tree-based navigation
- Full interactive flow

**B) Search mode** (3+ chars):
```bash
start assets add git-workflow    # Matches directory
start assets add commit-review   # Matches leaf items
start assets add pre             # Partial match
```

**Results Display:**
```
Found 2 matches for 'git':

1. 📁 assets/tasks/git-workflow/
   ├── pre-commit-review.toml - Review staged changes
   ├── pr-ready.toml - Complete PR preparation
   └── commit-message.toml - Generate commit message

2. 📄 assets/tasks/debugging/git-story.toml - Code archaeology

Select [1-2]: _
```

**Directory Selection Flow:**
```
Selected: assets/tasks/git-workflow/ (4 assets)

Install all 4 assets or select individually?
  1) Install all to global config
  2) Install all to local config
  3) Select individually

Choice [1-3]: _
```

**Flags:**
- `--local` - Sets default to local config but doesn't skip interactive prompts

**Direct Install:**
```bash
start assets add tasks/git-workflow/pre-commit-review
start assets add git-workflow/pre-commit-review  # Type can be inferred
```

**Behavior:**
- Downloads from GitHub catalog
- Caches to `~/.config/start/assets/`
- Adds configuration entry to global or local config
- Never searches local/global/cache - only GitHub catalog

**Exit Codes:**
- 0 - Success (asset added)
- 1 - Network error (catalog unavailable)
- 2 - Asset not found
- 3 - User cancelled

---

#### 3. `start assets search <term>`

**Purpose:** Non-interactive terminal search of GitHub catalog

**Behavior:**
- Substring match (same as `add`)
- Terminal output only - no interactive prompts
- Just lists results and exits
- Tree structure output

**Minimum Query Length:**
- 3 characters required
- < 3 chars → error/warning message

**Output Format:**
```
Found 3 matches for 'commit':

📁 assets/tasks/git-workflow/
  └── pre-commit-review.toml - Review staged changes

📄 assets/tasks/git-workflow/commit-message.toml - Generate commit message
```

**Exit Codes:**
- 0 - Matches found
- 1 - No matches or error (< 3 chars, network error)

---

#### 4. `start assets info <query>`

**Purpose:** Show detailed metadata for specific asset

**Matching:**
- Substring matching (same as search/add)
- Multiple matches → interactive selection
- Single match → show info directly

**Display:**
```
Asset: pre-commit-review
Type: task
Category: git-workflow
Description: Review staged changes before committing
Tags: git, review, quality, pre-commit
SHA: a1b2c3d4e5f6...
Size: 2.4 KB
Created: 2025-01-10
Updated: 2025-01-10
Path: assets/tasks/git-workflow/pre-commit-review.toml

Installation Status:
✓ Installed in global config
✓ Cached locally
```

**Status Display:**
- "Installed in global config"
- "Installed in local config"
- "Cached but not configured"
- "Not installed"

**Exit Codes:**
- 0 - Success (info displayed)
- 1 - Asset not found
- 2 - Network error

---

#### 5. `start assets update [query]`

**Purpose:** Update cached assets by comparing SHAs

**Replaces:** `start update` (top-level command completely removed)

**Behavior:**
- No arguments → update all cached assets
- With query → substring match (categories, names, paths, etc.)
- Multiple matches → interactive selection
- Uses same SHA-based detection as DR-037
- Updates cache only, never modifies user config files

**Examples:**
```bash
start assets update                    # Update all
start assets update pre-commit-review  # Update specific asset
start assets update git-workflow       # Update all in category
```

**Uses:** Substring matching DR (to be written)

**Exit Codes:**
- 0 - Success (updates applied)
- 1 - Network error
- 2 - File system error
- 3 - Partial failure

---

#### 6. `start assets clean [--force]`

**Purpose:** Selective cache cleanup with per-asset prompts

**Behavior - Interactive Mode:**

1. Scan global and local configs for asset references
2. Prompt for each asset config entry
3. Delete from cache based on confirmations
4. Auto-delete cached assets not in any config

**Flow:**
```
Found 5 asset references in configs:

Global config (3):
  Remove task 'pre-commit-review'? [y/N]: y
  Remove role 'code-reviewer'? [y/N]: n
  Remove task 'pr-ready'? [y/N]: y

Local config (2):
  Remove task 'git-story'? [y/N]: y
  Remove role 'go-expert'? [y/N]: y

Found 3 cached assets not in any config:
  - tasks/find-bugs
  - tasks/quick-wins
  - roles/rubber-duck

Deleting from cache:
  ✓ tasks/pre-commit-review (removed from config)
  ✓ tasks/pr-ready (removed from config)
  ✓ tasks/git-story (removed from config)
  ✓ roles/go-expert (removed from config)
  ✓ tasks/find-bugs (not in config)
  ✓ tasks/quick-wins (not in config)
  ✓ roles/rubber-duck (not in config)

Kept in cache:
  - roles/code-reviewer (still in global config)

✓ Cleaned 7 assets, kept 1
```

**Logic:**
- User says YES → remove from config + delete from cache
- User says NO → keep in config + keep in cache
- Not in config → auto-delete from cache (no prompt)

**Flag - `--force`:**
- Skip all prompts
- Delete EVERYTHING (cache + all asset configs)
- Nuclear option for complete cleanup

**Order of Operations:**
1. Prompt for config removals first
2. Then delete from cache based on decisions
3. Configs are backed up before removal

**Exit Codes:**
- 0 - Success (cleaned or user cancelled)
- 3 - File system error (backup failed, etc.)

---

### Design Records Needed

**New DRs to write:**

1. **DR-0XX: Substring Matching Algorithm**
   - Define matching behavior for add/search/info/update
   - 3-character minimum
   - Substring match against: name, full path, description, tags
   - Multiple match → interactive selection
   - Case-insensitive
   - Example matching scenarios

2. **DR-0XX: Asset Command Reorganization**
   - Migrate from `start config <type> add` to `start assets add`
   - Remove `start update` top-level command
   - Semantic separation rationale
   - Migration guide for users
   - Update all documentation references

**DRs to update:**

- **DR-017**: CLI command reorganization - add `start assets` commands
- **DR-031**: Catalog-based assets - reference new commands
- **DR-033**: Asset resolution algorithm - note that `start assets` only searches GitHub
- **DR-035**: Interactive browsing - update to reference `start assets add`
- **DR-037**: Asset updates - update to reference `start assets update`

### Documentation Updates Needed

**CLI docs to create:**
- `docs/cli/start-assets.md` - Main command documentation
- `docs/cli/start-assets-browse.md`
- `docs/cli/start-assets-add.md`
- `docs/cli/start-assets-search.md`
- `docs/cli/start-assets-info.md`
- `docs/cli/start-assets-update.md`
- `docs/cli/start-assets-clean.md`

**CLI docs to update:**
- `docs/cli/start-config-task.md` - Remove `add` subcommand, reference `start assets`
- `docs/cli/start-config-role.md` - Remove `add` subcommand, reference `start assets`
- `docs/cli/start-config-agent.md` - Remove `add` subcommand, reference `start assets`
- `docs/cli/start-update.md` - Mark as removed, redirect to `start assets update`
- `docs/cli/start-show.md` - Remove `assets` subcommand, reference `start assets info`

**Config docs to update:**
- `docs/config.md` - Update asset catalog references

### Implementation Notes

**Core Principles:**
- `start assets` ONLY interacts with GitHub catalog
- Never searches local config, global config, or cache
- Asset resolution algorithm (DR-033) remains separate for runtime resolution
- Cache remains invisible (DR-036) - no manual cache inspection commands

**Semantic Separation:**
- `start config` = Manage YOUR configuration (edit, test, remove custom assets)
- `start assets` = Shop THE catalog (browse, add, update from GitHub)

**User Mental Model:**
- "I want to see what's available" → `start assets browse` or `start assets search`
- "I want to see what I have installed" → `start config task list`, `start config role list`, etc.
- "I want to add something" → `start assets add <query>`
- "I want to update my stuff" → `start assets update`
- "I want to clean up" → `start assets clean`
- "I want to configure my stuff" → `start config <type> ...`

### Breaking Changes

**Commands removed:**
- `start update` (replaced by `start assets update`)
- `start config task add` (replaced by `start assets add`)
- `start config role add` (replaced by `start assets add`)
- `start config agent add` (replaced by `start assets add`)
- `start show assets` (use `start assets info` instead)

**Migration path:**
- Document changes in release notes
- Add deprecation warnings in v1 (if we implement deprecation period)
- Or just hard cutover since project is in design phase

### Success Criteria

- [ ] All 6 commands implemented and tested
- [ ] Substring matching DR written and implemented
- [ ] All CLI documentation created
- [ ] All existing docs updated
- [ ] `start config <type> add` commands removed
- [ ] `start update` command removed
- [ ] Integration tests cover all commands
- [ ] User workflows documented with examples
