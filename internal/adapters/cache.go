package adapters

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/grantcarthew/start/internal/domain"
	"github.com/pelletier/go-toml/v2"
)

// FileCache implements the Cache interface using the filesystem
// Cache structure: ~/.config/start/assets/{type}/{category}/{name}.toml
type FileCache struct {
	FS   domain.FileSystem
	Base string // e.g., ~/.config/start/assets
}

// NewFileCache creates a new file-based cache
func NewFileCache(fs domain.FileSystem, basePath string) *FileCache {
	return &FileCache{
		FS:   fs,
		Base: basePath,
	}
}

// Get retrieves an asset from the cache by type and name
// Returns the asset content (.toml file)
// Uses glob pattern to find asset across all categories since category is unknown
func (c *FileCache) Get(assetType, name string) ([]byte, error) {
	// Pattern: ~/.config/start/assets/{type}/*/{name}.toml
	pattern := filepath.Join(c.Base, assetType, "*", name+".toml")
	matches, err := c.FS.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search cache: %w", err)
	}

	if len(matches) == 0 {
		return nil, os.ErrNotExist
	}

	// Return first match (should only be one)
	content, err := c.FS.ReadFile(matches[0])
	if err != nil {
		return nil, fmt.Errorf("failed to read cached asset: %w", err)
	}

	return content, nil
}

// Set stores an asset in the cache with its metadata
// Writes both the asset content and .meta.toml sidecar file
func (c *FileCache) Set(assetType, name string, content []byte, meta domain.AssetMeta) error {
	// Create directory: ~/.config/start/assets/{type}/{category}/
	dir := filepath.Join(c.Base, assetType, meta.Category)
	if err := c.FS.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Write asset content: {name}.toml
	assetPath := filepath.Join(dir, name+".toml")
	if err := c.FS.WriteFile(assetPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write asset to cache: %w", err)
	}

	// Write metadata: {name}.meta.toml
	metaPath := filepath.Join(dir, name+".meta.toml")
	metaBytes, err := marshalMetadata(meta)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	if err := c.FS.WriteFile(metaPath, metaBytes, 0644); err != nil {
		return fmt.Errorf("failed to write metadata to cache: %w", err)
	}

	return nil
}

// List returns all cached assets of a given type
// Scans all category subdirectories for assets
func (c *FileCache) List(assetType string) ([]domain.CachedAsset, error) {
	// Pattern: ~/.config/start/assets/{type}/*/*.toml (excluding .meta.toml)
	pattern := filepath.Join(c.Base, assetType, "*", "*.toml")
	matches, err := c.FS.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to list cached assets: %w", err)
	}

	var assets []domain.CachedAsset
	for _, path := range matches {
		// Skip .meta.toml files
		if strings.HasSuffix(path, ".meta.toml") {
			continue
		}

		// Extract name and category from path
		// Path format: {base}/{type}/{category}/{name}.toml
		relPath := strings.TrimPrefix(path, c.Base+string(filepath.Separator))
		parts := strings.Split(relPath, string(filepath.Separator))
		if len(parts) < 3 {
			continue
		}

		category := parts[1]
		nameWithExt := parts[2]
		name := strings.TrimSuffix(nameWithExt, ".toml")

		// Try to read metadata
		meta, err := c.readMetadata(assetType, category, name)
		if err != nil {
			// If metadata missing, create minimal entry
			meta = domain.AssetMeta{
				Type:     assetType,
				Category: category,
				Name:     name,
			}
		}

		assets = append(assets, domain.CachedAsset{
			Type:     assetType,
			Category: category,
			Name:     name,
			Meta:     meta,
		})
	}

	return assets, nil
}

// Delete removes an asset and its metadata from the cache
func (c *FileCache) Delete(assetType, name string) error {
	// Find the asset using glob (category unknown)
	pattern := filepath.Join(c.Base, assetType, "*", name+".toml")
	matches, err := c.FS.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to search cache: %w", err)
	}

	if len(matches) == 0 {
		return os.ErrNotExist
	}

	// Delete asset file and metadata file
	assetPath := matches[0]
	metaPath := strings.TrimSuffix(assetPath, ".toml") + ".meta.toml"

	if err := c.FS.Remove(assetPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete asset: %w", err)
	}

	if err := c.FS.Remove(metaPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete metadata: %w", err)
	}

	return nil
}

// readMetadata reads the .meta.toml file for an asset
func (c *FileCache) readMetadata(assetType, category, name string) (domain.AssetMeta, error) {
	metaPath := filepath.Join(c.Base, assetType, category, name+".meta.toml")
	content, err := c.FS.ReadFile(metaPath)
	if err != nil {
		return domain.AssetMeta{}, err
	}

	var wrapper struct {
		Metadata domain.AssetMeta `toml:"metadata"`
	}

	if err := toml.Unmarshal(content, &wrapper); err != nil {
		return domain.AssetMeta{}, fmt.Errorf("failed to parse metadata: %w", err)
	}

	// Fill in derived fields
	wrapper.Metadata.Type = assetType
	wrapper.Metadata.Category = category
	wrapper.Metadata.Name = name

	return wrapper.Metadata, nil
}

// marshalMetadata converts AssetMeta to TOML bytes
func marshalMetadata(meta domain.AssetMeta) ([]byte, error) {
	wrapper := struct {
		Metadata domain.AssetMeta `toml:"metadata"`
	}{
		Metadata: meta,
	}

	bytes, err := toml.Marshal(wrapper)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
