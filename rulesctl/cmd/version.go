package cmd

import (
	"fmt"

	"github.com/choigawoon/rulesctl/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display rulesctl version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("rulesctl version: %s\n", version.Version)
		if verbose {
			fmt.Printf("Build time: %s\n", version.BuildTime)
			fmt.Printf("Git commit: %s\n", version.GitCommit)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
} 