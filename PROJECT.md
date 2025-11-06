# Project: start - Configuration Design Phase

**Status:** Design Phase - Configuration Structure & Patterns
**Date Started:** 2025-01-03
**Current Phase:** Unified Template Design (UTD) implementation

## Overview

Context-aware AI agent launcher that detects project context, builds intelligent prompts, and launches AI development tools with proper configuration.

**Links:** [Vision](./docs/vision.md) | [Config Reference](./docs/config.md) | [UTD](./docs/design/unified-template-design.md) | [Design Decisions](./docs/design/design-record.md) (13 DRs) | [Tasks](./docs/tasks.md)

## Command Status

### Core Commands

- ✅ `start` - [docs/cli/start.md](./docs/cli/start.md) - Launch with all context
- ✅ `start prompt [text]` - [docs/cli/start-prompt.md](./docs/cli/start-prompt.md) - Launch with required context + optional custom prompt
- ✅ `start init` - [docs/cli/start-init.md](./docs/cli/start-init.md) - Interactive wizard, GitHub fetch, agent detection
- ✅ `start task <name> [instructions]` - [docs/cli/start-task.md](./docs/cli/start-task.md) - Predefined workflows with roles and content commands

### Management Commands

- ✅ `start agent add|list|test|remove|edit|default` - [docs/cli/start-agent.md](./docs/cli/start-agent.md) - Agent configuration management
- ✅ `start config show|edit|path|validate` - [docs/cli/start-config.md](./docs/cli/start-config.md) - Config file management
- ✅ `start doctor` - [docs/cli/start-doctor.md](./docs/cli/start-doctor.md) - Diagnose installation and configuration health
- ✅ `start update` - [docs/cli/start-update.md](./docs/cli/start-update.md) - Update asset library from GitHub

### Optional Commands (TBD)

- ❓ `start context` - Context document management (needed or just edit config?)
- ❓ `start role` - Role template management (needed or just file operations?)

## Architecture Decisions

### Completed

- Configuration: TOML with global + local merge
- Context documents: Named with `file`, `command`, `prompt`, `required` fields using UTD
- Document order: Config definition order (TOML preserves order)
- Agents: Global + local (team standardization), flexible model aliases, metadata fields (description, url, models_url)
- Agent management: list, add, test, edit, remove, default subcommands
- Tasks: Full UTD pattern for system_prompt and task prompt, auto-includes required contexts
- Command pattern: Positional scope arguments (`start init [scope]`, `start config edit [scope]`)
- CLI Framework: Cobra with dynamic task loading

### Key Design

- `start` includes ALL documents (required + optional)
- `start prompt` includes ONLY required documents
- Model aliases user-defined per agent (not hardcoded tiers)
- Missing files show but don't error
- Verbosity: quiet/normal/verbose/debug

## Open Questions

### High Priority

1. ~~**Task listing:** Subcommand (`start task list`) or flag (`start task --list`)?~~ ✅ Resolved: `start task` with no args + `--help`
2. ~~**Agent testing:** What does `start agent test` actually validate?~~ ✅ Resolved: Binary availability (exec.LookPath), config validation, dry-run display
3. ~~**Config editing:** Validation behavior on save - error or warn?~~ ✅ Resolved: Soft warnings (non-blocking, user already saved)

### Medium Priority

4. ~~**JSON output:** Which commands should support `--json` flag?~~ ✅ Resolved: Not needed, human-facing tool
5. ~~**Task structure:** Finalize task config with system_prompt_* fields and UTD~~ ✅ Resolved: Full UTD for both system_prompt and task prompt
6. **Context management:** Build `start context` commands or skip?
7. **Role management:** Build `start role` commands or skip? (roles section not currently used)

### Low Priority

7. **Shell completion:** Generate for bash/zsh/fish?
8. **Non-interactive mode:** What flags needed for CI/automation?

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

### Asset Management System (2025-01-06)

**Asset Library Design:**
- ✅ Assets fetched from GitHub repository (not embedded in binary)
- ✅ Stored in `~/.config/start/assets/` directory
- ✅ Version tracked in `.asset-version` file (commit SHA + timestamp)
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
- [ ] **Task 12a:** Decide GitHub download strategy (API calls, manifest format, file discovery)
- [ ] **Task 12b:** Design atomic update mechanism (temp directory, swap, partial failure recovery)
- [ ] **Task 12c:** Define asset discovery system (how start knows what asset types exist)

*Version Tracking & Checking:*
- [ ] **Task 13a:** Define binary version source (build-time injection strategy)
- [ ] **Task 13b:** Design GitHub version checking (API endpoint, rate limiting, caching)
- [ ] **Task 13c:** Define commit SHA retrieval strategy (releases vs commits)

*Doctor Implementation:*
- [ ] **Task 14a:** Design asset staleness checking (local-only vs GitHub comparison)
- [ ] **Task 14b:** Define exit code priority system (multiple simultaneous issues)
- [ ] **Task 14c:** Design automatic check frequency and caching strategy

*Integration & Offline Support:*
- [ ] **Task 15a:** Define start init + start update relationship (does init call update?)
- [ ] **Task 15b:** Design offline fallback strategy (manual asset installation)
- [ ] **Task 15c:** Define behavior when network unavailable

*Task Merging Implementation:*
- [ ] **Task 16a:** Design task loading and merging algorithm (assets + user config)
- [ ] **Task 16b:** Define source metadata tracking ([default] vs [user] labels)
- [ ] **Task 16c:** Specify precedence rules implementation details

*Security & Trust:*
- [ ] **Task 17a:** Define trust model for downloaded assets
- [ ] **Task 17b:** Decide on signature verification (if any)
- [ ] **Task 17c:** Design commit/tag pinning strategy

*Remaining High-Level Design:*
- [ ] **Task 18:** Evaluate `start context` command necessity (vs direct config editing)
- [ ] **Task 19:** Evaluate `start role` command necessity (roles section not currently used)
- [ ] **Task 20:** Determine shell completion requirements (bash/zsh/fish)
- [ ] **Task 21:** Define non-interactive mode flags for CI/automation
