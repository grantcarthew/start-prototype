package assert

import (
	"strings"
	"testing"
)

// NoError fails the test if err is not nil
func NoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Error fails the test if err is nil
func Error(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// Equal fails the test if expected != actual
func Equal[T comparable](t *testing.T, expected, actual T) {
	t.Helper()
	if expected != actual {
		t.Errorf("got %v, want %v", actual, expected)
	}
}

// NotEqual fails the test if expected == actual
func NotEqual[T comparable](t *testing.T, expected, actual T) {
	t.Helper()
	if expected == actual {
		t.Errorf("got %v, want different value", actual)
	}
}

// Contains fails the test if s does not contain substr
func Contains(t *testing.T, s, substr string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Errorf("string %q does not contain %q", s, substr)
	}
}

// NotContains fails the test if s contains substr
func NotContains(t *testing.T, s, substr string) {
	t.Helper()
	if strings.Contains(s, substr) {
		t.Errorf("string %q should not contain %q", s, substr)
	}
}

// True fails the test if condition is false
func True(t *testing.T, condition bool, msg string) {
	t.Helper()
	if !condition {
		t.Errorf("condition is false: %s", msg)
	}
}

// False fails the test if condition is true
func False(t *testing.T, condition bool, msg string) {
	t.Helper()
	if condition {
		t.Errorf("condition is true: %s", msg)
	}
}

// Len fails the test if len(v) != expected
func Len(t *testing.T, expected int, v interface{}) {
	t.Helper()
	var actual int
	switch val := v.(type) {
	case string:
		actual = len(val)
	case []interface{}:
		actual = len(val)
	default:
		// Use reflection for other slice types
		t.Fatal("Len() called with unsupported type")
		return
	}
	if actual != expected {
		t.Errorf("length is %d, want %d", actual, expected)
	}
}
