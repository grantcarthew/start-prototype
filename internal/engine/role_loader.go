package engine

import (
	"fmt"

	"github.com/grantcarthew/start/internal/domain"
)

// RoleLoader loads and processes roles
type RoleLoader struct {
	utdProcessor *UTDProcessor
	fs           domain.FileSystem
}

// NewRoleLoader creates a new role loader
func NewRoleLoader(utdProcessor *UTDProcessor, fs domain.FileSystem) *RoleLoader {
	return &RoleLoader{
		utdProcessor: utdProcessor,
		fs:           fs,
	}
}

// LoadedRole represents a processed role
type LoadedRole struct {
	Name     string
	Content  string   // Resolved role content (for {role} placeholder)
	FilePath string   // Path for {role_file} placeholder (original or temp file)
	IsTemp   bool     // True if FilePath points to a temporary file
	Warnings []string // Warnings during processing
}

// LoadRole loads and processes a role through UTD
func (l *RoleLoader) LoadRole(
	role domain.Role,
	defaultShell string,
	defaultTimeout int,
) (LoadedRole, error) {
	result := LoadedRole{
		Name:     role.Name,
		Warnings: []string{},
	}

	// Process through UTD
	utdInput := UTDInput{
		File:           role.File,
		Command:        role.Command,
		Prompt:         role.Prompt,
		Shell:          role.Shell,
		CommandTimeout: role.CommandTimeout,
	}

	utdResult := l.utdProcessor.Process(utdInput, defaultShell, defaultTimeout)

	// Check if processing was skipped
	if utdResult.Skipped {
		return result, fmt.Errorf("role processing failed: %v", utdResult.Warnings)
	}

	result.Content = utdResult.Content
	result.Warnings = utdResult.Warnings

	// Determine file path for {role_file} placeholder
	// Simple role (file only) -> use original file path
	// Complex role (UTD) -> create temp file
	if role.File != "" && role.Command == "" && role.Prompt == "" {
		// Simple role - use original file path
		result.FilePath = utdResult.FilePath
		result.IsTemp = false
	} else {
		// Complex role or non-file role - create temp file
		tempPath, err := l.fs.TempFile("start-role-*.md")
		if err != nil {
			return result, fmt.Errorf("failed to create temp file for role: %w", err)
		}

		// Write resolved content to temp file
		if err := l.fs.WriteFile(tempPath, []byte(result.Content), 0600); err != nil {
			l.fs.Remove(tempPath) // Clean up on error
			return result, fmt.Errorf("failed to write temp file for role: %w", err)
		}

		result.FilePath = tempPath
		result.IsTemp = true
	}

	return result, nil
}

// CleanupRole removes temporary files if needed
func (l *RoleLoader) CleanupRole(role LoadedRole) error {
	if role.IsTemp && role.FilePath != "" {
		return l.fs.Remove(role.FilePath)
	}
	return nil
}
