package engine

import (
	"testing"

	"github.com/grantcarthew/start/internal/domain"
)

func TestTaskResolver_Resolve(t *testing.T) {
	globalTasks := map[string]domain.Task{
		"global-task": {
			Alias:       "gt",
			Description: "Global task",
			Prompt:      "Global prompt",
		},
		"shared-task": {
			Alias:       "st",
			Description: "Shared global",
			Prompt:      "Global version",
		},
	}

	localTasks := map[string]domain.Task{
		"local-task": {
			Alias:       "lt",
			Description: "Local task",
			Prompt:      "Local prompt",
		},
		"shared-task": {
			Alias:       "st-local",
			Description: "Shared local",
			Prompt:      "Local version",
		},
	}

	tests := []struct {
		name        string
		input       string
		wantTask    string // task name
		wantPrompt  string // to verify we got right task
		wantErr     bool
	}{
		{
			name:       "resolve local task by name",
			input:      "local-task",
			wantTask:   "local-task",
			wantPrompt: "Local prompt",
			wantErr:    false,
		},
		{
			name:       "resolve local task by alias",
			input:      "lt",
			wantTask:   "local-task",
			wantPrompt: "Local prompt",
			wantErr:    false,
		},
		{
			name:       "resolve global task by name",
			input:      "global-task",
			wantTask:   "global-task",
			wantPrompt: "Global prompt",
			wantErr:    false,
		},
		{
			name:       "resolve global task by alias",
			input:      "gt",
			wantTask:   "global-task",
			wantPrompt: "Global prompt",
			wantErr:    false,
		},
		{
			name:       "local task overrides global (by name)",
			input:      "shared-task",
			wantTask:   "shared-task",
			wantPrompt: "Local version",
			wantErr:    false,
		},
		{
			name:       "local alias overrides global alias",
			input:      "st-local",
			wantTask:   "shared-task",
			wantPrompt: "Local version",
			wantErr:    false,
		},
		{
			name:    "task not found",
			input:   "nonexistent",
			wantErr: true,
		},
		{
			name:    "alias not found",
			input:   "nx",
			wantErr: true,
		},
	}

	resolver := NewTaskResolver()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := resolver.Resolve(tt.input, localTasks, globalTasks)

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

			if task.Name != tt.wantTask {
				t.Errorf("Task name mismatch: got %q, want %q", task.Name, tt.wantTask)
			}

			if task.Prompt != tt.wantPrompt {
				t.Errorf("Task prompt mismatch: got %q, want %q", task.Prompt, tt.wantPrompt)
			}
		})
	}
}

func TestTaskResolver_ListAllTasks(t *testing.T) {
	globalTasks := map[string]domain.Task{
		"global-only": {
			Description: "Global only task",
		},
		"shared": {
			Description: "Global version",
			Prompt:      "global-prompt",
		},
	}

	localTasks := map[string]domain.Task{
		"local-only": {
			Description: "Local only task",
		},
		"shared": {
			Description: "Local version",
			Prompt:      "local-prompt",
		},
	}

	resolver := NewTaskResolver()
	allTasks := resolver.ListAllTasks(localTasks, globalTasks)

	// Should have 3 tasks: global-only, local-only, shared (local version)
	if len(allTasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(allTasks))
	}

	// Check global-only exists
	if task, ok := allTasks["global-only"]; !ok {
		t.Errorf("Missing global-only task")
	} else if task.Name != "global-only" {
		t.Errorf("global-only task name not set")
	}

	// Check local-only exists
	if task, ok := allTasks["local-only"]; !ok {
		t.Errorf("Missing local-only task")
	} else if task.Name != "local-only" {
		t.Errorf("local-only task name not set")
	}

	// Check shared task is local version
	if task, ok := allTasks["shared"]; !ok {
		t.Errorf("Missing shared task")
	} else {
		if task.Prompt != "local-prompt" {
			t.Errorf("Shared task should be local version, got prompt: %q", task.Prompt)
		}
		if task.Name != "shared" {
			t.Errorf("shared task name not set")
		}
	}
}

func TestTaskResolver_EmptyConfigs(t *testing.T) {
	resolver := NewTaskResolver()

	// Test with empty maps
	_, err := resolver.Resolve("any-task", map[string]domain.Task{}, map[string]domain.Task{})
	if err == nil {
		t.Errorf("Expected error with empty configs, got nil")
	}

	// Test list with empty maps
	allTasks := resolver.ListAllTasks(map[string]domain.Task{}, map[string]domain.Task{})
	if len(allTasks) != 0 {
		t.Errorf("Expected 0 tasks with empty configs, got %d", len(allTasks))
	}
}
