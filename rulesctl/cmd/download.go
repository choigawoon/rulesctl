package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/choigawoon/rulesctl/internal/gist"
	"github.com/choigawoon/rulesctl/pkg/config"
)

var (
	gistID string
)

var downloadCmd = &cobra.Command{
	Use:   "download [title]",
	Short: "Download rules from GIST",
	Long: `Download rules from GIST to .cursor/rules directory using title or Gist ID.
Use --force option to overwrite existing files.

Examples:
  # Download by title (search in your Gists)
  rulesctl download "python-linting-rules"
  rulesctl download "python-linting-rules" --force

  # Download by Gist ID (public Gist)
  rulesctl download --gistid abc123
  rulesctl download --gistid abc123 --force`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		if cfg.Token == "" {
			cmd.SilenceUsage = true
			return fmt.Errorf("GitHub token not set. Please run 'rulesctl auth' to set your token")
		}

		var targetGistID string

		if gistID != "" {
			// Download by Gist ID
			targetGistID = gistID
		} else {
			// Download by title
			if len(args) == 0 {
				cmd.SilenceUsage = true
				return fmt.Errorf("please specify a title or use --gistid option")
			}
			title := args[0]

			// Fetch all Gists
			gists, err := gist.FetchUserGists(nil)
			if err != nil {
				cmd.SilenceUsage = true
				return fmt.Errorf("failed to fetch Gist list: %w", err)
			}

			// Find Gist by title
			found := false
			for _, g := range gists {
				if g.Description == title {
					targetGistID = g.ID
					found = true
					break
				}
			}

			if !found {
				cmd.SilenceUsage = true
				return fmt.Errorf("no Gist found with title: %s", title)
			}
		}

		// Fetch Gist
		g, err := gist.FetchGist(cfg.Token, targetGistID)
		if err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("failed to fetch Gist: %w", err)
		}

		// Check .rulesctl.meta.json file
		metaFile, exists := g.Files[gist.MetaFileName]
		if !exists {
			cmd.SilenceUsage = true
			return fmt.Errorf("this Gist is not managed by rulesctl (no metadata file)")
		}

		// Parse metadata
		meta, err := gist.ParseMetadataFromGist(metaFile.Content)
		if err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("failed to parse metadata: %w", err)
		}

		// Check for file conflicts
		if !force {
			conflicts, err := gist.CheckConflicts(meta)
			if err != nil {
				cmd.SilenceUsage = true
				return fmt.Errorf("failed to check conflicts: %w", err)
			}
			if len(conflicts) > 0 {
				fmt.Println("The following files already exist:")
				for _, path := range conflicts {
					fmt.Printf("  - %s\n", path)
				}
				cmd.SilenceUsage = true
				return fmt.Errorf("file conflicts detected. Use --force option to overwrite")
			}
		}

		// Download files
		fmt.Printf("Downloading rules... (Gist ID: %s)\n", targetGistID)
		if err := gist.DownloadFiles(cfg.Token, targetGistID, meta, force); err != nil {
			return fmt.Errorf("failed to download: %w", err)
		}

		fmt.Println("Download completed successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVar(&gistID, "gistid", "", "Gist ID to download")
} 