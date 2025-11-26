package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/grantcarthew/start/internal/config"
	"github.com/grantcarthew/start/internal/domain"
	"github.com/spf13/cobra"
)

// NewConfigRoleCommand creates the config role command
func NewConfigRoleCommand(configLoader *config.Loader, validator *config.Validator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role",
		Short: "Manage role configurations",
		Long:  "Commands for managing role configurations in global and local config files",
	}

	// Add subcommands
	cmd.AddCommand(NewConfigRoleListCommand(configLoader))
	cmd.AddCommand(NewConfigRoleShowCommand(configLoader))
	cmd.AddCommand(NewConfigRoleTestCommand(configLoader))
	cmd.AddCommand(NewConfigRoleNewCommand(configLoader))
	cmd.AddCommand(NewConfigRoleEditCommand(configLoader))
	cmd.AddCommand(NewConfigRoleRemoveCommand(configLoader))
	cmd.AddCommand(NewConfigRoleDefaultCommand(configLoader))

	return cmd
}

// NewConfigRoleListCommand creates the config role list command
func NewConfigRoleListCommand(configLoader *config.Loader) *cobra.Command {
	var localOnly bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Display all configured roles",
		Long:  "List all roles defined in global and/or local configuration files",
		RunE: func(cmd *cobra.Command, args []string) error {
			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			var roles map[string]domain.Role
			var scope string

			if localOnly {
				// Load local only
				localCfg, err := configLoader.LoadLocal(workDir)
				if err != nil {
					return fmt.Errorf("failed to load local config: %w", err)
				}
				roles = localCfg.Roles
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
					roles = mergedCfg.Roles
					scope = "merged"
				} else {
					roles = globalCfg.Roles
					scope = "global"
				}
			}

			if len(roles) == 0 {
				fmt.Println("No roles configured.")
				fmt.Println()
				fmt.Println("Run 'start init' to set up roles, or")
				fmt.Println("use 'start assets add' to install from catalog or 'start config role new' to create custom.")
				return nil
			}

			// Sort roles by name for consistent output
			names := make([]string, 0, len(roles))
			for name := range roles {
				names = append(names, name)
			}
			sort.Strings(names)

			fmt.Printf("Configured roles (%s):\n", scope)
			fmt.Println("═══════════════════════════════════════════════════════════")
			fmt.Println()
			for _, name := range names {
				role := roles[name]
				fmt.Printf("%s\n", name)
				if role.Description != "" {
					fmt.Printf("  %s\n", role.Description)
				}

				// Show source type
				sourceType := getRoleSourceType(role)
				fmt.Printf("  Type: %s\n", sourceType)

				if role.File != "" {
					fmt.Printf("  File: %s\n", role.File)
				}
				if role.Command != "" {
					fmt.Printf("  Command: %s\n", role.Command)
				}

				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&localOnly, "local", "l", false, "List local roles only")

	return cmd
}

