package cli

import (
	"fmt"
	"os"

	"github.com/grantcarthew/start/internal/config"
	"github.com/spf13/cobra"
)

// NewConfigAgentDefaultCommand creates the config agent default command
func NewConfigAgentDefaultCommand(configLoader *config.Loader) *cobra.Command {
	var localOnly bool

	cmd := &cobra.Command{
		Use:   "default [name]",
		Short: "Set or show default agent",
		Long:  "Set the default agent in settings, or show current default if no name provided",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tomlHelper := config.NewTOMLHelper(configLoader.GetFS())
			backupHelper := config.NewBackupHelper(configLoader.GetFS())

			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			// Determine target directory
			var targetDir string
			var scope string

			if localOnly {
				targetDir = tomlHelper.GetLocalDir(workDir)
				scope = "local"
			} else {
				targetDir, err = tomlHelper.GetGlobalDir()
				if err != nil {
					return err
				}
				scope = "global"
			}

			// If no agent name provided, show current default
			if len(args) == 0 {
				settings, err := tomlHelper.ReadSettingsFile(targetDir)
				if err != nil {
					return fmt.Errorf("failed to read settings: %w", err)
				}

				if settings.DefaultAgent == "" {
					fmt.Printf("No default agent set in %s config.\n", scope)
					fmt.Println()
					fmt.Println("Set default agent:")
					fmt.Println("  start config agent default <name>")
				} else {
					fmt.Printf("Default agent (%s): %s\n", scope, settings.DefaultAgent)
				}

				return nil
			}

			// Validate agent exists
			agentName := args[0]
			agents, err := tomlHelper.ReadAgentsFile(targetDir)
			if err != nil {
				return fmt.Errorf("failed to read agents: %w", err)
			}

			if _, ok := agents[agentName]; !ok {
				return fmt.Errorf("agent '%s' not found in %s config.\n\nUse 'start config agent list' to see available agents.", agentName, scope)
			}

			// Read current settings
			settings, err := tomlHelper.ReadSettingsFile(targetDir)
			if err != nil {
				return fmt.Errorf("failed to read settings: %w", err)
			}

			// Check if already default
			if settings.DefaultAgent == agentName {
				fmt.Printf("Agent '%s' is already the default agent in %s config.\n", agentName, scope)
				return nil
			}

			// Create backup if config file exists
			configPath := tomlHelper.GetConfigPath(targetDir)
			if tomlHelper.GetFS().Exists(configPath) {
				backupPath, err := backupHelper.CreateBackup(configPath)
				if err != nil {
					return fmt.Errorf("failed to create backup: %w", err)
				}
				fmt.Printf("✓ Backup created: %s\n", backupPath)
			}

			// Update settings
			settings.DefaultAgent = agentName

			// Write settings
			if err := tomlHelper.WriteSettingsFile(targetDir, settings); err != nil {
				return fmt.Errorf("failed to write settings: %w", err)
			}

			fmt.Printf("✓ Default agent set to '%s' in %s config\n", agentName, scope)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&localOnly, "local", "l", false, "Set default in local config")

	return cmd
}
