# DR-026: Offline Fallback and Network Unavailable Behavior

- Date: 2025-01-07
- Status: Accepted
- Category: Asset Management

## Problem

The CLI needs to handle network unavailability scenarios. The strategy must address:

- Initial setup when network is unavailable (can user run init?)
- Asset download failures (what happens when GitHub is unreachable?)
- Command execution without assets (which commands work offline?)
- Update operations when offline (should update fail or provide offline mode?)
- Doctor health checks without network (show partial results or fail?)
- Manual asset installation (support air-gapped environments or require network?)
- Security concerns with manual asset files (trust, validation)
- Implementation complexity (fallback chains, offline bundles)

## Decision

The CLI handles network unavailability gracefully with clear error messages. Network access is required only for asset download operations. No manual asset installation support.

Network-only approach:

- No manual asset installation (no --asset-path flag, no manual file dropping)
- No ZIP/tarball asset bundles for offline distribution
- Users must have network access for asset downloads
- Assets downloaded only from GitHub catalog on-demand

Network failure behavior by command:

- start init: Warns about network failure, creates config anyway, succeeds (exit 0)
- start assets update: Fails immediately with clear error (exit 1)
- start doctor: Shows partial results, skips network-dependent checks (exit based on local checks)
- start / start prompt / start task: Work without assets (warn if referenced assets missing)

## Why

Network-only approach simplifies implementation and security:

- Assets are optional content, not critical for core functionality
- Manual installation adds complexity (validation, security checks, UX)
- Offline scenarios are rare for development tools (developers usually have network)
- One way to install assets (network download) follows KISS principle
- No security concerns with manually-injected asset files
- Simpler implementation (no bundle extraction, validation, or manual file handling)

Init succeeds without network because:

- Config file creation is the primary goal
- User can create config manually if needed
- Assets can be downloaded later via start assets update
- Blocking init on network would be frustrating for offline users
- Config creation is more important than asset download

Update fails without network because:

- Downloading assets is the explicit purpose of the command
- No point continuing if network is unavailable
- User expectation: update means download or nothing
- Clear failure better than partial success

Doctor shows partial results offline because:

- Local checks (config, agents, contexts) are valuable without network
- Skipping network checks is better than failing completely
- User gets useful information even offline
- Clear messaging about what was skipped

Core functionality works offline:

- Launching agents doesn't require assets
- Custom user configurations work without catalog
- Assets enhance UX (templates, default roles/tasks)
- Missing assets means reduced convenience, not broken tool

## Trade-offs

Accept:

