package config

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/grantcarthew/start/internal/domain"
)

// BackupHelper handles creating timestamped backups of config files
type BackupHelper struct {
	fs domain.FileSystem
}

// NewBackupHelper creates a new backup helper
func NewBackupHelper(fs domain.FileSystem) *BackupHelper {
	return &BackupHelper{fs: fs}
}

// CreateBackup creates a timestamped backup of a config file
// Format: <filename>.YYYY-MM-DD-HHMMSS.toml
// Returns the backup path on success
func (b *BackupHelper) CreateBackup(configPath string) (string, error) {
	// Check if file exists
	if !b.fs.Exists(configPath) {
		return "", fmt.Errorf("config file does not exist: %s", configPath)
	}

	// Read current file
	data, err := b.fs.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("failed to read config file: %w", err)
	}

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("2006-01-02-150405")
	dir := filepath.Dir(configPath)
	base := filepath.Base(configPath)

	// Remove .toml extension if present
	if filepath.Ext(base) == ".toml" {
		base = base[:len(base)-5]
	}

	backupName := fmt.Sprintf("%s.%s.toml", base, timestamp)
	backupPath := filepath.Join(dir, backupName)

	// Write backup
	if err := b.fs.WriteFile(backupPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write backup file: %w", err)
	}

	return backupPath, nil
}
