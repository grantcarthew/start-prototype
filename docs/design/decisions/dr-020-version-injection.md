# DR-020: Binary Version Injection Strategy

**Date:** 2025-01-06
**Status:** Accepted
**Category:** Build & Distribution

## Decision

Use Go's `-ldflags` at build time to inject version information into a dedicated package variable.

## Implementation

### Version Package

Location: `internal/version/version.go`

```go
package version

// Injected at build time via -ldflags
var (
    Version   = "dev"           // Semantic version (e.g., "1.2.3")
    Commit    = "unknown"       // Git commit SHA (short, 7 chars)
    BuildDate = "unknown"       // ISO 8601 timestamp
    GoVersion = "unknown"       // Go compiler version
)

// Full returns a complete version string
func Full() string {
    return fmt.Sprintf("%s (commit: %s, built: %s, go: %s)",
        Version, Commit, BuildDate, GoVersion)
}

// Short returns just the semantic version
func Short() string {
    return Version
}
```

### Build Command

Using ldflags to inject values:

```bash
go build -ldflags "\
  -X 'github.com/grantcarthew/start/internal/version.Version=$(git describe --tags --always --dirty)' \
  -X 'github.com/grantcarthew/start/internal/version.Commit=$(git rev-parse --short HEAD)' \
  -X 'github.com/grantcarthew/start/internal/version.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)' \
  -X 'github.com/grantcarthew/start/internal/version.GoVersion=$(go version | awk '{print $3}')' \
" -o start ./cmd/start
```

### Makefile Integration

```makefile
VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GO_VERSION := $(shell go version | awk '{print $$3}')

LDFLAGS := -X 'github.com/grantcarthew/start/internal/version.Version=$(VERSION)' \
           -X 'github.com/grantcarthew/start/internal/version.Commit=$(COMMIT)' \
           -X 'github.com/grantcarthew/start/internal/version.BuildDate=$(BUILD_DATE)' \
           -X 'github.com/grantcarthew/start/internal/version.GoVersion=$(GO_VERSION)'

build:
	go build -ldflags "$(LDFLAGS)" -o start ./cmd/start

install:
	go install -ldflags "$(LDFLAGS)" ./cmd/start
```

### GoReleaser Integration

```yaml
# .goreleaser.yaml
builds:
  - main: ./cmd/start
    binary: start
    ldflags:
      - -s -w
      - -X github.com/grantcarthew/start/internal/version.Version={{.Version}}
      - -X github.com/grantcarthew/start/internal/version.Commit={{.ShortCommit}}
      - -X github.com/grantcarthew/start/internal/version.BuildDate={{.Date}}
      - -X github.com/grantcarthew/start/internal/version.GoVersion={{.Env.GOVERSION}}
```

## Version Sources

### Git Tags (Production)

- **Format:** Semantic versioning (v1.2.3)
- **Command:** `git describe --tags --always --dirty`
- **Examples:**
  - Clean tagged release: `v1.2.3`
  - Post-release commit: `v1.2.3-5-gabc1234`
  - Uncommitted changes: `v1.2.3-dirty`
  - No tags yet: `abc1234` (commit SHA)

### Development Builds

Default values when built without ldflags:
- `Version = "dev"`
- `Commit = "unknown"`
- `BuildDate = "unknown"`
- `GoVersion = "unknown"`

This allows:
```bash
go run ./cmd/start --version
# Output: start dev (commit: unknown, built: unknown, go: unknown)
```

## Version Display

### Version Flag Output

```bash
$ start --version
start 1.2.3 (commit: abc1234, built: 2025-01-06T10:30:00Z, go: go1.22.0)

$ start version
start 1.2.3 (commit: abc1234, built: 2025-01-06T10:30:00Z, go: go1.22.0)
```

### Doctor Command Output

```bash
$ start doctor
Version Information:
  CLI Version:     1.2.3
  Commit:          abc1234
  Build Date:      2025-01-06T10:30:00Z
  Go Version:      go1.22.0

Asset Information:
  Asset Version:   1.1.0 (commit: def5678)
  Last Updated:    2 days ago
  Status:          ✓ Up to date
```

## Benefits

- ✅ **Zero runtime dependencies** - No files to read, no network calls
- ✅ **Standard Go practice** - Widely used and understood
- ✅ **Single source of truth** - Git tags define version
- ✅ **Build tool agnostic** - Works with make, goreleaser, mise, etc.
- ✅ **Reproducible builds** - Version encoded in binary
- ✅ **Dev-friendly** - Works without ldflags (shows "dev")

## Trade-offs Accepted

- ❌ Requires Git during build (acceptable for standard workflow)
- ❌ Cannot change version without rebuild (expected behavior)
- ❌ Slightly longer build commands (mitigated by Makefile)

## Rationale

This approach is the Go ecosystem standard because:
1. No external dependencies at runtime
2. Version is immutable and verifiable
3. Works seamlessly with all build tools
4. Simple to implement and maintain
5. Supports both release and development builds

## Related Decisions

- [DR-011](./dr-011-asset-distribution.md) - Asset distribution (separate versioning)
- [DR-014](./dr-014-github-tree-api.md) - Asset version tracking (different mechanism)

## Implementation Notes

### Package Design

The `internal/version` package should be minimal:
- No external dependencies
- Pure data + simple formatters
- Used by `cmd/start/main.go` for `--version` flag
- Used by `cmd/start/doctor.go` for version reporting

### Version Command

Add both flag and subcommand:
```go
// Global flag
rootCmd.Version = version.Full()
rootCmd.SetVersionTemplate("{{.Version}}\n")

// Explicit subcommand (for consistency with other commands)
var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Show version information",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println(version.Full())
    },
}
```

### Build Integration Checklist

- [ ] Create `internal/version/version.go`
- [ ] Add Makefile with version injection
- [ ] Configure goreleaser with ldflags
- [ ] Add `start version` subcommand
- [ ] Add `--version` global flag
- [ ] Include version in `start doctor` output
- [ ] Document build process in README
