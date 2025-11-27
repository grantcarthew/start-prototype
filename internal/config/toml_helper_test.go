package config

import (
	"testing"

	"github.com/grantcarthew/start/internal/domain"
	"github.com/grantcarthew/start/test/mocks"
)

func TestTOMLHelper_ReadWriteAgents(t *testing.T) {
	fs := mocks.NewMockFileSystem()
	helper := NewTOMLHelper(fs)

	// Test writing agents
	agents := map[string]domain.Agent{
		"claude": {
			Name:        "claude",
			Bin:         "claude",
			Description: "Claude AI",
			Command:     "{bin} {prompt}",
			Models: map[string]string{
				"sonnet": "claude-sonnet-4-20250929",
			},
			DefaultModel: "sonnet",
		},
	}

	dir := "/test"
	err := helper.WriteAgentsFile(dir, agents)
	if err != nil {
		t.Fatalf("WriteAgentsFile failed: %v", err)
	}

	// Verify file was created
	path := dir + "/agents.toml"
	if !fs.Exists(path) {
		t.Fatal("agents.toml was not created")
	}

	// Test reading agents back
	readAgents, err := helper.ReadAgentsFile(dir)
	if err != nil {
		t.Fatalf("ReadAgentsFile failed: %v", err)
	}

	if len(readAgents) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(readAgents))
	}

	claude, ok := readAgents["claude"]
	if !ok {
		t.Fatal("claude agent not found")
	}

	if claude.Bin != "claude" {
		t.Errorf("expected bin=claude, got %s", claude.Bin)
	}

	if claude.Description != "Claude AI" {
		t.Errorf("expected description=Claude AI, got %s", claude.Description)
	}

	if len(claude.Models) != 1 {
		t.Errorf("expected 1 model, got %d", len(claude.Models))
	}
}

func TestTOMLHelper_ReadAgentsFile_Empty(t *testing.T) {
	fs := mocks.NewMockFileSystem()
	helper := NewTOMLHelper(fs)

	// Read from non-existent file should return empty map
	agents, err := helper.ReadAgentsFile("/nonexistent")
	if err != nil {
		t.Fatalf("expected no error for non-existent file, got: %v", err)
	}

	if len(agents) != 0 {
		t.Errorf("expected empty map, got %d agents", len(agents))
	}
}

func TestTOMLHelper_ReadWriteSettings(t *testing.T) {
	fs := mocks.NewMockFileSystem()
	helper := NewTOMLHelper(fs)

	settings := domain.Settings{
		DefaultAgent: "claude",
		DefaultRole:  "default",
		Shell:        "bash",
	}

	dir := "/test"
	err := helper.WriteSettingsFile(dir, settings)
	if err != nil {
		t.Fatalf("WriteSettingsFile failed: %v", err)
	}

	// Read back
	readSettings, err := helper.ReadSettingsFile(dir)
	if err != nil {
		t.Fatalf("ReadSettingsFile failed: %v", err)
	}

	if readSettings.DefaultAgent != "claude" {
		t.Errorf("expected DefaultAgent=claude, got %s", readSettings.DefaultAgent)
	}

	if readSettings.DefaultRole != "default" {
		t.Errorf("expected DefaultRole=default, got %s", readSettings.DefaultRole)
	}
}

func TestTOMLHelper_ReadWriteRoles(t *testing.T) {
	fs := mocks.NewMockFileSystem()
	helper := NewTOMLHelper(fs)

	// Test writing roles
	roles := map[string]domain.Role{
		"code-reviewer": {
			Name:        "code-reviewer",
			Description: "Code review expert",
			File:        "~/roles/code-reviewer.md",
			Prompt:      "{file_contents}\n\nFocus on security.",
		},
		"inline": {
			Name:   "inline",
			Prompt: "You are a helpful assistant.",
		},
	}

	dir := "/test"
	err := helper.WriteRolesFile(dir, roles)
	if err != nil {
		t.Fatalf("WriteRolesFile failed: %v", err)
	}

	// Verify file was created
	path := dir + "/roles.toml"
	if !fs.Exists(path) {
		t.Fatal("roles.toml was not created")
	}

	// Test reading roles back
	readRoles, err := helper.ReadRolesFile(dir)
	if err != nil {
		t.Fatalf("ReadRolesFile failed: %v", err)
	}

	if len(readRoles) != 2 {
		t.Fatalf("expected 2 roles, got %d", len(readRoles))
	}

	reviewer, ok := readRoles["code-reviewer"]
	if !ok {
		t.Fatal("code-reviewer role not found")
	}

	if reviewer.Description != "Code review expert" {
		t.Errorf("expected description=Code review expert, got %s", reviewer.Description)
	}

	if reviewer.File != "~/roles/code-reviewer.md" {
		t.Errorf("expected file=~/roles/code-reviewer.md, got %s", reviewer.File)
	}

	inline, ok := readRoles["inline"]
	if !ok {
		t.Fatal("inline role not found")
	}

	if inline.Prompt != "You are a helpful assistant." {
		t.Errorf("expected prompt=You are a helpful assistant., got %s", inline.Prompt)
	}
}

func TestTOMLHelper_ReadRolesFile_Empty(t *testing.T) {
	fs := mocks.NewMockFileSystem()
	helper := NewTOMLHelper(fs)

	// Read from non-existent file should return empty map
	roles, err := helper.ReadRolesFile("/nonexistent")
	if err != nil {
		t.Fatalf("expected no error for non-existent file, got: %v", err)
	}

	if len(roles) != 0 {
		t.Errorf("expected empty map, got %d roles", len(roles))
	}
}

