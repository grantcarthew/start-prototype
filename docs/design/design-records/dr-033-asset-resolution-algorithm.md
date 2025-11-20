# DR-033: Asset Resolution Algorithm

- Date: 2025-01-10
- Status: Accepted
- Category: Asset Management

## Problem

When executing commands that reference assets (tasks, roles, agents, contexts), the CLI needs to find the asset definition. The resolution strategy must address:

- Search order (which locations to check and in what priority)
- Merging behavior (combine sources or first-match-wins)
- Download control (when to download from GitHub automatically)
- User configuration (global vs local config preference)
- Cache lookup (how to find assets in cache without knowing category)
- Offline behavior (work without network when possible)
- Error messaging (clear feedback when asset not found)
- Flag precedence (command-line override vs config setting)
- Discovery vs resolution (different algorithms for different purposes)

## Decision

Assets resolve in priority order: local config → global config → cache → GitHub catalog, with automatic download controlled by asset_download setting and --asset-download flag. First match wins with no merging.

Resolution priority order:

When user runs start task <name> (or role, agent, context):

1. Local config (.start/tasks.toml)
2. Global config (~/.config/start/tasks.toml)
3. Asset cache (~/.config/start/assets/tasks/)
4. GitHub catalog (query catalog, download if allowed)
5. Error (not found anywhere)

First match wins - no merging, no fallback to next level once found.

Download control:

Configuration setting (config.toml):

```toml
[settings]
asset_download = true  # Default: auto-download from GitHub
```

Behavior:

- asset_download = true: Auto-download if asset not found in config/cache
- asset_download = false: Fail if asset not found in config/cache

Command-line flag override:

```bash
start task <name> --asset-download[=bool] --local
```

Flag precedence:

- Flag explicitly set: use flag value
- Flag not set: use asset_download setting
- Setting not set: default to true

Scope control:

Downloaded assets added to config:

- Default: global config (~/.config/start/tasks.toml)
- With --local flag: local config (.start/tasks.toml)

Cache behavior:

Cached assets:

- Used immediately without prompting
- Cache is transparent implementation detail
- No user interaction needed

Resolution vs discovery:

This resolution applies to command execution (running tasks, roles, agents, contexts).

Asset discovery commands (start assets search/browse/info) work differently:

