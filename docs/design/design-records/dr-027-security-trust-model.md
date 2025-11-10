# DR-027: Security and Trust Model for Assets

**Date:** 2025-01-07
**Status:** Accepted
**Category:** Asset Management

## Decision

Trust GitHub's HTTPS infrastructure and the specific repository (`grantcarthew/start`). No cryptographic signatures. No commit pinning. Always fetch latest from main branch.

## What This Means

### Trust Model (17a)

**We trust:**
- GitHub's infrastructure and HTTPS transport security
- The specific repository: `github.com/grantcarthew/start`
- Repository maintainer's account security (2FA, strong credentials)

**We do NOT trust:**
- Arbitrary repositories (no user-configurable repo URLs)
- Unsigned/unverified sources
- Manual asset installations (per DR-026)

**Protection mechanisms:**
- HTTPS prevents man-in-the-middle attacks during download
- SHA-based caching (DR-014) detects file tampering after first download
- Hardcoded repository URL in binary (no config override)

### No Signature Verification (17b)

**No cryptographic signatures on assets:**
- No GPG signing of commits
- No individual file signatures
- No signature verification during download
- Trust GitHub's HTTPS as sufficient

**Rationale:**
- Assets are not executables (markdown and TOML files)
- Users review assets before using (templates, roles, tasks)
- HTTPS + GitHub infrastructure provides adequate security
- Signatures add complexity without proportional security benefit
- No automatic execution means low risk

### No Commit Pinning (17c)

**Always latest from main branch:**
- No commit SHA pinning in config
- No git tag/release pinning
- Users always get latest when running `start update`
- Matches DR-022 (assets from main branch)

**No configuration options for:**
- Custom repository URLs
- Commit SHA pins
- Git tag/release references
- Asset source override

## Rationale

### Pragmatic Security Posture

**Assets are low-risk content:**
- Markdown files (role definitions) - human-readable text
- TOML files (agent configs, tasks) - configuration data
- No executables, no scripts auto-executed
- Users review templates before using them
- Users can inspect asset files anytime

**Protection is proportional:**
- HTTPS prevents network interception
- SHA checking detects post-download tampering
- Hardcoded repo prevents typosquatting
- Simple implementation reduces attack surface

### HTTPS is Sufficient

**GitHub's HTTPS provides:**
- Transport encryption (TLS)
- Server authentication (certificate verification)
- Data integrity (TLS checksums)
- Industry-standard security

**Additional signatures would add:**
- Key management complexity
- Public key distribution problem
- Verification failure handling
- Minimal security improvement for non-executable content

### Simplicity Over Paranoia

**Complexity costs:**
- More code = more bugs = more attack surface
- Key management is hard to get right
- Signature verification can fail (expired keys, clock skew, etc.)
- User confusion when verification fails

**Benefits of simplicity:**
- Less code to audit
- Fewer failure modes
- Clearer security model
- Easier to reason about

## Threat Model

### Threats We Mitigate

