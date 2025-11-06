# Project: start - Configuration Design Phase

**Status:** Design Phase - Implementation Details & Command Specs
**Date Started:** 2025-01-03
**Current Phase:** Asset management system design, CLI reorganization, implementation details

**Working Mode:** Interactive design sessions - walk through each implementation task collaboratively, discussing design choices before creating design records.

## Overview

Context-aware AI agent launcher that detects project context, builds intelligent prompts, and launches AI development tools with proper configuration.

**Links:** [Vision](./docs/vision.md) | [Config Reference](./docs/config.md) | [UTD](./docs/design/unified-template-design.md) | [Design Decisions](./docs/design/design-record.md) (25 DRs) | [Tasks](./docs/tasks.md)

## Command Status

### Execution Commands

- ✅ `start` - [docs/cli/start.md](./docs/cli/start.md) - Launch with all context
- ✅ `start prompt [text]` - [docs/cli/start-prompt.md](./docs/cli/start-prompt.md) - Launch with required context + custom prompt
- ✅ `start task <name> [inst]` - [docs/cli/start-task.md](./docs/cli/start-task.md) - Run predefined task

### Configuration Management

**File Operations:**
- ✅ `start config show` - [docs/cli/start-config.md](./docs/cli/start-config.md) - View merged config
- ✅ `start config edit [scope]` - [docs/cli/start-config.md](./docs/cli/start-config.md) - Edit config file
- ✅ `start config path` - [docs/cli/start-config.md](./docs/cli/start-config.md) - Show config paths
- ✅ `start config validate` - [docs/cli/start-config.md](./docs/cli/start-config.md) - Validate config

**Configuration Sections:**
- ✅ `start config agent` - [docs/cli/start-agent.md](./docs/cli/start-agent.md) - Manage agents (MOVED from `start agent`)
- ✅ `start config context` - [docs/cli/start-config-context.md](./docs/cli/start-config-context.md) - Manage contexts
- ✅ `start config task` - [docs/cli/start-config-task.md](./docs/cli/start-config-task.md) - Manage tasks
- ✅ `start config role` - [docs/cli/start-config-role.md](./docs/cli/start-config-role.md) - Manage system prompts

### Utility Commands

- ✅ `start init [scope]` - [docs/cli/start-init.md](./docs/cli/start-init.md) - Initialize configuration
- ✅ `start doctor` - [docs/cli/start-doctor.md](./docs/cli/start-doctor.md) - Diagnose installation
- ✅ `start update` - [docs/cli/start-update.md](./docs/cli/start-update.md) - Update asset library

## Architecture Decisions

### Completed (25 Design Records)

**Core Configuration (DR-001 to DR-008):**
- DR-001: TOML for configuration format
- DR-002: Global + local config file structure with merge
- DR-003: Named context documents (not arrays)
- DR-004: Agents in both global and local configs
- DR-005: System prompt handling (separate and optional)
- DR-006: Cobra CLI with subcommands
- DR-007: Command interpolation and placeholders
- DR-008: Context file detection and handling

**Tasks & Commands (DR-009 to DR-013):**
- DR-009: Task structure and placeholders
- DR-010: Four default interactive review tasks
- DR-011: Asset distribution and update system (GitHub-fetched)
- DR-012: Context document required field and order
- DR-013: Agent configuration distribution via GitHub

**Asset Management & CLI (DR-014 to DR-025):**
- DR-014: GitHub Tree API with SHA-based caching for incremental updates
- DR-015: Atomic update mechanism with rollback capability
- DR-016: Asset discovery - each feature checks its own directory
- DR-017: CLI reorganization - `start config` for all configuration management
- DR-018: Init and update command integration - shared implementation, no conditional logic
- DR-019: Task loading and merging - global + local only, assets as templates, local precedence
- DR-020: Binary version injection - ldflags at build time with git-based versioning
- DR-021: GitHub version checking - Releases API with no caching, integrated into doctor and update commands
- DR-022: Asset branch strategy - pull from main branch (not releases) for rapid iteration
- DR-023: Asset staleness checking - GitHub commit comparison with no caching
- DR-024: Doctor exit code system - simple binary exit codes (0 = healthy, 1 = issues)
- DR-025: No automatic checks or caching - all health checks are user-initiated only

**Implementation:**
- Unified Template Design (UTD): `file`, `command`, `prompt` pattern across all sections
- Document order: Config definition order (TOML preserves order)
- Command pattern: Positional scope arguments (`start init [scope]`, `start config edit [scope]`)
- CLI Framework: Cobra with dynamic task loading
- Asset file: `asset-version.toml` tracks commit SHA + file SHAs

### Key Design

