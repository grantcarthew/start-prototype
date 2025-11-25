package adapters

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// RealGitHubClient implements the GitHubClient interface using net/http
type RealGitHubClient struct {
	client *http.Client
}

// NewRealGitHubClient creates a new GitHub client
func NewRealGitHubClient() *RealGitHubClient {
	return &RealGitHubClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchIndex downloads the catalog index.csv file from raw.githubusercontent.com
// repo format: "owner/repo" (e.g., "grantcarthew/start")
// branch: typically "main"
func (c *RealGitHubClient) FetchIndex(ctx context.Context, repo, branch string) ([]byte, error) {
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/assets/index.csv", repo, branch)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add GitHub token if available (increases rate limits for API calls, though not needed for raw URLs)
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch index: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch index: HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read index response: %w", err)
	}

	return data, nil
}

// FetchAsset downloads an asset file from raw.githubusercontent.com
// repo format: "owner/repo" (e.g., "grantcarthew/start")
// branch: typically "main"
// path: relative path in repo (e.g., "assets/tasks/git-workflow/pre-commit-review.toml")
func (c *RealGitHubClient) FetchAsset(ctx context.Context, repo, branch, path string) ([]byte, error) {
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", repo, branch, path)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add GitHub token if available
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch asset: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch asset: HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read asset response: %w", err)
	}

	return data, nil
}
