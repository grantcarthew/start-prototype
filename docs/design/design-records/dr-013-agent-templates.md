# DR-013: Agent Configuration Distribution via GitHub

- Date: 2025-01-04
- Status: Accepted
- Category: Distribution

## Problem

Agent configurations need to stay current as model names, command flags, and new agents evolve. The system must:

- Support multiple AI CLI tools (Claude, Gemini, GPT, etc.)
- Keep model names current as providers release new versions
- Update agent command templates as flags change
- Add new agents without requiring code changes or releases
- Allow users to get current configs without waiting for tool updates
- Enable community contributions for new agent configs
- Provide agent auto-detection based on installed binaries
- Separate configuration data from application code

## Decision

Distribute agent configurations via GitHub catalog as downloadable assets. Agent TOML files include `bin` field for auto-detection and follow standardized structure.

Agent configurations are part of the asset catalog and follow the same resolution as other assets (local config → global config → cache → GitHub catalog).

## Why

Fetch instead of embed:

- Model names change frequently (claude-3-5 → claude-3-7 → claude-4)
- Agent command flags evolve over time (new features, deprecated options)
- New agents emerge regularly (new AI providers, new tools)
- Embedding means stale configs until next binary release
- Users get current configs from catalog without waiting
- Configuration data separated from code (different update cycles)

Auto-detection via `bin` field:

- Users don't manually configure agents they have installed
- Tool can detect which agents are available on the system
- Simplifies initial setup (automatic discovery)
- `bin` field provides executable name or path to check

Standardized structure:

- Consistent configuration across all agents
- Validates required fields (bin, command, models)
- Enforces placeholder usage in command templates
- Clear pattern for community contributions

Catalog-based distribution:

- Agents available on-demand via catalog
- Updates without binary releases
- Community can contribute new agent configs
- Living documentation of current best practices

## Trade-offs

Accept:

- Requires network for agent discovery and download
- Dependency on GitHub availability
- Users must understand agent TOML structure
- Command templates must use specific placeholder syntax
- No offline agent discovery (manual config required)

Gain:

- Always-current model names and command flags
- New agents available immediately via catalog
- No code changes needed for agent updates
- Easy community contributions (PR a TOML file)
- Clear separation: code vs configuration data
- Agent auto-detection simplifies setup

## Alternatives

Embed agent configs in binary:

- Pro: Always available offline
- Pro: No network dependency
- Con: Stale configs until next release
- Con: Model name changes require binary updates
- Con: New agents require code changes
- Con: Users get outdated configs
- Rejected: Configuration churn too high for binary embedding

Manual agent configuration only:

- Pro: Maximum flexibility for users
- Pro: No auto-detection complexity
- Con: Poor onboarding experience
- Con: Users must research agent commands and flags
- Con: Error-prone (typos, wrong placeholders)
- Con: No sharing of best practices
- Rejected: User experience too poor

Package manager distribution (apt, brew, etc):

- Pro: Integrated with system package management
- Pro: Offline after initial install
- Con: Separate package for each OS
- Con: Slower update cycle (package review process)
- Con: Requires maintaining multiple packages
- Con: Users must add custom package repositories
- Rejected: Too much overhead for simple TOML files

## Structure

Agent TOML file structure:

```toml
[agents.claude]
bin = "claude"  # Binary path or name (required)
command = "{bin} --model {model} '{prompt}'"  # Must contain {bin} and {model}
description = "Claude Code by Anthropic"
default_model = "sonnet"

  [agents.claude.models]
  haiku = "claude-3-5-haiku-20241022"
  sonnet = "claude-3-7-sonnet-20250219"
  opus = "claude-opus-4-20250514"
```

Required fields:

- `bin` - Binary path or name for agent detection, used to check if agent is installed
- `command` - Command template, must contain `{bin}` and `{model}` placeholders
- `description` - Human-readable description
- `models` - Model name mappings (user-friendly names to full identifiers)

Required placeholders in command:

- `{bin}` - References the bin field, enforces DRY (don't repeat binary name)
- `{model}` - Model selection placeholder

Recommended placeholders:

- `{prompt}` - Composed prompt (warns if missing during validation)
- `{role}` or `{role_file}` - System prompt (for agents supporting roles)

Optional fields:

- `default_model` - Default model name when `--model` flag not provided

## Agent Discovery and Installation

Agent configurations are assets in the catalog (same as tasks, roles, contexts).

Resolution follows standard asset resolution (DR-033):

1. Local config (`.start/agents.toml`)
2. Global config (`~/.config/start/agents.toml`)
3. Asset cache (`~/.config/start/assets/agents/`)
4. GitHub catalog (query, download if `asset_download = true`)

Discovery commands:

```bash
# Browse agent catalog
start assets add  # Interactive TUI browser

# Search for agents
start assets search "claude"

# Direct installation
start assets add agents/claude
```

Auto-detection during init:

- `start init` can optionally detect installed agents by checking for executables
- Prompts user to download configs for detected agents
- Not a bulk download - user confirms which agents to add

## Usage Examples

Agent in catalog:

```toml
[agents.claude]
bin = "claude"
command = "{bin} --model {model} --append-system-prompt '{role}' '{prompt}'"
description = "Claude Code by Anthropic"
default_model = "sonnet"

  [agents.claude.models]
  haiku = "claude-3-5-haiku-20241022"
  sonnet = "claude-3-7-sonnet-20250219"
  opus = "claude-opus-4-20250514"
```

User adds agent to config:

```bash
# Interactive browser
start assets add

# Or direct
start assets add agents/claude
```

Result - added to `~/.config/start/agents.toml`:

```toml
[agents.claude]
bin = "claude"
command = "{bin} --model {model} --append-system-prompt '{role}' '{prompt}'"
description = "Claude Code by Anthropic"
default_model = "sonnet"

  [agents.claude.models]
  haiku = "claude-3-5-haiku-20241022"
  sonnet = "claude-3-7-sonnet-20250219"
  opus = "claude-opus-4-20250514"
```

Using the agent:

```bash
# Use default model
start --agent claude

# Use specific model
start --agent claude --model haiku

# Task with agent override
start task code-review --agent claude --model sonnet
```

## Update Workflow

When model names or commands change:

1. Update TOML file in GitHub catalog
2. Users run `start assets update` to refresh cached agents
3. Users manually update their `agents.toml` if desired (or re-download from catalog)

No code changes or binary releases needed.

## Community Contributions

Adding a new agent:

1. Create TOML file following structure (e.g., `agents/gemini.toml`)
2. Add metadata file (e.g., `agents/gemini.meta.toml`)
3. Submit PR to catalog repository
4. After merge, agent available via `start assets add`

Clear separation: code vs configuration data enables community contributions.

## Validation

At configuration load:

- Agent `bin` field must be present
- Agent `command` must contain `{bin}` placeholder (enforces DRY)
- Agent `command` must contain `{model}` placeholder
- Agent `models` section must have at least one model
- `default_model` (if present) must reference a model in `models` section
- Warn if `command` missing `{prompt}` placeholder

At execution time:

- Check if `bin` executable exists (warn if not found)
- Validate model name resolves to full identifier
- Ensure all command placeholders can be resolved

## Updates

- 2025-01-17: Updated to align with catalog-based asset system (DR-031) - agents are catalog assets, not bulk downloaded during init