- `start` includes ALL documents (required + optional)
- `start prompt` includes ONLY required documents
- Model aliases user-defined per agent (not hardcoded tiers)
- Missing files show but don't error
- Verbosity: quiet/normal/verbose/debug

## Open Questions

### Resolved
1. ✅ **Task listing:** `start task` with no args + `--help`
2. ✅ **Agent testing:** Binary availability, config validation, dry-run display
3. ✅ **Config editing:** Soft warnings (non-blocking)
4. ✅ **JSON output:** Not needed, human-facing tool
5. ✅ **Task structure:** Full UTD for system_prompt and task prompt
6. ✅ **Context management:** `start config context` (DR-017)
7. ✅ **Role management:** `start config role` (DR-017)

### Remaining

8. **Shell completion:** Generate for bash/zsh/fish?
9. **Non-interactive mode:** What flags needed for CI/automation?
10. **Version tracking:** Build-time injection strategy
11. **Asset staleness:** Local-only check vs GitHub comparison
12. **Security:** Trust model for downloaded assets

## Success Criteria

CLI design is complete when:

- [x] All core commands fully specified (start, prompt, init, task)
- [x] All management commands specified (agent ✅, config ✅)
- [x] All high-priority questions resolved (3 of 3 done)
- [x] Error cases documented across commands
- [x] Output formats specified consistently
- [x] Patterns consistent across all commands

## Reference

**Doc Template:** See [docs/cli/start.md](./docs/cli/start.md) for complete documentation structure

**Key Documents:**

- [docs/vision.md](./docs/vision.md) - Product vision and goals
- [docs/config.md](./docs/config.md) - Complete configuration reference
- [docs/unified-template-design.md](./docs/unified-template-design.md) - UTD pattern (file/command/prompt)
- [docs/design-record.md](./docs/design-record.md) - All design decisions (DR-001 through DR-013+)
- [docs/task.md](./docs/task.md) - Task configuration details
- [docs/archive/](./docs/archive/) - Design discussion history

## Recent Progress

### Task Loading Algorithm Design - DR-019 (2025-01-06)

**Task 16a-c Completed:**
- ✅ Designed complete task loading and merging algorithm
- ✅ Created DR-019: Task loading and merging
- ✅ Clarified asset tasks as templates (not runtime loaded)

**Key Decisions:**
- Runtime tasks load from global + local configs only (assets NOT auto-loaded)
- Asset tasks serve as templates for `start config task add` workflow
- Local completely replaces global for same task name (no field merging)
- Resolution priority: local name → local alias → global name → global alias
- Source metadata tracked for transparency/security
- `start config task list` shows all three: global, local, and available assets

**Task Adding Workflow:**
- Interactive: "Start from template or create new?"
- If template: prompts show asset values as defaults (easy to accept)
- Conflict warnings for name/alias shadowing

**Alias Handling:**
- Local alias wins over global alias (silent override with warning)
- Conflict warning during add: "Your task will shadow the global alias"
- Runtime warning: "Global alias 'cr' shadowed by local task"

**Benefits:**
- Clear mental model: assets are templates, not active tasks
- Predictable: local always wins
- Secure: source shown at runtime (detect malicious local configs)
- Discoverable: asset templates visible in list

**Implementation Ready:**
- `internal/config/TaskRegistry` with Load/Resolve/CheckConflicts
- Clear algorithm for name/alias resolution
- Metadata tracking structure defined

### Init/Update Integration Design - DR-018 (2025-01-06)

**Task 15a Completed:**
- ✅ Defined relationship between `start init` and `start update` commands
- ✅ Created DR-018: Init and update command integration
- ✅ Simplified approach: No conditional logic, no special flags

**Key Decisions:**
- `start init` **always** invokes asset update logic (no staleness checks)
- Both commands share identical SHA-based update implementation from DR-014
- No flags: `--skip-assets`, `--force`, etc. (KISS principle)
- Network failure during init: warns but succeeds (can update later)
- Network failure during update: fails (its explicit purpose)

**Implementation Strategy:**
```go
// Package: internal/assets
func UpdateAssets() error {
    // Single source of truth for both init and update
    // Uses SHA comparison to download only changed files
}
```

**Benefits:**
- Self-optimizing: SHA comparison prevents redundant downloads
- Predictable: init always tries to update, update requires network
- Offline-friendly: init can proceed without assets
- Simple: removed all conditional logic and staleness checks

**Next Steps:**
- Task 15b: Design offline fallback strategy
- Task 15c: Define behavior when network unavailable

### Documentation Updates - config.md (2025-01-06)

**Agent Section Updates:**
- Updated to reflect DR-004 change allowing agents in both global and local configs
- Documented merge behavior: local overrides global for same agent name
- Added use case: teams can commit `.start/` with standard configs
- Updated all scope references and validation rules

