# DR-042: Missing Asset Restoration

- Date: 2025-01-18
- Status: Accepted
- Category: Asset Management

## Problem

When a configuration file (global or local) references assets in the asset cache (`~/.config/start/assets/...`), the referenced content files may be missing. This happens when:

1. **Cloning a project**: Local config (`./.start/`) exists, but cache is empty.
2. **Restoring backups**: Global config (`~/.config/start/`) restored, but cache excluded.
3. **Sharing configs**: Copying a `tasks.toml` to another machine.
4. **Cache clearing**: Accidental deletion of the `assets/` directory.

This issue affects any referenced asset file, including:

- Task prompt files (`prompt_file` in tasks)
- Role prompt files (system prompts referenced by tasks or defaults)
- Context templates (if referenced from cache)

Currently, standard file handling (DR-008) would report "file not found" and fail.

## Decision

Implement "Asset Restoration" logic in the low-level file resolution layer.

Whenever the CLI attempts to resolve a file path (for a task, role, or context) that:

1. Is missing from the filesystem
2. Is located within the configured `asset_path` (default `~/.config/start/assets/`)

The CLI will automatically attempt to restore the asset from the GitHub catalog before proceeding.

**Restoration Logic:**

1. **Intercept Missing File:** Detect that a requested file path does not exist.
2. **Check Path Location:** Verify the path is a subdirectory of the global `asset_path`.
3. **Extract Identity:** Parse the path to extract asset identity.
   - Pattern: `{asset_path}/{type}/{category}/{name}.{ext}`
   - Example: `.../assets/roles/general/code-reviewer.md` -> type=`roles`, category=`general`, name=`code-reviewer`
4. **Verify Settings:** Ensure `asset_download` setting is enabled.
5. **Catalog Lookup:** Check the GitHub catalog index for this asset.
6. **Restore:** Download the asset files (content + metadata) to the cache location.
7. **Retry:** Return the path to the caller, which now exists.

If restoration fails (not in catalog, network error), fall back to the standard "file not found" error or warning behavior defined in DR-008.

## Why

Universal fix for all file-backed asset types:

- Works for Tasks, Roles, and Contexts.
- Works for **Global** and **Local** configurations equally.
- Solves the problem at the root (file access) rather than per-command.

Enables seamless workflows:

- **Team Standardization:** Clone a repo, run a task -> assets restore automatically.
- **Backup Recovery:** Restore config files only -> assets restore on first use.
- **Portability:** Copy `tasks.toml` to a new machine -> assets restore automatically.

Maintains separation of concerns:

- Config remains lightweight and committable.
- Content remains in global cache (deduplicated).
- "Lazy loading" extends to "lazy restoration" of specific files.

Safe and predictable:

- Only attempts restoration for files inside the managed `asset_path`.
- Uses standard catalog resolution.
- Does not modify the configuration file (restores the _content_ the config points to).

## Trade-offs

Accept:

- Slight delay on first run (network check).
- Implicit behavior (magic "healing").
- Requires standard path structure in cache to reverse-engineer asset identity.

Gain:

- "It just works" experience for new team members.
- Zero-setup onboarding for shared projects and backup restoration.
- Robustness against cache deletion.

## Example

**Scenario:**

1. `tasks.toml` has:

   ```toml
   [tasks.code-review]
   prompt_file = "~/.config/start/assets/tasks/git/review.md"
   role = "reviewer"
   ```

2. `roles.toml` (referenced by "reviewer") has:

   ```toml
   [roles.reviewer]
   file = "~/.config/start/assets/roles/general/reviewer.md"
   ```

3. Both `.md` files are missing from disk.

**Behavior:**

1. User runs `start task code-review`.
2. **Task Load:** CLI resolves `review.md`. Missing!
   - Restoration triggers. Parses path. Downloads `tasks/git/review`.
   - File now exists. Task loads.
3. **Role Load:** Task references "reviewer" role. CLI resolves `reviewer.md`. Missing!
   - Restoration triggers. Parses path. Downloads `roles/general/reviewer`.
   - File now exists. Role loads.
4. Task executes successfully.
