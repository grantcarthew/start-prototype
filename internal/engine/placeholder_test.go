package engine_test

import (
	"regexp"
	"testing"

	"github.com/grantcarthew/start/internal/engine"
	"github.com/grantcarthew/start/test/assert"
)

func TestPlaceholderResolver_Resolve_BasicPlaceholders(t *testing.T) {
	resolver := engine.NewPlaceholderResolver()

	template := "Run {bin} with model {model} and prompt: {prompt}"
	values := map[string]string{
		"bin":    "smith",
		"model":  "test-model",
		"prompt": "hello world",
	}

	result := resolver.Resolve(template, values)

	expected := "Run smith with model test-model and prompt: hello world"
	assert.Equal(t, expected, result)
}

func TestPlaceholderResolver_Resolve_DatePlaceholder(t *testing.T) {
	resolver := engine.NewPlaceholderResolver()

	template := "Current date: {date}"
	values := map[string]string{}

	result := resolver.Resolve(template, values)

	// Should contain ISO 8601 format date
	assert.Contains(t, result, "Current date: ")

	// Verify it's a valid ISO 8601 timestamp (YYYY-MM-DDTHH:MM:SS...)
	matched, _ := regexp.MatchString(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`, result)
	assert.True(t, matched, "Result should contain ISO 8601 timestamp")
}

func TestPlaceholderResolver_Resolve_MissingPlaceholder(t *testing.T) {
	resolver := engine.NewPlaceholderResolver()

	template := "Model: {model}, Missing: {missing}"
	values := map[string]string{
		"model": "test-model",
	}

	result := resolver.Resolve(template, values)

	// Missing placeholder should remain as-is
	assert.Contains(t, result, "Model: test-model")
	assert.Contains(t, result, "Missing: {missing}")
}

func TestPlaceholderResolver_Resolve_EmptyTemplate(t *testing.T) {
	resolver := engine.NewPlaceholderResolver()

	template := ""
	values := map[string]string{
		"model": "test-model",
	}

	result := resolver.Resolve(template, values)

	assert.Equal(t, "", result)
}

func TestPlaceholderResolver_Resolve_EmptyValues(t *testing.T) {
	resolver := engine.NewPlaceholderResolver()

	template := "Date: {date}"
	values := map[string]string{}

	result := resolver.Resolve(template, values)

	// Should still replace {date}
	assert.Contains(t, result, "Date: ")
	assert.NotContains(t, result, "{date}")
}

func TestPlaceholderResolver_Resolve_MultipleOccurrences(t *testing.T) {
	resolver := engine.NewPlaceholderResolver()

	template := "{bin} {bin} {model}"
	values := map[string]string{
		"bin":   "smith",
		"model": "test",
	}

	result := resolver.Resolve(template, values)

	assert.Equal(t, "smith smith test", result)
}
