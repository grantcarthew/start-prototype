# start update

## Name

start update - Update asset library from GitHub

## Synopsis

```bash
start update
start update [flags]
```

## Description

Downloads the latest asset library from GitHub repository, replacing the local asset directory (`~/.config/start/assets/`). Does not modify user configuration files - only updates the shared asset library.

**What gets updated:**

- **Agents** - Agent configuration templates
- **Roles** - System prompt markdown files
- **Tasks** - Default task configurations
- **Examples** - Configuration examples

**What doesn't change:**

- User's `config.toml` (global or local)
- Custom files outside asset directory
- Installed agent binaries

**Update behavior:**

- Fetches from GitHub repository default branch
- Overwrites existing asset directory completely
- Updates `asset-version.toml` file with timestamp and commit SHA
- Reports what changed (new, updated, removed files)

**Use cases:**

- Get latest agent configurations
- Update role templates with improvements
- Get new default tasks
- Refresh after repository updates

**Safety:**

- Non-destructive: Only updates asset library
- User config untouched
- Can re-run safely anytime
- No backup needed (assets are from repo)

## Flags

**--verbose**, **-v**
: Show detailed download progress and file-by-file changes.

**--quiet**, **-q**
: Suppress progress output, show only summary and errors.

**--help**, **-h**
: Show help for this command.

## Behavior

### Normal Update

```bash
start update
```

**Process:**

1. Check network connectivity
2. Fetch asset manifest from GitHub
3. Download all asset files
4. Overwrite `~/.config/start/assets/`
5. Update `asset-version.toml` file
6. Report changes

**Output:**

```
Updating asset library...
═══════════════════════════════════════════════════════════

Fetching from GitHub:
  Repository: github.com/grantcarthew/start
  Branch: main
  Path: /assets

Downloading assets:
  agents/     [████████████████████████████████] 8 files
  roles/      [████████████████████████████████] 12 files
  tasks/      [████████████████████████████████] 6 files
  examples/   [████████████████████████████████] 2 files

Installed to: ~/.config/start/assets/

What's new:
  + 2 new agents: openai, deepseek
  ~ 3 updated roles: code-reviewer, security-reviewer, doc-reviewer
  + 1 new task: security-review
  - 1 removed: legacy-task

Asset version: abc123def456 (2025-01-06T10:30:00Z)

✓ Update complete

To use updated assets:
  - New agents: start agent add
  - Updated roles: Already in use if referenced
  - New tasks: Available in 'start task' list

Run 'start doctor' to verify.
```

### Verbose Mode

```bash
start update --verbose
```

Shows file-by-file progress:

```
Updating asset library...
═══════════════════════════════════════════════════════════

Connecting to GitHub...
  API endpoint: https://api.github.com/repos/grantcarthew/start
  Checking repository: github.com/grantcarthew/start
  ✓ Repository accessible
  Latest commit: abc123def456

Fetching asset manifest...
  ✓ Found 28 files to download

Downloading agents/ (8 files):
  ✓ claude.toml (2.3 KB)
  ✓ gemini.toml (1.8 KB)
  ✓ aichat.toml (2.1 KB)
  + openai.toml (2.5 KB) - NEW
  + deepseek.toml (1.9 KB) - NEW
  ✓ anthropic.toml (2.2 KB)
  ✓ google.toml (1.7 KB)
  ✓ local.toml (1.2 KB)

Downloading roles/ (12 files):
  ~ code-reviewer.md (4.2 KB) - UPDATED
  ✓ doc-reviewer.md (3.1 KB)
  ~ security-reviewer.md (5.8 KB) - UPDATED
  ✓ architect.md (3.7 KB)
  ~ bug-fixer.md (2.9 KB) - UPDATED
  [... 7 more files ...]

Downloading tasks/ (6 files):
  ✓ code-review.toml (892 bytes)
  ✓ git-diff-review.toml (1.1 KB)
  ✓ comment-tidy.toml (745 bytes)
  ✓ doc-review.toml (823 bytes)
  + security-review.toml (1.3 KB) - NEW
  ✓ refactor-review.toml (967 bytes)

Downloading examples/ (2 files):
  ✓ global-config.toml (2.8 KB)
  ✓ local-config.toml (1.2 KB)

Installing assets:
  Target: /Users/grant/.config/start/assets/
  ✓ Removed old assets
  ✓ Created directory structure
  ✓ Installed 28 files (45.2 KB total)

Updating version file:
  ✓ Written: /Users/grant/.config/start/asset-version.toml
  Commit: abc123def456
  Timestamp: 2025-01-06T10:30:00Z

Summary:
  Total files: 28
  New: 3 files
  Updated: 5 files
  Unchanged: 19 files
  Removed: 1 file

✓ Update complete
```

### Quiet Mode

```bash
start update --quiet
```

Minimal output:

