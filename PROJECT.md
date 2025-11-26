# PROJECT: start - AI Agent CLI Implementation

**Status:** Implementation Phase
**Current Phase:** Phase 8a (Diagnostics & CLI UX) - Complete
**Started:** 2025-11-24
**Last Updated:** 2025-11-26

---

## Quick Links

- **Architecture:** [docs/architecture.md](docs/architecture.md)
- **Testing Strategy:** [docs/testing.md](docs/testing.md)
- **Implementation Phases:** [docs/implementation/](docs/implementation/)

---

## Documentation Index

### Design Records

See [docs/design/design-records/README.md](docs/design/design-records/README.md) for the complete index of all Design Records (DR-001 through DR-041).

**Key Design Records by Phase:**

- **Configuration:** DR-001 (TOML), DR-002 (Config Merge), DR-003 (Named Documents), DR-004 (Agent Scope), DR-005 (Roles), DR-007 (Placeholders), DR-008 (File Handling), DR-012 (Context Required)
- **CLI Design:** DR-006 (Cobra), DR-017 (CLI Organization), DR-024 (Doctor Exit Codes), DR-025 (No Auto Checks), DR-028 (Shell Completion), DR-030 (Prefix Matching), DR-038 (Flag Resolution), DR-041 (Asset Commands)
- **Tasks:** DR-009 (Task Structure), DR-019 (Task Loading), DR-029 (Task Agent Field)
- **Asset Management:** DR-031 (Catalog Architecture), DR-032 (Metadata Schema), DR-033 (Resolution Algorithm), DR-034 (GitHub API), DR-035 (Interactive Browsing), DR-036 (Cache), DR-037 (Updates), DR-039 (Index), DR-040 (Substring Matching)
- **Build & Distribution:** DR-020 (Version Injection), DR-021 (GitHub Version Check), DR-022 (Asset Branch Strategy)
- **Runtime Behavior:** DR-008 (File Handling), DR-026 (Offline Behavior), DR-027 (Security Trust Model)

### CLI Documentation

Main commands:

- [start.md](docs/cli/start.md) - Main command reference and interactive mode
- [start-init.md](docs/cli/start-init.md) - Initialize configuration wizard
- [start-task.md](docs/cli/start-task.md) - Task execution
- [start-prompt.md](docs/cli/start-prompt.md) - Direct prompt execution
- [start-doctor.md](docs/cli/start-doctor.md) - System health diagnostics
- [start-show.md](docs/cli/start-show.md) - Display configuration

Asset management:

- [start-assets.md](docs/cli/start-assets.md) - Asset management overview
- [start-assets-add.md](docs/cli/start-assets-add.md) - Add assets from catalog
- [start-assets-browse.md](docs/cli/start-assets-browse.md) - Browse available assets
- [start-assets-search.md](docs/cli/start-assets-search.md) - Search catalog
- [start-assets-info.md](docs/cli/start-assets-info.md) - View asset details
- [start-assets-index.md](docs/cli/start-assets-index.md) - List cached assets
- [start-assets-update.md](docs/cli/start-assets-update.md) - Update cached assets

Configuration:

- [start-config.md](docs/cli/start-config.md) - Config management overview
- [start-config-agent.md](docs/cli/start-config-agent.md) - Agent configuration
- [start-config-role.md](docs/cli/start-config-role.md) - Role configuration
- [start-config-context.md](docs/cli/start-config-context.md) - Context configuration
- [start-config-task.md](docs/cli/start-config-task.md) - Task configuration

### Architecture & Design Documents

- [docs/architecture.md](docs/architecture.md) - System architecture (Hexagonal pattern)
- [docs/testing.md](docs/testing.md) - Testing strategy and smith agent
- [docs/design/unified-template-design.md](docs/design/unified-template-design.md) - UTD pattern (file/command/prompt)
- [docs/config.md](docs/config.md) - Complete configuration reference
- [examples/](examples/) - Configuration examples (minimal, complete, real-world)

---

## Overview

`start` is a command-line orchestrator for AI agents that manages prompt composition, context injection, and workflow automation. It wraps various AI CLI tools (Claude, Gemini, GPT, etc.) with configurable roles, reusable tasks, and project-aware context documents.

**What it does:**

- Loads project-specific context documents automatically
- Applies role-based system prompts
- Executes predefined workflow tasks
- Manages asset catalog from GitHub with lazy loading
- Supports multiple AI agents with unified configuration

---

## Technology Stack

