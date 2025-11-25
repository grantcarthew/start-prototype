package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/grantcarthew/start/internal/assets"
	"github.com/spf13/cobra"
)

// AssetsCommand wraps all asset-related commands
type AssetsCommand struct {
	resolver *assets.Resolver
}

// NewAssetsCommand creates the assets command and its subcommands
func NewAssetsCommand(resolver *assets.Resolver) *cobra.Command {
	ac := &AssetsCommand{
		resolver: resolver,
	}

	cmd := &cobra.Command{
		Use:   "assets",
		Short: "Manage asset catalog",
		Long:  "Browse, search, and install assets from the catalog",
	}

	// Add subcommands
	cmd.AddCommand(ac.newSearchCommand())
	cmd.AddCommand(ac.newAddCommand())

	return cmd
}

// newSearchCommand creates the 'start assets search' command
func (ac *AssetsCommand) newSearchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search the asset catalog",
		Long:  "Search assets by name, description, or tags",
		Args:  cobra.MaximumNArgs(1),
		RunE:  ac.runSearch,
	}

	return cmd
}

// newAddCommand creates the 'start assets add' command
func (ac *AssetsCommand) newAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <type> <name>",
		Short: "Download and cache an asset",
		Long:  "Download an asset from the catalog and add it to the cache",
		Args:  cobra.ExactArgs(2),
		RunE:  ac.runAdd,
	}

	return cmd
}

// runSearch executes the search command
func (ac *AssetsCommand) runSearch(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Get query from args (empty string if no args)
	query := ""
	if len(args) > 0 {
		query = args[0]
	}

	// Get repo from config (use default if not set)
	repo := os.Getenv("ASSET_REPO")
	if repo == "" {
		repo = "grantcarthew/start"
	}

	fmt.Println("Searching catalog...")

	// Search catalog
	results, err := ac.resolver.SearchCatalog(ctx, query, repo)
	if err != nil {
		return fmt.Errorf("failed to search catalog: %w", err)
	}

	// Display results
	if len(results) == 0 {
		fmt.Println("No assets found")
		return nil
	}

	fmt.Printf("\nFound %d asset(s):\n\n", len(results))

	for _, asset := range results {
		// Format: type/category/name
		fullPath := fmt.Sprintf("%s/%s/%s", asset.Type, asset.Category, asset.Name)
		fmt.Printf("%s\n", fullPath)
		fmt.Printf("  Description: %s\n", asset.Description)

		// Display tags if present
		if asset.Tags != "" {
			tags := strings.ReplaceAll(asset.Tags, ";", ", ")
			fmt.Printf("  Tags: %s\n", tags)
		}

		fmt.Println()
	}

	return nil
}

// runAdd executes the add command
func (ac *AssetsCommand) runAdd(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	assetType := args[0]
	name := args[1]

	// Validate asset type
	validTypes := []string{"tasks", "roles", "agents", "contexts"}
	isValid := false
	for _, t := range validTypes {
		if assetType == t {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("invalid asset type '%s': must be one of %v", assetType, validTypes)
	}

	// Get repo from config (use default if not set)
	repo := os.Getenv("ASSET_REPO")
	if repo == "" {
		repo = "grantcarthew/start"
	}

	fmt.Printf("Downloading %s/%s from catalog...\n", assetType, name)

	// Download and cache asset
	if err := ac.resolver.DownloadAsset(ctx, assetType, name, repo); err != nil {
		return fmt.Errorf("failed to download asset: %w", err)
	}

	fmt.Printf("✓ Downloaded and cached %s/%s\n", assetType, name)
	fmt.Printf("\nAsset cached to ~/.config/start/assets/%s/\n", assetType)
	fmt.Printf("\nTo use this asset, add it to your configuration:\n")
	fmt.Printf("  start config %s add %s\n", strings.TrimSuffix(assetType, "s"), name)

	return nil
}
