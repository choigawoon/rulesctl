package cmd

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/choigawoon/rulesctl/internal/gist"
	"github.com/choigawoon/rulesctl/pkg/config"
)

var testTime = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func getTestToken() string {
	file, err := os.Open(".env.local")
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "GITHUB_PERSONAL_ACCESS_TOKEN=") {
			return strings.TrimPrefix(line, "GITHUB_PERSONAL_ACCESS_TOKEN=")
		}
	}
	return ""
}

func TestListCmd(t *testing.T) {
	token := getTestToken()
	if token == "" {
		t.Skip("테스트 토큰이 없습니다. .env.local 파일에 GITHUB_PERSONAL_ACCESS_TOKEN을 설정해주세요.")
	}

	// 테스트용 설정 생성
	testConfig := &config.Config{
		Token: token,
	}

	// 테스트용 Gist 데이터 생성
	testGists := []gist.Gist{
		{
			ID:          "1",
			Description: "테스트 Gist 1",
			Public:      true,
			UpdatedAt:   testTime,
			Files: map[string]struct {
				Filename string `json:"filename"`
				Type     string `json:"type"`
				Language string `json:"language"`
				RawURL   string `json:"raw_url"`
				Size     int    `json:"size"`
			}{
				"test1.mdc": {
					Filename: "test1.mdc",
					Type:     "text/plain",
					Language: "Markdown",
					RawURL:   "http://example.com/test1.mdc",
					Size:     100,
				},
			},
		},
	}

	// 테스트 케이스
	tests := []struct {
		name           string
		config         *config.Config
		gists          []gist.Gist
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "정상적인 Gist 목록 출력",
			config:         testConfig,
			gists:          testGists,
			expectedOutput: "테스트 Gist 1",
			expectError:    false,
		},
		{
			name:           "토큰이 없는 경우",
			config:         &config.Config{},
			gists:          nil,
			expectedOutput: "GitHub 토큰이 설정되지 않았습니다",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 출력 버퍼 생성
			var buf bytes.Buffer
			listCmd.SetOut(&buf)

			// 테스트 실행
			err := listCmd.RunE(listCmd, []string{})

			// 에러 검증
			if (err != nil) != tt.expectError {
				t.Errorf("예상된 에러: %v, 실제: %v", tt.expectError, err)
			}

			// 출력 검증
			output := buf.String()
			if !tt.expectError && !bytes.Contains([]byte(output), []byte(tt.expectedOutput)) {
				t.Errorf("예상된 출력: %s, 실제: %s", tt.expectedOutput, output)
			}
		})
	}
} 