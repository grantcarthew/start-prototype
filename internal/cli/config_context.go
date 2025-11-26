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

// NewConfigContextCommand creates the config context command
func NewConfigContextCommand(configLoader *config.Loader, validator *config.Validator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Manage context document configurations",
		Long:  "Commands for managing context document configurations in global and local config files",
	}

	// Add subcommands
	cmd.AddCommand(NewConfigContextListCommand(configLoader))
	cmd.AddCommand(NewConfigContextShowCommand(configLoader))
	cmd.AddCommand(NewConfigContextTestCommand(configLoader))
	cmd.AddCommand(NewConfigContextNewCommand(configLoader))
	cmd.AddCommand(NewConfigContextEditCommand(configLoader))
	cmd.AddCommand(NewConfigContextRemoveCommand(configLoader))

	return cmd
}

// NewConfigContextListCommand creates the config context list command
func NewConfigContextListCommand(configLoader *config.Loader) *cobra.Command {
	var localOnly bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Display all configured contexts",
		Long:  "List all context documents defined in global and/or local configuration files",
		RunE: func(cmd *cobra.Command, args []string) error {
			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			var contexts map[string]domain.Context
			var scope string

			if localOnly {
				// Load local only
				localCfg, err := configLoader.LoadLocal(workDir)
				if err != nil {
					return fmt.Errorf("failed to load local config: %w", err)
				}
				contexts = localCfg.Contexts
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
					contexts = mergedCfg.Contexts
					scope = "merged"
				} else {
					contexts = globalCfg.Contexts
					scope = "global"
				}
			}

			if len(contexts) == 0 {
				fmt.Println("No contexts configured.")
				fmt.Println()
				fmt.Println("Create contexts: start config context new")
				return nil
			}

			// Separate into required and optional
			var requiredNames, optionalNames []string
			for name, ctx := range contexts {
				if ctx.Required {
					requiredNames = append(requiredNames, name)
				} else {
					optionalNames = append(optionalNames, name)
				}
			}
			sort.Strings(requiredNames)
			sort.Strings(optionalNames)

			fmt.Printf("Configured contexts (%s):\n", scope)
			fmt.Println("═══════════════════════════════════════════════════════════")
			fmt.Println()

			if len(requiredNames) > 0 {
				fmt.Printf("Required contexts (%d):\n", len(requiredNames))
				for _, name := range requiredNames {
					ctx := contexts[name]
					fmt.Printf("  %s\n", name)
					if ctx.Description != "" {
						fmt.Printf("    %s\n", ctx.Description)
					}
					if ctx.File != "" {
						fmt.Printf("    File: %s\n", ctx.File)
					}
					if ctx.Command != "" {
						fmt.Printf("    Command: %s\n", ctx.Command)
					}
					if ctx.Prompt != "" && ctx.File == "" && ctx.Command == "" {
						// Show prompt preview for inline prompts
						preview := ctx.Prompt
						if len(preview) > 60 {
							preview = preview[:57] + "..."
						}
						fmt.Printf("    Prompt: %s\n", preview)
					}
					fmt.Println()
				}
			}

			if len(optionalNames) > 0 {
				fmt.Printf("Optional contexts (%d):\n", len(optionalNames))
				for _, name := range optionalNames {
					ctx := contexts[name]
					fmt.Printf("  %s\n", name)
					if ctx.Description != "" {
						fmt.Printf("    %s\n", ctx.Description)
					}
					if ctx.File != "" {
						fmt.Printf("    File: %s\n", ctx.File)
					}
					if ctx.Command != "" {
						fmt.Printf("    Command: %s\n", ctx.Command)
					}
					if ctx.Prompt != "" && ctx.File == "" && ctx.Command == "" {
						preview := ctx.Prompt
						if len(preview) > 60 {
							preview = preview[:57] + "..."
						}
						fmt.Printf("    Prompt: %s\n", preview)
					}
					fmt.Println()
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&localOnly, "local", "l", false, "List local contexts only")

	return cmd
}

// NewConfigContextShowCommand creates the config context show command
func NewConfigContextShowCommand(configLoader *config.Loader) *cobra.Command {
	var localOnly bool

	cmd := &cobra.Command{
		Use:   "show <name>",
		Short: "Display context configuration",
		Long:  "Show detailed configuration for a specific context document",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			contextName := args[0]
			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			var contexts map[string]domain.Context
			var scope string

			if localOnly {
				localCfg, err := configLoader.LoadLocal(workDir)
				if err != nil {
					return fmt.Errorf("failed to load local config: %w", err)
				}
				contexts = localCfg.Contexts
				scope = "local"
			} else {
				globalCfg, err := configLoader.LoadGlobal()
				if err != nil {
					return fmt.Errorf("failed to load global config: %w", err)
				}

				localCfg, err := configLoader.LoadLocal(workDir)
				if err == nil {
					mergedCfg := config.Merge(globalCfg, localCfg)
					contexts = mergedCfg.Contexts

					// Determine scope
					if _, hasLocal := localCfg.Contexts[contextName]; hasLocal {
						scope = "local"
					} else {
						scope = "global"
					}
				} else {
					contexts = globalCfg.Contexts
					scope = "global"
				}
			}

			ctx, exists := contexts[contextName]
			if !exists {
				fmt.Printf("No context '%s' found in configuration.\n\n", contextName)
				fmt.Println("Configure: start config context new")
				return fmt.Errorf("context not found")
			}

			// Display context configuration
			fmt.Printf("Context configuration: %s (%s)\n", contextName, scope)
			fmt.Println("═══════════════════════════════════════════════════════════")
			fmt.Println()

			if ctx.Description != "" {
				fmt.Printf("Description: %s\n", ctx.Description)
			}
			fmt.Printf("Required: %v\n", ctx.Required)
			fmt.Println()

			sourceType := getContextSourceType(ctx)
			fmt.Printf("Source: %s\n", sourceType)
			fmt.Println()

			if ctx.File != "" {
				fmt.Println("File:")
				fmt.Printf("  Path: %s\n", ctx.File)

				// Try to resolve and check file
				resolved := resolvePath(ctx.File, workDir)
				fmt.Printf("  Resolved: %s\n", resolved)

				if fileInfo, err := os.Stat(resolved); err == nil {
					fmt.Printf("  ✓ File exists (%.1f KB)\n", float64(fileInfo.Size())/1024)
				} else {
					fmt.Println("  ✗ File not found")
				}
				fmt.Println()
			}

			if ctx.Command != "" {
				fmt.Println("Command:")
				if ctx.Shell != "" {
					fmt.Printf("  Shell: %s\n", ctx.Shell)
				} else {
					fmt.Println("  Shell: (default)")
				}
				if ctx.CommandTimeout > 0 {
					fmt.Printf("  Timeout: %d seconds\n", ctx.CommandTimeout)
				} else {
					fmt.Println("  Timeout: (default)")
				}
				fmt.Printf("  Command: %s\n", ctx.Command)
				fmt.Println()
			}

			if ctx.Prompt != "" {
				fmt.Println("Prompt template:")
				// Show full prompt
				lines := strings.Split(ctx.Prompt, "\n")
				for _, line := range lines {
					fmt.Printf("  %s\n", line)
				}
				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&localOnly, "local", "l", false, "Show local context only")

	return cmd
}

// NewConfigContextTestCommand creates the config context test command
func NewConfigContextTestCommand(configLoader *config.Loader) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test <name>",
		Short: "Test context configuration and file availability",
		Long:  "Validate context configuration without using it",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			contextName := args[0]
			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			// Load merged config
			globalCfg, err := configLoader.LoadGlobal()
			if err != nil {
				return fmt.Errorf("failed to load global config: %w", err)
			}

			var contexts map[string]domain.Context
			var scope string
			localCfg, err := configLoader.LoadLocal(workDir)
			if err == nil {
				mergedCfg := config.Merge(globalCfg, localCfg)
				contexts = mergedCfg.Contexts

				if _, hasLocal := localCfg.Contexts[contextName]; hasLocal {
					scope = "local"
				} else {
					scope = "global"
				}
			} else {
				contexts = globalCfg.Contexts
				scope = "global"
			}

			ctx, exists := contexts[contextName]
			if !exists {
				fmt.Fprintf(os.Stderr, "Error: Context '%s' not found in configuration.\n\n", contextName)
				fmt.Fprintln(os.Stderr, "Use 'start config context list' to see available contexts.")
				return fmt.Errorf("context not found")
			}

			fmt.Printf("Testing context: %s\n", contextName)
			fmt.Println("─────────────────────────────────────────────────")
			fmt.Println()

			fmt.Println("Configuration:")
			fmt.Printf("  Scope: %s\n", scope)
			if ctx.Description != "" {
				fmt.Printf("  Description: %s\n", ctx.Description)
			}
			if ctx.Required {
				fmt.Println("  Required: yes")
			} else {
				fmt.Println("  Required: no")
			}
			fmt.Printf("  Type: %s\n", getContextSourceType(ctx))
			fmt.Println()

			hasErrors := false
			hasWarnings := false

			// Check file availability
			if ctx.File != "" {
				fmt.Println("File:")
				fmt.Printf("  Path: %s\n", ctx.File)
				resolved := resolvePath(ctx.File, workDir)
				fmt.Printf("  Resolved: %s\n", resolved)

				if fileInfo, err := os.Stat(resolved); err == nil {
					fmt.Printf("  ✓ File exists (%.1f KB)\n", float64(fileInfo.Size())/1024)
				} else {
					fmt.Println("  ✗ File not found")
					if ctx.Required {
						hasErrors = true
					} else {
						hasWarnings = true
					}
				}
				fmt.Println()
			}

			// Check command execution
			if ctx.Command != "" {
				fmt.Println("Command:")
				shell := ctx.Shell
				if shell == "" {
					shell = "sh"
				}
				fmt.Printf("  Shell: %s\n", shell)

				timeout := ctx.CommandTimeout
				if timeout == 0 {
					timeout = 30
				}
				fmt.Printf("  Timeout: %d seconds\n", timeout)
				fmt.Printf("  Command: %s\n", ctx.Command)

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
			if ctx.Prompt != "" {
				fmt.Println("Prompt template:")

				// Check for placeholders
				placeholders := findPlaceholders(ctx.Prompt)
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
					if strings.Contains(ctx.Prompt, "{file}") || strings.Contains(ctx.Prompt, "{file_contents}") {
						if ctx.File == "" {
							fmt.Println("  ⚠ Prompt uses {file} or {file_contents} but no file configured")
							hasWarnings = true
						} else {
							fmt.Println("  ✓ Uses {file} placeholder (matches file field)")
						}
					}
					if strings.Contains(ctx.Prompt, "{command}") || strings.Contains(ctx.Prompt, "{command_output}") {
						if ctx.Command == "" {
							fmt.Println("  ⚠ Prompt uses {command} or {command_output} but no command configured")
							hasWarnings = true
						} else {
							fmt.Println("  ✓ Uses {command} placeholder (matches command field)")
						}
					}
				} else {
					fmt.Printf("  ✓ Valid inline prompt (%d characters)\n", len(ctx.Prompt))
				}
				fmt.Println()
			}

			// Check UTD requirement
			if ctx.File == "" && ctx.Command == "" && ctx.Prompt == "" {
				fmt.Println("✗ No content source defined")
				fmt.Println("  At least one field required: file, command, or prompt")
				hasErrors = true
				fmt.Println()
			}

			// Summary
			if hasErrors {
				fmt.Printf("✗ Context '%s' has errors\n", contextName)
				fmt.Println("  Fix configuration: start config context edit", contextName)
				return fmt.Errorf("configuration errors")
			} else if hasWarnings {
				fmt.Printf("⚠ Context '%s' has warnings\n", contextName)
				fmt.Println("  File will generate warning and be skipped at runtime")
			} else {
				fmt.Printf("✓ Context '%s' is configured correctly\n", contextName)
			}

			return nil
		},
	}

	return cmd
}

// getContextSourceType returns a human-readable source type for a context
func getContextSourceType(ctx domain.Context) string {
	hasFile := ctx.File != ""
	hasCommand := ctx.Command != ""
	hasPrompt := ctx.Prompt != ""

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
