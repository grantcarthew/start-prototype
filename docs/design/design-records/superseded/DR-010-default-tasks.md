# DR-010: Default Task Definitions

- Date: 2025-01-03 (original), 2025-01-07 (updated), 2025-01-17 (superseded)
- Status: Superseded by DR-031 (Catalog-Based Assets)
- Category: Tasks

## Superseded

This design record is superseded by DR-031 (Catalog-Based Assets). The concept of "default tasks" is obsolete in the catalog-based system.

The catalog system eliminates the need for pre-shipped default tasks because:

- All catalog tasks are equally discoverable via `start assets add` (interactive TUI browser)
- Tasks can be searched via `start assets search "query"`
- Lazy loading: Running `start task <name>` auto-downloads from catalog if found (with `asset_download = true`)
- No distinction between "default" and "other" tasks - all catalog tasks are on-demand

Users can discover and install any task from the catalog interactively or via auto-download on first use. There is no concept of tasks being "shipped" or "default" - everything is catalog-driven and on-demand.

## Problem

The tool needs to provide useful out-of-the-box functionality for common development workflows. The system must:

- Give users immediate value without requiring custom task configuration
- Demonstrate task system capabilities and patterns
- Focus on interactive review workflows (not file writing or orchestration)
- Work across different programming languages and project types
- Stay true to the "launcher only" vision (not a workflow orchestrator)
- Allow user customization and override of defaults
- Support common development activities without requiring complex setup

## Decision

Ship four interactive review tasks as defaults, each using role-based configuration:

1. code-review (alias: cr) - general code quality review
2. git-diff-review (alias: gdr) - review git diff output
3. comment-tidy (alias: ct) - review and tidy code comments
4. doc-review (alias: dr) - review and improve documentation

All default tasks are stored in the asset library (not embedded in binary) and reference roles from the asset library.

## Why

Interactive review focus:

- All tasks launch agent for interactive chat and review
- No file writing or orchestration required
- Agent provides review, user decides on actions
- Stays true to "launcher only" vision
- Clear separation: tool launches, agent advises, user acts

These four specific tasks:

- Cover common development workflows (code review, git diffs, comments, docs)
- Useful immediately without project-specific customization
- Demonstrate task system patterns (with/without commands, different roles)
- Work across any programming language or project type
- Simple enough to understand and customize

Role references instead of inline:

- Roles can be updated independently via asset updates
- Users can customize role behavior globally
- Consistent with role-based design
- Enables role reuse across tasks
- Separates task structure from role content

Asset library distribution:

- Not embedded in binary (can be updated without recompiling)
- Downloaded during `start init` and `start assets update`
- Users can override by defining same task name in their config
- Clear separation between tool code and default content

Tasks NOT included preserve vision:

- No commit-message task (requires git orchestration after generation)
- No gitignore task (requires writing to .gitignore file)
- No update-changelog task (requires file I/O to CHANGELOG.md)
- These violate "launcher only" principle
- Users can add as custom tasks if they want orchestration

## Trade-offs

Accept:

- Limited to four default tasks (not comprehensive coverage of all workflows)
- All tasks are interactive reviews (no automated file writing or orchestration)
- Requires downloading asset library during init
- Tasks depend on role files being present in asset library
- Users must customize if they want non-interactive or file-writing tasks

Gain:

- Immediate value out-of-the-box
- Clear demonstration of task patterns
- Stays true to "launcher only" vision
- Easy to understand and customize
- No file I/O complexity or git orchestration
- Asset-based means updateable without binary changes
- Users control all file changes (tool never writes without user action)

## Alternatives

Embedded tasks in binary:

- Pro: Always available, no asset download needed
- Pro: Guaranteed to be present
- Con: Cannot update without recompiling binary
- Con: Harder to customize (must override via config)
- Con: Mixes task content with tool code
- Rejected: Asset library provides better update and customization path

Include orchestration tasks (commit-message, gitignore, changelog):

- Pro: More comprehensive default task set
- Pro: Covers more common workflows
- Con: Requires file I/O and git orchestration
- Con: Violates "launcher only" principle
- Con: Tool would write files without explicit user action
- Con: More complex error handling and rollback
- Rejected: Interactive review tasks maintain clear boundaries

Larger default task set (10+ tasks):

- Pro: More comprehensive coverage
- Pro: More examples for users
- Con: Overwhelming for new users
- Con: Harder to maintain quality
- Con: More assets to download and store
- Rejected: Four focused tasks provide good balance

No default tasks:

- Pro: Simpler - users define all tasks
- Pro: No opinionated defaults
- Con: No immediate value for new users
- Con: Steeper learning curve
- Con: Users must understand task system before getting value
- Rejected: Default tasks provide important onboarding and examples

Inline role definitions (not role references):

- Pro: Self-contained task definitions
- Pro: No dependency on role files
- Con: Cannot update role behavior independently
- Con: Duplicates role content if multiple tasks use same role
- Con: Inconsistent with role-based design
- Rejected: Role references enable reuse and independent updates

