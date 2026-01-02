// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/schema"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3033{})
}

// E3033 checks that string properties have valid lengths.
type E3033 struct{}

func (r *E3033) ID() string { return "E3033" }

func (r *E3033) ShortDesc() string {
	return "String length out of range"
}

func (r *E3033) Description() string {
	return "Checks that string property values meet minimum and maximum length constraints from CloudFormation schemas."
}

func (r *E3033) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3033"
}

func (r *E3033) Tags() []string {
	return []string{"resources", "properties", "string", "length"}
}

func (r *E3033) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if !schema.HasConstraints(res.Type) {
			continue
		}

		for propName, propValue := range res.Properties {
			// Skip non-string values and intrinsic functions
			strValue, ok := propValue.(string)
			if !ok {
				continue
			}

			constraints := schema.GetPropertyConstraints(res.Type, propName)
			if constraints == nil {
				continue
			}

			strLen := len(strValue)

			// Check minimum length
			if constraints.MinLength != nil && strLen < *constraints.MinLength {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Property '%s' in resource '%s' (%s): string length is %d, minimum is %d",
						propName, resName, res.Type, strLen, *constraints.MinLength,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", propName},
				})
			}

			// Check maximum length
			if constraints.MaxLength != nil && strLen > *constraints.MaxLength {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Property '%s' in resource '%s' (%s): string length is %d, maximum is %d",
						propName, resName, res.Type, strLen, *constraints.MaxLength,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", propName},
				})
			}
		}
	}

	return matches
}
