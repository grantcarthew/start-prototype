package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/grantcarthew/start/internal/config"
	"github.com/grantcarthew/start/internal/version"
	"github.com/spf13/cobra"
)

// DoctorCommand handles health checks
type DoctorCommand struct {
	configLoader *config.Loader
	validator    *config.Validator
	version      string
	quiet        bool
	verbose      bool
}

// NewDoctorCommand creates the doctor command
func NewDoctorCommand(
	configLoader *config.Loader,
	validator *config.Validator,
	versionString string,
) *cobra.Command {
	dc := &DoctorCommand{
		configLoader: configLoader,
		validator:    validator,
		version:      versionString,
	}

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Diagnose start installation and configuration",
		Long: `Performs comprehensive health check of start installation, configuration, and environment.

Health checks performed:
- Version check (current vs latest release)
- Asset library (age and availability)
- Configuration validation
- Agent diagnostics
- Context verification
- Environment checks

Exit codes:
  0 - All checks passed (no issues)
  1 - Issues found (errors or warnings)`,
		RunE: dc.run,
	}

	cmd.Flags().BoolVarP(&dc.quiet, "quiet", "q", false, "Quiet mode (only show issues)")
	cmd.Flags().BoolVarP(&dc.verbose, "verbose", "v", false, "Verbose output")

	return cmd
}

// run executes the doctor command
func (dc *DoctorCommand) run(cmd *cobra.Command, args []string) error {
	var errors []string
	var warnings []string

	if !dc.quiet {
		fmt.Println("Diagnosing start installation...")
		fmt.Println("═══════════════════════════════════════════════════════════")
		fmt.Println()
	}

	// 1. Version check
	if !dc.quiet {
		fmt.Println("Version")
	}
	versionWarnings := dc.checkVersion()
	warnings = append(warnings, versionWarnings...)

	// 2. Asset library check
	if !dc.quiet {
		fmt.Println()
		fmt.Println("Assets")
	}
	assetWarnings := dc.checkAssets()
	warnings = append(warnings, assetWarnings...)

	// 3. Configuration validation
	if !dc.quiet {
		fmt.Println()
		fmt.Println("Configuration")
	}
	configErrors, configWarnings := dc.checkConfiguration()
	errors = append(errors, configErrors...)
	warnings = append(warnings, configWarnings...)

	// 4. Agent diagnostics
	if !dc.quiet {
		fmt.Println()
		fmt.Println("Agents")
	}
	agentErrors := dc.checkAgents()
	errors = append(errors, agentErrors...)

	// 5. Context verification
	if !dc.quiet {
		fmt.Println()
		fmt.Println("Contexts")
	}
	contextWarnings := dc.checkContexts()
	warnings = append(warnings, contextWarnings...)

	// 6. Environment check
	if !dc.quiet {
		fmt.Println()
		fmt.Println("Environment")
	}
	envErrors := dc.checkEnvironment()
	errors = append(errors, envErrors...)

	// Summary
	if !dc.quiet {
		fmt.Println()
		fmt.Println("Summary")
		fmt.Println("───────────────────────────────────────────────────────────")
	}

	hasIssues := len(errors) > 0 || len(warnings) > 0

	if !hasIssues {
		if !dc.quiet {
			fmt.Println("  ✓ No issues found")
			fmt.Println()
			fmt.Println("Everything looks good!")
		}
		return nil
	}

	// Show issues
	fmt.Printf("  %d errors, %d warnings found\n\n", len(errors), len(warnings))

	if len(errors) > 0 {
		fmt.Println("Errors:")
		for _, err := range errors {
			fmt.Printf("  ✗ %s\n", err)
		}
		fmt.Println()
	}

	if len(warnings) > 0 {
		fmt.Println("Warnings:")
		for _, warn := range warnings {
			fmt.Printf("  ⚠ %s\n", warn)
		}
		fmt.Println()
	}

	// Exit with error code 1 if any issues found
	os.Exit(1)
	return nil
}

