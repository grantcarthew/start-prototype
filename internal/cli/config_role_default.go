package cli

import (
	"fmt"
	"os"

	"github.com/grantcarthew/start/internal/config"
	"github.com/spf13/cobra"
)

// NewConfigRoleDefaultCommand creates the config role default command
func NewConfigRoleDefaultCommand(configLoader *config.Loader) *cobra.Command {
	var localOnly bool

	cmd := &cobra.Command{
		Use:   "default [name]",
		Short: "Set or show default role",
		Long:  "Set the default role in settings, or show current default if no name provided",
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

			// If no role name provided, show current default
			if len(args) == 0 {
				settings, err := tomlHelper.ReadSettingsFile(targetDir)
				if err != nil {
					return fmt.Errorf("failed to read settings: %w", err)
				}

				if settings.DefaultRole == "" {
					fmt.Printf("No default role set in %s config.\n", scope)
					fmt.Println()
					fmt.Println("Set default role:")
					fmt.Println("  start config role default <name>")
				} else {
					fmt.Printf("Default role (%s): %s\n", scope, settings.DefaultRole)
				}

				return nil
			}

			// Validate role exists
			roleName := args[0]
			roles, err := tomlHelper.ReadRolesFile(targetDir)
			if err != nil {
				return fmt.Errorf("failed to read roles: %w", err)
			}

			if _, ok := roles[roleName]; !ok {
				return fmt.Errorf("role '%s' not found in %s config.\n\nUse 'start config role list' to see available roles.", roleName, scope)
			}

			// Read current settings
			settings, err := tomlHelper.ReadSettingsFile(targetDir)
			if err != nil {
				return fmt.Errorf("failed to read settings: %w", err)
			}

			// Check if already default
			if settings.DefaultRole == roleName {
				fmt.Printf("Role '%s' is already the default role in %s config.\n", roleName, scope)
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
			settings.DefaultRole = roleName

			// Write settings
			if err := tomlHelper.WriteSettingsFile(targetDir, settings); err != nil {
				return fmt.Errorf("failed to write settings: %w", err)
			}

			fmt.Printf("✓ Default role set to '%s' in %s config\n", roleName, scope)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&localOnly, "local", "l", false, "Set default in local config")

	return cmd
}
