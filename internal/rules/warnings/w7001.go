// Package warnings contains warning-level rules (Wxxx).
package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W7001{})
}

// W7001 warns about mappings that are defined but never used.
type W7001 struct{}

func (r *W7001) ID() string { return "W7001" }

func (r *W7001) ShortDesc() string {
	return "Unused mapping"
}

func (r *W7001) Description() string {
	return "Warns when a mapping is defined but never referenced via Fn::FindInMap."
}

func (r *W7001) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/mappings-section-structure.html"
}

func (r *W7001) Tags() []string {
	return []string{"warnings", "mappings", "unused"}
}

func (r *W7001) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	if len(tmpl.Mappings) == 0 {
		return matches
	}

	// Find all mapping references
	usedMappings := make(map[string]bool)

	// Check Resources
	for _, res := range tmpl.Resources {
		findMappingRefs(res.Properties, usedMappings)
	}

	// Check Outputs
	for _, out := range tmpl.Outputs {
		findMappingRefs(out.Value, usedMappings)
	}

	// Check Conditions
	for _, cond := range tmpl.Conditions {
		findMappingRefs(cond.Expression, usedMappings)
	}

	// Report unused mappings
	for mappingName := range tmpl.Mappings {
		if !usedMappings[mappingName] {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Mapping '%s' is defined but never used", mappingName),
				Path:    []string{"Mappings", mappingName},
			})
		}
	}

	return matches
}

func findMappingRefs(v any, usedMappings map[string]bool) {
	switch val := v.(type) {
	case map[string]any:
		// Check for Fn::FindInMap
		if findInMap, ok := val["Fn::FindInMap"].([]any); ok {
			if len(findInMap) >= 1 {
				if mappingName, ok := findInMap[0].(string); ok {
					usedMappings[mappingName] = true
				}
			}
		}
		// Recurse into children
		for _, child := range val {
			findMappingRefs(child, usedMappings)
		}
	case []any:
		for _, child := range val {
			findMappingRefs(child, usedMappings)
		}
	}
}
