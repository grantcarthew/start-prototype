package engine

import (
	"testing"

	"github.com/grantcarthew/start/test/mocks"
)

func TestUTDProcessor_FileOnly(t *testing.T) {
	fs := mocks.NewMockFileSystem()
	fs.Files["/test/role.md"] = "Role content"

	cmdRunner := mocks.NewMockCommandRunner()
	processor := NewUTDProcessor(fs, cmdRunner, "/test")

	input := UTDInput{
		File: "role.md",
	}

	result := processor.Process(input, "bash", 30)

	if result.Skipped {
		t.Errorf("Expected not skipped, got skipped with warnings: %v", result.Warnings)
	}
	if result.Content != "Role content" {
		t.Errorf("Expected 'Role content', got %q", result.Content)
	}
	if result.FilePath != "/test/role.md" {
		t.Errorf("Expected '/test/role.md', got %q", result.FilePath)
	}
}

func TestUTDProcessor_CommandOnly(t *testing.T) {
	fs := mocks.NewMockFileSystem()
	cmdRunner := mocks.NewMockCommandRunner()
	cmdRunner.SetOutput("command output", nil)

	processor := NewUTDProcessor(fs, cmdRunner, "/test")

	input := UTDInput{
		Command: "git status",
	}

	result := processor.Process(input, "bash", 30)

	if result.Skipped {
		t.Errorf("Expected not skipped, got skipped")
	}
	if result.Content != "command output" {
		t.Errorf("Expected 'command output', got %q", result.Content)
	}
}

func TestUTDProcessor_PromptOnly(t *testing.T) {
	fs := mocks.NewMockFileSystem()
	cmdRunner := mocks.NewMockCommandRunner()
	processor := NewUTDProcessor(fs, cmdRunner, "/test")

	input := UTDInput{
		Prompt: "Static prompt text",
	}

	result := processor.Process(input, "bash", 30)

	if result.Skipped {
		t.Errorf("Expected not skipped, got skipped")
	}
	if result.Content != "Static prompt text" {
		t.Errorf("Expected 'Static prompt text', got %q", result.Content)
	}
}

func TestUTDProcessor_FileWithPrompt(t *testing.T) {
	fs := mocks.NewMockFileSystem()
	fs.Files["/test/role.md"] = "Role content"

	cmdRunner := mocks.NewMockCommandRunner()
	processor := NewUTDProcessor(fs, cmdRunner, "/test")

	input := UTDInput{
		File:   "role.md",
		Prompt: "Read {file_contents} for context.",
	}

	result := processor.Process(input, "bash", 30)

	if result.Skipped {
		t.Errorf("Expected not skipped, got skipped")
	}
	expected := "Read Role content for context."
	if result.Content != expected {
		t.Errorf("Expected %q, got %q", expected, result.Content)
	}
}

func TestUTDProcessor_CommandWithPrompt(t *testing.T) {
	fs := mocks.NewMockFileSystem()
	cmdRunner := mocks.NewMockCommandRunner()
	cmdRunner.SetOutput("git output", nil)

	processor := NewUTDProcessor(fs, cmdRunner, "/test")

	input := UTDInput{
		Command: "git status",
		Prompt:  "Status:\n{command_output}",
	}

	result := processor.Process(input, "bash", 30)

	if result.Skipped {
		t.Errorf("Expected not skipped, got skipped")
	}
	expected := "Status:\ngit output"
	if result.Content != expected {
		t.Errorf("Expected %q, got %q", expected, result.Content)
	}
}

func TestUTDProcessor_AllThreeFields(t *testing.T) {
	fs := mocks.NewMockFileSystem()
	fs.Files["/test/doc.md"] = "Doc content"

	cmdRunner := mocks.NewMockCommandRunner()
	cmdRunner.SetOutput("cmd output", nil)

	processor := NewUTDProcessor(fs, cmdRunner, "/test")

	input := UTDInput{
		File:    "doc.md",
		Command: "git log",
		Prompt:  "Docs:\n{file_contents}\n\nLog:\n{command_output}",
	}

	result := processor.Process(input, "bash", 30)

	if result.Skipped {
		t.Errorf("Expected not skipped, got skipped")
	}
	expected := "Docs:\nDoc content\n\nLog:\ncmd output"
	if result.Content != expected {
		t.Errorf("Expected %q, got %q", expected, result.Content)
	}
}

func TestUTDProcessor_MissingFile(t *testing.T) {
	fs := mocks.NewMockFileSystem()
	cmdRunner := mocks.NewMockCommandRunner()
	processor := NewUTDProcessor(fs, cmdRunner, "/test")

	input := UTDInput{
		File:   "missing.md",
		Prompt: "Read {file_contents}",
	}

	result := processor.Process(input, "bash", 30)

	if !result.Skipped {
		t.Errorf("Expected skipped due to missing file")
	}
	if len(result.Warnings) == 0 {
		t.Errorf("Expected warnings about missing file")
	}
}

func TestUTDProcessor_EmptySection(t *testing.T) {
	fs := mocks.NewMockFileSystem()
	cmdRunner := mocks.NewMockCommandRunner()
	processor := NewUTDProcessor(fs, cmdRunner, "/test")

	input := UTDInput{}

	result := processor.Process(input, "bash", 30)

	if !result.Skipped {
		t.Errorf("Expected skipped for empty section")
	}
	if len(result.Warnings) == 0 {
		t.Errorf("Expected warnings about empty section")
	}
}

func TestUTDProcessor_TildeExpansion(t *testing.T) {
	fs := mocks.NewMockFileSystem()
	// Mock will handle tilde expansion in resolvePath
	cmdRunner := mocks.NewMockCommandRunner()
	processor := NewUTDProcessor(fs, cmdRunner, "/test")

	input := UTDInput{
		File: "~/role.md",
	}

	result := processor.Process(input, "bash", 30)

	// FilePath should have tilde expanded (mock doesn't actually expand, but we test the logic)
	if result.FilePath == "~/role.md" {
		t.Errorf("Expected tilde to be expanded, got %q", result.FilePath)
	}
}
