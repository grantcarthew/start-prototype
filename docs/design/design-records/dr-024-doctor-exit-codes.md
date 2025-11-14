# DR-024: Doctor Exit Code System

**Date:** 2025-01-06
**Status:** Accepted
**Category:** CLI Design

## Decision

`start doctor` uses simple binary exit codes: `0` for healthy, `1` for any issues found.

## Exit Code Strategy

**Simple Binary:**
- `0` = Everything is healthy (no issues)
- `1` = One or more issues found (warnings or errors)

**Not severity-based:**
- Don't distinguish between warning severity vs error severity in exit code
- Output shows severity (errors vs warnings), but exit code is simple
- Appropriate for a user-facing tool (not a CI/monitoring tool)

## Issue Categories

Issues are categorized in output for clarity, but all result in exit code `1`:

### Errors (Critical - Prevent Operation)

**Configuration Issues:**
- Config file has invalid TOML syntax
- Config file references non-existent files
- Required context documents missing

**Asset Issues:**
- Assets not initialized (`start init` never run)
- Asset directory corrupted or empty

**Agent Issues:**
- Agent binary not found
- Agent binary not executable

**Environment Issues:**
- Required environment variables not set (if specified in config)

### Warnings (Non-Critical - Informational)

**Update Availability:**
- CLI update available
- Asset updates available

**Configuration Warnings:**
- Optional environment variables not set
- Deprecated configuration fields used

**Environment Warnings:**
- GH_TOKEN not set (API rate limits apply)
- EDITOR not set (manual editing required)

## Priority Rules

When multiple issues are found:

1. **Show all issues** - Display both errors and warnings in output
2. **Categorize by severity** - Group errors together, warnings together
3. **Exit code is 1** - Any issue (error or warning) results in exit code 1
4. **Overall status reflects worst** - If any errors, show "Critical issues found"

## Output Format

### All Healthy

```bash
$ start doctor

Version Information:
  CLI Version:     1.3.0
  Commit:          abc1234
  Build Date:      2025-01-06T10:30:00Z
  Go Version:      go1.22.0

Configuration:
  Global Config:   ✓ ~/.config/start/config.toml
  Local Config:    ✓ .start/config.toml
  Validation:      ✓ Valid

Asset Information:
  Current Commit:  abc1234 (2 days ago)
  Latest Commit:   abc1234
  Status:          ✓ Up to date

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
Exit Code:        0
```

```bash
$ echo $?
0
```

### Warnings Only

```bash
$ start doctor

# ... other checks ...

Asset Information:
  Current Commit:  abc1234 (45 days ago)
  Latest Commit:   def5678 (2 hours ago)
  Status:          ⚠ Updates available
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
Exit Code:        1
```

```bash
$ echo $?
1
```

### Errors Present

```bash
$ start doctor

# ... other checks ...

Configuration:
  Global Config:   ✓ ~/.config/start/config.toml
  Local Config:    ✗ Invalid TOML syntax at line 15
  Validation:      ✗ Failed

Asset Information:
  Status:          ✗ Not initialized
  Action:          Run 'start init' to download assets

Agents:
  claude:          ✗ Binary not found
  aider:           ✓ Binary found at ~/.local/bin/aider

Overall Status:   ✗ Critical issues found
Exit Code:        1
```

```bash
$ echo $?
1
```

### Mixed (Errors + Warnings)

```bash
$ start doctor

# ... checks ...

Configuration:
  Global Config:   ✓ ~/.config/start/config.toml
  Local Config:    ✗ Invalid TOML syntax at line 15
  Validation:      ✗ Failed

Asset Information:
  Current Commit:  abc1234 (45 days ago)
  Latest Commit:   def5678 (2 hours ago)
  Status:          ⚠ Updates available (but cannot use until config fixed)

Agents:
  claude:          ✗ Binary not found
  aider:           ✓ Binary found at ~/.local/bin/aider

Issues Found:

Errors (2):
  • Local config has invalid TOML syntax
  • Agent 'claude' binary not found

Warnings (1):
  • Asset updates available

Overall Status:   ✗ Critical issues found
Exit Code:        1
```

```bash
$ echo $?
1
```

**Note:** When both errors and warnings exist, overall status shows "Critical issues found" (errors take precedence in messaging).