// NewConfigRoleShowCommand creates the config role show command
func NewConfigRoleShowCommand(configLoader *config.Loader) *cobra.Command {
	var localOnly bool

	cmd := &cobra.Command{
		Use:   "show <name>",
		Short: "Display role configuration",
		Long:  "Show detailed configuration for a specific role",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			roleName := args[0]
			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			var roles map[string]domain.Role
			var scope string

			if localOnly {
				localCfg, err := configLoader.LoadLocal(workDir)
				if err != nil {
					return fmt.Errorf("failed to load local config: %w", err)
				}
				roles = localCfg.Roles
				scope = "local"
			} else {
				globalCfg, err := configLoader.LoadGlobal()
				if err != nil {
					return fmt.Errorf("failed to load global config: %w", err)
				}

				localCfg, err := configLoader.LoadLocal(workDir)
				if err == nil {
					mergedCfg := config.Merge(globalCfg, localCfg)
					roles = mergedCfg.Roles

					// Determine scope
					if _, hasLocal := localCfg.Roles[roleName]; hasLocal {
						scope = "local (overrides global)"
					} else {
						scope = "global"
					}
				} else {
					roles = globalCfg.Roles
					scope = "global"
				}
			}

			role, exists := roles[roleName]
			if !exists {
				return fmt.Errorf("role '%s' not found in configuration.\n\nUse 'start config role list' to see available roles.", roleName)
			}

			// Display role configuration
			fmt.Printf("Role configuration: %s (%s)\n", roleName, scope)
			fmt.Println("═══════════════════════════════════════════════════════════")
			fmt.Println()

			if role.Description != "" {
				fmt.Printf("Description: %s\n", role.Description)
				fmt.Println()
			}

			sourceType := getRoleSourceType(role)
			fmt.Printf("Type: %s\n", sourceType)
			fmt.Println()

			if role.File != "" {
				fmt.Println("File:")
				fmt.Printf("  Path: %s\n", role.File)

				// Try to resolve and check file
				resolved := resolvePath(role.File, workDir)
				fmt.Printf("  Resolved: %s\n", resolved)

				if fileInfo, err := os.Stat(resolved); err == nil {
					fmt.Printf("  ✓ File exists (%d bytes)\n", fileInfo.Size())
				} else {
					fmt.Println("  ✗ File not found")
				}
				fmt.Println()
			}

			if role.Command != "" {
				fmt.Println("Command:")
				fmt.Printf("  %s\n", role.Command)
				if role.Shell != "" {
					fmt.Printf("  Shell: %s\n", role.Shell)
				}
				if role.CommandTimeout > 0 {
					fmt.Printf("  Timeout: %d seconds\n", role.CommandTimeout)
				}
				fmt.Println()
			}

			if role.Prompt != "" {
				fmt.Println("Prompt template:")
				// Show first few lines
				lines := strings.Split(role.Prompt, "\n")
				maxLines := 10
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

	cmd.Flags().BoolVarP(&localOnly, "local", "l", false, "Show local role only")

	return cmd
}

// NewConfigRoleTestCommand creates the config role test command
func NewConfigRoleTestCommand(configLoader *config.Loader) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test <name>",
		Short: "Test role configuration and file availability",
		Long:  "Validate role configuration without executing it",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			roleName := args[0]
			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			// Load merged config
			globalCfg, err := configLoader.LoadGlobal()
			if err != nil {
				return fmt.Errorf("failed to load global config: %w", err)
			}

			var roles map[string]domain.Role
			var scope string
			localCfg, err := configLoader.LoadLocal(workDir)
			if err == nil {
				mergedCfg := config.Merge(globalCfg, localCfg)
				roles = mergedCfg.Roles

				if _, hasLocal := localCfg.Roles[roleName]; hasLocal {
					scope = "local (overrides global)"
				} else {
					scope = "global"
				}
			} else {
				roles = globalCfg.Roles
				scope = "global"
			}

			role, exists := roles[roleName]
			if !exists {
				fmt.Fprintf(os.Stderr, "Error: Role '%s' not found in configuration.\n\n", roleName)
				fmt.Fprintln(os.Stderr, "Use 'start config role list' to see available roles.")
				fmt.Fprintln(os.Stderr, "Use 'start assets add' to install from catalog or 'start config role new' to create custom.")
				return fmt.Errorf("role not found")
			}

			fmt.Printf("Testing role: %s\n", roleName)
			fmt.Println("─────────────────────────────────────────────────")
			fmt.Println()

			fmt.Printf("Effective configuration:\n")
			fmt.Printf("  Scope: %s\n", scope)
			fmt.Printf("  Type: %s\n", getRoleSourceType(role))
			fmt.Println()

			hasErrors := false
			hasWarnings := false

			// Check file availability
			if role.File != "" {
				fmt.Println("File:")
				fmt.Printf("  Path: %s\n", role.File)
				resolved := resolvePath(role.File, workDir)
				fmt.Printf("  Resolved: %s\n", resolved)

				if fileInfo, err := os.Stat(resolved); err == nil {
					fmt.Printf("  ✓ File exists (%d bytes)\n", fileInfo.Size())
					fmt.Printf("  Modified: %s\n", fileInfo.ModTime().Format("2006-01-02 15:04:05"))
				} else {
					fmt.Println("  ✗ File not found")
					hasErrors = true
				}
				fmt.Println()
			}

			// Check command execution
			if role.Command != "" {
				fmt.Println("Command:")
				shell := role.Shell
				if shell == "" {
					shell = "sh"
				}
				fmt.Printf("  Shell: %s\n", shell)

				timeout := role.CommandTimeout
				if timeout == 0 {
					timeout = 30
				}
				fmt.Printf("  Timeout: %d seconds\n", timeout)
				fmt.Printf("  Command: %s\n", role.Command)

				// Try to execute command
				shellBin, err := exec.LookPath(shell)
				if err != nil {
					fmt.Printf("  ✗ Shell not found: %s\n", shell)
					hasErrors = true
				} else {
					fmt.Printf("  ✓ Shell found: %s\n", shellBin)

					// Don't actually execute for test, just validate it's runnable
					fmt.Println("  ℹ Command validation: Not executed (use runtime to test)")
				}
				fmt.Println()
			}

			// Check prompt template
			if role.Prompt != "" {
				fmt.Println("Prompt template:")

				// Check for placeholders
				placeholders := findPlaceholders(role.Prompt)
				if len(placeholders) > 0 {
					fmt.Printf("  ✓ Uses placeholders: %s\n", strings.Join(placeholders, ", "))

					// Validate placeholders
					validPlaceholders := []string{"{file}", "{file_contents}", "{command}", "{command_output}", "{date}"}
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
							fmt.Println("    Valid: {file}, {file_contents}, {command}, {command_output}, {date}")
							hasWarnings = true
						}
					}

					// Check if placeholders match configuration
					if strings.Contains(role.Prompt, "{file_contents}") && role.File == "" {
						fmt.Println("  ⚠ Prompt uses {file_contents} but no file configured")
						hasWarnings = true
					}
					if strings.Contains(role.Prompt, "{command_output}") && role.Command == "" {
						fmt.Println("  ⚠ Prompt uses {command_output} but no command configured")
						hasWarnings = true
					}
				} else {
					fmt.Println("  ✓ Valid prompt (no placeholders)")
				}
				fmt.Println()
			}

			// Check UTD requirement
			if role.File == "" && role.Command == "" && role.Prompt == "" {
				fmt.Println("✗ No content source defined")
				fmt.Println("  At least one field required: file, command, or prompt")
				hasErrors = true
			}

			// Summary
			fmt.Println()
			if hasErrors {
				fmt.Printf("✗ Role '%s' has errors\n", roleName)
				fmt.Println("  Fix: start config role edit", roleName)
				return fmt.Errorf("configuration errors")
			} else if hasWarnings {
				fmt.Printf("⚠ Role '%s' has warnings (see above)\n", roleName)
			} else {
				fmt.Printf("✓ Role '%s' is configured correctly\n", roleName)
			}

			return nil
		},
	}

	return cmd
}

