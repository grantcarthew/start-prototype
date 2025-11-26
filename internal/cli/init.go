package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/grantcarthew/start/internal/assets"
	"github.com/grantcarthew/start/internal/domain"
	"github.com/spf13/cobra"
)

// InitCommand handles the init command
type InitCommand struct {
	resolver *assets.Resolver
}

// NewInitCommand creates the 'start init' command
func NewInitCommand(resolver *assets.Resolver) *cobra.Command {
	ic := &InitCommand{
		resolver: resolver,
	}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize start configuration",
		Long:  "Interactive wizard to create start configuration files with auto-detected agents",
		Args:  cobra.NoArgs,
		RunE:  ic.runInit,
	}

	cmd.Flags().BoolP("local", "l", false, "Create local config in ./.start/")
	cmd.Flags().BoolP("force", "f", false, "Skip all prompts, auto-configure detected agents")

	return cmd
}

// runInit executes the init command
func (ic *InitCommand) runInit(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Get flags
	local, _ := cmd.Flags().GetBool("local")
	force, _ := cmd.Flags().GetBool("force")

	fmt.Println("Initialize start configuration")
	fmt.Println()

	// Determine target location
	var targetPath string
	var locationName string

	if local {
		targetPath = "./.start"
		locationName = "local"
		if !force {
			fmt.Printf("Creating local config at %s...\n\n", targetPath)
		}
	} else if force {
		home, _ := os.UserHomeDir()
		targetPath = filepath.Join(home, ".config", "start")
		locationName = "global"
	} else {
		// Interactive location selection
		fmt.Println("Where should this configuration be created?")
		fmt.Println("  1) Global (~/.config/start/)")
		fmt.Println("     Personal config across all projects")
		fmt.Println("  2) Local (./.start/)")
		fmt.Println("     Project config (can be committed to git)")
		fmt.Println()
		fmt.Print("Select [1-2] (default: 1): ")

		var input string
		fmt.Scanln(&input)
		fmt.Println()

		if input == "2" {
			targetPath = "./.start"
			locationName = "local"
		} else {
			home, _ := os.UserHomeDir()
			targetPath = filepath.Join(home, ".config", "start")
			locationName = "global"
		}
	}

	// Check for existing config
	configExists := false
	configPath := filepath.Join(targetPath, "config.toml")
	if _, err := os.Stat(configPath); err == nil {
		configExists = true
	}

	// Handle existing config
	if configExists {
		if force {
			// Auto-backup
			if err := ic.backupConfig(targetPath); err != nil {
				return fmt.Errorf("failed to backup existing config: %w", err)
			}
			if !force {
				fmt.Println("✓ Backed up existing config")
				fmt.Println()
			}
		} else {
			// Interactive backup prompt
			fmt.Printf("Existing config found: %s\n\n", targetPath)
			fmt.Print("Backup and reinitialize? [y/N]: ")

			var confirm string
			fmt.Scanln(&confirm)
			fmt.Println()

			if confirm != "y" && confirm != "Y" {
				fmt.Println("Cancelled. No changes made.")
				return nil
			}

			// Backup
			fmt.Println("Backing up config files...")
			if err := ic.backupConfig(targetPath); err != nil {
				return fmt.Errorf("failed to backup config: %w", err)
			}
			fmt.Println()
		}
	}

	if !force {
		fmt.Println("Welcome to start!")
		fmt.Println()
	}

	// Fetch catalog index
	repo := os.Getenv("ASSET_REPO")
	if repo == "" {
		repo = "grantcarthew/start"
	}

	if !force {
		fmt.Println("Fetching latest agent configurations from GitHub...")
	}

	allAssets, err := ic.resolver.SearchCatalog(ctx, "", repo)
	if err != nil {
		fmt.Println()
		fmt.Println("Error: Failed to fetch agent configurations from GitHub.")
		fmt.Println()
		fmt.Println("Check your network connection and try again.")
		fmt.Println()
		fmt.Println("See https://github.com/grantcarthew/start#configuration for manual setup.")
		return err
	}

	// Filter for agents
	var agentAssets []domain.AssetMeta
	for _, asset := range allAssets {
		if asset.Type == "agents" {
			agentAssets = append(agentAssets, asset)
		}
	}

	if !force {
		fmt.Printf("✓ Found %d agent configurations\n\n", len(agentAssets))
	}

	// Auto-detect installed agents
	if !force {
		fmt.Println("Detecting installed agents...")
	}

	detectedAgents := ic.detectInstalledAgents(agentAssets)

	if !force {
		if len(detectedAgents) > 0 {
			for _, agent := range detectedAgents {
				fmt.Printf("✓ %s (%s)\n", agent.Name, agent.Description)
			}
		} else {
			fmt.Println("✗ No agents detected")
		}
		fmt.Println()
	}

	// Select agents to configure
	var selectedAgents []domain.AssetMeta

	if force {
		// Auto-mode: use all detected agents
		selectedAgents = detectedAgents
		if len(selectedAgents) == 0 {
			fmt.Println("No agents detected.")
			fmt.Println()
			fmt.Println("To configure custom agents, see the documentation:")
			fmt.Println("https://github.com/grantcarthew/start#configuration")
			return nil
		}
	} else {
		// Interactive mode: allow selecting additional agents
		selectedAgents = detectedAgents

		// TODO: Implement interactive selection of additional agents
		// For now, just use detected agents
	}

	// Select default agent
	var defaultAgent string
	if len(selectedAgents) > 0 {
		if force {
			// Auto-select first agent (priority: claude > gemini > others)
			defaultAgent = ic.selectDefaultAgent(selectedAgents)
		} else {
			// Interactive selection
			if len(selectedAgents) == 1 {
				defaultAgent = selectedAgents[0].Name
			} else {
				fmt.Println("Select default agent:")
				for i, agent := range selectedAgents {
					fmt.Printf("  %d) %s\n", i+1, agent.Name)
				}
				fmt.Print("Default [1]: ")

				var input string
				fmt.Scanln(&input)
				fmt.Println()

				if input == "" || input == "1" {
					defaultAgent = selectedAgents[0].Name
				} else {
					// Parse selection
					var sel int
					fmt.Sscanf(input, "%d", &sel)
					if sel > 0 && sel <= len(selectedAgents) {
						defaultAgent = selectedAgents[sel-1].Name
					} else {
						defaultAgent = selectedAgents[0].Name
					}
				}
			}
		}
	}

	// Create config directory
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Generate and write config files
	if !force {
		fmt.Printf("Creating configuration at %s...\n", targetPath)
	} else {
		fmt.Printf("Creating %s config at %s...\n", locationName, targetPath)
	}

	if err := ic.writeConfigFiles(targetPath, selectedAgents, defaultAgent); err != nil {
		return fmt.Errorf("failed to write config files: %w", err)
	}

	if !force {
		fmt.Println("✓ config.toml created")
		fmt.Println("✓ agents.toml created")
		fmt.Println("✓ roles.toml created")
		fmt.Println("✓ contexts.toml created")
		fmt.Println("✓ tasks.toml created")
		fmt.Println()

		fmt.Println("Default context documents configured:")
		fmt.Println("  ~/reference/ENVIRONMENT.md (required)")
		fmt.Println("  ~/reference/INDEX.csv")
		fmt.Println("  ./AGENTS.md")
		fmt.Println("  ./PROJECT.md")
		fmt.Println()

		if locationName == "local" {
			fmt.Println("Local config created. This can be committed to git for team consistency.")
		}

		fmt.Println("Run 'start config show' to see your configuration.")
		fmt.Println("Run 'start' to launch!")
	} else {
		if len(selectedAgents) > 0 {
			agentNames := make([]string, len(selectedAgents))
			for i, a := range selectedAgents {
				agentNames[i] = a.Name
			}
			fmt.Printf("✓ Detected and configured: %s\n", strings.Join(agentNames, ", "))
			fmt.Printf("✓ Default agent: %s\n", defaultAgent)
		}
		fmt.Println("✓ Config created successfully")
	}

	return nil
}

