// Package outputs contains output validation rules (E6xxx).
package outputs

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
	"gopkg.in/yaml.v3"
)

func init() {
	rules.Register(&E6003{})
}

// E6003 checks that output property types are valid.
type E6003 struct{}

func (r *E6003) ID() string { return "E6003" }

func (r *E6003) ShortDesc() string {
	return "Output property type error"
}

func (r *E6003) Description() string {
	return "Checks that output properties have correct types (Description is string, Export is object)."
}

func (r *E6003) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E6003"
}

func (r *E6003) Tags() []string {
	return []string{"outputs", "types"}
}

func (r *E6003) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for name, out := range tmpl.Outputs {
		if out.Node == nil || out.Node.Kind != yaml.MappingNode {
			continue
		}

		for i := 0; i < len(out.Node.Content); i += 2 {
			key := out.Node.Content[i]
			value := out.Node.Content[i+1]

			switch key.Value {
			case "Description":
				// Description must be a string
				if value.Kind != yaml.ScalarNode {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Output '%s' Description must be a string", name),
						Line:    value.Line,
						Column:  value.Column,
						Path:    []string{"Outputs", name, "Description"},
					})
				}

			case "Condition":
				// Condition must be a string (condition name)
				if value.Kind != yaml.ScalarNode {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Output '%s' Condition must be a string (condition name)", name),
						Line:    value.Line,
						Column:  value.Column,
						Path:    []string{"Outputs", name, "Condition"},
					})
				}

			case "Export":
				// Export must be an object with Name key
				if value.Kind != yaml.MappingNode {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Output '%s' Export must be an object with 'Name' property", name),
						Line:    value.Line,
						Column:  value.Column,
						Path:    []string{"Outputs", name, "Export"},
					})
				} else {
					// Check for Name property in Export
					hasName := false
					for j := 0; j < len(value.Content); j += 2 {
						if value.Content[j].Value == "Name" {
							hasName = true
							break
						}
					}
					if !hasName {
						matches = append(matches, rules.Match{
							Message: fmt.Sprintf("Output '%s' Export is missing required 'Name' property", name),
							Line:    value.Line,
							Column:  value.Column,
							Path:    []string{"Outputs", name, "Export"},
						})
					}
				}
			}
		}
	}

	return matches
}
