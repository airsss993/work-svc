package client

type RepositoryContentResp struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	DownloadURL string `json:"download_url"`
}
