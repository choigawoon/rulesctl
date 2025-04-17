package gist

import (
	"context"
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
	Owner       struct {
		Login string `json:"login"`
	} `json:"owner"`
	Files       map[string]struct {
		Filename string `json:"filename"`
		Type     string `json:"type"`
		Language string `json:"language"`
		RawURL   string `json:"raw_url"`
		Size     int    `json:"size"`
		Content  string `json:"content"`
	} `json:"files"`
}

// GetGist는 지정된 ID의 Gist 정보를 가져옵니다
func GetGist(gistID string) (*Gist, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	gist, _, err := client.Gists.Get(ctx, gistID)
	if err != nil {
		return nil, fmt.Errorf("Gist 조회 실패: %w", err)
	}

	// API 응답을 우리의 Gist 구조체로 변환
	result := &Gist{
		ID:          gist.GetID(),
		Description: gist.GetDescription(),
		Public:      gist.GetPublic(),
		UpdatedAt:   gist.GetUpdatedAt().Time,
	}
	result.Owner.Login = gist.GetOwner().GetLogin()

	return result, nil
}

// IsOwnedByCurrentUser는 주어진 Gist가 현재 인증된 사용자의 것인지 확인합니다
func IsOwnedByCurrentUser(gistID string) (bool, error) {
	client, err := getClient()
	if err != nil {
		return false, err
	}

	ctx := context.Background()
	
	// 현재 사용자 정보 가져오기
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return false, fmt.Errorf("사용자 정보 조회 실패: %w", err)
	}
	currentUser := user.GetLogin()

	// Gist 정보 가져오기
	gist, err := GetGist(gistID)
	if err != nil {
		return false, err
	}

	return gist.Owner.Login == currentUser, nil
}

// FetchUserGists는 사용자의 Gist를 가져옵니다
// since가 지정된 경우 해당 시간 이후의 Gist만 가져옵니다
func FetchUserGists(since *time.Time) ([]Gist, error) {
	client := &http.Client{}
	
	url := fmt.Sprintf("%s/gists", baseURL)
	if since != nil {
		url = fmt.Sprintf("%s?since=%s", url, since.Format(time.RFC3339))
	}
	
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

// DeleteGist는 지정된 ID의 Gist를 삭제합니다.
func DeleteGist(gistID string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	_, err = client.Gists.Delete(ctx, gistID)
	return err
} 