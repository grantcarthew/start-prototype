# DR-030: Prefix Matching for Commands

- Date: 2025-01-10
- Status: Accepted
- Category: CLI Design

## Problem

CLI usability can be improved by allowing shorter command input. The design must address:

- Command typing speed (full commands are verbose for frequent use)
- User experience (balance between brevity and clarity)
- Ambiguity handling (what happens when prefix matches multiple commands)
- Script stability (automated scripts should not break)
- Implementation complexity (build custom vs leverage framework)
- Command naming constraints (avoiding conflicting prefixes)
- Breaking changes (adding new commands may conflict with shortcuts)
- Documentation clarity (examples should be understandable)

## Decision

Enable Cobra's built-in prefix matching globally via cobra.EnablePrefixMatching = true. Users can type unambiguous prefixes of commands at all levels instead of full command names.

How it works:

- User types partial command (e.g., con)
- Cobra checks all subcommands for prefix match
- If exactly one matches: uses that command
- If zero or multiple match: returns error with suggestions
- Works at all command levels automatically

Examples:

Top-level commands:
- start d → start doctor
- start con → start config
- start t mytask → start task mytask

Nested commands:
- start con ag l → start config agent list
- start ass a role1 → start assets add role1
- start con r e myRole → start config role edit myRole

Ambiguous prefixes fail with helpful error.

## Why

Faster typing for power users:

- Significantly reduces keystrokes for frequent commands
- Allows natural development of personal shortcuts
- Less cognitive load during interactive sessions
- Muscle memory develops for common prefix patterns
- Better UX for daily usage

Forgiving and discoverable:

- Typos often still match (sta → start, doe → doctor)
- Progressive disclosure with Tab completion
- Can explore commands by typing prefixes
- Reduces friction for learning the CLI

Zero implementation cost:

- Single line of code: cobra.EnablePrefixMatching = true
- Cobra's findNext() method handles all logic
- Automatic at all command levels
- No custom code to maintain

Industry standard pattern:

- kubectl uses this extensively (get po, desc no)
- Familiar to developer audience
- Expected feature for modern CLIs
- Works well with shell completion

Complements shell completion:

- Tab completion for full commands
- Prefix matching for quick shortcuts
- Together provide flexible UX
- Users can choose their preferred style

## Trade-offs

Accept:

- Ambiguity risk when adding new commands (new command may conflict with existing shortcuts)
- Script fragility (scripts using shortcuts may break if new commands added)
- Reduced explicitness (start con ag l less readable than full command)
- Command naming constraints (must avoid similar prefixes where possible)
- Cobra warns this can be "dangerous" (acknowledge the risk for better UX)
- Documentation complexity (must explain shortcuts vs full commands)

Gain:

- Significantly faster typing for power users (fewer keystrokes, better flow)
- Forgiving user experience (typos often still work, less frustration)
- Zero implementation cost (one line of code, Cobra handles everything)
- Works with shell completion (complementary features, not redundant)
- Consistent everywhere (all commands, all levels, no special cases)
- Familiar pattern (kubectl users expect this, industry standard for complex CLIs)
- Progressive disclosure (helps users explore command structure)

## Alternatives

No prefix matching:

Example: Require full command names always
```bash
start config agent list  # Always required
start con ag l          # Error: unknown command
```

Pros:
- No ambiguity possible (commands always explicit)
- Scripts never break from new commands
- Documentation always clear
- No naming constraints
- Simple and predictable

Cons:
- Verbose for frequent use (lots of typing)
- Poor UX for power users (tedious and slow)
- Less forgiving (typos always fail)
- Missing industry-standard feature
- Users will create shell aliases anyway

Rejected: Poor UX for interactive use. Power users expect prefix matching in modern CLIs. Implementation is trivial.

Selective prefix matching:

Example: Enable only for specific command levels
```go
// Only top-level commands
rootCmd.EnablePrefixMatching = true
// Nested commands require full names
```

Pros:
- Reduces ambiguity scope (fewer places for conflicts)
- Scripts could use shortcuts at top level safely
- More predictable which shortcuts work
- Gradual adoption possible

