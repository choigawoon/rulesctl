package gist

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var baseURL = "https://api.github.com"

const MetaFileName = ".rulesctl.meta.json"

// Gist는 GitHub Gist의 기본 정보를 나타냅니다
type Gist struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Public      bool      `json:"public"`
	UpdatedAt   time.Time `json:"updated_at"`
	Files       map[string]struct {
		Filename string `json:"filename"`
		Type     string `json:"type"`
		Language string `json:"language"`
		RawURL   string `json:"raw_url"`
		Size     int    `json:"size"`
		Content  string `json:"content"`
	} `json:"files"`
}

// FetchUserGists는 사용자의 모든 Gist를 가져옵니다
func FetchUserGists(token string, since time.Time) ([]Gist, error) {
	client := &http.Client{}
	
	// since 파라미터 추가
	url := fmt.Sprintf("%s/gists?since=%s", baseURL, since.Format(time.RFC3339))
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("요청 생성 실패: %w", err)
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API 요청 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 요청 실패: %s", resp.Status)
	}

	var gists []Gist
	if err := json.NewDecoder(resp.Body).Decode(&gists); err != nil {
		return nil, fmt.Errorf("응답 파싱 실패: %w", err)
	}

	// .rulesctl.meta.json 파일이 있는 Gist만 필터링
	var rulesctlGists []Gist
	for _, g := range gists {
		if _, hasRulesctlMeta := g.Files[MetaFileName]; hasRulesctlMeta {
			rulesctlGists = append(rulesctlGists, g)
		}
	}

	return rulesctlGists, nil
} 