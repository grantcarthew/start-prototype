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
| Metadata format | Sidecar `.meta.toml` files (6 required fields) | ✅ Decided | DR-032 |
| Config structure | Multi-file (config, tasks, agents, contexts) | ✅ Decided | DR-031 |
| Installation model | Lazy loading - download on first use | ✅ Decided | DR-031 |
| State management | Filesystem = state, no tracking file | ✅ Decided | DR-031 |
| Browse/search | Numbered selection (v1), TUI future | ✅ Decided | DR-035 |
| Update strategy | Manual only, SHA-based detection | ✅ Decided | DR-037 |
| Versioning | SHA-based (git blob SHA) | ✅ Decided | DR-032 |
| Offline behavior | Manual config only, no GitHub access | ✅ Decided | DR-026 |
| Resolution order | local → global → cache → GitHub → error | ✅ Decided | DR-033 |
| GitHub API | Tree API + raw.githubusercontent.com | ✅ Decided | DR-034 |
| Cache management | Invisible, no commands, manual delete only | ✅ Decided | DR-036 |
| Config integration | Inline assets completely, global by default | ✅ Decided | DR-033 |
| Asset download | `asset_download = true` setting + flag | ✅ Decided | DR-033 |
| Minimal viable set | 30 assets (8 roles, 12 tasks, 6 agents, 2 templates, 2 contexts) | ✅ Decided | DR-031 |

## High-Priority Questions - RESOLVED ✅

1. **Sidecar metadata schema** - ✅ RESOLVED (DR-032)
   - 6 required fields: name, description, tags, sha, created, updated
   - No optional fields (KISS)
   - Category derived from filesystem (no drift)
   - SHA is 40-char git blob hex

2. **GitHub API strategy** - ✅ RESOLVED (DR-034)
   - Tree API for catalog browsing (in-memory cache)
   - raw.githubusercontent.com for downloads (no rate limit!)
   - Contents API as fallback
   - Recommend GITHUB_TOKEN env var (5000/hr vs 60/hr)

3. **Asset resolution algorithm** - ✅ RESOLVED (DR-033)
   - Priority: local → global → cache → GitHub
   - `asset_download = true` setting (default)
   - `--asset-download[=bool]` and `--local` flags
   - Auto-add to global unless `--local`

4. **Cache management** - ✅ RESOLVED (DR-036)
   - No `start cache` command - cache is invisible
   - Manual deletion only: `rm -rf ~/.config/start/assets`
   - Files are tiny (< 1MB for hundreds of assets)
   - No cleanup needed

5. **Config integration** - ✅ RESOLVED (DR-031, DR-033)
   - Multi-file config: config.toml, tasks.toml, agents.toml, contexts.toml
   - Downloaded assets inlined completely (no source tracking)
   - Cache updated via `start update`, user config never auto-updated
   - Global by default, `--local` to override

## Open Questions

### Medium Priority (Resolved for v1, Future Enhancement)

1. **TUI library choice** - ✅ RESOLVED (DR-035)
   - v1: Numbered selection (no dependencies)
   - Future: TUI library if users request it
   - Works everywhere (SSH, containers, CI)

2. **Asset validation** - ✅ RESOLVED (DR-032)
   - TOML parsing validates structure
   - Required field checks in Go code
   - SHA format validation (40-char hex)
   - No schema file (keep simple)

3. **Update notification** - ✅ RESOLVED (DR-037)
   - Manual `start update` only (DR-025 compliant)
   - No automatic checks
   - No background processes
   - User explicitly opts in

4. **Error recovery** - ✅ RESOLVED (DR-037)
   - Partial failures reported clearly
   - User can retry `start update`
   - Cache is disposable (can delete and re-download)

5. **Asset dependencies** - DEFERRED (not in v1)
    - Add if users request it
    - Would use `dependencies` field in metadata
    - Auto-install would require dependency resolution

### Low Priority (Nice to Have)

