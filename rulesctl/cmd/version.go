package cmd

import (
	"fmt"

	"github.com/choigawoon/rulesctl/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "rulesctl 버전 정보 출력",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("rulesctl 버전: %s\n", version.Version)
		if verbose {
			fmt.Printf("빌드 시간: %s\n", version.BuildTime)
			fmt.Printf("Git 커밋: %s\n", version.GitCommit)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
} 