package config

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/grantcarthew/start/test/mocks"
)

func TestBackupHelper_CreateBackup(t *testing.T) {
	tests := []struct {
		name        string
		setupFS     func(*mocks.MockFileSystem)
		configPath  string
		wantErr     bool
		wantBackup  bool
		errContains string
	}{
		{
			name: "creates backup successfully",
			setupFS: func(fs *mocks.MockFileSystem) {
				fs.Files["/test/agents.toml"] = "[agents.test]\nbin = \"test\""
			},
			configPath: "/test/agents.toml",
			wantErr:    false,
			wantBackup: true,
		},
		{
			name:        "returns error when file does not exist",
			setupFS:     func(fs *mocks.MockFileSystem) {},
			configPath:  "/test/nonexistent.toml",
			wantErr:     true,
			wantBackup:  false,
			errContains: "does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := mocks.NewMockFileSystem()
			tt.setupFS(fs)

			helper := NewBackupHelper(fs)
			backupPath, err := helper.CreateBackup(tt.configPath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error should contain %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tt.wantBackup {
				// Check backup path format (contains "agents" from "agents.toml")
				baseName := strings.TrimSuffix(filepath.Base(tt.configPath), ".toml")
				if !strings.Contains(backupPath, baseName) {
					t.Errorf("backup path should contain base name %q, got %q", baseName, backupPath)
				}

				// Check backup was created
				if !fs.Exists(backupPath) {
					t.Errorf("backup file was not created at %s", backupPath)
				}

				// Check content matches
				originalContent := fs.Files[tt.configPath]
				backupContent := fs.Files[backupPath]
				if originalContent != backupContent {
					t.Errorf("backup content mismatch")
				}
			}
		})
	}
}
