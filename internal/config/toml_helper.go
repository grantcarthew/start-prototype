package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/grantcarthew/start/internal/domain"
	"github.com/pelletier/go-toml/v2"
)

// TOMLHelper provides utilities for reading and writing TOML files
type TOMLHelper struct {
	fs domain.FileSystem
}

// NewTOMLHelper creates a new TOML helper
func NewTOMLHelper(fs domain.FileSystem) *TOMLHelper {
	return &TOMLHelper{fs: fs}
}

// ReadAgentsFile reads the agents.toml file from a directory
func (h *TOMLHelper) ReadAgentsFile(dir string) (map[string]domain.Agent, error) {
	path := filepath.Join(dir, "agents.toml")
	data, err := h.fs.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]domain.Agent), nil
		}
		return nil, fmt.Errorf("failed to read agents file: %w", err)
	}

	var parsed struct {
		Agents map[string]domain.Agent `toml:"agents"`
	}

	if err := toml.Unmarshal(data, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse agents file: %w", err)
	}

	if parsed.Agents == nil {
		return make(map[string]domain.Agent), nil
	}

	return parsed.Agents, nil
}

// WriteAgentsFile writes agents to the agents.toml file
func (h *TOMLHelper) WriteAgentsFile(dir string, agents map[string]domain.Agent) error {
	// Ensure directory exists
	if err := h.fs.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Prepare structure for marshaling
	tomlData := struct {
		Agents map[string]domain.Agent `toml:"agents"`
	}{
		Agents: agents,
	}

	// Marshal to TOML
	data, err := toml.Marshal(tomlData)
	if err != nil {
		return fmt.Errorf("failed to marshal agents: %w", err)
	}

	// Write file
	path := filepath.Join(dir, "agents.toml")
	if err := h.fs.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write agents file: %w", err)
	}

	return nil
}

// ReadSettingsFile reads the config.toml file (settings section)
func (h *TOMLHelper) ReadSettingsFile(dir string) (domain.Settings, error) {
	path := filepath.Join(dir, "config.toml")
	data, err := h.fs.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return domain.Settings{}, nil
		}
		return domain.Settings{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var parsed struct {
		Settings domain.Settings `toml:"settings"`
	}

	if err := toml.Unmarshal(data, &parsed); err != nil {
		return domain.Settings{}, fmt.Errorf("failed to parse config file: %w", err)
	}

	return parsed.Settings, nil
}

// WriteSettingsFile writes settings to the config.toml file
func (h *TOMLHelper) WriteSettingsFile(dir string, settings domain.Settings) error {
	// Ensure directory exists
	if err := h.fs.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Prepare structure for marshaling
	tomlData := struct {
		Settings domain.Settings `toml:"settings"`
	}{
		Settings: settings,
	}

	// Marshal to TOML
	data, err := toml.Marshal(tomlData)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	// Write file
	path := filepath.Join(dir, "config.toml")
	if err := h.fs.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetGlobalDir returns the global config directory path
func (h *TOMLHelper) GetGlobalDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".config", "start"), nil
}

// GetLocalDir returns the local config directory path for the given working directory
func (h *TOMLHelper) GetLocalDir(workDir string) string {
	return filepath.Join(workDir, ".start")
}

// GetConfigPath returns the path to config.toml in the given directory
func (h *TOMLHelper) GetConfigPath(dir string) string {
	return filepath.Join(dir, "config.toml")
}

// GetFS returns the filesystem used by this helper
func (h *TOMLHelper) GetFS() domain.FileSystem {
	return h.fs
}

// ReadRolesFile reads the roles.toml file from a directory
func (h *TOMLHelper) ReadRolesFile(dir string) (map[string]domain.Role, error) {
	path := filepath.Join(dir, "roles.toml")
	data, err := h.fs.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]domain.Role), nil
		}
		return nil, fmt.Errorf("failed to read roles file: %w", err)
	}

	var parsed struct {
		Roles map[string]domain.Role `toml:"roles"`
	}

	if err := toml.Unmarshal(data, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse roles file: %w", err)
	}

	if parsed.Roles == nil {
		return make(map[string]domain.Role), nil
	}

	return parsed.Roles, nil
}

// WriteRolesFile writes roles to the roles.toml file
func (h *TOMLHelper) WriteRolesFile(dir string, roles map[string]domain.Role) error {
	// Ensure directory exists
	if err := h.fs.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Prepare structure for marshaling
	tomlData := struct {
		Roles map[string]domain.Role `toml:"roles"`
	}{
		Roles: roles,
	}

	// Marshal to TOML
	data, err := toml.Marshal(tomlData)
	if err != nil {
		return fmt.Errorf("failed to marshal roles: %w", err)
	}

	// Write file
	path := filepath.Join(dir, "roles.toml")
	if err := h.fs.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write roles file: %w", err)
	}

	return nil
}
