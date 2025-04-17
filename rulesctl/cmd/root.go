package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// 전역 플래그
	verbose bool
	force   bool
)

// rootCmd는 기본 명령어를 나타냅니다
var rootCmd = &cobra.Command{
	Use:   "rulesctl",
	Short: "Cursor Rules 관리 CLI 도구",
	Long: `rulesctl은 Cursor Rules를 효율적으로 관리하기 위한 CLI 도구입니다.
GitHub Gist를 통해 규칙 세트를 저장하고 공유할 수 있습니다.`,
}

// Execute는 root 명령어를 실행합니다
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// 전역 플래그 설정
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "상세 로깅 활성화")
	rootCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "충돌 시 강제 덮어쓰기")
} 