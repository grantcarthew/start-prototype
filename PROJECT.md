# Project: start - CLI Design Phase

**Status:** Design Phase - Command-Line Interface Specification
**Date Started:** 2025-01-03
**Current Phase:** Core command specifications

## Overview

Context-aware AI agent launcher that detects project context, builds intelligent prompts, and launches AI development tools with proper configuration.

**Links:** [Vision](./docs/vision.md) | [Design Decisions](./docs/design-record.md) (13 DRs) | [Tasks](./docs/task.md)

## Command Status

### Core Commands
- ✅ `start` - [docs/cli/start.md](./docs/cli/start.md) - Launch with all context
- ✅ `start prompt [text]` - [docs/cli/start-prompt.md](./docs/cli/start-prompt.md) - Launch with required context + optional custom prompt
- ✅ `start init` - [docs/cli/start-init.md](./docs/cli/start-init.md) - Interactive wizard, GitHub fetch, agent detection
- ✅ `start task <name> [instructions]` - [docs/cli/start-task.md](./docs/cli/start-task.md) - Predefined workflows with roles and content commands

### Management Commands
- 🚧 `start agent add|list|test|remove|edit` - Agent configuration management
- 🚧 `start config show|edit|path|validate` - Config file management

### Optional Commands (TBD)
- ❓ `start context` - Context document management (needed or just edit config?)
- ❓ `start role` - Role template management (needed or just file operations?)

## Architecture Decisions

### Completed
- Configuration: TOML with global + local merge
- Context documents: Named with `path`, `prompt`, `required` fields
- Document order: Config definition order (TOML preserves order)
- Agents: Global only, flexible model aliases
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
2. **Agent testing:** What does `start agent test` actually validate?
3. **Config editing:** Validation behavior on save - error or warn?

### Medium Priority
4. **JSON output:** Which commands should support `--json` flag?
5. **Context management:** Build `start context` commands or skip?
6. **Role management:** Build `start role` commands or skip?

### Low Priority
7. **Shell completion:** Generate for bash/zsh/fish?
8. **Non-interactive mode:** What flags needed for CI/automation?

## Success Criteria

CLI design is complete when:
- [ ] All core commands fully specified (start, prompt, init, task)
- [ ] All management commands specified (agent, config)
- [ ] All high-priority questions resolved
- [ ] Error cases documented across commands
- [ ] Output formats specified consistently
- [ ] Patterns consistent across all commands

## Reference

**Doc Template:** See [docs/cli/start.md](./docs/cli/start.md) for complete documentation structure

**Key Documents:**
- [docs/vision.md](./docs/vision.md) - Product vision and goals
- [docs/design-record.md](./docs/design-record.md) - All design decisions (DR-001 through DR-012)
- [docs/task.md](./docs/task.md) - Task configuration details
- [docs/archive/](./docs/archive/) - Design discussion history
