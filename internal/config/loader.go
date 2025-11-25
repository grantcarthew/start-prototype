package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/grantcarthew/start/internal/domain"
	"github.com/pelletier/go-toml/v2"
)

// Loader handles loading configuration from files
type Loader struct {
	fs domain.FileSystem
}

// NewLoader creates a new config loader
func NewLoader(fs domain.FileSystem) *Loader {
	return &Loader{
		fs: fs,
	}
}

// LoadGlobal loads configuration from global directory (~/.config/start/)
func (l *Loader) LoadGlobal() (domain.Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return domain.Config{}, fmt.Errorf("failed to get home directory: %w", err)
	}

	globalDir := filepath.Join(homeDir, ".config", "start")
	return l.loadFromDir(globalDir)
}

// LoadLocal loads configuration from local directory (./.start/)
func (l *Loader) LoadLocal(workDir string) (domain.Config, error) {
	localDir := filepath.Join(workDir, ".start")
	return l.loadFromDir(localDir)
}

// loadFromDir loads all config files from a directory
func (l *Loader) loadFromDir(dir string) (domain.Config, error) {
	config := domain.Config{
		Agents:   make(map[string]domain.Agent),
		Roles:    make(map[string]domain.Role),
		Contexts: make(map[string]domain.Context),
		Tasks:    make(map[string]domain.Task),
	}

	// Load config.toml (settings)
	if err := l.loadSettings(dir, &config); err != nil {
		// Settings file is optional, only error if file exists but fails to parse
		if !os.IsNotExist(err) {
			return config, fmt.Errorf("failed to load settings: %w", err)
		}
	}

	// Load agents.toml
	if err := l.loadAgents(dir, &config); err != nil {
		if !os.IsNotExist(err) {
			return config, fmt.Errorf("failed to load agents: %w", err)
		}
	}

	// Load roles.toml
	if err := l.loadRoles(dir, &config); err != nil {
		if !os.IsNotExist(err) {
			return config, fmt.Errorf("failed to load roles: %w", err)
		}
	}

	// Load contexts.toml
	if err := l.loadContexts(dir, &config); err != nil {
		if !os.IsNotExist(err) {
			return config, fmt.Errorf("failed to load contexts: %w", err)
		}
	}

	// Load tasks.toml
	if err := l.loadTasks(dir, &config); err != nil {
		if !os.IsNotExist(err) {
			return config, fmt.Errorf("failed to load tasks: %w", err)
		}
	}

	return config, nil
}

// loadSettings loads settings from config.toml
func (l *Loader) loadSettings(dir string, config *domain.Config) error {
	path := filepath.Join(dir, "config.toml")
	data, err := l.fs.ReadFile(path)
	if err != nil {
		return err
	}

	var parsed struct {
		Settings domain.Settings `toml:"settings"`
	}

	if err := toml.Unmarshal(data, &parsed); err != nil {
		return fmt.Errorf("failed to parse %s: %w", path, err)
	}

	config.Settings = parsed.Settings
	return nil
}

// loadAgents loads agents from agents.toml
func (l *Loader) loadAgents(dir string, config *domain.Config) error {
	path := filepath.Join(dir, "agents.toml")
	data, err := l.fs.ReadFile(path)
	if err != nil {
		return err
	}

	var parsed struct {
		Agents map[string]domain.Agent `toml:"agents"`
	}

	if err := toml.Unmarshal(data, &parsed); err != nil {
		return fmt.Errorf("failed to parse %s: %w", path, err)
	}

	// Set the Name field for each agent (it's the map key)
	for name, agent := range parsed.Agents {
		agent.Name = name
		config.Agents[name] = agent
	}

	return nil
}

// loadRoles loads roles from roles.toml
func (l *Loader) loadRoles(dir string, config *domain.Config) error {
	path := filepath.Join(dir, "roles.toml")
	data, err := l.fs.ReadFile(path)
	if err != nil {
		return err
	}

	var parsed struct {
		Roles map[string]domain.Role `toml:"roles"`
	}

	if err := toml.Unmarshal(data, &parsed); err != nil {
		return fmt.Errorf("failed to parse %s: %w", path, err)
	}

	// Set the Name field for each role (it's the map key)
	for name, role := range parsed.Roles {
		role.Name = name
		config.Roles[name] = role
	}

	return nil
}

// loadContexts loads contexts from contexts.toml
func (l *Loader) loadContexts(dir string, config *domain.Config) error {
	path := filepath.Join(dir, "contexts.toml")
	data, err := l.fs.ReadFile(path)
	if err != nil {
		return err
	}

	var parsed struct {
		Contexts map[string]domain.Context `toml:"contexts"`
	}

	// Use decoder to preserve order (go-toml v2 preserves map iteration order)
	if err := toml.Unmarshal(data, &parsed); err != nil {
		return fmt.Errorf("failed to parse %s: %w", path, err)
	}

	// Set the Name field for each context and preserve order
	// Note: go-toml/v2 preserves the order of map keys during iteration
	for name, ctx := range parsed.Contexts {
		ctx.Name = name
		config.Contexts[name] = ctx
		config.ContextOrder = append(config.ContextOrder, name)
	}

	return nil
}

// loadTasks loads tasks from tasks.toml
func (l *Loader) loadTasks(dir string, config *domain.Config) error {
	path := filepath.Join(dir, "tasks.toml")
	data, err := l.fs.ReadFile(path)
	if err != nil {
		return err
	}

	var parsed struct {
		Tasks map[string]domain.Task `toml:"tasks"`
	}

	if err := toml.Unmarshal(data, &parsed); err != nil {
		return fmt.Errorf("failed to parse %s: %w", path, err)
	}

	// Set the Name field for each task (it's the map key)
	for name, task := range parsed.Tasks {
		task.Name = name
		config.Tasks[name] = task
	}

	return nil
}
