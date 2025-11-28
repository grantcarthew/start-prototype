# DR-044: Shell Quote Escaping for Placeholder Substitution

- Date: 2025-11-28
- Status: Accepted
- Category: Runtime Behavior

## Problem

Agent command templates use placeholder substitution to inject dynamic values (role content, prompts, model names) into shell commands. When these values contain shell metacharacters (quotes, spaces, newlines, `$()`, etc.), the resulting command breaks or behaves unexpectedly.

Current implementation uses naive `strings.ReplaceAll` without escaping:

```toml
[agents.claude]
command = "{bin} --model {model} --append-system-prompt '{role}' '{prompt}'"
```

When `{role} = "You're a Go expert"` is substituted:

```bash
claude --model sonnet --append-system-prompt 'You're a Go expert' 'prompt'
                                                    ^-- SYNTAX ERROR
```

The single quote in "You're" terminates the bash string early, breaking command parsing.

### Forces at Play

**Template Flexibility**: Users need to write templates for any agent with any quoting style.

**Multi-Shell Support**: The `shell` setting can be any interpreter: bash, zsh, fish, python, node, ruby, perl. Each has different quoting rules.

**User Empowerment**: Users should be able to use shell features like `$()` command substitution and `$VAR` expansion in their content if they choose.

**Correctness over Security**: The user controls both template and content. The problem is mechanical correctness (making it work), not security (preventing malicious injection).

**Portability**: Solution must work across POSIX shells (bash, sh, zsh, fish) while allowing programming languages (python, node, ruby) to use their own syntax.

**Clear Errors**: When quotes clash or content is unsafe, provide compiler-style errors with file locations, character positions, and actionable fixes.

## Decision

Implement **context-aware quote escaping** for POSIX shells only:

### For POSIX Shells (bash, sh, zsh, fish)

Parse command template to detect quote context around each placeholder:

1. **Single-quoted context** `'{placeholder}'`:
   - Escape single quotes in value using bash `'\''` pattern
   - Example: `You're` → `You'\''re`
   - Result: `'You'\''re a Go expert'` (valid bash)

2. **Double-quoted context** `"{placeholder}"`:
   - Escape double quotes and backslashes only
   - Leave `$()`, `$VAR`, and backticks unescaped (user feature)
   - Example: `He said "hello"` → `He said \"hello\"`
   - Allows: `"Today is $(date)"` executes command substitution

3. **Unquoted context** `{placeholder}`:
   - Validate value contains only shell-safe characters: `a-z A-Z 0-9 _ - . / :`
   - **Error** if value contains spaces, quotes, `$()`, or other metacharacters
   - Provide detailed error with position, problematic characters, and suggested fixes

### For Programming Languages (python, node, ruby, perl, deno, bun)

**No automatic escaping**:
- User writes language-specific syntax
- User is responsible for correct quoting
- Example Python: `"import subprocess; subprocess.run(['{bin}', '{model}'])"`

### Shell Detection

Classify shells by type:

```go
func IsPOSIXShell(shell string) bool {
    posixShells := []string{"bash", "sh", "zsh", "fish"}
    // Returns true for POSIX shells, false for programming languages
}
```

### Error Reporting

When validation fails (unquoted placeholder with unsafe value), generate detailed error:

```
Error: Unsafe placeholder substitution detected

Agent: claude
Config: ~/.config/start/agents.toml
Field: command
Template: {bin} --model {model} --system {role} --prompt "{prompt}"
                                         ^-----^
Position: characters 28-33

Placeholder: {role}
Quote context: unquoted
Value: "You're a Go expert"
       ^   ^-- contains single quote (') which will break shell parsing

Problematic characters found in value:
  - Position 3: ' (single quote)
  - Position 5: ' (single quote)
  - Position 9: ' ' (space)

This placeholder must be quoted in the template to safely contain this value.

Suggested fixes:
1. Use single quotes in template:  '{role}'
   (will auto-escape single quotes as '\'' in content)

2. Use double quotes in template: "{role}"
   (will auto-escape " and \, allows $(commands) to execute)

3. Ensure value only contains: a-z A-Z 0-9 _ - . / :
   (these are safe in unquoted context)

Current template:
  {bin} --model {model} --system {role} --prompt "{prompt}"

Suggested template:
  {bin} --model {model} --system '{role}' --prompt "{prompt}"
```

