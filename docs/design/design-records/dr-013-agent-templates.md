# DR-013: Agent Configuration Distribution via GitHub

**Date:** 2025-01-04
**Status:** Accepted
**Category:** Distribution

## Decision

Fetch agent configurations from GitHub during `start init` rather than embedding in binary

## Agent TOML Structure

Each agent configuration file includes a `bin` field for auto-detection:

```toml
[agents.claude]
bin = "claude"  # Binary name for command -v detection (required)
command = "{bin} --model {model} --prompt {prompt}"  # Must contain {bin} and {model}
description = "Claude Code by Anthropic"

[agents.claude.models]
sonnet = "claude-sonnet-4-5"
opus = "claude-opus-4"
```

**Required fields:**
- `bin` - Binary name for auto-detection via `command -v`
- `command` - Command template, must contain `{bin}` and `{model}` placeholders
- `description` - Human-readable description
- `models` - Model name mappings

**Required placeholders in command:**
- `{bin}` - References the bin field, enforces consistency
- `{model}` - Model selection placeholder

**Recommended placeholders:**
- `{prompt}` - Composed prompt (warns if missing)

## Init Behavior

1. Fetch `assets/index.csv` from GitHub (contains agent list with `bin` column)
2. Auto-detect installed agents using `command -v <bin>` from index
3. Download TOML files only for detected/selected agents (lazy loading)
4. Merge into user's `~/.config/start/agents.toml`

## Technical Details

- API endpoint: `https://api.github.com/repos/grantcarthew/start/contents/assets/agents`
- Timeout: 10 seconds
- No caching between runs
- Network required (error if offline)
- Rate limit: 60 requests/hour (unauthenticated)

## Rationale

**Why fetch instead of embed:**

- Model names change frequently (claude-3-5 → claude-3-7 → claude-4)
- Agent command flags evolve over time
- New agents emerge regularly
- Embedding means stale configs until next release
- Users get current configs without waiting for release

**Trade-offs accepted:**

- Requires network during init (acceptable for one-time setup)
- Dependency on GitHub availability (agents need network anyway)
- No offline init (manual config documented as alternative)

## Update Workflow

- Model names stale? Update TOML file in repo
- New agent released? Add new config file
- Flag changes? Update command template
- No code changes or releases needed

## Community Benefits

- Easy to contribute new agent configs (PR a TOML file)
- Clear separation: code vs configuration data
- Living documentation (configs show current best practices)

## Related Decisions

- [DR-011](./dr-011-asset-distribution.md) - Asset distribution strategy
- [DR-014](./dr-014-github-tree-api.md) - GitHub download mechanism
