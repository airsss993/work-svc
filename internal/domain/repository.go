package domain

import "time"

type RepoContent struct {
	Items []RepoContentItem `json:"items"`
}
type RepoContentItem struct {
	Type         string    `json:"type"`
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	DownloadURL  string    `json:"download_url"`
	LastModified time.Time `json:"last_modified"`
}
