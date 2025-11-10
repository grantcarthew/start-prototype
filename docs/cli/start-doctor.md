# start doctor

## Name

start doctor - Diagnose start installation and configuration

## Synopsis

```bash
start doctor
start doctor [flags]
```

## Description

Performs comprehensive health check of `start` installation, configuration, and environment. Reports issues, warnings, and suggestions. Does not modify any files - diagnostic only.

**Health checks performed:**

- **Version** - Current version vs latest release
- **Assets** - Age and availability of asset library
- **Configuration** - TOML syntax and semantic validation
- **Agents** - Binary availability and configuration validity
- **Contexts** - Required context file existence
- **Environment** - Shell, permissions, directory structure

**Use cases:**

- Troubleshooting `start` issues
- Verifying installation after setup
- Checking for outdated assets
- Quick health check before important work
- CI/CD validation (with `--quiet` flag)

## Flags

This command supports the standard global flags for controlling output verbosity and showing help: `--verbose`, `--quiet`, and `--help`. See `start --help` for more details.

## Behavior

Runs all diagnostic checks in order and reports results. Non-blocking - shows all issues before exiting.

**Check order:**

1. Version check
2. Asset library check
3. Configuration validation
4. Agent diagnostics
5. Role validation
6. Task validation (including agent and role references)
7. Context verification
8. Environment check

**Exit codes:**

- 0 - All checks passed (no errors)
- 1 - Configuration errors (broken config)
- 2 - Missing dependencies (agents not installed)
- 3 - Asset issues (outdated or missing)
- 4 - Multiple issues (combination of above)

## Output

### Normal Mode

```bash
start doctor
```

Output:

```
Diagnosing start installation...
═══════════════════════════════════════════════════════════

Version
  start v1.2.3
  ✓ Latest version (checked 2025-01-06)

Assets
  ⚠ Asset library is 45 days old
  Last updated: 2024-11-22

  Run 'start update' to download latest assets.

Configuration
  ✓ Global config valid (~/.config/start/config.toml)
  ✓ Local config valid (./.start/config.toml)
  ✓ No TOML syntax errors
  ✓ All required fields present
  ✓ Merge behavior correct

Agents (3 configured)
  ✓ claude
    Binary: /usr/local/bin/claude
    Models: 3 configured
    Default: sonnet

  ✓ gemini
    Binary: /usr/local/bin/gemini
    Models: 2 configured
    Default: flash

  ✗ aichat
    Binary: NOT FOUND
    Install: brew install aichat
    Or remove: start agent remove aichat

Tasks (4 configured)
  ✓ code-review (cr)
    Agent: claude (default)

  ✓ git-diff-review (gdr)
    Agent: claude (default)

  ✗ go-review (gor)
    Agent: go-expert (NOT FOUND)
    Fix: start config agent add go-expert

  ✓ security-audit (sec)
    Agent: claude (from task config)

Contexts (2 required, 1 optional)
  Required:
    ✓ environment - ~/reference/ENVIRONMENT.md
    ✓ index - ~/reference/INDEX.csv

  Optional:
    ⚠ project - ./PROJECT.md (not found)

Environment
  ✓ Shell: /bin/bash
  ✓ Config directory: ~/.config/start/ (writable)
  ✓ Working directory: /Users/grant/Projects/myapp

Summary
───────────────────────────────────────────────────────────
  2 errors, 2 warnings found

Issues:
  ✗ Agent 'aichat' binary not found
  ✗ Task 'go-review' references undefined agent 'go-expert'
  ⚠ Assets outdated (45 days old)
  ⚠ Optional context 'project' missing

Recommendations:
  1. Install aichat: brew install aichat
  2. Add go-expert agent: start config agent add go-expert
  3. Update assets: start update
  4. Optional: Create PROJECT.md in current directory
```

### Verbose Mode

```bash
start doctor --verbose
```

Adds detailed information:

```
Diagnosing start installation...
═══════════════════════════════════════════════════════════

Version
  Binary: /usr/local/bin/start
  Version: v1.2.3
  Built: 2024-12-15T10:30:00Z
  Go version: go1.22.0
  Platform: darwin/arm64

  Latest release check:
    URL: https://api.github.com/repos/grantcarthew/start/releases/latest
    Latest: v1.2.3
    ✓ Up to date

Assets
  Location: ~/.config/start/assets/
  Version file: ~/.config/start/asset-version.toml
  Last updated: 2024-11-22T14:23:10Z (45 days ago)
  Commit: abc123def456

  Asset inventory:
    agents/ - 8 files
    roles/ - 12 files
    tasks/ - 6 files
    examples/ - 2 files

  ⚠ Assets are 45 days old (recommended: < 30 days)
  Run 'start update' to download latest assets.

Configuration
  Global config:
    Path: /Users/grant/.config/start/config.toml
    Size: 2.3 KB
    Modified: 2025-01-05T09:15:32Z
    ✓ TOML syntax valid
    ✓ Schema valid

  Local config:
    Path: /Users/grant/Projects/myapp/.start/config.toml
    Size: 456 bytes
    Modified: 2025-01-06T10:12:05Z
    ✓ TOML syntax valid
    ✓ Schema valid

  Merged configuration:
    ✓ No section conflicts
    ✓ Settings merge: 3 from global, 1 from local override
    ✓ Contexts: 2 global + 1 local = 3 total
    ✓ Agents: 3 (global only - correct)

[... detailed agent, context, environment checks ...]
```

### Quiet Mode (CI/CD)

```bash
start doctor --quiet
echo $?
```

No output on success (exit code 0).

Output on failure:

```
Error: Agent 'aichat' binary not found
Warning: Assets outdated (45 days old)
```

