# start config agent

## Name

start config agent - Manage AI agent configurations

## Synopsis

```bash
start config agent list [scope]
start config agent new [scope]
start config agent show [name] [scope]
start config agent test <name>
start config agent edit [name] [scope]
start config agent remove [name] [scope]
start config agent default [name]
```

## Description

Manages AI agent configurations in both global (`~/.config/start/agents.toml`) and local (`./.start/agents.toml`) files. Agents define how `start` delegates to different AI tools (claude, gemini, aichat, etc.).

**Agent management operations:**

- **list** - Display all configured agents with details
- **new** - Create new agent interactively
- **show** - Display agent configuration structure
- **test** - Test agent configuration and availability
- **edit** - Modify existing agent configuration
- **remove** - Delete agent from configuration
- **default** - Set or show default agent

To install agents from the catalog, use `start assets add`.

**Note:** Agents can be defined in both global and local configs. These commands can manage either scope using the `[scope]` argument. If scope is omitted, the command prompts interactively.

## Agent Configuration Structure

Agents are defined in the configuration with the following fields:

```toml
[agents.claude]
description = "Anthropic's Claude AI assistant via Claude Code CLI"
url = "https://docs.claude.com/claude-code"
models_url = "https://docs.anthropic.com/en/docs/about-claude/models"
bin = "claude"
command = "{bin} --model {model} --append-system-prompt '{role}' '{prompt}'"
default_model = "sonnet"

  [agents.claude.models]
  haiku = "claude-3-5-haiku-20241022"
  sonnet = "claude-3-7-sonnet-20250219"
  opus = "claude-opus-4-20250514"
```

**Fields:**

**bin** (required)
: Binary name or path to execute (e.g., `claude` or `/usr/local/bin/claude`). Used for binary detection (`start doctor`), the `{bin}` placeholder in command templates, and catalog index generation.

**command** (required)
: Command template to execute the agent. Should contain `{prompt}` placeholder. Supports placeholders: `{bin}`, `{model}`, `{role}`, `{role_file}`, `{prompt}`, `{date}`.

**description** (optional)
: Human-readable description of the agent. Displayed in `start config agent list` and help output.

**url** (optional)
: Documentation or homepage URL for the agent tool.

**models_url** (optional)
: URL to model documentation, helping users understand available models and their capabilities.

**default_model** (optional)
: Model name to use when `--model` flag not provided. If omitted, the first model in the `models` table is used.

**models** (optional)
: Table of user-defined model names mapping to full model identifiers. Each agent can define its own names.

## Subcommands

### start config agent list

Display all configured agents with their details.

**Synopsis:**

```bash
start config agent list          # Select scope interactively
start config agent list global   # List global agents only
start config agent list local    # List local agents only
start config agent list merged   # Show merged view (global + local)
```

**Behavior:**

Lists all agents defined in the selected scope(s) with:

- Agent name
- Description
- Documentation URL
- Default model (full identifier and model name)
- All available models (full identifier and model name)
- Model documentation URL

Missing optional fields are omitted from display.

**Output:**

```
Configured agents:

claude
  Anthropic's Claude AI assistant via Claude Code CLI
  https://docs.claude.com/claude-code
  Command: {bin} --model {model} --append-system-prompt '{role}' '{prompt}'
  Default model: claude-3-7-sonnet-20250219 (sonnet)
  Models:
    - claude-3-5-haiku-20241022 (haiku)
    - claude-3-7-sonnet-20250219 (sonnet)
    - claude-opus-4-20250514 (opus)
  Model docs: https://docs.anthropic.com/en/docs/about-claude/models

gemini
  Google's Gemini AI via CLI
  https://github.com/example/gemini-cli
  Command: {bin} --model {model} '{prompt}'
  Default model: gemini-2.0-flash-exp (flash)
  Models:
    - gemini-2.0-flash-exp (flash)
    - gemini-2.0-pro-exp (pro-exp)
  Model docs: https://ai.google.dev/models/gemini

aichat
  All-in-one multi-provider AI chat tool
  https://github.com/sigoden/aichat
  Command: {bin} --model {model} '{prompt}'
  Default model: gpt-4o-mini (gpt4-mini)
  Models:
    - gpt-4o-mini (gpt4-mini)
    - gpt-4o (gpt4)
    - claude-3-5-sonnet-20241022 (claude)
```

