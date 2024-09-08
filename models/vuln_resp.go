package models

type Vulnerability struct {
	Aliases  []string   `json:"aliases"`
	Summary  string     `json:"summary"`
	Affected []Affected `json:"affected"`
}

type VerboseResp struct {
	Vulns []Vulnerability `json:"vulns"`
}

type Affected struct {
	Package Package `json:"package"`
	Ranges  []Range `json:"ranges"`
}

type Package struct {
	Name string `json:"name"`
}

type Range struct {
	Type   string  `json:"type"`
	Events []Event `json:"events"`
}

type Event struct {
	Introduced string `json:"introduced,omitempty"`
	Fixed      string `json:"fixed,omitempty"`
}
