# start assets browse

## Name

start assets browse - Open GitHub asset catalog in web browser

## Synopsis

```bash
start assets browse
```

## Description

Opens the GitHub asset catalog in your default web browser for visual exploration and discovery. This provides a graphical interface for browsing available roles, tasks, agents, and contexts organized by category.

The command opens the catalog's `assets/` directory in GitHub's web interface, allowing you to:

- Browse assets by type and category
- Read file contents directly
- View commit history and recent changes
- Access metadata files (`.meta.toml`)
- Explore the catalog structure visually

Uses the `[settings] asset_repo` configuration value (default: `grantcarthew/start`) to determine which repository to open.

**Workflow pattern:**

1. Browse catalog visually in browser
2. Identify assets of interest
3. Return to terminal and install with `start assets add <query>`

## Behavior

**URL format:**

```
https://github.com/{asset_repo}/tree/main/assets
```

Where `{asset_repo}` comes from `[settings] asset_repo` in `config.toml`.

**Default repository:**

- Repository: `grantcarthew/start`
- Opens: `https://github.com/grantcarthew/start/tree/main/assets`

**Custom repository:**
If you've configured a custom asset repository:

```toml
[settings]
asset_repo = "myorg/custom-assets"
```

Opens: `https://github.com/myorg/custom-assets/tree/main/assets`

**Browser detection:**
The command uses platform-specific defaults to open URLs:

- macOS: `open <url>`
- Linux: `xdg-open <url>`

**If browser fails to open:**

```
Error: Could not open browser

URL: https://github.com/grantcarthew/start/tree/main/assets

Copy and paste this URL into your browser to view the catalog.
```

The command exits with code 0 even if the browser fails to open, since the URL is displayed.

## Output

**Successful browser launch:**

```bash
$ start assets browse

Opening GitHub catalog in browser...
✓ https://github.com/grantcarthew/start/tree/main/assets
```

Browser opens to the catalog.

**Browser launch failed (fallback):**

```bash
$ start assets browse

Opening GitHub catalog in browser...

⚠ Could not open browser automatically

URL: https://github.com/grantcarthew/start/tree/main/assets

Copy and paste this URL into your browser to view the catalog.
```

URL is displayed for manual access.

**Custom repository:**

```bash
$ start assets browse

Opening GitHub catalog in browser...
✓ https://github.com/myorg/custom-assets/tree/main/assets
```

## Catalog Structure

The catalog is organized by type and category:

```
assets/
├── agents/
│   ├── anthropic/
│   │   └── claude.toml
│   └── google/
│       └── gemini.toml
├── contexts/
│   └── project/
│       └── project-summary.toml
├── roles/
│   ├── general/
│   │   ├── code-reviewer.md
│   │   └── code-reviewer.meta.toml
│   └── languages/
│       └── go-expert.md
└── tasks/
    ├── git-workflow/
    │   ├── commit-message.toml
    │   ├── pre-commit-review.toml
    │   └── pre-commit-review.md
    └── quality/
        └── code-review.toml
```

**File types:**

- `.toml` - Asset configuration files
- `.md` - Documentation or role content
- `.meta.toml` - Asset metadata (description, tags, etc.)

## Exit Codes

**0** - Success (browser opened or URL displayed)

## Use Cases

### Visual Discovery

**Problem:** Want to see all available assets graphically.

```bash
start assets browse
```

Opens GitHub web interface where you can:

- Click through categories
- Read asset descriptions
- View file contents
- Explore related assets

### First-Time Exploration

**Problem:** New to the catalog, want to understand what's available.

```bash
start assets browse
```

Browse the catalog structure to understand:

- What types of assets exist
- How assets are categorized
- Naming conventions
- Asset documentation style

### Asset Documentation

**Problem:** Want to read detailed documentation for an asset.

```bash
start assets browse
```

In browser:

1. Navigate to asset of interest
2. Read `.md` files for detailed documentation
3. View `.meta.toml` for metadata
4. Check commit history for changes

