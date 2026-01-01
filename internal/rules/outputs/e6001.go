// Package outputs contains output validation rules (E6xxx).
package outputs

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
	"gopkg.in/yaml.v3"
)

func init() {
	rules.Register(&E6001{})
}

// E6001 checks that outputs have valid property structure.
type E6001 struct{}

func (r *E6001) ID() string { return "E6001" }

func (r *E6001) ShortDesc() string {
	return "Output property structure error"
}

func (r *E6001) Description() string {
	return "Checks that each output has valid properties (Value, Description, Export, Condition)."
}

func (r *E6001) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E6001"
}

func (r *E6001) Tags() []string {
	return []string{"outputs", "configuration"}
}

// Valid output properties per CloudFormation spec
var validOutputProperties = map[string]bool{
	"Value":       true,
	"Description": true,
	"Export":      true,
	"Condition":   true,
}

func (r *E6001) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for name, out := range tmpl.Outputs {
		if out.Node == nil || out.Node.Kind != yaml.MappingNode {
			continue
		}

		// Check for invalid properties
		for i := 0; i < len(out.Node.Content); i += 2 {
			propKey := out.Node.Content[i]
			if !validOutputProperties[propKey.Value] {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Output '%s' has invalid property '%s'", name, propKey.Value),
					Line:    propKey.Line,
					Column:  propKey.Column,
					Path:    []string{"Outputs", name, propKey.Value},
				})
			}
		}
	}

	return matches
}
