package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/grantcarthew/start/internal/config"
	"github.com/grantcarthew/start/internal/domain"
	"github.com/grantcarthew/start/internal/engine"
	"github.com/spf13/cobra"
)

// TaskCommand holds the dependencies for the task command
type TaskCommand struct {
	configLoader   *config.Loader
	validator      *config.Validator
	executor       *engine.Executor
	roleSelector   *engine.RoleSelector
	roleLoader     *engine.RoleLoader
	contextLoader  *engine.ContextLoader
	taskLoader     *engine.TaskLoader
	taskResolver   *engine.TaskResolver
}

// NewTaskCommand creates the task command
func NewTaskCommand(
	configLoader *config.Loader,
	validator *config.Validator,
	executor *engine.Executor,
	roleSelector *engine.RoleSelector,
	roleLoader *engine.RoleLoader,
	contextLoader *engine.ContextLoader,
	taskLoader *engine.TaskLoader,
	taskResolver *engine.TaskResolver,
) *cobra.Command {
	tc := &TaskCommand{
		configLoader:   configLoader,
		validator:      validator,
		executor:       executor,
		roleSelector:   roleSelector,
		roleLoader:     roleLoader,
		contextLoader:  contextLoader,
		taskLoader:     taskLoader,
		taskResolver:   taskResolver,
	}

	cmd := &cobra.Command{
		Use:   "task [name] [instructions]",
		Short: "Run predefined AI workflow tasks",
		Long: `Executes predefined AI workflow tasks configured in tasks.toml.

Tasks are reusable workflows with optional role overrides, automatic required context
inclusion, and dynamic content from shell commands.

Examples:
  start task                              # List all tasks
  start task code-review                  # Run task
  start task gdr "focus on security"      # Run with instructions
  start task code-review --agent gemini   # Override agent`,
		RunE: tc.run,
		Args: cobra.ArbitraryArgs,
	}

	return cmd
}

// run executes the task command
func (tc *TaskCommand) run(cmd *cobra.Command, args []string) error {
	// Load configuration
	globalCfg, err := tc.configLoader.LoadGlobal()
	if err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}

	// Get current working directory for local config
	workDir := "."
	localCfg, err := tc.configLoader.LoadLocal(workDir)
	if err != nil {
		// Local config is optional, use empty config
		localCfg = globalCfg
	}

	// Merge configs
	cfg := config.Merge(globalCfg, localCfg)

	// Validate merged config
	if err := tc.validator.Validate(cfg); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	// If no arguments, list tasks
	if len(args) == 0 {
		return tc.listTasks(cfg, globalCfg, localCfg)
	}

	// Resolve task
	taskName := args[0]
	task, err := tc.taskResolver.Resolve(taskName, localCfg.Tasks, globalCfg.Tasks)
	if err != nil {
		return tc.taskNotFoundError(taskName, cfg)
	}

	// Get instructions (remaining args joined)
	instructions := ""
	if len(args) > 1 {
		instructions = strings.Join(args[1:], " ")
	}

	// Get flags
	agentFlag, _ := cmd.Flags().GetString("agent")
	modelFlag, _ := cmd.Flags().GetString("model")
	roleFlag, _ := cmd.Flags().GetString("role")

	// Select agent with precedence: --agent flag > task agent > default_agent > first agent
	var agentName string
	if agentFlag != "" {
		agentName = agentFlag
	} else if task.Agent != "" {
		agentName = task.Agent
	} else if cfg.Settings.DefaultAgent != "" {
		agentName = cfg.Settings.DefaultAgent
	} else {
		return fmt.Errorf("no agent specified and no default agent configured")
	}

	// Get agent from config
	agent, ok := cfg.Agents[agentName]
	if !ok {
		return fmt.Errorf("agent %q not found in configuration", agentName)
	}
	agent.Name = agentName

	// Select model
	var modelID string
	if modelFlag != "" {
		// Check if it's a model name (needs resolution) or full ID
		if fullID, ok := agent.Models[modelFlag]; ok {
			modelID = fullID
		} else {
			// Assume it's a full model ID
			modelID = modelFlag
		}
	} else if agent.DefaultModel != "" {
		// Use default model
		if fullID, ok := agent.Models[agent.DefaultModel]; ok {
			modelID = fullID
		} else {
			return fmt.Errorf("default model %q not found in agent %q models", agent.DefaultModel, agentName)
		}
	} else {
		return fmt.Errorf("no model specified and no default model for agent %q", agentName)
	}

	// Select role with precedence: --role flag > task role > default_role > first role
	selectionCtx := engine.SelectionContext{
		RoleFlag:    roleFlag,
		TaskRole:    task.Role,
		DefaultRole: cfg.Settings.DefaultRole,
	}
	role, err := tc.roleSelector.Select(selectionCtx, cfg.Roles)
	if err != nil {
		return fmt.Errorf("role selection failed: %w", err)
	}

	// Get shell and timeout settings
	shell := cfg.Settings.Shell
	if shell == "" {
		shell = engine.DetectShell()
	}
	timeout := cfg.Settings.CommandTimeout
	if timeout == 0 {
		timeout = 30 // Default 30 seconds
	}

	// Load role
	loadedRole, err := tc.roleLoader.LoadRole(role, shell, timeout)
	if err != nil {
		return fmt.Errorf("failed to load role: %w", err)
	}
	// Cleanup temp role file if needed (deferred)
	defer tc.roleLoader.CleanupRole(loadedRole)

	// Load required contexts only (tasks use CommandTypeTask)
	contexts := tc.contextLoader.LoadContexts(
		cfg.Contexts,
		cfg.ContextOrder,
		engine.CommandTypeTask,
		shell,
		timeout,
	)

	// Load task with instructions
	loadedTask, err := tc.taskLoader.LoadTask(task, instructions, shell, timeout)
	if err != nil {
		return fmt.Errorf("failed to load task: %w", err)
	}

	// Execute agent (replaces current process, never returns on success)
	execParams := engine.ExecuteParams{
		Agent:        agent,
		Model:        modelID,
		UserPrompt:   loadedTask.Prompt,
		RoleContent:  loadedRole.Content,
		RoleFilePath: loadedRole.FilePath,
		Contexts:     contexts,
		Shell:        shell,
	}

	if err := tc.executor.Execute(execParams); err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	return nil
}