| Component | Technology | Reason |
|-----------|------------|--------|
| **Language** | Go 1.23+ | Modern, fast, cross-platform |
| **TOML** | `github.com/pelletier/go-toml/v2` | Fast, strict mode, v1.0 compliant |
| **CLI** | `github.com/spf13/cobra` | Industry standard, feature-rich |
| **HTTP** | Standard library `net/http` | No external deps needed |
| **Build** | Go toolchain only | No make, simple bash scripts |

---

## Distribution

**Primary:** Homebrew via custom tap
**Repository:** <https://github.com/grantcarthew/homebrew-tap>

**Installation:**

```bash
brew tap grantcarthew/tap
brew install start
```

**Versioning:**

- **v0.0.x** - Early development
- **v0.1.0** - First usable release (Phase 5)
- **v0.5.0** - Asset catalog working (Phase 6)
- **v1.0.0** - Production-ready (Phase 9)

---

## Architecture

See [docs/architecture.md](docs/architecture.md) for complete architecture details.

**Pattern:** Hexagonal Architecture (Ports and Adapters)

```
CLI Layer (Cobra)
    ↓
Engine Layer (Business Logic)
    ↓
Domain Layer (Interfaces + Models)
    ↓
Adapters Layer (Concrete Implementations)
```

**Key Principles:**

- Interface-based dependency injection
- Test-first development
- Domain-driven design
- Idiomatic Go patterns

---

## Testing Strategy

See [docs/testing.md](docs/testing.md) for complete testing strategy.

**Approach:**

1. **Unit Tests:** Mock interfaces, test components in isolation
2. **Integration Tests:** Real binary + smith agent
3. **Manual Testing:** After each phase

**Smith Test Agent:**

- Deterministic test double for real AI agents
- Captures args and prompts to files
- No network calls, no API costs
- Fast and reliable

**Coverage Goals:**

- 80%+ for core logic
- 100% for critical paths

---

## Implementation Phases

### Phase Summary

| Phase | Name | Status | Effort | Link |
|-------|------|--------|--------|------|
| 0 | Foundation & Smith | ✅ Complete | 2-3h | [phase-0.md](docs/implementation/phase-0.md) |
| 1 | Config Loading & Validation | ✅ Complete | 4-6h | [phase-1.md](docs/implementation/phase-1.md) |
| 2 | Simple Agent Execution | ✅ Complete | 4-6h | [phase-2.md](docs/implementation/phase-2.md) |
| 3 | Roles & Contexts | ✅ Complete | 5-7h | [phase-3.md](docs/implementation/phase-3.md) |
| 4 | UTD Pattern Processing | ✅ Complete | 6-8h | [phase-4.md](docs/implementation/phase-4.md) |
| 5 | Tasks | ✅ Complete | 4-6h | [phase-5.md](docs/implementation/phase-5.md) |
| 6 | Asset Catalog & Lazy Loading | ✅ Complete | 8-10h | [phase-6.md](docs/implementation/phase-6.md) |
| 7 | Init & Asset Management | ✅ Complete | 8-10h | [phase-7.md](docs/implementation/phase-7.md) |
| 8a | Diagnostics & CLI UX | ✅ Complete | 3-4h | [phase-8.md](docs/implementation/phase-8.md) |
| 8b | Config Management Commands | Not Started | 8-10h | [phase-8.md](docs/implementation/phase-8.md) |
| 9 | Polish & Documentation | Not Started | 6-8h | [phase-9.md](docs/implementation/phase-9.md) |

### Phase Descriptions

**Phase 0:** Project scaffolding, smith agent, domain models, test harness

**Phase 1:** TOML config loading, merging (global + local), validation
- ✅ Config loader (global + local directories)
- ✅ Config merger (local overrides global)
- ✅ Comprehensive validation (agents, roles, contexts, tasks, settings)
- ✅ CLI structure with `start config show` command
- ✅ 13 unit tests + 3 integration tests (all passing)
- ✅ Fixed example configs (minimal, complete, real-world)

**Phase 2:** Basic agent execution with simple placeholder resolution
- ✅ RealRunner adapter (os/exec wrapper)
- ✅ Placeholder resolver ({bin}, {model}, {prompt}, {date})
- ✅ Executor (command construction and execution)
- ✅ Root command RunE (agent/model selection, execution)
- ✅ 14 unit tests + 2 integration tests (all passing)

**Phase 3:** Role system prompts and context document loading
- ✅ UTD processor (file/command/prompt resolution)
- ✅ CommandRunner adapter (command execution with output capture)
- ✅ Role selector (precedence: flag → task → default)
- ✅ Role loader (UTD processing, temp file management)
- ✅ Context loader (required field filtering, order preservation)
- ✅ Updated executor (ExecuteParams, {role}/{role_file} placeholders)
- ✅ Root command integration (--role flag, role/context loading)
- ✅ Context order preservation (ContextOrder field)
- ✅ 9 UTD unit tests + updated executor/integration tests (34 tests passing)

