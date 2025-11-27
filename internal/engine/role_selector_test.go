package engine

import (
	"testing"

	"github.com/grantcarthew/start/internal/domain"
)

func TestRoleSelector_Select_Precedence(t *testing.T) {
	selector := NewRoleSelector()

	roles := map[string]domain.Role{
		"flag-role":    {Prompt: "Flag role"},
		"task-role":    {Prompt: "Task role"},
		"default-role": {Prompt: "Default role"},
	}

	tests := []struct {
		name          string
		ctx           SelectionContext
		expectedRole  string
		expectError   bool
		errorContains string
	}{
		{
			name: "precedence 1: CLI flag takes priority",
			ctx: SelectionContext{
				RoleFlag:    "flag-role",
				TaskRole:    "task-role",
				DefaultRole: "default-role",
			},
			expectedRole: "flag-role",
			expectError:  false,
		},
		{
			name: "precedence 2: task role when no flag",
			ctx: SelectionContext{
				RoleFlag:    "",
				TaskRole:    "task-role",
				DefaultRole: "default-role",
			},
			expectedRole: "task-role",
			expectError:  false,
		},
		{
			name: "precedence 3: default role when no flag or task",
			ctx: SelectionContext{
				RoleFlag:    "",
				TaskRole:    "",
				DefaultRole: "default-role",
			},
			expectedRole: "default-role",
			expectError:  false,
		},
		{
			name: "error when no role specified",
			ctx: SelectionContext{
				RoleFlag:    "",
				TaskRole:    "",
				DefaultRole: "",
			},
			expectError:   true,
			errorContains: "no role specified",
		},
		{
			name: "error when role not found - flag",
			ctx: SelectionContext{
				RoleFlag:    "nonexistent",
				TaskRole:    "",
				DefaultRole: "",
			},
			expectError:   true,
			errorContains: "not found",
		},
		{
			name: "error when role not found - task",
			ctx: SelectionContext{
				RoleFlag:    "",
				TaskRole:    "nonexistent",
				DefaultRole: "",
			},
			expectError:   true,
			errorContains: "not found",
		},
		{
			name: "error when role not found - default",
			ctx: SelectionContext{
				RoleFlag:    "",
				TaskRole:    "",
				DefaultRole: "nonexistent",
			},
			expectError:   true,
			errorContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			role, err := selector.Select(tt.ctx, roles)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else if tt.errorContains != "" && !contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain %q, got: %v", tt.errorContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if role.Name != tt.expectedRole {
					t.Errorf("Expected role %q, got %q", tt.expectedRole, role.Name)
				}
			}
		})
	}
}

func TestRoleSelector_Select_NoRolesDefined(t *testing.T) {
	selector := NewRoleSelector()

	roles := map[string]domain.Role{}

	ctx := SelectionContext{
		RoleFlag:    "any-role",
		TaskRole:    "",
		DefaultRole: "",
	}

	_, err := selector.Select(ctx, roles)

	if err == nil {
		t.Error("Expected error when no roles defined")
	}

	if !contains(err.Error(), "no roles defined") {
		t.Errorf("Expected error about no roles defined, got: %v", err)
	}
}

func TestRoleSelector_Select_RoleContentPreserved(t *testing.T) {
	selector := NewRoleSelector()

	roles := map[string]domain.Role{
		"test-role": {
			Prompt:      "Test prompt",
			File:        "/path/to/file.md",
			Description: "Test description",
		},
	}

	ctx := SelectionContext{
		RoleFlag: "test-role",
	}

	role, err := selector.Select(ctx, roles)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if role.Name != "test-role" {
		t.Errorf("Expected role name %q, got %q", "test-role", role.Name)
	}

	if role.Prompt != "Test prompt" {
		t.Errorf("Expected prompt %q, got %q", "Test prompt", role.Prompt)
	}

	if role.File != "/path/to/file.md" {
		t.Errorf("Expected file %q, got %q", "/path/to/file.md", role.File)
	}

	if role.Description != "Test description" {
		t.Errorf("Expected description %q, got %q", "Test description", role.Description)
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && indexOfSubstring(s, substr) >= 0))
}

func indexOfSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
