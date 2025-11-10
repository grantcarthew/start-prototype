# DR-017: CLI Command Reorganization

**Date:** 2025-01-06
**Status:** Accepted
**Category:** CLI Design

## Decision

Reorganize commands by purpose - configuration management under `start config`, execution at top level

## Problem

Original structure was inconsistent:
- `start agent add` - Configuration management
- `start task code-review` - Execution

Different purposes, similar command structure = confusing.

## New Structure

**Execution commands (top-level):**
```bash
start                        # Launch agent with all context
start prompt [text]          # Launch with required context + custom prompt
start task <name> [inst]     # Run predefined task
```

**Configuration management:**
```bash
start config show            # View merged configuration
start config edit [scope]    # Edit config file
start config path            # Show config file paths
start config validate        # Validate configuration

start config agent <sub>     # Manage agents (moved from start agent)
start config context <sub>   # Manage contexts (new)
start config task <sub>      # Manage tasks (new)
start config role <sub>      # Manage system prompts (new)
```

**Utility commands:**
```bash
start init [scope]           # Initialize configuration
start doctor                 # Diagnose installation
start update                 # Update asset library
```

## Configuration Subcommands

All follow consistent pattern:

```bash
start config agent list
start config agent add
start config agent edit [name]
start config agent remove [name]
start config agent test <name>
start config agent default [name]

start config context list
start config context add
start config context edit [name]
start config context remove [name]
start config context test <name>

start config task list
start config task add
start config task edit [name]
start config task remove [name]

start config role list
start config role add
start config role edit [name]
start config role remove [name]
```

## Benefits

- ✅ **Clear mental model:** `start config X` = managing config, `start X` = executing
- ✅ **Consistent:** All configuration under one command
- ✅ **Discoverable:** Easy to find: "How do I manage X? → start config X"
- ✅ **Extensible:** New config sections fit the pattern naturally
- ✅ **No ambiguity:** Command purpose clear from structure

## Breaking Change

- Design phase only - no existing users
- `start agent` → `start config agent`

## Rationale

Clear separation between configuration (managing settings) and execution (running the tool) provides better mental model and allows consistent expansion of configuration management without top-level command pollution.

## Related Decisions

- [DR-006](./dr-006-cobra-cli.md) - CLI framework and structure
