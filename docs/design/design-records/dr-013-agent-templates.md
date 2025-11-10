# DR-013: Agent Configuration Distribution via GitHub

**Date:** 2025-01-04
**Status:** Accepted
**Category:** Distribution

## Decision

Fetch agent configurations from GitHub during `start init` rather than embedding in binary

## Init Behavior

1. Fetch agent list from GitHub API: `GET /repos/grantcarthew/start/contents/assets/agents`
2. Auto-detect installed agents using `command -v`
3. Download config files for selected agents
4. Merge into user's `~/.config/start/config.toml`

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
