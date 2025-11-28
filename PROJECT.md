# Current Work

This document tracks current development status and planned work.

## Completed (2025-11-28)

### Bug Fixes
- ✅ Fixed `start config edit` command not working
  - Created: `internal/cli/config_edit.go`
  - Added subcommand to config command structure
  - Supports `--local` flag for editing local vs global config
  - Validates configuration after editing with helpful errors

- ✅ Fixed no-prompt interactive session
  - Modified: `internal/cli/root.go`
  - Changed usage from `start [prompt]` to `start [flags]`
  - Removed error requiring prompt argument
  - Allows `start` to run without arguments for interactive sessions

- ✅ Fixed agent command template quoting issues
  - Modified: `~/.config/start/agents.toml`
  - Changed from backticks to single quotes in command templates
  - Added proper quoting around `{role}` and `{prompt}` placeholders

### Enhancements
- ✅ Improved agent configuration creation wizard
  - Modified: `internal/cli/config_agent_new.go`
  - Added examples showing correct quoting patterns
  - Added warning about bash safety and quote usage

### Design Decisions
- ✅ Created DR-044: Shell Quote Escaping for Placeholder Substitution
  - Documented context-aware escaping strategy
  - POSIX shells get auto-escaping, programming languages don't
  - Detailed error messages for quote conflicts
  - Preserves user power for intentional `$()` usage

## In Progress

### Shell Quote Escaping Implementation
Implementation of DR-044 design:
- [ ] Create shell classifier: `IsPOSIXShell(shell string) bool`
- [ ] Implement quote parser state machine for bash/POSIX shells
- [ ] Add escaping functions: `EscapeSingleQuote`, `EscapeDoubleQuote`, `ValidateUnquoted`
- [ ] Create detailed error message generator
- [ ] Add shell reporting to execution output
- [ ] Write tests for all quoting scenarios

## Known Issues

None currently tracked.
