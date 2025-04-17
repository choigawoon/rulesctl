package cmd

import (
	"fmt"
	"strings"

	"github.com/choigawoon/rulesctl/internal/gist"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:           "delete [name]",
	Short:         "Delete a rule set",
	Long: `Delete a rule set stored in GitHub Gist.
Rule sets are deleted by searching for their title.
Only rule sets uploaded within the last month can be deleted.

Examples:
  rulesctl delete "my-python-ruleset"    # Delete by title
  rulesctl delete "my-ruleset" --force   # Delete without confirmation`,
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")
		title := args[0]

		// Fetch all Gists
		gists, err := gist.FetchUserGists(nil)
		if err != nil {
			return fmt.Errorf("failed to fetch Gist list: %w", err)
		}

		// Find Gist by title
		var targetGist gist.Gist
		found := false
		for _, g := range gists {
			if g.Description == title {
				targetGist = g
				found = true
				break
			}
		}

		if !found {
			cmd.SilenceUsage = true
			return fmt.Errorf("rule set not found: %s", title)
		}

		// Confirm before deletion
		if !force {
			fmt.Printf("Are you sure you want to delete rule set '%s'? (y/N): ", title)
			var response string
			fmt.Scanln(&response)
			if !strings.EqualFold(response, "y") {
				fmt.Println("Deletion cancelled.")
				return nil
			}
		}

		if err := gist.DeleteGist(targetGist.ID); err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("failed to delete Gist: %w", err)
		}

		fmt.Printf("Rule set '%s' has been deleted.\n", title)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().Bool("force", false, "Delete without confirmation")
} 