**Phase 4:** UTD pattern processing (completed with Phase 3, shell auto-detection added)
- ✅ UTD file/command/prompt parsing (completed in Phase 3)
- ✅ Command execution with shell integration (completed in Phase 3)
- ✅ UTD placeholders ({file}, {file_contents}, {command_output}) (completed in Phase 3)
- ✅ Shell configuration and timeouts (completed in Phase 3)
- ✅ Temporary file handling (completed in Phase 3)
- ✅ Shell auto-detection (DetectShell function - bash if available, otherwise sh)
- ✅ All tests passing (35 tests total)

**Phase 5:** Task system with {instructions} placeholder and auto-contexts
- ✅ TaskLoader (UTD processing with {instructions} placeholder)
- ✅ TaskResolver (name/alias resolution with local override)
- ✅ Task command (start task <name> [instructions])
- ✅ Task listing (start task with no arguments)
- ✅ Agent/role selection precedence (CLI flag → task field → default → first)
- ✅ Auto-include required contexts only (CommandTypeTask)
- ✅ {instructions} placeholder (defaults to "None" when empty)
- ✅ 13 unit tests + 4 integration tests (43 tests total passing)

**Phase 6:** GitHub catalog integration, asset resolution, lazy loading, caching
- ✅ RealGitHubClient adapter (HTTP client for raw.githubusercontent.com)
- ✅ FileCache adapter (filesystem-based cache with .meta.toml sidecar files)
- ✅ Catalog index parser (CSV to AssetMeta, search, filter functions)
- ✅ Asset resolver (DR-033 resolution algorithm: local → global → cache → GitHub)
- ✅ CLI commands (start assets search, start assets add)
- ✅ Integration into main.go (dependency injection wiring)
- ✅ 20+ unit tests for catalog parser (all passing)
- ✅ All existing tests passing (43 tests total)

**Phase 7:** `start init` wizard, agent detection, complete asset management commands
- ✅ InitCommand (interactive wizard with location selection, agent detection, config generation)
- ✅ Auto-detect installed agents (exec.LookPath integration)
- ✅ Priority-based default agent selection (claude > gemini > others)
- ✅ Timestamped config backups (YYYY-MM-DD-HHMMSS format)
- ✅ Multi-file config generation (config.toml, agents.toml, roles.toml, contexts.toml, tasks.toml)
- ✅ Default context documents (ENVIRONMENT.md, INDEX.csv, AGENTS.md, PROJECT.md)
- ✅ start assets browse (open GitHub catalog in browser)
- ✅ start assets info <query> (detailed asset information with interactive selection)
- ✅ start assets update [query] (SHA-based update detection, selective updates)
- ✅ start assets index (catalog index generation for maintainers)
- ✅ Updated start assets add (query-based search per DR-041, interactive selection)
- ✅ Flags: --local (project config), --force (fully automatic), --yes (skip prompts)
- ✅ Unit tests for init helper functions (4 tests, all passing)
- ✅ Unit tests for catalog functions (20+ tests, all passing)
- ✅ Integration tests for start init command (4 tests: force, local, backup, help)
- ✅ Integration tests for start assets commands (12 tests: search, add, browse, info, update, index, help, validation)
- ✅ Query length validation (minimum 3 characters per DR-040)
- ✅ DR-035 updated to match implementation (query-based mode documented)
- ✅ All tests passing or skipping gracefully (116 total: 91 unit tests + 25 integration tests)

**Phase 8a:** Diagnostics & CLI UX enhancements (prefix matching, completion, doctor, version checking)
- ✅ Prefix matching (DR-030: cobra.EnablePrefixMatching for abbreviated commands)
- ✅ Shell completion (DR-028: bash/zsh/fish completion generation and auto-install)
- ✅ Version checker (DR-021: GitHub API integration with semver comparison, rate limiting, GH_TOKEN support)
- ✅ Doctor command (DR-024: comprehensive health checks with exit codes, --quiet and --verbose modes)
- ✅ Health checks: version, assets age, config validation, agent binaries, contexts, environment
- ✅ Version comparison with development build detection
- ✅ Installation method detection (Homebrew, go install, manual)
- ✅ Unit tests for version checker (10 tests, all passing)
- ✅ Integration tests for Phase 8a (8 tests: prefix matching, completion output, doctor modes)
- ✅ All tests passing (124 total: 103 unit tests + 21 integration tests)

