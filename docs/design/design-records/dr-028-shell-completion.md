# DR-028: Shell Completion Support

- Date: 2025-01-07
- Status: Accepted
- Category: CLI Design

## Problem

Modern CLI tools provide shell completion for improved usability. The CLI needs a completion strategy that addresses:

- Shell support (which shells to support and why?)
- Installation methods (how do users enable completion?)
- Completion scope (what gets completed: commands, flags, dynamic values?)
- Implementation complexity (leverage framework or build custom?)
- Dynamic completions (complete agent names, task names from config)
- Cross-platform paths (different completion locations per OS and shell)
- User experience (manual vs auto-install, system-wide vs user-only)
- Maintenance burden (supporting multiple shells and completion types)

## Decision

Support shell completion for bash, zsh, and fish using Cobra's built-in completion system. Provide both manual output and auto-install commands. Complete commands, flags, agent names, task names, and scope arguments.

Supported shells:

- bash (most common on Linux)
- zsh (macOS default since Catalina, popular on Linux)
- fish (growing popularity, different syntax)
- Not PowerShell (Windows-focused, not planning Windows support)

Installation patterns:

Manual output (print to stdout):
- start completion bash
- start completion zsh
- start completion fish
- User redirects to file and sources in shell config

Auto-install (convenience):
- start completion install bash
- start completion install zsh
- start completion install fish
- Automatically installs to standard location for shell
- Detects OS (macOS vs Linux) for path selection
- User directory by default (no sudo required)
- Optional --system flag for system-wide install
- Optional --path flag for custom location

Completion tiers:

Tier 1 - Static (free from Cobra):
- Commands and subcommands
- Flags and short flags
- No custom code required

Tier 2 - Dynamic (high-value, implement):
- Agent names for --agent flag (from config)
- Task names for start task command (from config and catalog)
- Scope arguments (global, local, merged)
- Custom ValidArgsFunction implementations

