package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	outputDir := os.Getenv("SMITH_OUTPUT_DIR")

	if outputDir == "" {
		// Manual mode: print prompt to stdout
		if len(os.Args) > 1 {
			fmt.Println(os.Args[len(os.Args)-1])
		}
		os.Exit(0)
	}

	// Testing mode: write args and prompt to files
	argsFile := filepath.Join(outputDir, "args.txt")
	promptFile := filepath.Join(outputDir, "prompt.md")

	// Write args (one per line)
	argsContent := strings.Join(os.Args, "\n")
	if err := os.WriteFile(argsFile, []byte(argsContent), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing args: %v\n", err)
		os.Exit(1)
	}

	// Write prompt (last arg)
	if len(os.Args) > 1 {
		prompt := os.Args[len(os.Args)-1]
		if err := os.WriteFile(promptFile, []byte(prompt), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing prompt: %v\n", err)
			os.Exit(1)
		}
	}

	os.Exit(0)
}
