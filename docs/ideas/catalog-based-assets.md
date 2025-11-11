# Catalog-Based Asset System

**Date:** 2025-01-10
**Status:** Proposed Architecture
**Context:** Session brainstorming for Task 22 (out-of-box assets) revealed a better model

## Overview

Transform `start` from a tool with bundled assets to a **catalog-driven system** where assets are discovered, downloaded on-demand, and cached locally. Think Homebrew for AI development workflows.

## Core Concept

**GitHub as Database**
- Asset catalog lives in GitHub repository
- Users browse/search via CLI
- Download on first use (lazy loading)
- Cache locally for offline use
- Update check compares SHAs

**State Management**
- Filesystem IS the state
- No tracking files needed
- If file exists in cache → it's available
- If `.meta.toml` exists → we can check for updates

## Asset Types

### Currently Planned (v1)

1. **Roles** - System prompt templates (`.md` files)
   - Define agent behavior and expertise
   - Examples: code-reviewer, pair-programmer, go-expert

2. **Tasks** - Workflow definitions (`.toml` files)
   - Pre-configured commands with prompts
   - Examples: pre-commit-review, pr-ready, security-scan

3. **Agents** - AI tool configurations (`.toml` files)
   - Provider-specific settings
   - Examples: claude/sonnet, openai/gpt-4

4. **Templates** - Full config examples (`.toml` files)
   - Complete starter configurations
   - Examples: solo-developer, team-project

### Future Possibilities

5. **Contexts** - Document templates (`.md` files)
   - Standard documentation to include
   - Examples: api-guidelines, security-checklist

6. **Metaprompts** - Reusable prompt components (`.toml` files)
   - Mix-and-match behaviors
   - Examples: output-as-checklist, minimal-changes

7. **Snippets** - Common command patterns (`.toml` files)
   - Reusable shell commands
   - Examples: git-workflows, test-runners

8. **Workflows** - Multi-step task chains (`.toml` files)
   - Sequential task execution
   - Examples: full-pr-review = pre-commit + tests + docs + security

## Directory Structure

