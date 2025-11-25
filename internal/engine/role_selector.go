package engine

import (
	"fmt"

	"github.com/grantcarthew/start/internal/domain"
)

// RoleSelector selects roles based on precedence rules
type RoleSelector struct{}

// NewRoleSelector creates a new role selector
func NewRoleSelector() *RoleSelector {
	return &RoleSelector{}
}

// SelectionContext contains the context for role selection
type SelectionContext struct {
	RoleFlag    string // --role flag value
	TaskRole    string // role field from task (if executing task)
	DefaultRole string // default_role from settings
}

// Select chooses a role based on precedence rules:
// 1. CLI --role flag (highest priority)
// 2. Task role field (if executing a task)
// 3. default_role setting
// 4. First role in config (TOML order)
func (s *RoleSelector) Select(ctx SelectionContext, roles map[string]domain.Role) (domain.Role, error) {
	if len(roles) == 0 {
		return domain.Role{}, fmt.Errorf("no roles defined in configuration")
	}

	var selectedName string

	// Precedence 1: CLI --role flag
	if ctx.RoleFlag != "" {
		selectedName = ctx.RoleFlag
	} else if ctx.TaskRole != "" {
		// Precedence 2: Task role field
		selectedName = ctx.TaskRole
	} else if ctx.DefaultRole != "" {
		// Precedence 3: default_role setting
		selectedName = ctx.DefaultRole
	} else {
		// Precedence 4: First role in config
		// Note: Go map iteration is not ordered, but the config loader
		// should preserve TOML order in the map or provide a separate
		// ordering mechanism. For now, we'll require default_role.
		return domain.Role{}, fmt.Errorf("no role specified: use --role flag or set default_role in settings")
	}

	// Get the selected role
	role, ok := roles[selectedName]
	if !ok {
		return domain.Role{}, fmt.Errorf("role %q not found in configuration", selectedName)
	}

	// Set the name field
	role.Name = selectedName

	return role, nil
}
