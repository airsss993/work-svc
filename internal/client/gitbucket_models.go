package client

import "time"

type RepositoryContentResp struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	DownloadURL string `json:"download_url"`
}

type RepositoryInfoResp struct {
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	DefaultBranch string `json:"default_branch"`
}

type CommitListResp []CommitItem

type CommitItem struct {
	Commit CommitInfo   `json:"commit"`
	Files  []CommitFile `json:"files"`
}

type CommitInfo struct {
	Message string       `json:"message"`
	Author  CommitAuthor `json:"author"`
}

type CommitAuthor struct {
	Date *time.Time `json:"date"`
}

type CommitFile struct {
	Filename string `json:"filename"`
}
