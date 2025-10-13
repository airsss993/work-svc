package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/airsss993/work-svc/internal/config"
)

type GitBucketClient struct {
	httpClient *http.Client
	cfg        *config.Config
}

func NewGitBucketClient(cfg *config.Config) *GitBucketClient {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	return &GitBucketClient{
		httpClient: httpClient,
		cfg:        cfg,
	}
}

func (c *GitBucketClient) GetRepositoryContent(ctx context.Context, userID, path string) ([]RepositoryContentResp, error) {
	baseURL := fmt.Sprintf("%s/api/v3/repos/%s/Work/contents", c.cfg.GitBucket.URL, userID)

	if path != "" {
		baseURL = baseURL + "/" + path
	}

	url := baseURL

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "token "+c.cfg.GitBucket.APIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gitbucket API returned status %d: %s", resp.StatusCode, string(body))
	}

	var content []RepositoryContentResp
	if err := json.Unmarshal(body, &content); err != nil {
		return nil, err
	}

	return content, nil
}