```
Updating assets from github.com/grantcarthew/start...
✓ Update complete (28 files, 3 new, 5 updated)
```

### First Time (No Existing Assets)

```bash
start update
```

When `~/.config/start/assets/` doesn't exist:

```
Updating asset library...
═══════════════════════════════════════════════════════════

No asset library found. This will download initial assets.

Fetching from GitHub:
  Repository: github.com/grantcarthew/start
  Branch: main

Downloading assets:
  agents/     [████████████████████████████████] 8 files
  roles/      [████████████████████████████████] 12 files
  tasks/      [████████████████████████████████] 6 files
  examples/   [████████████████████████████████] 2 files

Installed to: ~/.config/start/assets/

Initial assets installed:
  + 8 agents
  + 12 roles
  + 6 tasks
  + 2 examples

✓ Update complete

Next steps:
  - Add agents: start agent add
  - View tasks: start task
  - Run diagnostics: start doctor
```

## Output Details

### Change Detection

**New files (+):**
```
+ openai.toml - NEW
```

File didn't exist in local assets, downloaded from repo.

**Updated files (~):**
```
~ code-reviewer.md - UPDATED
```

File existed locally, content changed in repo.

**Unchanged files (✓):**
```
✓ gemini.toml
```

File exists locally and in repo, content identical (no download).

**Removed files (-):**
```
- legacy-task.toml - REMOVED
```

File existed locally but not in repo, deleted from local assets.

### Progress Indicators

**Download progress:**
```
agents/ [████████████████████████████████] 8/8 files
```

Shows progress bar and file count.

**Verbose progress:**
```
Downloading agents/ (8 files):
  ✓ claude.toml (2.3 KB)
  ✓ gemini.toml (1.8 KB)
  ...
```

Shows each file name and size as downloaded.

## Exit Codes

**0** - Success (assets updated)

**1** - Network error
- Cannot reach GitHub
- Repository not found
- API rate limit exceeded

**2** - File system error
- Cannot create asset directory
- Permission denied writing files
- Disk full

**3** - Invalid repository
- Asset structure incorrect
- Manifest parsing failed
- Missing required directories

**4** - Partial failure
- Some files downloaded, some failed
- Asset directory in inconsistent state
- Version file not updated

## Error Handling

### Network Errors

**Cannot reach GitHub:**

```
Error: Cannot connect to GitHub

  Network error: dial tcp: no route to host

Check your internet connection and try again.
```

Exit code: 1

**Repository not found:**

```
Error: Repository not found

  URL: github.com/grantcarthew/start
  HTTP 404: Not Found

Check repository configuration.
```

Exit code: 1

**Rate limit exceeded:**

```
Error: GitHub API rate limit exceeded

  Limit: 60 requests/hour (anonymous)
  Reset: 2025-01-06 11:30:00 (in 45 minutes)

Try again after rate limit resets.
Or authenticate with GH_TOKEN for higher limits.
```

Exit code: 1

### File System Errors

**Permission denied:**

```
Error: Cannot write to asset directory

  Path: ~/.config/start/assets/
  Error: permission denied

Check directory permissions:
  chmod 755 ~/.config/start
```

Exit code: 2

**Disk full:**

```
Error: Insufficient disk space

  Required: ~500 KB
  Available: 124 KB

Free up disk space and try again.
```

Exit code: 2

### Repository Structure Errors

**Invalid asset structure:**

```
Error: Invalid repository structure

  Missing required directories:
    - /assets/agents/
    - /assets/roles/

Repository may be incompatible with this version.
```

Exit code: 3

**Manifest parsing failed:**

```
Error: Cannot parse asset manifest

  File: /assets/manifest.json
  Error: invalid JSON at line 5

Repository content may be corrupted.
```

Exit code: 3

### Partial Failures

**Some files failed:**

```
Warning: Update partially failed

Downloaded: 25/28 files
Failed:
  - agents/openai.toml (network timeout)
  - roles/new-role.md (404 not found)
  - tasks/beta-task.toml (network timeout)

Asset directory may be incomplete.
Re-run 'start update' to retry.
```

Exit code: 4

Asset version file NOT updated (prevents false "up to date").

## Asset Version File

### Format

`asset-version.toml` tracks downloaded assets:

```toml
# Asset version tracking - managed by 'start update'
# Last updated: 2025-01-06T10:30:00Z

commit = "abc123def456"
timestamp = "2025-01-06T10:30:00Z"
repository = "github.com/grantcarthew/start"
branch = "main"

[files]
"agents/claude.toml" = "a1b2c3d4e5f6..."
"agents/gemini.toml" = "e5f6g7h8i9j0..."
"roles/code-reviewer.md" = "i9j0k1l2m3n4..."
```

### Purpose

- Track when assets were last updated
- Display in `start doctor`
- Determine if assets are stale (> 30 days)
- Debugging and support

