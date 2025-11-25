# Configuration Examples

This directory contains example configurations demonstrating different use cases for the `start` CLI tool.

## Directory Structure

```
examples/
├── minimal/          Quick start with bare minimum configuration
├── complete/         Complete reference showing all possible fields
└── real-world/       Practical working configuration for everyday use
```

Each example contains both **global** and **local** configuration directories showing the multi-file structure.

## Configuration File Structure

All examples use the same multi-file structure:

**Global configuration** (`~/.config/start/`):
- `config.toml` - Tool settings and defaults
- `agents.toml` - AI agent configurations
- `roles.toml` - Role (system prompt) definitions
- `contexts.toml` - Context document configurations
- `tasks.toml` - Task workflow definitions

**Local configuration** (`./.start/`):
- Same structure as global
- Local configs merge with global configs
- Local completely replaces global for same-named items

## Examples

### Minimal Example (`minimal/`)

**Purpose:** Quick start guide showing the bare minimum to get started.

**Use case:** Learning the basic structure and getting a working setup quickly.

**Features:**
- Single test agent (`smith`)
- Basic role configuration
- One required context
- Simple task example
- **Only global config** - local follows same structure

**How to use:**
```bash
# Copy to your global config
mkdir -p ~/.config/start
cp examples/minimal/global/* ~/.config/start/

# Verify configuration
start doctor

# For project-specific config, create ./.start/ with same structure
# See complete/ or real-world/ examples for local config patterns
```

### Complete Reference (`complete/`)

**Purpose:** Comprehensive documentation of every possible configuration field.

**Use case:** Reference when configuring advanced features or understanding all options.

**Features:**
- ALL possible fields documented with comments
- Multiple UTD pattern examples
- Validation rules explained
- Placeholder usage documented
- Merge behavior demonstrated
- Shell and timeout overrides shown

**How to use:**
- Read for understanding all available options
- Reference when adding advanced features
- Copy specific sections as needed

**Note:** Not intended as a working configuration - use as documentation reference.

### Real-World Example (`real-world/`)

**Purpose:** Practical, ready-to-use configuration for everyday development.

**Use case:** Starting point for actual projects with real AI tools.

**Features:**
- Real agents: Claude, Gemini, aichat
- Practical roles: code-reviewer, go-expert, security-auditor
- Useful contexts: environment, project, architecture docs
- Common tasks: git-diff-review, commit-message, security-scan
- Project-specific overrides in local config

**How to use:**
```bash
# Copy to your global config
cp -r examples/real-world/global/* ~/.config/start/

# For project-specific config (in your project root)
mkdir -p .start
cp -r examples/real-world/local/* ./.start/

# Customize for your needs
# Edit ~/.config/start/agents.toml to match your installed AI tools
# Edit ./.start/roles.toml for project-specific roles

# Verify configuration
start doctor
```

## Key Concepts

### Multi-File Structure

Configuration is split into 5 files per scope (global/local):

1. **config.toml** - Settings (defaults, timeouts, shell)
2. **agents.toml** - AI agent definitions
3. **roles.toml** - System prompt (role) definitions
4. **contexts.toml** - Context documents
5. **tasks.toml** - Reusable workflow tasks

**Why?** Separation of concerns, selective version control, easier management.

### Global vs Local

**Global** (`~/.config/start/`):
- Personal defaults across all projects
- Installed AI agents and preferred models
- Personal roles and contexts
- Reusable tasks

**Local** (`./.start/`):
- Project-specific overrides
- Project conventions and guidelines
- Team-shared configurations (version controlled)
- Project-specific tasks

**Merge behavior:**
- Local + global configs are combined
- Same name in local **completely replaces** global (no field merging)
- Different names are additive

### Unified Template Design (UTD)

Roles, contexts, and tasks use UTD pattern with three optional fields:

- **file** - Path to file (supports `{file}` and `{file_contents}` placeholders)
- **command** - Shell command to execute (supports `{command}` and `{command_output}` placeholders)
- **prompt** - Template text with placeholders

**At least one field required.**

**Examples:**

```toml
# File only
[roles.simple]
file = "~/roles/reviewer.md"

# Inline prompt only
[roles.inline]
prompt = "You are a helpful assistant. Date: {date}"

# File + command + prompt (full UTD)
[contexts.project-state]
file = "./PROJECT.md"
command = "git status --short"
prompt = """
Project: {file_contents}
Status: {command_output}
"""
```

### Placeholders

**Universal** (available everywhere):
- `{date}` - Current timestamp

**Agent commands only:**
- `{bin}`, `{model}`, `{prompt}`, `{role}`, `{role_file}`