func TestTOMLHelper_ReadWriteContexts(t *testing.T) {
	fs := mocks.NewMockFileSystem()
	helper := NewTOMLHelper(fs)

	// Test writing contexts
	contexts := map[string]domain.Context{
		"project-info": {
			Name:        "project-info",
			Description: "Project context information",
			File:        "PROJECT.md",
			Prompt:      "Project info:\n{file_contents}",
			Required:    true,
		},
		"environment": {
			Name:    "environment",
			Command: "env | grep PATH",
			Prompt:  "Environment:\n{command_output}",
		},
	}

	dir := "/test"
	err := helper.WriteContextsFile(dir, contexts)
	if err != nil {
		t.Fatalf("WriteContextsFile failed: %v", err)
	}

	// Verify file was created
	path := dir + "/contexts.toml"
	if !fs.Exists(path) {
		t.Fatal("contexts.toml was not created")
	}

	// Test reading contexts back
	readContexts, err := helper.ReadContextsFile(dir)
	if err != nil {
		t.Fatalf("ReadContextsFile failed: %v", err)
	}

	if len(readContexts) != 2 {
		t.Fatalf("expected 2 contexts, got %d", len(readContexts))
	}

	projectInfo, ok := readContexts["project-info"]
	if !ok {
		t.Fatal("project-info context not found")
	}

	if projectInfo.Description != "Project context information" {
		t.Errorf("expected description=Project context information, got %s", projectInfo.Description)
	}

	if projectInfo.File != "PROJECT.md" {
		t.Errorf("expected file=PROJECT.md, got %s", projectInfo.File)
	}

	if !projectInfo.Required {
		t.Error("expected Required=true, got false")
	}

	env, ok := readContexts["environment"]
	if !ok {
		t.Fatal("environment context not found")
	}

	if env.Command != "env | grep PATH" {
		t.Errorf("expected command=env | grep PATH, got %s", env.Command)
	}
}

func TestTOMLHelper_ReadContextsFile_Empty(t *testing.T) {
	fs := mocks.NewMockFileSystem()
	helper := NewTOMLHelper(fs)

	// Read from non-existent file should return empty map
	contexts, err := helper.ReadContextsFile("/nonexistent")
	if err != nil {
		t.Fatalf("expected no error for non-existent file, got: %v", err)
	}

	if len(contexts) != 0 {
		t.Errorf("expected empty map, got %d contexts", len(contexts))
	}
}

func TestTOMLHelper_ReadWriteTasks(t *testing.T) {
	fs := mocks.NewMockFileSystem()
	helper := NewTOMLHelper(fs)

	// Test writing tasks
	tasks := map[string]domain.Task{
		"code-review": {
			Name:        "code-review",
			Alias:       "cr",
			Description: "Review code changes",
			Role:        "code-reviewer",
			Agent:       "claude",
			Command:     "git diff --staged",
			Prompt:      "Review:\n{instructions}\n\nChanges:\n{command_output}",
		},
		"quick-help": {
			Name:   "quick-help",
			Alias:  "qh",
			Prompt: "Help with: {instructions}",
		},
	}

	dir := "/test"
	err := helper.WriteTasksFile(dir, tasks)
	if err != nil {
		t.Fatalf("WriteTasksFile failed: %v", err)
	}

	// Verify file was created
	path := dir + "/tasks.toml"
	if !fs.Exists(path) {
		t.Fatal("tasks.toml was not created")
	}

	// Test reading tasks back
	readTasks, err := helper.ReadTasksFile(dir)
	if err != nil {
		t.Fatalf("ReadTasksFile failed: %v", err)
	}

	if len(readTasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(readTasks))
	}

	codeReview, ok := readTasks["code-review"]
	if !ok {
		t.Fatal("code-review task not found")
	}

	if codeReview.Alias != "cr" {
		t.Errorf("expected alias=cr, got %s", codeReview.Alias)
	}

	if codeReview.Description != "Review code changes" {
		t.Errorf("expected description=Review code changes, got %s", codeReview.Description)
	}

	if codeReview.Role != "code-reviewer" {
		t.Errorf("expected role=code-reviewer, got %s", codeReview.Role)
	}

	if codeReview.Agent != "claude" {
		t.Errorf("expected agent=claude, got %s", codeReview.Agent)
	}

	if codeReview.Command != "git diff --staged" {
		t.Errorf("expected command=git diff --staged, got %s", codeReview.Command)
	}

	quickHelp, ok := readTasks["quick-help"]
	if !ok {
		t.Fatal("quick-help task not found")
	}

	if quickHelp.Alias != "qh" {
		t.Errorf("expected alias=qh, got %s", quickHelp.Alias)
	}

	if quickHelp.Prompt != "Help with: {instructions}" {
		t.Errorf("expected prompt=Help with: {instructions}, got %s", quickHelp.Prompt)
	}
}

func TestTOMLHelper_ReadTasksFile_Empty(t *testing.T) {
	fs := mocks.NewMockFileSystem()
	helper := NewTOMLHelper(fs)

	// Read from non-existent file should return empty map
	tasks, err := helper.ReadTasksFile("/nonexistent")
	if err != nil {
		t.Fatalf("expected no error for non-existent file, got: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("expected empty map, got %d tasks", len(tasks))
	}
}
