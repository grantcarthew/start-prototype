package cli

import (
	"fmt"
	"os"

	"github.com/grantcarthew/start/internal/config"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

// NewConfigCommand creates the config command
func NewConfigCommand(configLoader *config.Loader, validator *config.Validator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Long:  "Commands for viewing and managing start configuration",
	}

	// Add subcommands
	cmd.AddCommand(NewConfigShowCommand(configLoader, validator))
	cmd.AddCommand(NewConfigAgentCommand(configLoader, validator))
	cmd.AddCommand(NewConfigRoleCommand(configLoader, validator))

	return cmd
}

// NewConfigShowCommand creates the config show command
func NewConfigShowCommand(configLoader *config.Loader, validator *config.Validator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show merged configuration",
		Long:  "Display the merged configuration from global and local config files",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get current working directory
			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			// Load global config
			globalCfg, err := configLoader.LoadGlobal()
			if err != nil {
				return fmt.Errorf("failed to load global config: %w", err)
			}

			// Load local config
			localCfg, err := configLoader.LoadLocal(workDir)
			if err != nil {
				// Local config is optional, so we just use global if it doesn't exist
				localCfg = globalCfg
				globalCfg = config.Merge(globalCfg, localCfg)
			} else {
				// Merge configs
				globalCfg = config.Merge(globalCfg, localCfg)
			}

			// Validate merged config
			if err := validator.Validate(globalCfg); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Configuration validation errors:\n%v\n\n", err)
			}

			// Marshal to TOML for display
			output, err := toml.Marshal(globalCfg)
			if err != nil {
				return fmt.Errorf("failed to marshal config: %w", err)
			}

			fmt.Println(string(output))
			return nil
		},
	}

	return cmd
}
