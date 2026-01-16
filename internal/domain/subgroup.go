package domain

type SubgroupsResponse struct {
	English  []string `json:"english"`
	Profiles []string `json:"profiles,omitempty"`
	Course   int      `json:"course"`
}
