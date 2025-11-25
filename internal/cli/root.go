package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/grantcarthew/start/internal/config"
	"github.com/grantcarthew/start/internal/engine"
	"github.com/spf13/cobra"
)

// RootCommand holds the dependencies for the root command
type RootCommand struct {
	configLoader *config.Loader
	validator    *config.Validator
	executor     *engine.Executor
	version      string
}

// NewRootCommand creates the root command
func NewRootCommand(configLoader *config.Loader, validator *config.Validator, executor *engine.Executor, version string) *cobra.Command {
	rc := &RootCommand{
		configLoader: configLoader,
		validator:    validator,
		executor:     executor,
		version:      version,
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

	// Get agent flag
	agentFlag, _ := cmd.Flags().GetString("agent")

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

	// Get model flag
	modelFlag, _ := cmd.Flags().GetString("model")

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

	// Assemble prompt from arguments
	prompt := strings.Join(args, " ")
	if prompt == "" {
		return fmt.Errorf("no prompt provided")
	}

	// Execute agent
	ctx := context.Background()
	if err := rc.executor.Execute(ctx, agent, modelID, prompt); err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	return nil
}
