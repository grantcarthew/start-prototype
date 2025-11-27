package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/grantcarthew/start/internal/config"
	"github.com/grantcarthew/start/internal/domain"
	"github.com/spf13/cobra"
)

// NewConfigTaskRemoveCommand creates the config task remove command
func NewConfigTaskRemoveCommand(configLoader *config.Loader) *cobra.Command {
	var localOnly bool

	cmd := &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove task from configuration",
		Long:  "Remove a task from the configuration file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskName := args[0]
			prompter := NewPromptHelper()
			tomlHelper := config.NewTOMLHelper(configLoader.GetFS())
			backupHelper := config.NewBackupHelper(configLoader.GetFS())

			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			// Determine which config file contains the task
			var targetDir string
			var scope string
			var tasks map[string]domain.Task

			globalDir, err := tomlHelper.GetGlobalDir()
			if err != nil {
				return err
			}
			localDir := tomlHelper.GetLocalDir(workDir)

			if localOnly {
				// Remove from local only
				tasks, err = tomlHelper.ReadTasksFile(localDir)
				if err != nil {
					return fmt.Errorf("failed to read local tasks: %w", err)
				}
				targetDir = localDir
				scope = "local"
			} else {
				// Check both global and local
				globalTasks, err := tomlHelper.ReadTasksFile(globalDir)
				if err != nil {
					return fmt.Errorf("failed to read global tasks: %w", err)
				}

				localTasks, err := tomlHelper.ReadTasksFile(localDir)
				if err != nil {
					return fmt.Errorf("failed to read local tasks: %w", err)
				}

				// Check where task exists
				_, inGlobal := globalTasks[taskName]
				_, inLocal := localTasks[taskName]

				if !inGlobal && !inLocal {
					fmt.Printf("Error: Task '%s' not found in configuration.\n\n", taskName)
					fmt.Println("Use 'start config task list' to see available tasks.")
					return fmt.Errorf("task not found")
				}

				if inGlobal && inLocal {
					// Exists in both - ask which to remove
					fmt.Println("Task exists in both global and local configs.")
					fmt.Println()
					choice, err := prompter.AskChoice("Select scope to remove from:", []string{
						"global",
						"local",
						"both",
					})
					if err != nil {
						return err
					}

					if choice == "both" {
						// Remove from both
						if err := removeTaskFromScope(taskName, globalDir, "global", globalTasks, prompter, tomlHelper, backupHelper); err != nil {
							return err
						}
						return removeTaskFromScope(taskName, localDir, "local", localTasks, prompter, tomlHelper, backupHelper)
					} else if choice == "global" {
						targetDir = globalDir
						tasks = globalTasks
						scope = "global"
					} else {
						targetDir = localDir
						tasks = localTasks
						scope = "local"
					}
				} else if inLocal {
					targetDir = localDir
					tasks = localTasks
					scope = "local"
				} else {
					targetDir = globalDir
					tasks = globalTasks
					scope = "global"
				}
			}

			return removeTaskFromScope(taskName, targetDir, scope, tasks, prompter, tomlHelper, backupHelper)
		},
	}

	cmd.Flags().BoolVarP(&localOnly, "local", "l", false, "Remove from local config only")

	return cmd
}

// removeTaskFromScope removes a task from a specific scope
func removeTaskFromScope(
	taskName string,
	targetDir string,
	scope string,
	tasks map[string]domain.Task,
	prompter *PromptHelper,
	tomlHelper *config.TOMLHelper,
	backupHelper *config.BackupHelper,
) error {
	// Check if task exists
	task, exists := tasks[taskName]
	if !exists {
		return fmt.Errorf("task '%s' not found in %s config", taskName, scope)
	}

	// Show task details
	if task.Alias != "" {
		fmt.Printf("Task: %s (alias: %s)\n", taskName, task.Alias)
	} else {
		fmt.Printf("Task: %s\n", taskName)
	}
	if task.Description != "" {
		fmt.Printf("Description: %s\n", task.Description)
	}
	fmt.Println()

	// Confirm removal
	confirmed, err := prompter.AskYesNo(fmt.Sprintf("Remove task '%s' from %s config?", taskName, scope), false)
	if err != nil {
		return err
	}

	if !confirmed {
		fmt.Println()
		fmt.Printf("Task '%s' not removed.\n", taskName)
		return nil
	}

	// Backup existing config
	fmt.Println()
	tasksPath := filepath.Join(targetDir, "tasks.toml")
	fmt.Println("Backing up config to tasks.YYYY-MM-DD-HHMMSS.toml...")
	backupPath, err := backupHelper.CreateBackup(tasksPath)
	if err != nil {
		return fmt.Errorf("failed to backup config: %w", err)
	}
	prompter.PrintSuccess(fmt.Sprintf("Backup created: %s", filepath.Base(backupPath)))
	fmt.Println()

	// Remove task and save
	delete(tasks, taskName)
	if err := tomlHelper.WriteTasksFile(targetDir, tasks); err != nil {
		return fmt.Errorf("failed to write tasks file: %w", err)
	}

	fmt.Printf("Removing task '%s' from %s...\n", taskName, tasksPath)
	prompter.PrintSuccess(fmt.Sprintf("Task '%s' removed successfully", taskName))
	fmt.Println()
	fmt.Printf("Use 'start config task list' to see remaining tasks.\n")

	return nil
}
