# DR-032: Asset Metadata Schema

- Date: 2025-01-10
- Status: Accepted
- Category: Asset Management

## Problem

Catalog assets need metadata for discovery, searching, and update detection. The metadata strategy must address:

- Storage approach (embedded frontmatter vs sidecar files)
- Schema design (which fields, required vs optional)
- Category tracking (explicit field vs derived from path)
- Update detection (version numbers vs content hashing)
- Validation requirements (ensure data quality)
- Searchability (enable filtering and discovery)
- Drift prevention (ensure metadata matches reality)
- File management (minimize clutter in content files)
- Consistency (same approach for all asset types)

## Decision

Use sidecar .meta.toml files for asset metadata with 6 required fields, deriving category from filesystem path to prevent drift.

Sidecar metadata approach:

Each catalog asset has a separate metadata file:

```
assets/tasks/git-workflow/
├── pre-commit-review.toml        # Asset content
├── pre-commit-review.md          # Prompt file (if using UTD)
└── pre-commit-review.meta.toml   # Metadata (sidecar)
```

Metadata schema (all fields required):

```toml
[metadata]
name = "pre-commit-review"
description = "Review staged changes before committing"
tags = ["git", "review", "quality", "pre-commit"]
sha = "a1b2c3d4e5f6789012345678901234567890abcd"
created = "2025-01-10T00:00:00Z"
updated = "2025-01-10T00:00:00Z"
```

Field definitions:

- name (string): Asset identifier, must match filename without extensions
- description (string): One-line summary shown in catalog browsing and search
- tags (array of strings): Keywords for searching and filtering, at least 1 required
- sha (string): Git blob SHA of content files (40-char hex), used for update detection
- created (ISO 8601 timestamp): When asset was first created
- updated (ISO 8601 timestamp): When asset was last modified

Category derived from filesystem:

No category field in metadata - derive from directory structure:

```
assets/tasks/git-workflow/pre-commit-review.toml
       ^     ^
       |     └─ category = "git-workflow"
       └─────── asset type = "tasks"
```

SHA generation:

- Single-file assets: SHA of that file
- Multi-file assets: SHA of files concatenated in sorted order (excludes .meta.toml)
- 40-character hexadecimal string
- Used for update detection via comparison

## Why

Sidecar files keep content clean:

- Content files stay clean (no metadata clutter for AI to process)
- Metadata separate from what agent sees
- Easy to parse independently (no frontmatter parsing)
- No risk of corrupting content when updating metadata
- Consistent format across asset types (roles, tasks, agents, contexts)

Required fields ensure data quality:

- No optional fields means no missing data
- Clear validation rules (simpler to check)
- Catalog always has complete information
- No need to handle missing field cases

Category derived from filesystem prevents drift:

