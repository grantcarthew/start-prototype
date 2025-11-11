# DR-032: Asset Metadata Schema

**Date:** 2025-01-10
**Status:** Accepted
**Category:** Asset Management

## Decision

Use sidecar `.meta.toml` files for asset metadata with 6 required fields, deriving category from filesystem path to prevent drift.

## What This Means

### Sidecar Metadata Files

**Each catalog asset has a metadata file:**
```
assets/tasks/git-workflow/
├── pre-commit-review.toml        # Asset content
├── pre-commit-review.md          # Prompt file (if using UTD)
└── pre-commit-review.meta.toml   # Metadata (sidecar)
```

**Why sidecar, not frontmatter?**
- ✅ Content files stay clean (no metadata clutter for AI to process)
- ✅ Metadata separate from what agent sees
- ✅ Easy to parse independently
- ✅ No risk of corrupting content when updating metadata
- ✅ Consistent format across asset types (roles, tasks, agents, contexts)
- ❌ More files to manage (acceptable trade-off for cleanliness)

### Metadata Schema

**All fields required (no optional fields):**

```toml
# pre-commit-review.meta.toml
name = "pre-commit-review"
description = "Review staged changes before committing"
tags = ["git", "review", "quality", "pre-commit"]
sha = "a1b2c3d4e5f6789012345678901234567890abcd"
created = "2025-01-10T00:00:00Z"
updated = "2025-01-10T00:00:00Z"
```

### Field Definitions

**name** (string)
- Asset identifier (must match filename without extensions)
- Example: `"pre-commit-review"` for `pre-commit-review.toml`
- Used for lookups and validation

**description** (string)
- One-line summary of what the asset does
- Shown in catalog browsing and search results
- Example: `"Review staged changes before committing"`

**tags** (array of strings)
- Keywords for searching and filtering
- At least 1 tag required
- Example: `["git", "review", "quality", "pre-commit"]`

**sha** (string)
- Git blob SHA of the content file (not metadata file)
- 40-character hexadecimal string
- Used for update detection via comparison
- Example: `"a1b2c3d4e5f6789012345678901234567890abcd"`

**created** (ISO 8601 timestamp)
- When asset was first created
- Example: `"2025-01-10T00:00:00Z"`

**updated** (ISO 8601 timestamp)
- When asset was last modified
- Updated whenever content changes (and SHA changes)
- Example: `"2025-01-10T12:30:00Z"`

### Category Derived from Filesystem

**No category field** - Derive from directory structure:

```
assets/tasks/git-workflow/pre-commit-review.toml
       ^     ^
       |     └─ category = "git-workflow"
       └─────── asset type = "tasks"
```

**Rationale:**
- Eliminates potential drift between category field and actual path
- Filesystem IS the source of truth
- Simpler validation (one less field to check)
- Category can't be wrong if it's derived

### Validation Rules

**Required field validation:**
```go
func ValidateMetadata(meta *AssetMetadata, filePath string) error {
    if meta.Name == "" { return ErrMissingName }
    if meta.Description == "" { return ErrMissingDescription }
    if len(meta.Tags) == 0 { return ErrMissingTags }
    if meta.SHA == "" { return ErrMissingSHA }
    if meta.Created.IsZero() { return ErrMissingCreated }
    if meta.Updated.IsZero() { return ErrMissingUpdated }
    return nil
}
```

**SHA format validation:**
```go
func ValidateSHA(sha string) error {
    if len(sha) != 40 {
        return fmt.Errorf("SHA must be 40 characters, got %d", len(sha))
    }
    if !regexp.MustCompile(`^[a-f0-9]{40}$`).MatchString(sha) {
        return fmt.Errorf("SHA must be hexadecimal: %s", sha)
    }
    return nil
}
```

**Name match validation:**
```go
func ValidateNameMatch(meta *AssetMetadata, filePath string) error {
    filename := filepath.Base(filePath)
    expected := strings.TrimSuffix(filename, ".meta.toml")
    if meta.Name != expected {
        return fmt.Errorf("name %q does not match filename %q", meta.Name, expected)
    }
    return nil
}
```

**Timestamp validation:**
```go
func ValidateTimestamps(meta *AssetMetadata) error {
    if meta.Updated.Before(meta.Created) {
        return fmt.Errorf("updated timestamp (%v) before created timestamp (%v)",
            meta.Updated, meta.Created)
    }
    return nil
}
```

## Implementation

### Go Struct

```go
type AssetMetadata struct {
    Name        string    `toml:"name"`
    Description string    `toml:"description"`
    Tags        []string  `toml:"tags"`
    SHA         string    `toml:"sha"`
    Created     time.Time `toml:"created"`
    Updated     time.Time `toml:"updated"`

    // Derived from filesystem (not in TOML)
    Category    string    `toml:"-"`
    AssetType   string    `toml:"-"`
    FilePath    string    `toml:"-"`  // Full path to .meta.toml
}
```

