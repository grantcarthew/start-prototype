package cli

import (
	"github.com/grantcarthew/start/internal/config"
	"github.com/spf13/cobra"
)

// NewRootCommand creates the root command
func NewRootCommand(configLoader *config.Loader, validator *config.Validator, version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Short:   "AI agent CLI orchestrator",
		Long:    "start is a command-line orchestrator for AI agents that manages prompt composition, context injection, and workflow automation.",
		Version: version,
		// RunE will be implemented in later phases
	}

	// Add subcommands
	cmd.AddCommand(NewConfigCommand(configLoader, validator))

	return cmd
}
