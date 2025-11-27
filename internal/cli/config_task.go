package cli

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/grantcarthew/start/internal/config"
	"github.com/grantcarthew/start/internal/domain"
	"github.com/spf13/cobra"
)

// NewConfigTaskCommand creates the config task command
func NewConfigTaskCommand(configLoader *config.Loader, validator *config.Validator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "Manage task configurations",
		Long:  "Commands for managing predefined workflow task configurations in global and local config files",
	}

	// Add subcommands
	cmd.AddCommand(NewConfigTaskListCommand(configLoader))
	cmd.AddCommand(NewConfigTaskShowCommand(configLoader))
	cmd.AddCommand(NewConfigTaskTestCommand(configLoader))
	cmd.AddCommand(NewConfigTaskNewCommand(configLoader))
	cmd.AddCommand(NewConfigTaskEditCommand(configLoader))
	cmd.AddCommand(NewConfigTaskRemoveCommand(configLoader))

	return cmd
}

// NewConfigTaskListCommand creates the config task list command
func NewConfigTaskListCommand(configLoader *config.Loader) *cobra.Command {
	var localOnly bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Display all configured tasks",
		Long:  "List all tasks defined in global and/or local configuration files",
		RunE: func(cmd *cobra.Command, args []string) error {
			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			var tasks map[string]domain.Task
			var scope string

			if localOnly {
				// Load local only
				localCfg, err := configLoader.LoadLocal(workDir)
				if err != nil {
					return fmt.Errorf("failed to load local config: %w", err)
				}
				tasks = localCfg.Tasks
				scope = "local"
			} else {
				// Load and merge global + local
				globalCfg, err := configLoader.LoadGlobal()
				if err != nil {
					return fmt.Errorf("failed to load global config: %w", err)
				}

				localCfg, err := configLoader.LoadLocal(workDir)
				if err == nil {
					mergedCfg := config.Merge(globalCfg, localCfg)
					tasks = mergedCfg.Tasks
					scope = "merged"
				} else {
					tasks = globalCfg.Tasks
					scope = "global"
				}
			}

			if len(tasks) == 0 {
				fmt.Println("No tasks configured.")
				fmt.Println()
				fmt.Println("Create task: start config task new")
				fmt.Println("Install from catalog: start assets add")
				return nil
			}

			// Sort tasks by name for consistent output
			names := make([]string, 0, len(tasks))
			for name := range tasks {
				names = append(names, name)
			}
			sort.Strings(names)

			fmt.Printf("Configured tasks (%s):\n", scope)
			fmt.Println("═══════════════════════════════════════════════════════════")
			fmt.Println()

			for _, name := range names {
				task := tasks[name]

				// Task name and alias
				if task.Alias != "" {
					fmt.Printf("%s (%s)\n", name, task.Alias)
				} else {
					fmt.Printf("%s\n", name)
				}

				// Description
				if task.Description != "" {
					fmt.Printf("  %s\n", task.Description)
				}

				// Role selection
				if task.Role != "" {
					fmt.Printf("  Role: %s\n", task.Role)
				} else {
					fmt.Println("  Role: (default)")
				}

				// Agent selection
				if task.Agent != "" {
					fmt.Printf("  Agent: %s\n", task.Agent)
				}

				// Task type
				sourceType := getTaskSourceType(task)
				fmt.Printf("  Task: %s\n", sourceType)

				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&localOnly, "local", "l", false, "List local tasks only")

	return cmd
}

// NewConfigTaskShowCommand creates the config task show command
func NewConfigTaskShowCommand(configLoader *config.Loader) *cobra.Command {
	var localOnly bool

	cmd := &cobra.Command{
		Use:   "show <name>",
		Short: "Display task configuration",
		Long:  "Show detailed configuration for a specific task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskName := args[0]
			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			var tasks map[string]domain.Task
			var scope string

			if localOnly {
				localCfg, err := configLoader.LoadLocal(workDir)
				if err != nil {
					return fmt.Errorf("failed to load local config: %w", err)
				}
				tasks = localCfg.Tasks
				scope = "local"
			} else {
				globalCfg, err := configLoader.LoadGlobal()
				if err != nil {
					return fmt.Errorf("failed to load global config: %w", err)
				}

				localCfg, err := configLoader.LoadLocal(workDir)
				if err == nil {
					mergedCfg := config.Merge(globalCfg, localCfg)
					tasks = mergedCfg.Tasks

					// Determine scope
					if _, hasLocal := localCfg.Tasks[taskName]; hasLocal {
						scope = "local"
					} else {
						scope = "global"
					}
				} else {
					tasks = globalCfg.Tasks
					scope = "global"
				}
			}

			task, exists := tasks[taskName]
			if !exists {
				fmt.Printf("No task '%s' found in configuration.\n\n", taskName)
				fmt.Println("Create task: start config task new")
				fmt.Println("Install from catalog: start assets add")
				return fmt.Errorf("task not found")
			}

			// Display task configuration
			fmt.Printf("Task configuration: %s (%s)\n", taskName, scope)
			fmt.Println("═══════════════════════════════════════════════════════════")
			fmt.Println()

			if task.Alias != "" {
				fmt.Printf("Alias: %s\n", task.Alias)
			}
			if task.Description != "" {
				fmt.Printf("Description: %s\n", task.Description)
			}
			fmt.Println()

			// Role and agent
			if task.Role != "" {
				fmt.Printf("Role: %s\n", task.Role)
			} else {
				fmt.Println("Role: (default)")
			}
			if task.Agent != "" {
				fmt.Printf("Agent: %s\n", task.Agent)
			} else {
				fmt.Println("Agent: (default)")
			}
			fmt.Println()

			sourceType := getTaskSourceType(task)
			fmt.Printf("Task prompt type: %s\n", sourceType)
			fmt.Println()

			if task.File != "" {
				fmt.Println("File:")
				fmt.Printf("  Path: %s\n", task.File)

				// Try to resolve and check file
				resolved := resolvePath(task.File, workDir)
				fmt.Printf("  Resolved: %s\n", resolved)

				if fileInfo, err := os.Stat(resolved); err == nil {
					fmt.Printf("  ✓ File exists (%.1f KB)\n", float64(fileInfo.Size())/1024)
				} else {
					fmt.Println("  ✗ File not found")
				}
				fmt.Println()
			}

			if task.Command != "" {
				fmt.Println("Command:")
				if task.Shell != "" {
					fmt.Printf("  Shell: %s\n", task.Shell)
				} else {
					fmt.Println("  Shell: (default)")
				}
				if task.CommandTimeout > 0 {
					fmt.Printf("  Timeout: %d seconds\n", task.CommandTimeout)
				} else {
					fmt.Println("  Timeout: (default)")
				}
				fmt.Printf("  Command: %s\n", task.Command)
				fmt.Println()
			}

			if task.Prompt != "" {
				fmt.Println("Prompt template:")
				// Show full prompt
				lines := strings.Split(task.Prompt, "\n")
				maxLines := 20
				if len(lines) > maxLines {
					for _, line := range lines[:maxLines] {
						fmt.Printf("  %s\n", line)
					}
					fmt.Printf("  ... (%d more lines)\n", len(lines)-maxLines)
				} else {
					for _, line := range lines {
						fmt.Printf("  %s\n", line)
					}
				}
				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&localOnly, "local", "l", false, "Show local task only")

	return cmd
}

// NewConfigTaskTestCommand creates the config task test command
func NewConfigTaskTestCommand(configLoader *config.Loader) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test <name>",
		Short: "Test task configuration and command execution",
		Long:  "Validate task configuration without executing the full task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskName := args[0]
			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			// Load merged config
			globalCfg, err := configLoader.LoadGlobal()
			if err != nil {
				return fmt.Errorf("failed to load global config: %w", err)
			}

			var tasks map[string]domain.Task
			var scope string
			localCfg, err := configLoader.LoadLocal(workDir)
			if err == nil {
				mergedCfg := config.Merge(globalCfg, localCfg)
				tasks = mergedCfg.Tasks

				if _, hasLocal := localCfg.Tasks[taskName]; hasLocal {
					scope = "local"
				} else {
					scope = "global"
				}
			} else {
				tasks = globalCfg.Tasks
				scope = "global"
			}

			task, exists := tasks[taskName]
			if !exists {
				fmt.Fprintf(os.Stderr, "Error: Task '%s' not found in configuration.\n\n", taskName)
				fmt.Fprintln(os.Stderr, "Use 'start config task list' to see available tasks.")
				return fmt.Errorf("task not found")
			}

			fmt.Printf("Testing task: %s\n", taskName)
			fmt.Println("─────────────────────────────────────────────────")
			fmt.Println()

			fmt.Println("Configuration:")
			fmt.Printf("  Scope: %s\n", scope)
			if task.Alias != "" {
				fmt.Printf("  Alias: %s\n", task.Alias)
			}
			if task.Description != "" {
				fmt.Printf("  Description: %s\n", task.Description)
			}
			if task.Role != "" {
				fmt.Printf("  Role: %s\n", task.Role)
			} else {
				fmt.Println("  Role: (default)")
			}
			if task.Agent != "" {
				fmt.Printf("  Agent: %s\n", task.Agent)
			} else {
				fmt.Println("  Agent: (default)")
			}
			fmt.Printf("  Type: %s\n", getTaskSourceType(task))
			fmt.Println()

			hasErrors := false
			hasWarnings := false

			// Check file availability
			if task.File != "" {
				fmt.Println("File:")
				fmt.Printf("  Path: %s\n", task.File)
				resolved := resolvePath(task.File, workDir)
				fmt.Printf("  Resolved: %s\n", resolved)

				if fileInfo, err := os.Stat(resolved); err == nil {
					fmt.Printf("  ✓ File exists (%.1f KB)\n", float64(fileInfo.Size())/1024)
				} else {
					fmt.Println("  ✗ File not found")
					hasWarnings = true
				}
				fmt.Println()
			}

			// Check command execution
			if task.Command != "" {
				fmt.Println("Command:")
				shell := task.Shell
				if shell == "" {
					shell = "sh"
				}
				fmt.Printf("  Shell: %s\n", shell)

				timeout := task.CommandTimeout
				if timeout == 0 {
					timeout = 30
				}
				fmt.Printf("  Timeout: %d seconds\n", timeout)
				fmt.Printf("  Command: %s\n", task.Command)

				// Try to find shell binary
				shellBin, err := exec.LookPath(shell)
				if err != nil {
					fmt.Printf("  ✗ Shell not found: %s\n", shell)
					hasErrors = true
				} else {
					fmt.Printf("  ✓ Shell available: %s\n", shellBin)
				}
				fmt.Println()
			}

			// Check prompt template
			if task.Prompt != "" {
				fmt.Println("Prompt template:")

				// Check for placeholders
				placeholders := findPlaceholders(task.Prompt)
				if len(placeholders) > 0 {
					fmt.Printf("  ✓ Uses placeholders: %s\n", strings.Join(placeholders, ", "))

					// Validate placeholders
					validPlaceholders := []string{"{file}", "{file_contents}", "{command}", "{command_output}", "{instructions}", "{date}"}
					for _, ph := range placeholders {
						isValid := false
						for _, valid := range validPlaceholders {
							if ph == valid {
								isValid = true
								break
							}
						}
						if !isValid {
							fmt.Printf("  ⚠ Unknown placeholder %s\n", ph)
							fmt.Println("    Valid: {file}, {file_contents}, {command}, {command_output}, {instructions}, {date}")
							hasWarnings = true
						}
					}

					// Check if {instructions} placeholder is present (recommended for tasks)
					if !strings.Contains(task.Prompt, "{instructions}") {
						fmt.Println("  ℹ Template doesn't use {instructions} placeholder")
						fmt.Println("    Tasks typically use {instructions} for dynamic user input")
					}

					// Check if placeholders match configuration
					if strings.Contains(task.Prompt, "{file}") || strings.Contains(task.Prompt, "{file_contents}") {
						if task.File == "" {
							fmt.Println("  ⚠ Prompt uses {file} or {file_contents} but no file configured")
							hasWarnings = true
						}
					}
					if strings.Contains(task.Prompt, "{command}") || strings.Contains(task.Prompt, "{command_output}") {
						if task.Command == "" {
							fmt.Println("  ⚠ Prompt uses {command} or {command_output} but no command configured")
							hasWarnings = true
						}
					}
				} else {
					fmt.Printf("  ✓ Valid inline prompt (%d characters)\n", len(task.Prompt))
				}
				fmt.Println()
			}

			// Check UTD requirement
			if task.File == "" && task.Command == "" && task.Prompt == "" {
				fmt.Println("✗ No task prompt defined")
				fmt.Println("  At least one field required: file, command, or prompt")
				hasErrors = true
				fmt.Println()
			}

			// Summary
			if hasErrors {
				fmt.Printf("✗ Task '%s' has errors\n", taskName)
				fmt.Println("  Fix configuration: start config task edit", taskName)
				return fmt.Errorf("configuration errors")
			} else if hasWarnings {
				fmt.Printf("⚠ Task '%s' has warnings (see above)\n", taskName)
			} else {
				fmt.Printf("✓ Task '%s' is configured correctly\n", taskName)
			}

			return nil
		},
	}

	return cmd
}

// getTaskSourceType returns a human-readable source type for a task
func getTaskSourceType(task domain.Task) string {
	hasFile := task.File != ""
	hasCommand := task.Command != ""
	hasPrompt := task.Prompt != ""

	if hasFile && hasCommand && hasPrompt {
		return "Combination (file + command + template)"
	} else if hasFile && hasCommand {
		return "Combination (file + command)"
	} else if hasFile && hasPrompt {
		return "File-based"
	} else if hasCommand && hasPrompt {
		return "Command-based"
	} else if hasFile {
		return "File only"
	} else if hasCommand {
		return "Command only"
	} else if hasPrompt {
		return "Inline prompt"
	}
	return "Invalid (no UTD fields)"
}