**Phase 8b:** Config management commands (interactive TOML manipulation)
- start config agent (list/new/show/test/edit/remove/default) - DR-038
- start config role (list/new/show/test/edit/remove/default)
- start config context (list/new/show/test/edit/remove/default)
- start config task (list/new/show/test/edit/remove/default)
- Interactive wizards with validation, backup creation, scope selection (global vs local)
- TOML file manipulation with preserve-format editing

**Phase 9:** Error messages, output formatting, performance, documentation, v1.0.0 release

---

## Development Workflow

### Build

```bash
# Build main binary
go build -o bin/start cmd/start/main.go

# Build with version
go build -ldflags "-X main.version=0.0.1" -o bin/start cmd/start/main.go

# Build smith
go build -o bin/smith cmd/smith/main.go
```

### Test

```bash
# All tests
./test.sh

# Unit tests only
go test -short ./...

# Integration tests only
go test ./test/integration/...

# With coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Git Workflow

**Branches:**

- `main` - stable code
- `phase-N-description` - phase implementation branches

**Commits:**

- Format: `type: description`
- Types: feat, fix, docs, test, refactor, chore

**Example:**

```bash
git checkout -b phase-0-foundation
# ... implement phase 0 ...
git commit -m "feat: add smith agent"
git commit -m "feat: add domain models"
git commit -m "test: add mock implementations"
./test.sh
git checkout main
git merge phase-0-foundation
git tag phase-0-complete
```

### Release Workflow

1. Update version in code
2. Update CHANGELOG.md
3. Commit: `git commit -m "chore: bump version to X.Y.Z"`
4. Tag: `git tag vX.Y.Z`
5. Push: `git push origin main --tags`
6. Build binaries (macOS, Linux)
7. Create GitHub release
8. Update Homebrew formula

---

## Progress Tracking

### Current Status

**Phase:** 8a (Diagnostics & CLI UX) - ✅ Complete
**Last Completed:** Phase 8a (2025-11-26)
**Next Phase:** Phase 8b (Config Management Commands) or Phase 9 (Polish & Documentation)
**Next Milestone:** v0.5.0 (Phase 6-8a complete) - Full feature set ✅ READY

### Phase Checklist

- [x] Phase 0: Foundation & Smith
- [x] Phase 1: Config Loading & Validation
- [x] Phase 2: Simple Agent Execution
- [x] Phase 3: Roles & Contexts
- [x] Phase 4: UTD Pattern Processing
- [x] Phase 5: Tasks
- [x] Phase 6: Asset Catalog & Lazy Loading
- [x] Phase 7: Init & Asset Management
- [x] Phase 8a: Diagnostics & CLI UX
- [ ] Phase 8b: Config Management Commands
- [ ] Phase 9: Polish & Documentation

### Milestone Targets

- **v0.0.1:** Phase 2-3 complete (basic execution with roles/contexts) - First usable version
- **v0.1.0:** Phase 5 complete (tasks working) - Core functionality
- **v0.5.0:** Phase 6-8a complete (catalog, init, diagnostics) - Full feature set ✅ READY
- **v1.0.0:** Phase 9 complete (production ready) - Stable release

---

## Design Records Reference

All implementation aligns with Design Records in `docs/design/design-records/`.

**Key DRs:**

- DR-001: TOML format
- DR-002: Config merge strategy
- DR-005: Role configuration
- DR-006: Cobra CLI structure
- DR-007: Placeholders
- DR-009: Task structure
- DR-031: Catalog-based assets
- DR-033: Asset resolution algorithm

See [docs/design/design-records/README.md](docs/design/design-records/README.md) for complete index.

---

## Getting Started

### For Implementation

1. Read [docs/architecture.md](docs/architecture.md)
2. Read [docs/testing.md](docs/testing.md)
3. Start with [Phase 0](docs/implementation/phase-0.md)
4. Follow phases in order
5. Test thoroughly after each phase

### For Understanding

1. Read AGENTS.md for repository context
2. Review Design Records for decisions
3. Check CLI documentation in `docs/cli/`
4. Understand UTD pattern in `docs/design/unified-template-design.md`

---

## Notes

**This is implementation guidance.** Each phase document contains detailed tasks, testing criteria, and completion checklists.

**Stop points:** After each phase, stop for manual testing before proceeding.

**Documentation-driven:** All behavior specified in Design Records. If unclear, update DR first.

---

_Document Status: In Progress_
_Last Updated: 2025-11-26_
_Phase 0 Complete: Yes_
_Phase 1 Complete: Yes_
_Phase 2 Complete: Yes_
_Phase 3 Complete: Yes_
_Phase 4 Complete: Yes_
_Phase 5 Complete: Yes_
_Phase 6 Complete: Yes_
_Phase 7 Complete: Yes_
_Phase 8a Complete: Yes_
_Phase 8b Complete: No_
