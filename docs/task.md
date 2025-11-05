# Tasks

Tasks are predefined AI workflows with specific roles and context documents.

## Purpose

Tasks allow you to run common AI-assisted workflows with a single command. Each task defines:

- A specific role/system prompt for the AI
- Which context documents to include
- A prompt template with placeholders
- Optional dynamic content (e.g., git diff output)
- An optional alias for quick access

## Configuration

Tasks are defined in `config.toml`:

````toml
[task.code-review]
alias = "cr"
description = "Review code changes"
role = "./roles/code-reviewer.md"
documents = ["environment", "agents"]
prompt = "Review the code for quality, bugs, and best practices."

[task.commit-message]
alias = "cm"
description = "Generate git commit message"
role = """
You are an expert at writing clear, concise git commit messages.
Follow conventional commits format.
Focus on the 'why' not the 'what'.
"""
documents = ["agents"]
content_command = "git diff --staged"
prompt = """
Generate a commit message for the following changes.

## Special Instructions

{instructions}

## Staged Changes

```diff
{content}
````

"""

[task.git-diff-review]
alias = "gdr"
description = "Review git diff changes"
role = "./roles/code-reviewer.md"
documents = ["environment", "agents"]
content_command = "git diff --staged"
prompt = """
Analyze the following git diff and act as a code reviewer.

## Special Instructions

{instructions}

## Git Diff

```diff
{content}
```

"""

````

## Configuration Fields

### Required Fields

- **Name** - The task identifier (e.g., `code-review`)
- **role** - System prompt for the task (file path or inline text)
- **prompt** - Prompt template (file path or inline text)

### Optional Fields

- **alias** - Short name for quick access (e.g., `cr`)
- **description** - Help text shown in `start task --help`
- **documents** - Array of named context documents to include
- **content_command** - Shell command to run, output becomes `{content}` placeholder

## Role Definition

The `role` field supports two formats:

### File Path

```toml
[task.code-review]
role = "./roles/code-reviewer.md"
````

If the value is a valid file path, the file contents are used as the system prompt.

### Inline Text

```toml
[task.commit-msg]
role = """
You are an expert at writing git commit messages.
Use conventional commits format.
"""
```

Use TOML triple-quote syntax for multi-line inline roles.

## Prompt Definition

The `prompt` field defines the template for the task's prompt. It supports two formats:

### File Path

```toml
[task.git-diff-review]
prompt = "./prompts/gdr-prompt.md"
```

If the value is a valid file path, the file contents are used as the prompt template.

### Inline Text

```toml
[task.code-review]
prompt = """
Review this code for quality and best practices.

{instructions}
"""
```

Use TOML triple-quote syntax for multi-line inline prompts.

## Placeholders

Task prompts support placeholders that get replaced at runtime:

### Task-Specific Placeholders

- **{instructions}** - User's command-line arguments after task name

  - If no arguments: replaced with "None"
  - Usage: `start task gdr "focus on security"` → `{instructions}` = "focus on security"

- **{content}** - Output from `content_command`
  - Empty string if no `content_command` configured
  - Example: Output of `git diff --staged`

### Global Placeholders

All global placeholders are also available:

- **{model}** - Model name
- **{system_prompt}** - System prompt content
- **{date}** - Current timestamp

## Content Command

The optional `content_command` field runs a shell command and makes its output available via the `{content}` placeholder:

````toml
[task.git-diff-review]
content_command = "git diff --staged"
prompt = """
Review this diff:

```diff
{content}
````

"""

````

The command runs in the working directory before the AI agent is launched.

## Document References

The `documents` array references named documents from `context.documents`:

```toml
# Context documents defined globally
[context.documents.environment]
path = "~/reference/ENVIRONMENT.md"
suffix = "Read {file} for environment context."

[context.documents.agents]
path = "./AGENTS.md"
suffix = "Read {file} for repository context."

# Task references them by name
[task.code-review]
role = "./roles/code-reviewer.md"
documents = ["environment", "agents"]  # Include these two
````

## Usage

```bash
# Run by full name
start task code-review

# Run by alias
start task cr

# With instructions
start task git-diff-review "focus on security issues"
start task gdr "ignore comment changes"

# With specific agent
start task code-review --agent gemini

# Task with agent and instructions
start task gdr --agent gemini "check performance"

# List all available tasks
start task --help
```

## Default Tasks

`start` ships with default tasks for common workflows. Users can:

- Override defaults by defining tasks with the same name
- Add custom tasks
- Remove tasks by not including them in config

## Example: Complete Task Configuration

```toml
[task.documentation-review]
alias = "dr"
description = "Review and improve documentation"
role = """
You are a technical documentation expert.
Review documentation for:
- Clarity and accuracy
- Completeness
- Code examples
- Common user questions
"""
documents = ["environment", "agents", "project"]
prompt = """
Review the following documentation.

## Special Instructions

{instructions}

## Documentation Files

Review all markdown files in the current directory for quality and completeness.
"""

[task.update-changelog]
alias = "ucl"
description = "Generate changelog entries from git commits"
role = "./roles/changelog-writer.md"
documents = ["agents"]
content_command = "git log --oneline --no-merges $(git describe --tags --abbrev=0)..HEAD"
prompt = """
Generate changelog entries for the following commits.

## Special Instructions

{instructions}

## Recent Commits

```

{content}

```

Format the output as markdown suitable for CHANGELOG.md.
"""
```
