# DR-029: Task Agent Field

**Date:** 2025-01-07
**Status:** Accepted
**Category:** Tasks

## Decision

Tasks can specify a preferred agent using an optional `agent` field.

## Configuration

```toml
[tasks.go-review]
agent = "go-expert"              # Optional: Preferred agent for this task
role = "go-reviewer"             # Optional: Preferred role for this task
alias = "gor"
description = "Review Go code with specialized agent"

command = "git diff --staged"
prompt = "Review this Go code: {command}\n\n{instructions}"
```

## Agent Selection Precedence

When executing a task, agent selection follows this priority order:

1. **CLI flag** - `--agent` flag (highest priority, explicit user override)
2. **Task agent** - `agent` field in task configuration
3. **Default agent** - `default_agent` from `[settings]` section
4. **First agent** - First agent in config (TOML order)

**Example:**

```toml
# config.toml
[settings]
default_agent = "claude"

[tasks.go-review]
agent = "go-expert"
# ...
```

```bash
# Uses go-expert (from task)
start task go-review

# Uses gemini (CLI flag overrides task)
start task go-review --agent gemini

# Uses claude (task has no agent field, falls to default_agent)
start task code-review

# Uses first agent in config (no CLI flag, no task agent, no default_agent)
# (First agent defined in config TOML order)
start task simple-task
```

## Field Specification

**agent** (string, optional)
: The name of an agent defined in `[agents.<name>]` configuration. Must match an existing agent name.

**Validation:**
- Field is optional (tasks without it use default agent)
- Value must reference an existing agent name
- Validation occurs at task execution time (not load time)
- Validation also performed by `start doctor` and `start config validate`

## Error Handling

### Agent Not Found (Execution Time)

```
Error: Agent 'go-expert' not found (required by task 'go-review').

Configured agents:
  claude
  opencode

Add agent: start config agent add go-expert
Or override: start task go-review --agent claude
```

Exit code: 2

### Agent Not Found (Doctor Check)

```bash
start doctor
```

```
Configuration Issues:
  ✗ Task 'go-review' references undefined agent 'go-expert'
    Fix: start config agent add go-expert
    Or: Remove agent field from task configuration
```

Exit code: 1

### Agent Not Found (Config Validate)

```bash
start config validate
```

```
Validation errors:

Tasks:
  ✗ go-review: Agent 'go-expert' not found in configuration
  ✗ security-scan: Agent 'security-bot' not found in configuration

Fix: Add agents or remove agent fields from tasks
```

Exit code: 1

## Use Cases

**1. Specialized Agents**

```toml
[tasks.go-review]
agent = "go-expert"
description = "Review Go code with specialized agent"
```

**2. Different Model Perspectives**

```toml
[tasks.alternative-review]
agent = "gemini"
description = "Get second opinion from different model"
```

**3. Performance Optimization**

```toml
[tasks.quick-check]
agent = "haiku-agent"
description = "Fast code check with lightweight agent"
```

**4. Tool-Specific Features**

```toml
[tasks.visual-review]
agent = "claude-with-vision"
description = "Review that may involve images/diagrams"
```

## Scope and Merge Behavior

Task agent field follows standard task merge behavior:

- Tasks defined in local config completely override global tasks
- If local task defines same name, entire task replaced (including agent)
- No field-level merging between global and local

**Example:**

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

## Implementation Notes

**Task Execution Flow:**

1. Load and merge configuration (global + local)
2. Find task by name or alias
3. Determine agent:
   - Check for `--agent` CLI flag → use if present
   - Otherwise check task `agent` field → use if present
   - Otherwise use `default_agent` from settings
4. Validate agent exists in configuration → error if not found
5. Proceed with task execution using selected agent

**Doctor/Validate Integration:**

- Check all tasks with `agent` field
- Verify each agent name exists in merged `[agents.<name>]` config
- Report errors with actionable fix suggestions

## Rationale

**Why optional?**
- Most tasks work with any agent
- Simple tasks don't need agent specification
- Default agent is usually sufficient

**Why execution-time validation?**
- Allows task configs to be shared across machines
- Different machines may have different agents configured
- Task config valid even if agent not yet installed
- Doctor/validate can still catch issues proactively

**Why simple string field?**
- Agent name reference is sufficient
- All agent configuration lives in `[agents.<name>]` section
- Keeps task config simple and focused
- Similar to how `role` field references role names (DR-005)

**Why agent and role fields parallel?**
- Tasks can specify both agent and role preferences
- Both follow same pattern: string reference to named entity
- Agent controls execution tool, role controls AI persona
- Both can be overridden via CLI flags (--agent, --role)
- Contexts remain global (required contexts auto-included per DR-012)

## Alternatives Considered

**1. Agent field in settings section only**
- Rejected: No per-task customization
- Can't optimize specific workflows for specific agents

**2. Load-time validation**
- Rejected: Prevents sharing task configs across machines
- Better to fail at execution with clear error message

**3. Complex agent configuration in task**
- Rejected: Would duplicate agent config in multiple places
- Keep agent details in `[agents.<name>]`, tasks just reference by name

**4. Allow multiple agents (fallback chain)**
- Rejected: Over-engineered for current needs
- Single agent + CLI override is sufficient
- Can add later if needed

## Related Decisions

- [DR-005](./dr-005-role-configuration.md) - Role configuration (parallel to agent field)
- [DR-009](./dr-009-task-structure.md) - Task structure and placeholders (includes role field)
- [DR-004](./dr-004-agent-scope.md) - Agent configuration scope
- [DR-019](./dr-019-task-loading.md) - Task loading and merging

## Documentation Updates Required

- [x] Create DR-029
- [x] Update `docs/tasks.md` - Add `agent` field to configuration section
- [x] Update `docs/cli/start-task.md` - Document agent selection precedence
- [x] Update `docs/cli/start-doctor.md` - Add task agent validation check
- [x] Update `docs/config.md` - Add `agent` field to tasks section
