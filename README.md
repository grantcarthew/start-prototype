# start - AI Agent CLI Orchestrator

> **PROTOTYPE STATUS - NOT ACTIVELY DEVELOPED**
>
> This is a research prototype that validated the core concepts for an AI agent CLI orchestrator. Development has moved to a **CUE-based implementation** that addresses fundamental limitations discovered during this prototype phase.
>
> **Key Findings:**
> - TOML's lack of table ordering and limited validation made it unsuitable for context injection requirements
> - Custom GitHub-based asset distribution added unnecessary complexity
> - The CUE version uses native package registry and built-in validation
>
> **Current Status:** This prototype contains 44 Design Records documenting the research and design decisions. It remains valuable as reference documentation for the architectural exploration, but no new development will occur here.
>
> **Active Development:** [github.com/grantcarthew/start](https://github.com/grantcarthew/start) (CUE-based)

---

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://go.dev)

**start** is a command-line orchestrator for AI agents that manages prompt composition, context injection, and workflow automation. It wraps various AI CLI tools (Claude, Gemini, GPT, etc.) with configurable roles, reusable tasks, and project-aware context documents.

---

## Quick Start

```bash
# Install via Homebrew
brew tap grantcarthew/tap
brew install start

# Initialize configuration
start init

# Start interactive session
start

# Run a task
start task review "check security"

# Execute with specific role
start --role go-expert "optimize this function"
```

---

## Features

### 🎯 Context Management
- **Automatic context loading** - Project-specific documents loaded on every invocation
- **Global + local configs** - Personal defaults + project overrides
- **Required vs optional** - Control which contexts appear in tasks vs full sessions

### 🎭 Role System
- **System prompts as roles** - Define AI behavior with reusable roles
- **File, command, or inline** - Flexible role definition (Unified Template Design)
- **Dynamic content** - Execute commands to generate role context

### 📋 Task Workflows
- **Reusable tasks** - Define common workflows once, use everywhere
- **Aliases for speed** - Quick shortcuts for frequent tasks
- **Agent/role selection** - Tasks can specify which agent and role to use

### 🔌 Multi-Agent Support
- **Claude, Gemini, GPT, aichat** - Works with any CLI-based AI tool
- **Custom agents** - Define your own agent configurations
- **Model management** - Named models with aliases (e.g., "sonnet", "pro")

### 📦 Asset Catalog
- **GitHub-backed assets** - Downloadable roles, tasks, and configs
- **Lazy loading** - Assets fetched on-demand and cached locally
- **Offline-friendly** - Cached assets work without network

### 🛠️ Configuration Commands
- **Interactive wizards** - Create agents, roles, contexts, and tasks interactively
- **Validation** - Test configurations before using them
- **Doctor diagnostics** - Health checks for your setup

---

## Installation

### Homebrew (Recommended)

```bash
brew tap grantcarthew/tap
brew install start
```

### From Source

```bash
# Clone repository
git clone https://github.com/grantcarthew/start.git
cd start

# Build
go build -o bin/start cmd/start/main.go

# Install
sudo mv bin/start /usr/local/bin/
```

### Verify Installation

```bash
start --version
start doctor
```

---

## Configuration

### Quick Setup

```bash
# Interactive setup wizard
start init

# Or copy example configs
mkdir -p ~/.config/start
cp examples/real-world/global/* ~/.config/start/
```

### Configuration Structure

**Global** (`~/.config/start/`):
- `config.toml` - Settings and defaults
- `agents.toml` - AI agent configurations
- `roles.toml` - Role (system prompt) definitions
- `contexts.toml` - Context document configurations
- `tasks.toml` - Task workflow definitions

**Local** (`./.start/`):
- Same structure as global
- Project-specific overrides
- Local completely replaces global for same-named items

---

## Core Concepts

### Contexts

Context documents are automatically loaded and prepended to every prompt:

```toml
# ~/.config/start/contexts.toml
[contexts.environment]
file = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
required = true

[contexts.project]
file = "./PROJECT.md"
prompt = "Read {file} for project overview."
required = false  # Only in full sessions, not tasks/prompts
```

**Required vs Optional:**
- `required = true` - Included in `start`, `start prompt`, `start task`
- `required = false` - Included only in `start` (full interactive sessions)

### Roles

Roles define AI agent behavior (system prompts):

```toml
# ~/.config/start/roles.toml
[roles.go-expert]
file = "~/.config/start/roles/go-expert.md"

[roles.code-reviewer]
prompt = """
You are an expert code reviewer.
Focus on: security, performance, maintainability.
Date: {date}
"""

[roles.dynamic]
file = "./ROLE.md"
command = "git log -1 --oneline"
prompt = "{file_contents}\n\nLast commit: {command_output}"
```

**Unified Template Design (UTD):**
- `file` - Path to file (supports `{file}` and `{file_contents}` placeholders)
- `command` - Shell command to execute (supports `{command}` and `{command_output}` placeholders)
- `prompt` - Template with placeholders
- At least one field required, all three can be combined

### Tasks

Tasks are reusable workflows:

```toml
# ~/.config/start/tasks.toml or ./.start/tasks.toml
[tasks.review]
alias = "r"
description = "Code review with git diff"
role = "code-reviewer"
agent = "claude"
command = "git diff --staged"
prompt = """
Review the following changes:

{command_output}

Focus on: {instructions}
"""
```

**Usage:**
```bash
start task review "check for security issues"
start task r "verify error handling"  # Using alias
```

**Placeholders:**
- `{instructions}` - User's command-line arguments (defaults to "None")
- All UTD placeholders available (`{file}`, `{command_output}`, etc.)

### Agents

Agent configurations define how to invoke AI CLIs:

```toml
# ~/.config/start/agents.toml
[agents.claude]
bin = "claude"
command = "{bin} --model {model} {prompt}"
description = "Anthropic's Claude AI assistant"
url = "https://claude.ai"
models_url = "https://docs.anthropic.com/claude/docs/models-overview"
default_model = "sonnet"

  [agents.claude.models]
  haiku = "claude-haiku-4-20250409"
  sonnet = "claude-sonnet-4-20250929"
  opus = "claude-opus-4-20250514"
```

**Command placeholders:**
- `{bin}` - Agent binary name
- `{model}` - Resolved model name
- `{prompt}` - Assembled prompt (contexts + role + user input)
- `{role}` - Role content (for inline system prompts)
- `{role_file}` - Temporary file with role content
- `{date}` - Current ISO 8601 timestamp

---

## Usage Examples

### Interactive Session

```bash
# Start with default role
start

# Start with specific role
start --role go-expert

# Start with specific agent
start --agent gemini

# Combine flags
start --agent claude --role code-reviewer --model opus
```

### Direct Prompts

```bash
# Quick one-off prompt
start prompt "explain this error message"

# With role
start --role security-expert prompt "review this authentication code"
```

### Task Execution

```bash
# Run task with instructions
start task review "check for race conditions"

# Use task alias
start task r "verify input validation"

# Task with no instructions
start task commit-message
```

### Asset Management

```bash
# Search catalog
start assets search golang

# Browse catalog in browser
start assets browse

# View asset details
start assets info code-reviewer

# Add asset from catalog
start assets add go-expert

# Update cached assets
start assets update

# List cached assets
start assets index
```

### Configuration Management

```bash
# List configured items
start config agent list
start config role list
start config context list
start config task list

# Create new items (interactive wizards)
start config agent new
start config role new
start config context new
start config task new

# Show specific item
start config agent show claude
start config role show go-expert

# Test configuration
start config agent test claude
start config role test code-reviewer

# Set defaults
start config agent default claude
start config role default assistant

# Edit items (opens in editor)
start config agent edit claude
start config role edit go-expert

# Remove items
start config agent remove old-agent
start config role remove unused-role

# Show merged configuration
start config show
```

### Diagnostics

```bash
# Full health check
start doctor

# Quiet mode (exit code only)
start doctor --quiet

# Verbose mode
start doctor --verbose
```

### Shell Completion

```bash
# Generate completion script
start completion bash > /etc/bash_completion.d/start
start completion zsh > ~/.zsh/completion/_start
start completion fish > ~/.config/fish/completions/start.fish

# Or source directly
source <(start completion bash)
```

---

## Command Reference

### Main Commands

| Command | Description |
|---------|-------------|
| `start` | Start interactive session with contexts and role |
| `start prompt <text>` | Execute one-off prompt with required contexts |
| `start task <name> [instructions]` | Run predefined workflow task |
| `start init` | Initialize configuration with wizard |
| `start doctor` | Run health checks and diagnostics |
| `start config show` | Display merged configuration |
| `start completion <shell>` | Generate shell completion script |

### Asset Commands

| Command | Description |
|---------|-------------|
| `start assets search <query>` | Search asset catalog (min 3 chars) |
| `start assets browse` | Open catalog in browser |
| `start assets info <query>` | Show detailed asset information |
| `start assets add <query>` | Download and install asset |
| `start assets update [query]` | Update cached assets |
| `start assets index` | List all cached assets |

### Config Commands

| Command | Description |
|---------|-------------|
| `start config agent list/new/show/edit/remove/test/default` | Manage agents |
| `start config role list/new/show/edit/remove/test/default` | Manage roles |
| `start config context list/new/show/edit/remove/test` | Manage contexts |
| `start config task list/new/show/edit/remove/test` | Manage tasks |

### Global Flags

| Flag | Description |
|------|-------------|
| `-a, --agent <name>` | Agent to use (overrides config default) |
| `-m, --model <name>` | Model to use (overrides agent default) |
| `-r, --role <name>` | Role to use (overrides config default) |

---

## Architecture

**Pattern:** Hexagonal Architecture (Ports and Adapters)

```
CLI Layer (Cobra)
    ↓
Engine Layer (Business Logic)
    ↓
Domain Layer (Interfaces + Models)
    ↓
Adapters Layer (Concrete Implementations)
```

**Key Principles:**
- Interface-based dependency injection
- Test-first development
- Domain-driven design
- Zero external dependencies in core logic

**See:** [docs/architecture.md](docs/architecture.md) for complete details

---

## Examples

Comprehensive examples in `examples/`:

- **minimal/** - Quick start with bare minimum
- **complete/** - Full reference showing all possible fields
- **real-world/** - Practical working configuration

**See:** [examples/README.md](examples/README.md) for detailed guide

---

## Documentation

### User Guides
- [Configuration Reference](docs/config.md) - Complete TOML specification
- [CLI Commands](docs/cli/) - Detailed command documentation
- [Examples](examples/README.md) - Configuration examples

### Developer Guides
- [Architecture](docs/architecture.md) - System design and patterns
- [Testing Strategy](docs/testing.md) - Test approach and smith agent
- [Design Records](docs/design/design-records/) - Design decisions (DR-001 to DR-043)
- [Implementation Phases](docs/implementation/) - Development roadmap

### Key Design Decisions
- [DR-001: TOML Format](docs/design/design-records/dr-001-toml-format.md)
- [DR-002: Config Merge](docs/design/design-records/dr-002-config-merge.md)
- [DR-005: Role Configuration](docs/design/design-records/dr-005-role-configuration.md)
- [DR-031: Catalog Architecture](docs/design/design-records/dr-031-catalog-based-assets.md)
- [DR-042: Missing Asset Restoration](docs/design/design-records/dr-042-missing-asset-restoration.md)
- [DR-043: Process Replacement](docs/design/design-records/dr-043-process-replacement.md)

---

## Development

### Building

```bash
# Build main binary
go build -o bin/start cmd/start/main.go

# Build with version
go build -ldflags "-X main.version=1.0.0" -o bin/start cmd/start/main.go

# Build smith test agent
go build -o bin/smith cmd/smith/main.go
```

### Testing

```bash
# All tests (unit + integration)
./test.sh

# Unit tests only
go test -short ./...

# Integration tests only
go test ./test/integration/...

# With coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Project Structure

```
start/
├── cmd/
│   ├── start/              # Main entry point
│   └── smith/              # Test agent
├── internal/
│   ├── domain/             # Core models and interfaces
│   ├── config/             # TOML configuration
│   ├── engine/             # Business logic
│   ├── assets/             # Asset management
│   ├── cli/                # Cobra commands
│   └── adapters/           # External integrations
├── test/
│   ├── integration/        # End-to-end tests
│   ├── mocks/              # Mock implementations
│   └── assert/             # Test helpers
├── docs/                   # Documentation
├── examples/               # Configuration examples
└── assets/                 # Asset catalog
```

---

## Contributing

Contributions welcome! Please:

1. Read the [Architecture](docs/architecture.md) and [Design Records](docs/design/design-records/)
2. Follow existing code patterns
3. Add tests for new features
4. Update documentation
5. Submit pull request

**For major changes:** Open an issue first to discuss the approach.

---

## Troubleshooting

### Configuration not loading

```bash
# Check configuration paths and validation
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

# Test agent configuration
start config agent test claude
```

### Task validation errors

```bash
# List available tasks
start config task list

# Show task details
start config task show review

# Test task configuration
start config task test review
```

### Context file not found

Check file paths in contexts - use absolute paths or tilde expansion:

```toml
# ✓ Good
file = "~/reference/ENVIRONMENT.md"
file = "/absolute/path/to/file.md"
file = "./relative/to/working/dir.md"

# ✗ Bad (will not find file)
file = "relative-without-prefix.md"
```

### Network issues with assets

```bash
# Check GitHub connectivity
start assets search test

# Use cached assets only (disable downloads)
# Edit ~/.config/start/config.toml:
[settings]
asset_download = false
```

---

## Roadmap

### v1.0.0 (Current)
- ✅ Core CLI with agent execution
- ✅ Role and context system
- ✅ Task workflows
- ✅ Asset catalog with lazy loading
- ✅ Init wizard and config commands
- ✅ Shell completion
- ✅ Doctor diagnostics

### Future Enhancements
- Process replacement (syscall.Exec) for seamless handoff
- Prompt templates with variables
- Environment-specific configurations
- Asset versioning and pinning
- Configuration validation on save
- Enhanced error messages with suggestions

---

## License

MIT License - see [LICENSE](LICENSE) for details

---

## Acknowledgments

- Inspired by the need for consistent AI agent workflows
- Built with [Cobra](https://github.com/spf13/cobra) CLI framework
- Uses [go-toml](https://github.com/pelletier/go-toml) for configuration
- Follows Document Driven Development principles

---

## Links

- **GitHub:** https://github.com/grantcarthew/start
- **Issues:** https://github.com/grantcarthew/start/issues
- **Homebrew Tap:** https://github.com/grantcarthew/homebrew-tap

---

**Built with ❤️ by Grant Carthew**
