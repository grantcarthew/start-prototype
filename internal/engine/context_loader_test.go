package engine

import (
	"testing"

	"github.com/grantcarthew/start/internal/domain"
)

func TestContextLoader_LoadContexts_Interactive(t *testing.T) {
	fs := newMockFileSystem()
	fs.files["/ctx1.md"] = "Context 1"
	fs.files["/ctx2.md"] = "Context 2"
	fs.files["/ctx3.md"] = "Context 3"

	cmdRunner := &mockCommandRunner{}
	utdProcessor := NewUTDProcessor(fs, cmdRunner, "/workdir")
	loader := NewContextLoader(utdProcessor)

	contexts := map[string]domain.Context{
		"ctx1": {File: "/ctx1.md", Required: true},
		"ctx2": {File: "/ctx2.md", Required: false},
		"ctx3": {File: "/ctx3.md", Required: true},
	}

	contextOrder := []string{"ctx1", "ctx2", "ctx3"}

	// Interactive mode - should load all contexts
	results := loader.LoadContexts(contexts, contextOrder, CommandTypeInteractive, "bash", 30)

	if len(results) != 3 {
		t.Errorf("Expected 3 contexts, got %d", len(results))
	}

	// Verify order is preserved
	if results[0].Name != "ctx1" {
		t.Errorf("Expected first context to be ctx1, got %s", results[0].Name)
	}
	if results[1].Name != "ctx2" {
		t.Errorf("Expected second context to be ctx2, got %s", results[1].Name)
	}
	if results[2].Name != "ctx3" {
		t.Errorf("Expected third context to be ctx3, got %s", results[2].Name)
	}

	// Verify content
	if results[0].Content != "Context 1" {
		t.Errorf("Expected content %q, got %q", "Context 1", results[0].Content)
	}
}

func TestContextLoader_LoadContexts_Prompt(t *testing.T) {
	fs := newMockFileSystem()
	fs.files["/ctx1.md"] = "Required context"
	fs.files["/ctx2.md"] = "Optional context"

	cmdRunner := &mockCommandRunner{}
	utdProcessor := NewUTDProcessor(fs, cmdRunner, "/workdir")
	loader := NewContextLoader(utdProcessor)

	contexts := map[string]domain.Context{
		"ctx1": {File: "/ctx1.md", Required: true},
		"ctx2": {File: "/ctx2.md", Required: false},
	}

	contextOrder := []string{"ctx1", "ctx2"}

	// Prompt mode - should load only required contexts
	results := loader.LoadContexts(contexts, contextOrder, CommandTypePrompt, "bash", 30)

	if len(results) != 1 {
		t.Errorf("Expected 1 context in prompt mode, got %d", len(results))
	}

	if results[0].Name != "ctx1" {
		t.Errorf("Expected context ctx1, got %s", results[0].Name)
	}

	if results[0].Content != "Required context" {
		t.Errorf("Expected content %q, got %q", "Required context", results[0].Content)
	}
}

func TestContextLoader_LoadContexts_Task(t *testing.T) {
	fs := newMockFileSystem()
	fs.files["/ctx1.md"] = "Required context"
	fs.files["/ctx2.md"] = "Optional context"
	fs.files["/ctx3.md"] = "Another required"

	cmdRunner := &mockCommandRunner{}
	utdProcessor := NewUTDProcessor(fs, cmdRunner, "/workdir")
	loader := NewContextLoader(utdProcessor)

	contexts := map[string]domain.Context{
		"ctx1": {File: "/ctx1.md", Required: true},
		"ctx2": {File: "/ctx2.md", Required: false},
		"ctx3": {File: "/ctx3.md", Required: true},
	}

	contextOrder := []string{"ctx1", "ctx2", "ctx3"}

	// Task mode - should load only required contexts
	results := loader.LoadContexts(contexts, contextOrder, CommandTypeTask, "bash", 30)

	if len(results) != 2 {
		t.Errorf("Expected 2 contexts in task mode, got %d", len(results))
	}

	// Verify order preserved (ctx1, then ctx3)
	if results[0].Name != "ctx1" {
		t.Errorf("Expected first context to be ctx1, got %s", results[0].Name)
	}
	if results[1].Name != "ctx3" {
		t.Errorf("Expected second context to be ctx3, got %s", results[1].Name)
	}
}

func TestContextLoader_LoadContexts_EmptyList(t *testing.T) {
	fs := newMockFileSystem()
	cmdRunner := &mockCommandRunner{}
	utdProcessor := NewUTDProcessor(fs, cmdRunner, "/workdir")
	loader := NewContextLoader(utdProcessor)

	contexts := map[string]domain.Context{}
	contextOrder := []string{}

	results := loader.LoadContexts(contexts, contextOrder, CommandTypeInteractive, "bash", 30)

	if len(results) != 0 {
		t.Errorf("Expected 0 contexts, got %d", len(results))
	}
}

func TestContextLoader_LoadContexts_FileNotFound(t *testing.T) {
	fs := newMockFileSystem()
	// Don't add file to fs

	cmdRunner := &mockCommandRunner{}
	utdProcessor := NewUTDProcessor(fs, cmdRunner, "/workdir")
	loader := NewContextLoader(utdProcessor)

	contexts := map[string]domain.Context{
		"ctx1": {File: "/nonexistent.md", Required: true},
	}

	contextOrder := []string{"ctx1"}

	// Should still return result with warnings
	results := loader.LoadContexts(contexts, contextOrder, CommandTypeInteractive, "bash", 30)

	if len(results) != 1 {
		t.Errorf("Expected 1 context result, got %d", len(results))
	}

	if results[0].Name != "ctx1" {
		t.Errorf("Expected context name ctx1, got %s", results[0].Name)
	}

	// Content should be empty since file not found
	if results[0].Content != "" {
		t.Errorf("Expected empty content for missing file, got %q", results[0].Content)
	}

	// Should have warnings
	if len(results[0].Warnings) == 0 {
		t.Error("Expected warnings for missing file")
	}
}

