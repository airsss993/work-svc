package service

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/airsss993/work-svc/internal/client"
	"golang.org/x/net/html"
)

type ProxyService struct {
	gitbucketClient *client.GitBucketClient
}

func NewProxyService(gitbucketClient *client.GitBucketClient) *ProxyService {
	return &ProxyService{
		gitbucketClient: gitbucketClient,
	}
}

func (p *ProxyService) GetHTMLWithBase(ctx context.Context, owner, repo, ref, filePath, baseURL string) ([]byte, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path is required")
	}

	content, err := p.gitbucketClient.GetFileContent(ctx, owner, repo, ref, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file content: %w", err)
	}

	doc, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	dirPath := path.Dir(filePath)
	if dirPath == "." {
		dirPath = ""
	}

	var baseHref string
	if dirPath != "" {
		baseHref = fmt.Sprintf("%s/%s/", baseURL, dirPath)
	} else {
		baseHref = fmt.Sprintf("%s/", baseURL)
	}

	if err := injectBaseTag(doc, baseHref); err != nil {
		return nil, fmt.Errorf("failed to inject base tag: %w", err)
	}

	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return nil, fmt.Errorf("failed to render HTML: %w", err)
	}

	return buf.Bytes(), nil
}

func (p *ProxyService) GetRawFile(ctx context.Context, owner, repo, ref, filePath string) ([]byte, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path is required")
	}

	content, err := p.gitbucketClient.GetFileContent(ctx, owner, repo, ref, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file content: %w", err)
	}

	return content, nil
}

func injectBaseTag(n *html.Node, baseHref string) error {
	head := findHead(n)
	if head == nil {
		return fmt.Errorf("no <head> tag found in HTML")
	}

	baseTag := &html.Node{
		Type: html.ElementNode,
		Data: "base",
		Attr: []html.Attribute{
			{Key: "href", Val: baseHref},
		},
	}

	if head.FirstChild != nil {
		head.InsertBefore(baseTag, head.FirstChild)
	} else {
		head.AppendChild(baseTag)
	}

	return nil
}

func findHead(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && strings.ToLower(n.Data) == "head" {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findHead(c); result != nil {
			return result
		}
	}

	return nil
}
