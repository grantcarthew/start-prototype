# DR-024: Doctor Exit Code System

- Date: 2025-01-06
- Status: Accepted
- Category: CLI Design

## Problem

The `start doctor` command performs system health checks and reports issues with configuration, assets, agents, and environment. The tool needs an exit code strategy that:

- Indicates whether issues were found for scripting and automation
- Works in CI/CD pipelines and health check scripts
- Balances simplicity with informativeness
- Remains appropriate for a user-facing CLI tool (not a monitoring system)
- Handles both critical errors and informational warnings
- Provides clear feedback without complex exit code interpretation

## Decision

`start doctor` uses simple binary exit codes: `0` for healthy (no issues), `1` for any issues found (errors or warnings).

Exit codes do not distinguish between error severity and warning severity. Output categorizes issues by severity (errors vs warnings), but the exit code remains binary.

## Why

Simple binary exit codes match user expectations for CLI tools:

- Users need to know "is there a problem?" not "what kind of problem?"
- Output provides full details about severity and specific issues
- Script writers can use simple `if [ $? -eq 0 ]` checks
- No need to memorize exit code meanings (3 = config error, 4 = agent error, etc.)
- Appropriate for user-facing tools, not monitoring systems

Binary approach simplifies automation:

- CI/CD scripts: health check passes or fails
- Pre-flight checks: proceed or investigate
- No complex exit code interpretation needed
- Works with standard shell error handling

All issues result in exit code 1 because:

- Even warnings indicate suboptimal state
- Users should investigate warnings before critical work
- Outdated CLI or assets may cause unexpected behavior
- Better to be conservative (flag issues) than permissive (ignore warnings)

Output categorization provides severity information:

- Errors clearly marked as critical (prevent operation)
- Warnings marked as informational (suboptimal but functional)
- Overall status message reflects worst issue type
- Users get severity details from output, not exit code

## Trade-offs

Accept:

- Cannot distinguish warning vs error from exit code alone (must read output)
- Cannot detect specific issue type from exit code (config vs agent vs asset)
- Scripts cannot differentiate "update available" from "config broken"
- All issues treated equally in exit code (warning = error = 1)

Gain:

- Extremely simple exit code interpretation (0 or 1)
- Works with standard shell error handling and automation
- No need to memorize or document multiple exit codes
- Output provides full severity and category information
- Appropriate simplicity for user-facing CLI tool
- Easy to use in scripts and CI/CD pipelines
- Consistent with common CLI tool behavior

## Alternatives

Severity-based exit codes (0, 1, 2):

- Exit 0: Healthy (no issues)
- Exit 1: Warnings only (update available, optional env vars missing)
- Exit 2: Errors (config broken, agent missing, critical issues)

Pros:
- Scripts could differentiate warnings from errors
- Could auto-proceed on warnings, stop on errors
- More nuanced status reporting

Cons:
- More complex for users to understand
- Scripts must handle three cases instead of two
- Boundary between "warning" and "error" can be fuzzy
- Users still need to read output for details
- Overkill for user-facing tool

Rejected: Complexity outweighs benefits for a diagnostic command. Users always read output anyway.

Category-based exit codes (0-7):

- Exit 0: Healthy
- Exit 1: Configuration issues
- Exit 2: Asset issues
- Exit 3: Agent issues
- Exit 4: Environment issues
- Exit 5: Multiple categories
- Exit 6+: Reserved

Pros:
- Scripts could react differently to each issue type
- Specific automation based on category
- Very detailed status encoding

Cons:
- Complex to document and remember
- Scripts rarely need category-specific handling
- Users must check multiple exit codes
- Overkill for diagnostic tool
- Better suited for monitoring systems, not CLI tools

Rejected: Far too complex for user-facing diagnostic command. Category information available in output.

Always exit 0 (never fail):

Pros:
- Never breaks scripts or automation
- Doctor always runs successfully

Cons:
- Cannot use in CI/CD health checks
- Scripts cannot detect issues programmatically
- Must parse output text (fragile)
- Users may ignore doctor output if it never "fails"

Rejected: Makes doctor useless for automation and health checks.

