// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"
	"sort"
	"strings"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/schema"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3018{})
}

// E3018 checks that exactly one property from oneOf sets is present.
type E3018 struct{}

func (r *E3018) ID() string { return "E3018" }

func (r *E3018) ShortDesc() string {
	return "Required oneOf property missing or multiple specified"
}

func (r *E3018) Description() string {
	return "Checks that exactly one property from each oneOf set is specified."
}

func (r *E3018) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3018"
}

func (r *E3018) Tags() []string {
	return []string{"resources", "properties", "oneof"}
}

func (r *E3018) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		constraints := schema.GetResourceConstraints(res.Type)
		if constraints == nil || len(constraints.OneOf) == 0 {
			continue
		}

		// Check each oneOf set
		for _, oneOfSet := range constraints.OneOf {
			// Find which properties from the set are present
			present := []string{}
			for _, prop := range oneOfSet {
				if hasProperty(res.Properties, prop) {
					present = append(present, prop)
				}
			}

			// Exactly one must be present
			if len(present) == 0 {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s' (%s) must have exactly one of: %s",
						resName, res.Type, strings.Join(oneOfSet, ", "),
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties"},
				})
			} else if len(present) > 1 {
				sort.Strings(present)
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s' (%s) has multiple oneOf properties but only one is allowed: %s",
						resName, res.Type, strings.Join(present, ", "),
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
