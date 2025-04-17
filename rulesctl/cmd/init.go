package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create default rules directory",
	Long: `Create .cursor/rules directory in the current directory.
Use --sample flag to create example rule files.

Created files:
- .cursor/rules/                  : Rules directory
- .cursor/rules/hello.mdc         : (with --sample) Basic greeting rule`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check current directory
		workDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		// Create .cursor/rules directory
		rulesDir := filepath.Join(workDir, ".cursor", "rules")
		if err := os.MkdirAll(rulesDir, 0755); err != nil {
			return fmt.Errorf("failed to create .cursor/rules directory: %w", err)
		}

		fmt.Printf(".cursor/rules directory created: %s\n", rulesDir)

		// Create example files only if --sample flag is used
		sample, _ := cmd.Flags().GetBool("sample")
		if !sample {
			return nil
		}

		// Create hello.mdc file
		helloPath := filepath.Join(rulesDir, "hello.mdc")
		helloContent := `# Hello Rule

This rule makes every conversation start with "hello, rulesctl".

<instructions>
1. Output "hello, rulesctl" at the start of every conversation.
2. This greeting must come before any other response content.
3. Add a blank line after the greeting to separate it from the main content.
</instructions>

<examples>
Example 1:
User: Create a file
Assistant: hello, rulesctl

What kind of file would you like to create?

Example 2:
User: Review my code
Assistant: hello, rulesctl

I'll review your code.
</examples>

<rules>
- Every response must start with "hello, rulesctl"
- There must be a blank line between the greeting and the main content
- The greeting should be output as plain text without any formatting
</rules>`

		if err := os.WriteFile(helloPath, []byte(helloContent), 0644); err != nil {
			return fmt.Errorf("failed to create hello.mdc file: %w", err)
		}

		fmt.Println("Example rule files created:")
		fmt.Printf("- %s\n", helloPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().Bool("sample", false, "Create example rule files")
} 