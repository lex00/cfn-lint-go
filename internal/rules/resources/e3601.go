// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"encoding/json"
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3601{})
}

// E3601 validates the structure of a StateMachine definition.
type E3601 struct{}

func (r *E3601) ID() string { return "E3601" }

func (r *E3601) ShortDesc() string {
	return "Validate StateMachine definition structure"
}

func (r *E3601) Description() string {
	return "Validates that AWS::StepFunctions::StateMachine resources have a valid State Machine definition structure conforming to Amazon States Language specification."
}

func (r *E3601) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3601"
}

func (r *E3601) Tags() []string {
	return []string{"resources", "properties", "stepfunctions", "statemachine"}
}

func (r *E3601) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if res.Type != "AWS::StepFunctions::StateMachine" {
			continue
		}

		// Check DefinitionString or Definition
		defString, hasDefString := res.Properties["DefinitionString"]
		definition, hasDefinition := res.Properties["Definition"]

		if !hasDefString && !hasDefinition {
			continue
		}

		// Validate Definition if present (DefinitionString is just a string, can't validate structure)
		if hasDefinition {
			// Skip intrinsic functions
			if isIntrinsicFunction(definition) {
				continue
			}

			defMap, ok := definition.(map[string]any)
			if !ok {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': StateMachine Definition must be an object",
						resName,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", "Definition"},
				})
				continue
			}

			// Validate required fields according to Amazon States Language
			if _, hasStates := defMap["States"]; !hasStates {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': StateMachine Definition must have 'States' field",
						resName,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", "Definition"},
				})
			}

			if _, hasStartAt := defMap["StartAt"]; !hasStartAt {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': StateMachine Definition must have 'StartAt' field",
						resName,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", "Definition"},
				})
			}
		}

		// If using DefinitionString, try to parse as JSON
		if hasDefString {
			defStr, ok := defString.(string)
			if !ok || isIntrinsicFunction(defString) {
				continue
			}

			var defMap map[string]any
			if err := json.Unmarshal([]byte(defStr), &defMap); err != nil {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': StateMachine DefinitionString must be valid JSON: %v",
						resName, err,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", "DefinitionString"},
				})
				continue
			}

			// Validate required fields
			if _, hasStates := defMap["States"]; !hasStates {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': StateMachine DefinitionString must have 'States' field",
						resName,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", "DefinitionString"},
				})
			}

			if _, hasStartAt := defMap["StartAt"]; !hasStartAt {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s': StateMachine DefinitionString must have 'StartAt' field",
						resName,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", "DefinitionString"},
				})
			}
		}
	}

	return matches
}
