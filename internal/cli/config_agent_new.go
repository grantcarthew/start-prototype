package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/grantcarthew/start/internal/config"
	"github.com/grantcarthew/start/internal/domain"
	"github.com/spf13/cobra"
)

// NewConfigAgentNewCommand creates the config agent new command
func NewConfigAgentNewCommand(configLoader *config.Loader) *cobra.Command {
	var localOnly bool

	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create new agent interactively",
		Long:  "Interactive wizard to create a new agent configuration",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			prompter := NewPromptHelper()
			tomlHelper := config.NewTOMLHelper(configLoader.GetFS())
			backupHelper := config.NewBackupHelper(configLoader.GetFS())

			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			prompter.PrintHeader("Add new agent")

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

			// Read existing agents to check for duplicates
			existingAgents, err := tomlHelper.ReadAgentsFile(targetDir)
			if err != nil {
				return fmt.Errorf("failed to read existing agents: %w", err)
			}

			// Agent name
			var agentName string
			for {
				agentName, err = prompter.AskValidatedName("\nAgent name: ")
				if err != nil {
					return err
				}

				// Check for duplicate
				if _, exists := existingAgents[agentName]; exists {
					prompter.PrintError(fmt.Sprintf("Agent '%s' already exists in %s config.", agentName, scope))
					fmt.Println()
					fmt.Println("Use 'start config agent edit", agentName, "' to modify existing agent.")
					fmt.Println()
					continue
				}

				break
			}

			// Create agent struct
			agent := domain.Agent{
				Name:   agentName,
				Models: make(map[string]string),
			}

			// Description (optional)
			description, err := prompter.AskOptional("\nDescription")
			if err != nil {
				return err
			}
			agent.Description = description

			// Binary name
			bin, err := prompter.Ask("\nBinary name (e.g., claude, openai): ")
			if err != nil {
				return err
			}
			agent.Bin = bin

			// Command template
			fmt.Println("\nCommand template")
			fmt.Println("Available placeholders: {bin}, {model}, {role}, {role_file}, {prompt}, {date}")
			fmt.Println()
			fmt.Println("Example for Claude:")
			fmt.Println(`  {bin} --model {model} --append-system-prompt '{role}' '{prompt}'`)
			fmt.Println()
			fmt.Println("Example for Gemini (file-based role):")
			fmt.Println(`  GEMINI_SYSTEM_MD="{role_file}" {bin} --model {model} --prompt-interactive '{prompt}'`)
			fmt.Println()
			fmt.Println("Important: Use single quotes around '{role}' and '{prompt}' for bash safety")
			fmt.Println()
			command, err := prompter.Ask("Command: ")
			if err != nil {
				return err
			}
			agent.Command = command

			// URL (optional)
			url, err := prompter.AskOptional("\nURL")
			if err != nil {
				return err
			}
			agent.URL = url

			// Models URL (optional)
			modelsURL, err := prompter.AskOptional("Models URL")
			if err != nil {
				return err
			}
			agent.ModelsURL = modelsURL

			// Add models
			addModels, err := prompter.AskYesNo("\nAdd models?", true)
			if err != nil {
				return err
			}

			if addModels {
				fmt.Println("\nEnter models (one per line, format: name=full-model-id)")
				fmt.Println("Examples:")
				fmt.Println("  sonnet=claude-sonnet-4-20250929")
				fmt.Println("  opus=claude-opus-4-20250514")
				fmt.Println("Press enter with empty line to finish.")
				fmt.Println()

				for {
					model, err := prompter.Ask("Model (name=id): ")
					if err != nil {
						return err
					}

					if model == "" {
						break
					}

					// Parse name=id
					parts := splitOnce(model, "=")
					if len(parts) != 2 {
						fmt.Println("Invalid format. Use: name=full-model-id")
						continue
					}

					modelName := parts[0]
					modelID := parts[1]

					if err := prompter.ValidateName(modelName); err != nil {
						fmt.Printf("Invalid model name: %v\n", err)
						continue
					}

					agent.Models[modelName] = modelID
					prompter.PrintSuccess(fmt.Sprintf("Added model: %s = %s", modelName, modelID))
				}

				// Set default model
				if len(agent.Models) > 0 {
					defaultModel, err := prompter.Ask("\nDefault model name (or enter for first): ")
					if err != nil {
						return err
					}

					if defaultModel != "" {
						if _, ok := agent.Models[defaultModel]; !ok {
							prompter.PrintWarning(fmt.Sprintf("Model '%s' not found, using first model", defaultModel))
						} else {
							agent.DefaultModel = defaultModel
						}
					}
				}
			}

			// Create backup if file exists
			configPath := filepath.Join(targetDir, "agents.toml")
			if tomlHelper.GetFS().Exists(configPath) {
				fmt.Println()
				backupPath, err := backupHelper.CreateBackup(configPath)
				if err != nil {
					return fmt.Errorf("failed to create backup: %w", err)
				}
				prompter.PrintSuccess(fmt.Sprintf("Backup created: %s", filepath.Base(backupPath)))
			}

			// Add to existing agents
			existingAgents[agentName] = agent

			// Write agents file
			if err := tomlHelper.WriteAgentsFile(targetDir, existingAgents); err != nil {
				return fmt.Errorf("failed to write agents file: %w", err)
			}

			fmt.Println()
			prompter.PrintSuccess(fmt.Sprintf("Agent '%s' added to %s config", agentName, scope))
			fmt.Println()
			fmt.Printf("Use 'start config agent list' to see all agents.\n")
			fmt.Printf("Use 'start config agent test %s' to verify.\n", agentName)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&localOnly, "local", "l", false, "Add to local config")

	return cmd
}

// splitOnce splits a string on the first occurrence of sep
func splitOnce(s, sep string) []string {
	parts := make([]string, 0, 2)
	idx := 0
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			parts = append(parts, s[idx:i])
			parts = append(parts, s[i+len(sep):])
			return parts
		}
	}
	return []string{s}
}
