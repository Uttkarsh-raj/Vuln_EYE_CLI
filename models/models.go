package models

type DepsPayload struct {
	Version string `json:"version"`
	Package struct {
		Name string `json:"name"`
	} `json:"package"`
}