- No air-gapped environment support (users must have network for asset downloads)
- No manual asset installation (users can't drop asset files manually)
- Init might succeed with incomplete setup (config created but no assets downloaded)
- Update always requires network (no offline mode)
- Doctor shows partial results offline (may miss important issues)

Gain:

- Simple implementation (one way to get assets: network download)
- Secure (no manual asset injection concerns or validation complexity)
- Clear network requirements (users know when network is needed)
- Graceful degradation (commands work without assets at reduced functionality)
- User-controlled init (succeeds even if assets fail, clear warnings)
- Explicit update contract (update requires network by design)
- Offline usability for core commands (start, start prompt, start task)

## Alternatives

Manual asset installation with --asset-path flag:

Example:
```bash
start assets update --from-directory /path/to/assets
start assets update --from-bundle assets-v1.2.3.tar.gz
```

Pros:
- Supports air-gapped environments
- Users can install assets without network
- Useful for corporate environments with restricted internet access

Cons:
- Adds complexity (bundle extraction, validation, directory structure checking)
- Security concerns (must validate manually-provided assets, prevent code injection)
- Two installation methods to maintain and test
- UX complexity (how do users get the bundle in the first place?)
- Rare use case (development tools typically have network access)
- Bundle distribution and versioning complexity

Rejected: Complexity and security concerns outweigh rare use case benefits. Development tools typically have network access.

Fail init when network unavailable:

Example: start init fails with "Network required" error

Pros:
- Ensures complete setup (config and assets both succeed or both fail)
- No partial setup states
- Clear all-or-nothing contract

Cons:
- Frustrating for users who just want to create config
- Config creation is more important than asset download
- Users can't proceed with manual config creation
- Blocks workflow unnecessarily (assets can be downloaded later)

Rejected: Config creation is valuable even without assets. Init should succeed with warnings.

Cache assets for offline use with full mirroring:

Example: Download all catalog assets during init, use cached versions when offline

Pros:
- Completely offline-capable after initial setup
- All catalog assets available without network
- No partial functionality degradation

Cons:
- Large download on init (entire asset catalog)
- Storage overhead (cache all assets even if unused)
- Conflicts with lazy loading design (download only what's needed)
- Asset updates require full re-download
- Complexity managing full catalog mirror

Rejected: Conflicts with catalog-based lazy loading architecture. Downloads should be on-demand.

## Structure

Network requirements by command:

Require network (fail if unavailable):
- start assets update - Explicit download operation, fails with clear error (exit 1)

Optional network (work without):
- start init - Creates config, warns on asset download failure (exit 0)
- start doctor - Shows partial results, skips network checks (exit based on local checks)
- start - Launches agent, works without assets
- start prompt - Launches agent with custom prompt, works without assets
- start task - Executes task, warns if referenced assets missing

Network failure exit codes:

- start init: Exit 0 (success, config created despite asset failure)
- start assets update: Exit 1 (network error, command failed)
- start doctor: Exit 0 or 1 based on local checks only (network checks skipped)
- Other commands: Varies by command (missing asset references may fail at runtime)

Missing assets behavior:

Pattern 1: Warn and continue (core commands)
- start / start prompt / start task work without assets
- Role references fail at runtime with clear error (expected behavior)

Pattern 2: Reduced functionality (asset commands)
- start assets add can't browse catalog but manual creation works
- Interactive prompts skip template selection

## Usage Examples

Init with network failure:

```bash
$ start init

Creating config at ~/.config/start/...

Warning: Unable to download assets from GitHub.
  Network error: dial tcp: no route to host

Config created successfully.

You can create your own configuration manually:
  start config edit global

Or try downloading assets later:
  start assets update
```

Exit code: 0 (success)

Update with network failure:

```bash
$ start assets update

Error: Cannot connect to GitHub

  Network error: dial tcp: no route to host

Update requires network access.
Asset library not modified.

Check your internet connection and try again.
```

Exit code: 1 (network error)

Doctor offline:

```bash
$ start doctor

Version Information:
  CLI Version:     1.3.0
  Update Check:    ⚠ Skipped (network unavailable)

Configuration:
  Global Config:
    settings.toml: ✓ ~/.config/start/settings.toml
    agents.toml:   ✓ ~/.config/start/agents.toml
    tasks.toml:    ✓ ~/.config/start/tasks.toml
    roles.toml:    ✓ ~/.config/start/roles.toml
    contexts.toml: ✓ ~/.config/start/contexts.toml
  Local Config:    - Not found
  Validation:      ✓ Valid

Asset Information:
  Cache Status:    ✓ Initialized
  Update Check:    ⚠ Skipped (network unavailable)

Agents:
  claude:          ✓ Binary found at /usr/local/bin/claude
  gemini:          ✓ Binary found at /usr/local/bin/gemini

Environment:
  GH_TOKEN:        ✓ Set
  EDITOR:          ✓ Set (vim)

Overall Status:   ✓ No local issues found
                  ⚠ Some checks skipped (network unavailable)

Run 'start doctor' with network access for complete check.
```

Exit code: 0 (local checks healthy)

Assets add without network:

```bash
$ start assets add

Add new asset
─────────────────────────────────────────────────

⚠ Asset catalog unavailable (network error)

Cannot browse catalog. You can still create assets manually.

[continues with manual input prompts, no templates offered]
```

Exit code: 0 (asset added successfully via manual input)

Task execution with missing role:

```bash
$ start task code-review

Error: Role 'code-reviewer' not found

The task references role 'code-reviewer' but it is not defined.

Available roles:
  (none)

To download catalog roles:
  start assets update

To create a custom role:
  start config role add
```

Exit code: 1 (missing role error)

Network error guidance:

```
Error: Cannot connect to GitHub

Check your internet connection and try again.

If you're behind a proxy, configure HTTP_PROXY and HTTPS_PROXY.
```

## Updates

- 2025-01-17: Removed references to superseded DR-011, DR-014, DR-018; updated to catalog-based asset system
