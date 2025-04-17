package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	verbose bool
	force   bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "rulesctl",
	Short: "CLI tool for managing Cursor Rules",
	Long: `rulesctl is a CLI tool for efficiently managing Cursor Rules.
You can store and share rule sets through GitHub Gist.`,
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Set global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Force overwrite on conflicts")
} 