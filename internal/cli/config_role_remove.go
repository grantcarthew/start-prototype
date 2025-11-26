package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/grantcarthew/start/internal/config"
	"github.com/spf13/cobra"
)

// NewConfigRoleRemoveCommand creates the config role remove command
func NewConfigRoleRemoveCommand(configLoader *config.Loader) *cobra.Command {
	var localOnly bool

	cmd := &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove role from configuration",
		Long:  "Delete a role configuration with backup and confirmation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			prompter := NewPromptHelper()
			tomlHelper := config.NewTOMLHelper(configLoader.GetFS())
			backupHelper := config.NewBackupHelper(configLoader.GetFS())

			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			roleName := args[0]

			// Determine scope
			var targetDir string
			var scope string

			if localOnly {
				targetDir = tomlHelper.GetLocalDir(workDir)
				scope = "local"
			} else {
				// Check which config has the role
				globalDir, err := tomlHelper.GetGlobalDir()
				if err != nil {
					return err
				}
				localDir := tomlHelper.GetLocalDir(workDir)

				globalRoles, _ := tomlHelper.ReadRolesFile(globalDir)
				localRoles, _ := tomlHelper.ReadRolesFile(localDir)

				hasGlobal := false
				hasLocal := false

				if _, ok := globalRoles[roleName]; ok {
					hasGlobal = true
				}
				if _, ok := localRoles[roleName]; ok {
					hasLocal = true
				}

				if !hasGlobal && !hasLocal {
					return fmt.Errorf("role '%s' not found in configuration.\n\nUse 'start config role list' to see available roles.", roleName)
				}

				// If exists in both, prompt for scope
				if hasGlobal && hasLocal {
					choice, err := prompter.AskChoice("Role exists in multiple scopes. Select scope to remove from:", []string{
						"global",
						"local",
						"both",
					})
					if err != nil {
						return err
					}

					switch choice {
					case "global":
						targetDir = globalDir
						scope = "global"
					case "local":
						targetDir = localDir
						scope = "local"
					case "both":
						// Remove from both
						if err := removeRole(roleName, globalDir, "global", tomlHelper, backupHelper, prompter); err != nil {
							return err
						}
						return removeRole(roleName, localDir, "local", tomlHelper, backupHelper, prompter)
					}
				} else if hasGlobal {
					targetDir = globalDir
					scope = "global"
				} else {
					targetDir = localDir
					scope = "local"
				}
			}

			return removeRole(roleName, targetDir, scope, tomlHelper, backupHelper, prompter)
		},
	}

	cmd.Flags().BoolVarP(&localOnly, "local", "l", false, "Remove from local config")

	return cmd
}

// removeRole removes a role from the specified config directory
func removeRole(name, dir, scope string, tomlHelper *config.TOMLHelper, backupHelper *config.BackupHelper, prompter *PromptHelper) error {
	// Read current roles
	roles, err := tomlHelper.ReadRolesFile(dir)
	if err != nil {
		return fmt.Errorf("failed to read roles: %w", err)
	}

	// Check if role exists
	if _, ok := roles[name]; !ok {
		return fmt.Errorf("role '%s' not found in %s config", name, scope)
	}

	// Confirm removal
	confirmed, err := prompter.AskYesNo(fmt.Sprintf("Remove role '%s' from %s config?", name, scope), false)
	if err != nil {
		return err
	}

	if !confirmed {
		fmt.Printf("\nRole '%s' not removed.\n", name)
		return nil
	}

	// Create backup
	configPath := filepath.Join(dir, "roles.toml")
	backupPath, err := backupHelper.CreateBackup(configPath)
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	fmt.Printf("\n✓ Backup created: %s\n", filepath.Base(backupPath))

	// Remove role
	delete(roles, name)

	// Write updated config
	if err := tomlHelper.WriteRolesFile(dir, roles); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Printf("✓ Role '%s' removed from %s config\n", name, scope)
	fmt.Printf("\nUse 'start config role list' to see remaining roles.\n")

	return nil
}
