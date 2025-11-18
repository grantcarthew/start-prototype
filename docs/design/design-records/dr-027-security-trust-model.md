# DR-027: Security and Trust Model for Assets

- Date: 2025-01-07
- Status: Accepted
- Category: Asset Management

## Problem

The CLI downloads assets (roles, tasks, agent templates) from a remote source. The security strategy must address:

- Trust model (what sources do we trust and why?)
- Transport security (how to prevent man-in-the-middle attacks?)
- Content authenticity (how to verify assets are legitimate?)
- Tampering detection (how to detect modified assets?)
- Repository security (how to prevent typosquatting or malicious repos?)
- Key management (if using signatures, how to distribute and verify keys?)
- Attack surface (how to minimize security-related code complexity?)
- Threat priorities (what threats matter most for non-executable content?)
- User control (how to let users verify and inspect assets?)
- Maintainability (how to keep security simple and auditable?)

## Decision

Trust GitHub's HTTPS infrastructure and the specific hardcoded repository. No cryptographic signatures. No commit pinning. Always fetch latest from main branch.

Trust model:

- Trust GitHub's infrastructure and HTTPS transport security
- Trust the specific hardcoded repository: github.com/grantcarthew/start
- Trust repository maintainer's account security (2FA, strong credentials)
- Do not trust arbitrary repositories (no user-configurable repo URLs)
- Do not trust manual asset installations (network-only per DR-026)

No signature verification:

- No GPG signing of commits
- No individual file signatures
- No signature verification during download
- Trust GitHub's HTTPS as sufficient

No commit pinning:

- Always latest from main branch (no commit SHA pinning in config)
- No git tag/release pinning
- Users always get latest when running start assets update
- No configuration options for custom repository URLs, commit SHAs, or tags

Protection mechanisms:

- HTTPS prevents man-in-the-middle attacks during download
- Hardcoded repository URL in binary (no config override)
- Assets cached locally after download
- Users can inspect all assets before and after use

## Why

Assets are low-risk content:

- Markdown files (role definitions) - human-readable text
- TOML files (agent configs, tasks) - configuration data
- No executables, no auto-executed scripts
- Users review templates before using them
- Users can inspect asset files anytime
- Content is visible and auditable

HTTPS provides adequate security:

- Transport encryption via TLS
- Server authentication via certificate verification
- Data integrity via TLS checksums
- Industry-standard security trusted by GitHub
- System certificate store for verification (no custom pinning needed)

Hardcoded repository prevents attacks:

- No typosquatting (users can't misconfigure malicious repo URL)
- No user configuration of asset source (compile-time constant)
- Clear single source of truth
- Reduces attack surface (no URL parsing or validation)

Simplicity reduces attack surface:

- Less code means fewer bugs and vulnerabilities
- No key management complexity (distribution, rotation, revocation)
- No signature verification failure handling
- Easier to audit and reason about
- Clearer security model for users and developers

User control provides safety:

- Users initiate all updates (no automatic downloads)
- Users can inspect assets before use
- Users can review asset changes via GitHub commit history
- Users decide when to update (control timing)

Protection is proportional to risk:

- Non-executable content has lower risk than binaries
- Users review templates before using (not blind execution)
- HTTPS + hardcoded repo prevents common attacks
- Additional signatures provide minimal benefit for text files
- Complexity costs outweigh marginal security improvements

## Trade-offs

Accept:

- No defense against compromised maintainer account (maintainer responsibility, same as any open-source project)
- No cryptographic proof of authenticity beyond HTTPS (HTTPS + hardcoded repo is sufficient for text files)
- No commit pinning for reproducibility (user controls update timing, can inspect changes)
- Trust GitHub infrastructure completely (acceptable dependency, standard practice)
- Accept supply chain risk from bad updates (users control timing, can inspect, no auto-execution)

Gain:

- Extremely simple implementation (no key management, signature verification, or pinning logic)
- Secure enough for content distribution (HTTPS + hardcoded repo prevents common attacks)
- Transparent to users (all assets inspectable locally)
- Low risk for non-executable content (text files, user-reviewed templates)
- Maintainable security model (less code, fewer dependencies, clearer design)
- Industry-standard approach (similar to Homebrew, npm, cargo for content)
- No verification failure modes (signatures can fail for many reasons)

## Alternatives

Cryptographic signature verification:

Example approaches:
- GPG-signed commits (git verify-commit)
- Individual file signatures (detached .sig files)
- Signed release tags

Pros:
- Cryptographic proof of authenticity
- Detect compromised repository or man-in-the-middle attacks
- Industry best practice for security-critical software

Cons:
- Key management complexity (distribution, storage, rotation, revocation)
- Public key distribution problem (how do users get the trusted key?)
- Signature verification can fail (expired keys, clock skew, key rotation issues)
- Users confused when verification fails
- Minimal security improvement for non-executable text files
- More code to audit and maintain
- Additional dependencies and failure modes

Rejected: Complexity and failure modes outweigh benefits for non-executable content. HTTPS provides adequate security for markdown and TOML files.

Commit pinning with SHA references:

Example: Config specifies exact commit SHA to download
```toml
[settings]
asset_commit = "abc123def456..."
```

Pros:
- Reproducible asset versions (same SHA = same content)
- Users can choose when to update (explicit SHA change)
- Protection against unexpected changes

Cons:
- Users must manually update commit SHAs (friction)
- Reduces benefit of "always latest" design
- False sense of security (still trusting GitHub and HTTPS)
- Complexity in config management
- Conflicts with main branch strategy
- Doesn't prevent attacks, just controls timing

Rejected: User controls update timing already (user-initiated updates only). SHA pinning adds complexity without meaningful security benefit.

Multiple trusted repositories with fallbacks:

Example: Try official repo, fallback to mirrors
```go
repos := []string{
    "github.com/grantcarthew/start",
    "gitlab.com/grantcarthew/start-mirror",
    "codeberg.org/grantcarthew/start-mirror",
}
```

Pros:
- Availability if one source is down
- Reduced dependency on single platform
- Geographic distribution

Cons:
- Multiple sources to trust and secure
- Synchronization complexity (which is canonical?)
- Authentication complexity (trust all equally?)
- Attack surface increases (compromise any mirror = success)
- Maintenance burden (keep mirrors in sync)
- Overkill for optional content (assets not critical)

Rejected: Single trusted source is simpler and more secure. GitHub availability is sufficient for optional content.

## Structure

Hardcoded repository configuration:

Repository constants (compile-time, not configurable):
- Repository: github.com/grantcarthew/start
- Branch: main
- No environment variable override
- No config file override
- No command-line flag override

Threat model:

Threats mitigated:

1. Man-in-the-middle attacks:
   - Threat: Attacker intercepts network and injects malicious assets
   - Mitigation: HTTPS with certificate verification
   - Residual risk: Very low (requires compromising GitHub TLS)

2. Typosquatting:
   - Threat: User misconfigures repo URL to malicious look-alike
   - Mitigation: No user-configurable repo URL (hardcoded)
   - Residual risk: None (no configuration option exists)

Threats accepted:

1. Compromised GitHub account:
   - Threat: Attacker gains access to grantcarthew/start repository
   - Mitigation: Maintainer uses 2FA, strong credentials, GitHub's security
   - Residual risk: Low (GitHub account security is maintainer responsibility)
   - Accepted: Yes (same risk as any open-source project, standard practice)

2. Compromised GitHub infrastructure:
   - Threat: GitHub itself compromised, serves malicious content
   - Mitigation: None (if GitHub is compromised, entire ecosystem affected)
   - Residual risk: Very low (GitHub has strong security practices)
   - Accepted: Yes (acceptable dependency, industry standard)

3. Malicious asset content:
   - Threat: Malicious TOML/markdown content in assets
   - Mitigation: Users review before using, no auto-execution, inspectable content
   - Residual risk: Low (content is visible, not auto-executed)
   - Accepted: Yes (users responsible for reviewing templates)

4. Supply chain attack via bad update:
   - Threat: Malicious commit pushed to main, users update
   - Mitigation: Users control update timing, can inspect assets, user-initiated only
   - Residual risk: Low (no automatic updates, user reviews changes)
   - Accepted: Yes (same risk as any auto-updating software with user control)

Download process security:

HTTPS download:
- System certificate store for verification (standard library http.Client)
- No custom certificate pinning
- 30-second timeout for network calls
- GitHub API over HTTPS

Atomic installation:
- Download to temporary location
- Verify download completed successfully
- Move to final location atomically
- Prevents partial/corrupted installations

User inspection capabilities:

Asset transparency:
- All assets visible in local cache directory
- Users can read any asset file
- Users can review GitHub commit history
- Users can compare local vs remote versions

Update control:
- User-initiated updates only (no automatic downloads)
- Users decide when to update
- Users can inspect changes before updating
- Users can rollback by not updating

## Usage Examples

Inspecting downloaded assets:

```bash
# View asset cache directory
ls -R ~/.cache/start/

# Inspect specific role
cat ~/.cache/start/roles/code-reviewer.md

# Compare with GitHub version
start assets diff code-reviewer
```

Reviewing changes before update:

```bash
# Check for available updates
start doctor

# Review what changed on GitHub
start assets changes

# Update when ready
start assets update
```

Comparison with other tools:

Similar trust model (HTTPS only):
- Homebrew: Trusts GitHub, HTTPS, no signatures on formulae
- npm: Trusts registry, HTTPS (optional signatures available)
- cargo: Trusts crates.io, HTTPS, checksums only
- pip: Trusts PyPI, HTTPS (optional GPG verification)

More paranoid (signatures required):
- apt/yum: GPG signatures on packages (system packages, executables)
- Arch pacman: Package signing required (system binaries)
- Signal: Binary transparency, reproducible builds (privacy-critical)
- Tor: Multi-signature releases (security-critical software)

Our position:
- More like Homebrew (content packages, not executables)
- Less like apt (text templates, not system binaries)
- Appropriate for markdown/TOML asset distribution

## Updates

- 2025-01-17: Removed references to superseded DR-011, DR-014, DR-015; updated asset paths to cache directory for catalog system
