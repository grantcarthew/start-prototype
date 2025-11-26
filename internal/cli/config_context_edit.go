package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/grantcarthew/start/internal/config"
	"github.com/grantcarthew/start/internal/domain"
	"github.com/spf13/cobra"
)

// NewConfigContextEditCommand creates the config context edit command
func NewConfigContextEditCommand(configLoader *config.Loader) *cobra.Command {
	var localOnly bool

	cmd := &cobra.Command{
		Use:   "edit <name>",
		Short: "Edit context configuration interactively",
		Long:  "Interactive wizard to modify an existing context document configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			contextName := args[0]
			prompter := NewPromptHelper()
			tomlHelper := config.NewTOMLHelper(configLoader.GetFS())
			backupHelper := config.NewBackupHelper(configLoader.GetFS())

			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			// Determine which config file contains the context
			var targetDir string
			var scope string
			var contexts map[string]domain.Context

			globalDir, err := tomlHelper.GetGlobalDir()
			if err != nil {
				return err
			}
			localDir := tomlHelper.GetLocalDir(workDir)

			if localOnly {
				// Edit local only
				contexts, err = tomlHelper.ReadContextsFile(localDir)
				if err != nil {
					return fmt.Errorf("failed to read local contexts: %w", err)
				}
				targetDir = localDir
				scope = "local"
			} else {
				// Check both global and local
				globalContexts, err := tomlHelper.ReadContextsFile(globalDir)
				if err != nil {
					return fmt.Errorf("failed to read global contexts: %w", err)
				}

				localContexts, err := tomlHelper.ReadContextsFile(localDir)
				if err != nil {
					return fmt.Errorf("failed to read local contexts: %w", err)
				}

				// Check where context exists
				_, inGlobal := globalContexts[contextName]
				_, inLocal := localContexts[contextName]

				if !inGlobal && !inLocal {
					fmt.Printf("Error: Context '%s' not found in configuration.\n\n", contextName)
					fmt.Println("Use 'start config context list' to see available contexts.")
					return fmt.Errorf("context not found")
				}

				if inGlobal && inLocal {
					// Exists in both - ask which to edit
					fmt.Println("Context exists in both global and local configs.")
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
						contexts = globalContexts
						scope = "global"
					} else {
						targetDir = localDir
						contexts = localContexts
						scope = "local"
					}
				} else if inLocal {
					targetDir = localDir
					contexts = localContexts
					scope = "local"
				} else {
					targetDir = globalDir
					contexts = globalContexts
					scope = "global"
				}
			}

			// Get existing context
			existingContext, exists := contexts[contextName]
			if !exists {
				return fmt.Errorf("context '%s' not found in %s config", contextName, scope)
			}

			prompter.PrintHeader(fmt.Sprintf("Edit context: %s (%s)", contextName, scope))

			// Show current configuration
			fmt.Println("Current configuration:")
			if existingContext.Description != "" {
				fmt.Printf("  Description: %s\n", existingContext.Description)
			}
			if existingContext.File != "" {
				fmt.Printf("  File: %s\n", existingContext.File)
			}
			if existingContext.Command != "" {
				fmt.Printf("  Command: %s\n", existingContext.Command)
			}
			if existingContext.Prompt != "" {
				promptPreview := existingContext.Prompt
				if len(promptPreview) > 60 {
					promptPreview = promptPreview[:57] + "..."
				}
				fmt.Printf("  Prompt: %s\n", promptPreview)
			}
			if existingContext.Required {
				fmt.Println("  Required: yes")
			} else {
				fmt.Println("  Required: no")
			}
			if existingContext.Shell != "" {
				fmt.Printf("  Shell: %s\n", existingContext.Shell)
			} else {
				fmt.Println("  Shell: (default)")
			}
			if existingContext.CommandTimeout > 0 {
				fmt.Printf("  Timeout: %d\n", existingContext.CommandTimeout)
			} else {
				fmt.Println("  Timeout: (default)")
			}
			fmt.Println()
			fmt.Println("Press enter to keep current value, or type new value:")
			fmt.Println()

			// Create updated context (start with existing)
			updatedContext := existingContext

			// Description
			descPrompt := "Description"
			if existingContext.Description != "" {
				descPrompt = fmt.Sprintf("Description [%s]", existingContext.Description)
			}
			description, err := prompter.Ask(descPrompt + ": ")
			if err != nil {
				return err
			}
			description = strings.TrimSpace(description)
			if description != "" {
				updatedContext.Description = description
			}

			// File path
			filePrompt := "File path"
			if existingContext.File != "" {
				filePrompt = fmt.Sprintf("File path [%s]", existingContext.File)
			}
			filePath, err := prompter.Ask(filePrompt + ": ")
			if err != nil {
				return err
			}
			filePath = strings.TrimSpace(filePath)
			if filePath != "" {
				updatedContext.File = filePath

				// Check if file exists (warning only)
				resolved := resolvePath(filePath, workDir)
				if _, err := os.Stat(resolved); err != nil {
					fmt.Printf("⚠ Warning: File does not exist: %s\n", filePath)
					cont, err := prompter.AskYesNo("Continue anyway?", true)
					if err != nil {
						return err
					}
					if !cont {
						updatedContext.File = existingContext.File
					}
				}
			}

			// Command
			fmt.Println()
			addCommand, err := prompter.AskYesNo("Add/modify command for dynamic content?", existingContext.Command != "")
			if err != nil {
				return err
			}

			if addCommand {
				cmdPrompt := "Command"
				if existingContext.Command != "" {
					cmdPrompt = fmt.Sprintf("Command [%s]", existingContext.Command)
				}
				cmdStr, err := prompter.Ask(cmdPrompt + ": ")
				if err != nil {
					return err
				}
				cmdStr = strings.TrimSpace(cmdStr)
				if cmdStr != "" {
					updatedContext.Command = cmdStr
					prompter.PrintSuccess("Valid command")
				}
			} else {
				// Clear command
				updatedContext.Command = ""
			}

			// Prompt template
			fmt.Println()
			promptPrompt := "Prompt template"
			if existingContext.Prompt != "" {
				promptPreview := existingContext.Prompt
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
				updatedContext.Prompt = promptTemplate

				// Validate template
				if updatedContext.File != "" && (strings.Contains(promptTemplate, "{file}") || strings.Contains(promptTemplate, "{file_contents}")) {
					prompter.PrintSuccess("Valid template (uses file placeholders)")
				}
				if updatedContext.Command != "" && (strings.Contains(promptTemplate, "{command}") || strings.Contains(promptTemplate, "{command_output}")) {
					prompter.PrintSuccess("Valid template (uses command placeholders)")
				}
			}

			// Required
			fmt.Println()
			requiredPrompt := fmt.Sprintf("Required [%v]", existingContext.Required)
			requiredStr, err := prompter.Ask(requiredPrompt + ": ")
			if err != nil {
				return err
			}
			requiredStr = strings.TrimSpace(strings.ToLower(requiredStr))
			if requiredStr != "" {
				if requiredStr == "y" || requiredStr == "yes" || requiredStr == "true" {
					updatedContext.Required = true
				} else if requiredStr == "n" || requiredStr == "no" || requiredStr == "false" {
					updatedContext.Required = false
				}
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
				if existingContext.Shell != "" {
					shellPrompt = fmt.Sprintf("Shell override [%s]", existingContext.Shell)
				}
				shell, err := prompter.Ask(shellPrompt + " (or enter for default): ")
				if err != nil {
					return err
				}
				shell = strings.TrimSpace(shell)
				if shell != "" {
					updatedContext.Shell = shell
				}

				// Command timeout
				timeoutPrompt := "Command timeout in seconds"
				if existingContext.CommandTimeout > 0 {
					timeoutPrompt = fmt.Sprintf("Command timeout in seconds [%d]", existingContext.CommandTimeout)
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
					updatedContext.CommandTimeout = timeoutInt
				}
			}

			// Check if anything changed
			if updatedContext == existingContext {
				fmt.Println()
				fmt.Println("No changes detected.")
				fmt.Println()
				fmt.Printf("Context '%s' not modified.\n", contextName)
				return nil
			}

			// Backup existing config
			fmt.Println()
			contextsPath := filepath.Join(targetDir, "contexts.toml")
			fmt.Println("Backing up config to contexts.YYYY-MM-DD-HHMMSS.toml...")
			backupPath, err := backupHelper.CreateBackup(contextsPath)
			if err != nil {
				return fmt.Errorf("failed to backup config: %w", err)
			}
			prompter.PrintSuccess(fmt.Sprintf("Backup created: %s", filepath.Base(backupPath)))
			fmt.Println()

			// Update context and save
			contexts[contextName] = updatedContext
			if err := tomlHelper.WriteContextsFile(targetDir, contexts); err != nil {
				return fmt.Errorf("failed to write contexts file: %w", err)
			}

			fmt.Printf("Saving changes to %s...\n", contextsPath)
			prompter.PrintSuccess(fmt.Sprintf("Context '%s' updated successfully", contextName))
			fmt.Println()
			fmt.Printf("Use 'start config context list' to see changes.\n")
			fmt.Printf("Use 'start config context test %s' to validate.\n", contextName)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&localOnly, "local", "l", false, "Edit in local config only")

	return cmd
}
