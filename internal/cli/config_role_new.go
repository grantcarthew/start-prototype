package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/grantcarthew/start/internal/config"
	"github.com/grantcarthew/start/internal/domain"
	"github.com/spf13/cobra"
)

// NewConfigRoleNewCommand creates the config role new command
func NewConfigRoleNewCommand(configLoader *config.Loader) *cobra.Command {
	var localOnly bool

	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create new role interactively",
		Long:  "Interactive wizard to create a new role configuration",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			prompter := NewPromptHelper()
			tomlHelper := config.NewTOMLHelper(configLoader.GetFS())
			backupHelper := config.NewBackupHelper(configLoader.GetFS())

			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			prompter.PrintHeader("Add new role")

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

			// Read existing roles to check for duplicates
			existingRoles, err := tomlHelper.ReadRolesFile(targetDir)
			if err != nil {
				return fmt.Errorf("failed to read existing roles: %w", err)
			}

			// Role name
			var roleName string
			for {
				roleName, err = prompter.AskValidatedName("\nRole name: ")
				if err != nil {
					return err
				}

				// Check for duplicate
				if _, exists := existingRoles[roleName]; exists {
					prompter.PrintError(fmt.Sprintf("Role '%s' already exists in %s config.", roleName, scope))
					fmt.Println()
					fmt.Println("Use 'start config role edit", roleName, "' to modify existing role.")
					fmt.Println()
					continue
				}

				break
			}

			// Create role struct
			role := domain.Role{
				Name: roleName,
			}

			// Description (optional)
			description, err := prompter.AskOptional("\nDescription")
			if err != nil {
				return err
			}
			role.Description = description

			// Content source
			fmt.Println()
			contentChoice, err := prompter.AskChoice("Content source:", []string{
				"File path",
				"Command",
				"Inline prompt",
				"Combination",
			})
			if err != nil {
				return err
			}

			switch contentChoice {
			case "File path":
				// File only
				filePath, err := prompter.Ask("\nFile path: ")
				if err != nil {
					return err
				}
				role.File = filePath

				// Check if file exists
				resolved := resolvePath(filePath, workDir)
				if _, err := os.Stat(resolved); err == nil {
					prompter.PrintSuccess("File exists")
				} else {
					prompter.PrintWarning("File does not exist")
					fmt.Println("  Role will fail at runtime if file is not found.")
				}

				// Optional template
				useTemplate, err := prompter.AskYesNo("\nUse prompt template to frame file content?", false)
				if err != nil {
					return err
				}

				if useTemplate {
					fmt.Println("\nEnter prompt template (use {file_contents} placeholder):")
					fmt.Println("Press Ctrl+D or enter empty line to finish.")
					fmt.Println()

					prompt, err := readMultilineInput(prompter)
					if err != nil {
						return err
					}

					if prompt != "" {
						role.Prompt = prompt
						prompter.PrintSuccess("Template configured")
					}
				} else {
					prompter.PrintSuccess("Will use file content directly")
				}

			case "Command":
				// Command only
				command, err := prompter.Ask("\nCommand: ")
				if err != nil {
					return err
				}
				role.Command = command

				// Optional template
				useTemplate, err := prompter.AskYesNo("\nUse prompt template to frame command output?", false)
				if err != nil {
					return err
				}

				if useTemplate {
					fmt.Println("\nEnter prompt template (use {command_output} placeholder):")
					fmt.Println("Press Ctrl+D or enter empty line to finish.")
					fmt.Println()

					prompt, err := readMultilineInput(prompter)
					if err != nil {
						return err
					}

					if prompt != "" {
						role.Prompt = prompt
						prompter.PrintSuccess("Template configured")
					}
				}

			case "Inline prompt":
				// Prompt only
				fmt.Println("\nEnter prompt text:")
				fmt.Println("Press Ctrl+D or enter empty line to finish.")
				fmt.Println()

				prompt, err := readMultilineInput(prompter)
				if err != nil {
					return err
				}

				if prompt == "" {
					return fmt.Errorf("prompt cannot be empty")
				}

				role.Prompt = prompt
				prompter.PrintSuccess(fmt.Sprintf("Valid prompt (%d characters)", len(prompt)))

			case "Combination":
				// File + command + prompt
				filePath, err := prompter.Ask("\nFile path: ")
				if err != nil {
					return err
				}
				role.File = filePath

				// Check if file exists
				resolved := resolvePath(filePath, workDir)
				if _, err := os.Stat(resolved); err == nil {
					prompter.PrintSuccess("File exists")
				} else {
					prompter.PrintWarning("File does not exist")
				}

				// Add command
				addCommand, err := prompter.AskYesNo("\nAdd command for dynamic content?", false)
				if err != nil {
					return err
				}

				if addCommand {
					command, err := prompter.Ask("\nCommand: ")
					if err != nil {
						return err
					}
					role.Command = command
					prompter.PrintSuccess("Command configured")
				}

				// Prompt template
				fmt.Println("\nEnter prompt template:")
				fmt.Println("Available placeholders: {file_contents}, {command_output}, {date}")
				fmt.Println("Press Ctrl+D or enter empty line to finish.")
				fmt.Println()

				prompt, err := readMultilineInput(prompter)
				if err != nil {
					return err
				}

				if prompt != "" {
					role.Prompt = prompt
					prompter.PrintSuccess("Template configured")
				}
			}

			// Advanced options
			advanced, err := prompter.AskYesNo("\nAdvanced options?", false)
			if err != nil {
				return err
			}

			if advanced {
				// Shell override
				shell, err := prompter.AskOptional("\nShell override (or enter for default)")
				if err != nil {
					return err
				}
				if shell != "" {
					role.Shell = shell
				}

				// Command timeout
				timeoutStr, err := prompter.AskOptional("Command timeout in seconds (or enter for default)")
				if err != nil {
					return err
				}
				if timeoutStr != "" {
					var timeout int
					if _, err := fmt.Sscanf(timeoutStr, "%d", &timeout); err == nil && timeout > 0 {
						role.CommandTimeout = timeout
					}
				}
			}

			// Create backup if file exists
			configPath := filepath.Join(targetDir, "roles.toml")
			if tomlHelper.GetFS().Exists(configPath) {
				fmt.Println()
				backupPath, err := backupHelper.CreateBackup(configPath)
				if err != nil {
					return fmt.Errorf("failed to create backup: %w", err)
				}
				prompter.PrintSuccess(fmt.Sprintf("Backup created: %s", filepath.Base(backupPath)))
			}

			// Add to existing roles
			existingRoles[roleName] = role

			// Write roles file
			if err := tomlHelper.WriteRolesFile(targetDir, existingRoles); err != nil {
				return fmt.Errorf("failed to write roles file: %w", err)
			}

			fmt.Println()
			prompter.PrintSuccess(fmt.Sprintf("Role '%s' added to %s config", roleName, scope))
			fmt.Println()
			fmt.Printf("Use 'start config role list' to see all roles.\n")
			fmt.Printf("Use 'start config role test %s' to verify.\n", roleName)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&localOnly, "local", "l", false, "Add to local config")

	return cmd
}

// readMultilineInput reads multiple lines of input until empty line or EOF
func readMultilineInput(prompter *PromptHelper) (string, error) {
	var lines []string

	for {
		line, err := prompter.Ask("")
		if err != nil {
			// EOF (Ctrl+D) - finish input
			break
		}

		if line == "" {
			// Empty line - finish input
			break
		}

		lines = append(lines, line)
	}

	if len(lines) == 0 {
		return "", nil
	}

	result := ""
	for i, line := range lines {
		if i > 0 {
			result += "\n"
		}
		result += line
	}

	return result, nil
}
