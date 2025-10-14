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

func (c *GitBucketClient) GetRepositoryContent(ctx context.Context, owner, path string) ([]RepositoryContentResp, error) {
	baseURL := fmt.Sprintf("%s/api/v3/repos/%s/Work/contents", c.cfg.GitBucket.URL, owner)

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

func (c *GitBucketClient) GetFileContent(ctx context.Context, owner, repo, ref, path string) ([]byte, error) {
	url := fmt.Sprintf("%s/api/v3/repos/%s/%s/raw/%s/%s", c.cfg.GitBucket.URL, owner, repo, ref, path)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "token "+c.cfg.GitBucket.APIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gitbucket API returned status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func (c *GitBucketClient) GetRepositoryInfo(ctx context.Context, owner, repo string) (*RepositoryInfoResp, error) {
	url := fmt.Sprintf("%s/api/v3/repos/%s/%s", c.cfg.GitBucket.URL, owner, repo)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "token "+c.cfg.GitBucket.APIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gitbucket API returned status %d: %s", resp.StatusCode, string(body))
	}

	var repoInfo RepositoryInfoResp
	if err := json.Unmarshal(body, &repoInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &repoInfo, nil
}

func (c *GitBucketClient) GetCommitsList(ctx context.Context, owner, repo string, perPage, page int) (CommitListResp, error) {
	url := fmt.Sprintf("%s/api/v3/repos/%s/%s/commits?per_page=%d&page=%d", c.cfg.GitBucket.URL, owner, repo, perPage, page)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return CommitListResp{}, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "token "+c.cfg.GitBucket.APIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return CommitListResp{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CommitListResp{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return CommitListResp{}, fmt.Errorf("gitbucket API returned status %d: %s", resp.StatusCode, string(body))
	}

	var commits CommitListResp
	if err := json.Unmarshal(body, &commits); err != nil {
		return CommitListResp{}, fmt.Errorf("failed to unmarshal commits: %w", err)
	}
	return commits, nil
}
