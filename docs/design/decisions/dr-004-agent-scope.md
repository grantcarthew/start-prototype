# DR-004: Agent Configuration Scope

**Date:** 2025-01-03, Updated 2025-01-05
**Status:** Accepted
**Category:** Configuration

## Decision

Agents can be defined in both global and local configs with merge behavior

## Structure

```toml
[settings]
default_agent = "claude"

[agents.claude]
command = "claude --model {model} --append-system-prompt '{system_prompt}' '{prompt}'"
default_model = "sonnet"  # Default alias to use

  [agents.claude.models]
  haiku = "claude-3-5-haiku-20241022"
  sonnet = "claude-3-7-sonnet-20250219"
  opus = "claude-opus-4-20250514"
```

## Model Alias Behavior

- Alias names are user-defined (not hardcoded tiers)
- Each agent has its own set of model aliases
- `default_model` specifies which alias to use when `--model` flag not provided
- Users can use `--model <alias>` or `--model <full-model-name>`

## Scope

**Global agents:** `~/.config/start/config.toml`
- Personal agent configurations
- Individual preferences (model aliases, default models)
- Private configurations

**Local agents:** `./.start/config.toml`
- Team-standardized configurations (can be committed to git)
- Project-specific agent wrappers or custom tools
- Consistent team experience (clone and go)

## Merge Behavior

- Global + local agents are combined
- Same agent name: **local overrides global**
- Enables team standardization while allowing personal overrides
- Local config in version control ensures consistent team setup

## Security Note

Don't commit secrets in local agent configs. Use environment variable references:

```toml
# Bad
[agents.custom.env]
API_KEY = "sk-1234567890"  # DON'T COMMIT

# Good
[agents.custom.env]
API_KEY = "${CUSTOM_API_KEY}"  # Reference user's env var
```

## Updates

- **2025-01-04:** Changed from hardcoded tier names to flexible user-defined aliases
- **2025-01-05:** Changed from global-only to both global + local support

## Rationale

- Agent names are actual tool names (claude, gemini, opencode) not arbitrary aliases
- Self-documenting - clear which agents are available
- Flexible aliases allow users to name models meaningfully for their workflow
- Team standardization via committed `.start/` directory

## Related Decisions

- [DR-002](./dr-002-config-merge.md) - Merge strategy
- [DR-013](./dr-013-agent-templates.md) - Agent templates from GitHub
