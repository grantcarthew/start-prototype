package assets

import (
	"testing"

	"github.com/grantcarthew/start/internal/domain"
)

func TestParseCatalogIndex(t *testing.T) {
	tests := []struct {
		name      string
		csvData   string
		wantCount int
		wantErr   bool
	}{
		{
			name: "valid index with multiple assets",
			csvData: `type,category,name,description,tags,bin,sha,size,created,updated
tasks,git-workflow,pre-commit-review,Review staged changes before committing,git;review;quality,,abc123def456,2048,2025-01-10T00:00:00Z,2025-01-10T12:30:00Z
agents,anthropic,claude,Anthropic Claude AI via Claude Code CLI,claude;anthropic;ai,claude,def456abc789,1024,2025-01-10T00:00:00Z,2025-01-10T00:00:00Z`,
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "empty index",
			csvData: `type,category,name,description,tags,bin,sha,size,created,updated`,
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "invalid header",
			csvData:   `wrong,header,format`,
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "invalid row column count",
			csvData: `type,category,name,description,tags,bin,sha,size,created,updated
tasks,git-workflow,pre-commit-review`,
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "invalid size",
			csvData: `type,category,name,description,tags,bin,sha,size,created,updated
tasks,git-workflow,pre-commit-review,Review staged changes,git,,abc123,not-a-number,2025-01-10T00:00:00Z,2025-01-10T00:00:00Z`,
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "invalid timestamp",
			csvData: `type,category,name,description,tags,bin,sha,size,created,updated
tasks,git-workflow,pre-commit-review,Review staged changes,git,,abc123,2048,invalid-date,2025-01-10T00:00:00Z`,
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assets, err := ParseCatalogIndex([]byte(tt.csvData))

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCatalogIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(assets) != tt.wantCount {
				t.Errorf("ParseCatalogIndex() returned %d assets, want %d", len(assets), tt.wantCount)
			}
		})
	}
}

func TestParseCatalogIndex_ValidData(t *testing.T) {
	csvData := `type,category,name,description,tags,bin,sha,size,created,updated
tasks,git-workflow,pre-commit-review,Review staged changes before committing,git;review;quality,,abc123def456,2048,2025-01-10T00:00:00Z,2025-01-10T12:30:00Z`

	assets, err := ParseCatalogIndex([]byte(csvData))
	if err != nil {
		t.Fatalf("ParseCatalogIndex() unexpected error: %v", err)
	}

	if len(assets) != 1 {
		t.Fatalf("ParseCatalogIndex() returned %d assets, want 1", len(assets))
	}

	asset := assets[0]

	if asset.Type != "tasks" {
		t.Errorf("asset.Type = %q, want %q", asset.Type, "tasks")
	}
	if asset.Category != "git-workflow" {
		t.Errorf("asset.Category = %q, want %q", asset.Category, "git-workflow")
	}
	if asset.Name != "pre-commit-review" {
		t.Errorf("asset.Name = %q, want %q", asset.Name, "pre-commit-review")
	}
	if asset.Description != "Review staged changes before committing" {
		t.Errorf("asset.Description = %q, want %q", asset.Description, "Review staged changes before committing")
	}
	if asset.Tags != "git;review;quality" {
		t.Errorf("asset.Tags = %q, want %q", asset.Tags, "git;review;quality")
	}
	if asset.SHA != "abc123def456" {
		t.Errorf("asset.SHA = %q, want %q", asset.SHA, "abc123def456")
	}
	if asset.Size != 2048 {
		t.Errorf("asset.Size = %d, want %d", asset.Size, 2048)
	}
}

func TestSearchAssets(t *testing.T) {
	assets := []domain.AssetMeta{
		{
			Type:        "tasks",
			Category:    "git-workflow",
			Name:        "pre-commit-review",
			Description: "Review staged changes before committing",
			Tags:        "git;review;quality",
		},
		{
			Type:        "tasks",
			Category:    "git-workflow",
			Name:        "commit-message",
			Description: "Generate conventional commit message",
			Tags:        "git;commit;conventional",
		},
		{
			Type:        "roles",
			Category:    "general",
			Name:        "code-reviewer",
			Description: "Expert code reviewer focusing on security",
			Tags:        "review;security;quality",
		},
	}

	tests := []struct {
		name      string
		query     string
		wantCount int
	}{
		{
			name:      "search by name",
			query:     "commit",
			wantCount: 2, // pre-commit-review and commit-message
		},
		{
			name:      "search by description",
			query:     "conventional",
			wantCount: 1, // commit-message
		},
		{
			name:      "search by tag",
			query:     "security",
			wantCount: 1, // code-reviewer
		},
		{
			name:      "search case insensitive",
			query:     "REVIEW",
			wantCount: 2, // pre-commit-review and code-reviewer
		},
		{
			name:      "no matches",
			query:     "nonexistent",
			wantCount: 0,
		},
		{
			name:      "empty query returns all",
			query:     "",
			wantCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := SearchAssets(assets, tt.query)

			if len(results) != tt.wantCount {
				t.Errorf("SearchAssets() returned %d results, want %d", len(results), tt.wantCount)
			}
		})
	}
}

func TestFilterAssetsByType(t *testing.T) {
	assets := []domain.AssetMeta{
		{Type: "tasks", Name: "task1"},
		{Type: "tasks", Name: "task2"},
		{Type: "roles", Name: "role1"},
		{Type: "agents", Name: "agent1"},
	}

	tests := []struct {
		name      string
		assetType string
		wantCount int
	}{
		{
			name:      "filter tasks",
			assetType: "tasks",
			wantCount: 2,
		},
		{
			name:      "filter roles",
			assetType: "roles",
			wantCount: 1,
		},
		{
			name:      "filter nonexistent type",
			assetType: "contexts",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := FilterAssetsByType(assets, tt.assetType)

			if len(results) != tt.wantCount {
				t.Errorf("FilterAssetsByType() returned %d results, want %d", len(results), tt.wantCount)
			}
		})
	}
}

func TestFindAssetByName(t *testing.T) {
	assets := []domain.AssetMeta{
		{Type: "tasks", Name: "task1"},
		{Type: "tasks", Name: "task2"},
		{Type: "roles", Name: "role1"},
	}

	tests := []struct {
		name      string
		assetType string
		assetName string
		wantFound bool
	}{
		{
			name:      "find existing task",
			assetType: "tasks",
			assetName: "task1",
			wantFound: true,
		},
		{
			name:      "find existing role",
			assetType: "roles",
			assetName: "role1",
			wantFound: true,
		},
		{
			name:      "wrong type",
			assetType: "roles",
			assetName: "task1",
			wantFound: false,
		},
		{
			name:      "nonexistent asset",
			assetType: "tasks",
			assetName: "nonexistent",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, found := FindAssetByName(assets, tt.assetType, tt.assetName)

			if found != tt.wantFound {
				t.Errorf("FindAssetByName() found = %v, want %v", found, tt.wantFound)
			}
		})
	}
}
