package main

import (
	"os"
	"path/filepath"

	"github.com/grantcarthew/start/internal/adapters"
	"github.com/grantcarthew/start/internal/assets"
	"github.com/grantcarthew/start/internal/cli"
	"github.com/grantcarthew/start/internal/config"
	"github.com/grantcarthew/start/internal/engine"
)

var version = "dev" // Injected at build time via -ldflags

func main() {
	// Create adapters (real implementations)
	fs := &adapters.RealFileSystem{}
	runner := &adapters.RealRunner{}
	commandRunner := adapters.NewRealCommandRunner()
	githubClient := adapters.NewRealGitHubClient()

	// Create config loader
	configLoader := config.NewLoader(fs)

	// Create validator
	validator := config.NewValidator()

	// Get working directory
	workDir, err := os.Getwd()
	if err != nil {
		workDir = "."
	}

	// Create cache
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	cacheBase := filepath.Join(home, ".config", "start", "assets")
	cache := adapters.NewFileCache(fs, cacheBase)

	// Create engine components
	placeholderResolver := engine.NewPlaceholderResolver()
	utdProcessor := engine.NewUTDProcessor(fs, commandRunner, workDir)
	roleSelector := engine.NewRoleSelector()
	roleLoader := engine.NewRoleLoader(utdProcessor, fs)
	contextLoader := engine.NewContextLoader(utdProcessor)
	taskLoader := engine.NewTaskLoader(utdProcessor, placeholderResolver)
	taskResolver := engine.NewTaskResolver()
	executor := engine.NewExecutor(runner, placeholderResolver)

	// Create asset resolver
	assetResolver := assets.NewResolver(fs, cache, githubClient, configLoader)

	// Create root command with dependencies
	rootCmd := cli.NewRootCommand(
		configLoader,
		validator,
		executor,
		roleSelector,
		roleLoader,
		contextLoader,
		taskLoader,
		taskResolver,
		assetResolver,
		version,
	)

	// Execute
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
