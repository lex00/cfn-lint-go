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
	rules.Register(&E3014{})
}

// E3014 checks for mutually exclusive properties.
type E3014 struct{}

func (r *E3014) ID() string { return "E3014" }

func (r *E3014) ShortDesc() string {
	return "Mutually exclusive properties"
}

func (r *E3014) Description() string {
	return "Checks that mutually exclusive properties are not specified together."
}

func (r *E3014) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3014"
}

func (r *E3014) Tags() []string {
	return []string{"resources", "properties", "exclusive"}
}

func (r *E3014) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		constraints := schema.GetResourceConstraints(res.Type)
		if constraints == nil || len(constraints.MutuallyExclusive) == 0 {
			continue
		}

		// Check each mutually exclusive set
		for _, exclusiveSet := range constraints.MutuallyExclusive {
			// Find which properties from the set are present
			present := []string{}
			for _, prop := range exclusiveSet {
				if hasProperty(res.Properties, prop) {
					present = append(present, prop)
				}
			}

			// If more than one is present, report error
			if len(present) > 1 {
				sort.Strings(present)
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s' (%s) has mutually exclusive properties: %s",
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

// hasProperty checks if a property exists, supporting dot notation for nested properties.
func hasProperty(props map[string]any, propPath string) bool {
	parts := strings.Split(propPath, ".")
	current := props

	for i, part := range parts {
		val, ok := current[part]
		if !ok {
			return false
		}

		// Last part - just check existence
		if i == len(parts)-1 {
			return true
		}

		// Navigate deeper
		nested, ok := val.(map[string]any)
		if !ok {
			return false
		}
		current = nested
	}

	return true
}