**Minimal agent (only required fields):**

```
opencode
  Command: {bin} '{prompt}'
  Default model: (first model in config)
```

**No agents configured:**

```
No agents configured.

Run 'start init' to set up agents, or
use 'start assets add' to install from catalog or 'start config agent new' to create custom.
```

**Exit codes:**

- 0 - Success (agents listed)
- 1 - No config file exists

### start config agent new

Interactively add a new agent to the configuration.

**Synopsis:**

```bash
start config agent new          # Select scope interactively
start config agent new global   # Add to global config
start config agent new local    # Add to local config
```

**Behavior:**

Prompts for agent details and adds to the selected configuration file:

1. **Agent name** (required)

   - Validation: lowercase alphanumeric with hyphens
   - Pattern: `/^[a-z0-9]+(-[a-z0-9]+)*$/`
   - Must be unique (not already exist)
   - Examples: `claude`, `gemini`, `my-custom-agent`

2. **Description** (optional)

   - Human-readable description
   - Press enter to skip

3. **URL** (optional)

   - Documentation or homepage URL
   - Press enter to skip

4. **Models URL** (optional)

   - Model documentation URL
   - Press enter to skip

5. **Binary** (required)
   - Binary name or path to execute (e.g., `claude`)
   - Used for detection, `{bin}` placeholder, and catalog index generation

6. **Command template** (required)

   - Should contain `{prompt}` placeholder (warns if missing)
   - Warns on unknown placeholders (typos)
   - Valid placeholders: {bin}, {model}, {role}, {role_file}, {prompt}, {date}

7. **Add models?** (yes/no)

   - If yes, loop to add multiple models
   - Each model: model name + full model identifier
   - Type "done" to finish adding models

8. **Default model** (if models added)

   - Shows numbered list of added models
   - Select which to use as default
   - Can skip (uses first model)

9. **Backup and save**
   - Backs up existing config to `config.YYYY-MM-DD-HHMMSS.toml`
   - Writes new agent to config
   - Shows success message

**Interactive flow:**

```
Add new agent
─────────────────────────────────────────────────

Agent name: my-agent
Description (optional): My custom AI agent
URL (optional): https://example.com/my-agent
Models URL (optional): https://example.com/models

Binary [my-agent]: my-agent
Command template: {bin} --model {model} '{prompt}'
✓ Valid command template

Add models? [y/N]: y

Model name: fast
Full model identifier: my-agent-fast-v1
✓ Model added: fast → my-agent-fast-v1

Model name: best
Full model identifier: my-agent-best-v2
✓ Model added: best → my-agent-best-v2

Model name: done

Select default model:
  1) my-agent-fast-v1 (fast)
  2) my-agent-best-v2 (best)
  [skip to use first model]
Default: 1

Backing up config to config.2025-01-04-143022.toml...
✓ Backup created

Saving agent 'my-agent' to ~/.config/start/agents.toml...
✓ Agent added successfully

Use 'start config agent list' to see all agents.
Use 'start --agent my-agent' to test.
```

**Minimal agent (no optional fields):**

```
Add new agent
─────────────────────────────────────────────────

Agent name: simple-agent
Description (optional):
URL (optional):
Models URL (optional):

Binary [simple-agent]: simple-agent
Command template: {bin} '{prompt}'
✓ Valid command template

Add models? [y/N]: n

Backing up config to config.2025-01-04-143105.toml...
✓ Backup created

Saving agent 'simple-agent' to ~/.config/start/agents.toml...
✓ Agent added successfully

Use 'start config agent list' to see all agents.
Use 'start --agent simple-agent' to test.
```

**Resulting config (full agent):**

```toml
[agents.my-agent]
description = "My custom AI agent"
url = "https://example.com/my-agent"
models_url = "https://example.com/models"
bin = "my-agent"
command = "{bin} --model {model} '{prompt}'"
default_model = "fast"

  [agents.my-agent.models]
  fast = "my-agent-fast-v1"
  best = "my-agent-best-v2"
```

**Resulting config (minimal agent):**

