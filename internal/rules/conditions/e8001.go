// Package conditions contains condition validation rules (E8xxx).
package conditions

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E8001{})
}

// E8001 checks that conditions have valid configuration.
type E8001 struct{}

func (r *E8001) ID() string { return "E8001" }

func (r *E8001) ShortDesc() string {
	return "Condition configuration error"
}

func (r *E8001) Description() string {
	return "Checks that each condition has a valid condition function (Fn::Equals, Fn::And, Fn::Or, Fn::Not, Fn::If, Condition)."
}

func (r *E8001) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E8001"
}

func (r *E8001) Tags() []string {
	return []string{"conditions", "configuration"}
}

// Valid top-level condition functions
var validConditionFunctions = map[string]bool{
	"Fn::Equals": true,
	"Fn::And":    true,
	"Fn::Or":     true,
	"Fn::Not":    true,
	"Fn::If":     true,
	"Condition":  true,
}

func (r *E8001) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for name, cond := range tmpl.Conditions {
		if cond.Expression == nil {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Condition '%s' has no expression", name),
				Line:    cond.Node.Line,
				Column:  cond.Node.Column,
				Path:    []string{"Conditions", name},
			})
			continue
		}

		// Check if expression is a valid condition function
		exprMap, ok := cond.Expression.(map[string]any)
		if !ok {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Condition '%s' must be a condition function (Fn::Equals, Fn::And, Fn::Or, Fn::Not, Condition)", name),
				Line:    cond.Node.Line,
				Column:  cond.Node.Column,
				Path:    []string{"Conditions", name},
			})
			continue
		}

		// Check for valid function key
		hasValidFunc := false
		for key := range exprMap {
			if validConditionFunctions[key] {
				hasValidFunc = true
				break
			}
		}

		if !hasValidFunc {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Condition '%s' must use a valid condition function (Fn::Equals, Fn::And, Fn::Or, Fn::Not, Condition)", name),
				Line:    cond.Node.Line,
				Column:  cond.Node.Column,
				Path:    []string{"Conditions", name},
			})
		}
	}

	return matches
}
