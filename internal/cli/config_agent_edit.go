package cli

import (
	"fmt"
	"os"

	"github.com/grantcarthew/start/internal/config"
	"github.com/spf13/cobra"
)

// NewConfigAgentEditCommand creates the config agent edit command
func NewConfigAgentEditCommand(configLoader *config.Loader) *cobra.Command {
	var localOnly bool

	cmd := &cobra.Command{
		Use:   "edit [name]",
		Short: "Edit agent configuration",
		Long:  "Edit an existing agent configuration (currently requires manual editing)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("agent name required\n\nUsage: start config agent edit <name>")
			}

			agentName := args[0]
			tomlHelper := config.NewTOMLHelper(configLoader.GetFS())

			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			// Determine which config has the agent
			globalDir, err := tomlHelper.GetGlobalDir()
			if err != nil {
				return err
			}
			localDir := tomlHelper.GetLocalDir(workDir)

			globalAgents, _ := tomlHelper.ReadAgentsFile(globalDir)
			localAgents, _ := tomlHelper.ReadAgentsFile(localDir)

			hasGlobal := false
			hasLocal := false
			var configPath string
			var scope string

			if _, ok := globalAgents[agentName]; ok {
				hasGlobal = true
				configPath = tomlHelper.GetConfigPath(globalDir)
				configPath = globalDir + "/agents.toml"
				scope = "global"
			}
			if _, ok := localAgents[agentName]; ok {
				hasLocal = true
				if !hasGlobal || localOnly {
					configPath = localDir + "/agents.toml"
					scope = "local"
				}
			}

			if !hasGlobal && !hasLocal {
				return fmt.Errorf("agent '%s' not found in configuration.\n\nUse 'start config agent list' to see available agents.", agentName)
			}

			// For Phase 8b, provide guidance for manual editing
			fmt.Printf("Edit agent: %s\n", agentName)
			fmt.Println()
			fmt.Printf("Agent found in %s config.\n", scope)
			fmt.Printf("Config file: %s\n", configPath)
			fmt.Println()
			fmt.Println("To edit this agent:")
			fmt.Printf("  1. Open the config file in your editor: %s\n", configPath)
			fmt.Printf("  2. Find the [agents.%s] section\n", agentName)
			fmt.Println("  3. Make your changes and save")
			fmt.Printf("  4. Test your changes: start config agent test %s\n", agentName)
			fmt.Println()
			fmt.Println("Or recreate the agent:")
			fmt.Printf("  start config agent remove %s\n", agentName)
			fmt.Printf("  start config agent new\n")
			fmt.Println()
			fmt.Println("Note: Full interactive editing will be added in a future version.")

			return nil
		},
	}

	cmd.Flags().BoolVarP(&localOnly, "local", "l", false, "Edit local agent")

	return cmd
}
