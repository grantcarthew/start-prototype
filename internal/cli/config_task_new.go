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

// NewConfigTaskNewCommand creates the config task new command
func NewConfigTaskNewCommand(configLoader *config.Loader) *cobra.Command {
	var localOnly bool

	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create new task interactively",
		Long:  "Interactive wizard to create a new task configuration",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			prompter := NewPromptHelper()
			tomlHelper := config.NewTOMLHelper(configLoader.GetFS())
			backupHelper := config.NewBackupHelper(configLoader.GetFS())

			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			prompter.PrintHeader("Add new task")

			// Determine scope
			var targetDir string
			var scope string

			if !localOnly {
				choice, err := prompter.AskChoice("Select scope:", []string{
					"global (all projects)",
					"local (this project only)",
				})
				if err != nil {
					return err
				}

				if choice == "global (all projects)" {
					targetDir, err = tomlHelper.GetGlobalDir()
					if err != nil {
						return err
					}
					scope = "global"
				} else {
					targetDir = tomlHelper.GetLocalDir(workDir)
					scope = "local"

					// Check if local directory exists
					if !tomlHelper.GetFS().Exists(targetDir) {
						fmt.Printf("\n✗ Local config directory doesn't exist: %s\n", targetDir)
						fmt.Printf("  Create it first: mkdir -p %s\n\n", targetDir)
						fmt.Println("Or add to global config instead.")
						return fmt.Errorf("local config directory doesn't exist")
					}
				}
			} else {
				targetDir = tomlHelper.GetLocalDir(workDir)
				scope = "local"

				if !tomlHelper.GetFS().Exists(targetDir) {
					return fmt.Errorf("local config directory doesn't exist: %s\nCreate it first: mkdir -p %s", targetDir, targetDir)
				}
			}

			// Read existing tasks to check for duplicates
			existingTasks, err := tomlHelper.ReadTasksFile(targetDir)
			if err != nil {
				return fmt.Errorf("failed to read existing tasks: %w", err)
			}

			// Task name
			var taskName string
			for {
				taskName, err = prompter.AskValidatedName("\nTask name: ")
				if err != nil {
					return err
				}

				// Check for duplicate
				if _, exists := existingTasks[taskName]; exists {
					prompter.PrintError(fmt.Sprintf("Task '%s' already exists in %s config.", taskName, scope))
					fmt.Println()
					fmt.Println("Use 'start config task edit", taskName, "' to modify existing task.")
					fmt.Println()
					continue
				}

				break
			}

			// Create task struct
			newTask := domain.Task{}

			// Alias (optional)
			fmt.Println()
			alias, err := prompter.Ask("Alias (optional): ")
			if err != nil {
				return err
			}
			alias = strings.TrimSpace(alias)
			if alias != "" {
				// Validate alias
				if err := prompter.ValidateName(alias); err != nil {
					fmt.Printf("⚠ Warning: Invalid alias format: %v\n", err)
					fmt.Println("  Alias not set.")
				} else {
					// Check for duplicate alias
					duplicate := false
					for _, t := range existingTasks {
						if t.Alias == alias {
							fmt.Printf("⚠ Warning: Alias '%s' already in use\n", alias)
							fmt.Println("  Alias not set.")
							duplicate = true
							break
						}
					}
					if !duplicate {
						newTask.Alias = alias
					}
				}
			}

			// Description (optional)
			description, err := prompter.Ask("Description (optional): ")
			if err != nil {
				return err
			}
			newTask.Description = strings.TrimSpace(description)

			// Role selection (optional)
			fmt.Println()
			selectRole, err := prompter.AskYesNo("Select role?", false)
			if err != nil {
				return err
			}

			if selectRole {
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
					fmt.Println("  Create roles with: start config role new")
					prompter.PrintSuccess("Will use default role")
				} else {
					// Show available roles
					roleNames := make([]string, 0, len(roles))
					for name := range roles {
						roleNames = append(roleNames, name)
					}
					sort.Strings(roleNames)

					fmt.Println("\nAvailable roles:")
					options := make([]string, 0, len(roleNames)+1)
					for i, name := range roleNames {
						role := roles[name]
						if role.Description != "" {
							fmt.Printf("  %d) %s - %s\n", i+1, name, role.Description)
						} else {
							fmt.Printf("  %d) %s\n", i+1, name)
						}
						options = append(options, name)
					}
					options = append(options, "(skip - use default)")
					fmt.Printf("  %d) (skip - use default)\n", len(options))
					fmt.Println()

					choice, err := prompter.AskChoice("Select role:", options)
					if err != nil {
						return err
					}

					if choice != "(skip - use default)" {
						newTask.Role = choice
						prompter.PrintSuccess(fmt.Sprintf("Selected role: %s", choice))
					} else {
						prompter.PrintSuccess("Will use default role")
					}
				}
			} else {
				prompter.PrintSuccess("Will use default role")
			}

			// Agent selection (optional)
			fmt.Println()
			selectAgent, err := prompter.AskYesNo("Select agent?", false)
			if err != nil {
				return err
			}

			if selectAgent {
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
					fmt.Println("  Configure agents with: start init")
					prompter.PrintSuccess("Will use default agent")
				} else {
					// Show available agents
					agentNames := make([]string, 0, len(agents))
					for name := range agents {
						agentNames = append(agentNames, name)
					}
					sort.Strings(agentNames)

					fmt.Println("\nAvailable agents:")
					options := make([]string, 0, len(agentNames)+1)
					for i, name := range agentNames {
						agent := agents[name]
						if agent.Description != "" {
							fmt.Printf("  %d) %s - %s\n", i+1, name, agent.Description)
						} else {
							fmt.Printf("  %d) %s\n", i+1, name)
						}
						options = append(options, name)
					}
					options = append(options, "(skip - use default)")
					fmt.Printf("  %d) (skip - use default)\n", len(options))
					fmt.Println()

					choice, err := prompter.AskChoice("Select agent:", options)
					if err != nil {
						return err
					}

					if choice != "(skip - use default)" {
						newTask.Agent = choice
						prompter.PrintSuccess(fmt.Sprintf("Selected agent: %s", choice))
					} else {
						prompter.PrintSuccess("Will use default agent")
					}
				}
			} else {
				prompter.PrintSuccess("Will use default agent")
			}

			// Task prompt
			fmt.Println()
			fmt.Println("Task prompt:")
			sourceChoice, err := prompter.AskChoice("Content source:", []string{
				"File path",
				"Command",
				"Inline prompt",
				"Combination",
			})
			if err != nil {
				return err
			}

			switch sourceChoice {
			case "File path":
				// File-based task
				for {
					filePath, err := prompter.Ask("\nFile path: ")
					if err != nil {
						return err
					}
					filePath = strings.TrimSpace(filePath)

					if filePath == "" {
						prompter.PrintError("File path cannot be empty.")
						continue
					}

					newTask.File = filePath

					// Check if file exists (warning only)
					resolved := resolvePath(filePath, workDir)
					if _, err := os.Stat(resolved); err != nil {
						fmt.Printf("⚠ Warning: File does not exist: %s\n", filePath)
						cont, err := prompter.AskYesNo("Continue anyway?", false)
						if err != nil {
							return err
						}
						if !cont {
							newTask.File = ""
							continue
						}
					} else {
						prompter.PrintSuccess("File exists")
					}

					break
				}

				// Prompt template
				fmt.Println()
				promptTemplate, err := prompter.Ask("Prompt template: ")
				if err != nil {
					return err
				}
				newTask.Prompt = strings.TrimSpace(promptTemplate)

				if newTask.Prompt != "" && strings.Contains(newTask.Prompt, "{instructions}") {
					prompter.PrintSuccess("Valid template (uses {instructions} placeholder)")
				}

			case "Command":
				// Command-based task
				for {
					cmdStr, err := prompter.Ask("\nCommand: ")
					if err != nil {
						return err
					}
					cmdStr = strings.TrimSpace(cmdStr)

					if cmdStr == "" {
						prompter.PrintError("Command cannot be empty.")
						continue
					}

					newTask.Command = cmdStr
					prompter.PrintSuccess("Valid command")
					break
				}

				// Prompt template
				fmt.Println()
				promptTemplate, err := prompter.Ask("Prompt template: ")
				if err != nil {
					return err
				}
				newTask.Prompt = strings.TrimSpace(promptTemplate)

				if newTask.Prompt != "" && strings.Contains(newTask.Prompt, "{instructions}") {
					prompter.PrintSuccess("Valid template (uses {instructions} placeholder)")
				}

			case "Inline prompt":
				// Inline prompt only
				for {
					promptText, err := prompter.Ask("\nPrompt text: ")
					if err != nil {
						return err
					}
					promptText = strings.TrimSpace(promptText)

					if promptText == "" {
						prompter.PrintError("Prompt text cannot be empty.")
						continue
					}

					newTask.Prompt = promptText
					if strings.Contains(promptText, "{instructions}") {
						prompter.PrintSuccess("Valid prompt (uses {instructions} placeholder)")
					} else {
						prompter.PrintSuccess("Valid prompt")
					}
					break
				}

			case "Combination":
				// File + Command combination
				// File path
				fmt.Println()
				filePath, err := prompter.Ask("File path (optional, press Enter to skip): ")
				if err != nil {
					return err
				}
				filePath = strings.TrimSpace(filePath)
				if filePath != "" {
					newTask.File = filePath

					// Check if file exists (warning only)
					resolved := resolvePath(filePath, workDir)
					if _, err := os.Stat(resolved); err != nil {
						fmt.Printf("⚠ Warning: File does not exist: %s\n", filePath)
						cont, err := prompter.AskYesNo("Continue anyway?", true)
						if err != nil {
							return err
						}
						if !cont {
							newTask.File = ""
						}
					} else {
						prompter.PrintSuccess("File exists")
					}
				}

				// Command
				fmt.Println()
				cmdStr, err := prompter.Ask("Command (optional, press Enter to skip): ")
				if err != nil {
					return err
				}
				cmdStr = strings.TrimSpace(cmdStr)
				if cmdStr != "" {
					newTask.Command = cmdStr
					prompter.PrintSuccess("Valid command")
				}

				// Prompt template
				fmt.Println()
				promptTemplate, err := prompter.Ask("Prompt template: ")
				if err != nil {
					return err
				}
				newTask.Prompt = strings.TrimSpace(promptTemplate)

				// Validate at least one source
				if newTask.File == "" && newTask.Command == "" && newTask.Prompt == "" {
					return fmt.Errorf("at least one content source is required (file, command, or prompt)")
				}

				if newTask.Prompt != "" && strings.Contains(newTask.Prompt, "{instructions}") {
					prompter.PrintSuccess("Valid template (uses {instructions} placeholder)")
				}
			}

			// Advanced options?
			fmt.Println()
			advanced, err := prompter.AskYesNo("Advanced options?", false)
			if err != nil {
				return err
			}

			if advanced {
				// Shell override
				fmt.Println()
				shell, err := prompter.Ask("Shell override (or enter for default): ")
				if err != nil {
					return err
				}
				newTask.Shell = strings.TrimSpace(shell)

				// Command timeout
				timeout, err := prompter.Ask("Command timeout in seconds (or enter for default): ")
				if err != nil {
					return err
				}
				if timeout != "" {
					var timeoutInt int
					_, err := fmt.Sscanf(timeout, "%d", &timeoutInt)
					if err != nil {
						return fmt.Errorf("invalid timeout value: %w", err)
					}
					newTask.CommandTimeout = timeoutInt
				}
			}

			// Backup existing config
			fmt.Println()
			tasksPath := filepath.Join(targetDir, "tasks.toml")
			if tomlHelper.GetFS().Exists(tasksPath) {
				fmt.Println("Backing up config to tasks.YYYY-MM-DD-HHMMSS.toml...")
				backupPath, err := backupHelper.CreateBackup(tasksPath)
				if err != nil {
					return fmt.Errorf("failed to backup config: %w", err)
				}
				prompter.PrintSuccess(fmt.Sprintf("Backup created: %s", filepath.Base(backupPath)))
				fmt.Println()
			}

			// Add to tasks map and save
			existingTasks[taskName] = newTask
			if err := tomlHelper.WriteTasksFile(targetDir, existingTasks); err != nil {
				return fmt.Errorf("failed to write tasks file: %w", err)
			}

			fmt.Printf("Saving task '%s' to %s...\n", taskName, tasksPath)
			prompter.PrintSuccess("Task added successfully")
			fmt.Println()
			fmt.Printf("Use 'start config task list' to see all tasks.\n")
			if newTask.Alias != "" {
				fmt.Printf("Use 'start task %s \"instructions\"' to run.\n", newTask.Alias)
			} else {
				fmt.Printf("Use 'start task %s \"instructions\"' to run.\n", taskName)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&localOnly, "local", "l", false, "Create in local config only")

	return cmd
}