// checkVersion checks the CLI version
func (dc *DoctorCommand) checkVersion() []string {
	var warnings []string

	currentVersion := dc.version
	if currentVersion == "" || currentVersion == "dev" {
		if !dc.quiet {
			fmt.Printf("  start %s (development build)\n", currentVersion)
		}
		return warnings
	}

	// Check latest release
	checker := version.NewChecker("grantcarthew/start")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	latest, err := checker.CheckLatestRelease(ctx)
	if err != nil {
		if !dc.quiet {
			fmt.Printf("  start v%s\n", currentVersion)
			fmt.Printf("  ⚠ Could not check for updates: %v\n", err)
		}
		return warnings
	}

	status, message := version.CompareVersions(currentVersion, latest.TagName)

	if !dc.quiet {
		fmt.Printf("  start v%s\n", currentVersion)
		if status == "Up to date" {
			fmt.Printf("  ✓ %s\n", message)
		} else if status == "Update available" {
			fmt.Printf("  ⚠ %s\n", message)
			updateCmd := version.DetectInstallMethod()
			fmt.Printf("  Update: %s\n", updateCmd)
		} else {
			fmt.Printf("  ℹ %s\n", message)
		}
	}

	if status == "Update available" {
		warnings = append(warnings, fmt.Sprintf("CLI update available (%s)", message))
	}

	return warnings
}

// checkAssets checks the asset library
func (dc *DoctorCommand) checkAssets() []string {
	var warnings []string

	home, err := os.UserHomeDir()
	if err != nil {
		if !dc.quiet {
			fmt.Println("  ✗ Could not determine home directory")
		}
		return warnings
	}

	assetDir := filepath.Join(home, ".config", "start", "assets")
	info, err := os.Stat(assetDir)
	if err != nil {
		if !dc.quiet {
			fmt.Println("  ✗ Asset library not initialized")
			fmt.Println("  Action: Assets will download on-demand or run 'start init'")
		}
		warnings = append(warnings, "Asset library not initialized")
		return warnings
	}

	// Check age
	age := time.Since(info.ModTime())
	days := int(age.Hours() / 24)

	if !dc.quiet {
		if days <= 30 {
			fmt.Printf("  ✓ Asset library up to date (updated %d days ago)\n", days)
		} else if days <= 90 {
			fmt.Printf("  ⚠ Assets are %d days old\n", days)
			fmt.Println("  Run 'start assets update' to refresh")
		} else {
			fmt.Printf("  ⚠ Assets are very old (%d days)\n", days)
			fmt.Println("  Run 'start assets update' to refresh")
		}
	}

	if days > 30 {
		warnings = append(warnings, fmt.Sprintf("Assets outdated (%d days old)", days))
	}

	return warnings
}

// checkConfiguration validates configuration files
func (dc *DoctorCommand) checkConfiguration() ([]string, []string) {
	var errors []string
	var warnings []string

	// Try to load configuration
	globalCfg, err := dc.configLoader.LoadGlobal()
	if err != nil {
		if !dc.quiet {
			fmt.Printf("  ✗ Failed to load global config: %v\n", err)
		}
		errors = append(errors, fmt.Sprintf("Global config error: %v", err))
		return errors, warnings
	}

	localCfg, err := dc.configLoader.LoadLocal(".")
	if err != nil {
		// Local config is optional
		localCfg = globalCfg
	}

	// Merge configs
	cfg := config.Merge(globalCfg, localCfg)

	// Validate
	if err := dc.validator.Validate(cfg); err != nil {
		if !dc.quiet {
			fmt.Printf("  ✗ Configuration validation failed: %v\n", err)
		}
		errors = append(errors, fmt.Sprintf("Configuration validation: %v", err))
		return errors, warnings
	}

	if !dc.quiet {
		fmt.Println("  ✓ Configuration valid")
	}

	return errors, warnings
}

