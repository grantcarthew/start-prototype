# start assets clean

## Name

start assets clean - Remove unused cached assets

## Synopsis

```bash
start assets clean             # Interactive cleanup
start assets clean --force     # Delete all (no prompts)
start assets clean --dry-run   # Preview without deleting
```

## Description

Remove unused cached assets to free disk space. Scans global and local configuration files to identify referenced assets, then prompts for removal of each configured asset. Automatically deletes cached assets not referenced in any configuration.

**Two-phase cleanup:**

**Phase 1: Configured assets** (interactive prompts)
- Scan global and local configs for asset references
- Prompt user for each configured asset
- YES → remove from config + delete from cache
- NO → keep in config + keep in cache

**Phase 2: Orphaned assets** (automatic)
- Find cached assets not in any config
- Delete automatically (no prompts)
- Frees space from old/removed assets

**Safety features:**
- Interactive prompts by default
- Config files backed up before modification
- Dry-run mode for preview
- Clear summary of actions taken

## Behavior

### Interactive Mode (Default)

```
1. Scan ~/.config/start/*.toml for asset references
2. Scan ./.start/*.toml for asset references
3. For each configured asset:
   - Display asset name and scope (global/local)
   - Prompt: Remove? [y/N]
   - If YES: mark for config removal + cache deletion
   - If NO: keep as-is
4. Find cached assets not in any config
5. Delete all marked/orphaned assets
6. Create config backups before modification
7. Remove assets from configs
8. Delete from cache
9. Display summary
```

### Force Mode (--force)

```
1. Skip all prompts
2. Remove ALL asset references from configs
3. Delete ALL cached assets
4. Create config backups
5. Display summary
```

**Nuclear option** - Use with caution.

### Dry Run Mode (--dry-run)

```
1. Perform full analysis
2. Show what would be deleted
3. Make no actual changes
4. Exit
```

## Output

### Interactive Mode

```bash
$ start assets clean

Scanning configurations...
  Global: ~/.config/start/
  Local:  ./.start/

Found 5 asset references:

Global config (3 assets):
  Remove task 'pre-commit-review'? [y/N]: y
  Remove role 'code-reviewer'? [y/N]: n
  Remove task 'commit-message'? [y/N]: y

Local config (2 assets):
  Remove task 'security-audit'? [y/N]: y
  Remove role 'go-expert'? [y/N]: n

Found 3 cached assets not in any config:
  - tasks/experimental/old-task (1.2 KB)
  - tasks/testing/deprecated-test (850 bytes)
  - roles/unused/old-role (2.1 KB)

Summary of planned deletions:
  From config: 3 assets (pre-commit-review, commit-message, security-audit)
  Orphaned: 3 assets
  Total: 6 assets (4.15 KB)

Creating config backups...
  ✓ ~/.config/start/tasks.toml.2025-01-13-143052.backup
  ✓ ./.start/tasks.toml.2025-01-13-143052.backup

Updating configurations...
  ✓ Removed 2 tasks from global config
  ✓ Removed 1 task from local config

Deleting from cache...
  ✓ tasks/git-workflow/pre-commit-review
  ✓ tasks/git-workflow/commit-message
  ✓ tasks/quality/security-audit
  ✓ tasks/experimental/old-task
  ✓ tasks/testing/deprecated-test
  ✓ roles/unused/old-role

✓ Cleaned 6 assets (4.15 KB freed)

Kept in cache:
  - roles/general/code-reviewer (in global config)
  - roles/languages/go-expert (in local config)
```

### Force Mode

```bash
$ start assets clean --force

⚠ FORCE MODE: This will delete ALL assets and config references

Continue? [y/N]: y

Scanning configurations...
  Global: 3 asset references
  Local:  2 asset references
  Cached: 8 assets total

Creating config backups...
  ✓ ~/.config/start/tasks.toml.2025-01-13-143100.backup
  ✓ ~/.config/start/roles.toml.2025-01-13-143100.backup
  ✓ ./.start/tasks.toml.2025-01-13-143100.backup
  ✓ ./.start/roles.toml.2025-01-13-143100.backup

Removing all asset references from configs...
  ✓ Cleared 2 tasks, 1 role from global config
  ✓ Cleared 1 task, 1 role from local config

Deleting all cached assets...
  ✓ Deleted 8 assets (12.4 KB)

✓ Force clean complete
  Deleted: 8 assets (12.4 KB freed)
  Configs: Backed up and cleared

Backups saved. To restore:
  mv ~/.config/start/tasks.toml.2025-01-13-143100.backup ~/.config/start/tasks.toml
```

