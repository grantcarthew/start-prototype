# DR-008: Context File Detection and Handling

**Date:** 2025-01-03
**Status:** Accepted
**Category:** Runtime Behavior

## Decision

Relative paths resolve to working directory; missing files generate warnings and are skipped

## Path Resolution

- Relative paths (e.g., `./AGENTS.md` or `AGENTS.md`) resolve to working directory
- Working directory defaults to current directory (`pwd`)
- Override with `--directory` flag
- Absolute paths and `~` paths resolve independently of working directory

## Path Equivalence

```toml
file = "./AGENTS.md"   # Same
file = "AGENTS.md"     # Same
file = "/absolute/path/file.md"  # Absolute
file = "~/reference/file.md"     # Home-relative
```

## Missing File Behavior

- Files that don't exist generate warnings
- Warnings displayed with ⚠ symbol
- No error, no exit - execution continues
- Missing files excluded from prompt
- Status displayed to user before execution

## Output Format

```
Starting AI Agent
===============================================================================================
Agent: claude (model: claude-sonnet-4-5@20250929)

Context documents:
  ✓ environment     ~/reference/ENVIRONMENT.md
  ✓ index          ~/reference/INDEX.csv
  ⚠ agents         ./AGENTS.md (not found, skipped)
  ⚠ project        ./PROJECT.md (not found, skipped)

System prompt: ./ROLE.md

Executing command...
❯ claude --model claude-sonnet-4-5@20250929 --append-system-prompt '...' '2025-11-03...'
```

## Rationale

- Users see exactly what context is being used
- Can diagnose path issues easily via warnings
- Warnings help catch configuration errors (typos, wrong paths)
- Follows standard CLI conventions (git, npm, make warn about missing configured resources)
- Optional files work naturally (not all projects have all documents)
- Warnings displayed but don't stop execution
- Command display truncates system prompt ('...') to avoid noise
- Full prompt visible in agent chat once started

## Working Directory Examples

```bash
# Default - uses pwd
cd ~/my-project
start  # Looks for ~/my-project/AGENTS.md

# Override working directory
start --directory ~/my-project

# From anywhere
cd ~
start --directory ~/my-project  # Still finds ~/my-project/AGENTS.md
```

## Related Decisions

- [DR-012](./dr-012-context-required.md) - Required vs optional contexts