// detectInstalledAgents checks which agents from the catalog are installed
func (ic *InitCommand) detectInstalledAgents(agents []domain.AssetMeta) []domain.AssetMeta {
	var detected []domain.AssetMeta

	for _, agent := range agents {
		if agent.Bin == "" {
			continue
		}

		// Check if binary is in PATH
		_, err := exec.LookPath(agent.Bin)
		if err == nil {
			detected = append(detected, agent)
		}
	}

	return detected
}

// selectDefaultAgent selects the default agent with priority
func (ic *InitCommand) selectDefaultAgent(agents []domain.AssetMeta) string {
	// Priority: claude > gemini > others
	priority := []string{"claude", "gemini"}

	for _, p := range priority {
		for _, agent := range agents {
			if agent.Name == p {
				return agent.Name
			}
		}
	}

	// Return first if no priority match
	if len(agents) > 0 {
		return agents[0].Name
	}

	return ""
}

// backupConfig creates timestamped backups of existing config files
func (ic *InitCommand) backupConfig(targetPath string) error {
	timestamp := time.Now().Format("2006-01-02-150405")

	configFiles := []string{
		"config.toml",
		"agents.toml",
		"roles.toml",
		"contexts.toml",
		"tasks.toml",
	}

	for _, filename := range configFiles {
		sourcePath := filepath.Join(targetPath, filename)
		if _, err := os.Stat(sourcePath); err == nil {
			// File exists, back it up
			backupName := strings.TrimSuffix(filename, ".toml") + "." + timestamp + ".toml"
			backupPath := filepath.Join(targetPath, backupName)

			data, err := os.ReadFile(sourcePath)
			if err != nil {
				return fmt.Errorf("failed to read %s: %w", filename, err)
			}

			if err := os.WriteFile(backupPath, data, 0644); err != nil {
				return fmt.Errorf("failed to write backup %s: %w", backupName, err)
			}

			fmt.Printf("✓ %s\n", backupName)
		}
	}

	return nil
}

