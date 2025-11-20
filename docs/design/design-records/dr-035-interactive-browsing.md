# DR-035: Interactive Asset Browsing

- Date: 2025-01-10
- Status: Accepted
- Category: Asset Management

## Problem

Users need to discover and install assets from the catalog interactively. The browsing strategy must address:

- Discovery mechanism (how users browse available assets)
- Navigation approach (category-based vs flat list vs search-first)
- TUI complexity (full framework vs simple prompts vs numbered selection)
- Terminal compatibility (SSH, containers, CI environments)
- Non-interactive support (automation and scripts)
- Dependency footprint (external libraries vs native Go)
- Error handling (network failures, user cancellation)
- Command organization (asset-focused vs type-specific commands)
- Confirmation workflow (prevent accidental installations)

## Decision

Provide interactive terminal-based catalog browsing via start assets add (no arguments) using native numbered selection with category-first navigation and confirmation prompts.

Key aspects:

Command: start assets add (no arguments)

- Interactive TUI browser for all asset types
- Downloads index.csv for catalog metadata
- Groups assets by category
- Numbered selection interface

Navigation: Category-first approach

- Step 1: Select category (with asset counts)
- Step 2: Select asset within category (with descriptions)
- Step 3: Confirm download with metadata preview
- Option to view all assets (skip category filtering)

Implementation: Native numbered selection

- No TUI library dependency (no bubbletea, promptui, survey)
- Standard input/output (print numbered lists, read integers)
- Works in all terminal environments
- KISS principle for v1

Confirmation workflow:

- Show asset metadata before downloading (name, description, tags)
- Prompt: "Download and add to config? [Y/n]"
- User can cancel anytime (no changes made)
- Downloads and adds to config only after confirmation

Non-interactive mode:

- --yes flag skips confirmation prompts (for automation)
- --category flag skips category selection
- --local flag adds to local config instead of global
- Direct installation: start assets add <path>

## Why

Native numbered selection provides universal compatibility:

- Works in all terminal environments (SSH, containers, CI)
- No terminal capability detection needed
- No external dependencies (keeps binary small)
- Simple implementation (standard input/output)
- Predictable behavior across platforms
- Fast build times (no complex dependencies)

Category-first navigation improves organization:

- Easier to scan than flat list (grouped by domain)
- Scales to hundreds of assets (smaller lists per screen)
- Clear mental model (type → category → asset)
- Reduces cognitive load (focus on one category at a time)
- Better for discovery (browse related assets together)

Confirmation prompts prevent accidents:

- User sees what they're downloading (name, description, tags)
- Explicit consent required (Y/n prompt)
- Can review before committing (metadata preview)
- No surprise downloads (always prompted)
- Can cancel anytime (no changes made if cancelled)

KISS principle for v1:

- Don't add dependencies speculatively
- Ship simple version first (numbered selection proven pattern)
- Add TUI library later if users request (based on real feedback)
- Keeps binary small and builds fast
- Avoids over-engineering for initial release

Flexible modes support different workflows:

- Interactive browsing (discovery and exploration)
- Direct installation (when user knows asset path)
- Non-interactive mode (automation and scripts)
- Category filtering (skip navigation if category known)

## Trade-offs

Accept:

