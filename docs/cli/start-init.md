# start init

## Name

start init - Initialize start configuration

## Synopsis

```bash
start init [scope]
```

## Description

Interactive wizard to create `start` configuration files. Detects installed AI agents, fetches current agent configurations from GitHub, and creates a working multi-file configuration.

**Scopes:**
- **global** - Personal config at `~/.config/start/` (default, recommended)
- **local** - Project config at `./.start/` (for team use)

**Configuration files created** (per DR-031):
- `config.toml` - Settings only
- `agents.toml` - Agent configurations
- `roles.toml` - Role definitions
- `contexts.toml` - Context document references
- `tasks.toml` - Task definitions

**What init does:**

1. Checks for existing configuration (offers backup if found)
2. Asks where to create config (global or local) if scope not specified
3. Fetches latest agent configurations from GitHub
4. Auto-detects installed agents in PATH
5. Prompts for additional agents to configure
6. Prompts for default agent selection
7. Creates config at chosen location
8. Adds default context document configuration

**When to run init:**

- First time setup
- Reset configuration to defaults
- Create local project configuration

**Network requirement:**
Init attempts to download asset catalog from GitHub. If offline, creates config files without downloading assets (per DR-026).

## Arguments

**[scope]** (optional)
: Where to create the configuration. If omitted, init will ask interactively.

- `global` - Create config files in `~/.config/start/` (personal, default)
- `local` - Create config files in `./.start/` (project, can be committed)

```bash
start init          # Interactive: asks where to create
start init global   # Create global config
start init local    # Create local config
```

## Flags

**--force**, **-f**
: Skip backup confirmation prompt. Automatically backs up existing config and continues.

```bash
start init --force
```

**--verbose**, **-v**
: Show detailed output including GitHub API calls and file operations.

**--debug**
: Debug mode. Shows all internal operations, API responses, and config generation.

**--help**, **-h**
: Show help text.

## Behavior

### Smart Behavior (No Scope Argument)

When run without scope argument, init detects existing configs and asks where to create.

**Scenario 1: No global, no local (first-time user)**

```
Initialize start configuration

No existing configuration found.

Where should this configuration be created?
  1) Global (~/.config/start/) [RECOMMENDED]
     Personal config across all projects
  2) Local (./.start/)
     Project config (for team use)

Select [1-2] (default: 1):
```

Press Enter → Creates global (default)

**Scenario 2: Global exists, no local**

```
Initialize start configuration

Existing global config found: ~/.config/start/

What would you like to do?
  1) Replace global config [BACKUP WILL BE CREATED]
  2) Create local config for this project
  3) Cancel

Select [1-3] (default: 1):
```

Option 1 → Backup config files to `*.YYYY-MM-DD-HHMMSS.toml`, create new global
Option 2 → Create local, keep global untouched
Option 3 → Exit

**Scenario 3: Global exists, local exists**

```
Initialize start configuration

Existing configs found:
  Global: ~/.config/start/
  Local:  ./.start/

Which would you like to replace? [BACKUP WILL BE CREATED]
  1) Replace global config
  2) Replace local config
  3) Cancel

Select [1-3] (default: 1):
```

**Scenario 4: No global, local exists**

```
Initialize start configuration

Existing local config found: ./.start/

What would you like to do?
  1) Create global config [RECOMMENDED]
  2) Replace local config [BACKUP WILL BE CREATED]
  3) Cancel

Select [1-3] (default: 1):
```

### Explicit Scope Behavior

When scope argument provided, skip prompts and create specified config:

```bash
start init global    # Force global (backup if exists)
start init local     # Force local (backup if exists)
```

Both create backups automatically if replacing existing config.

**Main wizard:**

1. Fetch agent configs from GitHub (`assets/agents/*.toml`)
   - Timeout: 10 seconds
   - Endpoint: `https://api.github.com/repos/grantcarthew/start/contents/assets/agents`
2. Auto-detect installed agents using `command -v`
   - Checks for: claude, gemini, aichat, opencode, codex, aider
3. Auto-configure all detected agents
4. Prompt for additional agents (from fetched configs)
5. Prompt for default agent selection
6. Create multi-file configuration:
   - `config.toml` - Settings (default_agent, default_role, etc.)
   - `agents.toml` - Agent configurations for each selected agent
   - `roles.toml` - Default role definitions
   - `contexts.toml` - Context document references (4 default documents)
   - `tasks.toml` - Default task definitions
7. Write all config files to chosen directory
8. Display success message

**Default context documents:**
These documents are always added to the config:

1. `~/reference/ENVIRONMENT.md` (required = true)
2. `~/reference/INDEX.csv`
3. `./AGENTS.md`
4. `./PROJECT.md`

**Default role:**
A `code-reviewer` role referencing `./ROLE.md` is always added to config.

Files don't need to exist - runtime gracefully handles missing files (see DR-008).

### Agent Detection

Init uses `command -v` to detect installed agents:

