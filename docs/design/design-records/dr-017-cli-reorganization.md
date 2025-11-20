# DR-017: CLI Command Reorganization

- Date: 2025-01-06
- Status: Accepted
- Category: CLI Design

## Problem

CLI command structure needs clear organization as the tool grows. The system must:

- Support multiple command types (execution, configuration management, asset management, utilities)
- Provide clear mental model for users (what goes where?)
- Avoid top-level command pollution as features are added
- Make commands discoverable (users can find what they need)
- Follow consistent patterns across similar operations
- Separate concerns (execution vs management vs discovery)
- Allow extensibility without restructuring
- Work well with shell completion

Original inconsistent structure:

- `start agent add` - configuration management
- `start task code-review` - execution

Different purposes with similar command structure creates confusion.

## Decision

Reorganize commands by purpose: execution at top level, configuration management under `start config`, asset management under `start assets`.

Execution commands (top-level):

- `start` - launch agent with all context
- `start prompt [text]` - launch with required context + custom prompt
- `start task <name> [inst]` - run predefined task

Configuration management (under `start config`):

- `start config show` - view merged configuration
- `start config edit [scope]` - edit config file
- `start config path` - show config file paths
- `start config validate` - validate configuration
- `start config agent <sub>` - manage agents
- `start config context <sub>` - manage contexts
- `start config task <sub>` - manage tasks
- `start config role <sub>` - manage roles

Asset management (under `start assets`):

- `start assets browse` - open GitHub catalog in browser
- `start assets search <query>` - search by name/description/tags
- `start assets add [query]` - interactive TUI or direct install
- `start assets info <query>` - show asset details
- `start assets update [query]` - update cached assets
- `start assets clean` - remove unused cache

Utility commands (top-level):

- `start init [--local]` - initialize configuration
- `start doctor` - diagnose installation

## Why

Clear separation by purpose:

- `start <thing>` - execute something (run tasks, launch sessions)
- `start config <type>` - manage configuration (agents, tasks, roles, contexts)
- `start assets <action>` - discover and install from catalog
- `start utility` - setup and diagnostics (init, doctor)

Top-level for execution:

- Primary use case is running tasks or launching sessions
- Shortest path for most common operations
- Natural language: "start task X", "start prompt Y"
- No namespace pollution (execution commands are finite)

Grouped configuration management:

- All config operations under one namespace
- Consistent subcommand pattern across types
- Easy discovery: "How do I manage X? → start config X"
- Scales without top-level pollution
- Clear mental model: config = managing settings

Asset management separation:

- Distinct workflow from config management
- Catalog browsing and discovery separate from direct config editing
- Unified interface for all asset types
- Type-agnostic installation

Consistent subcommand patterns:

- All config types support: list, add, new, show, test, edit, remove
- Predictable interface across agents, tasks, roles, contexts
- Easy to learn one, apply to all

## Trade-offs

Accept:

- Longer paths for configuration commands (`start config agent list` vs `start agent list`)
- More nesting (three levels for some commands)
- Initial learning curve for command structure
- Some commands are longer to type
- Breaking change from earlier design iterations

Gain:

- Clear mental model (execution vs configuration vs assets)
- Consistent patterns across all config types
- Scalable (new config types fit naturally)
- Discoverable (logical grouping)
- No top-level pollution (only execution commands at root)
- Better shell completion organization
- Separation of concerns

## Alternatives

Flat structure (everything at top level):

```bash
start task <name>
start agent-list
start agent-add
start task-list
start task-add
start role-list
```

- Pro: Shorter commands, fewer levels
- Pro: Faster to type
- Con: Top-level pollution (dozens of commands)
- Con: Inconsistent naming (dash vs no dash)
- Con: Hard to discover (no logical grouping)
- Con: Poor shell completion (all mixed together)
- Rejected: Doesn't scale, hard to maintain consistency

Type-first grouping:

```bash
start agent list
start agent add
start agent test
start task list
start task add
start task test
```

- Pro: Shorter than config namespace
- Pro: Type-based grouping
- Con: Mixes execution and configuration (`start task run` vs `start task list`)
- Con: Ambiguous purpose (is `start agent` running or configuring?)
- Con: Still pollutes top level as types grow
- Rejected: Confuses execution with configuration

Separate binaries:

```bash
start task <name>        # Execution binary
start-config agent list  # Configuration binary
start-assets search      # Asset binary
```

- Pro: Complete separation
- Pro: Each binary focused
- Con: Multiple binaries to install and maintain
- Con: Harder to discover (users may not know about start-config)
- Con: Inconsistent experience (different binaries)
- Con: Distribution complexity
- Rejected: Over-engineering, poor discoverability

Action-first grouping:

```bash
start list agents
start list tasks
start add agent
start add task
```

- Pro: Action-oriented
- Pro: Consistent verb-first pattern
- Con: Awkward for execution (`start run task` vs `start task`)
- Con: Doesn't group related types together
- Con: Poor for tab completion (all verbs mixed)
- Rejected: Doesn't match natural language or user mental models

## Command Structure

Execution commands (top-level):

```bash
start                        # Launch agent with all context
start prompt [text]          # Launch with required context + custom prompt
start task <name> [inst]     # Run predefined task
```

Configuration management:

```bash
start config show            # View merged configuration
start config edit [scope]    # Edit config file
start config path            # Show config file paths
start config validate        # Validate configuration

start config agent <sub>     # Manage agents
start config context <sub>   # Manage contexts
start config task <sub>      # Manage tasks
start config role <sub>      # Manage system prompts
```

Asset management (per DR-041):

```bash
start assets browse          # Open GitHub catalog in browser
start assets search <query>  # Search by name/description/tags
start assets add [query]     # Interactive TUI or direct install
start assets info <query>    # Show detailed asset information
start assets update [query]  # Update cached assets
start assets clean           # Remove unused cached assets
start assets index           # Generate catalog index (contributors)
```

Utility commands:

```bash
start init [--local]         # Initialize configuration
start doctor                 # Diagnose installation
```

## Configuration Subcommands

All config types follow consistent pattern:

Agent management:

```bash
start config agent list [scope]
start config agent add [name] [scope]  # DEPRECATED - use 'start assets add'
start config agent new [scope]
start config agent show [name] [scope]
start config agent test <name>
start config agent edit [name] [scope]
start config agent remove [name] [scope]
start config agent default [name]
```

Context management:

```bash
start config context list [scope]
start config context add [name] [scope]
start config context new [scope]
start config context show [name] [scope]
start config context test <name>
start config context edit [name] [scope]
start config context remove [name] [scope]
```

Task management:

```bash
start config task list [scope]
start config task add [name] [scope]  # DEPRECATED - use 'start assets add'
start config task new [scope]
start config task show [name] [scope]
start config task test <name>
start config task edit [name] [scope]
start config task remove [name] [scope]
```

Role management:

```bash
start config role list [scope]
start config role add [path]  # DEPRECATED - use 'start assets add'
start config role new [scope]
start config role show [scope]
start config role test
start config role edit [scope]
start config role remove [scope]
start config role default [name]
```

## Scope

Scope parameter applies to configuration commands:

- `[scope]` is optional, defaults to global
- Valid values: `global`, `local`
- Examples:
  - `start config agent list` - shows global agents
  - `start config agent list local` - shows local agents
  - `start config agent add claude global` - adds to global config
  - `start config agent add claude local` - adds to local config

## Usage Examples

Execution (primary use cases):

```bash
# Interactive session
start

# One-off query
start prompt "What's the current date?"

# Run task
start task code-review "focus on security"
start task gdr "check error handling"
```

Configuration management:

```bash
# View configuration
start config show
start config path

# Manage agents
start config agent list
start config agent test claude
start config agent default claude

# Manage tasks
start config task list
start config task show code-review
start config task remove old-task

# Manage roles
start config role list
start config role default code-reviewer
```

Asset management:

```bash
# Browse catalog
start assets browse

# Search and install
start assets search "commit"
start assets add "pre-commit-review"

# Update cached assets
start assets update
start assets update git
```

Utilities:

```bash
# Initialize config
start init
start init --local

# Diagnose issues
start doctor
```

## Breaking Changes

This is a design-phase reorganization (no existing users):

- Previous: `start agent <sub>` → New: `start config agent <sub>`
- Asset commands consolidated under `start assets` (per DR-041)

## Updates

- 2025-01-10: Added `start assets` command suite per DR-041 (unified asset management)
- 2025-01-10: Deprecated `start config [type] add` commands in favor of `start assets add`