### Loading Metadata

```go
func LoadMetadata(metaPath string) (*AssetMetadata, error) {
    // Parse TOML
    var meta AssetMetadata
    if err := toml.DecodeFile(metaPath, &meta); err != nil {
        return nil, fmt.Errorf("parse metadata: %w", err)
    }

    // Derive category and type from path
    // e.g., "assets/tasks/git-workflow/pre-commit-review.meta.toml"
    parts := strings.Split(filepath.Dir(metaPath), string(filepath.Separator))
    if len(parts) >= 2 {
        meta.AssetType = parts[len(parts)-2]  // "tasks"
        meta.Category = parts[len(parts)-1]   // "git-workflow"
    }
    meta.FilePath = metaPath

    // Validate
    if err := ValidateMetadata(&meta, metaPath); err != nil {
        return nil, err
    }

    return &meta, nil
}
```

### SHA Generation

**For asset creators:**
```bash
# Generate SHA for content file
sha=$(git hash-object pre-commit-review.toml)
echo "sha = \"$sha\"" >> pre-commit-review.meta.toml
```

**Or in Go:**
```go
func GenerateContentSHA(filePath string) (string, error) {
    content, err := os.ReadFile(filePath)
    if err != nil {
        return "", err
    }

    // Git blob SHA: "blob <size>\0<content>"
    header := fmt.Sprintf("blob %d\x00", len(content))
    data := append([]byte(header), content...)

    hash := sha1.Sum(data)
    return hex.EncodeToString(hash[:]), nil
}
```

## Examples

### Task Metadata

```toml
# assets/tasks/git-workflow/pre-commit-review.meta.toml
name = "pre-commit-review"
description = "Review staged changes before committing"
tags = ["git", "review", "quality", "pre-commit"]
sha = "a1b2c3d4e5f6789012345678901234567890abcd"
created = "2025-01-10T00:00:00Z"
updated = "2025-01-10T12:30:00Z"
```

### Role Metadata

```toml
# assets/roles/general/code-reviewer.meta.toml
name = "code-reviewer"
description = "Strict quality and security focused code reviewer"
tags = ["review", "quality", "security", "strict"]
sha = "b2c3d4e5f6789012345678901234567890abcde1"
created = "2025-01-10T00:00:00Z"
updated = "2025-01-10T00:00:00Z"
```

### Agent Metadata

```toml
# assets/agents/claude/sonnet.meta.toml
name = "sonnet"
description = "Balanced Claude Sonnet model for general use"
tags = ["claude", "balanced", "recommended", "default"]
sha = "c3d4e5f6789012345678901234567890abcde12"
created = "2025-01-10T00:00:00Z"
updated = "2025-01-10T00:00:00Z"
```

## Benefits

**Simplicity:**
- ✅ Only 6 fields, all required
- ✅ No complex optional fields or conditionals
- ✅ Clear validation rules

**No drift:**
- ✅ Category derived from filesystem (can't be wrong)
- ✅ Name must match filename (enforced validation)
- ✅ SHA tracks actual content

**Update detection:**
- ✅ SHA comparison for reliable change detection
- ✅ No version numbers to maintain
- ✅ Content hash is the version

**Searchability:**
- ✅ Tags enable filtering without full-text indexing
- ✅ Description for quick scanning
- ✅ Structured data for catalog browsing

## Trade-offs Accepted

**No semantic versioning:**
- ❌ SHA is not human-readable like semver
- **Mitigation:** `updated` timestamp provides context, SHA is reliable

**No author/license fields:**
- ❌ Can't track who created asset
- **Mitigation:** Add in future if community contributions happen

**No minimum version requirement:**
- ❌ Can't specify "requires start >= 0.2.0"
- **Mitigation:** Add when needed, keep v1 simple

**Timestamps required:**
- ❌ Manual creation requires setting timestamps
- **Mitigation:** Simple ISO 8601 format, tooling can help

## Future Considerations

**If we need more metadata:**

```toml
# Potential future fields (not in v1)
author = "start-project"           # Creator
license = "MIT"                    # License type
min_start_version = "0.2.0"        # Minimum CLI version
deprecated = false                 # Mark as deprecated
replacement = "new-task-name"      # If deprecated
dependencies = ["roles/reviewer"]  # Required assets
```

**Current stance:** KISS for v1. Add fields when actually needed, not speculatively.

## Related Decisions

- [DR-031](./dr-031-catalog-based-assets.md) - Catalog architecture (metadata usage context)
- [DR-034](./dr-034-github-catalog-api.md) - GitHub API (how SHA is obtained)
- [DR-037](./dr-037-asset-updates.md) - Update mechanism (SHA comparison)
