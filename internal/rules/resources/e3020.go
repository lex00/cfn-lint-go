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
	rules.Register(&E3020{})
}

// E3020 checks for dependent exclusions (properties that must not be present together).
type E3020 struct{}

func (r *E3020) ID() string { return "E3020" }

func (r *E3020) ShortDesc() string {
	return "Dependent property exclusion violation"
}

func (r *E3020) Description() string {
	return "Checks that properties which exclude other properties are not specified together."
}

func (r *E3020) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3020"
}

func (r *E3020) Tags() []string {
	return []string{"resources", "properties", "dependencies"}
}

func (r *E3020) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		constraints := schema.GetResourceConstraints(res.Type)
		if constraints == nil || len(constraints.DependentExcluded) == 0 {
			continue
		}

		// Check each dependent exclusion rule
		for triggerProp, excludedProps := range constraints.DependentExcluded {
			// If trigger property is not present, skip
			if !hasProperty(res.Properties, triggerProp) {
				continue
			}

			// Check if any excluded properties are present
			foundExcluded := []string{}
			for _, prop := range excludedProps {
				if hasProperty(res.Properties, prop) {
					foundExcluded = append(foundExcluded, prop)
				}
			}

			if len(foundExcluded) > 0 {
				sort.Strings(foundExcluded)
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s' (%s): property '%s' cannot be used with: %s",
						resName, res.Type, triggerProp, strings.Join(foundExcluded, ", "),
					),
					Line:   res.Node.Line,
					Column: res.Node.Column,
					Path:   []string{"Resources", resName, "Properties", triggerProp},
				})
			}
		}
	}

	return matches
}
