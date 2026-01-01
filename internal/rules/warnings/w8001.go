// Package warnings contains warning-level rules (Wxxx).
package warnings

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W8001{})
}

// W8001 warns about conditions that are defined but never used.
type W8001 struct{}

func (r *W8001) ID() string { return "W8001" }

func (r *W8001) ShortDesc() string {
	return "Unused condition"
}

func (r *W8001) Description() string {
	return "Warns when a condition is defined but never used in resources, outputs, or other conditions."
}

func (r *W8001) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/conditions-section-structure.html"
}

func (r *W8001) Tags() []string {
	return []string{"warnings", "conditions", "unused"}
}

func (r *W8001) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	if len(tmpl.Conditions) == 0 {
		return matches
	}

	// Find all condition references
	usedConditions := make(map[string]bool)

	// Check Resources for Condition attribute
	for _, res := range tmpl.Resources {
		if res.Condition != "" {
			usedConditions[res.Condition] = true
		}
		// Also check for Fn::If in properties
		findConditionRefs(res.Properties, usedConditions)
	}

	// Check Outputs for Condition attribute
	for _, out := range tmpl.Outputs {
		if out.Condition != "" {
			usedConditions[out.Condition] = true
		}
		findConditionRefs(out.Value, usedConditions)
	}

	// Check Conditions section for Fn::Condition references
	for _, cond := range tmpl.Conditions {
		findConditionRefs(cond.Expression, usedConditions)
	}

	// Report unused conditions
	for condName := range tmpl.Conditions {
		if !usedConditions[condName] {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Condition '%s' is defined but never used", condName),
				Path:    []string{"Conditions", condName},
			})
		}
	}

	return matches
}

func findConditionRefs(v any, usedConditions map[string]bool) {
	switch val := v.(type) {
	case map[string]any:
		// Check for Fn::If
		if fnIf, ok := val["Fn::If"].([]any); ok {
			if len(fnIf) >= 1 {
				if condName, ok := fnIf[0].(string); ok {
					usedConditions[condName] = true
				}
			}
		}
		// Check for Condition (used in Fn::And, Fn::Or, Fn::Not)
		if condRef, ok := val["Condition"].(string); ok {
			usedConditions[condRef] = true
		}
		// Recurse into children
		for _, child := range val {
			findConditionRefs(child, usedConditions)
		}
	case []any:
		for _, child := range val {
			findConditionRefs(child, usedConditions)
		}
	}
}