## Structure

Exit code values:

- 0: Everything healthy (no issues found)
- 1: One or more issues found (warnings or errors)

Issue categories for output organization:

Errors (critical - prevent operation):

- Configuration issues (invalid TOML, missing required fields, file reference errors)
- Asset issues (cache empty, catalog unreachable when needed)
- Agent issues (binary not found, not executable)
- Environment issues (required environment variables missing)

Warnings (non-critical - informational):

- Update availability (CLI update available, asset updates available)
- Configuration warnings (deprecated fields, optional issues)
- Environment warnings (GH_TOKEN not set, EDITOR not set)

Priority rules when multiple issues found:

1. Display all issues (both errors and warnings)
2. Categorize by severity (group errors together, warnings together)
3. Exit code is 1 for any issue (error or warning)
4. Overall status message reflects worst severity (errors take precedence)

Check definitions:

Version Information:
- Source: Version injection from build process
- Always succeeds (informational only)
- Displays: CLI version, commit hash, build date, Go version

Configuration:
- Check: TOML syntax validation for each file (settings.toml, agents.toml, tasks.toml, roles.toml, contexts.toml)
- Check: File references exist
- Check: Required fields present
- Error if: Invalid syntax, missing required fields
- Warning if: Deprecated fields used, optional issues found

Asset Information:
- Check: Asset cache directory exists and accessible
- Check: Catalog connectivity (GitHub API reachable)
- Check: Compare local cache with latest catalog index
- Error if: Cache corrupted or inaccessible
- Warning if: Updates available in catalog

CLI Version Check:
- Check: Latest release from GitHub
- Warning if: Update available
- Never errors (informational only)

