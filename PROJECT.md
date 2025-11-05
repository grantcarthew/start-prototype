# Project: start - CLI Design Phase

**Status:** Design Phase - Command-Line Interface Specification
**Date Started:** 2025-01-03
**Current Phase:** Core command specifications

## Overview

Context-aware AI agent launcher that detects project context, builds intelligent prompts, and launches AI development tools with proper configuration.

**Links:** [Vision](./docs/vision.md) | [Config Reference](./docs/config.md) | [Design Decisions](./docs/design-record.md) (13 DRs) | [Tasks](./docs/task.md) | [Agent Docs](./docs/cli/start-agent.md)

## Command Status

### Core Commands

- ✅ `start` - [docs/cli/start.md](./docs/cli/start.md) - Launch with all context
- ✅ `start prompt [text]` - [docs/cli/start-prompt.md](./docs/cli/start-prompt.md) - Launch with required context + optional custom prompt
- ✅ `start init` - [docs/cli/start-init.md](./docs/cli/start-init.md) - Interactive wizard, GitHub fetch, agent detection
- ✅ `start task <name> [instructions]` - [docs/cli/start-task.md](./docs/cli/start-task.md) - Predefined workflows with roles and content commands

### Management Commands

- ✅ `start agent add|list|test|remove|edit|default` - [docs/cli/start-agent.md](./docs/cli/start-agent.md) - Agent configuration management
- ✅ `start config show|edit|path|validate` - [docs/cli/start-config.md](./docs/cli/start-config.md) - Config file management

### Optional Commands (TBD)

- ❓ `start context` - Context document management (needed or just edit config?)
- ❓ `start role` - Role template management (needed or just file operations?)

## Architecture Decisions

### Completed

- Configuration: TOML with global + local merge
- Context documents: Named with `path`, `prompt`, `required` fields
- Document order: Config definition order (TOML preserves order)
- Agents: Global only, flexible model aliases, metadata fields (description, url, models_url)
- Agent management: list, add, test, edit, remove, default subcommands
- Tasks: Role + prompt template + optional content_command
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

4. **JSON output:** Which commands should support `--json` flag?
5. **Context management:** Build `start context` commands or skip?
6. **Role management:** Build `start role` commands or skip?

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
- [docs/design-record.md](./docs/design-record.md) - All design decisions (DR-001 through DR-012)
- [docs/task.md](./docs/task.md) - Task configuration details
- [docs/archive/](./docs/archive/) - Design discussion history