Exit code 4 (multiple issues).

### All Healthy

```bash
start doctor
```

Output:

```
Diagnosing start installation...
═══════════════════════════════════════════════════════════

Version
  start v1.2.3
  ✓ Latest version

Assets
  ✓ Asset library up to date (updated 2 days ago)

Configuration
  ✓ Global config valid
  ✓ Local config valid

Agents (3 configured)
  ✓ claude - /usr/local/bin/claude
  ✓ gemini - /usr/local/bin/gemini
  ✓ aichat - /usr/local/bin/aichat

Contexts (2 required)
  ✓ environment - ~/reference/ENVIRONMENT.md
  ✓ index - ~/reference/INDEX.csv

Environment
  ✓ All checks passed

Summary
───────────────────────────────────────────────────────────
  ✓ No issues found

Everything looks good!
```

## Diagnostic Details

### Version Check

Compares local version against latest GitHub release.

**What it checks:**
- Binary version from build metadata
- Latest release from GitHub API (with caching)
- Update recommendation if behind

**Output examples:**

```
✓ Latest version (v1.2.3)
⚠ Update available (v1.2.3 → v1.3.0)
  Update: go install github.com/grantcarthew/start@latest
✗ Very outdated (v1.0.5 → v1.3.0, 6 months behind)
```

### Asset Library Check

Verifies age and completeness of asset library.

**What it checks:**
- `asset-version.toml` exists and is readable
- Asset age (warns if > 30 days)
- Asset directory structure intact
- All expected subdirectories present

**Output examples:**

```
✓ Asset library up to date (updated 2 days ago)
⚠ Assets are 45 days old
  Run 'start update' to refresh
✗ Asset library missing
  Run 'start init' to initialize
```

### Configuration Validation

Comprehensive config validation.

**What it checks:**
- TOML syntax (parse errors)
- Required fields present
- Field types correct
- Agent command templates valid
- Context paths resolvable
- Merge behavior (global + local)

Reuses validation from `start config validate`.

### Agent Diagnostics

Tests all configured agents.

**What it checks:**
- Binary availability (`exec.LookPath`)
- Command template syntax
- Model configuration
- Default model set

Similar to `start agent test` but for all agents.

**Output example:**

```
✓ claude
  Binary: /usr/local/bin/claude
  Models: 3 configured
  Default: sonnet

✗ gemini
  Binary: NOT FOUND
  Install: https://github.com/example/gemini-cli
```

### Context Verification

Checks context document files exist.

**What it checks:**
- Required contexts exist (error if missing)
- Optional contexts exist (warning if missing)
- File permissions readable
- Path resolution (~ expansion, relative paths)

**Output example:**

```
Required:
  ✓ environment - ~/reference/ENVIRONMENT.md
  ✗ index - ~/reference/INDEX.csv (not found)

Optional:
  ⚠ project - ./PROJECT.md (not found)
  ✓ agents - ./AGENTS.md
```

### Environment Check

Verifies runtime environment.

**What it checks:**
- Shell availability (from config or auto-detect)
- Config directory writable
- Asset directory exists and writable
- Working directory accessible

**Output example:**

```
✓ Shell: /bin/bash
✓ Config directory: ~/.config/start/ (writable)
✓ Asset directory: ~/.config/start/assets/ (writable)
✓ Working directory: /Users/grant/Projects/myapp
```

## Exit Codes

**0** - Healthy (all checks passed, warnings OK)

**1** - Configuration errors
- TOML syntax errors
- Missing required fields
- Invalid configuration values

**2** - Missing dependencies
- Agent binaries not found
- Required contexts missing
- Shell not available

**3** - Asset issues
- Asset library missing
- Asset version file corrupted
- Asset directory not writable

**4** - Multiple issues (combination of above)

**In quiet mode:** Exit code indicates severity. Check output for details.

## Examples

### Basic Health Check

```bash
start doctor
```

Run all diagnostics and show results.

### CI/CD Validation

```bash
start doctor --quiet
if [ $? -ne 0 ]; then
  echo "start configuration has issues"
  exit 1
fi
```

Exit code only, no output on success.

### Detailed Diagnostics

```bash
start doctor --verbose
```

Show all details for troubleshooting.

### Check Before Important Work

```bash
start doctor && start task code-review
```

Only run task if environment is healthy.

## Notes

### Automatic Checks

`start doctor` runs automatically in limited form:

- Asset age check runs on every `start` invocation
- If assets > 30 days: "⚠ Assets outdated. Run 'start update'"
- Non-blocking: warning only, doesn't prevent execution
- Cached for 24 hours to avoid GitHub API rate limits

Full `start doctor` command runs all checks manually.

### Asset Age Thresholds

- **< 30 days** - OK
- **30-90 days** - Warning
- **> 90 days** - Strong warning

Assets don't expire, but recommendations:
- Update every 30 days for new features
- Critical fixes announced in release notes

### GitHub API Rate Limiting

Version check uses GitHub API:
- Anonymous: 60 requests/hour
- Cached for 24 hours per machine
- Cache location: `~/.config/start/.version-check-cache`
- On rate limit: Shows cached version or skips check

### Privacy

`start doctor` makes these network requests:
- GitHub API: Check latest release version
- No telemetry, no tracking, no data sent

Can run offline (skips version check).

### Performance

Doctor runs in < 1 second typically:
- Version check: 50-200ms (cached) or 200-500ms (fresh)
- Config validation: 10-50ms
- Agent checks: 10ms per agent
- File checks: 5ms per file

Total: ~100-800ms for typical installation.

## See Also

- start-update(1) - Update asset library
- start-config(1) - Configuration management
- start-agent(1) - Agent management
- start-init(1) - Initialize configuration
