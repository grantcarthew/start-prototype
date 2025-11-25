package engine

import (
	"fmt"
	"strings"

	"github.com/grantcarthew/start/internal/domain"
)

// TaskLoader loads and processes tasks
type TaskLoader struct {
	utdProcessor *UTDProcessor
	resolver     *PlaceholderResolver
}

// NewTaskLoader creates a new task loader
func NewTaskLoader(utdProcessor *UTDProcessor, resolver *PlaceholderResolver) *TaskLoader {
	return &TaskLoader{
		utdProcessor: utdProcessor,
		resolver:     resolver,
	}
}

// LoadedTask represents a processed task
type LoadedTask struct {
	Name         string
	Prompt       string   // Final task prompt with placeholders resolved
	CommandExec  string   // Command that was executed (for display)
	Warnings     []string // Warnings during processing
}

// LoadTask loads and processes a task through UTD with instructions
func (l *TaskLoader) LoadTask(
	task domain.Task,
	instructions string,
	defaultShell string,
	defaultTimeout int,
) (LoadedTask, error) {
	result := LoadedTask{
		Name:     task.Name,
		Warnings: []string{},
	}

	// Process through UTD to get file contents and command output
	utdInput := UTDInput{
		File:           task.File,
		Command:        task.Command,
		Prompt:         task.Prompt,
		Shell:          task.Shell,
		CommandTimeout: task.CommandTimeout,
	}

	utdResult := l.utdProcessor.Process(utdInput, defaultShell, defaultTimeout)

	// Check if processing was skipped
	if utdResult.Skipped {
		return result, fmt.Errorf("task processing failed: %v", utdResult.Warnings)
	}

	result.Warnings = utdResult.Warnings
	result.CommandExec = task.Command

	// Now resolve task-specific placeholders
	// The UTD processor already handled {file}, {file_contents}, {command}, {command_output}
	// We need to handle {instructions} placeholder
	prompt := utdResult.Content

	// Handle {instructions} placeholder
	// Default to "None" if instructions are empty (per DR-009)
	instructionsValue := instructions
	if instructionsValue == "" {
		instructionsValue = "None"
	}
	prompt = strings.ReplaceAll(prompt, "{instructions}", instructionsValue)

	// Handle universal placeholders ({date}) through resolver
	// Pass empty values map since we only need {date} which is automatic
	prompt = l.resolver.Resolve(prompt, map[string]string{})

	result.Prompt = prompt
	return result, nil
}
