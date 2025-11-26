package version

import (
	"testing"
)

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name           string
		current        string
		latest         string
		expectedStatus string
	}{
		{
			name:           "exact match",
			current:        "1.2.3",
			latest:         "1.2.3",
			expectedStatus: "Up to date",
		},
		{
			name:           "exact match with v prefix",
			current:        "v1.2.3",
			latest:         "v1.2.3",
			expectedStatus: "Up to date",
		},
		{
			name:           "update available",
			current:        "1.2.3",
			latest:         "1.3.0",
			expectedStatus: "Update available",
		},
		{
			name:           "ahead of release",
			current:        "1.4.0",
			latest:         "1.3.0",
			expectedStatus: "Ahead of latest release",
		},
		{
			name:           "development build up to date",
			current:        "1.2.3-5-gabc1234",
			latest:         "1.2.3",
			expectedStatus: "Up to date",
		},
		{
			name:           "development build needs update",
			current:        "1.2.3-5-gabc1234",
			latest:         "1.3.0",
			expectedStatus: "Update available",
		},
		{
			name:           "development build ahead",
			current:        "1.4.0-dirty",
			latest:         "1.3.0",
			expectedStatus: "Ahead of latest release",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, _ := CompareVersions(tt.current, tt.latest)
			if status != tt.expectedStatus {
				t.Errorf("CompareVersions(%q, %q) status = %q, want %q",
					tt.current, tt.latest, status, tt.expectedStatus)
			}
		})
	}
}

func TestCompareVersions_Messages(t *testing.T) {
	tests := []struct {
		name            string
		current         string
		latest          string
		expectedMessage string
	}{
		{
			name:            "update available message",
			current:         "1.2.3",
			latest:          "1.3.0",
			expectedMessage: "v1.2.3 → v1.3.0",
		},
		{
			name:            "development build message",
			current:         "1.2.3-5-gabc1234",
			latest:          "1.2.3",
			expectedMessage: "v1.2.3 (development build, base version matches release)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, message := CompareVersions(tt.current, tt.latest)
			if message != tt.expectedMessage {
				t.Errorf("CompareVersions(%q, %q) message = %q, want %q",
					tt.current, tt.latest, message, tt.expectedMessage)
			}
		})
	}
}

func TestDetectInstallMethod(t *testing.T) {
	// Just test that it returns something
	method := DetectInstallMethod()
	if method == "" {
		t.Error("DetectInstallMethod() returned empty string")
	}
}