```
assets/
├── roles/
│   ├── general/
│   │   ├── default.md
│   │   ├── default.meta.toml
│   │   ├── code-reviewer.md
│   │   ├── code-reviewer.meta.toml
│   │   ├── pair-programmer.md
│   │   ├── pair-programmer.meta.toml
│   │   └── explainer.md
│   │       explainer.meta.toml
│   ├── languages/
│   │   ├── go-expert.{md,meta.toml}
│   │   ├── python-expert.{md,meta.toml}
│   │   ├── rust-expert.{md,meta.toml}
│   │   └── typescript-expert.{md,meta.toml}
│   ├── security/
│   │   ├── security-focused.{md,meta.toml}
│   │   └── penetration-tester.{md,meta.toml}
│   ├── specialized/
│   │   ├── architect.{md,meta.toml}
│   │   ├── performance-optimizer.{md,meta.toml}
│   │   └── accessibility-advocate.{md,meta.toml}
│   └── creative/
│       ├── rubber-duck.{md,meta.toml}
│       └── socratic-teacher.{md,meta.toml}
│
├── tasks/
│   ├── git-workflow/
│   │   ├── pre-commit-review.{toml,meta.toml}
│   │   ├── pr-ready.{toml,meta.toml}
│   │   ├── commit-message.{toml,meta.toml}
│   │   └── explain-changes.{toml,meta.toml}
│   ├── code-quality/
│   │   ├── find-bugs.{toml,meta.toml}
│   │   ├── quick-wins.{toml,meta.toml}
│   │   ├── naming-review.{toml,meta.toml}
│   │   └── test-suggestions.{toml,meta.toml}
│   ├── security/
│   │   ├── security-scan.{toml,meta.toml}
│   │   ├── dependency-audit.{toml,meta.toml}
│   │   └── threat-modeling.{toml,meta.toml}
│   ├── performance/
│   │   ├── performance-check.{toml,meta.toml}
│   │   ├── profiling-analysis.{toml,meta.toml}
│   │   └── optimization-suggestions.{toml,meta.toml}
│   ├── architecture/
│   │   ├── api-review.{toml,meta.toml}
│   │   ├── breaking-changes.{toml,meta.toml}
│   │   └── architectural-review.{toml,meta.toml}
│   ├── documentation/
│   │   ├── doc-review.{toml,meta.toml}
│   │   ├── onboarding-guide.{toml,meta.toml}
│   │   └── api-docs-gen.{toml,meta.toml}
│   └── debugging/
│       ├── debug-help.{toml,meta.toml}
│       ├── git-story.{toml,meta.toml}
│       └── root-cause-analysis.{toml,meta.toml}
│
├── agents/
│   ├── claude/
│   │   ├── sonnet.{toml,meta.toml}
│   │   ├── opus.{toml,meta.toml}
│   │   └── haiku.{toml,meta.toml}
│   ├── openai/
│   │   ├── gpt-4.{toml,meta.toml}
│   │   └── gpt-4-turbo.{toml,meta.toml}
│   ├── google/
│   │   └── gemini-pro.{toml,meta.toml}
│   └── local/
│       └── ollama.{toml,meta.toml}
│
├── contexts/
│   ├── standards/
│   │   ├── api-design-guidelines.{md,meta.toml}
│   │   ├── testing-standards.{md,meta.toml}
│   │   └── security-checklist.{md,meta.toml}
│   ├── team/
│   │   ├── code-style.{md,meta.toml}
│   │   └── pr-template.{md,meta.toml}
│   └── frameworks/
│       ├── react-patterns.{md,meta.toml}
│       └── go-idioms.{md,meta.toml}
│
├── metaprompts/
│   ├── output/
│   │   ├── checklist.{toml,meta.toml}
│   │   ├── diff.{toml,meta.toml}
│   │   ├── table.{toml,meta.toml}
│   │   └── json.{toml,meta.toml}
│   ├── behavior/
│   │   ├── no-execution.{toml,meta.toml}
│   │   ├── minimal-changes.{toml,meta.toml}
│   │   ├── ask-first.{toml,meta.toml}
│   │   └── test-driven.{toml,meta.toml}
│   └── lens/
│       ├── security.{toml,meta.toml}
│       ├── performance.{toml,meta.toml}
│       ├── accessibility.{toml,meta.toml}
│       └── maintainability.{toml,meta.toml}
│
├── snippets/
│   ├── git/
│   │   ├── common-commands.{toml,meta.toml}
│   │   └── advanced-workflows.{toml,meta.toml}
│   ├── testing/
│   │   ├── go-test.{toml,meta.toml}
│   │   └── pytest.{toml,meta.toml}
│   └── build/
│       └── make-targets.{toml,meta.toml}
│
└── templates/
    ├── projects/
    │   ├── solo-developer.{toml,meta.toml}
    │   ├── team-project.{toml,meta.toml}
    │   └── open-source.{toml,meta.toml}
    └── agents/
        └── multi-agent-compare.{toml,meta.toml}
```

## Metadata Format

**Sidecar `.meta.toml` files** - Keeps content clean, metadata separate

```toml
# pre-commit-review.meta.toml
name = "pre-commit-review"
category = "git-workflow"
description = "Review staged changes before committing"
tags = ["git", "review", "quality", "pre-commit"]
sha = "a1b2c3d4e5f6..."  # Content hash or git blob SHA
version = "1.0.0"
author = "start-project"
created = "2025-01-10T00:00:00Z"
updated = "2025-01-10T00:00:00Z"
```

**Why sidecar?**
- ✅ Content files stay clean (no frontmatter clutter for AI)
- ✅ Metadata separate from what agent sees
- ✅ Easy to parse independently
- ✅ No risk of corrupting content when updating metadata
- ❌ More files to manage (acceptable trade-off)

## User Workflows

### Workflow 1: Browse and Install

