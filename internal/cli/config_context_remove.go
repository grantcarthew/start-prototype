package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/grantcarthew/start/internal/config"
	"github.com/grantcarthew/start/internal/domain"
	"github.com/spf13/cobra"
)

// NewConfigContextRemoveCommand creates the config context remove command
func NewConfigContextRemoveCommand(configLoader *config.Loader) *cobra.Command {
	var localOnly bool

	cmd := &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove context from configuration",
		Long:  "Remove a context document from the configuration file",
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
				// Remove from local only
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
					// Exists in both - ask which to remove
					fmt.Println("Context exists in both global and local configs.")
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
						if err := removeContextFromScope(contextName, globalDir, "global", globalContexts, prompter, tomlHelper, backupHelper); err != nil {
							return err
						}
						return removeContextFromScope(contextName, localDir, "local", localContexts, prompter, tomlHelper, backupHelper)
					} else if choice == "global" {
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

			return removeContextFromScope(contextName, targetDir, scope, contexts, prompter, tomlHelper, backupHelper)
		},
	}

	cmd.Flags().BoolVarP(&localOnly, "local", "l", false, "Remove from local config only")

	return cmd
}

// removeContextFromScope removes a context from a specific scope
func removeContextFromScope(
	contextName string,
	targetDir string,
	scope string,
	contexts map[string]domain.Context,
	prompter *PromptHelper,
	tomlHelper *config.TOMLHelper,
	backupHelper *config.BackupHelper,
) error {
	// Check if context exists
	ctx, exists := contexts[contextName]
	if !exists {
		return fmt.Errorf("context '%s' not found in %s config", contextName, scope)
	}

	// Warn if required
	if ctx.Required {
		fmt.Printf("⚠ Warning: '%s' is marked as required context.\n", contextName)
		fmt.Println("  Removing it may affect agent behavior.")
		fmt.Println()
	}

	// Confirm removal
	confirmed, err := prompter.AskYesNo(fmt.Sprintf("Remove context '%s' from %s config?", contextName, scope), false)
	if err != nil {
		return err
	}

	if !confirmed {
		fmt.Println()
		fmt.Printf("Context '%s' not removed.\n", contextName)
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

	// Remove context and save
	delete(contexts, contextName)
	if err := tomlHelper.WriteContextsFile(targetDir, contexts); err != nil {
		return fmt.Errorf("failed to write contexts file: %w", err)
	}

	fmt.Printf("Removing context '%s' from %s...\n", contextName, contextsPath)
	prompter.PrintSuccess(fmt.Sprintf("Context '%s' removed successfully", contextName))
	fmt.Println()
	fmt.Printf("Use 'start config context list' to see remaining contexts.\n")

	return nil
}
