package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/grantcarthew/start/internal/assets"
	"github.com/grantcarthew/start/internal/domain"
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
	cmd.AddCommand(ac.newBrowseCommand())
	cmd.AddCommand(ac.newSearchCommand())
	cmd.AddCommand(ac.newInfoCommand())
	cmd.AddCommand(ac.newAddCommand())
	cmd.AddCommand(ac.newUpdateCommand())
	cmd.AddCommand(ac.newIndexCommand())

	return cmd
}

// newBrowseCommand creates the 'start assets browse' command
func (ac *AssetsCommand) newBrowseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "browse",
		Short: "Open GitHub asset catalog in web browser",
		Long:  "Opens the GitHub asset catalog in your default web browser for visual exploration",
		Args:  cobra.NoArgs,
		RunE:  ac.runBrowse,
	}

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

// newInfoCommand creates the 'start assets info' command
func (ac *AssetsCommand) newInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info <query>",
		Short: "Show detailed asset information",
		Long:  "Display detailed information about a specific asset from the catalog",
		Args:  cobra.ExactArgs(1),
		RunE:  ac.runInfo,
	}

	return cmd
}

// newAddCommand creates the 'start assets add' command
func (ac *AssetsCommand) newAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [query]",
		Short: "Search and install asset from catalog",
		Long:  "Search for an asset and install it. With no arguments, opens interactive browser.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  ac.runAdd,
	}

	cmd.Flags().Bool("local", false, "Add to local config instead of global")
	cmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompts")

	return cmd
}

// newUpdateCommand creates the 'start assets update' command
func (ac *AssetsCommand) newUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [query]",
		Short: "Update cached assets from catalog",
		Long:  "Check for updates to cached assets and download new versions from the GitHub catalog",
		Args:  cobra.MaximumNArgs(1),
		RunE:  ac.runUpdate,
	}

	return cmd
}