1. **Community assets** - How can users contribute?
    - PR process to main repo?
    - User-hosted repos?
    - Namespacing: `user/repo/path`

2. **Search functionality** - Full-text search across assets?
    - Search by description, tags, name
    - Fuzzy matching
    - Integration with TUI

3. **Analytics** - Track popular assets (privacy-respecting)?
    - Opt-in telemetry
    - Help prioritize development
    - No PII collection

4. **Asset preview** - Preview asset before downloading?
    - Show metadata (description, tags)
    - Show actual content
    - Diff between versions

5. **Bulk operations** - Install multiple assets at once?
    - Install entire category
    - Install from list/file
    - Workspace templates

## Design Records - COMPLETED ✅

### New DRs (All Written)

1. ✅ **DR-031: Catalog-Based Asset Architecture**
   - Overall system design and core principles
   - Multi-file configuration structure
   - Lazy loading model and GitHub as database
   - Supersedes DR-014, DR-015, DR-016; Updates DR-019, DR-023

2. ✅ **DR-032: Asset Metadata Schema**
   - Sidecar `.meta.toml` format with 6 required fields
   - Category derived from filesystem (no drift)
   - SHA-based versioning, no semver needed
   - Validation rules and Go struct

3. ✅ **DR-033: Asset Resolution Algorithm**
   - Priority order: local → global → cache → GitHub
   - `asset_download` setting and `--asset-download` flag
   - `--local` flag for local config (global by default)
   - Error handling for each failure case

4. ✅ **DR-034: GitHub Catalog API Strategy**
   - Tree API for browsing (in-memory cache)
   - raw.githubusercontent.com for downloads (no rate limit)
   - Contents API as fallback
   - Rate limiting and authentication strategy

5. ✅ **DR-035: Interactive Asset Browsing**
   - Numbered selection for v1 (no dependencies)
   - Category navigation workflow
   - Future TUI library consideration
   - Non-interactive mode for scripts

6. ✅ **DR-036: Cache Management**
   - Cache is invisible (no user commands)
   - Manual deletion only: `rm -rf ~/.config/start/assets`
   - No size limits or cleanup policies
   - Cache structure and operations

7. ✅ **DR-037: Asset Update Mechanism**
   - SHA-based update detection via Tree API
   - Manual `start update` only (DR-025 compliant)
   - Updates cache only, never user config
   - Per-asset updates with clear reporting

### Updates to Existing DRs (All Completed)

1. ✅ **DR-014: GitHub Tree API** → Status: Superseded by DR-031
2. ✅ **DR-015: Atomic Updates** → Status: Superseded by DR-031
3. ✅ **DR-016: Asset Discovery** → Status: Superseded by DR-031
4. ✅ **DR-019: Task Loading** → Status: Updated by DR-031
5. ✅ **DR-023: Staleness Checking** → Status: Superseded by DR-031

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

### Command Specs to Update

- [x] `start-config-agent.md` - Add interactive browsing flow
- [x] `start-config-task.md` - Add interactive browsing flow
- [x] `start-config-role.md` - Add interactive browsing flow (NEW)
- [ ] `start-task.md` - Add lazy loading behavior
- [ ] `start-update.md` - Update to per-asset model

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
asset_download = true                           # Auto-download from GitHub if not found
asset_path = "~/.config/start/assets"           # Cache location
asset_repo = "grantcarthew/start"               # GitHub repository
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
├── contexts/
│   └── dev/
│       ├── go.toml
│       └── go.meta.toml
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

### Contexts (2)