**1. Man-in-the-Middle (MITM) attacks:**
- **Threat:** Attacker intercepts network traffic and injects malicious assets
- **Mitigation:** HTTPS with certificate verification
- **Residual risk:** Very low (requires compromising GitHub's TLS infrastructure)

**2. Post-download tampering:**
- **Threat:** Local asset files modified after download
- **Mitigation:** SHA-based caching detects changes on next update
- **Residual risk:** Low (user would notice broken behavior before next update)

**3. Typosquatting:**
- **Threat:** User misconfigures repo URL to malicious look-alike
- **Mitigation:** No user-configurable repo URL (hardcoded in binary)
- **Residual risk:** None (no configuration option)

### Threats We Accept

**1. Compromised GitHub account:**
- **Threat:** Attacker gains access to grantcarthew/start repository
- **Mitigation:** Repository maintainer uses 2FA, strong credentials, GitHub's security
- **Residual risk:** Low (GitHub account security is maintainer's responsibility)
- **Accepted:** Yes (same risk as any open-source project)

**2. Compromised GitHub infrastructure:**
- **Threat:** GitHub itself is compromised, serves malicious content
- **Mitigation:** None (if GitHub is compromised, we have bigger problems)
- **Residual risk:** Very low (GitHub has strong security)
- **Accepted:** Yes (acceptable dependency on GitHub)

**3. Malicious asset content:**
- **Threat:** Malicious TOML/markdown content in assets
- **Mitigation:** Users review before using, no auto-execution
- **Residual risk:** Low (content is inspectable, not auto-executed)
- **Accepted:** Yes (users responsible for reviewing templates)

**4. Supply chain attack via bad update:**
- **Threat:** Malicious commit pushed to main, users auto-update
- **Mitigation:** Users control update timing (DR-025), can inspect assets
- **Residual risk:** Low (user-initiated updates only)
- **Accepted:** Yes (same risk as any auto-updating software)

## Implementation

### Hardcoded Repository

**Binary contains:**
```go
const (
    AssetRepository = "github.com/grantcarthew/start"
    AssetBranch     = "main"
    AssetPath       = "/assets"
)
```

**No configuration options:**
- No `asset_repository` in config file
- No `--repo` flag on `start update`
- No environment variable override
- Repository URL is compile-time constant

### Download Process

**Security checks during download:**
```go
func downloadAssets() error {
    // 1. HTTPS URL construction (hardcoded repo)
    url := fmt.Sprintf("https://api.github.com/repos/%s/...", AssetRepository)

    // 2. Standard library HTTP client (respects system cert store)
    client := &http.Client{
        Timeout: 30 * time.Second,
        // Uses system certificate verification automatically
    }

    // 3. SHA verification on subsequent downloads
    if existingSHA != downloadedSHA {
        // Update needed
    }

    // 4. Atomic install (DR-015)
    return atomicInstall(downloadedAssets)
}
```

**What we DON'T do:**
- No custom certificate pinning
- No signature verification
- No checksum files from separate source
- No manual verification prompts

### User Transparency

**Users can inspect assets:**
```bash
# View all downloaded assets
ls -R ~/.config/start/assets/

# Inspect specific role
cat ~/.config/start/assets/roles/code-reviewer.md

# Check asset version
cat ~/.config/start/asset-version.toml
```

**Users control updates:**
- Per DR-025: No automatic updates
- `start update` is user-initiated
- `start doctor` shows asset age
- Users decide when to update

## Security Best Practices

### For Repository Maintainer

**Account security:**
- Enable 2FA on GitHub account
- Use strong, unique password
- Review account access regularly
- Use SSH keys for git operations

**Repository security:**
- Enable branch protection on main
- Require code review for PRs (if team)
- Monitor repository access logs
- Keep dependencies updated

**Asset review:**
- Review all asset changes carefully
- Test assets before committing to main
- Document asset changes in commits
- Consider security implications of templates

### For Users

**Asset hygiene:**
- Run `start doctor` periodically to check for updates
- Review asset changes after `start update` (check git history)
- Inspect role/task templates before using
- Report suspicious content to maintainer

**Update timing:**
- Update when convenient (user-controlled per DR-025)
- Check release notes/commit history before updating
- Test after updates to detect issues
- Keep backups of working configurations

## Comparison with Other Tools

**Similar trust model (HTTPS only):**
- **Homebrew:** Trusts GitHub, HTTPS, no signatures on formulae
- **npm:** Trusts registry, HTTPS (supports optional signatures)
- **pip:** Trusts PyPI, HTTPS (optional GPG verification)
- **cargo:** Trusts crates.io, HTTPS (checksums only)

**More paranoid (signatures required):**
- **apt/yum:** GPG signatures on packages
- **Arch pacman:** Package signing required
- **Signal:** Binary transparency, reproducible builds
- **Tor:** Multi-signature releases

**Our position:**
- More like Homebrew (content packages, not executables)
- Less like apt (system packages, higher stakes)
- Appropriate for markdown/TOML asset distribution

## Benefits

- ✅ **Simple** - No key management, signature verification, or pinning logic
- ✅ **Secure enough** - HTTPS + hardcoded repo prevents common attacks
- ✅ **Transparent** - Users can inspect all assets locally
- ✅ **Low risk** - Non-executable content, user-reviewed templates
- ✅ **Maintainable** - Less code, fewer dependencies, clearer model
- ✅ **Industry standard** - Similar to other content distribution tools

## Trade-offs Accepted

- ❌ No defense against compromised maintainer account (acceptable: maintainer responsibility)
- ❌ No cryptographic proof of authenticity (acceptable: HTTPS + hardcoded repo is sufficient)
- ❌ No commit pinning for reproducibility (acceptable: user controls update timing)
- ❌ Trust GitHub infrastructure (acceptable: standard dependency)

## Future Considerations

**If threat model changes, we could add:**

**Option 1: Commit signing verification**
```bash
# Verify commits are signed by maintainer
git verify-commit abc123def456
```

**Option 2: Asset checksums from separate source**
```bash
# Checksums published via DNS TXT record or separate repo
start update --verify-checksums
```

**Option 3: Reproducible asset builds**
```bash
# Assets generated deterministically, users can reproduce
start update --verify-reproducible
```

**Current stance:** Don't implement unless threat landscape changes or users request it. Current model is appropriate for content distribution.

## Related Decisions

- [DR-011](./dr-011-asset-distribution.md) - GitHub-fetched assets (establishes network dependency)
- [DR-014](./dr-014-github-tree-api.md) - SHA-based caching (file integrity checking)
- [DR-022](./dr-022-asset-branch-strategy.md) - Main branch strategy (no release tags)
- [DR-025](./dr-025-no-automatic-checks.md) - User-initiated updates (user controls timing)
- [DR-026](./dr-026-offline-behavior.md) - Network-only approach (no manual injection)

## Documentation

### User Documentation

**README security section:**
```markdown
## Security

`start` downloads assets from GitHub over HTTPS:
- Repository: github.com/grantcarthew/start
- Transport: HTTPS (encrypted, authenticated)
- Content: Markdown and TOML files (non-executable)

You can inspect downloaded assets:
- Location: ~/.config/start/assets/
- View: cat ~/.config/start/assets/roles/code-reviewer.md

Updates are user-initiated only:
- Run `start update` when you want to update
- Review changes: git log in the repository
```

### Developer Documentation

**Security guidelines for contributors:**
```markdown
## Asset Security

When adding/modifying assets:
- Assets are trusted content (users will use them)
- Review all template content carefully
- Avoid suggesting dangerous commands in templates
- Document template purpose and usage
- Test thoroughly before committing to main

Repository security:
- Enable 2FA on your GitHub account
- Use strong credentials
- Review access permissions regularly
```
