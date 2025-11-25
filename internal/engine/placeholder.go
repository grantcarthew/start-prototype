package engine

import (
	"strings"
	"time"
)

// PlaceholderResolver resolves placeholders in templates
type PlaceholderResolver struct{}

// NewPlaceholderResolver creates a new placeholder resolver
func NewPlaceholderResolver() *PlaceholderResolver {
	return &PlaceholderResolver{}
}

// Resolve replaces placeholders in a template with provided values
// Supports: {bin}, {model}, {prompt}, {date}
func (r *PlaceholderResolver) Resolve(template string, values map[string]string) string {
	result := template

	// Replace all provided values
	for key, value := range values {
		placeholder := "{" + key + "}"
		result = strings.ReplaceAll(result, placeholder, value)
	}

	// Always replace {date} with current timestamp
	result = strings.ReplaceAll(result, "{date}", r.getCurrentTimestamp())

	return result
}

// getCurrentTimestamp returns the current timestamp in ISO 8601 format with timezone
func (r *PlaceholderResolver) getCurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}
