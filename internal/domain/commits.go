package domain

import "time"

type CommitResp struct {
	Count   int      `json:"count"`
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