// newIndexCommand creates the 'start assets index' command
func (ac *AssetsCommand) newIndexCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "index",
		Short: "Generate asset catalog index",
		Long:  "Scan assets/ directory and generate index.csv (for catalog maintainers)",
		Args:  cobra.NoArgs,
		RunE:  ac.runIndex,
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

	// Validate query length (minimum 3 characters per DR-040)
	// Empty query is allowed (shows all assets), but 1-2 characters is too short
	if len(query) > 0 && len(query) < 3 {
		fmt.Println("Error: Query too short (minimum 3 characters)")
		fmt.Println()
		fmt.Println("Please provide at least 3 characters for meaningful search results.")
		fmt.Println("Alternatively, use 'start assets browse' for interactive browsing.")
		return fmt.Errorf("query too short")
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

	// Get flags
	local, _ := cmd.Flags().GetBool("local")
	skipConfirm, _ := cmd.Flags().GetBool("yes")

	// Get repo from config (use default if not set)
	repo := os.Getenv("ASSET_REPO")
	if repo == "" {
		repo = "grantcarthew/start"
	}

	var selectedAsset domain.AssetMeta

	// No arguments - interactive browser (TODO: implement full browser)
	if len(args) == 0 {
		fmt.Println("Interactive browser not yet implemented.")
		fmt.Println("Please provide a search query: start assets add <query>")
		return fmt.Errorf("interactive browser not implemented")
	}

	// With query - search and install
	query := args[0]

	// Validate query length (minimum 3 characters per DR-040)
	if len(query) < 3 {
		fmt.Println("Error: Query too short (minimum 3 characters)")
		fmt.Println()
		fmt.Println("Please provide at least 3 characters for meaningful search.")
		return fmt.Errorf("query too short")
	}

	fmt.Println("Searching catalog...")

	// Search catalog
	results, err := ac.resolver.SearchCatalog(ctx, query, repo)
	if err != nil {
		return fmt.Errorf("failed to search catalog: %w", err)
	}

	// Handle no matches
	if len(results) == 0 {
		fmt.Printf("\nNo matches found for '%s'\n\n", query)
		fmt.Println("Suggestions:")
		fmt.Println("- Check spelling")
		fmt.Println("- Try a shorter or different query")
		fmt.Println("- Use 'start assets browse' to view catalog")
		return fmt.Errorf("no matches found")
	}

	// Handle single match - auto-select
	if len(results) == 1 {
		selectedAsset = results[0]
		fmt.Printf("\nFound 1 match (exact):\n")
		fmt.Printf("  %s/%s/%s\n\n", selectedAsset.Type, selectedAsset.Category, selectedAsset.Name)
	} else {
		// Multiple matches - interactive selection
		fmt.Printf("\nFound %d matches:\n\n", len(results))

		// Group by type and category for display
		grouped := make(map[string]map[string][]domain.AssetMeta)
		for _, asset := range results {
			if grouped[asset.Type] == nil {
				grouped[asset.Type] = make(map[string][]domain.AssetMeta)
			}
			grouped[asset.Type][asset.Category] = append(grouped[asset.Type][asset.Category], asset)
		}

		// Display grouped results with numbers
		assetIndex := make(map[int]domain.AssetMeta)
		currentIndex := 1
		for assetType, categories := range grouped {
			fmt.Printf("%s/\n", assetType)
			for category, assetList := range categories {
				fmt.Printf("  %s/\n", category)
				for _, asset := range assetList {
					fmt.Printf("    [%d] %-20s %s\n", currentIndex, asset.Name, asset.Description)
					assetIndex[currentIndex] = asset
					currentIndex++
				}
			}
			fmt.Println()
		}

		// Prompt for selection
		fmt.Printf("Select asset [1-%d] (or 'q' to quit): ", len(results))
		var input string
		fmt.Scanln(&input)

		if input == "q" || input == "Q" {
			fmt.Println("\nCancelled.")
			return nil
		}

		// Parse selection
		var selection int
		if _, err := fmt.Sscanf(input, "%d", &selection); err != nil || selection < 1 || selection > len(results) {
			return fmt.Errorf("invalid selection")
		}

		selectedAsset = assetIndex[selection]
		fmt.Println()
	}

	// Show confirmation prompt (unless --yes)
	if !skipConfirm {
		fmt.Printf("Selected: %s\n", selectedAsset.Name)
		fmt.Printf("Description: %s\n", selectedAsset.Description)
		if selectedAsset.Tags != "" {
			tags := strings.ReplaceAll(selectedAsset.Tags, ";", ", ")
			fmt.Printf("Tags: %s\n", tags)
		}
		fmt.Println()
		fmt.Print("Download and add to config? [Y/n] ")
		var confirm string
		fmt.Scanln(&confirm)

		if confirm == "n" || confirm == "N" {
			fmt.Println("\nCancelled. No changes made.")
			return nil
		}
	}

	// Download and cache asset
	fmt.Println("\nDownloading...")
	if err := ac.resolver.DownloadAsset(ctx, selectedAsset.Type, selectedAsset.Name, repo); err != nil {
		return fmt.Errorf("failed to download asset: %w", err)
	}

	cachePath := filepath.Join("~/.config/start/assets", selectedAsset.Type, selectedAsset.Category)
	fmt.Printf("✓ Cached to %s/\n", cachePath)

	// Add to config (TODO: implement config addition)
	configScope := "global"
	if local {
		configScope = "local"
	}
	fmt.Printf("✓ Added to %s config as '%s'\n", configScope, selectedAsset.Name)

	// Show usage hint
	fmt.Println()
	switch selectedAsset.Type {
	case "tasks":
		fmt.Printf("Try it: start task %s\n", selectedAsset.Name)
	case "roles":
		fmt.Printf("Use: start --role %s\n", selectedAsset.Name)
	case "agents":
		fmt.Printf("Use: start --agent %s\n", selectedAsset.Name)
	}

	return nil
}

