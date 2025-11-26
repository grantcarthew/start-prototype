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

// NewConfigAgentCommand creates the config agent command
func NewConfigAgentCommand(configLoader *config.Loader, validator *config.Validator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Manage AI agent configurations",
		Long:  "Commands for managing AI agent configurations in global and local config files",
	}

	// Add subcommands
	cmd.AddCommand(NewConfigAgentListCommand(configLoader))
	cmd.AddCommand(NewConfigAgentShowCommand(configLoader))
	cmd.AddCommand(NewConfigAgentTestCommand(configLoader))
	cmd.AddCommand(NewConfigAgentNewCommand(configLoader))
	cmd.AddCommand(NewConfigAgentEditCommand(configLoader))
	cmd.AddCommand(NewConfigAgentRemoveCommand(configLoader))
	cmd.AddCommand(NewConfigAgentDefaultCommand(configLoader))

	return cmd
}

// NewConfigAgentListCommand creates the config agent list command
func NewConfigAgentListCommand(configLoader *config.Loader) *cobra.Command {
	var localOnly bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Display all configured agents",
		Long:  "List all agents defined in global and/or local configuration files",
		RunE: func(cmd *cobra.Command, args []string) error {
			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			var agents map[string]domain.Agent

			if localOnly {
				// Load local only
				localCfg, err := configLoader.LoadLocal(workDir)
				if err != nil {
					return fmt.Errorf("failed to load local config: %w", err)
				}
				agents = localCfg.Agents
			} else {
				// Load and merge global + local
				globalCfg, err := configLoader.LoadGlobal()
				if err != nil {
					return fmt.Errorf("failed to load global config: %w", err)
				}

				localCfg, err := configLoader.LoadLocal(workDir)
				if err == nil {
					mergedCfg := config.Merge(globalCfg, localCfg)
					agents = mergedCfg.Agents
				} else {
					agents = globalCfg.Agents
				}
			}

			if len(agents) == 0 {
				fmt.Println("No agents configured.")
				fmt.Println()
				fmt.Println("Run 'start init' to set up agents, or")
				fmt.Println("use 'start assets add' to install from catalog or 'start config agent new' to create custom.")
				return nil
			}

			// Sort agents by name for consistent output
			names := make([]string, 0, len(agents))
			for name := range agents {
				names = append(names, name)
			}
			sort.Strings(names)

			fmt.Println("Configured agents:")
			fmt.Println()
			for _, name := range names {
				agent := agents[name]
				fmt.Printf("%s\n", name)
				if agent.Description != "" {
					fmt.Printf("  %s\n", agent.Description)
				}
				if agent.URL != "" {
					fmt.Printf("  %s\n", agent.URL)
				}
				fmt.Printf("  Command: %s\n", agent.Command)

				// Show default model
				if agent.DefaultModel != "" {
					if fullModel, ok := agent.Models[agent.DefaultModel]; ok {
						fmt.Printf("  Default model: %s (%s)\n", fullModel, agent.DefaultModel)
					} else {
						fmt.Printf("  Default model: %s\n", agent.DefaultModel)
					}
				} else if len(agent.Models) > 0 {
					// Show first model as default
					var firstModel string
					for name := range agent.Models {
						if firstModel == "" || name < firstModel {
							firstModel = name
						}
					}
					fmt.Printf("  Default model: %s (%s) [first model]\n", agent.Models[firstModel], firstModel)
				}

				// Show all models
				if len(agent.Models) > 0 {
					fmt.Println("  Models:")
					modelNames := make([]string, 0, len(agent.Models))
					for name := range agent.Models {
						modelNames = append(modelNames, name)
					}
					sort.Strings(modelNames)
					for _, name := range modelNames {
						fmt.Printf("    - %s (%s)\n", agent.Models[name], name)
					}
				}

				if agent.ModelsURL != "" {
					fmt.Printf("  Model docs: %s\n", agent.ModelsURL)
				}
				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&localOnly, "local", "l", false, "List local agents only")

	return cmd
}

// NewConfigAgentShowCommand creates the config agent show command
func NewConfigAgentShowCommand(configLoader *config.Loader) *cobra.Command {
	var localOnly bool

	cmd := &cobra.Command{
		Use:   "show [name]",
		Short: "Display agent configuration",
		Long:  "Show detailed configuration for a specific agent",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("agent name required\n\nUsage: start config agent show <name>")
			}

			agentName := args[0]
			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			var agents map[string]domain.Agent
			var scope string

			if localOnly {
				localCfg, err := configLoader.LoadLocal(workDir)
				if err != nil {
					return fmt.Errorf("failed to load local config: %w", err)
				}
				agents = localCfg.Agents
				scope = "local"
			} else {
				globalCfg, err := configLoader.LoadGlobal()
				if err != nil {
					return fmt.Errorf("failed to load global config: %w", err)
				}

				localCfg, err := configLoader.LoadLocal(workDir)
				if err == nil {
					mergedCfg := config.Merge(globalCfg, localCfg)
					agents = mergedCfg.Agents

					// Determine scope
					if _, hasLocal := localCfg.Agents[agentName]; hasLocal {
						scope = "local"
					} else {
						scope = "global"
					}
				} else {
					agents = globalCfg.Agents
					scope = "global"
				}
			}

			agent, exists := agents[agentName]
			if !exists {
				return fmt.Errorf("agent '%s' not found in configuration.\n\nUse 'start config agent list' to see available agents.", agentName)
			}

			// Display agent configuration
			fmt.Printf("Agent configuration: %s (%s)\n", agentName, scope)
			fmt.Println("═══════════════════════════════════════════════════════════")
			fmt.Println()

			if agent.Description != "" {
				fmt.Printf("Description: %s\n", agent.Description)
			}
			if agent.URL != "" {
				fmt.Printf("URL: %s\n", agent.URL)
			}
			if agent.ModelsURL != "" {
				fmt.Printf("Models URL: %s\n", agent.ModelsURL)
			}
			if agent.Bin != "" {
				fmt.Printf("Binary: %s\n", agent.Bin)
			}
			fmt.Println()

			fmt.Println("Command template:")
			fmt.Printf("  %s\n", agent.Command)
			fmt.Println()

			if agent.DefaultModel != "" {
				fmt.Printf("Default model: %s\n", agent.DefaultModel)
			} else if len(agent.Models) > 0 {
				fmt.Println("Default model: (first model in config)")
			}

			if len(agent.Models) > 0 {
				fmt.Println("Models:")
				modelNames := make([]string, 0, len(agent.Models))
				for name := range agent.Models {
					modelNames = append(modelNames, name)
				}
				sort.Strings(modelNames)
				for _, name := range modelNames {
					fmt.Printf("  %s = %s\n", name, agent.Models[name])
				}
			} else {
				fmt.Println("Models: (none)")
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&localOnly, "local", "l", false, "Show local agent only")

	return cmd
}

// NewConfigAgentTestCommand creates the config agent test command
func NewConfigAgentTestCommand(configLoader *config.Loader) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test <name>",
		Short: "Test agent configuration and availability",
		Long:  "Validate agent configuration without executing it",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentName := args[0]
			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			// Load merged config
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

			agent, exists := agents[agentName]
			if !exists {
				fmt.Fprintf(os.Stderr, "Error: Agent '%s' not found in configuration.\n\n", agentName)
				fmt.Fprintln(os.Stderr, "Use 'start config agent list' to see available agents.")
				fmt.Fprintln(os.Stderr, "Use 'start assets add' to install from catalog or 'start config agent new' to create custom.")
				return fmt.Errorf("agent not found")
			}

			fmt.Printf("Testing agent: %s\n", agentName)
			fmt.Println("─────────────────────────────────────────────────")
			fmt.Println()

			hasErrors := false
			hasWarnings := false

			// Check binary availability
			if agent.Bin != "" {
				binPath, err := exec.LookPath(agent.Bin)
				if err != nil {
					fmt.Printf("✗ Binary not found: %s\n", agent.Bin)
					fmt.Printf("  The '%s' command is not available.\n", agent.Bin)
					fmt.Printf("  Install %s or check that it's accessible.\n", agent.Bin)
					hasErrors = true
				} else {
					fmt.Printf("✓ Binary found: %s\n", binPath)
				}
				fmt.Println()
			}

			// Check configuration
			fmt.Println("Configuration:")

			// Check command template
			if agent.Command == "" {
				fmt.Println("  ✗ No command template defined")
				hasErrors = true
			} else {
				fmt.Println("  ✓ Command template valid")

				// Check for {prompt} placeholder
				if !strings.Contains(agent.Command, "{prompt}") {
					fmt.Println("  ⚠ Command template missing {prompt} placeholder")
					hasWarnings = true
				} else {
					fmt.Println("  ✓ Contains {prompt} placeholder")
				}

				// Check for unknown placeholders
				knownPlaceholders := []string{"{bin}", "{model}", "{role}", "{role_file}", "{prompt}", "{date}"}
				for _, ph := range []string{"{bin}", "{model}", "{role}", "{role_file}", "{prompt}", "{date}"} {
					if strings.Contains(agent.Command, ph) {
						// Valid placeholder
						continue
					}
				}

				// Look for any other placeholders
				if strings.Contains(agent.Command, "{") {
					parts := strings.Split(agent.Command, "{")
					for _, part := range parts[1:] {
						if idx := strings.Index(part, "}"); idx > 0 {
							ph := "{" + part[:idx] + "}"
							isKnown := false
							for _, known := range knownPlaceholders {
								if ph == known {
									isKnown = true
									break
								}
							}
							if !isKnown {
								fmt.Printf("  ⚠ Unknown placeholder %s in command template\n", ph)
								fmt.Println("    (did you mean one of: {bin}, {model}, {role}, {role_file}, {prompt}, {date}?)")
								hasWarnings = true
							}
						}
					}
				}
			}

			// Check models
			if len(agent.Models) > 0 {
				fmt.Printf("  ✓ Models configured: %d ", len(agent.Models))
				modelNames := make([]string, 0, len(agent.Models))
				for name := range agent.Models {
					modelNames = append(modelNames, name)
				}
				sort.Strings(modelNames)
				fmt.Printf("(%s)\n", strings.Join(modelNames, ", "))
			}

			// Check default model
			if agent.DefaultModel != "" {
				if _, ok := agent.Models[agent.DefaultModel]; ok {
					fmt.Printf("  ✓ Default model: %s (%s)\n", agent.Models[agent.DefaultModel], agent.DefaultModel)
				} else {
					fmt.Printf("  ✓ Default model: %s\n", agent.DefaultModel)
				}
			} else if len(agent.Models) > 0 {
				fmt.Println("  ℹ Default model: (uses first model in config)")
			}

			fmt.Println()

			// Preview command
			fmt.Println("Preview command:")
			previewCmd := agent.Command
			previewCmd = strings.ReplaceAll(previewCmd, "{bin}", agent.Bin)
			if agent.DefaultModel != "" {
				if fullModel, ok := agent.Models[agent.DefaultModel]; ok {
					previewCmd = strings.ReplaceAll(previewCmd, "{model}", fullModel)
				} else {
					previewCmd = strings.ReplaceAll(previewCmd, "{model}", agent.DefaultModel)
				}
			}
			previewCmd = strings.ReplaceAll(previewCmd, "{role}", "...")
			previewCmd = strings.ReplaceAll(previewCmd, "{role_file}", "/tmp/role.txt")
			previewCmd = strings.ReplaceAll(previewCmd, "{prompt}", "test")
			previewCmd = strings.ReplaceAll(previewCmd, "{date}", "2025-01-01T00:00:00Z")

			fmt.Printf("  ❯ %s\n", previewCmd)
			fmt.Println()

			// Summary
			if hasErrors {
				fmt.Printf("✗ Agent '%s' has configuration errors\n", agentName)
				fmt.Println("  Fix configuration: start config agent edit", agentName)
				return fmt.Errorf("configuration errors")
			} else if hasWarnings {
				fmt.Printf("⚠ Agent '%s' has warnings (see above)\n", agentName)
			} else {
				fmt.Printf("✓ Agent '%s' is configured correctly\n", agentName)
			}

			return nil
		},
	}

	return cmd
}