**Task Section Updates:**
- Completely rewrote to use Unified Template Design (UTD) pattern
- Documented system prompt override: `system_prompt_file`, `system_prompt_command`, `system_prompt`
- Updated task prompt fields: `file`, `command`, `prompt` (UTD)
- Added shell configuration: `shell`, `command_timeout`
- Documented auto-inclusion of `required = true` contexts (no `documents` array)
- Removed old field references: `role`, `documents`, `content_command`

**Consistency Fixes:**
- Fixed placeholder documentation (`{command}` not `{content}` for tasks)
- Updated path examples (consistent use of `file` field)
- Updated validation rules for contexts (UTD pattern)
- Verified merge behavior sections throughout

### Configuration Design (2025-01-05)

**Unified Template Design (UTD):**
- Created `docs/unified-template-design.md` - Consistent pattern for file/command/prompt across all sections
- Fields: `file`, `command`, `prompt` with `{file}` and `{command}` placeholders
- Shell configuration: Global `shell` setting, per-section override, supports bash/node/python/bun/deno/etc
- Command timeout: Global `command_timeout`, per-section override

**Config Sections Completed:**
- ✅ `[settings]` - default_agent, log_level, shell, command_timeout
- ✅ `[agents.<name>]` - Full design with models, env, validation (global + local)
- ✅ `[system_prompt]` - Uses UTD pattern (file, command, prompt)
- ✅ `[context.<name>]` - Uses UTD pattern with required/description fields
- ✅ `[tasks.<name>]` - Full UTD for system_prompt_* and task prompt, auto-includes required contexts

**Config Section Naming:**
- Changed `[context.documents.<name>]` → `[context.<name>]`
- Renamed `path` attribute → `file` (UTD standard)
- Renamed `verbosity` → `log_level`

**Updated Documentation:**
- `docs/config.md` now references UTD, removed duplication
- All examples updated to use new field names

### Tasks Design & Documentation Updates (2025-01-05)

**Tasks Configuration Finalized:**
- ✅ Full UTD pattern for `system_prompt_*` fields (file, command, prompt)
- ✅ Full UTD pattern for task prompt fields (file, command, prompt)
- ✅ Auto-includes contexts where `required = true` (no `documents` array)
- ✅ `{instructions}` placeholder only in task prompts, not system prompts
- ✅ Tasks can be in both global and local configs (merge/override)
- ✅ Alias conflict resolution: First in TOML order wins
- ✅ Updated `docs/tasks.md` with complete specification

**Agent Scope Updated (DR-004):**
- ✅ Changed from global-only to global + local support
- ✅ Enables team standardization via committed `.start/` directory
- ✅ Local agents override global for same name (merge behavior)
- ✅ Security note: Don't commit secrets, use env var references
- ✅ Updated `docs/design/design-record.md` DR-004

**Command Pattern Finalized:**
- ✅ Positional scope arguments: `start init [scope]` and `start config edit [scope]`
- ✅ Scopes: `global` (default) or `local`
- ✅ Smart behavior when no scope: Interactive prompts with recommendations
- ✅ Explicit scope skips prompts for scripting/automation

**Smart Init Behavior:**
- ✅ Scenario 1: No configs → Ask, default to global
- ✅ Scenario 2: Global exists → Ask to replace global or create local
- ✅ Scenario 3: Both exist → Ask which to replace
- ✅ Scenario 4: Local exists → Ask to create global or replace local
- ✅ Always recommends global as default/safe choice

### Asset Management System & CLI Reorganization (2025-01-06)

**CLI Command Reorganization (DR-017):**
- ✅ Identified inconsistency: `start agent` (config) vs `start task` (execution)
- ✅ Decided: Configuration management under `start config`, execution at top level
- ✅ New structure: `start config agent|context|task|role` for all config management
- ✅ Benefits: Clear separation of purpose, consistent patterns, better discoverability
- ✅ Breaking change acceptable (design phase, no existing users)

**Asset Library Design:**
- ✅ Assets fetched from GitHub repository (not embedded in binary)
- ✅ Stored in `~/.config/start/assets/` directory
- ✅ Version tracked in `asset-version.toml` file (commit SHA + file SHAs)
- ✅ `start init` performs initial asset download
- ✅ `start update` refreshes asset library from GitHub
- ✅ `start doctor` checks asset age and reports if stale (> 30 days)

**Asset Types:**
- ✅ **Agents** (`assets/agents/*.toml`) - Templates used during `start agent add`
- ✅ **Roles** (`assets/roles/*.md`) - Referenced in config, updates flow automatically
- ✅ **Tasks** (`assets/tasks/*.toml`) - Merged with user tasks, user tasks take precedence
- ✅ **Examples** (`assets/examples/*.toml`) - Reference configs, not auto-loaded