### Dry Run Mode

```bash
$ start assets clean --dry-run

Scanning configurations (DRY RUN)...
  Global: ~/.config/start/
  Local:  ./.start/

Would prompt for removal:

Global config (3 assets):
  - task 'pre-commit-review'
  - role 'code-reviewer'
  - task 'commit-message'

Local config (2 assets):
  - task 'security-audit'
  - role 'go-expert'

Would auto-delete (not in any config):
  - tasks/experimental/old-task (1.2 KB)
  - tasks/testing/deprecated-test (850 bytes)
  - roles/unused/old-role (2.1 KB)

If you answered YES to all prompts:
  Would delete: 8 assets (12.4 KB)
  Would backup: 4 config files

No changes made (dry run).
Run without --dry-run to perform cleanup.
```

### Nothing to Clean

```bash
$ start assets clean

Scanning configurations...
  Global: ~/.config/start/
  Local:  ./.start/

Found 3 asset references (all in use)
Found 0 orphaned cached assets

Nothing to clean.
All cached assets are referenced in configurations.
```

Exit code: 0

### No Cached Assets

```bash
$ start assets clean

Scanning configurations...

No cached assets found.

Use 'start assets add <query>' to install assets.
```

Exit code: 0

### User Cancellation

```bash
$ start assets clean --force

⚠ FORCE MODE: This will delete ALL assets and config references

Continue? [y/N]: n

Cancelled.
```

Exit code: 3

## Exit Codes

**0** - Success (cleaned or nothing to clean)

**1** - File system error (cannot read/write configs or cache)

**2** - Backup failed

**3** - User cancelled

## Flags

**--force**, **-f**
: Delete all assets and config references without prompts. Nuclear option.

**--dry-run**, **-n**
: Show what would be deleted without making changes.

**--yes**, **-y**
: Answer YES to all prompts (faster than interactive, safer than --force).

**--keep-cache**
: Remove from configs but keep files in cache.

## Examples

### Basic Interactive Cleanup

```bash
$ start assets clean

Scanning configurations...

Found 5 asset references:

Global config (3 assets):
  Remove task 'pre-commit-review'? [y/N]: y
  Remove role 'code-reviewer'? [y/N]: n
  Remove task 'commit-message'? [y/N]: n

Local config (2 assets):
  Remove task 'security-audit'? [y/N]: y
  Remove role 'go-expert'? [y/N]: n

Found 2 orphaned cached assets:
  - tasks/old/deprecated (500 bytes)
  - roles/test/experimental (1.1 KB)

Deleting 3 assets...

✓ Cleaned 3 assets (1.6 KB freed)
```

### Preview with Dry Run

```bash
$ start assets clean --dry-run

Would prompt for 5 configured assets
Would auto-delete 2 orphaned assets
Maximum deletion: 7 assets (8.2 KB)

No changes made.
```

### Force Delete All

```bash
$ start assets clean --force

⚠ FORCE MODE: Delete ALL assets and references?
Continue? [y/N]: y

Deleting 8 assets and clearing configs...

✓ Force clean complete
  Deleted: 8 assets (12.4 KB freed)
  Backups: 4 files saved
```

### Auto-Answer YES

```bash
$ start assets clean --yes

Scanning configurations...

Auto-answering YES to all prompts...

Global config:
  ✓ Removing task 'pre-commit-review'
  ✓ Removing role 'code-reviewer'
  ✓ Removing task 'commit-message'

Local config:
  ✓ Removing task 'security-audit'
  ✓ Removing role 'go-expert'

Deleting 5 assets + 2 orphaned...

✓ Cleaned 7 assets (9.1 KB freed)
```

### Keep Cache Files

```bash
$ start assets clean --keep-cache

Scanning configurations...

[... prompts for config removal ...]

Removing from configs only (keeping cache files)...

✓ Removed 3 assets from configs
✓ Cache files preserved

Cache still contains 8 assets (12.4 KB).
Use 'start assets clean' without --keep-cache to delete.
```

## Use Cases

