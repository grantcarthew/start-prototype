# start init

## Name

start init - Initialize start configuration

## Synopsis

```bash
start init [-l|--local] [-f|--force]
```

## Description

Interactive wizard to create `start` configuration files. Detects installed AI agents, fetches current agent configurations from GitHub, and creates a working multi-file configuration.

**Interactive by default** - Prompts for configuration choices unless `--force` flag is used.

**Target locations:**

- **global** - Personal config at `~/.config/start/` (default)
- **local** - Project config at `./.start/` (use `--local` flag)

**Configuration files created**:

- `config.toml` - Settings only
- `agents.toml` - Agent configurations
- `roles.toml` - Role definitions
- `contexts.toml` - Context document references
- `tasks.toml` - Task definitions

**What init does:**

1. Asks for target location (global or local) unless `--local` or `--force` provided
2. Checks for existing configuration at target location (prompts for backup unless `--force`)
3. Fetches latest agent configurations from GitHub
4. Auto-detects installed agents
5. Interactive mode: Prompts for additional agents and default agent selection
6. Automatic mode (`--force`): Configures all detected agents, uses first as default
7. Creates multi-file config at target location
8. Adds default context document configuration

**When to run init:**

- First time setup
- Create local project configuration

**Network requirement:**
Init attempts to download asset catalog from GitHub. If offline, creates config files without downloading assets.

## Flags

**--local**, **-l**
: Create config in `./.start/` (project config). Without this flag, init will ask interactively or default to global with `--force`.

```bash
start init          # Interactive: asks global or local
start init --local  # Create local config (still interactive for wizard)
```

**--force**, **-f**
: Fully automatic mode. Skips all prompts (location choice, backup confirmation, agent wizard). Auto-configures all detected agents, uses first detected as default, auto-backs up if config exists.

```bash
start init --force         # Automatic global config
start init --local --force # Automatic local config (perfect for CI/CD)
```

**--version**, **-v**
: Display version information and exit.

This command also supports the standard global flags for verbosity and help: `--verbose`, `--debug`, and `--help`.

## Behavior

### Interactive Mode (Default)

**`start init` with no flags:**

```
Initialize start configuration

Where should this configuration be created?
  1) Global (~/.config/start/)
     Personal config across all projects
  2) Local (./.start/)
     Project config (can be committed to git)

Select [1-2] (default: 1):
```

Then continues with wizard (see Main Wizard section below).

**If config exists at chosen location:**

```
Existing config found: ~/.config/start/

Backup and reinitialize? [y/N]:
```

- Answer `y` → Backup config files to `*.YYYY-MM-DD-HHMMSS.toml`, continue with wizard
- Answer `N` or press Enter → Exit gracefully (exit code 0)

### Partially Interactive (--local)

**`start init --local`:**

Skips location question, creates local config at `./.start/`.
Still runs backup prompt (if exists) and agent wizard.

```bash
start init --local
```

```
Initialize start configuration

Creating local config at ./.start/...

[...agent wizard prompts...]
```

### Fully Automatic (--force)

**`start init --force`:**

Zero interaction. Creates global config with smart defaults:

- Defaults to global (`~/.config/start/`)
- Auto-backs up if config exists (no prompt)
- Auto-detects and configures all discovered agents
- Sets first detected agent as default (priority: claude > gemini > aichat > others)
- No wizard prompts

```bash
start init --force
```

```
Initialize start configuration

Creating global config at ~/.config/start/...
✓ Backed up existing config
✓ Fetched agent configs from GitHub
✓ Detected and configured: claude, gemini
✓ Default agent: claude
✓ Config created successfully
```

**`start init --local --force`:**

Same automatic behavior, but creates local config at `./.start/`.
Perfect for CI/CD pipelines or scripting.

```bash
start init --local --force
```

### Main Wizard (Interactive Mode)

In interactive mode (without `--force`), init runs this wizard:

1. Fetch catalog index from GitHub (`assets/index.csv`)
   - Timeout: 10 seconds
   - URL: `https://raw.githubusercontent.com/grantcarthew/start/main/assets/index.csv`
   - Filters for `type=agents`, extracts `bin` column
