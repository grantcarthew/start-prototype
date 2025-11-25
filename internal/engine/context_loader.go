package engine

import (
	"github.com/grantcarthew/start/internal/domain"
)

// ContextLoader loads and processes context documents
type ContextLoader struct {
	utdProcessor *UTDProcessor
}

// NewContextLoader creates a new context loader
func NewContextLoader(utdProcessor *UTDProcessor) *ContextLoader {
	return &ContextLoader{
		utdProcessor: utdProcessor,
	}
}

// CommandType represents the type of command being executed
type CommandType string

const (
	CommandTypeInteractive CommandType = "interactive" // start (includes all contexts)
	CommandTypePrompt      CommandType = "prompt"      // start prompt (required only)
	CommandTypeTask        CommandType = "task"        // start task (required only)
)

// LoadedContext represents a processed context
type LoadedContext struct {
	Name     string
	Content  string
	FilePath string // For display purposes
	Warnings []string
}

// LoadContexts loads and processes contexts based on command type
// Returns loaded contexts in definition order
func (l *ContextLoader) LoadContexts(
	contexts map[string]domain.Context,
	contextOrder []string,
	commandType CommandType,
	defaultShell string,
	defaultTimeout int,
) []LoadedContext {
	var result []LoadedContext

	// Process contexts in order
	for _, name := range contextOrder {
		ctx, ok := contexts[name]
		if !ok {
			continue
		}

		// Filter by required field based on command type
		if commandType != CommandTypeInteractive && !ctx.Required {
			// Skip optional contexts for prompt and task commands
			continue
		}

		// Process through UTD
		utdInput := UTDInput{
			File:           ctx.File,
			Command:        ctx.Command,
			Prompt:         ctx.Prompt,
			Shell:          ctx.Shell,
			CommandTimeout: ctx.CommandTimeout,
		}

		utdResult := l.utdProcessor.Process(utdInput, defaultShell, defaultTimeout)

		// Skip if UTD processing failed
		if utdResult.Skipped {
			result = append(result, LoadedContext{
				Name:     name,
				Warnings: utdResult.Warnings,
			})
			continue
		}

		// Add successfully loaded context
		result = append(result, LoadedContext{
			Name:     name,
			Content:  utdResult.Content,
			FilePath: utdResult.FilePath,
			Warnings: utdResult.Warnings,
		})
	}

	return result
}
