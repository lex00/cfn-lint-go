// Package outputs contains output validation rules (E6xxx).
package outputs

import (
	"fmt"

	"github.com/lex00/cfn-lint-go/pkg/rules"
	"github.com/lex00/cfn-lint-go/pkg/template"
)

func init() {
	rules.Register(&E6002{})
}

// E6002 checks that outputs have the required Value property.
type E6002 struct{}

func (r *E6002) ID() string { return "E6002" }

func (r *E6002) ShortDesc() string {
	return "Outputs have required properties"
}

func (r *E6002) Description() string {
	return "Checks that all Outputs have the required Value property."
}

func (r *E6002) Source() string {
	return "https://github.com/aws-cloudformation/cfn-lint/blob/main/docs/rules.md#E6002"
}

func (r *E6002) Tags() []string {
	return []string{"outputs", "required"}
}

func (r *E6002) Match(tmpl *template.Template) []rules.Match {
	var matches []rules.Match

	for name, out := range tmpl.Outputs {
		if out.Value == nil {
			matches = append(matches, rules.Match{
				Message: fmt.Sprintf("Output '%s' is missing required property 'Value'", name),
				Line:    out.Node.Line,
				Column:  out.Node.Column,
				Path:    []string{"Outputs", name},
			})
		}
	}

	return matches
}