```bash
$ start config task add

Fetching catalog from GitHub...
✓ Found 42 tasks across 7 categories

Select category:
  1. git-workflow (4 tasks)
  2. code-quality (4 tasks)
  3. security (2 tasks)
  4. debugging (2 tasks)
  5. [view all]
> 1

git-workflow tasks:
  1. pre-commit-review - Review staged changes before commit
  2. pr-ready - Complete PR preparation
  3. commit-message - Generate conventional commit message
  4. explain-changes - Understand what changed in commits
> 1

Download 'pre-commit-review'...
✓ Cached to ~/.config/start/assets/tasks/git-workflow/

Add to config? [Y/n] y
Scope: [g]lobal or [l]ocal? [g]

✓ Added to global config as 'pre-commit-review' (alias: pcr)

Try it: start task pre-commit-review
    or: start task pcr
```

### Workflow 2: Direct Install

```bash
# Skip browsing, install directly
$ start config task add git-workflow/pre-commit-review --global

Downloading git-workflow/pre-commit-review...
✓ Cached to ~/.config/start/assets/tasks/git-workflow/
✓ Added to global config

Try it: start task pre-commit-review (pcr)
```

### Workflow 3: Lazy Loading (Just Run It)

```bash
# User runs task that doesn't exist
$ start task pre-commit-review

Resolving task 'pre-commit-review'...
  ✗ Not in local config
  ✗ Not in global config
  ✗ Not in asset cache
  ✓ Found in GitHub catalog: tasks/git-workflow/pre-commit-review.toml

Download and cache? [Y/n] y

Downloading...
✓ Cached to ~/.config/start/assets/tasks/git-workflow/
Running task 'pre-commit-review'...

[task executes]

Tip: Add to config for faster access: start config task add git-workflow/pre-commit-review
```

### Workflow 4: Update Cached Assets

```bash
# Interactive update (default)
$ start update

Checking 12 cached assets for updates...
  ✓ role: code-reviewer (up to date)
  ⚠ task: pre-commit-review (update available)
    Current: v1.0.0 (sha: abc123...)
    Latest:  v1.1.0 (sha: def456...)
    Changes: Improved error handling, updated prompts
  Update? [y/N] y
  ✓ Updated pre-commit-review

  ✓ task: pr-ready (up to date)

Summary: Updated 1/12 cached assets

# Automatic update
$ start update --auto

Updating all cached assets...
  ✓ code-reviewer (up to date)
  ✓ pre-commit-review (updated)
  ✓ pr-ready (up to date)

Updated 1/12 cached assets
```

### Workflow 5: List Assets

```bash
# Show what's configured vs cached vs available
$ start config task list

Configured tasks:
  global:
    - pre-commit-review (pcr) [cached: git-workflow/]
    - pr-ready (pr) [cached: git-workflow/]
  local:
    - custom-review (cr) [user-defined]

Cached (not in config):
  - explain-changes [git-workflow/]
  - find-bugs [code-quality/]

Available in catalog: 42 tasks (2 configured, 4 cached)
Run 'start config task add' to browse
```

## Resolution Order

When user runs `start task <name>`:

1. **Local config** (`.start/config.toml`)
2. **Global config** (`~/.config/start/config.toml`)
3. **Asset cache** (`~/.config/start/assets/`)
4. **GitHub catalog** (lazy fetch and cache)
5. **Error:** "Task 'xyz' not found locally or in GitHub catalog"

This enables lazy loading while respecting user customization.

## Settings Configuration

```toml
[settings]
default_agent = "claude"
default_role = "default"
log_level = "normal"
shell = "bash"
command_timeout = 30

# Asset management
asset_download = true                           # Auto-download from GitHub if not found
asset_path = "~/.config/start/assets"           # Where assets are cached
github_token_env = "GITHUB_TOKEN"               # Env var for GitHub API
asset_repo = "start-project/start-assets"       # GitHub repo
```

**GitHub Token:**
- Recommended for all users (prevents rate limiting)
- Anonymous: 60 requests/hour
- Authenticated: 5,000 requests/hour
- Set via: `export GITHUB_TOKEN=ghp_xxx`

