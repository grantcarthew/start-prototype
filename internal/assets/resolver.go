package assets

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/grantcarthew/start/internal/config"
	"github.com/grantcarthew/start/internal/domain"
	"github.com/pelletier/go-toml/v2"
)

// Resolver implements the asset resolution algorithm from DR-033
// Resolution order: local config → global config → cache → GitHub catalog
type Resolver struct {
	fs           domain.FileSystem
	cache        domain.Cache
	github       domain.GitHubClient
	configLoader *config.Loader
}

// NewResolver creates a new asset resolver
func NewResolver(fs domain.FileSystem, cache domain.Cache, github domain.GitHubClient, configLoader *config.Loader) *Resolver {
	return &Resolver{
		fs:           fs,
		cache:        cache,
		github:       github,
		configLoader: configLoader,
	}
}

// ResolveTask resolves a task by name following the resolution algorithm
// Returns the task definition and whether it was found
func (r *Resolver) ResolveTask(ctx context.Context, name string, cfg domain.Config, downloadAllowed bool) (domain.Task, bool, error) {
	// 1. Check local config
	if task, found := cfg.Tasks[name]; found {
		return task, true, nil
	}

	// 2. Check global config (already merged into cfg by config.Loader)
	// Already checked above since cfg contains merged local+global

	// 3. Check asset cache
	cacheData, err := r.cache.Get("tasks", name)
	if err == nil {
		// Found in cache, parse and return
		task, err := r.parseTaskFromTOML(name, cacheData)
		if err != nil {
			return domain.Task{}, false, fmt.Errorf("failed to parse cached task: %w", err)
		}
		return task, true, nil
	} else if !os.IsNotExist(err) {
		return domain.Task{}, false, fmt.Errorf("failed to check cache: %w", err)
	}

	// 4. Check if downloads allowed
	if !downloadAllowed {
		return domain.Task{}, false, nil
	}

	// 5. Query GitHub catalog
	return r.downloadTaskFromCatalog(ctx, name, cfg.Settings.AssetRepo)
}

// downloadTaskFromCatalog downloads a task from the GitHub catalog
func (r *Resolver) downloadTaskFromCatalog(ctx context.Context, name, repo string) (domain.Task, bool, error) {
	// Default repo if not set
	if repo == "" {
		repo = "grantcarthew/start"
	}

	// Fetch catalog index
	indexData, err := r.github.FetchIndex(ctx, repo, "main")
	if err != nil {
		return domain.Task{}, false, fmt.Errorf("failed to fetch catalog index: %w", err)
	}

	// Parse index
	assets, err := ParseCatalogIndex(indexData)
	if err != nil {
		return domain.Task{}, false, fmt.Errorf("failed to parse catalog index: %w", err)
	}

	// Find task in index
	meta, found := FindAssetByName(assets, "tasks", name)
	if !found {
		return domain.Task{}, false, nil
	}

	// Download task asset
	assetPath := fmt.Sprintf("assets/tasks/%s/%s.toml", meta.Category, name)
	assetData, err := r.github.FetchAsset(ctx, repo, "main", assetPath)
	if err != nil {
		return domain.Task{}, false, fmt.Errorf("failed to download task: %w", err)
	}

	// Cache the asset
	if err := r.cache.Set("tasks", name, assetData, meta); err != nil {
		// Log error but don't fail (caching is optional)
		// In a real implementation, we'd use a logger here
		_ = err
	}

	// Parse task
	task, err := r.parseTaskFromTOML(name, assetData)
	if err != nil {
		return domain.Task{}, false, fmt.Errorf("failed to parse downloaded task: %w", err)
	}

	return task, true, nil
}

// parseTaskFromTOML parses a task from TOML bytes
func (r *Resolver) parseTaskFromTOML(name string, data []byte) (domain.Task, error) {
	var wrapper struct {
		Task domain.Task `toml:"task"`
	}

	if err := toml.Unmarshal(data, &wrapper); err != nil {
		return domain.Task{}, err
	}

	wrapper.Task.Name = name
	return wrapper.Task, nil
}

// DownloadAsset downloads an asset from the catalog and caches it
// Used by the "start assets add" command
func (r *Resolver) DownloadAsset(ctx context.Context, assetType, name, repo string) error {
	// Default repo if not set
	if repo == "" {
		repo = "grantcarthew/start"
	}

	// Fetch catalog index
	indexData, err := r.github.FetchIndex(ctx, repo, "main")
	if err != nil {
		return fmt.Errorf("failed to fetch catalog index: %w", err)
	}

	// Parse index
	assets, err := ParseCatalogIndex(indexData)
	if err != nil {
		return fmt.Errorf("failed to parse catalog index: %w", err)
	}

	// Find asset in index
	meta, found := FindAssetByName(assets, assetType, name)
	if !found {
		return fmt.Errorf("asset not found in catalog: %s/%s", assetType, name)
	}

	// Download asset content
	assetPath := fmt.Sprintf("assets/%s/%s/%s.toml", assetType, meta.Category, name)
	assetData, err := r.github.FetchAsset(ctx, repo, "main", assetPath)
	if err != nil {
		return fmt.Errorf("failed to download asset: %w", err)
	}

	// Download additional files if they exist (e.g., .md files for tasks/roles)
	if assetType == "tasks" || assetType == "roles" {
		mdPath := fmt.Sprintf("assets/%s/%s/%s.md", assetType, meta.Category, name)
		mdData, err := r.github.FetchAsset(ctx, repo, "main", mdPath)
		if err == nil {
			// MD file exists, cache it too
			mdCachePath := filepath.Join(r.getCachePath(assetType, meta.Category), name+".md")
			if err := r.fs.WriteFile(mdCachePath, mdData, 0644); err != nil {
				// Don't fail if MD file can't be cached
				_ = err
			}
		}
	}

	// Cache the asset
	if err := r.cache.Set(assetType, name, assetData, meta); err != nil {
		return fmt.Errorf("failed to cache asset: %w", err)
	}

	return nil
}

// SearchCatalog searches the GitHub catalog for assets matching the query
func (r *Resolver) SearchCatalog(ctx context.Context, query, repo string) ([]domain.AssetMeta, error) {
	// Default repo if not set
	if repo == "" {
		repo = "grantcarthew/start"
	}

	// Fetch catalog index
	indexData, err := r.github.FetchIndex(ctx, repo, "main")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch catalog index: %w", err)
	}

	// Parse index
	assets, err := ParseCatalogIndex(indexData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse catalog index: %w", err)
	}

	// Search assets
	results := SearchAssets(assets, query)
	return results, nil
}

// getCachePath returns the cache path for an asset
func (r *Resolver) getCachePath(assetType, category string) string {
	// This should match the cache base path, but we need access to it
	// For now, we'll use a hardcoded path
	// In production, this should be configurable
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "start", "assets", assetType, category)
}
