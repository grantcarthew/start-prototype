# DR-010: Default Task Definitions

**Date:** 2025-01-03 (original), 2025-01-07 (updated)
**Status:** Accepted
**Category:** Tasks

## Decision

Ship four interactive review tasks as defaults, each using role-based configuration.

## Default Tasks

1. **code-review** (alias: `cr`) - General code quality review
2. **git-diff-review** (alias: `gdr`) - Review git diff output
3. **comment-tidy** (alias: `ct`) - Review and tidy code comments
4. **doc-review** (alias: `dr`) - Review and improve documentation

## Task Definitions

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

**Usage:**
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
{command}
```

Provide specific feedback on:
- Quality of changes
- Potential issues introduced
- Missing edge cases
- Testing considerations
"""
`````

**Usage:**
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

**Usage:**
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

{file}

Specific improvements needed: {instructions}

Focus on:
- Clarity and readability
- Completeness of information
- Accuracy of examples
- Proper formatting
- User-focused explanations
"""
```

**Usage:**
```bash
start task doc-review "add installation instructions"
start task dr "improve examples"
```

## Role Requirements

These tasks reference two roles that should be defined in the default asset library:

**code-reviewer:**
```toml
[roles.code-reviewer]
description = "Expert code reviewer focused on quality and best practices"
file = "~/.config/start/assets/roles/code-reviewer.md"
```

**documentation-writer:**
```toml
[roles.documentation-writer]
description = "Technical documentation specialist"
file = "~/.config/start/assets/roles/documentation-writer.md"
```

These role files are part of the asset library (see [DR-011](./dr-011-asset-distribution.md)).

## Rationale

**Why these four tasks:**
- All are **interactive reviews** - user works with agent in chat
- Common development workflows
- No file writing or orchestration required
- Useful out-of-the-box for most projects

**Why interactive only:**
- Stays true to vision: launcher only, not a workflow orchestrator
- No tasks that write files or require post-processing
- Agent provides review, user decides on actions
- Clear separation of concerns

**Why role references (not inline):**
- Roles can be updated independently via asset updates
- Users can customize role behavior globally
- Consistent with role-based design (DR-005)
- Enables role reuse across tasks

**Why these specific roles:**
- **code-reviewer**: General purpose for code analysis
- **documentation-writer**: Specialized for doc tasks
- Both available in default asset library
- Users can override with their own roles

## Tasks NOT Included

**commit-message** - Requires committing after generation (orchestration)
- Would need to either:
  - Write commit and execute `git commit` (file I/O + git orchestration)
  - Copy to clipboard and expect user to commit (complex, platform-dependent)
- Violates "launcher only" principle

**gitignore** - Requires saving to file (file I/O)
- Would need to write to `.gitignore` file
- Violates "launcher only" principle
- User can easily copy/paste from interactive review

**update-changelog** - Requires writing to CHANGELOG.md (file I/O)
- Same issues as gitignore
- Better as interactive review where user controls changes

**Note:** Users can add these as custom tasks if they want orchestration in their workflow.

## User Customization

Users can customize default tasks in several ways:

**Override task completely:**
```toml
# User's config
[tasks.code-review]
role = "go-expert"  # Use different role
command = "git diff HEAD~1"  # Add command
prompt = "Review last commit: {command}"
```

**Add new tasks:**
```toml
[tasks.security-audit]
alias = "sec"
role = "security-auditor"
description = "Security-focused review"
command = "git diff --staged"
prompt = "Security audit: {command}"
```

**Override role behavior:**
```toml
# User's config - customize the code-reviewer role
[roles.code-reviewer]
file = "~/my-custom-reviewer.md"
```

**Remove defaults:**
- Default tasks are in assets (not embedded in binary)
- User can simply not use them
- Or override with minimal implementation

## Implementation

**Task location:**
- Default tasks are in the asset library: `~/.config/start/assets/tasks/`
- Not embedded in binary
- Downloaded during `start init` and `start assets update`

**Loading order:**
1. Load global config tasks
2. Load local config tasks
3. Load asset tasks (as templates, not runtime tasks per DR-019)
4. Merge: User config overrides assets

**See [DR-019](./dr-019-task-loading.md) for task loading algorithm.**

## Examples

### Using Default Tasks

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

### Overriding Default Role

```bash
# Use different role for code review
start task code-review --role security-auditor

# Use different role for doc review
start task doc-review --role technical-writer
```

### Overriding Default Agent

```bash
# Use Gemini for code review
start task code-review --agent gemini

# Use different model
start task gdr --model haiku "quick check"
```

## Breaking Changes from Original

This updates the original DR-010 with:

1. **Updated:** All task examples use `role` field (not inline system prompts)
2. **Added:** Role requirement section (code-reviewer, documentation-writer)
3. **Clarified:** Tasks are in asset library (not embedded in binary)
4. **Updated:** Loading order to reflect DR-019 (assets as templates)
5. **Added:** Complete task configurations with examples
6. **Updated:** All examples to use role-based design

## Related Decisions

- [DR-005](./dr-005-role-configuration.md) - Role configuration
- [DR-009](./dr-009-task-structure.md) - Task structure
- [DR-011](./dr-011-asset-distribution.md) - Asset distribution
- [DR-019](./dr-019-task-loading.md) - Task loading (assets are templates, not runtime loaded)
