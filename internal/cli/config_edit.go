package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/grantcarthew/start/internal/config"
	"github.com/grantcarthew/start/internal/domain"
	"github.com/spf13/cobra"
)

// NewConfigEditCommand creates the config edit command
func NewConfigEditCommand(configLoader *config.Loader, validator *config.Validator) *cobra.Command {
	var localFlag bool

	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit configuration file",
		Long:  "Open settings configuration file (config.toml) in editor. For editing other config files, use: 'start config agent edit', 'start config context edit', etc.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			tomlHelper := config.NewTOMLHelper(configLoader.GetFS())

			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			// Get config directories
			globalDir, err := tomlHelper.GetGlobalDir()
			if err != nil {
				return err
			}
			localDir := tomlHelper.GetLocalDir(workDir)

			globalPath := tomlHelper.GetConfigPath(globalDir)
			localPath := tomlHelper.GetConfigPath(localDir)

			// Check which configs exist
			globalExists := fileExists(globalPath)
			localExists := fileExists(localPath)

			var configPath string
			var scope string

			// Determine which config to edit
			if localFlag {
				configPath = localPath
				scope = "local"
				if !localExists {
					// Create the directory if it doesn't exist
					if err := os.MkdirAll(localDir, 0755); err != nil {
						return fmt.Errorf("failed to create local config directory: %w", err)
					}
				}
			} else if globalExists && localExists {
				// Ask which to edit
				selection, err := promptConfigSelection(globalPath, localPath)
				if err != nil {
					return err
				}
				if selection == 1 {
					configPath = globalPath
					scope = "global"
				} else {
					configPath = localPath
					scope = "local"
				}
			} else if globalExists {
				configPath = globalPath
				scope = "global"
			} else if localExists {
				configPath = localPath
				scope = "local"
			} else {
				// Neither exists, ask which to create
				selection, err := promptConfigCreation()
				if err != nil {
					return err
				}
				if selection == 1 {
					configPath = globalPath
					scope = "global"
					// Create the directory if it doesn't exist
					if err := os.MkdirAll(globalDir, 0755); err != nil {
						return fmt.Errorf("failed to create global config directory: %w", err)
					}
				} else {
					configPath = localPath
					scope = "local"
					// Create the directory if it doesn't exist
					if err := os.MkdirAll(localDir, 0755); err != nil {
						return fmt.Errorf("failed to create local config directory: %w", err)
					}
				}
			}

			// Detect editor
			editor := os.Getenv("VISUAL")
			showEditorMessage := false
			if editor == "" {
				editor = os.Getenv("EDITOR")
				if editor == "" {
					editor = "vi"
					showEditorMessage = true
				}
			}

			// Show opening message
			fmt.Printf("Opening %s in %s...\n", configPath, editor)
			if showEditorMessage {
				fmt.Println("Set $EDITOR to use your preferred editor.")
			}
			fmt.Println()

			// Open editor
			editorCmd := exec.Command(editor, configPath)
			editorCmd.Stdin = os.Stdin
			editorCmd.Stdout = os.Stdout
			editorCmd.Stderr = os.Stderr

			if err := editorCmd.Run(); err != nil {
				return fmt.Errorf("editor failed: %w", err)
			}

			// Validate the config after editing
			fmt.Println()
			fmt.Println("Validating configuration...")

			// Load and validate the edited config
			var cfg domain.Config
			if scope == "global" {
				cfg, err = configLoader.LoadGlobal()
			} else {
				cfg, err = configLoader.LoadLocal(workDir)
			}

			if err != nil {
				fmt.Fprintf(os.Stderr, "\n⚠ Configuration has errors:\n%v\n\n", err)
				fmt.Fprintf(os.Stderr, "Use 'start config edit%s' to fix the errors.\n",
					map[bool]string{true: " --local", false: ""}[scope == "local"])
				return nil // Don't return error, file is already saved
			}

			// Run validation
			if err := validator.Validate(cfg); err != nil {
				fmt.Printf("\n⚠ Warnings found:\n\n%v\n\n", err)
				fmt.Printf("Changes saved to %s\n", configPath)
				fmt.Println()
				fmt.Println("Note: Warnings don't prevent using start, but may affect functionality.")
				return nil
			}

			fmt.Println("✓ Configuration is valid")
			fmt.Println()
			fmt.Printf("Changes saved to %s\n", configPath)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&localFlag, "local", "l", false, "Edit local configuration")

	return cmd
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// promptConfigSelection prompts user to select which config to edit
func promptConfigSelection(globalPath, localPath string) (int, error) {
	fmt.Println("Edit configuration")
	fmt.Println("─────────────────────────────────────────────────")
	fmt.Println()
	fmt.Println("Both global and local configs exist:")
	fmt.Printf("  1) Global: %s\n", globalPath)
	fmt.Printf("  2) Local:  %s\n", localPath)
	fmt.Println()
	fmt.Print("Select [1-2] (or 'q' to quit): ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimSpace(input)

	if input == "q" || input == "Q" {
		return 0, fmt.Errorf("cancelled by user")
	}

	switch input {
	case "1":
		return 1, nil
	case "2":
		return 2, nil
	default:
		return 0, fmt.Errorf("invalid selection: %s", input)
	}
}

// promptConfigCreation prompts user to select which config to create
func promptConfigCreation() (int, error) {
	fmt.Println("Create configuration")
	fmt.Println("─────────────────────────────────────────────────")
	fmt.Println()
	fmt.Println("No configuration files found:")
	fmt.Println("  1) Create global config")
	fmt.Println("  2) Create local config")
	fmt.Println()
	fmt.Print("Select [1-2] (or 'q' to quit): ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimSpace(input)

	if input == "q" || input == "Q" {
		return 0, fmt.Errorf("cancelled by user")
	}

	switch input {
	case "1":
		return 1, nil
	case "2":
		return 2, nil
	default:
		return 0, fmt.Errorf("invalid selection: %s", input)
	}
}