2. Auto-detect installed agents by checking if `<bin>` is executable
   - Checks each binary from index (e.g., claude, gemini, aichat, opencode, codex, aider)
   - Binaries that are found and executable are marked as detected
3. Download agent asset files only for detected agents (lazy loading)
   - Fetches `assets/agents/{category}/{name}.toml` and `.meta.toml` via raw.githubusercontent.com
   - Only downloads what's needed
4. Display detected agents
5. **Prompt**: Additional agents to configure (from catalog index)
6. Download asset files for any additional selected agents
7. **Prompt**: Default agent selection
8. Create multi-file configuration:
   - `config.toml` - Settings (default_agent, default_role, etc.)
   - `agents.toml` - Agent configurations for each selected agent
   - `roles.toml` - Default role definitions
   - `contexts.toml` - Context document references (4 default documents)
   - `tasks.toml` - Default task definitions
9. Write all config files to chosen directory
10. Display success message

### Automatic Mode (--force)

In automatic mode, init does the same steps but **skips all prompts**:

1. Fetch catalog index from GitHub (same)
2. Auto-detect installed agents by checking if `<bin>` from index is executable (same)
3. Download agent asset files for detected agents (same lazy loading)
4. **Auto**: Configure ALL detected agents (no prompt for additional agents)
5. **Auto**: Use first detected agent as default (priority order: claude > gemini > aichat > others)
6. Create config files (same structure)
7. Display success message

If no agents detected: Creates config with empty `[agents]` section (user can add later).

**Default context documents:**
These documents are always added to the config:

1. `~/reference/ENVIRONMENT.md` (required = true)
2. `~/reference/INDEX.csv`
3. `./AGENTS.md`
4. `./PROJECT.md`

**Default role:**
A `code-reviewer` role referencing `./ROLE.md` is always added to config.

Files don't need to exist - runtime gracefully handles missing files.

### Agent Detection

Init uses the catalog index to detect installed agents:

**Step 1: Fetch index**
```
GET https://raw.githubusercontent.com/grantcarthew/start/main/assets/index.csv
```

Parses CSV to extract all agent entries with their `bin` field.

**Step 2: Auto-detect**
```
Check if 'claude' is executable  → Detected
Check if 'gemini' is executable  → Detected
Check if 'aichat' is executable  → Not found
```

**Step 3: Lazy download**

Only downloads asset files for detected/selected agents:
```
GET https://raw.githubusercontent.com/grantcarthew/start/main/assets/agents/{category}/{name}.toml
GET https://raw.githubusercontent.com/grantcarthew/start/main/assets/agents/{category}/{name}.meta.toml
```

Each agent asset consists of 2 files: configuration (`.toml`) and metadata (`.meta.toml`).