### Free Disk Space

**Problem:** Cache growing too large.

```bash
# Preview space savings
start assets clean --dry-run

# Clean interactively
start assets clean
```

### Remove Unused Experiments

**Problem:** Tried many assets, only using a few.

```bash
start assets clean
# Answer NO to assets you use
# Answer YES to assets you don't need
```

### Fresh Start

**Problem:** Want to remove all assets and start over.

```bash
# Nuclear option
start assets clean --force
```

Removes everything, backups configs.

### Selective Cleanup

**Problem:** Want to remove only orphaned assets (not configured ones).

```bash
start assets clean --yes --keep-cache
# Removes from config but preserves cache
```

### Audit Before Cleanup

**Problem:** Want to see what would be deleted.

```bash
start assets clean --dry-run > cleanup-plan.txt
# Review plan
start assets clean
```

## Comparison with Other Commands

### vs `start assets update`

**`start assets clean`** - Removes unused assets
```bash
start assets clean
# Deletes assets to free space
```

**`start assets update`** - Updates existing assets
```bash
start assets update
# Downloads new versions
```

Clean removes, update refreshes.

### vs manual deletion

**Manual:**
```bash
rm -rf ~/.config/start/assets/*
# Dangerous, no backups, breaks configs
```

**`start assets clean`:**
```bash
start assets clean
# Safe, interactive, backs up configs
```

Always use clean command for safety.

## Configuration

**No configuration required.** Operates on standard paths:
- Cache: `~/.config/start/assets/`
- Global config: `~/.config/start/*.toml`
- Local config: `./.start/*.toml`

## Notes

### What Gets Deleted

**Deleted from cache:**
- Asset files (`.toml`, `.md`, etc.)
- Metadata files (`.meta.toml`)
- Empty category directories

**Deleted from configs:**
- Asset entries in `tasks.toml`, `roles.toml`, etc.

**Preserved:**
- Config backups (timestamped)
- User-created custom assets (if selected NO)

### Config Backups

**Automatic backups before modification:**
```
~/.config/start/tasks.toml
  → ~/.config/start/tasks.toml.2025-01-13-143052.backup

./.start/roles.toml
  → ./.start/roles.toml.2025-01-13-143052.backup
```

**Restore from backup:**
```bash
mv ~/.config/start/tasks.toml.2025-01-13-143052.backup \
   ~/.config/start/tasks.toml
```

### Orphaned Assets

**Definition:** Cached assets not referenced in global or local config.

**Common causes:**
- Removed from config manually
- Deleted config entry
- Experimental assets not added to config

**Auto-deleted** without prompts (safe, they're not used).

### Force Mode Safety

**--force is destructive:**
- Deletes ALL assets
- Removes ALL config references
- Only confirmation prompt prevents data loss

**Use cases:**
- Complete reset
- Starting fresh
- Removing everything

**Alternative:** Use `--yes` for faster cleanup without being nuclear.

### Interactive vs Batch

**Interactive (default):**
- Prompts for each asset
- Full control
- Slower but safer

**Batch (--yes):**
- Auto-answers YES
- No prompts
- Faster, still creates backups

**Force (--force):**
- Nuclear option
- Single confirmation
- Deletes everything

### Dry Run Before Clean

**Best practice:**
```bash
# 1. Preview
start assets clean --dry-run

# 2. Review what would be deleted

# 3. Execute if acceptable
start assets clean
```

### Disk Space Reclaimed

**Typical savings:**
- Small catalog: 1-5 MB
- Medium catalog: 5-20 MB
- Large catalog: 20-100 MB

**Depends on:**
- Number of cached assets
- Asset file sizes
- Orphaned assets

### Custom Assets

**User-created assets** are treated same as catalog assets:
- Prompt for removal if in config
- Auto-delete if not in config

**To preserve custom assets:**
- Answer NO when prompted
- Or don't run clean command

### Re-downloading After Clean

**After cleaning, assets are gone:**
```bash
start assets clean --force  # Deletes all

# Later, re-download if needed
start assets add "pre-commit-review"
```

Re-downloading is safe and easy.

## See Also

- start-assets(1) - Asset management overview
- start-assets-update(1) - Update cached assets
- start-assets-add(1) - Install new assets
- start-config(1) - Manage configuration
- DR-036 - Cache management