```toml
[agents.simple-agent]
bin = "simple-agent"
command = "{bin} '{prompt}'"
```

**Exit codes:**

- 0 - Success (agent added)
- 1 - Validation error (invalid name, duplicate agent, invalid command)
- 3 - File system error (cannot write config, backup failed)

**Error handling:**

**Invalid agent name:**

```
Agent name: My-Agent
✗ Invalid agent name. Use lowercase alphanumeric with hyphens.
  Examples: claude, gemini, my-agent

Agent name: my-agent
✓ Valid name
```

**Duplicate agent:**

```
Agent name: claude
✗ Agent 'claude' already exists.

Use 'start config agent edit claude' to modify existing agent.
```

Exit code: 1

**Command template missing {prompt} (warning):**

```
Command template: my-agent --model {model}
⚠ Warning: Command doesn't contain {prompt} placeholder.
  Composed prompt won't be passed to agent. Continue? [y/N]: n

Command template: my-agent --model {model} '{prompt}'
✓ Valid command template
```

**Unknown placeholder warning:**

```
Command template: my-agent --model {mdoel} '{prompt}'
⚠ Warning: Unknown placeholder {mdoel} (did you mean {model}?)
  Valid placeholders: {model}, {role}, {role_file}, {prompt}, {date}

Continue anyway? [y/N]: n

Command template: my-agent --model {model} '{prompt}'
✓ Valid command template
```

**Backup failed:**

```
Backing up config to config.2025-01-04-143022.toml...
✗ Failed to backup config: permission denied

Existing config preserved at: ~/.config/start/agents.toml
Agent not added.
```

Exit code: 3

### start config agent show

Display current agent configuration.

**Synopsis:**

```bash
start config agent show                 # Select agent and scope interactively
start config agent show <name>          # Select scope for named agent
start config agent show <name> global   # Show global agent only
start config agent show <name> local    # Show local agent only
```

**Behavior:**

Displays agent configuration from the selected scope with:

- Scope (global or local)
- Command template
- Description (if configured)
- URL (if configured)
- Models URL (if configured)
- Default model (if configured)
- All configured models

**Output (global agent):**

```
Agent configuration: claude (global)
═══════════════════════════════════════════════════════════

Description: Anthropic's Claude AI assistant via Claude Code CLI
URL: https://docs.claude.com/claude-code
Models URL: https://docs.anthropic.com/en/docs/about-claude/models
Binary: claude

Command template:
  {bin} --model {model} --append-system-prompt '{role}' '{prompt}'

Default model: sonnet
Models:
  haiku = claude-3-5-haiku-20241022
  sonnet = claude-3-7-sonnet-20250219
  opus = claude-opus-4-20250514
```

**Output (local agent):**

```bash
start config agent show custom-agent local
```

```
Agent configuration: custom-agent (local)
═══════════════════════════════════════════════════════════

Binary: custom-agent
Command template:
  {bin} '{prompt}'

Default model: (first model in config)
Models: (none)
```

**Output (minimal agent):**

```
Agent configuration: simple-agent (global)
═══════════════════════════════════════════════════════════

Binary: simple-agent
Command template:
  {bin} '{prompt}'
```

**No agent configured:**

```
No agent 'nonexistent' found in global config.

Configure: start config agent new global
```

**Interactive selection:**

```bash
start config agent show
```

```
Show agent configuration
─────────────────────────────────────────────────

Select agent:
  1) claude
  2) gemini
  3) aichat

Select [1-3]: 1

Select scope:
  1) global
  2) local

Scope [1-2]: 1

(displays agent configuration)
```

**Exit codes:**

- 0 - Success (agent shown)
- 1 - No agent configured
- 2 - Invalid scope argument
- 3 - Agent not found

**Error handling:**

**Agent not found:**

```
Error: Agent 'nonexistent' not found in configuration.

Use 'start config agent list' to see available agents.
```

Exit code: 3

### start config agent test

Test agent configuration and availability.

**Synopsis:**

```bash
start config agent test <name>
```

**Behavior:**

Validates agent configuration without executing it. Performs three checks:

1. **Binary availability**

   - Checks if the configured `bin` executable is discoverable
   - Reports: found (with path) or not found

