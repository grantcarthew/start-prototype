package config

import "github.com/grantcarthew/start/internal/domain"

// Merge merges local config into global config
// Local config takes precedence over global config
func Merge(global, local domain.Config) domain.Config {
	result := domain.Config{
		Settings: mergeSettings(global.Settings, local.Settings),
		Agents:   mergeAgents(global.Agents, local.Agents),
		Roles:    mergeRoles(global.Roles, local.Roles),
		Contexts: mergeContexts(global.Contexts, local.Contexts),
		Tasks:    mergeTasks(global.Tasks, local.Tasks),
	}

	return result
}

// mergeSettings merges settings - local overrides global per-field
func mergeSettings(global, local domain.Settings) domain.Settings {
	result := global

	// Override individual fields if set in local
	if local.DefaultAgent != "" {
		result.DefaultAgent = local.DefaultAgent
	}
	if local.DefaultRole != "" {
		result.DefaultRole = local.DefaultRole
	}
	if local.LogLevel != "" {
		result.LogLevel = local.LogLevel
	}
	if local.Shell != "" {
		result.Shell = local.Shell
	}
	if local.CommandTimeout != 0 {
		result.CommandTimeout = local.CommandTimeout
	}
	// AssetDownload is a bool, so we need to check if it was explicitly set
	// For now, we always take the local value (false is a valid override)
	if local.AssetDownload != global.AssetDownload || local.AssetDownload {
		result.AssetDownload = local.AssetDownload
	}
	if local.AssetRepo != "" {
		result.AssetRepo = local.AssetRepo
	}
	if local.AssetPath != "" {
		result.AssetPath = local.AssetPath
	}

	return result
}

// mergeAgents combines agents from both configs
// Local agent replaces global agent with same name
func mergeAgents(global, local map[string]domain.Agent) map[string]domain.Agent {
	result := make(map[string]domain.Agent)

	// Copy all global agents
	for name, agent := range global {
		result[name] = agent
	}

	// Override/add local agents
	for name, agent := range local {
		result[name] = agent
	}

	return result
}

// mergeRoles combines roles from both configs
// Local role replaces global role with same name
func mergeRoles(global, local map[string]domain.Role) map[string]domain.Role {
	result := make(map[string]domain.Role)

	// Copy all global roles
	for name, role := range global {
		result[name] = role
	}

	// Override/add local roles
	for name, role := range local {
		result[name] = role
	}

	return result
}

// mergeContexts combines contexts from both configs
// Local context replaces global context with same name
func mergeContexts(global, local map[string]domain.Context) map[string]domain.Context {
	result := make(map[string]domain.Context)

	// Copy all global contexts
	for name, ctx := range global {
		result[name] = ctx
	}

	// Override/add local contexts
	for name, ctx := range local {
		result[name] = ctx
	}

	return result
}

// mergeTasks combines tasks from both configs
// Local task replaces global task with same name
func mergeTasks(global, local map[string]domain.Task) map[string]domain.Task {
	result := make(map[string]domain.Task)

	// Copy all global tasks
	for name, task := range global {
		result[name] = task
	}

	// Override/add local tasks
	for name, task := range local {
		result[name] = task
	}

	return result
}
