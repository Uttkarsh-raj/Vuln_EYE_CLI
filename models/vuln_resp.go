package models

type Vulnerability struct {
	Aliases []string `json:"aliases"`
	Summary string   `json:"summary"`
}

type VerboseResp struct {
	Vulns []Vulnerability `json:"vulns"`
}