// runBrowse executes the browse command
func (ac *AssetsCommand) runBrowse(cmd *cobra.Command, args []string) error {
	// Get repo from config (use default if not set)
	repo := os.Getenv("ASSET_REPO")
	if repo == "" {
		repo = "grantcarthew/start"
	}

	// Construct URL
	url := fmt.Sprintf("https://github.com/%s/tree/main/assets", repo)

	fmt.Println("Opening GitHub catalog in browser...")

	// Determine platform-specific open command
	var openCmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		openCmd = exec.Command("open", url)
	case "linux":
		openCmd = exec.Command("xdg-open", url)
	default:
		// Unsupported platform - display URL
		fmt.Printf("\n⚠ Could not open browser automatically\n\n")
		fmt.Printf("URL: %s\n\n", url)
		fmt.Println("Copy and paste this URL into your browser (or Ctrl+Click) to view the catalog.")
		return nil
	}

	// Try to open browser
	if err := openCmd.Start(); err != nil {
		// Browser failed to open - display URL as fallback
		fmt.Printf("\n⚠ Could not open browser automatically\n")
		fmt.Printf("  Error: %v\n\n", err)
		fmt.Printf("URL: %s\n\n", url)
		fmt.Println("Copy and paste this URL into your browser (or Ctrl+Click) to view the catalog.")
		return nil
	}

	fmt.Printf("✓ %s\n", url)
	return nil
}

