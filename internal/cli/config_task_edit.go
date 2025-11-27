package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/grantcarthew/start/internal/config"
	"github.com/grantcarthew/start/internal/domain"
	"github.com/spf13/cobra"
)

// NewConfigTaskEditCommand creates the config task edit command
func NewConfigTaskEditCommand(configLoader *config.Loader) *cobra.Command {
	var localOnly bool

	cmd := &cobra.Command{
		Use:   "edit <name>",
		Short: "Edit task configuration interactively",
		Long:  "Interactive wizard to modify an existing task configuration",
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
				// Edit local only
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
					// Exists in both - ask which to edit
					fmt.Println("Task exists in both global and local configs.")
					fmt.Println()
					choice, err := prompter.AskChoice("Select scope to edit:", []string{
						"global",
						"local",
					})
					if err != nil {
						return err
					}

					if choice == "global" {
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

			// Get existing task
			existingTask, exists := tasks[taskName]
			if !exists {
				return fmt.Errorf("task '%s' not found in %s config", taskName, scope)
			}

			prompter.PrintHeader(fmt.Sprintf("Edit task: %s (%s)", taskName, scope))

			// Show current configuration
			fmt.Println("Current configuration:")
			if existingTask.Alias != "" {
				fmt.Printf("  Alias: %s\n", existingTask.Alias)
			}
			if existingTask.Description != "" {
				fmt.Printf("  Description: %s\n", existingTask.Description)
			}
			if existingTask.Role != "" {
				fmt.Printf("  Role: %s\n", existingTask.Role)
			} else {
				fmt.Println("  Role: (default)")
			}
			if existingTask.Agent != "" {
				fmt.Printf("  Agent: %s\n", existingTask.Agent)
			} else {
				fmt.Println("  Agent: (default)")
			}
			if existingTask.File != "" {
				fmt.Printf("  File: %s\n", existingTask.File)
			}
			if existingTask.Command != "" {
				fmt.Printf("  Command: %s\n", existingTask.Command)
			}
			if existingTask.Prompt != "" {
				promptPreview := existingTask.Prompt
				if len(promptPreview) > 60 {
					promptPreview = promptPreview[:57] + "..."
				}
				fmt.Printf("  Prompt: %s\n", promptPreview)
			}
			fmt.Println()
			fmt.Println("Press enter to keep current value, or type new value:")
			fmt.Println()

			// Create updated task (start with existing)
			updatedTask := existingTask

			// Alias
			aliasPrompt := "Alias"
			if existingTask.Alias != "" {
				aliasPrompt = fmt.Sprintf("Alias [%s]", existingTask.Alias)
			}
			alias, err := prompter.Ask(aliasPrompt + ": ")
			if err != nil {
				return err
			}
			alias = strings.TrimSpace(alias)
			if alias != "" {
				// Validate and check for duplicates
				if err := prompter.ValidateName(alias); err != nil {
					fmt.Printf("⚠ Warning: Invalid alias format: %v\n", err)
				} else {
					// Check for duplicate alias (excluding current task)
					duplicate := false
					for name, t := range tasks {
						if name != taskName && t.Alias == alias {
							fmt.Printf("⚠ Warning: Alias '%s' already in use\n", alias)
							duplicate = true
							break
						}
					}
					if !duplicate {
						updatedTask.Alias = alias
					}
				}
			}

			// Description
			descPrompt := "Description"
			if existingTask.Description != "" {
				descPrompt = fmt.Sprintf("Description [%s]", existingTask.Description)
			}
			description, err := prompter.Ask(descPrompt + ": ")
			if err != nil {
				return err
			}
			description = strings.TrimSpace(description)
			if description != "" {
				updatedTask.Description = description
			}

			// Role selection
			fmt.Println()
			changeRole, err := prompter.AskYesNo("Change role selection?", false)
			if err != nil {
				return err
			}

			if changeRole {
				// Load roles from merged config
				globalCfg, err := configLoader.LoadGlobal()
				if err != nil {
					return fmt.Errorf("failed to load global config: %w", err)
				}

				var roles map[string]domain.Role
				localCfg, err := configLoader.LoadLocal(workDir)
				if err == nil {
					mergedCfg := config.Merge(globalCfg, localCfg)
					roles = mergedCfg.Roles
				} else {
					roles = globalCfg.Roles
				}

				if len(roles) == 0 {
					fmt.Println("⚠ No roles configured")
					updatedTask.Role = ""
				} else {
					roleNames := make([]string, 0, len(roles))
					for name := range roles {
						roleNames = append(roleNames, name)
					}
					sort.Strings(roleNames)

					options := append(roleNames, "(clear - use default)")
					choice, err := prompter.AskChoice("Select role:", options)
					if err != nil {
						return err
					}

					if choice == "(clear - use default)" {
						updatedTask.Role = ""
					} else {
						updatedTask.Role = choice
					}
				}
			}

			// Agent selection
			fmt.Println()
			changeAgent, err := prompter.AskYesNo("Change agent selection?", false)
			if err != nil {
				return err
			}

			if changeAgent {
				// Load agents from merged config
				globalCfg, err := configLoader.LoadGlobal()
				if err != nil {
					return fmt.Errorf("failed to load global config: %w", err)
				}

				var agents map[string]domain.Agent
				localCfg, err := configLoader.LoadLocal(workDir)
				if err == nil {
					mergedCfg := config.Merge(globalCfg, localCfg)
					agents = mergedCfg.Agents
				} else {
					agents = globalCfg.Agents
				}

				if len(agents) == 0 {
					fmt.Println("⚠ No agents configured")
					updatedTask.Agent = ""
				} else {
					agentNames := make([]string, 0, len(agents))
					for name := range agents {
						agentNames = append(agentNames, name)
					}
					sort.Strings(agentNames)

					options := append(agentNames, "(clear - use default)")
					choice, err := prompter.AskChoice("Select agent:", options)
					if err != nil {
						return err
					}

					if choice == "(clear - use default)" {
						updatedTask.Agent = ""
					} else {
						updatedTask.Agent = choice
					}
				}
			}

			// File path
			fmt.Println()
			filePrompt := "File path"
			if existingTask.File != "" {
				filePrompt = fmt.Sprintf("File path [%s]", existingTask.File)
			}
			filePath, err := prompter.Ask(filePrompt + ": ")
			if err != nil {
				return err
			}
			filePath = strings.TrimSpace(filePath)
			if filePath != "" {
				updatedTask.File = filePath
			}

			// Command
			fmt.Println()
			cmdPrompt := "Command"
			if existingTask.Command != "" {
				cmdPrompt = fmt.Sprintf("Command [%s]", existingTask.Command)
			}
			cmdStr, err := prompter.Ask(cmdPrompt + ": ")
			if err != nil {
				return err
			}
			cmdStr = strings.TrimSpace(cmdStr)
			if cmdStr != "" {
				updatedTask.Command = cmdStr
			}

			// Prompt template
			fmt.Println()
			promptPrompt := "Prompt template"
			if existingTask.Prompt != "" {
				promptPreview := existingTask.Prompt
				if len(promptPreview) > 40 {
					promptPreview = promptPreview[:37] + "..."
				}
				promptPrompt = fmt.Sprintf("Prompt template [%s]", promptPreview)
			}
			promptTemplate, err := prompter.Ask(promptPrompt + ": ")
			if err != nil {
				return err
			}
			promptTemplate = strings.TrimSpace(promptTemplate)
			if promptTemplate != "" {
				updatedTask.Prompt = promptTemplate
			}

			// Advanced options
			fmt.Println()
			advanced, err := prompter.AskYesNo("Advanced options?", false)
			if err != nil {
				return err
			}

			if advanced {
				// Shell override
				fmt.Println()
				shellPrompt := "Shell override"
				if existingTask.Shell != "" {
					shellPrompt = fmt.Sprintf("Shell override [%s]", existingTask.Shell)
				}
				shell, err := prompter.Ask(shellPrompt + " (or enter for default): ")
				if err != nil {
					return err
				}
				shell = strings.TrimSpace(shell)
				if shell != "" {
					updatedTask.Shell = shell
				}

				// Command timeout
				timeoutPrompt := "Command timeout in seconds"
				if existingTask.CommandTimeout > 0 {
					timeoutPrompt = fmt.Sprintf("Command timeout in seconds [%d]", existingTask.CommandTimeout)
				}
				timeout, err := prompter.Ask(timeoutPrompt + " (or enter for default): ")
				if err != nil {
					return err
				}
				timeout = strings.TrimSpace(timeout)
				if timeout != "" {
					var timeoutInt int
					_, err := fmt.Sscanf(timeout, "%d", &timeoutInt)
					if err != nil {
						return fmt.Errorf("invalid timeout value: %w", err)
					}
					updatedTask.CommandTimeout = timeoutInt
				}
			}

			// Check if anything changed
			if updatedTask == existingTask {
				fmt.Println()
				fmt.Println("No changes detected.")
				fmt.Println()
				fmt.Printf("Task '%s' not modified.\n", taskName)
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

			// Update task and save
			tasks[taskName] = updatedTask
			if err := tomlHelper.WriteTasksFile(targetDir, tasks); err != nil {
				return fmt.Errorf("failed to write tasks file: %w", err)
			}

			fmt.Printf("Saving changes to %s...\n", tasksPath)
			prompter.PrintSuccess(fmt.Sprintf("Task '%s' updated successfully", taskName))
			fmt.Println()
			fmt.Printf("Use 'start config task list' to see changes.\n")
			fmt.Printf("Use 'start config task test %s' to validate.\n", taskName)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&localOnly, "local", "l", false, "Edit in local config only")

	return cmd
}