Tier 3 - Skip for v1:
- Model names for --model flag (requires parsing agent's model table)
- Context names, role names
- File path completion (shells handle this already)

## Why

Shell completion improves usability:

- Faster command entry (Tab completion vs typing full names)
- Discoverability (users discover subcommands and flags via Tab)
- Fewer typos (completion prevents command and flag mistakes)
- Professional expectation (users expect modern CLIs to support completion)
- Learning aid (seeing available options helps users learn the CLI)

Cobra provides free implementation:

- Built-in completion generation for bash, zsh, fish
- ValidArgsFunction hooks for dynamic completions
- Standard patterns that work across shells
- Less code to write and maintain
- Proven, tested completion logic

Dynamic completions add value:

- Agent names from config (users Tab through configured agents)
- Task names from config (shows all available tasks with aliases)
- Scope values (global, local, merged for various commands)
- Context-aware (reads current config state)
- Practical for daily usage

Auto-install reduces friction:

- One command to enable completion (vs multi-step manual process)
- Detects correct paths for shell and OS
- Creates directories if needed
- Shows reload instructions
- Lower barrier to adoption

Three shells cover most users:

- bash, zsh, fish cover >95% of Unix/macOS users
- All three provided free by Cobra
- PowerShell requires Windows support not planned
- Sufficient coverage without excessive maintenance

## Trade-offs

Accept:

- Requires one-time setup by user (completion not enabled by default)
- Different installation steps per shell (auto-install abstracts this)
- Doesn't complete model names (less frequently used, Tier 3)
- No PowerShell support (Windows not planned, niche for Unix CLI)
- Dynamic completions slower than static (config loading overhead acceptable)

Gain:

- Better UX (faster command entry, fewer mistakes, discoverability)
- Professional polish (expected feature for modern CLI tools)
- Free implementation (Cobra provides heavy lifting)
- Dynamic agent and task completion (context-aware from config)
- Easy installation (auto-install removes complexity)
- Standard patterns (follows shell completion conventions)
- Maintainable (Cobra handles cross-shell differences)

## Alternatives

No shell completion:

Pros:
- No implementation required
- No maintenance burden
- No cross-shell compatibility concerns

Cons:
- Poor UX (users must type full commands and flags)
- No discoverability (users can't Tab to explore)
- Unprofessional (users expect completion in modern CLIs)
- More typos and errors
- Harder to learn CLI

Rejected: Shell completion is an expected feature. Modern CLIs without completion feel outdated.

Manual installation only (no auto-install):

Example: Only provide start completion bash, user handles installation
```bash
start completion bash > ~/.bash_completion.d/start
echo 'source ~/.bash_completion.d/start' >> ~/.bashrc
```

Pros:
- Simpler implementation (no path detection, no OS-specific logic)
- Maximum flexibility for users
- Less code to maintain

Cons:
- Higher friction (multi-step process, users must know paths)
- OS and shell-specific knowledge required
- Many users won't bother (lower adoption)
- More support questions about installation

Rejected: Auto-install significantly improves adoption and reduces friction. Implementation complexity is manageable.

Support all completion types including Tier 3:

Example: Complete model names, context names, role names, file paths
```bash
start --agent claude --model <Tab>
# Shows: haiku, sonnet, opus

start config context edit <Tab>
# Shows: environment, index, project, agents
```

Pros:
- Most complete experience
- Every value completable
- Maximum polish

Cons:
- Requires parsing agent model tables (complex, error-prone)
- File path completion redundant (shells do this)
- Diminishing returns (model/context names less frequently used)
- More maintenance burden
- More code complexity

Rejected: Tier 1 + Tier 2 provides sufficient value. Tier 3 adds complexity for marginal benefit.

## Structure

Completion command structure:

Root command:
- start completion <shell> - Print completion script to stdout
- start completion install <shell> [flags] - Auto-install to standard location

Supported shells:
- bash
- zsh
- fish

Install flags:
- --user (default) - Install to user directory
- --system - Install system-wide (requires sudo)
- --path <path> - Install to custom path

Standard installation paths:

bash:
- User: ~/.bash_completion
- System: /etc/bash_completion.d/start

zsh:
- User: ~/.zsh/completion/_start
- System: /usr/local/share/zsh/site-functions/_start

fish:
- User: ~/.config/fish/completions/start.fish
- System: /usr/share/fish/vendor_completions.d/start.fish

Completion tiers:

Tier 1 (static, free from Cobra):
- All commands and subcommands
- All flags (long and short forms)
- Help text for each

Tier 2 (dynamic, custom implementations):
- Agent names: Read from config (global + local merge)
- Task names: Read from config and catalog (show alias in parentheses)
- Scope values: global, local (or global, local, merged for list commands)

Tier 3 (not implemented):
- Model names (requires parsing agent model tables)
- Context names, role names
- Custom file path filtering

## Usage Examples

Manual installation (bash):

```bash
$ start completion bash > ~/.bash_completion/start

# Add to .bashrc
$ echo 'source ~/.bash_completion/start' >> ~/.bashrc

# Reload
$ source ~/.bashrc

# Test
$ start <Tab>
prompt  task  init  config  doctor  assets  completion  help
```

Auto-install (zsh):

```bash
$ start completion install zsh

Creating directory: /Users/grant/.zsh/completion
✓ Completion installed to: /Users/grant/.zsh/completion/_start

Reload your shell:
  source ~/.zshrc

Or start a new terminal session.

$ source ~/.zshrc

# Test
$ start --agent <Tab>
claude  gemini  aichat
```

System-wide install (fish):

```bash
$ sudo start completion install fish --system

✓ Completion installed to: /usr/share/fish/vendor_completions.d/start.fish

Reload your shell:
  source ~/.config/fish/config.fish

Fish completion is now available for all users.
```

Command completion:

```bash
$ start con<Tab>
$ start config   # Completed

$ start config ag<Tab>
$ start config agent   # Completed

$ start config agent li<Tab>
$ start config agent list   # Completed
```

Flag completion:

```bash
$ start --ag<Tab>
$ start --agent   # Completed

$ start --agent <Tab>
claude  gemini  aichat

$ start --agent cl<Tab>
$ start --agent claude   # Completed
```

Task completion with aliases:

```bash
$ start task <Tab>
code-review (cr)     git-diff-review (gdr)     comment-tidy (ct)     doc-review (dr)

$ start task co<Tab>
$ start task code-review   # Completed

# Alias matching works too
$ start task cr<Tab>
$ start task code-review   # Completed (matched alias)
```

Scope completion:

```bash
$ start init <Tab>
global  local

$ start config edit <Tab>
global  local

$ start config task list <Tab>
global  local  merged
```

Custom path installation:

```bash
$ start completion install zsh --path ~/.my-completions/_start

✓ Completion installed to: /Users/grant/.my-completions/_start

Add to your .zshrc:
  fpath=(~/.my-completions $fpath)
  autoload -Uz compinit && compinit

Reload your shell:
  source ~/.zshrc
```

## Updates

- 2025-01-17: Initial version aligned with schema