// writeConfigFiles generates and writes all config files
func (ic *InitCommand) writeConfigFiles(targetPath string, agents []domain.AssetMeta, defaultAgent string) error {
	// config.toml
	configContent := fmt.Sprintf(`[settings]
default_agent = "%s"
default_role = "code-reviewer"
log_level = "info"
asset_download = true
asset_repo = "grantcarthew/start"
`, defaultAgent)

	if err := os.WriteFile(filepath.Join(targetPath, "config.toml"), []byte(configContent), 0644); err != nil {
		return err
	}

	// agents.toml
	agentsContent := "# Agent configurations\n\n"
	for _, agent := range agents {
		// TODO: Download and use actual agent config from GitHub
		// For now, create a basic template
		agentsContent += fmt.Sprintf(`[agents.%s]
bin = "%s"
description = "%s"

`, agent.Name, agent.Bin, agent.Description)
	}

	if err := os.WriteFile(filepath.Join(targetPath, "agents.toml"), []byte(agentsContent), 0644); err != nil {
		return err
	}

	// roles.toml
	rolesContent := `# Role definitions

[roles.code-reviewer]
description = "Expert code reviewer focusing on quality and best practices"
file = "./ROLE.md"
`

	if err := os.WriteFile(filepath.Join(targetPath, "roles.toml"), []byte(rolesContent), 0644); err != nil {
		return err
	}

	// contexts.toml
	contextsContent := `# Context documents

[contexts.environment]
description = "Environment and system information"
file = "~/reference/ENVIRONMENT.md"
required = true

[contexts.index]
description = "Documentation index"
file = "~/reference/INDEX.csv"

[contexts.agents]
description = "Agent configuration and usage"
file = "./AGENTS.md"

[contexts.project]
description = "Project overview and context"
file = "./PROJECT.md"
`

	if err := os.WriteFile(filepath.Join(targetPath, "contexts.toml"), []byte(contextsContent), 0644); err != nil {
		return err
	}

	// tasks.toml
	tasksContent := `# Task definitions
#
# This file is created by 'start init'.
# Add custom tasks here or install them from the
# asset catalog using 'start assets add'.
#
# Example task:
#
# [tasks.code-review]
# description = "Review code for quality"
# role = "code-reviewer"
# prompt = "Review the following code..."
`

	if err := os.WriteFile(filepath.Join(targetPath, "tasks.toml"), []byte(tasksContent), 0644); err != nil {
		return err
	}

	return nil
}
