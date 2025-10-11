package domain

type RepoContent struct {
	Items []RepoContentItem
}

type RepoContentItem struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	DownloadURL string `json:"download_url"`
}
