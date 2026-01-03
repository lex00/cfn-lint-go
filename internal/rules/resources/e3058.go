// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/schema"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3058{})
}

// E3058 validates that at least one required property is specified.
type E3058 struct{}

func (r *E3058) ID() string { return "E3058" }

func (r *E3058) ShortDesc() string {
	return "At least one property required"
}

func (r *E3058) Description() string {
	return "Validates that resources specify at least one property from a required set (OR logic for required properties)."
}

func (r *E3058) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3058"
}

func (r *E3058) Tags() []string {
	return []string{"resources", "properties", "required"}
}

func (r *E3058) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		constraints := schema.GetResourceConstraints(res.Type)
		if constraints == nil || len(constraints.AnyOf) == 0 {
			continue
		}

		// Check each AnyOf group (at least one property required)
		for _, group := range constraints.AnyOf {
			hasAtLeastOne := false
			for _, prop := range group {
				if _, exists := res.Properties[prop]; exists {
					hasAtLeastOne = true
					break
				}
			}

			if !hasAtLeastOne {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s' (%s): At least one of the following properties must be specified: %v",
						resName, res.Type, group,
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties"},
				})
			}
		}
	}

	return matches
}
