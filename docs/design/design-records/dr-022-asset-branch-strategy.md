# DR-022: Asset Branch Strategy

**Date:** 2025-01-06
**Status:** Accepted
**Category:** Asset Management

## Decision

Asset library updates (`start update` and `start init`) always pull from the **latest commit on the main branch**, not from GitHub Releases.

## Rationale

**Assets are content, not code:**
- Roles: Markdown prompt documents
- Agents: TOML configuration templates
- Tasks: TOML task definitions
- Contexts: TOML context configurations
- Examples: TOML reference configurations

Content updates have lower risk than binary code updates and benefit from rapid iteration.

**Decoupled update cycles:**
- **CLI binary**: Tied to GitHub Releases (DR-021)
  - Users update via: `brew upgrade` or `go install`
  - Versioned, tested, stable releases
  - Changed via your manual release process

- **Asset library**: Tied to main branch commits
  - Users update via: `start update`
  - Can iterate independently of CLI releases
  - Improvements available immediately

**Benefits:**
- Iterate on tasks/roles/agents/contexts without cutting CLI releases
- Users get content improvements immediately
- Simple mental model: "Code = releases, Content = main branch"
- Faster feedback loop for asset improvements

## Implementation

### Asset Update Target

Using DR-014's Tree API mechanism:

```bash
# Always pull from main branch
GET /repos/grantcarthew/start/git/trees/main?recursive=1
```

**Not from releases:**
```bash
# ❌ Don't do this for assets
GET /repos/grantcarthew/start/releases/latest  # Get tag
GET /repos/grantcarthew/start/git/trees/{tag}?recursive=1
```

### Version Tracking

The `asset-version.toml` file tracks the main branch commit SHA:

```toml
# Asset version tracking - managed by 'start update'
# Last updated: 2025-01-06T10:30:00Z

commit = "abc123def456"  # Latest commit SHA from main branch
timestamp = "2025-01-06T10:30:00Z"
repository = "github.com/grantcarthew/start"
branch = "main"  # Always main

[files]
"agents/claude.toml" = "a1b2c3d4e5f6..."
"agents/gemini.toml" = "e5f6g7h8i9j0..."
"roles/code-reviewer.md" = "i9j0k1l2m3n4..."
# ... etc
```

### Update Flow

```bash
$ start update

Checking asset library...
  Current commit:  abc1234 (3 days ago)
  Latest commit:   def5678 (2 hours ago)
  Branch:          main

Downloading updates...
  ✓ roles/golang-expert.md (updated)
  ✓ tasks/code-review.toml (updated)
  ✓ agents/claude.toml (unchanged, skipped)

✓ Asset library updated to commit def5678

CLI Version Check:
  Current:         v1.2.3
  Latest Release:  v1.3.0
  Status:          ⚠ Update available
  Update via:      brew upgrade grantcarthew/tap/start
```

## Asset Stability Strategy

**To prevent broken assets reaching users:**

1. **Keep main branch stable**
   - Test assets before merging to main
   - Don't commit broken TOML or malformed markdown
   - Review asset changes like code

2. **Optional: Use develop branch** (future consideration)
   - Develop new assets in `develop` branch
   - Only merge to `main` when tested
   - Could add `--branch develop` flag later if needed

3. **Validation on update** (future consideration)
   - `start update` could validate TOML syntax before installing
   - Rollback on validation failure (DR-015)
   - Show warnings for deprecated fields

## Comparison: CLI vs Assets

| Aspect | CLI Binary | Asset Library |
|--------|------------|---------------|
| **Source** | GitHub Releases | Main branch commits |
| **Version** | Semantic (v1.2.3) | Git commit SHA |
| **Update Command** | `brew upgrade` / `go install` | `start update` |
| **Update Frequency** | Manual releases | Every commit to main |
| **Testing** | Full release process | Commit to main |
| **Stability** | High (release process) | Medium (review before merge) |
| **Risk** | High (executable code) | Low (config templates) |

## Future Considerations

**Branch selection flag** (not implementing now):
```bash
# Could add later if needed
start update --branch develop  # Bleeding edge
start update --branch main     # Default
start update --tag v1.2.0      # Pin to specific version
```

**Asset channels** (not implementing now):
```toml
# Could allow users to choose stability level
[settings]
asset_channel = "stable"   # From releases
asset_channel = "latest"   # From main (default)
asset_channel = "develop"  # From develop branch
```

These are **not** part of this decision - keeping it simple with main branch only.

## Benefits

- ✅ **Faster iteration** - Improve assets without CLI releases
- ✅ **Simple** - One branch, no channel complexity
- ✅ **Immediate availability** - Users get improvements on next update
- ✅ **Decoupled** - Asset improvements independent of CLI changes
- ✅ **Lower risk** - Content files vs executable binaries

## Trade-offs Accepted

- ❌ Broken assets could reach users (mitigated: keep main stable)
- ❌ No "stable" vs "bleeding edge" choice (acceptable: assets are templates)
- ❌ No version pinning for assets (acceptable: users control update timing)

## Rationale Summary

Assets are content that benefits from rapid iteration. Pulling from main branch allows:
- Fixing typos in role prompts without a release
- Adding new task templates as they're created
- Updating agent configs for new model versions
- Adding new context configurations for common patterns
- Iterating on examples based on user feedback

Users control when they get updates (via `start update`), so they're not forced to take changes immediately. The CLI binary remains tied to stable releases for safety.

## Related Decisions

- [DR-014](./dr-014-github-tree-api.md) - GitHub Tree API mechanism (how we fetch)
- [DR-015](./dr-015-atomic-updates.md) - Atomic updates with rollback
- [DR-018](./dr-018-init-update-integration.md) - Init and update integration
- [DR-021](./dr-021-github-version-check.md) - CLI version checking (releases)

## Implementation Notes

### Code Location

Update `internal/assets/updater.go`:

```go
const (
    AssetRepository = "grantcarthew/start"
    AssetBranch     = "main"  // Always main
    AssetsBasePath  = "assets"
)

func FetchLatestAssets(ctx context.Context) error {
    // GET /repos/grantcarthew/start/git/trees/main?recursive=1
    // Use Tree API from DR-014
}
```

### Asset Version Display

`start doctor` shows asset commit and age:

```bash
Asset Information:
  Commit:        abc1234
  Branch:        main
  Last Updated:  2 days ago
  Status:        ✓ Up to date
```

No "version number" for assets - just commit SHA and timestamp.

## Implementation Checklist

- [ ] Update `internal/assets/updater.go` to target main branch
- [ ] Ensure `asset-version.toml` tracks branch field
- [ ] Update `start doctor` to show asset commit and age
- [ ] Update `start update` output to show current/latest commits
- [ ] Document in README that assets come from main branch
- [ ] Add validation for TOML syntax before installing assets (optional)