2. **Configuration validation**

   - Command template checked for `{prompt}` placeholder (warns if missing)
   - Unknown placeholders detected (likely typos)
   - Model names defined (if `{model}` used in template)
   - Default model configured or first model available
   - TOML syntax valid

3. **Dry-run command construction**
   - Builds command with placeholders resolved
   - Uses default model
   - Uses test prompt: "test"
   - Displays command (system prompt truncated to `'...'`)

**Does NOT execute the agent** - only validates and shows what would run.

**Output (success):**

```
Testing agent: claude
─────────────────────────────────────────────────

✓ Binary found: /usr/local/bin/claude

Configuration:
  ✓ Command template valid
  ✓ Contains {prompt} placeholder
  ✓ Default model: claude-3-7-sonnet-20250219 (sonnet)
  ✓ Models configured: 3 (haiku, sonnet, opus)

Dry-run command:
  ❯ claude --model claude-3-7-sonnet-20250219 --append-system-prompt '...' 'test'

✓ Agent 'claude' is configured correctly
```

**Output (warnings):**

```
Testing agent: my-agent
─────────────────────────────────────────────────

✓ Binary found: /usr/local/bin/my-agent

Configuration:
  ✓ Command template valid
  ✓ Contains {prompt} placeholder
  ⚠ Unknown placeholder {mdoel} in command template
    (did you mean {model}?)
  ✓ Default model: my-agent-v1 (default)

Dry-run command:
  ❯ my-agent --model {mdoel} '{prompt}'

⚠ Agent 'my-agent' has warnings (see above)
```

**Output (binary not found):**

```
Testing agent: gemini
─────────────────────────────────────────────────

✗ Binary not found: gemini
  The 'gemini' command is not available.
  Install gemini or check that it's accessible.

Configuration:
  ✓ Command template valid
  ✓ Contains {prompt} placeholder
  ✓ Default model: gemini-2.0-flash-exp (flash)
  ✓ Models configured: 2 (flash, pro-exp)

Dry-run command:
  ❯ gemini --model gemini-2.0-flash-exp 'test'

✗ Agent 'gemini' is not available (binary not found)
```

**Output (configuration error):**

```
Testing agent: broken-agent
─────────────────────────────────────────────────

✓ Binary found: /usr/local/bin/broken-agent

Configuration:
  ⚠ Command template missing {prompt} placeholder
  ✗ No models configured but {model} used in template
  ✗ No default_model specified

✗ Agent 'broken-agent' has configuration errors
  Fix configuration: start config agent edit broken-agent
```

**Verbose output:**

```bash
start config agent test claude --verbose
```

```
Testing agent: claude
─────────────────────────────────────────────────

Checking binary availability...
  Command: claude
  Search PATH: /usr/local/bin:/usr/bin:/bin
  ✓ Found: /usr/local/bin/claude

Validating configuration...
  Config file: ~/.config/start/agents.toml
  Agent section: [agents.claude]

  Command template:
    claude --model {model} --append-system-prompt '{role}' '{prompt}'

  Placeholder analysis:
    ✓ {model} - valid
    ✓ {role} / {role_file} - valid
    ✓ {prompt} - valid

  Model configuration:
    ✓ default_model: sonnet
    ✓ Alias 'sonnet' defined: claude-3-7-sonnet-20250219
    ✓ Total models: 3
      - haiku: claude-3-5-haiku-20241022
      - sonnet: claude-3-7-sonnet-20250219
      - opus: claude-opus-4-20250514

Building dry-run command...
  Model: claude-3-7-sonnet-20250219 (using default_model: sonnet)
  System prompt: (truncated to '...')
  Prompt: 'test'
  Date: 2025-01-04T14:35:12+10:00

Dry-run command:
  ❯ claude --model claude-3-7-sonnet-20250219 --append-system-prompt '...' 'test'

✓ Agent 'claude' is configured correctly
```

**Exit codes:**

- 0 - Success (agent valid and binary found)
- 1 - Configuration error (invalid config)
- 2 - Agent not found in config
- 4 - Binary not found (config valid but tool not installed)

**Error handling:**

**Agent not in config:**

