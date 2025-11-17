# start - AI Agent CLI Tool

`start` is a command-line orchestrator for AI agents that manages prompt composition, context injection, and workflow automation. It wraps various AI CLI tools (Claude, Gemini, GPT, etc.) with configurable roles, reusable tasks, and project-aware context documents.

If you need a better understanding of the project, read the docs/cli/*.md documents.

> **Note:** This project is currently in the design and documentation phase. No implementation exists yet. All documentation describes the planned system architecture and behavior. There are no backward compatibility concerns or migration requirements at this stage.

## Core Concepts

- **Roles**: Define AI agent behavior and expertise (e.g., go-expert, code-reviewer)
- **Tasks**: Reusable prompts for common workflows (e.g., pre-commit-review, debug-help)
- **Contexts**: Environment-specific information loaded at runtime
- **Agents**: AI model configurations (Claude, GPT, Gemini, etc.)
- **Assets**: Downloadable catalog of roles, tasks, and configurations from GitHub

## Quick Reference

```bash
start                           # Start interactive session with default role
start --role go-expert          # Start with specific role
start task pre-commit-review    # Run a specific task
start init                      # Initialize configuration
start assets add                # Browse and install assets from catalog
```

## Architecture

- **Catalog-driven**: Assets stored in GitHub, downloaded on-demand, cached locally
- **Multi-file config**: Separate files for config, tasks, agents, roles, contexts
- **Lazy loading**: Downloads assets from GitHub when first needed
- **Offline-friendly**: Cached assets work without network access

## Documentation

Complete documentation is in the `docs/` directory:

- `docs/cli/` - Command reference for all CLI commands
- `docs/design/` - Design decisions and architecture
- `docs/config.md` - Configuration reference

For detailed information about commands, configuration, and design decisions, refer to the documentation files.

---

## Document Driven Development (DDD)

This project uses Document Driven Development. Design decisions are documented in Design Records (DRs) before implementation.

**For complete DR writing guidelines:** See [docs/design/dr-writing-guide.md](docs/design/dr-writing-guide.md)

**Location:** `docs/design/design-records/`

**Quick reference:**

- Create DRs for: Architectural decisions, algorithms, breaking changes, API/CLI structure
- Get next DR number from: `docs/design/design-records/README.md`
- Follow reconciliation process after 5-10 DRs (see dr-writing-guide.md)