### Execution Visibility

Always report shell being used during execution:

```
Executing with shell: bash
❯ claude --model sonnet --append-system-prompt 'You'\''re a Go expert' 'prompt'
```

Or for non-POSIX:

```
Executing with shell: python
❯ python -c "import subprocess; subprocess.run(['claude', 'sonnet'])"
```

## Why

**Correctness is the Goal**: The issue is mechanical (making templates work), not security (preventing attacks). Users control both templates and content, so there's no injection threat—only syntax errors.

**Preserve User Power**: By leaving `$()` and `$VAR` unescaped in double quotes, users can intentionally use shell features. This is a feature, not a bug. Example: `"Today is $(date)"` in a prompt executes the date command.

**POSIX Shell Standardization**: POSIX shells (bash, sh, zsh, fish) share similar quoting rules. A single escaping strategy works across all of them.

**Programming Language Flexibility**: Python, Node, Ruby have completely different syntax. Don't try to parse or escape—let users write native code for their interpreter.

**Shell Detection is Reliable**: The `shell` setting explicitly declares the interpreter. We can confidently classify as POSIX vs programming language.

**Detailed Errors Guide Users**: Compiler-style errors with positions, problematic characters, and suggested fixes teach users the correct pattern without frustration.

**Execution Visibility Prevents Confusion**: Showing which shell is executing and the exact command makes behavior transparent and debuggable.

## Trade-offs

### Accept

**Escaping Complexity**: Must implement quote-aware parser that tracks bash state machine (in quotes, not in quotes, escaped characters). This is non-trivial but solvable.

**POSIX-Only Auto-Escaping**: Programming languages get no help. Users must understand their language's quoting. This is acceptable since they chose to use that language.

**Escaping Edge Cases**: Bash has complex rules (heredocs, ANSI-C quoting, etc.). We only handle the common cases: single quotes, double quotes, unquoted. Rare cases may still break.

**Performance Overhead**: Parsing templates and validating values adds computation. Acceptable since this happens once at startup, not in hot path.

**No Cross-Shell Validation**: If user sets `shell = "python"` but writes bash syntax in template, we can't detect it. User gets runtime error from Python interpreter.

### Gain

**Templates Work Correctly**: Single quotes in content (`You're`) no longer break commands. Auto-escaping makes it just work.

**User Empowerment**: Users can intentionally use `$(date)` or `$HOME` in double-quoted contexts for dynamic behavior.

**Clear Error Messages**: When things fail validation, users get actionable guidance instead of cryptic bash errors.

**Flexibility Preserved**: Users can still use any shell or programming language. No restrictions on template patterns.

**Debuggability**: Showing the exact command being executed makes it easy to understand what's happening and debug issues.

**Future-Proof**: Works with shells that don't exist yet, as long as they follow POSIX conventions or are explicitly classified.

## Alternatives

### Always Escape Everything (Reject All Special Characters)

**Approach**: Reject or escape `$()`, `$VAR`, and all metacharacters, even in double quotes.

**Pros**:
- Simplest implementation
- Most predictable behavior
- No unexpected command execution

