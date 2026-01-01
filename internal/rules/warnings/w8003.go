// Package warnings contains warning-level rules (Wxxx).
package warnings

import (
	"fmt"
	"reflect"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&W8003{})
}

// W8003 warns when Fn::Equals compares identical values (always true/false).
type W8003 struct{}

func (r *W8003) ID() string { return "W8003" }

func (r *W8003) ShortDesc() string {
	return "Fn::Equals with static result"
}

func (r *W8003) Description() string {
	return "Warns when Fn::Equals compares two identical static values, resulting in a condition that is always true or always false."
}

func (r *W8003) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-conditions.html#intrinsic-function-reference-conditions-equals"
}

func (r *W8003) Tags() []string {
	return []string{"warnings", "conditions", "equals"}
}

func (r *W8003) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for condName, cond := range tmpl.Conditions {
		staticEquals := findStaticEquals(cond.Expression, condName)
		matches = append(matches, staticEquals...)
	}

	return matches
}

func findStaticEquals(v any, condName string) []rules.Match {
	var matches []rules.Match

	switch val := v.(type) {
	case map[string]any:
		if equals, ok := val["Fn::Equals"].([]any); ok {
			if len(equals) == 2 {
				// Check if both values are static and identical
				if isStaticValue(equals[0]) && isStaticValue(equals[1]) {
					if reflect.DeepEqual(equals[0], equals[1]) {
						matches = append(matches, rules.Match{
							Message: fmt.Sprintf("Fn::Equals in condition '%s' compares identical values (always true)", condName),
							Path:    []string{"Conditions", condName},
						})
					} else {
						// Both are static but different - always false
						matches = append(matches, rules.Match{
							Message: fmt.Sprintf("Fn::Equals in condition '%s' compares different static values (always false)", condName),
							Path:    []string{"Conditions", condName},
						})
					}
				}
			}
		}
		// Recurse into nested conditions
		for _, child := range val {
			matches = append(matches, findStaticEquals(child, condName)...)
		}
	case []any:
		for _, child := range val {
			matches = append(matches, findStaticEquals(child, condName)...)
		}
	}

	return matches
}

func isStaticValue(v any) bool {
	switch val := v.(type) {
	case string, int, int64, float64, bool:
		return true
	case map[string]any:
		// Check if it's an intrinsic function
		for key := range val {
			if key == "Ref" || key == "Fn::GetAtt" || key == "Fn::Sub" ||
				key == "Fn::Join" || key == "Fn::Select" || key == "Fn::If" ||
				key == "Fn::GetAZs" || key == "Fn::ImportValue" || key == "Fn::FindInMap" {
				return false
			}
		}
		return true
	case []any:
		// Arrays could be static if all elements are static
		for _, elem := range val {
			if !isStaticValue(elem) {
				return false
			}
		}
		return true
	default:
		return true
	}
}
