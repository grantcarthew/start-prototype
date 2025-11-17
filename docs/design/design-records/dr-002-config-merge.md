# DR-002: Configuration File Structure

- Date: 2025-01-03
- Status: Accepted
- Category: Configuration

## Problem

The tool needs to support both user-wide defaults and project-specific configurations. Users should be able to:

- Define global defaults that apply to all projects
- Override settings per-project without duplicating everything
- Keep different types of config organized (agents vs tasks vs roles)
- Use version control for project configs without exposing personal settings

## Decision

Multi-file configuration structure with global and local scopes, using merge strategy.

Files in each scope:

- `settings.toml` - Tool settings and defaults
- `agents.toml` - AI agent configurations
- `tasks.toml` - Task definitions
- `roles.toml` - Role definitions
- `contexts.toml` - Context configurations

Locations:

- Global: `~/.config/start/` (user-wide defaults)
- Local: `./.start/` (project-specific overrides)

## Why

Multi-file structure:

- Separates concerns (agents separate from tasks, roles, contexts)
- Easier to manage and edit (smaller files, focused content)
- Better for version control (commit tasks.toml without exposing agents.toml with API keys)
- CLI commands can target specific config types

Merge behavior:

- Local configs merge with global configs
- Same keys in local override global values
- New keys in local are added
- Omitted keys use global defaults
- Allows both defaults and project-specific overrides

Path choices:

- Global uses `~/.config/start/` following XDG Base Directory specification
- Local uses `./.start/` following project-level tool convention (like `.vscode/`, `.github/`, `.docker/`, `.git/`)
- Not `./.config/start/` - the `.config/` pattern is for user-level configs, not project-level
- Follows established pattern where tools use `.<toolname>/` at project root

## Trade-offs

Accept:

- More files to manage (5 files per scope instead of 1)
- Slightly more complex structure for users to understand initially
- CLI must load and merge 5 files instead of 1

Gain:

- Separation of concerns (agents/tasks/roles/contexts/settings isolated)
- Selective version control (commit project tasks without personal agent configs)
- Easier to edit (smaller, focused files)
- Better organization (clear purpose for each file)
- CLI can validate each file type separately

## Alternatives

Single configuration file:

- Pro: Simpler mental model (one file to find)
- Pro: Easier to load (single file read)
- Con: Everything mixed together (settings, agents, tasks, roles, contexts)
- Con: Harder to version control selectively (can't exclude personal agent config)
- Con: Large file becomes unwieldy as config grows
- Rejected: Lack of separation makes version control and organization difficult

Environment variables only:

- Pro: No files to manage
- Pro: Works well in containerized environments
- Con: Cannot handle complex nested structures (agent models, task configurations)
- Con: No comments or documentation inline
- Con: Difficult to manage many settings
- Rejected: Too limited for complex configuration needs

Multiple single-file scopes:

- Pro: Simple structure (one config.toml per scope)
- Con: Still mixes all config types in one file
- Con: Same version control issues as single file approach
- Rejected: Doesn't solve the separation of concerns problem