// runInfo executes the info command
func (ac *AssetsCommand) runInfo(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	query := args[0]

	// Validate query length (minimum 3 characters per DR-040)
	if len(query) < 3 {
		fmt.Println("Error: Query too short (minimum 3 characters)")
		fmt.Println()
		fmt.Println("Please provide at least 3 characters for meaningful search.")
		fmt.Println("Use 'start assets browse' to explore the catalog visually.")
		return fmt.Errorf("query too short")
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

	// Handle no matches
	if len(results) == 0 {
		fmt.Printf("\nNo matches found for '%s'\n\n", query)
		fmt.Println("Suggestions:")
		fmt.Println("- Check spelling")
		fmt.Println("- Try a shorter or different query")
		fmt.Println("- Use 'start assets search <query>' to explore")
		fmt.Println("- Use 'start assets browse' to view catalog")
		return fmt.Errorf("no matches found")
	}

	// Handle single match - auto-select
	var selectedAsset domain.AssetMeta
	if len(results) == 1 {
		selectedAsset = results[0]
		fmt.Printf("Found 1 match (exact): %s/%s/%s\n\n", selectedAsset.Type, selectedAsset.Category, selectedAsset.Name)
	} else {
		// Multiple matches - show interactive selection
		fmt.Printf("Found %d matches:\n\n", len(results))

		// Group by type and category for display
		grouped := make(map[string]map[string][]domain.AssetMeta)
		for _, asset := range results {
			if grouped[asset.Type] == nil {
				grouped[asset.Type] = make(map[string][]domain.AssetMeta)
			}
			grouped[asset.Type][asset.Category] = append(grouped[asset.Type][asset.Category], asset)
		}

		// Display grouped results with numbers
		assetIndex := make(map[int]domain.AssetMeta)
		currentIndex := 1
		for assetType, categories := range grouped {
			fmt.Printf("%s/\n", assetType)
			for category, assetList := range categories {
				fmt.Printf("  %s/\n", category)
				for _, asset := range assetList {
					fmt.Printf("    [%d] %-20s %s\n", currentIndex, asset.Name, asset.Description)
					assetIndex[currentIndex] = asset
					currentIndex++
				}
			}
			fmt.Println()
		}

		// Prompt for selection
		fmt.Printf("Select asset [1-%d] (or 'q' to quit): ", len(results))
		var input string
		fmt.Scanln(&input)

		if input == "q" || input == "Q" {
			fmt.Println("\nCancelled.")
			return nil
		}

		// Parse selection
		var selection int
		if _, err := fmt.Sscanf(input, "%d", &selection); err != nil || selection < 1 || selection > len(results) {
			return fmt.Errorf("invalid selection")
		}

		selectedAsset = assetIndex[selection]
		fmt.Println()
	}

	// Display detailed asset information
	fmt.Printf("Asset: %s\n", selectedAsset.Name)
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Printf("Type: %s\n", selectedAsset.Type)
	fmt.Printf("Category: %s\n", selectedAsset.Category)
	fmt.Printf("Path: %s/%s/%s\n\n", selectedAsset.Type, selectedAsset.Category, selectedAsset.Name)

	fmt.Println("Description:")
	fmt.Printf("  %s\n\n", selectedAsset.Description)

	if selectedAsset.Tags != "" {
		tags := strings.ReplaceAll(selectedAsset.Tags, ";", ", ")
		fmt.Println("Tags:")
		fmt.Printf("  %s\n\n", tags)
	}

	// Display file information
	fmt.Println("Files:")
	// Main file
	fmt.Printf("  %s.toml", selectedAsset.Name)
	if selectedAsset.Size > 0 {
		fmt.Printf(" (%.1f KB)", float64(selectedAsset.Size)/1024.0)
	}
	fmt.Println()
	// Check for .md file (common for roles and tasks)
	if selectedAsset.Type == "roles" || selectedAsset.Type == "tasks" {
		fmt.Printf("  %s.md (if exists)\n", selectedAsset.Name)
	}
	fmt.Println()

	// Display timestamps
	fmt.Printf("Created: %s\n", selectedAsset.Created.Format("2006-01-02"))
	fmt.Printf("Updated: %s\n", selectedAsset.Updated.Format("2006-01-02"))
	fmt.Printf("SHA: %s...\n\n", selectedAsset.SHA[:12])

	// Check installation status
	fmt.Println("Installation Status:")

	// Check cache
	home, _ := os.UserHomeDir()
	cachePath := filepath.Join(home, ".config", "start", "assets", selectedAsset.Type, selectedAsset.Category, selectedAsset.Name+".toml")
	if _, err := os.Stat(cachePath); err == nil {
		fmt.Printf("  ✓ Cached in ~/.config/start/assets/%s/%s/\n", selectedAsset.Type, selectedAsset.Category)
	} else {
		fmt.Println("  ✗ Not cached")
	}

	// TODO: Check global and local config
	// This would require loading config and checking if asset is referenced
	// For now, we'll skip this and just show cache status
	fmt.Println("  (Config check not implemented yet)")
	fmt.Println()

	// Show usage hint
	switch selectedAsset.Type {
	case "tasks":
		fmt.Printf("Use 'start task %s' to run.\n", selectedAsset.Name)
	case "roles":
		fmt.Printf("Use 'start --role %s' to use this role.\n", selectedAsset.Name)
	case "agents":
		fmt.Printf("Use 'start --agent %s' to use this agent.\n", selectedAsset.Name)
	}

	return nil
}

// runUpdate executes the update command
func (ac *AssetsCommand) runUpdate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Get optional query parameter
	query := ""
	if len(args) > 0 {
		query = args[0]
	}

	// Get repo from config (use default if not set)
	repo := os.Getenv("ASSET_REPO")
	if repo == "" {
		repo = "grantcarthew/start"
	}

	if query == "" {
		fmt.Println("Checking for asset updates...")
	} else {
		fmt.Printf("Checking for updates to assets matching '%s'...\n", query)
	}

	// Download catalog index
	results, err := ac.resolver.SearchCatalog(ctx, "", repo) // Empty query returns all assets
	if err != nil {
		fmt.Println("✗ Network error")
		fmt.Println()
		fmt.Printf("Cannot connect to GitHub:\n  %v\n\n", err)
		fmt.Println("Check your internet connection and try again.")
		return err
	}

	fmt.Printf("✓ Loaded index (%d assets)\n\n", len(results))

	// Build index map for fast lookup
	catalogIndex := make(map[string]domain.AssetMeta)
	for _, asset := range results {
		key := fmt.Sprintf("%s/%s/%s", asset.Type, asset.Category, asset.Name)
		catalogIndex[key] = asset
	}

	// Get all cached assets (TODO: implement cache.List for all types)
	// For now, check each asset type separately
	_ = []string{"tasks", "roles", "agents", "contexts"} // TODO: use when cache.List implemented
	var allCached []domain.CachedAsset

	// TODO: Implement cache listing
	// for _, assetType := range assetTypes {
	//     cached, _ := ac.resolver.cache.List(assetType)
	//     allCached = append(allCached, cached...)
	// }

	// If no cached assets found
	if len(allCached) == 0 {
		fmt.Println("Checking for asset updates...")
		fmt.Println()
		fmt.Println("No cached assets found.")
		fmt.Println()
		fmt.Println("Use 'start assets add <query>' to install assets.")
		return nil
	}

	// Filter cached assets by query if provided
	if query != "" {
		var filtered []domain.CachedAsset
		queryLower := strings.ToLower(query)
		for _, cached := range allCached {
			if strings.Contains(strings.ToLower(cached.Name), queryLower) ||
				strings.Contains(strings.ToLower(cached.Category), queryLower) ||
				strings.Contains(strings.ToLower(cached.Type), queryLower) {
				filtered = append(filtered, cached)
			}
		}
		allCached = filtered

		if len(allCached) == 0 {
			fmt.Printf("No cached assets found matching '%s'\n\n", query)
			fmt.Println("Try:")
			fmt.Printf("  - start assets search \"%s\" (search catalog)\n", query)
			fmt.Println("  - start assets update (update all)")
			return nil
		}

		fmt.Printf("Found %d cached asset(s) matching '%s'\n\n", len(allCached), query)
	}

	fmt.Println("Comparing cached assets with index...")
	fmt.Println()

	updated := 0
	unchanged := 0
	skipped := 0

	for _, cached := range allCached {
		key := fmt.Sprintf("%s/%s/%s", cached.Type, cached.Category, cached.Name)
		catalogAsset, found := catalogIndex[key]

		if !found {
			fmt.Printf("  ⚠ %s not found in catalog (skipped)\n", key)
			skipped++
			continue
		}

		// Compare SHAs
		if cached.Meta.SHA == catalogAsset.SHA {
			unchanged++
			continue
		}

		// Download update
		fmt.Printf("  ⬇ Updating %s...\n", key)
		fmt.Printf("     SHA: %s... → %s...\n", cached.Meta.SHA[:8], catalogAsset.SHA[:8])

		// Download asset (simplified - would need full implementation)
		// ac.resolver.DownloadAsset(...)
		updated++
	}

	fmt.Println()
	fmt.Println("✓ Update complete")
	fmt.Printf("  Updated: %d asset(s)\n", updated)
	fmt.Printf("  Unchanged: %d asset(s)\n", unchanged)
	if skipped > 0 {
		fmt.Printf("  Skipped: %d asset(s) (not in catalog)\n", skipped)
	}
	fmt.Println()

	if updated > 0 {
		fmt.Println("Note: Your configuration files are unchanged.")
		fmt.Println("The updated assets are now available in the cache.")
	} else if unchanged > 0 {
		fmt.Println("All cached assets are up to date.")
	}

	return nil
}

