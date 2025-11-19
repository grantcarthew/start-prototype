# DR-012: Context Document Required Field and Order

- Date: 2025-01-04
- Status: Accepted
- Category: Configuration

## Problem

Context documents need different inclusion rules for different use cases. The system must:

- Support "essential" context that should always be included (environment info, repository overview)
- Support "optional" context that adds value in full sessions but adds noise in focused queries
- Allow users to control which contexts are included in different scenarios
- Maintain consistent ordering of context documents in prompts
- Work with different command types (interactive sessions, one-off prompts, tasks)
- Give users explicit control over context priority and ordering

## Decision

Add optional `required` field to context documents to control inclusion behavior. Documents appear in config definition order (TOML preserves declaration order).

Context inclusion by command:

- `start` (interactive session) - includes ALL documents (required + optional)
- `start prompt` (one-off query) - includes ONLY required documents
- `start task` (task execution) - includes ONLY required documents

Default value: `required = false` (optional document)

## Why

Required field for inclusion control:

- `start` provides full context for comprehensive interactive sessions
- `start prompt` provides minimal context for focused queries (reduces noise)
- `start task` provides essential context without overwhelming task-specific prompts
- Users designate "essential" vs "nice-to-have" context
- Clear opt-in mechanism for critical context

Definition order for predictability:

- TOML preserves declaration order within sections
- Users control order by arranging config file
- No alphabetical or automatic sorting (explicit, not magic)
- Predictable and consistent across all commands
- First-defined = first in prompt (clear priority)

Default to optional (false):

- Conservative default: don't include unless explicitly required
- Users opt-in to always-included context
- Prevents accidental context bloat
- Explicit is better than implicit

Automatic inclusion in tasks:

- Tasks automatically include required contexts (simplifies task configuration)
- No per-task context management needed
- Ensures critical context always present
- Consistent behavior across all tasks

## Trade-offs

Accept:

- Users must understand required vs optional distinction
- Two-tier context system adds configuration complexity
- Definition order matters (users must arrange config file carefully)
- Cannot exclude required contexts from tasks (no per-task override)

Gain:

- Focused queries with `start prompt` avoid context noise
- Interactive sessions get full context with `start`
- Tasks get essential context without configuration
- Users control context priority via file order
- Clear, predictable behavior across commands
- Essential context never accidentally omitted

## Alternatives

Single inclusion rule (all contexts always included):

- Pro: Simpler - no required field to understand
- Pro: Consistent behavior across all commands
- Con: One-off queries get noisy with unnecessary context
- Con: No way to distinguish essential vs optional
- Con: Task prompts become cluttered
- Rejected: Different commands need different context levels

Per-command context selection:

```toml
[contexts.project]
file = "./PROJECT.md"
prompt = "..."
include_in = ["start", "task"]  # Not in "prompt"
```

- Pro: Fine-grained control per command type
- Pro: Explicit about where each context appears
- Con: More complex configuration
- Con: Three fields to manage per context
- Con: Harder to reason about (must check three boolean-like values)
- Rejected: Required/optional distinction is simpler

Tasks with `documents` array field:

```toml
[tasks.review]
documents = ["environment", "project"]
prompt = "Review code"
```

- Pro: Per-task control over contexts
- Pro: Can cherry-pick contexts
- Con: Must manage context list for every task
- Con: Easy to forget critical contexts
- Con: Inconsistent context inclusion across tasks
- Con: More configuration burden
- Rejected: Automatic required context inclusion is simpler (per DR-009)

Alphabetical ordering:

- Pro: Predictable without users arranging file
- Pro: Consistent across edits
- Con: Users lose control over priority
- Con: "Important" context might come last alphabetically
- Con: Magic behavior (not explicit)
- Rejected: Definition order gives users explicit control

## Structure

Context configuration with required field:

```toml
[contexts.environment]  # First in prompt
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
required = true    # Always included

[contexts.project]      # Second in prompt
file = "./PROJECT.md"
prompt = "Read {file}. Respond with summary."
required = false   # Optional (default)
```

Behavior by command:

- `start` (root) - includes ALL documents (required + optional)
- `start prompt` - includes ONLY required documents
- `start task` - includes ONLY required documents (automatically, per DR-009)

Default value:

If `required` field is omitted, defaults to `false` (optional document).

Document order:

- Documents appear in prompt in the order defined in config file
- TOML preserves declaration order within sections
- Users control order by arranging config file
- Predictable and explicit - no alphabetical or automatic sorting
- Consistent across all commands (start, start prompt, tasks)

## Usage Examples

Essential context (always included):

```toml
[contexts.environment]
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
required = true  # Included in: start, start prompt, start task

[contexts.index]
file = "~/reference/INDEX.csv"
prompt = "Documentation index: {file}"
required = true  # Included in: start, start prompt, start task

[contexts.agents]
file = "./AGENTS.md"
prompt = "Read {file} for repository context."
required = true  # Included in: start, start prompt, start task
```

Optional context (interactive sessions only):

```toml
[contexts.project]
file = "./PROJECT.md"
prompt = "Read {file} for project context."
required = false  # Included in: start only (not start prompt, not start task)
```

Use cases:

- `~/reference/ENVIRONMENT.md` marked required: Always provides user/environment context (first)
- `~/reference/INDEX.csv` marked required: Always provides documentation index (second)
- `AGENTS.md` marked required: Always provides repository overview (third)
- `PROJECT.md` marked optional: Included for full sessions, excluded for quick queries and tasks

Commands with different context levels:

```bash
# Interactive session - gets ALL contexts (required + optional)
start

# One-off query - gets ONLY required contexts
start prompt "What's the current date?"

# Task execution - gets ONLY required contexts
start task code-review "focus on security"
```

## Updates

- 2025-01-17: Fixed task behavior - tasks automatically include required contexts (no `documents` array field per DR-009)
