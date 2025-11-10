# DR-035: Interactive Asset Browsing

**Date:** 2025-01-10
**Status:** Accepted
**Category:** Asset Management

## Decision

Provide interactive catalog browsing via `start config <type> add` with TUI library support, falling back to numbered selection for non-interactive environments.

## What This Means

### Interactive Browsing Commands

**Browse and install assets:**
```bash
start config task add           # Interactive browser for tasks
start config role add           # Interactive browser for roles
start config agent add          # Interactive browser for agents
```

**Direct installation (bypass browsing):**
```bash
start config task add git-workflow/pre-commit-review
start config role add general/code-reviewer
start config agent add claude/sonnet
```

### User Experience Flow

**1. Invoke browser:**
```bash
$ start config task add

Fetching catalog from GitHub...
✓ Found 12 tasks across 4 categories
```

**2. Select category:**
```
Select category:
  1. git-workflow (4 tasks)
  2. code-quality (4 tasks)
  3. security (2 tasks)
  4. debugging (2 tasks)
  5. [view all tasks]

> _
```

**3. Select asset:**
```
git-workflow tasks:
  1. pre-commit-review - Review staged changes before commit
  2. pr-ready - Complete PR preparation checklist
  3. commit-message - Generate conventional commit message
  4. explain-changes - Understand what changed in commits

  [b] back  [q] quit

> _
```

**4. Confirm and download:**
```
Selected: pre-commit-review
Description: Review staged changes before commit
Tags: git, review, quality, pre-commit

Download and add to config? [Y/n] _
```

**5. Success:**
```
Downloading...
✓ Cached to ~/.config/start/assets/tasks/git-workflow/
✓ Added to global config as 'pre-commit-review'

Try it: start task pre-commit-review
```

### TUI Library Selection

**Evaluation criteria:**
- Ease of use (API simplicity)
- Dependency footprint
- Maintenance status
- Feature completeness

**Options considered:**

**1. bubbletea (charmbracelet/bubbletea)**
- Full-featured TUI framework
- Beautiful rendering
- Complex API (Elm architecture)
- Large dependency tree

**2. promptui (manifoldco/promptui)**
- Simple prompts and selects
- Minimal API
- Small dependency footprint
- Good for basic interactions

**3. survey (AlecAivazis/survey)**
- Rich prompt library
- Medium complexity
- Active maintenance
- Good balance

**4. Native numbered selection (no dependency)**
- Print numbered list
- Read user input
- No external dependencies
- Works everywhere

**Decision for v1:** Start with **native numbered selection**, evaluate TUI library later based on user feedback.

**Rationale:**
- KISS principle - don't add dependencies speculatively
- Numbered selection works everywhere (SSH, containers, CI)
- Can add TUI library in v2 if users request it
- Keeps binary small and build fast

### Implementation (Numbered Selection)

```go
func BrowseTaskCatalog() (*Asset, error) {
    // Fetch catalog
    catalog := getCatalog()
    tasks := catalog.FilterByType("tasks")

    // Group by category
    categories := tasks.GroupByCategory()

    // Show categories
    fmt.Println("\nSelect category:")
    for i, cat := range categories {
        fmt.Printf("  %d. %s (%d tasks)\n", i+1, cat.Name, len(cat.Assets))
    }
    fmt.Printf("  %d. [view all tasks]\n", len(categories)+1)

    // Read category choice
    choice := readInt(fmt.Sprintf("\nChoice [1-%d]: ", len(categories)+1))
    if choice == len(categories)+1 {
        // Show all
        return showAllTasks(tasks)
    }

    selectedCategory := categories[choice-1]

    // Show tasks in category
    fmt.Printf("\n%s tasks:\n", selectedCategory.Name)
    for i, asset := range selectedCategory.Assets {
        fmt.Printf("  %d. %s - %s\n", i+1, asset.Name, asset.Description)
    }

    // Read task choice
    taskChoice := readInt(fmt.Sprintf("\nChoice [1-%d]: ", len(selectedCategory.Assets)))
    selectedAsset := selectedCategory.Assets[taskChoice-1]

    // Confirm
    fmt.Printf("\nSelected: %s\n", selectedAsset.Name)
    fmt.Printf("Description: %s\n", selectedAsset.Description)
    fmt.Printf("Tags: %s\n", strings.Join(selectedAsset.Tags, ", "))

    if !confirm("Download and add to config?") {
        return nil, ErrUserCancelled
    }

    return selectedAsset, nil
}
```

### Filtering and Search

**Category filtering (v1):**
```bash
start config task add --category git-workflow

# Shows only git-workflow tasks, skip category selection
```

**Search by keyword (future):**
```bash
start config task add --search commit

# Shows all tasks matching "commit" in name, description, or tags
```

