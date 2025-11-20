# DR-029: Task Agent Field

- Date: 2025-01-07
- Status: Accepted
- Category: Tasks

## Problem

Tasks need a way to specify preferred agents. The design must address:

- Task-specific agent preferences (some tasks work better with specific agents)
- Agent selection precedence (how to combine task preference with user override)
- Validation timing (when to check if agent exists)
- Field design (simple reference vs embedded config)
- Parallel design with role field (consistent patterns)
- Cross-machine sharing (tasks with agents not yet installed)
- Default behavior (what happens when no agent specified)
- Override capability (user must be able to override task's choice)

## Decision

Tasks can specify a preferred agent using an optional agent field that references an agent name.

Agent field specification:

- Field name: agent (string, optional)
- Value: Name reference to agent defined in [agents.<name>]
- Example: agent = "go-expert"
- Must match existing agent name
- Validation at execution time (not load time)

Agent selection precedence (highest to lowest):

1. CLI flag: --agent flag (explicit user override)
2. Task agent: agent field in task configuration
3. Default agent: default_agent from [settings] section
4. First agent: First agent in config (TOML order)

Validation timing:

- Execution time: Error if agent not found when task runs
- Doctor check: Warns about undefined agents in tasks
- Config validate: Reports tasks with undefined agents
- Not at load time: Allows sharing tasks across machines

## Why

Tasks need agent preferences for optimization:

- Specialized agents for specific workflows (go-expert for Go code)
- Different model perspectives (get second opinion from different model)
- Performance optimization (fast agent for quick checks)
- Tool-specific features (vision models for image review)

Optional field keeps tasks simple:

- Most tasks work with any agent
- Simple tasks don't need agent specification
- Default agent is usually sufficient
- Only specify when needed

Execution-time validation enables sharing:

- Task configs can be shared across machines
- Different machines may have different agents configured
- Task config valid even if agent not yet installed
- Doctor and validate can still catch issues proactively

Simple string reference keeps config clean:

- Agent name reference is sufficient (not embedded config)
- All agent configuration lives in [agents.<name>] section
- Keeps task config simple and focused
- Similar to how role field references role names

Parallel design with role field provides consistency:

- Tasks can specify both agent and role preferences
- Both follow same pattern: string reference to named entity
- Agent controls execution tool, role controls AI persona
- Both can be overridden via CLI flags (--agent, --role)
- Contexts remain global (required contexts auto-included)

CLI override maintains user control:

- User can always override task's agent choice
- Explicit --agent flag takes highest precedence
- Useful for experimentation or temporary changes
- Task preference is suggestion, not requirement

## Trade-offs

Accept:

- Execution-time validation (task may fail if agent not installed)
- No per-task agent configuration (can't embed agent details in task)
- Agent reference can break (if agent removed from config)
- No fallback chain (single agent choice, not multiple options)
- Tasks may specify agents users don't have (shareable configs)

Gain:

- Task-specific agent optimization (better results for specialized workflows)
- Simple string reference (clean config, no duplication)
- User override capability (--agent flag always wins)
- Shareable task configs (work across machines with different agents)
- Consistent with role field design (parallel patterns)
- Proactive validation available (doctor and validate catch issues)
- Default behavior clear (precedence order well-defined)

## Alternatives

Agent field in settings section only:

Example: Only default_agent in settings, no per-task agent

```toml
[settings]
default_agent = "claude"
```

Pros:

- Simpler configuration (one place to set agent)
- Consistent agent usage across all tasks
- No per-task complexity

Cons:

- No per-task customization or optimization
- Can't use specialized agents for specific workflows
- Must use --agent flag for every specialized task
- Less flexible

Rejected: Per-task optimization is valuable. Some tasks genuinely work better with specific agents.

Load-time validation:

Example: Validate agent exists when loading task config

- Fail to load tasks.toml if agent reference invalid
- Error message immediately on start

Pros:

- Immediate feedback about configuration problems
- Clear error before attempting to use task

Cons:

- Prevents sharing task configs across machines
- Breaks if agent not yet installed (even if not using that task)
- Must install all referenced agents before config loads
- Less flexible for team configurations

Rejected: Execution-time validation enables sharing. Better to fail at use with clear error than block all config loading.

Complex agent configuration in task:

Example: Embed agent details directly in task

```toml
[tasks.go-review]
agent.bin = "claude"
agent.command = "claude --model sonnet '{prompt}'"
agent.models.sonnet = "claude-3-7-sonnet-20250219"
```

Pros:

- Self-contained task definition
- No dependency on agent configuration
- Explicit about what agent is used

Cons:

- Duplicates agent config in multiple places
- Hard to maintain (change agent = update all tasks)
- Verbose and complex task definitions
- Breaks separation of concerns

Rejected: Keep agent details in [agents.<name>] section. Tasks just reference by name (DRY principle).

Multiple agents with fallback chain:

Example: Specify primary and fallback agents

```toml
[tasks.go-review]
agents = ["go-expert", "claude", "gemini"]  # Try in order
```

Pros:

- Automatic fallback if preferred agent not available
- More resilient to missing agents
- Tasks can work across different machine configs

Cons:

- Over-engineered for current needs
- Adds complexity to precedence rules
- Unclear which agent actually ran
- Makes debugging harder

Rejected: Single agent + CLI override is sufficient. Can add later if users request fallback chains.

## Structure

Agent field in task configuration:

Field specification:

- Field name: agent
- Type: string (optional)
- Value: Name reference to [agents.<name>] section
- Example: agent = "go-expert"
- Validation: Must match existing agent name
- When validated: Execution time, doctor, config validate

Complete task example:

```toml
[tasks.go-review]
agent = "go-expert"              # Optional: Preferred agent
role = "go-reviewer"             # Optional: Preferred role
alias = "gor"
description = "Review Go code with specialized agent"
command = "git diff --staged"
prompt = "Review this Go code: {command_output}\n\n{instructions}"
```

Agent selection precedence:

Priority order (highest to lowest):

1. CLI flag: start task go-review --agent gemini
   - Explicit user override
   - Always takes precedence
   - Temporary for this execution

2. Task agent: agent = "go-expert" in task config
   - Task's preferred agent
   - Used if no CLI flag
   - Persistent preference

3. Default agent: default_agent = "claude" in [settings]
   - Global default for all tasks
   - Used if task has no agent field
   - Fallback when no task preference

4. First agent: First [agents.<name>] in config
   - TOML declaration order
   - Ultimate fallback
   - Ensures agent always selected

Scope and merge behavior:

Task agent field follows standard task merge:

- Local task completely overrides global task (same name)
- Entire task replaced, including agent field
- No field-level merging between global and local

Example:

```toml
# Global: ~/.config/start/tasks.toml
[tasks.code-review]
agent = "claude"
prompt = "Review code: {instructions}"

# Local: ./.start/tasks.toml
[tasks.code-review]
agent = "go-expert"
prompt = "Review Go code: {instructions}"
```

Result: Local task completely replaces global (agent is "go-expert")

Validation behavior:

Execution time (start task go-review):

- Check if agent exists when task runs
- Error if agent not found
- Exit code: 2 (agent error)

Doctor check (start doctor):

- Validates all tasks with agent field
- Reports undefined agents with fix suggestions
- Exit code: 1 if validation errors

Config validate (start config validate):

- Reports all tasks with undefined agents
- Shows which agents are missing
- Exit code: 1 if validation errors

## Usage Examples

Basic task with agent preference:

```toml
[tasks.go-review]
agent = "go-expert"
description = "Review Go code with specialized agent"
command = "git diff --staged"
prompt = "Review: {command_output}\n\n{instructions}"
```

```bash
# Uses go-expert (from task)
start task go-review

# Uses gemini (CLI flag overrides task)
start task go-review --agent gemini
```

Multiple tasks with different agents:

```toml
[tasks.quick-check]
agent = "haiku-agent"
description = "Fast code check with lightweight agent"
prompt = "Quick review: {instructions}"

[tasks.deep-review]
agent = "opus-agent"
description = "Thorough code review with advanced agent"
prompt = "Comprehensive review: {instructions}"

[tasks.alternative-view]
agent = "gemini"
description = "Get second opinion from different model"
prompt = "Alternative perspective: {instructions}"
```

Agent selection precedence examples:

```toml
# settings.toml
[settings]
default_agent = "claude"

# tasks.toml
[tasks.go-review]
agent = "go-expert"

[tasks.code-review]
# No agent field
```

```bash
# Uses go-expert (from task)
start task go-review

# Uses gemini (CLI flag overrides task)
start task go-review --agent gemini

# Uses claude (task has no agent field, falls to default_agent)
start task code-review

# Uses first agent in config (no CLI flag, no task agent, no default_agent)
start task simple-task
```

Error handling - agent not found at execution:

```bash
$ start task go-review

Error: Agent 'go-expert' not found (required by task 'go-review').

Configured agents:
  claude
  opencode

Add agent: start assets update
Or override: start task go-review --agent claude
```

Exit code: 2

Error handling - doctor check:

```bash
$ start doctor

Configuration:
  Global Config:
    settings.toml: ✓
    agents.toml:   ✓
    tasks.toml:    ✓
    roles.toml:    ✓
    contexts.toml: ✓
  Validation:      ✗ Issues found

Configuration Issues:
  ✗ Task 'go-review' references undefined agent 'go-expert'
    Fix: start assets update
    Or: Remove agent field from task configuration

Overall Status:   ✗ Critical issues found
```

Exit code: 1

Error handling - config validate:

```bash
$ start config validate

Validation errors:

Tasks:
  ✗ go-review: Agent 'go-expert' not found in configuration
  ✗ security-scan: Agent 'security-bot' not found in configuration

Fix: Add agents or remove agent fields from tasks
```

Exit code: 1

Use case - specialized agent:

```toml
[tasks.go-review]
agent = "go-expert"
description = "Review Go code with specialized agent"
```

Use case - performance optimization:

```toml
[tasks.quick-check]
agent = "haiku-agent"
description = "Fast code check with lightweight agent"
```

Use case - tool-specific features:

```toml
[tasks.visual-review]
agent = "claude-with-vision"
description = "Review that may involve images/diagrams"
```

## Updates

- 2025-01-17: Initial version aligned with schema
