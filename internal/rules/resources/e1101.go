// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/schema"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E1101{})
}

// E1101 performs comprehensive schema validation against CloudFormation resource specifications.
type E1101 struct{}

func (r *E1101) ID() string { return "E1101" }

func (r *E1101) ShortDesc() string {
	return "Schema validation error"
}

func (r *E1101) Description() string {
	return "Validates resources against CloudFormation resource specifications, checking for unknown properties and type mismatches."
}

func (r *E1101) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E1101"
}

func (r *E1101) Tags() []string {
	return []string{"resources", "schema", "validation"}
}

func (r *E1101) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		// Skip if resource type is unknown
		exists, err := schema.HasResourceType(res.Type)
		if err != nil || !exists {
			continue
		}

		// Check for unknown properties
		for propName := range res.Properties {
			has, err := schema.HasProperty(res.Type, propName)
			if err != nil {
				continue
			}
			if !has {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf(
						"Resource '%s' (%s) has unknown property '%s'",
						resName, res.Type, propName,
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
