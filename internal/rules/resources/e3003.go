// Package resources contains resource validation rules (E3xxx).
package resources

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/schema"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E3003{})
}

// E3003 checks that resources have required properties.
// Uses CloudFormation resource schemas from cloudformation-schema-go.
type E3003 struct{}

func (r *E3003) ID() string { return "E3003" }

func (r *E3003) ShortDesc() string {
	return "Required properties are present"
}

func (r *E3003) Description() string {
	return "Checks that resources have their required properties based on CloudFormation resource schemas."
}

func (r *E3003) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E3003"
}

func (r *E3003) Tags() []string {
	return []string{"resources", "required", "properties"}
}

func (r *E3003) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for resName, res := range tmpl.Resources {
		required, err := schema.GetRequiredProperties(res.Type)
		if err != nil {
			// Schema loading error - skip validation for this resource
			continue
		}
		if required == nil {
			// Unknown resource type - skip validation
			continue
		}

		for _, prop := range required {
			if _, exists := res.Properties[prop]; !exists {
				matches = append(matches, rules.Match{
					Message: fmt.Sprintf("Resource '%s' (%s) is missing required property '%s'", resName, res.Type, prop),
					Line:    res.Node.Line,
					Column:  res.Node.Column,
					Path:    []string{"Resources", resName, "Properties"},
				})
			}
		}
	}

	return matches
}