```
Error: Agent 'nonexistent' not found in configuration.

Use 'start config agent list' to see available agents.
Use 'start assets add' to install from catalog or 'start config agent new' to create custom.
```

Exit code: 2

**Multiple errors:**

```
Testing agent: broken
─────────────────────────────────────────────────

✗ Binary not found: broken

Configuration:
  ⚠ Command template missing {prompt} placeholder
  ⚠ Unknown placeholder {foo} in command template
  ✗ No default_model specified and no models configured

✗ Agent 'broken' has multiple errors:
  - Binary not found
  - Invalid command template
  - Missing model configuration
```

Exit code: 1 (configuration errors take precedence over binary not found)

### start config agent edit

Edit agent configuration interactively.

**Synopsis:**

```bash
start config agent edit                  # Select agent and scope
start config agent edit <name>           # Select scope for named agent
start config agent edit <name> global    # Edit in global config
start config agent edit <name> local     # Edit in local config
```

**Behavior:**

**Without agent name (interactive selection):**

Shows list of configured agents for selection:

```bash
start config agent edit
```

Output:

```
Edit agent
─────────────────────────────────────────────────

Select agent to edit:
  1) claude
  2) gemini
  3) aichat
  4) my-custom-agent

Select [1-4] (or 'q' to quit): 1

(continues to interactive edit flow for 'claude')
```

**With agent name (interactive edit):**

Interactive prompts to edit specific agent. Shows current values as defaults - press enter to keep current value.

1. **Description** - Current value shown in brackets
2. **URL** - Current value shown in brackets
3. **Models URL** - Current value shown in brackets
4. **Binary** - Current value shown in brackets
5. **Command template** - Current value shown in brackets
   - Validates: should contain `{prompt}` placeholder (warns if missing)
   - Warns on unknown placeholders
6. **Edit models?** - Show current models, ask to modify
   - Add new models
   - Remove existing models
   - Modify model values
7. **Default model** - Select from available models
8. **Backup and save** - Backs up to `config.YYYY-MM-DD-HHMMSS.toml`

**Interactive flow (edit specific agent):**

```
Edit agent: claude
─────────────────────────────────────────────────

Current configuration:
  Description: Anthropic's Claude AI assistant via Claude Code CLI
  URL: https://docs.claude.com/claude-code
  Models URL: https://docs.anthropic.com/en/docs/about-claude/models
  Binary: claude
  Command: {bin} --model {model} --append-system-prompt '{role}' '{prompt}'
  Default model: sonnet
  Models: 3 (haiku, sonnet, opus)

Press enter to keep current value, or type new value:

Description [Anthropic's Claude AI assistant via Claude Code CLI]:
URL [https://docs.claude.com/claude-code]:
Models URL [https://docs.anthropic.com/en/docs/about-claude/models]:
Binary [claude]:
Command template [{bin} --model {model} --append-system-prompt '{role}' '{prompt}']:

Current models:
  haiku = claude-3-5-haiku-20241022
  sonnet = claude-3-7-sonnet-20250219
  opus = claude-opus-4-20250514

Edit models? [y/N]: y

Add model (or "done" to finish):
Model name: haiku2
Full model identifier: claude-3-5-haiku-20241022-v2
✓ Model added: haiku2 → claude-3-5-haiku-20241022-v2

Add model (or "done" to finish): done

Remove models? [y/N]: n

Select default model:
  1) haiku = claude-3-5-haiku-20241022
  2) sonnet = claude-3-7-sonnet-20250219
  3) opus = claude-opus-4-20250514
  4) haiku2 = claude-3-5-haiku-20241022-v2
Current: sonnet [2]
Default [2]: 4

Backing up config to config.2025-01-04-144512.toml...
✓ Backup created

Saving changes to ~/.config/start/agents.toml...
✓ Agent 'claude' updated successfully

Use 'start config agent list' to see changes.
Use 'start config agent test claude' to validate.
```

**Interactive flow (minimal changes):**

