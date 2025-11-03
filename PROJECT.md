# Project: start - CLI Design Phase

**Status:** Design Phase - Command-Line Interface Specification
**Date Started:** 2025-01-03
**Current Focus:** Defining all CLI commands, subcommands, arguments, and flags

## Overview

`start` is a context-aware AI agent launcher that detects project context, builds intelligent prompts, and launches AI development tools with proper configuration.

**Vision:** See [docs/vision.md](./docs/vision.md)
**Design Decisions:** See [docs/design-record.md](./docs/design-record.md) (11 decisions locked in)

## What's Been Decided

### Architecture (Complete)

- ✅ Configuration format: TOML
- ✅ Config structure: Global (~/.config/start/config.toml) + Local (./.start/config.toml) merge
- ✅ Named documents for context (override/add pattern)
- ✅ Agents: Global only, named by actual tool (claude, gemini, etc.)
- ✅ System prompt: Separate from context documents, optional
- ✅ Placeholders: {model}, {system_prompt}, {prompt}, {date}, {instructions}, {content}
- ✅ Tasks: Role + prompt template + optional content_command + documents
- ✅ Default tasks: 4 interactive review tasks (cr, gdr, ct, dr)
- ✅ Asset distribution: Embedded in binary via go:embed
- ✅ Init wizard: Interactive agent setup with multi-select

### CLI Framework (Decided)

- ✅ Using Cobra for command structure
- ✅ Pattern: `start <subcommand> [args] [flags]`
- ✅ Persistent flags work across all commands
- ✅ Tasks loaded dynamically from config

## What Needs to Be Designed

### CLI Command Specifications

The following commands and subcommands need complete specification:

#### 1. Root Command: `start`

**File:** `docs/cli/start.md`

**Needs Definition:**

- Can user pass custom prompt as argument? `start "custom prompt here"`
- Or is custom prompt a flag? `start --prompt "..."`
- How does it interact with context documents?
- What does verbose/debug output look like?
- Error handling (no config, no agents, etc.)

**Global Flags:**

- `--agent <name>` - Which agent to use (override default)
- `--model <tier>` - Model tier (fast/mid/pro) OR specific model name?
- `--directory <path>` - Working directory (default: pwd)
- `--verbose` / `-v` - Verbose output?
- `--debug` - Debug mode?
- `--quiet` / `-q` - Suppress non-essential output?
- `--help` / `-h` - Help text
- `--version` - Version info

**Open Questions:**

1. Does `--model` accept tier names (fast/mid/pro) OR full model names OR both?
2. Custom prompt: argument or flag?
3. What flags are truly global vs command-specific?
4. Output verbosity levels (normal/verbose/debug/quiet)?

#### 2. Init Command: `start init`

**File:** `docs/cli/start-init.md`

**Needs Definition:**

- Interactive wizard flow (exact prompts)
- Agent selection (multi-select UI)
- Model configuration per agent
- Document detection and setup
- Backup behavior details
- Force/non-interactive mode flags?

**Possible Flags:**

- `--force` / `-f` - Overwrite without backup?
- `--non-interactive` - Skip wizard, use defaults?
- `--agent <name>` - Pre-select agent(s)?

**Open Questions:**

1. Can you re-run init to add agents without losing config?
2. Should there be `start init --minimal` for quick setup?
3. How to handle partial failures (e.g., agent tool not found)?

#### 3. Task Command: `start task <name> [instructions]`

**File:** `docs/cli/start-task.md`

**Needs Definition:**

- How to list available tasks
- Task help text format
- Instructions argument handling (quoted strings, multiple args?)
- Task-specific flags?
- What if content_command fails?

**Subcommands Needed?**

- `start task list` - List all tasks
- `start task info <name>` - Show task details
- OR just `start task --help` shows all?

**Possible Flags:**

- `--list` - List available tasks
- `--info <name>` - Show task configuration
- Plus global flags (--agent, --model)

**Open Questions:**