**Update Flow:**
- ✅ User runs `start doctor` → Sees "Assets 45 days old"
- ✅ User runs `start update` → Downloads latest from GitHub
- ✅ Role file references automatically use updated content
- ✅ New tasks immediately available in `start task` list
- ✅ User config never modified automatically

**Design Decisions:**
- ✅ Updated DR-011 to reflect GitHub-fetched assets
- ✅ Separation of binary (code) vs content (assets)
- ✅ Users control update timing (not forced)
- ✅ Offline work after initial download
- ✅ Network dependency acceptable for updates

**Documentation Updates (Complete - 6/6):**
- ✅ `docs/cli/start-init.md` - Added scope argument, smart behavior, local support
- ✅ `docs/cli/start-config.md` - Changed flags to positional args, updated agent scope info
- ✅ `docs/config.md` - Updated agent section (local scope support), task section (UTD fields), merge behaviors, validation rules
- ✅ `docs/cli/start-task.md` - Updated to UTD pattern (system_prompt_*, command, {command} placeholder, auto-include required contexts)
- ✅ `docs/cli/start-doctor.md` - Comprehensive health check command (version, assets, config, agents, contexts, environment)
- ✅ `docs/cli/start-update.md` - Asset library update from GitHub (agents, roles, tasks, examples)
- ✅ `docs/design/design-record.md` DR-011 - Updated to reflect GitHub-fetched assets (not embedded)

**High-Level Design Complete:**
- ✅ Review/update `docs/cli/start-task.md` for task UTD fields
- ✅ Design `start doctor` and `start update` commands (user-facing behavior)

**Implementation Details to Design:**

*Asset Update Mechanism:*
- [x] **Task 12a:** Decide GitHub download strategy → DR-014: GitHub Tree API with SHA caching
- [x] **Task 12b:** Design atomic update mechanism → DR-015: SHA-filtered incremental + batch atomic install
- [x] **Task 12c:** Define asset discovery system → DR-016: No discovery system, each feature checks its directory

*Command Reorganization:*
- [x] **Task 12d:** CLI command reorganization → DR-017: `start config` for all configuration management
- [x] **Task 12e:** Update `start-agent.md` to reflect `start config agent` (path change only)
- [x] **Task 12f:** Create `start-config-context.md` spec (NEW command)
- [x] **Task 12g:** Create `start-config-task.md` spec (NEW command)
- [x] **Task 12h:** Create `start-config-role.md` spec (NEW command)

*Version Tracking & Checking:*
- [x] **Task 13a:** Define binary version source (build-time injection strategy) → DR-020: ldflags with git tags
- [x] **Task 13b:** Design GitHub version checking (API endpoint, rate limiting, caching) → DR-021: Releases API, no caching
- [x] **Task 13c:** Define commit SHA retrieval strategy (releases vs commits) → DR-022: Main branch for assets

*Doctor Implementation:*
- [x] **Task 14a:** Design asset staleness checking (local-only vs GitHub comparison) → DR-023: GitHub commit comparison
- [x] **Task 14b:** Define exit code priority system (multiple simultaneous issues) → DR-024: Simple binary exit codes
- [x] **Task 14c:** Design automatic check frequency and caching strategy → DR-025: No automatic checks or caching

*Integration & Offline Support:*
- [x] **Task 15a:** Define start init + start update relationship → DR-018: Shared implementation, no conditional logic
- [ ] **Task 15b:** Design offline fallback strategy (manual asset installation)
- [ ] **Task 15c:** Define behavior when network unavailable

*Task Merging Implementation:*
- [x] **Task 16a:** Design task loading and merging algorithm (assets + user config) → DR-019: Global + local only, assets as templates
- [x] **Task 16b:** Define source metadata tracking ([default] vs [user] labels) → DR-019: Source and SourcePath metadata
- [x] **Task 16c:** Specify precedence rules implementation details → DR-019: Local > global, name/alias resolution priority

*Security & Trust:*
- [ ] **Task 17a:** Define trust model for downloaded assets
- [ ] **Task 17b:** Decide on signature verification (if any)
- [ ] **Task 17c:** Design commit/tag pinning strategy

*Remaining High-Level Design:*
- [x] **Task 18:** Evaluate `start context` command necessity → Resolved: `start config context` created (Task 12f)
- [x] **Task 19:** Evaluate `start role` command necessity → Resolved: `start config role` created (Task 12h)
- [ ] **Task 20:** Determine shell completion requirements (bash/zsh/fish)
- [ ] **Task 21:** Define non-interactive mode flags for CI/automation
