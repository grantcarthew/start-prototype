# DR-008: Context File Detection and Handling

- Date: 2025-01-03
- Status: Accepted
- Category: Runtime Behavior

## Problem

The tool needs to handle file paths in configurations (roles, contexts, tasks) that may or may not exist at runtime. The system must:

- Resolve relative paths consistently (relative to what?)
- Handle missing files gracefully (error? warn? ignore?)
- Provide visibility into what files are being used
- Support optional files (not all projects have all documents)
- Help users diagnose path configuration errors
- Not break execution when optional files are missing

## Decision

Relative paths resolve to working directory; missing files generate warnings and are skipped.

Path resolution:

- Relative paths (`./AGENTS.md` or `AGENTS.md`) resolve to working directory
- Working directory defaults to current directory (pwd)
- Override with `--directory` flag
- Absolute paths and tilde paths resolve independently of working directory

Missing file behavior:

- Files that don't exist generate warnings
- Warnings displayed with ⚠ symbol
- No error, no exit - execution continues
- Missing files excluded from prompt
- Status displayed to user before execution

## Why

Warnings instead of errors:

- Optional files work naturally (not all projects have all documents)
- Execution continues even if some contexts missing
- Users can start with minimal config and add files incrementally
- Follows standard CLI conventions (git, npm, make warn about missing configured resources)

Warnings instead of silent skip:

- Users see exactly what context is being used
- Can diagnose path issues easily via warnings
- Warnings help catch configuration errors (typos, wrong paths)
- Visible feedback prevents confusion about why context isn't included

Working directory for relative paths:

- Intuitive - files relative to where you run the command
- `./AGENTS.md` means "in current directory" to users
- Matches shell behavior and user expectations
- Works with `--directory` flag to override when needed

Status display before execution:

- User confirmation of what's being sent to agent
- Clear visibility into which contexts loaded successfully
- Easy to spot misconfigured paths
- Command display shows exactly what's executed

## Trade-offs

Accept:

- Warnings add noise to output (every missing file shows a warning)
- Users must check warnings to spot configuration errors
- Relative paths depend on working directory (could be confusing if running from wrong location)
- Path resolution rules must be understood (relative vs absolute vs tilde)

Gain:

- Optional files work without special configuration
- Easy to diagnose path issues (warnings show exactly what's missing)
- Execution doesn't fail due to missing optional files
- Users can see exactly what context is being used
- Standard CLI behavior (warnings, not errors)
- Flexible - works with partial configs

## Path Resolution

Relative paths (`./file.md` or `file.md`):

- Resolve to working directory
- Working directory defaults to pwd
- Override with `--directory` flag

Path equivalence:

```toml
file = "./AGENTS.md"              # Relative to working directory
file = "AGENTS.md"                # Relative to working directory (same as above)
file = "/absolute/path/file.md"   # Absolute path
file = "~/reference/file.md"      # Tilde expands to home directory
```

Working directory examples:

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

## Missing File Behavior

When a file doesn't exist:

1. Generate warning with ⚠ symbol
2. Display file path that was attempted
3. Mark as "not found, skipped"
4. Continue execution (no error, no exit)
5. Exclude file from prompt composition

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

Output details:

- ✓ for files found and loaded
- ⚠ for files not found (warning)
- Full path shown for clarity
- Command display truncates system prompt ('...') to avoid noise
- Full prompt visible in agent chat once started

## Alternatives

Error on missing files:

- Pro: Forces users to fix path issues immediately
- Pro: Prevents silent failures if file is actually required
- Con: Breaks execution for optional files
- Con: Users must remove optional file references or create empty files
- Con: Prevents incremental config (can't add file references before files exist)
- Con: Too strict for flexible workflow
- Rejected: Too rigid, prevents optional files pattern

Silent skip (no warnings):

- Pro: Clean output with no warnings
- Pro: Works seamlessly with optional files
- Con: Users don't know why context isn't included
- Con: Typos in paths go unnoticed
- Con: Hard to diagnose configuration issues
- Con: Surprising behavior (file configured but not used, no indication why)
- Rejected: Lack of visibility creates confusion

Strict mode flag (--strict-files):

- Pro: Users can choose behavior (warn or error)
- Pro: Flexible for different workflows
- Con: Another flag to learn
- Con: Another mode to document and test
- Con: Most users would need to figure out which mode to use
- Con: Adds complexity for minimal benefit
- Rejected: Warnings are the right default, adding flag is over-engineering

Environment variable for working directory:

- Pro: Could set working directory via WORKDIR env var
- Pro: Consistent across multiple invocations
- Con: Less discoverable than --directory flag
- Con: Environment variables less visible than CLI flags
- Con: Adds another configuration method
- Rejected: --directory flag is clearer and more standard

Relative to config file location:

- Pro: Files relative to config location (like many tools)
- Pro: Config bundles nicely with related files
- Con: Confusing with global vs local configs (which config's location?)
- Con: Breaks with merged configs (files in different locations)
- Con: Doesn't match shell behavior (users expect relative to pwd)
- Rejected: Too confusing with multi-file, multi-scope config structure
