package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// NewCompletionCommand creates the completion command
func NewCompletionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish]",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for bash, zsh, or fish.

Examples:
  # Output completion script for bash
  start completion bash

  # Install completion for zsh
  start completion install zsh

  # Install completion for bash to custom path
  start completion install bash --path ~/.bash_completion/start`,
		ValidArgs: []string{"bash", "zsh", "fish"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE:      runCompletion,
	}

	// Add install subcommand
	cmd.AddCommand(NewCompletionInstallCommand())

	return cmd
}

// runCompletion generates completion scripts
func runCompletion(cmd *cobra.Command, args []string) error {
	shell := args[0]
	rootCmd := cmd.Root()

	switch shell {
	case "bash":
		return rootCmd.GenBashCompletion(os.Stdout)
	case "zsh":
		return rootCmd.GenZshCompletion(os.Stdout)
	case "fish":
		return rootCmd.GenFishCompletion(os.Stdout, true)
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}
}

// NewCompletionInstallCommand creates the install subcommand
func NewCompletionInstallCommand() *cobra.Command {
	var pathFlag string
	var systemFlag bool

	cmd := &cobra.Command{
		Use:   "install [bash|zsh|fish]",
		Short: "Install completion scripts to standard locations",
		Long: `Auto-install completion scripts to standard shell locations.

Examples:
  # Install for zsh (user directory)
  start completion install zsh

  # Install for bash system-wide (requires sudo)
  start completion install bash --system

  # Install to custom path
  start completion install fish --path ~/.config/fish/completions/start.fish`,
		ValidArgs: []string{"bash", "zsh", "fish"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCompletionInstall(cmd, args[0], pathFlag, systemFlag)
		},
	}

	cmd.Flags().StringVar(&pathFlag, "path", "", "Custom installation path")
	cmd.Flags().BoolVar(&systemFlag, "system", false, "Install system-wide (requires sudo)")

	return cmd
}

// runCompletionInstall installs completion to standard location
func runCompletionInstall(cmd *cobra.Command, shell string, customPath string, system bool) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	var targetPath string
	var reloadInstructions string

	// Determine target path
	if customPath != "" {
		targetPath = customPath
	} else {
		targetPath = getStandardCompletionPath(shell, home, system)
	}

	// Create directory if needed
	dir := filepath.Dir(targetPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Create completion file
	file, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", targetPath, err)
	}
	defer file.Close()

	// Generate completion
	rootCmd := cmd.Root()
	switch shell {
	case "bash":
		if err := rootCmd.GenBashCompletion(file); err != nil {
			return fmt.Errorf("failed to generate bash completion: %w", err)
		}
		reloadInstructions = "source ~/.bashrc"
	case "zsh":
		if err := rootCmd.GenZshCompletion(file); err != nil {
			return fmt.Errorf("failed to generate zsh completion: %w", err)
		}
		reloadInstructions = "source ~/.zshrc"
	case "fish":
		if err := rootCmd.GenFishCompletion(file, true); err != nil {
			return fmt.Errorf("failed to generate fish completion: %w", err)
		}
		reloadInstructions = "source ~/.config/fish/config.fish"
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}

	// Success message
	fmt.Printf("✓ Completion installed to: %s\n\n", targetPath)
	fmt.Printf("Reload your shell:\n  %s\n\n", reloadInstructions)
	fmt.Println("Or start a new terminal session.")

	return nil
}

// getStandardCompletionPath returns the standard completion path for a shell
func getStandardCompletionPath(shell, home string, system bool) string {
	if system {
		switch shell {
		case "bash":
			return "/etc/bash_completion.d/start"
		case "zsh":
			return "/usr/local/share/zsh/site-functions/_start"
		case "fish":
			return "/usr/share/fish/vendor_completions.d/start.fish"
		}
	}

	// User-level paths
	switch shell {
	case "bash":
		return filepath.Join(home, ".bash_completion")
	case "zsh":
		return filepath.Join(home, ".zsh", "completion", "_start")
	case "fish":
		return filepath.Join(home, ".config", "fish", "completions", "start.fish")
	}

	return ""
}