// getRoleSourceType returns a human-readable source type for a role
func getRoleSourceType(role domain.Role) string {
	hasFile := role.File != ""
	hasCommand := role.Command != ""
	hasPrompt := role.Prompt != ""

	if hasFile && hasCommand && hasPrompt {
		return "Combination (file + command + template)"
	} else if hasFile && hasCommand {
		return "Combination (file + command)"
	} else if hasFile && hasPrompt {
		return "File with template"
	} else if hasCommand && hasPrompt {
		return "Command with template"
	} else if hasFile {
		return "File only"
	} else if hasCommand {
		return "Command only"
	} else if hasPrompt {
		return "Inline prompt"
	}
	return "Invalid (no UTD fields)"
}

// resolvePath resolves ~ and relative paths
func resolvePath(path, workDir string) string {
	if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(workDir, path)
}

// findPlaceholders finds all {placeholder} patterns in a string
func findPlaceholders(text string) []string {
	var placeholders []string
	seen := make(map[string]bool)

	parts := strings.Split(text, "{")
	for _, part := range parts[1:] {
		if idx := strings.Index(part, "}"); idx > 0 {
			ph := "{" + part[:idx] + "}"
			if !seen[ph] {
				placeholders = append(placeholders, ph)
				seen[ph] = true
			}
		}
	}

	return placeholders
}
