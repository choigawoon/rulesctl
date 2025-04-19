package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/choigawoon/rulesctl/internal/gist"
	"github.com/choigawoon/rulesctl/pkg/config"
	"github.com/spf13/cobra"
)

const (
	titleWidth = 25    // Title max width
	dateWidth  = 19    // Date width
	idWidth    = 32    // Gist ID width
	revWidth   = 8     // Revision width
	typeWidth  = 8     // Type width (Public/Private)
	separator  = "..."
)

// truncateString truncates a string to the specified width and adds ellipsis if needed
func truncateString(s string, width int) string {
	if utf8.RuneCountInString(s) <= width {
		return s + strings.Repeat(" ", width-utf8.RuneCountInString(s))
	}
	return string([]rune(s)[:width-len(separator)]) + separator
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List rules stored in GIST",
	Long: `List all rules stored in GIST.
By default, outputs in [Type] [Title] [Last Modified] [Gist ID] format.
Use --detail flag to include revision information.

Examples:
  rulesctl list          # Show basic information
  rulesctl list --detail # Show with revision information`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		if config.Token == "" {
			return fmt.Errorf("GitHub token not set. Please run 'rulesctl auth' to set your token")
		}

		// Show token source
		if os.Getenv("GITHUB_TOKEN") != "" {
			fmt.Println("GitHub Token: Loaded from environment variable")
		} else {
			fmt.Println("GitHub Token: Loaded from config file")
		}

		// Only fetch Gists from the last month
		since := time.Now().AddDate(0, -1, 0)
		gists, err := gist.FetchUserGists(&since)
		if err != nil {
			return fmt.Errorf("failed to fetch Gist list: %w", err)
		}

		// Sort Gists by last modified time
		sort.Slice(gists, func(i, j int) bool {
			return gists[i].UpdatedAt.After(gists[j].UpdatedAt)
		})

		// Check detail mode
		detail, _ := cmd.Flags().GetBool("detail")

		// Print table header
		typeHeader := truncateString("Type", typeWidth)
		titleHeader := truncateString("Title", titleWidth)
		dateHeader := truncateString("Last Modified", dateWidth)
		idHeader := truncateString("Gist ID", idWidth)
		
		if detail {
			revHeader := truncateString("Rev", revWidth)
			fmt.Printf("%s  %s  %s  %s  %s\n", typeHeader, titleHeader, dateHeader, idHeader, revHeader)
			fmt.Println(strings.Repeat("-", typeWidth+titleWidth+dateWidth+idWidth+revWidth+8))
		} else {
			fmt.Printf("%s  %s  %s  %s\n", typeHeader, titleHeader, dateHeader, idHeader)
			fmt.Println(strings.Repeat("-", typeWidth+titleWidth+dateWidth+idWidth+6))
		}

		// Print each Gist information
		for _, g := range gists {
			description := g.Description
			if description == "" {
				description = "(No title)"
			}
			
			gistType := "Private"
			if g.Public {
				gistType = "Public"
			}
			typeStr := truncateString(gistType, typeWidth)
			title := truncateString(description, titleWidth)
			date := truncateString(g.UpdatedAt.Format("2006-01-02 15:04:05"), dateWidth)
			id := truncateString(g.ID, idWidth)

			if detail {
				// Fetch detailed information and history
				gistDetail, err := gist.FetchGistWithHistory(config.Token, g.ID)
				if err != nil {
					continue // Skip if history fetch fails
				}
				rev := truncateString(fmt.Sprintf("%d", gistDetail.RevisionNumber), revWidth)
				fmt.Printf("%s  %s  %s  %s  %s\n", typeStr, title, date, id, rev)
			} else {
				fmt.Printf("%s  %s  %s  %s\n", typeStr, title, date, id)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().Bool("detail", false, "Show detailed information including revision")
} 