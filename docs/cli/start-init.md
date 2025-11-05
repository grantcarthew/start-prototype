# start init

## Name

start init - Initialize start configuration

## Synopsis

```bash
start init [flags]
```

## Description

Interactive wizard to create initial `start` configuration. Detects installed AI agents, fetches current agent configurations from GitHub, and creates a working config file.

**What init does:**

1. Checks for existing configuration (offers backup if found)
2. Fetches latest agent configurations from GitHub
3. Auto-detects installed agents in PATH
4. Prompts for additional agents to configure
5. Prompts for default agent selection
6. Creates config at `~/.config/start/config.toml`
7. Adds default context document configuration

**When to run init:**

- First time setup
- Reset configuration to defaults
- Add newly installed agents

**Network requirement:**
Init requires network access to fetch agent configurations from GitHub. If offline, see manual configuration documentation.

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

### Execution Flow

**With existing config:**

1. Check if `~/.config/start/config.toml` exists
2. Prompt: "Backup and reinitialize? [y/N]"
   - If N: Exit with message
   - If Y (or `--force`): Continue
3. Backup existing config to `config.YYYY-MM-DD-HHMMSS.toml`
4. Continue to main wizard

**Main wizard:**

1. Fetch agent configs from GitHub (`assets/agents/*.toml`)
   - Timeout: 10 seconds
   - Endpoint: `https://api.github.com/repos/grantcarthew/start/contents/assets/agents`
2. Auto-detect installed agents using `command -v`
   - Checks for: claude, gemini, aichat, opencode, codex, aider
3. Auto-configure all detected agents
4. Prompt for additional agents (from fetched configs)
5. Prompt for default agent selection
6. Create config file structure:
   - `[settings]` section with default_agent
   - `[agents.*]` sections for each selected agent
   - `[context.system_prompt]` section
   - `[context.documents.*]` sections (4 default documents)
7. Write config to `~/.config/start/config.toml`
8. Display success message

**Default context documents:**
These documents are always added to the config:

1. `~/reference/ENVIRONMENT.md` (required = true)
2. `~/reference/INDEX.csv`
3. `./AGENTS.md`
4. `./PROJECT.md`

**Default system prompt:**
`./ROLE.md` is always added to config.

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

### First Time Setup

```bash
start init
```

Output:

```
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

Creating configuration at ~/.config/start/config.toml...
✓ Configuration created

Default context documents configured:
  ~/reference/ENVIRONMENT.md (required)
  ~/reference/INDEX.csv
  ./AGENTS.md
  ./PROJECT.md

Run 'start config show' to see your configuration.
Run 'start' to launch!
```

### Reinitialize Existing Config

```bash
start init
```

Output:

```
Configuration already exists at ~/.config/start/config.toml

Backup and reinitialize? [y/N]: y

Backing up to config.2025-01-04-103045.toml...
✓ Backup created

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
  Path: ~/.config/start/config.toml
  Exists: true

Prompting for backup confirmation...

Backing up existing config...
  From: ~/.config/start/config.toml
  To: ~/.config/start/config.2025-01-04-103045.toml
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

**Cannot write config file:**

```
Error: Failed to write config file: ~/.config/start/config.toml

Permission denied. Check file permissions.
```

Exit code: 3

**Backup fails:**

```
Error: Failed to backup existing config.

Check permissions: ~/.config/start/
Existing config preserved at: ~/.config/start/config.toml
```

Exit code: 3

Does not proceed with initialization. Existing config remains untouched.

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

- Shown only if `~/.config/start/config.toml` exists
- Skipped with `--force` flag
- Answer 'N' exits gracefully (exit code 0)

**Backup naming:**
Format: `config.YYYY-MM-DD-HHMMSS.toml`

Example: `config.2025-01-04-103045.toml`

Multiple backups accumulate (not overwritten).

### Config Directory Creation

If `~/.config/start/` doesn't exist, init creates it automatically.

Standard permissions: `0755` (drwxr-xr-x)

### Generated Config Structure

Example generated config:

```toml
[settings]
default_agent = "claude"

[agents.claude]
command = "claude --model {model} --append-system-prompt '{system_prompt}' '{prompt}'"
default_model = "sonnet"

  [agents.claude.models]
  haiku = "claude-3-5-haiku-20241022"
  sonnet = "claude-3-7-sonnet-20250219"
  opus = "claude-opus-4-20250514"

[agents.gemini]
command = "gemini --model {model} '{prompt}'"
default_model = "flash"

  [agents.gemini.models]
  flash = "gemini-2.0-flash-exp"
  pro-exp = "gemini-2.0-pro-exp"

[context.system_prompt]
path = "./ROLE.md"

[context.documents.environment]
path = "~/reference/ENVIRONMENT.md"
prompt = "Read {file} for environment context."
required = true

[context.documents.index]
path = "~/reference/INDEX.csv"
prompt = "Read {file} for documentation index."

[context.documents.agents]
path = "./AGENTS.md"
prompt = "Read {file} for repository context."

[context.documents.project]
path = "./PROJECT.md"
prompt = "Read {file} for current project status."
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

If init doesn't meet your needs (offline, custom setup, etc.), create config manually:

```bash
mkdir -p ~/.config/start
$EDITOR ~/.config/start/config.toml
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
