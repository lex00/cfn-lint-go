// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/schema"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3040{})
}

// E3040 checks that read-only properties are not specified.
type E3040 struct{}

func (r *E3040) ID() string { return "E3040" }

func (r *E3040) ShortDesc() string {
	return "Read-only property specified"
}

func (r *E3040) Description() string {
	return "Checks that read-only properties (returned by CloudFormation but not settable) are not specified in templates."
}

func (r *E3040) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3040"
}

func (r *E3040) Tags() []string {
	return []string{"resources", "properties", "readonly"}
}

func (r *E3040) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		constraints := schema.GetResourceConstraints(res.Type)
		if constraints == nil || len(constraints.ReadOnlyProperties) == 0 {
			continue
		}

		// Build set of read-only properties
		readOnly := make(map[string]bool)
		for _, prop := range constraints.ReadOnlyProperties {
			readOnly[prop] = true
		}

		// Check each property
		for propName := range res.Properties {
			if readOnly[propName] {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Property '%s' in resource '%s' (%s) is read-only and cannot be specified",
						propName, resName, res.Type,
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
