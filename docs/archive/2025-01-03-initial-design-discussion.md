# Initial Design Discussion - 2025-01-03

**Session Type:** Design Discussion
**Duration:** Extended session
**Participants:** Grant Carthew, Claude (Sonnet 4.5)
**Outcome:** Complete architectural design for `start` tool

## Session Overview

This session established the complete architectural design for converting the bash `start` script into a standalone Go binary. The discussion covered vision, configuration, CLI structure, tasks, and distribution strategy.

## Context

Started with the question of whether to convert the bash `start` script to Go. Initial answer: maybe not worth it for current functionality. But when the goal shifted to distribution and adding features (wizards, tasks), Go became the clear choice.

**Key realization:** Not just converting a script - building a proper tool for distribution with complementary GitHub projects (reference-template, snag, kagi).

## Major Discussion Topics

### 1. Vision Definition

Created `docs/vision.md` to establish:

- The problem: Context setup friction for AI development
- The solution: Context-aware launcher
- The pattern: Optional markdown files (ROLE.md, AGENTS.md, PROJECT.md)
- Non-goals: Not replacing agents, not making API calls, not orchestrating workflows
- Success criteria: Easy install, configure, launch

**Key refinement:** Made "the pattern" configurable rather than prescriptive. Files are examples, not requirements.

### 2. Configuration Design

**Major decision point:** TOML vs YAML

- Explored both formats with examples
- User preference: "I like yaml, but I hate yaml..."
- Landed on TOML: no whitespace sensitivity, good Go support, used by mise

**Configuration structure evolution:**

- Started with single file idea
- Discussed global vs local
- Settled on merge strategy: global (~/.config/start/) + local (./.start/) merge
- Named documents (not arrays) enable override/add patterns

**Clever insight:** Named document sections solve the override problem

```toml
[context.documents.project]  # Global defines this
path = "./PROJECT.md"

[context.documents.project]  # Local can override same name
path = "~/multi-repo/BIG-PROJECT.md"
```

### 3. Agent Configuration

**Evolution of thinking:**

- Started with alpha/beta/gamma aliases (from bash script)
- Realized: actual agent names (claude, gemini, opencode) are clearer
- Decision: Agents are global-only (no per-project agent definitions)
- Rationale: Users have 3-5 agents, manageable in global config

### 4. Tasks Discovery

**Turning point:** Discussing how to support gdr (git-diff-review) with special instructions.

Started with simple task idea, evolved to:

- Tasks have roles (system prompts)
- Tasks have prompt templates
- Tasks can run content_command (e.g., git diff)
- Tasks support {instructions} and {content} placeholders
- Tasks reference named context documents

**Example evolution:**

```bash
# Simple
start task code-review

# With instructions
start task gdr "focus on security"

# Full power
start task gdr --agent gemini "check performance"
```

**Major insight:** Tasks solve the "common workflows" problem without becoming a workflow orchestrator.

### 5. Default Tasks

**Discussion:** Which tasks to ship?

From bash scripts identified:

- code-review, git-diff-review, commit-message, comment-review, doc-review, update-changelog, gitignore

**Decision point:** commit-message and gitignore need file I/O (saving results).

**Resolution:** Ship only interactive review tasks (4 total):

1. code-review (cr)
2. git-diff-review (gdr)
3. comment-tidy (ct) - renamed from comment-review
4. doc-review (dr)

**Rationale:** Stays true to vision - launcher, not orchestrator. Users can add non-interactive tasks themselves.

### 6. Placeholder System

Discussed several times, refined incrementally:

**Global placeholders:**

- {model} - Model name
- {system_prompt} - Role file contents
- {prompt} - Built prompt
- {date} - Timestamp

**Task placeholders:**

- {instructions} - User's extra arguments ("None" if empty)
- {content} - Output from content_command

**Considered but rejected:**

- {env:VAR} - Agents inherit environment naturally
- {home} - Just use ~ instead
- {cwd} - Use --directory flag
- {{double braces}} - Single braces simpler

### 7. CLI Structure

**Framework choice:** Cobra (like kubectl/git)

**Pattern established:**

```bash
start <subcommand> [args] [flags]
```

**Subcommands identified:**

- start (root) - Launch session
- start task <name> - Run predefined task
- start agent add/list/test - Manage agents
- start config show/edit - Manage config
- start init - First-time setup

**Key insight:** Persistent flags in Cobra work everywhere:

```bash
start --agent gemini task gdr "focus on security"
```