// listTasks displays all configured tasks
func (tc *TaskCommand) listTasks(cfg, globalCfg, localCfg domain.Config) error {
	allTasks := tc.taskResolver.ListAllTasks(localCfg.Tasks, globalCfg.Tasks)

	if len(allTasks) == 0 {
		fmt.Println("No tasks configured.")
		fmt.Println()
		fmt.Println("Add tasks to your configuration:")
		fmt.Println("  start config edit")
		return nil
	}

	// Sort task names
	var names []string
	for name := range allTasks {
		names = append(names, name)
	}
	sort.Strings(names)

	fmt.Println("Available tasks:")
	for _, name := range names {
		task := allTasks[name]
		aliasStr := ""
		if task.Alias != "" {
			aliasStr = fmt.Sprintf(" (%s)", task.Alias)
		}
		descStr := ""
		if task.Description != "" {
			descStr = fmt.Sprintf(" - %s", task.Description)
		}
		fmt.Printf("  %s%s%s\n", name, aliasStr, descStr)
	}
	fmt.Println()
	fmt.Println("Use 'start task <name>' to run a task.")

	return nil
}

// taskNotFoundError returns a helpful error when task is not found
func (tc *TaskCommand) taskNotFoundError(taskName string, cfg domain.Config) error {
	msg := fmt.Sprintf("Task %q not found.", taskName)

	// Show available tasks if any exist
	if len(cfg.Tasks) > 0 {
		msg += "\n\nAvailable tasks:"
		var names []string
		for name := range cfg.Tasks {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			task := cfg.Tasks[name]
			aliasStr := ""
			if task.Alias != "" {
				aliasStr = fmt.Sprintf(" (%s)", task.Alias)
			}
			msg += fmt.Sprintf("\n  %s%s", name, aliasStr)
		}
		msg += "\n\nUse 'start task' to see all tasks."
	}

	return fmt.Errorf("%s", msg)
}
