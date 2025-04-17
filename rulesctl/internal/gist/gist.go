package gist

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	
	"github.com/choigawoon/rulesctl/pkg/config"
)

var baseURL = "https://api.github.com"

const MetaFileName = ".rulesctl.meta.json"

// Gist represents GitHub Gist information
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
	History []struct {
		Version   string    `json:"version"`
		CommitID  string    `json:"commit_id"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"history"`
	RevisionNumber int // Same as GitHub web UI numbering (latest is 1)
}

// FetchUserGists fetches user's Gists
// If since is specified, only fetches Gists after that time
func FetchUserGists(since *time.Time) ([]Gist, error) {
	// Load token from config
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	
	if cfg.Token == "" {
		return nil, fmt.Errorf("GitHub token not set")
	}
	
	client := &http.Client{}
	
	url := fmt.Sprintf("%s/gists", baseURL)
	if since != nil {
		url = fmt.Sprintf("%s?since=%s", url, since.Format(time.RFC3339))
	}
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "token "+cfg.Token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed: %s", resp.Status)
	}

	var gists []Gist
	if err := json.NewDecoder(resp.Body).Decode(&gists); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Filter Gists with .rulesctl.meta.json file
	var rulesctlGists []Gist
	for _, g := range gists {
		if _, hasRulesctlMeta := g.Files[MetaFileName]; hasRulesctlMeta {
			rulesctlGists = append(rulesctlGists, g)
		}
	}

	return rulesctlGists, nil
}

// FetchGistWithHistory fetches detailed information and history of a specific Gist
func FetchGistWithHistory(token, gistID string) (*Gist, error) {
	url := fmt.Sprintf("%s/gists/%s", baseURL, gistID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed: %s", resp.Status)
	}
	
	var gist Gist
	if err := json.NewDecoder(resp.Body).Decode(&gist); err != nil {
		return nil, err
	}

	// Set RevisionNumber based on History length
	// Latest becomes Rev 1
	gist.RevisionNumber = len(gist.History)
	
	return &gist, nil
}

// DeleteGist deletes a Gist with the specified ID
func DeleteGist(gistID string) error {
	// Load token from config
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	
	if cfg.Token == "" {
		return fmt.Errorf("GitHub token not set")
	}
	
	client := &http.Client{}
	url := fmt.Sprintf("%s/gists/%s", baseURL, gistID)
	
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", "token "+cfg.Token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make API request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete Gist: %s", resp.Status)
	}
	
	return nil
} 