### Non-Interactive Mode

**For scripts/automation:**
```bash
# Direct installation (no interaction)
start config task add git-workflow/pre-commit-review --yes

# Error if not found
start config task add nonexistent/task --yes
# Exit code 1
```

**Flags:**
```
--yes, -y          Skip confirmation prompts
--category <cat>   Filter by category
--search <term>    Search by keyword (future)
```

### Error Handling

**Network unavailable:**
```
$ start config task add

Error: Cannot fetch catalog from GitHub

Network error: dial tcp: no route to host

To resolve:
- Check internet connection
- Use cached assets: start task <name>
- Add custom task: start config task add my-task
```

**No assets found:**
```
$ start config task add --category nonexistent

Error: No tasks found in category 'nonexistent'

Available categories:
  - git-workflow
  - code-quality
  - security
  - debugging

Try: start config task add --category git-workflow
```

**User cancels:**
```
$ start config task add

[... user navigates and selects task ...]

Download and add to config? [Y/n] n

Cancelled. No changes made.
```

## Alternative UX Patterns

### Pattern 1: Flat List (Simple)

```
Available tasks:
  1. git-workflow/pre-commit-review - Review staged changes
  2. git-workflow/pr-ready - Complete PR preparation
  3. code-quality/find-bugs - Find potential bugs
  4. code-quality/quick-wins - Low-hanging refactoring
  [... 8 more ...]

Choice [1-12]: _
```

**Pros:** Simplest, no category navigation
**Cons:** Long list, hard to scan

### Pattern 2: Category First (Chosen)

```
Categories:
  1. git-workflow (4 tasks)
  2. code-quality (4 tasks)

Choice: 1

Tasks in git-workflow:
  1. pre-commit-review - Review staged changes
  2. pr-ready - Complete PR preparation

Choice: 1
```

**Pros:** Organized, easier to scan
**Cons:** Extra step

### Pattern 3: Search First

```
Search for task (or press Enter to browse all): commit

Found 3 tasks:
  1. pre-commit-review - Review staged changes
  2. commit-message - Generate conventional commit
  3. explain-changes - Understand what changed

Choice: _
```

**Pros:** Fast for users who know what they want
**Cons:** Requires implementation of search

**Decision:** Use Pattern 2 (category first) for v1, can add search later.

## Implementation Phases

### Phase 1: Core Browsing (v1)
- ✅ Numbered category selection
- ✅ Asset list in category
- ✅ Confirmation prompt
- ✅ Download and add to config

### Phase 2: Enhanced Navigation (v2)
- Category filtering flag
- Back/quit navigation
- Asset preview (show content)
- Diff between local and catalog

### Phase 3: TUI Enhancement (v3)
- Evaluate user feedback
- Add TUI library if requested
- Arrow key navigation
- Fuzzy search

### Phase 4: Advanced Features (future)
- Full-text search
- Tag filtering
- Bulk installation
- Asset ratings/popularity

## Benefits

**Discoverable:**
- ✅ Browse all available assets
- ✅ Organized by category
- ✅ See descriptions before downloading

**Accessible:**
- ✅ Works in all terminals (numbered selection)
- ✅ Works over SSH
- ✅ Works in CI/containers (with --yes flag)

**Simple:**
- ✅ No external dependencies for v1
- ✅ Straightforward numbered lists
- ✅ Clear prompts and confirmations

**Flexible:**
- ✅ Interactive or direct installation
- ✅ Filter by category
- ✅ Non-interactive mode for automation

## Trade-offs Accepted

**No fancy TUI in v1:**
- ❌ No arrow key navigation
- ❌ No fuzzy search
- ❌ No multi-select
- **Mitigation:** Can add later if users request, numbered selection works well

**Category navigation adds a step:**
- ❌ Extra prompt compared to flat list
- **Mitigation:** Better organization, easier to scan large catalogs

**No search in v1:**
- ❌ Can't filter by keyword
- **Mitigation:** Category filter helps, can add search in v2

## Related Decisions

- [DR-031](./dr-031-catalog-based-assets.md) - Catalog architecture (browsing context)
- [DR-034](./dr-034-github-catalog-api.md) - GitHub API (catalog source)
- [DR-033](./dr-033-asset-resolution-algorithm.md) - Resolution (post-download behavior)

## Future Considerations

**TUI library addition:**
```go
// If user feedback requests enhanced UX
import "github.com/charmbracelet/bubbletea"

func BrowseTasks() {
    if isInteractive() && config.EnableTUI {
        return browseTUI()  // Fancy interface
    }
    return browseNumbered()  // Fallback
}
```

**Current stance:** Ship v1 with numbered selection. Monitor user feedback. Add TUI in v2 if requested.
