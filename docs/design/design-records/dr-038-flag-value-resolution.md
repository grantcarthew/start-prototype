# DR-038: Flag Value Resolution and Prefix Matching

- Date: 2025-01-11
- Status: Accepted
- Category: CLI Design

## Problem

CLI flags (--agent, --role, --task, --model) need a resolution strategy. The approach must address:

- Exact vs prefix matching (require full names or allow shortcuts)
- Ambiguity handling (multiple matches for same prefix)
- Interactive vs non-interactive environments (TTY vs scripts)
- Source priority (local, global, cache, GitHub)
- Network dependency (when to query GitHub catalog)
- Model passthrough (allow unknown model names for new models)
- Performance (minimize network calls)
- Error messages (clear feedback for users)
- Short-circuit evaluation (avoid unnecessary searches)

## Decision

Implement two-phase resolution (exact → prefix) with short-circuit evaluation across sources, ambiguity detection, and interactive selection in TTY environments.

Phase 1: Exact match across all sources
- Check local config (.start/)
- Check global config (~/.config/start/)
- Check cache (~/.config/start/assets/)
- Check GitHub catalog (lazy fetch if exact match)
- If found: use immediately
- If not found: proceed to Phase 2

Phase 2: Prefix match with short-circuit
- Check local config for prefix matches
  - Single match: use it
  - Multiple matches: handle ambiguity
  - No matches: continue to global
- Check global config for prefix matches
  - Single match: use it
  - Multiple matches: handle ambiguity
  - No matches: continue to cache
- Check cache for prefix matches
  - Single match: use it
  - Multiple matches: handle ambiguity
  - No matches: continue to GitHub
- Check GitHub catalog for prefix matches
  - Single match: lazy fetch and use
  - Multiple matches: handle ambiguity
  - No matches: error "not found"

Short-circuit: Stop at first source with matches (single or multiple), only proceed to next if zero matches.

Ambiguity handling:

Interactive (TTY detected):
- Show numbered selection menu
- User selects by number
- Asset loaded

Non-interactive (piped, scripted, or --non-interactive flag):
- Error with list of matches
- Require exact name or longer prefix
- Exit with error code

Flag-specific behaviors:

--model: Resolution exact → prefix → passthrough
- Sources: Agent configuration file only (in-memory)
- Passthrough: If no match, pass literal value to underlying CLI
- Allows using brand-new models before updating config

--agent, --role, --task: Resolution exact → prefix (no passthrough)
- Sources: local → global → cache → GitHub catalog
- No passthrough: Must be valid asset file
- Invalid name = error

## Why

Prefix matching improves usability:

- Quick shortcuts (--agent anth instead of --agent anthropic)
- Less typing for common operations
- Familiar pattern (like command line tools, git, docker)
- Reduces errors (less typing = fewer typos)
- Discoverable (users can experiment with prefixes)

Short-circuit evaluation optimizes performance:

- Stops at first source with matches (avoids unnecessary checks)
- Local matches fastest (no network)
- Global next (still no network)
- Cache before GitHub (avoid network if possible)
- GitHub only when needed (minimal network calls)
- Typical case: zero network calls

Ambiguity detection prevents errors:

- Catches prefix conflicts early (before execution)
- Interactive selection when multiple matches (TTY)
- Explicit errors in scripts (prevents silent failures)
- Clear error messages (user knows what to fix)
- Safe for automation (non-TTY errors out)

Two-phase resolution (exact → prefix) provides predictability:

- Exact match always wins (no ambiguity)
- Prefix only when no exact match (clear fallback)
- Consistent behavior (same algorithm for all flags)
- Simple mental model (try exact first, then prefix)

Model passthrough enables flexibility:

- Use new models immediately (no config update needed)
- Temporary model names (experimentation)
- Forward compatibility (new models released)
- Still validates known models (catches typos in configured names)

Interactive selection improves discoverability:

- Shows all matches (user learns what's available)
- Source location visible (local vs cache vs GitHub)
- User chooses explicitly (no guessing)
- Teaches users about ambiguity (encourages longer prefixes)

## Trade-offs

Accept:

- Network dependency for uncached assets (prefix matching may require GitHub catalog query, but only if not found locally/globally/cache, typically <200ms)
- Alphabetical bias in selection (interactive list shows alphabetical order not priority order, show source location to help)
- Complexity over exact-only (more complex implementation than simple exact matching, but better UX worth it)
- Model passthrough ambiguity (unknown if model name typo or intentional new model, warning message helps)
- No fuzzy matching (can't suggest "did you mean?", add if users request)

Gain:

- Usability shortcuts (quick prefixes, less typing, familiar pattern, reduces errors)
- Performance optimization (short-circuit evaluation, minimal network calls, zero API calls typical case)
- Safe ambiguity handling (interactive selection in TTY, explicit errors in scripts, clear messages)
- Predictable resolution (exact first then prefix, consistent across flags, simple mental model)
- Model flexibility (passthrough for new models, temporary experimentation, forward compatible)
- Network-aware (lazy fetch, only when needed, respects user's network context)

## Alternatives

Exact match only (no prefix matching):

Example: Require full asset names always
```bash
start --agent anthropic  # Works
start --agent anth       # Error: not found
```

Pros:
- Simple implementation (just string comparison)
- No ambiguity (exact match or error)
- Predictable (always know what you'll get)
- Fast (no prefix search needed)

Cons:
- More typing (full names required)
- Less usable (common shortcuts impossible)
- Unfamiliar pattern (most CLIs support prefixes)
- User frustration (typing long names repeatedly)

Rejected: Prefix matching significantly improves usability. Ambiguity detection handles conflicts. Worth the complexity.

Prefix match with no exact-first phase:

Example: Always do prefix matching, skip exact matching phase
```bash
start --agent anthropic  # Matches as prefix (still works)
start --agent anth       # Matches as prefix
```

Pros:
- Simpler (one-phase instead of two)
- Still supports shortcuts
- Less code (no exact phase)

Cons:
- Ambiguity more common (exact matches become ambiguous if prefix of another)
- Slower (always search all items even if exact match exists)
- Counterintuitive (exact match should be instant)
- No performance benefit (can't short-circuit on exact)

Rejected: Two-phase is better. Exact match should be instant and unambiguous. Prefix is fallback.

Always error on ambiguity (no interactive selection):

Example: Ambiguous prefix always errors, even in TTY
```bash
start --agent a  # Error (even in interactive terminal)
Error: Ambiguous prefix 'a' matches multiple agents
```

Pros:
- Consistent behavior (same in TTY and scripts)
- Forces explicit names (users learn full names)
- Simpler implementation (no TTY detection, no selection UI)

Cons:
- Poor UX in interactive use (forces re-typing instead of selecting)
- Frustrating (user knows what they want, can't choose)
- Less discoverable (doesn't show what's available)
- Inconsistent with modern CLI patterns (most tools offer interactive selection)

Rejected: Interactive selection better UX. TTY detection standard practice. Worth the complexity.

No model passthrough:

Example: --model must match configured model, no passthrough
```bash
start --model gpt-5-new  # Error: not found (even if valid for underlying CLI)
```

Pros:
- Catches typos (all model names validated)
- Consistent (same behavior as --agent)
- Simple (no special case)

Cons:
- Prevents using new models (must update config first)
- Breaks experimentation (can't try temporary model names)
- Forward compatibility issues (new models released, can't use immediately)
- User frustration (unnecessary config updates for temporary use)

Rejected: Model passthrough important for flexibility. Warning message prevents silent typos. New models common use case.

Cache GitHub catalog listings:

Example: Cache Tree API results for N minutes
```go
// Cache catalog listing in memory
catalogCache = fetchGitHubCatalog()  // Cache for 5 minutes
```

Pros:
- Fewer network calls (reuse catalog across invocations)
- Faster repeated operations (no refetch)
- Works offline after first fetch (within TTL)

Cons:
- Stale catalog possible (user sees old asset list)
- Cache invalidation complexity (when to refresh?)
- Confusing (new assets added, not visible until cache expires)
- More state to manage (TTL, expiry, refresh logic)

Rejected: Always-fresh catalog preferred. Network call acceptable (<200ms). Staleness worse than latency.

## Structure

Resolution algorithm:

Phase 1: Exact match
1. Check local config (.start/{type}.toml) for exact name match
2. Check global config (~/.config/start/{type}.toml) for exact name match
3. Check cache (~/.config/start/assets/{type}/*/*.toml) for exact name match
4. Check GitHub catalog (query index.csv) for exact name match
5. If found anywhere: use immediately (lazy fetch from GitHub if needed)
6. If not found: proceed to Phase 2

Phase 2: Prefix match (short-circuit)
1. Search local config for prefix matches
   - If single match: use it
   - If multiple matches: handle ambiguity (see below)
   - If zero matches: continue to step 2

2. Search global config for prefix matches
   - If single match: use it
   - If multiple matches: handle ambiguity
   - If zero matches: continue to step 3

3. Search cache for prefix matches
   - If single match: use it
   - If multiple matches: handle ambiguity
   - If zero matches: continue to step 4

4. Query GitHub catalog for prefix matches
   - Download index.csv
   - Filter by prefix
   - If single match: lazy fetch and use
   - If multiple matches: handle ambiguity
   - If zero matches: error "not found"

Ambiguity handling:

Interactive mode (TTY detected and --non-interactive NOT set):
1. Show numbered list of matches with source locations
2. Prompt user for selection
3. Load selected asset
4. Continue execution

Non-interactive mode (piped, scripted, or --non-interactive flag):
1. Show error message with list of matches
2. Suggest using exact name or longer prefix
3. Exit with error code 1

TTY detection:
- Check if stdout is terminal: os.Stdout.Stat() & os.ModeCharDevice
- Check --non-interactive flag is NOT set
- Both must be true for interactive mode

Flag-specific resolution:

--model (from agent config):
- Resolution: exact → prefix → passthrough
- Sources: Agent configuration [models] section only
- Passthrough behavior:
  - If no match found: warn and pass literal value to CLI
  - Warning: "Model '{name}' not in config, passing through to CLI"
  - Allows new models without config update

--agent, --role, --task (asset files):
- Resolution: exact → prefix (no passthrough)
- Sources: local → global → cache → GitHub (priority order)
- No passthrough: Must be valid asset
- Error if not found: "agent '{name}' not found"

Lazy fetch behavior:

When match found in GitHub catalog:
1. Download asset via raw.githubusercontent.com
2. Save to cache (~/.config/start/assets/{type}/{category}/)
3. Save metadata (.meta.toml)
4. Load and use immediately
5. Available offline on subsequent runs

GitHub catalog query:

For prefix matching against GitHub:
- Download index.csv from raw.githubusercontent.com
- Parse CSV into memory
- Filter by prefix match on name field
- Return matching asset names
- Zero API rate limit impact (raw URL)

## Usage Examples

Exact match (instant):

```bash
$ start --agent anthropic
[loads anthropic agent immediately]
```

Unique prefix (instant if cached):

```bash
$ start --agent anth
[loads anthropic agent via prefix match]
```

Ambiguous prefix - interactive (TTY):

```bash
$ start --agent a

Multiple agents match 'a':
  1. anthropic (local)
  2. azure-openai (cache)
  3. aws-bedrock (GitHub)

Select agent [1-3]: 2
[loads azure-openai agent]
```

Ambiguous prefix - non-interactive (script):

```bash
$ start --agent a --non-interactive

Error: Ambiguous prefix 'a' matches multiple agents:
  - anthropic (local)
  - azure-openai (cache)
  - aws-bedrock (GitHub)

Use exact name or longer prefix.

$ echo $?
1
```

Not found:

```bash
$ start --agent fake

Error: agent 'fake' not found

Available agents:
  - anthropic
  - azure-openai
  - aws-bedrock

Try: start assets browse
```

Model exact match:

```bash
$ start --model claude-sonnet-4
[uses configured claude-sonnet-4]
```

Model prefix match (ambiguous):

```bash
$ start --model claude

Error: Ambiguous prefix 'claude' matches:
  - claude-opus
  - claude-sonnet-4

Use exact name or longer prefix.
```

Model passthrough (not found):

```bash
$ start --model gpt-5-experimental

⚠ Model 'gpt-5-experimental' not in config, passing through to CLI
[passes literal string to underlying tool]
```

Lazy fetch from GitHub:

```bash
$ start --agent aws-bedrock

Agent 'aws-bedrock' not in cache.
Found in catalog: agents/cloud/aws-bedrock

Downloading...
✓ Cached to ~/.config/start/assets/agents/cloud/
✓ Loaded aws-bedrock

[executes with aws-bedrock agent]
```

Role prefix matching:

```bash
$ start --role go-expert
[exact match, loads immediately]

$ start --role go
[prefix match, loads go-expert]

$ start --role code
Multiple roles match 'code':
  1. code-reviewer
  2. code-auditor

Select role [1-2]: _
```

Task prefix matching:

```bash
$ start task pre-commit-review
[exact match]

$ start task pre
[prefix matches pre-commit-review]

$ start task code
Multiple tasks match 'code':
  1. code-quality/find-bugs
  2. code-quality/quick-wins
  3. git-workflow/code-review

Select task [1-3]: _
```

## Updates

- 2025-01-17: Initial version aligned with schema; removed implementation code, Related Decisions, and Future Considerations sections; fixed paths to use local (.start/), global (~/.config/start/), and cache (~/.config/start/assets/)