- Purpose: Find assets in catalog (explore what's available)
- Search sources: GitHub catalog only (not local/global/cache)
- Match algorithm: Substring matching (name, description, tags)
- Uses index.csv: Yes (for fast searching)

Resolution searches all sources (local → global → cache → GitHub) with exact/prefix matching.

Discovery searches GitHub only with substring matching for exploration.

Recursive Resolution & Content Restoration:

Asset resolution works in conjunction with file restoration (DR-042) to ensure dependencies are available:

1. **Recursive Resolution (Config)**
   - When a Task references a Role or Agent not found in the config, the Resolution Algorithm is invoked recursively.
   - Example: `role = "reviewer"` missing? → Search catalog for "reviewer" → Download & Use.
   - Allows lazy loading of entire dependency chains (Task → Role).

2. **Content Restoration (Files)**
   - When a resolved asset config references a file path (e.g., `prompt_file = ".../assets/..."`) that is missing from disk.
   - The low-level file loader detects the missing file in the asset cache path.
   - Automatically restores the file from the catalog (DR-042).
   - Handles cases where config exists (e.g., shared in git) but content is missing.

## Why

Clear priority order simplifies user mental model:

- Local config takes precedence (project-specific overrides)
- Global config next (personal defaults)
- Cache before GitHub (use what's already downloaded)
- GitHub as final source (discover and download new assets)
- First match wins (no complex merging logic)

Download control provides flexibility:

- Global setting for default behavior (asset_download)
- Per-command flag override (--asset-download)
- Explicit control over network usage
- Users choose when downloads happen

Offline-friendly design:

- Works without network if asset configured or cached
- Clear error messages when network needed
- Graceful degradation
- No unexpected network calls

Cache transparency improves UX:

- Cached assets used immediately (no prompts)
- User doesn't manage cache manually
- Cache is implementation detail
- Reduces network calls automatically

Separate discovery and resolution serves different needs:

- Discovery explores full catalog (GitHub only makes sense)
- Resolution finds configured assets (checks all sources)
- Discovery uses substring matching (broad exploration)
- Resolution uses exact/prefix matching (precise lookup)
- Clear distinction between "what's available" vs "what do I have"

User config always wins:

- Local and global config checked first
- User customizations never overridden by catalog
- Explicit configuration takes precedence
- Predictable behavior

## Trade-offs

Accept:

- Network dependency for new assets (first use of catalog asset requires network, but clear error messages)
- No automatic config detection (can't auto-detect global vs local preference, sensible default with --local flag)
- Cache glob on every lookup (must glob all categories to find asset, filesystem cache is small and fast)
- No merging between sources (first match wins, can't combine local + global, simpler and more predictable)
- Discovery doesn't search local/global (only GitHub catalog, but local/global are already known to user)

Gain:

- Simple and predictable (clear priority order, first match wins, user config always wins)
- Flexible download control (global setting, per-command override, explicit network usage control)
- Offline-friendly (works without network if asset configured or cached, clear errors when needed)
- Cache transparency (cached assets used immediately, no prompts, automatic efficiency)
- Clear discovery vs resolution (different algorithms for different purposes, intuitive separation)
- User config precedence (customizations never overridden, predictable behavior)
- Fast cache lookup (filesystem cache small, glob fast enough)

## Alternatives

Merge config sources instead of first-match-wins:

Example: Combine local + global configs

- Local config provides overrides (command, prompt)
- Global config provides defaults (agent, role)
- Merge fields from both sources

Pros:

- More flexible (partial overrides possible)
- Could reuse global defaults with local tweaks
- Finer-grained control

Cons:

- Complex merging logic (which fields merge, which override?)
- Unclear precedence (is local command used with global agent?)
- Harder to reason about (where did this value come from?)
- Debugging nightmare (multiple sources for single asset)
- Conflicts are confusing (what wins when both specify same field?)

Rejected: First-match-wins is simpler and more predictable. Users can copy entire asset to override.

Always auto-download without control:

Example: No asset_download setting, always download

- Asset not found in config/cache: download automatically
- No user control over downloads
- Always fetch from GitHub when needed

Pros:

- Simpler (no setting to configure)
- Always works (users get assets automatically)
- No "download disabled" errors

Cons:

- Unexpected network calls (privacy/security concern)
- No control for restricted environments (corporate firewalls)
- Can't disable for offline work
- Surprising behavior (network when not expected)

Rejected: User control over network important. Some environments require explicit opt-in for downloads.

Never auto-download:

Example: Always require explicit start assets add

- No automatic downloads ever
- Users must browse catalog and explicitly add
- More manual workflow

Pros:

- Very explicit (user controls everything)
- No surprise downloads
- Clear when network used

Cons:

- Friction for common workflows (must add before use)
- Worse UX (extra steps for common tasks)
- Can't discover assets during execution
- More tedious for exploration

Rejected: Lazy loading with download control provides better UX while maintaining user control.

Search all sources for discovery:

Example: start assets search queries local + global + cache + GitHub

- Shows assets from all locations
- Comprehensive search results

Pros:

- See everything in one search
- Discover what's already configured
- Comprehensive results

Cons:

- Confusing (user already knows what's configured)
- Mixing exploration with status check
- Different mental models (discovering vs reviewing)
- Cluttered results (known + unknown mixed)

Rejected: Discovery is about exploring catalog. Local/global configs are already known. Separate concerns are clearer.

## Structure

Resolution algorithm:

Priority order (first match wins):

1. Check local config (.start/tasks.toml)
   - Exact match by name
   - Then check alias
   - If found: return asset, stop
   - If not found: continue to step 2

2. Check global config (~/.config/start/tasks.toml)
   - Exact match by name
   - Then check alias
   - If found: return asset, stop
   - If not found: continue to step 3

3. Check asset cache (~/.config/start/assets/tasks/)
   - Glob pattern: ~/.config/start/assets/{type}/*/{name}.toml
   - Any category (don't know category in advance)
   - If found: load from cache, return asset, stop
   - If not found: continue to step 4

4. Check if downloads allowed
   - Check asset_download setting from config.toml [settings] (default: true)
   - Check --asset-download flag (overrides setting)
   - If disabled: return error (download disabled)
   - If enabled: continue to step 5

5. Query GitHub catalog
   - Check network connectivity
   - If offline: return error (network required)
   - Download index.csv
   - Search for asset by name
   - If not found: return error (not in catalog)
   - If found: download asset, cache it, add to config, return asset

Download behavior:

When asset found in GitHub:

- Download all asset files (.toml, .md, .meta.toml)
- Cache to ~/.config/start/assets/{type}/{category}/
- Add to config (global or local based on --local flag)
- Log download and config addition
- Return asset for immediate use

Cache lookup:

Pattern: ~/.config/start/assets/{type}/*/{name}.toml

- Type is known (tasks, roles, agents, contexts)
- Category is unknown (user doesn't specify it)
- Glob all categories to find match
- Return first match (should only be one)

Config addition:

Scope determination:

- Default: global (~/.config/start/tasks.toml)
- With --local flag: local (.start/tasks.toml)
- Local requires existing .start/ directory

Add to config:

- Load existing config TOML
- Add asset entry (inline all content)
- Write back to file
- Preserve existing entries

Flag precedence:

asset_download determination:

1. If --asset-download flag provided: use flag value
2. Else if config.toml [settings] asset_download set: use setting value
3. Else: default to true

Examples:

- No setting, no flag: true (default)
- Setting false, no flag: false
- Setting false, flag true: true (flag wins)
- Setting true, flag false: false (flag wins)

Resolution vs discovery distinction:

Resolution (this DR):

- Purpose: Execute configured asset
- Commands: start task <name>, start --role <name>
- Search sources: Local → Global → Cache → GitHub
- Match algorithm: Exact match, then alias
- Search fields: Name only
- Uses index.csv: No

Discovery (DR-039, DR-040, DR-041):

- Purpose: Find assets in catalog
- Commands: start assets search, start assets browse
- Search sources: GitHub catalog only
- Match algorithm: Substring matching
- Search fields: Name, path, description, tags
- Uses index.csv: Yes

## Usage Examples

Default behavior (asset_download = true):

```bash
$ start task pre-commit-review

Task 'pre-commit-review' not found locally.
Found in GitHub catalog: tasks/git-workflow/pre-commit-review
Downloading...

✓ Cached to ~/.config/start/assets/tasks/git-workflow/
✓ Added to global config as 'pre-commit-review'

Running task 'pre-commit-review'...
[task executes]
```

Add to local config:

```bash
$ start task pre-commit-review --local

Task 'pre-commit-review' not found locally.
Found in GitHub catalog: tasks/git-workflow/pre-commit-review
Downloading...

✓ Cached to ~/.config/start/assets/tasks/git-workflow/
✓ Added to local config as 'pre-commit-review'

Running task 'pre-commit-review'...
[task executes]
```

Downloads disabled (asset_download = false):

```bash
$ start task pre-commit-review

Error: Task 'pre-commit-review' not found

  ✗ Not in local config (.start/tasks.toml)
  ✗ Not in global config (~/.config/start/tasks.toml)
  ✗ Not in asset cache (~/.config/start/assets/)
  ⚠ GitHub download disabled (asset_download = false)

To resolve:
  - Enable: start task pre-commit-review --asset-download
  - Add manually: start assets add
  - Browse catalog: start assets browse
```

Override with flag:

```bash
$ start task pre-commit-review --asset-download=false

Error: Task 'pre-commit-review' not found

  ✗ Not in local config
  ✗ Not in global config
  ✗ Not in asset cache
  ⚠ Download disabled by --asset-download=false flag

To resolve:
  - Remove flag to allow download
  - Add manually: start assets add
```

Found in cache:

```bash
$ start task pre-commit-review

Using cached asset: pre-commit-review
Running task 'pre-commit-review'...
[task executes immediately]
```

Note: Cached assets used immediately without prompting.

Offline (no network):

```bash
$ start task pre-commit-review

Error: Task 'pre-commit-review' not found

  ✗ Not in local config
  ✗ Not in global config
  ✗ Not in asset cache
  ⚠ Cannot check GitHub catalog (offline)

To resolve:
  - Check spelling: 'pre-commit-review'
  - Add manually when online: start assets add
  - Use a configured task: start config task list
```

Not in catalog:

```bash
$ start task nonexistent-task

Error: Task 'nonexistent-task' not found

  ✗ Not in local config
  ✗ Not in global config
  ✗ Not in asset cache
  ✗ Not in GitHub catalog

To resolve:
  - Check spelling: 'nonexistent-task'
  - Browse available: start assets browse
  - Create custom: start config task add
```

Found in local config (priority):

```bash
$ start task my-task

Running task 'my-task'...
[executes from .start/tasks.toml immediately]
```

Note: Local config checked first, no cache or GitHub lookup needed.

Behavior matrix examples:

Setting true, no flag:

```bash
# asset_download = true in config.toml [settings]
$ start task new-task
# Downloads from GitHub if not found
```

Setting false, no flag:

```bash
# asset_download = false in config.toml [settings]
$ start task new-task
# Error: download disabled
```

Setting false, flag true:

```bash
# asset_download = false in config.toml [settings]
$ start task new-task --asset-download
# Downloads (flag overrides setting)
```

Setting true, flag false:

```bash
# asset_download = true in config.toml [settings]
$ start task new-task --asset-download=false
# Error: download disabled (flag overrides setting)
```

Add to local with flag:

```bash
$ start task new-task --asset-download --local
# Downloads and adds to .start/tasks.toml
```

## Updates

- 2025-01-17: Initial version aligned with schema; corrected to config.toml with [settings] section