func TestContextLoader_LoadContexts_WithPrompt(t *testing.T) {
	fs := newMockFileSystem()
	cmdRunner := &mockCommandRunner{}
	utdProcessor := NewUTDProcessor(fs, cmdRunner, "/workdir")
	loader := NewContextLoader(utdProcessor)

	contexts := map[string]domain.Context{
		"ctx1": {Prompt: "Static prompt text", Required: true},
	}

	contextOrder := []string{"ctx1"}

	results := loader.LoadContexts(contexts, contextOrder, CommandTypeInteractive, "bash", 30)

	if len(results) != 1 {
		t.Errorf("Expected 1 context, got %d", len(results))
	}

	if results[0].Content != "Static prompt text" {
		t.Errorf("Expected content %q, got %q", "Static prompt text", results[0].Content)
	}
}

func TestContextLoader_LoadContexts_WithCommand(t *testing.T) {
	fs := newMockFileSystem()
	cmdRunner := &mockCommandRunner{output: "Command output"}
	utdProcessor := NewUTDProcessor(fs, cmdRunner, "/workdir")
	loader := NewContextLoader(utdProcessor)

	contexts := map[string]domain.Context{
		"ctx1": {Command: "echo test", Required: true},
	}

	contextOrder := []string{"ctx1"}

	results := loader.LoadContexts(contexts, contextOrder, CommandTypeInteractive, "bash", 30)

	if len(results) != 1 {
		t.Errorf("Expected 1 context, got %d", len(results))
	}

	if results[0].Content != "Command output" {
		t.Errorf("Expected content %q, got %q", "Command output", results[0].Content)
	}
}

func TestContextLoader_LoadContexts_OrderPreserved(t *testing.T) {
	fs := newMockFileSystem()
	fs.files["/a.md"] = "A"
	fs.files["/b.md"] = "B"
	fs.files["/c.md"] = "C"

	cmdRunner := &mockCommandRunner{}
	utdProcessor := NewUTDProcessor(fs, cmdRunner, "/workdir")
	loader := NewContextLoader(utdProcessor)

	contexts := map[string]domain.Context{
		"ctx-a": {File: "/a.md", Required: true},
		"ctx-b": {File: "/b.md", Required: true},
		"ctx-c": {File: "/c.md", Required: true},
	}

	// Test different orders
	contextOrder := []string{"ctx-c", "ctx-a", "ctx-b"}

	results := loader.LoadContexts(contexts, contextOrder, CommandTypeInteractive, "bash", 30)

	if len(results) != 3 {
		t.Fatalf("Expected 3 contexts, got %d", len(results))
	}

	// Verify order matches contextOrder
	if results[0].Name != "ctx-c" {
		t.Errorf("Expected first context to be ctx-c, got %s", results[0].Name)
	}
	if results[1].Name != "ctx-a" {
		t.Errorf("Expected second context to be ctx-a, got %s", results[1].Name)
	}
	if results[2].Name != "ctx-b" {
		t.Errorf("Expected third context to be ctx-b, got %s", results[2].Name)
	}
}

func TestContextLoader_LoadContexts_MissingFromMap(t *testing.T) {
	fs := newMockFileSystem()
	fs.files["/ctx1.md"] = "Context 1"

	cmdRunner := &mockCommandRunner{}
	utdProcessor := NewUTDProcessor(fs, cmdRunner, "/workdir")
	loader := NewContextLoader(utdProcessor)

	contexts := map[string]domain.Context{
		"ctx1": {File: "/ctx1.md", Required: true},
	}

	// Context order includes a context not in the map
	contextOrder := []string{"ctx1", "ctx-nonexistent", "ctx2"}

	results := loader.LoadContexts(contexts, contextOrder, CommandTypeInteractive, "bash", 30)

	// Should only load ctx1 (others not in map are skipped)
	if len(results) != 1 {
		t.Errorf("Expected 1 context, got %d", len(results))
	}

	if results[0].Name != "ctx1" {
		t.Errorf("Expected context ctx1, got %s", results[0].Name)
	}
}

func TestContextLoader_LoadContexts_UTDWithFilePath(t *testing.T) {
	fs := newMockFileSystem()
	fs.files["/context.md"] = "File content"

	cmdRunner := &mockCommandRunner{}
	utdProcessor := NewUTDProcessor(fs, cmdRunner, "/workdir")
	loader := NewContextLoader(utdProcessor)

	contexts := map[string]domain.Context{
		"ctx1": {File: "/context.md", Required: true},
	}

	contextOrder := []string{"ctx1"}

	results := loader.LoadContexts(contexts, contextOrder, CommandTypeInteractive, "bash", 30)

	if len(results) != 1 {
		t.Fatalf("Expected 1 context, got %d", len(results))
	}

	// Verify file path is set
	if results[0].FilePath != "/context.md" {
		t.Errorf("Expected file path %q, got %q", "/context.md", results[0].FilePath)
	}
}
