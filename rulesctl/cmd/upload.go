package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/choigawoon/rulesctl/internal/fileutils"
	"github.com/choigawoon/rulesctl/internal/gist"
	"github.com/choigawoon/rulesctl/pkg/config"
	"github.com/spf13/cobra"
)

var (
	forceUpload bool
	preview     bool
)

var uploadCmd = &cobra.Command{
	Use:   "upload [name]",
	Short: "Upload local rules to GIST",
	Long: `Upload rule files from local .cursor/rules directory to GIST.
The rule set name should be enclosed in quotes.

Use --preview flag to preview metadata without actual upload.`,
	Args: cobra.ExactArgs(1),
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

		if len(args) == 0 {
			cmd.SilenceUsage = true
			return fmt.Errorf("please specify a title")
		}

		title := args[0]

		// Check and create rules directory
		if err := fileutils.EnsureRulesDir(); err != nil {
			return fmt.Errorf("failed to create rules directory: %v", err)
		}

		// Preview metadata
		meta, err := gist.PreviewMetadataFromWorkingDir()
		if err != nil {
			// Show guidance if no files found
			rulesDir, _ := fileutils.GetRulesDirPath()
			fmt.Printf("Note: %v\n", err)
			fmt.Printf("Current rules directory: %s\n", rulesDir)
			fmt.Printf("You can add rule files with these commands:\n")
			fmt.Printf("  mkdir -p %s/python\n", rulesDir)
			fmt.Printf("  echo \"Python linting rules\" > %s/python/linting.mdc\n", rulesDir)
			return fmt.Errorf("no rule files to upload")
		}

		// Handle case with no files
		if len(meta.Files) == 0 {
			rulesDir, _ := fileutils.GetRulesDirPath()
			fmt.Printf("Current rules directory: %s\n", rulesDir)
			fmt.Printf("You can add rule files with these commands:\n")
			fmt.Printf("  mkdir -p %s/python\n", rulesDir)
			fmt.Printf("  echo \"Python linting rules\" > %s/python/linting.mdc\n", rulesDir)
			return fmt.Errorf("no rule files to upload")
		}

		// Preview mode only shows metadata
		if preview {
			fmt.Printf("Files to upload (total %d):\n", len(meta.Files))
			for _, file := range meta.Files {
				fmt.Printf("  - %s\n", file.Path)
			}
			return nil
		}

		// Read file contents and create Gist file map
		files := make(map[string]gist.File)
		rulesDir, err := fileutils.GetRulesDirPath()
		if err != nil {
			return fmt.Errorf("failed to get rules directory path: %v", err)
		}

		for _, fileInfo := range meta.Files {
			fullPath := filepath.Join(rulesDir, fileInfo.Path)
			content, err := os.ReadFile(fullPath)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %v", fileInfo.Path, err)
			}

			files[fileInfo.GistName] = gist.File{
				Content: string(content),
			}
		}

		// Add metadata file
		metaContent, err := meta.ToJSON()
		if err != nil {
			return fmt.Errorf("failed to generate metadata JSON: %v", err)
		}
		files[gist.MetaFileName] = gist.File{
			Content: string(metaContent),
		}

		// Initialize Gist client
		client, err := gist.NewClient()
		if err != nil {
			return fmt.Errorf("failed to initialize Gist client: %v", err)
		}

		// Create or update Gist
		gistID, err := client.CreateOrUpdateGist(title, files, forceUpload)
		if err != nil {
			return fmt.Errorf("failed to upload Gist: %v", err)
		}

		fmt.Printf("Rules successfully uploaded. Gist ID: %s\n", gistID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().BoolVarP(&forceUpload, "force", "f", false, "Force upload when conflicts exist")
	uploadCmd.Flags().BoolVarP(&preview, "preview", "p", false, "Preview metadata before upload")
} 