## Check Definitions

### Version Information
- **Source:** `internal/version` package (DR-020)
- **Always succeeds** (informational only)

### Configuration
- **Check:** TOML syntax validation
- **Check:** File references exist
- **Check:** Required fields present
- **Error if:** Invalid syntax, missing required fields
- **Warning if:** Deprecated fields, optional issues

### Asset Information
- **Check:** `asset-version.toml` exists and valid
- **Check:** Asset files exist in `~/.config/start/assets/`
- **Check:** Latest commit from GitHub (DR-023)
- **Error if:** Not initialized
- **Warning if:** Updates available

### CLI Version Check
- **Check:** Latest release from GitHub (DR-021)
- **Warning if:** Update available
- **Never errors** (informational only)

### Agents
- **Check:** Each configured agent binary is discoverable
- **Check:** Binary is executable
- **Error if:** Binary not found or not executable
- **Never warns** (binary either works or doesn't)

### Environment
- **Check:** Required environment variables (if specified in config)
- **Check:** Optional environment variables (EDITOR, GH_TOKEN)
- **Error if:** Required env var missing
- **Warning if:** Optional env var missing

## Scripting Examples

### Basic Health Check

```bash
#!/bin/bash
if start doctor > /dev/null 2>&1; then
    echo "System healthy"
else
    echo "Issues found, run 'start doctor' for details"
fi
```

### Pre-flight Check

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

### Automation

```bash
#!/bin/bash
# CI/automation script
if ! start doctor --quiet; then
    echo "::error::Health check failed"
    exit 1
fi
```

## Implementation Notes

### Check Order

Run checks in this order (user-facing priority):

1. Version Information (always shown first)
2. Configuration (critical - affects everything)
3. Asset Information (needed for tasks/roles)
4. CLI Version Check (informational)
5. Agents (needed for execution)
6. Environment (helpful context)

### Error Accumulation

```go
type DoctorResult struct {
    Errors   []Issue
    Warnings []Issue
}

type Issue struct {
    Category    string  // "Configuration", "Assets", "Agents", etc.
    Description string
    Action      string  // Suggested fix
}

// Determine exit code
func (r *DoctorResult) ExitCode() int {
    if len(r.Errors) > 0 || len(r.Warnings) > 0 {
        return 1
    }
    return 0
}

// Overall status message
func (r *DoctorResult) Status() string {
    if len(r.Errors) > 0 {
        return "✗ Critical issues found"
    }
    if len(r.Warnings) > 0 {
        return "⚠ Updates available"
    }
    return "✓ Healthy"
}
```

### Quiet Mode

For scripting, support `--quiet` flag:

```bash
$ start doctor --quiet
# Only shows issues, no verbose output
# Exit code still 0 or 1

# Or if healthy, no output at all
$ start doctor --quiet
$ echo $?
0
```

## Benefits

- ✅ **Simple** - Binary exit code, easy to understand
- ✅ **Scriptable** - Works in automation/CI
- ✅ **Informative** - Output categorizes severity
- ✅ **User-friendly** - Appropriate for CLI tool
- ✅ **Consistent** - All issues treated equally in exit code

## Trade-offs Accepted

- ❌ Can't distinguish warning vs error in exit code (acceptable: output shows severity)
- ❌ Can't detect "which issue" from exit code alone (acceptable: run doctor to see)

## Rationale

A simple binary exit code (0 or 1) is appropriate for a user-facing tool. Users need to know "is there an issue?" and can read the output for details. Complex exit codes (0, 1, 2, 3...) are better suited for monitoring tools or CI systems, which this is not.

## Related Decisions

- [DR-021](./dr-021-github-version-check.md) - CLI version checking
- [DR-023](./dr-023-asset-staleness-check.md) - Asset staleness checking

## Implementation Checklist

- [ ] Create `internal/doctor/checker.go` with `DoctorResult` struct
- [ ] Implement each check function (config, assets, agents, env)
- [ ] Accumulate errors and warnings during checks
- [ ] Display results with severity categorization
- [ ] Return exit code based on `DoctorResult.ExitCode()`
- [ ] Add `--quiet` flag for minimal output
- [ ] Add unit tests for each check type
- [ ] Document exit codes in `start doctor --help`
