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

// NewConfigContextNewCommand creates the config context new command
func NewConfigContextNewCommand(configLoader *config.Loader) *cobra.Command {
	var localOnly bool

	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create new context interactively",
		Long:  "Interactive wizard to create a new context document configuration",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			prompter := NewPromptHelper()
			tomlHelper := config.NewTOMLHelper(configLoader.GetFS())
			backupHelper := config.NewBackupHelper(configLoader.GetFS())

			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			prompter.PrintHeader("Add new context")

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

			// Read existing contexts to check for duplicates
			existingContexts, err := tomlHelper.ReadContextsFile(targetDir)
			if err != nil {
				return fmt.Errorf("failed to read existing contexts: %w", err)
			}

			// Context name
			var contextName string
			for {
				contextName, err = prompter.AskValidatedName("\nContext name: ")
				if err != nil {
					return err
				}

				// Check for duplicate
				if _, exists := existingContexts[contextName]; exists {
					prompter.PrintError(fmt.Sprintf("Context '%s' already exists in %s config.", contextName, scope))
					fmt.Println()
					fmt.Println("Use 'start config context edit", contextName, "' to modify existing context.")
					fmt.Println()
					continue
				}

				break
			}

			// Create context struct
			newContext := domain.Context{}

			// Description (optional)
			description, err := prompter.Ask("Description (optional): ")
			if err != nil {
				return err
			}
			newContext.Description = strings.TrimSpace(description)

			// Content source
			fmt.Println()
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
				// File-based context
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

					newContext.File = filePath

					// Check if file exists (warning only)
					resolved := resolvePath(filePath, workDir)
					if _, err := os.Stat(resolved); err != nil {
						fmt.Printf("⚠ Warning: File does not exist: %s\n", filePath)
						fmt.Println("  Context will be skipped at runtime if file is not found.")
						fmt.Println()
						cont, err := prompter.AskYesNo("Continue anyway?", false)
						if err != nil {
							return err
						}
						if !cont {
							newContext.File = ""
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
				newContext.Prompt = strings.TrimSpace(promptTemplate)

				if newContext.Prompt != "" {
					// Validate template uses {file} placeholder
					if strings.Contains(newContext.Prompt, "{file}") || strings.Contains(newContext.Prompt, "{file_contents}") {
						prompter.PrintSuccess("Valid template (uses {file} or {file_contents} placeholder)")
					} else {
						fmt.Println("⚠ Template doesn't use {file} or {file_contents} placeholder")
					}
				}

			case "Command":
				// Command-based context
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

					newContext.Command = cmdStr
					prompter.PrintSuccess("Valid command")
					break
				}

				// Prompt template
				fmt.Println()
				promptTemplate, err := prompter.Ask("Prompt template: ")
				if err != nil {
					return err
				}
				newContext.Prompt = strings.TrimSpace(promptTemplate)

				if newContext.Prompt != "" {
					// Validate template uses {command} or {command_output} placeholder
					if strings.Contains(newContext.Prompt, "{command}") || strings.Contains(newContext.Prompt, "{command_output}") {
						prompter.PrintSuccess("Valid template (uses {command} or {command_output} placeholder)")
					} else {
						fmt.Println("⚠ Template doesn't use {command} or {command_output} placeholder")
					}
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

					newContext.Prompt = promptText
					prompter.PrintSuccess("Valid prompt")
					break
				}

			case "Combination":
				// File + Command combination
				// File path
				for {
					filePath, err := prompter.Ask("\nFile path (optional, press Enter to skip): ")
					if err != nil {
						return err
					}
					filePath = strings.TrimSpace(filePath)

					if filePath != "" {
						newContext.File = filePath

						// Check if file exists (warning only)
						resolved := resolvePath(filePath, workDir)
						if _, err := os.Stat(resolved); err != nil {
							fmt.Printf("⚠ Warning: File does not exist: %s\n", filePath)
							fmt.Println("  Context will be skipped at runtime if file is not found.")
							fmt.Println()
							cont, err := prompter.AskYesNo("Continue anyway?", false)
							if err != nil {
								return err
							}
							if !cont {
								newContext.File = ""
								continue
							}
						} else {
							prompter.PrintSuccess("File exists")
						}
					}

					break
				}

				// Command
				fmt.Println()
				cmdStr, err := prompter.Ask("Command (optional, press Enter to skip): ")
				if err != nil {
					return err
				}
				cmdStr = strings.TrimSpace(cmdStr)
				if cmdStr != "" {
					newContext.Command = cmdStr
					prompter.PrintSuccess("Valid command")
				}

				// Prompt template
				fmt.Println()
				promptTemplate, err := prompter.Ask("Prompt template: ")
				if err != nil {
					return err
				}
				newContext.Prompt = strings.TrimSpace(promptTemplate)

				// Validate at least one source
				if newContext.File == "" && newContext.Command == "" && newContext.Prompt == "" {
					return fmt.Errorf("at least one content source is required (file, command, or prompt)")
				}
			}

			// Required context?
			fmt.Println()
			required, err := prompter.AskYesNo("Required context?", false)
			if err != nil {
				return err
			}
			newContext.Required = required

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
				newContext.Shell = strings.TrimSpace(shell)

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
					newContext.CommandTimeout = timeoutInt
				}
			}

			// Backup existing config
			fmt.Println()
			contextsPath := filepath.Join(targetDir, "contexts.toml")
			if tomlHelper.GetFS().Exists(contextsPath) {
				fmt.Println("Backing up config to contexts.YYYY-MM-DD-HHMMSS.toml...")
				backupPath, err := backupHelper.CreateBackup(contextsPath)
				if err != nil {
					return fmt.Errorf("failed to backup config: %w", err)
				}
				prompter.PrintSuccess(fmt.Sprintf("Backup created: %s", filepath.Base(backupPath)))
				fmt.Println()
			}

			// Add to contexts map and save
			existingContexts[contextName] = newContext
			if err := tomlHelper.WriteContextsFile(targetDir, existingContexts); err != nil {
				return fmt.Errorf("failed to write contexts file: %w", err)
			}

			fmt.Printf("Saving context '%s' to %s...\n", contextName, contextsPath)
			prompter.PrintSuccess("Context added successfully")
			fmt.Println()
			fmt.Printf("Use 'start config context list' to see all contexts.\n")
			fmt.Printf("Use 'start config context test %s' to verify.\n", contextName)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&localOnly, "local", "l", false, "Create in local config only")

	return cmd
}
