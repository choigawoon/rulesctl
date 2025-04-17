package cmd

import (
	"fmt"

	"github.com/choigawoon/rulesctl/pkg/config"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "GitHub Personal Access Token 설정",
	Long: `GitHub Personal Access Token을 설정합니다.
이 토큰은 Gist API에 접근하는 데 사용됩니다.

토큰은 다음 권한이 필요합니다:
- Gist (read/write) 권한
- repo 권한 (규칙 파일 목록 접근용)

토큰은 ~/.rulesctl/config.json 파일에 안전하게 저장됩니다.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		token, _ := cmd.Flags().GetString("token")
		if token == "" {
			fmt.Print("GitHub Personal Access Token을 입력하세요: ")
			fmt.Scanln(&token)
		}

		if token == "" {
			return fmt.Errorf("토큰이 입력되지 않았습니다")
		}

		if err := config.SaveToken(token); err != nil {
			return fmt.Errorf("토큰을 저장할 수 없습니다: %w", err)
		}

		fmt.Println("토큰이 성공적으로 저장되었습니다.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.Flags().StringP("token", "t", "", "GitHub Personal Access Token")
} 