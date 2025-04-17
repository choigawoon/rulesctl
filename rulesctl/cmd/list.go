package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/choigawoon/rulesctl/internal/gist"
	"github.com/choigawoon/rulesctl/pkg/config"
	"github.com/spf13/cobra"
)

const (
	titleWidth = 25    // 제목 최대 너비
	dateWidth  = 19    // 날짜 너비
	idWidth    = 32    // Gist ID 너비
	revWidth   = 8     // Revision 너비
	separator  = "..."
)

// truncateString은 문자열을 지정된 너비로 자르고 필요한 경우 말줄임표를 추가합니다.
func truncateString(s string, width int) string {
	if utf8.RuneCountInString(s) <= width {
		return s + strings.Repeat(" ", width-utf8.RuneCountInString(s))
	}
	return string([]rune(s)[:width-len(separator)]) + separator
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "GIST에 저장된 규칙 목록 출력",
	Long: `GIST에 저장된 모든 규칙 목록을 출력합니다.
기본적으로 [제목] [최종수정시각] [Gist ID] 형식으로 출력됩니다.
--detail 플래그를 사용하면 revision 정보도 함께 표시됩니다.

사용 예시:
  rulesctl list          # 기본 정보만 출력
  rulesctl list --detail # revision 정보 포함하여 출력`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("설정을 로드할 수 없습니다: %w", err)
		}

		if config.Token == "" {
			return fmt.Errorf("GitHub 토큰이 설정되지 않았습니다. 'rulesctl auth' 명령어로 토큰을 설정해주세요.")
		}

		// 토큰 소스 표시
		if os.Getenv("GITHUB_TOKEN") != "" {
			fmt.Println("GitHub 토큰: 환경 변수에서 로드됨")
		} else {
			fmt.Println("GitHub 토큰: 설정 파일에서 로드됨")
		}

		// 최근 1달 이내의 Gist만 조회
		since := time.Now().AddDate(0, -1, 0)
		gists, err := gist.FetchUserGists(&since)
		if err != nil {
			return fmt.Errorf("Gist 목록 조회 실패: %w", err)
		}

		// GIST를 최종 수정 시간 기준으로 정렬
		sort.Slice(gists, func(i, j int) bool {
			return gists[i].UpdatedAt.After(gists[j].UpdatedAt)
		})

		// 상세 모드 여부 확인
		detail, _ := cmd.Flags().GetBool("detail")

		// 테이블 헤더 출력
		titleHeader := truncateString("제목", titleWidth)
		dateHeader := truncateString("최종 수정", dateWidth)
		idHeader := truncateString("Gist ID", idWidth)
		
		if detail {
			revHeader := truncateString("Rev", revWidth)
			fmt.Printf("%s  %s  %s  %s\n", titleHeader, dateHeader, idHeader, revHeader)
			fmt.Println(strings.Repeat("-", titleWidth+dateWidth+idWidth+revWidth+6))
		} else {
			fmt.Printf("%s  %s  %s\n", titleHeader, dateHeader, idHeader)
			fmt.Println(strings.Repeat("-", titleWidth+dateWidth+idWidth+4))
		}

		// 각 GIST 정보 출력
		for _, g := range gists {
			description := g.Description
			if description == "" {
				description = "(제목 없음)"
			}
			
			title := truncateString(description, titleWidth)
			date := truncateString(g.UpdatedAt.Format("2006-01-02 15:04:05"), dateWidth)
			id := truncateString(g.ID, idWidth)

			if detail {
				// 상세 정보와 히스토리 조회
				gistDetail, err := gist.FetchGistWithHistory(config.Token, g.ID)
				if err != nil {
					continue // 히스토리 조회 실패 시 건너뛰기
				}
				rev := truncateString(fmt.Sprintf("%d", gistDetail.Version), revWidth)
				fmt.Printf("%s  %s  %s  %s\n", title, date, id, rev)
			} else {
				fmt.Printf("%s  %s  %s\n", title, date, id)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().Bool("detail", false, "revision 정보를 포함한 상세 정보 출력")
} 