# Project: start - Catalog-Based Asset System Design

**Status:** Design Phase - Catalog Architecture
**Date Started:** 2025-01-10
**Previous Phase:** CLI Design (completed, see docs/archive/2025-01-10-cli-configuration-design-phase.md)

**Working Mode:** Design then implement - Make key architectural decisions, document thoroughly, then build.

## Overview

Redesign the asset system from bulk download/merge to a **catalog-driven model** where assets are discovered via GitHub, downloaded on-demand, cached locally, and lazily loaded at runtime.

**Key Shift:** GitHub as database, filesystem as state, lazy loading, on-demand installation.

**Links:** [Catalog Ideas](./docs/ideas/catalog-based-assets.md) | [Previous Phase: CLI Design](./docs/archive/2025-01-10-cli-configuration-design-phase.md)

## Vision

```bash
# Browse catalog interactively
$ start config task add
> [Beautiful TUI shows categories and tasks]
> Downloads on selection, caches locally

# Just run it (lazy loading)
$ start task pre-commit-review
> Task not found locally
> Queries GitHub, downloads, caches, runs
> Next time: instant from cache

# Keep it fresh
$ start update
> Checks cached assets for updates
> Interactive or automatic mode
```

## Core Principles

1. **GitHub as Source of Truth** - All assets in GitHub repo, not bundled
2. **Lazy Loading** - Download on first use, not upfront
3. **Filesystem = State** - No tracking files, cache structure is the state
4. **On-Demand Installation** - Browse and install what you need
5. **Hash-Based Versioning** - SHA comparison for updates
6. **Offline-Friendly** - Cached assets work offline, manual config always possible

## Key Design Decisions

| Decision | Choice | Status | DR |
|----------|--------|--------|----|
| Metadata format | Sidecar `.meta.toml` files | ✅ Decided | TBD |
| Installation model | Lazy loading - download on first use | ✅ Decided | TBD |
| State management | Filesystem = state, no tracking file | ✅ Decided | TBD |
| Browse/search | Go TUI libraries (bubbletea/promptui) | ✅ Decided | TBD |
| Update strategy | Interactive (default) or `--auto` flag | ✅ Decided | TBD |
| Versioning | SHA-based (git blob or content hash) | ✅ Decided | TBD |
| Offline behavior | Manual config only, no GitHub access | ✅ Decided | DR-026 |
| Resolution order | local → global → cache → GitHub → error | ✅ Decided | TBD |
| GitHub token | Recommend `GITHUB_TOKEN` env var | ✅ Decided | TBD |
| Asset types (v1) | roles, tasks, agents, templates | ✅ Decided | TBD |
| Minimal viable set | 28 assets (8 roles, 12 tasks, 6 agents, 2 templates) | ✅ Decided | TBD |

## Open Questions

### High Priority (Blockers)

1. **Sidecar metadata schema** - Exact fields in `.meta.toml`?
   - Required: name, category, description, sha, tags
   - Optional: version, author, created, updated, dependencies
   - Schema validation?

2. **GitHub API strategy** - Which endpoints, rate limiting, caching?
   - Tree API for browsing vs Contents API for downloading
   - How to handle rate limits gracefully
   - Cache catalog index locally?

3. **Asset resolution algorithm** - Exact priority order and fallback behavior?
   - local config → global config → cache → GitHub
   - When to prompt user vs auto-download
   - Error messages for each failure case

4. **Cache management** - When to clean up, size limits, manual control?
   - `start cache clean` command?
   - Age-based cleanup?
   - Size-based cleanup?

5. **Config integration** - How do cached assets get added to config?
   - Automatic on first use?
   - Manual add command?
   - Prompt user each time?

### Medium Priority (Important)

6. **TUI library choice** - bubbletea vs promptui vs basic numbered selection?
   - Feature comparison
   - Dependency size
   - Fallback strategy

7. **Asset validation** - How to validate downloaded assets?
   - TOML validation
   - Schema checking
   - Security considerations

8. **Update notification** - When to check for updates?
   - On `start update` only?
   - Weekly background check?
   - Never automatic (DR-025 compliance)