// checkAgents checks agent binary availability
func (dc *DoctorCommand) checkAgents() []string {
	var errors []string

	// Load configuration
	globalCfg, err := dc.configLoader.LoadGlobal()
	if err != nil {
		return errors
	}

	localCfg, err := dc.configLoader.LoadLocal(".")
	if err != nil {
		localCfg = globalCfg
	}

	cfg := config.Merge(globalCfg, localCfg)

	if len(cfg.Agents) == 0 {
		if !dc.quiet {
			fmt.Println("  ⚠ No agents configured")
		}
		return errors
	}

	// Check each agent
	for name, agent := range cfg.Agents {
		bin := agent.Bin
		if bin == "" {
			bin = name
		}

		path, err := exec.LookPath(bin)
		if err != nil {
			if !dc.quiet {
				fmt.Printf("  ✗ %s - Binary not found (%s)\n", name, bin)
			}
			errors = append(errors, fmt.Sprintf("Agent '%s' binary not found", name))
		} else {
			if !dc.quiet {
				if dc.verbose {
					fmt.Printf("  ✓ %s - %s\n", name, path)
				} else {
					fmt.Printf("  ✓ %s\n", name)
				}
			}
		}
	}

	return errors
}

// checkContexts verifies context files
func (dc *DoctorCommand) checkContexts() []string {
	var warnings []string

	// Load configuration
	globalCfg, err := dc.configLoader.LoadGlobal()
	if err != nil {
		return warnings
	}

	localCfg, err := dc.configLoader.LoadLocal(".")
	if err != nil {
		localCfg = globalCfg
	}

	cfg := config.Merge(globalCfg, localCfg)

	if len(cfg.Contexts) == 0 {
		if !dc.quiet {
			fmt.Println("  ℹ No contexts configured")
		}
		return warnings
	}

	// Check each context
	requiredCount := 0
	optionalCount := 0

	for name, ctx := range cfg.Contexts {
		if ctx.File == "" {
			continue
		}

		// Expand home directory
		path := ctx.File
		if strings.HasPrefix(path, "~/") {
			home, _ := os.UserHomeDir()
			path = filepath.Join(home, path[2:])
		}

		_, err := os.Stat(path)
		exists := err == nil

		if ctx.Required {
			requiredCount++
			if !exists {
				if !dc.quiet {
					fmt.Printf("  ✗ %s (required) - File not found: %s\n", name, ctx.File)
				}
				warnings = append(warnings, fmt.Sprintf("Required context '%s' file not found", name))
			} else {
				if !dc.quiet && dc.verbose {
					fmt.Printf("  ✓ %s (required) - %s\n", name, ctx.File)
				}
			}
		} else {
			optionalCount++
			if !exists {
				if !dc.quiet && dc.verbose {
					fmt.Printf("  ⚠ %s (optional) - File not found: %s\n", name, ctx.File)
				}
			} else {
				if !dc.quiet && dc.verbose {
					fmt.Printf("  ✓ %s (optional) - %s\n", name, ctx.File)
				}
			}
		}
	}

	if !dc.quiet && !dc.verbose {
		fmt.Printf("  %d required, %d optional contexts configured\n", requiredCount, optionalCount)
	}

	return warnings
}

// checkEnvironment verifies environment
func (dc *DoctorCommand) checkEnvironment() []string {
	var errors []string

	// Check config directory
	home, err := os.UserHomeDir()
	if err != nil {
		if !dc.quiet {
			fmt.Println("  ✗ Could not determine home directory")
		}
		errors = append(errors, "Could not determine home directory")
		return errors
	}

	configDir := filepath.Join(home, ".config", "start")
	if info, err := os.Stat(configDir); err != nil {
		if !dc.quiet {
			fmt.Printf("  ⚠ Config directory not found: %s\n", configDir)
		}
	} else if !info.IsDir() {
		if !dc.quiet {
			fmt.Printf("  ✗ Config path is not a directory: %s\n", configDir)
		}
		errors = append(errors, "Config path is not a directory")
	} else {
		if !dc.quiet {
			fmt.Printf("  ✓ Config directory: %s\n", configDir)
		}
	}

	// Check working directory
	workDir, err := os.Getwd()
	if err != nil {
		if !dc.quiet {
			fmt.Println("  ✗ Could not determine working directory")
		}
		errors = append(errors, "Could not determine working directory")
	} else {
		if !dc.quiet && dc.verbose {
			fmt.Printf("  ✓ Working directory: %s\n", workDir)
		}
	}

	// Check shell
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "sh"
	}
	if !dc.quiet && dc.verbose {
		fmt.Printf("  ✓ Shell: %s\n", shell)
	}

	return errors
}