**Benefits:**
- Efficient: 1 index download + N agent downloads (only what's needed)
- Fast: No rate limits on raw.githubusercontent.com
- Lazy: Only downloads agents you actually use

**Auto-configuration:**

- All detected agents are automatically configured
- Downloaded configs provide: `bin` field, command template with `{bin}` placeholder, model names, default model

**Unknown agents:**

If a binary is installed but not in the catalog index:
- Not auto-detected (index is source of truth for available agents)
- User can manually add later with `start assets add`

### GitHub Catalog Details

**Index file:**
```
URL: https://raw.githubusercontent.com/grantcarthew/start/main/assets/index.csv
Format: CSV with columns: type,category,name,description,tags,bin,sha,size,created,updated
Filter: type=agents
Extract: bin column for detection
```

**Agent asset files:**
```
URL patterns:
  https://raw.githubusercontent.com/grantcarthew/start/main/assets/agents/{category}/{name}.toml
  https://raw.githubusercontent.com/grantcarthew/start/main/assets/agents/{category}/{name}.meta.toml

TOML contains: bin, command (with {bin} placeholder), description, models, default_model
Meta contains: SHA, size, timestamps for cache management
```

**Rate limits:**
- raw.githubusercontent.com has no rate limits
- Fast and reliable for init workflow

**Timeout:**
- 10 seconds for index fetch
- 5 seconds per agent TOML download
- Error and exit if timeout reached

## Examples

### Interactive Setup

```bash
start init
```

Output:

```
Initialize start configuration

Where should this configuration be created?
  1) Global (~/.config/start/)
     Personal config across all projects
  2) Local (./.start/)
     Project config (can be committed to git)

Select [1-2] (default: 1): 1

Welcome to start!

Fetching latest agent configurations from GitHub...
✓ Found 6 agent configurations

Detecting installed agents...
✓ claude (Claude Code by Anthropic)
✓ gemini (Gemini CLI by Google)

Additional agents available (not detected):
  [ ] aichat - All-in-one multi-provider CLI
  [ ] opencode - Open-source coding agent
  [ ] codex - OpenAI Codex CLI
  [ ] aider - Popular coding assistant

Select additional agents to configure (space to select, enter to continue):

Select default agent:
  1) claude
  2) gemini
Default [1]: 1

Creating configuration at ~/.config/start/...
✓ config.toml created
✓ agents.toml created
✓ roles.toml created
✓ contexts.toml created
✓ tasks.toml created

Default context documents configured:
  ~/reference/ENVIRONMENT.md (required)
  ~/reference/INDEX.csv
  ./AGENTS.md
  ./PROJECT.md

Run 'start config show' to see your configuration.
Run 'start' to launch!
```

### Local Config (Interactive)

```bash
start init --local
```

Skips location question, creates local config, still runs wizard:

```
Initialize start configuration

Creating local config at ./.start/...

Welcome to start!

Fetching latest agent configurations from GitHub...
✓ Found 6 agent configurations

Detecting installed agents...
✓ claude (Claude Code by Anthropic)
✓ gemini (Gemini CLI by Google)

[...wizard prompts for additional agents and default selection...]

Writing configuration files...
✓ config.toml created
✓ agents.toml created
✓ roles.toml created
✓ contexts.toml created
✓ tasks.toml created

Default context documents configured:
  ~/reference/ENVIRONMENT.md (required)
  ~/reference/INDEX.csv
  ./AGENTS.md
  ./PROJECT.md

Local config created. This can be committed to git for team consistency.
Run 'start config show' to see your configuration.
Run 'start' to launch!
```

### Automatic Mode

```bash
start init --force
```

Zero interaction, smart defaults:

```
Initialize start configuration

Creating global config at ~/.config/start/...

Fetching latest agent configurations from GitHub...
✓ Found 6 agent configurations

Detecting installed agents...
✓ claude (Claude Code by Anthropic)
✓ gemini (Gemini CLI by Google)

Auto-configuring detected agents...
✓ claude configured
✓ gemini configured
✓ Default agent: claude

Writing configuration files...
✓ config.toml created
✓ agents.toml created
✓ roles.toml created
✓ contexts.toml created
✓ tasks.toml created

Configuration created successfully.
Run 'start' to launch!
```

### Automatic Local Config (CI/CD)

```bash
start init --local --force
```

Perfect for automated setup in CI/CD or team onboarding scripts:

```
Initialize start configuration

Creating local config at ./.start/...
✓ Detected and configured: claude, gemini
✓ Default agent: claude
✓ Config created successfully
```

### Reinitialize Existing Config

```bash
start init
```

Output:

```
Initialize start configuration

Where should this configuration be created?
  1) Global (~/.config/start/)
     Personal config across all projects
  2) Local (./.start/)
     Project config (can be committed to git)

Select [1-2] (default: 1): 1

Existing config found: ~/.config/start/

Backup and reinitialize? [y/N]: y

Backing up config files...
✓ config.2025-01-14-143045.toml
✓ agents.2025-01-14-143045.toml
✓ roles.2025-01-14-143045.toml
✓ contexts.2025-01-14-143045.toml
✓ tasks.2025-01-14-143045.toml

Welcome to start!
[...continues with wizard...]
```

### Force Reinitialize

```bash
start init --force
```

Automatically backs up and reinitializes without any prompts:

```
Initialize start configuration

Creating global config at ~/.config/start/...
✓ Backed up existing config (config.2025-01-14-143045.toml, ...)
✓ Detected and configured: claude, gemini
✓ Default agent: claude
✓ Config created successfully
```

### Verbose Mode

```bash
start init --verbose
```

Output:

```
Checking for existing configuration...
  Path: ~/.config/start/
  Exists: true
  Files: config.toml, agents.toml, roles.toml, contexts.toml

Prompting for backup confirmation...

Backing up existing config files...
  config.toml → config.2025-01-04-103045.toml
  agents.toml → agents.2025-01-04-103045.toml
  roles.toml → roles.2025-01-04-103045.toml
  contexts.toml → contexts.2025-01-04-103045.toml
  ✓ Backup successful

Fetching agent configurations from GitHub...
  Endpoint: https://api.github.com/repos/grantcarthew/start/contents/assets/agents
  Timeout: 10s
  ✓ Response: 200 OK
  ✓ Found 6 files: claude.toml, gemini.toml, aichat.toml, opencode.toml, codex.toml, aider.toml

Downloading agent configs...
  ✓ claude.toml (1.2 KB)
  ✓ gemini.toml (980 B)
  [...]

Detecting installed agents...
  Checking: claude
    ✓ Found: /usr/local/bin/claude
  Checking: gemini
    ✓ Found: /usr/local/bin/gemini
  [...]

[...continues with wizard...]
```

## Output

### No Agents Detected

```bash
Detecting installed agents...
✗ No agents detected

Available agents:
  [ ] claude - Claude Code by Anthropic
  [ ] gemini - Gemini CLI by Google
  [ ] aichat - All-in-one multi-provider CLI
  [ ] opencode - Open-source coding agent
  [ ] codex - OpenAI Codex CLI
  [ ] aider - Popular coding assistant
  [ ] Other...

Select agents to configure (space to select, enter to continue):
```

### User Selects "Other..." Only

```bash
No agents configured.

To configure custom agents, see the documentation:
https://github.com/grantcarthew/start#configuration

Create your config manually:
  start config edit
```

No config file created. Exit code: 0

### User Selects No Agents

```bash
No agents configured.

To configure custom agents, see the documentation:
https://github.com/grantcarthew/start#configuration

Create your config manually:
  start config edit
```

No config file created. Exit code: 0

## Exit Codes

**0** - Success (config created or user chose not to proceed)

**1** - Configuration error

- Invalid TOML from GitHub
- Config validation failed

**3** - File system error

- Cannot create config directory
- Cannot write config file
- Backup failed

**4** - Network/runtime error

- Cannot reach GitHub
- API timeout (10 seconds)
- GitHub API rate limit exceeded
- Invalid API response

## Error Handling

### Network Errors

**Cannot reach GitHub:**

```
Error: Failed to fetch agent configurations from GitHub.

Check your network connection and try again.

See https://github.com/grantcarthew/start#configuration for manual setup.
```

Exit code: 4

**API timeout:**

```
Error: Request to GitHub timed out (10 seconds).

Check your network connection and try again.
```

Exit code: 4

**GitHub API rate limit:**

```
Error: GitHub API rate limit exceeded.

Try again in 42 minutes, or authenticate to increase limits.
See: https://docs.github.com/en/rest/overview/rate-limits
```

Exit code: 4

### File System Errors

**Cannot create config directory:**

```
Error: Failed to create config directory: ~/.config/start/

Permission denied. Check directory permissions.
```

Exit code: 3

**Cannot write config files:**

```
Error: Failed to write config files: ~/.config/start/

Permission denied. Check file permissions.
```

Exit code: 3

**Backup fails:**

```
Error: Failed to backup existing config files.

Check permissions: ~/.config/start/
Existing config files preserved.
```

Exit code: 3

Does not proceed with initialization. Existing config files remain untouched.

### Invalid Agent Config from GitHub

If a fetched agent config contains invalid TOML:

```
Warning: Failed to parse claude.toml from GitHub (invalid TOML).
Skipping claude configuration.

Continuing with remaining agents...
```

Init continues with other agents. Does not exit.

## Notes

### First Time vs Reinitialize

**Prompt behavior summary:**

| Command | Location prompt | Backup prompt | Agent wizard |
|---------|----------------|---------------|--------------|
| `start init` | Yes | Yes (if exists) | Yes |
| `start init --local` | No (→ local) | Yes (if exists) | Yes |
| `start init --force` | No (→ global) | No (auto-backup) | No (auto-config) |
| `start init --local --force` | No (→ local) | No (auto-backup) | No (auto-config) |

**Backup prompt:**
- Shown in interactive mode if config exists at target location
- Skipped with `--force` flag (auto-backs up instead)
- Answer 'N' exits gracefully (exit code 0)

**Backup naming:**
Format: `<filename>.YYYY-MM-DD-HHMMSS.toml`

Examples:

- `config.2025-01-04-103045.toml`
- `agents.2025-01-04-103045.toml`
- `roles.2025-01-04-103045.toml`
- `contexts.2025-01-04-103045.toml`

Multiple backups accumulate (not overwritten).

### Config Directory Creation

If `~/.config/start/` doesn't exist, init creates it automatically.

Standard permissions: `0755` (drwxr-xr-x)

### Generated Config Structure

Init creates 5 configuration files. Example content shown below:

**config.toml** (settings):

```toml
[settings]
default_agent = "claude"
default_role = "code-reviewer"
# ...
```

**agents.toml** (agent configurations):

```toml
[agents.claude]
bin = "claude"
command = "{bin} --model {model} ..."
# ...
```

**roles.toml** (role definitions):

```toml
[roles.code-reviewer]
description = "Expert code reviewer"
file = "./ROLE.md"
```

**contexts.toml** (context documents):

```toml
[contexts.environment]
file = "~/reference/ENVIRONMENT.md"
required = true
# ...
```

**tasks.toml** (task definitions):

```toml
# This file is created by 'start init'.
# Add custom tasks here or install them from the
# asset catalog using 'start assets add'.
```

**Note:** Agent model names and command templates come from GitHub. Above is illustrative only.

### Post-Init Steps

After running init:

1. **Edit model names** - Update model identifiers to current versions:

   ```bash
   start config edit
   ```

2. **Create context documents** - Add the files referenced in config:

   ```bash
   # Example
   mkdir -p ~/reference
   echo "# Environment" > ~/reference/ENVIRONMENT.md
   echo "# Agents" > AGENTS.md
   ```

3. **Test configuration**:
   ```bash
   start config show     # View config
   start --agent claude  # Test launch
   ```

### Manual Configuration Alternative

If init doesn't meet your needs (offline, custom setup, etc.), create config files manually:

```bash
mkdir -p ~/.config/start
$EDITOR ~/.config/start/config.toml
$EDITOR ~/.config/start/agents.toml
$EDITOR ~/.config/start/roles.toml
$EDITOR ~/.config/start/contexts.toml
```

See: https://github.com/grantcarthew/start#configuration

### Re-running Init

Running `start init` multiple times is safe:

- Always prompts for backup (unless `--force`)
- Previous backups preserved (timestamped)
- Fetches latest agent configs from GitHub

**Common reasons to re-run:**

- Add newly installed agents (`start init --force` for quick refresh)
- Reset to defaults after config errors
- Update agent templates from GitHub
- Create team config for new project (`start init --local --force`)

### Agent Detection Limitations

Agent detection only finds executables that are discoverable. If an agent is not found:

- Not auto-detected
- In interactive mode: Select manually from "Additional agents" list
- In automatic mode (`--force`): Not configured (add later with `start assets add`)

### GitHub Repository Structure

Init fetches from:

```
https://github.com/grantcarthew/start/tree/main/assets/agents/
```

Repository structure:

```
start/
├── assets/
│   ├── agents/
│   │   ├── anthropic/
│   │   │   └── claude.toml
│   │   ├── google/
│   │   │   └── gemini.toml
│   │   └── open/
│   │       ├── aichat.toml
│   │       └── opencode.toml
│   ├── tasks/
│   └── roles/
```

New agents added to repository become available immediately (no start release needed).

## See Also

- start-assets(1) - Manage catalog assets
- start-config(1) - Manage configuration
- start-config-agent(1) - Manage agents
- start(1) - Launch with context
