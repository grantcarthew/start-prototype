# DR-009: Task Structure and Placeholders

**Date:** 2025-01-03
**Status:** Accepted
**Category:** Tasks

## Decision

Tasks include prompt templates with {instructions} and {content} placeholders

## Complete Task Configuration

```toml
[task.git-diff-review]
alias = "gdr"
description = "Review git diff changes"
role = "./roles/code-reviewer.md"
documents = ["environment", "agents"]
command = "git diff --staged"
prompt = """
Analyze the following git diff and act as a code reviewer.

## Special Instructions

{instructions}

## Git Diff

```diff
{content}
```

"""
```

## Task Fields

- **alias** (optional) - Short name for quick access
- **description** (optional) - Help text
- **role** (required) - System prompt (file path or inline text)
- **documents** (optional) - Array of named context documents to include
- **command** (optional) - Shell command to run, output becomes {content}
- **prompt** (required) - Prompt template (file path or inline text)

## Task-Specific Placeholders

- `{instructions}` - Command-line arguments after task name
  - Value: User's arguments or "None" if not provided
  - Usage: `start task gdr "focus on security"` → `{instructions}` = "focus on security"
  - Usage: `start task gdr` → `{instructions}` = "None"

- `{content}` - Output from `command`
  - Value: Command output or empty string if no command
  - Example: `git diff --staged` output

## All Placeholders Available in Task Prompts

- Task-specific: `{instructions}`, `{content}`
- Global: `{model}`, `{system_prompt}`, `{prompt}`, `{date}`

## Usage Examples

```bash
# Basic task
start task git-diff-review

# Task with instructions
start task git-diff-review "only focus on security issues"
start task gdr "ignore comment changes"

# Task with agent override
start task gdr --agent gemini "check for performance issues"
```

## Rationale

- Mirrors existing bash script pattern (gdr, ucl, etc.)
- Flexible: simple tasks don't need command or instructions
- Clear placeholder names ({instructions}, {content})
- "None" default matches existing bash script behavior
- Supports both dynamic content (git diff) and static prompts

## Related Decisions

- [DR-007](./dr-007-placeholders.md) - Global placeholders
- [DR-010](./dr-010-default-tasks.md) - Default task definitions
- [DR-019](./dr-019-task-loading.md) - Task loading algorithm