- No fancy TUI in v1 (no arrow key navigation, no fuzzy search, no multi-select, but can add in v2 if users request)
- Category navigation adds step (extra prompt compared to flat list, but better organization for large catalogs)
- No inline search in TUI (can't filter by keyword while browsing, use start assets search or --category flag instead)
- Numbered input only (typing numbers instead of arrow keys, but works everywhere)
- Manual category selection (can't auto-detect user intent, but --category flag available)

Gain:

- Universal compatibility (works in SSH, containers, CI, all terminals without feature detection)
- Zero dependencies (no external TUI libraries, keeps binary small and builds fast)
- Simple and predictable (numbered lists, clear prompts, straightforward flow)
- Flexible modes (interactive browsing or direct installation with --yes for automation)
- Category organization (easier to scan, better for large catalogs, scales well)
- Confirmation before changes (user sees metadata before installing, can cancel anytime)
- Non-interactive support (automation-friendly with --yes, --category, and direct path installation)

## Alternatives

Use TUI library (bubbletea, promptui, survey):

Example: Integrate charmbracelet/bubbletea for enhanced UX

- Full-featured TUI framework
- Arrow key navigation instead of numbered selection
- Fuzzy search, multi-select capabilities
- Beautiful rendering with colors and boxes
- Modern terminal UI patterns

Pros:

- Better UX (arrow keys more intuitive than typing numbers)
- Richer interactions (fuzzy search, multi-select)
- Modern appearance (colors, borders, animations)
- Familiar pattern (like fzf, kubectl)

Cons:

- External dependency (larger binary, slower builds)
- Terminal compatibility issues (some terminals don't support features)
- Complex API (Elm architecture for bubbletea, steeper learning curve)
- May not work in all environments (SSH, containers, limited terminals)
- Adds maintenance burden (dependency updates, API changes)

Rejected: KISS principle - ship v1 with numbered selection, add TUI library in v2 if users request. Numbered selection works everywhere and avoids complexity.

Flat list without categories:

Example: Show all assets in one long list

```
Available tasks:
  1. git-workflow/pre-commit-review - Review staged changes
  2. git-workflow/pr-ready - Complete PR preparation
  3. code-quality/find-bugs - Find potential bugs
  4. code-quality/quick-wins - Low-hanging refactoring
  [... 42 more ...]
Choice [1-46]: _
```

Pros:

- Simpler (one step instead of two)
- Faster (no category navigation)
- Less code (no grouping logic)

Cons:

- Long list (hard to scan with hundreds of assets)
- Poor organization (no grouping by domain)
- Doesn't scale (overwhelming as catalog grows)
- Cognitive overload (too many choices at once)

Rejected: Category-first navigation better for large catalogs, easier to scan, clearer organization.

Search-first pattern:

Example: Prompt for search query before browsing

```
Search for task (or press Enter to browse all): commit

Found 3 tasks:
  1. pre-commit-review - Review staged changes
  2. commit-message - Generate conventional commit
  3. explain-changes - Understand what changed
```

Pros:

- Fast for users who know what they want
- Reduces browsing (direct to relevant assets)
- Good for power users (targeted search)

Cons:

- Extra step for exploration (must type something)
- Not good for discovery (browsing to see what exists)
- Requires search implementation upfront
- Poor for new users (don't know what to search for)

Rejected: Category browsing better for discovery. Search available via separate start assets search command for targeted queries.

Direct installation without browsing:

Example: Require explicit asset path, no interactive mode

```
start assets add git-workflow/pre-commit-review --yes
```

- No interactive browsing at all
- User must know asset path upfront
- Always requires explicit path

Pros:

- Very explicit (user knows exactly what they're getting)
- No interactive navigation needed
- Simple implementation (no TUI code)

Cons:

- Poor discoverability (must know asset path upfront)
- No exploration (can't browse available assets)
- Requires external documentation (users look up paths elsewhere)
- Worse UX for new users (high barrier to discovery)

Rejected: Interactive browsing critical for discovery. Direct installation available as alternative mode, not replacement.

## Structure

Interactive TUI browser:

Command: start assets add (no arguments)

- Downloads index.csv from GitHub
- Groups assets by type and category
- Shows numbered category list
- User selects category (or "view all")
- Shows numbered asset list with descriptions
- User selects asset
- Shows confirmation prompt with metadata
- Downloads and adds to config on confirmation

Navigation flow:

1. Fetch catalog
   - Download index.csv from raw.githubusercontent.com
   - Parse into in-memory structure
   - Group by type and category

2. Category selection
   - Display: "Select category:"
   - List categories with asset counts (e.g., "1. git-workflow (4 tasks)")
   - Include option: "[view all tasks]" to skip filtering
   - Read user input (integer choice)

3. Asset selection
   - Display: "{category} tasks:"
   - List assets with name and description
   - Read user input (integer choice)

4. Confirmation prompt
   - Display selected asset metadata:
     - Name
     - Description
     - Tags
   - Prompt: "Download and add to config? [Y/n]"
   - Read user input (Y/n)

5. Download and install
   - Download asset files from GitHub (raw URLs)
   - Cache to ~/.cache/start/{type}/{category}/
   - Add to config (global or local based on --local flag)
   - Display success message with usage hint

Non-interactive mode:

Flags:

- --yes, -y: Skip confirmation prompts (for automation)
- --category <cat>: Filter by category (skip category selection step)
- --local: Add to local config instead of global

Direct installation:

- Syntax: start assets add <path>
- Example: start assets add git-workflow/pre-commit-review
- Skips category and asset selection steps
- Shows confirmation prompt (unless --yes provided)
- Downloads and adds to config

Error handling:

Network unavailable:

- Message: "Cannot fetch catalog from GitHub"
- Show network error details
- Suggestions:
  - Check internet connection
  - Use cached assets: start task <name>
  - Add custom task manually

No assets found:

- Message: "No tasks found in category 'X'"
- Show available categories
- Suggest: Try different category

User cancels:

- Message: "Cancelled. No changes made."
- Exit code 0 (user choice, not error)

Invalid input:

- Re-prompt for valid choice
- Show valid range (e.g., "Choice [1-5]:")

## Usage Examples

Interactive browsing:

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
✓ Cached to ~/.cache/start/tasks/git-workflow/
✓ Added to global config as 'pre-commit-review'

Try it: start task pre-commit-review
```

Category filtering (skip category selection):

```bash
$ start assets add --category git-workflow

Fetching catalog from GitHub...
✓ Found 4 assets in git-workflow

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
✓ Cached to ~/.cache/start/tasks/git-workflow/
✓ Added to global config as 'pre-commit-review'
```

Direct installation:

```bash
$ start assets add git-workflow/pre-commit-review

Found asset: pre-commit-review
Description: Review staged changes before commit
Tags: git, review, quality, pre-commit

Download and add to config? [Y/n] y

Downloading...
✓ Cached to ~/.cache/start/tasks/git-workflow/
✓ Added to global config as 'pre-commit-review'
```

Non-interactive mode for automation:

```bash
$ start assets add git-workflow/pre-commit-review --yes

✓ Downloaded and cached pre-commit-review
✓ Added to global config

$ echo $?
0
```

Add to local config:

```bash
$ start assets add --local

Fetching catalog from GitHub...
✓ Found 46 assets across 4 types and 12 categories

Select category:
  1. git-workflow (4 tasks)
  [...]

> 1

git-workflow tasks:
  1. pre-commit-review - Review staged changes before commit
  [...]

> 1

[... confirmation ...]

Downloading...
✓ Cached to ~/.cache/start/tasks/git-workflow/
✓ Added to local config as 'pre-commit-review'
```

Error handling - network unavailable:

```bash
$ start assets add

Error: Cannot fetch catalog from GitHub

Network error: dial tcp: no route to host

To resolve:
- Check internet connection
- Use cached assets: start task <name>
- Add custom task manually
```

Error handling - no assets found:

```bash
$ start assets add --category nonexistent

Error: No tasks found in category 'nonexistent'

Available categories:
  - git-workflow
  - code-quality
  - security
  - debugging

Try: start assets add --category git-workflow
```

Error handling - user cancels:

```bash
$ start assets add

Fetching catalog from GitHub...
✓ Found 46 assets across 4 types and 12 categories

[... user navigates and selects task ...]

Selected: pre-commit-review
Description: Review staged changes before commit
Tags: git, review, quality, pre-commit

Download and add to config? [Y/n] n

Cancelled. No changes made.
```

View all assets (skip category filtering):

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

> 5

All tasks:
  1. git-workflow/pre-commit-review - Review staged changes before commit
  2. git-workflow/pr-ready - Complete PR preparation checklist
  3. git-workflow/commit-message - Generate conventional commit message
  4. git-workflow/explain-changes - Understand what changed in commits
  5. code-quality/find-bugs - Find potential bugs in code
  [... 8 more ...]

> 1

[... confirmation and download ...]
```

## Updates

- 2025-01-17: Initial version aligned with schema; removed implementation code, Related Decisions, and Future Considerations sections; command changed from start config <type> add to start assets add per command reorganization
