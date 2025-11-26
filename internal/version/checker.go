package version

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// ReleaseInfo contains information about a GitHub release
type ReleaseInfo struct {
	TagName     string    `json:"tag_name"`
	PublishedAt time.Time `json:"published_at"`
	HTMLURL     string    `json:"html_url"`
}

// RateLimit contains GitHub API rate limit information
type RateLimit struct {
	Remaining int `json:"remaining"`
	Limit     int `json:"limit"`
}

// RateLimitResponse contains the rate limit API response
type RateLimitResponse struct {
	Resources struct {
		Core RateLimit `json:"core"`
	} `json:"resources"`
}

// Checker handles version checking
type Checker struct {
	httpClient *http.Client
	repo       string
}

// NewChecker creates a new version checker
func NewChecker(repo string) *Checker {
	return &Checker{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		repo:       repo,
	}
}

// CheckLatestRelease queries GitHub for the latest release
func (c *Checker) CheckLatestRelease(ctx context.Context) (*ReleaseInfo, error) {
	// Check rate limit first
	canCheck, err := c.checkRateLimit(ctx)
	if err != nil {
		return nil, fmt.Errorf("rate limit check failed: %w", err)
	}
	if !canCheck {
		return nil, fmt.Errorf("rate limited (set GH_TOKEN for higher limits)")
	}

	// Query releases API
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", c.repo)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add GH_TOKEN if available
	if token := os.Getenv("GH_TOKEN"); token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var release ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &release, nil
}

// checkRateLimit checks if we have enough API quota remaining
func (c *Checker) checkRateLimit(ctx context.Context) (bool, error) {
	url := "https://api.github.com/rate_limit"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	// Add GH_TOKEN if available
	if token := os.Getenv("GH_TOKEN"); token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var rateLimit RateLimitResponse
	if err := json.NewDecoder(resp.Body).Decode(&rateLimit); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	// Allow check if we have at least 10 requests remaining
	return rateLimit.Resources.Core.Remaining >= 10, nil
}

// CompareVersions compares current and latest versions
// Returns status and message
func CompareVersions(current, latest string) (status, message string) {
	// Strip 'v' prefix
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")

	// Handle development builds (v1.2.3-5-gabc1234)
	isDev := false
	if strings.Contains(current, "-") {
		parts := strings.Split(current, "-")
		current = parts[0]
		isDev = true
	}

	// Parse versions (simple major.minor.patch)
	currentVer, err := parseVersion(current)
	if err != nil {
		return "Unknown", fmt.Sprintf("Invalid current version: %s", current)
	}

	latestVer, err := parseVersion(latest)
	if err != nil {
		return "Unknown", fmt.Sprintf("Invalid latest version: %s", latest)
	}

	// Compare
	cmp := compareVersion(currentVer, latestVer)
	if cmp < 0 {
		if isDev {
			return "Update available", fmt.Sprintf("v%s (development build) → v%s", current, latest)
		}
		return "Update available", fmt.Sprintf("v%s → v%s", current, latest)
	}

	if cmp == 0 {
		if isDev {
			return "Up to date", fmt.Sprintf("v%s (development build, base version matches release)", current)
		}
		return "Up to date", "Latest version"
	}

	// Current is ahead of latest
	return "Ahead of latest release", fmt.Sprintf("v%s > v%s (local build)", current, latest)
}

// version represents a semantic version
type version struct {
	major int
	minor int
	patch int
}

// parseVersion parses a version string like "1.2.3"
func parseVersion(s string) (version, error) {
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return version{}, fmt.Errorf("invalid version format: %s", s)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return version{}, fmt.Errorf("invalid major version: %s", parts[0])
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return version{}, fmt.Errorf("invalid minor version: %s", parts[1])
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return version{}, fmt.Errorf("invalid patch version: %s", parts[2])
	}

	return version{major: major, minor: minor, patch: patch}, nil
}

// compareVersion compares two versions
// Returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func compareVersion(v1, v2 version) int {
	if v1.major != v2.major {
		if v1.major < v2.major {
			return -1
		}
		return 1
	}

	if v1.minor != v2.minor {
		if v1.minor < v2.minor {
			return -1
		}
		return 1
	}

	if v1.patch != v2.patch {
		if v1.patch < v2.patch {
			return -1
		}
		return 1
	}

	return 0
}

// DetectInstallMethod returns the appropriate update command
func DetectInstallMethod() string {
	// Check for Homebrew
	if _, err := exec.LookPath("brew"); err == nil {
		if output, err := exec.Command("brew", "list", "--formula").Output(); err == nil {
			if strings.Contains(string(output), "grantcarthew/tap/start") {
				return "brew upgrade grantcarthew/tap/start"
			}
		}
	}

	// Check for go install
	if gopath := os.Getenv("GOPATH"); gopath != "" {
		return "go install github.com/grantcarthew/start/cmd/start@latest"
	}

	// Default fallback
	return "See https://github.com/grantcarthew/start#installation"
}
