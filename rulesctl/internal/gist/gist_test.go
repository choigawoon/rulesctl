package gist

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetchUserGists(t *testing.T) {
	// 테스트 서버 생성
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 인증 헤더 확인
		if r.Header.Get("Authorization") != "token test-token" {
			t.Errorf("잘못된 인증 헤더: %s", r.Header.Get("Authorization"))
		}

		// 테스트용 Gist 데이터 생성
		gists := []Gist{
			{
				ID:          "1",
				Description: "테스트 Gist 1",
				Public:      true,
				UpdatedAt:   time.Now(),
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

		// JSON 응답 생성
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(gists)
	}))
	defer ts.Close()

	// 테스트 서버 URL로 API 클라이언트 설정
	oldBaseURL := baseURL
	baseURL = ts.URL
	defer func() { baseURL = oldBaseURL }()

	// Gist 목록 가져오기 테스트
	gists, err := FetchUserGists("test-token")
	if err != nil {
		t.Errorf("Gist 목록 가져오기 실패: %v", err)
	}

	// 결과 검증
	if len(gists) != 1 {
		t.Errorf("예상된 Gist 수: 1, 실제: %d", len(gists))
	}

	if gists[0].Description != "테스트 Gist 1" {
		t.Errorf("예상된 설명: '테스트 Gist 1', 실제: '%s'", gists[0].Description)
	}
} 