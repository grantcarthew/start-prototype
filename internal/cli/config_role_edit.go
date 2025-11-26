package cli

import (
	"fmt"
	"os"

	"github.com/grantcarthew/start/internal/config"
	"github.com/spf13/cobra"
)

// NewConfigRoleEditCommand creates the config role edit command
func NewConfigRoleEditCommand(configLoader *config.Loader) *cobra.Command {
	var localOnly bool

	cmd := &cobra.Command{
		Use:   "edit <name>",
		Short: "Edit role configuration",
		Long:  "Edit an existing role configuration (currently requires manual editing)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			roleName := args[0]
			tomlHelper := config.NewTOMLHelper(configLoader.GetFS())

			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			// Determine which config has the role
			globalDir, err := tomlHelper.GetGlobalDir()
			if err != nil {
				return err
			}
			localDir := tomlHelper.GetLocalDir(workDir)

			globalRoles, _ := tomlHelper.ReadRolesFile(globalDir)
			localRoles, _ := tomlHelper.ReadRolesFile(localDir)

			hasGlobal := false
			hasLocal := false
			var configPath string
			var scope string

			if _, ok := globalRoles[roleName]; ok {
				hasGlobal = true
				configPath = globalDir + "/roles.toml"
				scope = "global"
			}
			if _, ok := localRoles[roleName]; ok {
				hasLocal = true
				if !hasGlobal || localOnly {
					configPath = localDir + "/roles.toml"
					scope = "local"
				}
			}

			if !hasGlobal && !hasLocal {
				return fmt.Errorf("role '%s' not found in configuration.\n\nUse 'start config role list' to see available roles.", roleName)
			}

			// Provide guidance for manual editing
			fmt.Printf("Edit role: %s\n", roleName)
			fmt.Println()
			fmt.Printf("Role found in %s config.\n", scope)
			fmt.Printf("Config file: %s\n", configPath)
			fmt.Println()
			fmt.Println("To edit this role:")
			fmt.Printf("  1. Open the config file in your editor: %s\n", configPath)
			fmt.Printf("  2. Find the [roles.%s] section\n", roleName)
			fmt.Println("  3. Make your changes and save")
			fmt.Printf("  4. Test your changes: start config role test %s\n", roleName)
			fmt.Println()
			fmt.Println("Or recreate the role:")
			fmt.Printf("  start config role remove %s\n", roleName)
			fmt.Printf("  start config role new\n")
			fmt.Println()
			fmt.Println("Note: Full interactive editing will be added in a future version.")

			return nil
		},
	}

	cmd.Flags().BoolVarP(&localOnly, "local", "l", false, "Edit local role")

	return cmd
}