**UTD pattern** (roles/contexts/tasks):
- `{file}`, `{file_contents}`, `{command}`, `{command_output}`

**Tasks only:**
- `{instructions}` - User's command-line arguments

### Required vs Optional Contexts

Contexts have a `required` field controlling inclusion:

- `required = true` - Included in: `start`, `start prompt`, `start task`
- `required = false` - Included in: `start` only (excluded from prompts and tasks)

**Use required contexts** for essential information (environment, repo guidelines).

**Use optional contexts** for nice-to-have information that adds noise in focused queries.

## Customization Guide

### 1. Choose a Starting Point

- **New to `start`?** Start with `minimal/`
- **Setting up real project?** Use `real-world/`
- **Need specific feature?** Reference `complete/`

### 2. Install Global Config

```bash
# Choose your example
cp -r examples/real-world/global/* ~/.config/start/
```

### 3. Customize Agents

Edit `~/.config/start/agents.toml`:

- Remove agents you don't have installed
- Update model names to current versions
- Adjust default models
- Add custom agents

### 4. Customize Roles

Edit `~/.config/start/roles.toml`:

- Modify role prompts for your needs
- Add domain-specific roles
- Reference external role files if preferred

### 5. Set Up Required Contexts

Edit `~/.config/start/contexts.toml`:

- Update paths to your reference documents
- Mark essential contexts as `required = true`
- Add your own context documents

### 6. Add Useful Tasks

Edit `~/.config/start/tasks.toml`:

- Add tasks for your common workflows
- Set up aliases for quick access
- Reference appropriate roles and agents

### 7. Project-Specific Config

In your project root:

```bash
mkdir -p .start
cp -r examples/real-world/local/* ./.start/

# Customize for your project
$EDITOR .start/roles.toml
$EDITOR .start/contexts.toml
$EDITOR .start/tasks.toml
```

### 8. Verify Configuration

```bash
start doctor
start config show
```

## Common Patterns

### Pattern: Project README Role

```toml
# ./.start/roles.toml
[roles.project-expert]
file = "./ROLE.md"
```

Create `ROLE.md` in project root with project-specific guidelines.

### Pattern: Required Environment Context

```toml
# ~/.config/start/contexts.toml
[contexts.environment]
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
required = true
```

### Pattern: Git Diff Review Task

```toml
# ~/.config/start/tasks.toml or ./.start/tasks.toml
[tasks.review]
alias = "r"
role = "code-reviewer"
command = "git diff --staged"
prompt = "Review: {command_output}\n\nFocus: {instructions}"
```

Usage: `start task review "check for security issues"`

### Pattern: Dynamic Role with Command

```toml
# ~/.config/start/roles.toml
[roles.go-expert]
file = "~/.config/start/roles/go-base.md"
command = "go version"
prompt = "{file_contents}\n\nGo: {command_output}"
```

### Pattern: Optional Project Context

```toml
# ./.start/contexts.toml
[contexts.project]
file = "./PROJECT.md"
prompt = "Read {file} for project overview."
required = false
```

Included in `start` (full sessions), excluded from `start prompt` and `start task`.

## Troubleshooting

### Configuration not loading

```bash
# Check configuration paths
start doctor

# Show merged configuration
start config show
```

### Agent not found

```bash
# List configured agents
start config agent list

# Verify binary exists
which claude
which gemini
```

### Task validation errors

```bash
# Validate all configurations
start config validate

# Check task references
start config task list
```

### Context file not found

Check paths in contexts - use absolute paths or tilde expansion:

```toml
# Good
file = "~/reference/ENVIRONMENT.md"
file = "/absolute/path/to/file.md"
file = "./relative/to/working/dir.md"

# Bad (will not find file)
file = "relative-without-prefix.md"
```

## Next Steps

1. **Choose an example** that matches your needs
2. **Copy to appropriate directory** (global or local)
3. **Customize** agents, roles, contexts, and tasks
4. **Verify** with `start doctor`
5. **Test** with `start --help` and `start task help`
6. **Iterate** based on your workflow

## Additional Resources

- [Configuration Reference](../docs/config.md) - Complete config documentation
- [Architecture](../docs/architecture.md) - System design and patterns
- [Design Records](../docs/design/design-records/) - Design decisions
- [CLI Documentation](../docs/cli/) - Command reference

## Contributing

Found a useful pattern? Consider contributing:

1. Add example to appropriate directory
2. Document in this README
3. Submit pull request

Examples should be:
- Clear and well-documented
- Solve real problems
- Follow existing patterns
- Include comments explaining why