```
Edit agent: simple-agent
─────────────────────────────────────────────────

Current configuration:
  Description: (none)
  URL: (none)
  Models URL: (none)
  Binary: simple-agent
  Command: {bin} '{prompt}'
  Default model: (none - uses first model)
  Models: (none)

Press enter to keep current value, or type new value:

Description []: Simple AI agent for testing
URL []: https://example.com/simple-agent
Models URL []:
Binary [simple-agent]:
Command template [{bin} '{prompt}']:

Current models: (none)

Add models? [y/N]: n

Backing up config to config.2025-01-04-144623.toml...
✓ Backup created

Saving changes to ~/.config/start/agents.toml...
✓ Agent 'simple-agent' updated successfully
```

**Exit codes:**

- 0 - Success (agent edited)
- 1 - Validation error (invalid name, invalid command)
- 2 - Agent not found
- 3 - File system error (cannot write config, backup failed)

**Error handling:**

**Agent not found:**

```
Error: Agent 'nonexistent' not found in configuration.

Use 'start config agent list' to see available agents.
Use 'start assets add' to install from catalog or 'start config agent new' to create custom.
```

Exit code: 2

**Command template missing {prompt} (warning):**

```
Command template [claude --model {model}]: claude --other-flag
⚠ Warning: Command doesn't contain {prompt} placeholder.
  Composed prompt won't be passed to agent. Continue? [y/N]: n

Command template [claude --model {model}]: claude --model {model} '{prompt}'
✓ Valid command template
```

**Unknown placeholder warning:**

```
Command template [claude '{prompt}']: claude --model {mdoel} '{prompt}'
⚠ Warning: Unknown placeholder {mdoel} (did you mean {model}?)
  Valid placeholders: {bin}, {model}, {role}, {role_file}, {prompt}, {date}

Continue anyway? [y/N]: n

Command template [claude '{prompt}']: {bin} --model {model} '{prompt}'
✓ Valid command template
```

**Model management details:**

**Adding models:**

- Validates model name (same rules as agent names: lowercase, alphanumeric, hyphens)
- Doesn't validate full model identifier (too variable across agents)
- Detects duplicate model names

**Removing models:**

```
Remove models? [y/N]: y

Select models to remove (space to select, enter to continue):
  [ ] haiku = claude-3-5-haiku-20241022
  [x] sonnet = claude-3-7-sonnet-20250219
  [ ] opus = claude-opus-4-20250514

✓ Removed: sonnet
```

If default_model is removed, user must select new default from remaining models.

**No changes made:**

```
No changes detected.

Agent 'claude' not modified.
```

Exit code: 0 (no backup created, no write)

### start config agent remove

Remove agent from configuration.

**Synopsis:**

```bash
start config agent remove                  # Select agent and scope
start config agent remove <name>           # Select scope for named agent
start config agent remove <name> global    # Remove from global config
start config agent remove <name> local     # Remove from local config
```

**Behavior:**

**Without agent name:**
Shows list of configured agents for selection:

```
Remove agent
─────────────────────────────────────────────────

Select agent to remove:
  1) claude
  2) gemini
  3) aichat
  4) my-custom-agent

Select [1-4] (or 'q' to quit): 2

Remove agent 'gemini'? [y/N]: y

Backing up config to config.2025-01-04-150212.toml...
✓ Backup created

Removing agent 'gemini' from ~/.config/start/agents.toml...
✓ Agent 'gemini' removed successfully

Use 'start config agent list' to see remaining agents.
```

**With agent name:**
Removes specific agent directly (with confirmation):

```bash
start config agent remove gemini
```

Output:

```
Remove agent 'gemini'? [y/N]: y

Backing up config to config.2025-01-04-150245.toml...
✓ Backup created

Removing agent 'gemini' from ~/.config/start/agents.toml...
✓ Agent 'gemini' removed successfully

Use 'start config agent list' to see remaining agents.
```

**Removing default agent:**

If removing the agent that's set as `default_agent`:

```bash
start config agent remove claude
```

Output:

```
⚠ Warning: 'claude' is currently your default agent.

Remove agent 'claude'? [y/N]: y

Backing up config to config.2025-01-04-150312.toml...
✓ Backup created

Removing agent 'claude' from ~/.config/start/agents.toml...
✓ Agent 'claude' removed successfully
✓ Default agent setting cleared

Your default agent is now the first configured agent: gemini

Use 'start config agent default <name>' to set a new default.
Use 'start config agent list' to see remaining agents.
```