1. Are task subcommands needed or just flags?
2. How are multi-word instructions passed? `start task gdr "word1 word2"` or `start task gdr word1 word2`?
3. Should task show what content_command will run before executing?

#### 4. Agent Management: `start agent <subcommand>`

**File:** `docs/cli/start-agent.md`

**Subcommands Needed:**

##### 4.1. `start agent add`

**File:** `docs/cli/start-agent-add.md`

**Needs Definition:**

- Interactive wizard vs flags
- Custom agent configuration
- Model tier setup
- Command template input
- Environment variable setup

**Possible Flags:**

- `--name <name>` - Agent name
- `--command <template>` - Command template
- `--model-fast <model>` - Fast tier model
- `--model-mid <model>` - Mid tier model
- `--model-pro <model>` - Pro tier model
- `--interactive` - Force wizard mode

**Open Questions:**

1. Support both wizard and flag-based addition?
2. How to configure agent.env in CLI?
3. Validation before saving?

##### 4.2. `start agent list`

**File:** `docs/cli/start-agent-list.md`

**Needs Definition:**

- Output format (table, list, JSON?)
- Show which is default?
- Show model tiers?
- Indicate if tool is installed?

**Possible Flags:**

- `--json` - JSON output
- `--verbose` - Show full config

##### 4.3. `start agent test <name>`

**File:** `docs/cli/start-agent-test.md`

**Needs Definition:**

- What does "test" mean?
  - Check if binary exists?
  - Run with test prompt?
  - Validate command template?
  - Check model availability?

**Open Questions:**

1. What constitutes a successful test?
2. Test all models or just default?

##### 4.4. `start agent remove <name>`

**File:** `docs/cli/start-agent-remove.md`

**Needs Definition:**

- Confirmation required?
- What if it's the default agent?
- Cascade delete related config?

**Possible Flags:**

- `--force` / `-f` - Skip confirmation

##### 4.5. `start agent edit <name>`

**File:** `docs/cli/start-agent-edit.md`

**Needs Definition:**

- Open in $EDITOR?
- Interactive update wizard?
- Specific field updates?

**Possible Flags:**

- `--field <name>` - Update specific field
- `--editor` - Open in editor

#### 5. Config Management: `start config <subcommand>`

**File:** `docs/cli/start-config.md`

**Subcommands Needed:**

##### 5.1. `start config show`

**File:** `docs/cli/start-config-show.md`

**Needs Definition:**

- Show merged config or just user config?
- Output format (TOML, JSON, pretty-print?)
- Show all sections or specific?

**Possible Flags:**

- `--global` - Show only global config
- `--local` - Show only local config
- `--merged` - Show merged result (default?)
- `--json` - JSON output
- `--section <name>` - Show specific section (agents, context, tasks)

##### 5.2. `start config edit`

**File:** `docs/cli/start-config-edit.md`

**Needs Definition:**

- Which config to edit (global vs local)?
- What if $EDITOR not set?
- Validate on save?

**Possible Flags:**

- `--global` - Edit global config (default?)
- `--local` - Edit local config
- `--editor <cmd>` - Specify editor

##### 5.3. `start config path`

**File:** `docs/cli/start-config-path.md`

**Needs Definition:**

- Print which path? Global, local, both?
- Check if exists?

**Possible Flags:**

- `--global` - Global config path
- `--local` - Local config path

##### 5.4. `start config validate`

**File:** `docs/cli/start-config-validate.md`

**Needs Definition:**

- What validations to run?
  - TOML syntax
  - Required fields
  - Agent command templates
  - File paths exist?
  - Placeholder usage
- Output format (errors, warnings)

**Possible Flags:**

- `--strict` - Error on warnings
- `--json` - JSON output

##### 5.5. `start config init`

**File:** Should this be separate from `start init`?

**Open Questions:**

1. Is `start config init` redundant with `start init`?
2. Should `start init` be an alias for `start config init`?

#### 6. Context/Document Management

**Should this exist?**

Possible commands:

- `start context list` - Show configured documents
- `start context add <path>` - Add document to config
- `start context remove <name>` - Remove document
- `start context test` - Check which documents exist

**Open Questions:**

1. Is this needed or just edit config manually?
2. Global vs local context management?

#### 7. Role Management

**Should this exist?**

Possible commands:

- `start role list` - List available role templates
- `start role create` - Interactive role creation
- `start role edit <name>` - Edit role file

**Open Questions:**

1. Is this needed or just file operations?
2. Role templates in assets?

## Command Documentation Structure

Each `docs/cli/<command>.md` file should include:

### Standard Sections

1. **Name** - Command name and brief description
2. **Synopsis** - Usage pattern with all options
3. **Description** - Detailed explanation of what it does
4. **Arguments** - Positional arguments (required/optional)
5. **Flags** - All flags with descriptions
6. **Examples** - Common usage examples
7. **Output** - What output looks like (with examples)
8. **Exit Codes** - 0 = success, 1 = error, etc.
9. **See Also** - Related commands

### Example Template

```markdown
# start agent list

## Name

start agent list - List configured AI agents

## Synopsis

start agent list [flags]

## Description

Lists all configured AI agents with their details...

## Flags

--json
Output in JSON format

--verbose, -v
Show full agent configuration

## Examples

# List all agents

start agent list

# JSON output

start agent list --json

## Output

[Example output here]

## Exit Codes

0 - Success
1 - Configuration error

## See Also

- start-agent-add(1)
- start-agent-remove(1)
```

## Design Approach

Work through commands in this order:

1. **Root command first** (`start`) - This affects everything else

   - Define global flags
   - Custom prompt handling
   - Output format

2. **Core workflow commands**

   - `start init` - First-run experience
   - `start task` - Primary use case

3. **Management commands**

   - `start agent *` - Agent CRUD
   - `start config *` - Config management

4. **Optional/Advanced commands**
   - Context management (if needed)
   - Role management (if needed)

## Open Design Questions

### High Priority

1. **Custom prompts:** How does user override prompt for one-off use?
2. **Model selection:** Accept tiers (fast/mid/pro) OR full model names OR both?
3. **Task listing:** Subcommand or flag?
4. **Agent testing:** What does "test" actually do?
5. **Config editing:** How to handle validation and errors?

### Medium Priority

6. **Verbosity levels:** How many? (quiet/normal/verbose/debug)
7. **Output formats:** Where do we support --json?
8. **Non-interactive mode:** For CI/automation, what flags needed?
9. **Context management:** Build it or skip it?
10. **Role management:** Build it or skip it?

### Low Priority

11. **Shell completion:** Generate for bash/zsh/fish?
12. **Aliases:** Command aliases support?
13. **Config profiles:** Multiple named configs?

## Success Criteria

CLI design is complete when:

- [ ] Every command has a complete specification document
- [ ] All flags are defined with types and defaults
- [ ] All open questions are resolved
- [ ] Examples cover common use cases
- [ ] Error cases are documented
- [ ] Output formats are specified
- [ ] Validation rules are clear

## Next Steps

1. **Start with root command:** Define `start` command completely

   - Resolve custom prompt handling
   - Finalize global flags
   - Document output format

2. **Create template:** Use first command as template for others

3. **Work through subcommands:** One at a time, document completely

4. **Review consistency:** Ensure patterns consistent across all commands

5. **Update design-record.md:** Record final CLI decisions

## Reference Documents

- [Vision](./docs/vision.md) - Product vision and goals
- [Design Record](./docs/design-record.md) - All design decisions (DR-001 through DR-011)
- [Task Documentation](./docs/task.md) - Task configuration details
- [Thoughts](./docs/thoughts.md) - Design ideas and considerations

## Notes for Next Session

- Focus on one command at a time
- Ask questions when ambiguous
- Provide examples for each decision
- Consider both interactive and scripting use cases
- Think about error messages and user experience
- Consider what Cobra makes easy vs hard