### Location

`~/.config/start/asset-version.toml`

Same directory as `config.toml`, not inside `assets/`.

## How Assets Are Used

### Agent Templates

Located in `~/.config/start/assets/agents/`:

**During `start agent add`:**

```bash
start agent add

Add new agent
─────────────────────────────────────────────────

Use a template? [Y/n]: y

Available templates:
  1) claude - Anthropic's Claude AI
  2) gemini - Google's Gemini AI
  3) openai - OpenAI GPT models (NEW)
  4) deepseek - DeepSeek coding models (NEW)
  5) Custom (manual entry)

Select [1-5]: 3

Loading template: openai.toml...
Agent name [openai]:
Description: OpenAI GPT models
[... continues with template values ...]
```

Templates pre-fill agent configuration.

### Role Files

Located in `~/.config/start/assets/roles/`:

**Referenced in config:**

```toml
[system_prompt]
file = "~/.config/start/assets/roles/code-reviewer.md"

[tasks.review]
system_prompt_file = "~/.config/start/assets/roles/security-reviewer.md"
```

**Auto-update behavior:**

When you run `start update`:
1. `code-reviewer.md` is updated in assets directory
2. Next `start` invocation reads updated file
3. **No config change needed** - reference is to file path

### Task Definitions

Located in `~/.config/start/assets/tasks/`:

**Merged into task list:**

```bash
start task

Available tasks:
  code-review (cr)        - Review code quality [default]
  git-diff-review (gdr)   - Review git diff [default]
  comment-tidy (ct)       - Tidy code comments [default]
  doc-review (dr)         - Review documentation [default]
  security-review (sr)    - Security-focused review [default, NEW]
  my-custom-task          - Custom task [user]
```

Tasks from assets marked `[default]`, user tasks marked `[user]`.

**Override behavior:**

If user defines task with same name in their config:

```toml
[tasks.code-review]  # User's version
# Overrides default code-review from assets
```

User's version takes precedence.

### Example Configs

Located in `~/.config/start/assets/examples/`:

**Used as reference:**

Users can view examples:
```bash
cat ~/.config/start/assets/examples/global-config.toml
cat ~/.config/start/assets/examples/local-config.toml
```

Copy sections into their own config.

Not automatically loaded - reference only.

## Network Requirements

### GitHub Access

**Required:**
- HTTPS access to `github.com`
- HTTPS access to `api.github.com`
- Outbound TCP 443

**Optional:**
- `GH_TOKEN` environment variable (for higher rate limits)

### Offline Behavior

If network unavailable:

```
Error: Cannot connect to GitHub

  Network error: no internet connection

Update requires network access.
Asset library not modified.
```

**Workaround:**

Assets are in GitHub repository - can manually clone:

```bash
git clone https://github.com/grantcarthew/start /tmp/start-repo
cp -r /tmp/start-repo/assets ~/.config/start/assets/
```

## Performance

**Typical update:**
- Initial download: 2-5 seconds (28 files, ~50 KB)
- Incremental update: 1-2 seconds (3-5 changed files)
- Network speed dependent

**Bandwidth:**
- Full download: ~50-100 KB
- Incremental: ~5-20 KB

**Disk space:**
- Asset directory: ~500 KB
- Per update (no cleanup needed)

## Notes

### Update Frequency

**Recommended:**
- Update every 30 days for new features
- Update when `start doctor` warns
- Update after repository announcements

**Not required:**
- Assets don't expire
- Old assets continue working
- Update is optional, not mandatory

### Breaking Changes

Asset updates are backward compatible:

- New fields added (optional)
- Old fields preserved (deprecated gradually)
- Major changes announced in release notes

If breaking change needed:
- Announced 90 days in advance
- Migration guide provided
- `start doctor` will warn

### Repository Configuration

Default repository: `github.com/grantcarthew/start`

Future: Support custom repositories:
```bash
start update --repo github.com/myorg/start-assets
```

Not implemented yet.

### Privacy

`start update` makes these network requests:
- GitHub API: Repository metadata
- GitHub raw content: Asset files

No telemetry, no tracking, no data sent.

## Examples

### Regular Update

```bash
start update
```

Download latest assets from GitHub.

### Quiet Update (CI/CD)

```bash
start update --quiet
if [ $? -ne 0 ]; then
  echo "Asset update failed"
  exit 1
fi
```

Minimal output, check exit code.

### Verbose Troubleshooting

```bash
start update --verbose
```

See exactly what's being downloaded and why.

### Check Then Update

```bash
start doctor
# See warning about outdated assets
start update
start doctor
# Verify update successful
```

## See Also

- start-doctor(1) - Diagnose installation
- start-init(1) - Initialize configuration
- start-agent(1) - Manage agents
- start-task(1) - Run tasks