## Task Definitions

Complete task configurations:

### 1. code-review (cr)

```toml
[tasks.code-review]
alias = "cr"
role = "code-reviewer"
description = "General code quality review"
prompt = """
Review the code in this project for quality and best practices.

Focus areas: {instructions}

Provide specific, actionable feedback on:
- Code quality and maintainability
- Potential bugs or issues
- Best practices and improvements
- Security considerations
- Performance implications
"""
```

Usage:

```bash
start task code-review "focus on error handling"
start task cr "check security"
```

### 2. git-diff-review (gdr)

`````toml
[tasks.git-diff-review]
alias = "gdr"
role = "code-reviewer"
description = "Review staged git changes"
command = "git diff --staged"
shell = "bash"
prompt = """
Review the following staged changes:

## Instructions

{instructions}

## Staged Changes

```diff
{command_output}
```

Provide specific feedback on:
- Quality of changes
- Potential issues introduced
- Missing edge cases
- Testing considerations
"""
`````

Usage:

```bash
start task git-diff-review "focus on security"
start task gdr "ignore formatting changes"
```

### 3. comment-tidy (ct)

```toml
[tasks.comment-tidy]
alias = "ct"
role = "code-reviewer"
description = "Review and improve code comments"
prompt = """
Review the code comments in this project.

Special focus: {instructions}

Analyze:
- Comment clarity and accuracy
- Missing comments for complex code
- Outdated or misleading comments
- Documentation completeness
- Comment quality and style

Provide specific suggestions for improvement.
"""
```

Usage:

```bash
start task comment-tidy "check function documentation"
start task ct "verify accuracy"
```

### 4. doc-review (dr)

```toml
[tasks.doc-review]
alias = "dr"
role = "documentation-writer"
description = "Review and improve documentation"
file = "./README.md"
prompt = """
Review this documentation for clarity and completeness:

{file_contents}

Specific improvements needed: {instructions}

Focus on:
- Clarity and readability
- Completeness of information
- Accuracy of examples
- Proper formatting
- User-focused explanations
"""
```

Usage:

```bash
start task doc-review "add installation instructions"
start task dr "improve examples"
```

## Role Requirements

These tasks reference two roles that should be defined in the default asset library:

code-reviewer:

```toml
[roles.code-reviewer]
description = "Expert code reviewer focused on quality and best practices"
file = "~/.config/start/assets/roles/code-reviewer.md"
```

documentation-writer:

```toml
[roles.documentation-writer]
description = "Technical documentation specialist"
file = "~/.config/start/assets/roles/documentation-writer.md"
```

These role files are part of the asset library.

## Structure

Task location:

- Default tasks are in the asset library: `~/.config/start/assets/tasks/`
- Not embedded in binary
- Downloaded during `start init` and `start assets update`

Loading order:

1. Load global config tasks
2. Load local config tasks
3. Load asset tasks (as templates, not runtime tasks)
4. Merge: user config overrides assets

## Usage Examples

Using default tasks:

```bash
# General code review
start task code-review

# Git diff review
start task git-diff-review "focus on security"
start task gdr "check error handling"

# Comment review
start task comment-tidy

# Documentation review
start task doc-review "add examples"
start task dr "improve clarity"
```

Overriding default role:

```bash
# Use different role for code review
start task code-review --role security-auditor

# Use different role for doc review
start task doc-review --role technical-writer
```

Overriding default agent:

```bash
# Use Gemini for code review
start task code-review --agent gemini

# Use different model
start task gdr --model haiku "quick check"
```

User customization - override task completely:

```toml
# User's config
[tasks.code-review]
role = "go-expert"  # Use different role
command = "git diff HEAD~1"  # Add command
prompt = "Review last commit: {command_output}"
```

User customization - add new tasks:

```toml
[tasks.security-audit]
alias = "sec"
role = "security-auditor"
description = "Security-focused review"
command = "git diff --staged"
prompt = "Security audit: {command_output}"
```

User customization - override role behavior:

```toml
# User's config - customize the code-reviewer role
[roles.code-reviewer]
file = "~/my-custom-reviewer.md"
```

Remove defaults:

- Default tasks are in assets (not embedded in binary)
- User can simply not use them
- Or override with minimal implementation

## Breaking Changes from Original

This updates the original DR-010 with:

1. Updated: All task examples use `role` field (not inline system prompts)
2. Added: Role requirement section (code-reviewer, documentation-writer)
3. Clarified: Tasks are in asset library (not embedded in binary)
4. Updated: Loading order to reflect asset-based loading
5. Added: Complete task configurations with examples
6. Updated: All examples to use role-based design

## Updates

- 2025-01-17: Fixed git-diff-review placeholder - changed {command} to {command_output} to show actual diff output (not command string)
- 2025-01-17: Fixed doc-review placeholder - changed {file} to {file_contents} to include file contents in prompt
