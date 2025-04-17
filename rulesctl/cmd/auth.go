package cmd

import (
	"fmt"
	"os"
	"golang.org/x/term"

	"github.com/choigawoon/rulesctl/pkg/config"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Set GitHub Personal Access Token",
	Long: `Set GitHub Personal Access Token.
This token is used to access the Gist API.

Required permissions:
- Gist (read/write) permission
- repo permission (for accessing rule file list)

The token is securely stored in ~/.rulesctl/config.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		token, _ := cmd.Flags().GetString("token")
		if token == "" {
			fmt.Print("Enter GitHub Personal Access Token: ")
			tokenBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println() // Add newline
			if err != nil {
				return fmt.Errorf("failed to read token: %w", err)
			}
			token = string(tokenBytes)
		}

		if token == "" {
			return fmt.Errorf("token not provided")
		}

		if err := config.SaveToken(token); err != nil {
			return fmt.Errorf("failed to save token: %w", err)
		}

		fmt.Println("Token saved successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.Flags().StringP("token", "t", "", "GitHub Personal Access Token")
} 