Cons:
- Inconsistent UX (works sometimes, not others)
- Users confused about where prefixes work
- More complex implementation (per-command config)
- Limits benefits to partial use cases
- Still have same risks at top level

Rejected: Inconsistency is worse than consistent behavior. If we enable prefix matching, enable it everywhere.

Custom alias system:

Example: User-defined command aliases
```toml
[aliases]
d = "doctor"
ca = "config agent"
cal = "config agent list"
```

Pros:
- Users define their own shortcuts
- No ambiguity (explicit mappings)
- Scripts can use aliases safely (defined in config)
- Full control over abbreviations

Cons:
- Requires configuration (not automatic)
- Users must set up aliases manually
- More implementation work
- Aliases not portable (machine-specific)
- Doesn't help with typo forgiveness
- Must maintain alias system

Rejected: More work to implement and configure. Cobra's prefix matching provides better UX with zero configuration.

## Structure

Matching logic:

How Cobra handles prefix matching:

1. User types partial command (e.g., con)
2. Cobra checks all subcommands for prefix match
3. If exactly one matches: uses that command
4. If zero or multiple match: returns error with suggestions
5. Works at all command levels automatically

Ambiguity handling:

Exactly one match:
- Command executes normally
- User sees expected behavior

Multiple matches (ambiguous):
- Error message lists matching commands
- Suggests full command names
- Exit with error code

No matches (unknown command):
- Error message shows similar commands
- Suggests using --help
- Exit with error code

Usage guidelines:

For users in interactive sessions:
- Use shortcuts freely (start con ag l)
- Experiment to find what works
- Develop personal muscle memory

For scripts and CI:
- Always use full commands (start config agent list)
- Prevents breakage when new commands added
- More explicit and maintainable

For documentation and sharing:
- Prefer full commands for clarity
- Can mention shortcuts as optional feature
- Examples should be understandable

For development and command naming:

Avoid similar prefixes:
- Don't add both configure and config
- Think about common single-letter prefixes
- Consider existing shortcuts when naming

Version stability:
- Adding commands is a minor change (may break shortcuts)
- Users should expect new commands to affect shortcuts
- Document command additions in release notes

## Usage Examples

Top-level command shortcuts:

```bash
$ start d
# Executes: start doctor

$ start con
# Executes: start config

$ start t mytask
# Executes: start task mytask

$ start ass
# Executes: start assets
```

Nested command shortcuts:

```bash
$ start con ag l
# Executes: start config agent list

$ start con r e myRole
# Executes: start config role edit myRole

$ start ass a task1
# Executes: start assets add task1
```

Ambiguous prefix error:

```bash
$ start c

Error: command "c" is ambiguous, matches:
  - config
  - completion

Try:
  start --help   # See all commands
```

Exit code: 1

Unknown command with suggestions:

```bash
$ start xyz

Error: unknown command "xyz"

Did you mean one of these?
  - config
  - task

Run 'start --help' for usage.
```

Exit code: 1

Script usage (full commands):

```bash
#!/bin/bash
# CI script - always use full commands

start config validate
start doctor
start task code-review "check security"

# Don't use shortcuts in scripts:
# start con val     # Bad - may break if new command added
# start d           # Bad - ambiguous if "debug" command added
# start t review    # Bad - could break with new commands
```

Interactive usage (shortcuts):

```bash
# Power user session - shortcuts are fine
$ start d
$ start con ag l
$ start t cr "check this"

# These develop naturally with muscle memory
```

Comparison with other CLIs:

kubectl (extensive prefix matching):
```bash
kubectl get po       # kubectl get pods
kubectl desc no      # kubectl describe nodes
kubectl apply -f     # kubectl apply -f (exact)
```

git (manual aliases, no built-in prefix):
```bash
git st    # Error (unless aliased)
git co    # Error (unless aliased)
git br    # Error (unless aliased)
```

docker (short exact commands, no prefix):
```bash
docker ps    # Exact command
docker co    # Error (not a prefix match)
```

Our approach follows kubectl:
- Similar target audience (developers)
- Similar command complexity (nested subcommands)
- Better UX for frequent use
- Industry standard for developer CLIs

## Updates

- 2025-01-17: Initial version aligned with schema
