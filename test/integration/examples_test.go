package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/grantcarthew/start/internal/adapters"
	"github.com/grantcarthew/start/internal/config"
	"github.com/grantcarthew/start/internal/domain"
	"github.com/pelletier/go-toml/v2"
)

// Helper to load config from a directory (examples are organized as if they're already in .start/ or .config/start/)
func loadExampleConfig(fs *adapters.RealFileSystem, dir string) (domain.Config, error) {
	cfg := domain.Config{
		Agents:   make(map[string]domain.Agent),
		Roles:    make(map[string]domain.Role),
		Contexts: make(map[string]domain.Context),
		Tasks:    make(map[string]domain.Task),
	}

	// Try to load each file
	files := []string{"config.toml", "agents.toml", "roles.toml", "contexts.toml", "tasks.toml"}
	for _, file := range files {
		path := filepath.Join(dir, file)
		if _, err := os.Stat(path); err == nil {
			data, err := fs.ReadFile(path)
			if err != nil {
				return cfg, err
			}

			// Parse based on file type
			switch file {
			case "config.toml":
				var parsed struct {
					Settings domain.Settings `toml:"settings"`
				}
				if err := toml.Unmarshal(data, &parsed); err != nil {
					return cfg, err
				}
				cfg.Settings = parsed.Settings
			case "agents.toml":
				var parsed struct {
					Agents map[string]domain.Agent `toml:"agents"`
				}
				if err := toml.Unmarshal(data, &parsed); err != nil {
					return cfg, err
				}
				for name, agent := range parsed.Agents {
					agent.Name = name
					cfg.Agents[name] = agent
				}
			case "roles.toml":
				var parsed struct {
					Roles map[string]domain.Role `toml:"roles"`
				}
				if err := toml.Unmarshal(data, &parsed); err != nil {
					return cfg, err
				}
				for name, role := range parsed.Roles {
					role.Name = name
					cfg.Roles[name] = role
				}
			case "contexts.toml":
				var parsed struct {
					Contexts map[string]domain.Context `toml:"contexts"`
				}
				if err := toml.Unmarshal(data, &parsed); err != nil {
					return cfg, err
				}
				for name, ctx := range parsed.Contexts {
					ctx.Name = name
					cfg.Contexts[name] = ctx
					cfg.ContextOrder = append(cfg.ContextOrder, name)
				}
			case "tasks.toml":
				var parsed struct {
					Tasks map[string]domain.Task `toml:"tasks"`
				}
				if err := toml.Unmarshal(data, &parsed); err != nil {
					return cfg, err
				}
				for name, task := range parsed.Tasks {
					task.Name = name
					cfg.Tasks[name] = task
				}
			}
		}
	}

	return cfg, nil
}

func TestMinimalExample(t *testing.T) {
	fs := &adapters.RealFileSystem{}
	validator := config.NewValidator()

	// Load minimal global config
	globalCfg, err := loadExampleConfig(fs, "../../examples/minimal/global")
	if err != nil {
		t.Fatalf("Failed to load minimal global config: %v", err)
	}

	// Validate
	if err := validator.Validate(globalCfg); err != nil {
		t.Errorf("Minimal global config validation failed: %v", err)
	}

	// Check some basic expectations
	if len(globalCfg.Agents) == 0 {
		t.Error("Expected at least one agent in minimal config")
	}
}

func TestCompleteExample(t *testing.T) {
	fs := &adapters.RealFileSystem{}
	validator := config.NewValidator()

	// Load complete global config
	globalCfg, err := loadExampleConfig(fs, "../../examples/complete/global")
	if err != nil {
		t.Fatalf("Failed to load complete global config: %v", err)
	}

	// Load complete local config
	localCfg, err := loadExampleConfig(fs, "../../examples/complete/local")
	if err != nil {
		t.Fatalf("Failed to load complete local config: %v", err)
	}

	// Merge
	merged := config.Merge(globalCfg, localCfg)

	// Validate
	if err := validator.Validate(merged); err != nil {
		t.Errorf("Complete merged config validation failed: %v", err)
	}

	// Check some expectations
	if len(merged.Agents) == 0 {
		t.Error("Expected at least one agent in complete config")
	}

	if len(merged.Roles) == 0 {
		t.Error("Expected at least one role in complete config")
	}

	if len(merged.Contexts) == 0 {
		t.Error("Expected at least one context in complete config")
	}
}

func TestRealWorldExample(t *testing.T) {
	fs := &adapters.RealFileSystem{}
	validator := config.NewValidator()

	// Load real-world global config
	globalCfg, err := loadExampleConfig(fs, "../../examples/real-world/global")
	if err != nil {
		t.Fatalf("Failed to load real-world global config: %v", err)
	}

	// Load real-world local config
	localCfg, err := loadExampleConfig(fs, "../../examples/real-world/local")
	if err != nil {
		t.Fatalf("Failed to load real-world local config: %v", err)
	}

	// Merge
	merged := config.Merge(globalCfg, localCfg)

	// Validate
	if err := validator.Validate(merged); err != nil {
		t.Errorf("Real-world merged config validation failed: %v", err)
	}
}
