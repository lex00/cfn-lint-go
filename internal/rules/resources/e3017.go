// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/schema"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3017{})
}

// E3017 checks that at least one property from anyOf sets is present.
type E3017 struct{}

func (r *E3017) ID() string { return "E3017" }

func (r *E3017) ShortDesc() string {
	return "Required anyOf properties missing"
}

func (r *E3017) Description() string {
	return "Checks that at least one property from each anyOf set is specified."
}

func (r *E3017) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3017"
}

func (r *E3017) Tags() []string {
	return []string{"resources", "properties", "anyof"}
}

func (r *E3017) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		constraints := schema.GetResourceConstraints(res.Type)
		if constraints == nil || len(constraints.AnyOf) == 0 {
			continue
		}

		// Check each anyOf set
		for _, anyOfSet := range constraints.AnyOf {
			// Check if at least one property from the set is present
			hasAny := false
			for _, prop := range anyOfSet {
				if hasProperty(res.Properties, prop) {
					hasAny = true
					break
				}
			}

			if !hasAny {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s' (%s) must have at least one of: %s",
						resName, res.Type, strings.Join(anyOfSet, ", "),
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
