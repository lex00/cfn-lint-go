// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/schema"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3032{})
}

// E3032 checks that array properties have valid lengths.
type E3032 struct{}

func (r *E3032) ID() string { return "E3032" }

func (r *E3032) ShortDesc() string {
	return "Array length out of range"
}

func (r *E3032) Description() string {
	return "Checks that array property values meet minimum and maximum length constraints from CloudFormation schemas."
}

func (r *E3032) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3032"
}

func (r *E3032) Tags() []string {
	return []string{"resources", "properties", "array", "length"}
}

func (r *E3032) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		if !schema.HasConstraints(res.Type) {
			continue
		}

		for propName, propValue := range res.Properties {
			// Only check array values
			arr, ok := propValue.([]any)
			if !ok {
				continue
			}

			constraints := schema.GetPropertyConstraints(res.Type, propName)
			if constraints == nil {
				continue
			}

			arrLen := len(arr)

			// Check minimum items
			if constraints.MinItems != nil && arrLen < *constraints.MinItems {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Property '%s' in resource '%s' (%s): array has %d items, minimum is %d",
						propName, resName, res.Type, arrLen, *constraints.MinItems,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", propName},
				})
			}

			// Check maximum items
			if constraints.MaxItems != nil && arrLen > *constraints.MaxItems {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Property '%s' in resource '%s' (%s): array has %d items, maximum is %d",
						propName, resName, res.Type, arrLen, *constraints.MaxItems,
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
