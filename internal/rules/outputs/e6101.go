// Package outputs contains output validation rules (E6xxx).
package outputs

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E6101{})
}

// E6101 checks that output values resolve to strings.
type E6101 struct{}

func (r *E6101) ID() string { return "E6101" }

func (r *E6101) ShortDesc() string {
	return "Output Value must be a string"
}

func (r *E6101) Description() string {
	return "Checks that output Value properties resolve to string types."
}

func (r *E6101) Source() string {
	return "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/outputs-section-structure.html"
}

func (r *E6101) Tags() []string {
	return []string{"outputs", "types"}
}

func (r *E6101) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for name, out := range tmpl.Outputs {
		if out.Value == nil {
			continue // E6002 handles missing Value
		}

		if !isStringOrIntrinsic(out.Value) {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Output '%s' Value must be a string, got %T", name, out.Value),
				Line:    out.Node.Line,
				Column:  out.Node.Column,
				Path:    []string{"Outputs", name, "Value"},
			})
		}
	}

	return matches
}

// isStringOrIntrinsic checks if a value is a string or a CloudFormation intrinsic
// function that could resolve to a string at runtime.
func isStringOrIntrinsic(v any) bool {
	switch val := v.(type) {
	case string:
		return true
	case map[string]any:
		// CloudFormation intrinsic functions that return strings
		stringFunctions := []string{
			"Ref",
			"Fn::Base64",
			"Fn::GetAtt",
			"Fn::GetAZs",
			"Fn::ImportValue",
			"Fn::Join",
			"Fn::Select",
			"Fn::Sub",
			"Fn::If",
			"Fn::FindInMap",
		}
		for _, fn := range stringFunctions {
			if _, ok := val[fn]; ok {
				return true
			}
		}
		return false
	default:
		// Numbers, booleans, arrays, etc. are not valid output values
		return false
	}
}
