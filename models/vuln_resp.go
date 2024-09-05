package models

import (
	"encoding/json"
	"fmt"
)

type Vulnerability struct {
	Id       string   `json:"id"`
	Summary  string   `json:"summary"`
	Details  string   `json:"details"`
	Severity string   `json:"severity"`
	Aliases  []string `json:"aliases"`
}

type VulnerabilityResult struct {
	Vulnerabilities []Vulnerability `json:"vulns"`
}

func NewVulnerabilityResult(jsonData []byte) (*VulnerabilityResult, error) {
	var rawData struct {
		Vulns []struct {
			ID               string   `json:"id"`
			Summary          string   `json:"summary"`
			Details          string   `json:"details"`
			Aliases          []string `json:"aliases"`
			DatabaseSpecific struct {
				Severity string `json:"severity"`
			} `json:"database_specific"`
		} `json:"vulns"`
	}

	if err := json.Unmarshal(jsonData, &rawData); err != nil {
		return nil, err
	}

	var vulnerabilities []Vulnerability
	for _, rawVulnerability := range rawData.Vulns {
		if rawVulnerability.ID != "" && rawVulnerability.Summary != "" && rawVulnerability.DatabaseSpecific.Severity != "" {
			fmt.Println()
			vulnerabilities = append(vulnerabilities, Vulnerability{
				Id:       rawVulnerability.ID,
				Summary:  rawVulnerability.Summary,
				Details:  rawVulnerability.Details,
				Severity: rawVulnerability.DatabaseSpecific.Severity,
				Aliases:  rawVulnerability.Aliases,
			})
		}
	}

	vulnerabilityResult := &VulnerabilityResult{
		Vulnerabilities: vulnerabilities,
	}

	return vulnerabilityResult, nil
}
