// Package output provides various output formats for cfn-lint results.
package output

import (
	"encoding/json"
	"io"

	"github.com/lex00/cfn-lint-go/pkg/lint"
)

// SARIF represents the SARIF 2.1.0 output format.
// Reference: https://docs.oasis-open.org/sarif/sarif/v2.1.0/sarif-v2.1.0.html
type SARIF struct {
	Schema  string     `json:"$schema"`
	Version string     `json:"version"`
	Runs    []SARIFRun `json:"runs"`
}

// SARIFRun represents a single run in SARIF format.
type SARIFRun struct {
	Tool    SARIFTool     `json:"tool"`
	Results []SARIFResult `json:"results"`
}

// SARIFTool represents the tool metadata.
type SARIFTool struct {
	Driver SARIFDriver `json:"driver"`
}

// SARIFDriver represents the tool driver.
type SARIFDriver struct {
	Name           string      `json:"name"`
	InformationURI string      `json:"informationUri"`
	Version        string      `json:"version"`
	Rules          []SARIFRule `json:"rules,omitempty"`
}

// SARIFRule represents a rule in SARIF format.
type SARIFRule struct {
	ID               string          `json:"id"`
	ShortDescription SARIFMessage    `json:"shortDescription"`
	FullDescription  SARIFMessage    `json:"fullDescription,omitempty"`
	HelpURI          string          `json:"helpUri,omitempty"`
	Properties       SARIFProperties `json:"properties,omitempty"`
}

// SARIFMessage represents a message in SARIF format.
type SARIFMessage struct {
	Text string `json:"text"`
}

// SARIFProperties represents rule properties.
type SARIFProperties struct {
	Tags []string `json:"tags,omitempty"`
}

// SARIFResult represents a single result in SARIF format.
type SARIFResult struct {
	RuleID    string          `json:"ruleId"`
	Level     string          `json:"level"`
	Message   SARIFMessage    `json:"message"`
	Locations []SARIFLocation `json:"locations"`
}

// SARIFLocation represents a location in SARIF format.
type SARIFLocation struct {
	PhysicalLocation SARIFPhysicalLocation `json:"physicalLocation"`
}

// SARIFPhysicalLocation represents a physical location.
type SARIFPhysicalLocation struct {
	ArtifactLocation SARIFArtifactLocation `json:"artifactLocation"`
	Region           SARIFRegion           `json:"region"`
}

// SARIFArtifactLocation represents an artifact location.
type SARIFArtifactLocation struct {
	URI string `json:"uri"`
}

// SARIFRegion represents a region in the source.
type SARIFRegion struct {
	StartLine   int `json:"startLine"`
	StartColumn int `json:"startColumn"`
	EndLine     int `json:"endLine"`
	EndColumn   int `json:"endColumn"`
}

// WriteSARIF writes matches in SARIF 2.1.0 format.
func WriteSARIF(w io.Writer, matches []lint.Match, version string) error {
	// Collect unique rules
	rulesMap := make(map[string]lint.MatchRule)
	for _, m := range matches {
		rulesMap[m.Rule.ID] = m.Rule
	}

	// Build SARIF rules
	var sarifRules []SARIFRule
	for _, rule := range rulesMap {
		sarifRules = append(sarifRules, SARIFRule{
			ID: rule.ID,
			ShortDescription: SARIFMessage{
				Text: rule.ShortDescription,
			},
			FullDescription: SARIFMessage{
				Text: rule.Description,
			},
			HelpURI: rule.Source,
		})
	}

	// Build SARIF results
	var results []SARIFResult
	for _, m := range matches {
		level := "error"
		switch m.Level {
		case "Warning":
			level = "warning"
		case "Informational":
			level = "note"
		}

		results = append(results, SARIFResult{
			RuleID: m.Rule.ID,
			Level:  level,
			Message: SARIFMessage{
				Text: m.Message,
			},
			Locations: []SARIFLocation{
				{
					PhysicalLocation: SARIFPhysicalLocation{
						ArtifactLocation: SARIFArtifactLocation{
							URI: m.Location.Filename,
						},
						Region: SARIFRegion{
							StartLine:   m.Location.Start.LineNumber,
							StartColumn: m.Location.Start.ColumnNumber,
							EndLine:     m.Location.End.LineNumber,
							EndColumn:   m.Location.End.ColumnNumber,
						},
					},
				},
			},
		})
	}

	sarif := SARIF{
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		Version: "2.1.0",
		Runs: []SARIFRun{
			{
				Tool: SARIFTool{
					Driver: SARIFDriver{
						Name:           "cfn-lint-go",
						InformationURI: "https://github.com/lex00/cfn-lint-go",
						Version:        version,
						Rules:          sarifRules,
					},
				},
				Results: results,
			},
		},
	}

	// Ensure we output [] for empty slice, not null
	if sarif.Runs[0].Results == nil {
		sarif.Runs[0].Results = []SARIFResult{}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(sarif)
}
