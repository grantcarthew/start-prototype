# DR-004: Agent Configuration Scope

- Date: 2025-01-03
- Status: Accepted
- Category: Configuration

## Problem

Agent configurations need to support different use cases:

- Personal preferences (individual model names, defaults)
- Team standardization (shared agent configs in version control)
- Project-specific agent wrappers or custom tools
- Allow users to override team settings with personal preferences

The configuration must balance individual flexibility with team consistency.

## Decision

Agents can be defined in both global and local configs with merge behavior.

Locations:

- Global agents: `~/.config/start/agents.toml` (personal configurations)
- Local agents: `./.start/agents.toml` (project/team configurations)

Merge behavior:

- Global + local agents are combined
- Same agent name in local overrides global definition completely
- Different agent names are additive

## Why

Global scope for personal preferences:

- Users define their own model names and defaults
- Private configurations not shared with team
- Personal agent wrappers and customizations

Local scope for team standardization:

- Team-standardized configurations committed to git
- Ensures consistent team experience (clone and go)
- Project-specific agent configurations

Merge enables override pattern:

- Teams can define standard agent configs
- Individuals can override with personal preferences
- Best of both worlds: consistency + flexibility

## Trade-offs

Accept:

- Same agent name in both scopes means local completely replaces global (no per-field merge)
- Users must understand precedence (local wins)
- Potential confusion if team config overrides personal config unexpectedly

Gain:

- Team standardization via version control
- Individual flexibility to customize
- Clear, predictable override behavior
- Project-specific agent configurations possible

## Structure

```toml
[agents.claude]
bin = "claude"
command = "{bin} --model {model} '{prompt}'"
default_model = "sonnet"

  [agents.claude.models]
  haiku = "claude-3-5-haiku-20241022"
  sonnet = "claude-3-7-sonnet-20250219"
  opus = "claude-opus-4-20250514"
```

Model name behavior:

- Model names are user-defined (not hardcoded tiers like "fast/standard/advanced")
- Each agent has its own set of model names
- `default_model` specifies which name to use when `--model` flag not provided
- Users can use `--model <name>` or `--model <full-identifier>`

Example merge:

Global `~/.config/start/agents.toml`:
```toml
[agents.claude]
bin = "claude"
command = "{bin} --model {model} '{prompt}'"
default_model = "haiku"  # Personal preference: use cheap model

  [agents.claude.models]
  haiku = "claude-3-5-haiku-20241022"
  sonnet = "claude-3-7-sonnet-20250219"
```

Local `./.start/agents.toml`:
```toml
[agents.claude]
bin = "claude"
command = "{bin} --model {model} '{prompt}'"
default_model = "sonnet"  # Team standard: use better model

  [agents.claude.models]
  sonnet = "claude-3-7-sonnet-20250219"
  opus = "claude-opus-4-20250514"
```

Result: Local definition completely replaces global (default_model = "sonnet", only sonnet/opus models available)

## Alternatives

Global-only agent configuration:

- Pro: Simpler - only one place to look
- Pro: No merge complexity or override confusion
- Con: No team standardization possible
- Con: Every team member must configure agents identically manually
- Con: No project-specific agent configurations
- Rejected: Eliminates team collaboration benefits

Per-field merge instead of complete override:

- Pro: Could keep global models and override just default_model
- Pro: More granular control
- Con: Very complex merge logic (deep merging nested structures)
- Con: Unpredictable behavior when mixing global and local models
- Con: Harder to reason about final configuration
- Rejected: Complexity outweighs benefits, complete override is clearer

Separate team config file (not local):

- Pro: Clear distinction between team and personal configs
- Con: Adds third scope (global/team/local) - too complex
- Con: Unclear precedence order
- Con: More files to manage
- Rejected: Two scopes (global/local) is sufficient

## Updates

- 2025-01-04: Changed from hardcoded tier names to flexible user-defined model names
- 2025-01-05: Changed from global-only to both global + local support
