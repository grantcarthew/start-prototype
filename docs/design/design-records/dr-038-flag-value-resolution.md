# DR-038: Flag Value Resolution and Prefix Matching

**Date:** 2025-01-11
**Status:** Accepted
**Category:** CLI Design

## Decision

Implement intelligent prefix matching for flag values (--agent, --role, --task, --model) with ambiguity detection and interactive resolution in TTY environments.

## Resolution Algorithm

All flags follow a two-phase resolution process:

### Phase 1: Exact Match (All Sources)
Search for exact filename match across all sources in priority order:
1. Local config (`~/.config/start/`)
2. Global config (`/etc/start/` or system-wide)
3. Cache (`~/.cache/start/catalog/`)
4. GitHub catalog (lazy fetch if exact match found)

If exact match found, use it immediately (with lazy fetch from GitHub if needed).

### Phase 2: Prefix Match (Short-Circuit)
If no exact match, search for prefix matches with **short-circuit evaluation**:

```
1. Check local config for prefix matches
   - Single match: use it
   - Multiple matches: handle ambiguity (see below)
   - No matches: continue to step 2

2. Check global config for prefix matches
   - Single match: use it
   - Multiple matches: handle ambiguity
   - No matches: continue to step 3

3. Check cache for prefix matches
   - Single match: use it
   - Multiple matches: handle ambiguity
   - No matches: continue to step 4

4. Fetch GitHub Tree API for catalog/[type]/ directory
   - Parse filenames (remove .toml extension)
   - Filter by prefix
   - Single match: lazy fetch and use
   - Multiple matches: handle ambiguity
   - No matches: error "not found"
```

**Key:** Stop at first source that has matches (single or multiple). Only proceed to next source if current source has zero matches.

## Ambiguity Handling

When multiple assets match a prefix, behavior depends on environment:

### Interactive (TTY detected)
Show numbered selection menu:

```bash
$ start --agent a

Multiple agents match 'a':
  1. anthropic (local)
  2. azure-openai (cache)

Select agent [1-2]: _
```

User selects by number, asset is loaded.

### Non-Interactive (piped, scripted, or --non-interactive flag)
Error with list of matches:

```bash
$ start --agent a

Error: Ambiguous prefix 'a' matches multiple agents:
  - anthropic (local)
  - azure-openai (cache)

Use exact name or longer prefix.
```

Script exits with error code, forcing explicit specification.

## Flag-Specific Behaviors

### --model (from agent config)

**Resolution:** exact → prefix → passthrough

**Sources:** Agent configuration file only (in-memory)

**Passthrough:** If no exact or prefix match, pass the literal value to the underlying CLI tool. This allows using brand-new models before updating config.

**Example:**
```bash
start --model gpt         → Matches "gpt-4o" (prefix)
start --model gpt-5-new   → Passthrough (no match, sends to CLI)
start --model claude      → Ambiguous: "claude-opus", "claude-sonnet"
```

### --agent (asset file)

**Resolution:** exact → prefix (no passthrough)

**Sources:** local → global → cache → GitHub catalog

**No passthrough:** Must be valid asset file. Invalid agent name = error.

**Example:**
```bash
start --agent anthropic   → Exact match, lazy fetch if needed
start --agent anth        → Prefix match "anthropic"
start --agent a           → Ambiguous or interactive selection
start --agent fake        → Error: "agent 'fake' not found"
```

### --role (asset file)

**Resolution:** exact → prefix (no passthrough)

**Sources:** local → global → cache → GitHub catalog

**Behavior:** Identical to --agent

**Example:**
```bash
start --role go-expert       → Exact match
start --role go              → Prefix match "go-expert"
start --role code            → Ambiguous: "code-reviewer", "code-auditor"
```

### --task (asset file)

**Resolution:** exact → prefix (no passthrough)

**Sources:** local → global → cache → GitHub catalog

**Behavior:** Identical to --agent and --role

**Example:**
```bash
start task pre-commit-review   → Exact match
start task pre                 → Prefix match "pre-commit-review"
start task code                → Ambiguous: multiple code-* tasks
```

## Implementation Details

### TTY Detection

```go
func isTTY() bool {
    fileInfo, _ := os.Stdout.Stat()
    return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
```

Check both:
- `isTTY()` returns true
- `--non-interactive` flag is NOT set

If either condition false, use non-interactive error behavior.

### GitHub Tree API Call

For prefix matching against uncached GitHub assets:

```go
// Only called if no matches found in local/global/cache
func fetchGitHubAssetNames(assetType string) ([]string, error) {
    // GET https://api.github.com/repos/{org}/{repo}/git/trees/{branch}?recursive=false
    // Parse response, filter by catalog/{assetType}/*.toml
    // Return list of asset names (without .toml extension)
}
```

