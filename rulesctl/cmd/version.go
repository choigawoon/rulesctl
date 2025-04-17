package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version   = "0.1.0"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "rulesctl 버전 정보 출력",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("rulesctl 버전: %s\n", Version)
		if verbose {
			fmt.Printf("빌드 시간: %s\n", BuildTime)
			fmt.Printf("Git 커밋: %s\n", GitCommit)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
} 