**Behavior when removing default:**

- Removes agent from config
- Removes `default_agent` setting from `[settings]` section
- `start` command will use first agent in config (TOML order)
- User can set new default with `start config agent default <name>`

**Declining confirmation:**

```
Remove agent 'gemini'? [y/N]: n

Agent 'gemini' not removed.
```

Exit code: 0

**Exit codes:**

- 0 - Success (agent removed, or user declined)
- 1 - No agents configured
- 2 - Agent not found
- 3 - File system error (cannot write config, backup failed)

**Error handling:**

**Agent not found:**

```
Error: Agent 'nonexistent' not found in configuration.

Use 'start config agent list' to see available agents.
```

Exit code: 2

**No agents configured:**

```
No agents configured.

Use 'start assets add' to install from catalog or 'start config agent new' to create custom.
```

Exit code: 1

**Only one agent configured:**

```
Warning: 'claude' is the only configured agent.

Remove agent 'claude'? [y/N]: y

Backing up config to config.2025-01-04-150412.toml...
✓ Backup created

Removing agent 'claude' from ~/.config/start/agents.toml...
✓ Agent 'claude' removed successfully
⚠ No agents remaining in configuration

Use 'start assets add' to install from catalog or 'start config agent new' to create custom.
Use 'start init' to set up agents automatically.
```

**Backup failed:**

```
Remove agent 'gemini'? [y/N]: y

Backing up config to config.2025-01-04-150445.toml...
✗ Failed to backup config: permission denied

Existing config preserved at: ~/.config/start/agents.toml
Agent not removed.
```

Exit code: 3

### start config agent default

Set default agent interactively or directly.

**Synopsis:**

```bash
start config agent default          # Select from list
start config agent default <name>   # Set specific default
```

**Behavior:**

**Without agent name (interactive selection):**

Shows list of configured agents for selection:

```bash
start config agent default
```

Output:

```
Set default agent
─────────────────────────────────────────────────

Current default: claude

Select new default agent:
  1) claude (current)
  2) gemini
  3) aichat
  4) my-custom-agent

Select [1-4] (or 'q' to quit): 2

Backing up config to config.2025-01-04-151523.toml...
✓ Backup created

Setting default agent to 'gemini'...
✓ Default agent changed: claude → gemini

Use 'start' to launch with new default.
```

**If no default_agent currently set:**

```
Set default agent
─────────────────────────────────────────────────

Current default: gemini (first in config)

No default_agent configured in [settings].

Select default agent to set explicitly:
  1) gemini (current fallback)
  2) claude
  3) aichat

Select [1-3] (or 'q' to quit): 2

Backing up config to config.2025-01-04-151556.toml...
✓ Backup created

Setting default agent to 'claude'...
✓ Default agent set to 'claude'

Use 'start' to launch with new default.
```

**Quitting without changes:**

```
Select [1-4] (or 'q' to quit): q

Default agent not changed.
```

Exit code: 0

**With agent name (set specific default):**

Sets the default agent in `[settings]` section directly:

```bash
start config agent default gemini
```

Output:

```
Backing up config to config.2025-01-04-151023.toml...
✓ Backup created

Setting default agent to 'gemini'...
✓ Default agent set to 'gemini'

Use 'start' to launch with default agent.
Use 'start config agent default' to confirm.
```

**Updating existing default:**

```bash
start config agent default opus
```

Output:

```
Current default: claude

Backing up config to config.2025-01-04-151056.toml...
✓ Backup created

Setting default agent to 'opus'...
✓ Default agent changed: claude → opus

Use 'start' to launch with new default.
```

**Exit codes:**

- 0 - Success (default shown or set)
- 1 - No agents configured
- 2 - Agent not found
- 3 - File system error (cannot write config, backup failed)

**Error handling:**

**Agent not found:**

```
Error: Agent 'nonexistent' not found in configuration.

Available agents:
  - claude
  - gemini
  - aichat

Use 'start config agent list' for details.
```

Exit code: 2

**No agents configured:**

```
Error: No agents configured.

Use 'start assets add' to install from catalog or 'start config agent new' to create custom.
Use 'start init' to set up agents automatically.
```

