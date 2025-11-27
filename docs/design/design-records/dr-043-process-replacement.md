# DR-043: Process Replacement Execution Model

- Date: 2025-11-25
- Status: Accepted
- Category: Runtime Behavior

## Problem

The `start` tool needs to execute AI agent CLI tools (claude, gemini, aichat, etc.) after assembling prompts and resolving placeholders. The execution model must:

- Launch the AI agent with the final command
- Allow the AI agent to interact directly with the user's terminal
- Not interfere with streaming responses from AI agents
- Exit cleanly when the AI agent exits
- Be simple and efficient

## Decision

Use **process replacement** via `syscall.Exec` to replace the `start` process with the AI agent process.

The `start` tool:
1. Loads configuration
2. Resolves placeholders
3. Assembles the final command string
4. Replaces itself with the AI agent process
5. The AI agent takes over the terminal completely

After step 4, the `start` process no longer exists - the AI agent has replaced it.

## Why

**`start` is a launcher, not a wrapper:**

- The tool's purpose is to prepare context and launch agents, not manage them
- Once the agent is launched, `start` has no more work to do
- The AI agent should own the terminal interaction completely

**Process replacement advantages:**

- AI agent inherits stdin/stdout/stderr automatically - no stream management needed
- No parent process overhead - `start` doesn't sit idle while agent runs
- No buffering issues - agent controls all output
- No timeout management - agent runs as long as needed
- Simpler code - no subprocess management logic
- Standard Unix pattern - same as shell `exec` builtin

**User experience:**

- User sees AI agent responses in real-time (streaming)
- AI agent can use interactive features (prompts, colors, etc.)
- Process tree is clean - only the AI agent is running
- Terminal signals (Ctrl+C, etc.) go directly to AI agent

## Trade-offs

Accept:

- Cannot capture agent output for logging/analysis
- Cannot implement timeouts or resource limits on agent
- Cannot clean up or perform actions after agent exits
- Process replacement is platform-specific (Unix syscall.Exec)

Gain:

- Simpler implementation - no subprocess management
- Better user experience - direct terminal control
- No buffering or streaming issues
- No idle parent process
- Minimal overhead

## Implementation

**Runner Interface:**

```go
// Runner abstracts command execution
type Runner interface {
    // Exec replaces the current process with the command
    // After successful exec, this function never returns
    // Only returns on error (before exec)
    Exec(shell, command string) error
}
```

**RealRunner Implementation:**

```go
type RealRunner struct{}

func (r *RealRunner) Exec(shell, command string) error {
    // Find shell binary
    shellPath, err := exec.LookPath(shell)
    if err != nil {
        return fmt.Errorf("shell not found: %w", err)
    }

    // Replace current process with shell running command
    // This never returns on success
    return syscall.Exec(shellPath, []string{shell, "-c", command}, os.Environ())
}
```

**Executor:**

```go
func (e *Executor) Execute(agent domain.Agent, model, prompt string) error {
    // Resolve placeholders
    values := map[string]string{
        "bin":    agent.Bin,
        "model":  model,
        "prompt": prompt,
    }
    command := e.resolver.Resolve(agent.Command, values)

    // Get shell from settings (or default)
    shell := e.getShell()

    // Replace process with agent
    // This never returns on success
    return e.runner.Exec(shell, command)
}
```

**Testing:**

- Unit tests use MockRunner that simulates exec without actually replacing process
- Integration tests verify command construction but cannot test actual replacement
- Manual testing required to verify real execution behavior

## Windows Compatibility

`syscall.Exec` is Unix-specific. For Windows:

- Option A: Use `cmd.Run()` and wait (subprocess model)
- Option B: Use `syscall.StartProcess` with creative process management
- Option C: Document as Unix-only feature initially

Decision: Start with Unix-only (Option C), add Windows support in later phase if needed.

## Executor Changes

**Before (Phase 2 - buffered subprocess):**

```go
type Runner interface {
    Run(ctx context.Context, shell, command string, timeout time.Duration) (string, string, error)
}

func (e *Executor) Execute(ctx context.Context, agent domain.Agent, model, prompt string) error {
    // ... resolve placeholders ...
    stdout, stderr, err := e.runner.Run(ctx, "bash", command, 2*time.Minute)
    fmt.Print(stdout)
    fmt.Print(stderr)
    return err
}
```

**After (Phase 2.5 - process replacement):**

```go
type Runner interface {
    Exec(shell, command string) error
}

func (e *Executor) Execute(agent domain.Agent, model, prompt, shell string) error {
    // ... resolve placeholders ...
    return e.runner.Exec(shell, command)
    // Never reaches here on success
}
```

## Context Removal

The `context.Context` parameter is removed from execution:

- Process replacement happens immediately
- No timeout management needed
- No cancellation needed
- Context was only useful for subprocess management

If graceful shutdown is needed in the future, it would be implemented differently (signal handling, not context).

## Settings Integration

The executor needs access to `Settings.Shell`:

```go
func NewExecutor(runner domain.Runner, resolver *PlaceholderResolver) *Executor {
    return &Executor{
        runner:   runner,
        resolver: resolver,
    }
}

func (e *Executor) Execute(agent domain.Agent, model, prompt, shell string) error {
    // Shell passed as parameter from settings
}
```

## Alternatives

**Subprocess with stream pass-through:**

```go
cmd := exec.Command(shell, "-c", command)
cmd.Stdin = os.Stdin
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
return cmd.Run()
```

- Pro: Cross-platform (works on Windows)
- Pro: Can perform cleanup after agent exits
- Pro: Can implement timeouts/resource limits
- Con: Parent process sits idle while agent runs
- Con: Extra process in process tree
- Con: Slightly more complex signal handling
- Rejected: Unnecessary complexity for a launcher

**Subprocess with output capture:**

- Pro: Can log agent output
- Pro: Can parse agent responses
- Con: Breaks streaming for interactive agents
- Con: Memory issues with large outputs
- Con: Wrong model - we're a launcher, not a wrapper
- Rejected: This is what Phase 2 had, and it's wrong for a launcher

**Background launch (fork and exit):**

```go
cmd := exec.Command(shell, "-c", command)
cmd.Start()
// Exit immediately without waiting
```

- Pro: `start` exits immediately
- Con: Agent becomes orphaned/detached
- Con: No way to return agent exit code
- Con: Confusing for user (agent still running after command "finishes")
- Rejected: Wrong semantics for CLI tool

## Breaking Changes

From Phase 2 implementation:

1. `Runner.Run()` → `Runner.Exec()` (different signature)
2. `context.Context` parameter removed from execution path
3. `timeout` parameter removed (no longer applicable)
4. Return values changed (no stdout/stderr strings)
5. Executor no longer handles output (agent controls terminal)

## Updates

- 2025-11-25: Initial decision for process replacement execution model
