package main

import (
	"os"

	"github.com/grantcarthew/start/internal/adapters"
	"github.com/grantcarthew/start/internal/cli"
	"github.com/grantcarthew/start/internal/config"
	"github.com/grantcarthew/start/internal/engine"
)

var version = "dev" // Injected at build time via -ldflags

func main() {
	// Create adapters (real implementations)
	fs := &adapters.RealFileSystem{}
	runner := &adapters.RealRunner{}

	// Create config loader
	configLoader := config.NewLoader(fs)

	// Create validator
	validator := config.NewValidator()

	// Create engine components
	placeholderResolver := engine.NewPlaceholderResolver()
	executor := engine.NewExecutor(runner, placeholderResolver)

	// Create root command with dependencies
	rootCmd := cli.NewRootCommand(configLoader, validator, executor, version)

	// Execute
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
