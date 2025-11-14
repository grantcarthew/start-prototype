# DR-030: Prefix Matching for Commands

**Date:** 2025-01-10
**Status:** Accepted
**Category:** CLI Design

## Decision

Enable Cobra's built-in prefix matching globally via `cobra.EnablePrefixMatching = true`. Users can type unambiguous prefixes of commands at all levels instead of full command names.

## What This Means

### User Experience

Users can type partial commands as long as they're unambiguous:

**Top-level commands:**
```bash
start d          → start doctor
start u          → start (ambiguous: start-assets-update)
start con        → start config
start t mytask   → start task mytask
```

**Nested commands:**
```bash
start con ag l           → start config agent list
start ass a task1        → start assets add task1
start con r e myRole     → start config role edit myRole
```

**Ambiguous prefixes fail:**
```bash
start config t   → Error: "t" matches both "task" and nothing else currently
                   (or suggests "did you mean: task?")
```

### How Cobra Handles It

**Matching logic:**
1. User types partial command (e.g., `con`)
2. Cobra checks all subcommands for prefix match
3. If **exactly one** matches → uses that command
4. If **zero or multiple** match → returns error/suggestion

**Works at all levels automatically:**
- Root level: `start d`
- First level: `start con ag`
- Second level: `start config ag l`
- Everywhere!

### Implementation

**Single line in main.go:**
```go
func init() {
    cobra.EnablePrefixMatching = true
}
```

That's it. Cobra's `findNext()` method in `command.go` handles the rest.

## Benefits

**For power users:**
- ✅ **Faster typing** - `start con ag l` vs `start config agent list`
- ✅ **Muscle memory** - Develop personal shortcuts naturally
- ✅ **Less cognitive load** - Don't need to remember exact spelling mid-flow

**For all users:**
- ✅ **Forgiving** - Typos often still match (e.g., `sta` → `start`)
- ✅ **Progressive disclosure** - Can type `start c<Tab>` to see completions
- ✅ **Familiar pattern** - kubectl, docker, and other CLIs do this

**For the project:**
- ✅ **Zero implementation cost** - One line of code
- ✅ **Works with completion** - Tab completion + prefix matching = great UX
- ✅ **Consistent everywhere** - All commands, all levels

## Trade-offs Accepted

**Ambiguity risk:**
- ❌ Adding new commands could break existing shortcuts
- ❌ Example: If we add `start configure`, then `start con` becomes ambiguous
- **Mitigation:** Careful command naming, avoid similar prefixes

**Script fragility:**
- ❌ Scripts using shortcuts could break with new commands
- ❌ Example: CI script uses `start u` → we add `start upgrade` → breaks
- **Mitigation:** Document recommendation to use full commands in scripts

**Reduced explicitness:**
- ❌ `start con ag l` is less readable than full command
- ❌ Documentation examples harder to understand
- **Mitigation:** Always show full commands in docs, mention shortcuts as optional

**"Dangerous" per Cobra:**
- ❌ Cobra docs warn "Automatic prefix matching can be a dangerous thing"
- ❌ They know about the ambiguity and breaking change risks
- **Mitigation:** We accept the trade-off for better UX

## Guidelines

### For Users

**Interactive use:**
- Use shortcuts freely: `start con ag l`
- Experiment to find what works

**In scripts/CI:**
- Always use full commands: `start config agent list`
- Prevents breakage when new commands added

**Documentation/sharing:**
- Prefer full commands for clarity
- Can mention shortcuts: "You can also type `start con ag l`"

### For Development

**Command naming:**
- Avoid similar prefixes where possible
- Example: Don't add both `configure` and `config`
- Think about common prefixes: `a`, `c`, `d`, `l`, `s`, `t`, `u`

**Version stability:**
- Adding commands is a **minor** change (may break shortcuts)
- Users should expect new commands → update their shortcuts
- Document command additions in release notes

**Testing:**
- Test both full commands and common prefixes
- Verify ambiguity detection works
- Check error messages are helpful

## Documentation Updates

### README Quick Start

```markdown
**Pro tip:** You can use unambiguous prefixes:
```bash
start d              # start doctor
start con ag l       # start config agent list
start t code-review  # start task code-review
```

**For scripts:** Always use full command names to avoid breakage.
```

### CLI Documentation

Add note to command reference pages:

```markdown
## Prefix Matching

You can type unambiguous prefixes instead of full command names:
- `start con` → `start config`
- `start config ag` → `start config agent`

If your prefix matches multiple commands, `start` will show suggestions.

**Recommendation:** Use full commands in scripts and documentation.
```

### Error Messages

When ambiguous:
```
Error: command "t" is ambiguous, could be:
  - task
  - (if we add more commands starting with 't')

Try:
  start --help   # See all commands
```

When no match:
```
Error: unknown command "xyz"

Did you mean one of these?
  - config
  - task

Run 'start --help' for usage.
```

## Examples in the Wild

**kubectl (extensive use):**
```bash
kubectl get po       → kubectl get pods
kubectl desc no      → kubectl describe nodes
kubectl app version  → kubectl apply version
```

**git (limited use):**
```bash
git st    → Would work if aliased, not built-in
git co    → Not built-in, requires alias
git br    → Not built-in, requires alias
```
Git chose not to do this by default, requires explicit aliases.

**docker (no prefix matching):**
```bash
docker ps    → docker ps (exact command)
docker co    → Error
```
Docker uses short exact commands, not prefix matching.

**Our choice:** Follow kubectl's power-user approach, since:
- Similar target audience (developers)
- Similar command complexity (nested subcommands)
- Better UX for frequent use

## Related Decisions

- [DR-006](./dr-006-cobra-cli.md) - Cobra CLI framework (provides prefix matching)
- [DR-028](./dr-028-shell-completion.md) - Shell completion (complements prefix matching)

## Future Considerations

**If ambiguity becomes a problem:**
- Could disable for specific command levels
- Could add `--no-prefix` flag to force exact matching
- Could make it configurable via settings

**Current stance:** Enable everywhere, monitor for issues. One line to add, one line to remove if needed.
