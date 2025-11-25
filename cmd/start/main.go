package main

import (
	"os"

	"github.com/grantcarthew/start/internal/adapters"
	"github.com/grantcarthew/start/internal/cli"
	"github.com/grantcarthew/start/internal/config"
)

var version = "dev" // Injected at build time via -ldflags

func main() {
	// Create adapters (real implementations)
	fs := &adapters.RealFileSystem{}

	// Create config loader
	configLoader := config.NewLoader(fs)

	// Create validator
	validator := config.NewValidator()

	// Create root command with dependencies
	rootCmd := cli.NewRootCommand(configLoader, validator, version)

	// Execute
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