9. **Error recovery** - What if download fails mid-way?
   - Retry logic
   - Partial file cleanup
   - User notification

10. **Asset dependencies** - Can tasks require specific roles/agents?
    - Dependency declaration in metadata
    - Automatic installation of dependencies
    - Version compatibility

### Low Priority (Nice to Have)

11. **Community assets** - How can users contribute?
    - PR process to main repo?
    - User-hosted repos?
    - Namespacing: `user/repo/path`

12. **Search functionality** - Full-text search across assets?
    - Search by description, tags, name
    - Fuzzy matching
    - Integration with TUI

13. **Analytics** - Track popular assets (privacy-respecting)?
    - Opt-in telemetry
    - Help prioritize development
    - No PII collection

14. **Asset preview** - Preview asset before downloading?
    - Show metadata (description, tags)
    - Show actual content
    - Diff between versions

15. **Bulk operations** - Install multiple assets at once?
    - Install entire category
    - Install from list/file
    - Workspace templates

## Design Records Needed

### New DRs

1. **DR-031: Catalog-Based Asset Architecture** (High Priority)
   - Overall system design
   - Lazy loading model
   - GitHub as database
   - State management via filesystem
   - Supersedes aspects of DR-014, DR-015, DR-016, DR-019, DR-023

2. **DR-032: Asset Metadata Schema** (High Priority)
   - Sidecar `.meta.toml` format
   - Required vs optional fields
   - Versioning with SHA
   - Schema validation

3. **DR-033: Asset Resolution Algorithm** (High Priority)
   - Priority order: local → global → cache → GitHub
   - Lazy loading behavior
   - User prompts vs auto-download
   - Error handling and messages

4. **DR-034: GitHub Catalog API Strategy** (High Priority)
   - API endpoints used
   - Rate limiting handling
   - Authentication with GITHUB_TOKEN
   - Caching strategy

5. **DR-035: Interactive Asset Browsing** (Medium Priority)
   - TUI library choice
   - Category navigation
   - Search/filter functionality
   - Fallback for non-interactive environments

6. **DR-036: Cache Management** (Medium Priority)
   - Cache structure and location
   - Cleanup policies
   - Manual cache commands
   - Size management

7. **DR-037: Asset Update Mechanism** (Medium Priority)
   - SHA-based update detection
   - Interactive vs automatic updates
   - Update notifications (compliant with DR-025)
   - Rollback on failure

### Updates to Existing DRs

8. **DR-014: GitHub Tree API** → Add note: "See DR-031 for catalog model"
9. **DR-015: Atomic Updates** → Add note: "See DR-031 for per-asset updates"
10. **DR-016: Asset Discovery** → Add note: "See DR-031 for interactive browsing"
11. **DR-019: Task Loading** → Update with cache in resolution order
12. **DR-023: Staleness Checking** → Update with per-asset SHA comparison

## Command Updates Required

### New Behavior

**`start config task add`** - Interactive catalog browsing
```bash
start config task add                              # Browse catalog
start config task add git-workflow/pre-commit-review  # Direct install
start config task add --category git-workflow      # Filter by category
start config task add --search "commit"            # Search catalog
```

**`start config role add`** - Same pattern as tasks
**`start config agent add`** - Same pattern as tasks

**`start task <name>`** - Lazy loading
```bash
# If task not in config or cache, query GitHub
# Prompt to download and cache
# Run immediately
```

**`start update`** - Update cached assets
```bash
start update            # Interactive (default)
start update --auto     # Automatic mode
start update --check    # Check only, don't update
```

**`start cache`** - New command for cache management (optional)
```bash
start cache list        # Show cached assets
start cache clean       # Remove all cached assets
start cache clean --old # Remove assets older than 30 days
start cache info        # Show cache size, location
```

### Command Specs to Update

- [x] `start-config-agent.md` - Add interactive browsing flow
- [x] `start-config-task.md` - Add interactive browsing flow
- [ ] `start-config-role.md` - Add interactive browsing flow (NEW)
- [ ] `start-task.md` - Add lazy loading behavior
- [ ] `start-update.md` - Update to per-asset model
- [ ] `start-cache.md` - New cache management command (optional)