Agents:
- Check: Each configured agent binary is discoverable in PATH
- Check: Binary is executable
- Error if: Binary not found or not executable
- Never warns (binary either works or doesn't)

Environment:
- Check: Required environment variables (if specified in config)
- Check: Optional environment variables (EDITOR, GH_TOKEN)
- Error if: Required env var missing
- Warning if: Optional env var missing

## Usage Examples

All healthy output:

```bash
$ start doctor

Version Information:
  CLI Version:     1.3.0
  Commit:          abc1234
  Build Date:      2025-01-06T10:30:00Z
  Go Version:      go1.22.0

Configuration:
  Global Config:
    settings.toml: ✓ ~/.config/start/settings.toml
    agents.toml:   ✓ ~/.config/start/agents.toml
    tasks.toml:    ✓ ~/.config/start/tasks.toml
    roles.toml:    ✓ ~/.config/start/roles.toml
    contexts.toml: ✓ ~/.config/start/contexts.toml
  Local Config:
    settings.toml: ✓ .start/settings.toml
    agents.toml:   ✓ .start/agents.toml
    tasks.toml:    ✓ .start/tasks.toml
    roles.toml:    ✓ .start/roles.toml
    contexts.toml: ✓ .start/contexts.toml
  Validation:      ✓ Valid

Asset Information:
  Cache Status:    ✓ Initialized
  Latest Check:    2 hours ago
  Updates:         ✓ No updates available

CLI Version Check:
  Current:         v1.3.0
  Latest Release:  v1.3.0
  Status:          ✓ Up to date

Agents:
  claude:          ✓ Binary found at /usr/local/bin/claude
  aider:           ✓ Binary found at ~/.local/bin/aider

Environment:
  GH_TOKEN:        ✓ Set
  EDITOR:          ✓ Set (vim)

Overall Status:   ✓ Healthy
```

```bash
$ echo $?
0
```

Warnings only output:

```bash
$ start doctor

# ... other checks ...

Asset Information:
  Cache Status:    ✓ Initialized
  Latest Check:    45 days ago
  Updates:         ⚠ Updates available
  Action:          Run 'start assets update' to refresh

CLI Version Check:
  Current:         v1.2.3
  Latest Release:  v1.3.0
  Status:          ⚠ Update available
  Update via:      brew upgrade grantcarthew/tap/start

Environment:
  GH_TOKEN:        ⚠ Not set (API rate limits apply)
  EDITOR:          ✓ Set (vim)

Overall Status:   ⚠ Updates available
```

```bash
$ echo $?
1
```

Errors present output:

```bash
$ start doctor

# ... other checks ...

Configuration:
  Global Config:
    settings.toml: ✓ ~/.config/start/settings.toml
    agents.toml:   ✓ ~/.config/start/agents.toml
    tasks.toml:    ✓ ~/.config/start/tasks.toml
    roles.toml:    ✓ ~/.config/start/roles.toml
    contexts.toml: ✓ ~/.config/start/contexts.toml
  Local Config:
    settings.toml: ✗ Invalid TOML syntax at line 15
    agents.toml:   - Not found
    tasks.toml:    - Not found
    roles.toml:    - Not found
    contexts.toml: - Not found
  Validation:      ✗ Failed

Asset Information:
  Cache Status:    ✗ Not initialized
  Action:          Assets will download on-demand or run 'start init'

Agents:
  claude:          ✗ Binary not found
  aider:           ✓ Binary found at ~/.local/bin/aider

Overall Status:   ✗ Critical issues found
```

```bash
$ echo $?
1
```

Mixed errors and warnings output:

```bash
$ start doctor

# ... checks ...

Configuration:
  Global Config:
    settings.toml: ✓ ~/.config/start/settings.toml
    agents.toml:   ✓ ~/.config/start/agents.toml
    tasks.toml:    ✓ ~/.config/start/tasks.toml
    roles.toml:    ✓ ~/.config/start/roles.toml
    contexts.toml: ✓ ~/.config/start/contexts.toml
  Local Config:
    settings.toml: ✗ Invalid TOML syntax at line 15
    agents.toml:   - Not found
    tasks.toml:    - Not found
    roles.toml:    - Not found
    contexts.toml: - Not found
  Validation:      ✗ Failed

Asset Information:
  Cache Status:    ✓ Initialized
  Latest Check:    45 days ago
  Updates:         ⚠ Updates available (cannot use until config fixed)

Agents:
  claude:          ✗ Binary not found
  aider:           ✓ Binary found at ~/.local/bin/aider

Issues Found:

Errors (2):
  • Local settings.toml has invalid TOML syntax
  • Agent 'claude' binary not found

Warnings (1):
  • Asset updates available

Overall Status:   ✗ Critical issues found
```

```bash
$ echo $?
1
```

Note: When both errors and warnings exist, overall status shows "Critical issues found" (errors take precedence in messaging).

Basic health check script:

```bash
#!/bin/bash
if start doctor > /dev/null 2>&1; then
    echo "System healthy"
else
    echo "Issues found, run 'start doctor' for details"
fi
```

Pre-flight check script:

```bash
#!/bin/bash
# Run before starting work session
start doctor
if [ $? -ne 0 ]; then
    read -p "Issues found. Continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi
```

Automation script:

```bash
#!/bin/bash
# CI/automation script
if ! start doctor --quiet; then
    echo "::error::Health check failed"
    exit 1
fi
```

Quiet mode usage:

```bash
$ start doctor --quiet
# Only shows issues, no verbose output
# Exit code still 0 or 1

# If healthy, no output at all
$ start doctor --quiet
$ echo $?
0
```

## Execution Flow

Check execution order (user-facing priority):

1. Version Information (always shown first, informational)
2. Configuration (critical - affects everything else)
3. Asset Information (needed for tasks and roles)
4. CLI Version Check (informational)
5. Agents (needed for execution)
6. Environment (helpful context)

Error and warning accumulation:

- Run all checks regardless of failures (don't short-circuit)
- Collect errors in one list, warnings in another
- Display both lists in output with clear categorization
- Calculate exit code: 0 if both lists empty, 1 if either has items
- Generate overall status message based on worst issue type

Quiet mode behavior:

- Suppress verbose output (only show issues)
- If healthy: no output at all
- Exit code remains 0 or 1 regardless of quiet mode
- Use for scripting and automation

## Updates

- 2025-01-17: Updated Asset Information check to reflect catalog-based cache system (DR-031), removed asset-version.toml references