- Eliminates potential drift between category field and actual path
- Filesystem IS the source of truth (category can't be wrong)
- Simpler validation (one less field to check)
- Enforces correct directory structure

SHA-based versioning enables reliable updates:

- SHA comparison for reliable change detection
- No version numbers to maintain manually
- Content hash IS the version (always accurate)
- Simple comparison (local SHA vs GitHub SHA)

Tags enable rich searching:

- Enable filtering without full-text indexing
- Quick scanning by keywords
- Structured data for catalog browsing
- Users can search by domain-specific terms

Timestamps provide context:

- Created timestamp shows asset age
- Updated timestamp shows freshness
- Helps users understand asset maturity
- ISO 8601 standard format (widely supported)

## Trade-offs

Accept:

- More files to manage (sidecar file for each asset, but worth it for cleanliness)
- No semantic versioning (SHA not human-readable like v1.2.3, but more reliable)
- No author/license fields (can't track who created asset or licensing, add when community contributions happen)
- No minimum version requirement (can't specify "requires start >= 0.2.0", add when needed)
- Timestamps required manually (manual creation requires setting timestamps, tooling can help)
- All fields required (no flexibility for partial metadata, but ensures quality)

Gain:

- Clean content files (no metadata clutter, AI sees only relevant content)
- No drift possible (category derived from filesystem, name must match filename, SHA tracks actual content)
- Reliable update detection (SHA comparison simple and accurate, no version number maintenance)
- Searchability (tags enable filtering, description for quick scanning, structured data for browsing)
- Simple validation (all fields required, clear rules, easy to check)
- Consistent across asset types (same approach for roles, tasks, agents, contexts)
- Easy to parse independently (separate TOML file, no frontmatter complexity)

## Alternatives

Frontmatter in content files:

Example: Embed metadata in content files

```toml
# pre-commit-review.toml
[metadata]
description = "Review staged changes before committing"
tags = ["git", "review"]
# ... more metadata

[task]
# ... actual task definition
```

Pros:

- Single file per asset (simpler file management)
- Metadata and content together (no separate file to maintain)
- Fewer files in repository

Cons:

- Content files cluttered with metadata (AI processes metadata too)
- Risk of corrupting content when updating metadata
- Inconsistent across asset types (Markdown can't embed TOML frontmatter easily)
- Harder to parse independently (must parse entire file)
- Metadata mixed with what agent sees (not clean separation)

Rejected: Content cleanliness more important than file count. Sidecar files provide clear separation.

Category field in metadata:

Example: Explicit category field instead of deriving from path

```toml
[metadata]
category = "git-workflow"
# ... other fields
```

Pros:

- Explicit category declaration (clear in metadata file)
- Could move files without changing category (flexible organization)
- Self-documenting (category visible in metadata)

Cons:

- Potential drift (category field doesn't match actual path)
- Extra validation needed (must check field matches path)
- More complex (one more field to maintain and validate)
- Can be wrong (filesystem truth vs metadata claim conflict)

Rejected: Deriving from filesystem prevents drift. Single source of truth is simpler and more reliable.

Semantic versioning with version field:

Example: Use semver instead of SHA

```toml
[metadata]
version = "1.2.3"
# ... other fields
```

Pros:

- Human-readable versions (v1.2.3 easier to understand than SHA)
- Semantic meaning (major/minor/patch conveys change type)
- Familiar pattern (developers understand semver)

Cons:

- Manual maintenance (must remember to bump version)
- Can be wrong (version updated but content unchanged, or vice versa)
- No automatic detection (can't tell if version should change)
- Requires discipline (easy to forget to update)

Rejected: SHA is more reliable (automatic, always accurate, no manual maintenance). SHA IS the version.

Optional metadata fields:

Example: Make some fields optional

```toml
[metadata]
name = "my-task"
description = "..." # required
tags = [...] # optional - can be empty
created = "..." # optional - might be missing
```

Pros:

- Flexibility (easier to add assets with partial metadata)
- Less strict (don't need to fill out everything)
- Gradual completion (can add fields later)

Cons:

- Incomplete catalog data (missing tags means unsearchable)
- More complex validation (must handle missing field cases)
- Inconsistent experience (some assets have rich data, others don't)
- Quality degradation (easier to skip important metadata)

Rejected: All required fields ensures data quality. Better to have complete metadata for every asset.

## Structure

Metadata schema:

Required fields (all must be present):

name (string):

- Asset identifier
- Must match filename without extensions
- Example: "pre-commit-review" for pre-commit-review.toml
- Validated: name must equal filename base

description (string):

- One-line summary of what the asset does
- Shown in catalog browsing and search results
- Example: "Review staged changes before committing"
- No length limit but should be concise

tags (array of strings):

- Keywords for searching and filtering
- At least 1 tag required
- Example: ["git", "review", "quality", "pre-commit"]
- Used for substring matching in search

sha (string):

- Git blob SHA of content files
- 40-character hexadecimal string
- Excludes .meta.toml file itself
- Single-file: SHA of that file
- Multi-file: SHA of files concatenated in sorted order
- Example: "a1b2c3d4e5f6789012345678901234567890abcd"
- Used for update detection

created (ISO 8601 timestamp):

- When asset was first created
- Format: "2025-01-10T00:00:00Z"
- Never changes after initial creation

updated (ISO 8601 timestamp):

- When asset was last modified
- Format: "2025-01-10T12:30:00Z"
- Updated whenever content changes
- Must be >= created timestamp

Derived fields (not in TOML, computed from filesystem):

Category:

- Derived from directory structure
- Example: assets/tasks/git-workflow/ → category = "git-workflow"
- Prevents drift between metadata and actual location

AssetType:

- Derived from directory structure
- Example: assets/tasks/ → type = "tasks"
- One of: tasks, roles, agents, contexts

Validation rules:

Required field checks:

- All 6 fields must be present
- tags array must have at least 1 element
- No empty strings for name, description, sha

SHA format validation:

- Must be exactly 40 characters
- Must be hexadecimal [a-f0-9]
- Example valid: "a1b2c3d4..."
- Example invalid: "xyz123..." (non-hex)

Name match validation:

- name field must match filename base
- Example: pre-commit-review.meta.toml → name must be "pre-commit-review"
- Prevents naming mismatches

Timestamp validation:

- Must be valid ISO 8601 format
- updated must be >= created
- No future timestamps (should be <= now, but not enforced)

File structure:

Each asset has matching files:

```
assets/{type}/{category}/{name}.toml       # Required: asset content
assets/{type}/{category}/{name}.md         # Optional: prompt file for tasks/roles
assets/{type}/{category}/{name}.meta.toml  # Required: metadata
```

Metadata location:

- Same directory as asset content
- Same base filename with .meta.toml extension
- One metadata file per asset (not per file)

## Usage Examples

Task metadata:

```toml
# assets/tasks/git-workflow/pre-commit-review.meta.toml
[metadata]
name = "pre-commit-review"
description = "Review staged changes before committing"
tags = ["git", "review", "quality", "pre-commit"]
sha = "a1b2c3d4e5f6789012345678901234567890abcd"
created = "2025-01-10T00:00:00Z"
updated = "2025-01-10T12:30:00Z"
```

Corresponding files:

```
assets/tasks/git-workflow/
├── pre-commit-review.toml        # Task definition
├── pre-commit-review.md          # Prompt file
└── pre-commit-review.meta.toml   # This metadata
```

Role metadata:

```toml
# assets/roles/general/code-reviewer.meta.toml
[metadata]
name = "code-reviewer"
description = "Strict quality and security focused code reviewer"
tags = ["review", "quality", "security", "strict"]
sha = "b2c3d4e5f6789012345678901234567890abcde1"
created = "2025-01-10T00:00:00Z"
updated = "2025-01-10T00:00:00Z"
```

Agent metadata:

```toml
# assets/agents/anthropic/claude.meta.toml
[metadata]
name = "claude"
description = "Anthropic Claude AI via Claude Code CLI"
tags = ["claude", "anthropic", "ai", "recommended"]
sha = "c3d4e5f6789012345678901234567890abcde12"
created = "2025-01-10T00:00:00Z"
updated = "2025-01-10T00:00:00Z"
```

Context metadata:

```toml
# assets/contexts/reference/environment.meta.toml
[metadata]
name = "environment"
description = "System environment and context information"
tags = ["environment", "system", "context", "reference"]
sha = "d4e5f6a1b2c3d4e5f6789012345678901234abcd"
created = "2025-01-10T00:00:00Z"
updated = "2025-01-10T00:00:00Z"
```

Category derivation from path:

```
assets/tasks/git-workflow/pre-commit-review.meta.toml
→ AssetType: "tasks"
→ Category: "git-workflow"

assets/roles/general/code-reviewer.meta.toml
→ AssetType: "roles"
→ Category: "general"

assets/agents/anthropic/claude.meta.toml
→ AssetType: "agents"
→ Category: "anthropic"
```

Multi-file asset structure:

```
assets/tasks/documentation/review-docs.toml
assets/tasks/documentation/review-docs.md
assets/tasks/documentation/review-docs.meta.toml
```

SHA generated from:

```bash
# Concatenate all files except .meta.toml in sorted order
cat review-docs.md review-docs.toml | git hash-object --stdin
```

Single-file asset structure:

```
assets/agents/openai/gpt4.toml
assets/agents/openai/gpt4.meta.toml
```

SHA generated from:

```bash
# Just hash the single content file
git hash-object gpt4.toml
```

## Updates

- 2025-01-17: Initial version aligned with schema