Exit code: 1

## Global Flags

These flags work on all `start config agent` subcommands where applicable.

**--help**, **-h**
: Show help for the subcommand.

**--verbose**
: Verbose output. Shows config file paths and additional details.

**--debug**
: Debug mode. Shows all internal operations.

**--version**, **-v**
: Show version information.

## Examples

### List All Agents

```bash
start config agent list
```

Show all configured agents with details.

### List with Verbose Output

```bash
start config agent list --verbose
```

Output:

```
Loading configuration from: ~/.config/start/agents.toml

Configured agents: 3

claude
  Anthropic's Claude AI assistant via Claude Code CLI
  https://docs.claude.com/claude-code
  Default model: claude-3-7-sonnet-20250219 (sonnet)
  Models:
    - claude-3-5-haiku-20241022 (haiku)
    - claude-3-7-sonnet-20250219 (sonnet)
    - claude-opus-4-20250514 (opus)
  Model docs: https://docs.anthropic.com/en/docs/about-claude/models

[... other agents ...]
```

## Files

**~/.config/start/agents.toml**
: Global agent configurations file.

**./.start/agents.toml**
: Local project agent configurations file.

## Error Handling

### No Configuration File

```
Error: No agent configuration found in [scope] config.

Run 'start init' to create initial configuration.
```

Exit code: 1

### Invalid TOML Syntax

```
Error: Configuration file has invalid syntax.

File: ~/.config/start/agents.toml
Line 23: invalid TOML syntax

Fix the configuration file or restore from backup.
```

Exit code: 1

## Notes

### Agent Configuration Scope

Agents can be defined in both global and local configs with merge behavior:

**Global agents:** `~/.config/start/agents.toml`

- Personal agent configurations
- Managed by `start config agent` commands
- Individual preferences (model names, default models)

**Local agents:** `./.start/agents.toml`

- Team-standardized configurations (can be committed to git)
- Managed by `start config agent ... local` commands or manually edited
- Project-specific agent wrappers or custom tools

**Merge behavior:**

- Global + local agents are combined
- Same agent name: local overrides global
- Enables team standardization while allowing personal overrides

### Default Model Behavior

When `default_model` is omitted:

1. Uses first model in `[agents.<name>.models]` table
2. TOML preserves declaration order within tables
3. If no models defined, agent must be used with `--model <full-identifier>`

**Example:**

```toml
[agents.claude]
command = "claude --model {model} '{prompt}'"
# default_model omitted

  [agents.claude.models]
  haiku = "claude-3-5-haiku-20241022"    # First - becomes default
  sonnet = "claude-3-7-sonnet-20250219"
```

Default model: `claude-3-5-haiku-20241022` (haiku)

### Model Names

Model names are agent-specific and user-defined:

- Not hardcoded (no enforced tier names)
- Each agent defines its own model names
- Names can be any meaningful identifier (haiku, sonnet, opus, flash, quick, best, etc.)
- Full model identifiers can always be used with `--model` flag

### Command Template Placeholders

Agent commands support these placeholders:

- `{bin}` - Executable binary name or path (from `bin` field)
- `{model}` - Resolved model name
- `{role}` - The fully resolved role/system prompt content (inline).
- `{role_file}` - File path to role content. Simple roles: original file path. UTD roles: temp file with evaluated content (auto-created and cleaned up).
- `{prompt}` - Built prompt from context documents
- `{date}` - Current timestamp (ISO 8601)

**Example templates:**

```toml
# Placeholder in flag value
command = "{bin} --model {model} '{prompt}'"

# Multiple placeholders
command = "{bin} --model {model} --append-system-prompt '{role}' '{prompt}'"
```

### Metadata URLs

**url** - Agent tool documentation
: Helps users learn about the tool, installation, capabilities

**models_url** - Model documentation
: Helps users understand available models, pricing, context windows, capabilities

Both URLs are optional but recommended for discoverability and self-documentation.

## See Also

- start(1) - Launch with context
- start-init(1) - Initialize configuration
- start-config(1) - Manage configuration files
- start-config-context(1) - Manage context documents
- start-config-task(1) - Manage task configurations
- start-config-role(1) - Manage system prompts
