# DR-006: CLI Command Structure

- Date: 2025-01-03
- Status: Accepted
- Category: CLI Design

## Problem

The CLI needs a command structure that:

- Supports multiple commands (task execution, config management, asset management, diagnostics)
- Works with both simple and complex workflows
- Follows familiar patterns from other CLI tools
- Allows global flags (--agent, --role, --model) to work across all commands
- Enables dynamic command generation (tasks from config become subcommands)
- Provides good help and documentation support

## Decision

Use Cobra library with subcommand pattern and global flags.

Command pattern:

```bash
start <subcommand> [args] [flags]
```

Global flags available on all commands:

```bash
--agent <name>, -a <name>       # Which agent to use
--role <name>, -r <name>        # Which role to use
--model <name>, -m <name>       # Model name or full identifier
--directory <path>, -d <path>   # Working directory
--asset-download[=bool]         # Enable/disable asset download (default true)
--local, -l                     # Install assets to local config
--quiet, -q                     # Quiet mode
--verbose                       # Verbose output
--debug                         # Debug mode
--help, -h                      # Show help (automatic)
--version, -v                   # Show version
```

## Why

Cobra library benefits:

- Robust subcommand support (nested commands like `start config agent list`)
- Persistent flags work across all subcommands automatically
- Automatic help generation (--help flag, help command)
- Shell completion support built-in
- Standard pattern in Go ecosystem
- Well-tested and maintained

Familiar pattern:

- Follows kubectl/git/docker patterns (widely known)
- Developers already understand subcommand structure
- Natural fit for multiple command types (task, config, assets, doctor)

Dynamic task commands:

- Tasks defined in config become `start task <name>` commands
- No hardcoded task list in CLI code
- User-defined tasks work like built-in commands
- Alias support built into task loading

Extensibility:

- Easy to add new command groups (assets, config, doctor)
- Subcommands can nest (config agent, config task, config role)
- Consistent pattern scales as CLI grows

## Trade-offs

Accept:

- Dependency on Cobra library (external dependency)
- Slightly more complex initial setup than simple flag parsing
- Subcommand pattern adds one level of nesting (`start task <name>` vs `start <name>`)

Gain:

- Professional CLI structure following industry standards
- Automatic help generation and documentation
- Shell completion support out of box
- Persistent flags work everywhere
- Easy to extend with new commands
- Dynamic task generation from config
- Familiar pattern for users

## Command Structure

Root command (interactive session):

```bash
start                              # Launch default session
start --agent gemini               # Launch with specific agent
start --role code-reviewer         # Launch with specific role
```

Task execution:

```bash
start task <name>                  # Run predefined task
start task code-review             # By name
start task cr                      # By alias
start task cr --agent gemini       # With specific agent override
```

Configuration management:

```bash
start config show                  # Display current config
start config validate              # Validate config files
start config agent list            # List configured agents
start config agent test <name>     # Test agent configuration
start config role list             # List configured roles
start config task list             # List configured tasks
```

Asset management (updated in DR-017, DR-041):

```bash
start assets add                   # Add new asset (interactive)
start assets search <query>        # Search catalog
start assets update                # Update cached assets
start assets clean                 # Clean asset cache
```

Diagnostics:

```bash
start doctor                       # System diagnostics
start init                         # Initialize configuration
```

## Automatic Help Support

Cobra automatically adds help support to all commands:

Help flags (automatic):
```bash
-h, --help            # Show help for any command
```

Help command (automatic):
```bash
start help                    # Show root command help
start help config             # Show help for config subcommand
start help config agent       # Show help for nested subcommand
```

Works at all levels:
```bash
start --help
start config --help
start config agent --help
start config agent list --help
start task --help
start task <name> --help
```

## Alternatives

Simple flag-based CLI (no subcommands):

```bash
start --task code-review
start --config-show
start --agent-list
```

- Pro: Simpler structure, no subcommands
- Pro: Fewer concepts for users to learn
- Con: Flags become unwieldy as features grow (--config-agent-test)
- Con: No natural command grouping
- Con: Harder to provide contextual help
- Con: No standard pattern for dynamic task commands
- Rejected: Does not scale well with multiple features

Custom CLI parser (no Cobra):

- Pro: No external dependencies
- Pro: Complete control over parsing
- Con: Must implement subcommands manually
- Con: Must implement help generation manually
- Con: Must implement shell completion manually
- Con: Reinventing well-solved problem
- Con: More code to maintain and test
- Rejected: Significant development overhead for no real benefit

Click-style (Python approach with decorators):

- Pro: Very clean command definition
- Con: Go doesn't have decorators
- Con: Would require reflection-heavy approach
- Con: Less type-safe than Cobra
- Con: No established Go library for this pattern
- Rejected: Doesn't fit Go ecosystem

## Updates

- 2025-01-08: Removed task implementation details (moved to implementation phase)
- 2025-01-08: Added comprehensive command structure examples
