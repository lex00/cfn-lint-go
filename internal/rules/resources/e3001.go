// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
	"gopkg.in/yaml.v3"
)

func init() {
	rules.Register(&E3001{})
}

// E3001 checks that resources have valid configuration.
type E3001 struct{}

func (r *E3001) ID() string { return "E3001" }

func (r *E3001) ShortDesc() string {
	return "Resource configuration error"
}

func (r *E3001) Description() string {
	return "Checks that each resource has valid properties and required Type field."
}

func (r *E3001) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3001"
}

func (r *E3001) Tags() []string {
	return []string{"resources", "configuration"}
}

// Valid resource properties per CloudFormation spec
var validResourceProperties = map[string]bool{
	"Type":                true,
	"Properties":          true,
	"DependsOn":           true,
	"Condition":           true,
	"Metadata":            true,
	"DeletionPolicy":      true,
	"UpdatePolicy":        true,
	"UpdateReplacePolicy": true,
	"CreationPolicy":      true,
}

func (r *E3001) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for name, res := range tmpl.Resources {
		// Check for missing Type (required)
		if res.Type == "" {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Resource '%s' is missing required property 'Type'", name),
				Line:    res.Node.Line,
				Column:  res.Node.Column,
				Path:    []string{"Resources", name},
			})
		}

		// Check for invalid properties
		if res.Node != nil && res.Node.Kind == yaml.MappingNode {
			for i := 0; i < len(res.Node.Content); i += 2 {
				propKey := res.Node.Content[i]
				if !validResourceProperties[propKey.Value] {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Resource '%s' has invalid property '%s'", name, propKey.Value),
						Line:    propKey.Line,
						Column:  propKey.Column,
						Path:    []string{"Resources", name, propKey.Value},
					})
				}
			}
		}
	}

	return matches
}