**When called:**
- Phase 2, step 4 (prefix matching)
- Only if local/global/cache had zero matches
- Single API call per flag value resolution
- Results not cached (stateless per invocation)

### Lazy Fetch on Match

When exact or prefix match found in GitHub:

```go
// Fetch the specific asset file
func fetchGitHubAsset(assetType, name string) ([]byte, error) {
    // GET https://raw.githubusercontent.com/{org}/{repo}/{branch}/catalog/{assetType}/{name}.toml
    // Save to cache: ~/.cache/start/catalog/{assetType}/{name}.toml
    // Return content
}
```

## Performance Characteristics

**Best case (cached):** No network calls
- Exact match in local/global/cache
- Prefix match in local

**Typical case:** One Tree API call
- Prefix match requires checking GitHub catalog
- ~100-200ms for Tree API response

**Worst case:** Tree API + raw file fetch
- Prefix match found only in GitHub catalog
- Fetch specific asset file
- ~200-400ms total

**Optimization:** Short-circuit evaluation prevents unnecessary checks.

## Benefits

**For users:**
- ✅ Quick shortcuts: `--agent anth` instead of `--agent anthropic`
- ✅ Discover ambiguity before execution
- ✅ Interactive selection when multiple matches (TTY)
- ✅ Safe for scripts (explicit errors in non-TTY)
- ✅ Works with uncached assets (GitHub catalog aware)

**For implementation:**
- ✅ Consistent behavior across all asset-based flags
- ✅ Minimal network overhead (short-circuit + lazy fetch)
- ✅ Clear error messages for debugging
- ✅ Familiar pattern (similar to --model resolution)

## Trade-offs Accepted

**Network dependency:**
- ❌ Prefix matching uncached assets requires GitHub Tree API call
- **Mitigation:** Only called if not found in local/global/cache; results in <200ms typical

**Alphabetical bias:**
- ❌ Interactive selection shows alphabetical order, may not be priority order
- **Mitigation:** Show source location in list (local/cache/GitHub)

**Complexity:**
- ❌ More complex than simple exact-match-only
- **Mitigation:** Better UX worth the implementation complexity

## Examples

### Interactive Session (TTY)

```bash
# Exact match - instant
$ start --agent anthropic
[loads anthropic agent]

# Unique prefix - instant
$ start --agent anth
[loads anthropic agent]

# Ambiguous prefix - interactive
$ start --agent a
Multiple agents match 'a':
  1. anthropic (local)
  2. azure-openai (cache)
  3. aws-bedrock (GitHub)

Select agent [1-3]: 2
[loads azure-openai agent]

# Not found
$ start --agent fake
Error: agent 'fake' not found
```

### Script/CI (Non-TTY)

```bash
# Exact match - works
$ start --agent anthropic --non-interactive
[loads anthropic agent]

# Unique prefix - works
$ start --agent anth --non-interactive
[loads anthropic agent]

# Ambiguous prefix - errors
$ start --agent a --non-interactive
Error: Ambiguous prefix 'a' matches multiple agents:
  - anthropic (local)
  - azure-openai (cache)
  - aws-bedrock (GitHub)

Use exact name or longer prefix.
[exits with code 1]
```

### Model Passthrough

```bash
# Exact match
$ start --model claude-sonnet-4
[uses configured claude-sonnet-4]

# Prefix match
$ start --model claude
Error: Ambiguous prefix 'claude' matches:
  - claude-opus
  - claude-sonnet-4

# Passthrough (no match)
$ start --model gpt-5-experimental
⚠ Model 'gpt-5-experimental' not in config, passing through to CLI
[passes literal string to underlying tool]
```

## Related Decisions

- [DR-030](./dr-030-prefix-matching.md) - Command prefix matching (different: for commands, not flag values)
- [DR-033](./dr-033-asset-resolution-algorithm.md) - Asset resolution algorithm (exact match sources)
- [DR-034](./dr-034-github-catalog-api.md) - GitHub Tree API strategy
- [DR-035](./dr-035-interactive-browsing.md) - Interactive browsing UI (similar selection pattern)

## Future Considerations

**Fuzzy matching:**
- Could implement Levenshtein distance for "did you mean?" suggestions
- Example: `--agent antrhopic` → "Did you mean 'anthropic'?"

**Caching GitHub catalog listing:**
- Could cache Tree API results for N minutes
- Trade-off: Stale catalog vs fewer network calls
- Current: Always fresh, acceptable latency

**Prefix length minimum:**
- Could require minimum prefix length (e.g., 2+ chars)
- Prevents `--agent a` from being too broad
- Current: Any length accepted, let ambiguity detection handle it

**Current stance:** Ship with described behavior, iterate based on user feedback.