### Catalog Development

**Problem:** Contributing to the catalog, want to see current structure.

```bash
start assets browse
```

Review catalog structure before adding new assets to ensure consistency.

## Comparison with Other Commands

### vs `start assets search`

**`start assets browse`** - Visual exploration in browser

```bash
start assets browse
# Opens browser, graphical navigation
```

**`start assets search`** - Terminal-based keyword search

```bash
start assets search "commit"
# Prints matching assets in terminal
```

Use browse for exploration, search for finding specific assets.

### vs `start assets add`

**`start assets browse`** - View only, no installation

```bash
start assets browse
# Just opens browser, doesn't install anything
```

**`start assets add`** - Search and install

```bash
start assets add "pre-commit"
# Downloads and installs asset
```

Typical workflow:

1. `start assets browse` - Find interesting assets
2. Note asset names
3. `start assets add <name>` - Install them

## Configuration

**Setting the asset repository:**

In `~/.config/start/config.toml` or `./.start/config.toml`:

```toml
[settings]
asset_repo = "grantcarthew/start"    # Default
# asset_repo = "myorg/custom-assets"  # Custom repository
```

**No other configuration needed** - The command uses platform defaults for opening browsers.

## Examples

### Open Default Catalog

```bash
$ start assets browse

Opening GitHub catalog in browser...
✓ https://github.com/grantcarthew/start/tree/main/assets
```

Browser opens to the default catalog.

### With Custom Repository

```toml
# In config.toml
[settings]
asset_repo = "acme-corp/internal-assets"
```

```bash
$ start assets browse

Opening GitHub catalog in browser...
✓ https://github.com/acme-corp/internal-assets/tree/main/assets
```

Browser opens to your organization's private catalog.

### Browser Launch Failure

```bash
$ start assets browse

Opening GitHub catalog in browser...

⚠ Could not open browser automatically
  Error: exec: "xdg-open": executable file not found in $PATH

URL: https://github.com/grantcarthew/start/tree/main/assets

Copy and paste this URL into your browser to view the catalog.
```

URL is displayed for manual access. Exit code is still 0.

### Workflow Example

```bash
# 1. Browse catalog to discover assets
$ start assets browse

# (In browser: explore git-workflow tasks, find "pre-commit-review")

# 2. Return to terminal and install
$ start assets add "pre-commit-review"

Found 1 match:
  tasks/git-workflow/pre-commit-review

Downloading...
✓ Cached to ~/.config/start/assets/tasks/git-workflow/
✓ Added to global config

# 3. Use the installed asset
$ start task pre-commit-review
```

## Notes

### GitHub-Only Browsing

This command only works with GitHub-hosted catalogs:

- Requires GitHub repository structure
- Opens GitHub web interface
- No support for local file:// URLs

**For local catalogs:**

- Not supported via `start assets browse`
- Use file manager or `ls -R assets/` instead

### Network Required

The command constructs a URL but doesn't verify network connectivity. The browser will show an error if:

- Internet connection is down
- GitHub is unreachable
- Repository doesn't exist

**No validation performed:**

- Command always exits 0
- Browser handles connectivity issues

### Alternative Discovery

If browser-based browsing isn't suitable:

**Terminal-based alternatives:**

```bash
start assets search "keyword"    # Search by keyword
start assets add                 # Interactive TUI browser
```

**GitHub CLI:**

```bash
gh repo view grantcarthew/start
gh browse grantcarthew/start:assets
```

### Repository Permissions

**Public repositories:**

- Work without authentication
- Anyone can browse

**Private repositories:**

- Require GitHub authentication in browser
- User must be signed in to GitHub
- User must have repository access

`start assets browse` opens the URL - GitHub handles authentication.

## See Also

- start-assets(1) - Asset management overview
- start-assets-search(1) - Terminal-based catalog search
- start-assets-add(1) - Install asset from catalog
- start-assets-info(1) - Show asset details
