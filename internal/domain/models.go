package domain

import "time"

// Config represents the merged configuration
type Config struct {
	Settings     Settings
	Agents       map[string]Agent
	Roles        map[string]Role
	Contexts     map[string]Context
	ContextOrder []string // Preserves TOML definition order
	Tasks        map[string]Task
}

// Settings from config.toml [settings]
type Settings struct {
	DefaultAgent   string `toml:"default_agent"`
	DefaultRole    string `toml:"default_role"`
	LogLevel       string `toml:"log_level"`
	Shell          string `toml:"shell"`
	CommandTimeout int    `toml:"command_timeout"`
	AssetDownload  bool   `toml:"asset_download"`
	AssetRepo      string `toml:"asset_repo"`
	AssetPath      string `toml:"asset_path"`
}

// Agent from agents.toml [agents.<name>]
type Agent struct {
	Name         string
	Bin          string            `toml:"bin"`
	Command      string            `toml:"command"`
	Description  string            `toml:"description"`
	URL          string            `toml:"url"`
	ModelsURL    string            `toml:"models_url"`
	DefaultModel string            `toml:"default_model"`
	Models       map[string]string `toml:"models"`
}

// Role from roles.toml [roles.<name>] (UTD pattern)
type Role struct {
	Name           string
	Description    string `toml:"description"`
	File           string `toml:"file"`
	Command        string `toml:"command"`
	Prompt         string `toml:"prompt"`
	Shell          string `toml:"shell"`
	CommandTimeout int    `toml:"command_timeout"`
}

// Context from contexts.toml [contexts.<name>] (UTD pattern)
type Context struct {
	Name           string
	Description    string `toml:"description"`
	File           string `toml:"file"`
	Command        string `toml:"command"`
	Prompt         string `toml:"prompt"`
	Required       bool   `toml:"required"`
	Shell          string `toml:"shell"`
	CommandTimeout int    `toml:"command_timeout"`
}

// Task from tasks.toml [tasks.<name>] (UTD pattern)
type Task struct {
	Name           string
	Alias          string `toml:"alias"`
	Description    string `toml:"description"`
	Role           string `toml:"role"`
	Agent          string `toml:"agent"`
	File           string `toml:"file"`
	Command        string `toml:"command"`
	Prompt         string `toml:"prompt"`
	Shell          string `toml:"shell"`
	CommandTimeout int    `toml:"command_timeout"`
}

// AssetMeta from .meta.toml files
type AssetMeta struct {
	Type        string    `toml:"type"`
	Category    string    `toml:"category"`
	Name        string    `toml:"name"`
	Description string    `toml:"description"`
	Tags        string    `toml:"tags"`
	Bin         string    `toml:"bin"`
	SHA         string    `toml:"sha"`
	Size        int64     `toml:"size"`
	Created     time.Time `toml:"created"`
	Updated     time.Time `toml:"updated"`
}

// CachedAsset represents an asset in the cache
type CachedAsset struct {
	Type     string
	Category string
	Name     string
	Meta     AssetMeta
}
