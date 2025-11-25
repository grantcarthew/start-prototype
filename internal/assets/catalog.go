package assets

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/grantcarthew/start/internal/domain"
)

// ParseCatalogIndex parses the catalog index.csv file into a slice of AssetMeta
// CSV format: type,category,name,description,tags,bin,sha,size,created,updated
func ParseCatalogIndex(data []byte) ([]domain.AssetMeta, error) {
	reader := csv.NewReader(strings.NewReader(string(data)))

	// Read header row
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Verify header format
	expectedHeader := []string{"type", "category", "name", "description", "tags", "bin", "sha", "size", "created", "updated"}
	if len(header) != len(expectedHeader) {
		return nil, fmt.Errorf("invalid CSV header: expected %d columns, got %d", len(expectedHeader), len(header))
	}

	var assets []domain.AssetMeta

	// Read data rows
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV row: %w", err)
		}

		if len(record) != 10 {
			return nil, fmt.Errorf("invalid CSV row: expected 10 columns, got %d", len(record))
		}

		// Parse size
		size, err := strconv.ParseInt(record[7], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid size value '%s': %w", record[7], err)
		}

		// Parse timestamps
		created, err := time.Parse(time.RFC3339, record[8])
		if err != nil {
			return nil, fmt.Errorf("invalid created timestamp '%s': %w", record[8], err)
		}

		updated, err := time.Parse(time.RFC3339, record[9])
		if err != nil {
			return nil, fmt.Errorf("invalid updated timestamp '%s': %w", record[9], err)
		}

		asset := domain.AssetMeta{
			Type:        record[0],
			Category:    record[1],
			Name:        record[2],
			Description: record[3],
			Tags:        record[4], // Semicolon-separated string
			Bin:         record[5],
			SHA:         record[6],
			Size:        size,
			Created:     created,
			Updated:     updated,
		}

		assets = append(assets, asset)
	}

	return assets, nil
}

// SearchAssets performs substring matching on name, description, and tags
// Returns all assets that match the query (case-insensitive)
func SearchAssets(assets []domain.AssetMeta, query string) []domain.AssetMeta {
	if query == "" {
		return assets
	}

	query = strings.ToLower(query)
	var results []domain.AssetMeta

	for _, asset := range assets {
		// Check name (case-insensitive)
		if strings.Contains(strings.ToLower(asset.Name), query) {
			results = append(results, asset)
			continue
		}

		// Check description (case-insensitive)
		if strings.Contains(strings.ToLower(asset.Description), query) {
			results = append(results, asset)
			continue
		}

		// Check tags (case-insensitive, semicolon-separated)
		tags := strings.ToLower(asset.Tags)
		if strings.Contains(tags, query) {
			results = append(results, asset)
			continue
		}
	}

	return results
}

// FilterAssetsByType returns assets matching the given type
func FilterAssetsByType(assets []domain.AssetMeta, assetType string) []domain.AssetMeta {
	var results []domain.AssetMeta
	for _, asset := range assets {
		if asset.Type == assetType {
			results = append(results, asset)
		}
	}
	return results
}

// FindAssetByName finds an exact match by name within the given asset type
func FindAssetByName(assets []domain.AssetMeta, assetType, name string) (domain.AssetMeta, bool) {
	for _, asset := range assets {
		if asset.Type == assetType && asset.Name == name {
			return asset, true
		}
	}
	return domain.AssetMeta{}, false
}