## Minimal Viable Asset Set (v1)

Ship lean, prove value, iterate:

### Roles (8 total)
**general/** (4)
- default.md - Balanced, helpful, coding-focused
- code-reviewer.md - Strict quality/security review
- pair-programmer.md - Collaborative thinking
- explainer.md - Teaching mode, simplifies concepts

**languages/** (2)
- go-expert.md - Deep Go knowledge, idioms
- python-expert.md - Pythonic patterns

**specialized/** (2)
- security-focused.md - OWASP, paranoid mode
- rubber-duck.md - Only asks questions, helps YOU think 🦆

### Tasks (12 total)
**git-workflow/** (4)
- pre-commit-review.toml - Review staged changes
- pr-ready.toml - Complete PR preparation
- commit-message.toml - Generate conventional commit
- explain-changes.toml - Understand what changed

**code-quality/** (4)
- find-bugs.toml - Potential bugs and edge cases
- quick-wins.toml - Low-hanging refactoring fruit
- naming-review.toml - Better variable/function names
- test-suggestions.toml - What tests are missing

**security/** (2)
- security-scan.toml - Security-focused review
- dependency-audit.toml - Check dependencies

**debugging/** (2)
- debug-help.toml - Interactive debugging assistance
- git-story.toml - Code archaeology, why was it written this way

### Agents (6 total)
**claude/**
- sonnet.toml - Balanced (recommended default)
- opus.toml - Deep thinking
- haiku.toml - Fast iteration

**openai/**
- gpt-4.toml - Alternative provider
- gpt-4-turbo.toml - Faster GPT-4

**google/**
- gemini-pro.toml - Google's offering

### Templates (2 total)
**projects/**
- solo-developer.toml - Minimal config example
- team-project.toml - Full-featured config example

**Total v1: 28 assets** across 4 types

## Future Asset Ideas

### Roles (Beyond v1)

**Languages:**
- javascript-expert, typescript-expert, rust-expert
- java-expert, csharp-expert, kotlin-expert

**Specialized:**
- architect - Design patterns, SOLID, high-level
- performance-optimizer - Speed and efficiency
- accessibility-advocate - WCAG, inclusive design
- devops-expert - CI/CD, deployment, infrastructure
- database-expert - SQL, schema design, optimization
- api-designer - REST, GraphQL, API best practices

**Creative:**
- socratic-teacher - Teaches through questions
- devil's-advocate - Challenges assumptions
- minimalist - Simplest solution always

### Tasks (Beyond v1)

**Performance:**
- performance-check - Find bottlenecks
- profiling-analysis - Analyze profiling data
- optimization-suggestions - Algorithm improvements

**Architecture:**
- api-review - API design patterns
- breaking-changes - Detect breaking changes
- architectural-review - High-level design

**Documentation:**
- doc-review - Documentation quality
- onboarding-guide - Generate project onboarding
- api-docs-gen - Auto-generate API docs

**Advanced:**
- threat-modeling - Security threat analysis
- root-cause-analysis - Deep problem investigation
- migration-plan - Plan code migrations
- second-opinion - Use different agent/model for alternative perspective

### Contexts (Future)

**Standards:**
- api-design-guidelines
- testing-standards
- security-checklist
- code-review-checklist

**Team:**
- code-style-guide
- pr-template
- incident-response

**Frameworks:**
- react-patterns
- go-idioms
- rust-ownership-rules

### Metaprompts (Future)

**Output Formats:**
- checklist - Markdown checklist with [ ] items
- diff - Git-style diff format
- table - Markdown table
- json - Structured JSON for parsing

**Behavioral Constraints:**
- no-execution - Explain only, don't run code
- minimal-changes - Smallest possible edits
- ask-first - Ask permission before changes
- test-driven - Write tests first

**Lenses:**
- security - OWASP mindset
- performance - Speed and efficiency focus
- accessibility - WCAG compliance
- maintainability - Long-term maintenance view

## Technical Implementation Notes

### GitHub API Usage

**Required endpoints:**
- `GET /repos/{owner}/{repo}/git/trees/{sha}?recursive=1` - Get directory tree
- `GET /repos/{owner}/{repo}/contents/{path}` - Get file content
- `GET /repos/{owner}/{repo}/commits?path={path}` - Get file SHA/history

**Rate limiting:**
- Use `GITHUB_TOKEN` environment variable
- Check `X-RateLimit-Remaining` header
- Cache responses where possible

### Versioning Strategy

**Hash-based versioning:**
- Use Git blob SHA or content hash
- Store in `.meta.toml` when downloaded
- Update check: compare local SHA with remote SHA
- No need for semver (content hash is version)

**Update detection:**
```go
// Check if update available
localSHA := readMetadata("pre-commit-review.meta.toml").SHA
remoteSHA := githubAPI.getFileSHA("tasks/git-workflow/pre-commit-review.toml")
if localSHA != remoteSHA {
    // Update available
}
```

### Interactive Selection

**Go libraries:**
- `github.com/charmbracelet/bubbletea` - Beautiful TUIs
- `github.com/manifoldco/promptui` - Simple prompts
- Fallback: Numbered list selection (no dependencies)

**Selection flow:**
1. Fetch directory tree from GitHub
2. Group by category (parse directory structure)
3. Present categories for selection
4. Present items in category
5. Download selected item + metadata
6. Cache locally
7. Optionally add to config

### Cache Structure

```
~/.config/start/assets/
├── roles/
│   └── general/
│       ├── code-reviewer.md
│       └── code-reviewer.meta.toml
├── tasks/
│   └── git-workflow/
│       ├── pre-commit-review.toml
│       └── pre-commit-review.meta.toml
└── agents/
    └── claude/
        ├── sonnet.toml
        └── sonnet.meta.toml
```

**State management:**
- Filesystem IS the state
- No tracking files needed
- If `.meta.toml` exists → we can check for updates
- If content file exists → it's cached and ready

### Offline Behavior

Consistent with DR-026 (offline fallback):

**If online:**
- Browse catalog
- Download on-demand
- Check for updates

**If offline:**
- Use cached assets
- Use configured assets
- Cannot browse catalog
- Error with helpful message: "Configure manually or reconnect"

## Benefits

### For Users
- 🚀 **Immediate value** - Ship with 28 curated assets
- 🔍 **Discoverable** - Browse catalog interactively
- 📦 **On-demand** - Only download what you use
- 🔄 **Always fresh** - Check for updates anytime
- 🎨 **Customizable** - Mix catalog + custom assets
- 💾 **Offline-friendly** - Cached assets work offline

### For Project
- 🧩 **Extensible** - Add asset types easily
- 📈 **Scalable** - Can grow to hundreds of assets
- 🔧 **Maintainable** - Update assets without releases
- 🌍 **Community-ready** - Others can contribute assets
- 🎯 **Focused** - Binary is code, content is assets

## Questions to Resolve

1. **Asset submission process** - How can community contribute assets?
2. **Quality control** - How to review/approve community assets?
3. **Asset namespacing** - Support user repos? `user/repo/path`?
4. **Search functionality** - Full-text search across descriptions?
5. **Dependency tracking** - Can tasks depend on specific roles?
6. **Version constraints** - Can user pin to specific versions?
7. **Analytics** - Track most-used assets (privacy-respecting)?

## Related Design Decisions

This architecture impacts several existing DRs:
- DR-014: GitHub Tree API (now for browsing, not bulk download)
- DR-015: Atomic updates (now per-asset, not bulk)
- DR-016: Asset discovery (now interactive browsing)
- DR-019: Task loading (now includes cache in resolution)
- DR-023: Staleness checking (now per-asset SHA comparison)

## Next Steps

1. Create PROJECT-catalog-redesign.md to track design decisions
2. Design records needed for catalog architecture
3. Update impacted DRs with notes
4. Build minimal viable 28 assets
5. Implement catalog browsing and caching
6. Test with real GitHub asset repository
