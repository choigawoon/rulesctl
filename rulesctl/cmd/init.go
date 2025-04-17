package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "기본 규칙 파일 생성",
	Long: `현재 디렉토리에 .cursor/rules 디렉토리와 기본 규칙 파일을 생성합니다.
생성되는 파일:
- .cursor/rules/hello.mdc: 기본 인사말 규칙`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 현재 디렉토리 확인
		workDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("작업 디렉토리 확인 실패: %w", err)
		}

		// .cursor/rules 디렉토리 생성
		rulesDir := filepath.Join(workDir, ".cursor", "rules")
		if err := os.MkdirAll(rulesDir, 0755); err != nil {
			return fmt.Errorf(".cursor/rules 디렉토리 생성 실패: %w", err)
		}

		// hello.mdc 파일 생성
		helloPath := filepath.Join(rulesDir, "hello.mdc")
		helloContent := `# Hello Rule

이 규칙은 모든 대화에서 "hello, rulesctl"로 시작하도록 합니다.

<instructions>
1. 모든 대화의 시작에 "hello, rulesctl"을 출력합니다.
2. 이 인사말은 다른 응답 내용보다 먼저 나와야 합니다.
3. 인사말 다음에는 빈 줄을 추가하여 본문과 구분합니다.
</instructions>

<examples>
예시 1:
User: 파일을 생성해줘
Assistant: hello, rulesctl

네, 어떤 파일을 생성하시겠습니까?

예시 2:
User: 코드를 검토해줘
Assistant: hello, rulesctl

코드를 검토해드리겠습니다.
</examples>

<rules>
- 모든 응답은 "hello, rulesctl"로 시작해야 합니다.
- 인사말과 본문 사이에 빈 줄을 넣어야 합니다.
- 인사말은 다른 서식 없이 일반 텍스트로 출력합니다.
</rules>`

		if err := os.WriteFile(helloPath, []byte(helloContent), 0644); err != nil {
			return fmt.Errorf("hello.mdc 파일 생성 실패: %w", err)
		}

		fmt.Println("기본 규칙 파일이 생성되었습니다:")
		fmt.Printf("- %s\n", helloPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
} 