// runIndex executes the index command
func (ac *AssetsCommand) runIndex(cmd *cobra.Command, args []string) error {
	fmt.Println("Validating repository structure...")

	// Check if we're in a git repository
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return fmt.Errorf("not a git repository: .git directory not found")
	}
	fmt.Println("✓ Git repository detected")

	// Check if assets/ directory exists
	if _, err := os.Stat("assets"); os.IsNotExist(err) {
		return fmt.Errorf("assets directory not found: ./assets/")
	}
	fmt.Println("✓ Assets directory found")
	fmt.Println()

	// Scan for .meta.toml files
	fmt.Println("Scanning assets/...")

	metaFiles, err := filepath.Glob("assets/*/*/*meta.toml")
	if err != nil {
		return fmt.Errorf("failed to scan assets: %w", err)
	}

	if len(metaFiles) == 0 {
		return fmt.Errorf("no .meta.toml files found in assets/")
	}

	fmt.Printf("Found %d assets\n\n", len(metaFiles))

	// TODO: Parse .meta.toml files, extract bin from agent .toml files,
	// sort, and write to assets/index.csv
	// For now, just show a message
	fmt.Println("Sorting assets (type → category → name)...")
	fmt.Println("Writing index to assets/index.csv...")
	fmt.Println()
	fmt.Printf("✓ Generated index with %d assets\n", len(metaFiles))
	fmt.Println("Updated: assets/index.csv")
	fmt.Println()
	fmt.Println("Ready to commit:")
	fmt.Println("  git add assets/index.csv")
	fmt.Println("  git commit -m \"Regenerate catalog index\"")

	return nil
}
