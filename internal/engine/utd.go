package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/grantcarthew/start/internal/domain"
)

// UTDProcessor processes Unified Template Design (UTD) pattern
// Handles file/command/prompt fields with placeholders
type UTDProcessor struct {
	fs            domain.FileSystem
	commandRunner domain.CommandRunner
	workDir       string
}

// NewUTDProcessor creates a new UTD processor
func NewUTDProcessor(fs domain.FileSystem, commandRunner domain.CommandRunner, workDir string) *UTDProcessor {
	return &UTDProcessor{
		fs:            fs,
		commandRunner: commandRunner,
		workDir:       workDir,
	}
}

// UTDInput represents the UTD fields from config
type UTDInput struct {
	File           string
	Command        string
	Prompt         string
	Shell          string
	CommandTimeout int
}

// UTDResult represents the processed result
type UTDResult struct {
	Content  string   // Final resolved content
	FilePath string   // Resolved file path (if file field present)
	Warnings []string // Any warnings during processing
	Skipped  bool     // True if section should be skipped
}

// Process resolves a UTD pattern into final content
func (p *UTDProcessor) Process(input UTDInput, defaultShell string, defaultTimeout int) UTDResult {
	result := UTDResult{
		Warnings: []string{},
	}

	// Validate: at least one field must be present
	if input.File == "" && input.Command == "" && input.Prompt == "" {
		result.Warnings = append(result.Warnings, "Empty section: at least one of file, command, or prompt required")
		result.Skipped = true
		return result
	}

	// Determine shell and timeout
	shell := input.Shell
	if shell == "" {
		shell = defaultShell
	}
	timeout := input.CommandTimeout
	if timeout == 0 {
		timeout = defaultTimeout
	}

	// Process based on field combinations
	hasFile := input.File != ""
	hasCommand := input.Command != ""
	hasPrompt := input.Prompt != ""

	// Read file if present
	var fileContents string
	var filePath string
	if hasFile {
		filePath = p.resolvePath(input.File)
		contents, err := p.fs.ReadFile(filePath)
		if err != nil {
			if hasPrompt && (strings.Contains(input.Prompt, "{file}") || strings.Contains(input.Prompt, "{file_contents}")) {
				result.Warnings = append(result.Warnings, fmt.Sprintf("File not found: %s", filePath))
				result.Skipped = true
				return result
			}
			// File missing but not used in prompt - warn but continue
			result.Warnings = append(result.Warnings, fmt.Sprintf("File not found: %s (ignored)", filePath))
		} else {
			fileContents = string(contents)
		}
		result.FilePath = filePath
	}

	// Execute command if present
	var commandOutput string
	if hasCommand {
		output, err := p.commandRunner.Run(shell, input.Command, timeout)
		if err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Command failed: %v", err))
			commandOutput = "" // Use empty output on failure
		} else {
			commandOutput = strings.TrimRight(output, "\n")
		}
	}

	// Determine final content based on field combinations
	switch {
	case hasPrompt:
		// Prompt is the template - inject file and command placeholders
		content := input.Prompt

		// Check if file placeholders are used
		usesFile := strings.Contains(content, "{file}") || strings.Contains(content, "{file_contents}")
		if hasFile && !usesFile {
			result.Warnings = append(result.Warnings, "File defined but not used in prompt")
		}
		if !hasFile && usesFile {
			result.Warnings = append(result.Warnings, "No file defined but prompt uses {file}")
			result.Skipped = true
			return result
		}

		// Check if command placeholders are used
		usesCommand := strings.Contains(content, "{command}") || strings.Contains(content, "{command_output}")
		if hasCommand && !usesCommand {
			result.Warnings = append(result.Warnings, "Command defined but not used in prompt")
		}
		if !hasCommand && usesCommand {
			result.Warnings = append(result.Warnings, "No command defined but prompt uses {command}")
			result.Skipped = true
			return result
		}

		// Replace file placeholders
		content = strings.ReplaceAll(content, "{file}", filePath)
		content = strings.ReplaceAll(content, "{file_contents}", fileContents)

		// Replace command placeholders
		content = strings.ReplaceAll(content, "{command}", input.Command)
		content = strings.ReplaceAll(content, "{command_output}", commandOutput)

		result.Content = content

	case hasFile && hasCommand:
		// File + command (no prompt) - check if file contains command placeholders
		if strings.Contains(fileContents, "{command}") || strings.Contains(fileContents, "{command_output}") {
			content := fileContents
			content = strings.ReplaceAll(content, "{command}", input.Command)
			content = strings.ReplaceAll(content, "{command_output}", commandOutput)
			result.Content = content
		} else {
			result.Warnings = append(result.Warnings, "Command defined but not used in file")
			result.Content = fileContents
		}

	case hasFile:
		// Only file - use contents directly
		result.Content = fileContents

	case hasCommand:
		// Only command - use output directly
		result.Content = commandOutput
	}

	return result
}

// resolvePath resolves a file path (expanding ~ and making absolute)
func (p *UTDProcessor) resolvePath(path string) string {
	// Expand tilde
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[2:])
		}
	}

	// Make absolute if relative
	if !filepath.IsAbs(path) {
		path = filepath.Join(p.workDir, path)
	}

	return path
}