- **dev/** (2): go, python

**Total: 30 assets** - Enough to impress, not overwhelming

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

**Goal:** Ship with 30 high-quality assets

- [ ] Write 8 role prompts (Markdown)
- [ ] Write 12 task definitions (toml)
- [ ] Write 6 agent configs (toml)
- [ ] Write 2 project templates (toml)
- [ ] Write 2 context definitions (toml)
- [ ] Create all metadata files
- [ ] Validate all assets
- [ ] Test each asset end-to-end

**Validation:** All 30 assets work perfectly

### Phase 6: Polish (Quality)

**Goal:** Production-ready catalog system

- [ ] Search/filter functionality
- [ ] Asset validation and security
- [ ] Error recovery and retry logic
- [ ] Performance optimization
- [ ] Documentation and examples
- [ ] GitHub token setup guide

**Validation:** Smooth, reliable, professional experience

## Success Criteria

Catalog design is complete when:

- [x] All high-priority questions resolved (5 questions) - **DONE 2025-01-10**
- [x] All new DRs written (DR-031 through DR-037) - **DONE 2025-01-10**
- [x] Existing DRs updated with notes - **DONE 2025-01-10**
- [ ] Command specs updated for new behavior
- [x] Configuration schema updated - **DONE (multi-file config)**
- [x] Cache structure defined - **DONE (DR-036)**
- [x] Asset metadata schema finalized - **DONE (DR-032)**
- [x] Resolution algorithm specified - **DONE (DR-033)**
- [x] GitHub API strategy documented - **DONE (DR-034)**
- [x] Minimal viable asset set defined (30 assets) - **DONE**

Implementation is complete when:

- [ ] All 6 implementation phases done
- [ ] 30 assets created and tested
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

### Catalog Design Phase Complete (2025-01-10)

**Session 1: Catalog Architecture Brainstorming**

- Working on Task 22 (out-of-box assets) from CLI design phase
- Realized bulk download model was wrong
- Brainstormed catalog-driven architecture
- Defined minimal viable asset set (30 assets)
- Made initial architectural decisions

**Session 2: Design Resolution and Documentation**

- ✅ Resolved all 5 high-priority questions interactively
- ✅ Wrote 7 new Design Records (DR-031 through DR-037)
- ✅ Updated 5 existing DRs with catalog notes
- ✅ Tested GitHub API strategy (confirmed working)
- ✅ Documented all decisions in detail

**Key Final Decisions:**

1. ✅ Metadata: 6 required fields, category from filesystem
2. ✅ Config: Multi-file (config, tasks, agents, contexts)
3. ✅ API: Tree API + raw.githubusercontent.com (no rate limit)
4. ✅ Cache: Invisible, no commands, manual delete only
5. ✅ Resolution: local → global → cache → GitHub
6. ✅ Download: `asset_download = true` + `--asset-download` flag
7. ✅ Integration: Inline assets, global by default, `--local` to override
8. ✅ Updates: Manual only, SHA-based, cache-only (never touch config)
9. ✅ Browsing: Numbered selection (v1), TUI future

**Documents Created:**

- `docs/ideas/catalog-based-assets.md` - Complete vision
- `docs/design/decisions/dr-031-catalog-based-assets.md` - Overall architecture
- `docs/design/decisions/dr-032-asset-metadata-schema.md` - Metadata format
- `docs/design/decisions/dr-033-asset-resolution-algorithm.md` - Resolution logic
- `docs/design/decisions/dr-034-github-catalog-api.md` - GitHub API strategy
- `docs/design/decisions/dr-035-interactive-browsing.md` - Catalog browsing UX
- `docs/design/decisions/dr-036-cache-management.md` - Cache behavior
- `docs/design/decisions/dr-037-asset-updates.md` - Update mechanism

**Next Steps:**

1. Update command specs for new behavior
2. Begin Phase 1 implementation (Core Catalog Infrastructure)
3. Create minimal viable asset set (28 assets)
4. Test end-to-end workflows

## Notes

- This redesign does not invalidate the CLI design (PROJECT-cli-design.md)
- Commands remain the same, behavior changes
- Configuration structure remains the same, asset handling changes
- This is an enhancement, not a breaking change
- Users can still define everything manually (no GitHub required)
- GitHub token recommended but not required (graceful degradation)