**Cons**:
- Removes useful features (can't use `$(date)` in prompts)
- Overly restrictive
- Users lose shell power

**Rejected**: Kills flexibility and user empowerment. The goal is correctness, not sandboxing.

### Temporary Files for All Placeholders

**Approach**: Write all placeholder values to temp files, pass file paths to agent:

```bash
claude --model sonnet --role-file /tmp/role.txt --prompt-file /tmp/prompt.txt
```

**Pros**:
- No escaping needed
- Safe for arbitrary content
- Works across all shells

**Cons**:
- Not all agents support file arguments for prompts
- Extra I/O operations
- Temp file cleanup complexity
- Breaking change for existing configs

**Rejected**: Most agents expect inline arguments, not file paths. This would break existing templates.

### Environment Variables

**Approach**: Pass values via environment variables:

```toml
command = "START_ROLE='{role}' START_PROMPT='{prompt}' {bin} --model {model}"
```

Agent reads from env vars instead of args.

**Pros**:
- No escaping needed for content
- Clean separation of data and command

**Cons**:
- Requires agent support for env vars
- Claude/Gemini don't support this pattern
- Breaking change for all existing configs
- Still requires escaping the env var assignment

**Rejected**: Doesn't match how current agents work. Would require custom wrapper scripts.

### No Auto-Escaping (User Responsibility)

**Approach**: Never escape. Document that users must handle quoting themselves.

**Pros**:
- Simplest implementation (already exists)
- No magic behavior
- Users have full control

**Cons**:
- Error-prone (easy to forget quotes)
- Poor user experience (cryptic bash errors)
- Every user solves the same problem repeatedly

**Rejected**: Poor UX. Users shouldn't need to understand bash quoting rules in depth. Auto-escaping makes it just work.

### Restrict to Predefined Templates

**Approach**: Only allow templates from catalog or hardcoded list. No custom templates.

**Pros**:
- Known-safe patterns
- No escaping needed
- Simplest implementation

**Cons**:
- Kills flexibility (major design goal)
- Doesn't support custom agents
- Against project philosophy

**Rejected**: Template flexibility is a core feature. Users must be able to write custom templates.

## Implementation Notes

### Quote Parser State Machine

Track bash quote state while scanning template:

```
States: UNQUOTED, SINGLE_QUOTED, DOUBLE_QUOTED, ESCAPED
Transitions:
  UNQUOTED + ' → SINGLE_QUOTED
  UNQUOTED + " → DOUBLE_QUOTED
  UNQUOTED + \ → ESCAPED (next char)
  SINGLE_QUOTED + ' → UNQUOTED
  DOUBLE_QUOTED + " → UNQUOTED
  DOUBLE_QUOTED + \ → ESCAPED (next char in double context)
```

When encountering `{placeholder}`, record current state.

### Escaping Functions

```
EscapeSingleQuote(value string) string:
  Replace ' with '\''
  (close quote, escaped quote, open quote)

EscapeDoubleQuote(value string) string:
  Replace " with \"
  Replace \ with \\
  Leave $ ( ) ` unescaped

ValidateUnquoted(value string) error:
  Check value only contains: a-z A-Z 0-9 _ - . / :
  Return error listing problematic characters and positions
```

### Shell Classification

```
IsPOSIXShell(shell string) bool:
  Return true for: bash, sh, zsh, fish
  Return false for: python, python2, python3, node, nodejs, bun, deno, ruby, perl
  Default: false (programming languages safer to assume)
```

### Error Construction

Include in error message:
- Agent name
- Config file path
- Field name (command, command template in role/task)
- Template string
- Character position range (start-end)
- Placeholder name
- Quote context detected
- Value being substituted
- List of problematic characters with positions
- Suggested fixes (3 options: single quote, double quote, safe value)
- Current template (full)
- Suggested template (with fix applied)

## Validation

### Template Parsing

**Valid templates** (POSIX shells):

```toml
command = "{bin} --system '{role}' '{prompt}'"          # Single quotes
command = "{bin} --system \"{role}\" \"{prompt}\""      # Double quotes (escaped in TOML)
command = "{bin} --model {model}"                       # Unquoted (model is safe: "sonnet")
command = "{bin} --system '{role}' --data \"{prompt}\"" # Mixed quotes
```

**Invalid templates** (cause errors during execution):

```toml
# Unquoted with unsafe value (space in prompt)
command = "{bin} {prompt}"  # Error if prompt = "hello world"

# Unquoted with unsafe value (quote in role)
command = "{bin} {role}"    # Error if role = "You're an expert"
```

### Value Validation

**Safe for unquoted**:
- `sonnet`, `haiku`, `opus` (model names)
- `/usr/bin/claude` (paths)
- `192.168.1.1` (IPs)
- `my-model-v2` (kebab-case)

**Unsafe for unquoted** (require quotes):
- `You're a Go expert` (contains quote and spaces)
- `$HOME/path` (contains `$`)
- `echo $(date)` (contains special chars)
- `Line 1\nLine 2` (contains newline)

## Usage Examples

### Claude Agent (Single Quotes)

```toml
[agents.claude]
bin = "claude"
command = "{bin} --model {model} --append-system-prompt '{role}' '{prompt}'"
default_model = "sonnet"
```

**Role**: `"You're a Go expert with 10+ years experience"`

**After escaping**:
```bash
claude --model sonnet --append-system-prompt 'You'\''re a Go expert with 10+ years experience' 'prompt'
```

### Gemini Agent (Double Quotes for Command Substitution)

```toml
[agents.gemini]
bin = "gemini"
command = "{bin} --model {model} --system \"{role}\" --prompt \"{prompt}\""
default_model = "pro"
```

**Prompt**: `"What's the current date? $(date)"`

**After escaping** (minimal, preserves `$(date)`):
```bash
gemini --model pro --system "role" --prompt "What's the current date? $(date)"
```

The `$(date)` executes, injecting the actual date into the prompt.

### Python Shell (No Escaping)

```toml
[settings]
shell = "python"

[agents.custom]
bin = "custom_agent"
command = "import subprocess; subprocess.run(['{bin}', '--model', '{model}', '--prompt', '{prompt}'])"
```

**No auto-escaping**. User writes valid Python. The command is executed as:

```bash
python -c "import subprocess; subprocess.run(['custom_agent', '--model', 'sonnet', '--prompt', 'hello'])"
```

### Error Example (Unquoted with Spaces)

```toml
[agents.broken]
command = "{bin} {prompt}"  # Missing quotes around {prompt}
```

**Prompt**: `"hello world"`

**Error**:
```
Error: Unsafe placeholder substitution detected

Agent: broken
Config: ~/.config/start/agents.toml
Field: command
Template: {bin} {prompt}
                ^------^
Position: characters 6-14

Placeholder: {prompt}
Quote context: unquoted
Value: "hello world"
            ^-- contains space which will break shell parsing

Problematic characters found in value:
  - Position 5: ' ' (space)

This placeholder must be quoted in the template to safely contain this value.

Suggested fixes:
1. Use single quotes in template:  '{prompt}'
2. Use double quotes in template: "{prompt}"

Suggested template:
  {bin} '{prompt}'
```

## Security

This design is **not** about security (preventing injection attacks). It's about **correctness** (making templates work).

### Trust Model

**Assumption**: User controls both template and content. There is no adversary.

- Templates are in user-owned config files
- Role content is from user-written files or commands
- Prompts are from user-typed input or context documents
- Context documents are user-controlled files

**No Injection Threat**: Since the user controls all inputs, there's no one to protect against. If a user writes `$(rm -rf /)` in their role file, that's intentional (or a mistake), not an attack.

### Intentional Features vs Bugs

Leaving `$()` and `$VAR` unescaped in double quotes is **intentional**:

- User can write: `"Today is $(date)"` in a prompt
- This executes `date` command and injects result
- This is a **feature** for power users
- If unwanted, user can use single quotes: `'{prompt}'`

### Protection Scope

What this **does** protect against:
- Syntax errors from unintentional quote conflicts
- Broken commands from spaces in values
- Cryptic bash errors when content has metacharacters

What this **does not** protect against:
- User writing `$(rm -rf /)` in their own content
- User misconfiguring their own templates
- User choosing double quotes and getting command execution

The user is the admin. They can do anything. The goal is making it work, not sandboxing.

## Breaking Changes

None. This is a new feature, not a change to existing behavior.

**Current state**: Naive replacement breaks on quotes/spaces.

**New state**: Auto-escaping fixes breakage, adds validation errors.

Existing templates that work will continue to work. Existing templates that are broken (contain unquoted placeholders with unsafe values) will now produce helpful errors instead of cryptic bash failures.

## Updates

- 2025-11-28: Initial design