### 8. Init Wizard Design

**Evolution:**

- Started discussing GitHub resource fetching
- Considered caching with TTL
- Simplified to: fetch on init, cache indefinitely
- Further simplified to: embed in binary, no network

**Final approach:**

- Assets embedded via go:embed
- Interactive wizard on `start init`
- Multi-select agents from popular list
- Auto-detect context documents
- Backup existing config automatically

**Agent list from research:**

1. claude (Claude Code)
2. gemini (Gemini CLI)
3. aichat (multi-provider)
4. opencode (open-source)
5. codex (OpenAI)
6. aider (coding assistant)

### 9. File Detection and Output

**Output format design:**

```
Starting AI Agent
=================================
Agent: claude (model: claude-sonnet-4-5@20250929)

Context documents:
  ✓ environment     ~/reference/ENVIRONMENT.md
  ✓ index          ~/reference/INDEX.csv
  ✗ agents         ./AGENTS.md (not found)
  ✗ project        ./PROJECT.md (not found)

System prompt: ./ROLE.md

Executing command...
❯ claude --model ... --append-system-prompt '...' '2025-11-03...'
```

**Design principle:** Show what's happening, but missing files aren't errors.

### 10. Distribution Strategy

**Final decision:** Embed everything in binary

- No GitHub API calls
- No network dependency
- Works offline
- `go install` and `brew install` just work
- New release = new assets

## Key Design Principles Established

1. **Zero ceremony:** Just type `start` and go
2. **Configuration over code:** Users customize via TOML, not rebuilding
3. **Launcher, not orchestrator:** Delegates to existing AI tools
4. **Optional everything:** Missing files/config = graceful degradation
5. **Interactive when possible:** Wizards > manual config editing
6. **Visibility:** Show what's being used, what's missing

## Documents Created

1. **docs/vision.md** - Product vision and goals
2. **docs/design-record.md** - 11 design decisions with rationale
3. **docs/task.md** - Task configuration documentation
4. **docs/thoughts.md** - Started by Grant, tracked ideas

## Design Decisions (Summary)

**DR-001:** TOML for configuration
**DR-002:** Global + local config merge
**DR-003:** Named documents (not arrays)
**DR-004:** Agents are global-only
**DR-005:** System prompt separate and optional
**DR-006:** Cobra CLI with subcommands
**DR-007:** Single-brace placeholders
**DR-008:** Working directory path resolution, missing files skipped
**DR-009:** Tasks with {instructions} and {content} placeholders
**DR-010:** Four default interactive review tasks
**DR-011:** Assets embedded in binary

## Interesting Moments

**"I like yaml, but I hate yaml..."** - Led to TOML decision

**Named documents revelation** - Solved override problem elegantly

**"start is a good name"** - Almost changed it, decided it works for "starting tasks"

**Tasks evolution** - Went from simple to powerful without becoming complex

**"Maybe we should build the design decision document with placeholders?"** - Led to updating design-record.md as we went

**Distribution simplification** - Went from GitHub API → cache → embed (simpler each time)

## What Wasn't Decided

By end of session, realized CLI needs detailed specification. Created PROJECT.md to capture:

- All subcommands needing specification
- Open questions per command
- Documentation template
- Design approach and priorities

Next phase: Work through each command systematically.

## Insights for Implementation

1. **Start simple:** Root command first, then build out
2. **Use Cobra well:** Persistent flags, dynamic subcommands for tasks
3. **Embed smartly:** Assets directory mirrors config structure
4. **Validate early:** Config validation before execution
5. **Error well:** Missing tools, bad config = helpful messages
6. **Test paths:** Both interactive and flag-based flows

## Related Projects

- reference-template: Pattern this complements
- snag: Web fetching tool
- kagi: Search tool

All three are Go tools for developer workflows.

## Session Methodology

- One decision at a time (no information overload)
- Examples before abstractions
- Record decisions immediately
- Ask clarifying questions frequently
- Build on existing bash script patterns
- Consider both power users and newcomers

## Next Steps Identified

1. Design CLI commands completely (each in docs/cli/\*.md)
2. Resolve open questions about flags and arguments
3. Create project structure (cmd/, internal/, assets/)
4. Implement core functionality
5. Build init wizard
6. Test with real agents
7. Package for distribution

## Final State

Architecture complete and documented. Ready for detailed CLI specification phase. All major design questions resolved. Implementation can begin once CLI specification is complete.
