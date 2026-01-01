// Package parameters contains parameter validation rules (E2xxx).
package parameters

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
	"gopkg.in/yaml.v3"
)

func init() {
	rules.Register(&E2001{})
}

// E2001 checks that parameters have valid configuration.
type E2001 struct{}

func (r *E2001) ID() string { return "E2001" }

func (r *E2001) ShortDesc() string {
	return "Parameter configuration error"
}

func (r *E2001) Description() string {
	return "Checks that each parameter has valid properties and required Type field."
}

func (r *E2001) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E2001"
}

func (r *E2001) Tags() []string {
	return []string{"parameters", "configuration"}
}

// Valid parameter properties per CloudFormation spec
var validParamProperties = map[string]bool{
	"Type":                  true,
	"Default":               true,
	"AllowedPattern":        true,
	"AllowedValues":         true,
	"ConstraintDescription": true,
	"Description":           true,
	"MaxLength":             true,
	"MaxValue":              true,
	"MinLength":             true,
	"MinValue":              true,
	"NoEcho":                true,
}

func (r *E2001) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for name, param := range tmpl.Parameters {
		// Check for missing Type (required)
		if param.Type == "" {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Parameter '%s' is missing required property 'Type'", name),
				Line:    param.Node.Line,
				Column:  param.Node.Column,
				Path:    []string{"Parameters", name},
			})
		}

		// Check for invalid properties
		if param.Node != nil && param.Node.Kind == yaml.MappingNode {
			for i := 0; i < len(param.Node.Content); i += 2 {
				propKey := param.Node.Content[i]
				if !validParamProperties[propKey.Value] {
					matches = append(matches, rules.Match{
						Message: fmt.Sprintf("Parameter '%s' has invalid property '%s'", name, propKey.Value),
						Line:    propKey.Line,
						Column:  propKey.Column,
						Path:    []string{"Parameters", name, propKey.Value},
					})
				}
			}
		}
	}

	return matches
}
