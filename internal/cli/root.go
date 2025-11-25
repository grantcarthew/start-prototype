package cli

import (
	"fmt"
	"strings"

	"github.com/grantcarthew/start/internal/config"
	"github.com/grantcarthew/start/internal/engine"
	"github.com/spf13/cobra"
)

// RootCommand holds the dependencies for the root command
type RootCommand struct {
	configLoader   *config.Loader
	validator      *config.Validator
	executor       *engine.Executor
	roleSelector   *engine.RoleSelector
	roleLoader     *engine.RoleLoader
	contextLoader  *engine.ContextLoader
	version        string
}

// NewRootCommand creates the root command
func NewRootCommand(
	configLoader *config.Loader,
	validator *config.Validator,
	executor *engine.Executor,
	roleSelector *engine.RoleSelector,
	roleLoader *engine.RoleLoader,
	contextLoader *engine.ContextLoader,
	version string,
) *cobra.Command {
	rc := &RootCommand{
		configLoader:   configLoader,
		validator:      validator,
		executor:       executor,
		roleSelector:   roleSelector,
		roleLoader:     roleLoader,
		contextLoader:  contextLoader,
		version:        version,
	}

	cmd := &cobra.Command{
		Use:     "start [prompt]",
		Short:   "AI agent CLI orchestrator",
		Long:    "start is a command-line orchestrator for AI agents that manages prompt composition, context injection, and workflow automation.",
		Version: version,
		RunE:    rc.run,
		Args:    cobra.ArbitraryArgs,
	}

	// Add persistent flags
	cmd.PersistentFlags().StringP("agent", "a", "", "Agent to use")
	cmd.PersistentFlags().StringP("model", "m", "", "Model to use")
	cmd.PersistentFlags().StringP("role", "r", "", "Role to use")

	// Add subcommands
	cmd.AddCommand(NewConfigCommand(configLoader, validator))

	return cmd
}

// run executes the root command
func (rc *RootCommand) run(cmd *cobra.Command, args []string) error {
	// Load configuration
	globalCfg, err := rc.configLoader.LoadGlobal()
	if err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}

	// Get current working directory for local config
	workDir := "."
	localCfg, err := rc.configLoader.LoadLocal(workDir)
	if err != nil {
		// Local config is optional, use empty config
		localCfg = globalCfg
	}

	// Merge configs
	cfg := config.Merge(globalCfg, localCfg)

	// Validate merged config
	if err := rc.validator.Validate(cfg); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	// Get flags
	agentFlag, _ := cmd.Flags().GetString("agent")
	modelFlag, _ := cmd.Flags().GetString("model")
	roleFlag, _ := cmd.Flags().GetString("role")

	// Select agent
	var agentName string
	if agentFlag != "" {
		agentName = agentFlag
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

	// Select role
	selectionCtx := engine.SelectionContext{
		RoleFlag:    roleFlag,
		TaskRole:    "", // No task for root command
		DefaultRole: cfg.Settings.DefaultRole,
	}
	role, err := rc.roleSelector.Select(selectionCtx, cfg.Roles)
	if err != nil {
		return fmt.Errorf("role selection failed: %w", err)
	}

	// Get shell and timeout settings
	shell := cfg.Settings.Shell
	if shell == "" {
		shell = "bash"
	}
	timeout := cfg.Settings.CommandTimeout
	if timeout == 0 {
		timeout = 30 // Default 30 seconds
	}

	// Load role
	loadedRole, err := rc.roleLoader.LoadRole(role, shell, timeout)
	if err != nil {
		return fmt.Errorf("failed to load role: %w", err)
	}
	// Cleanup temp role file if needed (deferred)
	defer rc.roleLoader.CleanupRole(loadedRole)

	// Load contexts (interactive mode = all contexts)
	contexts := rc.contextLoader.LoadContexts(
		cfg.Contexts,
		cfg.ContextOrder,
		engine.CommandTypeInteractive,
		shell,
		timeout,
	)

	// Assemble prompt from arguments
	userPrompt := strings.Join(args, " ")
	if userPrompt == "" {
		return fmt.Errorf("no prompt provided")
	}

	// Execute agent (replaces current process, never returns on success)
	execParams := engine.ExecuteParams{
		Agent:        agent,
		Model:        modelID,
		UserPrompt:   userPrompt,
		RoleContent:  loadedRole.Content,
		RoleFilePath: loadedRole.FilePath,
		Contexts:     contexts,
		Shell:        shell,
	}

	if err := rc.executor.Execute(execParams); err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	return nil
}
