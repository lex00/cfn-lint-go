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
	rules.Register(&E3021{})
}

// E3021 checks for dependent requirements (properties that require other properties).
type E3021 struct{}

func (r *E3021) ID() string { return "E3021" }

func (r *E3021) ShortDesc() string {
	return "Dependent property requirement missing"
}

func (r *E3021) Description() string {
	return "Checks that properties which require other properties have those dependencies present."
}

func (r *E3021) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3021"
}

func (r *E3021) Tags() []string {
	return []string{"resources", "properties", "dependencies"}
}

func (r *E3021) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		constraints := schema.GetResourceConstraints(res.Type)
		if constraints == nil || len(constraints.DependentRequired) == 0 {
			continue
		}

		// Check each dependent requirement rule
		for triggerProp, requiredProps := range constraints.DependentRequired {
			// If trigger property is not present, skip
			if !hasProperty(res.Properties, triggerProp) {
				continue
			}

			// Check if all required properties are present
			missingProps := []string{}
			for _, prop := range requiredProps {
				if !hasProperty(res.Properties, prop) {
					missingProps = append(missingProps, prop)
				}
			}

			if len(missingProps) > 0 {
				sort.Strings(missingProps)
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s' (%s): property '%s' requires: %s",
						resName, res.Type, triggerProp, strings.Join(missingProps, ", "),
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
