# DR-006: CLI Command Structure

**Date:** 2025-01-03
**Status:** Accepted
**Category:** CLI Design

## Decision

Use Cobra with subcommand pattern and global flags

## Pattern

```bash
start <subcommand> [args] [flags]
```

## Core Commands

```bash
# Root command (no subcommand)
start                              # Launch default session with context
start --agent gemini               # Launch with specific agent

# Task subcommand
start task <name>                  # Run predefined task
start task code-review             # By name
start task cr                      # By alias
start task cr --agent gemini       # With specific agent

# Agent management (updated in DR-017, DR-041)
start assets add                   # Add new asset (interactive)
start config agent list            # List configured agents
start config agent test <name>     # Test agent configuration

# Config management
start config show                  # Display current config
start config init                  # Create default config
start config edit                  # Open config in editor
```

## Global Flags

Work on all commands:

```bash
--agent <name>, -a <name>       # Which agent to use (overrides default)
--role <name>, -r <name>        # Which role to use (overrides default)
--model <name>, -m <name>       # Model name or full model identifier
--directory <path>, -d <path>   # Working directory (default: pwd)
--quiet, -q                     # Quiet mode (no output)
--verbose                       # Verbose output (detailed info)
--debug                         # Debug mode (everything)
--help, -h                      # Show help (automatic from Cobra)
--version, -v                   # Show version information
```

## Automatic Help Support

Cobra automatically adds help support to all commands with zero configuration required:

**Help flags (automatic):**
```bash
-h, --help            # Show help for any command
```

**Help command (automatic):**
```bash
start help                    # Show root command help
start help config             # Show help for config subcommand
start help config agent       # Show help for nested subcommand
```

**Works at all levels:**
```bash
start --help
start config --help
start config agent --help
start config agent list --help
start task --help
start task code-review --help
```

Every command in the `start` CLI has help support built-in. No custom code required.

## Rationale

- Cobra provides robust subcommand support
- Persistent flags work across all subcommands
- Follows kubectl/git patterns (familiar to developers)
- Tasks discovered dynamically from config
- Extensible for future subcommands

## Task Implementation

- Tasks defined in config are loaded at startup
- Cobra subcommands generated dynamically
- Each task becomes `start task <name>` with alias support

## Related Decisions

- [DR-009](./dr-009-task-structure.md) - Task configuration details
- [DR-017](./dr-017-cli-reorganization.md) - CLI reorganization (start config)
