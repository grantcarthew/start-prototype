package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/grantcarthew/start/internal/config"
	"github.com/spf13/cobra"
)

// NewConfigAgentRemoveCommand creates the config agent remove command
func NewConfigAgentRemoveCommand(configLoader *config.Loader) *cobra.Command {
	var localOnly bool

	cmd := &cobra.Command{
		Use:   "remove [name]",
		Short: "Remove agent from configuration",
		Long:  "Delete an agent configuration with backup and confirmation",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			prompter := NewPromptHelper()
			tomlHelper := config.NewTOMLHelper(configLoader.GetFS())
			backupHelper := config.NewBackupHelper(configLoader.GetFS())

			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			// Determine scope
			var targetDir string
			var scope string

			if localOnly {
				targetDir = tomlHelper.GetLocalDir(workDir)
				scope = "local"
			} else {
				// Interactive selection or auto-detect
				globalDir, err := tomlHelper.GetGlobalDir()
				if err != nil {
					return err
				}
				localDir := tomlHelper.GetLocalDir(workDir)

				// If no name provided, return error (interactive selection not yet implemented)
				if len(args) == 0 {
					return fmt.Errorf("agent name required\n\nUsage: start config agent remove <name>")
				}

				// Check which config has the agent
				globalAgents, _ := tomlHelper.ReadAgentsFile(globalDir)
				localAgents, _ := tomlHelper.ReadAgentsFile(localDir)

				agentName := args[0]
				hasGlobal := false
				hasLocal := false

				if _, ok := globalAgents[agentName]; ok {
					hasGlobal = true
				}
				if _, ok := localAgents[agentName]; ok {
					hasLocal = true
				}

				if !hasGlobal && !hasLocal {
					return fmt.Errorf("agent '%s' not found in configuration.\n\nUse 'start config agent list' to see available agents.", agentName)
				}

				// If exists in both, prompt for scope
				if hasGlobal && hasLocal {
					choice, err := prompter.AskChoice("Agent exists in multiple scopes. Select scope to remove from:", []string{
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
						if err := removeAgent(agentName, globalDir, "global", tomlHelper, backupHelper, prompter); err != nil {
							return err
						}
						return removeAgent(agentName, localDir, "local", tomlHelper, backupHelper, prompter)
					}
				} else if hasGlobal {
					targetDir = globalDir
					scope = "global"
				} else {
					targetDir = localDir
					scope = "local"
				}
			}

			// Remove agent
			if len(args) == 0 {
				return fmt.Errorf("agent name required\n\nUsage: start config agent remove <name>")
			}

			return removeAgent(args[0], targetDir, scope, tomlHelper, backupHelper, prompter)
		},
	}

	cmd.Flags().BoolVarP(&localOnly, "local", "l", false, "Remove from local config")

	return cmd
}

// removeAgent removes an agent from the specified config directory
func removeAgent(name, dir, scope string, tomlHelper *config.TOMLHelper, backupHelper *config.BackupHelper, prompter *PromptHelper) error {
	// Read current agents
	agents, err := tomlHelper.ReadAgentsFile(dir)
	if err != nil {
		return fmt.Errorf("failed to read agents: %w", err)
	}

	// Check if agent exists
	if _, ok := agents[name]; !ok {
		return fmt.Errorf("agent '%s' not found in %s config", name, scope)
	}

	// Confirm removal
	confirmed, err := prompter.AskYesNo(fmt.Sprintf("Remove agent '%s' from %s config?", name, scope), false)
	if err != nil {
		return err
	}

	if !confirmed {
		fmt.Printf("\nAgent '%s' not removed.\n", name)
		return nil
	}

	// Create backup
	configPath := filepath.Join(dir, "agents.toml")
	backupPath, err := backupHelper.CreateBackup(configPath)
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	fmt.Printf("\n✓ Backup created: %s\n", filepath.Base(backupPath))

	// Remove agent
	delete(agents, name)

	// Write updated config
	if err := tomlHelper.WriteAgentsFile(dir, agents); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Printf("✓ Agent '%s' removed from %s config\n", name, scope)
	fmt.Printf("\nUse 'start config agent list' to see remaining agents.\n")

	return nil
}