```bash
command -v claude      # Detected
command -v gemini      # Detected
command -v aichat      # Not found
```

**Auto-configuration:**

- All detected agents are automatically configured
- No user prompt for detected agents
- Fetched configs from GitHub provide command templates and model aliases

**Unknown agents:**
If `command -v` finds a binary but no config exists in GitHub:

- Agent is skipped
- User can select "Other..." to see manual configuration docs

### GitHub Fetch Details

**API endpoint:**

```
GET https://api.github.com/repos/grantcarthew/start/contents/assets/agents
```

**Fetched for each agent:**

- Command template with placeholders
- Model aliases mapping
- Default model selection

**Rate limits:**

- Unauthenticated: 60 requests/hour
- Init uses 1-2 requests total

**Timeout:**

- 10 seconds for API calls
- Error and exit if timeout reached

## Examples

### First Time Setup (Interactive)

```bash
start init
```

Output:

```
Initialize start configuration

No existing configuration found.

Where should this configuration be created?
  1) Global (~/.config/start/) [RECOMMENDED]
     Personal config across all projects
  2) Local (./.start/)
     Project config (for team use)

Select [1-2] (default: 1): 1

Welcome to start!

Fetching latest agent configurations from GitHub...
✓ Found 6 agent configurations

Detecting installed agents...
✓ claude (Claude Code by Anthropic)
✓ gemini (Gemini CLI by Google)

Configuring detected agents...
✓ claude configured
✓ gemini configured

Additional agents available (not detected):
  [ ] aichat - All-in-one multi-provider CLI
  [ ] opencode - Open-source coding agent
  [ ] codex - OpenAI Codex CLI
  [ ] aider - Popular coding assistant
  [ ] Other...

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

Default context documents configured:
  ~/reference/ENVIRONMENT.md (required)
  ~/reference/INDEX.csv
  ./AGENTS.md
  ./PROJECT.md

Run 'start config show' to see your configuration.
Run 'start' to launch!
```

### Create Global Config (Explicit)

```bash
start init global
```

Skips scope prompt, creates global config directly.

### Create Local Config (Explicit)

```bash
start init local
```

Output:

```
Welcome to start!

Fetching latest agent configurations from GitHub...
✓ Found 6 agent configurations

[...agent detection and wizard...]

Creating configuration at ./.start/...
✓ config.toml created
✓ agents.toml created
✓ roles.toml created
✓ contexts.toml created

Default context documents configured:
  ~/reference/ENVIRONMENT.md (required)
  ~/reference/INDEX.csv
  ./AGENTS.md
  ./PROJECT.md

Local config created. This can be committed to git for team consistency.
Run 'start config show' to see your configuration.
Run 'start' to launch!
```

### Reinitialize Existing Config

```bash
start init
```

Output:

```
Configuration already exists at ~/.config/start/

Backup and reinitialize? [y/N]: y

Backing up config files...
✓ config.2025-01-04-103045.toml
✓ agents.2025-01-04-103045.toml
✓ roles.2025-01-04-103045.toml
✓ contexts.2025-01-04-103045.toml

Welcome to start!
[...continues with wizard...]
```

### Force Reinitialize

```bash
start init --force
```

Skips backup prompt, automatically backs up and continues.

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
  Checking: claude (command -v claude)
    ✓ Found: /usr/local/bin/claude
  Checking: gemini (command -v gemini)
    ✓ Found: /usr/local/bin/gemini
  [...]

[...continues with wizard...]
```

## Output

### No Agents Detected

```bash
Detecting installed agents...
✗ No agents detected in PATH

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

**Backup prompt behavior:**

- Shown only if `~/.config/start/` directory contains config files
- Skipped with `--force` flag
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
command = "claude --model {model} ..."
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
[context.environment]
file = "~/reference/ENVIRONMENT.md"
required = true
# ...
```

**tasks.toml** (task definitions):
```toml
# This file is created by 'start init'.
# Add custom tasks here or install them from the
# asset catalog using 'start config task add'.
```

**Note:** Agent model aliases and command templates come from GitHub. Above is illustrative only.

### Post-Init Steps

After running init:

1. **Edit model aliases** - Update model names to current versions:

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

- Add newly installed agents
- Reset to defaults after config errors
- Update agent templates from GitHub

### Agent Detection Limitations

`command -v` only finds agents in PATH. If an agent is installed but not in PATH:

- Not auto-detected
- Select manually from "Additional agents" list
- Or select "Other..." for custom configuration

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
│   │   ├── claude.toml
│   │   ├── gemini.toml
│   │   ├── aichat.toml
│   │   ├── opencode.toml
│   │   ├── codex.toml
│   │   └── aider.toml
│   ├── tasks/
│   └── roles/
```

New agents added to repository become available immediately (no start release needed).

## See Also

- start-config(1) - Manage configuration
- start-agent(1) - Manage agents
- start(1) - Launch with context
