package domain

import "time"

type CommitResp struct {
	Count   int      `json:"count"`
	Page    int      `json:"page"`
	PerPage int      `json:"per_page"`
	HasNext bool     `json:"has_next"`
	Commits []Commit `json:"commits"`
}

type Commit struct {
	Message string       `json:"message"`
	Date    *time.Time   `json:"date"`
	Files   []CommitFile `json:"files"`
}

type CommitFile struct {
	Filename string `json:"filename"`
}
