package engine

import (
	"strings"
	"testing"

	"github.com/grantcarthew/start/internal/domain"
	"github.com/grantcarthew/start/test/mocks"
)

func TestTaskLoader_LoadTask(t *testing.T) {
	tests := []struct {
		name         string
		task         domain.Task
		instructions string
		fileContents map[string]string
		cmdOutput    map[string]string
		wantPrompt   string
		wantErr      bool
	}{
		{
			name: "simple prompt with instructions",
			task: domain.Task{
				Name:   "help",
				Prompt: "Help me with: {instructions}",
			},
			instructions: "debugging this code",
			wantPrompt:   "Help me with: debugging this code",
			wantErr:      false,
		},
		{
			name: "prompt with no instructions defaults to None",
			task: domain.Task{
				Name:   "help",
				Prompt: "Help me with: {instructions}",
			},
			instructions: "",
			wantPrompt:   "Help me with: None",
			wantErr:      false,
		},
		{
			name: "prompt with file and instructions",
			task: domain.Task{
				Name:   "doc-review",
				File:   "README.md",
				Prompt: "Review this documentation:\n\n{file_contents}\n\nFocus: {instructions}",
			},
			fileContents: map[string]string{
				"README.md": "# Project\n\nDocumentation here.",
			},
			instructions: "clarity and examples",
			wantPrompt:   "Review this documentation:\n\n# Project\n\nDocumentation here.\n\nFocus: clarity and examples",
			wantErr:      false,
		},
		{
			name: "prompt with command and instructions",
			task: domain.Task{
				Name:    "git-review",
				Command: "git diff --staged",
				Prompt:  "Review changes:\n\n{command_output}\n\nInstructions: {instructions}",
			},
			cmdOutput: map[string]string{
				"git diff --staged": "diff --git a/file.go",
			},
			instructions: "focus on security",
			wantPrompt:   "Review changes:\n\ndiff --git a/file.go\n\nInstructions: focus on security",
			wantErr:      false,
		},
		{
			name: "prompt with file command and instructions",
			task: domain.Task{
				Name:    "complex-review",
				File:    "template.md",
				Command: "git log -1",
				Prompt:  "Template: {file_contents}\n\nLog: {command_output}\n\n{instructions}",
			},
			fileContents: map[string]string{
				"template.md": "Review Template",
			},
			cmdOutput: map[string]string{
				"git log -1": "commit abc123",
			},
			instructions: "check commit message",
			wantPrompt:   "Template: Review Template\n\nLog: commit abc123\n\ncheck commit message",
			wantErr:      false,
		},
		{
			name: "missing required field errors",
			task: domain.Task{
				Name: "empty-task",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			fs := mocks.NewMockFileSystem()
			cmdRunner := mocks.NewMockCommandRunner()

			// Setup file contents
			for path, content := range tt.fileContents {
				fs.Files[path] = content
			}

			// Setup command outputs
			for cmd, output := range tt.cmdOutput {
				cmdRunner.Outputs[cmd] = output
			}

			// Create components
			utdProcessor := NewUTDProcessor(fs, cmdRunner, ".")
			resolver := NewPlaceholderResolver()
			loader := NewTaskLoader(utdProcessor, resolver)

			// Execute
			result, err := loader.LoadTask(tt.task, tt.instructions, "bash", 30)

			// Check error
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Check prompt (trim for easier comparison)
			gotPrompt := strings.TrimSpace(result.Prompt)
			wantPrompt := strings.TrimSpace(tt.wantPrompt)
			if gotPrompt != wantPrompt {
				t.Errorf("Prompt mismatch\nGot:  %q\nWant: %q", gotPrompt, wantPrompt)
			}
		})
	}
}

func TestTaskLoader_DatePlaceholder(t *testing.T) {
	// Test that {date} placeholder is resolved
	task := domain.Task{
		Name:   "dated-task",
		Prompt: "Task created on {date}. Instructions: {instructions}",
	}

	fs := mocks.NewMockFileSystem()
	cmdRunner := mocks.NewMockCommandRunner()
	utdProcessor := NewUTDProcessor(fs, cmdRunner, ".")
	resolver := NewPlaceholderResolver()
	loader := NewTaskLoader(utdProcessor, resolver)

	result, err := loader.LoadTask(task, "test instructions", "bash", 30)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check that {date} was replaced (should not contain literal {date})
	if strings.Contains(result.Prompt, "{date}") {
		t.Errorf("Prompt still contains {date} placeholder: %q", result.Prompt)
	}

	// Check that instructions were replaced
	if !strings.Contains(result.Prompt, "test instructions") {
		t.Errorf("Prompt does not contain instructions: %q", result.Prompt)
	}
}