## Configuration Updates

### Settings Section

```toml
[settings]
default_agent = "claude"
default_role = "default"
log_level = "normal"
shell = "bash"
command_timeout = 30

# Asset management (NEW)
asset_path = "~/.config/start/assets"           # Cache location
github_token_env = "GITHUB_TOKEN"               # Env var for API token
asset_repo = "start-project/start-assets"       # GitHub repository
```

### Asset Cache Structure

```
~/.config/start/assets/
├── roles/
│   ├── general/
│   │   ├── code-reviewer.md
│   │   └── code-reviewer.meta.toml
│   └── languages/
│       ├── go-expert.md
│       └── go-expert.meta.toml
├── tasks/
│   └── git-workflow/
│       ├── pre-commit-review.toml
│       └── pre-commit-review.meta.toml
└── agents/
    └── claude/
        ├── sonnet.toml
        └── sonnet.meta.toml
```

## Minimal Viable Asset Set (v1)

### Roles (8)
- **general/** (4): default, code-reviewer, pair-programmer, explainer
- **languages/** (2): go-expert, python-expert
- **specialized/** (2): security-focused, rubber-duck

### Tasks (12)
- **git-workflow/** (4): pre-commit-review, pr-ready, commit-message, explain-changes
- **code-quality/** (4): find-bugs, quick-wins, naming-review, test-suggestions
- **security/** (2): security-scan, dependency-audit
- **debugging/** (2): debug-help, git-story

### Agents (6)
- **claude/** (3): sonnet, opus, haiku
- **openai/** (2): gpt-4, gpt-4-turbo
- **google/** (1): gemini-pro

### Templates (2)
- **projects/** (2): solo-developer, team-project

**Total: 28 assets** - Enough to impress, not overwhelming

See [catalog-based-assets.md](./docs/ideas/catalog-based-assets.md) for complete asset ideas and future possibilities.

## Implementation Phases

### Phase 1: Core Catalog Infrastructure (Foundation)

**Goal:** Basic GitHub browsing and caching

- [ ] DR-031: Catalog architecture design
- [ ] DR-032: Metadata schema design
- [ ] DR-033: Resolution algorithm design
- [ ] DR-034: GitHub API strategy design
- [ ] GitHub API client (tree, contents, auth)
- [ ] Cache management (read, write, validate)
- [ ] Metadata parsing and validation
- [ ] Asset resolution algorithm

**Validation:** Can browse GitHub catalog and cache assets locally

### Phase 2: Interactive Browsing (UX)

**Goal:** Beautiful catalog browsing experience

- [ ] DR-035: Interactive browsing design
- [ ] TUI library integration (bubbletea or promptui)
- [ ] Category navigation
- [ ] Asset selection and preview
- [ ] Download progress indicators
- [ ] Error handling and user feedback

**Validation:** Can interactively browse and install assets

### Phase 3: Lazy Loading (Smart Behavior)

**Goal:** Assets download on-demand when needed

- [ ] Task resolution with GitHub fallback
- [ ] User prompts for download confirmation
- [ ] Automatic caching after download
- [ ] Optional config integration
- [ ] Performance optimization (parallel downloads)

**Validation:** `start task <name>` works for uncached GitHub assets

### Phase 4: Update System (Freshness)

**Goal:** Keep cached assets up-to-date

- [ ] DR-036: Cache management design
- [ ] DR-037: Update mechanism design
- [ ] SHA-based update detection
- [ ] Interactive update flow
- [ ] Automatic update mode
- [ ] Update notifications (DR-025 compliant)

**Validation:** `start update` checks and updates cached assets

### Phase 5: Asset Content (Value)

**Goal:** Ship with 28 high-quality assets

- [ ] Write 8 role prompts (markdown)
- [ ] Write 12 task definitions (toml)
- [ ] Write 6 agent configs (toml)
- [ ] Write 2 project templates (toml)
- [ ] Create all metadata files
- [ ] Validate all assets
- [ ] Test each asset end-to-end

**Validation:** All 28 assets work perfectly

### Phase 6: Polish (Quality)

**Goal:** Production-ready catalog system

- [ ] Cache cleanup commands
- [ ] Search/filter functionality
- [ ] Asset validation and security
- [ ] Error recovery and retry logic
- [ ] Performance optimization
- [ ] Documentation and examples
- [ ] GitHub token setup guide

**Validation:** Smooth, reliable, professional experience

## Success Criteria

Catalog design is complete when:

- [ ] All high-priority questions resolved (5 questions)
- [ ] All new DRs written (DR-031 through DR-037)
- [ ] Existing DRs updated with notes
- [ ] Command specs updated for new behavior
- [ ] Configuration schema updated
- [ ] Cache structure defined
- [ ] Asset metadata schema finalized
- [ ] Resolution algorithm specified
- [ ] GitHub API strategy documented
- [ ] Minimal viable asset set defined (28 assets)

Implementation is complete when:

- [ ] All 6 implementation phases done
- [ ] 28 assets created and tested
- [ ] GitHub asset repository live
- [ ] Can browse catalog interactively
- [ ] Lazy loading works
- [ ] Update mechanism works
- [ ] Documentation complete
- [ ] Ready for first users

## Future Enhancements

**Beyond v1:**
- Community asset contributions
- User-hosted asset repositories
- Asset search and filtering
- Dependency management
- Asset preview and diff
- Bulk operations
- Context, metaprompt, snippet asset types
- Workflow chaining
- Analytics (opt-in)

See [catalog-based-assets.md](./docs/ideas/catalog-based-assets.md) for complete future vision.

## Reference

**Key Documents:**
- [docs/ideas/catalog-based-assets.md](./docs/ideas/catalog-based-assets.md) - Complete brainstorming and vision
- [docs/archive/2025-01-10-cli-configuration-design-phase.md](./docs/archive/2025-01-10-cli-configuration-design-phase.md) - Completed CLI design phase
- [docs/config.md](./docs/config.md) - Configuration reference (needs updates)
- [docs/design/design-record.md](./docs/design/design-record.md) - All design decisions

**Related DRs (to update):**
- DR-014: GitHub Tree API
- DR-015: Atomic Updates
- DR-016: Asset Discovery
- DR-019: Task Loading
- DR-023: Staleness Checking
- DR-026: Offline Behavior (already aligned)

## Recent Progress

### Catalog Architecture Brainstorming (2025-01-10)

**Session Context:**
- Working on Task 22 (out-of-box assets) from CLI design phase
- Realized bulk download model was wrong
- Brainstormed catalog-driven architecture
- Defined minimal viable asset set (28 assets)
- Made key architectural decisions

**Key Insights:**
- GitHub as database, not just distribution
- Lazy loading > bulk download
- Filesystem = state (no tracking files)
- Sidecar metadata keeps content clean
- On-demand installation is better UX

**Decisions Made:**
1. ✅ Metadata format: Sidecar `.meta.toml`
2. ✅ Installation model: Lazy loading
3. ✅ State management: Filesystem only
4. ✅ Browsing: Go TUI libraries
5. ✅ Updates: Interactive or `--auto`
6. ✅ Versioning: SHA-based
7. ✅ Offline: Manual config only
8. ✅ Minimal set: 28 assets

**Documents Created:**
- `docs/ideas/catalog-based-assets.md` - Complete vision and ideas
- `PROJECT-catalog-redesign.md` - This file, tracks new design phase

**Next Steps:**
1. Resolve high-priority open questions
2. Write DR-031 through DR-037
3. Update impacted DRs with notes
4. Update command specs
5. Begin implementation

## Notes

- This redesign does not invalidate the CLI design (PROJECT-cli-design.md)
- Commands remain the same, behavior changes
- Configuration structure remains the same, asset handling changes
- This is an enhancement, not a breaking change
- Users can still define everything manually (no GitHub required)
- GitHub token recommended but not required (graceful degradation)
